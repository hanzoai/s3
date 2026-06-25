// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// Native ZAP adapter for the HanzoIdentityAccessManagement service — the
// hand-written glue that binds the generated Dispatch/Client (iam_zap.go) to a
// real IAM backend over the canonical github.com/zap-proto/go transport. This
// is the IAM analogue of s3/object/service.go (object Get/Put): the SAME
// Channel/Client/Handler/Dispatch shape, the SAME zero-copy doctrine — the
// request/response payloads ARE the iamwire message buffers, shipped over a
// ZAP socket correlated by PromiseID. No struct marshaling, no protobuf on the
// wire, no gRPC.
//
// Layering (decomplected — one concern per seam):
//
//   IamStore        the backend contract, in pure ZAP terms. Every method takes
//                   the opaque request buffer and returns the opaque response
//                   buffer (exactly the generated Handler contract). The real
//                   credential manager satisfies this — typically by Wrap-ing
//                   the request, translating to/from the in-process iam_pb model
//                   via the bridge_pb.go helpers, and New-ing the response. The
//                   wire package does NOT import that backend; integration wires
//                   it. Tests use an in-memory store.
//   handler         adapts an IamStore to HanzoIdentityAccessManagementHandler.
//   Dispatch/Serve  server side: bind handler to transport.Listen.
//   Client          typed client internal services hold; owns one transport.Conn.

package iamwire

import (
	"github.com/zap-proto/go/transport"
)

// IamStore is the backend the ZAP IAM service delegates to. Each method mirrors
// one HanzoIdentityAccessManagement RPC; both the argument and the result are
// the opaque iamwire message buffers (decode the request with the matching
// Wrap*, build the reply with the matching New*). The real S3 IAM backend (the
// filer-side credential manager) implements this; tests use an in-memory store.
//
// Keeping the seam in []byte terms — rather than in iam_pb or in the typed
// views — is deliberate: it is the exact contract the generated dispatcher
// already speaks, so the handler below is a pure pass-through and the backend
// is free to translate to whatever in-process model it likes (the bridge_pb.go
// pb<->wire helpers are provided for the common iam_pb case) without this file
// taking a dependency on it.
type IamStore interface {
	// Configuration management.
	GetConfiguration(req []byte) (resp []byte, err error)
	PutConfiguration(req []byte) (resp []byte, err error)

	// User management.
	CreateUser(req []byte) (resp []byte, err error)
	GetUser(req []byte) (resp []byte, err error)
	UpdateUser(req []byte) (resp []byte, err error)
	DeleteUser(req []byte) (resp []byte, err error)
	ListUsers(req []byte) (resp []byte, err error)

	// Access-key management.
	CreateAccessKey(req []byte) (resp []byte, err error)
	DeleteAccessKey(req []byte) (resp []byte, err error)
	GetUserByAccessKey(req []byte) (resp []byte, err error)

	// Policy management.
	PutPolicy(req []byte) (resp []byte, err error)
	GetPolicy(req []byte) (resp []byte, err error)
	ListPolicies(req []byte) (resp []byte, err error)
	DeletePolicy(req []byte) (resp []byte, err error)

	// Service-account management.
	CreateServiceAccount(req []byte) (resp []byte, err error)
	UpdateServiceAccount(req []byte) (resp []byte, err error)
	DeleteServiceAccount(req []byte) (resp []byte, err error)
	GetServiceAccount(req []byte) (resp []byte, err error)
	ListServiceAccounts(req []byte) (resp []byte, err error)
	GetServiceAccountByAccessKey(req []byte) (resp []byte, err error)
}

// handler adapts an IamStore to the generated
// HanzoIdentityAccessManagementHandler. Each method is a direct pass-through:
// the store already speaks the opaque-buffer contract, so there is nothing to
// translate here — translation (if any) lives in the store implementation, in
// one place, shared by every transport.
type handler struct{ store IamStore }

func (h handler) GetConfiguration(req []byte) ([]byte, error) { return h.store.GetConfiguration(req) }
func (h handler) PutConfiguration(req []byte) ([]byte, error) { return h.store.PutConfiguration(req) }
func (h handler) CreateUser(req []byte) ([]byte, error)       { return h.store.CreateUser(req) }
func (h handler) GetUser(req []byte) ([]byte, error)          { return h.store.GetUser(req) }
func (h handler) UpdateUser(req []byte) ([]byte, error)       { return h.store.UpdateUser(req) }
func (h handler) DeleteUser(req []byte) ([]byte, error)       { return h.store.DeleteUser(req) }
func (h handler) ListUsers(req []byte) ([]byte, error)        { return h.store.ListUsers(req) }
func (h handler) CreateAccessKey(req []byte) ([]byte, error)  { return h.store.CreateAccessKey(req) }
func (h handler) DeleteAccessKey(req []byte) ([]byte, error)  { return h.store.DeleteAccessKey(req) }
func (h handler) GetUserByAccessKey(req []byte) ([]byte, error) {
	return h.store.GetUserByAccessKey(req)
}
func (h handler) PutPolicy(req []byte) ([]byte, error)    { return h.store.PutPolicy(req) }
func (h handler) GetPolicy(req []byte) ([]byte, error)    { return h.store.GetPolicy(req) }
func (h handler) ListPolicies(req []byte) ([]byte, error) { return h.store.ListPolicies(req) }
func (h handler) DeletePolicy(req []byte) ([]byte, error) { return h.store.DeletePolicy(req) }
func (h handler) CreateServiceAccount(req []byte) ([]byte, error) {
	return h.store.CreateServiceAccount(req)
}
func (h handler) UpdateServiceAccount(req []byte) ([]byte, error) {
	return h.store.UpdateServiceAccount(req)
}
func (h handler) DeleteServiceAccount(req []byte) ([]byte, error) {
	return h.store.DeleteServiceAccount(req)
}
func (h handler) GetServiceAccount(req []byte) ([]byte, error) {
	return h.store.GetServiceAccount(req)
}
func (h handler) ListServiceAccounts(req []byte) ([]byte, error) {
	return h.store.ListServiceAccounts(req)
}
func (h handler) GetServiceAccountByAccessKey(req []byte) ([]byte, error) {
	return h.store.GetServiceAccountByAccessKey(req)
}

// Dispatch is the HanzoIdentityAccessManagement server dispatch bound to store —
// pass it to transport.Listen (or compose it behind a TLS/QUIC listener).
func Dispatch(store IamStore) transport.Dispatch {
	h := handler{store: store}
	return func(envelope []byte) ([]byte, error) {
		return DispatchHanzoIdentityAccessManagement(h, envelope)
	}
}

// Serve starts the native ZAP IAM service on network/addr (e.g. "tcp",
// ":18801"), backed by store. Returns the running server; Close stops it. For
// the PQ-secured mesh, build the listener via the transport TLS/QUIC path and
// use Dispatch(store) as its dispatch.
func Serve(network, addr string, store IamStore) (*transport.Server, error) {
	return transport.Listen(network, addr, Dispatch(store))
}

// Client is the typed IAM ZAP service client internal services hold. It owns the
// transport connection to one S3 IAM endpoint and decodes every reply into the
// in-process iam_pb model (via the bridge_pb.go decoders) so call sites get the
// same shapes the retired gRPC stub returned.
type Client struct {
	conn transport.Conn
	rpc  *HanzoIdentityAccessManagementClient
}

// Dial connects to the S3 IAM ZAP service at addr over plain TCP. For the
// PQ-secured mesh, establish the transport.Conn via the TLS/QUIC path and wrap
// it with NewClient instead.
func Dial(network, addr string) (*Client, error) {
	conn, err := transport.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, rpc: NewHanzoIdentityAccessManagementClient(conn, nil)}, nil
}

// NewClient wraps an already-established transport.Conn (TCP, Unix, or PQ-TLS).
func NewClient(conn transport.Conn) *Client {
	return &Client{conn: conn, rpc: NewHanzoIdentityAccessManagementClient(conn, nil)}
}

// RPC exposes the underlying generated client for callers that want to drive raw
// message buffers or pipelined ("On") calls directly.
func (c *Client) RPC() *HanzoIdentityAccessManagementClient { return c.rpc }

// Close releases the underlying connection.
func (c *Client) Close() error { return c.conn.Close() }
