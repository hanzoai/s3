# Hanzo Storage

[![GitHub Stars](https://img.shields.io/github/stars/hanzoai/storage?style=flat-square)](https://github.com/hanzoai/storage)
[![License](https://img.shields.io/badge/license-AGPL%20v3-blue?style=flat-square)](https://github.com/hanzoai/storage/blob/main/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hanzoai/storage?style=flat-square)](https://github.com/hanzoai/storage)

High-performance, S3-compatible object storage for AI workloads, built on MinIO and optimized for the Hanzo ecosystem.

## Features

- **S3 API Compatible** -- Drop-in replacement for Amazon S3; works with any S3 client or SDK
- **Built for AI & Analytics** -- Optimized for large-scale model artifacts, training datasets, and data pipelines
- **High Performance** -- Designed to saturate modern NVMe and network hardware
- **Erasure Coding** -- Data protection with configurable redundancy across drives and nodes
- **Bucket Policies** -- Fine-grained access control with S3-compatible policy documents
- **Object Lifecycle Management** -- Automated expiration, transition, and tiering rules
- **Encryption** -- Server-side encryption (SSE-S3, SSE-KMS) for data at rest and in transit
- **Multi-Tenancy** -- Isolated namespaces and access boundaries for teams and services

## Quick Start

### Docker

Build and run a standalone Hanzo Storage server:

```sh
docker build -t hanzo-storage .
docker run -p 9000:9000 -p 9001:9001 \
  hanzo-storage server /data --console-address :9001
```

The web console is available at `http://127.0.0.1:9001`. Default credentials are `minioadmin:minioadmin` -- change these immediately in production.

### Install from Source

Requires Go 1.24 or later.

```sh
git clone https://github.com/hanzoai/storage.git
cd storage
go build -o hanzo-storage .
./hanzo-storage server /data --console-address :9001
```

Cross-compile for a specific target:

```sh
GOOS=linux GOARCH=arm64 go build -o hanzo-storage .
```

### Verify Connectivity

Use any S3-compatible client. With the MinIO Client (`mc`):

```sh
mc alias set hanzo http://localhost:9000 minioadmin minioadmin
mc admin info hanzo
mc mb hanzo/my-bucket
mc cp ~/data/model.safetensors hanzo/my-bucket/
mc ls hanzo/my-bucket/
```

## SDKs

Hanzo Storage is fully S3-compatible. Use any S3 SDK, or the purpose-built MinIO SDKs:

| Language | Package |
|----------|---------|
| Go       | [`github.com/minio/minio-go/v7`](https://github.com/minio/minio-go) |
| JavaScript / TypeScript | [`minio`](https://github.com/minio/minio-js) |
| Python   | [`minio`](https://github.com/minio/minio-py) |

Standard AWS SDKs (`aws-sdk-go`, `boto3`, `@aws-sdk/client-s3`) also work without modification.

## Documentation

Full documentation is available at [docs.hanzo.ai](https://docs.hanzo.ai).

## Attribution

Based on [MinIO](https://github.com/minio/minio). See the upstream [LICENSE](LICENSE) for attribution.

## License

Copyright (c) Hanzo AI Inc. Licensed under the [GNU Affero General Public License v3.0](LICENSE).
