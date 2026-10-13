FROM gcr.io/distroless/static-debian13:nonroot
COPY bin/prometheus-exporter /bin/prometheus-exporter

USER nobody
EXPOSE 9055
ENTRYPOINT ["/bin/prometheus-exporter"]
