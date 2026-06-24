// Code generated from mq_broker.proto; DO NOT EDIT.

package mq_brokerwire

import (
	zap "github.com/zap-proto/go"
)

// Data-plane messages (publish / subscribe / follow-me). Each oneof variant
// message is its own standalone type; the parent embeds the selected variant as
// a nested object (4-byte object pointer) selected by the parent's which-tag.
// schema_pb.Topic / schema_pb.Partition / schema_pb.PartitionOffset singular
// fields are stored as the child's own ZAP buffer (8-byte tail pointer); re-wrap
// with the matching mq_schemawire.Wrap*. schema_pb.OffsetType is a uint32 scalar
// (see mq_schemawire.OffsetType*).

// --- ControlMessage ---

const (
	controlMessageIsCloseOff       = 0 // bool
	controlMessagePublisherNameOff = 8 // string
	controlMessageSize             = 16
)

// ControlMessage is a zero-copy view into a ZAP-encoded ControlMessage message.
type ControlMessage struct{ o zap.Object }

// WrapControlMessage parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapControlMessage(b []byte) (ControlMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ControlMessage{}, err
	}
	return ControlMessage{o: m.Root()}, nil
}

// IsClose reads the is_close field (proto field 1, bool).
func (t ControlMessage) IsClose() bool { return t.o.Bool(controlMessageIsCloseOff) }

// PublisherName reads the publisher_name field (proto field 2, string).
func (t ControlMessage) PublisherName() string { return t.o.Text(controlMessagePublisherNameOff) }

// ControlMessageInput collects the field values for NewControlMessage.
type ControlMessageInput struct {
	IsClose       bool
	PublisherName string
}

// NewControlMessage builds a ZAP-encoded ControlMessage message from in and
// returns the bytes.
func NewControlMessage(in ControlMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(controlMessageSize)
	ob.SetBool(controlMessageIsCloseOff, in.IsClose)
	ob.SetText(controlMessagePublisherNameOff, in.PublisherName)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DataMessage ---

const (
	dataMessageKeyOff   = 0  // bytes
	dataMessageValueOff = 8  // bytes
	dataMessageTsNsOff  = 16 // int64
	dataMessageCtrlOff  = 24 // message ControlMessage
	dataMessageSize     = 32
)

// DataMessage is a zero-copy view into a ZAP-encoded DataMessage message.
type DataMessage struct{ o zap.Object }

// WrapDataMessage parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapDataMessage(b []byte) (DataMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DataMessage{}, err
	}
	return DataMessage{o: m.Root()}, nil
}

// Key reads the key field (proto field 1, bytes).
func (t DataMessage) Key() []byte { return t.o.Bytes(dataMessageKeyOff) }

// Value reads the value field (proto field 2, bytes).
func (t DataMessage) Value() []byte { return t.o.Bytes(dataMessageValueOff) }

// TsNs reads the ts_ns field (proto field 3, int64).
func (t DataMessage) TsNs() int64 { return t.o.Int64(dataMessageTsNsOff) }

// Ctrl reads the ctrl field (proto field 4, message ControlMessage). The bool
// is false when the field is absent.
func (t DataMessage) Ctrl() (ControlMessage, bool) {
	b := t.o.Bytes(dataMessageCtrlOff)
	if len(b) == 0 {
		return ControlMessage{}, false
	}
	c, err := WrapControlMessage(b)
	if err != nil {
		return ControlMessage{}, false
	}
	return c, true
}

// DataMessageInput collects the field values for NewDataMessage. Ctrl is the
// child message's own ZAP buffer (from NewControlMessage); nil encodes absent.
type DataMessageInput struct {
	Key   []byte
	Value []byte
	TsNs  int64
	Ctrl  []byte
}

// NewDataMessage builds a ZAP-encoded DataMessage message from in and returns
// the bytes.
func NewDataMessage(in DataMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(dataMessageSize)
	ob.SetBytes(dataMessageKeyOff, in.Key)
	ob.SetBytes(dataMessageValueOff, in.Value)
	ob.SetInt64(dataMessageTsNsOff, in.TsNs)
	ob.SetBytes(dataMessageCtrlOff, in.Ctrl)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PublishMessageRequestInitMessage ---
// (oneof message variant `init` of PublishMessageRequest)

const (
	publishMessageRequestInitMessageTopicOff          = 0  // message schema_pb.Topic
	publishMessageRequestInitMessagePartitionOff      = 8  // message schema_pb.Partition
	publishMessageRequestInitMessageAckIntervalOff    = 16 // int32
	publishMessageRequestInitMessageFollowerBrokerOff = 24 // string
	publishMessageRequestInitMessagePublisherNameOff  = 32 // string
	publishMessageRequestInitMessageSize              = 40
)

// PublishMessageRequestInitMessage is a zero-copy view into a ZAP-encoded
// PublishMessageRequest.InitMessage message.
type PublishMessageRequestInitMessage struct{ o zap.Object }

// WrapPublishMessageRequestInitMessage parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapPublishMessageRequestInitMessage(b []byte) (PublishMessageRequestInitMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PublishMessageRequestInitMessage{}, err
	}
	return PublishMessageRequestInitMessage{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t PublishMessageRequestInitMessage) Topic() []byte {
	return t.o.Bytes(publishMessageRequestInitMessageTopicOff)
}

// Partition reads the partition field (proto field 2, message
// schema_pb.Partition) as the child's raw ZAP buffer; nil when absent. Re-wrap
// with mq_schemawire.WrapPartition.
func (t PublishMessageRequestInitMessage) Partition() []byte {
	return t.o.Bytes(publishMessageRequestInitMessagePartitionOff)
}

// AckInterval reads the ack_interval field (proto field 3, int32).
func (t PublishMessageRequestInitMessage) AckInterval() int32 {
	return t.o.Int32(publishMessageRequestInitMessageAckIntervalOff)
}

// FollowerBroker reads the follower_broker field (proto field 4, string).
func (t PublishMessageRequestInitMessage) FollowerBroker() string {
	return t.o.Text(publishMessageRequestInitMessageFollowerBrokerOff)
}

// PublisherName reads the publisher_name field (proto field 5, string).
func (t PublishMessageRequestInitMessage) PublisherName() string {
	return t.o.Text(publishMessageRequestInitMessagePublisherNameOff)
}

// PublishMessageRequestInitMessageInput collects the field values for
// NewPublishMessageRequestInitMessage. Topic and Partition are each the child
// message's own ZAP buffer (from mq_schemawire.NewTopic / NewPartition).
type PublishMessageRequestInitMessageInput struct {
	Topic          []byte
	Partition      []byte
	AckInterval    int32
	FollowerBroker string
	PublisherName  string
}

// NewPublishMessageRequestInitMessage builds a ZAP-encoded
// PublishMessageRequest.InitMessage message from in and returns the bytes.
func NewPublishMessageRequestInitMessage(in PublishMessageRequestInitMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(publishMessageRequestInitMessageSize)
	ob.SetBytes(publishMessageRequestInitMessageTopicOff, in.Topic)
	ob.SetBytes(publishMessageRequestInitMessagePartitionOff, in.Partition)
	ob.SetInt32(publishMessageRequestInitMessageAckIntervalOff, in.AckInterval)
	ob.SetText(publishMessageRequestInitMessageFollowerBrokerOff, in.FollowerBroker)
	ob.SetText(publishMessageRequestInitMessagePublisherNameOff, in.PublisherName)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PublishMessageRequest ---

// PublishMessageRequest.message oneof tags (proto field numbers, stable wire
// ids). Zero means no variant is set. The selected variant message is embedded
// as a nested object pointed to by message_value.
const (
	PublishMessageRequestMessageNone uint32 = 0
	PublishMessageRequestMessageInit uint32 = 1
	PublishMessageRequestMessageData uint32 = 2
)

const (
	publishMessageRequestMessageWhichOff = 0 // oneof tag (uint32)
	publishMessageRequestMessageValueOff = 4 // embedded variant (obj ptr)
	publishMessageRequestSize            = 8
)

// PublishMessageRequest is a zero-copy view into a ZAP-encoded
// PublishMessageRequest message.
type PublishMessageRequest struct{ o zap.Object }

// WrapPublishMessageRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapPublishMessageRequest(b []byte) (PublishMessageRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PublishMessageRequest{}, err
	}
	return PublishMessageRequest{o: m.Root()}, nil
}

// WhichMessage returns the set oneof tag (PublishMessageRequestMessage*
// constants).
func (t PublishMessageRequest) WhichMessage() uint32 {
	return t.o.Uint32(publishMessageRequestMessageWhichOff)
}

// Init returns the embedded init variant.
func (t PublishMessageRequest) Init() PublishMessageRequestInitMessage {
	return PublishMessageRequestInitMessage{o: t.o.Object(publishMessageRequestMessageValueOff)}
}

// Data returns the embedded data variant.
func (t PublishMessageRequest) Data() DataMessage {
	return DataMessage{o: t.o.Object(publishMessageRequestMessageValueOff)}
}

// PublishMessageRequestInput collects the field values for
// NewPublishMessageRequest. MessageWhich selects the oneof variant; MessageValue
// is that variant's standalone buffer.
type PublishMessageRequestInput struct {
	MessageWhich uint32
	MessageValue []byte
}

// NewPublishMessageRequest builds a ZAP-encoded PublishMessageRequest message
// from in and returns the bytes.
func NewPublishMessageRequest(in PublishMessageRequestInput) []byte {
	b := zap.NewBuilder(512)
	valueOff := 0
	if len(in.MessageValue) > 0 {
		valueOff = embedMessage(b, in.MessageValue)
	}
	ob := b.StartObject(publishMessageRequestSize)
	ob.SetUint32(publishMessageRequestMessageWhichOff, in.MessageWhich)
	ob.SetObject(publishMessageRequestMessageValueOff, valueOff)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PublishMessageResponse ---

const (
	publishMessageResponseAckTsNsOff        = 0  // int64
	publishMessageResponseErrorOff          = 8  // string
	publishMessageResponseShouldCloseOff    = 16 // bool
	publishMessageResponseErrorCodeOff      = 20 // int32
	publishMessageResponseAssignedOffsetOff = 24 // int64
	publishMessageResponseSize              = 32
)

// PublishMessageResponse is a zero-copy view into a ZAP-encoded
// PublishMessageResponse message.
type PublishMessageResponse struct{ o zap.Object }

// WrapPublishMessageResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapPublishMessageResponse(b []byte) (PublishMessageResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PublishMessageResponse{}, err
	}
	return PublishMessageResponse{o: m.Root()}, nil
}

// AckTsNs reads the ack_ts_ns field (proto field 1, int64).
func (t PublishMessageResponse) AckTsNs() int64 { return t.o.Int64(publishMessageResponseAckTsNsOff) }

// Error reads the error field (proto field 2, string).
func (t PublishMessageResponse) Error() string { return t.o.Text(publishMessageResponseErrorOff) }

// ShouldClose reads the should_close field (proto field 3, bool).
func (t PublishMessageResponse) ShouldClose() bool {
	return t.o.Bool(publishMessageResponseShouldCloseOff)
}

// ErrorCode reads the error_code field (proto field 4, int32).
func (t PublishMessageResponse) ErrorCode() int32 {
	return t.o.Int32(publishMessageResponseErrorCodeOff)
}

// AssignedOffset reads the assigned_offset field (proto field 5, int64).
func (t PublishMessageResponse) AssignedOffset() int64 {
	return t.o.Int64(publishMessageResponseAssignedOffsetOff)
}

// PublishMessageResponseInput collects the field values for
// NewPublishMessageResponse.
type PublishMessageResponseInput struct {
	AckTsNs        int64
	Error          string
	ShouldClose    bool
	ErrorCode      int32
	AssignedOffset int64
}

// NewPublishMessageResponse builds a ZAP-encoded PublishMessageResponse message
// from in and returns the bytes.
func NewPublishMessageResponse(in PublishMessageResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(publishMessageResponseSize)
	ob.SetInt64(publishMessageResponseAckTsNsOff, in.AckTsNs)
	ob.SetText(publishMessageResponseErrorOff, in.Error)
	ob.SetBool(publishMessageResponseShouldCloseOff, in.ShouldClose)
	ob.SetInt32(publishMessageResponseErrorCodeOff, in.ErrorCode)
	ob.SetInt64(publishMessageResponseAssignedOffsetOff, in.AssignedOffset)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PublishFollowMeRequestInitMessage ---
// (oneof message variant `init` of PublishFollowMeRequest)

const (
	publishFollowMeRequestInitMessageTopicOff     = 0 // message schema_pb.Topic
	publishFollowMeRequestInitMessagePartitionOff = 8 // message schema_pb.Partition
	publishFollowMeRequestInitMessageSize         = 16
)

// PublishFollowMeRequestInitMessage is a zero-copy view into a ZAP-encoded
// PublishFollowMeRequest.InitMessage message.
type PublishFollowMeRequestInitMessage struct{ o zap.Object }

// WrapPublishFollowMeRequestInitMessage parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapPublishFollowMeRequestInitMessage(b []byte) (PublishFollowMeRequestInitMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PublishFollowMeRequestInitMessage{}, err
	}
	return PublishFollowMeRequestInitMessage{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t PublishFollowMeRequestInitMessage) Topic() []byte {
	return t.o.Bytes(publishFollowMeRequestInitMessageTopicOff)
}

// Partition reads the partition field (proto field 2, message
// schema_pb.Partition) as the child's raw ZAP buffer; nil when absent. Re-wrap
// with mq_schemawire.WrapPartition.
func (t PublishFollowMeRequestInitMessage) Partition() []byte {
	return t.o.Bytes(publishFollowMeRequestInitMessagePartitionOff)
}

// PublishFollowMeRequestInitMessageInput collects the field values for
// NewPublishFollowMeRequestInitMessage. Topic and Partition are each the child
// message's own ZAP buffer (from mq_schemawire.NewTopic / NewPartition).
type PublishFollowMeRequestInitMessageInput struct {
	Topic     []byte
	Partition []byte
}

// NewPublishFollowMeRequestInitMessage builds a ZAP-encoded
// PublishFollowMeRequest.InitMessage message from in and returns the bytes.
func NewPublishFollowMeRequestInitMessage(in PublishFollowMeRequestInitMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(publishFollowMeRequestInitMessageSize)
	ob.SetBytes(publishFollowMeRequestInitMessageTopicOff, in.Topic)
	ob.SetBytes(publishFollowMeRequestInitMessagePartitionOff, in.Partition)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PublishFollowMeRequestFlushMessage ---
// (oneof message variant `flush` of PublishFollowMeRequest)

const (
	publishFollowMeRequestFlushMessageTsNsOff = 0 // int64
	publishFollowMeRequestFlushMessageSize    = 8
)

// PublishFollowMeRequestFlushMessage is a zero-copy view into a ZAP-encoded
// PublishFollowMeRequest.FlushMessage message.
type PublishFollowMeRequestFlushMessage struct{ o zap.Object }

// WrapPublishFollowMeRequestFlushMessage parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapPublishFollowMeRequestFlushMessage(b []byte) (PublishFollowMeRequestFlushMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PublishFollowMeRequestFlushMessage{}, err
	}
	return PublishFollowMeRequestFlushMessage{o: m.Root()}, nil
}

// TsNs reads the ts_ns field (proto field 1, int64).
func (t PublishFollowMeRequestFlushMessage) TsNs() int64 {
	return t.o.Int64(publishFollowMeRequestFlushMessageTsNsOff)
}

// PublishFollowMeRequestFlushMessageInput collects the field values for
// NewPublishFollowMeRequestFlushMessage.
type PublishFollowMeRequestFlushMessageInput struct {
	TsNs int64
}

// NewPublishFollowMeRequestFlushMessage builds a ZAP-encoded
// PublishFollowMeRequest.FlushMessage message from in and returns the bytes.
func NewPublishFollowMeRequestFlushMessage(in PublishFollowMeRequestFlushMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(publishFollowMeRequestFlushMessageSize)
	ob.SetInt64(publishFollowMeRequestFlushMessageTsNsOff, in.TsNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PublishFollowMeRequestCloseMessage ---
// (oneof message variant `close` of PublishFollowMeRequest)

const publishFollowMeRequestCloseMessageSize = 0

// PublishFollowMeRequestCloseMessage is a zero-copy view into a ZAP-encoded
// PublishFollowMeRequest.CloseMessage message. It carries no fields.
type PublishFollowMeRequestCloseMessage struct{ o zap.Object }

// WrapPublishFollowMeRequestCloseMessage parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapPublishFollowMeRequestCloseMessage(b []byte) (PublishFollowMeRequestCloseMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PublishFollowMeRequestCloseMessage{}, err
	}
	return PublishFollowMeRequestCloseMessage{o: m.Root()}, nil
}

// PublishFollowMeRequestCloseMessageInput collects the field values for
// NewPublishFollowMeRequestCloseMessage. It has no fields.
type PublishFollowMeRequestCloseMessageInput struct{}

// NewPublishFollowMeRequestCloseMessage builds a ZAP-encoded
// PublishFollowMeRequest.CloseMessage message and returns the bytes.
func NewPublishFollowMeRequestCloseMessage(in PublishFollowMeRequestCloseMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(publishFollowMeRequestCloseMessageSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PublishFollowMeRequest ---

// PublishFollowMeRequest.message oneof tags (proto field numbers, stable wire
// ids). Zero means no variant is set. The selected variant message is embedded
// as a nested object pointed to by message_value.
const (
	PublishFollowMeRequestMessageNone  uint32 = 0
	PublishFollowMeRequestMessageInit  uint32 = 1
	PublishFollowMeRequestMessageData  uint32 = 2
	PublishFollowMeRequestMessageFlush uint32 = 3
	PublishFollowMeRequestMessageClose uint32 = 4
)

const (
	publishFollowMeRequestMessageWhichOff = 0 // oneof tag (uint32)
	publishFollowMeRequestMessageValueOff = 4 // embedded variant (obj ptr)
	publishFollowMeRequestSize            = 8
)

// PublishFollowMeRequest is a zero-copy view into a ZAP-encoded
// PublishFollowMeRequest message.
type PublishFollowMeRequest struct{ o zap.Object }

// WrapPublishFollowMeRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapPublishFollowMeRequest(b []byte) (PublishFollowMeRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PublishFollowMeRequest{}, err
	}
	return PublishFollowMeRequest{o: m.Root()}, nil
}

// WhichMessage returns the set oneof tag (PublishFollowMeRequestMessage*
// constants).
func (t PublishFollowMeRequest) WhichMessage() uint32 {
	return t.o.Uint32(publishFollowMeRequestMessageWhichOff)
}

// Init returns the embedded init variant.
func (t PublishFollowMeRequest) Init() PublishFollowMeRequestInitMessage {
	return PublishFollowMeRequestInitMessage{o: t.o.Object(publishFollowMeRequestMessageValueOff)}
}

// Data returns the embedded data variant.
func (t PublishFollowMeRequest) Data() DataMessage {
	return DataMessage{o: t.o.Object(publishFollowMeRequestMessageValueOff)}
}

// Flush returns the embedded flush variant.
func (t PublishFollowMeRequest) Flush() PublishFollowMeRequestFlushMessage {
	return PublishFollowMeRequestFlushMessage{o: t.o.Object(publishFollowMeRequestMessageValueOff)}
}

// Close returns the embedded close variant.
func (t PublishFollowMeRequest) Close() PublishFollowMeRequestCloseMessage {
	return PublishFollowMeRequestCloseMessage{o: t.o.Object(publishFollowMeRequestMessageValueOff)}
}

// PublishFollowMeRequestInput collects the field values for
// NewPublishFollowMeRequest. MessageWhich selects the oneof variant; MessageValue
// is that variant's standalone buffer.
type PublishFollowMeRequestInput struct {
	MessageWhich uint32
	MessageValue []byte
}

// NewPublishFollowMeRequest builds a ZAP-encoded PublishFollowMeRequest message
// from in and returns the bytes.
func NewPublishFollowMeRequest(in PublishFollowMeRequestInput) []byte {
	b := zap.NewBuilder(512)
	valueOff := 0
	if len(in.MessageValue) > 0 {
		valueOff = embedMessage(b, in.MessageValue)
	}
	ob := b.StartObject(publishFollowMeRequestSize)
	ob.SetUint32(publishFollowMeRequestMessageWhichOff, in.MessageWhich)
	ob.SetObject(publishFollowMeRequestMessageValueOff, valueOff)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PublishFollowMeResponse ---

const (
	publishFollowMeResponseAckTsNsOff = 0 // int64
	publishFollowMeResponseSize       = 8
)

// PublishFollowMeResponse is a zero-copy view into a ZAP-encoded
// PublishFollowMeResponse message.
type PublishFollowMeResponse struct{ o zap.Object }

// WrapPublishFollowMeResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapPublishFollowMeResponse(b []byte) (PublishFollowMeResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PublishFollowMeResponse{}, err
	}
	return PublishFollowMeResponse{o: m.Root()}, nil
}

// AckTsNs reads the ack_ts_ns field (proto field 1, int64).
func (t PublishFollowMeResponse) AckTsNs() int64 {
	return t.o.Int64(publishFollowMeResponseAckTsNsOff)
}

// PublishFollowMeResponseInput collects the field values for
// NewPublishFollowMeResponse.
type PublishFollowMeResponseInput struct {
	AckTsNs int64
}

// NewPublishFollowMeResponse builds a ZAP-encoded PublishFollowMeResponse
// message from in and returns the bytes.
func NewPublishFollowMeResponse(in PublishFollowMeResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(publishFollowMeResponseSize)
	ob.SetInt64(publishFollowMeResponseAckTsNsOff, in.AckTsNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscribeMessageRequestInitMessage ---
// (oneof message variant `init` of SubscribeMessageRequest)

const (
	subscribeMessageRequestInitMessageConsumerGroupOff     = 0  // string
	subscribeMessageRequestInitMessageConsumerIdOff        = 8  // string
	subscribeMessageRequestInitMessageClientIdOff          = 16 // string
	subscribeMessageRequestInitMessageTopicOff             = 24 // message schema_pb.Topic
	subscribeMessageRequestInitMessagePartitionOffsetOff   = 32 // message schema_pb.PartitionOffset
	subscribeMessageRequestInitMessageOffsetTypeOff        = 40 // enum schema_pb.OffsetType (uint32)
	subscribeMessageRequestInitMessageFilterOff            = 48 // string
	subscribeMessageRequestInitMessageFollowerBrokerOff    = 56 // string
	subscribeMessageRequestInitMessageSlidingWindowSizeOff = 64 // int32
	subscribeMessageRequestInitMessageSize                 = 72
)

// SubscribeMessageRequestInitMessage is a zero-copy view into a ZAP-encoded
// SubscribeMessageRequest.InitMessage message.
type SubscribeMessageRequestInitMessage struct{ o zap.Object }

// WrapSubscribeMessageRequestInitMessage parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapSubscribeMessageRequestInitMessage(b []byte) (SubscribeMessageRequestInitMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscribeMessageRequestInitMessage{}, err
	}
	return SubscribeMessageRequestInitMessage{o: m.Root()}, nil
}

// ConsumerGroup reads the consumer_group field (proto field 1, string).
func (t SubscribeMessageRequestInitMessage) ConsumerGroup() string {
	return t.o.Text(subscribeMessageRequestInitMessageConsumerGroupOff)
}

// ConsumerId reads the consumer_id field (proto field 2, string).
func (t SubscribeMessageRequestInitMessage) ConsumerId() string {
	return t.o.Text(subscribeMessageRequestInitMessageConsumerIdOff)
}

// ClientId reads the client_id field (proto field 3, string).
func (t SubscribeMessageRequestInitMessage) ClientId() string {
	return t.o.Text(subscribeMessageRequestInitMessageClientIdOff)
}

// Topic reads the topic field (proto field 4, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t SubscribeMessageRequestInitMessage) Topic() []byte {
	return t.o.Bytes(subscribeMessageRequestInitMessageTopicOff)
}

// PartitionOffset reads the partition_offset field (proto field 5, message
// schema_pb.PartitionOffset) as the child's raw ZAP buffer; nil when absent.
// Re-wrap with mq_schemawire.WrapPartitionOffset.
func (t SubscribeMessageRequestInitMessage) PartitionOffset() []byte {
	return t.o.Bytes(subscribeMessageRequestInitMessagePartitionOffsetOff)
}

// OffsetType reads the offset_type field (proto field 6, enum
// schema_pb.OffsetType). See mq_schemawire.OffsetType* for the values.
func (t SubscribeMessageRequestInitMessage) OffsetType() uint32 {
	return t.o.Uint32(subscribeMessageRequestInitMessageOffsetTypeOff)
}

// Filter reads the filter field (proto field 10, string).
func (t SubscribeMessageRequestInitMessage) Filter() string {
	return t.o.Text(subscribeMessageRequestInitMessageFilterOff)
}

// FollowerBroker reads the follower_broker field (proto field 11, string).
func (t SubscribeMessageRequestInitMessage) FollowerBroker() string {
	return t.o.Text(subscribeMessageRequestInitMessageFollowerBrokerOff)
}

// SlidingWindowSize reads the sliding_window_size field (proto field 12, int32).
func (t SubscribeMessageRequestInitMessage) SlidingWindowSize() int32 {
	return t.o.Int32(subscribeMessageRequestInitMessageSlidingWindowSizeOff)
}

// SubscribeMessageRequestInitMessageInput collects the field values for
// NewSubscribeMessageRequestInitMessage. Topic and PartitionOffset are each the
// child message's own ZAP buffer (from mq_schemawire.NewTopic /
// NewPartitionOffset).
type SubscribeMessageRequestInitMessageInput struct {
	ConsumerGroup     string
	ConsumerId        string
	ClientId          string
	Topic             []byte
	PartitionOffset   []byte
	OffsetType        uint32
	Filter            string
	FollowerBroker    string
	SlidingWindowSize int32
}

// NewSubscribeMessageRequestInitMessage builds a ZAP-encoded
// SubscribeMessageRequest.InitMessage message from in and returns the bytes.
func NewSubscribeMessageRequestInitMessage(in SubscribeMessageRequestInitMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(subscribeMessageRequestInitMessageSize)
	ob.SetText(subscribeMessageRequestInitMessageConsumerGroupOff, in.ConsumerGroup)
	ob.SetText(subscribeMessageRequestInitMessageConsumerIdOff, in.ConsumerId)
	ob.SetText(subscribeMessageRequestInitMessageClientIdOff, in.ClientId)
	ob.SetBytes(subscribeMessageRequestInitMessageTopicOff, in.Topic)
	ob.SetBytes(subscribeMessageRequestInitMessagePartitionOffsetOff, in.PartitionOffset)
	ob.SetUint32(subscribeMessageRequestInitMessageOffsetTypeOff, in.OffsetType)
	ob.SetText(subscribeMessageRequestInitMessageFilterOff, in.Filter)
	ob.SetText(subscribeMessageRequestInitMessageFollowerBrokerOff, in.FollowerBroker)
	ob.SetInt32(subscribeMessageRequestInitMessageSlidingWindowSizeOff, in.SlidingWindowSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscribeMessageRequestAckMessage ---
// (oneof message variant `ack` of SubscribeMessageRequest)

const (
	subscribeMessageRequestAckMessageTsNsOff = 0 // int64
	subscribeMessageRequestAckMessageKeyOff  = 8 // bytes
	subscribeMessageRequestAckMessageSize    = 16
)

// SubscribeMessageRequestAckMessage is a zero-copy view into a ZAP-encoded
// SubscribeMessageRequest.AckMessage message.
type SubscribeMessageRequestAckMessage struct{ o zap.Object }

// WrapSubscribeMessageRequestAckMessage parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapSubscribeMessageRequestAckMessage(b []byte) (SubscribeMessageRequestAckMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscribeMessageRequestAckMessage{}, err
	}
	return SubscribeMessageRequestAckMessage{o: m.Root()}, nil
}

// TsNs reads the ts_ns field (proto field 1, int64).
func (t SubscribeMessageRequestAckMessage) TsNs() int64 {
	return t.o.Int64(subscribeMessageRequestAckMessageTsNsOff)
}

// Key reads the key field (proto field 2, bytes).
func (t SubscribeMessageRequestAckMessage) Key() []byte {
	return t.o.Bytes(subscribeMessageRequestAckMessageKeyOff)
}

// SubscribeMessageRequestAckMessageInput collects the field values for
// NewSubscribeMessageRequestAckMessage.
type SubscribeMessageRequestAckMessageInput struct {
	TsNs int64
	Key  []byte
}

// NewSubscribeMessageRequestAckMessage builds a ZAP-encoded
// SubscribeMessageRequest.AckMessage message from in and returns the bytes.
func NewSubscribeMessageRequestAckMessage(in SubscribeMessageRequestAckMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(subscribeMessageRequestAckMessageSize)
	ob.SetInt64(subscribeMessageRequestAckMessageTsNsOff, in.TsNs)
	ob.SetBytes(subscribeMessageRequestAckMessageKeyOff, in.Key)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscribeMessageRequestSeekMessage ---
// (oneof message variant `seek` of SubscribeMessageRequest)

const (
	subscribeMessageRequestSeekMessageOffsetOff     = 0 // int64
	subscribeMessageRequestSeekMessageOffsetTypeOff = 8 // enum schema_pb.OffsetType (uint32)
	subscribeMessageRequestSeekMessageSize          = 16
)

// SubscribeMessageRequestSeekMessage is a zero-copy view into a ZAP-encoded
// SubscribeMessageRequest.SeekMessage message.
type SubscribeMessageRequestSeekMessage struct{ o zap.Object }

// WrapSubscribeMessageRequestSeekMessage parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapSubscribeMessageRequestSeekMessage(b []byte) (SubscribeMessageRequestSeekMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscribeMessageRequestSeekMessage{}, err
	}
	return SubscribeMessageRequestSeekMessage{o: m.Root()}, nil
}

// Offset reads the offset field (proto field 1, int64).
func (t SubscribeMessageRequestSeekMessage) Offset() int64 {
	return t.o.Int64(subscribeMessageRequestSeekMessageOffsetOff)
}

// OffsetType reads the offset_type field (proto field 2, enum
// schema_pb.OffsetType). See mq_schemawire.OffsetType* for the values.
func (t SubscribeMessageRequestSeekMessage) OffsetType() uint32 {
	return t.o.Uint32(subscribeMessageRequestSeekMessageOffsetTypeOff)
}

// SubscribeMessageRequestSeekMessageInput collects the field values for
// NewSubscribeMessageRequestSeekMessage.
type SubscribeMessageRequestSeekMessageInput struct {
	Offset     int64
	OffsetType uint32
}

// NewSubscribeMessageRequestSeekMessage builds a ZAP-encoded
// SubscribeMessageRequest.SeekMessage message from in and returns the bytes.
func NewSubscribeMessageRequestSeekMessage(in SubscribeMessageRequestSeekMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(subscribeMessageRequestSeekMessageSize)
	ob.SetInt64(subscribeMessageRequestSeekMessageOffsetOff, in.Offset)
	ob.SetUint32(subscribeMessageRequestSeekMessageOffsetTypeOff, in.OffsetType)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscribeMessageRequest ---

// SubscribeMessageRequest.message oneof tags (proto field numbers, stable wire
// ids). Zero means no variant is set. The selected variant message is embedded
// as a nested object pointed to by message_value.
const (
	SubscribeMessageRequestMessageNone uint32 = 0
	SubscribeMessageRequestMessageInit uint32 = 1
	SubscribeMessageRequestMessageAck  uint32 = 2
	SubscribeMessageRequestMessageSeek uint32 = 3
)

const (
	subscribeMessageRequestMessageWhichOff = 0 // oneof tag (uint32)
	subscribeMessageRequestMessageValueOff = 4 // embedded variant (obj ptr)
	subscribeMessageRequestSize            = 8
)

// SubscribeMessageRequest is a zero-copy view into a ZAP-encoded
// SubscribeMessageRequest message.
type SubscribeMessageRequest struct{ o zap.Object }

// WrapSubscribeMessageRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapSubscribeMessageRequest(b []byte) (SubscribeMessageRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscribeMessageRequest{}, err
	}
	return SubscribeMessageRequest{o: m.Root()}, nil
}

// WhichMessage returns the set oneof tag (SubscribeMessageRequestMessage*
// constants).
func (t SubscribeMessageRequest) WhichMessage() uint32 {
	return t.o.Uint32(subscribeMessageRequestMessageWhichOff)
}

// Init returns the embedded init variant.
func (t SubscribeMessageRequest) Init() SubscribeMessageRequestInitMessage {
	return SubscribeMessageRequestInitMessage{o: t.o.Object(subscribeMessageRequestMessageValueOff)}
}

// Ack returns the embedded ack variant.
func (t SubscribeMessageRequest) Ack() SubscribeMessageRequestAckMessage {
	return SubscribeMessageRequestAckMessage{o: t.o.Object(subscribeMessageRequestMessageValueOff)}
}

// Seek returns the embedded seek variant.
func (t SubscribeMessageRequest) Seek() SubscribeMessageRequestSeekMessage {
	return SubscribeMessageRequestSeekMessage{o: t.o.Object(subscribeMessageRequestMessageValueOff)}
}

// SubscribeMessageRequestInput collects the field values for
// NewSubscribeMessageRequest. MessageWhich selects the oneof variant; MessageValue
// is that variant's standalone buffer.
type SubscribeMessageRequestInput struct {
	MessageWhich uint32
	MessageValue []byte
}

// NewSubscribeMessageRequest builds a ZAP-encoded SubscribeMessageRequest
// message from in and returns the bytes.
func NewSubscribeMessageRequest(in SubscribeMessageRequestInput) []byte {
	b := zap.NewBuilder(512)
	valueOff := 0
	if len(in.MessageValue) > 0 {
		valueOff = embedMessage(b, in.MessageValue)
	}
	ob := b.StartObject(subscribeMessageRequestSize)
	ob.SetUint32(subscribeMessageRequestMessageWhichOff, in.MessageWhich)
	ob.SetObject(subscribeMessageRequestMessageValueOff, valueOff)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscribeMessageResponseSubscribeCtrlMessage ---
// (oneof message variant `ctrl` of SubscribeMessageResponse)

const (
	subscribeMessageResponseSubscribeCtrlMessageErrorOff         = 0 // string
	subscribeMessageResponseSubscribeCtrlMessageIsEndOfStreamOff = 8 // bool
	subscribeMessageResponseSubscribeCtrlMessageIsEndOfTopicOff  = 9 // bool
	subscribeMessageResponseSubscribeCtrlMessageSize             = 16
)

// SubscribeMessageResponseSubscribeCtrlMessage is a zero-copy view into a
// ZAP-encoded SubscribeMessageResponse.SubscribeCtrlMessage message.
type SubscribeMessageResponseSubscribeCtrlMessage struct{ o zap.Object }

// WrapSubscribeMessageResponseSubscribeCtrlMessage parses b and returns a typed
// view. Returns an error if the wire-level checks fail.
func WrapSubscribeMessageResponseSubscribeCtrlMessage(b []byte) (SubscribeMessageResponseSubscribeCtrlMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscribeMessageResponseSubscribeCtrlMessage{}, err
	}
	return SubscribeMessageResponseSubscribeCtrlMessage{o: m.Root()}, nil
}

// Error reads the error field (proto field 1, string).
func (t SubscribeMessageResponseSubscribeCtrlMessage) Error() string {
	return t.o.Text(subscribeMessageResponseSubscribeCtrlMessageErrorOff)
}

// IsEndOfStream reads the is_end_of_stream field (proto field 2, bool).
func (t SubscribeMessageResponseSubscribeCtrlMessage) IsEndOfStream() bool {
	return t.o.Bool(subscribeMessageResponseSubscribeCtrlMessageIsEndOfStreamOff)
}

// IsEndOfTopic reads the is_end_of_topic field (proto field 3, bool).
func (t SubscribeMessageResponseSubscribeCtrlMessage) IsEndOfTopic() bool {
	return t.o.Bool(subscribeMessageResponseSubscribeCtrlMessageIsEndOfTopicOff)
}

// SubscribeMessageResponseSubscribeCtrlMessageInput collects the field values
// for NewSubscribeMessageResponseSubscribeCtrlMessage.
type SubscribeMessageResponseSubscribeCtrlMessageInput struct {
	Error         string
	IsEndOfStream bool
	IsEndOfTopic  bool
}

// NewSubscribeMessageResponseSubscribeCtrlMessage builds a ZAP-encoded
// SubscribeMessageResponse.SubscribeCtrlMessage message from in.
func NewSubscribeMessageResponseSubscribeCtrlMessage(in SubscribeMessageResponseSubscribeCtrlMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(subscribeMessageResponseSubscribeCtrlMessageSize)
	ob.SetText(subscribeMessageResponseSubscribeCtrlMessageErrorOff, in.Error)
	ob.SetBool(subscribeMessageResponseSubscribeCtrlMessageIsEndOfStreamOff, in.IsEndOfStream)
	ob.SetBool(subscribeMessageResponseSubscribeCtrlMessageIsEndOfTopicOff, in.IsEndOfTopic)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscribeMessageResponse ---

// SubscribeMessageResponse.message oneof tags (proto field numbers, stable wire
// ids). Zero means no variant is set. The selected variant message is embedded
// as a nested object pointed to by message_value.
const (
	SubscribeMessageResponseMessageNone uint32 = 0
	SubscribeMessageResponseMessageCtrl uint32 = 1
	SubscribeMessageResponseMessageData uint32 = 2
)

const (
	subscribeMessageResponseMessageWhichOff = 0 // oneof tag (uint32)
	subscribeMessageResponseMessageValueOff = 4 // embedded variant (obj ptr)
	subscribeMessageResponseSize            = 8
)

// SubscribeMessageResponse is a zero-copy view into a ZAP-encoded
// SubscribeMessageResponse message.
type SubscribeMessageResponse struct{ o zap.Object }

// WrapSubscribeMessageResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapSubscribeMessageResponse(b []byte) (SubscribeMessageResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscribeMessageResponse{}, err
	}
	return SubscribeMessageResponse{o: m.Root()}, nil
}

// WhichMessage returns the set oneof tag (SubscribeMessageResponseMessage*
// constants).
func (t SubscribeMessageResponse) WhichMessage() uint32 {
	return t.o.Uint32(subscribeMessageResponseMessageWhichOff)
}

// Ctrl returns the embedded ctrl variant.
func (t SubscribeMessageResponse) Ctrl() SubscribeMessageResponseSubscribeCtrlMessage {
	return SubscribeMessageResponseSubscribeCtrlMessage{o: t.o.Object(subscribeMessageResponseMessageValueOff)}
}

// Data returns the embedded data variant.
func (t SubscribeMessageResponse) Data() DataMessage {
	return DataMessage{o: t.o.Object(subscribeMessageResponseMessageValueOff)}
}

// SubscribeMessageResponseInput collects the field values for
// NewSubscribeMessageResponse. MessageWhich selects the oneof variant;
// MessageValue is that variant's standalone buffer.
type SubscribeMessageResponseInput struct {
	MessageWhich uint32
	MessageValue []byte
}

// NewSubscribeMessageResponse builds a ZAP-encoded SubscribeMessageResponse
// message from in and returns the bytes.
func NewSubscribeMessageResponse(in SubscribeMessageResponseInput) []byte {
	b := zap.NewBuilder(512)
	valueOff := 0
	if len(in.MessageValue) > 0 {
		valueOff = embedMessage(b, in.MessageValue)
	}
	ob := b.StartObject(subscribeMessageResponseSize)
	ob.SetUint32(subscribeMessageResponseMessageWhichOff, in.MessageWhich)
	ob.SetObject(subscribeMessageResponseMessageValueOff, valueOff)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscribeFollowMeRequestInitMessage ---
// (oneof message variant `init` of SubscribeFollowMeRequest)

const (
	subscribeFollowMeRequestInitMessageTopicOff         = 0  // message schema_pb.Topic
	subscribeFollowMeRequestInitMessagePartitionOff     = 8  // message schema_pb.Partition
	subscribeFollowMeRequestInitMessageConsumerGroupOff = 16 // string
	subscribeFollowMeRequestInitMessageSize             = 24
)

// SubscribeFollowMeRequestInitMessage is a zero-copy view into a ZAP-encoded
// SubscribeFollowMeRequest.InitMessage message.
type SubscribeFollowMeRequestInitMessage struct{ o zap.Object }

// WrapSubscribeFollowMeRequestInitMessage parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapSubscribeFollowMeRequestInitMessage(b []byte) (SubscribeFollowMeRequestInitMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscribeFollowMeRequestInitMessage{}, err
	}
	return SubscribeFollowMeRequestInitMessage{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t SubscribeFollowMeRequestInitMessage) Topic() []byte {
	return t.o.Bytes(subscribeFollowMeRequestInitMessageTopicOff)
}

// Partition reads the partition field (proto field 2, message
// schema_pb.Partition) as the child's raw ZAP buffer; nil when absent. Re-wrap
// with mq_schemawire.WrapPartition.
func (t SubscribeFollowMeRequestInitMessage) Partition() []byte {
	return t.o.Bytes(subscribeFollowMeRequestInitMessagePartitionOff)
}

// ConsumerGroup reads the consumer_group field (proto field 3, string).
func (t SubscribeFollowMeRequestInitMessage) ConsumerGroup() string {
	return t.o.Text(subscribeFollowMeRequestInitMessageConsumerGroupOff)
}

// SubscribeFollowMeRequestInitMessageInput collects the field values for
// NewSubscribeFollowMeRequestInitMessage. Topic and Partition are each the child
// message's own ZAP buffer (from mq_schemawire.NewTopic / NewPartition).
type SubscribeFollowMeRequestInitMessageInput struct {
	Topic         []byte
	Partition     []byte
	ConsumerGroup string
}

// NewSubscribeFollowMeRequestInitMessage builds a ZAP-encoded
// SubscribeFollowMeRequest.InitMessage message from in and returns the bytes.
func NewSubscribeFollowMeRequestInitMessage(in SubscribeFollowMeRequestInitMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(subscribeFollowMeRequestInitMessageSize)
	ob.SetBytes(subscribeFollowMeRequestInitMessageTopicOff, in.Topic)
	ob.SetBytes(subscribeFollowMeRequestInitMessagePartitionOff, in.Partition)
	ob.SetText(subscribeFollowMeRequestInitMessageConsumerGroupOff, in.ConsumerGroup)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscribeFollowMeRequestAckMessage ---
// (oneof message variant `ack` of SubscribeFollowMeRequest)

const (
	subscribeFollowMeRequestAckMessageTsNsOff = 0 // int64
	subscribeFollowMeRequestAckMessageSize    = 8
)

// SubscribeFollowMeRequestAckMessage is a zero-copy view into a ZAP-encoded
// SubscribeFollowMeRequest.AckMessage message.
type SubscribeFollowMeRequestAckMessage struct{ o zap.Object }

// WrapSubscribeFollowMeRequestAckMessage parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapSubscribeFollowMeRequestAckMessage(b []byte) (SubscribeFollowMeRequestAckMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscribeFollowMeRequestAckMessage{}, err
	}
	return SubscribeFollowMeRequestAckMessage{o: m.Root()}, nil
}

// TsNs reads the ts_ns field (proto field 1, int64).
func (t SubscribeFollowMeRequestAckMessage) TsNs() int64 {
	return t.o.Int64(subscribeFollowMeRequestAckMessageTsNsOff)
}

// SubscribeFollowMeRequestAckMessageInput collects the field values for
// NewSubscribeFollowMeRequestAckMessage.
type SubscribeFollowMeRequestAckMessageInput struct {
	TsNs int64
}

// NewSubscribeFollowMeRequestAckMessage builds a ZAP-encoded
// SubscribeFollowMeRequest.AckMessage message from in and returns the bytes.
func NewSubscribeFollowMeRequestAckMessage(in SubscribeFollowMeRequestAckMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(subscribeFollowMeRequestAckMessageSize)
	ob.SetInt64(subscribeFollowMeRequestAckMessageTsNsOff, in.TsNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscribeFollowMeRequestCloseMessage ---
// (oneof message variant `close` of SubscribeFollowMeRequest)

const subscribeFollowMeRequestCloseMessageSize = 0

// SubscribeFollowMeRequestCloseMessage is a zero-copy view into a ZAP-encoded
// SubscribeFollowMeRequest.CloseMessage message. It carries no fields.
type SubscribeFollowMeRequestCloseMessage struct{ o zap.Object }

// WrapSubscribeFollowMeRequestCloseMessage parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapSubscribeFollowMeRequestCloseMessage(b []byte) (SubscribeFollowMeRequestCloseMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscribeFollowMeRequestCloseMessage{}, err
	}
	return SubscribeFollowMeRequestCloseMessage{o: m.Root()}, nil
}

// SubscribeFollowMeRequestCloseMessageInput collects the field values for
// NewSubscribeFollowMeRequestCloseMessage. It has no fields.
type SubscribeFollowMeRequestCloseMessageInput struct{}

// NewSubscribeFollowMeRequestCloseMessage builds a ZAP-encoded
// SubscribeFollowMeRequest.CloseMessage message and returns the bytes.
func NewSubscribeFollowMeRequestCloseMessage(in SubscribeFollowMeRequestCloseMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(subscribeFollowMeRequestCloseMessageSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscribeFollowMeRequest ---

// SubscribeFollowMeRequest.message oneof tags (proto field numbers, stable wire
// ids). Zero means no variant is set. The selected variant message is embedded
// as a nested object pointed to by message_value.
const (
	SubscribeFollowMeRequestMessageNone  uint32 = 0
	SubscribeFollowMeRequestMessageInit  uint32 = 1
	SubscribeFollowMeRequestMessageAck   uint32 = 2
	SubscribeFollowMeRequestMessageClose uint32 = 3
)

const (
	subscribeFollowMeRequestMessageWhichOff = 0 // oneof tag (uint32)
	subscribeFollowMeRequestMessageValueOff = 4 // embedded variant (obj ptr)
	subscribeFollowMeRequestSize            = 8
)

// SubscribeFollowMeRequest is a zero-copy view into a ZAP-encoded
// SubscribeFollowMeRequest message.
type SubscribeFollowMeRequest struct{ o zap.Object }

// WrapSubscribeFollowMeRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapSubscribeFollowMeRequest(b []byte) (SubscribeFollowMeRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscribeFollowMeRequest{}, err
	}
	return SubscribeFollowMeRequest{o: m.Root()}, nil
}

// WhichMessage returns the set oneof tag (SubscribeFollowMeRequestMessage*
// constants).
func (t SubscribeFollowMeRequest) WhichMessage() uint32 {
	return t.o.Uint32(subscribeFollowMeRequestMessageWhichOff)
}

// Init returns the embedded init variant.
func (t SubscribeFollowMeRequest) Init() SubscribeFollowMeRequestInitMessage {
	return SubscribeFollowMeRequestInitMessage{o: t.o.Object(subscribeFollowMeRequestMessageValueOff)}
}

// Ack returns the embedded ack variant.
func (t SubscribeFollowMeRequest) Ack() SubscribeFollowMeRequestAckMessage {
	return SubscribeFollowMeRequestAckMessage{o: t.o.Object(subscribeFollowMeRequestMessageValueOff)}
}

// Close returns the embedded close variant.
func (t SubscribeFollowMeRequest) Close() SubscribeFollowMeRequestCloseMessage {
	return SubscribeFollowMeRequestCloseMessage{o: t.o.Object(subscribeFollowMeRequestMessageValueOff)}
}

// SubscribeFollowMeRequestInput collects the field values for
// NewSubscribeFollowMeRequest. MessageWhich selects the oneof variant;
// MessageValue is that variant's standalone buffer.
type SubscribeFollowMeRequestInput struct {
	MessageWhich uint32
	MessageValue []byte
}

// NewSubscribeFollowMeRequest builds a ZAP-encoded SubscribeFollowMeRequest
// message from in and returns the bytes.
func NewSubscribeFollowMeRequest(in SubscribeFollowMeRequestInput) []byte {
	b := zap.NewBuilder(512)
	valueOff := 0
	if len(in.MessageValue) > 0 {
		valueOff = embedMessage(b, in.MessageValue)
	}
	ob := b.StartObject(subscribeFollowMeRequestSize)
	ob.SetUint32(subscribeFollowMeRequestMessageWhichOff, in.MessageWhich)
	ob.SetObject(subscribeFollowMeRequestMessageValueOff, valueOff)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscribeFollowMeResponse ---

const (
	subscribeFollowMeResponseAckTsNsOff = 0 // int64
	subscribeFollowMeResponseSize       = 8
)

// SubscribeFollowMeResponse is a zero-copy view into a ZAP-encoded
// SubscribeFollowMeResponse message.
type SubscribeFollowMeResponse struct{ o zap.Object }

// WrapSubscribeFollowMeResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapSubscribeFollowMeResponse(b []byte) (SubscribeFollowMeResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscribeFollowMeResponse{}, err
	}
	return SubscribeFollowMeResponse{o: m.Root()}, nil
}

// AckTsNs reads the ack_ts_ns field (proto field 1, int64).
func (t SubscribeFollowMeResponse) AckTsNs() int64 {
	return t.o.Int64(subscribeFollowMeResponseAckTsNsOff)
}

// SubscribeFollowMeResponseInput collects the field values for
// NewSubscribeFollowMeResponse.
type SubscribeFollowMeResponseInput struct {
	AckTsNs int64
}

// NewSubscribeFollowMeResponse builds a ZAP-encoded SubscribeFollowMeResponse
// message from in and returns the bytes.
func NewSubscribeFollowMeResponse(in SubscribeFollowMeResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(subscribeFollowMeResponseSize)
	ob.SetInt64(subscribeFollowMeResponseAckTsNsOff, in.AckTsNs)
	ob.FinishAsRoot()
	return b.Finish()
}
