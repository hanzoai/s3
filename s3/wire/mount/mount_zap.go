// Code generated from mount.proto; DO NOT EDIT.

package mountwire

import (
	"fmt"

	"github.com/zap-proto/go/rpc"
)

// Method ordinals for the HanzoMount service (stable 1-based wire ids,
// in mount.proto declaration order).
const (
	HanzoMountConfigureOrdinal uint32 = 1
)

// HanzoMountChannel ships one Call envelope and awaits its correlated Response.
type HanzoMountChannel interface {
	Call(envelope []byte) (rpc.Response, error)
}

// HanzoMountClient is a typed RPC client for the HanzoMount service over a ZAP
// call channel. Each call takes a fresh PromiseID from sess; the pipelined "On"
// form of a method sets Target to a prior call's Promise so the server chains
// them.
type HanzoMountClient struct {
	ch   HanzoMountChannel
	cap  []byte
	sess *rpc.Session
}

// NewHanzoMountClient returns a client that issues calls over ch, attaching cap
// (which may be nil) to every request.
func NewHanzoMountClient(ch HanzoMountChannel, capability []byte) *HanzoMountClient {
	return &HanzoMountClient{ch: ch, cap: capability, sess: rpc.NewSession()}
}

func (c *HanzoMountClient) Configure(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeConfigure(rpc.NoTarget, req)
}

// ConfigureOn issues Configure as a dependent call pipelined on the answer of on:
// the server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *HanzoMountClient) ConfigureOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeConfigure(on.ID, nil)
}

func (c *HanzoMountClient) invokeConfigure(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    HanzoMountConfigureOrdinal,
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
			return p, nil, fmt.Errorf("HanzoMount.Configure: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("HanzoMount.Configure: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// HanzoMountHandler is the server contract for the HanzoMount service. Implement
// each method, then route requests to it with DispatchHanzoMount.
type HanzoMountHandler interface {
	Configure(req []byte) ([]byte, error)
}

// DispatchHanzoMount decodes a Call envelope, routes it by method ordinal to h,
// and returns the response envelope. An unknown ordinal yields a StatusNotFound
// response; a handler error yields StatusInternal.
func DispatchHanzoMount(h HanzoMountHandler, envelope []byte) ([]byte, error) {
	call, err := rpc.ParseRequest(envelope)
	if err != nil {
		return nil, err
	}
	switch call.Method {
	case HanzoMountConfigureOrdinal:
		body, err := h.Configure(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	default:
		return rpc.BuildResponse(rpc.StatusNotFound, call.PromiseID, nil), nil
	}
}
