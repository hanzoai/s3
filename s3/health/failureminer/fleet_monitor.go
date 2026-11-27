package failureminer

import (
	"fmt"
	"sync"
)

type tenantDriveKey struct {
	tenantID string
	driveID  string
}

type monitoredStream struct {
	mu          sync.Mutex
	scorer      StreamScorer
	lastTS      int64
	sampleCount int64
	alertCount  int64
	lastPred    Prediction
}

// FleetMonitor provides thread-safe, multi-tenant isolated online scoring over
// concurrent NVMe telemetry streams. Each (tenant_id, drive_id) pair owns an
// independent StreamScorer and enforces strict ascending timestamp monotonicity
// to prevent out-of-order causal corruption.
type FleetMonitor struct {
	predictor Predictor

	mu      sync.RWMutex
	streams map[tenantDriveKey]*monitoredStream
}

// NewFleetMonitor creates a concurrent multi-tenant monitor backed by a fitted
// Predictor.
func NewFleetMonitor(predictor Predictor) *FleetMonitor {
	return &FleetMonitor{
		predictor: predictor,
		streams:   make(map[tenantDriveKey]*monitoredStream),
	}
}

// Ingest validates a sample, verifies strict per-(tenant, drive) temporal
// monotonicity, and emits a causal Prediction.
func (m *FleetMonitor) Ingest(sample DriveTelemetrySample) (Prediction, error) {
	if err := sample.Validate(); err != nil {
		return Prediction{}, err
	}
	if sample.TenantID == "" {
		return Prediction{}, fmt.Errorf("empty tenant_id")
	}

	key := tenantDriveKey{
		tenantID: sample.TenantID,
		driveID:  sample.DriveID,
	}

	ms := m.getOrCreateStream(key, sample.Meta())
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if ms.sampleCount > 0 && sample.TimestampUnix <= ms.lastTS {
		return Prediction{}, fmt.Errorf(
			"non-monotonic timestamp for tenant=%s drive=%s: got %d <= last %d",
			sample.TenantID, sample.DriveID, sample.TimestampUnix, ms.lastTS,
		)
	}

	pred := ms.scorer.Observe(sample)
	ms.lastTS = sample.TimestampUnix
	ms.sampleCount++
	if pred.Alert {
		ms.alertCount++
	}
	ms.lastPred = pred
	return pred, nil
}

// TenantSummary aggregates read-only fleet health statistics strictly scoped
// to a single tenant.
type TenantSummary struct {
	TenantID      string
	TrackedDrives int
	AlertingDrives int
	TotalSamples  int64
	TotalAlerts   int64
}

// SummaryForTenant returns telemetry statistics exclusively for tenantID.
func (m *FleetMonitor) SummaryForTenant(tenantID string) TenantSummary {
	m.mu.RLock()
	matched := make([]*monitoredStream, 0)
	for k, ms := range m.streams {
		if k.tenantID == tenantID {
			matched = append(matched, ms)
		}
	}
	m.mu.RUnlock()

	out := TenantSummary{
		TenantID:      tenantID,
		TrackedDrives: len(matched),
	}
	for _, ms := range matched {
		ms.mu.Lock()
		out.TotalSamples += ms.sampleCount
		out.TotalAlerts += ms.alertCount
		if ms.lastPred.Alert {
			out.AlertingDrives++
		}
		ms.mu.Unlock()
	}
	return out
}

func (m *FleetMonitor) getOrCreateStream(key tenantDriveKey, meta DriveMeta) *monitoredStream {
	m.mu.RLock()
	ms, ok := m.streams[key]
	m.mu.RUnlock()
	if ok {
		return ms
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if ms, ok = m.streams[key]; ok {
		return ms
	}
	ms = &monitoredStream{
		scorer: m.predictor.NewStream(meta),
	}
	m.streams[key] = ms
	return ms
}
