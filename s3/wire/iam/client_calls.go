// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// Typed Client call sites for the HanzoIdentityAccessManagement ZAP service.
// One method per RPC, in iam.proto order. Each builds the request with the
// matching New* iamwire builder, ships it over the generated client, and
// decodes the reply into the in-process iam_pb model using the bridge_pb.go
// decoders — so internal call sites get exactly the shapes the retired gRPC
// stub returned, with the bytes never leaving the zero-copy ZAP form on the
// wire. This is the IAM analogue of zapsvc Client.GetObject/PutObject.

package iamwire

import (
	"github.com/hanzoai/s3/s3/pb/iam_pb"
)

// ---- Configuration management ----

// GetConfiguration fetches the full S3 API configuration over ZAP.
func (c *Client) GetConfiguration() (*iam_pb.S3ApiConfiguration, error) {
	_, body, err := c.rpc.GetConfiguration(NewGetConfigurationRequest(GetConfigurationRequestInput{}))
	if err != nil {
		return nil, err
	}
	return GetConfigurationResp(body)
}

// PutConfiguration replaces the full S3 API configuration over ZAP.
func (c *Client) PutConfiguration(cfg *iam_pb.S3ApiConfiguration) error {
	_, _, err := c.rpc.PutConfiguration(NewPutConfigurationRequest(PutConfigurationRequestInput{
		Configuration: ConfigurationInputFromPB(cfg),
	}))
	return err
}

// ---- User management ----

// CreateUser creates a user from an identity.
func (c *Client) CreateUser(identity *iam_pb.Identity) error {
	_, _, err := c.rpc.CreateUser(NewCreateUserRequest(CreateUserRequestInput{
		Identity: IdentityInputFromPB(identity),
	}))
	return err
}

// GetUser fetches a user by username; a nil identity means no such user.
func (c *Client) GetUser(username string) (*iam_pb.Identity, error) {
	_, body, err := c.rpc.GetUser(NewGetUserRequest(GetUserRequestInput{Username: username}))
	if err != nil {
		return nil, err
	}
	return GetUserResp(body)
}

// UpdateUser replaces the identity stored for username.
func (c *Client) UpdateUser(username string, identity *iam_pb.Identity) error {
	_, _, err := c.rpc.UpdateUser(NewUpdateUserRequest(UpdateUserRequestInput{
		Username: username,
		Identity: IdentityInputFromPB(identity),
	}))
	return err
}

// DeleteUser removes a user by username.
func (c *Client) DeleteUser(username string) error {
	_, _, err := c.rpc.DeleteUser(NewDeleteUserRequest(DeleteUserRequestInput{Username: username}))
	return err
}

// ListUsers returns every username known to the IAM backend.
func (c *Client) ListUsers() ([]string, error) {
	_, body, err := c.rpc.ListUsers(NewListUsersRequest(ListUsersRequestInput{}))
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
func (c *Client) CreateAccessKey(username string, cred *iam_pb.Credential) error {
	_, _, err := c.rpc.CreateAccessKey(NewCreateAccessKeyRequest(CreateAccessKeyRequestInput{
		Username:   username,
		Credential: CredentialInputFromPB(cred),
	}))
	return err
}

// DeleteAccessKey removes accessKey from username.
func (c *Client) DeleteAccessKey(username, accessKey string) error {
	_, _, err := c.rpc.DeleteAccessKey(NewDeleteAccessKeyRequest(DeleteAccessKeyRequestInput{
		Username:  username,
		AccessKey: accessKey,
	}))
	return err
}

// GetUserByAccessKey resolves an access key to its owning identity; a nil
// identity means the key is unknown.
func (c *Client) GetUserByAccessKey(accessKey string) (*iam_pb.Identity, error) {
	_, body, err := c.rpc.GetUserByAccessKey(NewGetUserByAccessKeyRequest(GetUserByAccessKeyRequestInput{
		AccessKey: accessKey,
	}))
	if err != nil {
		return nil, err
	}
	return GetUserByAccessKeyResp(body)
}

// ---- Policy management ----

// PutPolicy stores a named managed policy (content is the JSON policy document).
func (c *Client) PutPolicy(name, content string) error {
	_, _, err := c.rpc.PutPolicy(NewPutPolicyRequest(PutPolicyRequestInput{Name: name, Content: content}))
	return err
}

// GetPolicy fetches a named managed policy, returning its name and JSON content.
func (c *Client) GetPolicy(name string) (gotName, content string, err error) {
	_, body, err := c.rpc.GetPolicy(NewGetPolicyRequest(GetPolicyRequestInput{Name: name}))
	if err != nil {
		return "", "", err
	}
	return GetPolicyResp(body)
}

// ListPolicies returns every managed policy.
func (c *Client) ListPolicies() ([]*iam_pb.Policy, error) {
	_, body, err := c.rpc.ListPolicies(NewListPoliciesRequest(ListPoliciesRequestInput{}))
	if err != nil {
		return nil, err
	}
	return ListPoliciesResp(body)
}

// DeletePolicy removes a named managed policy.
func (c *Client) DeletePolicy(name string) error {
	_, _, err := c.rpc.DeletePolicy(NewDeletePolicyRequest(DeletePolicyRequestInput{Name: name}))
	return err
}

// ---- Service-account management ----

// CreateServiceAccount creates a service account.
func (c *Client) CreateServiceAccount(sa *iam_pb.ServiceAccount) error {
	_, _, err := c.rpc.CreateServiceAccount(NewCreateServiceAccountRequest(CreateServiceAccountRequestInput{
		ServiceAccount: ServiceAccountInputFromPB(sa),
	}))
	return err
}

// UpdateServiceAccount replaces the service account stored under id.
func (c *Client) UpdateServiceAccount(id string, sa *iam_pb.ServiceAccount) error {
	_, _, err := c.rpc.UpdateServiceAccount(NewUpdateServiceAccountRequest(UpdateServiceAccountRequestInput{
		ID:             id,
		ServiceAccount: ServiceAccountInputFromPB(sa),
	}))
	return err
}

// DeleteServiceAccount removes a service account by id.
func (c *Client) DeleteServiceAccount(id string) error {
	_, _, err := c.rpc.DeleteServiceAccount(NewDeleteServiceAccountRequest(DeleteServiceAccountRequestInput{ID: id}))
	return err
}

// GetServiceAccount fetches a service account by id; a nil result means no such
// service account.
func (c *Client) GetServiceAccount(id string) (*iam_pb.ServiceAccount, error) {
	_, body, err := c.rpc.GetServiceAccount(NewGetServiceAccountRequest(GetServiceAccountRequestInput{ID: id}))
	if err != nil {
		return nil, err
	}
	return GetServiceAccountResp(body)
}

// ListServiceAccounts returns every service account.
func (c *Client) ListServiceAccounts() ([]*iam_pb.ServiceAccount, error) {
	_, body, err := c.rpc.ListServiceAccounts(NewListServiceAccountsRequest(ListServiceAccountsRequestInput{}))
	if err != nil {
		return nil, err
	}
	return ListServiceAccountsResp(body)
}

// GetServiceAccountByAccessKey resolves an access key to its service account; a
// nil result means the key is unknown.
func (c *Client) GetServiceAccountByAccessKey(accessKey string) (*iam_pb.ServiceAccount, error) {
	_, body, err := c.rpc.GetServiceAccountByAccessKey(NewGetServiceAccountByAccessKeyRequest(
		GetServiceAccountByAccessKeyRequestInput{AccessKey: accessKey}))
	if err != nil {
		return nil, err
	}
	return GetServiceAccountByAccessKeyResp(body)
}
