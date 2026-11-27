# HIP-1349 / HIP-1348 PolarStore Adaptive Object Compression — Measured Score Table

- **Environment**: `linux/amd64`, `go1.27.1`, `2` vCPUs
- **Corpus**: seed `20260224`, `32` objects/workload x 6 workloads = `192` objects/round (`20.01 MiB`/round), `5` measured rounds (after 1 warmup round)
- **Throughput unit**: decimal `MB/s` (`original_bytes / 1e6 / wall_seconds`)

## 1. Overall Mixed-Workload Comparison (5 Matched Strategies)

| Strategy | Compression Ratio | Space Savings | Compress Thru (MB/s) | Decompress Thru (MB/s) | Compress Wall (ms) | Compress CPU (ms) | p50 / p99 Put (us) | Alloc / Round (MiB) | ZK Ciphertext Bypass | SHA-256 Fidelity |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `uncompressed_baseline` | `1.000x` | `0.00%` | `272215.7` | `877596.0` | `0.39` | `0.49` | `0.0` / `0.1` | `0.01` | `100%` | 100% PASS |
| `always_gzip_baseline` | `1.898x` | `47.32%` | `125.4` | `268.3` | `836.82` | `845.88` | `530.2` / `3998.2` | `21.92` | `0%` | 100% PASS |
| `always_zstd_baseline` | `1.924x` | `48.01%` | `221.5` | `885.4` | `473.57` | `482.11` | `108.8` / `2260.2` | `35.94` | `0%` | 100% PASS |
| `polarstore_adaptive` | `1.782x` | `43.89%` | `478.9` | `1355.3` | `219.05` | `224.22` | `1.1` / `1258.7` | `10.36` | `100%` | 100% PASS |
| `post_encryption_negative` | `1.000x` | `-0.04%` | `493.0` | `644.5` | `212.77` | `224.48` | `170.1` / `2637.6` | `93.21` | `0%` | 100% PASS |

## 2. Per-Workload Breakdown (`polarstore_adaptive` vs Baselines)

| Workload | `uncompressed` Ratio / MB/s | `always_gzip` Ratio / MB/s | `always_zstd` Ratio / MB/s | `polarstore_adaptive` Ratio / MB/s | `polarstore_adaptive` Codec / Bypass | `post_encryption_negative` Ratio / MB/s |
| :--- | :---: | :---: | :---: | :---: | :--- | :---: |
| `json_logs` | `1.00x` / `1824086` | `6.56x` / `89` | `6.31x` / `161` | `6.31x` / `144` | `zstd:160` | `1.00x` / `336` |
| `source_tarball` | `1.00x` / `545707` | `47.21x` / `227` | `75.62x` / `1104` | `75.62x` / `985` | `zstd:160` | `1.00x` / `615` |
| `columnar_telemetry` | `1.00x` / `2596771` | `2.17x` / `50` | `2.34x` / `77` | `1.58x` / `407` | `s2:160` | `1.00x` / `505` |
| `precompressed_media` | `1.00x` / `2724448` | `1.00x` / `1232` | `1.00x` / `656` | `1.00x` / `210005` | `none:160 (precompressed_magic_header:160)` | `1.00x` / `605` |
| `aes_gcm_ciphertext_zk` | `1.00x` / `2652168` | `1.00x` / `335` | `1.00x` / `1220` | `1.00x` / `795661` | `none:160 (client_side_ciphertext_hip0302:160)` | `1.00x` / `535` |
| `tiny_metadata` | `1.00x` / `4683` | `1.88x` / `14` | `1.92x` / `28` | `1.00x` / `914` | `none:160 (below_min_size_threshold:160)` | `0.86x` / `44` |

## 3. CRIME / BREACH Mixed-Secret Oracle Negative Control

| Regime | Plaintext Bytes | Match Ciphertext (B) | Mismatch Ciphertext (B) | Length Delta (B) | Oracle Distinguishable? | Bypass Reason |
| :--- | :---: | :---: | :---: | :---: | :---: | :--- |
| Naive compress-then-encrypt (unprotected) | `1952` | `134` | `146` | `12` | **YES (LEAK)** | `none` |
| `polarstore_adaptive` (`ModeAttackerControlledMixedSecret`) | `1952` | `1980` | `1980` | `0` | **NO (ELIMINATED)** | `crime_breach_mixed_secret_oracle` |

## 4. Defense-in-Depth High-Entropy Probe (Unlabeled Ciphertext Fallback)

| Objects (`128 KiB` raw AES-GCM mislabeled as plaintext) | Bypassed (`high_entropy_incompressible`) | Mean Shannon Entropy (bits/B) | Mean Probe Latency (us) | Trial Compressions Run | Effective Bypass Thru (MB/s) | Ratio |
| :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `32` | `32 / 32 (100%)` | `7.977` | `15.04` | `0` | `8716.1` | `1.000x` |
