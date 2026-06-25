// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// Native ZAP service adapter for HanzoS3LifecycleInternal — the hand-written
// glue that binds the generated Dispatch/Client (s3_lifecycle_zap.go) to a real
// engine backend. This is the lifecycle peer of s3/svc/object/service.go for the
// object service: requests/responses are the zero-copy s3_lifecyclewire message
// buffers, carried over the canonical github.com/zap-proto/go transport and
// dispatched by the generated HanzoS3LifecycleInternal service. No struct
// marshaling, no protobuf, no HTTP — the bytes ARE the message.
//
// The backend is expressed as a Go interface (LifecycleDeleter) in pure wire
// terms so this file imports neither a protobuf domain model nor the
// real filer/volume/master engine. The S3 server wires its live
// (re-fetch + identity-CAS + object-lock + dispatch) implementation to this
// interface at integration time (s3/s3api/s3api_internal_lifecycle.go); tests
// use an in-memory fake (see zapsvc_roundtrip_test.go). This is the single,
// pb-free binding — request/response cross the wire as zero-copy
// s3_lifecyclewire buffers end to end.

package s3_lifecyclewire

import (
	"github.com/zap-proto/go/transport"
)

// LifecycleDeleteResult is the value a LifecycleDeleter returns: the verdict
// outcome (a LifecycleDeleteOutcome* enum value) plus a human-readable reason.
// It mirrors the LifecycleDeleteResponse fields without forcing the backend to
// hand-build a wire buffer.
type LifecycleDeleteResult struct {
	Outcome uint32 // one of the LifecycleDeleteOutcome* constants
	Reason  string
}

// IsNull reports whether this EntryIdentity view is the null nested object
// (the request carried no CAS witness). Reading any field accessor on a null
// view dereferences a nil backing message and panics, so callers MUST gate on
// IsNull before touching MtimeNs()/HeadFid()/etc. Hand-written companion to
// the generated zero-copy view (the .proto codegen omits null checks).
func (t EntryIdentity) IsNull() bool { return t.o.IsNull() }

// HasExpectedIdentity reports whether the request carried a CAS witness
// (a non-null expected_identity). Convenience over ExpectedIdentity().IsNull().
func (t LifecycleDeleteRequest) HasExpectedIdentity() bool {
	return !t.ExpectedIdentity().IsNull()
}

// LifecycleDeleter is the backend the ZAP service delegates to. The real S3
// engine implements this: it re-fetches the live entry, verifies the
// EntryIdentity CAS, runs object-lock protections, and dispatches to the
// internal helper for the given ActionKind. Tests use an in-memory fake.
//
// req is the zero-copy view over the wire LifecycleDeleteRequest; read its
// fields with the accessor methods (Bucket(), ObjectPath(), ActionKind(),
// ExpectedIdentity(), …). The view aliases the request buffer, so copy out any
// bytes the implementation needs to retain beyond the call.
type LifecycleDeleter interface {
	LifecycleDelete(req LifecycleDeleteRequest) (LifecycleDeleteResult, error)
}

// handler adapts a LifecycleDeleter to the generated
// HanzoS3LifecycleInternalHandler: it Wraps the request buffer into its
// zero-copy view, calls the backend, and builds the response buffer directly.
type handler struct{ deleter LifecycleDeleter }

func (h handler) LifecycleDelete(req []byte) ([]byte, error) {
	v, err := WrapLifecycleDeleteRequest(req)
	if err != nil {
		return nil, err
	}
	res, err := h.deleter.LifecycleDelete(v)
	if err != nil {
		return nil, err
	}
	return NewLifecycleDeleteResponse(LifecycleDeleteResponseInput{
		Outcome: res.Outcome,
		Reason:  res.Reason,
	}), nil
}

// Dispatch is the HanzoS3LifecycleInternal server dispatch bound to deleter —
// pass it to transport.Listen. Exposed so callers can compose it (e.g. behind
// PQ-TLS).
func Dispatch(deleter LifecycleDeleter) transport.Dispatch {
	h := handler{deleter: deleter}
	return func(envelope []byte) ([]byte, error) {
		return DispatchHanzoS3LifecycleInternal(h, envelope)
	}
}

// Serve starts the native ZAP lifecycle service on network/addr (e.g. "tcp",
// ":18902"), backed by deleter. Returns the running server; Close stops it.
func Serve(network, addr string, deleter LifecycleDeleter) (*transport.Server, error) {
	return transport.Listen(network, addr, Dispatch(deleter))
}

// Client is the typed lifecycle ZAP service client internal callers (the
// lifecycle worker / dailyrun pipeline) hold. It owns the transport connection
// to one S3 endpoint.
type Client struct {
	conn transport.Conn
	rpc  *HanzoS3LifecycleInternalClient
}

// Dial connects to the S3 lifecycle ZAP service at addr (e.g.
// "s3.hanzo.svc:18902") over plain TCP. For the PQ-secured mesh, establish a
// transport.Conn via transport.DialTLS with a transport.PQTLSConfig and use
// NewClient.
func Dial(network, addr string) (*Client, error) {
	conn, err := transport.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, rpc: NewHanzoS3LifecycleInternalClient(conn, nil)}, nil
}

// NewClient wraps an already-established transport.Conn (TCP, Unix, or PQ-TLS).
func NewClient(conn transport.Conn) *Client {
	return &Client{conn: conn, rpc: NewHanzoS3LifecycleInternalClient(conn, nil)}
}

// Close releases the underlying connection.
func (c *Client) Close() error { return c.conn.Close() }

// LifecycleDelete issues one (rule, action) verdict over ZAP and returns the
// server's outcome. in carries the routing, stream context, CAS witness, and
// snapshot id; the returned result holds the outcome enum and reason string.
func (c *Client) LifecycleDelete(in LifecycleDeleteRequestInput) (LifecycleDeleteResult, error) {
	_, body, err := c.rpc.LifecycleDelete(NewLifecycleDeleteRequest(in))
	if err != nil {
		return LifecycleDeleteResult{}, err
	}
	resp, err := WrapLifecycleDeleteResponse(body)
	if err != nil {
		return LifecycleDeleteResult{}, err
	}
	return LifecycleDeleteResult{Outcome: resp.Outcome(), Reason: resp.Reason()}, nil
}
