// Command discogc-bench runs the reproducible DisCoGC discard-versus-compaction
// storage economics benchmark across 5 matched GC policies, 4 S3 churn
// workloads, and 3 documented negative/boundary regimes using real OS
// segment files and Linux FALLOC_FL_PUNCH_HOLE / st_blocks / /proc/self/io
// measurements.
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

	"github.com/hanzoai/s3/s3/gc/discogc"
)

type BenchmarkReport struct {
	SchemaVersion   string                         `json:"schema_version"`
	GeneratedAtUTC  string                         `json:"generated_at_utc"`
	HostEnvironment HostInfo                       `json:"host_environment"`
	DriveProfile    discogc.DriveEconomics         `json:"drive_profile"`
	HonestyNotice   string                         `json:"honesty_notice"`
	Workloads       []WorkloadResultGroup          `json:"workloads"`
	BoundaryRegimes []discogc.BoundaryRegimeResult `json:"boundary_regimes"`
}

type HostInfo struct {
	GOOS      string `json:"goos"`
	GOARCH    string `json:"goarch"`
	GoVersion string `json:"go_version"`
	NumCPU    int    `json:"num_cpu"`
}

type WorkloadResultGroup struct {
	Spec    discogc.WorkloadSpec `json:"spec"`
	Results []discogc.Metrics    `json:"results"`
}

func main() {
	outDir := flag.String("out", "", "optional directory to write raw_results.json and score_table.md")
	flag.Parse()

	drive := discogc.DefaultEnterpriseNVMe()
	workloads := discogc.DefaultWorkloads()
	policies := discogc.AllPolicies()

	report := BenchmarkReport{
		SchemaVersion:  "hip-1349-discogc-v1",
		GeneratedAtUTC: time.Now().UTC().Format(time.RFC3339),
		HostEnvironment: HostInfo{
			GOOS:      runtime.GOOS,
			GOARCH:    runtime.GOARCH,
			GoVersion: runtime.Version(),
			NumCPU:    runtime.NumCPU(),
		},
		DriveProfile:  drive,
		HonestyNotice: "All measurements were collected on the local benchmark host using real OS files, syscall.Fallocate (FALLOC_FL_PUNCH_HOLE), syscall.Stat_t.Blocks, and /proc/self/io. Flash wear and $/TB-month TCO values combine measured Logical WAF and physical space amplification with explicit DriveEconomics model parameters; they are local host measurements and parameterized model outputs, not production Hanzo cluster guarantees until reproduced on target NVMe hardware.",
	}

	fmt.Printf("DisCoGC Discard-vs-Compaction Storage Economics Benchmark\n")
	fmt.Printf("Host: %s/%s (%d CPUs, %s) | Drive Model: %s\n\n",
		runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), runtime.Version(), drive.Name)

	for _, spec := range workloads {
		group := WorkloadResultGroup{Spec: spec}
		fmt.Printf("=== Workload: %s ===\n", spec.Name)
		fmt.Printf("    %s\n", spec.Description)
		for _, p := range policies {
			m, err := discogc.RunWorkload(spec, p, drive, "")
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR workload=%s policy=%s: %v\n", spec.Name, p, err)
				os.Exit(1)
			}
			group.Results = append(group.Results, m)
			fmt.Printf("  %-32s LWA=%.3fx  PWA=%.2fx  WA=%.3fx  AvgSA=%.2fx  FinalSA=%.2fx  Reclaimed=%6.2fMiB  Copied=%6.2fMiB  Life=%.2fyr  TCO=$%.2f/TB-mo\n",
				string(p),
				m.FlashWear.LogicalWAF,
				m.FlashWear.DeviceInternalPWA,
				m.FlashWear.EffectivePhysicalWAF,
				m.AvgSpaceAmplification,
				m.FinalSpaceAmplification,
				float64(m.OSBlocksReclaimedBytes)/(1024.0*1024.0),
				float64(m.GCCopiedLiveBytes)/(1024.0*1024.0),
				m.FlashWear.EstimatedDriveLifeYears,
				m.Cost.TotalUSDPerTBMonth,
			)
		}
		fmt.Println()
		report.Workloads = append(report.Workloads, group)
	}

	fmt.Printf("=== Boundary & Negative Regimes ===\n")
	regimes, err := discogc.RunBoundaryRegimes(drive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR boundary regimes: %v\n", err)
		os.Exit(1)
	}
	report.BoundaryRegimes = regimes
	for _, r := range regimes {
		fmt.Printf("  [%s]\n    %s\n", r.Name, r.Observation)
	}
	fmt.Println()

	if *outDir != "" {
		if err := os.MkdirAll(*outDir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR mkdir %s: %v\n", *outDir, err)
			os.Exit(1)
		}
		jsonPath := filepath.Join(*outDir, "raw_results.json")
		raw, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR marshal json: %v\n", err)
			os.Exit(1)
		}
		if err := os.WriteFile(jsonPath, append(raw, '\n'), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR write %s: %v\n", jsonPath, err)
			os.Exit(1)
		}

		tablePath := filepath.Join(*outDir, "score_table.md")
		if err := os.WriteFile(tablePath, []byte(renderScoreTable(report)), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR write %s: %v\n", tablePath, err)
			os.Exit(1)
		}
		fmt.Printf("Wrote %s and %s\n", jsonPath, tablePath)
	}
}

func renderScoreTable(r BenchmarkReport) string {
	var b strings.Builder
	b.WriteString("# DisCoGC Discard-versus-Compaction Economics — Measured Score Table\n\n")
	b.WriteString(fmt.Sprintf("- **Generated (UTC)**: `%s`\n", r.GeneratedAtUTC))
	b.WriteString(fmt.Sprintf("- **Host**: `%s/%s`, `%d` CPUs, `%s`\n",
		r.HostEnvironment.GOOS, r.HostEnvironment.GOARCH, r.HostEnvironment.NumCPU, r.HostEnvironment.GoVersion))
	b.WriteString(fmt.Sprintf("- **Drive Model**: `%s` ($%.0f / %.2f TB = $%.2f/TB, %.1f DWPD, %.0f-yr warranty, %.0f P/E cycles, %.1f logical TB/mo/TB)\n\n",
		r.DriveProfile.Name, r.DriveProfile.DriveCostUSD, r.DriveProfile.RawCapacityTB,
		r.DriveProfile.DriveCostUSD/r.DriveProfile.RawCapacityTB,
		r.DriveProfile.RatedDWPD, r.DriveProfile.WarrantyYears, r.DriveProfile.ProgramEraseCycles,
		r.DriveProfile.WorkloadLogicalTBPerMonthPerTB))
	b.WriteString("> **Measurement Scope**: Logical WAF, copied live bytes, OS physical blocks (`st_blocks * 512`), `FALLOC_FL_PUNCH_HOLE` reclaimed bytes, and `/proc/self/io` deltas are directly measured on the benchmark host. Flash wear (TBW, P/E cycles, drive life) and `$/TB-month` TCO combine those measurements with the explicit drive model above. Figures are reproducible local host measurements, not Hanzo production claims until replicated on target NVMe hardware.\n\n")

	for _, wg := range r.Workloads {
		b.WriteString(fmt.Sprintf("## Workload: `%s`\n\n", wg.Spec.Name))
		b.WriteString(fmt.Sprintf("%s (seed=`%d`, keys=`%d`, churn_ops=`%d`, obj=`%d–%d B`).\n\n",
			wg.Spec.Description, wg.Spec.Seed, wg.Spec.KeyCount, wg.Spec.ChurnOps, wg.Spec.MinObjBytes, wg.Spec.MaxObjBytes))
		b.WriteString("| Policy | Logical WAF (LWA) | Effective Physical WAF (WA) | Avg Space Amp | Final Space Amp | GC Copied Live (MiB) | OS Blocks Reclaimed (MiB) | Punch Calls | Drive Life (yr) | Cap Cost ($/TB-mo) | Wear Cost ($/TB-mo) | IO Tax ($/TB-mo) | Total ($/TB-mo) |\n")
		b.WriteString("|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|\n")
		for _, m := range wg.Results {
			b.WriteString(fmt.Sprintf("| `%s` | %.3fx | %.3fx | %.2fx | %.2fx | %.2f | %.2f | %d | %.2f | $%.2f | $%.2f | $%.2f | **$%.2f** |\n",
				m.Policy,
				m.FlashWear.LogicalWAF,
				m.FlashWear.EffectivePhysicalWAF,
				m.AvgSpaceAmplification,
				m.FinalSpaceAmplification,
				float64(m.GCCopiedLiveBytes)/(1024.0*1024.0),
				float64(m.OSBlocksReclaimedBytes)/(1024.0*1024.0),
				m.PunchHoleCalls,
				m.FlashWear.EstimatedDriveLifeYears,
				m.Cost.CapacityAmortizationUSDPerTBMonth,
				m.Cost.WearReplacementUSDPerTBMonth,
				m.Cost.CompactionIOTaxUSDPerTBMonth,
				m.Cost.TotalUSDPerTBMonth,
			))
		}
		b.WriteString("\n")
	}

	b.WriteString("## Documented Negative & Boundary Regimes\n\n")
	for _, reg := range r.BoundaryRegimes {
		b.WriteString(fmt.Sprintf("### `%s`\n\n", reg.Name))
		b.WriteString(fmt.Sprintf("- **Hypothesis**: %s\n", reg.Hypothesis))
		b.WriteString(fmt.Sprintf("- **Observed**: %s\n\n", reg.Observation))
		b.WriteString("| Policy / Variant | Logical WAF | Final Space Amp | Physical Disk (KiB) | OS Blocks Reclaimed (KiB) | Punch Calls | Compaction Cycles | Extent Fragments | Seq Scan (MiB/s) | Total ($/TB-mo) |\n")
		b.WriteString("|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|\n")
		for _, m := range reg.Results {
			b.WriteString(fmt.Sprintf("| `%s` | %.3fx | %.2fx | %.1f | %.1f | %d | %d | %d | %.1f | $%.2f |\n",
				m.Policy,
				m.FlashWear.LogicalWAF,
				m.FinalSpaceAmplification,
				float64(m.PhysicalDiskBytes)/1024.0,
				float64(m.OSBlocksReclaimedBytes)/1024.0,
				m.PunchHoleCalls,
				m.CompactionCycles,
				m.ExtentFragments,
				m.SeqScanThroughputMBs,
				m.Cost.TotalUSDPerTBMonth,
			))
		}
		b.WriteString("\n")
	}
	return b.String()
}
