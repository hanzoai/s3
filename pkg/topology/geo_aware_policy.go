package topology

import (
	"errors"
	"math/rand"
)

// GeoLevel represents the physical proximity of a replica
type GeoLevel int

const (
	LocalRack GeoLevel = iota
	LocalDC
	RemoteDC
)

// GeoAwarePlacement defines how replicas should be distributed based on ACOS principles
type GeoAwarePlacement struct {
	ReplicaCount int
	PreferLocal  bool
}

// DataNode represents a simplified storage node for placement decisions
type DataNode struct {
	ID         string
	Rack       string
	DataCenter string
	IsActive   bool
}

// ChooseNodes selects the best nodes for object placement based on geo-awareness
func (p *GeoAwarePlacement) ChooseNodes(availableNodes []DataNode, preferredDC string, preferredRack string) ([]DataNode, error) {
	if len(availableNodes) < p.ReplicaCount {
		return nil, errors.New("insufficient nodes available")
	}

	var selected []DataNode
	
	// 1. Try to find a node in the same Rack (Lowest Latency/Cost)
	for _, node := range availableNodes {
		if node.IsActive && node.DataCenter == preferredDC && node.Rack == preferredRack {
			selected = append(selected, node)
			if len(selected) >= p.ReplicaCount {
				return selected, nil
			}
		}
	}

	// 2. Try to find nodes in the same DataCenter (Medium Latency, Low Cost)
	for _, node := range availableNodes {
		if node.IsActive && node.DataCenter == preferredDC && node.Rack != preferredRack {
			// Check if already selected
			alreadySelected := false
			for _, s := range selected {
				if s.ID == node.ID { alreadySelected = true; break }
			}
			if !alreadySelected {
				selected = append(selected, node)
				if len(selected) >= p.ReplicaCount {
					return selected, nil
				}
			}
		}
	}

	// 3. Fallback to Remote DataCenters (High Latency, High Cost)
	// This is the "ACOS Failover" path
	rand.Shuffle(len(availableNodes), func(i, j int) {
		availableNodes[i], availableNodes[j] = availableNodes[j], availableNodes[i]
	})

	for _, node := range availableNodes {
		alreadySelected := false
		for _, s := range selected {
			if s.ID == node.ID { alreadySelected = true; break }
		}
		if !alreadySelected && node.IsActive {
			selected = append(selected, node)
			if len(selected) >= p.ReplicaCount {
				return selected, nil
			}
		}
	}

	return selected, nil
}
