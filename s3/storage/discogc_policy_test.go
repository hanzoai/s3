package storage

import (
	"testing"
)

func TestDisCoGCPolicy_HighInvalidRatio_DirectDiscard(t *testing.T) {
	cfg := DefaultDisCoGCConfig()
	policy := NewDisCoGCPolicy(cfg)

	// 100GB Volume with 60GB dead space (invalid ratio = 0.60)
	totalBytes := int64(100 * 1024 * 1024 * 1024)
	invalidBytes := int64(60 * 1024 * 1024 * 1024)

	runs := []ExtentRun{
		{Offset: 0, Length: 20 * 1024 * 1024 * 1024},
		{Offset: 40 * 1024 * 1024 * 1024, Length: 40 * 1024 * 1024 * 1024},
	}

	plan, err := policy.EvaluateVolumeEconomics(101, totalBytes, invalidBytes, runs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if plan.Decision != DecisionDirectDiscard {
		t.Errorf("expected DecisionDirectDiscard, got %v (%s)", plan.Decision, plan.Reason)
	}

	if plan.Metrics.CompactionWAF < 1.35 {
		t.Errorf("expected CompactionWAF ~1.40, got %f", plan.Metrics.CompactionWAF)
	}

	if plan.Metrics.DirectDiscardWAF > 1.05 {
		t.Errorf("expected DirectDiscardWAF ~1.01, got %f", plan.Metrics.DirectDiscardWAF)
	}

	if plan.Metrics.WAFReductionPct <= 0 {
		t.Errorf("expected positive WAF reduction percentage, got %f", plan.Metrics.WAFReductionPct)
	}
}

func TestDisCoGCPolicy_HighFragmentation_FallsBackToCompaction(t *testing.T) {
	cfg := DefaultDisCoGCConfig()
	policy := NewDisCoGCPolicy(cfg)

	totalBytes := int64(50 * 1024 * 1024 * 1024)
	invalidBytes := int64(25 * 1024 * 1024 * 1024) // 50% invalid

	// Many tiny fragments under 2MB min run length
	var runs []ExtentRun
	for i := 0; i < 20; i++ {
		runs = append(runs, ExtentRun{
			Offset: int64(i * 100 * 1024),
			Length: 64 * 1024, // 64KB tiny runs
		})
	}

	plan, err := policy.EvaluateVolumeEconomics(102, totalBytes, invalidBytes, runs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if plan.Decision != DecisionCompaction {
		t.Errorf("expected DecisionCompaction due to fragmentation, got %v (%s)", plan.Decision, plan.Reason)
	}
}

func TestDisCoGCPolicy_BelowThreshold_NoOp(t *testing.T) {
	cfg := DefaultDisCoGCConfig()
	policy := NewDisCoGCPolicy(cfg)

	totalBytes := int64(100 * 1024 * 1024 * 1024)
	invalidBytes := int64(10 * 1024 * 1024 * 1024) // 10% invalid (below 30% trigger)

	plan, err := policy.EvaluateVolumeEconomics(103, totalBytes, invalidBytes, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if plan.Decision != DecisionNoOp {
		t.Errorf("expected DecisionNoOp for 10%% invalid ratio, got %v (%s)", plan.Decision, plan.Reason)
	}
}
