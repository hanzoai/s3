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

## Master consensus: leaderless pinned-writer (raft deleted)
hashicorp/seaweedfs **raft is gone** (commit `8ac2c5465`) — no raft dep in `go.mod`/
`go.sum`, no `server/raft_hashicorp.go`, no election. Coordination is a leaderless,
deterministic **pinned writer**: the `topology.Coordinator` interface serializes the two
non-commutative allocations (volume ids, file ids) and pins exactly ONE writer per
cluster as a *pure function of the membership* (lowest-address live member), so every
master computes the same writer with zero rounds and re-pins the instant the live set
changes. Two impls, selected by env, never both:
- **`topology.LocalCoordinator`** (default) — in-process pinned-writer over the `-peers`
  universe, projected to the LIVE subset by `server.publishLiveMembers`. Standalone
  master = sole writer. `ErrNotWriter` replaces `raft.NotLeaderError` (clients back off +
  redirect to the advertised `Writer`).
- **Writer failure detector** — `server.ReconcileWriterMembership` is a 2s per-master
  ticker (multi-master only) that re-probes every peer via `pb.WithMasterClient` (pooled,
  honors grpc.master PQ-TLS — a raw plaintext probe would false-negative under mTLS and
  split-brain) and republishes the live set. This is the piece raft's heartbeat/election
  provided: peer-LEAVE events are disseminated by the writer, so a dead WRITER is
  invisible to followers via `OnPeerUpdate` alone; the independent probe lets every
  survivor notice it and deterministically re-pin the lowest live address (verified end
  to end: writer re-pins in ~2s, `test/multi_master`). `OnPeerUpdate` and the reconciler
  share `pingMaster`/`publishLiveMembers` — one liveness check, one publish path.
- **`s3server.schainCoordinator`** (production, `master.coordinator.schain_endpoint`) —
  composes `LocalCoordinator` for writer selection and routes id allocation through the
  schain storage VM's owner-gated, ML-DSA-verified `AllocateTx` over ZAP (single-writer
  *per range*, durable). **Fails CLOSED** (`ErrSchainNotWired`) until the schain ZAP
  service descriptor from the separate `luxfi/chains` module is wired in — never invents
  an id the chain has not agreed. **Remaining seam #1** (production durability).
- `/cluster/status` reports `IsWriter`/`Writer`/`Members` (not `IsLeader`). The
  `Raft{List,Add,Remove}*`/`RaftLeadershipTransfer` RPCs are retained as **advisory,
  leaderless** membership hints (wire-compat with the admin dashboard); transfer is
  refused (the writer is pinned, not hand-transferable).
- **Availability without quorum**: losing 2 of 3 masters does NOT block writes — the lone
  survivor drops the dead peers and pins itself. This is a deliberate trade vs raft's
  quorum; cross-partition id-collision safety is the schain VM's job, not the local pin.
- **TopologyId** is minted only by the writer (`EnsureTopologyId`) and is neither
  propagated to followers nor persisted; followers surface it by proxying `/dir/status`
  to the writer. So it is consistent at any settled moment but a new writer promoted on
  failover mints a FRESH id (not preserved across a writer change). Not a data-correctness
  issue — id uniqueness is the Coordinator's job. **Remaining seam #2**: preserve it by
  adding a field to `GetMasterConfigurationResponse`/`VolumeLocation`, having followers
  adopt it, and persisting it to `-mdir`.
- Tests: fast unit invariants in `s3/topology/coordinator_test.go` +
  `s3/server/{schain_coordinator,master_grpc_server_membership}_test.go` (deterministic
  pin, non-writer refusal, concurrent-allocation uniqueness, fail-closed schain,
  membership RPCs — 15 tests). The end-to-end suite `test/multi_master/` (writer failover,
  survivor-keeps-writing-without-quorum, full-restart, writer-consistency) is gated behind
  `MULTI_MASTER_IT=1` and builds the master `CGO_ENABLED=0 GOWORK=off`.

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
