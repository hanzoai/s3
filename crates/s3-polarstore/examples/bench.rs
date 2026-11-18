// SPDX-License-Identifier: Apache-2.0
//! Reproducible PolarStore benchmark runner.
use s3_polarstore::{BenchConfig, run_bench, write_results_json};
use std::path::PathBuf;

fn main() {
    let out = std::env::args().nth(1).map(PathBuf::from).unwrap_or_else(|| PathBuf::from("results/polarstore-results.json"));
    let cfg = BenchConfig::default();
    eprintln!("polarstore bench: seed={} objects={}", cfg.seed, cfg.n_objects);
    let report = run_bench(&cfg);
    write_results_json(&report, &out).expect("failed to write results");
    eprintln!("wrote {}", out.display());
    for run in &report.runs {
        println!("{}", serde_json::to_string(run).expect("serialize"));
    }
}
