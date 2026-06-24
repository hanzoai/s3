package s3server

import (
	"testing"

	"github.com/hanzoai/s3/s3/pb"
)

// TestConsensusServerWriterPin verifies the leaderless-consensus writer pin:
// with the same membership every node deterministically agrees on ONE writer
// (the lowest member address), so exactly one master passes the IsLeader()
// write-gate, and the pin fails over to the next-lowest when the writer leaves.
func TestConsensusServerWriterPin(t *testing.T) {
	addrs := []string{"10.0.0.3:9333", "10.0.0.1:9333", "10.0.0.2:9333"}
	nodes := make([]*ConsensusServer, len(addrs))
	for i, self := range addrs {
		peers := map[string]pb.ServerAddress{}
		for _, a := range addrs { // every node sees the full membership
			peers[a] = pb.ServerAddress(a)
		}
		nodes[i] = &ConsensusServer{serverAddr: pb.ServerAddress(self), peers: peers}
	}

	const wantWriter = "10.0.0.1:9333"
	writers := 0
	for _, n := range nodes {
		if got := n.Leader(); got != wantWriter {
			t.Errorf("node %s: Leader()=%q, want %q (every node must agree)", n.Name(), got, wantWriter)
		}
		if n.Name() == n.Leader() { // this is exactly the Topology.IsLeader() gate
			writers++
		}
	}
	if writers != 1 {
		t.Fatalf("nodes passing the IsLeader write-gate = %d, want exactly 1", writers)
	}

	// Failover: the writer leaves the membership; the next-lowest must pin.
	const wantWriter2 = "10.0.0.2:9333"
	for _, n := range nodes {
		if n.Name() == wantWriter {
			continue // departed node is no longer running
		}
		n.RemovePeer(wantWriter)
		if got := n.Leader(); got != wantWriter2 {
			t.Errorf("after failover node %s: Leader()=%q, want %q", n.Name(), got, wantWriter2)
		}
	}
}

// TestConsensusServerIsLogEmpty verifies IsLogEmpty reflects whether any
// finalized command has been applied — not replog.Pending(), which is the
// un-applied queue depth and drains back to 0 (the bug that disabled
// fastResume).
func TestConsensusServerIsLogEmpty(t *testing.T) {
	cs := &ConsensusServer{serverAddr: pb.ServerAddress("10.0.0.1:9333")}
	if !cs.IsLogEmpty() {
		t.Fatal("a fresh ConsensusServer must report IsLogEmpty()=true")
	}
	cs.applied.Store(true) // the FSM apply path sets this on the first command
	if cs.IsLogEmpty() {
		t.Fatal("after a finalized command is applied, IsLogEmpty() must be false")
	}
}
