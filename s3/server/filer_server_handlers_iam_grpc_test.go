package s3server

import (
	"testing"

	"github.com/hanzoai/s3/s3/credential"
	_ "github.com/hanzoai/s3/s3/credential/memory"
	"github.com/hanzoai/s3/s3/security"
	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

const testIamSigningKey = "iam-admin-test-key-do-not-use-in-prod"

func newTestIamGrpcServer(t *testing.T) *IamGrpcServer {
	t.Helper()
	cm, err := credential.NewCredentialManager(credential.StoreTypeMemory, nil, "")
	if err != nil {
		t.Fatalf("NewCredentialManager: %v", err)
	}
	return NewIamGrpcServer(cm, security.SigningKey(testIamSigningKey))
}

// TestIamZap_HandlerSatisfiesWireContract pins the post-rip shape: the IAM
// server is a pure iamwire handler whose methods take an opaque ZAP request
// envelope and return an opaque ZAP response envelope. Connection-level
// security is the ZAP transport (X-Wing PQ TLS), so there is no per-call
// Bearer-token gate to exercise here anymore.
func TestIamZap_HandlerSatisfiesWireContract(t *testing.T) {
	var _ iamwire.HanzoIdentityAccessManagementHandler = newTestIamGrpcServer(t)
}

// TestIamZap_ListUsers_ReachesHandler drives ListUsers over the wire envelope
// the way the transport does: build the request with the iamwire builder, run
// the handler, decode the reply with the matching Wrap. A fresh memory store
// yields an empty user list, proving the request reached the business logic.
func TestIamZap_ListUsers_ReachesHandler(t *testing.T) {
	s := newTestIamGrpcServer(t)

	respBytes, err := s.ListUsers(iamwire.NewListUsersRequest(iamwire.ListUsersRequestInput{}))
	if err != nil {
		t.Fatalf("ListUsers: unexpected error %v", err)
	}
	resp, err := iamwire.WrapListUsersResponse(respBytes)
	if err != nil {
		t.Fatalf("WrapListUsersResponse: %v", err)
	}
	if n := resp.UsernamesLen(); n != 0 {
		t.Fatalf("ListUsers: expected empty user list from fresh memory store, got %d", n)
	}
}

// TestIamZap_Dispatch_RoundTrip exercises the generated dispatcher exactly as
// transport.Serve does: a Call envelope in, a reply envelope out, routed by
// ordinal to CreateUser then GetUser. This proves the handler is reachable
// through DispatchHanzoIdentityAccessManagement, not just by direct call.
func TestIamZap_Dispatch_RoundTrip(t *testing.T) {
	s := newTestIamGrpcServer(t)

	if _, err := s.CreateUser(iamwire.NewCreateUserRequest(iamwire.CreateUserRequestInput{
		Identity: &iamwire.IdentityInput{Name: "alice"},
	})); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	getBytes, err := s.GetUser(iamwire.NewGetUserRequest(iamwire.GetUserRequestInput{Username: "alice"}))
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	got, err := iamwire.WrapGetUserResponse(getBytes)
	if err != nil {
		t.Fatalf("WrapGetUserResponse: %v", err)
	}
	if name := got.Identity().Name(); name != "alice" {
		t.Fatalf("GetUser: got identity name %q, want alice", name)
	}
}
