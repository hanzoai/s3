package s3server

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/luxfi/consensus"
	"github.com/luxfi/consensus/replog"
	"github.com/luxfi/ids"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/topology"
)

// ConsensusServer is the master's HA coordination on Lux consensus, replacing
// seaweedfs/raft (and hashicorp/raft). The master's entire replicated state
// machine is a monotonic max-volume-id plus the topology id; here those
// commands are committed through a linear consensus.Chain via replog instead
// of a Raft log. Versus Raft this is ZAP-native (block/vote gossip rides the
// zap-proto transport, no gRPC), post-quantum-final (Quasar BLS + ML-DSA),
// and needs no separate leader-election FSM.
//
// FINALITY MODEL: consensus.Chain finalizes a block when it has collected
// Config.Alpha votes (RecordVote). The accept quorum is sized to the membership
// in NewConsensusServer:
//   - SINGLE MASTER (the production default, -peers=none): quorum is 1, so the
//     node's own vote (cast in Do) finalizes every command immediately. Fully
//     functional — this is legitimate single-validator consensus.
//   - MULTI-MASTER (N>1): quorum is the BFT supermajority of N. Do casts this
//     node's one vote; the remaining votes must arrive from peer masters via
//     RecordPeerVote. GAP: the gossip that exchanges these votes over the ZAP
//     master transport is NOT yet wired here, so a multi-master set cannot
//     finalize cross-node commands until that lands. This is the live
//     consensus-engine swap's remaining multi-node work; the single-master path
//     (and all unit/integration coverage of it) is complete.
type ConsensusServer struct {
	serverAddr pb.ServerAddress
	topo       *topology.Topology

	chain *consensus.Chain
	log   *replog.Log

	// selfVoter is this node's stable consensus voter identity (derived from its
	// address), used to cast its own vote on blocks it proposes.
	selfVoter consensus.NodeID

	peersMu sync.RWMutex
	peers   map[string]pb.ServerAddress

	ctx    context.Context
	cancel context.CancelFunc

	// applied is set the first time a finalized command reaches the FSM. It is
	// the basis of IsLogEmpty: replog.Pending() is only the un-applied queue
	// depth, which drains back to 0, so it cannot distinguish a fresh log from
	// one that has already committed and applied entries.
	applied atomic.Bool
}

// NewConsensusServer starts the master coordination on Lux consensus. peers is
// the initial validator set (master peer name -> address) — INCLUDING self.
//
// The engine's accept quorum (Config.Alpha) is sized to the membership so the
// log finalizes correctly: for a single master the quorum is 1, so this node's
// own vote finalizes its proposals immediately (legitimate single-validator
// consensus); for N masters the quorum is the BFT supermajority of N.
func NewConsensusServer(serverAddr pb.ServerAddress, peers map[string]pb.ServerAddress, topo *topology.Topology) (*ConsensusServer, error) {
	if peers == nil {
		peers = make(map[string]pb.ServerAddress)
	}
	cfg := consensus.DefaultConfig()
	cfg.Alpha = acceptQuorum(len(peers))
	cfg.K = len(peers)
	if cfg.K < 1 {
		cfg.K = 1
	}
	chain := consensus.NewChain(cfg)
	cs := &ConsensusServer{
		serverAddr: serverAddr,
		topo:       topo,
		chain:      chain,
		selfVoter:  voterID(serverAddr),
		peers:      peers,
	}
	cs.log = replog.New(chain, cs.applyMaxVolumeId)
	cs.ctx, cs.cancel = context.WithCancel(context.Background())
	if err := cs.log.Start(cs.ctx); err != nil {
		return nil, fmt.Errorf("start consensus: %w", err)
	}
	// Drive the apply loop: finalized commands are applied to the topology.
	go cs.log.Run(cs.ctx, 20*time.Millisecond)
	glog.V(0).Infof("master %s coordinating on Lux consensus (Raft retired), accept quorum=%d of %d", serverAddr, cfg.Alpha, cfg.K)
	return cs, nil
}

// acceptQuorum is the number of votes needed to finalize a block for a
// membership of n voting masters: 1 for a single master (it finalizes its own
// proposals), otherwise the BFT supermajority ceil(2n/3)+1 capped at n.
func acceptQuorum(n int) int {
	if n <= 1 {
		return 1
	}
	q := (2*n)/3 + 1
	if q > n {
		q = n
	}
	return q
}

// voterID derives a stable consensus voter identity from a master's address, so
// the same master always casts votes under the same NodeID.
func voterID(addr pb.ServerAddress) consensus.NodeID {
	sum := sha256.Sum256([]byte(addr))
	id, _ := ids.ToNodeID(sum[:ids.NodeIDLen])
	return id
}

// applyMaxVolumeId is the master's entire replicated FSM: apply a finalized
// MaxVolumeIdCommand to the in-memory topology, on every replica, in order.
func (cs *ConsensusServer) applyMaxVolumeId(payload []byte) error {
	var cmd topology.MaxVolumeIdCommand
	if err := json.Unmarshal(payload, &cmd); err != nil {
		return err
	}
	cs.topo.UpAdjustMaxVolumeId(cmd.MaxVolumeId)
	if cmd.TopologyId != "" {
		cs.topo.SetTopologyId(cmd.TopologyId)
	}
	cs.applied.Store(true)
	return nil
}

// Do replicates a MaxVolumeIdCommand through consensus and blocks until it is
// finalized and applied on this replica. The drop-in for raft Server.Do.
//
// Flow: propose the command (replog.Submit), cast THIS node's vote on it, then
// drain the applied prefix (replog.Advance). For a single master (accept
// quorum 1) the self-vote finalizes the block immediately and Advance applies
// it before Do returns. For a multi-master set the self-vote is one of the
// supermajority needed; the rest arrive as peer votes — see RecordPeerVote and
// the gap note on the type. The shared background Run loop also drains, so a
// command finalized by later peer votes still applies; Do additionally drains
// synchronously so the single-master path is immediate and deterministic.
func (cs *ConsensusServer) Do(cmd *topology.MaxVolumeIdCommand) error {
	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	id, err := cs.log.Submit(cs.ctx, payload)
	if err != nil {
		return fmt.Errorf("submit consensus command: %w", err)
	}
	// Cast this node's vote on the block it just proposed.
	if err := cs.chain.RecordVote(cs.ctx, consensus.NewVote(id, consensus.VoteCommit, cs.selfVoter)); err != nil {
		return fmt.Errorf("self-vote consensus command: %w", err)
	}
	// Block until this command is finalized AND applied. Advance applies the
	// contiguous accepted prefix in order; commands are voted in submit order, so
	// once THIS block is StatusAccepted the prefix up to and including it has been
	// drained by an Advance — its payload is in the FSM. We gate on this block's
	// own status (not the aggregate Pending) so concurrent Do calls on this node
	// don't block each other. Both this loop and the background Run loop call
	// Advance; Advance is mutex-guarded, so the double drive is safe.
	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()
	for {
		if _, err := cs.log.Advance(cs.ctx); err != nil {
			return fmt.Errorf("apply consensus command: %w", err)
		}
		if cs.chain.GetStatus(id) == consensus.StatusAccepted {
			return nil
		}
		select {
		case <-cs.ctx.Done():
			return cs.ctx.Err()
		case <-ticker.C:
		}
	}
}

// RecordPeerVote records a vote cast by another master on a proposed block,
// contributing toward the accept quorum. In multi-master mode these votes ride
// the ZAP transport (the master service); the gossip wiring that delivers them
// is the documented remaining work for multi-master finality (see the type
// doc). Single-master never needs this — its own vote is the quorum.
//
// MULTI-MASTER GAP — two things the gossip wiring MUST add before N>1 is safe
// (neither is reachable on the shipped single-master path, where Do casts
// exactly one self-vote into an Alpha=1 chain):
//  1. Dedup by voter. consensus.Chain.RecordVote appends to the block's vote
//     slice and finalizes at len(votes) >= Alpha with NO dedup, so a
//     retransmitted peer vote would inflate the count and could finalize below
//     true quorum. The gossip layer (or this method) must drop duplicate
//     (blockID, voter) pairs.
//  2. BFT-floor-valid (K, Alpha). acceptQuorum(n) returns ceil(2n/3)+1, but the
//     consensus config's own validity floor (2*Alpha - K >= floor((K-1)/3)+1)
//     must hold for the chosen (K=n, Alpha) before the engine is trusted for
//     Byzantine safety — assert cfg.Valid() (or the BFT floor directly) when the
//     real validator set is wired, especially for small n in {2,3}.
func (cs *ConsensusServer) RecordPeerVote(blockID consensus.ID, voter consensus.NodeID) error {
	return cs.chain.RecordVote(cs.ctx, consensus.NewVote(blockID, consensus.VoteCommit, voter))
}

// Peers returns the known peer masters, for observability only. Lux
// consensus is PERMISSIONLESS — participation is not gated by this list.
func (cs *ConsensusServer) Peers() (members []string) {
	cs.peersMu.RLock()
	defer cs.peersMu.RUnlock()
	for name := range cs.peers {
		members = append(members, name)
	}
	return
}

// AddPeer / RemovePeer are ADVISORY. Lux consensus is permissionless:
// participation is open and not admitted/evicted through a fixed validator
// roster, so there is no membership-reconfiguration round (a Raft cost these
// methods used to incur). They only update the peer list Peers reports.
func (cs *ConsensusServer) AddPeer(name string, addr pb.ServerAddress) {
	cs.peersMu.Lock()
	cs.peers[name] = addr
	cs.peersMu.Unlock()
}

func (cs *ConsensusServer) RemovePeer(name string) {
	cs.peersMu.Lock()
	delete(cs.peers, name)
	cs.peersMu.Unlock()
}

// Name is this master's consensus node identity (its server address).
func (cs *ConsensusServer) Name() string { return string(cs.serverAddr) }

// --- raft.Server drop-in shims (so the master's call sites need no rewrite
// beyond the field type). Leader-election and log persistence are the
// consensus engine's job, so the Raft-era tuning knobs are no-ops. ---

// Leader returns the PINNED WRITER — the one master permitted to run the
// serialization-critical metadata ops the master gates behind IsLeader()
// (volume growth/assignment, vacuum/GC, sequence allocation). Lux consensus
// itself is LEADERLESS: any replica may propose and consensus totally-orders
// the replog, and the monotonic FSM (UpAdjustMaxVolumeId) reconciles
// *commutative* commands by order. But the gated ops are NOT commutative — N
// masters growing volumes, assigning IDs, and vacuuming at once would
// over-allocate, race volume-IDs, and double-GC. So the writer is an
// app-level role, orthogonal to the leaderless log, pinned deterministically
// to the lowest member address in the (consensus-agreed) membership: every
// node computes the same writer with no election and no leader-forwarding,
// and when the writer leaves the next-lowest takes over.
func (cs *ConsensusServer) Leader() string {
	cs.peersMu.RLock()
	defer cs.peersMu.RUnlock()
	writer := cs.Name()
	for _, addr := range cs.peers {
		if a := string(addr); a < writer {
			writer = a
		}
	}
	return writer
}

// Start is a no-op: the engine is already started by NewConsensusServer.
func (cs *ConsensusServer) Start() error { return nil }

// IsLogEmpty reports whether this replica has applied any finalized command —
// the drop-in for raft Server.IsLogEmpty (gates fast-resume and bootstrap). It
// is NOT Pending(): that is only the un-applied queue depth, which drains back
// to 0, so it would read "empty" even after the log has committed and applied
// entries (silently disabling fastResume at raft_server.go:185).
func (cs *ConsensusServer) IsLogEmpty() bool { return !cs.applied.Load() }

// LoadSnapshot is a no-op: consensus finality replaces Raft snapshotting.
func (cs *ConsensusServer) LoadSnapshot() error { return nil }

// SetElectionTimeout / SetHeartbeatInterval are Raft tuning knobs with no
// analogue in Lux consensus; kept as no-ops for call-site compatibility.
func (cs *ConsensusServer) SetElectionTimeout(time.Duration)   {}
func (cs *ConsensusServer) SetHeartbeatInterval(time.Duration) {}

// Stop tears down the consensus engine.
func (cs *ConsensusServer) Stop() {
	if cs.cancel != nil {
		cs.cancel()
	}
	if cs.log != nil {
		_ = cs.log.Stop()
	}
}
