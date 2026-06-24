// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	lifecycleDeleteResponseOutcomeOff = 0
	lifecycleDeleteResponseReasonOff  = 4
	lifecycleDeleteResponseSize       = 12
)

// LifecycleDeleteResponse is a zero-copy view into a ZAP-encoded LifecycleDeleteResponse message.
type LifecycleDeleteResponse struct{ o zap.Object }

// WrapLifecycleDeleteResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapLifecycleDeleteResponse(b []byte) (LifecycleDeleteResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LifecycleDeleteResponse{}, err
	}
	return LifecycleDeleteResponse{o: m.Root()}, nil
}

// Outcome reads the outcome field (proto field 1, enum LifecycleDeleteOutcome as uint32).
func (t LifecycleDeleteResponse) Outcome() uint32 {
	return t.o.Uint32(lifecycleDeleteResponseOutcomeOff)
}
func (t LifecycleDeleteResponse) Reason() string { return t.o.Text(lifecycleDeleteResponseReasonOff) }

// LifecycleDeleteResponseInput collects the field values for NewLifecycleDeleteResponse.
type LifecycleDeleteResponseInput struct {
	Outcome uint32
	Reason  string
}

// NewLifecycleDeleteResponse builds a ZAP-encoded LifecycleDeleteResponse message from in and returns the bytes.
func NewLifecycleDeleteResponse(in LifecycleDeleteResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(lifecycleDeleteResponseSize)
	ob.SetUint32(lifecycleDeleteResponseOutcomeOff, in.Outcome)
	ob.SetText(lifecycleDeleteResponseReasonOff, in.Reason)
	ob.FinishAsRoot()
	return b.Finish()
}
