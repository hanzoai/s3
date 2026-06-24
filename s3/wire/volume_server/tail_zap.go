// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

import (
	zap "github.com/zap-proto/go"
)

// --- VolumeTailSenderRequest ---

const (
	volumeTailSenderRequestVolumeIDOff           = 0
	volumeTailSenderRequestSinceNsOff            = 8
	volumeTailSenderRequestIdleTimeoutSecondsOff = 16
	volumeTailSenderRequestSize                  = 20
)

// VolumeTailSenderRequest is a zero-copy view into a ZAP-encoded
// VolumeTailSenderRequest message.
type VolumeTailSenderRequest struct{ o zap.Object }

// WrapVolumeTailSenderRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeTailSenderRequest(b []byte) (VolumeTailSenderRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeTailSenderRequest{}, err
	}
	return VolumeTailSenderRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeTailSenderRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeTailSenderRequestVolumeIDOff)
}

// SinceNs reads the since_ns field (proto field 2, uint64).
func (t VolumeTailSenderRequest) SinceNs() uint64 {
	return t.o.Uint64(volumeTailSenderRequestSinceNsOff)
}

// IdleTimeoutSeconds reads the idle_timeout_seconds field (proto field 3,
// uint32).
func (t VolumeTailSenderRequest) IdleTimeoutSeconds() uint32 {
	return t.o.Uint32(volumeTailSenderRequestIdleTimeoutSecondsOff)
}

// VolumeTailSenderRequestInput collects the field values for
// NewVolumeTailSenderRequest.
type VolumeTailSenderRequestInput struct {
	VolumeID           uint32
	SinceNs            uint64
	IdleTimeoutSeconds uint32
}

// NewVolumeTailSenderRequest builds a ZAP-encoded VolumeTailSenderRequest
// message from in and returns the bytes.
func NewVolumeTailSenderRequest(in VolumeTailSenderRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeTailSenderRequestSize)
	ob.SetUint32(volumeTailSenderRequestVolumeIDOff, in.VolumeID)
	ob.SetUint64(volumeTailSenderRequestSinceNsOff, in.SinceNs)
	ob.SetUint32(volumeTailSenderRequestIdleTimeoutSecondsOff, in.IdleTimeoutSeconds)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeTailSenderResponse ---

const (
	volumeTailSenderResponseNeedleHeaderOff = 0
	volumeTailSenderResponseNeedleBodyOff   = 8
	volumeTailSenderResponseIsLastChunkOff  = 16
	volumeTailSenderResponseVersionOff      = 20
	volumeTailSenderResponseSize            = 24
)

// VolumeTailSenderResponse is a zero-copy view into a ZAP-encoded
// VolumeTailSenderResponse message.
type VolumeTailSenderResponse struct{ o zap.Object }

// WrapVolumeTailSenderResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeTailSenderResponse(b []byte) (VolumeTailSenderResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeTailSenderResponse{}, err
	}
	return VolumeTailSenderResponse{o: m.Root()}, nil
}

// NeedleHeader reads the needle_header field (proto field 1, bytes).
func (t VolumeTailSenderResponse) NeedleHeader() []byte {
	return t.o.Bytes(volumeTailSenderResponseNeedleHeaderOff)
}

// NeedleBody reads the needle_body field (proto field 2, bytes).
func (t VolumeTailSenderResponse) NeedleBody() []byte {
	return t.o.Bytes(volumeTailSenderResponseNeedleBodyOff)
}

// IsLastChunk reads the is_last_chunk field (proto field 3, bool).
func (t VolumeTailSenderResponse) IsLastChunk() bool {
	return t.o.Bool(volumeTailSenderResponseIsLastChunkOff)
}

// Version reads the version field (proto field 4, uint32).
func (t VolumeTailSenderResponse) Version() uint32 {
	return t.o.Uint32(volumeTailSenderResponseVersionOff)
}

// VolumeTailSenderResponseInput collects the field values for
// NewVolumeTailSenderResponse.
type VolumeTailSenderResponseInput struct {
	NeedleHeader []byte
	NeedleBody   []byte
	IsLastChunk  bool
	Version      uint32
}

// NewVolumeTailSenderResponse builds a ZAP-encoded VolumeTailSenderResponse
// message from in and returns the bytes.
func NewVolumeTailSenderResponse(in VolumeTailSenderResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeTailSenderResponseSize)
	ob.SetBytes(volumeTailSenderResponseNeedleHeaderOff, in.NeedleHeader)
	ob.SetBytes(volumeTailSenderResponseNeedleBodyOff, in.NeedleBody)
	ob.SetBool(volumeTailSenderResponseIsLastChunkOff, in.IsLastChunk)
	ob.SetUint32(volumeTailSenderResponseVersionOff, in.Version)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeTailReceiverRequest ---

const (
	volumeTailReceiverRequestVolumeIDOff           = 0
	volumeTailReceiverRequestSinceNsOff            = 8
	volumeTailReceiverRequestIdleTimeoutSecondsOff = 16
	volumeTailReceiverRequestSourceVolumeServerOff = 24
	volumeTailReceiverRequestSize                  = 32
)

// VolumeTailReceiverRequest is a zero-copy view into a ZAP-encoded
// VolumeTailReceiverRequest message.
type VolumeTailReceiverRequest struct{ o zap.Object }

// WrapVolumeTailReceiverRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeTailReceiverRequest(b []byte) (VolumeTailReceiverRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeTailReceiverRequest{}, err
	}
	return VolumeTailReceiverRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeTailReceiverRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeTailReceiverRequestVolumeIDOff)
}

// SinceNs reads the since_ns field (proto field 2, uint64).
func (t VolumeTailReceiverRequest) SinceNs() uint64 {
	return t.o.Uint64(volumeTailReceiverRequestSinceNsOff)
}

// IdleTimeoutSeconds reads the idle_timeout_seconds field (proto field 3,
// uint32).
func (t VolumeTailReceiverRequest) IdleTimeoutSeconds() uint32 {
	return t.o.Uint32(volumeTailReceiverRequestIdleTimeoutSecondsOff)
}

// SourceVolumeServer reads the source_volume_server field (proto field 4,
// string).
func (t VolumeTailReceiverRequest) SourceVolumeServer() string {
	return t.o.Text(volumeTailReceiverRequestSourceVolumeServerOff)
}

// VolumeTailReceiverRequestInput collects the field values for
// NewVolumeTailReceiverRequest.
type VolumeTailReceiverRequestInput struct {
	VolumeID           uint32
	SinceNs            uint64
	IdleTimeoutSeconds uint32
	SourceVolumeServer string
}

// NewVolumeTailReceiverRequest builds a ZAP-encoded VolumeTailReceiverRequest
// message from in and returns the bytes.
func NewVolumeTailReceiverRequest(in VolumeTailReceiverRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeTailReceiverRequestSize)
	ob.SetUint32(volumeTailReceiverRequestVolumeIDOff, in.VolumeID)
	ob.SetUint64(volumeTailReceiverRequestSinceNsOff, in.SinceNs)
	ob.SetUint32(volumeTailReceiverRequestIdleTimeoutSecondsOff, in.IdleTimeoutSeconds)
	ob.SetText(volumeTailReceiverRequestSourceVolumeServerOff, in.SourceVolumeServer)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeTailReceiverResponse ---

const volumeTailReceiverResponseSize = 0

// VolumeTailReceiverResponse is a zero-copy view into a ZAP-encoded
// VolumeTailReceiverResponse message. VolumeTailReceiverResponse carries no
// fields.
type VolumeTailReceiverResponse struct{ o zap.Object }

// WrapVolumeTailReceiverResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeTailReceiverResponse(b []byte) (VolumeTailReceiverResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeTailReceiverResponse{}, err
	}
	return VolumeTailReceiverResponse{o: m.Root()}, nil
}

// VolumeTailReceiverResponseInput collects the field values for
// NewVolumeTailReceiverResponse. VolumeTailReceiverResponse has no fields.
type VolumeTailReceiverResponseInput struct{}

// NewVolumeTailReceiverResponse builds a ZAP-encoded VolumeTailReceiverResponse
// message and returns the bytes.
func NewVolumeTailReceiverResponse(in VolumeTailReceiverResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeTailReceiverResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}
