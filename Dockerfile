FROM golang:1.24-bookworm as build

COPY . /src
RUN set -ex \
 && cd /src \
 && CGO_ENABLED=0 go build -o /bin/prometheus-exporter \
 && strip /bin/prometheus-exporter

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /bin/prometheus-exporter /bin/prometheus-exporter

USER nobody
EXPOSE 9055
ENTRYPOINT [ "/bin/prometheus-exporter" ]
