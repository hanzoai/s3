# Hanzo Storage

S3-compatible object storage for AI assets.

## Overview

Hanzo Storage is a high-performance, S3-compatible object storage system designed for AI workloads. Store and retrieve model weights, training datasets, embeddings, and other large AI assets with enterprise-grade reliability.

## Features

- **S3 Compatible**: Full AWS S3 API compatibility
- **High Performance**: Optimized for large file transfers
- **Erasure Coding**: Data protection without replication overhead
- **Encryption**: Server-side and client-side encryption
- **Versioning**: Object versioning and lifecycle policies
- **Events**: Webhook notifications for object changes

## Quick Start

```bash
docker run -p 9000:9000 -p 9001:9001 hanzo/storage
```

## Documentation

See the [documentation](https://hanzo.ai/docs/storage) for detailed guides and API reference.

## License

MIT License - see [LICENSE](LICENSE) for details.
