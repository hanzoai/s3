package s3server

import (
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/topology"
)

// consensus_server_master.go is the ConsensusServer's master-integration surface:
// the methods the master server, topology, and /cluster HTTP handlers call. The
// core engine lives in consensus_server.go; this file faithfully mirrors the
// raft.Server / RaftServer semantics those call sites used to depend on, against
// the leaderless-with-pinned-writer Lux consensus engine.

// ConsensusServer satisfies topology.Consensus.
var _ topology.Consensus = (*ConsensusServer)(nil)

// raftServerID is the canonical consensus node id for a server address: its
// http address (host:httpPort), stable across grpc-port suffixes. The name is
// retained from the Raft era — it is now the consensus member id.
func raftServerID(server pb.ServerAddress) string {
	return server.ToHttpAddress()
}

// IsLeader reports whether THIS node is the pinned writer — the drop-in for the
// topology's old `RaftServer.State() == raft.Leader` gate. With Lux consensus the
// writer is pinned deterministically (lowest member address), so leadership is
// simply "am I the writer": Name() == Leader(). No election, no term.
func (cs *ConsensusServer) IsLeader() bool {
	return cs.Name() == cs.Leader()
}

// HasExistingState reports whether this node has already applied a finalized
// command — the drop-in for RaftServer.HasExistingState (gates whether startup
// needs an initial join). It mirrors raft's `!IsLogEmpty()`.
func (cs *ConsensusServer) HasExistingState() bool {
	return !cs.IsLogEmpty()
}

// DoJoinCommand initializes this node's participation in consensus — the drop-in
// for RaftServer.DoJoinCommand. seaweedfs/raft committed a DefaultJoinCommand to
// seed a fresh single-node log; Lux consensus is permissionless and needs no
// membership-admission round, so the equivalent "I am now participating and have
// state" is committing the current MaxVolumeId through the replog. That finalizes
// and applies on this node, which (a) makes HasExistingState()/IsLogEmpty()
// reflect a non-empty log exactly as the raft join did, and (b) publishes the
// current max-volume-id/topology-id into consensus.
func (cs *ConsensusServer) DoJoinCommand() {
	glog.V(0).Infof("master %s joining Lux consensus (Raft retired)", cs.serverAddr)
	cmd := topology.NewMaxVolumeIdCommand(cs.topo.GetMaxVolumeId(), cs.topo.GetTopologyId())
	if err := cs.Do(cmd); err != nil {
		glog.Errorf("fail to commit consensus join command: %v", err)
	}
}

// AddPeerAddress / RemovePeerByName are the membership operations the master's
// RaftAddServer/RaftRemoveServer RPCs map onto. Membership is advisory in
// permissionless consensus (it only changes the writer-pin set / observability),
// so these never block on a quorum round.
func (cs *ConsensusServer) AddPeerAddress(name string, addr pb.ServerAddress) {
	cs.AddPeer(name, addr)
}

func (cs *ConsensusServer) RemovePeerByName(name string) {
	cs.RemovePeer(name)
}

// AddPeerByGrpcAddress adds a peer from the master_pb RaftAddServer request shape
// (id = canonical http member id, grpcAddress = host:grpcPort). The id is the
// canonical member identity the writer-pin and observability use, so the peer is
// keyed and pinned by it.
func (cs *ConsensusServer) AddPeerByGrpcAddress(id, grpcAddress string) {
	cs.AddPeer(id, pb.ServerAddress(id))
}

// PeerAddresses returns the full membership (self + known peers) as name->address
// pairs, for the master's RaftListClusterServers RPC and the status UI.
func (cs *ConsensusServer) PeerAddresses() map[string]pb.ServerAddress {
	cs.peersMu.RLock()
	defer cs.peersMu.RUnlock()
	out := make(map[string]pb.ServerAddress, len(cs.peers)+1)
	out[cs.Name()] = cs.serverAddr
	for name, addr := range cs.peers {
		out[name] = addr
	}
	return out
}

// ConsensusStatusView is the status-UI projection of the consensus engine,
// shaped to the master.html template (which reads .Leader and ranges .Peers with
// $p.Name). It is the drop-in for the goraft raft.Server the template used to
// render.
type ConsensusStatusView struct {
	Leader string
	Peers  []ConsensusPeerView
}

type ConsensusPeerView struct {
	Name string
}

// StatusView returns the consensus engine projected for the status UI: the
// pinned-writer leader and the OTHER masters (self excluded, matching goraft's
// Peers() "other masters" semantics).
func (cs *ConsensusServer) StatusView() ConsensusStatusView {
	self := cs.Name()
	view := ConsensusStatusView{Leader: cs.Leader()}
	for name := range cs.PeerAddresses() {
		if name == self {
			continue
		}
		view.Peers = append(view.Peers, ConsensusPeerView{Name: name})
	}
	return view
}

// --- /cluster HTTP handlers (drop-ins for RaftServer.StatusHandler /
// HealthzHandler; same JSON shape and the same healthz semantics). ---

func (cs *ConsensusServer) StatusHandler(w http.ResponseWriter, r *http.Request) {
	ret := ClusterStatusResult{
		IsLeader:    cs.topo.IsLeader(),
		MaxVolumeId: cs.topo.GetMaxVolumeId(),
	}
	for name := range cs.PeerAddresses() {
		ret.Peers = append(ret.Peers, name)
	}
	if leader, e := cs.topo.Leader(); e == nil {
		ret.Leader = leader
	}
	writeJsonQuiet(w, r, http.StatusOK, ret)
}

func (cs *ConsensusServer) HealthzHandler(w http.ResponseWriter, r *http.Request) {
	leader, err := cs.topo.Leader()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if cs.serverAddr.Equals(leader) {
		expBackoff := backoff.NewExponentialBackOff()
		expBackoff.InitialInterval = 20 * time.Millisecond
		expBackoff.MaxInterval = 1 * time.Second
		expBackoff.MaxElapsedTime = 5 * time.Second
		isLocked, err := backoff.RetryWithData(cs.topo.IsChildLocked, expBackoff)
		if err != nil {
			glog.Errorf("HealthzHandler: %+v", err)
		}
		if isLocked {
			w.WriteHeader(http.StatusLocked)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
