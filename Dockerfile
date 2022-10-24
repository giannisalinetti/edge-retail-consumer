FROM registry.redhat.io/rhel8/go-toolset:1.17.12-3.1661802325 as builder

WORKDIR /app
USER root
COPY . .
RUN go build .

FROM registry.access.redhat.com/ubi8/ubi-minimal

COPY --from=builder /app/edge-retail-consumer /

CMD ["/edge-retail-consumer", "-help"]