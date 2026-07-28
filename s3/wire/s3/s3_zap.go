// Code generated from s3.proto; DO NOT EDIT.

package s3wire

import (
	"fmt"

	"github.com/zap-proto/go/rpc"
)

// Method ordinals for the HanzoS3IamCache service (stable 1-based wire ids,
// in s3.proto declaration order).
//
// Every method carries iam wire messages (package iamwire, at
// github.com/hanzoai/s3/s3/wire/iam). The request/response payloads cross the
// transport as opaque ZAP buffers; build and read them with the iamwire
// constructors named per method below. The iam types are NOT redefined here.
const (
	// PutIdentity: iamwire.PutIdentityRequest -> iamwire.PutIdentityResponse.
	HanzoS3IamCachePutIdentityOrdinal uint32 = 1
	// RemoveIdentity: iamwire.RemoveIdentityRequest -> iamwire.RemoveIdentityResponse.
	HanzoS3IamCacheRemoveIdentityOrdinal uint32 = 2
	// PutPolicy: iamwire.PutPolicyRequest -> iamwire.PutPolicyResponse.
	HanzoS3IamCachePutPolicyOrdinal uint32 = 3
	// GetPolicy: iamwire.GetPolicyRequest -> iamwire.GetPolicyResponse.
	HanzoS3IamCacheGetPolicyOrdinal uint32 = 4
	// ListPolicies: iamwire.ListPoliciesRequest -> iamwire.ListPoliciesResponse.
	HanzoS3IamCacheListPoliciesOrdinal uint32 = 5
	// DeletePolicy: iamwire.DeletePolicyRequest -> iamwire.DeletePolicyResponse.
	HanzoS3IamCacheDeletePolicyOrdinal uint32 = 6
	// PutGroup: iamwire.PutGroupRequest -> iamwire.PutGroupResponse.
	HanzoS3IamCachePutGroupOrdinal uint32 = 7
	// RemoveGroup: iamwire.RemoveGroupRequest -> iamwire.RemoveGroupResponse.
	HanzoS3IamCacheRemoveGroupOrdinal uint32 = 8
)

// HanzoS3IamCacheChannel ships one Call envelope and awaits its correlated
// Response.
type HanzoS3IamCacheChannel interface {
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

// HanzoS3IamCacheClient is a typed RPC client for the HanzoS3IamCache service
// over a ZAP call channel. Each call takes a fresh PromiseID from sess; the
// pipelined "On" form of a method sets Target to a prior call's Promise so the
// server chains them.
type HanzoS3IamCacheClient struct {
	ch   HanzoS3IamCacheChannel
	cap  []byte
}

// NewHanzoS3IamCacheClient returns a client that issues calls over ch,
// attaching cap (which may be nil) to every request.
func NewHanzoS3IamCacheClient(ch HanzoS3IamCacheChannel, capability []byte) *HanzoS3IamCacheClient {
	return &HanzoS3IamCacheClient{ch: ch, cap: capability}
}

func (c *HanzoS3IamCacheClient) PutIdentity(req []byte) (rpc.Promise, []byte, error) {
	return c.invokePutIdentity(rpc.NoTarget, req)
}

// PutIdentityOn issues PutIdentity as a dependent call pipelined on the answer
// of on: the server substitutes on's resolved result for this call's payload
// before dispatch, so it ships without waiting for on to round-trip.
func (c *HanzoS3IamCacheClient) PutIdentityOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokePutIdentity(on.ID, nil)
}

func (c *HanzoS3IamCacheClient) invokePutIdentity(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    HanzoS3IamCachePutIdentityOrdinal,
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
			return p, nil, fmt.Errorf("HanzoS3IamCache.PutIdentity: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("HanzoS3IamCache.PutIdentity: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *HanzoS3IamCacheClient) RemoveIdentity(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeRemoveIdentity(rpc.NoTarget, req)
}

// RemoveIdentityOn issues RemoveIdentity as a dependent call pipelined on the
// answer of on: the server substitutes on's resolved result for this call's
// payload before dispatch, so it ships without waiting for on to round-trip.
func (c *HanzoS3IamCacheClient) RemoveIdentityOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeRemoveIdentity(on.ID, nil)
}

func (c *HanzoS3IamCacheClient) invokeRemoveIdentity(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    HanzoS3IamCacheRemoveIdentityOrdinal,
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
			return p, nil, fmt.Errorf("HanzoS3IamCache.RemoveIdentity: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("HanzoS3IamCache.RemoveIdentity: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *HanzoS3IamCacheClient) PutPolicy(req []byte) (rpc.Promise, []byte, error) {
	return c.invokePutPolicy(rpc.NoTarget, req)
}

// PutPolicyOn issues PutPolicy as a dependent call pipelined on the answer of
// on: the server substitutes on's resolved result for this call's payload
// before dispatch, so it ships without waiting for on to round-trip.
func (c *HanzoS3IamCacheClient) PutPolicyOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokePutPolicy(on.ID, nil)
}

func (c *HanzoS3IamCacheClient) invokePutPolicy(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    HanzoS3IamCachePutPolicyOrdinal,
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
			return p, nil, fmt.Errorf("HanzoS3IamCache.PutPolicy: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("HanzoS3IamCache.PutPolicy: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *HanzoS3IamCacheClient) GetPolicy(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeGetPolicy(rpc.NoTarget, req)
}

// GetPolicyOn issues GetPolicy as a dependent call pipelined on the answer of
// on: the server substitutes on's resolved result for this call's payload
// before dispatch, so it ships without waiting for on to round-trip.
func (c *HanzoS3IamCacheClient) GetPolicyOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeGetPolicy(on.ID, nil)
}

func (c *HanzoS3IamCacheClient) invokeGetPolicy(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    HanzoS3IamCacheGetPolicyOrdinal,
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
			return p, nil, fmt.Errorf("HanzoS3IamCache.GetPolicy: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("HanzoS3IamCache.GetPolicy: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *HanzoS3IamCacheClient) ListPolicies(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeListPolicies(rpc.NoTarget, req)
}

// ListPoliciesOn issues ListPolicies as a dependent call pipelined on the answer
// of on: the server substitutes on's resolved result for this call's payload
// before dispatch, so it ships without waiting for on to round-trip.
func (c *HanzoS3IamCacheClient) ListPoliciesOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeListPolicies(on.ID, nil)
}

func (c *HanzoS3IamCacheClient) invokeListPolicies(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    HanzoS3IamCacheListPoliciesOrdinal,
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
			return p, nil, fmt.Errorf("HanzoS3IamCache.ListPolicies: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("HanzoS3IamCache.ListPolicies: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *HanzoS3IamCacheClient) DeletePolicy(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeDeletePolicy(rpc.NoTarget, req)
}

// DeletePolicyOn issues DeletePolicy as a dependent call pipelined on the answer
// of on: the server substitutes on's resolved result for this call's payload
// before dispatch, so it ships without waiting for on to round-trip.
func (c *HanzoS3IamCacheClient) DeletePolicyOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeDeletePolicy(on.ID, nil)
}

func (c *HanzoS3IamCacheClient) invokeDeletePolicy(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    HanzoS3IamCacheDeletePolicyOrdinal,
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
			return p, nil, fmt.Errorf("HanzoS3IamCache.DeletePolicy: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("HanzoS3IamCache.DeletePolicy: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *HanzoS3IamCacheClient) PutGroup(req []byte) (rpc.Promise, []byte, error) {
	return c.invokePutGroup(rpc.NoTarget, req)
}

// PutGroupOn issues PutGroup as a dependent call pipelined on the answer of on:
// the server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *HanzoS3IamCacheClient) PutGroupOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokePutGroup(on.ID, nil)
}

func (c *HanzoS3IamCacheClient) invokePutGroup(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    HanzoS3IamCachePutGroupOrdinal,
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
			return p, nil, fmt.Errorf("HanzoS3IamCache.PutGroup: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("HanzoS3IamCache.PutGroup: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *HanzoS3IamCacheClient) RemoveGroup(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeRemoveGroup(rpc.NoTarget, req)
}

// RemoveGroupOn issues RemoveGroup as a dependent call pipelined on the answer
// of on: the server substitutes on's resolved result for this call's payload
// before dispatch, so it ships without waiting for on to round-trip.
func (c *HanzoS3IamCacheClient) RemoveGroupOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeRemoveGroup(on.ID, nil)
}

func (c *HanzoS3IamCacheClient) invokeRemoveGroup(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    HanzoS3IamCacheRemoveGroupOrdinal,
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
			return p, nil, fmt.Errorf("HanzoS3IamCache.RemoveGroup: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("HanzoS3IamCache.RemoveGroup: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// HanzoS3IamCacheHandler is the server contract for the HanzoS3IamCache service.
// Implement each method, then route requests to it with DispatchHanzoS3IamCache.
// Each req/return value is an opaque iam wire buffer (package iamwire); decode
// it with the matching iamwire.Wrap* and build the reply with iamwire.New*.
type HanzoS3IamCacheHandler interface {
	PutIdentity(req []byte) ([]byte, error)
	RemoveIdentity(req []byte) ([]byte, error)
	PutPolicy(req []byte) ([]byte, error)
	GetPolicy(req []byte) ([]byte, error)
	ListPolicies(req []byte) ([]byte, error)
	DeletePolicy(req []byte) ([]byte, error)
	PutGroup(req []byte) ([]byte, error)
	RemoveGroup(req []byte) ([]byte, error)
}

// DispatchHanzoS3IamCache decodes a Call envelope, routes it by method ordinal
// to h, and returns the response envelope. An unknown ordinal yields a
// StatusNotFound response; a handler error yields StatusInternal.
func DispatchHanzoS3IamCache(h HanzoS3IamCacheHandler, envelope []byte) ([]byte, error) {
	call, err := rpc.ParseRequest(envelope)
	if err != nil {
		return nil, err
	}
	switch call.Method {
	case HanzoS3IamCachePutIdentityOrdinal:
		body, err := h.PutIdentity(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case HanzoS3IamCacheRemoveIdentityOrdinal:
		body, err := h.RemoveIdentity(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case HanzoS3IamCachePutPolicyOrdinal:
		body, err := h.PutPolicy(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case HanzoS3IamCacheGetPolicyOrdinal:
		body, err := h.GetPolicy(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case HanzoS3IamCacheListPoliciesOrdinal:
		body, err := h.ListPolicies(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case HanzoS3IamCacheDeletePolicyOrdinal:
		body, err := h.DeletePolicy(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case HanzoS3IamCachePutGroupOrdinal:
		body, err := h.PutGroup(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case HanzoS3IamCacheRemoveGroupOrdinal:
		body, err := h.RemoveGroup(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	default:
		return rpc.BuildResponse(rpc.StatusNotFound, call.PromiseID, nil), nil
	}
}
