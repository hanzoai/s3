// Code generated from mount.proto; DO NOT EDIT.

package mountwire

import (
	zap "github.com/zap-proto/go"
)

// --- ConfigureRequest ---

const (
	configureRequestCollectionCapacityOff = 0
	configureRequestSize                  = 8
)

// ConfigureRequest is a zero-copy view into a ZAP-encoded ConfigureRequest message.
type ConfigureRequest struct{ o zap.Object }

// WrapConfigureRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapConfigureRequest(b []byte) (ConfigureRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ConfigureRequest{}, err
	}
	return ConfigureRequest{o: m.Root()}, nil
}

// CollectionCapacity reads the collection_capacity field (proto field 1, int64).
func (t ConfigureRequest) CollectionCapacity() int64 {
	return t.o.Int64(configureRequestCollectionCapacityOff)
}

// ConfigureRequestInput collects the field values for NewConfigureRequest.
type ConfigureRequestInput struct {
	CollectionCapacity int64
}

// NewConfigureRequest builds a ZAP-encoded ConfigureRequest message from in and
// returns the bytes.
func NewConfigureRequest(in ConfigureRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(configureRequestSize)
	ob.SetInt64(configureRequestCollectionCapacityOff, in.CollectionCapacity)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ConfigureResponse ---

const configureResponseSize = 0

// ConfigureResponse is a zero-copy view into a ZAP-encoded ConfigureResponse
// message. ConfigureResponse carries no fields.
type ConfigureResponse struct{ o zap.Object }

// WrapConfigureResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapConfigureResponse(b []byte) (ConfigureResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ConfigureResponse{}, err
	}
	return ConfigureResponse{o: m.Root()}, nil
}

// ConfigureResponseInput collects the field values for NewConfigureResponse.
// ConfigureResponse has no fields.
type ConfigureResponseInput struct{}

// NewConfigureResponse builds a ZAP-encoded ConfigureResponse message and
// returns the bytes.
func NewConfigureResponse(in ConfigureResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(configureResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}
