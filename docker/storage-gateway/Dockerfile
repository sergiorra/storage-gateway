FROM golang:1.21.0-alpine AS builder
WORKDIR /src/storage-gateway
COPY . .
RUN go build ./cmd/storage-gateway

FROM docker:24.0.0-alpine3.18 AS runner

COPY --from=builder /src/storage-gateway/storage-gateway /opt/storage-gateway/bin/storage-gateway
COPY --from=builder /src/storage-gateway/config/config.docker.json /opt/storage-gateway/etc/config/

RUN apk add bash curl

ENTRYPOINT /opt/storage-gateway/bin/storage-gateway -conf /opt/storage-gateway/etc/config/config.docker.json