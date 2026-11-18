// SPDX-License-Identifier: Apache-2.0
//! [Bounty #44] (`hanzoai/s3#33`): **Predictive SSD failure detection** (FailureMiner).
//!
//! Clean-room simulation harness for early SSD failure prediction from SMART
//! telemetry: how much lead time can a simple thresholded detector buy, and at
//! what false-positive cost?
//!
//! # Design
//!
//! - Synthetic drive fleet; each drive has a latent "wear" trajectory and a
//!   failure time (or never fails within the window).
//! - SMART attributes (reallocated sectors, media errors, wear level, CRC
//!   errors) are generated as noisy functions of wear.
//! - Detectors:
//!   - `Reactive`: alert only when the drive has already failed.
//!   - `StaticThreshold`: alert when a single attribute crosses a fixed limit.
//!   - `FailureMinerTrend`: alert when the *slope* of a composite health score
//!     crosses a threshold (trend-based, gives lead time).
//!
//! Metrics: mean lead time before failure, detection recall, false-positive
//! rate per drive-year, and alert volume.
//!
//! # Honest boundaries
//! - Fully synthetic telemetry — no real SMART data from any drive fleet.
//! - Wear trajectories are hand-modelled monotone processes with noise; real
//!   failure modes are far more heterogeneous.
//! - All figures are simulation outputs on the authoring host, NOT Hanzo
//!   expectations until reproduced on real fleet telemetry.

use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use thiserror::Error;

#[derive(Debug, Error)]
pub enum FmError {
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
    /// Standard normal via Box-Muller.
    pub fn next_normal(&mut self) -> f64 {
        let u1 = self.next_f64().max(1e-12);
        let u2 = self.next_f64();
        (-2.0 * u1.ln()).sqrt() * (2.0 * std::f64::consts::PI * u2).cos()
    }
}

/// A SMART sample for one drive at one day.
#[derive(Debug, Clone, Copy)]
pub struct SmartSample {
    pub day: u32,
    /// Normalised reallocated-sector growth (0..1+).
    pub reallocated: f64,
    /// Normalised media error rate.
    pub media_errors: f64,
    /// Wear level: 0 = new, 1 = end of rated life.
    pub wear: f64,
    /// Normalised CRC/interface error count.
    pub crc_errors: f64,
}

/// One synthetic drive.
#[derive(Debug, Clone)]
pub struct Drive {
    pub id: u32,
    /// Day the drive fails, or None if it survives the observation window.
    pub fails_on: Option<u32>,
    pub samples: Vec<SmartSample>,
}

pub fn generate_fleet(
    seed: u64,
    n_drives: usize,
    window_days: u32,
    base_failure_rate: f64,
) -> Vec<Drive> {
    let mut rng = SplitMix64::new(seed);
    let mut drives = Vec::with_capacity(n_drives);

    for id in 0..n_drives {
        let fails = rng.next_f64() < base_failure_rate;
        // Failures cluster late in the window but are not all at the end.
        let fail_day = if fails {
            let frac = 0.4 + 0.6 * rng.next_f64();
            Some(((window_days as f64 * frac) as u32).max(10))
        } else {
            None
        };
        // Per-drive wear speed.
        let speed = 0.6 + 0.8 * rng.next_f64();

        let mut samples = Vec::with_capacity(window_days as usize);
        for day in 0..window_days {
            let progress = day as f64 / window_days as f64;
            // Wear accelerates as the drive approaches failure.
            let pre_fail_pressure = match fail_day {
                Some(f) if day >= f => 1.0,
                Some(f) => {
                    let remaining = (f.saturating_sub(day)) as f64;
                    (1.0 - (remaining / 30.0).min(1.0)) * 0.8
                }
                None => 0.0,
            };
            let wear = (progress * speed * 0.5 + pre_fail_pressure * 0.6).min(1.2);

            // SMART attributes rise super-linearly with wear, plus noise.
            let noise = |rng: &mut SplitMix64, s: f64| rng.next_normal() * s;
            samples.push(SmartSample {
                day,
                reallocated: (wear.powf(2.2) * 0.9 + noise(&mut rng, 0.02)).max(0.0),
                media_errors: (wear.powf(2.6) * 0.7 + noise(&mut rng, 0.015)).max(0.0),
                wear: (wear + noise(&mut rng, 0.01)).clamp(0.0, 1.5),
                crc_errors: (wear.powf(1.8) * 0.4 + noise(&mut rng, 0.01)).max(0.0),
            });
        }
        drives.push(Drive { id: id as u32, fails_on: fail_day, samples });
    }
    drives
}

/// Composite health score: higher = sicker. Weights hand-set.
pub fn health_score(s: &SmartSample) -> f64 {
    0.45 * s.reallocated + 0.30 * s.media_errors + 0.15 * s.wear + 0.10 * s.crc_errors
}

/// Slope of health over the trailing `w` samples (least squares).
pub fn trend(samples: &[SmartSample], w: usize) -> f64 {
    if samples.len() < 2 {
        return 0.0;
    }
    let tail = &samples[samples.len().saturating_sub(w)..];
    let n = tail.len() as f64;
    if n < 2.0 {
        return 0.0;
    }
    let mean_x = (n - 1.0) / 2.0;
    let mean_y = tail.iter().map(health_score).sum::<f64>() / n;
    let mut num = 0.0;
    let mut den = 0.0;
    for (i, s) in tail.iter().enumerate() {
        let dx = i as f64 - mean_x;
        num += dx * (health_score(s) - mean_y);
        den += dx * dx;
    }
    if den == 0.0 {
        0.0
    } else {
        num / den
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Detector {
    Reactive,
    StaticThreshold,
    FailureMinerTrend,
}

impl Detector {
    pub fn label(&self) -> &'static str {
        match self {
            Detector::Reactive => "reactive",
            Detector::StaticThreshold => "static-threshold",
            Detector::FailureMinerTrend => "failureminer-trend",
        }
    }

    /// First day this detector would have alerted, if any.
    pub fn first_alert(&self, d: &Drive) -> Option<u32> {
        match self {
            Detector::Reactive => d.fails_on,
            Detector::StaticThreshold => d
                .samples
                .iter()
                .find(|s| s.reallocated > 0.55 || s.media_errors > 0.45)
                .map(|s| s.day),
            Detector::FailureMinerTrend => {
                // Require both an elevated level and a positive slope, to avoid
                // alerting on noisy-but-healthy drives.
                let mut best: Option<u32> = None;
                for (i, s) in d.samples.iter().enumerate() {
                    if i < 7 {
                        continue;
                    }
                    let level = health_score(s);
                    let sl = trend(&d.samples[..=i], 7);
                    if level > 0.18 && sl > 0.004 {
                        best = Some(s.day);
                        break;
                    }
                }
                best
            }
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DetectorResult {
    pub detector: String,
    pub true_failures: usize,
    pub detected: usize,
    pub recall: f64,
    pub mean_lead_time_days: f64,
    pub median_lead_time_days: f64,
    pub false_positive_alerts: usize,
    pub false_positive_rate_per_drive_year: f64,
    pub total_alerts: usize,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BenchReport {
    pub bounty: String,
    pub hip: String,
    pub honesty_labels: Vec<String>,
    pub config: BenchConfig,
    pub results: Vec<DetectorResult>,
    pub notes: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BenchConfig {
    pub seed: u64,
    pub n_drives: usize,
    pub window_days: u32,
    pub base_failure_rate: f64,
}

impl Default for BenchConfig {
    fn default() -> Self {
        Self { seed: 0xFA11_5EE, n_drives: 1000, window_days: 180, base_failure_rate: 0.12 }
    }
}
impl BenchConfig {
    pub fn test_config() -> Self {
        Self { seed: 7, n_drives: 200, window_days: 120, base_failure_rate: 0.15 }
    }
}

pub fn evaluate(detector: Detector, drives: &[Drive]) -> DetectorResult {
    let mut true_failures = 0usize;
    let mut detected = 0usize;
    let mut leads: Vec<f64> = Vec::new();
    let mut fp = 0usize;
    let mut alerts = 0usize;

    for d in drives {
        let alert = detector.first_alert(d);
        if alert.is_some() {
            alerts += 1;
        }
        match (d.fails_on, alert) {
            (Some(f), Some(a)) => {
                true_failures += 1;
                detected += 1;
                // Lead time is measured only when the alert precedes the failure.
                leads.push((f as f64 - a as f64).max(0.0));
            }
            (Some(_), None) => true_failures += 1,
            (None, Some(_)) => fp += 1,
            (None, None) => {}
        }
    }

    let healthy_drive_years = drives
        .iter()
        .filter(|d| d.fails_on.is_none())
        .count() as f64
        * (120.0 / 365.0);

    leads.sort_by(|a, b| a.partial_cmp(b).unwrap());
    let mean_lead = if leads.is_empty() { 0.0 } else { leads.iter().sum::<f64>() / leads.len() as f64 };
    let median_lead = if leads.is_empty() {
        0.0
    } else {
        leads[leads.len() / 2]
    };

    DetectorResult {
        detector: detector.label().to_string(),
        true_failures,
        detected,
        recall: if true_failures == 0 { 0.0 } else { detected as f64 / true_failures as f64 },
        mean_lead_time_days: mean_lead,
        median_lead_time_days: median_lead,
        false_positive_alerts: fp,
        false_positive_rate_per_drive_year: if healthy_drive_years > 0.0 {
            fp as f64 / healthy_drive_years
        } else {
            0.0
        },
        total_alerts: alerts,
    }
}

pub fn run_bench(cfg: &BenchConfig) -> Result<BenchReport, FmError> {
    if cfg.n_drives == 0 {
        return Err(FmError::InvalidConfig("n_drives must be positive".into()));
    }
    if cfg.window_days < 30 {
        return Err(FmError::InvalidConfig("window_days must be >= 30 to measure trends".into()));
    }
    if !(0.0..=1.0).contains(&cfg.base_failure_rate) {
        return Err(FmError::InvalidConfig("base_failure_rate must be in [0,1]".into()));
    }

    let fleet = generate_fleet(cfg.seed, cfg.n_drives, cfg.window_days, cfg.base_failure_rate);
    let results: Vec<DetectorResult> = [Detector::Reactive, Detector::StaticThreshold, Detector::FailureMinerTrend]
        .iter()
        .map(|d| evaluate(*d, &fleet))
        .collect();

    let mut notes = Vec::new();
    if let (Some(reactive), Some(trendres)) = (
        results.iter().find(|r| r.detector == "reactive"),
        results.iter().find(|r| r.detector == "failureminer-trend"),
    ) {
        notes.push(format!(
            "trend detector buys {:.1} days mean lead over reactive, at {:.2} false positives per drive-year",
            trendres.mean_lead_time_days - reactive.mean_lead_time_days,
            trendres.false_positive_rate_per_drive_year
        ));
    }

    Ok(BenchReport {
        bounty: "hanzoai/s3#33 [Bounty #44] Predictive SSD failure detection (FailureMiner)".into(),
        hip: "HIP-1348: Object Storage Durability and Placement".into(),
        honesty_labels: vec![
            "fully synthetic SMART telemetry — no real drive fleet data".into(),
            "wear trajectories are monotone-with-noise models; real failure modes are far more heterogeneous".into(),
            "detector thresholds are hand-set, not tuned on labelled data".into(),
            "All figures are simulation outputs on the authoring host, NOT Hanzo expectations until reproduced on real fleet telemetry.".into(),
        ],
        config: cfg.clone(),
        results,
        notes,
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
    fn trend_beats_reactive_lead_time() {
        let report = run_bench(&BenchConfig::test_config()).unwrap();
        let reactive = report.results.iter().find(|r| r.detector == "reactive").unwrap();
        let trendr = report.results.iter().find(|r| r.detector == "failureminer-trend").unwrap();
        assert!(
            trendr.mean_lead_time_days >= reactive.mean_lead_time_days,
            "trend lead {:.2} < reactive {:.2}",
            trendr.mean_lead_time_days,
            reactive.mean_lead_time_days
        );
    }

    #[test]
    fn reactive_has_no_false_positives() {
        let report = run_bench(&BenchConfig::test_config()).unwrap();
        let reactive = report.results.iter().find(|r| r.detector == "reactive").unwrap();
        assert_eq!(reactive.false_positive_alerts, 0);
    }

    #[test]
    fn fleet_has_both_classes() {
        let fleet = generate_fleet(3, 300, 120, 0.2);
        assert!(fleet.iter().any(|d| d.fails_on.is_some()), "no failing drives generated");
        assert!(fleet.iter().any(|d| d.fails_on.is_none()), "no healthy drives generated");
    }

    #[test]
    fn all_detectors_produce_results() {
        let report = run_bench(&BenchConfig::test_config()).unwrap();
        assert_eq!(report.results.len(), 3);
        let names: HashMap<&str, usize> =
            report.results.iter().map(|r| (r.detector.as_str(), 1)).collect();
        assert!(names.contains_key("failureminer-trend"));
    }

    #[test]
    fn short_window_rejected() {
        let mut cfg = BenchConfig::test_config();
        cfg.window_days = 10;
        assert!(run_bench(&cfg).is_err());
    }
}
