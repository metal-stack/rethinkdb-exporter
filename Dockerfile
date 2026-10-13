FROM gcr.io/distroless/static-debian13:nonroot
COPY bin/rethinkdb-exporter /prometheus-exporter
USER nobody
EXPOSE 9055
ENTRYPOINT ["/prometheus-exporter"]
