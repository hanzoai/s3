// Code generated from iam.proto; DO NOT EDIT.

package iamwire

import (
	zap "github.com/zap-proto/go"
)

// Configuration management request/response messages.

// --- GetConfigurationRequest ---

const getConfigurationRequestSize = 0

// GetConfigurationRequest is a zero-copy view into a ZAP-encoded
// GetConfigurationRequest message. It carries no fields.
type GetConfigurationRequest struct{ o zap.Object }

// WrapGetConfigurationRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapGetConfigurationRequest(b []byte) (GetConfigurationRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetConfigurationRequest{}, err
	}
	return GetConfigurationRequest{o: m.Root()}, nil
}

// GetConfigurationRequestInput collects the field values for
// NewGetConfigurationRequest. GetConfigurationRequest has no fields.
type GetConfigurationRequestInput struct{}

// NewGetConfigurationRequest builds a ZAP-encoded GetConfigurationRequest
// message and returns the bytes.
func NewGetConfigurationRequest(in GetConfigurationRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getConfigurationRequestSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetConfigurationResponse ---

const (
	getConfigurationResponseConfigurationOff = 0
	getConfigurationResponseSize             = 4
)

// GetConfigurationResponse is a zero-copy view into a ZAP-encoded
// GetConfigurationResponse message.
type GetConfigurationResponse struct{ o zap.Object }

// WrapGetConfigurationResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapGetConfigurationResponse(b []byte) (GetConfigurationResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetConfigurationResponse{}, err
	}
	return GetConfigurationResponse{o: m.Root()}, nil
}

// Configuration reads the configuration field (proto field 1, S3ApiConfiguration).
func (t GetConfigurationResponse) Configuration() S3ApiConfiguration {
	return S3ApiConfiguration{o: t.o.Object(getConfigurationResponseConfigurationOff)}
}

// GetConfigurationResponseInput collects the field values for
// NewGetConfigurationResponse.
type GetConfigurationResponseInput struct {
	Configuration *S3ApiConfigurationInput
}

// NewGetConfigurationResponse builds a ZAP-encoded GetConfigurationResponse
// message from in and returns the bytes.
func NewGetConfigurationResponse(in GetConfigurationResponseInput) []byte {
	b := zap.NewBuilder(1024)
	configurationOff := 0
	if in.Configuration != nil {
		configurationOff = appendS3ApiConfiguration(b, *in.Configuration).Finish()
	}
	ob := b.StartObject(getConfigurationResponseSize)
	if configurationOff != 0 {
		ob.SetObject(getConfigurationResponseConfigurationOff, configurationOff)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PutConfigurationRequest ---

const (
	putConfigurationRequestConfigurationOff = 0
	putConfigurationRequestSize             = 4
)

// PutConfigurationRequest is a zero-copy view into a ZAP-encoded
// PutConfigurationRequest message.
type PutConfigurationRequest struct{ o zap.Object }

// WrapPutConfigurationRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapPutConfigurationRequest(b []byte) (PutConfigurationRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PutConfigurationRequest{}, err
	}
	return PutConfigurationRequest{o: m.Root()}, nil
}

// Configuration reads the configuration field (proto field 1, S3ApiConfiguration).
func (t PutConfigurationRequest) Configuration() S3ApiConfiguration {
	return S3ApiConfiguration{o: t.o.Object(putConfigurationRequestConfigurationOff)}
}

// PutConfigurationRequestInput collects the field values for
// NewPutConfigurationRequest.
type PutConfigurationRequestInput struct {
	Configuration *S3ApiConfigurationInput
}

// NewPutConfigurationRequest builds a ZAP-encoded PutConfigurationRequest
// message from in and returns the bytes.
func NewPutConfigurationRequest(in PutConfigurationRequestInput) []byte {
	b := zap.NewBuilder(1024)
	configurationOff := 0
	if in.Configuration != nil {
		configurationOff = appendS3ApiConfiguration(b, *in.Configuration).Finish()
	}
	ob := b.StartObject(putConfigurationRequestSize)
	if configurationOff != 0 {
		ob.SetObject(putConfigurationRequestConfigurationOff, configurationOff)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PutConfigurationResponse ---

const putConfigurationResponseSize = 0

// PutConfigurationResponse is a zero-copy view into a ZAP-encoded
// PutConfigurationResponse message. It carries no fields.
type PutConfigurationResponse struct{ o zap.Object }

// WrapPutConfigurationResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapPutConfigurationResponse(b []byte) (PutConfigurationResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PutConfigurationResponse{}, err
	}
	return PutConfigurationResponse{o: m.Root()}, nil
}

// PutConfigurationResponseInput collects the field values for
// NewPutConfigurationResponse. PutConfigurationResponse has no fields.
type PutConfigurationResponseInput struct{}

// NewPutConfigurationResponse builds a ZAP-encoded PutConfigurationResponse
// message and returns the bytes.
func NewPutConfigurationResponse(in PutConfigurationResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(putConfigurationResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}
