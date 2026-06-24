package grpc

import (
	"context"

	"github.com/hanzoai/s3/s3/credential"
	"github.com/hanzoai/s3/s3/pb/iam_pb"
	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

func (store *IamGrpcStore) LoadConfiguration(ctx context.Context) (*iam_pb.S3ApiConfiguration, error) {
	var config *iam_pb.S3ApiConfiguration
	err := store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, body, err := client.GetConfiguration(iamwire.NewGetConfigurationRequest(iamwire.GetConfigurationRequestInput{}))
		if err != nil {
			return err
		}
		config, err = iamwire.GetConfigurationResp(body)
		return err
	})
	return config, err
}

func (store *IamGrpcStore) SaveConfiguration(ctx context.Context, config *iam_pb.S3ApiConfiguration) error {
	return store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, _, err := client.PutConfiguration(iamwire.NewPutConfigurationRequest(iamwire.PutConfigurationRequestInput{
			Configuration: iamwire.ConfigurationInputFromPB(config),
		}))
		return err
	})
}

func (store *IamGrpcStore) CreateUser(ctx context.Context, identity *iam_pb.Identity) error {
	return store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, _, err := client.CreateUser(iamwire.NewCreateUserRequest(iamwire.CreateUserRequestInput{
			Identity: iamwire.IdentityInputFromPB(identity),
		}))
		return err
	})
}

func (store *IamGrpcStore) GetUser(ctx context.Context, username string) (*iam_pb.Identity, error) {
	var identity *iam_pb.Identity
	err := store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, body, err := client.GetUser(iamwire.NewGetUserRequest(iamwire.GetUserRequestInput{
			Username: username,
		}))
		if err != nil {
			return err
		}
		identity, err = iamwire.GetUserResp(body)
		return err
	})
	// The filer-side handler errors when the user is absent; over the ZAP
	// transport that surfaces as a generic call error (no gRPC status codes), so
	// translate any call failure to the package sentinel so callers can use
	// errors.Is(err, credential.ErrUserNotFound) uniformly across stores.
	if err != nil {
		return nil, credential.ErrUserNotFound
	}
	return identity, nil
}

func (store *IamGrpcStore) UpdateUser(ctx context.Context, username string, identity *iam_pb.Identity) error {
	return store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, _, err := client.UpdateUser(iamwire.NewUpdateUserRequest(iamwire.UpdateUserRequestInput{
			Username: username,
			Identity: iamwire.IdentityInputFromPB(identity),
		}))
		return err
	})
}

func (store *IamGrpcStore) DeleteUser(ctx context.Context, username string) error {
	return store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, _, err := client.DeleteUser(iamwire.NewDeleteUserRequest(iamwire.DeleteUserRequestInput{
			Username: username,
		}))
		return err
	})
}

func (store *IamGrpcStore) ListUsers(ctx context.Context) ([]string, error) {
	var usernames []string
	err := store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, body, err := client.ListUsers(iamwire.NewListUsersRequest(iamwire.ListUsersRequestInput{}))
		if err != nil {
			return err
		}
		resp, err := iamwire.WrapListUsersResponse(body)
		if err != nil {
			return err
		}
		n := resp.UsernamesLen()
		usernames = make([]string, n)
		for i := 0; i < n; i++ {
			usernames[i] = resp.UsernamesAt(i)
		}
		return nil
	})
	return usernames, err
}

func (store *IamGrpcStore) GetUserByAccessKey(ctx context.Context, accessKey string) (*iam_pb.Identity, error) {
	var identity *iam_pb.Identity
	err := store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, body, err := client.GetUserByAccessKey(iamwire.NewGetUserByAccessKeyRequest(iamwire.GetUserByAccessKeyRequestInput{
			AccessKey: accessKey,
		}))
		if err != nil {
			return err
		}
		identity, err = iamwire.GetUserByAccessKeyResp(body)
		return err
	})
	return identity, err
}

func (store *IamGrpcStore) CreateAccessKey(ctx context.Context, username string, cred *iam_pb.Credential) error {
	return store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, _, err := client.CreateAccessKey(iamwire.NewCreateAccessKeyRequest(iamwire.CreateAccessKeyRequestInput{
			Username:   username,
			Credential: iamwire.CredentialInputFromPB(cred),
		}))
		return err
	})
}

func (store *IamGrpcStore) DeleteAccessKey(ctx context.Context, username string, accessKey string) error {
	return store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, _, err := client.DeleteAccessKey(iamwire.NewDeleteAccessKeyRequest(iamwire.DeleteAccessKeyRequestInput{
			Username:  username,
			AccessKey: accessKey,
		}))
		return err
	})
}

// AttachUserPolicy attaches a managed policy to a user by policy name
func (store *IamGrpcStore) AttachUserPolicy(ctx context.Context, username string, policyName string) error {
	// Get current user
	identity, err := store.GetUser(ctx, username)
	if err != nil {
		return err
	}

	// Verify policy exists
	policy, err := store.GetPolicy(ctx, policyName)
	if err != nil {
		return err
	}
	if policy == nil {
		return credential.ErrPolicyNotFound
	}

	// Check if already attached
	for _, p := range identity.PolicyNames {
		if p == policyName {
			// Already attached - return success (idempotent)
			return nil
		}
	}

	identity.PolicyNames = append(identity.PolicyNames, policyName)
	return store.UpdateUser(ctx, username, identity)
}

// DetachUserPolicy detaches a managed policy from a user
func (store *IamGrpcStore) DetachUserPolicy(ctx context.Context, username string, policyName string) error {
	identity, err := store.GetUser(ctx, username)
	if err != nil {
		return err
	}

	found := false
	var newPolicies []string
	for _, p := range identity.PolicyNames {
		if p == policyName {
			found = true
		} else {
			newPolicies = append(newPolicies, p)
		}
	}

	if !found {
		return credential.ErrPolicyNotAttached
	}

	identity.PolicyNames = newPolicies
	return store.UpdateUser(ctx, username, identity)
}

// ListAttachedUserPolicies returns the list of policy names attached to a user
func (store *IamGrpcStore) ListAttachedUserPolicies(ctx context.Context, username string) ([]string, error) {
	identity, err := store.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	return identity.PolicyNames, nil
}
