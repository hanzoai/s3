package placement

import (
	"context"
	"sort"
)

type ReplicaLocality string

const (
	Local    ReplicaLocality = "local"
	Regional ReplicaLocality = "regional"
)

type Node struct {
	ID       string
	Locality ReplicaLocality
	Latency  int64 // ms
	Cost     float64
}

type PlacementPolicy interface {
	SelectNodes(count int, candidates []Node) []Node
}

type ACOS struct {
	MaxCost float64
}

func (p *ACOS) SelectNodes(count int, candidates []Node) []Node {
	// ACOS Logic: Prioritize at least 1 Local node, then fill with lowest cost Regional nodes
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Locality == Local && candidates[j].Locality != Local {
			return true
		}
		return candidates[i].Cost < candidates[j].Cost
	})

	if len(candidates) < count {
		return candidates
	}
	return candidates[:count]
}
