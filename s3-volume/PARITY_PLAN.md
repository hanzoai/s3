# Rust Volume Server Parity Plan

Generated: 2026-03-16

## Goal

Make `hanzo-volume` a drop-in replacement for the Go volume server by:

- comparing every Go volume-server code path against the Rust implementation,
- recording file-level ownership and verification status,
- closing verified behavior gaps one logic change per commit,
- extending tests so regressions are caught by Go parity suites and Rust unit/integration tests.

## Ground Truth

Primary Go sources:

- `s3/server/volume_server.go`
- `s3/server/volume_server_handlers*.go`
- `s3/server/volume_grpc_*.go`
- `s3/server/constants/volume.go`
- `s3/storage/store*.go`
- `s3/storage/disk_location*.go`
- `s3/storage/volume*.go`
- `s3/storage/needle/*.go`
- `s3/storage/idx/*.go`
- `s3/storage/needle_map*.go`
- `s3/storage/needle_map/*.go`
- `s3/storage/super_block/*.go`
- `s3/storage/erasure_coding/*.go`

Supporting Go dependencies that affect drop-in behavior:

- `s3/command/volume.go`
- `s3/security/*.go`
- `s3/images/*.go`
- `s3/stats/*.go`

Primary Rust sources:

- `hanzo-volume/src/main.rs`
- `hanzo-volume/src/config.rs`
- `hanzo-volume/src/security.rs`
- `hanzo-volume/src/images.rs`
- `hanzo-volume/src/server/*.rs`
- `hanzo-volume/src/storage/*.rs`
- `hanzo-volume/src/storage/needle/*.rs`
- `hanzo-volume/src/storage/idx/*.rs`
- `hanzo-volume/src/storage/erasure_coding/*.rs`
- `hanzo-volume/src/remote_storage/*.rs`

## Audit Method

For each Go file:

1. Map it to the Rust file or files that should own the same behavior.
2. Compare exported entry points, helper functions, state transitions, wire fields, and persistence side effects.
3. Mark each file `implemented`, `partial`, `missing`, or `needs verification`.
4. Link each behavior to an existing test or add a missing test.
5. Only treat a gap as closed after code review plus local verification.

## Acceptance Criteria

The Rust server is a drop-in replacement only when all of these hold:

- HTTP routes, status codes, headers, and body semantics match Go.
- gRPC RPCs match Go request validation, response fields, streaming behavior, and maintenance/read-only semantics.
- Master heartbeat and topology metadata match Go closely enough that the Go master treats Rust and Go volume servers the same.
- On-disk volume behavior matches Go for normal volumes, EC shards, tiering metadata, and readonly persistence.
- Startup flags and operational endpoints that affect production deployment behave equivalently or are explicitly documented as unsupported.
- Existing Go integration suites pass with `VOLUME_SERVER_IMPL=rust`.

## File Matrix

### HTTP server surface

| Go file | Rust counterpart | Status | Comparison focus |
| --- | --- | --- | --- |
| `s3/server/volume_server.go` | `hanzo-volume/src/main.rs`, `hanzo-volume/src/server/volume_server.rs`, `hanzo-volume/src/server/heartbeat.rs` | partial | startup wiring, routers, heartbeat, shutdown, metrics/debug listeners |
| `s3/server/volume_server_handlers.go` | `hanzo-volume/src/server/volume_server.rs`, `hanzo-volume/src/server/handlers.rs` | needs verification | method dispatch, OPTIONS behavior, public/admin split |
| `s3/server/volume_server_handlers_admin.go` | `hanzo-volume/src/server/handlers.rs` | implemented | `/status`, `/healthz`, stats, server headers |
| `s3/server/volume_server_handlers_helper.go` | `hanzo-volume/src/server/handlers.rs` | needs verification | JSON encoding, request parsing, helper parity |
| `s3/server/volume_server_handlers_read.go` | `hanzo-volume/src/server/handlers.rs` | needs verification | JWT, conditional reads, range reads, proxy/redirect, chunk manifests, image transforms |
| `s3/server/volume_server_handlers_ui.go` | `hanzo-volume/src/server/handlers.rs`, embedded assets | partial | UI payload and HTML parity |
| `s3/server/volume_server_handlers_write.go` | `hanzo-volume/src/server/handlers.rs`, `hanzo-volume/src/images.rs` | needs verification | multipart parsing, metadata, compression, ts, delete semantics |
| `s3/server/constants/volume.go` | `hanzo-volume/src/server/heartbeat.rs`, config defaults | needs verification | heartbeat timing, constants parity |

### gRPC server surface

| Go file | Rust counterpart | Status | Comparison focus |
| --- | --- | --- | --- |
| `s3/server/volume_grpc_admin.go` | `hanzo-volume/src/server/grpc_server.rs` | needs verification | readonly/writable, allocate/delete/configure/mount/unmount |
| `s3/server/volume_grpc_batch_delete.go` | `hanzo-volume/src/server/grpc_server.rs` | implemented | batch delete, EC delete path |
| `s3/server/volume_grpc_client_to_master.go` | `hanzo-volume/src/server/heartbeat.rs` | partial | heartbeat fields, leader changes, metrics settings from master |
| `s3/server/volume_grpc_copy.go` | `hanzo-volume/src/server/grpc_server.rs` | needs verification | full copy streams |
| `s3/server/volume_grpc_copy_incremental.go` | `hanzo-volume/src/server/grpc_server.rs` | needs verification | incremental copy binary search, timestamps |
| `s3/server/volume_grpc_erasure_coding.go` | `hanzo-volume/src/server/grpc_server.rs`, `hanzo-volume/src/storage/erasure_coding/*.rs` | needs verification | shard read/write/delete/mount/unmount/rebuild |
| `s3/server/volume_grpc_query.go` | `hanzo-volume/src/server/grpc_server.rs` | needs verification | query validation and error parity |
| `s3/server/volume_grpc_read_all.go` | `hanzo-volume/src/server/grpc_server.rs` | needs verification | read-all ordering and tail semantics |
| `s3/server/volume_grpc_read_write.go` | `hanzo-volume/src/server/grpc_server.rs`, `hanzo-volume/src/storage/*.rs` | needs verification | blob/meta/page reads, write blob semantics |
| `s3/server/volume_grpc_remote.go` | `hanzo-volume/src/server/grpc_server.rs`, `hanzo-volume/src/remote_storage/*.rs` | needs verification | remote fetch/write and tier metadata |
| `s3/server/volume_grpc_scrub.go` | `hanzo-volume/src/server/grpc_server.rs`, `hanzo-volume/src/storage/*.rs` | needs verification | scrub result semantics |
| `s3/server/volume_grpc_state.go` | `hanzo-volume/src/server/grpc_server.rs` | implemented | GetState/SetState/Status |
| `s3/server/volume_grpc_tail.go` | `hanzo-volume/src/server/grpc_server.rs` | needs verification | tail streaming and idle timeout |
| `s3/server/volume_grpc_tier_download.go` | `hanzo-volume/src/server/grpc_server.rs`, `hanzo-volume/src/remote_storage/*.rs` | needs verification | tier download stream/error paths |
| `s3/server/volume_grpc_tier_upload.go` | `hanzo-volume/src/server/grpc_server.rs`, `hanzo-volume/src/remote_storage/*.rs` | needs verification | tier upload stream/error paths |
| `s3/server/volume_grpc_vacuum.go` | `hanzo-volume/src/server/grpc_server.rs`, `hanzo-volume/src/storage/*.rs` | needs verification | compact/commit/cleanup progress and readonly transitions |

### Storage and persistence surface

| Go file group | Rust counterpart | Status | Comparison focus |
| --- | --- | --- | --- |
| `s3/storage/store.go`, `store_state.go` | `hanzo-volume/src/storage/store.rs`, `hanzo-volume/src/server/heartbeat.rs` | partial | topology metadata, disk tags, server id, state persistence |
| `s3/storage/store_vacuum.go` | `hanzo-volume/src/storage/store.rs`, `hanzo-volume/src/storage/volume.rs` | needs verification | vacuum sequencing |
| `s3/storage/store_ec.go`, `store_ec_delete.go`, `store_ec_scrub.go` | `hanzo-volume/src/storage/store.rs`, `hanzo-volume/src/storage/erasure_coding/*.rs` | needs verification | EC lifecycle and scrub behavior |
| `s3/storage/disk_location.go`, `disk_location_ec.go` | `hanzo-volume/src/storage/disk_location.rs`, `hanzo-volume/src/storage/store.rs` | partial | directory UUIDs, tags, load rules, disk space checks |
| `s3/storage/volume.go`, `volume_loading.go` | `hanzo-volume/src/storage/volume.rs` | needs verification | load/reload/readonly/remote metadata |
| `s3/storage/volume_super_block.go` | `hanzo-volume/src/storage/super_block.rs`, `hanzo-volume/src/storage/volume.rs` | implemented | super block parity |
| `s3/storage/volume_read.go`, `volume_read_all.go` | `hanzo-volume/src/storage/volume.rs`, `hanzo-volume/src/server/handlers.rs` | needs verification | full/meta/page reads, TTL, streaming |
| `s3/storage/volume_write.go` | `hanzo-volume/src/storage/volume.rs`, `hanzo-volume/src/server/write_queue.rs` | needs verification | dedup, sync/async writes, metadata flags |
| `s3/storage/volume_vacuum.go` | `hanzo-volume/src/storage/volume.rs` | needs verification | compact and commit parity |
| `s3/storage/volume_backup.go` | `hanzo-volume/src/storage/volume.rs`, `hanzo-volume/src/server/grpc_server.rs` | needs verification | backup/search logic |
| `s3/storage/volume_checking.go` | `hanzo-volume/src/storage/volume.rs`, `hanzo-volume/src/storage/idx/mod.rs`, `hanzo-volume/src/server/grpc_server.rs` | needs verification | scrub and integrity checks |
| `s3/storage/volume_info.go`, `volume_info/volume_info.go`, `volume_tier.go` | `hanzo-volume/src/storage/volume.rs`, `hanzo-volume/src/remote_storage/*.rs` | needs verification | `.vif` format and tiered file metadata |
| `s3/storage/needle/*.go` | `hanzo-volume/src/storage/needle/*.rs` | needs verification | needle parsing, CRC, TTL, multipart metadata |
| `s3/storage/idx/*.go` | `hanzo-volume/src/storage/idx/*.rs` | needs verification | index walking and binary search |
| `s3/storage/needle_map*.go`, `needle_map/*.go` | `hanzo-volume/src/storage/needle_map.rs` | needs verification | map kind parity, persistence, memory behavior |
| `s3/storage/super_block/*.go` | `hanzo-volume/src/storage/super_block.rs` | implemented | replica placement and TTL metadata |
| `s3/storage/erasure_coding/*.go` | `hanzo-volume/src/storage/erasure_coding/*.rs` | needs verification | EC shard placement, encode/decode, journal deletes |

### Supporting runtime surface

| Go file | Rust counterpart | Status | Comparison focus |
| --- | --- | --- | --- |
| `s3/command/volume.go` | `hanzo-volume/src/config.rs`, `hanzo-volume/src/main.rs` | partial | flags, metrics/debug listeners, startup behavior |
| `s3/security/*.go` | `hanzo-volume/src/security.rs`, `hanzo-volume/src/main.rs` | implemented | JWT and TLS loading |
| `s3/images/*.go` | `hanzo-volume/src/images.rs`, `hanzo-volume/src/server/handlers.rs` | implemented | JPEG orientation and transforms |
| `s3/stats/*.go` | `hanzo-volume/src/metrics.rs`, `hanzo-volume/src/server/handlers.rs` | partial | metrics endpoints, push-gateway integration |

## Verified Gaps As Of 2026-03-08

The startup/runtime gaps that were verified in the initial audit are now closed:

1. Heartbeat metadata parity
   Closed by `8ade1c51d` and retained in current HEAD.

2. Dedicated metrics/debug listener parity
   Closed by `fbe0e5829`.

3. Master-provided metrics push settings
   Closed by `fbe0e5829`.

4. Slow-read tuning parity
   Closed by `66e3900dc`.

There are no remaining verified gaps from the initial startup/runtime audit. The broader line-by-line comparison batches below are still required to either confirm parity or surface new gaps.

## Execution Status As Of 2026-03-16

The file-by-file comparison and verification work executed in this round was:

1. Startup and harness alignment
   Compared `s3/command/volume.go`, `test/volume_server/framework/cluster*.go`, `hanzo-volume/src/config.rs`, and `hanzo-volume/src/main.rs` to ensure the Rust server is invoked with Go-compatible flags and is rebuilt from the current source during parity runs.

2. HTTP admin surface
   Compared `s3/server/volume_server_handlers_admin.go` against `hanzo-volume/src/server/handlers.rs` with emphasis on `/status` payload shape, disk-status fields, and volume ordering.

3. gRPC admin surface
   Compared `s3/server/volume_grpc_admin.go` against `hanzo-volume/src/server/grpc_server.rs` with emphasis on `Ping`, `VolumeConfigure`, readonly/writable flows, and error wrapping.

4. Storage/index layout
   Compared Go index-entry defaults in `s3/storage/types` and `s3/storage/idx/*.go` against the Rust default feature set in `hanzo-volume/Cargo.toml` and the Rust index reader/writer paths to confirm default binaries use the same offset width.

5. End-to-end parity verification
   Re-ran the Go HTTP and gRPC integration suites with `VOLUME_SERVER_IMPL=rust` after each fix to confirm wire-level compatibility.

### Verified mismatches closed in this round

- Rust parity runs could reuse a stale `s3-volume` binary across test invocations, hiding source and feature changes from the Go harness.
- Rust defaulted to 5-byte index offsets, while the default Go `go build` path uses 4-byte offsets unless built with `-tags 5BytesOffset`.
- Rust `/status` omitted Go fields in both `Volumes` and `DiskStatuses`, and did not sort volumes by `Id`.
- Rust `Ping` treated an empty target as a self-ping and only performed a raw gRPC connect for filer targets; Go returns `remote_time_ns=0` for the empty request and performs a real filer `Ping` RPC.
- Rust `VolumeNeedleStatus` dropped stored TTL metadata and reported `data_size` instead of Go’s `Size` field.
- Rust multipart uploads ignored form fields such as `ts`, `ttl`, and `cm`, and also ignored part-level `Content-Encoding` and `Content-MD5`.
- Rust only treated `dl=true` and `dl=1` as truthy, while Go accepts the full `strconv.ParseBool` set such as `dl=t` and `dl=True`.

### Verification commands

- `VOLUME_SERVER_IMPL=rust go test -count=1 -timeout 1200s ./test/volume_server/http/...`
- `VOLUME_SERVER_IMPL=rust go test -count=1 -timeout 1200s ./test/volume_server/grpc/...`

## Execution Plan

### Batch 1: startup and heartbeat

- Compare `s3/command/volume.go`, `s3/server/volume_server.go`, `s3/server/volume_grpc_client_to_master.go`, `s3/storage/store.go`, and `s3/storage/disk_location.go`.
- Close metadata and startup parity gaps that affect master registration and deployment compatibility.
- Add Rust unit tests for heartbeat payloads and config wiring.

### Batch 2: HTTP read path

- Compare `volume_server_handlers_read.go`, `volume_server_handlers_helper.go`, and related storage read functions line by line.
- Verify JWT, path parsing, proxy/redirect, ranges, streaming, chunk manifests, image transforms, and response-header overrides.
- Extend `test/volume_server/http/...` and Rust handler tests where parity is not covered.

### Batch 3: HTTP write/delete path

- Compare `volume_server_handlers_write.go` and write-related storage functions.
- Verify multipart behavior, metadata, md5, compression, unchanged writes, delete edge cases, and timestamp handling.

### Batch 4: gRPC admin and lifecycle

- Compare `volume_grpc_admin.go`, `volume_grpc_state.go`, and `volume_grpc_vacuum.go`.
- Verify readonly/writable flows, maintenance mode, status payloads, mount/unmount/delete/configure, and vacuum transitions.

### Batch 5: gRPC data movement

- Compare `volume_grpc_read_write.go`, `copy*.go`, `read_all.go`, `tail.go`, `remote.go`, and `query.go`.
- Verify stream framing, binary search, idle timeout, and remote-storage semantics.

### Batch 6: storage internals

- Compare all `s3/storage` volume, needle, idx, needle map, and EC files line by line.
- Focus on persistence rules, readonly semantics, TTL, recovery/scrub, backup, and memory/disk map behavior.

## Commit Strategy

- One commit for the audit/plan document if the document itself changes.
- One commit per logic fix.
- Every logic commit must include the smallest test addition that proves the new parity claim.
