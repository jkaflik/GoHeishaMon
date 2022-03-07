FROM gcr.io/distroless/static
ENTRYPOINT ["/heatpump2mqtt"]
COPY heatpump2mqtt /
COPY topics.yaml /