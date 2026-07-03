package s3server

import (
	"context"
	"testing"

	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/master_pb"
	"github.com/hanzoai/s3/s3/sequence"
	"github.com/hanzoai/s3/s3/topology"
)

// The cluster RPC surface is retained after the raft rip so existing clients
// (the admin dashboard's master list) keep working — but it is now leaderless:
// Add/Remove are advisory membership hints and the "leader" is the pinned
// writer, a pure function of the membership. These tests pin that behavior.

func newMembershipTestServer(self string) *MasterServer {
	topo := topology.NewTopology("test", sequence.NewMemorySequencer(), 1<<30, 5, false)
	return &MasterServer{
		option:           &MasterOption{Master: pb.ServerAddress(self)},
		Topo:             topo,
		coordinatorPeers: make(map[string]pb.ServerAddress),
	}
}

func TestRaftMembershipRPCs_AddListPinsLowest(t *testing.T) {
	ctx := context.Background()
	ms := newMembershipTestServer("127.0.0.1:100")

	// Advisory adds widen the writer-eligible set.
	for _, id := range []string{"127.0.0.1:200", "127.0.0.1:100", "127.0.0.1:300"} {
		if _, err := ms.RaftAddServer(ctx, &master_pb.RaftAddServerRequest{Id: id}); err != nil {
			t.Fatalf("RaftAddServer(%s): %v", id, err)
		}
	}

	resp, err := ms.RaftListClusterServers(ctx, &master_pb.RaftListClusterServersRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.ClusterServers) != 3 {
		t.Fatalf("listed %d servers, want 3", len(resp.ClusterServers))
	}
	// Exactly one server is flagged the writer, and it is the lowest address.
	writers := 0
	for _, s := range resp.ClusterServers {
		if s.IsWriter {
			writers++
			if s.Id != "127.0.0.1:100" {
				t.Fatalf("writer flag on %s, want lowest 127.0.0.1:100", s.Id)
			}
		}
	}
	if writers != 1 {
		t.Fatalf("%d servers flagged writer, want exactly 1", writers)
	}
}

func TestRaftMembershipRPCs_RemoveWriterRepins(t *testing.T) {
	ctx := context.Background()
	ms := newMembershipTestServer("127.0.0.1:200")

	for _, id := range []string{"127.0.0.1:100", "127.0.0.1:200", "127.0.0.1:300"} {
		if _, err := ms.RaftAddServer(ctx, &master_pb.RaftAddServerRequest{Id: id}); err != nil {
			t.Fatal(err)
		}
	}
	// self=200 defers to writer 100.
	if ms.Topo.IsWriter() {
		t.Fatal("200 must defer to writer 100")
	}

	// Remove the writer 100 → the pin moves to the next-lowest member, 200 (self).
	if _, err := ms.RaftRemoveServer(ctx, &master_pb.RaftRemoveServerRequest{Id: "127.0.0.1:100"}); err != nil {
		t.Fatal(err)
	}
	if !ms.Topo.IsWriter() {
		t.Fatal("after removing 100, self 200 must become the writer")
	}
	if got := string(ms.Topo.Coordinator.Writer()); got != "127.0.0.1:200" {
		t.Fatalf("writer after removal = %q, want 127.0.0.1:200", got)
	}
}

func TestRaftMembershipRPCs_LeadershipTransferUnsupported(t *testing.T) {
	ctx := context.Background()
	ms := newMembershipTestServer("127.0.0.1:100")
	// There is no hand-transfer of leadership in the leaderless model.
	if _, err := ms.RaftLeadershipTransfer(ctx, &master_pb.RaftLeadershipTransferRequest{}); err == nil {
		t.Fatal("RaftLeadershipTransfer must return an error in the leaderless model")
	}
}
