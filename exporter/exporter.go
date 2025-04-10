package exporter

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// RethinkdbExporter is a prometheus exporter of the rethinkdb statistics
type RethinkdbExporter struct {
	rconn r.QueryExecutor

	collectTableStats bool

	listenAddress string
	mux           *http.ServeMux

	log     *slog.Logger
	metrics struct {
		clusterClientConnections *prometheus.Desc
		clusterDocsPerSecond     *prometheus.Desc

		serverClientConnections *prometheus.Desc
		serverQueriesPerSecond  *prometheus.Desc
		serverDocsPerSecond     *prometheus.Desc

		tableDocsPerSecond *prometheus.Desc
		tableRowsCount     *prometheus.Desc

		tableReplicaDocsPerSecond *prometheus.Desc
		tableReplicaCacheBytes    *prometheus.Desc
		tableReplicaIO            *prometheus.Desc
		tableReplicaDataBytes     *prometheus.Desc

		scrapeLatency *prometheus.Desc
		scrapeErrors  *prometheus.Desc
	}
}

type promHTTPLogger struct {
	log *slog.Logger
}

func (l promHTTPLogger) Println(v ...interface{}) {
	l.log.Error("promhttp", "msg", fmt.Sprint(v...))
}

// New creates a new instance of prometheus rethinkdb exporter
func New(
	log *slog.Logger,
	listenAddress string,
	telemetryPath string,
	rconn r.QueryExecutor,
	collectTableStats bool,
) (*RethinkdbExporter, error) {
	exporter := &RethinkdbExporter{
		listenAddress:     listenAddress,
		collectTableStats: collectTableStats,
		rconn:             rconn,
		log:               log,
	}

	exporter.initMetrics()

	prometheus.MustRegister(exporter)

	exporter.mux = http.NewServeMux()
	exporter.mux.Handle(telemetryPath,
		promhttp.InstrumentMetricHandler(
			prometheus.DefaultRegisterer,
			promhttp.HandlerFor(
				prometheus.DefaultGatherer,
				promhttp.HandlerOpts{
					ErrorLog: &promHTTPLogger{log: log},
				},
			),
		),
	)
	exporter.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
             <head><title>RethinkDB Exporter</title></head>
             <body>
             <h1>RethinkDB Exporter</h1>
             <p><a href='` + telemetryPath + `'>Metrics</a></p>
             <h2>Build</h2>
             <pre>` + version.Info() + ` ` + version.BuildContext() + `</pre>
             </body>
             </html>`))
	})
	exporter.mux.HandleFunc("/-/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "OK")
	})
	exporter.mux.HandleFunc("/-/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "OK")
	})

	return exporter, nil
}

// ListenAndServe runs prometheus http-server for exporting stats
func (e *RethinkdbExporter) ListenAndServe() error {
	serv := http.Server{Addr: e.listenAddress, Handler: e.mux, ReadHeaderTimeout: 10 * time.Second}
	return serv.ListenAndServe()
}
