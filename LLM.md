# Hanzo S3

## Overview
Go module: `github.com/hanzoai/s3`

**Upstream**: [SeaweedFS](https://github.com/seaweedfs/seaweedfs) (Apache-2.0). Hanzo S3 is the Apache-2.0 SeaweedFS fork, branded **Hanzo S3** — S3-compatible object storage. It replaced the former AGPL-3.0 MinIO fork (dropped). Attribution to SeaweedFS is retained in `NOTICE`.

The binary is **`s3`** (renamed from upstream `weed`). All internal code imports `github.com/hanzoai/s3`. The former external SeaweedFS-org utility libs are now Hanzo forks too: `github.com/hanzoai/goexif` (v1.0.4) and `github.com/hanzoai/go-fuse/v2` (v2.9.4) — so **no `seaweedfs`-org import paths remain**. Each fork retains its upstream license verbatim: `goexif` under **BSD-2-Clause** (Robert Carlsen), `go-fuse` under **BSD-3-Clause** (the Go-FUSE Authors / hanwen + AUTHORS) — neither has a `NOTICE`. (The s3 codebase itself is Apache-2.0 with a `NOTICE` crediting upstream SeaweedFS — that's separate from these two BSD libs.)

## Build & Run
```bash
CGO_ENABLED=0 go build -trimpath -o s3 ./s3
./s3 server -s3 -s3.port=9000 -dir=/data    # all-in-one: master + volume + filer + s3
```

## Image
`ghcr.io/hanzoai/s3` — built by `.github/workflows/docker-build.yml` via the shared `hanzoai/.github` reusable. `ENTRYPOINT ["s3"]`; default `CMD` runs the all-in-one server with the S3 gateway on `:9000`.

## S-Chain — engines on top

Hanzo S3 is the platform's distributed object-storage substrate, **"S-Chain"**. Query
and compute engines build on it, both storing their data here as objects:

- **Datastore** — the OLAP (columnar) engine, a ClickHouse fork
  ([hanzoai/datastore](https://github.com/hanzoai/datastore)). Its `MergeTree` parts
  live here, so stateless compute replicas share one zero-copy copy of the data
  (proven PoC: 5M-row table on S3, two compute nodes, one physical copy).
- **Decentralized SQL** — the planned OLTP engine (SQLite/Base), designed not yet
  built: one encrypted database file per user/project, also stored here.

Direction (not yet shipped): Quasar (post-quantum, leaderless) consensus for the small
commitments, with bulk data staying in S-Chain object storage. Hanzo S3's own transport
is already ZAP-native (see the gRPC-rip section below). See the Datastore design paper,
`hanzo-datastore/` in [hanzoai/papers](https://github.com/hanzoai/papers).

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
The filer, master, and **MQ broker** RPC services run **end to end over the native
ZAP transport** (`zap-proto/go`); gRPC is gone from those paths. Error semantics cross the wire as
the error STRING, not a gRPC status: each server tags a failure with its PascalCase
code name (e.g. `fmt.Errorf("Unavailable: ...")`, `fmt.Errorf("NotFound: ...")`) and
the ZAP dispatch ships `herr.Error()` verbatim, so the name survives to the client.
**Classify errors by `strings.Contains(err.Error(), "<CodeName>")`** — the one way,
already used in `s3/util/retry.go`. Do NOT reintroduce `google.golang.org/grpc/{codes,status}`
for filer/master error mapping. The filer's not-found is normalized back to the
`filer_pb.ErrNotFound` sentinel by `filer_pb.LookupEntry`, so `errors.Is(err, ErrNotFound)`
also works there. A round-trip proof lives in `s3/masterzap/error_propagation_test.go`.

The MQ broker rip mirrors filer/master exactly: `pb/mq_pb/mq_broker_grpc.pb.go` is
grpc-free (client iface returns `pb/rpc` stream seams, no `grpc.CallOption`; server
iface is a plain method set; all `Register*`/`ServiceDesc`/`_Handler`/`Unimplemented`-
via-grpc scaffolding deleted, replaced by a de-grpc'd `UnimplementedHanzoMessagingServer`
that returns a plain error). `mqzap/client.go` returns `rpc.{Bidi,Client,Server}Stream`;
`mqzap/stream_server.go` drops the `grpc.ServerStream` embed. `command/mq_broker.go`
serves ZAP (unary `Dispatch` + 7-stream `StreamHandler`); `pb.WithBrokerGrpcClient`
dials ZAP via `brokerPool`. The last gRPC broker dial (`admin/dash/admin_server.go`
`UpdateTopicRetention`) now uses `pb.WithBrokerGrpcClient` too. Broker engine methods
(`mq/broker/broker_grpc_*.go`) keep their `mq_pb.HanzoMessaging_*Server` signatures —
those names are now plain `Send/Recv/Context` interfaces, not grpc generic streams.

gRPC is still **load-bearing** in (these legitimately import grpc — leave them):
- **Volume server RPC** — still gRPC-served (`command/volume.go`, `server/volume_grpc_*.go`,
  `pb/volume_server_pb/volume_server_grpc.pb.go`). Its `status.Errorf(codes.X)` producers
  must stay so gRPC clients get real codes; `peer.FromContext` works over gRPC.
- **Raft transport** — hashicorp/seaweedfs raft rides the master's gRPC listener
  (`server/raft_hashicorp.go`, the `pb.NewGrpcServer` + `reflection.Register` in
  `command/master.go` and `master_follower.go`). The master's own service is ZAP; the
  gRPC listener carries raft only.
- **Strangler client seam** `volumezap/client.go` — implements the
  `volume_server_pb.*Client` interface, whose generated signatures use `grpc.{CallOption,
  *StreamingClient}` and `metadata.MD`. Bound to grpc until that `*_pb` interface is
  regenerated grpc-free. (`mqzap/client.go` is now grpc-free — see above.)
- **Core dial/server + request-id** — `pb/grpc_client_server.go` (volume dials,
  `NewGrpcServer`), `operation/grpc_client.go`, `server/common.go` (grpc metadata).
- **TLS cert hot-reload** — `security/tls.go`, `security/certreload/certreload.go`,
  `security/tls_reload.go`, `util/http/client/http_client.go` reuse grpc's
  `credentials/tls/certprovider` + `pemfile` as a generic file-watching cert provider;
  `tls.go` also returns `grpc.ServerOption` for the remaining gRPC servers.

Full `go.mod` grpc removal is unreachable until volume RPC is ported to ZAP (regen its
`*_pb` interface grpc-free, switch the server to ZAP dispatch) and cert hot-reload drops
grpc's certprovider. (Filer, master, and MQ broker are already ported.)

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
