// Command acos-bench executes the reproducible HIP-1349 / HIP-1348 geo-aware
// active object placement and recovery benchmark suite comparing five matched
// placement policies, three documented negative/boundary regimes, and tenant
// data-residency sovereignty enforcement.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hanzoai/s3/s3/placement/acos"
)

func main() {
	def := acos.DefaultSuiteConfig()
	seed := flag.Int64("seed", def.Seed, "deterministic PRNG seed for workload generation")
	objects := flag.Int("objects", def.ObjectCount, "number of active objects in the matched workload")
	payloadBytes := flag.Int("payload-bytes", def.PayloadSizeByte, "payload size in bytes per object")
	reads := flag.Int("reads", def.ReadOps, "number of read operations across all workload phases")
	outDir := flag.String("out", "", "optional directory to write raw_results.json and score_table.md")
	flag.Parse()

	cfg := acos.SuiteConfig{
		Seed:            *seed,
		ObjectCount:     *objects,
		PayloadSizeByte: *payloadBytes,
		ReadOps:         *reads,
		Config:          acos.DefaultACOSConfig(),
	}

	report, err := acos.RunBenchmarkSuite(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "acos-bench error: %v\n", err)
		os.Exit(1)
	}

	md := acos.RenderScoreTableMarkdown(report)
	fmt.Print(md)

	if *outDir != "" {
		if err := os.MkdirAll(*outDir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "acos-bench mkdir %s: %v\n", *outDir, err)
			os.Exit(1)
		}

		rawPath := filepath.Join(*outDir, "raw_results.json")
		rawBytes, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "acos-bench marshal json: %v\n", err)
			os.Exit(1)
		}
		rawBytes = append(rawBytes, '\n')
		if err := os.WriteFile(rawPath, rawBytes, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "acos-bench write %s: %v\n", rawPath, err)
			os.Exit(1)
		}

		tablePath := filepath.Join(*outDir, "score_table.md")
		if err := os.WriteFile(tablePath, []byte(md), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "acos-bench write %s: %v\n", tablePath, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "wrote %s and %s\n", rawPath, tablePath)
	}
}
