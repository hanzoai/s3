// Code generated from iam.proto; DO NOT EDIT.

package iamwire

import (
	zap "github.com/zap-proto/go"
)

// Access-key management and user-policy request/response messages.

// --- CreateAccessKeyRequest ---

const (
	createAccessKeyRequestUsernameOff   = 0
	createAccessKeyRequestCredentialOff = 8
	createAccessKeyRequestSize          = 16
)

// CreateAccessKeyRequest is a zero-copy view into a ZAP-encoded
// CreateAccessKeyRequest message.
type CreateAccessKeyRequest struct{ o zap.Object }

// WrapCreateAccessKeyRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapCreateAccessKeyRequest(b []byte) (CreateAccessKeyRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CreateAccessKeyRequest{}, err
	}
	return CreateAccessKeyRequest{o: m.Root()}, nil
}

// Username reads the username field (proto field 1, string).
func (t CreateAccessKeyRequest) Username() string {
	return t.o.Text(createAccessKeyRequestUsernameOff)
}

// Credential reads the credential field (proto field 2, Credential).
func (t CreateAccessKeyRequest) Credential() Credential {
	return Credential{o: t.o.Object(createAccessKeyRequestCredentialOff)}
}

// CreateAccessKeyRequestInput collects the field values for
// NewCreateAccessKeyRequest.
type CreateAccessKeyRequestInput struct {
	Username   string
	Credential *CredentialInput
}

// NewCreateAccessKeyRequest builds a ZAP-encoded CreateAccessKeyRequest message
// from in and returns the bytes.
func NewCreateAccessKeyRequest(in CreateAccessKeyRequestInput) []byte {
	b := zap.NewBuilder(512)
	credentialOff := 0
	if in.Credential != nil {
		credentialOff = appendCredential(b, *in.Credential).Finish()
	}
	ob := b.StartObject(createAccessKeyRequestSize)
	ob.SetText(createAccessKeyRequestUsernameOff, in.Username)
	if credentialOff != 0 {
		ob.SetObject(createAccessKeyRequestCredentialOff, credentialOff)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- CreateAccessKeyResponse ---

const createAccessKeyResponseSize = 0

// CreateAccessKeyResponse is a zero-copy view into a ZAP-encoded
// CreateAccessKeyResponse message. It carries no fields.
type CreateAccessKeyResponse struct{ o zap.Object }

// WrapCreateAccessKeyResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapCreateAccessKeyResponse(b []byte) (CreateAccessKeyResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CreateAccessKeyResponse{}, err
	}
	return CreateAccessKeyResponse{o: m.Root()}, nil
}

// CreateAccessKeyResponseInput collects the field values for
// NewCreateAccessKeyResponse. CreateAccessKeyResponse has no fields.
type CreateAccessKeyResponseInput struct{}

// NewCreateAccessKeyResponse builds a ZAP-encoded CreateAccessKeyResponse
// message and returns the bytes.
func NewCreateAccessKeyResponse(in CreateAccessKeyResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(createAccessKeyResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DeleteAccessKeyRequest ---

const (
	deleteAccessKeyRequestUsernameOff  = 0
	deleteAccessKeyRequestAccessKeyOff = 8
	deleteAccessKeyRequestSize         = 16
)

// DeleteAccessKeyRequest is a zero-copy view into a ZAP-encoded
// DeleteAccessKeyRequest message.
type DeleteAccessKeyRequest struct{ o zap.Object }

// WrapDeleteAccessKeyRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapDeleteAccessKeyRequest(b []byte) (DeleteAccessKeyRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DeleteAccessKeyRequest{}, err
	}
	return DeleteAccessKeyRequest{o: m.Root()}, nil
}

// Username reads the username field (proto field 1, string).
func (t DeleteAccessKeyRequest) Username() string {
	return t.o.Text(deleteAccessKeyRequestUsernameOff)
}

// AccessKey reads the access_key field (proto field 2, string).
func (t DeleteAccessKeyRequest) AccessKey() string {
	return t.o.Text(deleteAccessKeyRequestAccessKeyOff)
}

// DeleteAccessKeyRequestInput collects the field values for
// NewDeleteAccessKeyRequest.
type DeleteAccessKeyRequestInput struct {
	Username  string
	AccessKey string
}

// NewDeleteAccessKeyRequest builds a ZAP-encoded DeleteAccessKeyRequest message
// from in and returns the bytes.
func NewDeleteAccessKeyRequest(in DeleteAccessKeyRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(deleteAccessKeyRequestSize)
	ob.SetText(deleteAccessKeyRequestUsernameOff, in.Username)
	ob.SetText(deleteAccessKeyRequestAccessKeyOff, in.AccessKey)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DeleteAccessKeyResponse ---

const deleteAccessKeyResponseSize = 0

// DeleteAccessKeyResponse is a zero-copy view into a ZAP-encoded
// DeleteAccessKeyResponse message. It carries no fields.
type DeleteAccessKeyResponse struct{ o zap.Object }

// WrapDeleteAccessKeyResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapDeleteAccessKeyResponse(b []byte) (DeleteAccessKeyResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DeleteAccessKeyResponse{}, err
	}
	return DeleteAccessKeyResponse{o: m.Root()}, nil
}

// DeleteAccessKeyResponseInput collects the field values for
// NewDeleteAccessKeyResponse. DeleteAccessKeyResponse has no fields.
type DeleteAccessKeyResponseInput struct{}

// NewDeleteAccessKeyResponse builds a ZAP-encoded DeleteAccessKeyResponse
// message and returns the bytes.
func NewDeleteAccessKeyResponse(in DeleteAccessKeyResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(deleteAccessKeyResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetUserByAccessKeyRequest ---

const (
	getUserByAccessKeyRequestAccessKeyOff = 0
	getUserByAccessKeyRequestSize         = 8
)

// GetUserByAccessKeyRequest is a zero-copy view into a ZAP-encoded
// GetUserByAccessKeyRequest message.
type GetUserByAccessKeyRequest struct{ o zap.Object }

// WrapGetUserByAccessKeyRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapGetUserByAccessKeyRequest(b []byte) (GetUserByAccessKeyRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetUserByAccessKeyRequest{}, err
	}
	return GetUserByAccessKeyRequest{o: m.Root()}, nil
}

// AccessKey reads the access_key field (proto field 1, string).
func (t GetUserByAccessKeyRequest) AccessKey() string {
	return t.o.Text(getUserByAccessKeyRequestAccessKeyOff)
}

// GetUserByAccessKeyRequestInput collects the field values for
// NewGetUserByAccessKeyRequest.
type GetUserByAccessKeyRequestInput struct {
	AccessKey string
}

// NewGetUserByAccessKeyRequest builds a ZAP-encoded GetUserByAccessKeyRequest
// message from in and returns the bytes.
func NewGetUserByAccessKeyRequest(in GetUserByAccessKeyRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getUserByAccessKeyRequestSize)
	ob.SetText(getUserByAccessKeyRequestAccessKeyOff, in.AccessKey)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetUserByAccessKeyResponse ---

const (
	getUserByAccessKeyResponseIdentityOff = 0
	getUserByAccessKeyResponseSize        = 4
)

// GetUserByAccessKeyResponse is a zero-copy view into a ZAP-encoded
// GetUserByAccessKeyResponse message.
type GetUserByAccessKeyResponse struct{ o zap.Object }

// WrapGetUserByAccessKeyResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapGetUserByAccessKeyResponse(b []byte) (GetUserByAccessKeyResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetUserByAccessKeyResponse{}, err
	}
	return GetUserByAccessKeyResponse{o: m.Root()}, nil
}

// Identity reads the identity field (proto field 1, Identity).
func (t GetUserByAccessKeyResponse) Identity() Identity {
	return Identity{o: t.o.Object(getUserByAccessKeyResponseIdentityOff)}
}

// GetUserByAccessKeyResponseInput collects the field values for
// NewGetUserByAccessKeyResponse.
type GetUserByAccessKeyResponseInput struct {
	Identity *IdentityInput
}

// NewGetUserByAccessKeyResponse builds a ZAP-encoded GetUserByAccessKeyResponse
// message from in and returns the bytes.
func NewGetUserByAccessKeyResponse(in GetUserByAccessKeyResponseInput) []byte {
	b := zap.NewBuilder(512)
	identityOff := 0
	if in.Identity != nil {
		identityOff = appendIdentity(b, *in.Identity).Finish()
	}
	ob := b.StartObject(getUserByAccessKeyResponseSize)
	if identityOff != 0 {
		ob.SetObject(getUserByAccessKeyResponseIdentityOff, identityOff)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListAccessKeysRequest ---

const (
	listAccessKeysRequestUsernameOff = 0
	listAccessKeysRequestSize        = 8
)

// ListAccessKeysRequest is a zero-copy view into a ZAP-encoded
// ListAccessKeysRequest message.
type ListAccessKeysRequest struct{ o zap.Object }

// WrapListAccessKeysRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapListAccessKeysRequest(b []byte) (ListAccessKeysRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListAccessKeysRequest{}, err
	}
	return ListAccessKeysRequest{o: m.Root()}, nil
}

// Username reads the username field (proto field 1, string).
func (t ListAccessKeysRequest) Username() string {
	return t.o.Text(listAccessKeysRequestUsernameOff)
}

// ListAccessKeysRequestInput collects the field values for
// NewListAccessKeysRequest.
type ListAccessKeysRequestInput struct {
	Username string
}

// NewListAccessKeysRequest builds a ZAP-encoded ListAccessKeysRequest message
// from in and returns the bytes.
func NewListAccessKeysRequest(in ListAccessKeysRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(listAccessKeysRequestSize)
	ob.SetText(listAccessKeysRequestUsernameOff, in.Username)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListAccessKeysResponse ---

const (
	listAccessKeysResponseAccessKeysOff = 0
	listAccessKeysResponseSize          = 8
)

// ListAccessKeysResponse is a zero-copy view into a ZAP-encoded
// ListAccessKeysResponse message.
type ListAccessKeysResponse struct{ o zap.Object }

// WrapListAccessKeysResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapListAccessKeysResponse(b []byte) (ListAccessKeysResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListAccessKeysResponse{}, err
	}
	return ListAccessKeysResponse{o: m.Root()}, nil
}

// AccessKeysLen returns the number of access keys (proto field 1, repeated Credential).
func (t ListAccessKeysResponse) AccessKeysLen() int {
	return t.o.List(listAccessKeysResponseAccessKeysOff).Len()
}

// AccessKeysAt returns access key i (proto field 1, repeated Credential).
func (t ListAccessKeysResponse) AccessKeysAt(i int) Credential {
	return Credential{o: t.o.List(listAccessKeysResponseAccessKeysOff).ObjectAt(i)}
}

// ListAccessKeysResponseInput collects the field values for
// NewListAccessKeysResponse.
type ListAccessKeysResponseInput struct {
	AccessKeys []CredentialInput
}

// NewListAccessKeysResponse builds a ZAP-encoded ListAccessKeysResponse message
// from in and returns the bytes.
func NewListAccessKeysResponse(in ListAccessKeysResponseInput) []byte {
	b := zap.NewBuilder(512)
	ob := b.StartObject(listAccessKeysResponseSize)
	if len(in.AccessKeys) > 0 {
		lb := b.StartList(0)
		for i := range in.AccessKeys {
			lb.AddObjectBytes(NewCredential(in.AccessKeys[i]))
		}
		off, n := lb.Finish()
		ob.SetList(listAccessKeysResponseAccessKeysOff, off, n)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PutUserPolicyRequest ---

const (
	putUserPolicyRequestUsernameOff       = 0
	putUserPolicyRequestPolicyNameOff     = 8
	putUserPolicyRequestPolicyDocumentOff = 16
	putUserPolicyRequestSize              = 24
)

// PutUserPolicyRequest is a zero-copy view into a ZAP-encoded
// PutUserPolicyRequest message.
type PutUserPolicyRequest struct{ o zap.Object }

// WrapPutUserPolicyRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapPutUserPolicyRequest(b []byte) (PutUserPolicyRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PutUserPolicyRequest{}, err
	}
	return PutUserPolicyRequest{o: m.Root()}, nil
}

// Username reads the username field (proto field 1, string).
func (t PutUserPolicyRequest) Username() string { return t.o.Text(putUserPolicyRequestUsernameOff) }

// PolicyName reads the policy_name field (proto field 2, string).
func (t PutUserPolicyRequest) PolicyName() string {
	return t.o.Text(putUserPolicyRequestPolicyNameOff)
}

// PolicyDocument reads the policy_document field (proto field 3, string).
func (t PutUserPolicyRequest) PolicyDocument() string {
	return t.o.Text(putUserPolicyRequestPolicyDocumentOff)
}

// PutUserPolicyRequestInput collects the field values for NewPutUserPolicyRequest.
type PutUserPolicyRequestInput struct {
	Username       string
	PolicyName     string
	PolicyDocument string
}

// NewPutUserPolicyRequest builds a ZAP-encoded PutUserPolicyRequest message from
// in and returns the bytes.
func NewPutUserPolicyRequest(in PutUserPolicyRequestInput) []byte {
	b := zap.NewBuilder(512)
	ob := b.StartObject(putUserPolicyRequestSize)
	ob.SetText(putUserPolicyRequestUsernameOff, in.Username)
	ob.SetText(putUserPolicyRequestPolicyNameOff, in.PolicyName)
	ob.SetText(putUserPolicyRequestPolicyDocumentOff, in.PolicyDocument)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PutUserPolicyResponse ---

const putUserPolicyResponseSize = 0

// PutUserPolicyResponse is a zero-copy view into a ZAP-encoded
// PutUserPolicyResponse message. It carries no fields.
type PutUserPolicyResponse struct{ o zap.Object }

// WrapPutUserPolicyResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapPutUserPolicyResponse(b []byte) (PutUserPolicyResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PutUserPolicyResponse{}, err
	}
	return PutUserPolicyResponse{o: m.Root()}, nil
}

// PutUserPolicyResponseInput collects the field values for
// NewPutUserPolicyResponse. PutUserPolicyResponse has no fields.
type PutUserPolicyResponseInput struct{}

// NewPutUserPolicyResponse builds a ZAP-encoded PutUserPolicyResponse message
// and returns the bytes.
func NewPutUserPolicyResponse(in PutUserPolicyResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(putUserPolicyResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetUserPolicyRequest ---

const (
	getUserPolicyRequestUsernameOff   = 0
	getUserPolicyRequestPolicyNameOff = 8
	getUserPolicyRequestSize          = 16
)

// GetUserPolicyRequest is a zero-copy view into a ZAP-encoded
// GetUserPolicyRequest message.
type GetUserPolicyRequest struct{ o zap.Object }

// WrapGetUserPolicyRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapGetUserPolicyRequest(b []byte) (GetUserPolicyRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetUserPolicyRequest{}, err
	}
	return GetUserPolicyRequest{o: m.Root()}, nil
}

// Username reads the username field (proto field 1, string).
func (t GetUserPolicyRequest) Username() string { return t.o.Text(getUserPolicyRequestUsernameOff) }

// PolicyName reads the policy_name field (proto field 2, string).
func (t GetUserPolicyRequest) PolicyName() string {
	return t.o.Text(getUserPolicyRequestPolicyNameOff)
}

// GetUserPolicyRequestInput collects the field values for NewGetUserPolicyRequest.
type GetUserPolicyRequestInput struct {
	Username   string
	PolicyName string
}

// NewGetUserPolicyRequest builds a ZAP-encoded GetUserPolicyRequest message from
// in and returns the bytes.
func NewGetUserPolicyRequest(in GetUserPolicyRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getUserPolicyRequestSize)
	ob.SetText(getUserPolicyRequestUsernameOff, in.Username)
	ob.SetText(getUserPolicyRequestPolicyNameOff, in.PolicyName)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetUserPolicyResponse ---

const (
	getUserPolicyResponseUsernameOff       = 0
	getUserPolicyResponsePolicyNameOff     = 8
	getUserPolicyResponsePolicyDocumentOff = 16
	getUserPolicyResponseSize              = 24
)

// GetUserPolicyResponse is a zero-copy view into a ZAP-encoded
// GetUserPolicyResponse message.
type GetUserPolicyResponse struct{ o zap.Object }

// WrapGetUserPolicyResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapGetUserPolicyResponse(b []byte) (GetUserPolicyResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetUserPolicyResponse{}, err
	}
	return GetUserPolicyResponse{o: m.Root()}, nil
}

// Username reads the username field (proto field 1, string).
func (t GetUserPolicyResponse) Username() string {
	return t.o.Text(getUserPolicyResponseUsernameOff)
}

// PolicyName reads the policy_name field (proto field 2, string).
func (t GetUserPolicyResponse) PolicyName() string {
	return t.o.Text(getUserPolicyResponsePolicyNameOff)
}

// PolicyDocument reads the policy_document field (proto field 3, string).
func (t GetUserPolicyResponse) PolicyDocument() string {
	return t.o.Text(getUserPolicyResponsePolicyDocumentOff)
}

// GetUserPolicyResponseInput collects the field values for
// NewGetUserPolicyResponse.
type GetUserPolicyResponseInput struct {
	Username       string
	PolicyName     string
	PolicyDocument string
}

// NewGetUserPolicyResponse builds a ZAP-encoded GetUserPolicyResponse message
// from in and returns the bytes.
func NewGetUserPolicyResponse(in GetUserPolicyResponseInput) []byte {
	b := zap.NewBuilder(512)
	ob := b.StartObject(getUserPolicyResponseSize)
	ob.SetText(getUserPolicyResponseUsernameOff, in.Username)
	ob.SetText(getUserPolicyResponsePolicyNameOff, in.PolicyName)
	ob.SetText(getUserPolicyResponsePolicyDocumentOff, in.PolicyDocument)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DeleteUserPolicyRequest ---

const (
	deleteUserPolicyRequestUsernameOff   = 0
	deleteUserPolicyRequestPolicyNameOff = 8
	deleteUserPolicyRequestSize          = 16
)

// DeleteUserPolicyRequest is a zero-copy view into a ZAP-encoded
// DeleteUserPolicyRequest message.
type DeleteUserPolicyRequest struct{ o zap.Object }

// WrapDeleteUserPolicyRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapDeleteUserPolicyRequest(b []byte) (DeleteUserPolicyRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DeleteUserPolicyRequest{}, err
	}
	return DeleteUserPolicyRequest{o: m.Root()}, nil
}

// Username reads the username field (proto field 1, string).
func (t DeleteUserPolicyRequest) Username() string {
	return t.o.Text(deleteUserPolicyRequestUsernameOff)
}

// PolicyName reads the policy_name field (proto field 2, string).
func (t DeleteUserPolicyRequest) PolicyName() string {
	return t.o.Text(deleteUserPolicyRequestPolicyNameOff)
}

// DeleteUserPolicyRequestInput collects the field values for
// NewDeleteUserPolicyRequest.
type DeleteUserPolicyRequestInput struct {
	Username   string
	PolicyName string
}

// NewDeleteUserPolicyRequest builds a ZAP-encoded DeleteUserPolicyRequest
// message from in and returns the bytes.
func NewDeleteUserPolicyRequest(in DeleteUserPolicyRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(deleteUserPolicyRequestSize)
	ob.SetText(deleteUserPolicyRequestUsernameOff, in.Username)
	ob.SetText(deleteUserPolicyRequestPolicyNameOff, in.PolicyName)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DeleteUserPolicyResponse ---

const deleteUserPolicyResponseSize = 0

// DeleteUserPolicyResponse is a zero-copy view into a ZAP-encoded
// DeleteUserPolicyResponse message. It carries no fields.
type DeleteUserPolicyResponse struct{ o zap.Object }

// WrapDeleteUserPolicyResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapDeleteUserPolicyResponse(b []byte) (DeleteUserPolicyResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DeleteUserPolicyResponse{}, err
	}
	return DeleteUserPolicyResponse{o: m.Root()}, nil
}

// DeleteUserPolicyResponseInput collects the field values for
// NewDeleteUserPolicyResponse. DeleteUserPolicyResponse has no fields.
type DeleteUserPolicyResponseInput struct{}

// NewDeleteUserPolicyResponse builds a ZAP-encoded DeleteUserPolicyResponse
// message and returns the bytes.
func NewDeleteUserPolicyResponse(in DeleteUserPolicyResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(deleteUserPolicyResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}
