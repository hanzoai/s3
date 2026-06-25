package topology

import "errors"

// Consensus is the master's HA coordination engine as the topology sees it.
//
// It replaces the seaweedfs/raft (and hashicorp/raft) Server interface the
// topology used to hold. The implementation is s3server.ConsensusServer, which
// coordinates on the Lux consensus replog (ZAP-native, post-quantum) instead of
// a Raft log. The topology only needs these four operations; the master's
// membership/observability surface (peers, join, status handlers) is on the
// concrete *ConsensusServer in package s3server.
//
// Leaderless with a pinned writer: the engine totally-orders the replog with no
// election, and the non-commutative metadata ops the master gates behind
// IsLeader() (volume growth/assignment, vacuum/GC, sequence allocation) run only
// on the single deterministically-pinned writer. IsLeader() reports whether THIS
// node is that writer.
type Consensus interface {
	// IsLeader reports whether this node is the pinned writer (the one node
	// permitted to run the serialization-critical metadata ops).
	IsLeader() bool
	// Leader returns the pinned writer's node name (address).
	Leader() string
	// Name is this node's consensus identity (its server address).
	Name() string
	// Do replicates a MaxVolumeIdCommand through consensus and blocks until it is
	// finalized and applied on this node.
	Do(cmd *MaxVolumeIdCommand) error
}

// ErrNotLeader is returned by master RPCs that must run on the pinned writer when
// this node is not it. It replaces seaweedfs/raft's NotLeaderError so the master
// service carries no gRPC/raft dependency.
var ErrNotLeader = errors.New("not the leader (pinned writer)")
