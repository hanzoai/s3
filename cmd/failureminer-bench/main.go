// Command failureminer-bench runs the reproducible HIP-1349 / HIP-1348
// predictive SSD failure detection benchmark comparing 5 strictly causal
// methods on identical NVMe/SMART train/test splits, verifies temporal
// zero-leakage, executes boundary-regime stress tests, and emits raw JSON
// and markdown score tables.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/hanzoai/s3/s3/health/failureminer"
)

type honestyNotice struct {
	DataSource             string `json:"data_source"`
	RealHardwareExercised  bool   `json:"real_hardware_exercised"`
	AuthoringHostOnly      bool   `json:"authoring_host_only"`
	NotHanzoProductionSpec bool   `json:"not_hanzo_production_claim"`
	Note                   string `json:"note"`
}

type hostInfo struct {
	GoVersion string `json:"go_version"`
	GOOS      string `json:"goos"`
	GOARCH    string `json:"goarch"`
	NumCPU    int    `json:"num_cpu"`
}

type benchmarkReport struct {
	GeneratedAtUTC  string                              `json:"generated_at_utc"`
	Honesty         honestyNotice                       `json:"honesty_notice"`
	Host            hostInfo                            `json:"host"`
	DatasetName     string                              `json:"dataset_name"`
	Seed            int64                               `json:"seed"`
	TotalDrives     int                                 `json:"total_drives"`
	TrainDrives     int                                 `json:"train_drives"`
	TestDrives      int                                 `json:"test_drives"`
	HorizonDays     int                                 `json:"horizon_days"`
	StepHours       int                                 `json:"step_hours"`
	EvalConfig      failureminer.EvalConfig             `json:"eval_config"`
	MainResults     []failureminer.MethodMetrics        `json:"main_results"`
	CausalityAudits []failureminer.CausalityAudit       `json:"causality_audits"`
	BoundaryRegimes []failureminer.BoundaryRegimeResult `json:"boundary_regimes"`
}

func main() {
	var (
		seed          = flag.Int64("seed", 1349, "deterministic PRNG seed for synthetic NVMe/SMART fleet generation")
		numDrives     = flag.Int("drives", 1200, "number of NVMe drives in synthetic fleet")
		horizonDays   = flag.Int("horizon-days", 90, "observation window in days")
		stepHours     = flag.Int("step-hours", 6, "telemetry polling interval in hours")
		failureRate   = flag.Float64("failure-rate", 0.04, "fraction of drives that fail within horizon")
		hardNegRate   = flag.Float64("hard-negative-rate", 0.20, "fraction of healthy drives modeled as hot-tier hard negatives")
		suddenShare   = flag.Float64("sudden-death-share", 0.15, "fraction of failing drives with zero pre-failure precursors")
		externalCSV   = flag.String("external-smart-csv", "", "optional path to external Backblaze/Alibaba-compatible SMART CSV")
		externalJSON  = flag.String("external-smart-json", "", "optional path to external JSON FleetDataset")
		outDir        = flag.String("out-dir", "", "optional directory to write raw_results.json and score_table.md")
		skipBoundary  = flag.Bool("skip-boundary", false, "skip boundary-regime stress tests")
	)
	flag.Parse()

	var (
		ds         failureminer.FleetDataset
		err        error
		dataSource = "deterministic_seeded_nvme_smart_simulation"
		realHW     = false
	)

	switch {
	case *externalCSV != "":
		ds, err = failureminer.LoadSMARTCSV(*externalCSV)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading external SMART CSV: %v\n", err)
			os.Exit(1)
		}
		dataSource = "external_smart_csv:" + filepath.Base(*externalCSV)
		realHW = true
	case *externalJSON != "":
		ds, err = failureminer.LoadSMARTJSON(*externalJSON)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading external SMART JSON: %v\n", err)
			os.Exit(1)
		}
		dataSource = "external_smart_json:" + filepath.Base(*externalJSON)
		realHW = true
	default:
		cfg := failureminer.DefaultCorpusConfig(*seed)
		cfg.NumDrives = *numDrives
		cfg.HorizonDays = *horizonDays
		cfg.StepHours = *stepHours
		cfg.FailureRate = *failureRate
		cfg.HardNegativeRate = *hardNegRate
		cfg.SuddenDeathShare = *suddenShare
		ds = failureminer.GenerateFleet(cfg)
	}

	train, test := ds.SplitTrainTest(0.5)
	evalCfg := failureminer.DefaultEvalConfig()

	predictors := []failureminer.Predictor{
		failureminer.NewStaticSMARTThresholdPredictor(),
		failureminer.NewSlidingWindowZScorePredictor(failureminer.DefaultSlidingWindowZScoreConfig()),
		failureminer.NewEWMAHazardPredictor(failureminer.DefaultEWMAHazardConfig()),
		failureminer.NewFailureMinerPredictor(failureminer.DefaultFailureMinerConfig()),
		failureminer.NewOverSensitiveNegativePredictor(),
	}

	mainResults := make([]failureminer.MethodMetrics, 0, len(predictors))
	audits := make([]failureminer.CausalityAudit, 0, len(predictors))

	for _, pred := range predictors {
		res := failureminer.EvaluateMethod(pred, train, test, evalCfg)
		mainResults = append(mainResults, res)

		audit := failureminer.VerifyTemporalCausality(pred, train, test)
		audits = append(audits, audit)
		if !audit.TrainTestDisjoint || !audit.FuturePerturbationInvariant {
			fmt.Fprintf(os.Stderr, "FATAL: temporal causality audit failed for %s: %s\n", pred.Name(), audit.Details)
			os.Exit(2)
		}
	}

	var boundaries []failureminer.BoundaryRegimeResult
	if !*skipBoundary {
		boundaries = failureminer.RunBoundaryRegimes(*seed, evalCfg)
	}

	report := benchmarkReport{
		GeneratedAtUTC: time.Now().UTC().Format(time.RFC3339),
		Honesty: honestyNotice{
			DataSource:             dataSource,
			RealHardwareExercised:  realHW,
			AuthoringHostOnly:      true,
			NotHanzoProductionSpec: true,
			Note:                   "All figures are reproducible measurements from the evaluation corpus on the authoring host; they are not Hanzo production SLA claims until validated on production NVMe SMART logs.",
		},
		Host: hostInfo{
			GoVersion: runtime.Version(),
			GOOS:      runtime.GOOS,
			GOARCH:    runtime.GOARCH,
			NumCPU:    runtime.NumCPU(),
		},
		DatasetName:     ds.Name,
		Seed:            *seed,
		TotalDrives:     len(ds.Drives),
		TrainDrives:     len(train),
		TestDrives:      len(test),
		HorizonDays:     ds.HorizonDays,
		StepHours:       ds.StepHours,
		EvalConfig:      evalCfg,
		MainResults:     mainResults,
		CausalityAudits: audits,
		BoundaryRegimes: boundaries,
	}

	md := renderScoreTableMarkdown(report)
	fmt.Print(md)

	if *outDir != "" {
		if err := os.MkdirAll(*outDir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "error creating out-dir: %v\n", err)
			os.Exit(1)
		}
		rawPath := filepath.Join(*outDir, "raw_results.json")
		rawBytes, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error encoding raw json: %v\n", err)
			os.Exit(1)
		}
		rawBytes = append(rawBytes, '\n')
		if err := os.WriteFile(rawPath, rawBytes, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing %s: %v\n", rawPath, err)
			os.Exit(1)
		}
		mdPath := filepath.Join(*outDir, "score_table.md")
		if err := os.WriteFile(mdPath, []byte(md), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing %s: %v\n", mdPath, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Wrote %s and %s\n", rawPath, mdPath)
	}
}

func renderScoreTableMarkdown(r benchmarkReport) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# HIP-1349 / HIP-1348 FailureMiner Benchmark Score Table\n\n")
	fmt.Fprintf(&b, "- **Dataset**: `%s` (seed=`%d`, total_drives=`%d`, train=`%d`, test=`%d`, horizon=`%dd`, step=`%dh`)\n",
		r.DatasetName, r.Seed, r.TotalDrives, r.TrainDrives, r.TestDrives, r.HorizonDays, r.StepHours)
	fmt.Fprintf(&b, "- **Data Source**: `%s` (`real_hardware_exercised=%v`)\n", r.Honesty.DataSource, r.Honesty.RealHardwareExercised)
	fmt.Fprintf(&b, "- **Runtime Host**: `%s` (`%s/%s`, `%d` CPUs)\n", r.Host.GoVersion, r.Host.GOOS, r.Host.GOARCH, r.Host.NumCPU)
	fmt.Fprintf(&b, "- **Honesty Notice**: %s\n\n", r.Honesty.Note)

	b.WriteString("## 1. Matched 5-Method Fleet Comparison (Test Split)\n\n")
	b.WriteString("| Method | TP/FP/FN | Recall | Pred. Recall | Precision (Evac Acc) | F1 | Drive FPR | FAR / Drive-Day | Median Lead (h) | Mean Lead (h) | P10 / P90 Lead (h) | >=24h Recall | Hard-Neg FP | Cost Savings | Wasted Evac (TB) | Throughput (samp/s) | Latency (ns/samp) |\n")
	b.WriteString("|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|\n")
	for _, m := range r.MainResults {
		fmt.Fprintf(&b, "| `%s` | %d/%d/%d | %.1f%% | %.1f%% | %.1f%% | %.3f | %.2f%% | %.6f | %.1f | %.1f | %.1f / %.1f | %.1f%% | %d/%d | %+.1f%% | %.1f | %.0f | %.0f |\n",
			m.Method,
			m.TruePositives, m.FalsePositives, m.FalseNegatives,
			m.Recall*100, m.PredictableRecall*100, m.Precision*100, m.F1,
			m.DriveFPR*100, m.FalseAlarmsPerDriveDay,
			m.MedianLeadTimeHours, m.MeanLeadTimeHours,
			m.P10LeadTimeHours, m.P90LeadTimeHours,
			m.Actionable24hRecall*100,
			m.HardNegativeFalsePositives, m.HardNegativeDrives,
			m.CostSavingsPct, m.WastedEvacuationTB,
			m.ThroughputSamplesPerSec, m.LatencyNsPerSample,
		)
	}

	b.WriteString("\n## 2. Per-Failure-Mode Detection & Median Lead Time\n\n")
	b.WriteString("| Method | Failure Mode | Detected / Total | Rate | Median Lead (h) |\n")
	b.WriteString("|---|---|---:|---:|---:|\n")
	for _, m := range r.MainResults {
		for _, mb := range m.ModeBreakdown {
			fmt.Fprintf(&b, "| `%s` | `%s` | %d / %d | %.1f%% | %.1f |\n",
				m.Method, mb.Mode, mb.DetectedDrives, mb.TotalDrives, mb.DetectionRate*100, mb.MedianLeadTimeHours)
		}
	}

	b.WriteString("\n## 3. Strict Temporal Causality & Zero-Leakage Audit\n\n")
	b.WriteString("| Method | Train/Test Disjoint | Future Perturbation Invariant | Drives / Cutoffs / Observations |\n")
	b.WriteString("|---|---:|---:|---:|\n")
	for _, a := range r.CausalityAudits {
		fmt.Fprintf(&b, "| `%s` | %v | %v | %d / %d / %d |\n",
			a.Method, a.TrainTestDisjoint, a.FuturePerturbationInvariant,
			a.DrivesAudited, a.CutoffsAudited, a.ComparisonsChecked)
	}

	if len(r.BoundaryRegimes) > 0 {
		b.WriteString("\n## 4. Honest Failure & Boundary Regimes\n\n")
		for _, br := range r.BoundaryRegimes {
			fmt.Fprintf(&b, "### `%s`\n\n", br.RegimeName)
			fmt.Fprintf(&b, "- **Setup**: %s\n", br.Description)
			fmt.Fprintf(&b, "- **Finding**: %s\n\n", br.Takeaway)
			b.WriteString("| Method | TP/FP/FN | Recall | Precision | F1 | Drive FPR | Median Lead (h) | Cost Savings | Wasted Evac (TB) |\n")
			b.WriteString("|---|---:|---:|---:|---:|---:|---:|---:|---:|\n")
			for _, m := range br.Metrics {
				fmt.Fprintf(&b, "| `%s` | %d/%d/%d | %.1f%% | %.1f%% | %.3f | %.2f%% | %.1f | %+.1f%% | %.1f |\n",
					m.Method, m.TruePositives, m.FalsePositives, m.FalseNegatives,
					m.Recall*100, m.Precision*100, m.F1, m.DriveFPR*100,
					m.MedianLeadTimeHours, m.CostSavingsPct, m.WastedEvacuationTB)
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}
