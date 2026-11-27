package failureminer

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// VendorProfile parameterizes the healthy operating envelope of an NVMe SSD
// vendor/model family. Distinct profiles intentionally create cross-vendor
// distribution shift (different idle/load temperatures and tail latencies).
type VendorProfile struct {
	Vendor              string  `json:"vendor"`
	Model               string  `json:"model"`
	BaseTempC           float64 `json:"base_temp_c"`
	TempStdC            float64 `json:"temp_std_c"`
	BaseP99ReadUs       float64 `json:"base_p99_read_us"`
	BaseP99WriteUs      float64 `json:"base_p99_write_us"`
	LatencyCV           float64 `json:"latency_cv"`
	WearPerDayPct       float64 `json:"wear_per_day_pct"`
	CapacityTB          float64 `json:"capacity_tb"`
	BenignErrorProb     float64 `json:"benign_error_prob"`
	HotTierTempOffsetC  float64 `json:"hot_tier_temp_offset_c"`
	HotTierLatencyScale float64 `json:"hot_tier_latency_scale"`
}

// DefaultVendorProfiles returns three realistic enterprise NVMe SSD profiles
// spanning cool Gen4 TLC, higher-latency high-capacity QLC, and warm Gen5 TLC.
func DefaultVendorProfiles() []VendorProfile {
	return []VendorProfile{
		{
			Vendor:              "vendor_a_nvme",
			Model:               "nvme_gen4_tlc_7tb",
			BaseTempC:           37.0,
			TempStdC:            2.4,
			BaseP99ReadUs:       280.0,
			BaseP99WriteUs:      360.0,
			LatencyCV:           0.14,
			WearPerDayPct:       0.025,
			CapacityTB:          7.68,
			BenignErrorProb:     0.03,
			HotTierTempOffsetC:  10.0,
			HotTierLatencyScale: 1.45,
		},
		{
			Vendor:              "vendor_b_nvme",
			Model:               "nvme_gen4_qlc_15tb",
			BaseTempC:           45.5,
			TempStdC:            2.8,
			BaseP99ReadUs:       590.0,
			BaseP99WriteUs:      890.0,
			LatencyCV:           0.18,
			WearPerDayPct:       0.035,
			CapacityTB:          15.36,
			BenignErrorProb:     0.04,
			HotTierTempOffsetC:  9.5,
			HotTierLatencyScale: 1.50,
		},
		{
			Vendor:              "vendor_c_nvme",
			Model:               "nvme_gen5_tlc_3tb",
			BaseTempC:           51.5,
			TempStdC:            3.0,
			BaseP99ReadUs:       220.0,
			BaseP99WriteUs:      310.0,
			LatencyCV:           0.15,
			WearPerDayPct:       0.030,
			CapacityTB:          3.84,
			BenignErrorProb:     0.03,
			HotTierTempOffsetC:  10.5,
			HotTierLatencyScale: 1.40,
		},
	}
}

// CorpusConfig controls synthetic NVMe/SMART fleet generation.
type CorpusConfig struct {
	Name             string          `json:"name"`
	Seed             int64           `json:"seed"`
	NumDrives        int             `json:"num_drives"`
	HorizonDays      int             `json:"horizon_days"`
	StepHours        int             `json:"step_hours"`
	FailureRate      float64         `json:"failure_rate"`
	HardNegativeRate float64         `json:"hard_negative_rate"`
	SuddenDeathShare float64         `json:"sudden_death_share"`
	NumTenants       int             `json:"num_tenants"`
	Vendors          []VendorProfile `json:"vendors"`
}

// DefaultCorpusConfig returns the standard 90-day multi-vendor NVMe fleet
// evaluation configuration (1,200 drives, 6h cadence, 4% failure rate, 20%
// hard-negative hot-tier healthy drives, 15% sudden-death share of failures).
func DefaultCorpusConfig(seed int64) CorpusConfig {
	return CorpusConfig{
		Name:             "nvme_fleet_90d_multivendor",
		Seed:             seed,
		NumDrives:        1200,
		HorizonDays:      90,
		StepHours:        6,
		FailureRate:      0.04,
		HardNegativeRate: 0.20,
		SuddenDeathShare: 0.15,
		NumTenants:       8,
		Vendors:          DefaultVendorProfiles(),
	}
}

// GenerateFleet builds a deterministic, reproducible NVMe/SMART telemetry
// corpus with realistic class imbalance, multi-vendor operating envelopes,
// distinct progressive failure modes, hard-negative hot-tier healthy drives,
// transient benign datacenter noise, and unannounced sudden controller deaths.
func GenerateFleet(cfg CorpusConfig) FleetDataset {
	if cfg.NumDrives <= 0 {
		cfg.NumDrives = 200
	}
	if cfg.HorizonDays <= 0 {
		cfg.HorizonDays = 90
	}
	if cfg.StepHours <= 0 {
		cfg.StepHours = 6
	}
	if cfg.NumTenants <= 0 {
		cfg.NumTenants = 4
	}
	if len(cfg.Vendors) == 0 {
		cfg.Vendors = DefaultVendorProfiles()
	}
	if cfg.Name == "" {
		cfg.Name = "nvme_fleet_synthetic"
	}

	rng := rand.New(rand.NewPCG(uint64(cfg.Seed), uint64(cfg.Seed)^0x9e3779b97f4a7c15))

	totalSteps := (cfg.HorizonDays * 24) / cfg.StepHours
	stepSec := int64(cfg.StepHours) * 3600
	// Fixed deterministic epoch: 2026-01-01T00:00:00Z.
	const baseEpochUnix int64 = 1767225600

	numFail := int(math.Round(float64(cfg.NumDrives) * cfg.FailureRate))
	if cfg.FailureRate > 0 && numFail < 1 {
		numFail = 1
	}
	if numFail > cfg.NumDrives {
		numFail = cfg.NumDrives
	}

	// Partition failing drives deterministically across failure modes.
	numSudden := int(math.Round(float64(numFail) * cfg.SuddenDeathShare))
	if cfg.SuddenDeathShare >= 1.0 {
		numSudden = numFail
	}
	remFail := numFail - numSudden
	numMedia := int(math.Round(float64(remFail) * 0.40))
	numThermal := int(math.Round(float64(remFail) * 0.34))
	numCRC := remFail - numMedia - numThermal
	if numCRC < 0 {
		numCRC = 0
	}

	numHealthy := cfg.NumDrives - numFail
	numHardNeg := int(math.Round(float64(numHealthy) * cfg.HardNegativeRate))
	if numHardNeg > numHealthy {
		numHardNeg = numHealthy
	}

	modes := make([]FailureMode, 0, cfg.NumDrives)
	for i := 0; i < numMedia; i++ {
		modes = append(modes, FailureModeMediaErrorRamp)
	}
	for i := 0; i < numThermal; i++ {
		modes = append(modes, FailureModeThermalSpareDepletion)
	}
	for i := 0; i < numCRC; i++ {
		modes = append(modes, FailureModeCRCTimeoutBurst)
	}
	for i := 0; i < numSudden; i++ {
		modes = append(modes, FailureModeSuddenControllerDeath)
	}
	for i := 0; i < numHardNeg; i++ {
		modes = append(modes, FailureModeHardNegativeHot)
	}
	for len(modes) < cfg.NumDrives {
		modes = append(modes, FailureModeNone)
	}
	rng.Shuffle(len(modes), func(i, j int) { modes[i], modes[j] = modes[j], modes[i] })

	drives := make([]DriveSeries, 0, cfg.NumDrives)
	for i := 0; i < cfg.NumDrives; i++ {
		vp := cfg.Vendors[i%len(cfg.Vendors)]
		tenantID := fmt.Sprintf("tenant-%02d", (i%cfg.NumTenants)+1)
		driveID := fmt.Sprintf("%s-drv-%04d", vp.Vendor, i+1)
		mode := modes[i]

		series := simulateDrive(rng, DriveMeta{
			DriveID:  driveID,
			TenantID: tenantID,
			Vendor:   vp.Vendor,
			Model:    vp.Model,
		}, vp, mode, totalSteps, cfg.StepHours, stepSec, baseEpochUnix)
		drives = append(drives, series)
	}

	return FleetDataset{
		Name:        cfg.Name,
		Seed:        cfg.Seed,
		StepHours:   cfg.StepHours,
		HorizonDays: cfg.HorizonDays,
		Drives:      drives,
	}
}

func simulateDrive(
	rng *rand.Rand,
	meta DriveMeta,
	vp VendorProfile,
	mode FailureMode,
	totalSteps int,
	stepHours int,
	stepSec int64,
	baseEpochUnix int64,
) DriveSeries {
	isHot := mode == FailureModeHardNegativeHot
	failed := mode.IsFailure()

	// Drive-specific steady offsets around vendor profile.
	driveTempOffset := rng.NormFloat64() * 1.5
	if isHot {
		driveTempOffset += vp.HotTierTempOffsetC + rng.Float64()*3.0
	}
	driveReadScale := math.Max(0.75, 1.0+rng.NormFloat64()*0.11)
	driveWriteScale := math.Max(0.75, 1.0+rng.NormFloat64()*0.12)
	if isHot {
		driveReadScale *= vp.HotTierLatencyScale
		driveWriteScale *= vp.HotTierLatencyScale
	}

	powerOnHours := uint64(2000 + rng.IntN(22000))
	pctUsed := 5.0 + rng.Float64()*35.0
	if isHot {
		pctUsed = 35.0 + rng.Float64()*35.0
	}
	availSpare := 98.0 + rng.Float64()*2.0
	if isHot {
		availSpare = 89.0 + rng.Float64()*7.0
	}

	var (
		mediaErrors       uint64
		errLogEntries     uint64
		unsafeShutdowns   uint64 = uint64(rng.IntN(4))
		busyTime          uint64 = uint64(1000 + rng.IntN(8000))
		ioTimeouts        uint64
		crcErrors         uint64
		reallocated       uint64
		critWarn          uint8
		failureStep       int = -1
		onsetStep         int = totalSteps + 1
		failureTimeUnix   int64
		rampDurationSteps int
	)

	// Realistic static pre-existing benign counters on a minority of healthy
	// hot-tier drives (e.g., an old isolated single media error or historic
	// CRC error that never progresses).
	if isHot && rng.Float64() < 0.24 {
		mediaErrors = 1
		errLogEntries = 1
	} else if !failed && rng.Float64() < vp.BenignErrorProb {
		errLogEntries = 1
	}
	if isHot && rng.Float64() < 0.20 {
		crcErrors = 1 + uint64(rng.IntN(2))
	}

	// Select a minority of hot-tier healthy drives to undergo a non-fatal
	// transient datacenter stress episode (e.g., rack cooling fan maintenance or
	// heavy compaction + single retired flash block that stabilizes cleanly).
	hasBenignEpisode := isHot && rng.Float64() < 0.35
	benignEpisodeStart := 40 + rng.IntN(maxInt(10, totalSteps-70))
	benignEpisodeType := rng.IntN(3) // 0=heatwave+latency, 1=single-block-retire, 2=transient-cable-hiccup

	if failed {
		minFailStep := (25 * 24) / stepHours
		if minFailStep >= totalSteps-4 {
			minFailStep = totalSteps / 2
		}
		failureStep = minFailStep + rng.IntN(totalSteps-minFailStep)
		failureTimeUnix = baseEpochUnix + int64(failureStep)*stepSec

		switch mode {
		case FailureModeMediaErrorRamp:
			// 5 to 11 days of progressive precursor degradation.
			rampHours := 120 + rng.IntN(144)
			rampDurationSteps = maxInt(8, rampHours/stepHours)
			onsetStep = maxInt(8, failureStep-rampDurationSteps)
		case FailureModeThermalSpareDepletion:
			// 6 to 13 days of thermal + spare erosion.
			rampHours := 144 + rng.IntN(168)
			rampDurationSteps = maxInt(10, rampHours/stepHours)
			onsetStep = maxInt(8, failureStep-rampDurationSteps)
		case FailureModeCRCTimeoutBurst:
			// 4 to 9 days of intermittent escalating CRC/timeout bursts.
			rampHours := 96 + rng.IntN(120)
			rampDurationSteps = maxInt(8, rampHours/stepHours)
			onsetStep = maxInt(8, failureStep-rampDurationSteps)
		case FailureModeSuddenControllerDeath:
			onsetStep = failureStep + 1
		}
	}

	// Only emit telemetry up to failureStep-1 for failed drives (strictly
	// pre-failure observations).
	emitSteps := totalSteps
	if failed {
		emitSteps = failureStep
	}
	samples := make([]DriveTelemetrySample, 0, emitSteps)

	// Periodic hot-tier compaction bursts (1-2 intervals).
	compactionPeriod := 28 + rng.IntN(24)
	compactionPhase := rng.IntN(compactionPeriod)

	for step := 0; step < emitSteps; step++ {
		ts := baseEpochUnix + int64(step)*stepSec
		powerOnHours += uint64(stepHours)
		pctUsed = math.Min(100.0, pctUsed+vp.WearPerDayPct*(float64(stepHours)/24.0))

		busyInc := uint64(stepHours * (10 + rng.IntN(18)))
		if isHot {
			busyInc = uint64(stepHours * (32 + rng.IntN(22)))
		}
		busyTime += busyInc

		diurnal := 1.6 * math.Sin(2.0*math.Pi*float64(step*stepHours)/24.0)
		tempC := vp.BaseTempC + driveTempOffset + diurnal + rng.NormFloat64()*vp.TempStdC

		readLat := vp.BaseP99ReadUs * driveReadScale * math.Max(0.55, 1.0+rng.NormFloat64()*vp.LatencyCV)
		writeLat := vp.BaseP99WriteUs * driveWriteScale * math.Max(0.55, 1.0+rng.NormFloat64()*vp.LatencyCV)

		// Benign isolated workload spikes on healthy / hot-tier drives.
		if isHot && (step+compactionPhase)%compactionPeriod == 0 {
			tempC += 5.0 + rng.Float64()*3.5
			readLat *= 2.4 + rng.Float64()*1.3
			writeLat *= 2.8 + rng.Float64()*1.6
		} else if !failed && rng.Float64() < 0.004 {
			readLat *= 2.5 + rng.Float64()*1.1
			writeLat *= 2.2 + rng.Float64()*1.0
		}

		// Realistic non-fatal stress episodes on healthy hot-tier drives:
		// tests whether predictors separate isolated single-cause / self-healing
		// disturbances from progressive multi-signal SSD failure.
		if !failed && hasBenignEpisode {
			dtEp := step - benignEpisodeStart
			switch benignEpisodeType {
			case 0:
				// 3-step rack thermal + compaction latency surge with ZERO spare/media erosion.
				if dtEp >= 0 && dtEp < 3 {
					tempC += 8.5 + rng.Float64()*2.5
					writeLat *= 2.6 + rng.Float64()*1.2
					readLat *= 2.1 + rng.Float64()*0.9
					if dtEp == 1 {
						errLogEntries++ // single transient command abort
					}
				}
			case 1:
				// Single benign flash-block retirement (1 step) that stabilizes completely.
				if dtEp == 0 {
					reallocated += 1 + uint64(rng.IntN(2))
					errLogEntries++
					availSpare = math.Max(82.0, availSpare-1.0)
					writeLat *= 1.8
				}
			case 2:
				// Isolated single-interval backplane/PCIe retry hiccup that does NOT recur.
				if dtEp == 0 {
					crcErrors += 2
					ioTimeouts++
					errLogEntries++
					readLat *= 2.2
				}
			}
		}

		// Rare isolated unsafe shutdown on healthy drives.
		if !failed && rng.Float64() < 0.0015 {
			unsafeShutdowns++
		}

		// Rare isolated benign single-counter jitter on general healthy drives
		// (models normal flash ECC / backplane background noise in large fleets).
		if !failed && !isHot && rng.Float64() < 0.0005 {
			errLogEntries++
		}

		// Inject physically grounded progressive degradation once step >= onsetStep.
		if failed && step >= onsetStep {
			progress := float64(step-onsetStep+1) / float64(maxInt(1, failureStep-onsetStep))
			if progress > 1.0 {
				progress = 1.0
			}
			accel := progress * progress
			stepsLeft := failureStep - step

			switch mode {
			case FailureModeMediaErrorRamp:
				// Early precursors (days 1-5 of ramp): error-log growth, reallocated
				// sectors, and read latency drift; uncorrectable media_errors emerge
				// mid-to-late ramp so static SMART media_errors>0 lags behind motif detection.
				if progress >= 0.06 && rng.Float64() < (0.50+0.48*progress) {
					errLogEntries += 1 + uint64(math.Round(3.5*accel))
				}
				if progress >= 0.10 && rng.Float64() < (0.48+0.50*progress) {
					reallocated += 1 + uint64(math.Round(4.0*accel))
				}
				if progress >= 0.52 && rng.Float64() < (0.45+0.55*progress) {
					mediaErrors += 1 + uint64(math.Round(5.5*accel))
				}
				readLat *= 1.0 + 0.60*progress + 2.5*accel
				writeLat *= 1.0 + 0.35*progress + 1.5*accel
				tempC += 2.0*progress + 3.5*accel
				// Spare erodes gently until the final 24h before failure.
				spareDrop := 0.25 + 1.4*accel
				if stepsLeft <= 3 {
					spareDrop += 18.0
				}
				availSpare = math.Max(3.0, availSpare-spareDrop)
				if stepsLeft <= 2 || availSpare < 10.0 {
					critWarn |= 0x05
				}

			case FailureModeThermalSpareDepletion:
				// Sustained thermal creep + calibrated spare erosion (~98% -> ~18%
				// across the precursor window, crossing <10% only in the final 12-24h
				// for ~65% of drives, while ~35% lock read-only near 12-16% spare).
				tempC += 7.5 + 9.0*progress + rng.Float64()*2.5
				targetDropPerStep := 72.0 / float64(maxInt(8, rampDurationSteps))
				spareDrop := targetDropPerStep * (0.55 + 0.90*progress)
				if stepsLeft <= 3 && (step+int(powerOnHours))%3 != 0 {
					spareDrop += 8.0
				}
				availSpare = math.Max(5.0, availSpare-spareDrop)
				if rng.Float64() < (0.52 + 0.45*progress) {
					reallocated += 1 + uint64(math.Round(3.0*progress))
				}
				if progress > 0.30 && rng.Float64() < 0.45 {
					errLogEntries += 1 + uint64(math.Round(2.0*accel))
				}
				writeLat *= 1.0 + 0.95*progress + 2.8*accel
				readLat *= 1.0 + 0.35*progress + 1.1*accel
				if stepsLeft <= 3 && availSpare < 10.0 {
					critWarn |= 0x03
				}
				if stepsLeft <= 1 && rng.Float64() < 0.20 {
					mediaErrors++
				}

			case FailureModeCRCTimeoutBurst:
				// Intermittent escalating bursts of CRC + I/O timeouts + error log + read latency.
				burstActive := (step-onsetStep)%2 == 0 || progress > 0.55
				if burstActive {
					crcErrors += 1 + uint64(math.Round(4.0*progress))
					ioTimeouts += 1 + uint64(math.Round(3.5*accel))
					errLogEntries += 1 + uint64(math.Round(3.0*progress))
					readLat *= 1.65 + 2.6*progress
					writeLat *= 1.30 + 1.8*progress
					tempC += 2.5 * progress
				} else {
					readLat *= 1.15 + 0.5*progress
				}
				// Spare remains healthy (>25%) and media_errors stays 0 until final interval at most.
				availSpare = math.Max(30.0, availSpare-0.15*progress)
				if stepsLeft <= 1 && rng.Float64() < 0.25 {
					critWarn |= 0x04
				}
			}
		}

		samples = append(samples, DriveTelemetrySample{
			TimestampUnix:       ts,
			DriveID:             meta.DriveID,
			TenantID:            meta.TenantID,
			Vendor:              meta.Vendor,
			Model:               meta.Model,
			MediaErrors:         mediaErrors,
			NumErrLogEntries:    errLogEntries,
			AvailableSparePct:   math.Round(availSpare*100) / 100,
			PercentageUsed:      math.Round(pctUsed*100) / 100,
			CriticalWarningMask: critWarn,
			TemperatureC:        math.Round(tempC*100) / 100,
			UnsafeShutdowns:     unsafeShutdowns,
			ControllerBusyTime:  busyTime,
			IOTimeoutCount:      ioTimeouts,
			P99ReadLatencyUs:    math.Round(readLat*10) / 10,
			P99WriteLatencyUs:   math.Round(writeLat*10) / 10,
			CRCErrors:           crcErrors,
			ReallocatedSectors:  reallocated,
			PowerOnHours:        powerOnHours,
		})
	}

	return DriveSeries{
		Meta:            meta,
		Failed:          failed,
		FailureMode:     mode,
		FailureTimeUnix: failureTimeUnix,
		Samples:         samples,
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
