// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	predicateKeyShardOff    = 0
	predicateKeyPositionOff = 8
	predicateKeySize        = 12
)

// PredicateKey is a zero-copy view into a ZAP-encoded PredicateKey message.
type PredicateKey struct{ o zap.Object }

// WrapPredicateKey parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapPredicateKey(b []byte) (PredicateKey, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PredicateKey{}, err
	}
	return PredicateKey{o: m.Root()}, nil
}

func (t PredicateKey) Shard() string { return t.o.Text(predicateKeyShardOff) }

// Position reads the position field (proto field 2, message MessagePosition).
func (t PredicateKey) Position() MessagePosition {
	return MessagePosition{o: t.o.Object(predicateKeyPositionOff)}
}

// PredicateKeyInput collects the field values for NewPredicateKey.
type PredicateKeyInput struct {
	Shard    string
	Position *MessagePositionInput
}

// buildPredicateKey lays a PredicateKey object into b and returns its offset, for
// embedding as a nested field via ObjectBuilder.SetObject. A nil in yields
// offset 0 (a null nested object).
func buildPredicateKey(b *zap.Builder, in *PredicateKeyInput) int {
	if in == nil {
		return 0
	}
	pos := buildMessagePosition(b, in.Position)
	ob := b.StartObject(predicateKeySize)
	ob.SetText(predicateKeyShardOff, in.Shard)
	ob.SetObject(predicateKeyPositionOff, pos)
	return ob.Finish()
}

// NewPredicateKey builds a ZAP-encoded PredicateKey message from in and returns the bytes.
func NewPredicateKey(in PredicateKeyInput) []byte {
	b := zap.NewBuilder(256)
	pos := buildMessagePosition(b, in.Position)
	ob := b.StartObject(predicateKeySize)
	ob.SetText(predicateKeyShardOff, in.Shard)
	ob.SetObject(predicateKeyPositionOff, pos)
	ob.FinishAsRoot()
	return b.Finish()
}
