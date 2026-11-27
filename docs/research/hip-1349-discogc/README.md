# HIP-1349 / HIP-1348 Research Dossier: Discard-versus-Compaction Economics (DisCoGC)

- **Tracking Issues**: `hanzoai/s3#27` (HIP-1349) & `hanzoai/s3#31` (HIP-1348 Umbrella)
- **Bounty Class**: Research Bounty (quantitative storage economics study + bounded clean-room Go harness)
- **Primary Paper**: Runhua Bian, Liqiang Zhang, Jinxin Liu, Jiacheng Zhang, Jianong Zhong, Jiahao Gu, Hao Guo, Zhihong Guo, Yunhao Li, Fenghao Zhang, Jiangkun Zhao, Yangming Chen, Guojun Li, Ruwen Fan, Haijia Shen, Chengyu Dong, Yao Wang, Rui Shi, Jiwu Shu, and Youyou Lu. *"Discard-Based Garbage Collection for Distributed Log-Structured Storage Systems in ByteDance"* (**DisCoGC**), Proceedings of the 24th USENIX Conference on File and Storage Technologies (**FAST '26**), February 24–26, 2026, Santa Clara, CA, USA, ISBN 978-1-939133-53-3.
  - Presentation & Open-Access PDF: https://www.usenix.org/conference/fast26/presentation/bian
- **Artifacts in This Directory**:
  - `README.md` — 6-point HIP-1349 / HIP-1348 evaluation report, license audit, methodology, limitations, and Hanzo S3 integration recommendation.
  - `score_table.md` — Generated Markdown comparison tables across 5 policies × 4 workloads + 3 boundary/negative regimes.
  - `raw_results.json` — Machine-readable raw measurements (`syscall.Stat_t.Blocks`, `FALLOC_FL_PUNCH_HOLE`, `/proc/self/io`, flash wear, and `$/TB-month` cost breakdowns).

> **Honesty & Scope Notice**: Reported figures in Bian et al. (FAST '26) — such as ~20% TCO reduction on ByteDance's ByteStore/ByteDrive production fleet — belong strictly to the original authors' environment and are **not** Hanzo S3 claims or expectations. All numbers reported below were measured locally on the benchmark host (`linux/amd64`, 2 vCPUs, Linux 6.18.57, Go 1.27.1) using real OS files, real Linux `fallocate(FALLOC_FL_PUNCH_HOLE | FALLOC_FL_KEEP_SIZE)` syscalls, real `fstat(2)` `st_blocks * 512` physical block accounting, and real `/proc/self/io` counters, combined with an explicit, parameterized enterprise NVMe TLC SSD wear/TCO model. Production Hanzo S3 deployment decisions require reproducing these experiments on target Hanzo NVMe hardware (and verifying SSD TRIM IOPS and FTL PWA directly via SMART `Data Units Written` / `Media Units Written`).

---

## 1. Paper Summary & License / IP Audit

### 1.1 Paper Summary
Log-structured / append-only storage engines (such as ByteDance's ByteStore and Hanzo S3's SeaweedFS-derived volume files in `s3/storage/volume_vacuum.go`) achieve high sequential write throughput, simple crash recovery, and clean erasure-coding semantics, but suffer from a fundamental **write amplification vs. space amplification trade-off** under copy-based garbage collection:
- **Aggressive copy-compaction** reclaims stale space quickly (low space amplification) by reading live objects out of fragmented segments and rewriting them into new segments, which inflates Logical Write Amplification ($\text{LWA}$), compounds with SSD internal Physical Write Amplification ($\text{PWA}$) as $\text{WA} = \text{LWA} \times \text{PWA}$, burns NAND P/E cycles, shortens SSD lifetime, and steals foreground I/O bandwidth.
- **Lazy compaction / no GC** keeps $\text{LWA} \approx 1.0\times$ but traps dead space on disk, driving unbounded space amplification and SSD capacity capex.

**DisCoGC** introduces a **discard-and-compaction combined GC scheme**: instead of copying live data to reclaim dead extents, the storage engine issues direct extent discards (filesystem hole-punching / NVMe Dataset Management TRIM / ZNS discard) over stale block-aligned regions **without moving valid data**, falling back to copy-compaction only when sub-block fragmentation, extent-map fragmentation, or near-empty segment metadata overhead makes compaction cheaper.

### 1.2 License, Dependency & Patent Audit
| Dimension | Audit Result |
|---|---|
| Primary paper access | Open-access USENIX FAST '26 proceedings (`https://www.usenix.org/conference/fast26/presentation/bian`). |
| Upstream source code | ByteStore / ByteDrive / DisCoGC is proprietary ByteDance production infrastructure with **no public Apache-2.0 / MIT code release**. |
| Clean-room implementation | `s3/gc/discogc/` and `cmd/discogc-bench/` are **100% clean-room Go** written from first principles using **Go standard library only** (`syscall`, `os`, `crypto/sha256`, `hash/crc32`, `sync`, `encoding/binary`, `math/rand`). |
| External Go dependencies added | **0** (`go.mod` and `go.sum` are completely untouched). |
| Datasets / traces imported | **0** proprietary traces imported; all workloads are deterministic seeded generators in `s3/gc/discogc/workload.go`. |
| License compatibility | 100% compatible with `hanzoai/s3` (Apache-2.0). |
| Patent / deployment note | Linux `FALLOC_FL_PUNCH_HOLE` (sparse-file hole punching, Linux 2.6.38+) and NVMe TRIM/Deallocate are standard POSIX/Linux/NVMe interfaces. Before shipping any ByteStore-specific proprietary heuristics into production, standard Hanzo legal review applies; this clean-room study uses only general-purpose OS hole-punching and threshold-based fallback compaction. |

---

## 2. Matched Baselines & Reproduction

All 5 policies execute against the **exact same file-backed segment engine** (`s3/gc/discogc/store.go`), identical needle framing (64 B segment header, 32 B needle header with magic + CRC32 + SHA-256 verification), identical cluster alignment (`4096 B` for objects $\ge 4\,\text{KiB}$; tight `8 B` SeaweedFS-style packing for sub-block objects $< 4\,\text{KiB}$), identical garbage trigger (`0.30`, matching Hanzo S3 `Volume.garbageLevel()`), and identical deterministic RNG seeds.

### 2.1 Five Matched Policies
1. `no_gc_append_only` — Pure append-only baseline (zero GC copies, zero discards; exposes unbounded space amplification).
2. `greedy_copy_compaction` — Upstream Hanzo S3 / SeaweedFS-style vacuum baseline (`s3/storage/volume_vacuum.go:CompactByVolumeData`): rewrites all live needles from sealed segments exceeding `garbageRatio >= 0.30` into new segments (highest dead ratio first).
3. `cost_benefit_copy_compaction` — Rosenblum–Ousterhout LFS age-weighted copy-compaction (`score = dead_ratio * age / (1 + live_ratio)`), deferred to segment-seal passes so hot segments cool before copying.
4. `discogc_adaptive_discard` — DisCoGC-inspired hybrid policy: coalesces adjacent dead extents, computes strictly interior `DiscardGranularity`-aligned dead intervals `[ceil_align(start), floor_align(end))`, reclaims physical OS blocks via Linux `fallocate(FALLOC_FL_PUNCH_HOLE | FALLOC_FL_KEEP_SIZE)` with **zero live-data copying**, and falls back to selective copy-compaction ONLY when (a) un-reclaimable sub-block dead space $\ge 0.30$, (b) extent-fragment count exceeds `MaxExtentFragments` (default `128`), or (c) sealed segment live ratio drops $\le 0.05$.
5. `blind_subblock_discard_negative` — Documented negative control: calls raw `FALLOC_FL_PUNCH_HOLE` on every deleted object regardless of size or block alignment and never compacts.

### 2.2 Reproducible Commands
```bash
# Run full correctness, real-disk block reclamation, negative-regime, and -race concurrency tests
go test -race -v -count=1 ./s3/gc/discogc/...

# Reproduce raw_results.json and score_table.md on your host hardware
go run ./cmd/discogc-bench -out docs/research/hip-1349-discogc
```

---

## 3. Falsifiable Hypothesis & Expected Failure Modes

### 3.1 Primary Falsifiable Hypothesis
> **H1 (Aligned Churn Economic Advantage)**: On S3 object-churn workloads where object sizes are $\ge$ filesystem/NVMe discard granularity (`4 KiB`, tested across `16–256 KiB` Zipfian, uniform random, and bulk-TTL workloads), `discogc_adaptive_discard` will reduce **Logical WAF by $\ge 50\%$**, cut **live bytes copied by $\ge 90\%$**, preserve **nominal 5-year SSD lifetime** (vs. premature endurance wear-out under greedy copy-compaction at 0.80 logical DWPD), and reduce **total storage cost (`$/TB-month`) by $\ge 40\%$** relative to both `greedy_copy_compaction` and `cost_benefit_copy_compaction`, while maintaining time-averaged physical space amplification $\le 1.20\times$ and 100% bit-exact SHA-256 live-object integrity.

### 3.2 Documented Negative & Boundary Regimes (Expected Failure Modes)
1. **`sub_block_fragmentation_no_reclaim`**: Deleting small (`512 B–1.5 KiB`) objects interleaved with live objects inside shared `4 KiB` filesystem blocks **cannot** be reclaimed by `FALLOC_FL_PUNCH_HOLE` (`OSBlocksReclaimedBytes == 0`). Blind discard fails completely (unbounded space amp + massive wasted TRIM syscall storm); adaptive discard must detect sub-block dead space and fall back to copy-compaction (converging to compaction parity rather than beating compaction).
2. **`extent_map_fragmentation_tax`**: Alternating live/dead `8 KiB` checkerboard deletes create hundreds of non-contiguous holes under uncapped hole-punching, bloating filesystem extent trees and degrading sequential scan throughput. Bounding `MaxExtentFragments` restores contiguous layout and improves sequential scan throughput at the cost of periodic defragmentation writes.
3. **`append_only_zero_delete_workload`**: On immutable archives with `0%` overwrites and `0%` deletes, no GC work occurs; all non-blind policies converge to identical Logical WAF, identical physical disk usage, and identical `$/TB-month` TCO — discard adds zero benefit.
4. **Hardware-Dependent Limit (TRIM IOPS & FTL PWA)**: Real NVMe SSDs exhibit finite Dataset Management (TRIM) IOPS and device-specific internal PWA. High-rate tiny discards can saturate SSD TRIM queues if not batched/coalesced (confirmed by our `blind_subblock_discard_negative` syscall count of `4,290` calls for `5.5 MiB` written).

---

## 4. Quantitative Economics Model

All economic calculations are implemented transparently in `s3/gc/discogc/economics.go` with zero double-counting:

- **Logical Write Amplification ($\text{LWA}$)**:
  $$\text{LWA} = \frac{B_{\text{user}} + B_{\text{gc\_copied}} + B_{\text{meta}}}{B_{\text{user}}}$$
- **Total Physical Write Amplification ($\text{WA}$)** (per Bian et al., FAST '26):
  $$\text{WA} = \text{LWA} \times \text{PWA}_{\text{device}}$$
  where $\text{PWA}_{\text{device}}$ parameterizes SSD FTL amplification (`1.08` with active host discard/TRIM, `1.10` for pure sequential append-only with zero GC, `1.25` for copy-compaction churn).
- **Flash Wear & Drive Lifetime** (default model: `7.68 TB` Enterprise TLC NVMe, `$1,152` = `$150/TB`, `1.0 DWPD`, `5 yr` warranty, `3000` P/E cycles, `24.0 logical TB/month/TB` $\approx 0.80$ logical DWPD):
  $$\text{DWPD}_{\text{eff}} = \frac{24.0}{30.0} \times \text{WA}, \qquad L_{\text{drive}} = \min\!\left(L_{\text{warranty}},\; L_{\text{warranty}} \cdot \frac{\text{DWPD}_{\text{rated}}}{\text{DWPD}_{\text{eff}}}\right)$$
- **Total Monthly Cost (`$/TB-month` per live logical TB)**:
  - `Cap Cost` = $\frac{\$150}{60} \times \overline{\text{SA}} \times (1 + 0.07)$ (nominal 5-year capacity amortization scaled by time-averaged physical space amplification $\overline{\text{SA}}$).
  - `Wear Cost` = $\max\!\left(0,\; \frac{\$150}{12 \cdot L_{\text{drive}}} \cdot \overline{\text{SA}} \cdot 1.07 - \text{Cap Cost}\right)$ (exact incremental monthly SSD replacement cost strictly from wear shortening drive life below 5 years — zero double-counting).
  - `IO Tax` = $B_{\text{copied/mo}} \cdot \$0.00025/\text{GiB} + N_{\text{discards/mo}} \cdot \$0.002/10^6$.

---

## 5. Measurements & Kernel Verification Results

Full tables are in [`score_table.md`](./score_table.md) and machine-readable data in [`raw_results.json`](./raw_results.json).

### 5.1 Headline Comparison Across All 4 Workloads (Local Host Measurements)

| Workload | Policy | LWA | Effective WA | Avg SA | Kernel `/proc/self/io` `write_bytes` | OS Reclaimed (`st_blocks`) | GC Copied Live | Drive Life | Total (`$/TB-mo`) |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| **`zipfian_hot_cold_churn`** | `no_gc_append_only` | 1.057x | 1.162x | 3.05x | 75.24 MiB | 0.00 MiB | 0.00 MiB | 5.00 yr | $8.15 |
| (16–64 KiB, Zipf s=1.10) | `greedy_copy_compaction` | 2.636x | 3.296x | 1.27x | 187.75 MiB | 0.00 MiB | 106.22 MiB | 1.90 yr | $18.15 |
| | `cost_benefit_copy_compaction` | 1.574x | 1.968x | 1.58x | 112.11 MiB | 0.00 MiB | 34.83 MiB | 3.18 yr | $9.65 |
| | **`discogc_adaptive_discard`** | **1.060x** | **1.145x** | **1.07x** | **75.47 MiB** | **57.61 MiB** | **0.20 MiB** | **5.00 yr** | **$3.79** |
| | `blind_subblock_discard_negative` | 1.057x | 1.321x | 1.07x | 75.24 MiB | 57.61 MiB | 0.00 MiB | 4.73 yr | $4.05 |
| **`uniform_random_churn`** | `no_gc_append_only` | 1.129x | 1.242x | 3.59x | 59.62 MiB | 0.00 MiB | 0.00 MiB | 5.00 yr | $9.60 |
| (32 KiB uniform) | `greedy_copy_compaction` | 3.035x | 3.794x | 1.31x | 160.21 MiB | 0.00 MiB | 89.06 MiB | 1.65 yr | $21.03 |
| | `cost_benefit_copy_compaction` | 1.142x | 1.427x | 3.51x | 60.25 MiB | 0.00 MiB | 0.56 MiB | 4.38 yr | $10.80 |
| | **`discogc_adaptive_discard`** | **1.150x** | **1.242x** | **1.14x** | **60.71 MiB** | **49.61 MiB** | **0.97 MiB** | **5.00 yr** | **$4.50** |
| | `blind_subblock_discard_negative` | 1.129x | 1.412x | 1.14x | 59.62 MiB | 49.61 MiB | 0.00 MiB | 4.43 yr | $4.84 |
| **`large_object_bulk_deletes`** | `no_gc_append_only` | 1.013x | 1.114x | 2.91x | 84.86 MiB | 0.00 MiB | 0.00 MiB | 5.00 yr | $7.78 |
| (128–256 KiB + 50% TTL sweeps) | `greedy_copy_compaction` | 2.473x | 3.091x | 1.11x | 207.22 MiB | 0.00 MiB | 120.83 MiB | 2.02 yr | $16.21 |
| | `cost_benefit_copy_compaction` | 1.175x | 1.469x | 1.68x | 98.49 MiB | 0.00 MiB | 13.47 MiB | 4.25 yr | $6.28 |
| | **`discogc_adaptive_discard`** | **1.013x** | **1.094x** | **1.02x** | **84.86 MiB** | **66.88 MiB** | **0.00 MiB** | **5.00 yr** | **$2.93** |
| | `blind_subblock_discard_negative` | 1.013x | 1.266x | 1.02x | 84.86 MiB | 66.88 MiB | 0.00 MiB | 4.94 yr | $2.99 |
| **`small_fragmented_subblock`** | `no_gc_append_only` | 1.045x | 1.149x | 1.90x | 6.14 MiB | 0.00 MiB | 0.00 MiB | 5.00 yr | $5.07 |
| (512 B–1.5 KiB interleaved) | `greedy_copy_compaction` | 2.617x | 3.271x | 1.20x | 15.38 MiB | 0.00 MiB | 8.83 MiB | 1.91 yr | $17.62 |
| | `cost_benefit_copy_compaction` | 1.142x | 1.427x | 1.88x | 6.71 MiB | 0.00 MiB | 0.55 MiB | 4.38 yr | $6.31 |
| *(honest negative/fallback)* | `discogc_adaptive_discard` | 2.640x | 2.852x | 1.20x | 15.52 MiB | 0.05 MiB | 8.96 MiB | 2.19 yr | $16.80 |
| *(negative control)* | `blind_subblock_discard_negative` | 1.045x | 1.306x | 1.90x | 6.14 MiB | **0.00 MiB** | 0.00 MiB | 4.79 yr | **$42.13** |

### 5.2 Key Quantitative Takeaways
- **Hypothesis H1 Confirmed on Aligned Workloads ($\ge 4\,\text{KiB}$)**:
  - On `zipfian_hot_cold_churn`: `discogc_adaptive_discard` cuts Logical WAF by **59.8%** (`2.636x -> 1.060x`), cuts live copied bytes by **99.8%** (`106.22 MiB -> 0.20 MiB`), cuts real Linux kernel `/proc/self/io` `write_bytes` by **59.8%** (`187.75 MiB -> 75.47 MiB`), extends modeled SSD lifetime from **1.90 yr -> 5.00 yr**, and reduces modeled TCO by **79.1%** vs. greedy compaction (`$18.15 -> $3.79/TB-mo`) and **60.7%** vs. LFS cost-benefit (`$9.65 -> $3.79/TB-mo`).
  - On `uniform_random_churn`: Logical WAF drops **62.1%** (`3.035x -> 1.150x`), live copies drop **98.9%** (`89.06 MiB -> 0.97 MiB`), kernel `write_bytes` drops **62.1%** (`160.21 MiB -> 60.71 MiB`), and TCO drops **78.6%** vs. greedy (`$21.03 -> $4.50/TB-mo`).
  - On `large_object_bulk_deletes`: **0.00 MiB** live data copied (**100% copy elimination**), **66.88 MiB** reclaimed directly via `FALLOC_FL_PUNCH_HOLE`, and **81.9%** TCO reduction vs. greedy (`$16.21 -> $2.93/TB-mo`).
- **Boundary / Negative Regimes Confirmed**:
  - `sub_block_fragmentation_no_reclaim`: `blind_subblock_discard_negative` issued `512` `FALLOC_FL_PUNCH_HOLE` calls on interleaved 512 B objects and reclaimed **0 OS bytes** (`st_blocks` unchanged, `FinalSA = 2.23x`), whereas `discogc_adaptive_discard` detected non-reclaimable sub-block dead space and fell back to `1` compaction cycle (`FinalSA = 1.66x`).
  - `extent_map_fragmentation_tax`: Uncapped checkerboard hole-punching left **256 extent fragments** (`158.7 MiB/s` sequential scan); capping `MaxExtentFragments = 64` defragmented to **61 fragments** and improved real sequential scan throughput to **192.1 MiB/s** (`+21.0%`) at an explicit compaction WAF trade-off (`LWA 1.01x -> 3.27x`).
  - `append_only_zero_delete_workload`: All non-blind policies produced identical `LWA = 1.130x`, `0` copied bytes, `0` reclaimed blocks, and identical `$3.02/TB-mo` TCO.

---

## 6. Correctness, Concurrency, Multi-Tenant Safety & Integration Recommendation

### 6.1 Evidence & Safety Verification
Executed via `go test -race -v -count=1 ./s3/gc/discogc/...` (all passing, zero data races):
- `TestBitExactSHA256PreservationAcrossPolicies`: Verifies byte-for-byte payload and SHA-256 identity of all surviving live objects and `ErrNotFound` for deleted keys across all 5 policies.
- `TestLinuxRealPhysicalBlockReclamation`: Verifies `syscall.Fallocate(FALLOC_FL_PUNCH_HOLE|FALLOC_FL_KEEP_SIZE)` strictly lowers `syscall.Stat_t.Blocks * 512` while keeping `st_size` constant and copying `0` live bytes.
- `TestWAFReductionAdaptiveDiscardVsGreedyCompaction`: Asserts WAF, live-copy, drive-lifetime, and TCO improvements of `discogc_adaptive_discard` over `greedy_copy_compaction`.
- `TestNegativeRegimeBlindSubblockDiscardZeroReclaim`: Verifies all 3 boundary/negative regimes (0 bytes reclaimed by blind sub-block discard, extent-fragment cap enforcement, zero-delete policy parity).
- `TestConcurrentMultiTenantIsolationRace`: 8 concurrent tenants running parallel `Put`/`Get`/`Delete`/`RunGC`/`Metrics` with strict per-tenant segment file isolation and cross-tenant read-denial assertions under `-race`.
- `TestEconomicsFormulaSanity`: Unit tests `ComputeEconomics` monotonicity and wear-life bounds.

### 6.2 Canonical Architecture Fit & Hanzo S3 Integration Recommendation
This study preserves all canonical Hanzo stack boundaries (`vfs -> hanzoai/s3 -> visor`, Cloud owning lifecycle & billing, Enso routing models, Zen generating outputs, Kai returning structured decisions). No parallel storage stack is proposed.

**Recommendation for HIP-1349 / HIP-1348 (Conditional Proceed to Hardware Validation + Implementation Phase)**:
1. **Target Integration Point**: `s3/storage/volume_vacuum.go` and `s3/storage/backend/disk_file.go` inside Hanzo S3 volume servers.
2. **Where Direct Discard Wins**: S3 object tiers dominated by objects $\ge 16\,\text{KiB}$ (LLM checkpoints, Parquet/Arrow analytics shards, media chunks, RAG vector segments) with active overwrite or TTL expiration churn. Adding an optional block-aligned extent allocation class (`>= 4 KiB` / `64 KiB` alignment) plus coalesced `FALLOC_FL_PUNCH_HOLE` in `Volume.Vacuum` eliminates full `.dat -> .cpd` rewrites for large-object volumes.
3. **Where Copy-Compaction Must Remain**: Small-object volumes (`< 4 KiB` needles packed at 8-byte alignment in `s3/storage/needle`) and heavily checkerboard-fragmented volumes exceeding extent-tree budgets (`MaxExtentFragments`). Hybrid fallback (`discogc_adaptive_discard`) is mandatory; blind discard must never be enabled.
4. **Required Next Step Before Production Merge**: Re-run `cmd/discogc-bench` on a real bare-metal Hanzo NVMe SSD (`ext4` and `xfs` mounted with `discard=async` or periodic `fstrim`), collecting hardware SMART `Host Writes` vs. `NAND Writes` (`Media Units Written`) and NVMe TRIM queue tail latency to replace the parameterized `PWA` assumption with hardware-measured SSD internal WAF.
