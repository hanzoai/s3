package topology

import (
	"sync"
	"testing"

	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/sequence"
	"github.com/hanzoai/s3/s3/storage/needle"
)

func addr(s string) pb.ServerAddress { return pb.ServerAddress(s) }

func addrs(ss []string) []pb.ServerAddress {
	out := make([]pb.ServerAddress, len(ss))
	for i, s := range ss {
		out[i] = pb.ServerAddress(s)
	}
	return out
}

// This suite is the luxfi-native replacement for the deleted raft leader/log
// tests. It pins the correctness invariants the raft path used to guarantee —
// single-writer serialization of the two non-commutative allocations (volume
// ids, file ids), a globally-agreed writer, and clean failover — but now over
// the leaderless deterministic pin, with NO election.

func newTestCoordinator(t *testing.T) (*Topology, *LocalCoordinator) {
	t.Helper()
	topo := NewTopology("test", sequence.NewMemorySequencer(), 1<<30, 5, false)
	// NewTopology already installs a LocalCoordinator; use it so topo.Coordinator
	// and the returned coordinator are the same object.
	c, ok := topo.Coordinator.(*LocalCoordinator)
	if !ok {
		t.Fatalf("NewTopology installed %T, want *LocalCoordinator", topo.Coordinator)
	}
	return topo, c
}

// A standalone master (no membership configured) is the sole writer and
// allocates freely — the single-node / dev / test path.
func TestLocalCoordinator_StandaloneIsSoleWriter(t *testing.T) {
	_, c := newTestCoordinator(t)

	if !c.IsWriter() {
		t.Fatal("standalone coordinator must be the writer")
	}
	start, end, err := c.AllocateVolumeId(1)
	if err != nil {
		t.Fatalf("standalone AllocateVolumeId: %v", err)
	}
	if start != 1 || end != 1 {
		t.Fatalf("first volume id range = [%d,%d], want [1,1]", start, end)
	}
	if _, err := c.NextFileId(1); err != nil {
		t.Fatalf("standalone NextFileId: %v", err)
	}
}

// The writer is the lowest-address member — a pure function of the membership,
// so it is deterministic with zero rounds of election. Every master that holds
// the same set computes the same writer.
func TestLocalCoordinator_DeterministicPinLowestAddress(t *testing.T) {
	members := []string{"127.0.0.1:300", "127.0.0.1:100", "127.0.0.1:200"}

	// Every member computes the SAME writer from the SAME set, independently.
	for _, self := range members {
		_, c := newTestCoordinator(t)
		c.SetMembers(addr(self), addrs(members))
		if got := string(c.Writer()); got != "127.0.0.1:100" {
			t.Fatalf("member %s computed writer %q, want lowest 127.0.0.1:100", self, got)
		}
		wantWriter := self == "127.0.0.1:100"
		if c.IsWriter() != wantWriter {
			t.Fatalf("member %s IsWriter=%v, want %v", self, c.IsWriter(), wantWriter)
		}
	}
}

// A non-writer refuses the serialization-critical allocations with ErrNotWriter
// — the same contract raft gave with NotLeaderError. Clients back off and
// redirect to the advertised writer.
func TestLocalCoordinator_NonWriterRefusesAllocation(t *testing.T) {
	_, c := newTestCoordinator(t)
	c.SetMembers(addr("127.0.0.1:300"), addrs([]string{"127.0.0.1:100", "127.0.0.1:200", "127.0.0.1:300"}))

	if c.IsWriter() {
		t.Fatal("member 300 must not be the writer (100 is lowest)")
	}
	if _, _, err := c.AllocateVolumeId(1); err != ErrNotWriter {
		t.Fatalf("AllocateVolumeId on non-writer = %v, want ErrNotWriter", err)
	}
	if _, err := c.NextFileId(1); err != ErrNotWriter {
		t.Fatalf("NextFileId on non-writer = %v, want ErrNotWriter", err)
	}
	// A non-writer never mints the topology id.
	if id, err := c.EnsureTopologyId(); err != nil || id != "" {
		t.Fatalf("EnsureTopologyId on non-writer = (%q,%v), want (\"\",nil)", id, err)
	}
}

// Failover with no election: when the pinned writer leaves the membership, the
// next-lowest member becomes the writer the instant the set changes — a pure
// re-computation, no votes, no timeout. This is the leaderless analogue of raft
// leader election.
func TestLocalCoordinator_FailoverRepinsNextLowest(t *testing.T) {
	_, c := newTestCoordinator(t)
	full := []string{"127.0.0.1:100", "127.0.0.1:200", "127.0.0.1:300"}

	// This master is 200 — defers to 100 while 100 is present.
	c.SetMembers(addr("127.0.0.1:200"), addrs(full))
	if c.IsWriter() {
		t.Fatal("200 must defer to writer 100")
	}

	// 100 leaves. 200 is now the lowest address → it re-pins to writer instantly.
	c.SetMembers(addr("127.0.0.1:200"), addrs([]string{"127.0.0.1:200", "127.0.0.1:300"}))
	if !c.IsWriter() {
		t.Fatal("after 100 leaves, 200 must become the writer")
	}
	if got := string(c.Writer()); got != "127.0.0.1:200" {
		t.Fatalf("writer after failover = %q, want 127.0.0.1:200", got)
	}
	// The freshly-pinned writer allocates without any warmup round.
	if _, _, err := c.AllocateVolumeId(1); err != nil {
		t.Fatalf("newly-pinned writer AllocateVolumeId: %v", err)
	}
}

// Volume-id allocation is monotonic and non-overlapping across ranges, and it
// advances the durable topology floor.
func TestLocalCoordinator_AllocateVolumeIdMonotonic(t *testing.T) {
	topo, c := newTestCoordinator(t)

	s1, e1, err := c.AllocateVolumeId(5)
	if err != nil {
		t.Fatal(err)
	}
	if s1 != 1 || e1 != 5 {
		t.Fatalf("range 1 = [%d,%d], want [1,5]", s1, e1)
	}
	s2, e2, err := c.AllocateVolumeId(3)
	if err != nil {
		t.Fatal(err)
	}
	if s2 != 6 || e2 != 8 {
		t.Fatalf("range 2 = [%d,%d], want [6,8]", s2, e2)
	}
	if got := topo.GetMaxVolumeId(); got != 8 {
		t.Fatalf("topology MaxVolumeId = %d, want 8", got)
	}
}

// The writer serializes concurrent allocations: every id is globally unique,
// even under heavy contention. This is the guarantee that made raft's
// serialized log necessary; the coordinator's lock provides it directly.
func TestLocalCoordinator_ConcurrentAllocationUnique(t *testing.T) {
	_, c := newTestCoordinator(t)

	const n = 500
	var mu sync.Mutex
	seen := make(map[needle.VolumeId]struct{}, n)
	fileIds := make(map[uint64]struct{}, n)

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			vid, _, err := c.AllocateVolumeId(1)
			if err != nil {
				t.Errorf("AllocateVolumeId: %v", err)
				return
			}
			fid, err := c.NextFileId(1)
			if err != nil {
				t.Errorf("NextFileId: %v", err)
				return
			}
			mu.Lock()
			seen[vid] = struct{}{}
			fileIds[fid] = struct{}{}
			mu.Unlock()
		}()
	}
	wg.Wait()

	if len(seen) != n {
		t.Fatalf("got %d unique volume ids, want %d — allocation overlapped", len(seen), n)
	}
	if len(fileIds) != n {
		t.Fatalf("got %d unique file ids, want %d — allocation overlapped", len(fileIds), n)
	}
}

// ObserveMaxVolumeId / ObserveMaxFileId raise the allocation floor from state
// seen out of band (heartbeats, recovery) so a freshly-pinned writer never
// re-issues an id already in use — the property that let us drop the raft log.
func TestLocalCoordinator_ObserveRaisesFloor(t *testing.T) {
	_, c := newTestCoordinator(t)

	c.ObserveMaxVolumeId(100)
	start, _, err := c.AllocateVolumeId(1)
	if err != nil {
		t.Fatal(err)
	}
	if start != 101 {
		t.Fatalf("volume id after observing 100 = %d, want 101", start)
	}

	c.ObserveMaxFileId(5000)
	fid, err := c.NextFileId(1)
	if err != nil {
		t.Fatal(err)
	}
	if fid <= 5000 {
		t.Fatalf("file id after observing 5000 = %d, want > 5000", fid)
	}
}

// EnsureTopologyId mints the cluster identity exactly once on the writer and is
// idempotent thereafter.
func TestLocalCoordinator_EnsureTopologyIdIdempotent(t *testing.T) {
	_, c := newTestCoordinator(t)

	id1, err := c.EnsureTopologyId()
	if err != nil || id1 == "" {
		t.Fatalf("EnsureTopologyId = (%q,%v), want a non-empty id", id1, err)
	}
	id2, err := c.EnsureTopologyId()
	if err != nil {
		t.Fatal(err)
	}
	if id1 != id2 {
		t.Fatalf("EnsureTopologyId not idempotent: %q then %q", id1, id2)
	}
}

// Membership is deduplicated, empty-skipped, and sorted, so the pin is stable
// regardless of the order or noise in the configured set.
func TestLocalCoordinator_MembershipDedupSortSkipEmpty(t *testing.T) {
	_, c := newTestCoordinator(t)
	c.SetMembers(addr("127.0.0.1:200"), addrs([]string{
		"127.0.0.1:300", "", "127.0.0.1:100", "127.0.0.1:200", "127.0.0.1:100", "",
	}))

	got := c.Members()
	want := []string{"127.0.0.1:100", "127.0.0.1:200", "127.0.0.1:300"}
	if len(got) != len(want) {
		t.Fatalf("Members() = %v, want %v", got, want)
	}
	for i := range want {
		if string(got[i]) != want[i] {
			t.Fatalf("Members()[%d] = %q, want %q (full: %v)", i, got[i], want[i], got)
		}
	}
	if string(c.Writer()) != "127.0.0.1:100" {
		t.Fatalf("writer = %q, want 127.0.0.1:100", c.Writer())
	}
}
