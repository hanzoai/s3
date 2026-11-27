# HIP-1349 / HIP-1348 FailureMiner Benchmark Score Table

- **Dataset**: `nvme_fleet_90d_multivendor` (seed=`1349`, total_drives=`1200`, train=`601`, test=`599`, horizon=`90d`, step=`6h`)
- **Data Source**: `deterministic_seeded_nvme_smart_simulation` (`real_hardware_exercised=false`)
- **Runtime Host**: `go1.27.1` (`linux/amd64`, `2` CPUs)
- **Honesty Notice**: All figures are reproducible measurements from the evaluation corpus on the authoring host; they are not Hanzo production SLA claims until validated on production NVMe SMART logs.

## 1. Matched 5-Method Fleet Comparison (Test Split)

| Method | TP/FP/FN | Recall | Pred. Recall | Precision (Evac Acc) | F1 | Drive FPR | FAR / Drive-Day | Median Lead (h) | Mean Lead (h) | P10 / P90 Lead (h) | >=24h Recall | Hard-Neg FP | Cost Savings | Wasted Evac (TB) | Throughput (samp/s) | Latency (ns/samp) |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| `static_smart_threshold` | 14/24/9 | 60.9% | 70.0% | 36.8% | 0.459 | 4.17% | 0.000464 | 72.0 | 57.9 | 6.0 / 114.6 | 34.8% | 24/115 | +47.1% | 184.3 | 20376097 | 49 |
| `sliding_window_zscore` | 10/14/13 | 43.5% | 50.0% | 41.7% | 0.426 | 2.43% | 0.000271 | 114.0 | 99.0 | 70.8 / 127.8 | 39.1% | 10/115 | +34.8% | 107.5 | 2524914 | 396 |
| `ewma_hazard_score` | 20/3/3 | 87.0% | 100.0% | 87.0% | 0.870 | 0.52% | 0.000058 | 126.0 | 141.3 | 83.4 / 219.6 | 87.0% | 3/115 | +78.6% | 23.0 | 6745000 | 148 |
| `failureminer_temporal_pattern` | 20/1/3 | 87.0% | 100.0% | 95.2% | 0.909 | 0.17% | 0.000019 | 168.0 | 179.1 | 107.4 / 270.6 | 87.0% | 1/115 | +79.3% | 7.7 | 3121803 | 320 |
| `over_sensitive_negative` | 23/575/0 | 100.0% | 100.0% | 3.9% | 0.074 | 99.83% | 0.531202 | 1044.0 | 1093.3 | 562.8 / 1788.0 | 100.0% | 115/115 | -116.7% | 4416.0 | 13059072 | 77 |

## 2. Per-Failure-Mode Detection & Median Lead Time

| Method | Failure Mode | Detected / Total | Rate | Median Lead (h) |
|---|---|---:|---:|---:|
| `static_smart_threshold` | `progressive_media_error_ramp` | 8 / 8 | 100.0% | 99.0 |
| `static_smart_threshold` | `thermal_spare_depletion` | 5 / 7 | 71.4% | 6.0 |
| `static_smart_threshold` | `intermittent_crc_timeout_burst` | 1 / 5 | 20.0% | 6.0 |
| `static_smart_threshold` | `sudden_unannounced_controller_death` | 0 / 3 | 0.0% | 0.0 |
| `static_smart_threshold` | `hard_negative_hot_tier` | 24 / 115 | 20.9% | 0.0 |
| `static_smart_threshold` | `healthy` | 0 / 461 | 0.0% | 0.0 |
| `sliding_window_zscore` | `progressive_media_error_ramp` | 8 / 8 | 100.0% | 102.0 |
| `sliding_window_zscore` | `thermal_spare_depletion` | 2 / 7 | 28.6% | 129.0 |
| `sliding_window_zscore` | `intermittent_crc_timeout_burst` | 0 / 5 | 0.0% | 0.0 |
| `sliding_window_zscore` | `sudden_unannounced_controller_death` | 0 / 3 | 0.0% | 0.0 |
| `sliding_window_zscore` | `hard_negative_hot_tier` | 10 / 115 | 8.7% | 0.0 |
| `sliding_window_zscore` | `healthy` | 4 / 461 | 0.9% | 0.0 |
| `ewma_hazard_score` | `progressive_media_error_ramp` | 8 / 8 | 100.0% | 114.0 |
| `ewma_hazard_score` | `thermal_spare_depletion` | 7 / 7 | 100.0% | 204.0 |
| `ewma_hazard_score` | `intermittent_crc_timeout_burst` | 5 / 5 | 100.0% | 96.0 |
| `ewma_hazard_score` | `sudden_unannounced_controller_death` | 0 / 3 | 0.0% | 0.0 |
| `ewma_hazard_score` | `hard_negative_hot_tier` | 3 / 115 | 2.6% | 0.0 |
| `ewma_hazard_score` | `healthy` | 0 / 461 | 0.0% | 0.0 |
| `failureminer_temporal_pattern` | `progressive_media_error_ramp` | 8 / 8 | 100.0% | 162.0 |
| `failureminer_temporal_pattern` | `thermal_spare_depletion` | 7 / 7 | 100.0% | 258.0 |
| `failureminer_temporal_pattern` | `intermittent_crc_timeout_burst` | 5 / 5 | 100.0% | 108.0 |
| `failureminer_temporal_pattern` | `sudden_unannounced_controller_death` | 0 / 3 | 0.0% | 0.0 |
| `failureminer_temporal_pattern` | `hard_negative_hot_tier` | 1 / 115 | 0.9% | 0.0 |
| `failureminer_temporal_pattern` | `healthy` | 0 / 461 | 0.0% | 0.0 |
| `over_sensitive_negative` | `progressive_media_error_ramp` | 8 / 8 | 100.0% | 867.0 |
| `over_sensitive_negative` | `thermal_spare_depletion` | 7 / 7 | 100.0% | 1296.0 |
| `over_sensitive_negative` | `intermittent_crc_timeout_burst` | 5 / 5 | 100.0% | 894.0 |
| `over_sensitive_negative` | `sudden_unannounced_controller_death` | 3 / 3 | 100.0% | 1302.0 |
| `over_sensitive_negative` | `hard_negative_hot_tier` | 115 / 115 | 100.0% | 0.0 |
| `over_sensitive_negative` | `healthy` | 460 / 461 | 99.8% | 0.0 |

## 3. Strict Temporal Causality & Zero-Leakage Audit

| Method | Train/Test Disjoint | Future Perturbation Invariant | Drives / Cutoffs / Observations |
|---|---:|---:|---:|
| `static_smart_threshold` | true | true | 12 / 36 / 4751 |
| `sliding_window_zscore` | true | true | 12 / 36 / 4751 |
| `ewma_hazard_score` | true | true | 12 / 36 / 4751 |
| `failureminer_temporal_pattern` | true | true | 12 / 36 / 4751 |
| `over_sensitive_negative` | true | true | 12 / 36 / 4751 |

## 4. Honest Failure & Boundary Regimes

### `sudden_unannounced_controller_death`

- **Setup**: 100% of failing drives suffer abrupt ASIC/power controller death with zero pre-failure SMART/NVMe telemetry drift.
- **Finding**: Proves the physical recall ceiling: calibrated causal models achieve ~0% recall (at ~0% FPR) because no precursor exists, whereas over-sensitive heuristics only catch sudden deaths accidentally by falsely flagging the majority of healthy drives. Proactive prediction complements, and never replaces, erasure coding / replication.

| Method | TP/FP/FN | Recall | Precision | F1 | Drive FPR | Median Lead (h) | Cost Savings | Wasted Evac (TB) |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `static_smart_threshold` | 0/4/9 | 0.0% | 0.0% | 0.000 | 2.84% | 0.0 | -3.7% | 30.7 |
| `failureminer_temporal_pattern` | 0/0/9 | 0.0% | 0.0% | 0.000 | 0.00% | 0.0 | +0.0% | 0.0 |
| `over_sensitive_negative` | 9/141/0 | 100.0% | 6.0% | 0.113 | 100.00% | 1314.0 | -38.9% | 1082.9 |

### `high_imbalance_base_rate_precision_trap`

- **Setup**: 1,000-drive fleet with 1.0% failure prevalence and 25% hot-tier hard negatives (stratified 50/50 split; ~500 test drives).
- **Finding**: Demonstrates Bayes' base-rate trap: at 1% failure prevalence, even a 1.2-6% drive FPR triggers more unnecessary evacuations than real failures (precision 12-40%). Tight FPR budgeting (<=0.5%) + multi-signal motif gating is mandatory before enabling automated volume drain.

| Method | TP/FP/FN | Recall | Precision | F1 | Drive FPR | Median Lead (h) | Cost Savings | Wasted Evac (TB) |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `static_smart_threshold` | 4/30/0 | 100.0% | 11.8% | 0.210 | 6.06% | 42.0 | +29.2% | 230.4 |
| `ewma_hazard_score` | 4/1/0 | 100.0% | 80.0% | 0.889 | 0.20% | 105.0 | +89.6% | 7.7 |
| `failureminer_loose_fpr_3pct` | 4/6/0 | 100.0% | 40.0% | 0.571 | 1.21% | 129.0 | +79.2% | 46.1 |
| `failureminer_tight_fpr_0.4pct` | 4/0/0 | 100.0% | 100.0% | 1.000 | 0.00% | 129.0 | +91.7% | 0.0 |
| `over_sensitive_negative` | 4/494/0 | 100.0% | 0.8% | 0.016 | 99.80% | 1305.0 | -937.5% | 3793.9 |

### `cross_vendor_distribution_shift`

- **Setup**: Models trained exclusively on Vendor A (cool 37C Gen4 TLC, 280us p99 read) and evaluated zero-shot on unseen Vendor B (QLC, 590us p99 read) and Vendor C (warm 51.5C Gen5 TLC).
- **Finding**: Thresholds calibrated globally on Vendor A collapse under cross-vendor thermal/latency shift (surge in false alarms on healthy Vendor B/C drives). Per-vendor calibration plus causal online per-drive burn-in normalization restores low FPR and high precision zero-shot.

| Method | TP/FP/FN | Recall | Precision | F1 | Drive FPR | Median Lead (h) | Cost Savings | Wasted Evac (TB) |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `ewma_hazard_score` | 10/27/2 | 83.3% | 27.0% | 0.408 | 11.84% | 156.0 | +57.6% | 207.4 |
| `failureminer_no_vendor_norm` | 11/141/1 | 91.7% | 7.2% | 0.134 | 61.84% | 276.0 | -13.9% | 1082.9 |
| `failureminer_with_causal_burnin_norm` | 10/0/2 | 83.3% | 100.0% | 0.909 | 0.00% | 150.0 | +76.4% | 0.0 |

