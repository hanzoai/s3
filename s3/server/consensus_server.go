package s3server

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/luxfi/consensus"
	"github.com/luxfi/consensus/replog"

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
type ConsensusServer struct {
	serverAddr pb.ServerAddress
	topo       *topology.Topology

	chain *consensus.Chain
	log   *replog.Log

	peersMu sync.RWMutex
	peers   map[string]pb.ServerAddress

	ctx    context.Context
	cancel context.CancelFunc
}

// NewConsensusServer starts the master coordination on Lux consensus. peers is
// the initial validator set (master peer name -> address).
func NewConsensusServer(serverAddr pb.ServerAddress, peers map[string]pb.ServerAddress, topo *topology.Topology) (*ConsensusServer, error) {
	if peers == nil {
		peers = make(map[string]pb.ServerAddress)
	}
	chain := consensus.NewChain(consensus.DefaultConfig())
	cs := &ConsensusServer{
		serverAddr: serverAddr,
		topo:       topo,
		chain:      chain,
		peers:      peers,
	}
	cs.log = replog.New(chain, cs.applyMaxVolumeId)
	cs.ctx, cs.cancel = context.WithCancel(context.Background())
	if err := cs.log.Start(cs.ctx); err != nil {
		return nil, fmt.Errorf("start consensus: %w", err)
	}
	// Drive the apply loop: finalized commands are applied to the topology.
	go cs.log.Run(cs.ctx, 20*time.Millisecond)
	glog.V(0).Infof("master %s coordinating on Lux consensus (Raft retired)", serverAddr)
	return cs, nil
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
	return nil
}

// Do replicates a MaxVolumeIdCommand through consensus and blocks until it is
// finalized and applied on this replica. The drop-in for raft Server.Do.
func (cs *ConsensusServer) Do(cmd *topology.MaxVolumeIdCommand) error {
	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return cs.log.Commit(cs.ctx, payload)
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

// IsLogEmpty reports whether the replicated log has no committed entries.
// A fresh consensus engine starts empty; once a command is applied it is not.
func (cs *ConsensusServer) IsLogEmpty() bool { return cs.log.Pending() == 0 }

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
