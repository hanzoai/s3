// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	originalKeyShardOff        = 0
	originalKeyDelaySecondsOff = 8
	originalKeyPositionOff     = 16
	originalKeySize            = 20
)

// OriginalKey is a zero-copy view into a ZAP-encoded OriginalKey message.
type OriginalKey struct{ o zap.Object }

// WrapOriginalKey parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapOriginalKey(b []byte) (OriginalKey, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return OriginalKey{}, err
	}
	return OriginalKey{o: m.Root()}, nil
}

func (t OriginalKey) Shard() string       { return t.o.Text(originalKeyShardOff) }
func (t OriginalKey) DelaySeconds() int64 { return t.o.Int64(originalKeyDelaySecondsOff) }

// Position reads the position field (proto field 3, message MessagePosition).
func (t OriginalKey) Position() MessagePosition {
	return MessagePosition{o: t.o.Object(originalKeyPositionOff)}
}

// OriginalKeyInput collects the field values for NewOriginalKey.
type OriginalKeyInput struct {
	Shard        string
	DelaySeconds int64
	Position     *MessagePositionInput
}

// buildOriginalKey lays an OriginalKey object into b and returns its offset, for
// embedding as a nested field via ObjectBuilder.SetObject. A nil in yields
// offset 0 (a null nested object).
func buildOriginalKey(b *zap.Builder, in *OriginalKeyInput) int {
	if in == nil {
		return 0
	}
	pos := buildMessagePosition(b, in.Position)
	ob := b.StartObject(originalKeySize)
	ob.SetText(originalKeyShardOff, in.Shard)
	ob.SetInt64(originalKeyDelaySecondsOff, in.DelaySeconds)
	ob.SetObject(originalKeyPositionOff, pos)
	return ob.Finish()
}

// NewOriginalKey builds a ZAP-encoded OriginalKey message from in and returns the bytes.
func NewOriginalKey(in OriginalKeyInput) []byte {
	b := zap.NewBuilder(256)
	pos := buildMessagePosition(b, in.Position)
	ob := b.StartObject(originalKeySize)
	ob.SetText(originalKeyShardOff, in.Shard)
	ob.SetInt64(originalKeyDelaySecondsOff, in.DelaySeconds)
	ob.SetObject(originalKeyPositionOff, pos)
	ob.FinishAsRoot()
	return b.Finish()
}
