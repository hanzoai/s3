package failure_prediction

import "time"

// TelemetryRecord represents a single SMART/NVMe snapshot
type TelemetryRecord struct {
	Timestamp  time.Time
	Serial     string
	Model      string
	Capacity   int64
	Failure    int // 1 if failed on this day, 0 otherwise
	Attributes map[int]int64 // SMART ID -> Value
}

// PredictionResult contains the outcome of a model's assessment
type PredictionResult struct {
	WillFail  bool
	LeadTime  time.Duration
	Confidence float64
}

// Predictor defines the interface for failure detection methods
type Predictor interface {
	Name() string
	Predict(history []TelemetryRecord) PredictionResult
}
