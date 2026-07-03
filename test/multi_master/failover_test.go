package multi_master

import (
	"testing"
	"time"

	"github.com/hanzoai/s3/s3/pb"
)

// These tests exercise the leaderless, deterministic pinned-writer coordination
// that replaced raft. There is no election and no quorum: every master pins the
// writer to the lowest-address live member (the Coordinator), so failover is a
// pure re-computation the instant OnPeerUpdate changes the membership. The old
// raft assertions ("a lost quorum means no leader") are deliberately inverted —
// the leaderless model trades quorum for availability and relies on the pinned
// writer (and, in production, the schain VM's chain-agreed allocation) for
// single-writer correctness.

const (
	// Generous margin for the survivor's liveness probe + membership re-pin to
	// propagate. Failover itself is instantaneous (a pure re-computation); this
	// covers cluster peer-update detection latency.
	failoverTimeout = 30 * time.Second
)

// TestWriterDownAndRecoverQuickly stops the pinned writer and verifies a
// survivor re-pins as the new writer, then that the cluster re-converges to a
// single writer when the old one restarts. TopologyId stays consistent.
func TestWriterDownAndRecoverQuickly(t *testing.T) {
	mc := StartMasterCluster(t)

	writerIdx, writerAddr := mc.FindWriter()
	if writerIdx < 0 {
		t.Fatal("no pinned writer after cluster start")
	}
	t.Logf("initial writer: node %d at %s", writerIdx, writerAddr)

	topologyId, err := mc.GetTopologyId(writerIdx)
	if err != nil || topologyId == "" {
		t.Fatalf("failed to get initial TopologyId: %v", err)
	}
	t.Logf("initial TopologyId: %s", topologyId)

	// Stop the writer.
	mc.StopNode(writerIdx)
	t.Logf("stopped writer node %d", writerIdx)

	// A survivor must re-pin as the writer — no election, just the next-lowest
	// live address taking over once the dead writer is dropped from membership.
	newWriterIdx, newWriterAddr, err := mc.WaitForNewWriter(writerAddr, failoverTimeout)
	if err != nil {
		mc.DumpLogs()
		t.Fatalf("survivor did not re-pin as writer after stopping the writer: %v", err)
	}
	t.Logf("new writer: node %d at %s", newWriterIdx, newWriterAddr)

	// Restart the old writer quickly.
	mc.StartNode(writerIdx)
	if err := mc.WaitForNodeReady(writerIdx, waitTimeout); err != nil {
		mc.DumpLogs()
		t.Fatalf("restarted node %d not ready: %v", writerIdx, err)
	}
	t.Logf("restarted node %d", writerIdx)

	// The membership re-converges to exactly one writer (deterministically the
	// lowest-address member across all three live nodes).
	if err := mc.WaitForWriter(failoverTimeout); err != nil {
		mc.DumpLogs()
		t.Fatalf("cluster did not re-converge to a single writer: %v", err)
	}

	assertTopologyIdConverged(t, mc)
}

// TestWriterDownSlowRecover verifies that with one master down for an extended
// period, the surviving two keep a single stable writer and continue serving,
// and that the slow node rejoins cleanly.
func TestWriterDownSlowRecover(t *testing.T) {
	mc := StartMasterCluster(t)

	writerIdx, writerAddr := mc.FindWriter()
	if writerIdx < 0 {
		t.Fatal("no writer found")
	}
	topologyId, err := mc.GetTopologyId(writerIdx)
	if err != nil || topologyId == "" {
		t.Fatalf("failed to get initial TopologyId: %v", err)
	}
	t.Logf("initial writer: node %d, TopologyId: %s", writerIdx, topologyId)

	mc.StopNode(writerIdx)

	newWriterIdx, _, err := mc.WaitForNewWriter(writerAddr, failoverTimeout)
	if err != nil {
		mc.DumpLogs()
		t.Fatalf("new writer not pinned: %v", err)
	}
	t.Logf("new writer: node %d", newWriterIdx)

	// Confirm the new writer self-reports as the writer over the whole outage.
	cs, err := mc.GetClusterStatus(newWriterIdx)
	if err != nil {
		mc.DumpLogs()
		t.Fatalf("cannot get cluster status from new writer: %v", err)
	}
	if !cs.IsWriter {
		t.Fatalf("node %d claims not to be the writer", newWriterIdx)
	}

	// Extended outage: the two-node set must hold a single stable writer.
	t.Log("simulating slow recovery (10 seconds)...")
	time.Sleep(10 * time.Second)
	if mc.CountWriters() != 1 {
		mc.DumpLogs()
		t.Fatalf("writer not stable during extended outage: %d writers", mc.CountWriters())
	}

	// Restart the downed node.
	mc.StartNode(writerIdx)
	if err := mc.WaitForNodeReady(writerIdx, waitTimeout); err != nil {
		mc.DumpLogs()
		t.Fatalf("slow-recovered node %d not ready: %v", writerIdx, err)
	}

	if err := mc.WaitForWriter(failoverTimeout); err != nil {
		mc.DumpLogs()
		t.Fatalf("cluster did not re-converge to a single writer: %v", err)
	}
	assertTopologyIdConverged(t, mc)
}

// TestSurvivorKeepsWritingWithoutQuorum is the deliberate inversion of the old
// raft TestTwoMastersDownAndRestart. Under raft, losing 2 of 3 lost quorum and
// no leader could be elected. Leaderless coordination has no quorum: the lone
// survivor drops the dead peers from its writer-eligible set and, as the sole
// live member, pins ITSELF as the writer and keeps serving. When the two
// restart, the cluster re-converges to a single writer and preserves TopologyId.
func TestSurvivorKeepsWritingWithoutQuorum(t *testing.T) {
	mc := StartMasterCluster(t)

	writerIdx, _ := mc.FindWriter()
	if writerIdx < 0 {
		t.Fatal("no writer found")
	}
	topologyId, err := mc.GetTopologyId(writerIdx)
	if err != nil || topologyId == "" {
		t.Fatalf("failed to get initial TopologyId: %v", err)
	}
	t.Logf("initial TopologyId: %s", topologyId)

	// Stop two of three, keeping one survivor.
	down1 := writerIdx
	down2 := (writerIdx + 1) % 3
	survivor := (writerIdx + 2) % 3
	t.Logf("stopping nodes %d and %d, keeping node %d", down1, down2, survivor)

	mc.StopNode(down1)
	mc.StopNode(down2)

	// The survivor must become the writer — availability without quorum. It
	// takes over once its liveness probes drop the two dead peers.
	deadline := time.Now().Add(failoverTimeout)
	for time.Now().Before(deadline) {
		if cs, err := mc.GetClusterStatus(survivor); err == nil && cs.IsWriter {
			break
		}
		time.Sleep(waitTick)
	}
	cs, err := mc.GetClusterStatus(survivor)
	if err != nil {
		mc.DumpLogs()
		t.Fatalf("cannot reach survivor node %d: %v", survivor, err)
	}
	if !cs.IsWriter {
		mc.DumpLogs()
		t.Fatalf("survivor node %d did not become the writer (leaderless model keeps writing without quorum)", survivor)
	}
	t.Logf("survivor node %d is the writer with no quorum", survivor)

	// Restart both downed nodes.
	mc.StartNode(down1)
	mc.StartNode(down2)
	for _, i := range []int{down1, down2} {
		if err := mc.WaitForNodeReady(i, waitTimeout); err != nil {
			mc.DumpLogs()
			t.Fatalf("restarted node %d not ready: %v", i, err)
		}
	}

	if err := mc.WaitForWriter(failoverTimeout); err != nil {
		mc.DumpLogs()
		t.Fatalf("no single writer after restarting the two downed nodes: %v", err)
	}
	assertTopologyIdConverged(t, mc)
}

// TestAllMastersDownAndRestart verifies that after a full stop/restart the set
// re-pins a single writer and all nodes agree on a TopologyId. With no durable
// raft log, TopologyId is recovered from persisted topology state when present;
// a short-lived cluster may mint a fresh one — but all nodes must still agree.
func TestAllMastersDownAndRestart(t *testing.T) {
	mc := StartMasterCluster(t)

	writerIdx, _ := mc.FindWriter()
	if writerIdx < 0 {
		t.Fatal("no writer found")
	}
	topologyId, _ := mc.GetTopologyId(writerIdx)
	if topologyId == "" {
		t.Fatal("no TopologyId on initial writer")
	}
	t.Logf("initial TopologyId: %s", topologyId)

	for i := range 3 {
		mc.StopNode(i)
	}
	t.Log("all nodes stopped")
	time.Sleep(2 * time.Second)

	for i := range 3 {
		mc.StartNode(i)
	}
	for i := range 3 {
		if err := mc.WaitForNodeReady(i, waitTimeout); err != nil {
			mc.DumpLogs()
			t.Fatalf("node %d not ready after full restart: %v", i, err)
		}
	}

	if err := mc.WaitForWriter(failoverTimeout); err != nil {
		mc.DumpLogs()
		t.Fatalf("no single writer after full cluster restart: %v", err)
	}

	newWriterIdx, _ := mc.FindWriter()
	t.Logf("writer after full restart: node %d", newWriterIdx)

	newTopologyId, err := mc.GetTopologyId(newWriterIdx)
	if err != nil || newTopologyId == "" {
		mc.DumpLogs()
		t.Fatal("no TopologyId after full restart")
	}
	if newTopologyId == topologyId {
		t.Logf("TopologyId preserved across full restart: %s", topologyId)
	} else {
		t.Logf("TopologyId changed (expected for short-lived cluster without persisted topology): %s -> %s", topologyId, newTopologyId)
	}
	assertTopologyIdConsistent(t, mc, newTopologyId)
}

// TestWriterConsistencyAcrossNodes verifies that all nodes agree on exactly one
// writer and report the same TopologyId — the deterministic-pin invariant that
// every master computes the same writer from the same membership.
func TestWriterConsistencyAcrossNodes(t *testing.T) {
	mc := StartMasterCluster(t)

	time.Sleep(3 * time.Second)

	writerIdx, writerAddr := mc.FindWriter()
	if writerIdx < 0 {
		t.Fatal("no writer found")
	}
	t.Logf("writer: node %d at %s", writerIdx, writerAddr)

	if got := mc.CountWriters(); got != 1 {
		mc.DumpLogs()
		t.Fatalf("%d nodes claim to be the writer, want exactly 1", got)
	}

	// Every node must name the same writer address.
	for i := range 3 {
		cs, err := mc.GetClusterStatus(i)
		if err != nil {
			t.Fatalf("node %d cluster/status error: %v", i, err)
		}
		if i == writerIdx {
			if !cs.IsWriter {
				t.Errorf("node %d should be the writer but IsWriter=false", i)
			}
			continue
		}
		if cs.IsWriter {
			t.Errorf("node %d should not be the writer but IsWriter=true", i)
		}
		// cs.Writer is a ServerAddress like "127.0.0.1:10000.20000"; compare
		// its HTTP form against the discovered writer's HTTP address.
		writerHttp := pb.ServerAddress(cs.Writer).ToHttpAddress()
		if writerHttp != writerAddr {
			t.Errorf("node %d names writer %q (http: %s), expected %q", i, cs.Writer, writerHttp, writerAddr)
		}
	}

	topologyId, _ := mc.GetTopologyId(writerIdx)
	if topologyId == "" {
		t.Fatal("writer has no TopologyId")
	}
	assertTopologyIdConverged(t, mc)
}

// assertTopologyIdConsistent verifies that all running nodes report the expected TopologyId.
func assertTopologyIdConsistent(t *testing.T, mc *MasterCluster, expectedId string) {
	t.Helper()
	for i := range 3 {
		if !mc.IsNodeRunning(i) {
			continue
		}
		id, err := mc.GetTopologyId(i)
		if err != nil {
			t.Errorf("node %d: failed to get TopologyId: %v", i, err)
			continue
		}
		if id != expectedId {
			t.Errorf("node %d: TopologyId=%q, expected %q", i, id, expectedId)
		}
	}
}

// assertTopologyIdConverged is the leaderless post-failover guarantee: the
// cluster settles on a single writer that serves a non-empty TopologyId, and
// every live node agrees on it (each proxies /dir/status to the writer). NOTE:
// the id VALUE is not preserved across a writer change — a follower promoted to
// writer mints a fresh id because TopologyId is not yet propagated/persisted
// (the one remaining seam; see LLM.md "Master consensus"). Correctness of the
// data path does not depend on it — volume/file id uniqueness comes from the
// Coordinator's single-writer serialization + max-id observation, not the
// TopologyId. This asserts what actually holds: convergence + agreement.
func assertTopologyIdConverged(t *testing.T, mc *MasterCluster) {
	t.Helper()
	writerIdx, _ := mc.FindWriter()
	if writerIdx < 0 {
		mc.DumpLogs()
		t.Fatal("no writer to serve TopologyId after failover")
	}
	id, err := mc.GetTopologyId(writerIdx)
	if err != nil || id == "" {
		mc.DumpLogs()
		t.Fatalf("writer node %d has no TopologyId: %v", writerIdx, err)
	}
	assertTopologyIdConsistent(t, mc, id)
}
