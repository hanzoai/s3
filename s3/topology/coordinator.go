package topology

import (
	"errors"
	"sort"
	"sync"

	"github.com/google/uuid"

	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/storage/needle"
)

// ErrNotWriter is returned by allocation and the metadata-mutating master RPCs
// when this master is NOT the pinned writer. It replaces raft.NotLeaderError:
// clients treat it identically — back off and redirect to the writer that the
// master advertises in the KeepConnected / Heartbeat responses (the Leader
// field, which now carries the pinned-writer address).
var ErrNotWriter = errors.New("master is not the pinned writer")

// Coordinator is the leaderless source of truth that replaces the raft
// leader + replicated log. It serializes the master's two non-commutative,
// correctness-critical allocations — volume ids and file ids — and pins
// exactly ONE writer per cluster deterministically, with no election.
//
// There is no raft, no isLeader, no serialized-writer FSM. Coordination is:
//   - allocation is the serialization point (AllocateVolumeId / NextFileId);
//     the coordinator — not an elected leader — guarantees global uniqueness.
//   - the writer is pinned by a pure function of the (agreed) membership, so
//     every master computes the same writer with zero rounds of election and
//     fails over to the next member the instant the writer leaves.
//
// Two implementations exist, selected by environment, never both:
//   - LocalCoordinator: in-process, deterministic pinned-writer over a static
//     membership. The default. A standalone master is the sole writer; a
//     multi-master set agrees on the lowest-address writer. Keeps the whole
//     suite green and a dev/standalone master fully functional.
//   - the schain-backed coordinator (s3server.schainCoordinator): delegates
//     allocation to the schain storage VM's owner-gated AllocateTx over the
//     native ZAP transport. The production backend: durable, range-parallel
//     (single-writer PER RANGE, not global), ML-DSA owner-verified.
type Coordinator interface {
	// AllocateVolumeId reserves `count` fresh, globally-unique, monotonically
	// increasing volume ids and returns the inclusive [start, end] range.
	// Replaces Topology.NextVolumeId + the raft MaxVolumeIdCommand. Returns
	// ErrNotWriter on a non-writer member.
	AllocateVolumeId(count uint64) (start, end needle.VolumeId, err error)

	// NextFileId reserves `count` fresh, globally-unique file ids and returns
	// the first id of the run. Replaces sequence.Sequencer.NextFileId funnelled
	// through the raft leader. Returns ErrNotWriter on a non-writer member.
	NextFileId(count uint64) (uint64, error)

	// ObserveMaxVolumeId / ObserveMaxFileId raise the allocation floor from
	// state seen out of band — volume-server heartbeats and recovery — so a
	// freshly-pinned writer never re-issues an id already in use.
	ObserveMaxVolumeId(needle.VolumeId)
	ObserveMaxFileId(uint64)

	// EnsureTopologyId returns the cluster identity, generating and recording
	// it on first call. Stable for the cluster's lifetime.
	EnsureTopologyId() (string, error)

	// IsWriter reports whether THIS master is the pinned writer — the single
	// member permitted to run the serialization-critical metadata ops the
	// master used to gate behind IsLeader(). Deterministic; no election.
	IsWriter() bool

	// Writer returns the pinned writer's address (for read redirect / proxy),
	// or empty if no membership is known yet. Members lists the known
	// coordination members for observability.
	Writer() pb.ServerAddress
	Members() []pb.ServerAddress

	// SetMembers updates the coordination membership. The pinned writer is a
	// pure function of this set, so every master that holds the same set
	// computes the same writer.
	SetMembers(self pb.ServerAddress, members []pb.ServerAddress)
}

// LocalCoordinator is the in-process, deterministic pinned-writer Coordinator.
// It is bound to its Topology: file ids come from the topology's Sequencer and
// volume ids advance the topology's monotonic MaxVolumeId, so a standalone
// master needs nothing external. Multi-master correctness rests on the pin:
// only the lowest-address member allocates; the others refuse with ErrNotWriter
// and clients redirect — exactly the single-writer guarantee raft gave, minus
// the election.
type LocalCoordinator struct {
	topo *Topology

	mu      sync.Mutex
	self    pb.ServerAddress
	members []pb.ServerAddress // sorted; empty ⇒ standalone (self is sole writer)
}

var _ Coordinator = (*LocalCoordinator)(nil)

// NewLocalCoordinator returns a coordinator bound to topo with no membership
// configured yet — a standalone master, the sole writer. command/master.go
// installs the real membership via SetMembers once peers are known.
func NewLocalCoordinator(topo *Topology) *LocalCoordinator {
	return &LocalCoordinator{topo: topo}
}

func (c *LocalCoordinator) SetMembers(self pb.ServerAddress, members []pb.ServerAddress) {
	sorted := make([]pb.ServerAddress, 0, len(members))
	seen := make(map[pb.ServerAddress]struct{}, len(members))
	for _, m := range members {
		if m == "" {
			continue
		}
		if _, ok := seen[m]; ok {
			continue
		}
		seen[m] = struct{}{}
		sorted = append(sorted, m)
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	c.mu.Lock()
	c.self = self
	c.members = sorted
	c.mu.Unlock()
}

// Writer is the lowest-address member — the deterministic pin. With no
// membership configured the writer is self (standalone). Pure function of the
// membership: every master that holds the same set returns the same writer.
func (c *LocalCoordinator) Writer() pb.ServerAddress {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.members) == 0 {
		return c.self
	}
	return c.members[0]
}

func (c *LocalCoordinator) Members() []pb.ServerAddress {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.members) == 0 {
		if c.self == "" {
			return nil
		}
		return []pb.ServerAddress{c.self}
	}
	out := make([]pb.ServerAddress, len(c.members))
	copy(out, c.members)
	return out
}

// IsWriter is true when this master is the pin. A coordinator with no
// membership configured (standalone / fresh test topology) is always the
// writer, so a single-node master and the whole test suite allocate freely.
func (c *LocalCoordinator) IsWriter() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.members) == 0 {
		return true
	}
	return c.self == c.members[0]
}

func (c *LocalCoordinator) AllocateVolumeId(count uint64) (start, end needle.VolumeId, err error) {
	if !c.IsWriter() {
		return 0, 0, ErrNotWriter
	}
	if count == 0 {
		count = 1
	}
	// Single-writer serialization: the lock makes the read-bump-publish atomic,
	// so concurrent allocations on the writer never overlap. The topology's
	// MaxVolumeId is the durable floor — raised here and by ObserveMaxVolumeId
	// from heartbeats — so a freshly-pinned writer continues above existing ids.
	c.mu.Lock()
	defer c.mu.Unlock()
	cur := c.topo.GetMaxVolumeId()
	start = cur + 1
	end = cur + needle.VolumeId(count)
	c.topo.UpAdjustMaxVolumeId(end)
	return start, end, nil
}

func (c *LocalCoordinator) NextFileId(count uint64) (uint64, error) {
	if !c.IsWriter() {
		return 0, ErrNotWriter
	}
	return c.topo.Sequence.NextFileId(count), nil
}

func (c *LocalCoordinator) ObserveMaxVolumeId(vid needle.VolumeId) {
	c.topo.UpAdjustMaxVolumeId(vid)
}

func (c *LocalCoordinator) ObserveMaxFileId(v uint64) {
	c.topo.Sequence.SetMax(v)
}

// EnsureTopologyId generates and records the cluster identity on first call.
// Only the writer mints it (the others learn it via heartbeat/observation);
// SetTopologyId is idempotent and panics on a split-brain mismatch.
func (c *LocalCoordinator) EnsureTopologyId() (string, error) {
	if id := c.topo.GetTopologyId(); id != "" {
		return id, nil
	}
	if !c.IsWriter() {
		return "", nil
	}
	id := uuid.New().String()
	c.topo.SetTopologyId(id)
	return c.topo.GetTopologyId(), nil
}
