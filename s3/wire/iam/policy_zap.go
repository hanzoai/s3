// Code generated from iam.proto; DO NOT EDIT.

package iamwire

import (
	zap "github.com/zap-proto/go"
)

// Managed-policy management request/response messages.

// --- PutPolicyRequest ---

const (
	putPolicyRequestNameOff    = 0
	putPolicyRequestContentOff = 8
	putPolicyRequestSize       = 16
)

// PutPolicyRequest is a zero-copy view into a ZAP-encoded PutPolicyRequest message.
type PutPolicyRequest struct{ o zap.Object }

// WrapPutPolicyRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapPutPolicyRequest(b []byte) (PutPolicyRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PutPolicyRequest{}, err
	}
	return PutPolicyRequest{o: m.Root()}, nil
}

// Name reads the name field (proto field 1, string).
func (t PutPolicyRequest) Name() string { return t.o.Text(putPolicyRequestNameOff) }

// Content reads the content field (proto field 2, string).
func (t PutPolicyRequest) Content() string { return t.o.Text(putPolicyRequestContentOff) }

// PutPolicyRequestInput collects the field values for NewPutPolicyRequest.
type PutPolicyRequestInput struct {
	Name    string
	Content string
}

// NewPutPolicyRequest builds a ZAP-encoded PutPolicyRequest message from in and
// returns the bytes.
func NewPutPolicyRequest(in PutPolicyRequestInput) []byte {
	b := zap.NewBuilder(512)
	ob := b.StartObject(putPolicyRequestSize)
	ob.SetText(putPolicyRequestNameOff, in.Name)
	ob.SetText(putPolicyRequestContentOff, in.Content)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PutPolicyResponse ---

const putPolicyResponseSize = 0

// PutPolicyResponse is a zero-copy view into a ZAP-encoded PutPolicyResponse
// message. It carries no fields.
type PutPolicyResponse struct{ o zap.Object }

// WrapPutPolicyResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapPutPolicyResponse(b []byte) (PutPolicyResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PutPolicyResponse{}, err
	}
	return PutPolicyResponse{o: m.Root()}, nil
}

// PutPolicyResponseInput collects the field values for NewPutPolicyResponse.
// PutPolicyResponse has no fields.
type PutPolicyResponseInput struct{}

// NewPutPolicyResponse builds a ZAP-encoded PutPolicyResponse message and
// returns the bytes.
func NewPutPolicyResponse(in PutPolicyResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(putPolicyResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetPolicyRequest ---

const (
	getPolicyRequestNameOff = 0
	getPolicyRequestSize    = 8
)

// GetPolicyRequest is a zero-copy view into a ZAP-encoded GetPolicyRequest message.
type GetPolicyRequest struct{ o zap.Object }

// WrapGetPolicyRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapGetPolicyRequest(b []byte) (GetPolicyRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetPolicyRequest{}, err
	}
	return GetPolicyRequest{o: m.Root()}, nil
}

// Name reads the name field (proto field 1, string).
func (t GetPolicyRequest) Name() string { return t.o.Text(getPolicyRequestNameOff) }

// GetPolicyRequestInput collects the field values for NewGetPolicyRequest.
type GetPolicyRequestInput struct {
	Name string
}

// NewGetPolicyRequest builds a ZAP-encoded GetPolicyRequest message from in and
// returns the bytes.
func NewGetPolicyRequest(in GetPolicyRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getPolicyRequestSize)
	ob.SetText(getPolicyRequestNameOff, in.Name)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetPolicyResponse ---

const (
	getPolicyResponseNameOff    = 0
	getPolicyResponseContentOff = 8
	getPolicyResponseSize       = 16
)

// GetPolicyResponse is a zero-copy view into a ZAP-encoded GetPolicyResponse message.
type GetPolicyResponse struct{ o zap.Object }

// WrapGetPolicyResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapGetPolicyResponse(b []byte) (GetPolicyResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetPolicyResponse{}, err
	}
	return GetPolicyResponse{o: m.Root()}, nil
}

// Name reads the name field (proto field 1, string).
func (t GetPolicyResponse) Name() string { return t.o.Text(getPolicyResponseNameOff) }

// Content reads the content field (proto field 2, string).
func (t GetPolicyResponse) Content() string { return t.o.Text(getPolicyResponseContentOff) }

// GetPolicyResponseInput collects the field values for NewGetPolicyResponse.
type GetPolicyResponseInput struct {
	Name    string
	Content string
}

// NewGetPolicyResponse builds a ZAP-encoded GetPolicyResponse message from in
// and returns the bytes.
func NewGetPolicyResponse(in GetPolicyResponseInput) []byte {
	b := zap.NewBuilder(512)
	ob := b.StartObject(getPolicyResponseSize)
	ob.SetText(getPolicyResponseNameOff, in.Name)
	ob.SetText(getPolicyResponseContentOff, in.Content)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListPoliciesRequest ---

const listPoliciesRequestSize = 0

// ListPoliciesRequest is a zero-copy view into a ZAP-encoded ListPoliciesRequest
// message. It carries no fields.
type ListPoliciesRequest struct{ o zap.Object }

// WrapListPoliciesRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapListPoliciesRequest(b []byte) (ListPoliciesRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListPoliciesRequest{}, err
	}
	return ListPoliciesRequest{o: m.Root()}, nil
}

// ListPoliciesRequestInput collects the field values for NewListPoliciesRequest.
// ListPoliciesRequest has no fields.
type ListPoliciesRequestInput struct{}

// NewListPoliciesRequest builds a ZAP-encoded ListPoliciesRequest message and
// returns the bytes.
func NewListPoliciesRequest(in ListPoliciesRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(listPoliciesRequestSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListPoliciesResponse ---

const (
	listPoliciesResponsePoliciesOff = 0
	listPoliciesResponseSize        = 8
)

// ListPoliciesResponse is a zero-copy view into a ZAP-encoded
// ListPoliciesResponse message.
type ListPoliciesResponse struct{ o zap.Object }

// WrapListPoliciesResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapListPoliciesResponse(b []byte) (ListPoliciesResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListPoliciesResponse{}, err
	}
	return ListPoliciesResponse{o: m.Root()}, nil
}

// PoliciesLen returns the number of policies (proto field 1, repeated Policy).
func (t ListPoliciesResponse) PoliciesLen() int {
	return t.o.List(listPoliciesResponsePoliciesOff).Len()
}

// PoliciesAt returns policy i (proto field 1, repeated Policy).
func (t ListPoliciesResponse) PoliciesAt(i int) Policy {
	return Policy{o: t.o.List(listPoliciesResponsePoliciesOff).ObjectAt(i)}
}

// ListPoliciesResponseInput collects the field values for NewListPoliciesResponse.
type ListPoliciesResponseInput struct {
	Policies []PolicyInput
}

// NewListPoliciesResponse builds a ZAP-encoded ListPoliciesResponse message from
// in and returns the bytes.
func NewListPoliciesResponse(in ListPoliciesResponseInput) []byte {
	b := zap.NewBuilder(512)
	ob := b.StartObject(listPoliciesResponseSize)
	if len(in.Policies) > 0 {
		lb := b.StartList(0)
		for i := range in.Policies {
			lb.AddObjectBytes(NewPolicy(in.Policies[i]))
		}
		off, n := lb.Finish()
		ob.SetList(listPoliciesResponsePoliciesOff, off, n)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DeletePolicyRequest ---

const (
	deletePolicyRequestNameOff = 0
	deletePolicyRequestSize    = 8
)

// DeletePolicyRequest is a zero-copy view into a ZAP-encoded DeletePolicyRequest message.
type DeletePolicyRequest struct{ o zap.Object }

// WrapDeletePolicyRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapDeletePolicyRequest(b []byte) (DeletePolicyRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DeletePolicyRequest{}, err
	}
	return DeletePolicyRequest{o: m.Root()}, nil
}

// Name reads the name field (proto field 1, string).
func (t DeletePolicyRequest) Name() string { return t.o.Text(deletePolicyRequestNameOff) }

// DeletePolicyRequestInput collects the field values for NewDeletePolicyRequest.
type DeletePolicyRequestInput struct {
	Name string
}

// NewDeletePolicyRequest builds a ZAP-encoded DeletePolicyRequest message from
// in and returns the bytes.
func NewDeletePolicyRequest(in DeletePolicyRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(deletePolicyRequestSize)
	ob.SetText(deletePolicyRequestNameOff, in.Name)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DeletePolicyResponse ---

const deletePolicyResponseSize = 0

// DeletePolicyResponse is a zero-copy view into a ZAP-encoded DeletePolicyResponse
// message. It carries no fields.
type DeletePolicyResponse struct{ o zap.Object }

// WrapDeletePolicyResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapDeletePolicyResponse(b []byte) (DeletePolicyResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DeletePolicyResponse{}, err
	}
	return DeletePolicyResponse{o: m.Root()}, nil
}

// DeletePolicyResponseInput collects the field values for NewDeletePolicyResponse.
// DeletePolicyResponse has no fields.
type DeletePolicyResponseInput struct{}

// NewDeletePolicyResponse builds a ZAP-encoded DeletePolicyResponse message and
// returns the bytes.
func NewDeletePolicyResponse(in DeletePolicyResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(deletePolicyResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}
