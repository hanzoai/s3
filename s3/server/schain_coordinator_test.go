package s3server

import (
	"context"
	"errors"
	"testing"

	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/sequence"
	"github.com/hanzoai/s3/s3/storage/needle"
	"github.com/hanzoai/s3/s3/topology"
)

// The schain-backed Coordinator is the production durable backend: it composes
// the deterministic pinned-writer LocalCoordinator (writer selection, membership)
// and routes only id allocation through the schain storage VM's owner-gated
// AllocateTx over ZAP. These tests pin the two safety-critical properties:
//   1. UNWIRED, it fails CLOSED — never invents an id the chain has not agreed.
//   2. WIRED (fake chain allocator injected), it returns chain-agreed ids and
//      raises the local topology floor, and still refuses on a non-writer.

func newSchainTestTopology() *topology.Topology {
	return topology.NewTopology("test", sequence.NewMemorySequencer(), 1<<30, 5, false)
}

// fakeAllocator stands in for the schain VM: it hands out a contiguous range
// starting at a deterministic per-kind cursor, exactly as the chain's AllocateTx
// does after block Accept.
type fakeAllocator struct {
	next map[string]uint64
}

func (f *fakeAllocator) Allocate(_ context.Context, kind string, count uint64) (uint64, error) {
	if f.next == nil {
		f.next = map[string]uint64{}
	}
	if f.next[kind] == 0 {
		f.next[kind] = 1
	}
	start := f.next[kind]
	f.next[kind] += count
	return start, nil
}

func TestSchainCoordinator_FailsClosedWhenUnwired(t *testing.T) {
	topo := newSchainTestTopology()
	// Empty endpoint ⇒ the ZAP allocator has nowhere to dial: it must fail closed.
	c := NewSchainCoordinator(topo, "", pb.DialOption{})

	// Standalone ⇒ this node IS the pinned writer, so it reaches the allocator.
	if !c.IsWriter() {
		t.Fatal("standalone schain coordinator must be the writer")
	}
	if _, _, err := c.AllocateVolumeId(1); !errors.Is(err, ErrSchainNotWired) {
		t.Fatalf("AllocateVolumeId unwired = %v, want ErrSchainNotWired", err)
	}
	if _, err := c.NextFileId(1); !errors.Is(err, ErrSchainNotWired) {
		t.Fatalf("NextFileId unwired = %v, want ErrSchainNotWired", err)
	}
	// The floor must NOT have moved — nothing was allocated.
	if got := topo.GetMaxVolumeId(); got != 0 {
		t.Fatalf("MaxVolumeId advanced to %d on a failed-closed allocation, want 0", got)
	}
}

func TestSchainCoordinator_WiredDelegatesToChain(t *testing.T) {
	topo := newSchainTestTopology()
	c := NewSchainCoordinator(topo, "schain.local:9000", pb.DialOption{})

	// Inject a fake chain allocator (the real one dials the schain ZAP endpoint).
	sc, ok := c.(*schainCoordinator)
	if !ok {
		t.Fatalf("NewSchainCoordinator returned %T, want *schainCoordinator", c)
	}
	sc.allocator = &fakeAllocator{}

	start, end, err := c.AllocateVolumeId(4)
	if err != nil {
		t.Fatalf("wired AllocateVolumeId: %v", err)
	}
	if start != needle.VolumeId(1) || end != needle.VolumeId(4) {
		t.Fatalf("volume range = [%d,%d], want [1,4]", start, end)
	}
	// The chain-agreed range must have raised the local floor for reads.
	if got := topo.GetMaxVolumeId(); got != 4 {
		t.Fatalf("MaxVolumeId = %d after chain allocation, want 4", got)
	}

	fid, err := c.NextFileId(2)
	if err != nil {
		t.Fatalf("wired NextFileId: %v", err)
	}
	if fid != 1 {
		t.Fatalf("first file id = %d, want 1 (independent per-kind chain cursor)", fid)
	}
}

func TestSchainCoordinator_NonWriterRefusesBeforeChain(t *testing.T) {
	topo := newSchainTestTopology()
	c := NewSchainCoordinator(topo, "schain.local:9000", pb.DialOption{})
	sc := c.(*schainCoordinator)
	// A fake that fails the test if ever called — a non-writer must short-circuit
	// with ErrNotWriter and NEVER reach the chain.
	sc.allocator = allocatorFunc(func(context.Context, string, uint64) (uint64, error) {
		t.Fatal("non-writer must not call the schain allocator")
		return 0, nil
	})

	// Pin the writer elsewhere: this node (300) defers to 100.
	c.SetMembers(pb.ServerAddress("127.0.0.1:300"),
		[]pb.ServerAddress{"127.0.0.1:100", "127.0.0.1:200", "127.0.0.1:300"})
	if c.IsWriter() {
		t.Fatal("member 300 must defer to writer 100")
	}
	if _, _, err := c.AllocateVolumeId(1); err != topology.ErrNotWriter {
		t.Fatalf("AllocateVolumeId on non-writer = %v, want ErrNotWriter", err)
	}
	if _, err := c.NextFileId(1); err != topology.ErrNotWriter {
		t.Fatalf("NextFileId on non-writer = %v, want ErrNotWriter", err)
	}
}

type allocatorFunc func(context.Context, string, uint64) (uint64, error)

func (f allocatorFunc) Allocate(ctx context.Context, kind string, count uint64) (uint64, error) {
	return f(ctx, kind, count)
}
