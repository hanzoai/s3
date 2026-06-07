FROM golang:1.26.4-bookworm AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG TARGETARCH
ARG VERSION=""
ARG COMMIT_ID=""

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -tags kqueue -trimpath \
    -ldflags "-s -w \
      -X github.com/hanzoai/s3/cmd.Version=${VERSION} \
      -X github.com/hanzoai/s3/cmd.CopyrightYear=2026 \
      -X github.com/hanzoai/s3/cmd.CommitID=${COMMIT_ID}" \
    -o /s3 .

FROM alpine:latest

RUN apk add --no-cache curl ca-certificates

COPY --from=build /s3 /usr/bin/minio

VOLUME ["/data"]
EXPOSE 9000 9001

ENTRYPOINT ["/usr/bin/minio"]
CMD ["server", "/data", "--console-address", ":9001"]
