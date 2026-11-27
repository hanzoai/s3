package failureminer

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestDeterministicCorpusReproducibility(t *testing.T) {
	cfg := DefaultCorpusConfig(4242)
	cfg.NumDrives = 120
	cfg.HorizonDays = 45

	ds1 := GenerateFleet(cfg)
	ds2 := GenerateFleet(cfg)

	if !reflect.DeepEqual(ds1, ds2) {
		t.Fatalf("GenerateFleet is not deterministic for identical seed %d", cfg.Seed)
	}
	if len(ds1.Drives) != cfg.NumDrives {
		t.Fatalf("got %d drives, want %d", len(ds1.Drives), cfg.NumDrives)
	}

	// Verify every generated sample passes validation and is strictly ascending in time.
	for _, d := range ds1.Drives {
		var prevTS int64
		for i, s := range d.Samples {
			if err := s.Validate(); err != nil {
				t.Fatalf("drive %s sample %d invalid: %v", d.Meta.DriveID, i, err)
			}
			if i > 0 && s.TimestampUnix <= prevTS {
				t.Fatalf("drive %s non-monotonic timestamps at %d: %d <= %d",
					d.Meta.DriveID, i, s.TimestampUnix, prevTS)
			}
			prevTS = s.TimestampUnix
		}
		if d.Failed && len(d.Samples) > 0 {
			lastTS := d.Samples[len(d.Samples)-1].TimestampUnix
			if lastTS >= d.FailureTimeUnix {
				t.Fatalf("drive %s has sample at/after failure time: %d >= %d",
					d.Meta.DriveID, lastTS, d.FailureTimeUnix)
			}
		}
	}
}

func TestTemporalCausalityZeroLeakage(t *testing.T) {
	cfg := DefaultCorpusConfig(1349)
	cfg.NumDrives = 300
	cfg.HorizonDays = 60
	ds := GenerateFleet(cfg)
	train, test := ds.SplitTrainTest(0.5)

	predictors := []Predictor{
		NewStaticSMARTThresholdPredictor(),
		NewSlidingWindowZScorePredictor(DefaultSlidingWindowZScoreConfig()),
		NewEWMAHazardPredictor(DefaultEWMAHazardConfig()),
		NewFailureMinerPredictor(DefaultFailureMinerConfig()),
		NewOverSensitiveNegativePredictor(),
	}

	for _, pred := range predictors {
		audit := VerifyTemporalCausality(pred, train, test)
		if !audit.TrainTestDisjoint {
			t.Fatalf("%s failed disjoint train/test check: %s", pred.Name(), audit.Details)
		}
		if !audit.FuturePerturbationInvariant {
			t.Fatalf("%s violated temporal causality: %s", pred.Name(), audit.Details)
		}
		if audit.ComparisonsChecked == 0 {
			t.Fatalf("%s checked 0 comparisons", pred.Name())
		}
	}
}

func TestFailureMinerBeatsStaticThreshold(t *testing.T) {
	cfg := DefaultCorpusConfig(1349)
	ds := GenerateFleet(cfg)
	train, test := ds.SplitTrainTest(0.5)
	evalCfg := DefaultEvalConfig()

	staticRes := EvaluateMethod(NewStaticSMARTThresholdPredictor(), train, test, evalCfg)
	fmRes := EvaluateMethod(NewFailureMinerPredictor(DefaultFailureMinerConfig()), train, test, evalCfg)

	if fmRes.Recall <= staticRes.Recall {
		t.Fatalf("expected failureminer recall (%.4f) > static recall (%.4f)", fmRes.Recall, staticRes.Recall)
	}
	if fmRes.Actionable24hRecall <= staticRes.Actionable24hRecall {
		t.Fatalf("expected failureminer >=24h recall (%.4f) > static (%.4f)",
			fmRes.Actionable24hRecall, staticRes.Actionable24hRecall)
	}
	if fmRes.MedianLeadTimeHours <= staticRes.MedianLeadTimeHours {
		t.Fatalf("expected failureminer median lead time (%.1fh) > static (%.1fh)",
			fmRes.MedianLeadTimeHours, staticRes.MedianLeadTimeHours)
	}
	if fmRes.DriveFPR >= staticRes.DriveFPR {
		t.Fatalf("expected failureminer drive FPR (%.4f) < static (%.4f)", fmRes.DriveFPR, staticRes.DriveFPR)
	}
	if fmRes.F1 <= staticRes.F1 {
		t.Fatalf("expected failureminer F1 (%.4f) > static F1 (%.4f)", fmRes.F1, staticRes.F1)
	}
	if fmRes.WastedEvacuationTB >= staticRes.WastedEvacuationTB {
		t.Fatalf("expected failureminer wasted evacuation (%.1f TB) < static (%.1f TB)",
			fmRes.WastedEvacuationTB, staticRes.WastedEvacuationTB)
	}
	// Honest recall ceiling: mixed fleet includes sudden controller deaths.
	if fmRes.SuddenDeathDrives == 0 {
		t.Fatalf("expected sudden-death drives in test split")
	}
	if fmRes.Recall >= 1.0 {
		t.Fatalf("expected honest recall ceiling < 1.0 on fleet with sudden deaths, got %.4f", fmRes.Recall)
	}
}

func TestSuddenDeathRecallCeilingAndNegativeControl(t *testing.T) {
	evalCfg := DefaultEvalConfig()

	// 1. Sudden unannounced controller death ceiling.
	suddenCfg := DefaultCorpusConfig(2026)
	suddenCfg.NumDrives = 240
	suddenCfg.FailureRate = 0.06
	suddenCfg.SuddenDeathShare = 1.0
	ds := GenerateFleet(suddenCfg)
	train, test := ds.SplitTrainTest(0.5)

	fmSudden := EvaluateMethod(NewFailureMinerPredictor(DefaultFailureMinerConfig()), train, test, evalCfg)
	negSudden := EvaluateMethod(NewOverSensitiveNegativePredictor(), train, test, evalCfg)

	if fmSudden.Recall > 0.15 {
		t.Fatalf("expected calibrated causal recall near 0 on zero-precursor sudden death, got %.4f", fmSudden.Recall)
	}
	if fmSudden.DriveFPR > 0.03 {
		t.Fatalf("expected bounded FPR <= 3%% on sudden death fleet, got %.4f", fmSudden.DriveFPR)
	}

	// 2. Documented negative control must collapse precision and operational cost.
	if negSudden.DriveFPR < 0.85 {
		t.Fatalf("expected over_sensitive_negative FPR >= 85%%, got %.4f", negSudden.DriveFPR)
	}
	if negSudden.Precision > 0.15 {
		t.Fatalf("expected over_sensitive_negative precision <= 15%%, got %.4f", negSudden.Precision)
	}
	if negSudden.CostSavingsPct >= 0 {
		t.Fatalf("expected over_sensitive_negative to produce negative cost savings, got %.2f%%", negSudden.CostSavingsPct)
	}
}

func TestBoundaryRegimesHighImbalanceAndVendorShift(t *testing.T) {
	evalCfg := DefaultEvalConfig()
	regimes := RunBoundaryRegimes(1349, evalCfg)
	if len(regimes) != 3 {
		t.Fatalf("expected 3 boundary regimes, got %d", len(regimes))
	}

	byName := make(map[string]BoundaryRegimeResult, len(regimes))
	for _, r := range regimes {
		byName[r.RegimeName] = r
	}

	// Verify high-imbalance tight FPR budget outperforms loose FPR budget on precision.
	imb, ok := byName["high_imbalance_base_rate_precision_trap"]
	if !ok {
		t.Fatalf("missing high_imbalance_base_rate_precision_trap regime")
	}
	var loosePrec, tightPrec float64
	for _, m := range imb.Metrics {
		switch m.Method {
		case "failureminer_loose_fpr_3pct":
			loosePrec = m.Precision
		case "failureminer_tight_fpr_0.4pct":
			tightPrec = m.Precision
		}
	}
	if tightPrec <= loosePrec {
		t.Fatalf("expected tight FPR precision (%.4f) > loose FPR precision (%.4f) under 1%% prevalence",
			tightPrec, loosePrec)
	}

	// Verify cross-vendor shift normalization dramatically reduces FPR and improves precision.
	shift, ok := byName["cross_vendor_distribution_shift"]
	if !ok {
		t.Fatalf("missing cross_vendor_distribution_shift regime")
	}
	var noNormFPR, normFPR, noNormPrec, normPrec float64
	for _, m := range shift.Metrics {
		switch m.Method {
		case "failureminer_no_vendor_norm":
			noNormFPR = m.DriveFPR
			noNormPrec = m.Precision
		case "failureminer_with_causal_burnin_norm":
			normFPR = m.DriveFPR
			normPrec = m.Precision
		}
	}
	if normFPR >= noNormFPR {
		t.Fatalf("expected causal burn-in norm FPR (%.4f) < unnormalized FPR (%.4f) under cross-vendor shift",
			normFPR, noNormFPR)
	}
	if normPrec <= noNormPrec {
		t.Fatalf("expected causal burn-in norm precision (%.4f) > unnormalized precision (%.4f)",
			normPrec, noNormPrec)
	}
}

func TestMultiTenantConcurrentFleetMonitor(t *testing.T) {
	cfg := DefaultCorpusConfig(77)
	cfg.NumDrives = 80
	cfg.HorizonDays = 30
	ds := GenerateFleet(cfg)
	train, test := ds.SplitTrainTest(0.5)

	fm := NewFailureMinerPredictor(DefaultFailureMinerConfig())
	fm.Fit(train)
	monitor := NewFleetMonitor(fm)

	// Concurrently stream all test drives across goroutines (-race verified).
	var wg sync.WaitGroup
	errCh := make(chan error, len(test))
	for _, d := range test {
		drive := d
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, s := range drive.Samples {
				if _, err := monitor.Ingest(s); err != nil {
					errCh <- err
					return
				}
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatalf("concurrent Ingest error: %v", err)
	}

	// Verify strict multi-tenant state isolation: identical drive_id under two
	// different tenants must not cross-contaminate.
	isoMonitor := NewFleetMonitor(fm)
	const sharedDriveID = "shared-nvme-slot-0"
	var baseTS int64 = 1767225600
	for step := 0; step < 20; step++ {
		ts := baseTS + int64(step)*6*3600
		healthySample := DriveTelemetrySample{
			TimestampUnix:     ts,
			DriveID:           sharedDriveID,
			TenantID:          "tenant-isolated-healthy",
			Vendor:            "vendor_a_nvme",
			Model:             "nvme_gen4_tlc_7tb",
			AvailableSparePct: 100.0,
			TemperatureC:      37.0,
			P99ReadLatencyUs:  280.0,
			P99WriteLatencyUs: 360.0,
			PowerOnHours:      uint64(1000 + step*6),
		}
		failingSample := DriveTelemetrySample{
			TimestampUnix:       ts,
			DriveID:             sharedDriveID,
			TenantID:            "tenant-isolated-failing",
			Vendor:              "vendor_a_nvme",
			Model:               "nvme_gen4_tlc_7tb",
			MediaErrors:         uint64(step * 4),
			NumErrLogEntries:    uint64(step * 6),
			ReallocatedSectors:  uint64(step * 3),
			AvailableSparePct:   math.Max(5.0, 98.0-float64(step)*5.0),
			CriticalWarningMask: 0x05,
			TemperatureC:        68.0,
			P99ReadLatencyUs:    1800.0,
			P99WriteLatencyUs:   2200.0,
			PowerOnHours:        uint64(1000 + step*6),
		}
		pHealthy, err := isoMonitor.Ingest(healthySample)
		if err != nil {
			t.Fatalf("healthy ingest failed: %v", err)
		}
		if pHealthy.Alert {
			t.Fatalf("healthy tenant alerted at step %d due to tenant bleed: %+v", step, pHealthy)
		}
		if _, err := isoMonitor.Ingest(failingSample); err != nil {
			t.Fatalf("failing ingest failed: %v", err)
		}
	}

	sumHealthy := isoMonitor.SummaryForTenant("tenant-isolated-healthy")
	sumFailing := isoMonitor.SummaryForTenant("tenant-isolated-failing")
	if sumHealthy.TotalAlerts != 0 {
		t.Fatalf("expected 0 alerts for isolated healthy tenant, got %d", sumHealthy.TotalAlerts)
	}
	if sumFailing.TotalAlerts == 0 {
		t.Fatalf("expected >0 alerts for isolated failing tenant")
	}

	// Verify rejection of out-of-order timestamps and invalid inputs.
	badOrder := DriveTelemetrySample{
		TimestampUnix:     baseTS,
		DriveID:           sharedDriveID,
		TenantID:          "tenant-isolated-healthy",
		AvailableSparePct: 100.0,
		TemperatureC:      37.0,
	}
	if _, err := isoMonitor.Ingest(badOrder); err == nil {
		t.Fatalf("expected error for non-monotonic timestamp, got nil")
	}

	nanSample := badOrder
	nanSample.TimestampUnix = baseTS + 999999
	nanSample.TemperatureC = math.NaN()
	if _, err := isoMonitor.Ingest(nanSample); err == nil {
		t.Fatalf("expected error for NaN temperature, got nil")
	}
}

func TestExternalSMARTCSVAndJSONLoader(t *testing.T) {
	// Backblaze / Alibaba-style CSV schema with flexible aliases and unsorted rows.
	csvData := strings.Join([]string{
		"date,serial_number,model,failure,smart_187_raw,smart_197_raw,smart_232_normalized,smart_194_raw,smart_5_raw,smart_199_raw,smart_9_raw,p99_read_latency_us,p99_write_latency_us",
		"2026-01-03,BB-DRV-01,NVME_A_7TB,0,0,1,98,39,1,0,1048,310,410",
		"2026-01-01,BB-DRV-01,NVME_A_7TB,0,0,0,100,37,0,0,1000,280,350",
		"2026-01-02,BB-DRV-01,NVME_A_7TB,0,0,0,99,38,0,0,1024,290,370",
		"2026-01-04,BB-DRV-01,NVME_A_7TB,0,2,5,92,43,4,0,1072,620,740",
		"2026-01-05,BB-DRV-01,NVME_A_7TB,1,12,18,65,49,14,0,1096,1900,2100",
		"2026-01-01,BB-DRV-02,NVME_A_7TB,0,0,0,100,38,0,0,2000,270,340",
		"2026-01-02,BB-DRV-02,NVME_A_7TB,0,0,0,100,37,0,0,2024,275,350",
		"2026-01-03,BB-DRV-02,NVME_A_7TB,0,0,0,100,38,0,0,2048,280,355",
		"2026-01-04,BB-DRV-02,NVME_A_7TB,0,0,0,100,39,0,0,2072,272,348",
		"2026-01-05,BB-DRV-02,NVME_A_7TB,0,0,0,100,38,0,0,2096,278,352",
	}, "\n")

	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "backblaze_sample.csv")
	if err := os.WriteFile(csvPath, []byte(csvData), 0o644); err != nil {
		t.Fatalf("write temp csv: %v", err)
	}

	ds, err := LoadSMARTCSV(csvPath)
	if err != nil {
		t.Fatalf("LoadSMARTCSV failed: %v", err)
	}
	if len(ds.Drives) != 2 {
		t.Fatalf("got %d drives, want 2", len(ds.Drives))
	}

	// BB-DRV-01 must be marked failed with 4 strictly pre-failure sorted samples.
	d0 := ds.Drives[0]
	if d0.Meta.DriveID != "BB-DRV-01" || !d0.Failed || len(d0.Samples) != 4 {
		t.Fatalf("unexpected failing drive parse: id=%s failed=%v samples=%d",
			d0.Meta.DriveID, d0.Failed, len(d0.Samples))
	}
	for i := 1; i < len(d0.Samples); i++ {
		if d0.Samples[i].TimestampUnix <= d0.Samples[i-1].TimestampUnix {
			t.Fatalf("CSV loader did not sort timestamps ascending")
		}
	}

	// Verify JSON round-trip.
	jsonPath := filepath.Join(tmpDir, "fleet.json")
	raw, err := json.Marshal(ds)
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	if err := os.WriteFile(jsonPath, raw, 0o644); err != nil {
		t.Fatalf("write json: %v", err)
	}
	dsJSON, err := LoadSMARTJSON(jsonPath)
	if err != nil {
		t.Fatalf("LoadSMARTJSON failed: %v", err)
	}
	if !reflect.DeepEqual(ds, dsJSON) {
		t.Fatalf("JSON round-trip mismatch")
	}
}
