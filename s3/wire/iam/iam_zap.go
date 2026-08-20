// Code generated from iam.proto; DO NOT EDIT.

package iamwire

import (
	"context"
	"fmt"

	"github.com/zap-proto/go/rpc"
)

// Method ordinals for the HanzoIdentityAccessManagement service (stable 1-based
// wire ids, in iam.proto declaration order).
const (
	HanzoIdentityAccessManagementGetConfigurationOrdinal             uint32 = 1
	HanzoIdentityAccessManagementPutConfigurationOrdinal             uint32 = 2
	HanzoIdentityAccessManagementCreateUserOrdinal                   uint32 = 3
	HanzoIdentityAccessManagementGetUserOrdinal                      uint32 = 4
	HanzoIdentityAccessManagementUpdateUserOrdinal                   uint32 = 5
	HanzoIdentityAccessManagementDeleteUserOrdinal                   uint32 = 6
	HanzoIdentityAccessManagementListUsersOrdinal                    uint32 = 7
	HanzoIdentityAccessManagementCreateAccessKeyOrdinal              uint32 = 8
	HanzoIdentityAccessManagementDeleteAccessKeyOrdinal              uint32 = 9
	HanzoIdentityAccessManagementGetUserByAccessKeyOrdinal           uint32 = 10
	HanzoIdentityAccessManagementPutPolicyOrdinal                    uint32 = 11
	HanzoIdentityAccessManagementGetPolicyOrdinal                    uint32 = 12
	HanzoIdentityAccessManagementListPoliciesOrdinal                 uint32 = 13
	HanzoIdentityAccessManagementDeletePolicyOrdinal                 uint32 = 14
	HanzoIdentityAccessManagementCreateServiceAccountOrdinal         uint32 = 15
	HanzoIdentityAccessManagementUpdateServiceAccountOrdinal         uint32 = 16
	HanzoIdentityAccessManagementDeleteServiceAccountOrdinal         uint32 = 17
	HanzoIdentityAccessManagementGetServiceAccountOrdinal            uint32 = 18
	HanzoIdentityAccessManagementListServiceAccountsOrdinal          uint32 = 19
	HanzoIdentityAccessManagementGetServiceAccountByAccessKeyOrdinal uint32 = 20
)

// HanzoIdentityAccessManagementChannel ships one Call envelope and awaits its
// correlated Response. CallContext is Call that also aborts when ctx is done
// (transport.Conn satisfies both).
type HanzoIdentityAccessManagementChannel interface {
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

// HanzoIdentityAccessManagementClient is a typed RPC client for the
// HanzoIdentityAccessManagement service over a ZAP call channel. Each call takes
// a fresh PromiseID from sess; the pipelined "On" form of a method sets Target
// to a prior call's Promise so the server chains them.
type HanzoIdentityAccessManagementClient struct {
	ch   HanzoIdentityAccessManagementChannel
	cap  []byte
}

// NewHanzoIdentityAccessManagementClient returns a client that issues calls over
// ch, attaching cap (which may be nil) to every request.
func NewHanzoIdentityAccessManagementClient(ch HanzoIdentityAccessManagementChannel, capability []byte) *HanzoIdentityAccessManagementClient {
	return &HanzoIdentityAccessManagementClient{ch: ch, cap: capability}
}

func (c *HanzoIdentityAccessManagementClient) GetConfiguration(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementGetConfigurationOrdinal, "GetConfiguration", rpc.NoTarget, req)
}

// GetConfigurationOn issues GetConfiguration as a dependent call pipelined on the
// answer of on: the server substitutes on's resolved result for this call's
// payload before dispatch, so it ships without waiting for on to round-trip.
func (c *HanzoIdentityAccessManagementClient) GetConfigurationOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementGetConfigurationOrdinal, "GetConfiguration", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) PutConfiguration(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementPutConfigurationOrdinal, "PutConfiguration", rpc.NoTarget, req)
}

// PutConfigurationOn issues PutConfiguration as a dependent call pipelined on the
// answer of on.
func (c *HanzoIdentityAccessManagementClient) PutConfigurationOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementPutConfigurationOrdinal, "PutConfiguration", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) CreateUser(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementCreateUserOrdinal, "CreateUser", rpc.NoTarget, req)
}

// CreateUserOn issues CreateUser as a dependent call pipelined on the answer of on.
func (c *HanzoIdentityAccessManagementClient) CreateUserOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementCreateUserOrdinal, "CreateUser", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) GetUser(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementGetUserOrdinal, "GetUser", rpc.NoTarget, req)
}

// GetUserOn issues GetUser as a dependent call pipelined on the answer of on.
func (c *HanzoIdentityAccessManagementClient) GetUserOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementGetUserOrdinal, "GetUser", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) UpdateUser(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementUpdateUserOrdinal, "UpdateUser", rpc.NoTarget, req)
}

// UpdateUserOn issues UpdateUser as a dependent call pipelined on the answer of on.
func (c *HanzoIdentityAccessManagementClient) UpdateUserOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementUpdateUserOrdinal, "UpdateUser", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) DeleteUser(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementDeleteUserOrdinal, "DeleteUser", rpc.NoTarget, req)
}

// DeleteUserOn issues DeleteUser as a dependent call pipelined on the answer of on.
func (c *HanzoIdentityAccessManagementClient) DeleteUserOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementDeleteUserOrdinal, "DeleteUser", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) ListUsers(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementListUsersOrdinal, "ListUsers", rpc.NoTarget, req)
}

// ListUsersOn issues ListUsers as a dependent call pipelined on the answer of on.
func (c *HanzoIdentityAccessManagementClient) ListUsersOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementListUsersOrdinal, "ListUsers", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) CreateAccessKey(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementCreateAccessKeyOrdinal, "CreateAccessKey", rpc.NoTarget, req)
}

// CreateAccessKeyOn issues CreateAccessKey as a dependent call pipelined on the
// answer of on.
func (c *HanzoIdentityAccessManagementClient) CreateAccessKeyOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementCreateAccessKeyOrdinal, "CreateAccessKey", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) DeleteAccessKey(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementDeleteAccessKeyOrdinal, "DeleteAccessKey", rpc.NoTarget, req)
}

// DeleteAccessKeyOn issues DeleteAccessKey as a dependent call pipelined on the
// answer of on.
func (c *HanzoIdentityAccessManagementClient) DeleteAccessKeyOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementDeleteAccessKeyOrdinal, "DeleteAccessKey", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) GetUserByAccessKey(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementGetUserByAccessKeyOrdinal, "GetUserByAccessKey", rpc.NoTarget, req)
}

// GetUserByAccessKeyOn issues GetUserByAccessKey as a dependent call pipelined on
// the answer of on.
func (c *HanzoIdentityAccessManagementClient) GetUserByAccessKeyOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementGetUserByAccessKeyOrdinal, "GetUserByAccessKey", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) PutPolicy(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementPutPolicyOrdinal, "PutPolicy", rpc.NoTarget, req)
}

// PutPolicyOn issues PutPolicy as a dependent call pipelined on the answer of on.
func (c *HanzoIdentityAccessManagementClient) PutPolicyOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementPutPolicyOrdinal, "PutPolicy", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) GetPolicy(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementGetPolicyOrdinal, "GetPolicy", rpc.NoTarget, req)
}

// GetPolicyOn issues GetPolicy as a dependent call pipelined on the answer of on.
func (c *HanzoIdentityAccessManagementClient) GetPolicyOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementGetPolicyOrdinal, "GetPolicy", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) ListPolicies(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementListPoliciesOrdinal, "ListPolicies", rpc.NoTarget, req)
}

// ListPoliciesOn issues ListPolicies as a dependent call pipelined on the answer
// of on.
func (c *HanzoIdentityAccessManagementClient) ListPoliciesOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementListPoliciesOrdinal, "ListPolicies", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) DeletePolicy(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementDeletePolicyOrdinal, "DeletePolicy", rpc.NoTarget, req)
}

// DeletePolicyOn issues DeletePolicy as a dependent call pipelined on the answer
// of on.
func (c *HanzoIdentityAccessManagementClient) DeletePolicyOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementDeletePolicyOrdinal, "DeletePolicy", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) CreateServiceAccount(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementCreateServiceAccountOrdinal, "CreateServiceAccount", rpc.NoTarget, req)
}

// CreateServiceAccountOn issues CreateServiceAccount as a dependent call
// pipelined on the answer of on.
func (c *HanzoIdentityAccessManagementClient) CreateServiceAccountOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementCreateServiceAccountOrdinal, "CreateServiceAccount", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) UpdateServiceAccount(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementUpdateServiceAccountOrdinal, "UpdateServiceAccount", rpc.NoTarget, req)
}

// UpdateServiceAccountOn issues UpdateServiceAccount as a dependent call
// pipelined on the answer of on.
func (c *HanzoIdentityAccessManagementClient) UpdateServiceAccountOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementUpdateServiceAccountOrdinal, "UpdateServiceAccount", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) DeleteServiceAccount(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementDeleteServiceAccountOrdinal, "DeleteServiceAccount", rpc.NoTarget, req)
}

// DeleteServiceAccountOn issues DeleteServiceAccount as a dependent call
// pipelined on the answer of on.
func (c *HanzoIdentityAccessManagementClient) DeleteServiceAccountOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementDeleteServiceAccountOrdinal, "DeleteServiceAccount", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) GetServiceAccount(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementGetServiceAccountOrdinal, "GetServiceAccount", rpc.NoTarget, req)
}

// GetServiceAccountOn issues GetServiceAccount as a dependent call pipelined on
// the answer of on.
func (c *HanzoIdentityAccessManagementClient) GetServiceAccountOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementGetServiceAccountOrdinal, "GetServiceAccount", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) ListServiceAccounts(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementListServiceAccountsOrdinal, "ListServiceAccounts", rpc.NoTarget, req)
}

// ListServiceAccountsOn issues ListServiceAccounts as a dependent call pipelined
// on the answer of on.
func (c *HanzoIdentityAccessManagementClient) ListServiceAccountsOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementListServiceAccountsOrdinal, "ListServiceAccounts", on.ID, nil)
}

func (c *HanzoIdentityAccessManagementClient) GetServiceAccountByAccessKey(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementGetServiceAccountByAccessKeyOrdinal, "GetServiceAccountByAccessKey", rpc.NoTarget, req)
}

// GetServiceAccountByAccessKeyOn issues GetServiceAccountByAccessKey as a
// dependent call pipelined on the answer of on.
func (c *HanzoIdentityAccessManagementClient) GetServiceAccountByAccessKeyOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(ctx, HanzoIdentityAccessManagementGetServiceAccountByAccessKeyOrdinal, "GetServiceAccountByAccessKey", on.ID, nil)
}

// invoke ships one call over the channel and unwraps its response. method is the
// human-readable name used only for the status-error message.
func (c *HanzoIdentityAccessManagementClient) invoke(ctx context.Context, ordinal uint32, method string, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
		Method:    ordinal,
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
			return p, nil, fmt.Errorf("HanzoIdentityAccessManagement.%s: %s", method, resp.Body)
		}
		return p, nil, fmt.Errorf("HanzoIdentityAccessManagement.%s: status %d", method, resp.Status)
	}
	return p, resp.Body, nil
}

// HanzoIdentityAccessManagementHandler is the server contract for the
// HanzoIdentityAccessManagement service. Implement each method, then route
// requests to it with DispatchHanzoIdentityAccessManagement.
type HanzoIdentityAccessManagementHandler interface {
	GetConfiguration(req []byte) ([]byte, error)
	PutConfiguration(req []byte) ([]byte, error)
	CreateUser(req []byte) ([]byte, error)
	GetUser(req []byte) ([]byte, error)
	UpdateUser(req []byte) ([]byte, error)
	DeleteUser(req []byte) ([]byte, error)
	ListUsers(req []byte) ([]byte, error)
	CreateAccessKey(req []byte) ([]byte, error)
	DeleteAccessKey(req []byte) ([]byte, error)
	GetUserByAccessKey(req []byte) ([]byte, error)
	PutPolicy(req []byte) ([]byte, error)
	GetPolicy(req []byte) ([]byte, error)
	ListPolicies(req []byte) ([]byte, error)
	DeletePolicy(req []byte) ([]byte, error)
	CreateServiceAccount(req []byte) ([]byte, error)
	UpdateServiceAccount(req []byte) ([]byte, error)
	DeleteServiceAccount(req []byte) ([]byte, error)
	GetServiceAccount(req []byte) ([]byte, error)
	ListServiceAccounts(req []byte) ([]byte, error)
	GetServiceAccountByAccessKey(req []byte) ([]byte, error)
}

// DispatchHanzoIdentityAccessManagement decodes a Call envelope, routes it by
// method ordinal to h, and returns the response envelope. An unknown ordinal
// yields a StatusNotFound response; a handler error yields StatusInternal.
func DispatchHanzoIdentityAccessManagement(h HanzoIdentityAccessManagementHandler, envelope []byte) ([]byte, error) {
	call, err := rpc.ParseRequest(envelope)
	if err != nil {
		return nil, err
	}
	var (
		body []byte
		herr error
	)
	switch call.Method {
	case HanzoIdentityAccessManagementGetConfigurationOrdinal:
		body, herr = h.GetConfiguration(call.Payload)
	case HanzoIdentityAccessManagementPutConfigurationOrdinal:
		body, herr = h.PutConfiguration(call.Payload)
	case HanzoIdentityAccessManagementCreateUserOrdinal:
		body, herr = h.CreateUser(call.Payload)
	case HanzoIdentityAccessManagementGetUserOrdinal:
		body, herr = h.GetUser(call.Payload)
	case HanzoIdentityAccessManagementUpdateUserOrdinal:
		body, herr = h.UpdateUser(call.Payload)
	case HanzoIdentityAccessManagementDeleteUserOrdinal:
		body, herr = h.DeleteUser(call.Payload)
	case HanzoIdentityAccessManagementListUsersOrdinal:
		body, herr = h.ListUsers(call.Payload)
	case HanzoIdentityAccessManagementCreateAccessKeyOrdinal:
		body, herr = h.CreateAccessKey(call.Payload)
	case HanzoIdentityAccessManagementDeleteAccessKeyOrdinal:
		body, herr = h.DeleteAccessKey(call.Payload)
	case HanzoIdentityAccessManagementGetUserByAccessKeyOrdinal:
		body, herr = h.GetUserByAccessKey(call.Payload)
	case HanzoIdentityAccessManagementPutPolicyOrdinal:
		body, herr = h.PutPolicy(call.Payload)
	case HanzoIdentityAccessManagementGetPolicyOrdinal:
		body, herr = h.GetPolicy(call.Payload)
	case HanzoIdentityAccessManagementListPoliciesOrdinal:
		body, herr = h.ListPolicies(call.Payload)
	case HanzoIdentityAccessManagementDeletePolicyOrdinal:
		body, herr = h.DeletePolicy(call.Payload)
	case HanzoIdentityAccessManagementCreateServiceAccountOrdinal:
		body, herr = h.CreateServiceAccount(call.Payload)
	case HanzoIdentityAccessManagementUpdateServiceAccountOrdinal:
		body, herr = h.UpdateServiceAccount(call.Payload)
	case HanzoIdentityAccessManagementDeleteServiceAccountOrdinal:
		body, herr = h.DeleteServiceAccount(call.Payload)
	case HanzoIdentityAccessManagementGetServiceAccountOrdinal:
		body, herr = h.GetServiceAccount(call.Payload)
	case HanzoIdentityAccessManagementListServiceAccountsOrdinal:
		body, herr = h.ListServiceAccounts(call.Payload)
	case HanzoIdentityAccessManagementGetServiceAccountByAccessKeyOrdinal:
		body, herr = h.GetServiceAccountByAccessKey(call.Payload)
	default:
		return rpc.BuildResponse(rpc.StatusNotFound, call.PromiseID, nil), nil
	}
	if herr != nil {
		// Carry the handler error message so the caller can reconstruct sentinels.
		return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(herr.Error())), nil
	}
	return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
}
