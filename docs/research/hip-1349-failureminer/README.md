# HIP-1349 / HIP-1348 Research Evaluation: Predictive SSD Failure Detection (FailureMiner)

- **Tracking Issues**: `hanzoai/s3#29` (HIP-1349) and `hanzoai/s3#33` (HIP-1348 umbrella)
- **Primary Paper**: Shuyang Wang et al., *FailureMiner* (USENIX FAST '26), <https://www.usenix.org/conference/fast26/presentation/wang-shuyang>
- **Deliverable Class**: Clean-room Go stdlib-only research harness, online causal scoring library, reproducible benchmark CLI, and raw experimental evidence.

> **Honesty & Scope Notice**: All measurements in `score_table.md` and `raw_results.json` were produced on the authoring host (`go1.27.1`, `linux/amd64`, 2 CPUs) using a deterministic seeded NVMe/SMART fleet simulation (`real_hardware_exercised=false`). **These figures are evaluation-harness outputs, not Hanzo production SLA expectations**, until reproduced on real Hanzo production NVMe telemetry via `--external-smart-csv` or `--external-smart-json`. Published numbers in the FAST '26 paper belong strictly to the original authors' datacenter environment.

---

## 1. Paper Summary & License Audit

### Paper Design Summary
FailureMiner (Wang et al., USENIX FAST '26) addresses large-scale cloud SSD failure prediction where static SMART vendor thresholds either trigger too late for safe volume evacuation or suffer severe false-alarm storms across heterogeneous drive models and hot-tier workloads. The key architectural insight is that impending flash and controller failures manifest as **co-occurring multi-signal temporal progression patterns** (e.g., error-log growth coupled with reallocated-sector acceleration and read tail-latency drift, or sustained thermal excess paired with negative available-spare slope and GC write-latency amplification) rather than isolated single-attribute threshold crossings.

### Clean-Room & License Audit
| Artifact / Component | Provenance | License / Compatibility | Status |
|---|---|---|---|
| `s3/health/failureminer/` | Clean-room Go implementation written from first principles in this PR | Apache-2.0 (repository license) | Cleared — zero upstream paper code copied |
| `cmd/failureminer-bench/` | Clean-room Go benchmark CLI | Apache-2.0 (repository license) | Cleared — Go stdlib only |
| External Go dependencies (`go.mod`) | None added | Unchanged | Zero new dependencies |
| Primary paper concepts | Publicly presented USENIX FAST '26 research methodology | Public academic publication; no code imported | Clean-room algorithmic inspiration only |
| External SMART datasets (optional CLI input) | Compatible with Backblaze Drive Stats / Alibaba SSD SMART CSV schemas via `--external-smart-csv` | Datasets are **not** vendored in this repo; operators supply local telemetry exports under their own license terms | Cleared — no third-party dataset blobs committed |
| Patent / IP check | Uses standard statistical signal processing (robust median/MAD normalization, causal OLS slope, rolling ring buffers, rule/motif persistence scoring, empirical quantile FPR budgeting) | Standard prior-art statistical methods; no proprietary firmware or model weights imported | Cleared for research & clean-room production use |

---

## 2. Matched Baselines & Reproducible Methodology

All five predictors implement the identical online `Predictor` / `StreamScorer` interface in [`s3/health/failureminer/types.go`](../../../s3/health/failureminer/types.go) and are evaluated on the exact same failure-mode-stratified 50/50 train/test split (`601` train drives, `599` test drives, `0` overlapping drive IDs):

1. **`static_smart_threshold`** — Standard NVMe/SMART static vendor rule: alert if `critical_warning_mask > 0 || available_spare_pct < 10% || media_errors > 0`.
2. **`sliding_window_zscore`** — Causal univariate anomaly detector over each drive's rolling 24-sample (6-day) `[t-W, t-1]` window of p99 read/write latency and media/error-log deltas (`MinConsecutive=2`, calibrated on healthy training drives).
3. **`ewma_hazard_score`** — Exponentially weighted moving average (EWMA) derivative/trend hazard model over SMART error deltas, spare depletion, and globally normalized thermal/latency signals, calibrated on healthy training drives to a 1.5% FPR budget.
4. **`failureminer_temporal_pattern`** — Clean-room FailureMiner-inspired multi-feature temporal pattern and progression-signature miner:
   - Per-vendor/per-model robust median/MAD normalization + causal online per-drive burn-in baseline adaptation (`<= t` only).
   - Multi-horizon causal windows (24h short / 72h medium / 168h long) extracting counter deltas, active-step recurrence counts, and OLS available-spare depletion slope (%/day).
   - Four co-occurring multi-signal degradation motifs (`motif_media_ramp`, `motif_thermal_spare_cascade`, `motif_crc_timeout_storm`, `motif_critical_hardware`) combined with a cross-window persistence accumulator and empirical FPR-budgeted threshold calibration (`TargetDriveFPR=1.0%`).
5. **`over_sensitive_negative`** (Documented Negative Control) — Ultra-low threshold alerting on moderate temperature (`> 52°C`), single latency spikes (`p99_read > 900us` or `p99_write > 1350us`), or any single error counter. Demonstrates false-alarm storms and evacuation cost collapse.

### Reproduction Commands
```bash
# Run the full 5-method comparison, causality audit, and 3 boundary regimes:
go run ./cmd/failureminer-bench --seed 1349 --drives 1200 --horizon-days 90 --step-hours 6 \
  --out-dir docs/research/hip-1349-failureminer

# Run identical evaluation on external Backblaze/Alibaba-compatible SMART CSV logs:
go run ./cmd/failureminer-bench --external-smart-csv /path/to/nvme_smart.csv

# Verify static analysis, build, and race-free unit/causality/isolation tests:
go vet ./s3/health/failureminer/... ./cmd/failureminer-bench/...
go build -o /dev/null ./cmd/failureminer-bench
go test -race -v -count=1 ./s3/health/failureminer/...
```

---

## 3. Falsifiable Hypothesis & Expected Failure Modes

### Upfront Falsifiable Hypothesis
On a 90-day multi-vendor NVMe telemetry fleet with realistic class imbalance (4% failing drives), 20% hot-tier hard-negative healthy drives, and 15% zero-precursor sudden controller deaths, `failureminer_temporal_pattern` will:
- **(H1 — Recall & Actionable Lead Time)**: Achieve higher overall recall, higher >=24h actionable recall, and longer median warning lead time than `static_smart_threshold` and `sliding_window_zscore`.
- **(H2 — False-Alarm Control on Hot-Tier Hard Negatives)**: Reduce drive-level FPR and wasted evacuation TB relative to `static_smart_threshold` by rejecting static legacy counters and isolated single-window workload/compaction spikes.
- **(H3 — Honest Causal Ceiling)**: Remain strictly below 100% recall (`<= 1 - sudden_death_share`) without leaking future information.

### Expected Failure & Boundary Regimes (Explicitly Tested)
1. **`sudden_unannounced_controller_death`**: Hardware ASIC/power-rail death produces zero pre-failure telemetry drift. Any calibrated causal predictor must exhibit ~0% recall on this subset; proactive prediction cannot replace erasure coding or replication.
2. **`high_imbalance_base_rate_precision_trap`**: At 1.0% failure prevalence, Bayes' rule magnifies even a 1.2–6% drive FPR into >60% false evacuations (precision 12–40%). Automated drain requires tight FPR budgeting (`<= 0.5%`).
3. **`cross_vendor_distribution_shift`**: Thresholds calibrated globally on cool Gen4 TLC (`vendor_a_nvme`, 37°C, 280us p99 read) suffer severe false-alarm collapse when evaluated zero-shot on unseen QLC (`vendor_b_nvme`, 590us) and warm Gen5 TLC (`vendor_c_nvme`, 51.5°C) unless per-vendor or causal online per-drive burn-in normalization is active.

---

## 4. Bounded PR Deliverable

| File | Purpose |
|---|---|
| [`s3/health/failureminer/types.go`](../../../s3/health/failureminer/types.go) | NVMe Log Page 02h / SMART `DriveTelemetrySample` schema, `FailureMode` taxonomy, stratified `SplitTrainTest`, validation, and strictly causal `Predictor` / `StreamScorer` interfaces |
| [`s3/health/failureminer/corpus.go`](../../../s3/health/failureminer/corpus.go) | Deterministic seeded multi-vendor NVMe fleet simulator with class imbalance, 4 failure modes, hard-negative hot-tier drives, and transient benign datacenter noise episodes |
| [`s3/health/failureminer/loader.go`](../../../s3/health/failureminer/loader.go) | External SMART CSV (Backblaze / Alibaba column aliases) and JSON loaders with strict chronological sorting and pre-failure prefix extraction |
| [`s3/health/failureminer/predictors.go`](../../../s3/health/failureminer/predictors.go) | Clean-room implementations of all 5 causal failure-prediction methods |
| [`s3/health/failureminer/eval.go`](../../../s3/health/failureminer/eval.go) | Evaluation metrics, operational cost model, future-perturbation temporal causality auditor, and 3 boundary-regime stress tests |
| [`s3/health/failureminer/fleet_monitor.go`](../../../s3/health/failureminer/fleet_monitor.go) | Thread-safe, multi-tenant isolated online streaming monitor enforcing strict per-drive timestamp monotonicity |
| [`s3/health/failureminer/failureminer_test.go`](../../../s3/health/failureminer/failureminer_test.go) | Unit, reproducibility, zero-leakage causality, boundary-regime, CSV/JSON loader, and `-race` concurrent multi-tenant tests |
| [`cmd/failureminer-bench/main.go`](../../../cmd/failureminer-bench/main.go) | Reproducible Go benchmark CLI emitting stdout summary, `score_table.md`, and `raw_results.json` |
| `docs/research/hip-1349-failureminer/{README.md,score_table.md,raw_results.json}` | Complete 6-point HIP-1349 / HIP-1348 research documentation and raw measurement artifacts |

No existing files in `s3/` or `go.mod` were modified, keeping this research PR strictly bounded and zero-risk to the production I/O path.

---

## 5. Correctness, Concurrency, Causality & Multi-Tenant Evidence

All checks pass cleanly with `go vet` and `go test -race -v -count=1 ./s3/health/failureminer/...`:

- **Deterministic Reproducibility (`TestDeterministicCorpusReproducibility`)**: Identical seeds produce byte-for-byte identical `FleetDataset` histories, strictly monotonic timestamps, and zero post-failure samples.
- **Temporal Causality & Zero Look-Ahead Leakage (`TestTemporalCausalityZeroLeakage`)**: Across all 5 predictors, 12 diverse test drives, 36 cutoffs, and **4,751 causal observations**, mutating all future telemetry (`> t`) to extreme failure states (`media_errors += 5000`, `temperature = 92°C`, `available_spare = 1%`, `critical_warning = 0x1F`) produces **bit-identical** `Prediction` structs at every timestamp `<= t`. Train and test drive ID sets are verified disjoint.
- **Concurrent Multi-Tenant Safety (`TestMultiTenantConcurrentFleetMonitor`)**:
  - 40 concurrent goroutines stream telemetry simultaneously into `FleetMonitor` under Go's race detector (`-race`) with zero data races.
  - Strict `(tenant_id, drive_id)` state partitioning ensures a failing drive in `tenant-isolated-failing` with the exact same `drive_id` causes **0 alerts** on `tenant-isolated-healthy`.
  - Input validation rejects NaN/Inf telemetry, out-of-range spare/temperature values, empty IDs, and non-monotonic (out-of-order) timestamps.
- **External Schema Interoperability (`TestExternalSMARTCSVAndJSONLoader`)**: Verifies Backblaze/Alibaba-style CSV parsing (`serial_number`, `date`, `smart_187_raw`, `smart_232_normalized`, `smart_194_raw`, `smart_5_raw`, `smart_199_raw`, `failure`) with unsorted rows, pre-failure truncation, and JSON round-trip equality.

---

## 6. Raw Measurements & HIP-1349 / HIP-1348 Integration Recommendation

### Headline Results (Test Split: 599 Drives — 20 Predictable Failures, 3 Sudden Deaths, 115 Hard-Negative Hot-Tier, 461 Healthy)

Full per-mode and raw JSON metrics are recorded in [`score_table.md`](./score_table.md) and [`raw_results.json`](./raw_results.json).

| Method | TP/FP/FN | Recall | Pred. Recall | Precision | F1 | Drive FPR | FAR / Drive-Day | Median Lead (h) | >=24h Recall | Hard-Neg FP | Cost Savings | Wasted Evac (TB) | Latency (ns/samp) |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| `static_smart_threshold` | 14/24/9 | 60.9% | 70.0% | 36.8% | 0.459 | 4.17% | 0.000464 | 72.0 | 34.8% | 24/115 | +47.1% | 184.3 | 49 |
| `sliding_window_zscore` | 10/14/13 | 43.5% | 50.0% | 41.7% | 0.426 | 2.43% | 0.000271 | 114.0 | 39.1% | 10/115 | +34.8% | 107.5 | 396 |
| `ewma_hazard_score` | 20/3/3 | 87.0% | 100.0% | 87.0% | 0.870 | 0.52% | 0.000058 | 126.0 | 87.0% | 3/115 | +78.6% | 23.0 | 148 |
| **`failureminer_temporal_pattern`** | **20/1/3** | **87.0%** | **100.0%** | **95.2%** | **0.909** | **0.17%** | **0.000019** | **168.0** | **87.0%** | **1/115** | **+79.3%** | **7.7** | **320** |
| `over_sensitive_negative` | 23/575/0 | 100.0% | 100.0% | 3.9% | 0.074 | 99.83% | 0.531202 | 1044.0 | 100.0% | 115/115 | -116.7% | 4416.0 | 77 |

### Key Takeaways & Boundary Findings
1. **Substantial Lead-Time & Precision Advantage over Static SMART**: `failureminer_temporal_pattern` increases overall recall from `60.9%` to `87.0%` (+26.1 pp), increases >=24h actionable recall from `34.8%` to `87.0%` (+52.2 pp), extends median warning lead time from `72.0h` to `168.0h` (7.0 days; `2.33x`), and cuts false evacuations from `24` drives (`184.3 TB`) down to `1` drive (`7.7 TB` — a 24x reduction).
2. **Motif Gating Beats Additive EWMA**: At higher precision (`95.2%` vs `87.0%`) and 3x lower FPR (`0.17%` vs `0.52%`), `failureminer_temporal_pattern` delivers `+42h` longer median lead time (`168.0h` vs `126.0h`) and `3x` less wasted evacuation bandwidth (`7.7 TB` vs `23.0 TB`) than `ewma_hazard_score`, at `~3.1M samples/sec` (`320 ns/sample`).
3. **Honest Sudden-Death Ceiling**: On `sudden_unannounced_controller_death` (100% zero-precursor failures), `failureminer_temporal_pattern` scores `0.0%` recall at `0.00%` FPR (`+0.0%` net cost delta), whereas `over_sensitive_negative` causes `-38.9%` cost loss and `1,082.9 TB` of wasted evacuations.
4. **Base-Rate Precision Trap Confirmed**: At `1.0%` failure prevalence, loosening the FPR budget (`failureminer_loose_fpr_3pct`, `1.21%` test FPR) drops precision to `40.0%` (`6` false evacuations vs `4` real failures), while tight FPR budgeting (`failureminer_tight_fpr_0.4pct`) holds `0.00%` test FPR and `100.0%` precision.
5. **Cross-Vendor Distribution Shift Confirmed**: Training solely on Vendor A and evaluating zero-shot on unseen Vendors B & C causes unnormalized FailureMiner (`failureminer_no_vendor_norm`) to collapse to `61.84%` FPR / `7.2%` precision (`-13.9%` cost savings) and `ewma_hazard_score` to drop to `11.84%` FPR / `27.0%` precision. Enabling causal online per-drive burn-in normalization (`failureminer_with_causal_burnin_norm`) recovers `0.00%` FPR, `100.0%` precision, and `150.0h` median lead time zero-shot.

### HIP-1349 / HIP-1348 Integration Recommendation
**Recommendation**: **Proceed to Phase-1 Shadow Mode** in `hanzoai/s3`, gated by real NVMe telemetry validation before enabling automated proactive drain.

Canonical architecture boundaries remain strictly intact (`vfs -> hanzoai/s3 -> visor`, Cloud owning lifecycle & billing, Enso routing models, Zen generating, Kai returning structured non-generative decisions):
1. **Telemetry Collection (`hanzoai/s3`)**: Extend the volume server's disk probe loop in [`s3/stats/disk.go`](../../../s3/stats/disk.go) (alongside `DiskIOProbeConfig` and `diskState`) with an opt-in Linux NVMe Log Page 02h (`nvme_smart_log`) reader emitting `failureminer.DriveTelemetrySample` at 15–60 minute cadence into `failureminer.FleetMonitor` (negligible CPU overhead at ~320 ns/sample).
2. **Shadow-Mode Validation (`visor` + `hanzoai/s3`)**: Export `RiskScore`, `Alert`, `LeadTimeEstimateHours`, and `ActiveMotifs` to Prometheus/Visor (`s3/stats/metrics.go`) for 30–60 days on real Hanzo storage clusters. Replay collected CSV logs through `cmd/failureminer-bench --external-smart-csv` to re-verify FPR `<= 0.5%` on real drive models.
3. **Structured Drain Decisions (`kai` -> `hanzoai/s3` workers)**: Once shadow-mode precision exceeds 80% on real hardware, route high-confidence motif alerts (with `TargetDriveFPR <= 0.005` and rate-limited concurrent drain caps) to Kai for structured drain approval, triggering existing `s3/worker/tasks/balance` and `s3/worker/tasks/ec_balance` tasks to evacuate volumes before degraded SSD stalls impact `vfs` / `hanzoai/s3` tail latency.
4. **HIP-1349 / HIP-1348 Status Update**: Research phase complete with reproducible harness, zero-leakage proof, negative controls, and explicit boundary-regime characterization. Ready for maintainer review toward an implementation bounty gated on real-cluster SMART replay.
