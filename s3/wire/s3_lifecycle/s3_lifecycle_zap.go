// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	"fmt"

	"github.com/zap-proto/go/rpc"
)

// Method ordinals for the HanzoS3LifecycleInternal service (stable 1-based wire
// ids, in s3_lifecycle.proto declaration order).
const (
	HanzoS3LifecycleInternalLifecycleDeleteOrdinal uint32 = 1
)

// HanzoS3LifecycleInternalChannel ships one Call envelope and awaits its
// correlated Response.
type HanzoS3LifecycleInternalChannel interface {
	Call(envelope []byte) (rpc.Response, error)
}

// HanzoS3LifecycleInternalClient is a typed RPC client for the
// HanzoS3LifecycleInternal service over a ZAP call channel. Each call takes a
// fresh PromiseID from sess; the pipelined "On" form of a method sets Target to
// a prior call's Promise so the server chains them.
type HanzoS3LifecycleInternalClient struct {
	ch   HanzoS3LifecycleInternalChannel
	cap  []byte
	sess *rpc.Session
}

// NewHanzoS3LifecycleInternalClient returns a client that issues calls over ch,
// attaching cap (which may be nil) to every request.
func NewHanzoS3LifecycleInternalClient(ch HanzoS3LifecycleInternalChannel, capability []byte) *HanzoS3LifecycleInternalClient {
	return &HanzoS3LifecycleInternalClient{ch: ch, cap: capability, sess: rpc.NewSession()}
}

func (c *HanzoS3LifecycleInternalClient) LifecycleDelete(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeLifecycleDelete(rpc.NoTarget, req)
}

// LifecycleDeleteOn issues LifecycleDelete as a dependent call pipelined on the
// answer of on: the server substitutes on's resolved result for this call's
// payload before dispatch, so it ships without waiting for on to round-trip.
func (c *HanzoS3LifecycleInternalClient) LifecycleDeleteOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeLifecycleDelete(on.ID, nil)
}

func (c *HanzoS3LifecycleInternalClient) invokeLifecycleDelete(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    HanzoS3LifecycleInternalLifecycleDeleteOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("HanzoS3LifecycleInternal.LifecycleDelete: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// HanzoS3LifecycleInternalHandler is the server contract for the
// HanzoS3LifecycleInternal service. Implement each method, then route requests
// to it with DispatchHanzoS3LifecycleInternal.
type HanzoS3LifecycleInternalHandler interface {
	LifecycleDelete(req []byte) ([]byte, error)
}

// DispatchHanzoS3LifecycleInternal decodes a Call envelope, routes it by method
// ordinal to h, and returns the response envelope. An unknown ordinal yields a
// StatusNotFound response; a handler error yields StatusInternal.
func DispatchHanzoS3LifecycleInternal(h HanzoS3LifecycleInternalHandler, envelope []byte) ([]byte, error) {
	call, err := rpc.ParseRequest(envelope)
	if err != nil {
		return nil, err
	}
	switch call.Method {
	case HanzoS3LifecycleInternalLifecycleDeleteOrdinal:
		body, err := h.LifecycleDelete(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	default:
		return rpc.BuildResponse(rpc.StatusNotFound, call.PromiseID, nil), nil
	}
}
