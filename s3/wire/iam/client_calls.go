// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// Typed Client call sites for the HanzoIdentityAccessManagement ZAP service.
// One method per RPC, in iam.proto order. Each builds the request with the
// matching New* iamwire builder, ships it over the generated client, and
// decodes the reply into the in-process iam_pb model using the bridge_pb.go
// decoders — so internal call sites get exactly the shapes the retired gRPC
// stub returned, with the bytes never leaving the zero-copy ZAP form on the
// wire. This is the IAM analogue of object Client.GetObject/PutObject.

package iamwire

import (
	"context"

	"github.com/hanzoai/s3/s3/pb/iam_pb"
)

// ---- Configuration management ----

// GetConfiguration fetches the full S3 API configuration over ZAP.
func (c *Client) GetConfiguration(ctx context.Context) (*iam_pb.S3ApiConfiguration, error) {
	_, body, err := c.rpc.GetConfiguration(ctx, NewGetConfigurationRequest(GetConfigurationRequestInput{}))
	if err != nil {
		return nil, err
	}
	return GetConfigurationResp(body)
}

// PutConfiguration replaces the full S3 API configuration over ZAP.
func (c *Client) PutConfiguration(ctx context.Context, cfg *iam_pb.S3ApiConfiguration) error {
	_, _, err := c.rpc.PutConfiguration(ctx, NewPutConfigurationRequest(PutConfigurationRequestInput{
		Configuration: ConfigurationInputFromPB(cfg),
	}))
	return err
}

// ---- User management ----

// CreateUser creates a user from an identity.
func (c *Client) CreateUser(ctx context.Context, identity *iam_pb.Identity) error {
	_, _, err := c.rpc.CreateUser(ctx, NewCreateUserRequest(CreateUserRequestInput{
		Identity: IdentityInputFromPB(identity),
	}))
	return err
}

// GetUser fetches a user by username; a nil identity means no such user.
func (c *Client) GetUser(ctx context.Context, username string) (*iam_pb.Identity, error) {
	_, body, err := c.rpc.GetUser(ctx, NewGetUserRequest(GetUserRequestInput{Username: username}))
	if err != nil {
		return nil, err
	}
	return GetUserResp(body)
}

// UpdateUser replaces the identity stored for username.
func (c *Client) UpdateUser(ctx context.Context, username string, identity *iam_pb.Identity) error {
	_, _, err := c.rpc.UpdateUser(ctx, NewUpdateUserRequest(UpdateUserRequestInput{
		Username: username,
		Identity: IdentityInputFromPB(identity),
	}))
	return err
}

// DeleteUser removes a user by username.
func (c *Client) DeleteUser(ctx context.Context, username string) error {
	_, _, err := c.rpc.DeleteUser(ctx, NewDeleteUserRequest(DeleteUserRequestInput{Username: username}))
	return err
}

// ListUsers returns every username known to the IAM backend.
func (c *Client) ListUsers(ctx context.Context) ([]string, error) {
	_, body, err := c.rpc.ListUsers(ctx, NewListUsersRequest(ListUsersRequestInput{}))
	if err != nil {
		return nil, err
	}
	resp, err := WrapListUsersResponse(body)
	if err != nil {
		return nil, err
	}
	n := resp.UsernamesLen()
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = resp.UsernamesAt(i)
	}
	return out, nil
}

// ---- Access-key management ----

// CreateAccessKey attaches a credential to username.
func (c *Client) CreateAccessKey(ctx context.Context, username string, cred *iam_pb.Credential) error {
	_, _, err := c.rpc.CreateAccessKey(ctx, NewCreateAccessKeyRequest(CreateAccessKeyRequestInput{
		Username:   username,
		Credential: CredentialInputFromPB(cred),
	}))
	return err
}

// DeleteAccessKey removes accessKey from username.
func (c *Client) DeleteAccessKey(ctx context.Context, username, accessKey string) error {
	_, _, err := c.rpc.DeleteAccessKey(ctx, NewDeleteAccessKeyRequest(DeleteAccessKeyRequestInput{
		Username:  username,
		AccessKey: accessKey,
	}))
	return err
}

// GetUserByAccessKey resolves an access key to its owning identity; a nil
// identity means the key is unknown.
func (c *Client) GetUserByAccessKey(ctx context.Context, accessKey string) (*iam_pb.Identity, error) {
	_, body, err := c.rpc.GetUserByAccessKey(ctx, NewGetUserByAccessKeyRequest(GetUserByAccessKeyRequestInput{
		AccessKey: accessKey,
	}))
	if err != nil {
		return nil, err
	}
	return GetUserByAccessKeyResp(body)
}

// ---- Policy management ----

// PutPolicy stores a named managed policy (content is the JSON policy document).
func (c *Client) PutPolicy(ctx context.Context, name, content string) error {
	_, _, err := c.rpc.PutPolicy(ctx, NewPutPolicyRequest(PutPolicyRequestInput{Name: name, Content: content}))
	return err
}

// GetPolicy fetches a named managed policy, returning its name and JSON content.
func (c *Client) GetPolicy(ctx context.Context, name string) (gotName, content string, err error) {
	_, body, err := c.rpc.GetPolicy(ctx, NewGetPolicyRequest(GetPolicyRequestInput{Name: name}))
	if err != nil {
		return "", "", err
	}
	return GetPolicyResp(body)
}

// ListPolicies returns every managed policy.
func (c *Client) ListPolicies(ctx context.Context) ([]*iam_pb.Policy, error) {
	_, body, err := c.rpc.ListPolicies(ctx, NewListPoliciesRequest(ListPoliciesRequestInput{}))
	if err != nil {
		return nil, err
	}
	return ListPoliciesResp(body)
}

// DeletePolicy removes a named managed policy.
func (c *Client) DeletePolicy(ctx context.Context, name string) error {
	_, _, err := c.rpc.DeletePolicy(ctx, NewDeletePolicyRequest(DeletePolicyRequestInput{Name: name}))
	return err
}

// ---- Service-account management ----

// CreateServiceAccount creates a service account.
func (c *Client) CreateServiceAccount(ctx context.Context, sa *iam_pb.ServiceAccount) error {
	_, _, err := c.rpc.CreateServiceAccount(ctx, NewCreateServiceAccountRequest(CreateServiceAccountRequestInput{
		ServiceAccount: ServiceAccountInputFromPB(sa),
	}))
	return err
}

// UpdateServiceAccount replaces the service account stored under id.
func (c *Client) UpdateServiceAccount(ctx context.Context, id string, sa *iam_pb.ServiceAccount) error {
	_, _, err := c.rpc.UpdateServiceAccount(ctx, NewUpdateServiceAccountRequest(UpdateServiceAccountRequestInput{
		ID:             id,
		ServiceAccount: ServiceAccountInputFromPB(sa),
	}))
	return err
}

// DeleteServiceAccount removes a service account by id.
func (c *Client) DeleteServiceAccount(ctx context.Context, id string) error {
	_, _, err := c.rpc.DeleteServiceAccount(ctx, NewDeleteServiceAccountRequest(DeleteServiceAccountRequestInput{ID: id}))
	return err
}

// GetServiceAccount fetches a service account by id; a nil result means no such
// service account.
func (c *Client) GetServiceAccount(ctx context.Context, id string) (*iam_pb.ServiceAccount, error) {
	_, body, err := c.rpc.GetServiceAccount(ctx, NewGetServiceAccountRequest(GetServiceAccountRequestInput{ID: id}))
	if err != nil {
		return nil, err
	}
	return GetServiceAccountResp(body)
}

// ListServiceAccounts returns every service account.
func (c *Client) ListServiceAccounts(ctx context.Context) ([]*iam_pb.ServiceAccount, error) {
	_, body, err := c.rpc.ListServiceAccounts(ctx, NewListServiceAccountsRequest(ListServiceAccountsRequestInput{}))
	if err != nil {
		return nil, err
	}
	return ListServiceAccountsResp(body)
}

// GetServiceAccountByAccessKey resolves an access key to its service account; a
// nil result means the key is unknown.
func (c *Client) GetServiceAccountByAccessKey(ctx context.Context, accessKey string) (*iam_pb.ServiceAccount, error) {
	_, body, err := c.rpc.GetServiceAccountByAccessKey(ctx, NewGetServiceAccountByAccessKeyRequest(
		GetServiceAccountByAccessKeyRequestInput{AccessKey: accessKey}))
	if err != nil {
		return nil, err
	}
	return GetServiceAccountByAccessKeyResp(body)
}
