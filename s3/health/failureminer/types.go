// Package failureminer provides a clean-room Go implementation and evaluation
// harness for predictive SSD failure detection inspired by FailureMiner
// (Wang et al., USENIX FAST '26, https://www.usenix.org/conference/fast26/presentation/wang-shuyang).
//
// All predictors implement an online, strictly causal streaming interface:
// predictions emitted at timestamp t depend exclusively on telemetry samples
// observed at timestamps <= t. Zero external dependencies beyond the Go
// standard library are used.
package failureminer

import (
	"fmt"
	"math"
	"sort"
)

// FailureMode classifies a drive's ground-truth lifecycle trajectory in the
// evaluation harness.
type FailureMode string

const (
	// FailureModeNone represents a standard healthy drive that does not fail
	// within the observation horizon.
	FailureModeNone FailureMode = "healthy"

	// FailureModeHardNegativeHot represents a healthy, high-utilization hot-tier
	// drive exhibiting sustained elevated temperature, transient latency spikes,
	// or static non-progressing single errors that does NOT fail.
	FailureModeHardNegativeHot FailureMode = "hard_negative_hot_tier"

	// FailureModeMediaErrorRamp represents progressive media error accumulation
	// accompanied by error-log growth, reallocated sectors, and read latency drift.
	FailureModeMediaErrorRamp FailureMode = "progressive_media_error_ramp"

	// FailureModeThermalSpareDepletion represents sustained thermal stress paired
	// with accelerating available-spare depletion and write-latency degradation.
	FailureModeThermalSpareDepletion FailureMode = "thermal_spare_depletion"

	// FailureModeCRCTimeoutBurst represents intermittent escalating bursts of CRC
	// errors, I/O timeouts, and tail-latency spikes prior to failure.
	FailureModeCRCTimeoutBurst FailureMode = "intermittent_crc_timeout_burst"

	// FailureModeSuddenControllerDeath represents an abrupt hardware/ASIC/power
	// failure with zero pre-failure SMART/NVMe precursor signals.
	FailureModeSuddenControllerDeath FailureMode = "sudden_unannounced_controller_death"
)

// IsFailure reports whether the failure mode corresponds to a drive that fails
// within the horizon.
func (m FailureMode) IsFailure() bool {
	return m != FailureModeNone && m != FailureModeHardNegativeHot
}

// IsPredictable reports whether the failure mode produces physical pre-failure
// telemetry precursors before failure.
func (m FailureMode) IsPredictable() bool {
	return m.IsFailure() && m != FailureModeSuddenControllerDeath
}

// DriveMeta identifies static hardware and tenancy metadata for a drive.
type DriveMeta struct {
	DriveID  string `json:"drive_id"`
	TenantID string `json:"tenant_id"`
	Vendor   string `json:"vendor"`
	Model    string `json:"model"`
}

// DriveTelemetrySample represents a single point-in-time NVMe Log Page 02h /
// SMART telemetry observation for an SSD.
type DriveTelemetrySample struct {
	TimestampUnix int64  `json:"timestamp_unix"`
	DriveID       string `json:"drive_id"`
	TenantID      string `json:"tenant_id"`
	Vendor        string `json:"vendor"`
	Model         string `json:"model"`

	// NVMe Log Page 02h / SMART health counters.
	MediaErrors         uint64  `json:"media_errors"`
	NumErrLogEntries    uint64  `json:"num_err_log_entries"`
	AvailableSparePct   float64 `json:"available_spare_pct"`
	PercentageUsed      float64 `json:"percentage_used"`
	CriticalWarningMask uint8   `json:"critical_warning_mask"`
	TemperatureC        float64 `json:"temperature_c"`
	UnsafeShutdowns     uint64  `json:"unsafe_shutdowns"`
	ControllerBusyTime  uint64  `json:"controller_busy_time"`
	IOTimeoutCount      uint64  `json:"io_timeout_count"`
	P99ReadLatencyUs    float64 `json:"p99_read_latency_us"`
	P99WriteLatencyUs   float64 `json:"p99_write_latency_us"`
	CRCErrors           uint64  `json:"crc_errors"`
	ReallocatedSectors  uint64  `json:"reallocated_sectors"`
	PowerOnHours        uint64  `json:"power_on_hours"`
}

// Meta extracts static drive metadata from a sample.
func (s DriveTelemetrySample) Meta() DriveMeta {
	return DriveMeta{
		DriveID:  s.DriveID,
		TenantID: s.TenantID,
		Vendor:   s.Vendor,
		Model:    s.Model,
	}
}

// Validate checks basic invariants and rejects NaN/Inf or out-of-bounds values.
func (s DriveTelemetrySample) Validate() error {
	if s.DriveID == "" {
		return fmt.Errorf("empty drive_id")
	}
	if s.TimestampUnix <= 0 {
		return fmt.Errorf("non-positive timestamp_unix: %d", s.TimestampUnix)
	}
	for name, v := range map[string]float64{
		"available_spare_pct":  s.AvailableSparePct,
		"percentage_used":      s.PercentageUsed,
		"temperature_c":        s.TemperatureC,
		"p99_read_latency_us":  s.P99ReadLatencyUs,
		"p99_write_latency_us": s.P99WriteLatencyUs,
	} {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return fmt.Errorf("invalid non-finite %s", name)
		}
	}
	if s.AvailableSparePct < 0 || s.AvailableSparePct > 100 {
		return fmt.Errorf("available_spare_pct out of range [0,100]: %v", s.AvailableSparePct)
	}
	if s.TemperatureC < -40 || s.TemperatureC > 150 {
		return fmt.Errorf("temperature_c out of physical range [-40,150]: %v", s.TemperatureC)
	}
	if s.P99ReadLatencyUs < 0 || s.P99WriteLatencyUs < 0 {
		return fmt.Errorf("negative latency sample")
	}
	return nil
}

// DriveSeries holds a chronologically ordered telemetry sequence for one drive
// along with its ground-truth lifecycle outcome.
type DriveSeries struct {
	Meta            DriveMeta              `json:"meta"`
	Failed          bool                   `json:"failed"`
	FailureMode     FailureMode            `json:"failure_mode"`
	FailureTimeUnix int64                  `json:"failure_time_unix"`
	Samples         []DriveTelemetrySample `json:"samples"`
}

// Clone returns a deep copy of a DriveSeries so perturbation tests never mutate
// shared slices.
func (d DriveSeries) Clone() DriveSeries {
	cp := d
	cp.Samples = make([]DriveTelemetrySample, len(d.Samples))
	copy(cp.Samples, d.Samples)
	return cp
}

// FleetDataset is a collection of labeled drive telemetry histories.
type FleetDataset struct {
	Name        string        `json:"name"`
	Seed        int64         `json:"seed"`
	StepHours   int           `json:"step_hours"`
	HorizonDays int           `json:"horizon_days"`
	Drives      []DriveSeries `json:"drives"`
}

// SplitTrainTest performs a deterministic, failure-mode-stratified split by
// drive ID so zero drives overlap between train and test splits and failure
// prevalence is preserved across splits.
func (f FleetDataset) SplitTrainTest(trainFraction float64) (train []DriveSeries, test []DriveSeries) {
	if trainFraction <= 0 {
		trainFraction = 0.5
	}
	if trainFraction >= 1 {
		trainFraction = 0.8
	}
	byMode := make(map[FailureMode][]DriveSeries)
	modes := make([]FailureMode, 0, 6)
	for _, d := range f.Drives {
		if _, ok := byMode[d.FailureMode]; !ok {
			modes = append(modes, d.FailureMode)
		}
		byMode[d.FailureMode] = append(byMode[d.FailureMode], d)
	}
	sort.Slice(modes, func(i, j int) bool { return modes[i] < modes[j] })
	for _, m := range modes {
		group := byMode[m]
		sort.Slice(group, func(i, j int) bool { return group[i].Meta.DriveID < group[j].Meta.DriveID })
		cut := int(math.Round(float64(len(group)) * trainFraction))
		if len(group) >= 2 {
			if cut < 1 {
				cut = 1
			}
			if cut >= len(group) {
				cut = len(group) - 1
			}
		}
		for i, d := range group {
			if i < cut {
				train = append(train, d)
			} else {
				test = append(test, d)
			}
		}
	}
	sort.Slice(train, func(i, j int) bool { return train[i].Meta.DriveID < train[j].Meta.DriveID })
	sort.Slice(test, func(i, j int) bool { return test[i].Meta.DriveID < test[j].Meta.DriveID })
	return train, test
}

// Prediction is the causal output produced for a single drive at timestamp t.
type Prediction struct {
	DriveID               string   `json:"drive_id"`
	TimestampUnix         int64    `json:"timestamp_unix"`
	RiskScore             float64  `json:"risk_score"`
	Alert                 bool     `json:"alert"`
	LeadTimeEstimateHours float64  `json:"lead_time_estimate_hours"`
	ActiveMotifs          []string `json:"active_motifs,omitempty"`
}

// StreamScorer processes a single drive's telemetry strictly in ascending
// timestamp order. Each Observe call may only inspect state derived from
// samples at timestamps <= sample.TimestampUnix.
type StreamScorer interface {
	Observe(sample DriveTelemetrySample) Prediction
}

// Predictor represents a trainable failure-prediction algorithm that constructs
// independent per-drive causal stream scorers.
type Predictor interface {
	Name() string
	Fit(train []DriveSeries)
	NewStream(meta DriveMeta) StreamScorer
}
