// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	filerShardCursorFilerIDOff  = 0
	filerShardCursorPositionOff = 8
	filerShardCursorSize        = 12
)

// FilerShardCursor is a zero-copy view into a ZAP-encoded FilerShardCursor message.
type FilerShardCursor struct{ o zap.Object }

// WrapFilerShardCursor parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapFilerShardCursor(b []byte) (FilerShardCursor, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FilerShardCursor{}, err
	}
	return FilerShardCursor{o: m.Root()}, nil
}

func (t FilerShardCursor) FilerID() string { return t.o.Text(filerShardCursorFilerIDOff) }
func (t FilerShardCursor) Position() MessagePosition {
	return MessagePosition{o: t.o.Object(filerShardCursorPositionOff)}
}

// FilerShardCursorInput collects the field values for NewFilerShardCursor.
type FilerShardCursorInput struct {
	FilerID  string
	Position *MessagePositionInput
}

// NewFilerShardCursor builds a ZAP-encoded FilerShardCursor message from in and returns the bytes.
func NewFilerShardCursor(in FilerShardCursorInput) []byte {
	b := zap.NewBuilder(256)
	pos := buildMessagePosition(b, in.Position)
	ob := b.StartObject(filerShardCursorSize)
	ob.SetText(filerShardCursorFilerIDOff, in.FilerID)
	ob.SetObject(filerShardCursorPositionOff, pos)
	ob.FinishAsRoot()
	return b.Finish()
}
