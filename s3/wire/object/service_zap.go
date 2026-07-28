// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// HanzoS3Object service — object-level Get/Put over the canonical ZAP RPC
// envelope (github.com/zap-proto/go/rpc). Same Channel/Client/Handler/Dispatch
// shape the iam/filer/master service wire emits. Request/response payloads are
// the opaque objectwire message buffers above; the transport ships them as ZAP
// bytes correlated by PromiseID.

package objectwire

import (
	"fmt"

	"github.com/zap-proto/go/rpc"
)

// Method ordinals for the HanzoS3Object service (stable 1-based wire ids).
const (
	// GetObject: GetObjectRequest -> GetObjectResponse.
	HanzoS3ObjectGetObjectOrdinal uint32 = 1
	// PutObject: PutObjectRequest -> PutObjectResponse.
	HanzoS3ObjectPutObjectOrdinal uint32 = 2
)

// HanzoS3ObjectChannel ships one Call envelope and awaits its correlated
// Response. transport.Conn satisfies this.
type HanzoS3ObjectChannel interface {
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

// HanzoS3ObjectClient is a typed RPC client for the HanzoS3Object service over a
// ZAP call channel. Each call takes a fresh PromiseID from sess.
type HanzoS3ObjectClient struct {
	ch   HanzoS3ObjectChannel
	cap  []byte
}

// NewHanzoS3ObjectClient returns a client that issues calls over ch, attaching
// capability (which may be nil) to every request.
func NewHanzoS3ObjectClient(ch HanzoS3ObjectChannel, capability []byte) *HanzoS3ObjectClient {
	return &HanzoS3ObjectClient{ch: ch, cap: capability}
}

func (c *HanzoS3ObjectClient) GetObject(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeGetObject(rpc.NoTarget, req)
}

// GetObjectOn issues GetObject as a dependent call pipelined on the answer of on.
func (c *HanzoS3ObjectClient) GetObjectOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeGetObject(on.ID, nil)
}

func (c *HanzoS3ObjectClient) invokeGetObject(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    HanzoS3ObjectGetObjectOrdinal,
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
			return p, nil, fmt.Errorf("HanzoS3Object.GetObject: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("HanzoS3Object.GetObject: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *HanzoS3ObjectClient) PutObject(req []byte) (rpc.Promise, []byte, error) {
	return c.invokePutObject(rpc.NoTarget, req)
}

// PutObjectOn issues PutObject as a dependent call pipelined on the answer of on.
func (c *HanzoS3ObjectClient) PutObjectOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokePutObject(on.ID, nil)
}

func (c *HanzoS3ObjectClient) invokePutObject(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    HanzoS3ObjectPutObjectOrdinal,
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
			return p, nil, fmt.Errorf("HanzoS3Object.PutObject: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("HanzoS3Object.PutObject: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// HanzoS3ObjectHandler is the server contract for the HanzoS3Object service.
// Each req/return value is an opaque objectwire buffer; decode with Wrap* and
// build the reply with New*.
type HanzoS3ObjectHandler interface {
	GetObject(req []byte) ([]byte, error)
	PutObject(req []byte) ([]byte, error)
}

// DispatchHanzoS3Object decodes a Call envelope, routes it by method ordinal to
// h, and returns the response envelope. An unknown ordinal yields StatusNotFound;
// a handler error yields StatusInternal.
func DispatchHanzoS3Object(h HanzoS3ObjectHandler, envelope []byte) ([]byte, error) {
	call, err := rpc.ParseRequest(envelope)
	if err != nil {
		return nil, err
	}
	switch call.Method {
	case HanzoS3ObjectGetObjectOrdinal:
		body, err := h.GetObject(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case HanzoS3ObjectPutObjectOrdinal:
		body, err := h.PutObject(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	default:
		return rpc.BuildResponse(rpc.StatusNotFound, call.PromiseID, nil), nil
	}
}
