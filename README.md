# Hanzo S3

Hanzo S3 is Hanzo AI's S3-compatible object store: a fast, simple, scalable
distributed storage system with a built-in S3 gateway, Filer, and FUSE mount.

It is an Apache-2.0 fork of [SeaweedFS](https://github.com/seaweedfs/seaweedfs).
The only differences from upstream are branding: the command-line binary is
shipped as `s3` (instead of `weed`), and user-facing help/branding reads
"Hanzo S3". All internal Go packages keep their upstream import paths, so the
fork tracks SeaweedFS with minimal divergence.

## Why this exists

The previous `hanzoai/s3` was an AGPL-3.0 MinIO fork. AGPL's network-copyleft
makes it impractical to embed and resell as part of a managed offering.
SeaweedFS is Apache-2.0, so Hanzo S3 can be embedded, modified, and offered as
a service without copyleft obligations.

## Build

Requires Go 1.25+.

```bash
make            # builds ./s3 at the repo root
# or
go build -o s3 ./weed
```

Full build with optional backends (elastic, sqlite, tikv, ydb, tarantool, rclone, gocdk):

```bash
make full_install
```

## Usage

The binary is `s3`. Common commands:

```bash
s3 server -s3 -filer -dir=/data    # all-in-one: master + volume + filer + S3 gateway
s3 master                          # start a master server
s3 volume -dir=/data -max=5        # start a volume server
s3 filer -master=localhost:9333    # start the Filer
s3 s3 -filer=localhost:8888        # start the S3 API gateway
s3 mount -filer=localhost:8888 -dir=/mnt/s3   # FUSE mount
s3 mini                            # single-process setup for S3 beginners / dev
s3 -h                              # full command list
s3 help <command>                  # help for a specific command
s3 version
```

The S3 gateway speaks the AWS S3 API; point any S3 SDK or `aws s3` client at it.

## Documentation

Hanzo S3 is API- and behavior-compatible with upstream SeaweedFS, so the
upstream documentation applies directly:

- SeaweedFS Wiki: https://github.com/seaweedfs/seaweedfs/wiki
- S3 API reference: https://github.com/seaweedfs/seaweedfs/wiki/Amazon-S3-API

Wherever the upstream docs say `weed`, run `s3` instead.

## License

Apache License 2.0 — see [LICENSE](./LICENSE).

This is a fork of SeaweedFS (Copyright Chris Lu and the SeaweedFS
contributors). See [NOTICE](./NOTICE) for attribution.
