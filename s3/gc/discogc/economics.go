package discogc

import "math"

// ComputeEconomics derives Logical WAF (LWA), SSD Device Internal WAF (PWA),
// Total Physical WAF (WA = LWA * PWA), flash wear (TBW, P/E cycles, effective
// DWPD, estimated drive lifetime), and total monthly storage cost ($/TB-month)
// from measured Store metrics and a DriveEconomics profile.
//
// Formula decomposition (following Bian et al., USENIX FAST '26):
//
//  1. Logical Write Amplification:
//     LWA = (UserBytesWritten + GCCopiedLiveBytes + MetadataBytesWritten) / UserBytesWritten
//
//  2. Total Physical Write Amplification:
//     WA = LWA * PWA
//     where PWA models internal SSD FTL amplification (lower when host issues
//     discard/TRIM on stale extents, higher under copy-compaction churn).
//
//  3. Flash Wear & Drive Lifetime:
//     - TBW per logical TB written = WA
//     - Annual physical TB written per TB stored = WorkloadLogicalTBPerMonthPerTB * 12 * WA
//     - Effective DWPD = (WorkloadLogicalTBPerMonthPerTB / 30.0) * WA
//     - Rated TBW per TB capacity = RatedDWPD * 365 * WarrantyYears
//     - PE cycles per logical TB written (per TB capacity) = (WA / RatedTBWPerTB) * ProgramEraseCycles
//     - Estimated drive life (years) = min(WarrantyYears, WarrantyYears * RatedDWPD / EffectiveDWPD)
//
//  4. Total Storage Economics ($/TB-month per live logical TB):
//     - CapacityAmortization = (DriveCostPerTB / (WarrantyYears * 12)) * AvgSpaceAmp * (1 + OP)
//     - WearReplacement = exact additional monthly SSD replacement capex strictly
//     caused by flash endurance shortening drive life below WarrantyYears:
//     effectiveSSDMonthly = (DriveCostPerTB / (EstimatedDriveLifeYears * 12)) * AvgSpaceAmp * (1 + OP)
//     WearReplacement = max(0, effectiveSSDMonthly - CapacityAmortization)
//     (Zero double-counting: CapacityAmortization + WearReplacement equals exact
//     endurance-limited SSD amortization.)
//     - CompactionIOTax = CompactionCPUCostPerGB * CopiedGBPerLogicalTBMonth +
//     DiscardSyscallCostPerMillion * DiscardsPerLogicalTBMonth / 1e6.
func ComputeEconomics(m *Metrics, d DriveEconomics) (FlashWearMetrics, CostBreakdown) {
	if d.RawCapacityTB <= 0 {
		d = DefaultEnterpriseNVMe()
	}

	var lwa float64 = 1.0
	if m.UserBytesWritten > 0 {
		totalLogical := float64(m.UserBytesWritten + m.GCCopiedLiveBytes + m.MetadataBytesWritten)
		lwa = totalLogical / float64(m.UserBytesWritten)
	}

	pwa := selectDevicePWA(m, d)
	wa := lwa * pwa

	ratedTBWPerTB := d.RatedDWPD * 365.0 * d.WarrantyYears
	if ratedTBWPerTB <= 0 {
		ratedTBWPerTB = 1825.0
	}

	logicalDWPD := d.WorkloadLogicalTBPerMonthPerTB / 30.0
	effectiveDWPD := logicalDWPD * wa

	driveLife := d.WarrantyYears
	if effectiveDWPD > 0 && d.RatedDWPD > 0 {
		wearLife := d.WarrantyYears * (d.RatedDWPD / effectiveDWPD)
		if wearLife < driveLife {
			driveLife = wearLife
		}
	}

	peCyclesPerLogicalTB := (wa / ratedTBWPerTB) * d.ProgramEraseCycles

	wear := FlashWearMetrics{
		LogicalWAF:              round4(lwa),
		DeviceInternalPWA:       round4(pwa),
		EffectivePhysicalWAF:    round4(wa),
		TBWPerLogicalTB:         round4(wa),
		PECyclesPerLogicalTB:    round4(peCyclesPerLogicalTB),
		EffectiveDWPD:           round4(effectiveDWPD),
		RatedDWPD:               round4(d.RatedDWPD),
		EstimatedDriveLifeYears: round4(driveLife),
		WarrantyYears:           round4(d.WarrantyYears),
	}

	// Space amplification factor: use time-averaged physical space amplification
	// when samples exist, otherwise final physical space amplification.
	spaceAmp := m.AvgSpaceAmplification
	if spaceAmp <= 0 {
		spaceAmp = m.FinalSpaceAmplification
	}
	if spaceAmp < 1.0 {
		spaceAmp = 1.0
	}

	driveCostPerTB := d.DriveCostUSD / d.RawCapacityTB
	opFactor := 1.0 + d.OverProvisioningFraction

	// 1. Capacity amortization at nominal warranty lifetime scaled by space amp.
	capAmort := (driveCostPerTB / (d.WarrantyYears * 12.0)) * spaceAmp * opFactor

	// 2. Wear replacement cost: exact additional monthly SSD capex strictly caused
	// by flash endurance shortening drive life below WarrantyYears (zero double-counting:
	// capAmort + wearCost == exact endurance-limited SSD hardware amortization).
	effectiveAmort := (driveCostPerTB / (driveLife * 12.0)) * spaceAmp * opFactor
	wearCost := math.Max(0, effectiveAmort-capAmort)

	// 3. Compaction & discard CPU / IO channel tax per live TB-month.
	var copiedRatio float64
	var discardsPerGB float64
	if m.UserBytesWritten > 0 {
		copiedRatio = float64(m.GCCopiedLiveBytes) / float64(m.UserBytesWritten)
		userGB := float64(m.UserBytesWritten) / (1024.0 * 1024.0 * 1024.0)
		if userGB > 0 {
			discardsPerGB = float64(m.PunchHoleCalls) / userGB
		}
	}
	logicalGBPerMonthPerTB := d.WorkloadLogicalTBPerMonthPerTB * 1024.0
	copiedGBPerMonthPerTB := logicalGBPerMonthPerTB * copiedRatio
	discardsPerMonthPerTB := logicalGBPerMonthPerTB * discardsPerGB

	ioTax := copiedGBPerMonthPerTB*d.CompactionCPUCostPerGBUSD +
		(discardsPerMonthPerTB/1e6)*d.DiscardSyscallCostPerMillionUSD

	cost := CostBreakdown{
		CapacityAmortizationUSDPerTBMonth: round4(capAmort),
		WearReplacementUSDPerTBMonth:      round4(wearCost),
		CompactionIOTaxUSDPerTBMonth:      round4(ioTax),
		TotalUSDPerTBMonth:                round4(capAmort + wearCost + ioTax),
	}

	return wear, cost
}

func selectDevicePWA(m *Metrics, d DriveEconomics) float64 {
	if m.GCCopiedLiveBytes == 0 && m.OSBlocksReclaimedBytes == 0 && m.PunchHoleCalls == 0 {
		if d.DeviceInternalWAFNoGC > 0 {
			return d.DeviceInternalWAFNoGC
		}
		return 1.10
	}
	switch m.Policy {
	case PolicyNoGC:
		if d.DeviceInternalWAFNoGC > 0 {
			return d.DeviceInternalWAFNoGC
		}
		return 1.10
	case PolicyAdaptiveDiscard:
		if m.OSBlocksReclaimedBytes > 0 && d.DeviceInternalWAFDiscard > 0 {
			return d.DeviceInternalWAFDiscard
		}
		if d.DeviceInternalWAFCompaction > 0 {
			return d.DeviceInternalWAFCompaction
		}
		return 1.20
	case PolicyBlindSubblockDiscard:
		// Sub-block zeroing writes partial dirty pages without freeing flash blocks.
		if d.DeviceInternalWAFCompaction > 0 {
			return d.DeviceInternalWAFCompaction
		}
		return 1.25
	default:
		if d.DeviceInternalWAFCompaction > 0 {
			return d.DeviceInternalWAFCompaction
		}
		return 1.25
	}
}

func round4(v float64) float64 {
	return math.Round(v*10000) / 10000
}
