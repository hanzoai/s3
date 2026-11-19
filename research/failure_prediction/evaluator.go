package failure_prediction

import (
	"fmt"
)

type Metrics struct {
	Recall         float64
	FalseAlarmRate float64
	Precision      float64
}

type Evaluator struct {
	Predictor Predictor
}

func (e *Evaluator) Evaluate(datasets map[string][]TelemetryRecord) Metrics {
	tp, fp, tn, fn := 0, 0, 0, 0

	for _, history := range datasets {
		actualFailure := false
		for _, rec := range history {
			if rec.Failure == 1 {
				actualFailure = true
				break
			}
		}

		prediction := e.Predictor.Predict(history)

		if prediction.WillFail && actualFailure {
			tp++
		} else if prediction.WillFail && !actualFailure {
			fp++
		} else if !prediction.WillFail && actualFailure {
			fn++
		} else {
			tn++
		}
	}

	return Metrics{
		Recall:         float64(tp) / float64(tp+fn+1e-9),
		FalseAlarmRate: float64(fp) / float64(fp+tn+1e-9),
		Precision:      float64(tp) / float64(tp+fp+1e-9),
	}
}

func PrintResults(name string, m Metrics) {
	fmt.Printf("--- Results for %s ---\n", name)
	fmt.Printf("Recall:         %.2f%%\n", m.Recall*100)
	fmt.Printf("False Alarm Rate: %.2f%%\n", m.FalseAlarmRate*100)
	fmt.Printf("Precision:      %.2f%%\n\n", m.Precision*100)
}
