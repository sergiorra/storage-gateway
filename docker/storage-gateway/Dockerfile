FROM golang:1.19
WORKDIR /mnt/homework
COPY . .
RUN go build ./cmd/object-gateway

# Docker is used as a base image so you can easily start playing around in the container using the Docker command line client.
FROM docker
COPY --from=0 /mnt/homework/object-gateway /usr/local/bin/homework-object-storage
RUN apk add bash curl

ENTRYPOINT ["homework-object-storage"]