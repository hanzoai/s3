// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package iamwire

import (
	"testing"

	"github.com/hanzoai/s3/s3/pb/iam_pb"
)

// memIam is an in-memory IamStore for the round-trip test. It implements the
// pure-[]byte IamStore contract the same way the real backend would: Wrap the
// request into its zero-copy view, translate to the in-process iam_pb model via
// the bridge_pb.go helpers, mutate state, and build the reply with the matching
// New* builder. This exercises the entire pb<->wire path on both ends, so the
// test proves the nested-struct + repeated-field builders survive a real socket.
type memIam struct {
	users    map[string]*iam_pb.Identity
	byKey    map[string]*iam_pb.Identity
	policies map[string]string
	sas      map[string]*iam_pb.ServiceAccount
	saByKey  map[string]*iam_pb.ServiceAccount
}

func newMemIam() *memIam {
	return &memIam{
		users:    map[string]*iam_pb.Identity{},
		byKey:    map[string]*iam_pb.Identity{},
		policies: map[string]string{},
		sas:      map[string]*iam_pb.ServiceAccount{},
		saByKey:  map[string]*iam_pb.ServiceAccount{},
	}
}

func (m *memIam) indexKeys(id *iam_pb.Identity) {
	for _, c := range id.Credentials {
		m.byKey[c.AccessKey] = id
	}
}

func (m *memIam) CreateUser(req []byte) ([]byte, error) {
	v, err := WrapCreateUserRequest(req)
	if err != nil {
		return nil, err
	}
	id := IdentityToPB(v.Identity())
	m.users[id.Name] = id
	m.indexKeys(id)
	return NewCreateUserResponse(CreateUserResponseInput{}), nil
}

func (m *memIam) GetUser(req []byte) ([]byte, error) {
	v, err := WrapGetUserRequest(req)
	if err != nil {
		return nil, err
	}
	id := m.users[v.Username()] // nil when absent — encodes as no identity field
	return NewGetUserResponse(GetUserResponseInput{Identity: IdentityInputFromPB(id)}), nil
}

func (m *memIam) UpdateUser(req []byte) ([]byte, error) {
	v, err := WrapUpdateUserRequest(req)
	if err != nil {
		return nil, err
	}
	id := IdentityToPB(v.Identity())
	m.users[v.Username()] = id
	m.indexKeys(id)
	return NewUpdateUserResponse(UpdateUserResponseInput{}), nil
}

func (m *memIam) DeleteUser(req []byte) ([]byte, error) {
	v, err := WrapDeleteUserRequest(req)
	if err != nil {
		return nil, err
	}
	delete(m.users, v.Username())
	return NewDeleteUserResponse(DeleteUserResponseInput{}), nil
}

func (m *memIam) ListUsers(req []byte) ([]byte, error) {
	names := make([]string, 0, len(m.users))
	for n := range m.users {
		names = append(names, n)
	}
	return NewListUsersResponse(ListUsersResponseInput{Usernames: names}), nil
}

func (m *memIam) CreateAccessKey(req []byte) ([]byte, error) {
	v, err := WrapCreateAccessKeyRequest(req)
	if err != nil {
		return nil, err
	}
	if id := m.users[v.Username()]; id != nil {
		cred := CredentialToPBExported(v.Credential())
		id.Credentials = append(id.Credentials, cred)
		m.byKey[cred.AccessKey] = id
	}
	return NewCreateAccessKeyResponse(CreateAccessKeyResponseInput{}), nil
}

func (m *memIam) DeleteAccessKey(req []byte) ([]byte, error) {
	v, err := WrapDeleteAccessKeyRequest(req)
	if err != nil {
		return nil, err
	}
	delete(m.byKey, v.AccessKey())
	return NewDeleteAccessKeyResponse(DeleteAccessKeyResponseInput{}), nil
}

func (m *memIam) GetUserByAccessKey(req []byte) ([]byte, error) {
	v, err := WrapGetUserByAccessKeyRequest(req)
	if err != nil {
		return nil, err
	}
	id := m.byKey[v.AccessKey()]
	return NewGetUserByAccessKeyResponse(GetUserByAccessKeyResponseInput{Identity: IdentityInputFromPB(id)}), nil
}

func (m *memIam) PutPolicy(req []byte) ([]byte, error) {
	v, err := WrapPutPolicyRequest(req)
	if err != nil {
		return nil, err
	}
	m.policies[v.Name()] = v.Content()
	return NewPutPolicyResponse(PutPolicyResponseInput{}), nil
}

func (m *memIam) GetPolicy(req []byte) ([]byte, error) {
	v, err := WrapGetPolicyRequest(req)
	if err != nil {
		return nil, err
	}
	return NewGetPolicyResponse(GetPolicyResponseInput{Name: v.Name(), Content: m.policies[v.Name()]}), nil
}

func (m *memIam) ListPolicies(req []byte) ([]byte, error) {
	in := ListPoliciesResponseInput{}
	for name, content := range m.policies {
		in.Policies = append(in.Policies, PolicyInput{Name: name, Content: content})
	}
	return NewListPoliciesResponse(in), nil
}

func (m *memIam) DeletePolicy(req []byte) ([]byte, error) {
	v, err := WrapDeletePolicyRequest(req)
	if err != nil {
		return nil, err
	}
	delete(m.policies, v.Name())
	return NewDeletePolicyResponse(DeletePolicyResponseInput{}), nil
}

func (m *memIam) GetConfiguration(req []byte) ([]byte, error) {
	cfg := &iam_pb.S3ApiConfiguration{}
	for _, id := range m.users {
		cfg.Identities = append(cfg.Identities, id)
	}
	return NewGetConfigurationResponse(GetConfigurationResponseInput{
		Configuration: ConfigurationInputFromPB(cfg),
	}), nil
}

func (m *memIam) PutConfiguration(req []byte) ([]byte, error) {
	v, err := WrapPutConfigurationRequest(req)
	if err != nil {
		return nil, err
	}
	cfg := ConfigurationToPB(v.Configuration())
	for _, id := range cfg.Identities {
		m.users[id.Name] = id
		m.indexKeys(id)
	}
	return NewPutConfigurationResponse(PutConfigurationResponseInput{}), nil
}

func (m *memIam) CreateServiceAccount(req []byte) ([]byte, error) {
	v, err := WrapCreateServiceAccountRequest(req)
	if err != nil {
		return nil, err
	}
	sa := ServiceAccountToPB(v.ServiceAccount())
	m.sas[sa.Id] = sa
	if sa.Credential != nil {
		m.saByKey[sa.Credential.AccessKey] = sa
	}
	return NewCreateServiceAccountResponse(CreateServiceAccountResponseInput{}), nil
}

func (m *memIam) UpdateServiceAccount(req []byte) ([]byte, error) {
	v, err := WrapUpdateServiceAccountRequest(req)
	if err != nil {
		return nil, err
	}
	sa := ServiceAccountToPB(v.ServiceAccount())
	m.sas[v.ID()] = sa
	return NewUpdateServiceAccountResponse(UpdateServiceAccountResponseInput{}), nil
}

func (m *memIam) DeleteServiceAccount(req []byte) ([]byte, error) {
	v, err := WrapDeleteServiceAccountRequest(req)
	if err != nil {
		return nil, err
	}
	delete(m.sas, v.ID())
	return NewDeleteServiceAccountResponse(DeleteServiceAccountResponseInput{}), nil
}

func (m *memIam) GetServiceAccount(req []byte) ([]byte, error) {
	v, err := WrapGetServiceAccountRequest(req)
	if err != nil {
		return nil, err
	}
	return NewGetServiceAccountResponse(GetServiceAccountResponseInput{
		ServiceAccount: ServiceAccountInputFromPB(m.sas[v.ID()]),
	}), nil
}

func (m *memIam) ListServiceAccounts(req []byte) ([]byte, error) {
	in := ListServiceAccountsResponseInput{}
	for _, sa := range m.sas {
		if v := ServiceAccountInputFromPB(sa); v != nil {
			in.ServiceAccounts = append(in.ServiceAccounts, *v)
		}
	}
	return NewListServiceAccountsResponse(in), nil
}

func (m *memIam) GetServiceAccountByAccessKey(req []byte) ([]byte, error) {
	v, err := WrapGetServiceAccountByAccessKeyRequest(req)
	if err != nil {
		return nil, err
	}
	return NewGetServiceAccountByAccessKeyResponse(GetServiceAccountByAccessKeyResponseInput{
		ServiceAccount: ServiceAccountInputFromPB(m.saByKey[v.AccessKey()]),
	}), nil
}

// TestRoundTrip proves the HanzoIdentityAccessManagement service end-to-end over
// the canonical github.com/zap-proto/go transport across a real TCP socket: the
// bytes cross the wire as ZAP RPC envelopes carrying zero-copy iamwire payloads,
// correlated by PromiseID — no HTTP, no protobuf on the wire, no gRPC.
func TestRoundTrip(t *testing.T) {
	store := newMemIam()

	srv, err := Serve("tcp", "127.0.0.1:0", store)
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	// A fully-populated identity: nested Account, two credentials, repeated
	// actions, repeated tags, and the bool/list scalars — every wire kind the
	// builders must round-trip.
	want := &iam_pb.Identity{
		Name: "alice",
		Credentials: []*iam_pb.Credential{
			{AccessKey: "AKIA1", SecretKey: "s1", Status: "Active"},
			{AccessKey: "AKIA2", SecretKey: "s2", Status: "Inactive"},
		},
		Actions:           []string{"Read", "Write", "Admin:*"},
		Account:           &iam_pb.Account{Id: "acct-1", DisplayName: "Alice", EmailAddress: "alice@hanzo.ai"},
		Disabled:          false,
		ServiceAccountIds: []string{"sa-1"},
		PolicyNames:       []string{"p-read", "p-write"},
		IsStatic:          true,
		Tags:              []*iam_pb.UserTag{{Key: "team", Value: "platform"}, {Key: "env", Value: "prod"}},
	}

	// CreateUser -> GetUser.
	if err := cli.CreateUser(want); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	got, err := cli.GetUser("alice")
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if got == nil {
		t.Fatalf("GetUser returned nil identity")
	}
	assertIdentityEqual(t, want, got)

	// GetUser of an unknown user yields a nil identity (gRPC nil-pointer
	// semantics preserved through the wire's absent-message encoding).
	missing, err := cli.GetUser("nobody")
	if err != nil {
		t.Fatalf("GetUser(missing): %v", err)
	}
	if missing != nil {
		t.Fatalf("GetUser(missing) = %+v, want nil", missing)
	}

	// GetUserByAccessKey resolves the indexed key back to alice.
	byKey, err := cli.GetUserByAccessKey("AKIA2")
	if err != nil {
		t.Fatalf("GetUserByAccessKey: %v", err)
	}
	if byKey == nil || byKey.Name != "alice" {
		t.Fatalf("GetUserByAccessKey = %+v, want alice", byKey)
	}

	// ListUsers sees exactly the one user.
	names, err := cli.ListUsers()
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	if len(names) != 1 || names[0] != "alice" {
		t.Fatalf("ListUsers = %v, want [alice]", names)
	}

	// PutPolicy -> GetPolicy round-trips the JSON document.
	const doc = `{"Version":"2012-10-17","Statement":[]}`
	if err := cli.PutPolicy("p-read", doc); err != nil {
		t.Fatalf("PutPolicy: %v", err)
	}
	gotName, gotDoc, err := cli.GetPolicy("p-read")
	if err != nil {
		t.Fatalf("GetPolicy: %v", err)
	}
	if gotName != "p-read" || gotDoc != doc {
		t.Fatalf("GetPolicy = (%q,%q), want (p-read,%q)", gotName, gotDoc, doc)
	}
	pols, err := cli.ListPolicies()
	if err != nil {
		t.Fatalf("ListPolicies: %v", err)
	}
	if len(pols) != 1 || pols[0].Name != "p-read" || pols[0].Content != doc {
		t.Fatalf("ListPolicies = %+v, want one p-read", pols)
	}

	// Service accounts: Create -> Get -> List, exercising the int64 + nested
	// Credential fields.
	sa := &iam_pb.ServiceAccount{
		Id:         "sa-1",
		ParentUser: "alice",
		Credential: &iam_pb.Credential{AccessKey: "SAKEY1", SecretKey: "sas1", Status: "Active"},
		Actions:    []string{"Read"},
		Expiration: 1893456000,
		CreatedAt:  1766708400,
		CreatedBy:  "alice",
	}
	if err := cli.CreateServiceAccount(sa); err != nil {
		t.Fatalf("CreateServiceAccount: %v", err)
	}
	gotSA, err := cli.GetServiceAccount("sa-1")
	if err != nil {
		t.Fatalf("GetServiceAccount: %v", err)
	}
	if gotSA == nil || gotSA.Id != "sa-1" || gotSA.ParentUser != "alice" ||
		gotSA.Credential == nil || gotSA.Credential.AccessKey != "SAKEY1" ||
		gotSA.Expiration != sa.Expiration || gotSA.CreatedAt != sa.CreatedAt {
		t.Fatalf("GetServiceAccount = %+v, want sa-1", gotSA)
	}
	byKeySA, err := cli.GetServiceAccountByAccessKey("SAKEY1")
	if err != nil {
		t.Fatalf("GetServiceAccountByAccessKey: %v", err)
	}
	if byKeySA == nil || byKeySA.Id != "sa-1" {
		t.Fatalf("GetServiceAccountByAccessKey = %+v, want sa-1", byKeySA)
	}
	sas, err := cli.ListServiceAccounts()
	if err != nil {
		t.Fatalf("ListServiceAccounts: %v", err)
	}
	if len(sas) != 1 || sas[0].Id != "sa-1" {
		t.Fatalf("ListServiceAccounts = %+v, want one sa-1", sas)
	}

	// GetConfiguration reflects the user we created (nested repeated tree).
	cfg, err := cli.GetConfiguration()
	if err != nil {
		t.Fatalf("GetConfiguration: %v", err)
	}
	if len(cfg.Identities) != 1 || cfg.Identities[0].Name != "alice" {
		t.Fatalf("GetConfiguration identities = %+v, want one alice", cfg.Identities)
	}

	// DeleteUser then confirm it is gone.
	if err := cli.DeleteUser("alice"); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	if gone, _ := cli.GetUser("alice"); gone != nil {
		t.Fatalf("GetUser after delete = %+v, want nil", gone)
	}
}

// assertIdentityEqual checks the fields that crossed the wire match, covering
// every wire kind (text, bool, repeated text, repeated nested struct, nested
// struct).
func assertIdentityEqual(t *testing.T, want, got *iam_pb.Identity) {
	t.Helper()
	if got.Name != want.Name {
		t.Fatalf("Name = %q, want %q", got.Name, want.Name)
	}
	if got.Disabled != want.Disabled || got.IsStatic != want.IsStatic {
		t.Fatalf("flags: Disabled=%v IsStatic=%v, want %v/%v", got.Disabled, got.IsStatic, want.Disabled, want.IsStatic)
	}
	if len(got.Credentials) != len(want.Credentials) {
		t.Fatalf("Credentials len = %d, want %d", len(got.Credentials), len(want.Credentials))
	}
	for i := range want.Credentials {
		if got.Credentials[i].AccessKey != want.Credentials[i].AccessKey ||
			got.Credentials[i].SecretKey != want.Credentials[i].SecretKey ||
			got.Credentials[i].Status != want.Credentials[i].Status {
			t.Fatalf("Credentials[%d] = %+v, want %+v", i, got.Credentials[i], want.Credentials[i])
		}
	}
	if len(got.Actions) != len(want.Actions) {
		t.Fatalf("Actions = %v, want %v", got.Actions, want.Actions)
	}
	for i := range want.Actions {
		if got.Actions[i] != want.Actions[i] {
			t.Fatalf("Actions[%d] = %q, want %q", i, got.Actions[i], want.Actions[i])
		}
	}
	if got.Account == nil || got.Account.Id != want.Account.Id ||
		got.Account.DisplayName != want.Account.DisplayName ||
		got.Account.EmailAddress != want.Account.EmailAddress {
		t.Fatalf("Account = %+v, want %+v", got.Account, want.Account)
	}
	if len(got.ServiceAccountIds) != len(want.ServiceAccountIds) {
		t.Fatalf("ServiceAccountIds = %v, want %v", got.ServiceAccountIds, want.ServiceAccountIds)
	}
	if len(got.PolicyNames) != len(want.PolicyNames) {
		t.Fatalf("PolicyNames = %v, want %v", got.PolicyNames, want.PolicyNames)
	}
	if len(got.Tags) != len(want.Tags) {
		t.Fatalf("Tags len = %d, want %d", len(got.Tags), len(want.Tags))
	}
	for i := range want.Tags {
		if got.Tags[i].Key != want.Tags[i].Key || got.Tags[i].Value != want.Tags[i].Value {
			t.Fatalf("Tags[%d] = %+v, want %+v", i, got.Tags[i], want.Tags[i])
		}
	}
}
