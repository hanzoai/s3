package main

import (
	"fmt"
	"time"
)

// Simulation of ACOS vs Standard Placement
func main() {
	fmt.Println("=== Hanzo S3 Geo-Aware Placement Benchmark ===")
	
	// Scenario: 100GB of "Hot Data"
	dataSizeGB := 100.0
	egressCostPerGB := 0.02 // $0.02 per GB for cross-region
	
	fmt.Printf("Scenario: Recovering %.2f GB of data after node failure\n", dataSizeGB)
	
	// Standard Replication (Random Cross-DC)
	start := time.Now()
	standardCost := dataSizeGB * egressCostPerGB
	standardLatency := 150 * time.Millisecond // Simulated cross-region RTT
	fmt.Printf("[Standard] Estimated Cost: $%.2f\n", standardCost)
	fmt.Printf("[Standard] Simulated Latency: %v\n", standardLatency)
	
	// Geo-Aware (ACOS Inspired)
	acosCost := (dataSizeGB * 0.1) * egressCostPerGB // Only 10% needs cross-region
	acosLatency := 15 * time.Millisecond // Local DC RTT
	fmt.Printf("[Geo-Aware] Estimated Cost: $%.2f\n", acosCost)
	fmt.Printf("[Geo-Aware] Simulated Latency: %v\n", acosLatency)
	
	improvement := ((standardCost - acosCost) / standardCost) * 100
	fmt.Printf("\nResult: Geo-Aware placement reduces egress costs by %.2f%%\n", improvement)
}
