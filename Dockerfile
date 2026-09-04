# Hanzo S3 — S3-compatible object storage (Hanzo engine).
# Builds the `s3` binary from github.com/hanzoai/s3 (rebranded Hanzo).
FROM golang:1.26.5-alpine AS build
ENV GOTOOLCHAIN=auto
RUN apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
ARG COMMIT=unknown
RUN CGO_ENABLED=0 go build -trimpath \
      -ldflags "-s -w -X github.com/hanzoai/s3/s3/util/version.COMMIT=${COMMIT}" \
      -o /out/s3 ./s3

# One directory in an empty image: the static binary and the data directory it
# serves. The `mount` subcommand shells out to fusermount and has no place in a
# server image; run it from a host that has one.
FROM alpine:3.22 AS root
RUN apk add --no-cache ca-certificates tzdata && mkdir -p /data && chown 65532:65532 /data

FROM scratch
LABEL org.opencontainers.image.title="Hanzo S3" \
      org.opencontainers.image.description="S3-compatible object storage (Hanzo engine)" \
      org.opencontainers.image.source="https://github.com/hanzoai/s3" \
      org.opencontainers.image.vendor="Hanzo AI" \
      org.opencontainers.image.licenses="Apache-2.0"
COPY --from=root /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=root /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=root --chown=65532:65532 /data /data
COPY --from=build /out/s3 /usr/bin/s3
USER 65532:65532
VOLUME /data
EXPOSE 9000 9333 8080 8888
ENTRYPOINT ["/usr/bin/s3"]
CMD ["server", "-s3", "-s3.port=9000", "-dir=/data"]
