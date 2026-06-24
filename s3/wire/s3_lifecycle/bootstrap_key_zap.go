// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	bootstrapKeyBucketOff     = 0
	bootstrapKeyObjectPathOff = 8
	bootstrapKeyVersionIDOff  = 16
	bootstrapKeyRuleHashOff   = 24
	bootstrapKeyActionKindOff = 32
	bootstrapKeySize          = 36
)

// BootstrapKey is a zero-copy view into a ZAP-encoded BootstrapKey message.
type BootstrapKey struct{ o zap.Object }

// WrapBootstrapKey parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapBootstrapKey(b []byte) (BootstrapKey, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return BootstrapKey{}, err
	}
	return BootstrapKey{o: m.Root()}, nil
}

func (t BootstrapKey) Bucket() string     { return t.o.Text(bootstrapKeyBucketOff) }
func (t BootstrapKey) ObjectPath() string { return t.o.Text(bootstrapKeyObjectPathOff) }
func (t BootstrapKey) VersionID() string  { return t.o.Text(bootstrapKeyVersionIDOff) }
func (t BootstrapKey) RuleHash() []byte   { return t.o.Bytes(bootstrapKeyRuleHashOff) }

// ActionKind reads the action_kind field (proto field 5, enum ActionKind as uint32).
func (t BootstrapKey) ActionKind() uint32 { return t.o.Uint32(bootstrapKeyActionKindOff) }

// BootstrapKeyInput collects the field values for NewBootstrapKey.
type BootstrapKeyInput struct {
	Bucket     string
	ObjectPath string
	VersionID  string
	RuleHash   []byte
	ActionKind uint32
}

// putBootstrapKey lays in's fields into ob (a fresh StartObject(bootstrapKeySize)).
func putBootstrapKey(ob *zap.ObjectBuilder, in BootstrapKeyInput) {
	ob.SetText(bootstrapKeyBucketOff, in.Bucket)
	ob.SetText(bootstrapKeyObjectPathOff, in.ObjectPath)
	ob.SetText(bootstrapKeyVersionIDOff, in.VersionID)
	ob.SetBytes(bootstrapKeyRuleHashOff, in.RuleHash)
	ob.SetUint32(bootstrapKeyActionKindOff, in.ActionKind)
}

// buildBootstrapKey lays a BootstrapKey object into b and returns its offset, for
// embedding as a nested field via ObjectBuilder.SetObject. A nil in yields
// offset 0 (a null nested object).
func buildBootstrapKey(b *zap.Builder, in *BootstrapKeyInput) int {
	if in == nil {
		return 0
	}
	ob := b.StartObject(bootstrapKeySize)
	putBootstrapKey(ob, *in)
	return ob.Finish()
}

// NewBootstrapKey builds a ZAP-encoded BootstrapKey message from in and returns the bytes.
func NewBootstrapKey(in BootstrapKeyInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(bootstrapKeySize)
	putBootstrapKey(ob, in)
	ob.FinishAsRoot()
	return b.Finish()
}
