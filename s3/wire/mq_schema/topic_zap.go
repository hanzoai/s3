// Code generated from mq_schema.proto; DO NOT EDIT.

package mq_schemawire

import (
	zap "github.com/zap-proto/go"
)

// --- Topic ---

const (
	topicNamespaceOff = 0
	topicNameOff      = 8
	topicSize         = 16
)

// Topic is a zero-copy view into a ZAP-encoded Topic message.
type Topic struct{ o zap.Object }

// WrapTopic parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapTopic(b []byte) (Topic, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Topic{}, err
	}
	return Topic{o: m.Root()}, nil
}

// Namespace reads the namespace field (proto field 1, string).
func (t Topic) Namespace() string { return t.o.Text(topicNamespaceOff) }

// Name reads the name field (proto field 2, string).
func (t Topic) Name() string { return t.o.Text(topicNameOff) }

// TopicInput collects the field values for NewTopic.
type TopicInput struct {
	Namespace string
	Name      string
}

// NewTopic builds a ZAP-encoded Topic message from in and returns the bytes.
func NewTopic(in TopicInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(topicSize)
	ob.SetText(topicNamespaceOff, in.Namespace)
	ob.SetText(topicNameOff, in.Name)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- Partition ---

const (
	partitionRingSizeOff   = 0
	partitionRangeStartOff = 4
	partitionRangeStopOff  = 8
	partitionUnixTimeNsOff = 16
	partitionSize          = 24
)

// Partition is a zero-copy view into a ZAP-encoded Partition message.
type Partition struct{ o zap.Object }

// WrapPartition parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapPartition(b []byte) (Partition, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Partition{}, err
	}
	return Partition{o: m.Root()}, nil
}

// RingSize reads the ring_size field (proto field 1, int32).
func (t Partition) RingSize() int32 { return t.o.Int32(partitionRingSizeOff) }

// RangeStart reads the range_start field (proto field 2, int32).
func (t Partition) RangeStart() int32 { return t.o.Int32(partitionRangeStartOff) }

// RangeStop reads the range_stop field (proto field 3, int32).
func (t Partition) RangeStop() int32 { return t.o.Int32(partitionRangeStopOff) }

// UnixTimeNs reads the unix_time_ns field (proto field 4, int64).
func (t Partition) UnixTimeNs() int64 { return t.o.Int64(partitionUnixTimeNsOff) }

// PartitionInput collects the field values for NewPartition.
type PartitionInput struct {
	RingSize   int32
	RangeStart int32
	RangeStop  int32
	UnixTimeNs int64
}

// NewPartition builds a ZAP-encoded Partition message from in and returns the bytes.
func NewPartition(in PartitionInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(partitionSize)
	ob.SetInt32(partitionRingSizeOff, in.RingSize)
	ob.SetInt32(partitionRangeStartOff, in.RangeStart)
	ob.SetInt32(partitionRangeStopOff, in.RangeStop)
	ob.SetInt64(partitionUnixTimeNsOff, in.UnixTimeNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PartitionOffset ---

const (
	partitionOffsetPartitionOff   = 0
	partitionOffsetStartTsNsOff   = 8
	partitionOffsetStartOffsetOff = 16
	partitionOffsetSize           = 24
)

// PartitionOffset is a zero-copy view into a ZAP-encoded PartitionOffset message.
type PartitionOffset struct{ o zap.Object }

// WrapPartitionOffset parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapPartitionOffset(b []byte) (PartitionOffset, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PartitionOffset{}, err
	}
	return PartitionOffset{o: m.Root()}, nil
}

// Partition reads the partition field (proto field 1, message Partition).
// The bool is false when the field is absent.
func (t PartitionOffset) Partition() (Partition, bool) {
	b := t.o.Bytes(partitionOffsetPartitionOff)
	if len(b) == 0 {
		return Partition{}, false
	}
	p, err := WrapPartition(b)
	if err != nil {
		return Partition{}, false
	}
	return p, true
}

// StartTsNs reads the start_ts_ns field (proto field 2, int64).
func (t PartitionOffset) StartTsNs() int64 { return t.o.Int64(partitionOffsetStartTsNsOff) }

// StartOffset reads the start_offset field (proto field 3, int64).
func (t PartitionOffset) StartOffset() int64 { return t.o.Int64(partitionOffsetStartOffsetOff) }

// PartitionOffsetInput collects the field values for NewPartitionOffset.
// Partition is the child message's own ZAP buffer (from NewPartition); nil
// encodes an absent partition.
type PartitionOffsetInput struct {
	Partition   []byte
	StartTsNs   int64
	StartOffset int64
}

// NewPartitionOffset builds a ZAP-encoded PartitionOffset message from in and
// returns the bytes.
func NewPartitionOffset(in PartitionOffsetInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(partitionOffsetSize)
	ob.SetBytes(partitionOffsetPartitionOff, in.Partition)
	ob.SetInt64(partitionOffsetStartTsNsOff, in.StartTsNs)
	ob.SetInt64(partitionOffsetStartOffsetOff, in.StartOffset)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- Offset ---

const (
	offsetTopicOff            = 0
	offsetPartitionOffsetsOff = 8
	offsetSize                = 16
)

// Offset is a zero-copy view into a ZAP-encoded Offset message.
type Offset struct{ o zap.Object }

// WrapOffset parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapOffset(b []byte) (Offset, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Offset{}, err
	}
	return Offset{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message Topic). The bool is
// false when the field is absent.
func (t Offset) Topic() (Topic, bool) {
	b := t.o.Bytes(offsetTopicOff)
	if len(b) == 0 {
		return Topic{}, false
	}
	tp, err := WrapTopic(b)
	if err != nil {
		return Topic{}, false
	}
	return tp, true
}

// PartitionOffsetsLen reports the number of partition_offsets elements
// (proto field 2, repeated PartitionOffset).
func (t Offset) PartitionOffsetsLen() int {
	return t.o.List(offsetPartitionOffsetsOff).Len()
}

// PartitionOffsetAt returns the i-th partition_offsets element. The bool is
// false when i is out of range or the element fails to parse.
func (t Offset) PartitionOffsetAt(i int) (PartitionOffset, bool) {
	o := t.o.List(offsetPartitionOffsetsOff).ObjectAt(i)
	if o.IsNull() {
		return PartitionOffset{}, false
	}
	return PartitionOffset{o: o}, true
}

// OffsetInput collects the field values for NewOffset. Topic is the child
// message's own ZAP buffer (from NewTopic); nil encodes an absent topic.
// PartitionOffsets are each a PartitionOffset buffer (from NewPartitionOffset).
type OffsetInput struct {
	Topic            []byte
	PartitionOffsets [][]byte
}

// NewOffset builds a ZAP-encoded Offset message from in and returns the bytes.
func NewOffset(in OffsetInput) []byte {
	b := zap.NewBuilder(256)
	lb := b.StartList(0)
	for _, e := range in.PartitionOffsets {
		lb.AddObjectBytes(e)
	}
	listOff, listLen := lb.Finish()
	ob := b.StartObject(offsetSize)
	ob.SetBytes(offsetTopicOff, in.Topic)
	ob.SetList(offsetPartitionOffsetsOff, listOff, listLen)
	ob.FinishAsRoot()
	return b.Finish()
}
