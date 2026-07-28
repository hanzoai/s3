// Code generated from worker.proto; DO NOT EDIT.

package workerwire

import (
	"context"
	"fmt"

	"github.com/zap-proto/go/rpc"
)

// Method ordinals for the WorkerService service (stable 1-based wire ids,
// in worker.proto declaration order).
const (
	WorkerServiceWorkerStreamOrdinal uint32 = 1
)

// WorkerServiceChannel ships one Call envelope and awaits its correlated Response.
// CallContext is Call that also aborts when ctx is done (transport.Conn satisfies
// both).
type WorkerServiceChannel interface {
	Call(envelope []byte) (rpc.Response, error)
	CallContext(ctx context.Context, envelope []byte) (rpc.Response, error)
	// NextPromiseID allocates a call's PromiseID from the CONNECTION rather than
	// from a per-client counter. PromiseID is the connection's demultiplexing
	// key: the transport keys its in-flight table by it, and OpenStream already
	// allocates stream IDs the same way. A client that numbered its own calls
	// would restart at 1 on every construction, so two clients sharing a pooled
	// conn collide — the transport overwrites the first waiter's slot and its
	// response is dropped, blocking that caller for as long as its context runs.
	NextPromiseID() uint32
}

// WorkerServiceClient is a typed RPC client for the WorkerService service over a
// ZAP call channel. Each call takes a fresh PromiseID from sess; the pipelined
// "On" form of a method sets Target to a prior call's Promise so the server
// chains them.
type WorkerServiceClient struct {
	ch   WorkerServiceChannel
	cap  []byte
}

// NewWorkerServiceClient returns a client that issues calls over ch, attaching
// cap (which may be nil) to every request.
func NewWorkerServiceClient(ch WorkerServiceChannel, capability []byte) *WorkerServiceClient {
	return &WorkerServiceClient{ch: ch, cap: capability}
}

// WorkerStream is the bidirectional worker<->admin stream: a stream of
// WorkerMessage requests multiplexed against a stream of AdminMessage responses.
//
// STREAMING: this is a bidi-streaming RPC (rpc WorkerStream(stream WorkerMessage)
// returns (stream AdminMessage)). The single-shot form below ships one
// WorkerMessage envelope (req, built with NewWorkerMessage) and returns one
// AdminMessage envelope (Wrap with WrapAdminMessage); the full duplex stream
// body lands when the transport streaming primitive ships. The request and
// response message schemas are already complete.
func (c *WorkerServiceClient) WorkerStream(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeWorkerStream(ctx, rpc.NoTarget, req)
}

// WorkerStreamOn issues WorkerStream as a dependent call pipelined on the answer
// of on: the server substitutes on's resolved result for this call's payload
// before dispatch, so it ships without waiting for on to round-trip.
//
// STREAMING: see WorkerStream.
func (c *WorkerServiceClient) WorkerStreamOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeWorkerStream(ctx, on.ID, nil)
}

func (c *WorkerServiceClient) invokeWorkerStream(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
		Method:    WorkerServiceWorkerStreamOrdinal,
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
			return p, nil, fmt.Errorf("WorkerService.WorkerStream: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("WorkerService.WorkerStream: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// WorkerServiceHandler is the server contract for the WorkerService service.
// Implement each method, then route requests to it with DispatchWorkerService.
//
// STREAMING: WorkerStream receives one WorkerMessage envelope (Wrap with
// WrapWorkerMessage) and returns one AdminMessage envelope (build with
// NewAdminMessage). The full duplex stream contract lands when the transport
// streaming primitive ships.
type WorkerServiceHandler interface {
	WorkerStream(req []byte) ([]byte, error)
}

// DispatchWorkerService decodes a Call envelope, routes it by method ordinal to
// h, and returns the response envelope. An unknown ordinal yields a
// StatusNotFound response; a handler error yields StatusInternal.
func DispatchWorkerService(h WorkerServiceHandler, envelope []byte) ([]byte, error) {
	call, err := rpc.ParseRequest(envelope)
	if err != nil {
		return nil, err
	}
	switch call.Method {
	case WorkerServiceWorkerStreamOrdinal:
		// STREAMING: single-shot dispatch of one WorkerMessage -> one
		// AdminMessage; the duplex stream loop lands with the transport
		// streaming primitive.
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
