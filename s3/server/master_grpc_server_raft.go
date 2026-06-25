package s3server

import (
	"context"
	"fmt"

	"github.com/hanzoai/s3/s3/cluster"
	"github.com/hanzoai/s3/s3/pb/master_pb"
)

// These four RPCs are the master service's cluster-membership surface
// (master_pb.HanzoServer). They were Raft management calls; they now operate on
// the Lux consensus engine (ConsensusServer). Membership is advisory in
// permissionless consensus — it changes the writer-pin set and observability, not
// a voted quorum — so add/remove never block on a configuration round.

func (ms *MasterServer) RaftListClusterServers(ctx context.Context, req *master_pb.RaftListClusterServersRequest) (*master_pb.RaftListClusterServersResponse, error) {
	resp := &master_pb.RaftListClusterServersResponse{}

	if ms.consensus == nil {
		return resp, nil
	}

	leader := ms.consensus.Leader()
	for name, addr := range ms.consensus.PeerAddresses() {
		resp.ClusterServers = append(resp.ClusterServers, &master_pb.RaftListClusterServersResponse_ClusterServers{
			Id:       name,
			Address:  addr.ToGrpcAddress(),
			Suffrage: "Voter",
			IsLeader: name == leader,
		})
	}
	return resp, nil
}

func (ms *MasterServer) RaftAddServer(ctx context.Context, req *master_pb.RaftAddServerRequest) (*master_pb.RaftAddServerResponse, error) {
	resp := &master_pb.RaftAddServerResponse{}

	if ms.consensus == nil {
		return resp, nil
	}
	if !ms.Topo.IsLeader() {
		return nil, fmt.Errorf("consensus add server %s failed: %s is not the pinned writer", req.Id, ms.consensus.Name())
	}
	ms.consensus.AddPeerByGrpcAddress(req.Id, req.Address)
	return resp, nil
}

func (ms *MasterServer) RaftRemoveServer(ctx context.Context, req *master_pb.RaftRemoveServerRequest) (*master_pb.RaftRemoveServerResponse, error) {
	resp := &master_pb.RaftRemoveServerResponse{}

	if ms.consensus == nil {
		return resp, nil
	}
	if !ms.Topo.IsLeader() {
		return nil, fmt.Errorf("consensus remove server %s failed: %s is not the pinned writer", req.Id, ms.consensus.Name())
	}

	if !req.Force {
		ms.clientChansLock.RLock()
		_, ok := ms.clientChans[fmt.Sprintf("%s@%s", cluster.MasterType, req.Id)]
		ms.clientChansLock.RUnlock()
		if ok {
			return resp, fmt.Errorf("consensus remove server %s failed: client connection to master exists", req.Id)
		}
	}

	ms.consensus.RemovePeerByName(req.Id)
	return resp, nil
}

func (ms *MasterServer) RaftLeadershipTransfer(ctx context.Context, req *master_pb.RaftLeadershipTransferRequest) (*master_pb.RaftLeadershipTransferResponse, error) {
	resp := &master_pb.RaftLeadershipTransferResponse{}

	// Lux consensus is leaderless with a DETERMINISTIC writer pin (lowest member
	// address): the writer cannot be hand-transferred — it moves only when the
	// current writer leaves the membership, at which point the next-lowest takes
	// over automatically. There is no manual leadership-transfer operation.
	if ms.consensus == nil {
		return nil, fmt.Errorf("consensus not initialized (single master mode)")
	}
	return resp, fmt.Errorf("leadership transfer is not supported: the Lux consensus writer is pinned deterministically and fails over automatically when the writer departs")
}
