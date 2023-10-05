# Use the golang:1.21.0-alpine image as the base image for the build stage
FROM golang:1.21.0-alpine AS builder

WORKDIR /src/storage-gateway
COPY . .

# Build the go application located at cmd/storage-gateway/main.go
RUN go build ./cmd/storage-gateway

# Switch to the docker:24.0.0-alpine3.18 image for the runtime stage
FROM docker:24.0.0-alpine3.18 AS runner

COPY --from=builder /src/storage-gateway/storage-gateway /opt/storage-gateway/bin/storage-gateway
COPY --from=builder /src/storage-gateway/config/config.docker.json /opt/storage-gateway/etc/config/

RUN apk add bash curl

# Set the entrypoint for the container to be the "storage-gateway" binary with the config file path as an argument
ENTRYPOINT /opt/storage-gateway/bin/storage-gateway -conf /opt/storage-gateway/etc/config/config.docker.json