// Code generated from iam.proto; DO NOT EDIT.

package iamwire

import (
	zap "github.com/zap-proto/go"
)

// S3 IAM cache-management messages (unidirectional propagation from filer to
// S3 servers).

// --- PutIdentityRequest ---

const (
	putIdentityRequestIdentityOff = 0
	putIdentityRequestSize        = 4
)

// PutIdentityRequest is a zero-copy view into a ZAP-encoded PutIdentityRequest message.
type PutIdentityRequest struct{ o zap.Object }

// WrapPutIdentityRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapPutIdentityRequest(b []byte) (PutIdentityRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PutIdentityRequest{}, err
	}
	return PutIdentityRequest{o: m.Root()}, nil
}

// Identity reads the identity field (proto field 1, Identity).
func (t PutIdentityRequest) Identity() Identity {
	return Identity{o: t.o.Object(putIdentityRequestIdentityOff)}
}

// PutIdentityRequestInput collects the field values for NewPutIdentityRequest.
type PutIdentityRequestInput struct {
	Identity *IdentityInput
}

// NewPutIdentityRequest builds a ZAP-encoded PutIdentityRequest message from in
// and returns the bytes.
func NewPutIdentityRequest(in PutIdentityRequestInput) []byte {
	b := zap.NewBuilder(512)
	identityOff := 0
	if in.Identity != nil {
		identityOff = appendIdentity(b, *in.Identity).Finish()
	}
	ob := b.StartObject(putIdentityRequestSize)
	if identityOff != 0 {
		ob.SetObject(putIdentityRequestIdentityOff, identityOff)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PutIdentityResponse ---

const putIdentityResponseSize = 0

// PutIdentityResponse is a zero-copy view into a ZAP-encoded PutIdentityResponse
// message. It carries no fields.
type PutIdentityResponse struct{ o zap.Object }

// WrapPutIdentityResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapPutIdentityResponse(b []byte) (PutIdentityResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PutIdentityResponse{}, err
	}
	return PutIdentityResponse{o: m.Root()}, nil
}

// PutIdentityResponseInput collects the field values for NewPutIdentityResponse.
// PutIdentityResponse has no fields.
type PutIdentityResponseInput struct{}

// NewPutIdentityResponse builds a ZAP-encoded PutIdentityResponse message and
// returns the bytes.
func NewPutIdentityResponse(in PutIdentityResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(putIdentityResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RemoveIdentityRequest ---

const (
	removeIdentityRequestUsernameOff = 0
	removeIdentityRequestSize        = 8
)

// RemoveIdentityRequest is a zero-copy view into a ZAP-encoded
// RemoveIdentityRequest message.
type RemoveIdentityRequest struct{ o zap.Object }

// WrapRemoveIdentityRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapRemoveIdentityRequest(b []byte) (RemoveIdentityRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RemoveIdentityRequest{}, err
	}
	return RemoveIdentityRequest{o: m.Root()}, nil
}

// Username reads the username field (proto field 1, string).
func (t RemoveIdentityRequest) Username() string {
	return t.o.Text(removeIdentityRequestUsernameOff)
}

// RemoveIdentityRequestInput collects the field values for
// NewRemoveIdentityRequest.
type RemoveIdentityRequestInput struct {
	Username string
}

// NewRemoveIdentityRequest builds a ZAP-encoded RemoveIdentityRequest message
// from in and returns the bytes.
func NewRemoveIdentityRequest(in RemoveIdentityRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(removeIdentityRequestSize)
	ob.SetText(removeIdentityRequestUsernameOff, in.Username)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RemoveIdentityResponse ---

const removeIdentityResponseSize = 0

// RemoveIdentityResponse is a zero-copy view into a ZAP-encoded
// RemoveIdentityResponse message. It carries no fields.
type RemoveIdentityResponse struct{ o zap.Object }

// WrapRemoveIdentityResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapRemoveIdentityResponse(b []byte) (RemoveIdentityResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RemoveIdentityResponse{}, err
	}
	return RemoveIdentityResponse{o: m.Root()}, nil
}

// RemoveIdentityResponseInput collects the field values for
// NewRemoveIdentityResponse. RemoveIdentityResponse has no fields.
type RemoveIdentityResponseInput struct{}

// NewRemoveIdentityResponse builds a ZAP-encoded RemoveIdentityResponse message
// and returns the bytes.
func NewRemoveIdentityResponse(in RemoveIdentityResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(removeIdentityResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PutGroupRequest ---

const (
	putGroupRequestGroupOff = 0
	putGroupRequestSize     = 4
)

// PutGroupRequest is a zero-copy view into a ZAP-encoded PutGroupRequest message.
type PutGroupRequest struct{ o zap.Object }

// WrapPutGroupRequest parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapPutGroupRequest(b []byte) (PutGroupRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PutGroupRequest{}, err
	}
	return PutGroupRequest{o: m.Root()}, nil
}

// Group reads the group field (proto field 1, Group).
func (t PutGroupRequest) Group() Group {
	return Group{o: t.o.Object(putGroupRequestGroupOff)}
}

// PutGroupRequestInput collects the field values for NewPutGroupRequest.
type PutGroupRequestInput struct {
	Group *GroupInput
}

// NewPutGroupRequest builds a ZAP-encoded PutGroupRequest message from in and
// returns the bytes.
func NewPutGroupRequest(in PutGroupRequestInput) []byte {
	b := zap.NewBuilder(512)
	groupOff := 0
	if in.Group != nil {
		groupOff = appendGroup(b, *in.Group).Finish()
	}
	ob := b.StartObject(putGroupRequestSize)
	if groupOff != 0 {
		ob.SetObject(putGroupRequestGroupOff, groupOff)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PutGroupResponse ---

const putGroupResponseSize = 0

// PutGroupResponse is a zero-copy view into a ZAP-encoded PutGroupResponse
// message. It carries no fields.
type PutGroupResponse struct{ o zap.Object }

// WrapPutGroupResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapPutGroupResponse(b []byte) (PutGroupResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PutGroupResponse{}, err
	}
	return PutGroupResponse{o: m.Root()}, nil
}

// PutGroupResponseInput collects the field values for NewPutGroupResponse.
// PutGroupResponse has no fields.
type PutGroupResponseInput struct{}

// NewPutGroupResponse builds a ZAP-encoded PutGroupResponse message and returns
// the bytes.
func NewPutGroupResponse(in PutGroupResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(putGroupResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RemoveGroupRequest ---

const (
	removeGroupRequestGroupNameOff = 0
	removeGroupRequestSize         = 8
)

// RemoveGroupRequest is a zero-copy view into a ZAP-encoded RemoveGroupRequest message.
type RemoveGroupRequest struct{ o zap.Object }

// WrapRemoveGroupRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapRemoveGroupRequest(b []byte) (RemoveGroupRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RemoveGroupRequest{}, err
	}
	return RemoveGroupRequest{o: m.Root()}, nil
}

// GroupName reads the group_name field (proto field 1, string).
func (t RemoveGroupRequest) GroupName() string { return t.o.Text(removeGroupRequestGroupNameOff) }

// RemoveGroupRequestInput collects the field values for NewRemoveGroupRequest.
type RemoveGroupRequestInput struct {
	GroupName string
}

// NewRemoveGroupRequest builds a ZAP-encoded RemoveGroupRequest message from in
// and returns the bytes.
func NewRemoveGroupRequest(in RemoveGroupRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(removeGroupRequestSize)
	ob.SetText(removeGroupRequestGroupNameOff, in.GroupName)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RemoveGroupResponse ---

const removeGroupResponseSize = 0

// RemoveGroupResponse is a zero-copy view into a ZAP-encoded RemoveGroupResponse
// message. It carries no fields.
type RemoveGroupResponse struct{ o zap.Object }

// WrapRemoveGroupResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapRemoveGroupResponse(b []byte) (RemoveGroupResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RemoveGroupResponse{}, err
	}
	return RemoveGroupResponse{o: m.Root()}, nil
}

// RemoveGroupResponseInput collects the field values for NewRemoveGroupResponse.
// RemoveGroupResponse has no fields.
type RemoveGroupResponseInput struct{}

// NewRemoveGroupResponse builds a ZAP-encoded RemoveGroupResponse message and
// returns the bytes.
func NewRemoveGroupResponse(in RemoveGroupResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(removeGroupResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}
