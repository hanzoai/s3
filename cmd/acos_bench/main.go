package main

import (
	"fmt"
	"time"
	"github.com/hanzoai/s3/pkg/storage/placement"
)

func main() {
	fmt.Println("Running ACOS Geo-aware Placement Experiment...")

	nodes := []placement.Node{
		{ID: "node-1", Locality: placement.Local, Latency: 5, Cost: 0.05},
		{ID: "node-2", Locality: placement.Regional, Latency: 100, Cost: 0.01},
		{ID: "node-3", Locality: placement.Regional, Latency: 105, Cost: 0.01},
		{ID: "node-4", Locality: placement.Regional, Latency: 110, Cost: 0.01},
	}

	policy := &placement.ACOS{MaxCost: 0.10}
	
	start := time.Now()
	selected := policy.SelectNodes(3, nodes)
	duration := time.Since(start)

	fmt.Printf("Experiment Results:\n")
	fmt.Printf("Time to decide: %v\n", duration)
	fmt.Printf("Selected Nodes for Active Replication:\n")
	for _, n := range selected {
		fmt.Printf("- %s (%s) Latency: %dms Cost: %f\n", n.ID, n.Locality, n.Latency, n.Cost)
	}

	// Validation
	hasLocal := false
	for _, n := range selected {
		if n.Locality == placement.Local {
			hasLocal = true
		}
	}

	if hasLocal {
		fmt.Println("Hypothesis Validated: Local replica prioritized for active data.")
	} else {
		fmt.Println("Hypothesis Failed: No local replica selected.")
	}
}
