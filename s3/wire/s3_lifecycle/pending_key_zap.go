// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	pendingKeyBucketOff     = 0
	pendingKeyRuleHashOff   = 8
	pendingKeyObjectPathOff = 16
	pendingKeyVersionIDOff  = 24
	pendingKeyActionKindOff = 32
	pendingKeySize          = 36
)

// PendingKey is a zero-copy view into a ZAP-encoded PendingKey message.
type PendingKey struct{ o zap.Object }

// WrapPendingKey parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapPendingKey(b []byte) (PendingKey, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PendingKey{}, err
	}
	return PendingKey{o: m.Root()}, nil
}

func (t PendingKey) Bucket() string     { return t.o.Text(pendingKeyBucketOff) }
func (t PendingKey) RuleHash() []byte   { return t.o.Bytes(pendingKeyRuleHashOff) }
func (t PendingKey) ObjectPath() string { return t.o.Text(pendingKeyObjectPathOff) }
func (t PendingKey) VersionID() string  { return t.o.Text(pendingKeyVersionIDOff) }

// ActionKind reads the action_kind field (proto field 5, enum ActionKind as uint32).
func (t PendingKey) ActionKind() uint32 { return t.o.Uint32(pendingKeyActionKindOff) }

// PendingKeyInput collects the field values for NewPendingKey.
type PendingKeyInput struct {
	Bucket     string
	RuleHash   []byte
	ObjectPath string
	VersionID  string
	ActionKind uint32
}

// putPendingKey lays in's fields into ob (a fresh StartObject(pendingKeySize)).
func putPendingKey(ob *zap.ObjectBuilder, in PendingKeyInput) {
	ob.SetText(pendingKeyBucketOff, in.Bucket)
	ob.SetBytes(pendingKeyRuleHashOff, in.RuleHash)
	ob.SetText(pendingKeyObjectPathOff, in.ObjectPath)
	ob.SetText(pendingKeyVersionIDOff, in.VersionID)
	ob.SetUint32(pendingKeyActionKindOff, in.ActionKind)
}

// buildPendingKey lays a PendingKey object into b and returns its offset, for
// embedding as a nested field via ObjectBuilder.SetObject. A nil in yields
// offset 0 (a null nested object).
func buildPendingKey(b *zap.Builder, in *PendingKeyInput) int {
	if in == nil {
		return 0
	}
	ob := b.StartObject(pendingKeySize)
	putPendingKey(ob, *in)
	return ob.Finish()
}

// NewPendingKey builds a ZAP-encoded PendingKey message from in and returns the bytes.
func NewPendingKey(in PendingKeyInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(pendingKeySize)
	putPendingKey(ob, in)
	ob.FinishAsRoot()
	return b.Finish()
}
