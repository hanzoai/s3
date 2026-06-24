// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

// StreamKeyWhich enumerates the StreamKey oneof "key" variants. The values
// match the proto oneof field numbers; StreamKeyWhichNone (0) means unset.
const (
	StreamKeyWhichNone      uint32 = 0
	StreamKeyWhichOriginal  uint32 = 1
	StreamKeyWhichPredicate uint32 = 2
	StreamKeyWhichBootstrap uint32 = 3
	StreamKeyWhichPending   uint32 = 4
)

const (
	streamKeyWhichOff = 0
	streamKeyValueOff = 4
	streamKeySize     = 8
)

// StreamKey is a zero-copy view into a ZAP-encoded StreamKey message. It models
// the proto oneof "key" as a discriminant (Which) plus a single object pointer
// to the selected variant; the variant accessors are valid only when Which
// matches.
type StreamKey struct{ o zap.Object }

// WrapStreamKey parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapStreamKey(b []byte) (StreamKey, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return StreamKey{}, err
	}
	return StreamKey{o: m.Root()}, nil
}

// Which reports which oneof variant is set (one of the StreamKeyWhich* values).
func (t StreamKey) Which() uint32 { return t.o.Uint32(streamKeyWhichOff) }

// Original returns the original variant; valid only when Which == StreamKeyWhichOriginal.
func (t StreamKey) Original() OriginalKey {
	return OriginalKey{o: t.o.Object(streamKeyValueOff)}
}

// Predicate returns the predicate variant; valid only when Which == StreamKeyWhichPredicate.
func (t StreamKey) Predicate() PredicateKey {
	return PredicateKey{o: t.o.Object(streamKeyValueOff)}
}

// Bootstrap returns the bootstrap variant; valid only when Which == StreamKeyWhichBootstrap.
func (t StreamKey) Bootstrap() BootstrapKey {
	return BootstrapKey{o: t.o.Object(streamKeyValueOff)}
}

// Pending returns the pending variant; valid only when Which == StreamKeyWhichPending.
func (t StreamKey) Pending() PendingKey {
	return PendingKey{o: t.o.Object(streamKeyValueOff)}
}

// StreamKeyInput collects the field values for NewStreamKey. Set Which to the
// active variant and supply that variant's Input; the other variant fields are
// ignored.
type StreamKeyInput struct {
	Which     uint32
	Original  *OriginalKeyInput
	Predicate *PredicateKeyInput
	Bootstrap *BootstrapKeyInput
	Pending   *PendingKeyInput
}

// streamKeyVariantOffset builds the variant object selected by in.Which into b
// (children-first) and returns its offset, or 0 if no variant is set.
func streamKeyVariantOffset(b *zap.Builder, in StreamKeyInput) int {
	switch in.Which {
	case StreamKeyWhichOriginal:
		return buildOriginalKey(b, in.Original)
	case StreamKeyWhichPredicate:
		return buildPredicateKey(b, in.Predicate)
	case StreamKeyWhichBootstrap:
		return buildBootstrapKey(b, in.Bootstrap)
	case StreamKeyWhichPending:
		return buildPendingKey(b, in.Pending)
	default:
		return 0
	}
}

// buildStreamKey lays a StreamKey object into b and returns its offset, for
// embedding as a nested field via ObjectBuilder.SetObject. A nil in yields
// offset 0 (a null nested object).
func buildStreamKey(b *zap.Builder, in *StreamKeyInput) int {
	if in == nil {
		return 0
	}
	variantOff := streamKeyVariantOffset(b, *in)
	ob := b.StartObject(streamKeySize)
	ob.SetUint32(streamKeyWhichOff, in.Which)
	ob.SetObject(streamKeyValueOff, variantOff)
	return ob.Finish()
}

// NewStreamKey builds a ZAP-encoded StreamKey message from in and returns the bytes.
func NewStreamKey(in StreamKeyInput) []byte {
	b := zap.NewBuilder(256)
	variantOff := streamKeyVariantOffset(b, in)
	ob := b.StartObject(streamKeySize)
	ob.SetUint32(streamKeyWhichOff, in.Which)
	ob.SetObject(streamKeyValueOff, variantOff)
	ob.FinishAsRoot()
	return b.Finish()
}
