# HIP-1349 & HIP-1348 Research Report: PolarStore-Inspired Adaptive Object Compression for Hanzo S3

- **Tracking Issues**: `hanzoai/s3#28` (HIP-1349 — Adaptive Object Compression) and `hanzoai/s3#32` (HIP-1348 — Encryption-Boundary Safe Selective Compression)
- **Deliverable Package**: `s3/compression/polarstore/`
- **Reproducible Benchmark CLI**: `cmd/polarstore-bench/main.go`
- **Measured Artifacts**: [`score_table.md`](./score_table.md), [`raw_results.json`](./raw_results.json)

---

## 1. Paper & License Audit

### Source Paper
- **Title**: *PolarStore: High-Performance Data Compression for Large-Scale Cloud-Native Databases*
- **Authors**: Qingda Hu, Xinjun (Jimmy) Yang, Feifei Li, Junru Li, Ya Lin, Yuqi Zhou, Yicong Zhu, Junwei Zhang, Rongbiao Xie, Ling Zhou, Bin Wu, and Wenchao Zhou (Alibaba Cloud / PolarDB)
- **Venue**: Proceedings of the 24th USENIX Conference on File and Storage Technologies (FAST '26), February 24–26, 2026, Santa Clara, CA, USA, ISBN 978-1-939133-53-3
- **URL**: https://www.usenix.org/conference/fast26/presentation/hu (Open Access PDF: https://www.usenix.org/system/files/fast26-hu.pdf)

### Key Architectural Takeaways Adapted to Hanzo S3
PolarStore demonstrates across >100 PB of cloud-native storage that uniform compression policies fail because storage workloads mix latency-critical paths, highly compressible cold/warm data, and incompressible payloads. Specifically, PolarStore introduces:
1. **Tiered / workload-aware algorithm selection (Opt#2)**: selecting heavier codecs (Zstandard) for cold/warm space-sensitive data while reserving lightweight codecs or bypass for latency-critical hot paths.
2. **Critical-path bypass (Opt#1)**: immediately bypassing compression on latency-sensitive small writes rather than paying codec setup and amplification penalties.
3. **Dual-layer flexibility**: combining hardware NVMe CSD compression (PolarCSD gzip-5) with software-layer compression mode flags and adaptive algorithm selection.

### License & Clean-Room Audit
- **Paper Access**: USENIX Open Access; underlying algorithms and architectural insights are non-copyrightable scientific methods.
- **Code Provenance**: PolarStore's production Alibaba PolarDB/PolarCSD implementation is proprietary and coupled to custom NVMe CSD hardware. **Zero third-party code was copied or translated.**
- **Clean-Room Go Implementation**: `s3/compression/polarstore` is a 100% original, software-only Go implementation built exclusively on Go standard library (`compress/gzip`, `crypto/aes`, `crypto/cipher`, `crypto/sha256`, `sync`, `sync/atomic`) and existing Apache-2.0/BSD-3-Clause dependencies already pinned in `hanzoai/s3/go.mod` (`github.com/klauspost/compress/s2` and `github.com/klauspost/compress/zstd`, `v1.18.6`). Zero new external dependencies were added to `go.mod`.

### Audit of Existing Hanzo S3 Compression Path (`s3/util/compression.go` & `s3/operation/upload_content.go`)
Inspecting `upstream/main` revealed four concrete gaps that motivate this HIP-1349 / HIP-1348 deliverable:
1. **Single legacy codec**: `s3/util/compression.go` only activates `gzip.BestSpeed` (`zstd` support is commented out), leaving both Zstd's higher compression ratio and S2's multi-GB/s fast-tier speed unused.
2. **Broken extension heuristic & 128-byte probe**: `s3/operation/upload_content.go` calls `util.IsCompressableFileType(filepath.Base(option.Filename), option.MimeType)` — passing the basename instead of `filepath.Ext(...)`, so extension checks like `".json"` or `".jpg"` never match — and falls back to `util.GzipData(data[0:128])`, where an 18-byte gzip header/trailer (14% overhead) biases the 90% threshold against real payloads.
3. **No encryption-boundary policy**: `upload_content.go` performs `encrypt(gzip(data))` when `option.Cipher` is true without distinguishing `ModePreEncryptionPlaintext`, `ModeClientSideCiphertext` (HIP-0302 zero-knowledge storage), or `ModeAttackerControlledMixedSecret` (CRIME/BREACH length-oracle risk).
4. **No per-tenant policy isolation or atomic telemetry**: no lock-free counters tracking bypass reasons, codec mix, or probe overhead.

---

## 2. Matched Baselines

All strategies are evaluated identically by `cmd/polarstore-bench` on the same deterministic seeded corpus (`seed=20260224`, 32 objects/workload x 6 workloads = 192 objects/round = `20.01 MiB`/round, 1 warmup round + 5 measured rounds = `100.05 MiB` processed per strategy):

| Strategy | Description |
| :--- | :--- |
| `uncompressed_baseline` | Identity pass-through storing raw bytes (`CodecNone`). Lower bound on CPU/latency; 0% space savings. |
| `always_gzip_baseline` | Blindly compresses every object with pooled `compress/gzip` (`DefaultCompression`) regardless of entropy, size, or encryption boundary. |
| `always_zstd_baseline` | Blindly compresses every object with pooled `klauspost/compress/zstd` (`SpeedDefault`) regardless of entropy, size, or encryption boundary. |
| `polarstore_adaptive` | Clean-room `polarstore.AdaptiveCompressor` with fail-closed `EncryptionBoundaryPolicy`, O(1) pre-compressed magic sniff, deterministic 8 KiB multi-window Shannon entropy screen, fast S2 trial probe, and access-tier/CPU-budget codec selection (`none` / `s2` / `zstd`). |
| `post_encryption_negative` | **Documented Negative Control #1**: encrypts every object with AES-256-GCM *first*, then attempts Zstd compression on ciphertext. Empirically proves post-encryption compression yields `<= 1.000x` ratio (net storage expansion) while wasting CPU and heap allocations. |

**Documented Negative Control #2 (`crime_breach_mixed_secret_oracle`)**: compares naive compress-then-encrypt vs `polarstore_adaptive` (`ModeAttackerControlledMixedSecret`) across 100 trials of matching vs mismatching attacker probes concatenated with a fixed marker prior to AES-256-GCM encryption.

**Six Realistic S3 Workloads**:
1. `json_logs` (`128 KiB`, `ModePreEncryptionPlaintext`, `TierWarm`): structured JSON service access logs.
2. `source_tarball` (`128 KiB`, `ModePreEncryptionPlaintext`, `TierCold`): POSIX tar-aligned Go source archive.
3. `columnar_telemetry` (`128 KiB`, `ModePreEncryptionPlaintext`, `TierHot`): Parquet-like 32-byte numeric telemetry rows.
4. `precompressed_media` (`128 KiB`, `ModePreEncryptionPlaintext`, `TierWarm`): JPEG-magic entropy blobs and real Zstd archives.
5. `aes_gcm_ciphertext_zk` (`128 KiB`, `ModeClientSideCiphertext`, `TierWarm`): HIP-0302 client-side AES-256-GCM ciphertext.
6. `tiny_metadata` (`256 B`, `ModePreEncryptionPlaintext`, `TierHot`): micro JSON object manifests.

---

## 3. Falsifiable Hypothesis

> **H1 (Throughput & CPU Efficiency)**: On a mixed S3 object corpus containing compressible plaintext, latency-sensitive hot telemetry, pre-compressed media, tiny manifests, and HIP-0302 client-side ciphertext, `polarstore_adaptive` will achieve **>= 1.8x higher aggregate compression throughput** and **>= 45% lower compression CPU time** than `always_zstd_baseline` (and **>= 3.0x higher throughput / >= 65% lower CPU** than `always_gzip_baseline`), while retaining **>= 90% of `always_zstd_baseline`'s total byte savings** and maintaining **100% SHA-256 bit-exact round-trip fidelity**.
>
> **H2 (Hard Encryption-Boundary Safety)**: `polarstore_adaptive` will enforce **100% bypass** on `ModeClientSideCiphertext` (HIP-0302 ZK), `ModeUnknown` (fail-closed), and `ModeAttackerControlledMixedSecret` (CRIME/BREACH oracle delta = `0 B`), AND its Stage 2 Shannon entropy screen will reject **100% of unlabeled AES-256-GCM ciphertext** in **< 50 us/object** with **0** trial compressions.
>
> **Falsification Criteria**: Any single SHA-256 mismatch, any `Codec != CodecNone` on `ModeClientSideCiphertext`/`ModeAttackerControlledMixedSecret`/`ModeUnknown`, any non-zero CRIME/BREACH ciphertext length delta under `ModeAttackerControlledMixedSecret`, any storage expansion (`stored > original`), or aggregate compression throughput `< 1.8x` of `always_zstd_baseline` falsifies the hypothesis.

**Outcome**: **CONFIRMED** on all quantitative thresholds (see Section 6).

---

## 4. Bounded PR Deliverable

This PR is strictly scoped to three self-contained additions with zero modifications to unrelated subsystems and zero new `go.mod` dependencies:

1. **Core Package (`s3/compression/polarstore/`)**:
   - `codec.go`: Pooled, concurrency-safe `CodecNone`, `CodecS2`, `CodecZstd`, `CodecGzip` encoders/decoders (`WithEncoderConcurrency(1)` / `WithDecoderConcurrency(1)` to avoid goroutine leaks) + 10-byte self-describing `PLST` envelope framing (`MarshalEnvelope` / `UnmarshalEnvelope`).
   - `policy.go`: Explicit `DataPlaneMode` (`ModeUnknown` fail-closed, `ModePreEncryptionPlaintext`, `ModeClientSideCiphertext`, `ModeAttackerControlledMixedSecret`), `AccessTier` (`TierHot`, `TierWarm`, `TierCold`), auditable `BypassReason`, and `EncryptionBoundaryPolicy` with non-weakening `EnforceSecurityFloor()`.
   - `probe.go`: Multi-stage compressibility probe: (1) O(1) container magic sniff (JPEG/PNG/Gzip/Zstd/ZIP/XZ/Bzip2/LZ4/S2/WebP/MP4/Ogg/FLAC), (2) deterministic head+mid+tail 8 KiB sample with zero-allocation `[256]uint32` stack histogram Shannon entropy + byte-diversity screen, (3) fast S2 trial probe, and (4) PolarStore-inspired tier & CPU-budget codec selector.
   - `compressor.go`: `AdaptiveCompressor` with lock-free `atomic.Int64` stats, `sync.RWMutex` per-tenant policy overrides clamped to the global security floor, non-expansion guarantee, and `CompressAndSeal` / `OpenAndDecompress` binding `TenantID|Bucket|Key|Codec|DataPlane` into AES-256-GCM AAD.
   - `workloads.go`: Deterministic seeded S3 corpus generator shared by unit tests and benchmark runner.
   - `compressor_test.go`: 9 test suites verified under `go test -race -v -count=1`.
2. **Benchmark Runner (`cmd/polarstore-bench/main.go`)**:
   - Runnable via `go run ./cmd/polarstore-bench -seed 20260224 -objects 32 -rounds 5` to reproduce wall-clock latency, `syscall.Getrusage` CPU time, `runtime.MemStats` heap allocations, p50/p95/p99 Put latency, SHA-256 fidelity, CRIME/BREACH oracle test, and defense-in-depth entropy check.
3. **Research Evidence (`docs/research/hip-1349-polarstore/`)**:
   - `README.md` (this report), `score_table.md`, and `raw_results.json`.

---

## 5. Correctness, Concurrency, Multi-Tenant & Encryption-Boundary Safety Evidence

### Verification Commands Executed
```bash
go vet ./s3/compression/polarstore/... ./cmd/polarstore-bench/...
go test -race -v -count=1 ./s3/compression/polarstore/...
go run ./cmd/polarstore-bench -seed 20260224 -objects 32 -rounds 5 \
  -json-out docs/research/hip-1349-polarstore/raw_results.json \
  -md-out docs/research/hip-1349-polarstore/score_table.md
```

### Verified Safety & Correctness Invariants
1. **Bit-Exact SHA-256 Fidelity**: `TestRoundTripSHA256FidelityAllCodecsAndSizes` verifies all 4 codecs and all 3 access tiers across 8 sizes (`0 B` to `1 MiB`), and `cmd/polarstore-bench` verifies **100% (`960 / 960` object-rounds per strategy)** bit-exact SHA-256 restoration across all 5 strategies.
2. **Fail-Closed Encryption Boundary (HIP-0302 & CRIME/BREACH)**:
   - `TestEncryptionBoundaryFailClosedAndTenantFloor` proves `ModeUnknown`, `ModeClientSideCiphertext`, and `ModeAttackerControlledMixedSecret` unconditionally bypass with `CodecNone`, exact `BypassReason`, and `RanTrialProbe == false` — even when a tenant attempts to submit a weakened policy (`EnforceSecurityFloor()` clamps all three security flags back to `true`).
   - `TestCrimeBreachLengthOracleMitigation` and `cmd/polarstore-bench` empirically prove naive compress-then-encrypt leaks a **`12 B`** length differential (`134 B` match vs `146 B` mismatch on a `1952 B` payload), while `ModeAttackerControlledMixedSecret` eliminates the oracle completely (**`1980 B` vs `1980 B`**, **`0 B` delta**).
   - `TestPostEncryptionNegativeControl` and `DefenseInDepthEntropy` confirm that even if raw AES-256-GCM ciphertext is accidentally labeled as `ModePreEncryptionPlaintext`, Stage 2 Shannon entropy (`7.977 bits/B` mean) rejects **32/32 (100%)** objects in **`15.04 us`** mean latency (**`8,716.1 MB/s`**) with **`0`** trial compressions.
3. **Authenticated Codec Binding (AAD Integrity)**: `TestCompressAndSealAADTamperDetection` verifies that mutating either `Codec` (downgrade/confusion attack) or `TenantID` (cross-tenant replay) causes AES-256-GCM authentication failure during `OpenAndDecompress`.
4. **Concurrency & Multi-Tenant Race Safety**: `TestConcurrentMultiTenantRaceSafety` runs 16 parallel goroutines executing concurrent `CompressAndSeal`, `OpenAndDecompress`, and dynamic `SetTenantPolicy` updates across distinct tenants under Go's race detector (`go test -race`) with **0 data races** and exact atomic accounting.

---

## 6. Before/After Measurements & HIP-1348/1349 Integration Recommendation

### Exact Measured Results (`linux/amd64`, `go1.27.1`, `2` vCPUs, seed `20260224`, `100.05 MiB` / strategy)

| Strategy | Compression Ratio | Space Savings | Compress Thru (MB/s) | Decompress Thru (MB/s) | Compress Wall (ms) | Compress CPU (ms) | p50 / p99 Put (us) | Alloc / Round (MiB) | ZK Ciphertext Bypass | SHA-256 Fidelity |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `uncompressed_baseline` | `1.000x` | `0.00%` | `272215.7` | `877596.0` | `0.39` | `0.49` | `0.0` / `0.1` | `0.01` | `100%` | 100% PASS |
| `always_gzip_baseline` | `1.898x` | `47.32%` | `125.4` | `268.3` | `836.82` | `845.88` | `530.2` / `3998.2` | `21.92` | `0%` | 100% PASS |
| `always_zstd_baseline` | `1.924x` | `48.01%` | `221.5` | `885.4` | `473.57` | `482.11` | `108.8` / `2260.2` | `35.94` | `0%` | 100% PASS |
| `polarstore_adaptive` | `1.782x` | `43.89%` | `478.9` | `1355.3` | `219.05` | `224.22` | `1.1` / `1258.7` | `10.36` | `100%` | 100% PASS |
| `post_encryption_negative` | `1.000x` | `-0.04%` | `493.0` | `644.5` | `212.77` | `224.48` | `170.1` / `2637.6` | `93.21` | `0%` | 100% PASS |

### Key Measured Deltas
- **vs. `always_gzip_baseline`** (closest to legacy Hanzo S3 gzip): **3.82x higher compression throughput** (`478.9` vs `125.4 MB/s`), **5.05x higher decompression throughput** (`1355.3` vs `268.3 MB/s`), **73.5% lower compression CPU** (`224.22` vs `845.88 ms`), **52.7% fewer heap allocations** (`10.36` vs `21.92 MiB/round`), **3.18x lower p99 Put latency** (`1258.7` vs `3998.2 us`), and **100% HIP-0302 ZK bypass** (vs `0%`).
- **vs. `always_zstd_baseline`**: **2.16x higher compression throughput** (`478.9` vs `221.5 MB/s`), **1.53x higher decompression throughput** (`1355.3` vs `885.4 MB/s`), **53.5% lower compression CPU** (`224.22` vs `482.11 ms`), **71.2% fewer heap allocations** (`10.36` vs `35.94 MiB/round`), **1.80x lower p99 Put latency** (`1258.7` vs `2260.2 us`), while preserving **91.4% of total space savings** (`43.89%` vs `48.01%`) and trading only the hot-tier telemetry ratio (`1.58x` S2 at `407 MB/s` vs `2.34x` Zstd at `77 MB/s` — a **5.29x hot-path speedup**).
- **vs. `post_encryption_negative`**: Post-encryption Zstd delivers **`-0.04%`** net space savings (`0.86x` / +16% bloat on tiny objects) while burning `224.48 ms` CPU and `93.21 MiB/round` of allocations (9.0x more memory churn than `polarstore_adaptive`), decisively ruling out post-encryption compression.

### HIP-1348 & HIP-1349 Integration Recommendation for Hanzo S3
1. **Adopt `s3/compression/polarstore` behind a feature flag in `s3/operation/upload_content.go` and `s3/s3api/`**:
   - Replace the broken `filepath.Base(...)` + 128-byte `util.GzipData` probe in `Uploader.doUploadData` with `AdaptiveCompressor.Compress`.
   - Map S3 request encryption headers explicitly to `DataPlaneMode`:
     - Client-side encrypted / HIP-0302 ZK objects (e.g., `x-amz-meta-x-amz-key*`, HIP-0302 ZK envelope flag) -> `ModeClientSideCiphertext` (hard bypass).
     - Buckets/objects flagged as containing attacker-influenced mixed secrets (or default-safe policy for SSE-C/SSE-KMS unless tenant explicitly opts into pre-encryption compression for homogeneous telemetry/log buckets) -> `ModeAttackerControlledMixedSecret` (hard bypass).
     - Explicitly opted-in server-side plaintext buckets (logs, CI artifacts, columnar analytics, unencrypted or homogeneous SSE-S3/SSE-KMS) -> `ModePreEncryptionPlaintext`.
2. **Persist Codec in Authenticated Chunk Metadata**: Store `Codec.String()` in `FileChunk` / S3 object system metadata and include it in SSE AAD (`s3/s3api/s3_sse_s3.go`, `s3/s3api/s3_sse_kms.go`) so legacy `gzip` chunks continue decoding seamlessly while new objects gain adaptive `s2`/`zstd`/`none` handling.
3. **Export Prometheus Counters**: Wire `AdaptiveCompressor.Stats()` into `s3/stats` to track per-bucket compression ratio, codec mix (`none`/`s2`/`zstd`), and `BypassReason` rates in production.
