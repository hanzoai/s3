// Code generated from iam.proto; DO NOT EDIT.

package iamwire

import (
	zap "github.com/zap-proto/go"
)

// User management request/response messages.

// --- CreateUserRequest ---

const (
	createUserRequestIdentityOff = 0
	createUserRequestSize        = 4
)

// CreateUserRequest is a zero-copy view into a ZAP-encoded CreateUserRequest message.
type CreateUserRequest struct{ o zap.Object }

// WrapCreateUserRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapCreateUserRequest(b []byte) (CreateUserRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CreateUserRequest{}, err
	}
	return CreateUserRequest{o: m.Root()}, nil
}

// Identity reads the identity field (proto field 1, Identity).
func (t CreateUserRequest) Identity() Identity {
	return Identity{o: t.o.Object(createUserRequestIdentityOff)}
}

// CreateUserRequestInput collects the field values for NewCreateUserRequest.
type CreateUserRequestInput struct {
	Identity *IdentityInput
}

// NewCreateUserRequest builds a ZAP-encoded CreateUserRequest message from in
// and returns the bytes.
func NewCreateUserRequest(in CreateUserRequestInput) []byte {
	b := zap.NewBuilder(512)
	identityOff := 0
	if in.Identity != nil {
		identityOff = appendIdentity(b, *in.Identity).Finish()
	}
	ob := b.StartObject(createUserRequestSize)
	if identityOff != 0 {
		ob.SetObject(createUserRequestIdentityOff, identityOff)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- CreateUserResponse ---

const createUserResponseSize = 0

// CreateUserResponse is a zero-copy view into a ZAP-encoded CreateUserResponse
// message. It carries no fields.
type CreateUserResponse struct{ o zap.Object }

// WrapCreateUserResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapCreateUserResponse(b []byte) (CreateUserResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CreateUserResponse{}, err
	}
	return CreateUserResponse{o: m.Root()}, nil
}

// CreateUserResponseInput collects the field values for NewCreateUserResponse.
// CreateUserResponse has no fields.
type CreateUserResponseInput struct{}

// NewCreateUserResponse builds a ZAP-encoded CreateUserResponse message and
// returns the bytes.
func NewCreateUserResponse(in CreateUserResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(createUserResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetUserRequest ---

const (
	getUserRequestUsernameOff = 0
	getUserRequestSize        = 8
)

// GetUserRequest is a zero-copy view into a ZAP-encoded GetUserRequest message.
type GetUserRequest struct{ o zap.Object }

// WrapGetUserRequest parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapGetUserRequest(b []byte) (GetUserRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetUserRequest{}, err
	}
	return GetUserRequest{o: m.Root()}, nil
}

// Username reads the username field (proto field 1, string).
func (t GetUserRequest) Username() string { return t.o.Text(getUserRequestUsernameOff) }

// GetUserRequestInput collects the field values for NewGetUserRequest.
type GetUserRequestInput struct {
	Username string
}

// NewGetUserRequest builds a ZAP-encoded GetUserRequest message from in and
// returns the bytes.
func NewGetUserRequest(in GetUserRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getUserRequestSize)
	ob.SetText(getUserRequestUsernameOff, in.Username)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetUserResponse ---

const (
	getUserResponseIdentityOff = 0
	getUserResponseSize        = 4
)

// GetUserResponse is a zero-copy view into a ZAP-encoded GetUserResponse message.
type GetUserResponse struct{ o zap.Object }

// WrapGetUserResponse parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapGetUserResponse(b []byte) (GetUserResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetUserResponse{}, err
	}
	return GetUserResponse{o: m.Root()}, nil
}

// Identity reads the identity field (proto field 1, Identity).
func (t GetUserResponse) Identity() Identity {
	return Identity{o: t.o.Object(getUserResponseIdentityOff)}
}

// GetUserResponseInput collects the field values for NewGetUserResponse.
type GetUserResponseInput struct {
	Identity *IdentityInput
}

// NewGetUserResponse builds a ZAP-encoded GetUserResponse message from in and
// returns the bytes.
func NewGetUserResponse(in GetUserResponseInput) []byte {
	b := zap.NewBuilder(512)
	identityOff := 0
	if in.Identity != nil {
		identityOff = appendIdentity(b, *in.Identity).Finish()
	}
	ob := b.StartObject(getUserResponseSize)
	if identityOff != 0 {
		ob.SetObject(getUserResponseIdentityOff, identityOff)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- UpdateUserRequest ---

const (
	updateUserRequestUsernameOff = 0
	updateUserRequestIdentityOff = 8
	updateUserRequestSize        = 16
)

// UpdateUserRequest is a zero-copy view into a ZAP-encoded UpdateUserRequest message.
type UpdateUserRequest struct{ o zap.Object }

// WrapUpdateUserRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapUpdateUserRequest(b []byte) (UpdateUserRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return UpdateUserRequest{}, err
	}
	return UpdateUserRequest{o: m.Root()}, nil
}

// Username reads the username field (proto field 1, string).
func (t UpdateUserRequest) Username() string { return t.o.Text(updateUserRequestUsernameOff) }

// Identity reads the identity field (proto field 2, Identity).
func (t UpdateUserRequest) Identity() Identity {
	return Identity{o: t.o.Object(updateUserRequestIdentityOff)}
}

// UpdateUserRequestInput collects the field values for NewUpdateUserRequest.
type UpdateUserRequestInput struct {
	Username string
	Identity *IdentityInput
}

// NewUpdateUserRequest builds a ZAP-encoded UpdateUserRequest message from in
// and returns the bytes.
func NewUpdateUserRequest(in UpdateUserRequestInput) []byte {
	b := zap.NewBuilder(512)
	identityOff := 0
	if in.Identity != nil {
		identityOff = appendIdentity(b, *in.Identity).Finish()
	}
	ob := b.StartObject(updateUserRequestSize)
	ob.SetText(updateUserRequestUsernameOff, in.Username)
	if identityOff != 0 {
		ob.SetObject(updateUserRequestIdentityOff, identityOff)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- UpdateUserResponse ---

const updateUserResponseSize = 0

// UpdateUserResponse is a zero-copy view into a ZAP-encoded UpdateUserResponse
// message. It carries no fields.
type UpdateUserResponse struct{ o zap.Object }

// WrapUpdateUserResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapUpdateUserResponse(b []byte) (UpdateUserResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return UpdateUserResponse{}, err
	}
	return UpdateUserResponse{o: m.Root()}, nil
}

// UpdateUserResponseInput collects the field values for NewUpdateUserResponse.
// UpdateUserResponse has no fields.
type UpdateUserResponseInput struct{}

// NewUpdateUserResponse builds a ZAP-encoded UpdateUserResponse message and
// returns the bytes.
func NewUpdateUserResponse(in UpdateUserResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(updateUserResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DeleteUserRequest ---

const (
	deleteUserRequestUsernameOff = 0
	deleteUserRequestSize        = 8
)

// DeleteUserRequest is a zero-copy view into a ZAP-encoded DeleteUserRequest message.
type DeleteUserRequest struct{ o zap.Object }

// WrapDeleteUserRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapDeleteUserRequest(b []byte) (DeleteUserRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DeleteUserRequest{}, err
	}
	return DeleteUserRequest{o: m.Root()}, nil
}

// Username reads the username field (proto field 1, string).
func (t DeleteUserRequest) Username() string { return t.o.Text(deleteUserRequestUsernameOff) }

// DeleteUserRequestInput collects the field values for NewDeleteUserRequest.
type DeleteUserRequestInput struct {
	Username string
}

// NewDeleteUserRequest builds a ZAP-encoded DeleteUserRequest message from in
// and returns the bytes.
func NewDeleteUserRequest(in DeleteUserRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(deleteUserRequestSize)
	ob.SetText(deleteUserRequestUsernameOff, in.Username)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DeleteUserResponse ---

const deleteUserResponseSize = 0

// DeleteUserResponse is a zero-copy view into a ZAP-encoded DeleteUserResponse
// message. It carries no fields.
type DeleteUserResponse struct{ o zap.Object }

// WrapDeleteUserResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapDeleteUserResponse(b []byte) (DeleteUserResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DeleteUserResponse{}, err
	}
	return DeleteUserResponse{o: m.Root()}, nil
}

// DeleteUserResponseInput collects the field values for NewDeleteUserResponse.
// DeleteUserResponse has no fields.
type DeleteUserResponseInput struct{}

// NewDeleteUserResponse builds a ZAP-encoded DeleteUserResponse message and
// returns the bytes.
func NewDeleteUserResponse(in DeleteUserResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(deleteUserResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListUsersRequest ---

const listUsersRequestSize = 0

// ListUsersRequest is a zero-copy view into a ZAP-encoded ListUsersRequest
// message. It carries no fields.
type ListUsersRequest struct{ o zap.Object }

// WrapListUsersRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapListUsersRequest(b []byte) (ListUsersRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListUsersRequest{}, err
	}
	return ListUsersRequest{o: m.Root()}, nil
}

// ListUsersRequestInput collects the field values for NewListUsersRequest.
// ListUsersRequest has no fields.
type ListUsersRequestInput struct{}

// NewListUsersRequest builds a ZAP-encoded ListUsersRequest message and returns
// the bytes.
func NewListUsersRequest(in ListUsersRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(listUsersRequestSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListUsersResponse ---

const (
	listUsersResponseUsernamesOff = 0
	listUsersResponseSize         = 8
)

// ListUsersResponse is a zero-copy view into a ZAP-encoded ListUsersResponse message.
type ListUsersResponse struct{ o zap.Object }

// WrapListUsersResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapListUsersResponse(b []byte) (ListUsersResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListUsersResponse{}, err
	}
	return ListUsersResponse{o: m.Root()}, nil
}

// UsernamesLen returns the number of usernames (proto field 1, repeated string).
func (t ListUsersResponse) UsernamesLen() int { return t.o.List(listUsersResponseUsernamesOff).Len() }

// UsernamesAt returns username i (proto field 1, repeated string).
func (t ListUsersResponse) UsernamesAt(i int) string {
	return string(t.o.List(listUsersResponseUsernamesOff).BytesAt(i))
}

// ListUsersResponseInput collects the field values for NewListUsersResponse.
type ListUsersResponseInput struct {
	Usernames []string
}

// NewListUsersResponse builds a ZAP-encoded ListUsersResponse message from in
// and returns the bytes.
func NewListUsersResponse(in ListUsersResponseInput) []byte {
	b := zap.NewBuilder(512)
	ob := b.StartObject(listUsersResponseSize)
	if len(in.Usernames) > 0 {
		lb := b.StartList(0)
		for _, s := range in.Usernames {
			lb.AddObjectBytes([]byte(s))
		}
		off, n := lb.Finish()
		ob.SetList(listUsersResponseUsernamesOff, off, n)
	}
	ob.FinishAsRoot()
	return b.Finish()
}
