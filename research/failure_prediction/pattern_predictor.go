package failure_prediction

import "math"

// PatternPredictor is a simplified implementation of the FailureMiner concept.
// It looks for "Failure Patterns" (FPs) - specific sequences or combinations 
// of attribute changes rather than absolute thresholds.
type PatternPredictor struct {
	WindowSize int
}

func (p *PatternPredictor) Name() string {
	return "FailureMiner_Pattern_Logic"
}

func (p *PatternPredictor) Predict(history []TelemetryRecord) PredictionResult {
	if len(history) < 2 {
		return PredictionResult{WillFail: false}
	}

	// FailureMiner Logic: Detect "Attribute Volatility"
	// Instead of val > 0, we look for the rate of change (ROC)
	// and correlation between multiple attributes (e.g., SMART 5 and 197 increasing together)
	
	latest := history[len(history)-1]
	prev := history[len(history)-2]
	
	volatilityCount := 0
	for id, val := range latest.Attributes {
		if prevVal, ok := prev.Attributes[id]; ok {
			if val > prevVal {
				volatilityCount++
			}
		}
	}

	// If more than 2 critical attributes are trending upwards simultaneously,
	// trigger a failure prediction even if absolute values are low.
	if volatilityCount >= 2 {
		return PredictionResult{
			WillFail:   true,
			Confidence: math.Min(1.0, float64(volatilityCount)/5.0),
		}
	}

	return PredictionResult{WillFail: false}
}
