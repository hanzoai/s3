// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// Native ZAP adapter for the HanzoMount service — the hand-written glue that
// binds the generated DispatchHanzoMount/HanzoMountClient (mount_zap.go) to a
// real mount backend. Mirrors s3/zapsvc/service.go exactly: a backend interface
// in engine terms, a handler that adapts it to the generated HanzoMountHandler,
// Dispatch/Serve for the server side, and Dial/NewClient/Client for callers.
//
// The bytes ARE the message: requests/responses cross the wire as zero-copy
// mountwire buffers over the canonical github.com/zap-proto/go transport — no
// struct, no marshal, no protobuf, no gRPC. The real S3 mount endpoint
// (HanzoMount gRPC server) implements Mounter; tests use an in-memory stub.

package mountwire

import (
	"github.com/zap-proto/go/transport"
)

// Mounter is the backend the ZAP mount service delegates to. It is the engine's
// HanzoMount.Configure in Go terms; the real mount server implements it, while
// tests supply an in-memory stub. We intentionally do NOT import the filer /
// master / volume backend here — integration wires a concrete Mounter later.
type Mounter interface {
	// Configure applies the mount configuration. collectionCapacity mirrors
	// the proto ConfigureRequest.collection_capacity field. ConfigureResponse
	// carries no fields, so a successful call returns only nil.
	Configure(collectionCapacity int64) error
}

// handler adapts a Mounter to the generated HanzoMountHandler: it Wraps each
// request buffer into its zero-copy view, calls the backend, and builds the
// response buffer directly.
type handler struct{ mounter Mounter }

func (h handler) Configure(req []byte) ([]byte, error) {
	v, err := WrapConfigureRequest(req)
	if err != nil {
		return nil, err
	}
	if err := h.mounter.Configure(v.CollectionCapacity()); err != nil {
		return nil, err
	}
	return NewConfigureResponse(ConfigureResponseInput{}), nil
}

// Dispatch is the HanzoMount server dispatch bound to mounter — pass it to
// transport.Listen. Exposed so callers can compose it (e.g. behind PQ-TLS).
func Dispatch(mounter Mounter) transport.Dispatch {
	h := handler{mounter: mounter}
	return func(envelope []byte) ([]byte, error) {
		return DispatchHanzoMount(h, envelope)
	}
}

// Serve starts the native ZAP HanzoMount service on network/addr (e.g. "tcp",
// ":18902"), backed by mounter. Returns the running server; Close stops it.
func Serve(network, addr string, mounter Mounter) (*transport.Server, error) {
	return transport.Listen(network, addr, Dispatch(mounter))
}

// Client is the typed HanzoMount ZAP service client internal services hold. It
// owns the transport connection to one mount endpoint.
type Client struct {
	conn *transport.Conn
	rpc  *HanzoMountClient
}

// Dial connects to the HanzoMount ZAP service at addr over plain TCP. For the
// PQ-secured mesh use transport.DialTLS with transport.PQTLSConfig and wrap the
// resulting *transport.Conn with NewClient.
func Dial(network, addr string) (*Client, error) {
	conn, err := transport.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, rpc: NewHanzoMountClient(conn, nil)}, nil
}

// NewClient wraps an already-established transport.Conn (TCP, Unix, or PQ-TLS).
func NewClient(conn *transport.Conn) *Client {
	return &Client{conn: conn, rpc: NewHanzoMountClient(conn, nil)}
}

// Close releases the underlying connection.
func (c *Client) Close() error { return c.conn.Close() }

// Configure applies a mount configuration over ZAP. ConfigureResponse carries no
// fields, so a successful call returns only nil.
func (c *Client) Configure(collectionCapacity int64) error {
	req := NewConfigureRequest(ConfigureRequestInput{CollectionCapacity: collectionCapacity})
	_, body, err := c.rpc.Configure(req)
	if err != nil {
		return err
	}
	// Decode the (empty) response to validate the wire frame before returning.
	if _, err := WrapConfigureResponse(body); err != nil {
		return err
	}
	return nil
}
