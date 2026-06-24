// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	tailDrainedKeyStreamKindOff   = 0
	tailDrainedKeyFilerIDOff      = 4
	tailDrainedKeyDelaySecondsOff = 12
	tailDrainedKeySize            = 20
)

// TailDrainedKey is a zero-copy view into a ZAP-encoded TailDrainedKey message.
type TailDrainedKey struct{ o zap.Object }

// WrapTailDrainedKey parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapTailDrainedKey(b []byte) (TailDrainedKey, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TailDrainedKey{}, err
	}
	return TailDrainedKey{o: m.Root()}, nil
}

// StreamKind reads the stream_kind field (proto field 1, enum StreamKind as uint32).
func (t TailDrainedKey) StreamKind() uint32 { return t.o.Uint32(tailDrainedKeyStreamKindOff) }
func (t TailDrainedKey) FilerID() string    { return t.o.Text(tailDrainedKeyFilerIDOff) }
func (t TailDrainedKey) DelaySeconds() int64 {
	return t.o.Int64(tailDrainedKeyDelaySecondsOff)
}

// TailDrainedKeyInput collects the field values for NewTailDrainedKey.
type TailDrainedKeyInput struct {
	StreamKind   uint32
	FilerID      string
	DelaySeconds int64
}

// NewTailDrainedKey builds a ZAP-encoded TailDrainedKey message from in and returns the bytes.
func NewTailDrainedKey(in TailDrainedKeyInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(tailDrainedKeySize)
	ob.SetUint32(tailDrainedKeyStreamKindOff, in.StreamKind)
	ob.SetText(tailDrainedKeyFilerIDOff, in.FilerID)
	ob.SetInt64(tailDrainedKeyDelaySecondsOff, in.DelaySeconds)
	ob.FinishAsRoot()
	return b.Finish()
}
