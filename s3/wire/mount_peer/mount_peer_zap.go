// Code generated from mount_peer.proto; DO NOT EDIT.

package mount_peerwire

import (
	"fmt"

	"github.com/zap-proto/go/rpc"
)

// Method ordinals for the MountPeer service (stable 1-based wire ids,
// in mount_peer.proto declaration order).
const (
	MountPeerChunkAnnounceOrdinal uint32 = 1
	MountPeerChunkLookupOrdinal   uint32 = 2
	MountPeerFetchChunkOrdinal    uint32 = 3
)

// MountPeerChannel ships one Call envelope and awaits its correlated Response.
type MountPeerChannel interface {
	Call(envelope []byte) (rpc.Response, error)
	// NextPromiseID allocates a call's PromiseID from the CONNECTION rather than
	// from a per-client counter. PromiseID is the connection's demultiplexing
	// key: the transport keys its in-flight table by it, and OpenStream already
	// allocates stream IDs the same way. A client that numbered its own calls
	// would restart at 1 on every construction, so two clients sharing a pooled
	// conn collide — the transport overwrites the first waiter's slot and its
	// response is dropped, blocking that caller for as long as its context runs.
	NextPromiseID() uint32
}

// MountPeerClient is a typed RPC client for the MountPeer service over a ZAP
// call channel. Each call takes a fresh PromiseID from sess; the pipelined "On"
// form of a method sets Target to a prior call's Promise so the server chains
// them.
type MountPeerClient struct {
	ch   MountPeerChannel
	cap  []byte
}

// NewMountPeerClient returns a client that issues calls over ch, attaching cap
// (which may be nil) to every request.
func NewMountPeerClient(ch MountPeerChannel, capability []byte) *MountPeerClient {
	return &MountPeerClient{ch: ch, cap: capability}
}

func (c *MountPeerClient) ChunkAnnounce(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeChunkAnnounce(rpc.NoTarget, req)
}

// ChunkAnnounceOn issues ChunkAnnounce as a dependent call pipelined on the
// answer of on: the server substitutes on's resolved result for this call's
// payload before dispatch, so it ships without waiting for on to round-trip.
func (c *MountPeerClient) ChunkAnnounceOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeChunkAnnounce(on.ID, nil)
}

func (c *MountPeerClient) invokeChunkAnnounce(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    MountPeerChunkAnnounceOrdinal,
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
			return p, nil, fmt.Errorf("MountPeer.ChunkAnnounce: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("MountPeer.ChunkAnnounce: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *MountPeerClient) ChunkLookup(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeChunkLookup(rpc.NoTarget, req)
}

// ChunkLookupOn issues ChunkLookup as a dependent call pipelined on the answer
// of on: the server substitutes on's resolved result for this call's payload
// before dispatch, so it ships without waiting for on to round-trip.
func (c *MountPeerClient) ChunkLookupOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeChunkLookup(on.ID, nil)
}

func (c *MountPeerClient) invokeChunkLookup(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    MountPeerChunkLookupOrdinal,
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
			return p, nil, fmt.Errorf("MountPeer.ChunkLookup: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("MountPeer.ChunkLookup: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: FetchChunk server-streams a sequence of FetchChunkResponse frames.
// This unary-shaped Call returns only the first frame's body; the full
// streaming client (frame iteration) lands when the transport streaming
// primitive ships. The request/response message schemas are already emitted in
// fetch_chunk_zap.go.
func (c *MountPeerClient) FetchChunk(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeFetchChunk(rpc.NoTarget, req)
}

// FetchChunkOn issues FetchChunk as a dependent call pipelined on the answer of
// on: the server substitutes on's resolved result for this call's payload
// before dispatch, so it ships without waiting for on to round-trip.
func (c *MountPeerClient) FetchChunkOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeFetchChunk(on.ID, nil)
}

func (c *MountPeerClient) invokeFetchChunk(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    MountPeerFetchChunkOrdinal,
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
			return p, nil, fmt.Errorf("MountPeer.FetchChunk: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("MountPeer.FetchChunk: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// MountPeerHandler is the server contract for the MountPeer service. Implement
// each method, then route requests to it with DispatchMountPeer.
//
// STREAMING: FetchChunk is server-streaming. Until the transport streaming
// primitive ships, the handler returns a single response body (one frame) and
// DispatchMountPeer routes it like a unary method.
type MountPeerHandler interface {
	ChunkAnnounce(req []byte) ([]byte, error)
	ChunkLookup(req []byte) ([]byte, error)
	FetchChunk(req []byte) ([]byte, error)
}

// DispatchMountPeer decodes a Call envelope, routes it by method ordinal to h,
// and returns the response envelope. An unknown ordinal yields a StatusNotFound
// response; a handler error yields StatusInternal.
func DispatchMountPeer(h MountPeerHandler, envelope []byte) ([]byte, error) {
	call, err := rpc.ParseRequest(envelope)
	if err != nil {
		return nil, err
	}
	switch call.Method {
	case MountPeerChunkAnnounceOrdinal:
		body, err := h.ChunkAnnounce(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case MountPeerChunkLookupOrdinal:
		body, err := h.ChunkLookup(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case MountPeerFetchChunkOrdinal:
		// STREAMING: unary-shaped dispatch of the first frame until the
		// transport streaming primitive ships.
		body, err := h.FetchChunk(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	default:
		return rpc.BuildResponse(rpc.StatusNotFound, call.PromiseID, nil), nil
	}
}
