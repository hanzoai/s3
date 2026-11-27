package failureminer

import (
	"math"
	"sort"
)

// ============================================================================
// Method 1: static_smart_threshold
// ============================================================================

// StaticSMARTThresholdPredictor implements the standard vendor SMART/NVMe
// static threshold baseline: alert if critical_warning > 0 OR
// available_spare < 10% OR media_errors > 0.
type StaticSMARTThresholdPredictor struct {
	SpareMinPct float64
}

// NewStaticSMARTThresholdPredictor constructs the static SMART baseline.
func NewStaticSMARTThresholdPredictor() *StaticSMARTThresholdPredictor {
	return &StaticSMARTThresholdPredictor{
		SpareMinPct: 10.0,
	}
}

func (p *StaticSMARTThresholdPredictor) Name() string { return "static_smart_threshold" }

func (p *StaticSMARTThresholdPredictor) Fit(_ []DriveSeries) {}

func (p *StaticSMARTThresholdPredictor) NewStream(_ DriveMeta) StreamScorer {
	return &staticSMARTStream{spareMinPct: p.SpareMinPct}
}

type staticSMARTStream struct {
	spareMinPct float64
}

func (s *staticSMARTStream) Observe(sample DriveTelemetrySample) Prediction {
	var motifs []string
	score := 0.0
	if sample.CriticalWarningMask > 0 {
		motifs = append(motifs, "nvme_critical_warning")
		score = 1.0
	}
	if sample.AvailableSparePct < s.spareMinPct {
		motifs = append(motifs, "available_spare_below_10pct")
		if score < 0.92 {
			score = 0.92
		}
	}
	if sample.MediaErrors > 0 {
		motifs = append(motifs, "media_errors_gt_zero")
		if score < 0.78 {
			score = 0.78
		}
	}
	alert := len(motifs) > 0
	leadEst := 0.0
	if alert {
		leadEst = 24.0
	}
	return Prediction{
		DriveID:               sample.DriveID,
		TimestampUnix:         sample.TimestampUnix,
		RiskScore:             score,
		Alert:                 alert,
		LeadTimeEstimateHours: leadEst,
		ActiveMotifs:          motifs,
	}
}

// ============================================================================
// Method 2: sliding_window_zscore
// ============================================================================

// SlidingWindowZScoreConfig configures the univariate causal z-score detector.
type SlidingWindowZScoreConfig struct {
	WindowSamples  int
	MinHistory     int
	ZThreshold     float64
	MinConsecutive int
	TargetDriveFPR float64
}

// DefaultSlidingWindowZScoreConfig returns standard defaults for the univariate
// sliding-window z-score baseline.
func DefaultSlidingWindowZScoreConfig() SlidingWindowZScoreConfig {
	return SlidingWindowZScoreConfig{
		WindowSamples:  24,
		MinHistory:     8,
		ZThreshold:     4.5,
		MinConsecutive: 2,
		TargetDriveFPR: 0.03,
	}
}

// SlidingWindowZScorePredictor is a univariate anomaly detector operating over
// each drive's rolling causal window of latency and media/error-log rates.
type SlidingWindowZScorePredictor struct {
	cfg       SlidingWindowZScoreConfig
	threshold float64
}

// NewSlidingWindowZScorePredictor constructs the sliding-window z-score baseline.
func NewSlidingWindowZScorePredictor(cfg SlidingWindowZScoreConfig) *SlidingWindowZScorePredictor {
	if cfg.WindowSamples <= 2 {
		cfg.WindowSamples = 24
	}
	if cfg.MinHistory <= 2 {
		cfg.MinHistory = 8
	}
	if cfg.ZThreshold <= 0 {
		cfg.ZThreshold = 4.5
	}
	if cfg.MinConsecutive <= 0 {
		cfg.MinConsecutive = 2
	}
	return &SlidingWindowZScorePredictor{
		cfg:       cfg,
		threshold: cfg.ZThreshold,
	}
}

func (p *SlidingWindowZScorePredictor) Name() string { return "sliding_window_zscore" }

func (p *SlidingWindowZScorePredictor) Fit(train []DriveSeries) {
	if len(train) == 0 || p.cfg.TargetDriveFPR <= 0 {
		p.threshold = p.cfg.ZThreshold
		return
	}
	// Calibrate univariate z-score threshold on healthy training drives using
	// 2-step consecutive peak z-score to respect MinConsecutive persistence.
	peaks := make([]float64, 0, len(train))
	for _, d := range train {
		if d.Failed {
			continue
		}
		st := &slidingZStream{
			cfg:       p.cfg,
			threshold: 1e9,
		}
		prevZ := 0.0
		maxPersistZ := 0.0
		for _, s := range d.Samples {
			z := st.observeRawZ(s)
			pairMin := math.Min(prevZ, z)
			if p.cfg.MinConsecutive <= 1 {
				pairMin = z
			}
			if pairMin > maxPersistZ {
				maxPersistZ = pairMin
			}
			prevZ = z
		}
		peaks = append(peaks, maxPersistZ)
	}
	if len(peaks) == 0 {
		p.threshold = p.cfg.ZThreshold
		return
	}
	q := quantile(peaks, 1.0-p.cfg.TargetDriveFPR)
	p.threshold = clamp(q, 3.8, 8.5)
}

func (p *SlidingWindowZScorePredictor) NewStream(_ DriveMeta) StreamScorer {
	return &slidingZStream{
		cfg:       p.cfg,
		threshold: p.threshold,
	}
}

type slidingZStream struct {
	cfg        SlidingWindowZScoreConfig
	threshold  float64
	readHist   []float64
	writeHist  []float64
	mediaHist  []float64
	errLogHist []float64
	hasPrev    bool
	prevMedia  uint64
	prevErrLog uint64
	consec     int
}

func (s *slidingZStream) observeRawZ(sample DriveTelemetrySample) float64 {
	var dMedia, dErrLog float64
	if s.hasPrev {
		dMedia = float64(safeSubUint64(sample.MediaErrors, s.prevMedia))
		dErrLog = float64(safeSubUint64(sample.NumErrLogEntries, s.prevErrLog))
	}
	s.hasPrev = true
	s.prevMedia = sample.MediaErrors
	s.prevErrLog = sample.NumErrLogEntries

	maxZ := 0.0
	if len(s.readHist) >= s.cfg.MinHistory {
		zRead := zScoreAgainstHistory(sample.P99ReadLatencyUs, s.readHist, 30.0)
		zWrite := zScoreAgainstHistory(sample.P99WriteLatencyUs, s.writeHist, 40.0)
		zMedia := zScoreAgainstHistory(dMedia, s.mediaHist, 0.35)
		zErrLog := zScoreAgainstHistory(dErrLog, s.errLogHist, 0.45)
		maxZ = maxFloat(maxFloat(zRead, zWrite), maxFloat(zMedia, zErrLog))
	}
	if sample.CriticalWarningMask > 0 || sample.AvailableSparePct < 10.0 {
		maxZ = math.Max(maxZ, s.threshold+2.0)
	}

	// Append current sample strictly AFTER computing its z-score against [t-W, t-1].
	s.readHist = appendCapped(s.readHist, sample.P99ReadLatencyUs, s.cfg.WindowSamples)
	s.writeHist = appendCapped(s.writeHist, sample.P99WriteLatencyUs, s.cfg.WindowSamples)
	s.mediaHist = appendCapped(s.mediaHist, dMedia, s.cfg.WindowSamples)
	s.errLogHist = appendCapped(s.errLogHist, dErrLog, s.cfg.WindowSamples)

	return maxZ
}

func (s *slidingZStream) Observe(sample DriveTelemetrySample) Prediction {
	z := s.observeRawZ(sample)
	if z >= s.threshold {
		s.consec++
	} else {
		s.consec = 0
	}
	alert := s.consec >= s.cfg.MinConsecutive
	risk := clamp(z/(s.threshold*1.35), 0.0, 1.0)
	if alert && risk < 0.70 {
		risk = 0.70
	}
	var motifs []string
	if alert {
		motifs = []string{"univariate_zscore_exceeded"}
	}
	leadEst := 0.0
	if alert {
		leadEst = 48.0
	}
	return Prediction{
		DriveID:               sample.DriveID,
		TimestampUnix:         sample.TimestampUnix,
		RiskScore:             math.Round(risk*10000) / 10000,
		Alert:                 alert,
		LeadTimeEstimateHours: leadEst,
		ActiveMotifs:          motifs,
	}
}

// ============================================================================
// Method 3: ewma_hazard_score
// ============================================================================

// EWMAHazardConfig configures the exponentially weighted trend/derivative
// hazard model.
type EWMAHazardConfig struct {
	Alpha          float64
	TargetDriveFPR float64
}

// DefaultEWMAHazardConfig returns standard defaults for the EWMA hazard baseline.
func DefaultEWMAHazardConfig() EWMAHazardConfig {
	return EWMAHazardConfig{
		Alpha:          0.28,
		TargetDriveFPR: 0.015,
	}
}

// EWMAHazardPredictor computes an exponentially weighted additive hazard score
// over SMART counter derivatives and global-normalized thermal/latency signals.
type EWMAHazardPredictor struct {
	cfg        EWMAHazardConfig
	meanTemp   float64
	stdTemp    float64
	meanRead   float64
	stdRead    float64
	meanWrite  float64
	stdWrite   float64
	threshold  float64
}

// NewEWMAHazardPredictor constructs the EWMA hazard baseline.
func NewEWMAHazardPredictor(cfg EWMAHazardConfig) *EWMAHazardPredictor {
	if cfg.Alpha <= 0 || cfg.Alpha >= 1 {
		cfg.Alpha = 0.28
	}
	if cfg.TargetDriveFPR <= 0 {
		cfg.TargetDriveFPR = 0.015
	}
	return &EWMAHazardPredictor{
		cfg:       cfg,
		meanTemp:  45.0,
		stdTemp:   7.5,
		meanRead:  400.0,
		stdRead:   180.0,
		meanWrite: 550.0,
		stdWrite:  260.0,
		threshold: 0.55,
	}
}

func (p *EWMAHazardPredictor) Name() string { return "ewma_hazard_score" }

func (p *EWMAHazardPredictor) Fit(train []DriveSeries) {
	if len(train) == 0 {
		return
	}
	var temps, reads, writes []float64
	for _, d := range train {
		if d.Failed {
			continue
		}
		for _, s := range d.Samples {
			temps = append(temps, s.TemperatureC)
			reads = append(reads, s.P99ReadLatencyUs)
			writes = append(writes, s.P99WriteLatencyUs)
		}
	}
	if len(temps) > 0 {
		p.meanTemp, p.stdTemp = meanStd(temps, 3.0)
		p.meanRead, p.stdRead = meanStd(reads, 50.0)
		p.meanWrite, p.stdWrite = meanStd(writes, 80.0)
	}

	peaks := make([]float64, 0, len(train))
	for _, d := range train {
		if d.Failed {
			continue
		}
		st := p.newStreamWithThreshold(1e9)
		peak := 0.0
		for _, s := range d.Samples {
			pred := st.Observe(s)
			if pred.RiskScore > peak {
				peak = pred.RiskScore
			}
		}
		peaks = append(peaks, peak)
	}
	if len(peaks) > 0 {
		q := quantile(peaks, 1.0-p.cfg.TargetDriveFPR)
		p.threshold = clamp(q+0.015, 0.30, 0.92)
	}
}

func (p *EWMAHazardPredictor) NewStream(_ DriveMeta) StreamScorer {
	return p.newStreamWithThreshold(p.threshold)
}

func (p *EWMAHazardPredictor) newStreamWithThreshold(thresh float64) *ewmaHazardStream {
	return &ewmaHazardStream{
		alpha:     p.cfg.Alpha,
		meanTemp:  p.meanTemp,
		stdTemp:   p.stdTemp,
		meanRead:  p.meanRead,
		stdRead:   p.stdRead,
		meanWrite: p.meanWrite,
		stdWrite:  p.stdWrite,
		threshold: thresh,
	}
}

type ewmaHazardStream struct {
	alpha     float64
	meanTemp  float64
	stdTemp   float64
	meanRead  float64
	stdRead   float64
	meanWrite float64
	stdWrite  float64
	threshold float64

	hasPrev      bool
	prevSample   DriveTelemetrySample
	ewmaMedia    float64
	ewmaErrLog   float64
	ewmaRealloc  float64
	ewmaCRCTime  float64
	ewmaSpareDep float64
	ewmaTempZ    float64
	ewmaLatZ     float64
}

func (s *ewmaHazardStream) Observe(sample DriveTelemetrySample) Prediction {
	var dMedia, dErrLog, dRealloc, dCRCTime, dSpareDep float64
	if s.hasPrev {
		dMedia = float64(safeSubUint64(sample.MediaErrors, s.prevSample.MediaErrors))
		dErrLog = float64(safeSubUint64(sample.NumErrLogEntries, s.prevSample.NumErrLogEntries))
		dRealloc = float64(safeSubUint64(sample.ReallocatedSectors, s.prevSample.ReallocatedSectors))
		dCRCTime = float64(safeSubUint64(sample.CRCErrors, s.prevSample.CRCErrors) +
			safeSubUint64(sample.IOTimeoutCount, s.prevSample.IOTimeoutCount))
		if s.prevSample.AvailableSparePct > sample.AvailableSparePct {
			dSpareDep = s.prevSample.AvailableSparePct - sample.AvailableSparePct
		}
	}
	s.hasPrev = true
	s.prevSample = sample

	tempZ := math.Max(0, (sample.TemperatureC-s.meanTemp)/s.stdTemp)
	readZ := math.Max(0, (sample.P99ReadLatencyUs-s.meanRead)/s.stdRead)
	writeZ := math.Max(0, (sample.P99WriteLatencyUs-s.meanWrite)/s.stdWrite)
	latZ := 0.5 * (readZ + writeZ)

	a := s.alpha
	s.ewmaMedia = (1-a)*s.ewmaMedia + a*math.Log1p(dMedia)
	s.ewmaErrLog = (1-a)*s.ewmaErrLog + a*math.Log1p(dErrLog)
	s.ewmaRealloc = (1-a)*s.ewmaRealloc + a*math.Log1p(dRealloc)
	s.ewmaCRCTime = (1-a)*s.ewmaCRCTime + a*math.Log1p(dCRCTime)
	s.ewmaSpareDep = (1-a)*s.ewmaSpareDep + a*math.Min(dSpareDep/4.0, 2.0)
	s.ewmaTempZ = (1-a)*s.ewmaTempZ + a*math.Min(tempZ/3.0, 1.5)
	s.ewmaLatZ = (1-a)*s.ewmaLatZ + a*math.Min(latZ/3.5, 1.5)

	// Additive EWMA hazard score (no explicit multi-signal motif gating or
	// per-vendor normalization).
	raw := 0.34*s.ewmaMedia +
		0.18*s.ewmaErrLog +
		0.20*s.ewmaRealloc +
		0.26*s.ewmaCRCTime +
		0.32*s.ewmaSpareDep +
		0.14*s.ewmaTempZ +
		0.16*s.ewmaLatZ

	if sample.CriticalWarningMask > 0 || sample.AvailableSparePct < 10.0 {
		raw = math.Max(raw, 2.5)
	}
	risk := 1.0 - math.Exp(-raw)
	alert := risk >= s.threshold

	var motifs []string
	if alert {
		motifs = []string{"ewma_hazard_threshold"}
	}
	leadEst := 0.0
	if alert {
		leadEst = clamp(168.0*(1.0-risk)+18.0, 12.0, 168.0)
	}
	return Prediction{
		DriveID:               sample.DriveID,
		TimestampUnix:         sample.TimestampUnix,
		RiskScore:             math.Round(risk*10000) / 10000,
		Alert:                 alert,
		LeadTimeEstimateHours: math.Round(leadEst*10) / 10,
		ActiveMotifs:          motifs,
	}
}

// ============================================================================
// Method 4: failureminer_temporal_pattern
// ============================================================================

// FailureMinerConfig configures the FailureMiner-inspired temporal pattern and
// progression-signature miner.
type FailureMinerConfig struct {
	// TargetDriveFPR sets the training false-positive budget over healthy drives
	// (e.g., 0.01 = 1% drive-level FPR, 0.005 = 0.5% under high imbalance).
	TargetDriveFPR float64
	// ShortWindowSamples (default 4 = 24h at 6h cadence).
	ShortWindowSamples int
	// MediumWindowSamples (default 12 = 72h at 6h cadence).
	MediumWindowSamples int
	// LongWindowSamples (default 28 = 168h / 7d at 6h cadence).
	LongWindowSamples int
	// BurnInSamples is the causal per-drive baseline adaptation window used to
	// normalize drive-specific hot-tier offsets and unseen vendor models.
	BurnInSamples int
	// DisableVendorNorm disables both per-vendor lookup and online per-drive
	// burn-in normalization, isolating cross-vendor distribution shift impact.
	DisableVendorNorm bool
	// CustomName optionally overrides Predictor.Name() for ablation runs.
	CustomName string
}

// DefaultFailureMinerConfig returns the default calibrated configuration.
func DefaultFailureMinerConfig() FailureMinerConfig {
	return FailureMinerConfig{
		TargetDriveFPR:      0.01,
		ShortWindowSamples:  4,
		MediumWindowSamples: 12,
		LongWindowSamples:   28,
		BurnInSamples:       12,
	}
}

type vendorStats struct {
	tempMean  float64
	tempStd   float64
	readMean  float64
	readStd   float64
	writeMean float64
	writeStd  float64
	count     int
}

// FailureMinerPredictor is a clean-room temporal pattern & progression-signature
// failure predictor inspired by FailureMiner (Wang et al., USENIX FAST '26).
//
// Design pillars:
//  1. Per-vendor/per-model + causal online per-drive burn-in normalization.
//  2. Multi-horizon causal rate-of-change slopes and acceleration (short/med/long).
//  3. Co-occurring multi-signal degradation motif mining (media ramp, thermal+spare
//     cascade, recurrent CRC/timeout storm, critical hardware flag).
//  4. Cross-window persistence gating & FPR-budgeted threshold calibration.
type FailureMinerPredictor struct {
	cfg           FailureMinerConfig
	vendorBaselines map[string]vendorStats
	globalBaseline  vendorStats
	threshold     float64
}

// NewFailureMinerPredictor constructs a FailureMiner temporal-pattern predictor.
func NewFailureMinerPredictor(cfg FailureMinerConfig) *FailureMinerPredictor {
	if cfg.TargetDriveFPR <= 0 {
		cfg.TargetDriveFPR = 0.01
	}
	if cfg.ShortWindowSamples <= 1 {
		cfg.ShortWindowSamples = 4
	}
	if cfg.MediumWindowSamples <= cfg.ShortWindowSamples {
		cfg.MediumWindowSamples = 12
	}
	if cfg.LongWindowSamples <= cfg.MediumWindowSamples {
		cfg.LongWindowSamples = 28
	}
	if cfg.BurnInSamples <= 2 {
		cfg.BurnInSamples = 12
	}
	return &FailureMinerPredictor{
		cfg:             cfg,
		vendorBaselines: make(map[string]vendorStats),
		globalBaseline: vendorStats{
			tempMean:  45.0,
			tempStd:   4.0,
			readMean:  360.0,
			readStd:   90.0,
			writeMean: 500.0,
			writeStd:  130.0,
		},
		threshold: 0.48,
	}
}

func (p *FailureMinerPredictor) Name() string {
	if p.cfg.CustomName != "" {
		return p.cfg.CustomName
	}
	return "failureminer_temporal_pattern"
}

func vendorKey(vendor, model string) string {
	return vendor + "|" + model
}

func (p *FailureMinerPredictor) Fit(train []DriveSeries) {
	if len(train) == 0 {
		return
	}

	type accum struct {
		temps  []float64
		reads  []float64
		writes []float64
	}
	byVendor := make(map[string]*accum)
	var global accum

	for _, d := range train {
		if d.Failed {
			continue
		}
		k := vendorKey(d.Meta.Vendor, d.Meta.Model)
		a, ok := byVendor[k]
		if !ok {
			a = &accum{}
			byVendor[k] = a
		}
		for _, s := range d.Samples {
			a.temps = append(a.temps, s.TemperatureC)
			a.reads = append(a.reads, s.P99ReadLatencyUs)
			a.writes = append(a.writes, s.P99WriteLatencyUs)

			global.temps = append(global.temps, s.TemperatureC)
			global.reads = append(global.reads, s.P99ReadLatencyUs)
			global.writes = append(global.writes, s.P99WriteLatencyUs)
		}
	}

	if len(global.temps) > 0 {
		tm, ts := robustMedianMAD(global.temps, 3.0)
		rm, rs := robustMedianMAD(global.reads, 60.0)
		wm, ws := robustMedianMAD(global.writes, 80.0)
		p.globalBaseline = vendorStats{
			tempMean:  tm,
			tempStd:   ts,
			readMean:  rm,
			readStd:   rs,
			writeMean: wm,
			writeStd:  ws,
			count:     len(global.temps),
		}
	}

	p.vendorBaselines = make(map[string]vendorStats, len(byVendor))
	for k, a := range byVendor {
		tm, ts := robustMedianMAD(a.temps, 2.5)
		rm, rs := robustMedianMAD(a.reads, 45.0)
		wm, ws := robustMedianMAD(a.writes, 60.0)
		p.vendorBaselines[k] = vendorStats{
			tempMean:  tm,
			tempStd:   ts,
			readMean:  rm,
			readStd:   rs,
			writeMean: wm,
			writeStd:  ws,
			count:     len(a.temps),
		}
	}

	// Calibrate alert threshold strictly on healthy training drives to meet
	// TargetDriveFPR without peeking at failing test trajectories.
	peaks := make([]float64, 0, len(train))
	for _, d := range train {
		if d.Failed {
			continue
		}
		st := p.newStreamWithThreshold(d.Meta, 2.0)
		peak := 0.0
		for _, s := range d.Samples {
			pred := st.Observe(s)
			if pred.RiskScore > peak {
				peak = pred.RiskScore
			}
		}
		peaks = append(peaks, peak)
	}
	if len(peaks) > 0 {
		q := quantile(peaks, 1.0-p.cfg.TargetDriveFPR)
		// For tight FPR budgets (<= 1.5%), add a small finite-sample safety
		// margin above the empirical quantile; for loose budgets, honor the
		// exact quantile so FPR-budget sweeps behave faithfully.
		margin := 0.0
		if p.cfg.TargetDriveFPR <= 0.006 {
			margin = 0.03
		} else if p.cfg.TargetDriveFPR <= 0.015 {
			margin = 0.015
		}
		p.threshold = clamp(q+margin, 0.14, 0.90)
	}
}

func (p *FailureMinerPredictor) NewStream(meta DriveMeta) StreamScorer {
	return p.newStreamWithThreshold(meta, p.threshold)
}

func (p *FailureMinerPredictor) newStreamWithThreshold(meta DriveMeta, threshold float64) *failureMinerStream {
	base := p.globalBaseline
	knownVendor := false
	if !p.cfg.DisableVendorNorm {
		if vb, ok := p.vendorBaselines[vendorKey(meta.Vendor, meta.Model)]; ok {
			base = vb
			knownVendor = true
		}
	}
	return &failureMinerStream{
		cfg:         p.cfg,
		vendorBase:  base,
		knownVendor: knownVendor,
		threshold:   threshold,
		buf:         make([]fmStep, 0, p.cfg.LongWindowSamples+1),
	}
}

type fmStep struct {
	ts           int64
	tempC        float64
	readLat      float64
	writeLat     float64
	sparePct     float64
	dMedia       float64
	dErrLog      float64
	dRealloc     float64
	dCRC         float64
	dTimeout     float64
	dSpareDrop   float64
	critWarn     uint8
	mediaTotal   uint64
	reallocTotal uint64
}

type failureMinerStream struct {
	cfg         FailureMinerConfig
	vendorBase  vendorStats
	knownVendor bool
	threshold   float64

	hasPrev    bool
	prevSample DriveTelemetrySample
	buf        []fmStep

	// Causal online per-drive burn-in baseline (accumulated strictly from early
	// observations <= t; never uses future samples).
	burnInCount  int
	burnTempSum  float64
	burnReadSum  float64
	burnWriteSum float64

	persistScore float64
}

func (s *failureMinerStream) effectiveBaseline() (tempMean, tempStd, readBase, writeBase float64) {
	if s.cfg.DisableVendorNorm {
		return s.vendorBase.tempMean, s.vendorBase.tempStd, s.vendorBase.readMean, s.vendorBase.writeMean
	}
	if !s.knownVendor {
		if s.burnInCount == 0 {
			return s.vendorBase.tempMean, math.Max(s.vendorBase.tempStd, 3.2), s.vendorBase.readMean, s.vendorBase.writeMean
		}
		driveTemp := s.burnTempSum / float64(s.burnInCount)
		driveRead := math.Max(60.0, s.burnReadSum/float64(s.burnInCount))
		driveWrite := math.Max(80.0, s.burnWriteSum/float64(s.burnInCount))
		return driveTemp, math.Max(s.vendorBase.tempStd, 3.2), driveRead, driveWrite
	}
	if s.burnInCount < 4 {
		return s.vendorBase.tempMean, s.vendorBase.tempStd, s.vendorBase.readMean, s.vendorBase.writeMean
	}
	driveTemp := s.burnTempSum / float64(s.burnInCount)
	driveRead := s.burnReadSum / float64(s.burnInCount)
	driveWrite := s.burnWriteSum / float64(s.burnInCount)

	// For known vendors, blend drive offset (capturing hot-tier steady state)
	// with the fitted vendor/model prior.
	const wDrive = 0.72
	tempMean = wDrive*driveTemp + (1-wDrive)*s.vendorBase.tempMean
	tempStd = math.Max(s.vendorBase.tempStd, 2.6)
	readBase = math.Max(60.0, wDrive*driveRead+(1-wDrive)*s.vendorBase.readMean)
	writeBase = math.Max(80.0, wDrive*driveWrite+(1-wDrive)*s.vendorBase.writeMean)
	return tempMean, tempStd, readBase, writeBase
}

func (s *failureMinerStream) Observe(sample DriveTelemetrySample) Prediction {
	var dMedia, dErrLog, dRealloc, dCRC, dTimeout, dSpareDrop float64
	if s.hasPrev {
		dMedia = float64(safeSubUint64(sample.MediaErrors, s.prevSample.MediaErrors))
		dErrLog = float64(safeSubUint64(sample.NumErrLogEntries, s.prevSample.NumErrLogEntries))
		dRealloc = float64(safeSubUint64(sample.ReallocatedSectors, s.prevSample.ReallocatedSectors))
		dCRC = float64(safeSubUint64(sample.CRCErrors, s.prevSample.CRCErrors))
		dTimeout = float64(safeSubUint64(sample.IOTimeoutCount, s.prevSample.IOTimeoutCount))
		if s.prevSample.AvailableSparePct > sample.AvailableSparePct {
			dSpareDrop = s.prevSample.AvailableSparePct - sample.AvailableSparePct
		}
	}
	s.hasPrev = true
	s.prevSample = sample

	// Update causal burn-in baseline only during the initial prefix when the
	// drive is not actively throwing critical warnings or rapid error bursts.
	if !s.cfg.DisableVendorNorm && s.burnInCount < s.cfg.BurnInSamples &&
		sample.CriticalWarningMask == 0 && dMedia == 0 && dTimeout == 0 && dSpareDrop < 1.5 {
		s.burnInCount++
		s.burnTempSum += sample.TemperatureC
		s.burnReadSum += sample.P99ReadLatencyUs
		s.burnWriteSum += sample.P99WriteLatencyUs
	}

	step := fmStep{
		ts:           sample.TimestampUnix,
		tempC:        sample.TemperatureC,
		readLat:      sample.P99ReadLatencyUs,
		writeLat:     sample.P99WriteLatencyUs,
		sparePct:     sample.AvailableSparePct,
		dMedia:       dMedia,
		dErrLog:      dErrLog,
		dRealloc:     dRealloc,
		dCRC:         dCRC,
		dTimeout:     dTimeout,
		dSpareDrop:   dSpareDrop,
		critWarn:     sample.CriticalWarningMask,
		mediaTotal:   sample.MediaErrors,
		reallocTotal: sample.ReallocatedSectors,
	}
	if len(s.buf) >= s.cfg.LongWindowSamples {
		copy(s.buf, s.buf[1:])
		s.buf[len(s.buf)-1] = step
	} else {
		s.buf = append(s.buf, step)
	}

	tempMean, tempStd, readBase, writeBase := s.effectiveBaseline()

	shortSlice := tailSteps(s.buf, s.cfg.ShortWindowSamples)
	medSlice := tailSteps(s.buf, s.cfg.MediumWindowSamples)
	longSlice := s.buf

	// Aggregate multi-horizon causal progression metrics.
	var (
		dMediaShort, dMediaMed       float64
		dErrLogShort, dErrLogMed     float64
		dReallocShort, dReallocMed   float64
		dCRCShort, dCRCMed           float64
		dTimeoutShort, dTimeoutMed   float64
		spareDropMed, spareDropLong  float64
		activeMediaSteps             int
		activeReallocSteps           int
		activeCRCTimeoutSteps        int
		activeSpareDropSteps         int
		tempZShortSum, tempZMedSum   float64
		readRatioShort, readRatioMed float64
		writeRatioShort              float64
	)

	for _, st := range shortSlice {
		dMediaShort += st.dMedia
		dErrLogShort += st.dErrLog
		dReallocShort += st.dRealloc
		dCRCShort += st.dCRC
		dTimeoutShort += st.dTimeout
		tempZShortSum += (st.tempC - tempMean) / tempStd
		readRatioShort += st.readLat / readBase
		writeRatioShort += st.writeLat / writeBase
	}
	if len(shortSlice) > 0 {
		n := float64(len(shortSlice))
		tempZShortSum /= n
		readRatioShort /= n
		writeRatioShort /= n
	}

	for _, st := range medSlice {
		dMediaMed += st.dMedia
		dErrLogMed += st.dErrLog
		dReallocMed += st.dRealloc
		dCRCMed += st.dCRC
		dTimeoutMed += st.dTimeout
		spareDropMed += st.dSpareDrop
		if st.dMedia > 0 {
			activeMediaSteps++
		}
		if st.dRealloc > 0 {
			activeReallocSteps++
		}
		if st.dCRC > 0 || st.dTimeout > 0 {
			activeCRCTimeoutSteps++
		}
		if st.dSpareDrop >= 0.4 {
			activeSpareDropSteps++
		}
		tempZMedSum += (st.tempC - tempMean) / tempStd
		readRatioMed += st.readLat / readBase
	}
	if len(medSlice) > 0 {
		n := float64(len(medSlice))
		tempZMedSum /= n
		readRatioMed /= n
	}
	for _, st := range longSlice {
		spareDropLong += st.dSpareDrop
	}

	// Causal OLS slope of available spare (% per day) over medium window.
	spareSlopePerDay := spareSlopeOLSPerDay(medSlice)

	// -------------------------------------------------------------------------
	// Mine 4 co-occurring multi-signal degradation motifs (FailureMiner core)
	// -------------------------------------------------------------------------
	var motifs []string
	var motifScores [4]float64

	// Motif 1: Progressive media-error / sector-reallocation ramp co-occurring
	// with error-log growth or sustained read-latency elevation.
	mediaProgressSteps := maxInt(activeMediaSteps, activeReallocSteps)
	if (dMediaMed >= 2.0 || (dMediaMed >= 1.0 && dReallocMed >= 2.0) || dReallocMed >= 3.0) &&
		mediaProgressSteps >= 2 &&
		(dErrLogMed >= 2.0 || readRatioShort >= 1.25 || spareDropMed >= 1.5) {
		evidence := 0.30*math.Log1p(dMediaMed+0.75*dReallocMed) +
			0.18*math.Log1p(dErrLogMed) +
			0.18*math.Max(0, readRatioShort-1.10) +
			0.14*float64(minInt(mediaProgressSteps, 6))/6.0
		if dMediaShort+dReallocShort >= 2.0 {
			evidence += 0.14
		}
		motifScores[0] = clamp(evidence, 0.0, 0.98)
		if motifScores[0] >= 0.30 {
			motifs = append(motifs, "motif_media_ramp")
		}
	}

	// Motif 2: Sustained thermal elevation + accelerating spare depletion +
	// write-latency / reallocation cascade.
	thermalExcess := tempZShortSum
	if s.cfg.DisableVendorNorm {
		if thermalExcess >= 1.6 && (writeRatioShort >= 1.35 || readRatioShort >= 1.35 || spareDropMed >= 0.5) {
			evidence := 0.26*math.Min(thermalExcess, 4.0) +
				0.24*math.Max(0, writeRatioShort-1.0) +
				0.20*math.Max(0, readRatioShort-1.0)
			motifScores[1] = clamp(evidence, 0.0, 0.96)
			if motifScores[1] >= 0.30 {
				motifs = append(motifs, "motif_thermal_unnormalized")
			}
		}
	}
	if (spareDropMed >= 3.0 || spareDropLong >= 5.5) &&
		activeSpareDropSteps >= 2 &&
		spareSlopePerDay <= -0.9 &&
		(thermalExcess >= 1.4 || writeRatioShort >= 1.30 || dReallocMed >= 2.0) {
		evidence := 0.26*math.Min(spareDropMed/6.0, 2.0) +
			0.20*math.Min(-spareSlopePerDay/4.0, 1.5) +
			0.20*math.Max(0, math.Min(thermalExcess/2.5, 1.5)) +
			0.18*math.Max(0, writeRatioShort-1.10) +
			0.10*math.Log1p(dReallocMed)
		if evidence > motifScores[1] {
			motifScores[1] = clamp(evidence, 0.0, 0.98)
		}
		if motifScores[1] >= 0.30 {
			motifs = append(motifs, "motif_thermal_spare_cascade")
		}
	}

	// Motif 3: Intermittent escalating CRC + I/O timeout burst co-occurring
	// with error-log accumulation or read tail-latency elevation.
	if activeCRCTimeoutSteps >= 2 &&
		(dCRCMed+dTimeoutMed) >= 3.0 &&
		(dErrLogMed >= 2.0 || readRatioShort >= 1.25 || readRatioMed >= 1.20) {
		evidence := 0.28*math.Log1p(dCRCMed+1.2*dTimeoutMed) +
			0.18*math.Log1p(dErrLogMed) +
			0.20*math.Max(0, readRatioShort-1.08) +
			0.14*float64(minInt(activeCRCTimeoutSteps, 6))/6.0
		if dCRCShort+dTimeoutShort >= 2.0 {
			evidence += 0.12
		}
		motifScores[2] = clamp(evidence, 0.0, 0.98)
		if motifScores[2] >= 0.30 {
			motifs = append(motifs, "motif_crc_timeout_storm")
		}
	}

	// Motif 4: Direct NVMe critical hardware flag / severe spare exhaustion.
	if sample.CriticalWarningMask > 0 || sample.AvailableSparePct < 10.5 || dMediaShort >= 8.0 {
		motifScores[3] = 0.98
		motifs = append(motifs, "motif_critical_hardware")
	}

	// Continuous low-level background stress score (bounded strictly below
	// confirmed multi-signal motifs) so FPR calibration sweeps have a smooth,
	// non-degenerate score distribution on healthy drives.
	bgScore := clamp(
		0.05*math.Max(0, thermalExcess-0.8)+
			0.07*math.Max(0, readRatioShort-1.15)+
			0.07*math.Max(0, writeRatioShort-1.15)+
			0.06*math.Log1p(dErrLogShort+dReallocShort+dCRCShort+dTimeoutShort),
		0.0, 0.34,
	)

	// Combine strongest motif with multi-motif synergy bonus.
	maxMotif := bgScore
	activeCount := 0
	for _, sc := range motifScores {
		if sc > maxMotif {
			maxMotif = sc
		}
		if sc >= 0.30 {
			activeCount++
		}
	}
	instantScore := maxMotif
	if activeCount >= 2 {
		instantScore = clamp(instantScore+0.12*float64(activeCount-1), 0.0, 1.0)
	}

	// Cross-window persistence gate: progressive degradation compounds across
	// consecutive windows, while isolated transient blips decay rapidly.
	if instantScore >= 0.32 {
		s.persistScore = 0.62*s.persistScore + 0.38*instantScore
	} else {
		s.persistScore = 0.42*s.persistScore + 0.15*instantScore
	}

	risk := clamp(0.60*instantScore+0.40*s.persistScore, 0.0, 1.0)
	if motifScores[3] >= 0.95 {
		risk = 0.99
	}

	alert := risk >= s.threshold

	// Estimate causal remaining lead time from spare depletion velocity or
	// progression score.
	leadEst := 0.0
	if alert {
		if spareSlopePerDay < -0.8 && sample.AvailableSparePct > 10.0 {
			daysTo10 := (sample.AvailableSparePct - 10.0) / (-spareSlopePerDay)
			leadEst = clamp(daysTo10*24.0, 12.0, 240.0)
		} else {
			leadEst = clamp(180.0*(1.0-0.65*risk), 18.0, 180.0)
		}
	}

	return Prediction{
		DriveID:               sample.DriveID,
		TimestampUnix:         sample.TimestampUnix,
		RiskScore:             math.Round(risk*10000) / 10000,
		Alert:                 alert,
		LeadTimeEstimateHours: math.Round(leadEst*10) / 10,
		ActiveMotifs:          motifs,
	}
}

// ============================================================================
// Method 5: over_sensitive_negative (Documented Negative Control)
// ============================================================================

// OverSensitiveNegativePredictor is a documented negative control that flags
// any moderate temperature bump, single latency spike, or isolated counter.
// It demonstrates the false-alarm storm and unnecessary evacuation cost
// collapse when FPR is not budgeted.
type OverSensitiveNegativePredictor struct {
	MaxTempC          float64
	MaxReadLatencyUs  float64
	MaxWriteLatencyUs float64
}

// NewOverSensitiveNegativePredictor constructs the documented negative control.
func NewOverSensitiveNegativePredictor() *OverSensitiveNegativePredictor {
	return &OverSensitiveNegativePredictor{
		MaxTempC:          52.0,
		MaxReadLatencyUs:  900.0,
		MaxWriteLatencyUs: 1350.0,
	}
}

func (p *OverSensitiveNegativePredictor) Name() string { return "over_sensitive_negative" }

func (p *OverSensitiveNegativePredictor) Fit(_ []DriveSeries) {}

func (p *OverSensitiveNegativePredictor) NewStream(_ DriveMeta) StreamScorer {
	return &overSensitiveStream{cfg: *p}
}

type overSensitiveStream struct {
	cfg OverSensitiveNegativePredictor
}

func (s *overSensitiveStream) Observe(sample DriveTelemetrySample) Prediction {
	var motifs []string
	if sample.TemperatureC > s.cfg.MaxTempC {
		motifs = append(motifs, "temp_gt_52c")
	}
	if sample.P99ReadLatencyUs > s.cfg.MaxReadLatencyUs {
		motifs = append(motifs, "read_p99_spike")
	}
	if sample.P99WriteLatencyUs > s.cfg.MaxWriteLatencyUs {
		motifs = append(motifs, "write_p99_spike")
	}
	if sample.MediaErrors > 0 || sample.CRCErrors > 0 || sample.IOTimeoutCount > 0 || sample.CriticalWarningMask > 0 {
		motifs = append(motifs, "any_nonzero_error")
	}
	alert := len(motifs) > 0
	score := 0.0
	leadEst := 0.0
	if alert {
		score = clamp(0.55+0.15*float64(len(motifs)), 0.0, 1.0)
		leadEst = 72.0
	}
	return Prediction{
		DriveID:               sample.DriveID,
		TimestampUnix:         sample.TimestampUnix,
		RiskScore:             score,
		Alert:                 alert,
		LeadTimeEstimateHours: leadEst,
		ActiveMotifs:          motifs,
	}
}

// ============================================================================
// Pure stdlib math & statistics helpers
// ============================================================================

func safeSubUint64(cur, prev uint64) uint64 {
	if cur >= prev {
		return cur - prev
	}
	return 0
}

func appendCapped(slice []float64, v float64, maxLen int) []float64 {
	if len(slice) < maxLen {
		return append(slice, v)
	}
	copy(slice, slice[1:])
	slice[maxLen-1] = v
	return slice
}

func tailSteps(buf []fmStep, n int) []fmStep {
	if len(buf) <= n {
		return buf
	}
	return buf[len(buf)-n:]
}

func zScoreAgainstHistory(x float64, hist []float64, minStd float64) float64 {
	mean, std := meanStd(hist, minStd)
	if x <= mean {
		return 0.0
	}
	return (x - mean) / std
}

func meanStd(vals []float64, minStd float64) (float64, float64) {
	if len(vals) == 0 {
		return 0, minStd
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	mean := sum / float64(len(vals))
	varSq := 0.0
	for _, v := range vals {
		d := v - mean
		varSq += d * d
	}
	std := math.Sqrt(varSq / float64(len(vals)))
	if std < minStd {
		std = minStd
	}
	return mean, std
}

func robustMedianMAD(vals []float64, minScale float64) (float64, float64) {
	if len(vals) == 0 {
		return 0, minScale
	}
	cp := make([]float64, len(vals))
	copy(cp, vals)
	sort.Float64s(cp)
	med := cp[len(cp)/2]

	devs := make([]float64, len(cp))
	for i, v := range cp {
		devs[i] = math.Abs(v - med)
	}
	sort.Float64s(devs)
	mad := devs[len(devs)/2] * 1.4826
	if mad < minScale {
		mad = minScale
	}
	return med, mad
}

func quantile(vals []float64, q float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	cp := make([]float64, len(vals))
	copy(cp, vals)
	sort.Float64s(cp)
	if q <= 0 {
		return cp[0]
	}
	if q >= 1 {
		return cp[len(cp)-1]
	}
	pos := q * float64(len(cp)-1)
	lo := int(math.Floor(pos))
	hi := int(math.Ceil(pos))
	if lo == hi {
		return cp[lo]
	}
	frac := pos - float64(lo)
	return cp[lo]*(1-frac) + cp[hi]*frac
}

func spareSlopeOLSPerDay(steps []fmStep) float64 {
	n := len(steps)
	if n < 3 {
		return 0.0
	}
	t0 := steps[0].ts
	var sumX, sumY, sumXX, sumXY float64
	for _, st := range steps {
		x := float64(st.ts-t0) / 86400.0 // days
		y := st.sparePct
		sumX += x
		sumY += y
		sumXX += x * x
		sumXY += x * y
	}
	fn := float64(n)
	denom := fn*sumXX - sumX*sumX
	if math.Abs(denom) < 1e-9 {
		return 0.0
	}
	return (fn*sumXY - sumX*sumY) / denom
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
