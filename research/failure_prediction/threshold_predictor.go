package failure_prediction

// ThresholdPredictor implements the baseline "Simple Threshold" method
// commonly used in industry (e.g., checking if critical SMART attributes > 0).
type ThresholdPredictor struct {
	CriticalAttributes []int
}

func NewThresholdPredictor() *ThresholdPredictor {
	return &ThresholdPredictor{
		// 5: Reallocated_Sector_Ct
		// 187: Reported_Uncorrectable_Errors
		// 197: Current_Pending_Sector_Ct
		// 198: Offline_Uncorrectable
		CriticalAttributes: []int{5, 187, 197, 198},
	}
}

func (p *ThresholdPredictor) Name() string {
	return "Static_Threshold_Baseline"
}

func (p *ThresholdPredictor) Predict(history []TelemetryRecord) PredictionResult {
	if len(history) == 0 {
		return PredictionResult{WillFail: false}
	}

	// Check the latest record
	latest := history[len(history)-1]
	for _, attrID := range p.CriticalAttributes {
		if val, ok := latest.Attributes[attrID]; ok && val > 0 {
			return PredictionResult{
				WillFail:   true,
				Confidence: 1.0,
			}
		}
	}

	return PredictionResult{WillFail: false}
}
