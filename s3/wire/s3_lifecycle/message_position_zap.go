// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	messagePositionTsNsOff   = 0
	messagePositionOffsetOff = 8
	messagePositionSize      = 16
)

// MessagePosition is a zero-copy view into a ZAP-encoded MessagePosition message.
type MessagePosition struct{ o zap.Object }

// WrapMessagePosition parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapMessagePosition(b []byte) (MessagePosition, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return MessagePosition{}, err
	}
	return MessagePosition{o: m.Root()}, nil
}

func (t MessagePosition) TsNs() int64   { return t.o.Int64(messagePositionTsNsOff) }
func (t MessagePosition) Offset() int64 { return t.o.Int64(messagePositionOffsetOff) }

// MessagePositionInput collects the field values for NewMessagePosition.
type MessagePositionInput struct {
	TsNs   int64
	Offset int64
}

// putMessagePosition lays in's fields into ob (a fresh StartObject(messagePositionSize)).
func putMessagePosition(ob *zap.ObjectBuilder, in MessagePositionInput) {
	ob.SetInt64(messagePositionTsNsOff, in.TsNs)
	ob.SetInt64(messagePositionOffsetOff, in.Offset)
}

// buildMessagePosition lays a MessagePosition object into b and returns its
// offset, for embedding as a nested field via ObjectBuilder.SetObject. A nil in
// yields offset 0 (a null nested object).
func buildMessagePosition(b *zap.Builder, in *MessagePositionInput) int {
	if in == nil {
		return 0
	}
	ob := b.StartObject(messagePositionSize)
	putMessagePosition(ob, *in)
	return ob.Finish()
}

// NewMessagePosition builds a ZAP-encoded MessagePosition message from in and returns the bytes.
func NewMessagePosition(in MessagePositionInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(messagePositionSize)
	putMessagePosition(ob, in)
	ob.FinishAsRoot()
	return b.Finish()
}
