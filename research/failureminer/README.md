# FailureMiner Research Module

This module implements the evaluation framework for predictive SSD failure detection as per HIP-1348.

## Methodology
1. **Data Collection**: Uses `smartctl` (via syscalls/exec) to capture NVMe telemetry.
2. **Baseline**: Compares current SMART attributes against static thresholds (e.g., `Media and Data Integrity Errors`, `Available Spare`).
3. **Evaluation**: Logs metrics to `failure_metrics.json` for recall/false-alarm analysis.

## How to Run
```bash
go run research/failureminer/main.go --device=/dev/nvme0n1
