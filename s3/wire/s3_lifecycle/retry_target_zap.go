// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	retryTargetKeyOff        = 0
	retryTargetBucketOff     = 4
	retryTargetObjectPathOff = 12
	retryTargetVersionIDOff  = 20
	retryTargetRuleHashOff   = 28
	retryTargetActionKindOff = 36
	retryTargetSize          = 40
)

// RetryTarget is a zero-copy view into a ZAP-encoded RetryTarget message.
type RetryTarget struct{ o zap.Object }

// WrapRetryTarget parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapRetryTarget(b []byte) (RetryTarget, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RetryTarget{}, err
	}
	return RetryTarget{o: m.Root()}, nil
}

// Key reads the key field (proto field 1, message StreamKey).
func (t RetryTarget) Key() StreamKey {
	return StreamKey{o: t.o.Object(retryTargetKeyOff)}
}
func (t RetryTarget) Bucket() string     { return t.o.Text(retryTargetBucketOff) }
func (t RetryTarget) ObjectPath() string { return t.o.Text(retryTargetObjectPathOff) }
func (t RetryTarget) VersionID() string  { return t.o.Text(retryTargetVersionIDOff) }
func (t RetryTarget) RuleHash() []byte   { return t.o.Bytes(retryTargetRuleHashOff) }

// ActionKind reads the action_kind field (proto field 6, enum ActionKind as uint32).
func (t RetryTarget) ActionKind() uint32 { return t.o.Uint32(retryTargetActionKindOff) }

// RetryTargetInput collects the field values for NewRetryTarget.
type RetryTargetInput struct {
	Key        *StreamKeyInput
	Bucket     string
	ObjectPath string
	VersionID  string
	RuleHash   []byte
	ActionKind uint32
}

// NewRetryTarget builds a ZAP-encoded RetryTarget message from in and returns the bytes.
func NewRetryTarget(in RetryTargetInput) []byte {
	b := zap.NewBuilder(256)
	key := buildStreamKey(b, in.Key)
	ob := b.StartObject(retryTargetSize)
	ob.SetObject(retryTargetKeyOff, key)
	ob.SetText(retryTargetBucketOff, in.Bucket)
	ob.SetText(retryTargetObjectPathOff, in.ObjectPath)
	ob.SetText(retryTargetVersionIDOff, in.VersionID)
	ob.SetBytes(retryTargetRuleHashOff, in.RuleHash)
	ob.SetUint32(retryTargetActionKindOff, in.ActionKind)
	ob.FinishAsRoot()
	return b.Finish()
}
