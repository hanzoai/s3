package failure_prediction

import (
	"testing"
	"time"
)

func TestFailurePredictionExperiment(t *testing.T) {
	// 1. Setup Synthetic Dataset (Simulating Hanzo Hardware Baseline)
	// In a real scenario, this would load from Backblaze or Hanzo's telemetry logs.
	dataset := make(map[string][]TelemetryRecord)

	// Drive A: Healthy
	dataset["drive_a"] = []TelemetryRecord{
		{Timestamp: time.Now(), Attributes: map[int]int64{5: 0, 187: 0}, Failure: 0},
	}

	// Drive B: Failed (Caught by Threshold)
	dataset["drive_b"] = []TelemetryRecord{
		{Timestamp: time.Now(), Attributes: map[int]int64{5: 10, 187: 0}, Failure: 1},
	}

	// Drive C: Failed (Caught by Pattern/FailureMiner, but missed by Threshold)
	// Threshold looks for > 0. Here we simulate a drive that fails with 0 critical errors
	// but high volatility in other attributes.
	dataset["drive_c"] = []TelemetryRecord{
		{Timestamp: time.Now().Add(-24 * time.Hour), Attributes: map[int]int64{7: 100, 12: 10}, Failure: 0},
		{Timestamp: time.Now(), Attributes: map[int]int64{7: 200, 12: 50}, Failure: 1},
	}

	// 2. Run Baseline (Threshold)
	baseline := &ThresholdPredictor{CriticalAttributes: []int{5, 187}}
	evalBaseline := Evaluator{Predictor: baseline}
	metricsBaseline := evalBaseline.Evaluate(dataset)
	PrintResults(baseline.Name(), metricsBaseline)

	// 3. Run FailureMiner Logic
	miner := &PatternPredictor{WindowSize: 2}
	evalMiner := Evaluator{Predictor: miner}
	metricsMiner := evalMiner.Evaluate(dataset)
	PrintResults(miner.Name(), metricsMiner)

	// 4. Assertions
	if metricsMiner.Recall < metricsBaseline.Recall {
		t.Errorf("FailureMiner logic should ideally have better or equal recall than baseline")
	}
}
