package s3server

import (
	"context"
	"fmt"

	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/master_pb"
)

// This file replaces the raft-backed cluster RPCs (the deleted
// master_grpc_server_raft.go). The master is leaderless: membership and the
// pinned writer come from the Coordinator, not a raft configuration. The RPC
// surface is retained so existing clients (the admin dashboard's master list)
// keep working; the mutation RPCs become advisory membership hints since there
// is no voter roster to reconfigure.

// RaftListClusterServers reports the coordination membership and which member is
// the pinned writer. Any master can answer — there is no leader to forward to.
func (ms *MasterServer) RaftListClusterServers(ctx context.Context, req *master_pb.RaftListClusterServersRequest) (*master_pb.RaftListClusterServersResponse, error) {
	resp := &master_pb.RaftListClusterServersResponse{}
	if ms.Topo.Coordinator == nil {
		return resp, nil
	}
	writer := ms.Topo.Coordinator.Writer()
	for _, m := range ms.Topo.Coordinator.Members() {
		resp.ClusterServers = append(resp.ClusterServers, &master_pb.RaftListClusterServersResponse_ClusterServers{
			Id:       string(m),
			Address:  m.ToGrpcAddress(),
			Suffrage: "Voter",
			IsWriter: m == writer,
		})
	}
	return resp, nil
}

// RaftAddServer is advisory: it widens the writer-eligible set. There is no
// voter promotion in the leaderless model — the pinned writer is a pure
// function of membership.
func (ms *MasterServer) RaftAddServer(ctx context.Context, req *master_pb.RaftAddServerRequest) (*master_pb.RaftAddServerResponse, error) {
	resp := &master_pb.RaftAddServerResponse{}
	if req.Id == "" {
		return resp, nil
	}
	ms.setCoordinatorMember(req.Id, pb.ServerAddress(req.Id), true)
	return resp, nil
}

// RaftRemoveServer is advisory: it drops a master from the writer-eligible set,
// re-pinning the writer to the next-lowest member if the removed one held it.
func (ms *MasterServer) RaftRemoveServer(ctx context.Context, req *master_pb.RaftRemoveServerRequest) (*master_pb.RaftRemoveServerResponse, error) {
	resp := &master_pb.RaftRemoveServerResponse{}
	if req.Id == "" {
		return resp, nil
	}
	ms.setCoordinatorMember(req.Id, "", false)
	return resp, nil
}

// RaftLeadershipTransfer has no analogue: the writer is pinned deterministically
// to the lowest-address member, so leadership cannot be hand-transferred.
func (ms *MasterServer) RaftLeadershipTransfer(ctx context.Context, req *master_pb.RaftLeadershipTransferRequest) (*master_pb.RaftLeadershipTransferResponse, error) {
	return nil, fmt.Errorf("leadership transfer is not supported: the master is leaderless; the writer is pinned to the lowest-address member")
}

// setCoordinatorMember adds or removes id from the writer-eligible set and
// republishes the membership to the Coordinator.
func (ms *MasterServer) setCoordinatorMember(id string, addr pb.ServerAddress, add bool) {
	ms.coordinatorPeersLock.Lock()
	if add {
		ms.coordinatorPeers[id] = addr
	} else {
		delete(ms.coordinatorPeers, id)
	}
	members := make([]pb.ServerAddress, 0, len(ms.coordinatorPeers))
	for _, p := range ms.coordinatorPeers {
		members = append(members, p)
	}
	ms.coordinatorPeersLock.Unlock()
	ms.Topo.Coordinator.SetMembers(ms.option.Master, members)
}
