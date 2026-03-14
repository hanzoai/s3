# Hanzo S3

[![GitHub Stars](https://img.shields.io/github/stars/hanzoai/s3?style=flat-square)](https://github.com/hanzoai/s3)
[![License](https://img.shields.io/badge/license-AGPL%20v3-blue?style=flat-square)](https://github.com/hanzoai/s3/blob/main/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hanzoai/s3?style=flat-square)](https://github.com/hanzoai/s3)

High-performance, S3-compatible object storage for AI workloads, optimized for the Hanzo ecosystem.

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

```sh
docker run -p 9000:9000 -p 9001:9001 \
  -e S3_ROOT_USER=admin \
  -e S3_ROOT_PASSWORD=changeme123 \
  ghcr.io/hanzoai/s3:latest server /data --console-address :9001
```

Console: `http://127.0.0.1:9001` -- change the default credentials immediately in production.

### Install from Source

Requires Go 1.24 or later.

```sh
git clone https://github.com/hanzoai/s3.git
cd s3
go build -o hanzo-s3 .
./hanzo-s3 server /data --console-address :9001
```

### Verify Connectivity

Use any S3-compatible client. With the Hanzo S3 CLI (`s3`):

```sh
s3 alias set hanzo http://localhost:9000 admin changeme123
s3 admin info hanzo
s3 mb hanzo/my-bucket
s3 cp ~/data/model.safetensors hanzo/my-bucket/
s3 ls hanzo/my-bucket/
```

## SDKs

Hanzo S3 is fully S3-compatible. Use any S3 SDK:

| Language | Package |
|----------|---------|
| Go       | [`@hanzo/s3-go`](https://github.com/hanzos3/go-sdk) |
| JavaScript / TypeScript | [`@hanzo/s3`](https://github.com/hanzos3/js-sdk) |
| Python   | [`hanzo-s3`](https://github.com/hanzos3/py-sdk) |

Standard AWS SDKs (`aws-sdk-go`, `boto3`, `@aws-sdk/client-s3`) also work without modification.

## Documentation

Full documentation at [docs.hanzo.ai/storage](https://docs.hanzo.ai/docs/services/s3).

## Demo Server

A public demo server is available at `s3-demo.hanzo.ai` for testing. Data is wiped hourly.

## Attribution

Based on [MinIO](https://github.com/minio/minio). See the upstream [LICENSE](LICENSE) for attribution.

## License

Copyright (c) Hanzo AI Inc. Licensed under the [GNU Affero General Public License v3.0](LICENSE).
