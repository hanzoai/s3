package s3server

import (
	"context"
	"errors"
	"fmt"

	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/storage/needle"
	"github.com/hanzoai/s3/s3/topology"
)

// ErrSchainNotWired is returned by the schain-backed Coordinator's allocation
// path when no live schain allocator has been installed. It fails CLOSED — an
// allocation that cannot be ordered by the chain is refused, never served from
// a local guess that could collide with another writer.
var ErrSchainNotWired = errors.New("schain coordinator: allocator endpoint not wired")

// SchainAllocator is the single integration seam to the schain storage VM. The
// VM serves an owner-gated, ML-DSA-verified AllocateTx that hands out a
// contiguous, globally-unique id range per (kind) — "volume" or "file". It is
// reached over the native ZAP transport (hanzo/s3 and luxfi/chains are separate
// modules by the org-separation rule, so this is a thin ZAP client to the
// chain's RPC, never a Go import of the VM).
//
// Wiring point: provide an implementation that dials the configured schain ZAP
// endpoint, submits an AllocateTx for `count` ids of `kind`, blocks on the
// cert-gated block Accept, and returns the first id of the accepted range. Until
// that implementation is installed, AllocateVolumeId / NextFileId fail closed
// with ErrSchainNotWired — the master must not invent ids the chain has not
// agreed.
type SchainAllocator interface {
	Allocate(ctx context.Context, kind string, count uint64) (start uint64, err error)
}

// schainCoordinator is the production Coordinator: durable, range-parallel
// (single-writer PER RANGE, not global), ML-DSA owner-verified. It composes the
// deterministic pinned-writer LocalCoordinator (writer selection, membership,
// max-id observation, topology id — all identical to the in-process backend)
// and overrides only the two allocation methods to route through the schain VM.
type schainCoordinator struct {
	*topology.LocalCoordinator
	topo      *topology.Topology
	endpoint  string
	allocator SchainAllocator
}

var _ topology.Coordinator = (*schainCoordinator)(nil)

// NewSchainCoordinator builds the schain-backed Coordinator for topo, pointed at
// the schain VM ZAP endpoint. dialOption carries the PQ-mTLS the master uses for
// its other ZAP links. The returned coordinator pins the writer locally and
// delegates id allocation to the chain.
func NewSchainCoordinator(topo *topology.Topology, endpoint string, dialOption pb.DialOption) topology.Coordinator {
	return &schainCoordinator{
		LocalCoordinator: topology.NewLocalCoordinator(topo),
		topo:             topo,
		endpoint:         endpoint,
		allocator:        newZapSchainAllocator(endpoint, dialOption),
	}
}

// AllocateVolumeId reserves `count` volume ids through the chain's AllocateTx,
// then raises the local topology floor so reads see the new max. Only the pinned
// writer proposes; non-writers redirect.
func (c *schainCoordinator) AllocateVolumeId(count uint64) (start, end needle.VolumeId, err error) {
	if !c.IsWriter() {
		return 0, 0, topology.ErrNotWriter
	}
	if count == 0 {
		count = 1
	}
	s, err := c.allocator.Allocate(context.Background(), "volume", count)
	if err != nil {
		return 0, 0, fmt.Errorf("schain allocate volume: %w", err)
	}
	start = needle.VolumeId(s)
	end = needle.VolumeId(s + count - 1)
	c.topo.UpAdjustMaxVolumeId(end)
	return start, end, nil
}

// NextFileId reserves `count` file ids through the chain's AllocateTx and returns
// the first id of the accepted range.
func (c *schainCoordinator) NextFileId(count uint64) (uint64, error) {
	if !c.IsWriter() {
		return 0, topology.ErrNotWriter
	}
	if count == 0 {
		count = 1
	}
	start, err := c.allocator.Allocate(context.Background(), "file", count)
	if err != nil {
		return 0, fmt.Errorf("schain allocate file: %w", err)
	}
	return start, nil
}

// zapSchainAllocator is the ZAP client to the schain VM. It holds the endpoint
// and the PQ-mTLS dial options shared with the master's other ZAP links. The
// request/response framing of the schain AllocateTx is the integration seam
// (see SchainAllocator): until it is wired against the chain's ZAP service
// descriptor, Allocate fails closed.
type zapSchainAllocator struct {
	endpoint   string
	dialOption pb.DialOption
}

func newZapSchainAllocator(endpoint string, dialOption pb.DialOption) SchainAllocator {
	return &zapSchainAllocator{endpoint: endpoint, dialOption: dialOption}
}

func (a *zapSchainAllocator) Allocate(ctx context.Context, kind string, count uint64) (uint64, error) {
	if a.endpoint == "" {
		return 0, ErrSchainNotWired
	}
	// Integration seam: dial a.endpoint over the ZAP transport, submit an
	// AllocateTx{kind, count} signed by this node's ML-DSA staking key, block on
	// the cert-gated Accept, and return the first id of the accepted range. The
	// schain ZAP service descriptor lives in the separate luxfi/chains module;
	// install a wired SchainAllocator on the coordinator to enable it.
	return 0, ErrSchainNotWired
}
