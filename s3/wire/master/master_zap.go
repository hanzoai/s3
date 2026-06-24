// Code generated from master.proto; DO NOT EDIT.

package masterwire

import (
	"fmt"

	"github.com/zap-proto/go/rpc"
)

// Method ordinals for the Hanzo service (stable 1-based wire ids, in
// master.proto declaration order).
//
// SendHeartbeat, KeepConnected and StreamAssign are bidirectional streaming
// RPCs (stream Request -> stream Response). Their ordinals, request/response
// message schemas and dispatch entries are emitted here for stability; the
// per-message streaming send/recv bodies land when the transport streaming
// primitive ships. Until then each is routed as a single request -> response
// exchange, identical to the unary methods.
const (
	HanzoSendHeartbeatOrdinal          uint32 = 1 // STREAMING (stream Heartbeat -> stream HeartbeatResponse)
	HanzoKeepConnectedOrdinal          uint32 = 2 // STREAMING (stream KeepConnectedRequest -> stream KeepConnectedResponse)
	HanzoLookupVolumeOrdinal           uint32 = 3
	HanzoAssignOrdinal                 uint32 = 4
	HanzoStreamAssignOrdinal           uint32 = 5 // STREAMING (stream AssignRequest -> stream AssignResponse)
	HanzoStatisticsOrdinal             uint32 = 6
	HanzoCollectionListOrdinal         uint32 = 7
	HanzoCollectionDeleteOrdinal       uint32 = 8
	HanzoVolumeListOrdinal             uint32 = 9
	HanzoLookupEcVolumeOrdinal         uint32 = 10
	HanzoVacuumVolumeOrdinal           uint32 = 11
	HanzoDisableVacuumOrdinal          uint32 = 12
	HanzoEnableVacuumOrdinal           uint32 = 13
	HanzoVolumeMarkReadonlyOrdinal     uint32 = 14
	HanzoGetMasterConfigurationOrdinal uint32 = 15
	HanzoListClusterNodesOrdinal       uint32 = 16
	HanzoLeaseAdminTokenOrdinal        uint32 = 17
	HanzoReleaseAdminTokenOrdinal      uint32 = 18
	HanzoPingOrdinal                   uint32 = 19
	HanzoRaftListClusterServersOrdinal uint32 = 20
	HanzoRaftAddServerOrdinal          uint32 = 21
	HanzoRaftRemoveServerOrdinal       uint32 = 22
	HanzoRaftLeadershipTransferOrdinal uint32 = 23
	HanzoVolumeGrowOrdinal             uint32 = 24
)

// HanzoChannel ships one Call envelope and awaits its correlated Response.
type HanzoChannel interface {
	Call(envelope []byte) (rpc.Response, error)
}

// HanzoClient is a typed RPC client for the Hanzo service over a ZAP call
// channel. Each call takes a fresh PromiseID from sess; the pipelined "On" form
// of a method sets Target to a prior call's Promise so the server chains them.
type HanzoClient struct {
	ch   HanzoChannel
	cap  []byte
	sess *rpc.Session
}

// NewHanzoClient returns a client that issues calls over ch, attaching cap
// (which may be nil) to every request.
func NewHanzoClient(ch HanzoChannel, capability []byte) *HanzoClient {
	return &HanzoClient{ch: ch, cap: capability, sess: rpc.NewSession()}
}

// invoke issues one request to method and returns the response body. Shared by
// every method's public entry points.
func (c *HanzoClient) invoke(method, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    method,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("Hanzo method %d: %s", method, resp.Body)
		}
		return p, nil, fmt.Errorf("Hanzo method %d: status %d", method, resp.Status)
	}
	return p, resp.Body, nil
}

// SendHeartbeat issues the SendHeartbeat RPC. STREAMING: until the transport
// streaming primitive ships this sends one Heartbeat and returns one
// HeartbeatResponse body.
func (c *HanzoClient) SendHeartbeat(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoSendHeartbeatOrdinal, rpc.NoTarget, req)
}

// SendHeartbeatOn issues SendHeartbeat pipelined on the answer of on.
func (c *HanzoClient) SendHeartbeatOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoSendHeartbeatOrdinal, on.ID, nil)
}

// KeepConnected issues the KeepConnected RPC. STREAMING: until the transport
// streaming primitive ships this sends one KeepConnectedRequest and returns one
// KeepConnectedResponse body.
func (c *HanzoClient) KeepConnected(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoKeepConnectedOrdinal, rpc.NoTarget, req)
}

// KeepConnectedOn issues KeepConnected pipelined on the answer of on.
func (c *HanzoClient) KeepConnectedOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoKeepConnectedOrdinal, on.ID, nil)
}

// LookupVolume issues the LookupVolume RPC.
func (c *HanzoClient) LookupVolume(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoLookupVolumeOrdinal, rpc.NoTarget, req)
}

// LookupVolumeOn issues LookupVolume pipelined on the answer of on.
func (c *HanzoClient) LookupVolumeOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoLookupVolumeOrdinal, on.ID, nil)
}

// Assign issues the Assign RPC.
func (c *HanzoClient) Assign(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoAssignOrdinal, rpc.NoTarget, req)
}

// AssignOn issues Assign pipelined on the answer of on.
func (c *HanzoClient) AssignOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoAssignOrdinal, on.ID, nil)
}

// StreamAssign issues the StreamAssign RPC. STREAMING: until the transport
// streaming primitive ships this sends one AssignRequest and returns one
// AssignResponse body.
func (c *HanzoClient) StreamAssign(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoStreamAssignOrdinal, rpc.NoTarget, req)
}

// StreamAssignOn issues StreamAssign pipelined on the answer of on.
func (c *HanzoClient) StreamAssignOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoStreamAssignOrdinal, on.ID, nil)
}

// Statistics issues the Statistics RPC.
func (c *HanzoClient) Statistics(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoStatisticsOrdinal, rpc.NoTarget, req)
}

// StatisticsOn issues Statistics pipelined on the answer of on.
func (c *HanzoClient) StatisticsOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoStatisticsOrdinal, on.ID, nil)
}

// CollectionList issues the CollectionList RPC.
func (c *HanzoClient) CollectionList(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoCollectionListOrdinal, rpc.NoTarget, req)
}

// CollectionListOn issues CollectionList pipelined on the answer of on.
func (c *HanzoClient) CollectionListOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoCollectionListOrdinal, on.ID, nil)
}

// CollectionDelete issues the CollectionDelete RPC.
func (c *HanzoClient) CollectionDelete(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoCollectionDeleteOrdinal, rpc.NoTarget, req)
}

// CollectionDeleteOn issues CollectionDelete pipelined on the answer of on.
func (c *HanzoClient) CollectionDeleteOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoCollectionDeleteOrdinal, on.ID, nil)
}

// VolumeList issues the VolumeList RPC.
func (c *HanzoClient) VolumeList(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoVolumeListOrdinal, rpc.NoTarget, req)
}

// VolumeListOn issues VolumeList pipelined on the answer of on.
func (c *HanzoClient) VolumeListOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoVolumeListOrdinal, on.ID, nil)
}

// LookupEcVolume issues the LookupEcVolume RPC.
func (c *HanzoClient) LookupEcVolume(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoLookupEcVolumeOrdinal, rpc.NoTarget, req)
}

// LookupEcVolumeOn issues LookupEcVolume pipelined on the answer of on.
func (c *HanzoClient) LookupEcVolumeOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoLookupEcVolumeOrdinal, on.ID, nil)
}

// VacuumVolume issues the VacuumVolume RPC.
func (c *HanzoClient) VacuumVolume(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoVacuumVolumeOrdinal, rpc.NoTarget, req)
}

// VacuumVolumeOn issues VacuumVolume pipelined on the answer of on.
func (c *HanzoClient) VacuumVolumeOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoVacuumVolumeOrdinal, on.ID, nil)
}

// DisableVacuum issues the DisableVacuum RPC.
func (c *HanzoClient) DisableVacuum(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoDisableVacuumOrdinal, rpc.NoTarget, req)
}

// DisableVacuumOn issues DisableVacuum pipelined on the answer of on.
func (c *HanzoClient) DisableVacuumOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoDisableVacuumOrdinal, on.ID, nil)
}

// EnableVacuum issues the EnableVacuum RPC.
func (c *HanzoClient) EnableVacuum(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoEnableVacuumOrdinal, rpc.NoTarget, req)
}

// EnableVacuumOn issues EnableVacuum pipelined on the answer of on.
func (c *HanzoClient) EnableVacuumOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoEnableVacuumOrdinal, on.ID, nil)
}

// VolumeMarkReadonly issues the VolumeMarkReadonly RPC.
func (c *HanzoClient) VolumeMarkReadonly(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoVolumeMarkReadonlyOrdinal, rpc.NoTarget, req)
}

// VolumeMarkReadonlyOn issues VolumeMarkReadonly pipelined on the answer of on.
func (c *HanzoClient) VolumeMarkReadonlyOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoVolumeMarkReadonlyOrdinal, on.ID, nil)
}

// GetMasterConfiguration issues the GetMasterConfiguration RPC.
func (c *HanzoClient) GetMasterConfiguration(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoGetMasterConfigurationOrdinal, rpc.NoTarget, req)
}

// GetMasterConfigurationOn issues GetMasterConfiguration pipelined on the answer
// of on.
func (c *HanzoClient) GetMasterConfigurationOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoGetMasterConfigurationOrdinal, on.ID, nil)
}

// ListClusterNodes issues the ListClusterNodes RPC.
func (c *HanzoClient) ListClusterNodes(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoListClusterNodesOrdinal, rpc.NoTarget, req)
}

// ListClusterNodesOn issues ListClusterNodes pipelined on the answer of on.
func (c *HanzoClient) ListClusterNodesOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoListClusterNodesOrdinal, on.ID, nil)
}

// LeaseAdminToken issues the LeaseAdminToken RPC.
func (c *HanzoClient) LeaseAdminToken(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoLeaseAdminTokenOrdinal, rpc.NoTarget, req)
}

// LeaseAdminTokenOn issues LeaseAdminToken pipelined on the answer of on.
func (c *HanzoClient) LeaseAdminTokenOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoLeaseAdminTokenOrdinal, on.ID, nil)
}

// ReleaseAdminToken issues the ReleaseAdminToken RPC.
func (c *HanzoClient) ReleaseAdminToken(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoReleaseAdminTokenOrdinal, rpc.NoTarget, req)
}

// ReleaseAdminTokenOn issues ReleaseAdminToken pipelined on the answer of on.
func (c *HanzoClient) ReleaseAdminTokenOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoReleaseAdminTokenOrdinal, on.ID, nil)
}

// Ping issues the Ping RPC.
func (c *HanzoClient) Ping(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoPingOrdinal, rpc.NoTarget, req)
}

// PingOn issues Ping pipelined on the answer of on.
func (c *HanzoClient) PingOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoPingOrdinal, on.ID, nil)
}

// RaftListClusterServers issues the RaftListClusterServers RPC.
func (c *HanzoClient) RaftListClusterServers(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoRaftListClusterServersOrdinal, rpc.NoTarget, req)
}

// RaftListClusterServersOn issues RaftListClusterServers pipelined on the answer
// of on.
func (c *HanzoClient) RaftListClusterServersOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoRaftListClusterServersOrdinal, on.ID, nil)
}

// RaftAddServer issues the RaftAddServer RPC.
func (c *HanzoClient) RaftAddServer(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoRaftAddServerOrdinal, rpc.NoTarget, req)
}

// RaftAddServerOn issues RaftAddServer pipelined on the answer of on.
func (c *HanzoClient) RaftAddServerOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoRaftAddServerOrdinal, on.ID, nil)
}

// RaftRemoveServer issues the RaftRemoveServer RPC.
func (c *HanzoClient) RaftRemoveServer(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoRaftRemoveServerOrdinal, rpc.NoTarget, req)
}

// RaftRemoveServerOn issues RaftRemoveServer pipelined on the answer of on.
func (c *HanzoClient) RaftRemoveServerOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoRaftRemoveServerOrdinal, on.ID, nil)
}

// RaftLeadershipTransfer issues the RaftLeadershipTransfer RPC.
func (c *HanzoClient) RaftLeadershipTransfer(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoRaftLeadershipTransferOrdinal, rpc.NoTarget, req)
}

// RaftLeadershipTransferOn issues RaftLeadershipTransfer pipelined on the answer
// of on.
func (c *HanzoClient) RaftLeadershipTransferOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoRaftLeadershipTransferOrdinal, on.ID, nil)
}

// VolumeGrow issues the VolumeGrow RPC.
func (c *HanzoClient) VolumeGrow(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoVolumeGrowOrdinal, rpc.NoTarget, req)
}

// VolumeGrowOn issues VolumeGrow pipelined on the answer of on.
func (c *HanzoClient) VolumeGrowOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoVolumeGrowOrdinal, on.ID, nil)
}

// HanzoHandler is the server contract for the Hanzo service. Implement each
// method, then route requests to it with DispatchHanzo.
//
// SendHeartbeat, KeepConnected and StreamAssign are STREAMING in the proto;
// their handler signatures use the same single request -> response body shape
// as the unary methods until the transport streaming primitive ships.
type HanzoHandler interface {
	SendHeartbeat(req []byte) ([]byte, error) // STREAMING
	KeepConnected(req []byte) ([]byte, error) // STREAMING
	LookupVolume(req []byte) ([]byte, error)
	Assign(req []byte) ([]byte, error)
	StreamAssign(req []byte) ([]byte, error) // STREAMING
	Statistics(req []byte) ([]byte, error)
	CollectionList(req []byte) ([]byte, error)
	CollectionDelete(req []byte) ([]byte, error)
	VolumeList(req []byte) ([]byte, error)
	LookupEcVolume(req []byte) ([]byte, error)
	VacuumVolume(req []byte) ([]byte, error)
	DisableVacuum(req []byte) ([]byte, error)
	EnableVacuum(req []byte) ([]byte, error)
	VolumeMarkReadonly(req []byte) ([]byte, error)
	GetMasterConfiguration(req []byte) ([]byte, error)
	ListClusterNodes(req []byte) ([]byte, error)
	LeaseAdminToken(req []byte) ([]byte, error)
	ReleaseAdminToken(req []byte) ([]byte, error)
	Ping(req []byte) ([]byte, error)
	RaftListClusterServers(req []byte) ([]byte, error)
	RaftAddServer(req []byte) ([]byte, error)
	RaftRemoveServer(req []byte) ([]byte, error)
	RaftLeadershipTransfer(req []byte) ([]byte, error)
	VolumeGrow(req []byte) ([]byte, error)
}

// DispatchHanzo decodes a Call envelope, routes it by method ordinal to h, and
// returns the response envelope. An unknown ordinal yields a StatusNotFound
// response; a handler error yields StatusInternal.
func DispatchHanzo(h HanzoHandler, envelope []byte) ([]byte, error) {
	call, err := rpc.ParseRequest(envelope)
	if err != nil {
		return nil, err
	}
	var body []byte
	var herr error
	switch call.Method {
	case HanzoSendHeartbeatOrdinal: // STREAMING
		body, herr = h.SendHeartbeat(call.Payload)
	case HanzoKeepConnectedOrdinal: // STREAMING
		body, herr = h.KeepConnected(call.Payload)
	case HanzoLookupVolumeOrdinal:
		body, herr = h.LookupVolume(call.Payload)
	case HanzoAssignOrdinal:
		body, herr = h.Assign(call.Payload)
	case HanzoStreamAssignOrdinal: // STREAMING
		body, herr = h.StreamAssign(call.Payload)
	case HanzoStatisticsOrdinal:
		body, herr = h.Statistics(call.Payload)
	case HanzoCollectionListOrdinal:
		body, herr = h.CollectionList(call.Payload)
	case HanzoCollectionDeleteOrdinal:
		body, herr = h.CollectionDelete(call.Payload)
	case HanzoVolumeListOrdinal:
		body, herr = h.VolumeList(call.Payload)
	case HanzoLookupEcVolumeOrdinal:
		body, herr = h.LookupEcVolume(call.Payload)
	case HanzoVacuumVolumeOrdinal:
		body, herr = h.VacuumVolume(call.Payload)
	case HanzoDisableVacuumOrdinal:
		body, herr = h.DisableVacuum(call.Payload)
	case HanzoEnableVacuumOrdinal:
		body, herr = h.EnableVacuum(call.Payload)
	case HanzoVolumeMarkReadonlyOrdinal:
		body, herr = h.VolumeMarkReadonly(call.Payload)
	case HanzoGetMasterConfigurationOrdinal:
		body, herr = h.GetMasterConfiguration(call.Payload)
	case HanzoListClusterNodesOrdinal:
		body, herr = h.ListClusterNodes(call.Payload)
	case HanzoLeaseAdminTokenOrdinal:
		body, herr = h.LeaseAdminToken(call.Payload)
	case HanzoReleaseAdminTokenOrdinal:
		body, herr = h.ReleaseAdminToken(call.Payload)
	case HanzoPingOrdinal:
		body, herr = h.Ping(call.Payload)
	case HanzoRaftListClusterServersOrdinal:
		body, herr = h.RaftListClusterServers(call.Payload)
	case HanzoRaftAddServerOrdinal:
		body, herr = h.RaftAddServer(call.Payload)
	case HanzoRaftRemoveServerOrdinal:
		body, herr = h.RaftRemoveServer(call.Payload)
	case HanzoRaftLeadershipTransferOrdinal:
		body, herr = h.RaftLeadershipTransfer(call.Payload)
	case HanzoVolumeGrowOrdinal:
		body, herr = h.VolumeGrow(call.Payload)
	default:
		return rpc.BuildResponse(rpc.StatusNotFound, call.PromiseID, nil), nil
	}
	if herr != nil {
		// Carry the handler error message so the caller can reconstruct sentinels.
		return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(herr.Error())), nil
	}
	return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
}
