// Code generated from plugin.proto; DO NOT EDIT.

package pluginwire

import (
	"context"
	"fmt"

	"github.com/zap-proto/go/rpc"
)

// Method ordinals for the PluginControlService service (stable 1-based wire
// ids, in plugin.proto declaration order).
const (
	PluginControlServiceWorkerStreamOrdinal uint32 = 1
)

// PluginControlServiceChannel ships one Call envelope and awaits its correlated
// Response. CallContext is Call that also aborts when ctx is done (transport.Conn
// satisfies both).
type PluginControlServiceChannel interface {
	Call(envelope []byte) (rpc.Response, error)
	CallContext(ctx context.Context, envelope []byte) (rpc.Response, error)
}

// PluginControlServiceClient is a typed RPC client for the PluginControlService
// service over a ZAP call channel. Each call takes a fresh PromiseID from sess;
// the pipelined "On" form of a method sets Target to a prior call's Promise so
// the server chains them.
type PluginControlServiceClient struct {
	ch   PluginControlServiceChannel
	cap  []byte
	sess *rpc.Session
}

// NewPluginControlServiceClient returns a client that issues calls over ch,
// attaching cap (which may be nil) to every request.
func NewPluginControlServiceClient(ch PluginControlServiceChannel, capability []byte) *PluginControlServiceClient {
	return &PluginControlServiceClient{ch: ch, cap: capability, sess: rpc.NewSession()}
}

// WorkerStream is a bidirectional streaming RPC: a stream of WorkerToAdminMessage
// frames in, a stream of AdminToWorkerMessage frames out.
//
// STREAMING: the duplex transport primitive has not shipped yet; this single-
// frame form ships one WorkerToAdminMessage and awaits one AdminToWorkerMessage
// (the per-frame unit the streaming transport will multiplex). The full
// stream-oriented body lands when the transport streaming primitive ships.
func (c *PluginControlServiceClient) WorkerStream(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeWorkerStream(ctx, rpc.NoTarget, req)
}

// WorkerStreamOn issues WorkerStream as a dependent call pipelined on the answer
// of on: the server substitutes on's resolved result for this call's payload
// before dispatch, so it ships without waiting for on to round-trip.
//
// STREAMING: see WorkerStream.
func (c *PluginControlServiceClient) WorkerStreamOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeWorkerStream(ctx, on.ID, nil)
}

func (c *PluginControlServiceClient) invokeWorkerStream(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
		Method:    PluginControlServiceWorkerStreamOrdinal,
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
			return p, nil, fmt.Errorf("PluginControlService.WorkerStream: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("PluginControlService.WorkerStream: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// PluginControlServiceHandler is the server contract for the
// PluginControlService service. Implement each method, then route requests to it
// with DispatchPluginControlService.
type PluginControlServiceHandler interface {
	// WorkerStream handles one WorkerToAdminMessage frame and returns one
	// AdminToWorkerMessage frame.
	//
	// STREAMING: the duplex transport primitive has not shipped yet; this
	// single-frame form is the per-frame unit the streaming transport will
	// multiplex once it lands.
	WorkerStream(req []byte) ([]byte, error)
}

// DispatchPluginControlService decodes a Call envelope, routes it by method
// ordinal to h, and returns the response envelope. An unknown ordinal yields a
// StatusNotFound response; a handler error yields StatusInternal.
func DispatchPluginControlService(h PluginControlServiceHandler, envelope []byte) ([]byte, error) {
	call, err := rpc.ParseRequest(envelope)
	if err != nil {
		return nil, err
	}
	switch call.Method {
	case PluginControlServiceWorkerStreamOrdinal:
		// STREAMING: single-frame dispatch until the duplex transport ships.
		body, err := h.WorkerStream(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	default:
		return rpc.BuildResponse(rpc.StatusNotFound, call.PromiseID, nil), nil
	}
}
