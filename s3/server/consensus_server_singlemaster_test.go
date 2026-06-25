package s3server

import (
	"testing"

	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/sequence"
	"github.com/hanzoai/s3/s3/topology"
)

// TestConsensusServerSingleMasterFormsAndCommits is the end-to-end single-master
// smoke test for the Raft->Lux-consensus cutover. It drives the REAL engine (not
// the bare-struct unit shims): NewConsensusServer starts the consensus replog and
// apply loop, the single node pins itself as the sole writer, and a MaxVolumeId
// command committed through Do() finalizes and APPLIES to the topology FSM. This
// is exactly the path command/master.go takes for `-peers=none`.
func TestConsensusServerSingleMasterFormsAndCommits(t *testing.T) {
	self := pb.NewServerAddress("127.0.0.1", 9333, 19333)
	topo := topology.NewTopology("test", sequence.NewMemorySequencer(), 0, 0, false)

	// Single-master membership: just self (what checkPeers returns for peers=none).
	peers := map[string]pb.ServerAddress{string(self): self}

	cs, err := NewConsensusServer(self, peers, topo)
	if err != nil {
		t.Fatalf("NewConsensusServer: %v", err)
	}
	defer cs.Stop()

	// Wire it into the topology exactly as SetRaftServer does, so the topology's
	// IsLeader()/NextVolumeId() gates run against the live engine.
	topo.RaftServerAccessLock.Lock()
	topo.RaftServer = cs
	topo.RaftServerAccessLock.Unlock()

	// The sole node is the pinned writer.
	if !cs.IsLeader() {
		t.Fatalf("single master must be the pinned writer: Leader()=%q Name()=%q", cs.Leader(), cs.Name())
	}
	if !topo.IsLeader() {
		t.Fatal("topology.IsLeader() must be true for the single pinned writer")
	}

	// Fresh engine: no command applied yet.
	if !cs.IsLogEmpty() {
		t.Fatal("a fresh ConsensusServer must report IsLogEmpty()=true")
	}
	if cs.HasExistingState() {
		t.Fatal("a fresh ConsensusServer must report HasExistingState()=false")
	}

	// Initialize participation (the single-master self-join path in command/master.go).
	cs.DoJoinCommand()

	// After a committed+applied command, the log is no longer empty.
	if cs.IsLogEmpty() {
		t.Fatal("after DoJoinCommand a finalized command must be applied: IsLogEmpty() should be false")
	}
	if !cs.HasExistingState() {
		t.Fatal("after DoJoinCommand HasExistingState() must be true")
	}

	// Commit a MaxVolumeId command and assert it COMMITS and APPLIES to the FSM
	// (the topology's monotonic max-volume-id advances and the topology id sticks).
	const wantMax = 4242
	const wantTopoId = "topo-abc"
	if err := cs.Do(topology.NewMaxVolumeIdCommand(wantMax, wantTopoId)); err != nil {
		t.Fatalf("Do(MaxVolumeIdCommand): %v", err)
	}
	if got := topo.GetMaxVolumeId(); uint32(got) != wantMax {
		t.Fatalf("MaxVolumeId not applied to topology FSM: got %d, want %d", got, wantMax)
	}
	if got := topo.GetTopologyId(); got != wantTopoId {
		t.Fatalf("TopologyId not applied to topology FSM: got %q, want %q", got, wantTopoId)
	}

	// NextVolumeId goes through the same consensus commit path and must advance
	// monotonically past the value we just set.
	next, err := topo.NextVolumeId()
	if err != nil {
		t.Fatalf("topology.NextVolumeId() over consensus: %v", err)
	}
	if uint32(next) != wantMax+1 {
		t.Fatalf("NextVolumeId()=%d, want %d (monotonic next after %d)", next, wantMax+1, wantMax)
	}

	// Membership: the single-master /cluster/status view lists exactly self.
	addrs := cs.PeerAddresses()
	if len(addrs) != 1 {
		t.Fatalf("single-master membership should be exactly self, got %d members: %v", len(addrs), addrs)
	}
	if _, ok := addrs[cs.Name()]; !ok {
		t.Fatalf("membership must contain self %q, got %v", cs.Name(), addrs)
	}
}
