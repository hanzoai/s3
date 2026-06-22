# Hanzo S3

## Overview
Go module: `github.com/hanzoai/s3`

**Upstream**: [Hanzo](https://github.com/hanzo/hanzo) (Apache-2.0). Branded as **Hanzo S3** — S3-compatible object storage. Replaced the former MinIO fork (MinIO killed; migrated to Hanzo).

The binary is **`s3`** (renamed from upstream `s3`). All imports are on `github.com/hanzoai/s3`; upstream `hanzo/{raft,goexif,go-fuse,cockroachdb-parser}` remain external deps.

## Build & Run
```bash
CGO_ENABLED=0 go build -trimpath -o s3 ./s3
./s3 server -s3 -s3.port=9000 -dir=/data    # all-in-one: master + volume + filer + s3
```

## Image
`ghcr.io/hanzoai/s3` — built by `.github/workflows/docker-build.yml` via the shared `hanzoai/.github` reusable. `ENTRYPOINT ["s3"]`; default `CMD` runs the all-in-one server with the S3 gateway on `:9000`.

## Rebrand notes (one way, don't re-litigate)
- `s3` → `s3`: binary + program name in `s3/s3.go` help. Internal package dir `s3/` retained (invisible to users; renaming it would churn 1900+ import paths for no user benefit).
- Display brand `Hanzo` → `Hanzo S3` in user-facing strings only (HTTP `Server` header, version banner, admin title).
- **Wire-protocol identifiers/headers stay `Hanzo*`** (`HanzoSessionTokenHeader`, SSE-S3/KMS key names, RDMA client, etc.) — renaming them breaks S3 compatibility. This is intentional.
- Version: upstream Hanzo 4.34; image injects `COMMIT` via ldflags (empty COMMIT panics `s3 version`).
