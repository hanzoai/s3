// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	readerStatePrimaryFilerEndpointOff   = 0
	readerStateLastProcessedOriginalOff  = 8
	readerStateLastProcessedPredicateOff = 16
	readerStateTailDrainedStreamsOff     = 24
	readerStateLastCheckpointNsOff       = 32
	readerStateSize                      = 40
)

// ReaderState is a zero-copy view into a ZAP-encoded ReaderState message.
type ReaderState struct{ o zap.Object }

// WrapReaderState parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapReaderState(b []byte) (ReaderState, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReaderState{}, err
	}
	return ReaderState{o: m.Root()}, nil
}

func (t ReaderState) PrimaryFilerEndpoint() string {
	return t.o.Text(readerStatePrimaryFilerEndpointOff)
}

// LastProcessedOriginal returns the proto map<int64, FilerShardCursorList> as a
// variable-element list of map entries. Read entry i with
// LastProcessedOriginalEntry{o: LastProcessedOriginal().ObjectAt(i)}.
func (t ReaderState) LastProcessedOriginal() zap.List {
	return t.o.List(readerStateLastProcessedOriginalOff)
}

// LastProcessedPredicate returns the repeated FilerShardCursor field as a
// variable-element list. Read element i with
// FilerShardCursor{o: LastProcessedPredicate().ObjectAt(i)}.
func (t ReaderState) LastProcessedPredicate() zap.List {
	return t.o.List(readerStateLastProcessedPredicateOff)
}

// TailDrainedStreams returns the repeated TailDrainedKey field as a
// variable-element list. Read element i with
// TailDrainedKey{o: TailDrainedStreams().ObjectAt(i)}.
func (t ReaderState) TailDrainedStreams() zap.List {
	return t.o.List(readerStateTailDrainedStreamsOff)
}

func (t ReaderState) LastCheckpointNs() int64 { return t.o.Int64(readerStateLastCheckpointNsOff) }

// ReaderStateInput collects the field values for NewReaderState. The list fields
// hold pre-built sub-buffers: LastProcessedOriginal of LastProcessedOriginalEntry,
// LastProcessedPredicate of FilerShardCursor, TailDrainedStreams of TailDrainedKey.
type ReaderStateInput struct {
	PrimaryFilerEndpoint   string
	LastProcessedOriginal  [][]byte
	LastProcessedPredicate [][]byte
	TailDrainedStreams     [][]byte
	LastCheckpointNs       int64
}

// NewReaderState builds a ZAP-encoded ReaderState message from in and returns the bytes.
func NewReaderState(in ReaderStateInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(readerStateSize)
	ob.SetText(readerStatePrimaryFilerEndpointOff, in.PrimaryFilerEndpoint)
	lastProcessedOriginalLB := b.StartList(0)
	for _, elem := range in.LastProcessedOriginal {
		lastProcessedOriginalLB.AddObjectBytes(elem)
	}
	ob.SetList(readerStateLastProcessedOriginalOff, lastProcessedOriginalLB.FinishOffset(), len(in.LastProcessedOriginal))
	lastProcessedPredicateLB := b.StartList(0)
	for _, elem := range in.LastProcessedPredicate {
		lastProcessedPredicateLB.AddObjectBytes(elem)
	}
	ob.SetList(readerStateLastProcessedPredicateOff, lastProcessedPredicateLB.FinishOffset(), len(in.LastProcessedPredicate))
	tailDrainedStreamsLB := b.StartList(0)
	for _, elem := range in.TailDrainedStreams {
		tailDrainedStreamsLB.AddObjectBytes(elem)
	}
	ob.SetList(readerStateTailDrainedStreamsOff, tailDrainedStreamsLB.FinishOffset(), len(in.TailDrainedStreams))
	ob.SetInt64(readerStateLastCheckpointNsOff, in.LastCheckpointNs)
	ob.FinishAsRoot()
	return b.Finish()
}
