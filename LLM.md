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

## Filer metadata store: zapdb (leveldb ripped)
The filer's on-disk metadata `FilerStore` is **zapdb** (`github.com/luxfi/zapdb`, a
transactional LSM-tree KV — Badger-derived; its Go package is `package badger`, so
import it aliased: `zapdb "github.com/luxfi/zapdb"`).
- Store: `s3/filer/zapdb/` (`zapdb_store.go` + `zapdb_store_kv.go`). `GetName()` =
  `"zapdb"`; registered in `init()` via `filer.Stores`. Key encoding is byte-for-byte
  identical to the old leveldb store (`dir\x00name`), so the on-disk key layout is the
  same shape.
- The FUSE mount holds a concrete `*zapdb.ZapDBStore` for batch writes:
  `s3/mount/meta_cache/meta_cache.go` (`zapdbStore` field + `openMetaStore` +
  `BatchInsertEntries`).
- Default store when no `filer.toml` is present: `s3/server/filer_server.go` sets
  `zapdb.enabled`/`zapdb.dir` (was `leveldb2`).
- `s3/filer/leveldb/` is **deleted**. goleveldb (`syndtr/goleveldb`) remains a dep ONLY
  for the volume server's needle-map index (`s3/storage/needle_map_leveldb.go`,
  `chunk_cache`) — a different subsystem; do not confuse it with the filer store.
- SQLite store (`-tags sqlite`) owns its own SQL generator now
  (`s3/filer/sqlite/sqlite_sql_gen.go`, `SqlGenSqlite`) — the deleted `filer/mysql`
  backend used to supply `SqlGenMysql`; that dangling import is gone.
- Dead DB drivers from the 20-backend rip were dropped via `go mod tidy`
  (go-sql-driver/mysql, mongo-driver*, gocql, arangodb, go-redis, tarantool, tikv,
  ydb, gohbase, olivere/elastic*). *elastic stays (meta-tail CLI); mongo-driver stays
  as a pure transitive dep.

## Rebrand notes (one way, don't re-litigate)
- `s3` → `s3`: binary + program name in `s3/s3.go` help. Internal package dir `s3/` retained (invisible to users; renaming it would churn 1900+ import paths for no user benefit).
- Display brand `Hanzo` → `Hanzo S3` in user-facing strings only (HTTP `Server` header, version banner, admin title).
- **Wire-protocol identifiers/headers stay `Hanzo*`** (`HanzoSessionTokenHeader`, SSE-S3/KMS key names, RDMA client, etc.) — renaming them breaks S3 compatibility. This is intentional.
- Version: upstream Hanzo 4.34; image injects `COMMIT` via ldflags (empty COMMIT panics `s3 version`).
