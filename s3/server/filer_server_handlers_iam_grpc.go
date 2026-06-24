package s3server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hanzoai/s3/s3/credential"
	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/s3api/policy_engine"
	"github.com/hanzoai/s3/s3/security"
	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

// IamGrpcServer implements the IAM service on the filer over the ZAP transport.
// It satisfies iamwire.HanzoIdentityAccessManagementHandler: every method takes
// a ZAP-encoded request envelope and returns a ZAP-encoded response.
//
// The legacy gRPC metadata Bearer-token check is gone: the IAM service now
// rides a dedicated ZAP listener whose transport layer (X-Wing PQ TLS) is the
// connection-level security boundary, so adminSigningKey is retained only for
// continuity of construction and is no longer consulted per-call.
type IamGrpcServer struct {
	credentialManager *credential.CredentialManager
	adminSigningKey   security.SigningKey
}

// Compile-time assertion that the server satisfies the wire handler contract.
var _ iamwire.HanzoIdentityAccessManagementHandler = (*IamGrpcServer)(nil)

// NewIamGrpcServer creates a new IAM server. adminSigningKey is retained for
// construction compatibility; connection security is provided by the ZAP
// transport.
func NewIamGrpcServer(credentialManager *credential.CredentialManager, adminSigningKey security.SigningKey) *IamGrpcServer {
	return &IamGrpcServer{
		credentialManager: credentialManager,
		adminSigningKey:   adminSigningKey,
	}
}

//////////////////////////////////////////////////
// Configuration Management

func (s *IamGrpcServer) GetConfiguration(req []byte) ([]byte, error) {
	if _, err := iamwire.WrapGetConfigurationRequest(req); err != nil {
		return nil, err
	}
	glog.V(4).Infof("GetConfiguration")

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	config, err := s.credentialManager.LoadConfiguration(context.Background())
	if err != nil {
		glog.Errorf("Failed to load configuration: %v", err)
		return nil, err
	}

	return iamwire.NewGetConfigurationResponse(iamwire.GetConfigurationResponseInput{
		Configuration: iamwire.ConfigurationInputFromPB(config),
	}), nil
}

func (s *IamGrpcServer) PutConfiguration(req []byte) ([]byte, error) {
	request, err := iamwire.WrapPutConfigurationRequest(req)
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("PutConfiguration")

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	config := iamwire.ConfigurationToPB(request.Configuration())
	if err := s.credentialManager.SaveConfiguration(context.Background(), config); err != nil {
		glog.Errorf("Failed to save configuration: %v", err)
		return nil, err
	}

	return iamwire.NewPutConfigurationResponse(iamwire.PutConfigurationResponseInput{}), nil
}

//////////////////////////////////////////////////
// User Management

func (s *IamGrpcServer) CreateUser(req []byte) ([]byte, error) {
	request, err := iamwire.WrapCreateUserRequest(req)
	if err != nil {
		return nil, err
	}
	identity := iamwire.IdentityToPB(request.Identity())
	glog.V(4).Infof("IAM: Filer.CreateUser %s", identity.Name)

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	if err := s.credentialManager.CreateUser(context.Background(), identity); err != nil {
		if err == credential.ErrUserAlreadyExists {
			return nil, fmt.Errorf("user %s already exists", identity.Name)
		}
		glog.Errorf("Failed to create user %s: %v", identity.Name, err)
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return iamwire.NewCreateUserResponse(iamwire.CreateUserResponseInput{}), nil
}

func (s *IamGrpcServer) GetUser(req []byte) ([]byte, error) {
	request, err := iamwire.WrapGetUserRequest(req)
	if err != nil {
		return nil, err
	}
	username := request.Username()
	glog.V(4).Infof("GetUser: %s", username)

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	identity, err := s.credentialManager.GetUser(context.Background(), username)
	if err != nil {
		if err == credential.ErrUserNotFound {
			// Fall back to static identities (loaded from -s3.config file)
			if si := s.credentialManager.GetStaticIdentity(username); si != nil {
				return iamwire.NewGetUserResponse(iamwire.GetUserResponseInput{
					Identity: iamwire.IdentityInputFromPB(si),
				}), nil
			}
			return nil, fmt.Errorf("user %s not found", username)
		}
		glog.Errorf("Failed to get user %s: %v", username, err)
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return iamwire.NewGetUserResponse(iamwire.GetUserResponseInput{
		Identity: iamwire.IdentityInputFromPB(identity),
	}), nil
}

func (s *IamGrpcServer) UpdateUser(req []byte) ([]byte, error) {
	request, err := iamwire.WrapUpdateUserRequest(req)
	if err != nil {
		return nil, err
	}
	username := request.Username()
	identity := iamwire.IdentityToPB(request.Identity())
	glog.V(4).Infof("IAM: Filer.UpdateUser %s", username)

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	if err := s.credentialManager.UpdateUser(context.Background(), username, identity); err != nil {
		if err == credential.ErrUserNotFound {
			return nil, fmt.Errorf("user %s not found", username)
		}
		glog.Errorf("Failed to update user %s: %v", username, err)
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return iamwire.NewUpdateUserResponse(iamwire.UpdateUserResponseInput{}), nil
}

func (s *IamGrpcServer) DeleteUser(req []byte) ([]byte, error) {
	request, err := iamwire.WrapDeleteUserRequest(req)
	if err != nil {
		return nil, err
	}
	username := request.Username()
	glog.V(4).Infof("IAM: Filer.DeleteUser %s", username)

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	if err := s.credentialManager.DeleteUser(context.Background(), username); err != nil {
		if err == credential.ErrUserNotFound {
			return nil, fmt.Errorf("user %s not found", username)
		}
		glog.Errorf("Failed to delete user %s: %v", username, err)
		return nil, fmt.Errorf("failed to delete user: %v", err)
	}

	return iamwire.NewDeleteUserResponse(iamwire.DeleteUserResponseInput{}), nil
}

func (s *IamGrpcServer) ListUsers(req []byte) ([]byte, error) {
	if _, err := iamwire.WrapListUsersRequest(req); err != nil {
		return nil, err
	}
	glog.V(4).Infof("ListUsers")

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	usernames, err := s.credentialManager.ListUsers(context.Background())
	if err != nil {
		glog.Errorf("Failed to list users: %v", err)
		return nil, err
	}

	// Merge static identities (from -s3.config file) into the result
	staticNames := s.credentialManager.GetStaticUsernames()
	if len(staticNames) > 0 {
		dynamicSet := make(map[string]bool, len(usernames))
		for _, name := range usernames {
			dynamicSet[name] = true
		}
		for _, name := range staticNames {
			if !dynamicSet[name] {
				usernames = append(usernames, name)
			}
		}
	}

	return iamwire.NewListUsersResponse(iamwire.ListUsersResponseInput{
		Usernames: usernames,
	}), nil
}

//////////////////////////////////////////////////
// Access Key Management

func (s *IamGrpcServer) CreateAccessKey(req []byte) ([]byte, error) {
	request, err := iamwire.WrapCreateAccessKeyRequest(req)
	if err != nil {
		return nil, err
	}
	username := request.Username()
	cred := iamwire.CredentialToPBExported(request.Credential())
	glog.V(4).Infof("CreateAccessKey for user: %s", username)

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	if err := s.credentialManager.CreateAccessKey(context.Background(), username, cred); err != nil {
		if err == credential.ErrUserNotFound {
			return nil, fmt.Errorf("user %s not found", username)
		}
		glog.Errorf("Failed to create access key for user %s: %v", username, err)
		return nil, fmt.Errorf("failed to create access key: %v", err)
	}

	return iamwire.NewCreateAccessKeyResponse(iamwire.CreateAccessKeyResponseInput{}), nil
}

func (s *IamGrpcServer) DeleteAccessKey(req []byte) ([]byte, error) {
	request, err := iamwire.WrapDeleteAccessKeyRequest(req)
	if err != nil {
		return nil, err
	}
	username := request.Username()
	accessKey := request.AccessKey()
	glog.V(4).Infof("DeleteAccessKey: %s for user: %s", accessKey, username)

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	if err := s.credentialManager.DeleteAccessKey(context.Background(), username, accessKey); err != nil {
		if err == credential.ErrUserNotFound {
			return nil, fmt.Errorf("user %s not found", username)
		}
		if err == credential.ErrAccessKeyNotFound {
			return nil, fmt.Errorf("access key %s not found", accessKey)
		}
		glog.Errorf("Failed to delete access key %s for user %s: %v", accessKey, username, err)
		return nil, fmt.Errorf("failed to delete access key: %v", err)
	}

	return iamwire.NewDeleteAccessKeyResponse(iamwire.DeleteAccessKeyResponseInput{}), nil
}

func (s *IamGrpcServer) GetUserByAccessKey(req []byte) ([]byte, error) {
	request, err := iamwire.WrapGetUserByAccessKeyRequest(req)
	if err != nil {
		return nil, err
	}
	accessKey := request.AccessKey()
	glog.V(4).Infof("GetUserByAccessKey: %s", accessKey)

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	identity, err := s.credentialManager.GetUserByAccessKey(context.Background(), accessKey)
	if err != nil {
		if err == credential.ErrAccessKeyNotFound {
			return nil, fmt.Errorf("access key %s not found", accessKey)
		}
		glog.Errorf("Failed to get user by access key %s: %v", accessKey, err)
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return iamwire.NewGetUserByAccessKeyResponse(iamwire.GetUserByAccessKeyResponseInput{
		Identity: iamwire.IdentityInputFromPB(identity),
	}), nil
}

//////////////////////////////////////////////////
// Policy Management

func (s *IamGrpcServer) PutPolicy(req []byte) ([]byte, error) {
	request, err := iamwire.WrapPutPolicyRequest(req)
	if err != nil {
		return nil, err
	}
	name := request.Name()
	content := request.Content()
	glog.V(4).Infof("IAM: Filer.PutPolicy %s", name)

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	if name == "" {
		return nil, fmt.Errorf("policy name is required")
	}
	if err := credential.ValidatePolicyName(name); err != nil {
		return nil, err
	}
	if content == "" {
		return nil, fmt.Errorf("policy content is required")
	}

	var policy policy_engine.PolicyDocument
	if err := json.Unmarshal([]byte(content), &policy); err != nil {
		glog.Errorf("Failed to unmarshal policy %s: %v", name, err)
		return nil, err
	}

	if err := s.credentialManager.PutPolicy(context.Background(), name, policy); err != nil {
		glog.Errorf("Failed to put policy %s: %v", name, err)
		return nil, err
	}

	return iamwire.NewPutPolicyResponse(iamwire.PutPolicyResponseInput{}), nil
}

func (s *IamGrpcServer) GetPolicy(req []byte) ([]byte, error) {
	request, err := iamwire.WrapGetPolicyRequest(req)
	if err != nil {
		return nil, err
	}
	name := request.Name()
	glog.V(4).Infof("GetPolicy: %s", name)

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	policy, err := s.credentialManager.GetPolicy(context.Background(), name)
	if err != nil {
		glog.Errorf("Failed to get policy %s: %v", name, err)
		return nil, err
	}

	if policy == nil {
		return nil, fmt.Errorf("policy %s not found", name)
	}

	jsonBytes, err := json.Marshal(policy)
	if err != nil {
		glog.Errorf("Failed to marshal policy %s: %v", name, err)
		return nil, err
	}

	return iamwire.NewGetPolicyResponse(iamwire.GetPolicyResponseInput{
		Name:    name,
		Content: string(jsonBytes),
	}), nil
}

func (s *IamGrpcServer) ListPolicies(req []byte) ([]byte, error) {
	if _, err := iamwire.WrapListPoliciesRequest(req); err != nil {
		return nil, err
	}
	glog.V(4).Infof("ListPolicies")

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	policiesData, err := s.credentialManager.GetPolicies(context.Background())
	if err != nil {
		glog.Errorf("Failed to list policies: %v", err)
		return nil, err
	}

	var policies []iamwire.PolicyInput
	for name, policy := range policiesData {
		jsonBytes, err := json.Marshal(policy)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal policy %s: %v", name, err)
		}
		policies = append(policies, iamwire.PolicyInput{
			Name:    name,
			Content: string(jsonBytes),
		})
	}

	return iamwire.NewListPoliciesResponse(iamwire.ListPoliciesResponseInput{
		Policies: policies,
	}), nil
}

func (s *IamGrpcServer) DeletePolicy(req []byte) ([]byte, error) {
	request, err := iamwire.WrapDeletePolicyRequest(req)
	if err != nil {
		return nil, err
	}
	name := request.Name()
	glog.V(4).Infof("DeletePolicy: %s", name)

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	if err := s.credentialManager.DeletePolicy(context.Background(), name); err != nil {
		glog.Errorf("Failed to delete policy %s: %v", name, err)
		return nil, err
	}

	return iamwire.NewDeletePolicyResponse(iamwire.DeletePolicyResponseInput{}), nil
}

//////////////////////////////////////////////////
// Service Account Management

func (s *IamGrpcServer) CreateServiceAccount(req []byte) ([]byte, error) {
	request, err := iamwire.WrapCreateServiceAccountRequest(req)
	if err != nil {
		return nil, err
	}
	sa := iamwire.ServiceAccountToPB(request.ServiceAccount())
	if err := credential.ValidateServiceAccountId(sa.Id); err != nil {
		return nil, err
	}
	glog.V(4).Infof("CreateServiceAccount: %s", sa.Id)

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	if err := s.credentialManager.CreateServiceAccount(context.Background(), sa); err != nil {
		glog.Errorf("Failed to create service account %s: %v", sa.Id, err)
		return nil, fmt.Errorf("failed to create service account: %v", err)
	}

	return iamwire.NewCreateServiceAccountResponse(iamwire.CreateServiceAccountResponseInput{}), nil
}

func (s *IamGrpcServer) UpdateServiceAccount(req []byte) ([]byte, error) {
	request, err := iamwire.WrapUpdateServiceAccountRequest(req)
	if err != nil {
		return nil, err
	}
	id := request.ID()
	sa := iamwire.ServiceAccountToPB(request.ServiceAccount())
	glog.V(4).Infof("UpdateServiceAccount: %s", id)

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	if err := s.credentialManager.UpdateServiceAccount(context.Background(), id, sa); err != nil {
		glog.Errorf("Failed to update service account %s: %v", id, err)
		return nil, fmt.Errorf("failed to update service account: %v", err)
	}

	return iamwire.NewUpdateServiceAccountResponse(iamwire.UpdateServiceAccountResponseInput{}), nil
}

func (s *IamGrpcServer) DeleteServiceAccount(req []byte) ([]byte, error) {
	request, err := iamwire.WrapDeleteServiceAccountRequest(req)
	if err != nil {
		return nil, err
	}
	id := request.ID()
	glog.V(4).Infof("DeleteServiceAccount: %s", id)

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	if err := s.credentialManager.DeleteServiceAccount(context.Background(), id); err != nil {
		if err == credential.ErrServiceAccountNotFound {
			return nil, fmt.Errorf("service account %s not found", id)
		}
		glog.Errorf("Failed to delete service account %s: %v", id, err)
		return nil, fmt.Errorf("failed to delete service account: %v", err)
	}

	return iamwire.NewDeleteServiceAccountResponse(iamwire.DeleteServiceAccountResponseInput{}), nil
}

func (s *IamGrpcServer) GetServiceAccount(req []byte) ([]byte, error) {
	request, err := iamwire.WrapGetServiceAccountRequest(req)
	if err != nil {
		return nil, err
	}
	id := request.ID()
	glog.V(4).Infof("GetServiceAccount: %s", id)

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	sa, err := s.credentialManager.GetServiceAccount(context.Background(), id)
	if err != nil {
		glog.Errorf("Failed to get service account %s: %v", id, err)
		return nil, fmt.Errorf("failed to get service account: %v", err)
	}

	if sa == nil {
		return nil, fmt.Errorf("service account %s not found", id)
	}

	return iamwire.NewGetServiceAccountResponse(iamwire.GetServiceAccountResponseInput{
		ServiceAccount: iamwire.ServiceAccountInputFromPB(sa),
	}), nil
}

func (s *IamGrpcServer) ListServiceAccounts(req []byte) ([]byte, error) {
	if _, err := iamwire.WrapListServiceAccountsRequest(req); err != nil {
		return nil, err
	}
	glog.V(4).Infof("ListServiceAccounts")

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	accounts, err := s.credentialManager.ListServiceAccounts(context.Background())
	if err != nil {
		glog.Errorf("Failed to list service accounts: %v", err)
		return nil, fmt.Errorf("failed to list service accounts: %v", err)
	}

	saInputs := make([]iamwire.ServiceAccountInput, 0, len(accounts))
	for _, sa := range accounts {
		if v := iamwire.ServiceAccountInputFromPB(sa); v != nil {
			saInputs = append(saInputs, *v)
		}
	}

	return iamwire.NewListServiceAccountsResponse(iamwire.ListServiceAccountsResponseInput{
		ServiceAccounts: saInputs,
	}), nil
}

func (s *IamGrpcServer) GetServiceAccountByAccessKey(req []byte) ([]byte, error) {
	request, err := iamwire.WrapGetServiceAccountByAccessKeyRequest(req)
	if err != nil {
		return nil, err
	}
	accessKey := request.AccessKey()
	glog.V(4).Infof("GetServiceAccountByAccessKey: %s", accessKey)
	if accessKey == "" {
		return nil, fmt.Errorf("access key is required")
	}

	if s.credentialManager == nil {
		return nil, fmt.Errorf("credential manager is not configured")
	}

	sa, err := s.credentialManager.GetStore().GetServiceAccountByAccessKey(context.Background(), accessKey)
	if err != nil {
		if err == credential.ErrAccessKeyNotFound {
			return nil, fmt.Errorf("access key %s not found", accessKey)
		}
		glog.Errorf("Failed to get service account by access key %s: %v", accessKey, err)
		return nil, fmt.Errorf("failed to get service account: %v", err)
	}

	return iamwire.NewGetServiceAccountByAccessKeyResponse(iamwire.GetServiceAccountByAccessKeyResponseInput{
		ServiceAccount: iamwire.ServiceAccountInputFromPB(sa),
	}), nil
}
