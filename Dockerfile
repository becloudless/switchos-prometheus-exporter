# syntax=docker/dockerfile:1
#
# Build note: this module depends on github.com/becloudless/switchos-client
# via a local `replace ... => ../switchos-client` directive (it's not a
# published module). The builder stage therefore expects that repository
# to be provided as a second build context named "switchos-client",
# checked out as a sibling of this repo (matching local development
# layout). Example:
#
#   docker buildx build \
#     --build-context switchos-client=../switchos-client \
#     -t switchos-prometheus-exporter .
#
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src
COPY . .
COPY --from=switchos-client . /switchos-client

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/switchos-prometheus-exporter .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 10001 exporter
COPY --from=builder /out/switchos-prometheus-exporter /usr/bin/switchos-prometheus-exporter

USER exporter
EXPOSE 9435
ENTRYPOINT ["/usr/bin/switchos-prometheus-exporter"]
CMD ["-config", "/etc/switchos-exporter.yaml"]
