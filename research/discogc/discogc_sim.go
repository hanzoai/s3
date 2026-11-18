package main

import (
	"fmt"
	"math"
)

// Constants representing typical Hanzo S3 hardware profiles (e.g., NVMe QLC)
const (
	SSD_Endurance_TBW = 1500.0 // 1.5PB for a 2TB drive
	Cost_Per_TB       = 80.0   // USD
	Years_Expected    = 5.0
)

type SimulationResult struct {
	Threshold      float64 // Liveness threshold (0.0 to 1.0)
	WAF            float64
	StorageWaste   float64
	AnnualCost     float64
	DataLossRatio  float64
}

// Simulate calculates the economics of DisCoGC
// livenessDist: distribution of live data across segments
func Simulate(threshold float64, livenessDist []float64) SimulationResult {
	var totalWrites float64
	var userWrites float64
	var wastedSpace float64
	var discardedData float64

	for _, liveness := range livenessDist {
		userWrites += 1.0 // Original write
		totalWrites += 1.0

		if liveness > threshold {
			// Traditional Compaction: Copy live data
			totalWrites += liveness
		} else {
			// DisCoGC: Discard segment
			discardedData += liveness
			wastedSpace += (1.0 - liveness) // This space is reclaimed but we lost the 'liveness'
		}
	}

	waf := totalWrites / userWrites
	// Economic calculation: Wear cost + Capacity cost
	wearCost := (totalWrites / SSD_Endurance_TBW) * Cost_Per_TB
	capacityCost := (wastedSpace / float64(len(livenessDist))) * Cost_Per_TB

	return SimulationResult{
		Threshold:     threshold,
		WAF:           waf,
		StorageWaste:  wastedSpace / float64(len(livenessDist)),
		AnnualCost:    (wearCost + capacityCost) / Years_Expected,
		DataLossRatio: discardedData / userWrites,
	}
}

func main() {
	// Simulated distribution: Most segments have 10-20% liveness after TTL/Deletes
	// Some have very low liveness (<5%)
	dist := []float64{0.01, 0.02, 0.05, 0.10, 0.20, 0.50, 0.005, 0.015}

	fmt.Println("DisCoGC Economic Analysis for Hanzo S3")
	fmt.Println("--------------------------------------")
	fmt.Printf("%-10s | %-6s | %-12s | %-10s\n", "Threshold", "WAF", "Annual Cost", "Data Loss")

	for t := 0.0; t <= 0.10; t += 0.02 {
		res := Simulate(t, dist)
		fmt.Printf("%-10.2f | %-6.2f | $%-11.2f | %-10.4f\n", 
			res.Threshold, res.WAF, res.AnnualCost, res.DataLossRatio)
	}
}
