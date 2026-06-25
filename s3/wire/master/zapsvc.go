// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// Native ZAP service adapter for the master (Hanzo) service — the hand-written
// glue that binds the generated DispatchHanzo / HanzoClient (master_zap.go) to
// the real master engine. This is the master analogue of s3/zapsvc/service.go:
// the SAME Channel/Client/Handler/Dispatch shape, the SAME zero-copy doctrine
// (the bytes ARE the message: masterwire.New* to build, masterwire.Wrap* to
// read — no struct, no marshal, no protobuf), carried over the canonical
// github.com/zap-proto/go transport.
//
// Layering (one and only one way):
//   - Backend is the engine contract, in engine terms — the real master
//     (topology + raft + vacuum) implements it. It does NOT live here; tests
//     and integration supply it. Each unary RPC is one Go method whose
//     arguments/returns are the message buffers (masterwire.New*/Wrap*), so the
//     backend stays zero-copy end to end without re-deriving 24 ad-hoc structs.
//   - backendHandler adapts a Backend to the generated HanzoHandler.
//   - Dispatch / Serve stand up the server; Dial / Client are what internal
//     callers (filer, volume server, shell, admin) hold instead of the killed
//     gRPC master client.
//
// Streaming: SendHeartbeat, KeepConnected and StreamAssign are bidirectional
// streams in master.proto. zapgen cannot express streaming and the transport's
// stream primitive is wired by hand (transport.OpenStream / ListenStream); until
// that lands for master they are NOT served here. The Backend interface marks
// them with the request/response message buffers but the dispatch path returns
// ErrStreamingNotWired rather than silently degrading a stream to one
// request->one response. See the STREAMING block below.

package masterwire

import (
	"errors"

	"github.com/zap-proto/go/transport"
)

// ErrStreamingNotWired is returned by the streaming RPCs (SendHeartbeat,
// KeepConnected, StreamAssign) until the hand-wired transport.OpenStream /
// transport.ListenStream path is added for master. It is deliberately NOT a
// degraded one-shot exchange — a streaming method must stream.
var ErrStreamingNotWired = errors.New("masterwire: streaming RPC not wired (use transport.OpenStream/ListenStream)")

// Backend is the master engine contract the ZAP service delegates to, expressed
// in wire-message terms: each method takes the decoded-on-arrival request buffer
// and returns the response buffer (build it with the matching masterwire.New*).
// The real master implements this against its topology/raft/vacuum state; tests
// use an in-memory fake. Methods may decode the request with the matching
// masterwire.Wrap* and never touch protobuf.
//
// The three streaming RPCs are present for completeness but the unary transport
// cannot carry them; a Backend should return ErrStreamingNotWired (the default
// dispatch already short-circuits them — see backendHandler).
type Backend interface {
	// --- streaming (NOT served over the unary path; see ErrStreamingNotWired) ---
	SendHeartbeat(req []byte) ([]byte, error)
	KeepConnected(req []byte) ([]byte, error)
	StreamAssign(req []byte) ([]byte, error)

	// --- unary ---
	LookupVolume(req []byte) ([]byte, error)
	Assign(req []byte) ([]byte, error)
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

// backendHandler adapts a Backend to the generated HanzoHandler. The unary
// methods pass straight through; the three streaming methods are forced to
// ErrStreamingNotWired so a misconfigured caller fails loudly instead of getting
// a truncated stream.
type backendHandler struct{ b Backend }

func (h backendHandler) SendHeartbeat(req []byte) ([]byte, error) { return nil, ErrStreamingNotWired }
func (h backendHandler) KeepConnected(req []byte) ([]byte, error) { return nil, ErrStreamingNotWired }
func (h backendHandler) StreamAssign(req []byte) ([]byte, error)  { return nil, ErrStreamingNotWired }

func (h backendHandler) LookupVolume(req []byte) ([]byte, error)   { return h.b.LookupVolume(req) }
func (h backendHandler) Assign(req []byte) ([]byte, error)         { return h.b.Assign(req) }
func (h backendHandler) Statistics(req []byte) ([]byte, error)     { return h.b.Statistics(req) }
func (h backendHandler) CollectionList(req []byte) ([]byte, error) { return h.b.CollectionList(req) }
func (h backendHandler) CollectionDelete(req []byte) ([]byte, error) {
	return h.b.CollectionDelete(req)
}
func (h backendHandler) VolumeList(req []byte) ([]byte, error)     { return h.b.VolumeList(req) }
func (h backendHandler) LookupEcVolume(req []byte) ([]byte, error) { return h.b.LookupEcVolume(req) }
func (h backendHandler) VacuumVolume(req []byte) ([]byte, error)   { return h.b.VacuumVolume(req) }
func (h backendHandler) DisableVacuum(req []byte) ([]byte, error)  { return h.b.DisableVacuum(req) }
func (h backendHandler) EnableVacuum(req []byte) ([]byte, error)   { return h.b.EnableVacuum(req) }
func (h backendHandler) VolumeMarkReadonly(req []byte) ([]byte, error) {
	return h.b.VolumeMarkReadonly(req)
}
func (h backendHandler) GetMasterConfiguration(req []byte) ([]byte, error) {
	return h.b.GetMasterConfiguration(req)
}
func (h backendHandler) ListClusterNodes(req []byte) ([]byte, error) {
	return h.b.ListClusterNodes(req)
}
func (h backendHandler) LeaseAdminToken(req []byte) ([]byte, error) {
	return h.b.LeaseAdminToken(req)
}
func (h backendHandler) ReleaseAdminToken(req []byte) ([]byte, error) {
	return h.b.ReleaseAdminToken(req)
}
func (h backendHandler) Ping(req []byte) ([]byte, error) { return h.b.Ping(req) }
func (h backendHandler) RaftListClusterServers(req []byte) ([]byte, error) {
	return h.b.RaftListClusterServers(req)
}
func (h backendHandler) RaftAddServer(req []byte) ([]byte, error) { return h.b.RaftAddServer(req) }
func (h backendHandler) RaftRemoveServer(req []byte) ([]byte, error) {
	return h.b.RaftRemoveServer(req)
}
func (h backendHandler) RaftLeadershipTransfer(req []byte) ([]byte, error) {
	return h.b.RaftLeadershipTransfer(req)
}
func (h backendHandler) VolumeGrow(req []byte) ([]byte, error) { return h.b.VolumeGrow(req) }

// NewHandler adapts a Backend to the generated HanzoHandler. Use it when you
// want to compose dispatch yourself (e.g. behind a custom listener).
func NewHandler(b Backend) HanzoHandler { return backendHandler{b: b} }

// Dispatch is the Hanzo server dispatch bound to b — pass it to
// transport.Listen. Exposed so callers can compose it (e.g. behind PQ-TLS).
func Dispatch(b Backend) transport.Dispatch {
	h := backendHandler{b: b}
	return func(envelope []byte) ([]byte, error) {
		return DispatchHanzo(h, envelope)
	}
}

// Serve starts the native ZAP master service on network/addr (e.g. "tcp",
// ":19333"), backed by b. Returns the running server; Close stops it. For the
// PQ-secured mesh, build a transport.Conn over PQ-TLS and serve dispatch on it.
func Serve(network, addr string, b Backend) (*transport.Server, error) {
	return transport.Listen(network, addr, Dispatch(b))
}

// Client is the typed master ZAP service client internal services hold. It owns
// the transport connection to one master endpoint and replaces the killed gRPC
// master client. Each method builds its request with masterwire.New* and decodes
// the reply with masterwire.Wrap*.
type Client struct {
	conn transport.Conn
	rpc  *HanzoClient
}

// Dial connects to the master ZAP service at addr over plain TCP. For the
// PQ-secured mesh use transport.DialTLS and NewClient.
func Dial(network, addr string) (*Client, error) {
	conn, err := transport.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, rpc: NewHanzoClient(conn, nil)}, nil
}

// NewClient wraps an already-established transport.Conn (TCP, Unix, or PQ-TLS).
func NewClient(conn transport.Conn) *Client {
	return &Client{conn: conn, rpc: NewHanzoClient(conn, nil)}
}

// Close releases the underlying connection.
func (c *Client) Close() error { return c.conn.Close() }

// RPC exposes the low-level generated client for the pipelined "On" forms and
// the streaming-ordinal calls the typed helpers below do not wrap.
func (c *Client) RPC() *HanzoClient { return c.rpc }

// --- typed unary calls (build request -> ship -> decode reply) ---

// LookupVolume resolves volume-or-file ids to their locations.
func (c *Client) LookupVolume(in LookupVolumeRequestInput) (LookupVolumeResponse, error) {
	_, body, err := c.rpc.LookupVolume(NewLookupVolumeRequest(in))
	if err != nil {
		return LookupVolumeResponse{}, err
	}
	return WrapLookupVolumeResponse(body)
}

// Assign allocates a file id / volume for a write.
func (c *Client) Assign(in AssignRequestInput) (AssignResponse, error) {
	_, body, err := c.rpc.Assign(NewAssignRequest(in))
	if err != nil {
		return AssignResponse{}, err
	}
	return WrapAssignResponse(body)
}

// Statistics returns cluster size/usage counters.
func (c *Client) Statistics(in StatisticsRequestInput) (StatisticsResponse, error) {
	_, body, err := c.rpc.Statistics(NewStatisticsRequest(in))
	if err != nil {
		return StatisticsResponse{}, err
	}
	return WrapStatisticsResponse(body)
}

// CollectionList lists collections.
func (c *Client) CollectionList(in CollectionListRequestInput) (CollectionListResponse, error) {
	_, body, err := c.rpc.CollectionList(NewCollectionListRequest(in))
	if err != nil {
		return CollectionListResponse{}, err
	}
	return WrapCollectionListResponse(body)
}

// CollectionDelete deletes a collection.
func (c *Client) CollectionDelete(in CollectionDeleteRequestInput) (CollectionDeleteResponse, error) {
	_, body, err := c.rpc.CollectionDelete(NewCollectionDeleteRequest(in))
	if err != nil {
		return CollectionDeleteResponse{}, err
	}
	return WrapCollectionDeleteResponse(body)
}

// VolumeList returns the full topology.
func (c *Client) VolumeList(in VolumeListRequestInput) (VolumeListResponse, error) {
	_, body, err := c.rpc.VolumeList(NewVolumeListRequest(in))
	if err != nil {
		return VolumeListResponse{}, err
	}
	return WrapVolumeListResponse(body)
}

// LookupEcVolume resolves the shard locations of an EC volume.
func (c *Client) LookupEcVolume(in LookupEcVolumeRequestInput) (LookupEcVolumeResponse, error) {
	_, body, err := c.rpc.LookupEcVolume(NewLookupEcVolumeRequest(in))
	if err != nil {
		return LookupEcVolumeResponse{}, err
	}
	return WrapLookupEcVolumeResponse(body)
}

// VacuumVolume triggers a vacuum.
func (c *Client) VacuumVolume(in VacuumVolumeRequestInput) (VacuumVolumeResponse, error) {
	_, body, err := c.rpc.VacuumVolume(NewVacuumVolumeRequest(in))
	if err != nil {
		return VacuumVolumeResponse{}, err
	}
	return WrapVacuumVolumeResponse(body)
}

// DisableVacuum disables automatic vacuum.
func (c *Client) DisableVacuum(in DisableVacuumRequestInput) (DisableVacuumResponse, error) {
	_, body, err := c.rpc.DisableVacuum(NewDisableVacuumRequest(in))
	if err != nil {
		return DisableVacuumResponse{}, err
	}
	return WrapDisableVacuumResponse(body)
}

// EnableVacuum enables automatic vacuum.
func (c *Client) EnableVacuum(in EnableVacuumRequestInput) (EnableVacuumResponse, error) {
	_, body, err := c.rpc.EnableVacuum(NewEnableVacuumRequest(in))
	if err != nil {
		return EnableVacuumResponse{}, err
	}
	return WrapEnableVacuumResponse(body)
}

// VolumeMarkReadonly marks a volume read-only.
func (c *Client) VolumeMarkReadonly(in VolumeMarkReadonlyRequestInput) (VolumeMarkReadonlyResponse, error) {
	_, body, err := c.rpc.VolumeMarkReadonly(NewVolumeMarkReadonlyRequest(in))
	if err != nil {
		return VolumeMarkReadonlyResponse{}, err
	}
	return WrapVolumeMarkReadonlyResponse(body)
}

// GetMasterConfiguration returns the master's runtime configuration.
func (c *Client) GetMasterConfiguration(in GetMasterConfigurationRequestInput) (GetMasterConfigurationResponse, error) {
	_, body, err := c.rpc.GetMasterConfiguration(NewGetMasterConfigurationRequest(in))
	if err != nil {
		return GetMasterConfigurationResponse{}, err
	}
	return WrapGetMasterConfigurationResponse(body)
}

// ListClusterNodes lists known cluster nodes of a client type.
func (c *Client) ListClusterNodes(in ListClusterNodesRequestInput) (ListClusterNodesResponse, error) {
	_, body, err := c.rpc.ListClusterNodes(NewListClusterNodesRequest(in))
	if err != nil {
		return ListClusterNodesResponse{}, err
	}
	return WrapListClusterNodesResponse(body)
}

// LeaseAdminToken leases the admin lock token.
func (c *Client) LeaseAdminToken(in LeaseAdminTokenRequestInput) (LeaseAdminTokenResponse, error) {
	_, body, err := c.rpc.LeaseAdminToken(NewLeaseAdminTokenRequest(in))
	if err != nil {
		return LeaseAdminTokenResponse{}, err
	}
	return WrapLeaseAdminTokenResponse(body)
}

// ReleaseAdminToken releases the admin lock token.
func (c *Client) ReleaseAdminToken(in ReleaseAdminTokenRequestInput) (ReleaseAdminTokenResponse, error) {
	_, body, err := c.rpc.ReleaseAdminToken(NewReleaseAdminTokenRequest(in))
	if err != nil {
		return ReleaseAdminTokenResponse{}, err
	}
	return WrapReleaseAdminTokenResponse(body)
}

// Ping round-trips a liveness probe.
func (c *Client) Ping(in PingRequestInput) (PingResponse, error) {
	_, body, err := c.rpc.Ping(NewPingRequest(in))
	if err != nil {
		return PingResponse{}, err
	}
	return WrapPingResponse(body)
}

// RaftListClusterServers lists the raft membership.
func (c *Client) RaftListClusterServers(in RaftListClusterServersRequestInput) (RaftListClusterServersResponse, error) {
	_, body, err := c.rpc.RaftListClusterServers(NewRaftListClusterServersRequest(in))
	if err != nil {
		return RaftListClusterServersResponse{}, err
	}
	return WrapRaftListClusterServersResponse(body)
}

// RaftAddServer adds a raft server.
func (c *Client) RaftAddServer(in RaftAddServerRequestInput) (RaftAddServerResponse, error) {
	_, body, err := c.rpc.RaftAddServer(NewRaftAddServerRequest(in))
	if err != nil {
		return RaftAddServerResponse{}, err
	}
	return WrapRaftAddServerResponse(body)
}

// RaftRemoveServer removes a raft server.
func (c *Client) RaftRemoveServer(in RaftRemoveServerRequestInput) (RaftRemoveServerResponse, error) {
	_, body, err := c.rpc.RaftRemoveServer(NewRaftRemoveServerRequest(in))
	if err != nil {
		return RaftRemoveServerResponse{}, err
	}
	return WrapRaftRemoveServerResponse(body)
}

// RaftLeadershipTransfer transfers raft leadership.
func (c *Client) RaftLeadershipTransfer(in RaftLeadershipTransferRequestInput) (RaftLeadershipTransferResponse, error) {
	_, body, err := c.rpc.RaftLeadershipTransfer(NewRaftLeadershipTransferRequest(in))
	if err != nil {
		return RaftLeadershipTransferResponse{}, err
	}
	return WrapRaftLeadershipTransferResponse(body)
}

// VolumeGrow grows writable volumes.
func (c *Client) VolumeGrow(in VolumeGrowRequestInput) (VolumeGrowResponse, error) {
	_, body, err := c.rpc.VolumeGrow(NewVolumeGrowRequest(in))
	if err != nil {
		return VolumeGrowResponse{}, err
	}
	return WrapVolumeGrowResponse(body)
}
