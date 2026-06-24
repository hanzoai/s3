// Code generated from mq_broker.proto; DO NOT EDIT.

package mq_brokerwire

import (
	zap "github.com/zap-proto/go"
)

// Subscriber<->coordinator messages. Each oneof variant message is its own
// standalone type; the parent embeds the selected variant as a nested object
// (4-byte object pointer) selected by the parent's which-tag. schema_pb.Topic /
// schema_pb.Partition singular fields inside the variants are stored as the
// child's own ZAP buffer (8-byte tail pointer); re-wrap with
// mq_schemawire.WrapTopic / WrapPartition.

// --- SubscriberToSubCoordinatorRequestInitMessage ---
// (oneof message variant `init` of SubscriberToSubCoordinatorRequest)

const (
	subscriberToSubCoordinatorRequestInitMessageConsumerGroupOff           = 0  // string
	subscriberToSubCoordinatorRequestInitMessageConsumerGroupInstanceIdOff = 8  // string
	subscriberToSubCoordinatorRequestInitMessageTopicOff                   = 16 // message schema_pb.Topic
	subscriberToSubCoordinatorRequestInitMessageMaxPartitionCountOff       = 24 // int32
	subscriberToSubCoordinatorRequestInitMessageRebalanceSecondsOff        = 28 // int32
	subscriberToSubCoordinatorRequestInitMessageSize                       = 32
)

// SubscriberToSubCoordinatorRequestInitMessage is a zero-copy view into a
// ZAP-encoded SubscriberToSubCoordinatorRequest.InitMessage message.
type SubscriberToSubCoordinatorRequestInitMessage struct{ o zap.Object }

// WrapSubscriberToSubCoordinatorRequestInitMessage parses b and returns a typed
// view. Returns an error if the wire-level checks (magic, version, size) fail.
func WrapSubscriberToSubCoordinatorRequestInitMessage(b []byte) (SubscriberToSubCoordinatorRequestInitMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscriberToSubCoordinatorRequestInitMessage{}, err
	}
	return SubscriberToSubCoordinatorRequestInitMessage{o: m.Root()}, nil
}

// ConsumerGroup reads the consumer_group field (proto field 1, string).
func (t SubscriberToSubCoordinatorRequestInitMessage) ConsumerGroup() string {
	return t.o.Text(subscriberToSubCoordinatorRequestInitMessageConsumerGroupOff)
}

// ConsumerGroupInstanceId reads the consumer_group_instance_id field (proto
// field 2, string).
func (t SubscriberToSubCoordinatorRequestInitMessage) ConsumerGroupInstanceId() string {
	return t.o.Text(subscriberToSubCoordinatorRequestInitMessageConsumerGroupInstanceIdOff)
}

// Topic reads the topic field (proto field 3, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t SubscriberToSubCoordinatorRequestInitMessage) Topic() []byte {
	return t.o.Bytes(subscriberToSubCoordinatorRequestInitMessageTopicOff)
}

// MaxPartitionCount reads the max_partition_count field (proto field 4, int32).
func (t SubscriberToSubCoordinatorRequestInitMessage) MaxPartitionCount() int32 {
	return t.o.Int32(subscriberToSubCoordinatorRequestInitMessageMaxPartitionCountOff)
}

// RebalanceSeconds reads the rebalance_seconds field (proto field 5, int32).
func (t SubscriberToSubCoordinatorRequestInitMessage) RebalanceSeconds() int32 {
	return t.o.Int32(subscriberToSubCoordinatorRequestInitMessageRebalanceSecondsOff)
}

// SubscriberToSubCoordinatorRequestInitMessageInput collects the field values
// for NewSubscriberToSubCoordinatorRequestInitMessage. Topic is the child
// message's own ZAP buffer (from mq_schemawire.NewTopic); nil encodes absent.
type SubscriberToSubCoordinatorRequestInitMessageInput struct {
	ConsumerGroup           string
	ConsumerGroupInstanceId string
	Topic                   []byte
	MaxPartitionCount       int32
	RebalanceSeconds        int32
}

// NewSubscriberToSubCoordinatorRequestInitMessage builds a ZAP-encoded
// SubscriberToSubCoordinatorRequest.InitMessage message from in and returns the
// bytes.
func NewSubscriberToSubCoordinatorRequestInitMessage(in SubscriberToSubCoordinatorRequestInitMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(subscriberToSubCoordinatorRequestInitMessageSize)
	ob.SetText(subscriberToSubCoordinatorRequestInitMessageConsumerGroupOff, in.ConsumerGroup)
	ob.SetText(subscriberToSubCoordinatorRequestInitMessageConsumerGroupInstanceIdOff, in.ConsumerGroupInstanceId)
	ob.SetBytes(subscriberToSubCoordinatorRequestInitMessageTopicOff, in.Topic)
	ob.SetInt32(subscriberToSubCoordinatorRequestInitMessageMaxPartitionCountOff, in.MaxPartitionCount)
	ob.SetInt32(subscriberToSubCoordinatorRequestInitMessageRebalanceSecondsOff, in.RebalanceSeconds)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscriberToSubCoordinatorRequestAckAssignmentMessage ---
// (oneof message variant `ack_assignment` of SubscriberToSubCoordinatorRequest)

const (
	subscriberToSubCoordinatorRequestAckAssignmentMessagePartitionOff = 0 // message schema_pb.Partition
	subscriberToSubCoordinatorRequestAckAssignmentMessageSize         = 8
)

// SubscriberToSubCoordinatorRequestAckAssignmentMessage is a zero-copy view into
// a ZAP-encoded SubscriberToSubCoordinatorRequest.AckAssignmentMessage message.
type SubscriberToSubCoordinatorRequestAckAssignmentMessage struct{ o zap.Object }

// WrapSubscriberToSubCoordinatorRequestAckAssignmentMessage parses b and returns
// a typed view. Returns an error if the wire-level checks fail.
func WrapSubscriberToSubCoordinatorRequestAckAssignmentMessage(b []byte) (SubscriberToSubCoordinatorRequestAckAssignmentMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscriberToSubCoordinatorRequestAckAssignmentMessage{}, err
	}
	return SubscriberToSubCoordinatorRequestAckAssignmentMessage{o: m.Root()}, nil
}

// Partition reads the partition field (proto field 1, message
// schema_pb.Partition) as the child's raw ZAP buffer; nil when absent. Re-wrap
// with mq_schemawire.WrapPartition.
func (t SubscriberToSubCoordinatorRequestAckAssignmentMessage) Partition() []byte {
	return t.o.Bytes(subscriberToSubCoordinatorRequestAckAssignmentMessagePartitionOff)
}

// SubscriberToSubCoordinatorRequestAckAssignmentMessageInput collects the field
// values for NewSubscriberToSubCoordinatorRequestAckAssignmentMessage. Partition
// is the child message's own ZAP buffer (from mq_schemawire.NewPartition).
type SubscriberToSubCoordinatorRequestAckAssignmentMessageInput struct {
	Partition []byte
}

// NewSubscriberToSubCoordinatorRequestAckAssignmentMessage builds a ZAP-encoded
// SubscriberToSubCoordinatorRequest.AckAssignmentMessage message from in.
func NewSubscriberToSubCoordinatorRequestAckAssignmentMessage(in SubscriberToSubCoordinatorRequestAckAssignmentMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(subscriberToSubCoordinatorRequestAckAssignmentMessageSize)
	ob.SetBytes(subscriberToSubCoordinatorRequestAckAssignmentMessagePartitionOff, in.Partition)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscriberToSubCoordinatorRequestAckUnAssignmentMessage ---
// (oneof message variant `ack_un_assignment` of SubscriberToSubCoordinatorRequest)

const (
	subscriberToSubCoordinatorRequestAckUnAssignmentMessagePartitionOff = 0 // message schema_pb.Partition
	subscriberToSubCoordinatorRequestAckUnAssignmentMessageSize         = 8
)

// SubscriberToSubCoordinatorRequestAckUnAssignmentMessage is a zero-copy view
// into a ZAP-encoded SubscriberToSubCoordinatorRequest.AckUnAssignmentMessage
// message.
type SubscriberToSubCoordinatorRequestAckUnAssignmentMessage struct{ o zap.Object }

// WrapSubscriberToSubCoordinatorRequestAckUnAssignmentMessage parses b and
// returns a typed view. Returns an error if the wire-level checks fail.
func WrapSubscriberToSubCoordinatorRequestAckUnAssignmentMessage(b []byte) (SubscriberToSubCoordinatorRequestAckUnAssignmentMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscriberToSubCoordinatorRequestAckUnAssignmentMessage{}, err
	}
	return SubscriberToSubCoordinatorRequestAckUnAssignmentMessage{o: m.Root()}, nil
}

// Partition reads the partition field (proto field 1, message
// schema_pb.Partition) as the child's raw ZAP buffer; nil when absent. Re-wrap
// with mq_schemawire.WrapPartition.
func (t SubscriberToSubCoordinatorRequestAckUnAssignmentMessage) Partition() []byte {
	return t.o.Bytes(subscriberToSubCoordinatorRequestAckUnAssignmentMessagePartitionOff)
}

// SubscriberToSubCoordinatorRequestAckUnAssignmentMessageInput collects the
// field values for NewSubscriberToSubCoordinatorRequestAckUnAssignmentMessage.
// Partition is the child message's own ZAP buffer (from
// mq_schemawire.NewPartition).
type SubscriberToSubCoordinatorRequestAckUnAssignmentMessageInput struct {
	Partition []byte
}

// NewSubscriberToSubCoordinatorRequestAckUnAssignmentMessage builds a
// ZAP-encoded SubscriberToSubCoordinatorRequest.AckUnAssignmentMessage message.
func NewSubscriberToSubCoordinatorRequestAckUnAssignmentMessage(in SubscriberToSubCoordinatorRequestAckUnAssignmentMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(subscriberToSubCoordinatorRequestAckUnAssignmentMessageSize)
	ob.SetBytes(subscriberToSubCoordinatorRequestAckUnAssignmentMessagePartitionOff, in.Partition)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscriberToSubCoordinatorRequest ---

// SubscriberToSubCoordinatorRequest.message oneof tags (proto field numbers,
// stable wire ids). Zero means no variant is set. The selected variant message
// is embedded as a nested object pointed to by message_value.
const (
	SubscriberToSubCoordinatorRequestMessageNone            uint32 = 0
	SubscriberToSubCoordinatorRequestMessageInit            uint32 = 1
	SubscriberToSubCoordinatorRequestMessageAckAssignment   uint32 = 2
	SubscriberToSubCoordinatorRequestMessageAckUnAssignment uint32 = 3
)

const (
	subscriberToSubCoordinatorRequestMessageWhichOff = 0 // oneof tag (uint32)
	subscriberToSubCoordinatorRequestMessageValueOff = 4 // embedded variant (obj ptr)
	subscriberToSubCoordinatorRequestSize            = 8
)

// SubscriberToSubCoordinatorRequest is a zero-copy view into a ZAP-encoded
// SubscriberToSubCoordinatorRequest message.
type SubscriberToSubCoordinatorRequest struct{ o zap.Object }

// WrapSubscriberToSubCoordinatorRequest parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapSubscriberToSubCoordinatorRequest(b []byte) (SubscriberToSubCoordinatorRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscriberToSubCoordinatorRequest{}, err
	}
	return SubscriberToSubCoordinatorRequest{o: m.Root()}, nil
}

// WhichMessage returns the set oneof tag
// (SubscriberToSubCoordinatorRequestMessage* constants).
func (t SubscriberToSubCoordinatorRequest) WhichMessage() uint32 {
	return t.o.Uint32(subscriberToSubCoordinatorRequestMessageWhichOff)
}

// Init returns the embedded init variant.
func (t SubscriberToSubCoordinatorRequest) Init() SubscriberToSubCoordinatorRequestInitMessage {
	return SubscriberToSubCoordinatorRequestInitMessage{o: t.o.Object(subscriberToSubCoordinatorRequestMessageValueOff)}
}

// AckAssignment returns the embedded ack_assignment variant.
func (t SubscriberToSubCoordinatorRequest) AckAssignment() SubscriberToSubCoordinatorRequestAckAssignmentMessage {
	return SubscriberToSubCoordinatorRequestAckAssignmentMessage{o: t.o.Object(subscriberToSubCoordinatorRequestMessageValueOff)}
}

// AckUnAssignment returns the embedded ack_un_assignment variant.
func (t SubscriberToSubCoordinatorRequest) AckUnAssignment() SubscriberToSubCoordinatorRequestAckUnAssignmentMessage {
	return SubscriberToSubCoordinatorRequestAckUnAssignmentMessage{o: t.o.Object(subscriberToSubCoordinatorRequestMessageValueOff)}
}

// SubscriberToSubCoordinatorRequestInput collects the field values for
// NewSubscriberToSubCoordinatorRequest. MessageWhich selects the oneof variant;
// MessageValue is that variant's standalone buffer.
type SubscriberToSubCoordinatorRequestInput struct {
	MessageWhich uint32
	MessageValue []byte
}

// NewSubscriberToSubCoordinatorRequest builds a ZAP-encoded
// SubscriberToSubCoordinatorRequest message from in and returns the bytes.
func NewSubscriberToSubCoordinatorRequest(in SubscriberToSubCoordinatorRequestInput) []byte {
	b := zap.NewBuilder(512)
	valueOff := 0
	if len(in.MessageValue) > 0 {
		valueOff = embedMessage(b, in.MessageValue)
	}
	ob := b.StartObject(subscriberToSubCoordinatorRequestSize)
	ob.SetUint32(subscriberToSubCoordinatorRequestMessageWhichOff, in.MessageWhich)
	ob.SetObject(subscriberToSubCoordinatorRequestMessageValueOff, valueOff)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscriberToSubCoordinatorResponseAssignment ---
// (oneof message variant `assignment` of SubscriberToSubCoordinatorResponse)

const (
	subscriberToSubCoordinatorResponseAssignmentPartitionAssignmentOff = 0 // message BrokerPartitionAssignment
	subscriberToSubCoordinatorResponseAssignmentSize                   = 8
)

// SubscriberToSubCoordinatorResponseAssignment is a zero-copy view into a
// ZAP-encoded SubscriberToSubCoordinatorResponse.Assignment message.
type SubscriberToSubCoordinatorResponseAssignment struct{ o zap.Object }

// WrapSubscriberToSubCoordinatorResponseAssignment parses b and returns a typed
// view. Returns an error if the wire-level checks fail.
func WrapSubscriberToSubCoordinatorResponseAssignment(b []byte) (SubscriberToSubCoordinatorResponseAssignment, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscriberToSubCoordinatorResponseAssignment{}, err
	}
	return SubscriberToSubCoordinatorResponseAssignment{o: m.Root()}, nil
}

// PartitionAssignment reads the partition_assignment field (proto field 1,
// message BrokerPartitionAssignment). The bool is false when the field is
// absent.
func (t SubscriberToSubCoordinatorResponseAssignment) PartitionAssignment() (BrokerPartitionAssignment, bool) {
	b := t.o.Bytes(subscriberToSubCoordinatorResponseAssignmentPartitionAssignmentOff)
	if len(b) == 0 {
		return BrokerPartitionAssignment{}, false
	}
	a, err := WrapBrokerPartitionAssignment(b)
	if err != nil {
		return BrokerPartitionAssignment{}, false
	}
	return a, true
}

// SubscriberToSubCoordinatorResponseAssignmentInput collects the field values
// for NewSubscriberToSubCoordinatorResponseAssignment. PartitionAssignment is
// the child message's own ZAP buffer (from NewBrokerPartitionAssignment).
type SubscriberToSubCoordinatorResponseAssignmentInput struct {
	PartitionAssignment []byte
}

// NewSubscriberToSubCoordinatorResponseAssignment builds a ZAP-encoded
// SubscriberToSubCoordinatorResponse.Assignment message from in.
func NewSubscriberToSubCoordinatorResponseAssignment(in SubscriberToSubCoordinatorResponseAssignmentInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(subscriberToSubCoordinatorResponseAssignmentSize)
	ob.SetBytes(subscriberToSubCoordinatorResponseAssignmentPartitionAssignmentOff, in.PartitionAssignment)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscriberToSubCoordinatorResponseUnAssignment ---
// (oneof message variant `un_assignment` of SubscriberToSubCoordinatorResponse)

const (
	subscriberToSubCoordinatorResponseUnAssignmentPartitionOff = 0 // message schema_pb.Partition
	subscriberToSubCoordinatorResponseUnAssignmentSize         = 8
)

// SubscriberToSubCoordinatorResponseUnAssignment is a zero-copy view into a
// ZAP-encoded SubscriberToSubCoordinatorResponse.UnAssignment message.
type SubscriberToSubCoordinatorResponseUnAssignment struct{ o zap.Object }

// WrapSubscriberToSubCoordinatorResponseUnAssignment parses b and returns a
// typed view. Returns an error if the wire-level checks fail.
func WrapSubscriberToSubCoordinatorResponseUnAssignment(b []byte) (SubscriberToSubCoordinatorResponseUnAssignment, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscriberToSubCoordinatorResponseUnAssignment{}, err
	}
	return SubscriberToSubCoordinatorResponseUnAssignment{o: m.Root()}, nil
}

// Partition reads the partition field (proto field 1, message
// schema_pb.Partition) as the child's raw ZAP buffer; nil when absent. Re-wrap
// with mq_schemawire.WrapPartition.
func (t SubscriberToSubCoordinatorResponseUnAssignment) Partition() []byte {
	return t.o.Bytes(subscriberToSubCoordinatorResponseUnAssignmentPartitionOff)
}

// SubscriberToSubCoordinatorResponseUnAssignmentInput collects the field values
// for NewSubscriberToSubCoordinatorResponseUnAssignment. Partition is the child
// message's own ZAP buffer (from mq_schemawire.NewPartition).
type SubscriberToSubCoordinatorResponseUnAssignmentInput struct {
	Partition []byte
}

// NewSubscriberToSubCoordinatorResponseUnAssignment builds a ZAP-encoded
// SubscriberToSubCoordinatorResponse.UnAssignment message from in.
func NewSubscriberToSubCoordinatorResponseUnAssignment(in SubscriberToSubCoordinatorResponseUnAssignmentInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(subscriberToSubCoordinatorResponseUnAssignmentSize)
	ob.SetBytes(subscriberToSubCoordinatorResponseUnAssignmentPartitionOff, in.Partition)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscriberToSubCoordinatorResponse ---

// SubscriberToSubCoordinatorResponse.message oneof tags (proto field numbers,
// stable wire ids). Zero means no variant is set. The selected variant message
// is embedded as a nested object pointed to by message_value.
const (
	SubscriberToSubCoordinatorResponseMessageNone         uint32 = 0
	SubscriberToSubCoordinatorResponseMessageAssignment   uint32 = 1
	SubscriberToSubCoordinatorResponseMessageUnAssignment uint32 = 2
)

const (
	subscriberToSubCoordinatorResponseMessageWhichOff = 0 // oneof tag (uint32)
	subscriberToSubCoordinatorResponseMessageValueOff = 4 // embedded variant (obj ptr)
	subscriberToSubCoordinatorResponseSize            = 8
)

// SubscriberToSubCoordinatorResponse is a zero-copy view into a ZAP-encoded
// SubscriberToSubCoordinatorResponse message.
type SubscriberToSubCoordinatorResponse struct{ o zap.Object }

// WrapSubscriberToSubCoordinatorResponse parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapSubscriberToSubCoordinatorResponse(b []byte) (SubscriberToSubCoordinatorResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscriberToSubCoordinatorResponse{}, err
	}
	return SubscriberToSubCoordinatorResponse{o: m.Root()}, nil
}

// WhichMessage returns the set oneof tag
// (SubscriberToSubCoordinatorResponseMessage* constants).
func (t SubscriberToSubCoordinatorResponse) WhichMessage() uint32 {
	return t.o.Uint32(subscriberToSubCoordinatorResponseMessageWhichOff)
}

// Assignment returns the embedded assignment variant.
func (t SubscriberToSubCoordinatorResponse) Assignment() SubscriberToSubCoordinatorResponseAssignment {
	return SubscriberToSubCoordinatorResponseAssignment{o: t.o.Object(subscriberToSubCoordinatorResponseMessageValueOff)}
}

// UnAssignment returns the embedded un_assignment variant.
func (t SubscriberToSubCoordinatorResponse) UnAssignment() SubscriberToSubCoordinatorResponseUnAssignment {
	return SubscriberToSubCoordinatorResponseUnAssignment{o: t.o.Object(subscriberToSubCoordinatorResponseMessageValueOff)}
}

// SubscriberToSubCoordinatorResponseInput collects the field values for
// NewSubscriberToSubCoordinatorResponse. MessageWhich selects the oneof variant;
// MessageValue is that variant's standalone buffer.
type SubscriberToSubCoordinatorResponseInput struct {
	MessageWhich uint32
	MessageValue []byte
}

// NewSubscriberToSubCoordinatorResponse builds a ZAP-encoded
// SubscriberToSubCoordinatorResponse message from in and returns the bytes.
func NewSubscriberToSubCoordinatorResponse(in SubscriberToSubCoordinatorResponseInput) []byte {
	b := zap.NewBuilder(512)
	valueOff := 0
	if len(in.MessageValue) > 0 {
		valueOff = embedMessage(b, in.MessageValue)
	}
	ob := b.StartObject(subscriberToSubCoordinatorResponseSize)
	ob.SetUint32(subscriberToSubCoordinatorResponseMessageWhichOff, in.MessageWhich)
	ob.SetObject(subscriberToSubCoordinatorResponseMessageValueOff, valueOff)
	ob.FinishAsRoot()
	return b.Finish()
}
