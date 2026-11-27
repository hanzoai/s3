# HIP-1349 / HIP-1348 Score Table: Geo-Aware Active Object Placement & Recovery (ACOS)

- **Reference**: Apple ACOS (Baron et al., USENIX FAST '26, https://www.usenix.org/conference/fast26/presentation/baron)
- **Seed / Workload**: `seed=20260414`, `400` objects x `16384` B real payloads, `4000` read events across 5 phases
- **Environment**: `go=go1.27.1 os=linux/amd64 cpus=2`
- **Honesty Disclosure**: All reported latency, cost, and TCO figures are reproducible harness outputs on the authoring host, NOT production Hanzo S3 hardware or billing claims until validated across live multi-region clusters.

## 1. Matched 5-Policy Comparison (Identical Multi-Region Workload)

| Policy | Read p50 / p95 / p99 (ms) | Write p50 / p99 (ms) | Failover p50 / p99 (ms) | RTO p95 (ms) | Overall / AZ / Blackout / Affected-Region Avail (%) | Storage Overhead & Cost ($/TB-mo) | Egress Cost ($/TB served) | Monthly TCO ($/TB-mo) | Decision Throughput (ops/s) | Migrations | Local AZ / Cross-WAN Recov | Violations / Corruptions |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `single_region_3x_local` | 2.32 / 197.37 / 197.60 | 3.21 / 3.72 | 3.21 / 3.21 | 0.00 | 87.81 / 100.00 / 38.88 / 0.00 | 3.00x ($74.49) | $14.98 | **$149.40** | 52692 | 0 | 407 / 0 | 0 / 0 |
| `naive_global_full_replication` | 1.91 / 69.27 / 143.27 | 169.37 / 269.58 | 1.81 / 3.21 | 16.81 | 100.00 / 100.00 / 100.00 / 100.00 | 4.80x ($120.30) | $18.43 | **$212.44** | 59723 | 0 | 200 / 489 | 0 / 0 |
| `static_cross_region_ec` | 109.22 / 166.42 / 169.22 | 109.73 / 169.64 | 144.52 / 166.42 | 159.52 | 100.00 / 100.00 / 100.00 / 100.00 | 1.50x ($35.92) | $15.97 | **$115.75** | 12360 | 0 | 0 / 1277 | 0 / 0 |
| `acos_geo_aware_adaptive` | 1.91 / 109.27 / 119.27 | 69.27 / 119.30 | 1.81 / 3.71 | 16.81 | 100.00 / 100.00 / 100.00 / 100.00 | 2.40x ($57.96) | $6.99 | **$92.92** | 28719 | 32 | 407 / 489 | 0 / 0 |
| `over_aggressive_migration_negative` | 2.01 / 218.55 / 278.55 | 69.27 / 119.30 | 3.71 / 3.71 | 16.81 | 100.00 / 100.00 / 100.00 / 100.00 | 2.40x ($58.24) | $25.02 | **$183.35** | 18560 | 680 | 387 / 32 | 0 / 0 |

## 2. Documented Negative & Boundary Regimes

### `over_aggressive_migration_thrashing`

- **Scenario**: Alternating 50/50 cross-region reads (us-east <-> ap-northeast) comparing hysteresis-guarded ACOS against zero-cooldown aggressive migration.
- **Measured Finding**: Zero-cooldown migration triggered 480 migrations and burned 11.84x cross-region egress cost ($131.00/TB vs $11.00/TB served) with p99 read latency 218.55 ms vs 109.27 ms for ACOS.

| Policy | Read p50 / p99 (ms) | Write p50 (ms) | Availability (%) | Egress ($/TB served) | Migrations | Metadata Overhead (B) | Corruptions |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `acos_geo_aware_adaptive` | 1.81 / 109.27 | 69.27 | 100.00% | $11.00 | 0 | 10752 | 0 |
| `over_aggressive_migration_negative` | 138.53 / 218.55 | 69.27 | 100.00% | $131.00 | 480 | 10752 | 0 |

### `simultaneous_dual_region_correlated_outage`

- **Scenario**: Simultaneous blackout of us-east and us-west exceeding f=1 fault-domain budget for us-east-homed objects.
- **Measured Finding**: Under f=2 simultaneous outage beyond f=1 budget, both acos_geo_aware_adaptive (0.0% avail, 0 corruptions) and static_cross_region_ec (0.0% avail, 0 corruptions) fail closed deterministically with zero silent data corruption; only 5x naive_global_full_replication survives (100.0% avail) at 2x+ storage cost.

| Policy | Read p50 / p99 (ms) | Write p50 (ms) | Availability (%) | Egress ($/TB served) | Migrations | Metadata Overhead (B) | Corruptions |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `acos_geo_aware_adaptive` | 0.00 / 0.00 | 69.23 | 0.00% | $0.00 | 0 | 22400 | 0 |
| `static_cross_region_ec` | 0.00 / 0.00 | 87.72 | 0.00% | $0.00 | 0 | 11200 | 0 |
| `naive_global_full_replication` | 1.91 / 1.91 | 169.33 | 100.00% | $80.00 | 0 | 9600 | 0 |

### `write_heavy_cold_once_objects`

- **Scenario**: Write-once rarely-read objects (200 writes, 20 single local reads) testing whether adaptive migration adds value on cold data.
- **Measured Finding**: ACOS triggers 0 migrations (0 ms read p50 gain: 1.91 ms vs 1.91 ms) while paying 27.18x write p50 latency for sync geo-standby (87.24 ms vs 3.21 ms) and 89600 B vs 25600 B metadata tracking overhead. Cold write-once archives belong in scheduled archival tiers (HIP-1342), not active adaptive placement.

| Policy | Read p50 / p99 (ms) | Write p50 (ms) | Availability (%) | Egress ($/TB served) | Migrations | Metadata Overhead (B) | Corruptions |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `single_region_3x_local` | 1.91 / 2.31 | 3.21 | 100.00% | $0.00 | 0 | 25600 | 0 |
| `acos_geo_aware_adaptive` | 1.91 / 2.31 | 87.24 | 100.00% | $520.00 | 0 | 89600 | 0 |

## 3. Tenant Data-Residency / Sovereignty Enforcement Audit

| Scenario | Regions | Audited Shards | Sovereignty Violations | Blocked Ops | Blackout Reads Succeeded | Outcome |
| :--- | :---: | :---: | :---: | :---: | :---: | :--- |
| `eu_only_single_eu_region_fail_closed` | 5 | 90 | 0 | 30 | 0 / 30 | PASS: fail-closed on sole EU region outage; 0 shards placed outside eu-* |
| `eu_only_multi_eu_region_sovereign_failover` | 6 | 120 | 0 | 0 | 30 / 30 | PASS: 100% failover strictly within eu-* (eu-central -> eu-west); 0 shards placed outside eu-* |

