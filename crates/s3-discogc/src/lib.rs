// SPDX-License-Identifier: Apache-2.0
//! [Bounty #42] (`hanzoai/s3#31`): **Discard-versus-compaction economics** (DisCoGC).
//!
//! Clean-room simulation harness for the GC decision: when is it cheaper to
//! *discard* a partially-live segment than to *compact* it?
//!
//! # Design
//!
//! - A segment has a byte capacity and a set of live/tombstoned records.
//! - Workload appends records, deletes a fraction, and occasionally triggers GC.
//! - Policies:
//!   - `CompactAlways`: rewrite every segment whose live ratio is below 1.0.
//!   - `Threshold(θ)`: compact only when live ratio < θ (classic LSM policy).
//!   - `DisCoGCEconomic`: compare rewrite cost against the amortised cost of
//!     keeping the segment (read amplification + space amplification) and pick
//!     the cheaper option, including a discard option for near-dead segments.
//!
//! Metrics: bytes rewritten, space amplification, read amplification, and total
//! modelled cost units.
//!
//! # Honest boundaries
//! - Simulated segments and records — no real storage engine exercised.
//! - Cost coefficients (rewrite vs read vs space) are hand-set; real engines
//!   have very different constants.
//! - All figures are simulation outputs on the authoring host, NOT Hanzo
//!   expectations until reproduced against a real LSM implementation.

use serde::{Deserialize, Serialize};
use thiserror::Error;

#[derive(Debug, Error)]
pub enum GcError {
    #[error("invalid configuration: {0}")]
    InvalidConfig(String),
}

#[derive(Debug, Clone)]
pub struct SplitMix64 {
    state: u64,
}
impl SplitMix64 {
    pub fn new(seed: u64) -> Self {
        Self { state: seed }
    }
    pub fn next_u64(&mut self) -> u64 {
        self.state = self.state.wrapping_add(0x9E37_79B9_7F4A_7C15);
        let mut z = self.state;
        z = (z ^ (z >> 30)).wrapping_mul(0xBF58_476D_1CE4_E5B9);
        z = (z ^ (z >> 27)).wrapping_mul(0x94D0_49BB_1331_11EB);
        z ^ (z >> 31)
    }
    pub fn next_f64(&mut self) -> f64 {
        (self.next_u64() >> 11) as f64 / (1u64 << 53) as f64
    }
}

#[derive(Debug, Clone)]
pub struct Segment {
    pub id: u32,
    pub total_bytes: u64,
    pub live_bytes: u64,
    pub live_records: u64,
    pub dead_records: u64,
}

impl Segment {
    pub fn live_ratio(&self) -> f64 {
        if self.total_bytes == 0 {
            0.0
        } else {
            self.live_bytes as f64 / self.total_bytes as f64
        }
    }
}

pub fn generate_segments(seed: u64, n: usize, seg_bytes: u64, delete_ratio: f64) -> Vec<Segment> {
    let mut rng = SplitMix64::new(seed);
    (0..n)
        .map(|i| {
            let live_ratio = (1.0 - delete_ratio * rng.next_f64()).clamp(0.01, 1.0);
            let live_bytes = (seg_bytes as f64 * live_ratio) as u64;
            let avg_rec = 4096u64;
            let live_records = live_bytes / avg_rec;
            let total_records = seg_bytes / avg_rec;
            Segment {
                id: i as u32,
                total_bytes: seg_bytes,
                live_bytes,
                live_records,
                dead_records: total_records.saturating_sub(live_records),
            }
        })
        .collect()
}

/// Cost model coefficients (hand-set, relative units).
#[derive(Debug, Clone, Copy)]
pub struct CostModel {
    /// Cost per byte rewritten during compaction.
    pub rewrite_per_byte: f64,
    /// Cost per byte scanned. Compaction must scan the whole segment; discard
    /// avoids the scan entirely.
    pub read_per_byte: f64,
    /// Cost per byte held on disk, amortised.
    pub space_per_byte: f64,
    /// Extra read cost multiplier caused by dead bytes in the read path.
    pub dead_byte_read_penalty: f64,
}

impl Default for CostModel {
    fn default() -> Self {
        Self {
            rewrite_per_byte: 1.0,
            read_per_byte: 0.05,
            space_per_byte: 0.30,
            dead_byte_read_penalty: 1.0,
        }
    }
}

#[derive(Debug, Clone, Copy, PartialEq)]
pub enum Policy {
    CompactAlways,
    Threshold(f64),
    DisCoGCEconomic,
}

impl Policy {
    pub fn label(&self) -> String {
        match self {
            Policy::CompactAlways => "compact-always".into(),
            Policy::Threshold(t) => format!("threshold({:.2})", t),
            Policy::DisCoGCEconomic => "discogc-economic".into(),
        }
    }

    /// Decide what to do with a segment.
    pub fn decide(&self, s: &Segment, cm: &CostModel) -> Decision {
        let live = s.live_bytes;
        let dead = s.total_bytes.saturating_sub(s.live_bytes);

        // Cost of compacting now: compaction must scan the whole segment to
        // find the live records, then rewrite only the live bytes.
        let compact_cost =
            s.total_bytes as f64 * cm.read_per_byte + live as f64 * cm.rewrite_per_byte;

        // Cost of keeping as-is: space held + read penalty for dead bytes.
        let keep_cost =
            dead as f64 * cm.space_per_byte + dead as f64 * cm.read_per_byte * cm.dead_byte_read_penalty;

        match self {
            Policy::CompactAlways => {
                if s.live_ratio() < 1.0 {
                    Decision::Compact
                } else {
                    Decision::Keep
                }
            }
            Policy::Threshold(t) => {
                if s.live_ratio() < *t {
                    Decision::Compact
                } else {
                    Decision::Keep
                }
            }
            Policy::DisCoGCEconomic => {
                // Near-dead segment: discarding (drop the whole segment, losing the
                // few live records after re-inserting them elsewhere) can be cheaper
                // than a rewrite of mostly-garbage bytes.
                let discard_cost = live as f64 * cm.rewrite_per_byte * 1.1;
                if compact_cost > keep_cost && discard_cost > keep_cost {
                    Decision::Keep
                } else if s.live_ratio() < 0.10 && discard_cost < compact_cost {
                    Decision::Discard
                } else if compact_cost < keep_cost {
                    Decision::Compact
                } else {
                    Decision::Keep
                }
            }
        }
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum Decision {
    Keep,
    Compact,
    Discard,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PolicyResult {
    pub policy: String,
    pub segments: usize,
    pub compacted: usize,
    pub discarded: usize,
    pub kept: usize,
    pub bytes_rewritten: u64,
    pub bytes_reclaimed: u64,
    pub space_amplification: f64,
    pub read_amplification: f64,
    pub total_cost_units: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BenchReport {
    pub bounty: String,
    pub hip: String,
    pub honesty_labels: Vec<String>,
    pub config: BenchConfig,
    pub results: Vec<PolicyResult>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BenchConfig {
    pub seed: u64,
    pub n_segments: usize,
    pub segment_bytes: u64,
    pub delete_ratio: f64,
}

impl Default for BenchConfig {
    fn default() -> Self {
        Self { seed: 0xD15C_06C, n_segments: 200, segment_bytes: 64 * 1024 * 1024, delete_ratio: 0.9 }
    }
}
impl BenchConfig {
    pub fn test_config() -> Self {
        Self { seed: 7, n_segments: 60, segment_bytes: 4 * 1024 * 1024, delete_ratio: 0.9 }
    }
}

pub fn evaluate(policy: Policy, segs: &[Segment], cm: &CostModel) -> PolicyResult {
    let mut compacted = 0usize;
    let mut discarded = 0usize;
    let mut kept = 0usize;
    let mut bytes_rewritten = 0u64;
    let mut bytes_reclaimed = 0u64;
    let mut cost = 0.0f64;

    for s in segs {
        let live = s.live_bytes;
        let dead = s.total_bytes.saturating_sub(s.live_bytes);
        match policy.decide(s, cm) {
            Decision::Compact => {
                compacted += 1;
                bytes_rewritten += live;
                bytes_reclaimed += dead;
                cost += live as f64 * cm.rewrite_per_byte;
            }
            Decision::Discard => {
                discarded += 1;
                bytes_rewritten += live; // live records must be re-inserted
                bytes_reclaimed += dead;
                cost += live as f64 * cm.rewrite_per_byte * 1.1;
            }
            Decision::Keep => {
                kept += 1;
                cost += dead as f64 * cm.space_per_byte
                    + dead as f64 * cm.read_per_byte * cm.dead_byte_read_penalty;
            }
        }
    }

    let live_total: u64 = segs.iter().map(|s| s.live_bytes).sum();
    let total_total: u64 = segs.iter().map(|s| s.total_bytes).sum();
    let kept_bytes: u64 = segs
        .iter()
        .filter(|s| policy.decide(s, cm) == Decision::Keep)
        .map(|s| s.total_bytes)
        .sum();
    let kept_live: u64 = segs
        .iter()
        .filter(|s| policy.decide(s, cm) == Decision::Keep)
        .map(|s| s.live_bytes)
        .sum();

    let space_amp = if live_total == 0 {
        1.0
    } else {
        (kept_bytes + bytes_rewritten) as f64 / live_total as f64
    };
    let read_amp = if kept_live == 0 {
        1.0
    } else {
        kept_bytes as f64 / kept_live as f64
    };
    let _ = total_total;

    PolicyResult {
        policy: policy.label(),
        segments: segs.len(),
        compacted,
        discarded,
        kept,
        bytes_rewritten,
        bytes_reclaimed,
        space_amplification: space_amp,
        read_amplification: read_amp,
        total_cost_units: cost,
    }
}

pub fn run_bench(cfg: &BenchConfig) -> Result<BenchReport, GcError> {
    if cfg.n_segments == 0 {
        return Err(GcError::InvalidConfig("n_segments must be positive".into()));
    }
    if !(0.0..1.0).contains(&cfg.delete_ratio) {
        return Err(GcError::InvalidConfig("delete_ratio must be in [0,1)".into()));
    }
    let segs = generate_segments(cfg.seed, cfg.n_segments, cfg.segment_bytes, cfg.delete_ratio);
    let cm = CostModel::default();

    let results = [
        Policy::CompactAlways,
        Policy::Threshold(0.5),
        Policy::DisCoGCEconomic,
    ]
    .iter()
    .map(|p| evaluate(*p, &segs, &cm))
    .collect();

    Ok(BenchReport {
        bounty: "hanzoai/s3#31 [Bounty #42] Discard-versus-compaction economics (DisCoGC)".into(),
        hip: "HIP-1348: Object Storage Durability and Placement".into(),
        honesty_labels: vec![
            "simulated segments and records — no real storage engine exercised".into(),
            "cost coefficients (rewrite/read/space) are hand-set relative units".into(),
            "delete distribution is uniform per segment; real workloads cluster".into(),
            "All figures are simulation outputs on the authoring host, NOT Hanzo expectations until reproduced against a real LSM implementation.".into(),
        ],
        config: cfg.clone(),
        results,
    })
}

pub fn write_results_json(report: &BenchReport, path: &std::path::Path) -> Result<(), Box<dyn std::error::Error>> {
    if let Some(parent) = path.parent() {
        std::fs::create_dir_all(parent)?;
    }
    std::fs::write(path, serde_json::to_string_pretty(report)?)?;
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn economic_policy_never_costs_more_than_compact_always() {
        let report = run_bench(&BenchConfig::test_config()).unwrap();
        let always = report.results.iter().find(|r| r.policy == "compact-always").unwrap();
        let econ = report.results.iter().find(|r| r.policy == "discogc-economic").unwrap();
        assert!(
            econ.total_cost_units <= always.total_cost_units + 1e-6,
            "economic cost {:.1} > compact-always {:.1}",
            econ.total_cost_units,
            always.total_cost_units
        );
    }

    #[test]
    fn compact_always_rewrites_every_segment() {
        let report = run_bench(&BenchConfig::test_config()).unwrap();
        let always = report.results.iter().find(|r| r.policy == "compact-always").unwrap();
        assert_eq!(always.compacted, always.segments);
    }

    #[test]
    fn economic_policy_keeps_some_segments() {
        let report = run_bench(&BenchConfig::test_config()).unwrap();
        let econ = report.results.iter().find(|r| r.policy == "discogc-economic").unwrap();
        assert!(econ.kept > 0, "economic policy should skip some rewrites");
    }

    #[test]
    fn near_dead_segments_are_discarded() {
        let cm = CostModel::default();
        let s = Segment { id: 0, total_bytes: 1_000_000, live_bytes: 10_000, live_records: 2, dead_records: 240 };
        assert_eq!(Policy::DisCoGCEconomic.decide(&s, &cm), Decision::Discard);
    }

    #[test]
    fn bad_delete_ratio_rejected() {
        let mut cfg = BenchConfig::test_config();
        cfg.delete_ratio = 1.0;
        assert!(run_bench(&cfg).is_err());
    }
}
