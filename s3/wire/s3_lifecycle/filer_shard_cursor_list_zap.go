// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	filerShardCursorListCursorsOff = 0
	filerShardCursorListSize       = 8
)

// FilerShardCursorList is a zero-copy view into a ZAP-encoded FilerShardCursorList message.
type FilerShardCursorList struct{ o zap.Object }

// WrapFilerShardCursorList parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapFilerShardCursorList(b []byte) (FilerShardCursorList, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FilerShardCursorList{}, err
	}
	return FilerShardCursorList{o: m.Root()}, nil
}

// Cursors returns the repeated FilerShardCursor field as a variable-element list.
// Read element i with FilerShardCursor{o: Cursors().ObjectAt(i)}.
func (t FilerShardCursorList) Cursors() zap.List { return t.o.List(filerShardCursorListCursorsOff) }

// FilerShardCursorListInput collects the field values for NewFilerShardCursorList.
// Cursors holds pre-built FilerShardCursor sub-buffers (see NewFilerShardCursor).
type FilerShardCursorListInput struct {
	Cursors [][]byte
}

// buildFilerShardCursorList lays a FilerShardCursorList object into b and returns
// its offset, for embedding as a nested field via ObjectBuilder.SetObject. A nil
// in yields offset 0 (a null nested object).
func buildFilerShardCursorList(b *zap.Builder, in *FilerShardCursorListInput) int {
	if in == nil {
		return 0
	}
	cursorsLB := b.StartList(0)
	for _, elem := range in.Cursors {
		cursorsLB.AddObjectBytes(elem)
	}
	listOff := cursorsLB.FinishOffset()
	ob := b.StartObject(filerShardCursorListSize)
	ob.SetList(filerShardCursorListCursorsOff, listOff, len(in.Cursors))
	return ob.Finish()
}

// NewFilerShardCursorList builds a ZAP-encoded FilerShardCursorList message from in and returns the bytes.
func NewFilerShardCursorList(in FilerShardCursorListInput) []byte {
	b := zap.NewBuilder(256)
	cursorsLB := b.StartList(0)
	for _, elem := range in.Cursors {
		cursorsLB.AddObjectBytes(elem)
	}
	listOff := cursorsLB.FinishOffset()
	ob := b.StartObject(filerShardCursorListSize)
	ob.SetList(filerShardCursorListCursorsOff, listOff, len(in.Cursors))
	ob.FinishAsRoot()
	return b.Finish()
}
