FROM golang:1.19-bullseye as build

COPY . /src
RUN set -ex \
 && cd /src \
 && CGO_ENABLED=0 go build -o /bin/prometheus-exporter \
 && strip /bin/prometheus-exporter

FROM alpine:3.17

COPY --from=build /bin/prometheus-exporter /bin/prometheus-exporter

USER nobody
EXPOSE     9055
ENTRYPOINT [ "/bin/prometheus-exporter" ]
