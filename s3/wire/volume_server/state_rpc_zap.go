// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

import (
	zap "github.com/zap-proto/go"
)

// --- GetStateRequest ---

const getStateRequestSize = 0

// GetStateRequest is a zero-copy view into a ZAP-encoded GetStateRequest
// message. GetStateRequest carries no fields.
type GetStateRequest struct{ o zap.Object }

// WrapGetStateRequest parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapGetStateRequest(b []byte) (GetStateRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetStateRequest{}, err
	}
	return GetStateRequest{o: m.Root()}, nil
}

// GetStateRequestInput collects the field values for NewGetStateRequest.
// GetStateRequest has no fields.
type GetStateRequestInput struct{}

// NewGetStateRequest builds a ZAP-encoded GetStateRequest message and returns
// the bytes.
func NewGetStateRequest(in GetStateRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getStateRequestSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetStateResponse ---

const (
	getStateResponseStateOff = 0
	getStateResponseSize     = 4
)

// GetStateResponse is a zero-copy view into a ZAP-encoded GetStateResponse
// message.
type GetStateResponse struct{ o zap.Object }

// WrapGetStateResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapGetStateResponse(b []byte) (GetStateResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetStateResponse{}, err
	}
	return GetStateResponse{o: m.Root()}, nil
}

// State reads the state field (proto field 1, message VolumeServerState). The
// bool is false when the field is absent.
func (t GetStateResponse) State() (VolumeServerState, bool) {
	o := t.o.Object(getStateResponseStateOff)
	if o.IsNull() {
		return VolumeServerState{}, false
	}
	return VolumeServerState{o: o}, true
}

// GetStateResponseInput collects the field values for NewGetStateResponse. State
// is the VolumeServerState message's own ZAP buffer (from NewVolumeServerState).
type GetStateResponseInput struct {
	State []byte
}

// NewGetStateResponse builds a ZAP-encoded GetStateResponse message from in and
// returns the bytes.
func NewGetStateResponse(in GetStateResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getStateResponseSize)
	ob.SetObject(getStateResponseStateOff, embedMessage(b, in.State))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SetStateRequest ---

const (
	setStateRequestStateOff = 0
	setStateRequestSize     = 4
)

// SetStateRequest is a zero-copy view into a ZAP-encoded SetStateRequest
// message.
type SetStateRequest struct{ o zap.Object }

// WrapSetStateRequest parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapSetStateRequest(b []byte) (SetStateRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SetStateRequest{}, err
	}
	return SetStateRequest{o: m.Root()}, nil
}

// State reads the state field (proto field 1, message VolumeServerState). The
// bool is false when the field is absent.
func (t SetStateRequest) State() (VolumeServerState, bool) {
	o := t.o.Object(setStateRequestStateOff)
	if o.IsNull() {
		return VolumeServerState{}, false
	}
	return VolumeServerState{o: o}, true
}

// SetStateRequestInput collects the field values for NewSetStateRequest. State
// is the VolumeServerState message's own ZAP buffer (from NewVolumeServerState).
type SetStateRequestInput struct {
	State []byte
}

// NewSetStateRequest builds a ZAP-encoded SetStateRequest message from in and
// returns the bytes.
func NewSetStateRequest(in SetStateRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(setStateRequestSize)
	ob.SetObject(setStateRequestStateOff, embedMessage(b, in.State))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SetStateResponse ---

const (
	setStateResponseStateOff = 0
	setStateResponseSize     = 4
)

// SetStateResponse is a zero-copy view into a ZAP-encoded SetStateResponse
// message.
type SetStateResponse struct{ o zap.Object }

// WrapSetStateResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapSetStateResponse(b []byte) (SetStateResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SetStateResponse{}, err
	}
	return SetStateResponse{o: m.Root()}, nil
}

// State reads the state field (proto field 1, message VolumeServerState). The
// bool is false when the field is absent.
func (t SetStateResponse) State() (VolumeServerState, bool) {
	o := t.o.Object(setStateResponseStateOff)
	if o.IsNull() {
		return VolumeServerState{}, false
	}
	return VolumeServerState{o: o}, true
}

// SetStateResponseInput collects the field values for NewSetStateResponse. State
// is the VolumeServerState message's own ZAP buffer (from NewVolumeServerState).
type SetStateResponseInput struct {
	State []byte
}

// NewSetStateResponse builds a ZAP-encoded SetStateResponse message from in and
// returns the bytes.
func NewSetStateResponse(in SetStateResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(setStateResponseSize)
	ob.SetObject(setStateResponseStateOff, embedMessage(b, in.State))
	ob.FinishAsRoot()
	return b.Finish()
}
