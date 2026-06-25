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

## gRPC rip — transport status
The filer and master RPC services run **end to end over the native ZAP transport**
(`zap-proto/go`); gRPC is gone from those paths. Error semantics cross the wire as
the error STRING, not a gRPC status: each server tags a failure with its PascalCase
code name (e.g. `fmt.Errorf("Unavailable: ...")`, `fmt.Errorf("NotFound: ...")`) and
the ZAP dispatch ships `herr.Error()` verbatim, so the name survives to the client.
**Classify errors by `strings.Contains(err.Error(), "<CodeName>")`** — the one way,
already used in `s3/util/retry.go`. Do NOT reintroduce `google.golang.org/grpc/{codes,status}`
for filer/master error mapping. The filer's not-found is normalized back to the
`filer_pb.ErrNotFound` sentinel by `filer_pb.LookupEntry`, so `errors.Is(err, ErrNotFound)`
also works there. A round-trip proof lives in `s3/masterzap/error_propagation_test.go`.

gRPC is still **load-bearing** in (these legitimately import grpc — leave them):
- **Volume server RPC** — still gRPC-served (`command/volume.go`, `server/volume_grpc_*.go`,
  `pb/volume_server_pb/volume_server_grpc.pb.go`). Its `status.Errorf(codes.X)` producers
  must stay so gRPC clients get real codes; `peer.FromContext` works over gRPC.
- **MQ broker RPC** — still gRPC-served (`command/mq_broker.go`, `mq/broker/broker_grpc_*.go`,
  `pb/mq_pb/mq_broker_grpc.pb.go`).
- **Raft transport** — hashicorp/seaweedfs raft rides the master's gRPC listener
  (`server/raft_hashicorp.go`, the `pb.NewGrpcServer` + `reflection.Register` in
  `command/master.go` and `master_follower.go`). The master's own service is ZAP; the
  gRPC listener carries raft only.
- **Strangler client seams** `mqzap/client.go` / `volumezap/client.go` — they implement
  the `*_pb.*Client` interfaces, whose generated signatures use `grpc.{CallOption,
  *StreamingClient}` and `metadata.MD`. Bound to grpc until those `*_pb` interfaces are
  regenerated grpc-free. Neither is wired into production yet.
- **Core dial/server + request-id** — `pb/grpc_client_server.go` (volume/broker dials,
  `NewGrpcServer`), `operation/grpc_client.go`, `server/common.go` + `admin/dash/admin_server.go`
  (grpc metadata / broker dial).
- **TLS cert hot-reload** — `security/tls.go`, `security/certreload/certreload.go`,
  `security/tls_reload.go`, `util/http/client/http_client.go` reuse grpc's
  `credentials/tls/certprovider` + `pemfile` as a generic file-watching cert provider;
  `tls.go` also returns `grpc.ServerOption` for the remaining gRPC servers.

Full `go.mod` grpc removal is unreachable until volume RPC + MQ broker RPC are ported to
ZAP (regen their `*_pb` interfaces grpc-free, switch servers to ZAP dispatch) and cert
hot-reload drops grpc's certprovider.

### ZAP listener ports — derive via `pb.ZapPort` (overflow-safe)
The master and IAM services serve ZAP on a port offset from their gRPC port by
`+10000` (i.e. httpPort+20000). That offset is `pb.ZapPort(grpcPort)`, the ONE
place the convention lives — both the client (`ServerAddress.ToMasterZapAddress`/
`ToIamZapAddress`) and the server listeners (`command/master.go`, `command/filer.go`
IAM listener) call it, so they always agree. It is a **rotation of [1,65535]**, not a
bare `grpcPort+10000`: an ephemeral high port (e.g. test ports from
`testutil.AllocateMiniPorts`, up to ~55000 → grpcPort ~65000) would otherwise push
`+10000` past 65535, yield an invalid port, and the master would Fatalf on its ZAP
listener — hanging cluster bring-up. `ZapPort` stays a plain `+10000` for the common
case (grpcPort ≤ 55535) and folds back into range above that, bijectively. Guard:
`pb/server_address_test.go` (`TestZapPort`).

### Test doubles for the deleted server scaffolding
The grpc rip deleted `filer_pb.UnimplementedHanzoFilerServer` / `master_pb.UnimplementedHanzoServer`
(no `mustEmbed`). Test doubles that need the full server method set embed
`filerstub.FilerServer` / `filerstub.MasterServer` from `s3/pb/filerstub` (one shared
base, every method stubbed to an unimplemented error) and override what they exercise —
the same shape `Unimplemented*` had. Doubles that serve over the wire instead embed the
`filerwire.Backend` / `masterwire.Backend` interface and implement the wire-level methods
(decode bytes → encode bytes), e.g. `test/plugin_workers/fake_master.go`. Client fakes
implement `filer_pb.HanzoFilerClient`/`master_pb.HanzoClient` directly: no
`opts ...grpc.CallOption`; streaming methods return `rpc.ServerStream[...]`/`rpc.BidiStream[...]`
(seam in `s3/pb/rpc`), and the fake stream type implements just `Recv`/`Send`/`CloseSend`.
DialOption in tests is `pb.DialOption{}` (the zero value), never `grpc.WithTransportCredentials`.

## Rebrand notes (one way, don't re-litigate)
- `s3` → `s3`: binary + program name in `s3/s3.go` help. Internal package dir `s3/` retained (invisible to users; renaming it would churn 1900+ import paths for no user benefit).
- Display brand `Hanzo` → `Hanzo S3` in user-facing strings only (HTTP `Server` header, version banner, admin title).
- **Wire-protocol identifiers/headers stay `Hanzo*`** (`HanzoSessionTokenHeader`, SSE-S3/KMS key names, RDMA client, etc.) — renaming them breaks S3 compatibility. This is intentional.
- Version: upstream Hanzo 4.34; image injects `COMMIT` via ldflags (empty COMMIT panics `s3 version`).
