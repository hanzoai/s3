# Hanzo S3 — S3-compatible object storage (SeaweedFS engine).
# Builds the `s3` binary from github.com/hanzoai/s3 (rebranded SeaweedFS).
FROM golang:1.26-alpine AS build
RUN apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
ARG COMMIT=unknown
RUN CGO_ENABLED=0 go build -trimpath \
      -ldflags "-s -w -X github.com/hanzoai/s3/weed/util/version.COMMIT=${COMMIT}" \
      -o /out/s3 ./weed

FROM alpine:3.21
RUN apk add --no-cache ca-certificates fuse3 curl && mkdir -p /data
COPY --from=build /out/s3 /usr/bin/s3
LABEL org.opencontainers.image.title="Hanzo S3" \
      org.opencontainers.image.description="S3-compatible object storage (SeaweedFS engine)" \
      org.opencontainers.image.source="https://github.com/hanzoai/s3" \
      org.opencontainers.image.vendor="Hanzo AI" \
      org.opencontainers.image.licenses="Apache-2.0"
VOLUME /data
# S3 API :9000  | master :9333 | volume :8080 | filer :8888
EXPOSE 9000 9333 8080 8888
ENTRYPOINT ["s3"]
CMD ["server", "-s3", "-s3.port=9000", "-dir=/data"]
