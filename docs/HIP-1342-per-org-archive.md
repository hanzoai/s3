# HIP-1342: Per-Org Scheduled Incremental Archive to Remote S3

## 1. Architecture and Design Rationale

`weed filer.backup` runs as a continuous cluster-wide event-log consumer over a single `-filerPath` prefix without per-tenant schedule cadences (`off`, `hourly`, `daily`, `weekly`), per-bucket byte/request billing counters, or verified point-in-time / post-retention bucket reconstruction. Under HIP-1342, multi-tenant Hanzo S3 organizations configure scheduled incremental cold-archive passes per bucket (`/buckets/<bucket>/`) into a remote S3 archive bucket while keeping hot NVMe volumes independent of cold-tier failures.

Key invariants implemented in `s3/filer/archive` and exposed via `weed filer.archive` (`s3/command/filer_archive.go`):

1. **Strict single-bucket tenant isolation**: `archive.ValidateBucketName` and `archive.SanitizeObjectRelativeKey` reject empty names, `.` / `..` traversal, path separators, backslashes, NUL bytes, and non-printable ASCII before any filer read or remote S3 key construction, ensuring a policy on `/buckets/<org>/` can never read or overwrite sibling tenant prefixes.
2. **Linearizable watermark checkpointing**: Per-bucket checkpoints (`archive.<fnv32a(bucket->remoteBucket:prefix)>`) advance strictly to `maxEventTsNs` (the highest committed metadata mutation timestamp durably uploaded during the pass), never to an uncommitted wall-clock timestamp and never on partial remote `503`/network failure. Concurrent writers committing in parallel with an archive pass are guaranteed never to be skipped by the next pass.
3. **Dual-layout incremental log + live mirror with date retention**:
   - `data/YYYY-MM-DD/<ts>_<seq>/<key>` and `log/YYYY-MM-DD/<ts>_<seq>.json` record immutable incremental mutations with SHA-256 payload digests.
   - `current/<key>` maintains the latest live mirror (propagating deletes when `PropagateDel = true`), allowing `PruneRetention` to delete expired `data/YYYY-MM-DD/` and `log/YYYY-MM-DD/` partitions older than `RetentionDays` while preserving full `RestoreBucket` capability for long-lived untouched objects.
4. **Per-org / per-bucket metering**: Lock-free `atomic.Int64` counters track `ArchivedBytes`, `RestoredBytes`, `PutRequests`, `DeleteRequests`, `ListRequests`, `GetRequests`, `PassesCompleted`, `PassesFailed`, `RetainedPrunedKeys`, `LastOffsetTsNs`, `LastLagNs`, and `LastPassDurationNs` per bucket.

## 2. Matched Baseline and Falsifiable Disaster-Recovery Protocol

Following the HIP-1342 benchmark contract, `TestPerOrgArchiveHotVolumeKillAndByteIdenticalRestoreWithConcurrentWorkload` (`s3/filer/archive/archive_test.go`) runs on identical `32 x 256 KiB` (`8.00 MiB`) deterministic binary shards:

1. **Baseline**: Sequential hot-tier NVMe file creation + `fsync()` (`os.File.Sync()`) on the local volume directory.
2. **Concurrent hot workload + scheduled incremental passes**: A background writer streams 32 shards plus an interleaving overwrite (`shard-002.bin`) and delete (`shard-001.bin`) while 3 concurrent incremental archive passes execute mid-flight, followed by 1 final catch-up pass (`lag_ns = 0`).
3. **Simulated hot-volume loss (`KillHotVolume`) + full restore (`RestoreBucket`)**: All hot bucket state and metadata logs are wiped (`0` live keys), `RestoreBucket` reconstructs the bucket from the remote S3 archive, and every restored object (`31` live shards; deleted `shard-001.bin` absent; overwritten `shard-002.bin` reflecting its latest payload) is verified byte-for-byte via `SHA-256` under `go test -race`.

## 3. Raw Measured Results (`GOWORK=off go test -v -race -count=1 ./s3/filer/archive/...`)

| Metric | Measured Value |
| :--- | :---: |
| Plain local NVMe `fsync` baseline (`32 x 256 KiB`) | `28.23 MiB/s` |
| Concurrent incremental archive pass throughput | `32.12 MiB/s` |
| Full bucket restore throughput (`RestoreBucket` + SHA-256 check) | `32.68 MiB/s` |
| Metered archived bytes (`data/` + `current/` + `log/`) | `17,311,347 B` (`16.51 MiB`) |
| Metered restored bytes (`31 x 256 KiB` live shards) | `8,126,464 B` (`7.75 MiB`) |
| Metered remote requests (`PUT` / `DELETE` / `GET` / `LIST`) | `100` / `1` / `65` / `2` |
| Post-pass replication lag (`LastLagNs`) | `0 ns` |
| SHA-256 byte-identical verification (`31/31` live + `0` deleted ghosts) | `100% PASS (-race)` |

## 4. CLI Usage (`weed filer.archive`)

```bash
# Single bounded incremental pass (e.g. Kubernetes CronJob or external scheduler)
weed filer.archive \
  -filer=localhost:8888 \
  -bucket=org-prod \
  -remoteBucket=hanzo-archive-s3 \
  -remotePrefix=hip1342 \
  -schedule=daily \
  -retentionDays=30 \
  -once

# Continuous scheduled daemon (off | hourly | daily | weekly)
weed filer.archive \
  -filer=localhost:8888 \
  -bucket=org-prod \
  -remoteBucket=hanzo-archive-s3 \
  -schedule=hourly \
  -retentionDays=30
```
