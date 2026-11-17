package storage

import (
	"errors"
	"sync"
	"time"
)

// ExtentReclamationDecision specifies the optimal space reclamation strategy.
type ExtentReclamationDecision int

const (
	// DecisionNoOp indicates volume does not meet reclamation threshold.
	DecisionNoOp ExtentReclamationDecision = iota
	// DecisionDirectDiscard indicates direct hole punching / extent discard without copying.
	DecisionDirectDiscard
	// DecisionCompaction indicates standard live-data copy compaction.
	DecisionCompaction
)

func (d ExtentReclamationDecision) String() string {
	switch d {
	case DecisionNoOp:
		return "NOOP"
	case DecisionDirectDiscard:
		return "DIRECT_DISCARD"
	case DecisionCompaction:
		return "COPY_COMPACTION"
	default:
		return "UNKNOWN"
	}
}

// DisCoGCConfig defines operational thresholds for the DisCoGC storage economics engine.
type DisCoGCConfig struct {
	MinInvalidRatio        float64 // Minimum dead space ratio to trigger reclamation (e.g. 0.30)
	DirectDiscardMinRatio  float64 // Ratio above which direct discard is preferred (e.g. 0.50)
	MinContiguousRunBytes  int64   // Minimum run length to justify kernel hole punch (e.g. 2MB)
	MaxFragmentationRatio  float64 // Upper bound of extent fragmentation before forcing compaction (e.g. 0.25)
	FlashCostPerTB         float64 // Cost per TB of enterprise NVMe (e.g. $120.00)
	DriveRatedTBW          float64 // Drive rated petabytes/terabytes written
}

// DefaultDisCoGCConfig returns production defaults aligned with HIP-1348.
func DefaultDisCoGCConfig() DisCoGCConfig {
	return DisCoGCConfig{
		MinInvalidRatio:       0.30,
		DirectDiscardMinRatio: 0.50,
		MinContiguousRunBytes: 2 * 1024 * 1024, // 2MB
		MaxFragmentationRatio: 0.25,
		FlashCostPerTB:        120.0,
		DriveRatedTBW:         3000.0,
	}
}

// ExtentRun represents a contiguous sequence of invalid/deleted bytes.
type ExtentRun struct {
	Offset int64
	Length int64
}

// VolumeEconomicsMetrics captures storage economics and write amplification analysis.
type VolumeEconomicsMetrics struct {
	TotalVolumeBytes       int64   `json:"total_volume_bytes"`
	LiveBytes              int64   `json:"live_bytes"`
	InvalidBytes           int64   `json:"invalid_bytes"`
	InvalidRatio           float64 `json:"invalid_ratio"`
	FragmentationRatio     float64 `json:"fragmentation_ratio"`
	CompactionWAF          float64 `json:"compaction_waf"`
	DirectDiscardWAF       float64 `json:"direct_discard_waf"`
	WAFReductionPct        float64 `json:"waf_reduction_pct"`
	ProjectedTbwSavingsGB  float64 `json:"projected_tbw_savings_gb"`
	ProjectedCostSavingsPct float64 `json:"projected_cost_savings_pct"`
}

// ReclamationPlan describes the execution plan for reclaiming volume space.
type ReclamationPlan struct {
	VolumeID      uint32
	Decision      ExtentReclamationDecision
	Metrics       VolumeEconomicsMetrics
	DiscardRuns   []ExtentRun
	EvaluatedAt   time.Time
	Reason        string
}

// DisCoGCPolicy evaluates storage volumes according to FAST '26 DisCoGC economics.
type DisCoGCPolicy struct {
	mu     sync.RWMutex
	config DisCoGCConfig
}

// NewDisCoGCPolicy constructs a new thread-safe policy instance.
func NewDisCoGCPolicy(cfg DisCoGCConfig) *DisCoGCPolicy {
	return &DisCoGCPolicy{
		config: cfg,
	}
}

// EvaluateVolumeEconomics calculates write amplification and selects the optimal reclamation mode.
func (p *DisCoGCPolicy) EvaluateVolumeEconomics(
	volumeID uint32,
	totalBytes int64,
	invalidBytes int64,
	runs []ExtentRun,
) (*ReclamationPlan, error) {
	if totalBytes <= 0 {
		return nil, errors.New("volume totalBytes must be positive")
	}

	p.mu.RLock()
	cfg := p.config
	p.mu.RUnlock()

	liveBytes := totalBytes - invalidBytes
	if liveBytes < 0 {
		liveBytes = 0
	}

	invalidRatio := float64(invalidBytes) / float64(totalBytes)

	// Calculate fragmentation: ratio of small discontinuous runs vs total runs
	var smallRunsCount int
	var qualifiedRuns []ExtentRun
	for _, run := range runs {
		if run.Length < cfg.MinContiguousRunBytes {
			smallRunsCount++
		} else {
			qualifiedRuns = append(qualifiedRuns, run)
		}
	}

	fragmentationRatio := 0.0
	if len(runs) > 0 {
		fragmentationRatio = float64(smallRunsCount) / float64(len(runs))
	}

	// Traditional Compaction WAF = 1 + (Live / Total)
	compactionWAF := 1.0 + (float64(liveBytes) / float64(totalBytes))
	// Direct Discard WAF = 1.00 + metadata update overhead (~0.01)
	directDiscardWAF := 1.01

	wafReductionPct := 0.0
	if compactionWAF > 0 {
		wafReductionPct = ((compactionWAF - directDiscardWAF) / compactionWAF) * 100.0
	}

	projectedTbwSavingsGB := float64(liveBytes) / (1024 * 1024 * 1024)
	projectedCostSavingsPct := wafReductionPct * 0.72 // Proportional amortized flash wear saving

	metrics := VolumeEconomicsMetrics{
		TotalVolumeBytes:        totalBytes,
		LiveBytes:               liveBytes,
		InvalidBytes:            invalidBytes,
		InvalidRatio:            invalidRatio,
		FragmentationRatio:      fragmentationRatio,
		CompactionWAF:           compactionWAF,
		DirectDiscardWAF:        directDiscardWAF,
		WAFReductionPct:         wafReductionPct,
		ProjectedTbwSavingsGB:   projectedTbwSavingsGB,
		ProjectedCostSavingsPct: projectedCostSavingsPct,
	}

	plan := &ReclamationPlan{
		VolumeID:    volumeID,
		Metrics:     metrics,
		EvaluatedAt: time.Now(),
	}

	// Decision Engine Logic
	if invalidRatio < cfg.MinInvalidRatio {
		plan.Decision = DecisionNoOp
		plan.Reason = "Invalid space ratio below trigger threshold"
		return plan, nil
	}

	// High fragmentation guardrail: fine-grained fragments force copy compaction
	if fragmentationRatio > cfg.MaxFragmentationRatio {
		plan.Decision = DecisionCompaction
		plan.Reason = "Excessive extent fragmentation; copy compaction selected to defragment storage layout"
		return plan, nil
	}

	// Coarse runs and sufficient dead ratio qualify for Direct-Discard
	if invalidRatio >= cfg.DirectDiscardMinRatio || len(qualifiedRuns) > 0 {
		plan.Decision = DecisionDirectDiscard
		plan.DiscardRuns = qualifiedRuns
		plan.Reason = "Coarse dead extent runs qualify for direct hole punch; avoids live rewrite"
		return plan, nil
	}

	plan.Decision = DecisionCompaction
	plan.Reason = "Moderate invalid ratio with sub-threshold run sizes; defaulting to standard compaction"
	return plan, nil
}
