# DisCoGC Discard-versus-Compaction Economics — Measured Score Table

- **Generated (UTC)**: `2026-11-25T13:13:42Z`
- **Host**: `linux/amd64`, `2` CPUs, `go1.27.1`
- **Drive Model**: `7.68TB Enterprise TLC NVMe (1.0 DWPD, 5yr)` ($1152 / 7.68 TB = $150.00/TB, 1.0 DWPD, 5-yr warranty, 3000 P/E cycles, 24.0 logical TB/mo/TB)

> **Measurement Scope**: Logical WAF, copied live bytes, OS physical blocks (`st_blocks * 512`), `FALLOC_FL_PUNCH_HOLE` reclaimed bytes, and `/proc/self/io` deltas are directly measured on the benchmark host. Flash wear (TBW, P/E cycles, drive life) and `$/TB-month` TCO combine those measurements with the explicit drive model above. Figures are reproducible local host measurements, not Hanzo production claims until replicated on target NVMe hardware.

## Workload: `zipfian_hot_cold_churn`

Zipfian (s=1.10) hot-cold S3 object churn (16-64 KiB objects, 75% overwrite / 15% delete / 10% fresh put) (seed=`20260224`, keys=`320`, churn_ops=`1800`, obj=`16384–65536 B`).

| Policy | Logical WAF (LWA) | Effective Physical WAF (WA) | Avg Space Amp | Final Space Amp | GC Copied Live (MiB) | OS Blocks Reclaimed (MiB) | Punch Calls | Drive Life (yr) | Cap Cost ($/TB-mo) | Wear Cost ($/TB-mo) | IO Tax ($/TB-mo) | Total ($/TB-mo) |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| `no_gc_append_only` | 1.057x | 1.162x | 3.05x | 4.58x | 0.00 | 0.00 | 0 | 5.00 | $8.15 | $0.00 | $0.00 | **$8.15** |
| `greedy_copy_compaction` | 2.636x | 3.296x | 1.27x | 1.33x | 106.22 | 0.00 | 0 | 1.90 | $3.41 | $5.57 | $9.16 | **$18.15** |
| `cost_benefit_copy_compaction` | 1.574x | 1.968x | 1.58x | 1.40x | 34.83 | 0.00 | 0 | 3.18 | $4.22 | $2.42 | $3.00 | **$9.65** |
| `discogc_adaptive_discard` | 1.060x | 1.145x | 1.07x | 1.07x | 0.20 | 57.61 | 1285 | 5.00 | $2.86 | $0.00 | $0.93 | **$3.79** |
| `blind_subblock_discard_negative` | 1.057x | 1.321x | 1.07x | 1.07x | 0.00 | 57.61 | 1466 | 4.73 | $2.85 | $0.16 | $1.04 | **$4.05** |

## Workload: `uniform_random_churn`

Uniform random overwrites & deletes over 32 KiB objects (worst-case mixed-age segments for copy-compaction) (seed=`20260225`, keys=`320`, churn_ops=`1600`, obj=`32768–32768 B`).

| Policy | Logical WAF (LWA) | Effective Physical WAF (WA) | Avg Space Amp | Final Space Amp | GC Copied Live (MiB) | OS Blocks Reclaimed (MiB) | Punch Calls | Drive Life (yr) | Cap Cost ($/TB-mo) | Wear Cost ($/TB-mo) | IO Tax ($/TB-mo) | Total ($/TB-mo) |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| `no_gc_append_only` | 1.129x | 1.242x | 3.59x | 6.86x | 0.00 | 0.00 | 0 | 5.00 | $9.60 | $0.00 | $0.00 | **$9.60** |
| `greedy_copy_compaction` | 3.035x | 3.794x | 1.31x | 1.38x | 89.06 | 0.00 | 0 | 1.65 | $3.51 | $7.15 | $10.37 | **$21.03** |
| `cost_benefit_copy_compaction` | 1.142x | 1.427x | 3.51x | 3.75x | 0.56 | 0.00 | 0 | 4.38 | $9.40 | $1.33 | $0.07 | **$10.80** |
| `discogc_adaptive_discard` | 1.150x | 1.242x | 1.14x | 1.14x | 0.97 | 49.61 | 1405 | 5.00 | $3.05 | $0.00 | $1.45 | **$4.50** |
| `blind_subblock_discard_negative` | 1.129x | 1.412x | 1.14x | 1.15x | 0.00 | 49.61 | 1468 | 4.43 | $3.05 | $0.39 | $1.40 | **$4.84** |

## Workload: `large_object_bulk_deletes`

Large 128-256 KiB S3 checkpoint/parquet objects with periodic 50% bulk TTL expiration sweeps (seed=`20260226`, keys=`120`, churn_ops=`360`, obj=`131072–262144 B`).

| Policy | Logical WAF (LWA) | Effective Physical WAF (WA) | Avg Space Amp | Final Space Amp | GC Copied Live (MiB) | OS Blocks Reclaimed (MiB) | Punch Calls | Drive Life (yr) | Cap Cost ($/TB-mo) | Wear Cost ($/TB-mo) | IO Tax ($/TB-mo) | Total ($/TB-mo) |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| `no_gc_append_only` | 1.013x | 1.114x | 2.91x | 4.82x | 0.00 | 0.00 | 0 | 5.00 | $7.78 | $0.00 | $0.00 | **$7.78** |
| `greedy_copy_compaction` | 2.473x | 3.091x | 1.11x | 1.13x | 120.83 | 0.00 | 0 | 2.02 | $2.97 | $4.38 | $8.86 | **$16.21** |
| `cost_benefit_copy_compaction` | 1.175x | 1.469x | 1.68x | 1.59x | 13.47 | 0.00 | 0 | 4.25 | $4.50 | $0.79 | $0.99 | **$6.28** |
| `discogc_adaptive_discard` | 1.013x | 1.094x | 1.02x | 1.02x | 0.00 | 66.88 | 349 | 5.00 | $2.72 | $0.00 | $0.21 | **$2.93** |
| `blind_subblock_discard_negative` | 1.013x | 1.266x | 1.02x | 1.02x | 0.00 | 66.88 | 392 | 4.94 | $2.72 | $0.03 | $0.24 | **$2.99** |

## Workload: `small_fragmented_subblock`

Small 512 B-1.5 KiB objects packed inside shared 4 KiB filesystem blocks with interleaved deletes (seed=`20260227`, keys=`2400`, churn_ops=`4800`, obj=`512–1536 B`).

| Policy | Logical WAF (LWA) | Effective Physical WAF (WA) | Avg Space Amp | Final Space Amp | GC Copied Live (MiB) | OS Blocks Reclaimed (MiB) | Punch Calls | Drive Life (yr) | Cap Cost ($/TB-mo) | Wear Cost ($/TB-mo) | IO Tax ($/TB-mo) | Total ($/TB-mo) |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| `no_gc_append_only` | 1.045x | 1.149x | 1.90x | 3.64x | 0.00 | 0.00 | 0 | 5.00 | $5.07 | $0.00 | $0.00 | **$5.07** |
| `greedy_copy_compaction` | 2.617x | 3.271x | 1.20x | 1.26x | 8.83 | 0.00 | 0 | 1.91 | $3.20 | $5.17 | $9.25 | **$17.62** |
| `cost_benefit_copy_compaction` | 1.142x | 1.427x | 1.88x | 2.35x | 0.55 | 0.00 | 0 | 4.38 | $5.03 | $0.71 | $0.57 | **$6.31** |
| `discogc_adaptive_discard` | 2.640x | 2.852x | 1.20x | 1.19x | 8.96 | 0.05 | 12 | 2.19 | $3.20 | $4.10 | $9.49 | **$16.80** |
| `blind_subblock_discard_negative` | 1.045x | 1.306x | 1.90x | 3.64x | 0.00 | 0.00 | 4290 | 4.79 | $5.07 | $0.23 | $36.83 | **$42.13** |

## Documented Negative & Boundary Regimes

### `sub_block_fragmentation_no_reclaim`

- **Hypothesis**: Deleting interleaved 512 B objects inside shared 4 KiB filesystem blocks yields 0 reclaimed OS blocks under direct hole-punching; only copy-compaction (or adaptive fallback) reclaims physical disk space.
- **Observed**: blind_subblock_discard_negative issued 512 PunchHole calls and reclaimed 0 OS bytes (final space amp 2.23x); discogc_adaptive_discard detected sub-block dead space and fell back to 1 compaction cycle(s) (final space amp 1.66x).

| Policy / Variant | Logical WAF | Final Space Amp | Physical Disk (KiB) | OS Blocks Reclaimed (KiB) | Punch Calls | Compaction Cycles | Extent Fragments | Seq Scan (MiB/s) | Total ($/TB-mo) |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| `blind_subblock_discard_negative` | 1.107x | 2.23x | 500.0 | 0.0 | 512 | 0 | 0 | 0.0 | $65.51 |
| `discogc_adaptive_discard` | 1.937x | 1.66x | 372.0 | 0.0 | 0 | 1 | 0 | 0.0 | $13.21 |
| `greedy_copy_compaction` | 1.882x | 1.55x | 348.0 | 0.0 | 0 | 1 | 0 | 0.0 | $12.12 |

### `extent_map_fragmentation_tax`

- **Hypothesis**: Checkerboard alternating 8 KiB deletes produce hundreds of non-contiguous hole extents when hole-punching is uncapped; imposing MaxExtentFragments bounds extent-tree fragmentation via periodic defragmentation at a small one-time copy cost.
- **Observed**: uncapped punch-hole left 256 extent fragments (LWA 1.01x, scan 158.7 MiB/s); bounded DisCoGC (cap=64) defragmented to 61 fragments (LWA 3.27x, scan 192.1 MiB/s).

| Policy / Variant | Logical WAF | Final Space Amp | Physical Disk (KiB) | OS Blocks Reclaimed (KiB) | Punch Calls | Compaction Cycles | Extent Fragments | Seq Scan (MiB/s) | Total ($/TB-mo) |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| `uncapped_punch_hole_no_defrag` | 1.009x | 1.01x | 2052.0 | 2048.0 | 256 | 0 | 256 | 158.7 | $5.95 |
| `discogc_bounded_defrag_cap64` | 3.268x | 1.01x | 2052.0 | 2048.0 | 256 | 3 | 61 | 192.1 | $24.62 |

### `append_only_zero_delete_workload`

- **Hypothesis**: When delete and overwrite rates are 0%, no GC work occurs and all 5 policies yield identical Logical WAF (~1.00x), zero reclaimed blocks, and identical space amplification.
- **Observed**: All 5 policies produced identical Logical WAF (1.1297x), 0 copied live bytes, 0 reclaimed blocks, and final space amplification 1.13x.

| Policy / Variant | Logical WAF | Final Space Amp | Physical Disk (KiB) | OS Blocks Reclaimed (KiB) | Punch Calls | Compaction Cycles | Extent Fragments | Seq Scan (MiB/s) | Total ($/TB-mo) |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| `no_gc_append_only` | 1.130x | 1.13x | 5784.0 | 0.0 | 0 | 0 | 0 | 219.6 | $3.02 |
| `greedy_copy_compaction` | 1.130x | 1.13x | 5784.0 | 0.0 | 0 | 0 | 0 | 206.7 | $3.02 |
| `cost_benefit_copy_compaction` | 1.130x | 1.13x | 5784.0 | 0.0 | 0 | 0 | 0 | 216.2 | $3.02 |
| `discogc_adaptive_discard` | 1.130x | 1.13x | 5784.0 | 0.0 | 0 | 0 | 0 | 210.2 | $3.02 |
| `blind_subblock_discard_negative` | 1.130x | 1.13x | 5784.0 | 0.0 | 6 | 0 | 0 | 211.3 | $3.47 |

