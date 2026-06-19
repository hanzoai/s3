# Hanzo S3

## Overview
Go module: `github.com/hanzoai/s3`

**Upstream**: [SeaweedFS](https://github.com/seaweedfs/seaweedfs) (Apache-2.0). Hanzo S3 is
an Apache-2.0 fork of SeaweedFS. It replaces the previous AGPL-3.0 MinIO fork so the
object store can be embedded and offered as a managed service without copyleft.

## What changed vs upstream
- **Binary/CLI rebranded** `weed` -> `s3`. Users only ever see `s3` (`s3 server`, `s3 filer`,
  `s3 master`, `s3 volume`, `s3 s3` for the S3 gateway, `s3 mount`, `s3 mini`).
- User-facing help/usage strings and the product name in CLI output read "Hanzo S3".
- `go.mod` module path is `github.com/hanzoai/s3` with a self-`replace` directive
  (`replace github.com/seaweedfs/seaweedfs => ./`) so the ~1800 internal Go files keep their
  upstream import paths (`github.com/seaweedfs/seaweedfs/...`) with **zero source churn**.
  Internal package identifiers and import paths are intentionally NOT renamed.

## Layout (key, from upstream)
- `weed/` — the `main` package (built as the `s3` binary) + all subsystems
  (`s3api/`, `filer/`, `server/`, `command/`, `storage/`, `topology/`, `mount/`, ...).
- `weed/command/` — cobra-style subcommands; this is where user-facing help text lives.
- `docker/`, `k8s/`, `helm` charts under `k8s/`, `terraform/` — deployment.

## Build & Run
```bash
make                  # builds ./s3 at the repo root (go build -o s3 ./weed)
go build -o s3 ./weed # equivalent
make full_install     # build with optional backends (sqlite, tikv, ydb, elastic, rclone, ...)
./s3 -h               # full command list
```

## License
Apache-2.0 (`LICENSE`). Attribution to SeaweedFS in `NOTICE`.
