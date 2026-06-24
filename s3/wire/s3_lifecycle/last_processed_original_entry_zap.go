// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

// LastProcessedOriginalEntry is the {key,value} entry object for the proto map
// field ReaderState.last_processed_original (map<int64, FilerShardCursorList>).
// The map is encoded on the wire as a variable-element list of these entries.

const (
	lastProcessedOriginalEntryKeyOff   = 0
	lastProcessedOriginalEntryValueOff = 8
	lastProcessedOriginalEntrySize     = 12
)

// LastProcessedOriginalEntry is a zero-copy view into a ZAP-encoded map entry.
type LastProcessedOriginalEntry struct{ o zap.Object }

// WrapLastProcessedOriginalEntry parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapLastProcessedOriginalEntry(b []byte) (LastProcessedOriginalEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LastProcessedOriginalEntry{}, err
	}
	return LastProcessedOriginalEntry{o: m.Root()}, nil
}

func (t LastProcessedOriginalEntry) Key() int64 { return t.o.Int64(lastProcessedOriginalEntryKeyOff) }

// Value reads the map value (FilerShardCursorList).
func (t LastProcessedOriginalEntry) Value() FilerShardCursorList {
	return FilerShardCursorList{o: t.o.Object(lastProcessedOriginalEntryValueOff)}
}

// LastProcessedOriginalEntryInput collects the field values for the map entry.
type LastProcessedOriginalEntryInput struct {
	Key   int64
	Value *FilerShardCursorListInput
}

// NewLastProcessedOriginalEntry builds a ZAP-encoded map entry from in and returns the bytes.
func NewLastProcessedOriginalEntry(in LastProcessedOriginalEntryInput) []byte {
	b := zap.NewBuilder(256)
	value := buildFilerShardCursorList(b, in.Value)
	ob := b.StartObject(lastProcessedOriginalEntrySize)
	ob.SetInt64(lastProcessedOriginalEntryKeyOff, in.Key)
	ob.SetObject(lastProcessedOriginalEntryValueOff, value)
	ob.FinishAsRoot()
	return b.Finish()
}
