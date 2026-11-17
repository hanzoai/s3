package topology

import (
	"testing"
)

func TestGeoAwarePlacement(t *testing.T) {
	nodes := []DataNode{
		{ID: "n1", Rack: "rack1", DataCenter: "dc1", IsActive: true},
		{ID: "n2", Rack: "rack1", DataCenter: "dc1", IsActive: true},
		{ID: "n3", Rack: "rack2", DataCenter: "dc1", IsActive: true},
		{ID: "n4", Rack: "rack1", DataCenter: "dc2", IsActive: true},
	}

	policy := &GeoAwarePlacement{ReplicaCount: 2, PreferLocal: true}

	// Test Local Preference
	selected, err := policy.ChooseNodes(nodes, "dc1", "rack1")
	if err != nil {
		t.Fatalf("Failed to choose nodes: %v", err)
	}

	if len(selected) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(selected))
	}

	if selected[0].Rack != "rack1" || selected[1].Rack != "rack1" {
		t.Errorf("Policy failed to prioritize local rack. Selected: %+v", selected)
	}
}

func TestFailoverToRemote(t *testing.T) {
	nodes := []DataNode{
		{ID: "n1", Rack: "rack1", DataCenter: "dc1", IsActive: false}, // Local down
		{ID: "n2", Rack: "rack2", DataCenter: "dc2", IsActive: true},  // Remote up
	}

	policy := &GeoAwarePlacement{ReplicaCount: 1, PreferLocal: true}
	selected, err := policy.ChooseNodes(nodes, "dc1", "rack1")

	if err != nil {
		t.Fatalf("Should have fallen back to remote node")
	}

	if selected[0].ID != "n2" {
		t.Errorf("Expected fallback to n2, got %s", selected[0].ID)
	}
}
