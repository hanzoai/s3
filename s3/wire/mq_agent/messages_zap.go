// Code generated from mq_agent.proto; DO NOT EDIT.

package mq_agentwire

import (
	zap "github.com/zap-proto/go"
)

// Messages for the HanzoMessagingAgent service. Singular message fields
// (schema_pb.Topic, schema_pb.RecordType, schema_pb.RecordValue,
// schema_pb.PartitionOffset) are stored as the child message's own ZAP buffer
// (8-byte tail pointer, read with Bytes / written with SetBytes); build those
// children with the mq_schema wire package and re-wrap them with that package's
// Wrap* readers. Repeated message fields are lists of out-of-line objects.
// schema_pb.OffsetType is a uint32 scalar (see mq_schemawire.OffsetType*).

// --- StartPublishSessionRequest ---

const (
	startPublishSessionRequestTopicOff          = 0  // message schema_pb.Topic
	startPublishSessionRequestPartitionCountOff = 8  // int32
	startPublishSessionRequestRecordTypeOff     = 16 // message schema_pb.RecordType
	startPublishSessionRequestPublisherNameOff  = 24 // string
	startPublishSessionRequestSize              = 32
)

// StartPublishSessionRequest is a zero-copy view into a ZAP-encoded
// StartPublishSessionRequest message.
type StartPublishSessionRequest struct{ o zap.Object }

// WrapStartPublishSessionRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapStartPublishSessionRequest(b []byte) (StartPublishSessionRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return StartPublishSessionRequest{}, err
	}
	return StartPublishSessionRequest{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t StartPublishSessionRequest) Topic() []byte {
	return t.o.Bytes(startPublishSessionRequestTopicOff)
}

// PartitionCount reads the partition_count field (proto field 2, int32).
func (t StartPublishSessionRequest) PartitionCount() int32 {
	return t.o.Int32(startPublishSessionRequestPartitionCountOff)
}

// RecordType reads the record_type field (proto field 3, message
// schema_pb.RecordType) as the child's raw ZAP buffer; nil when absent. Re-wrap
// with mq_schemawire.WrapRecordType.
func (t StartPublishSessionRequest) RecordType() []byte {
	return t.o.Bytes(startPublishSessionRequestRecordTypeOff)
}

// PublisherName reads the publisher_name field (proto field 4, string).
func (t StartPublishSessionRequest) PublisherName() string {
	return t.o.Text(startPublishSessionRequestPublisherNameOff)
}

// StartPublishSessionRequestInput collects the field values for
// NewStartPublishSessionRequest. Topic and RecordType are each the child
// message's own ZAP buffer (from mq_schemawire.NewTopic / NewRecordType); nil
// encodes an absent child.
type StartPublishSessionRequestInput struct {
	Topic          []byte
	PartitionCount int32
	RecordType     []byte
	PublisherName  string
}

// NewStartPublishSessionRequest builds a ZAP-encoded StartPublishSessionRequest
// message from in and returns the bytes.
func NewStartPublishSessionRequest(in StartPublishSessionRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(startPublishSessionRequestSize)
	ob.SetBytes(startPublishSessionRequestTopicOff, in.Topic)
	ob.SetInt32(startPublishSessionRequestPartitionCountOff, in.PartitionCount)
	ob.SetBytes(startPublishSessionRequestRecordTypeOff, in.RecordType)
	ob.SetText(startPublishSessionRequestPublisherNameOff, in.PublisherName)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- StartPublishSessionResponse ---

const (
	startPublishSessionResponseErrorOff     = 0 // string
	startPublishSessionResponseSessionIdOff = 8 // int64
	startPublishSessionResponseSize         = 16
)

// StartPublishSessionResponse is a zero-copy view into a ZAP-encoded
// StartPublishSessionResponse message.
type StartPublishSessionResponse struct{ o zap.Object }

// WrapStartPublishSessionResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapStartPublishSessionResponse(b []byte) (StartPublishSessionResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return StartPublishSessionResponse{}, err
	}
	return StartPublishSessionResponse{o: m.Root()}, nil
}

// Error reads the error field (proto field 1, string).
func (t StartPublishSessionResponse) Error() string {
	return t.o.Text(startPublishSessionResponseErrorOff)
}

// SessionId reads the session_id field (proto field 2, int64).
func (t StartPublishSessionResponse) SessionId() int64 {
	return t.o.Int64(startPublishSessionResponseSessionIdOff)
}

// StartPublishSessionResponseInput collects the field values for
// NewStartPublishSessionResponse.
type StartPublishSessionResponseInput struct {
	Error     string
	SessionId int64
}

// NewStartPublishSessionResponse builds a ZAP-encoded
// StartPublishSessionResponse message from in and returns the bytes.
func NewStartPublishSessionResponse(in StartPublishSessionResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(startPublishSessionResponseSize)
	ob.SetText(startPublishSessionResponseErrorOff, in.Error)
	ob.SetInt64(startPublishSessionResponseSessionIdOff, in.SessionId)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ClosePublishSessionRequest ---

const (
	closePublishSessionRequestSessionIdOff = 0 // int64
	closePublishSessionRequestSize         = 8
)

// ClosePublishSessionRequest is a zero-copy view into a ZAP-encoded
// ClosePublishSessionRequest message.
type ClosePublishSessionRequest struct{ o zap.Object }

// WrapClosePublishSessionRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapClosePublishSessionRequest(b []byte) (ClosePublishSessionRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ClosePublishSessionRequest{}, err
	}
	return ClosePublishSessionRequest{o: m.Root()}, nil
}

// SessionId reads the session_id field (proto field 1, int64).
func (t ClosePublishSessionRequest) SessionId() int64 {
	return t.o.Int64(closePublishSessionRequestSessionIdOff)
}

// ClosePublishSessionRequestInput collects the field values for
// NewClosePublishSessionRequest.
type ClosePublishSessionRequestInput struct {
	SessionId int64
}

// NewClosePublishSessionRequest builds a ZAP-encoded ClosePublishSessionRequest
// message from in and returns the bytes.
func NewClosePublishSessionRequest(in ClosePublishSessionRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(closePublishSessionRequestSize)
	ob.SetInt64(closePublishSessionRequestSessionIdOff, in.SessionId)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ClosePublishSessionResponse ---

const (
	closePublishSessionResponseErrorOff = 0 // string
	closePublishSessionResponseSize     = 8
)

// ClosePublishSessionResponse is a zero-copy view into a ZAP-encoded
// ClosePublishSessionResponse message.
type ClosePublishSessionResponse struct{ o zap.Object }

// WrapClosePublishSessionResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapClosePublishSessionResponse(b []byte) (ClosePublishSessionResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ClosePublishSessionResponse{}, err
	}
	return ClosePublishSessionResponse{o: m.Root()}, nil
}

// Error reads the error field (proto field 1, string).
func (t ClosePublishSessionResponse) Error() string {
	return t.o.Text(closePublishSessionResponseErrorOff)
}

// ClosePublishSessionResponseInput collects the field values for
// NewClosePublishSessionResponse.
type ClosePublishSessionResponseInput struct {
	Error string
}

// NewClosePublishSessionResponse builds a ZAP-encoded
// ClosePublishSessionResponse message from in and returns the bytes.
func NewClosePublishSessionResponse(in ClosePublishSessionResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(closePublishSessionResponseSize)
	ob.SetText(closePublishSessionResponseErrorOff, in.Error)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PublishRecordRequest ---

const (
	publishRecordRequestSessionIdOff = 0  // int64
	publishRecordRequestKeyOff       = 8  // bytes
	publishRecordRequestValueOff     = 16 // message schema_pb.RecordValue
	publishRecordRequestSize         = 24
)

// PublishRecordRequest is a zero-copy view into a ZAP-encoded
// PublishRecordRequest message.
type PublishRecordRequest struct{ o zap.Object }

// WrapPublishRecordRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapPublishRecordRequest(b []byte) (PublishRecordRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PublishRecordRequest{}, err
	}
	return PublishRecordRequest{o: m.Root()}, nil
}

// SessionId reads the session_id field (proto field 1, int64). session_id is
// required for the first record in a stream.
func (t PublishRecordRequest) SessionId() int64 {
	return t.o.Int64(publishRecordRequestSessionIdOff)
}

// Key reads the key field (proto field 2, bytes).
func (t PublishRecordRequest) Key() []byte {
	return t.o.Bytes(publishRecordRequestKeyOff)
}

// Value reads the value field (proto field 3, message schema_pb.RecordValue)
// as the child's raw ZAP buffer; nil when absent. Re-wrap with
// mq_schemawire.WrapRecordValue.
func (t PublishRecordRequest) Value() []byte {
	return t.o.Bytes(publishRecordRequestValueOff)
}

// PublishRecordRequestInput collects the field values for
// NewPublishRecordRequest. Value is the child message's own ZAP buffer (from
// mq_schemawire.NewRecordValue); nil encodes an absent value.
type PublishRecordRequestInput struct {
	SessionId int64
	Key       []byte
	Value     []byte
}

// NewPublishRecordRequest builds a ZAP-encoded PublishRecordRequest message
// from in and returns the bytes.
func NewPublishRecordRequest(in PublishRecordRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(publishRecordRequestSize)
	ob.SetInt64(publishRecordRequestSessionIdOff, in.SessionId)
	ob.SetBytes(publishRecordRequestKeyOff, in.Key)
	ob.SetBytes(publishRecordRequestValueOff, in.Value)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PublishRecordResponse ---

const (
	publishRecordResponseAckSequenceOff = 0  // int64
	publishRecordResponseErrorOff       = 8  // string
	publishRecordResponseBaseOffsetOff  = 16 // int64
	publishRecordResponseLastOffsetOff  = 24 // int64
	publishRecordResponseSize           = 32
)

// PublishRecordResponse is a zero-copy view into a ZAP-encoded
// PublishRecordResponse message.
type PublishRecordResponse struct{ o zap.Object }

// WrapPublishRecordResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapPublishRecordResponse(b []byte) (PublishRecordResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PublishRecordResponse{}, err
	}
	return PublishRecordResponse{o: m.Root()}, nil
}

// AckSequence reads the ack_sequence field (proto field 1, int64).
func (t PublishRecordResponse) AckSequence() int64 {
	return t.o.Int64(publishRecordResponseAckSequenceOff)
}

// Error reads the error field (proto field 2, string).
func (t PublishRecordResponse) Error() string {
	return t.o.Text(publishRecordResponseErrorOff)
}

// BaseOffset reads the base_offset field (proto field 3, int64). It is the
// first offset assigned to this batch.
func (t PublishRecordResponse) BaseOffset() int64 {
	return t.o.Int64(publishRecordResponseBaseOffsetOff)
}

// LastOffset reads the last_offset field (proto field 4, int64). It is the last
// offset assigned to this batch.
func (t PublishRecordResponse) LastOffset() int64 {
	return t.o.Int64(publishRecordResponseLastOffsetOff)
}

// PublishRecordResponseInput collects the field values for
// NewPublishRecordResponse.
type PublishRecordResponseInput struct {
	AckSequence int64
	Error       string
	BaseOffset  int64
	LastOffset  int64
}

// NewPublishRecordResponse builds a ZAP-encoded PublishRecordResponse message
// from in and returns the bytes.
func NewPublishRecordResponse(in PublishRecordResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(publishRecordResponseSize)
	ob.SetInt64(publishRecordResponseAckSequenceOff, in.AckSequence)
	ob.SetText(publishRecordResponseErrorOff, in.Error)
	ob.SetInt64(publishRecordResponseBaseOffsetOff, in.BaseOffset)
	ob.SetInt64(publishRecordResponseLastOffsetOff, in.LastOffset)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscribeRecordRequest.InitSubscribeRecordRequest (nested message) ---

const (
	initSubscribeRecordRequestConsumerGroupOff           = 0  // string
	initSubscribeRecordRequestConsumerGroupInstanceIdOff = 8  // string
	initSubscribeRecordRequestTopicOff                   = 16 // message schema_pb.Topic
	initSubscribeRecordRequestPartitionOffsetsOff        = 24 // repeated schema_pb.PartitionOffset
	initSubscribeRecordRequestOffsetTypeOff              = 32 // enum schema_pb.OffsetType (uint32)
	initSubscribeRecordRequestOffsetTsNsOff              = 40 // int64
	initSubscribeRecordRequestFilterOff                  = 48 // string
	initSubscribeRecordRequestMaxSubscribedPartitionsOff = 56 // int32
	initSubscribeRecordRequestSlidingWindowSizeOff       = 60 // int32
	initSubscribeRecordRequestSize                       = 64
)

// InitSubscribeRecordRequest is a zero-copy view into a ZAP-encoded
// SubscribeRecordRequest.InitSubscribeRecordRequest message.
type InitSubscribeRecordRequest struct{ o zap.Object }

// WrapInitSubscribeRecordRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapInitSubscribeRecordRequest(b []byte) (InitSubscribeRecordRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return InitSubscribeRecordRequest{}, err
	}
	return InitSubscribeRecordRequest{o: m.Root()}, nil
}

// ConsumerGroup reads the consumer_group field (proto field 1, string).
func (t InitSubscribeRecordRequest) ConsumerGroup() string {
	return t.o.Text(initSubscribeRecordRequestConsumerGroupOff)
}

// ConsumerGroupInstanceId reads the consumer_group_instance_id field (proto
// field 2, string).
func (t InitSubscribeRecordRequest) ConsumerGroupInstanceId() string {
	return t.o.Text(initSubscribeRecordRequestConsumerGroupInstanceIdOff)
}

// Topic reads the topic field (proto field 4, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t InitSubscribeRecordRequest) Topic() []byte {
	return t.o.Bytes(initSubscribeRecordRequestTopicOff)
}

// PartitionOffsetsLen reports the number of partition_offsets elements (proto
// field 5, repeated schema_pb.PartitionOffset).
func (t InitSubscribeRecordRequest) PartitionOffsetsLen() int {
	return t.o.List(initSubscribeRecordRequestPartitionOffsetsOff).Len()
}

// PartitionOffsetAt returns the i-th partition_offsets element as the child's
// raw ZAP buffer; nil when i is out of range. Re-wrap with
// mq_schemawire.WrapPartitionOffset.
func (t InitSubscribeRecordRequest) PartitionOffsetAt(i int) []byte {
	return t.o.List(initSubscribeRecordRequestPartitionOffsetsOff).BytesAt(i)
}

// OffsetType reads the offset_type field (proto field 6, enum
// schema_pb.OffsetType). See mq_schemawire.OffsetType* for the values.
func (t InitSubscribeRecordRequest) OffsetType() uint32 {
	return t.o.Uint32(initSubscribeRecordRequestOffsetTypeOff)
}

// OffsetTsNs reads the offset_ts_ns field (proto field 7, int64).
func (t InitSubscribeRecordRequest) OffsetTsNs() int64 {
	return t.o.Int64(initSubscribeRecordRequestOffsetTsNsOff)
}

// Filter reads the filter field (proto field 10, string).
func (t InitSubscribeRecordRequest) Filter() string {
	return t.o.Text(initSubscribeRecordRequestFilterOff)
}

// MaxSubscribedPartitions reads the max_subscribed_partitions field (proto
// field 11, int32).
func (t InitSubscribeRecordRequest) MaxSubscribedPartitions() int32 {
	return t.o.Int32(initSubscribeRecordRequestMaxSubscribedPartitionsOff)
}

// SlidingWindowSize reads the sliding_window_size field (proto field 12, int32).
func (t InitSubscribeRecordRequest) SlidingWindowSize() int32 {
	return t.o.Int32(initSubscribeRecordRequestSlidingWindowSizeOff)
}

// InitSubscribeRecordRequestInput collects the field values for
// NewInitSubscribeRecordRequest. Topic is the child message's own ZAP buffer
// (from mq_schemawire.NewTopic); nil encodes an absent topic. PartitionOffsets
// are each a PartitionOffset buffer (from mq_schemawire.NewPartitionOffset).
type InitSubscribeRecordRequestInput struct {
	ConsumerGroup           string
	ConsumerGroupInstanceId string
	Topic                   []byte
	PartitionOffsets        [][]byte
	OffsetType              uint32
	OffsetTsNs              int64
	Filter                  string
	MaxSubscribedPartitions int32
	SlidingWindowSize       int32
}

// NewInitSubscribeRecordRequest builds a ZAP-encoded
// InitSubscribeRecordRequest message from in and returns the bytes.
func NewInitSubscribeRecordRequest(in InitSubscribeRecordRequestInput) []byte {
	b := zap.NewBuilder(256)
	lb := b.StartList(0)
	for _, e := range in.PartitionOffsets {
		lb.AddObjectBytes(e)
	}
	listOff, listLen := lb.Finish()
	ob := b.StartObject(initSubscribeRecordRequestSize)
	ob.SetText(initSubscribeRecordRequestConsumerGroupOff, in.ConsumerGroup)
	ob.SetText(initSubscribeRecordRequestConsumerGroupInstanceIdOff, in.ConsumerGroupInstanceId)
	ob.SetBytes(initSubscribeRecordRequestTopicOff, in.Topic)
	ob.SetList(initSubscribeRecordRequestPartitionOffsetsOff, listOff, listLen)
	ob.SetUint32(initSubscribeRecordRequestOffsetTypeOff, in.OffsetType)
	ob.SetInt64(initSubscribeRecordRequestOffsetTsNsOff, in.OffsetTsNs)
	ob.SetText(initSubscribeRecordRequestFilterOff, in.Filter)
	ob.SetInt32(initSubscribeRecordRequestMaxSubscribedPartitionsOff, in.MaxSubscribedPartitions)
	ob.SetInt32(initSubscribeRecordRequestSlidingWindowSizeOff, in.SlidingWindowSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscribeRecordRequest ---

const (
	subscribeRecordRequestInitOff        = 0  // message InitSubscribeRecordRequest
	subscribeRecordRequestAckSequenceOff = 8  // int64
	subscribeRecordRequestAckKeyOff      = 16 // bytes
	subscribeRecordRequestSize           = 24
)

// SubscribeRecordRequest is a zero-copy view into a ZAP-encoded
// SubscribeRecordRequest message.
type SubscribeRecordRequest struct{ o zap.Object }

// WrapSubscribeRecordRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapSubscribeRecordRequest(b []byte) (SubscribeRecordRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscribeRecordRequest{}, err
	}
	return SubscribeRecordRequest{o: m.Root()}, nil
}

// Init reads the init field (proto field 1, message
// InitSubscribeRecordRequest) as the child's raw ZAP buffer; nil when absent.
// Re-wrap with WrapInitSubscribeRecordRequest.
func (t SubscribeRecordRequest) Init() []byte {
	return t.o.Bytes(subscribeRecordRequestInitOff)
}

// AckSequence reads the ack_sequence field (proto field 2, int64).
func (t SubscribeRecordRequest) AckSequence() int64 {
	return t.o.Int64(subscribeRecordRequestAckSequenceOff)
}

// AckKey reads the ack_key field (proto field 3, bytes).
func (t SubscribeRecordRequest) AckKey() []byte {
	return t.o.Bytes(subscribeRecordRequestAckKeyOff)
}

// SubscribeRecordRequestInput collects the field values for
// NewSubscribeRecordRequest. Init is the child message's own ZAP buffer (from
// NewInitSubscribeRecordRequest); nil encodes an absent init.
type SubscribeRecordRequestInput struct {
	Init        []byte
	AckSequence int64
	AckKey      []byte
}

// NewSubscribeRecordRequest builds a ZAP-encoded SubscribeRecordRequest message
// from in and returns the bytes.
func NewSubscribeRecordRequest(in SubscribeRecordRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(subscribeRecordRequestSize)
	ob.SetBytes(subscribeRecordRequestInitOff, in.Init)
	ob.SetInt64(subscribeRecordRequestAckSequenceOff, in.AckSequence)
	ob.SetBytes(subscribeRecordRequestAckKeyOff, in.AckKey)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscribeRecordResponse ---

const (
	subscribeRecordResponseKeyOff           = 0  // bytes
	subscribeRecordResponseValueOff         = 8  // message schema_pb.RecordValue
	subscribeRecordResponseTsNsOff          = 16 // int64
	subscribeRecordResponseErrorOff         = 24 // string
	subscribeRecordResponseIsEndOfStreamOff = 32 // bool
	subscribeRecordResponseIsEndOfTopicOff  = 33 // bool
	subscribeRecordResponseOffsetOff        = 40 // int64
	subscribeRecordResponseSize             = 48
)

// SubscribeRecordResponse is a zero-copy view into a ZAP-encoded
// SubscribeRecordResponse message.
type SubscribeRecordResponse struct{ o zap.Object }

// WrapSubscribeRecordResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapSubscribeRecordResponse(b []byte) (SubscribeRecordResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscribeRecordResponse{}, err
	}
	return SubscribeRecordResponse{o: m.Root()}, nil
}

// Key reads the key field (proto field 2, bytes).
func (t SubscribeRecordResponse) Key() []byte {
	return t.o.Bytes(subscribeRecordResponseKeyOff)
}

// Value reads the value field (proto field 3, message schema_pb.RecordValue)
// as the child's raw ZAP buffer; nil when absent. Re-wrap with
// mq_schemawire.WrapRecordValue.
func (t SubscribeRecordResponse) Value() []byte {
	return t.o.Bytes(subscribeRecordResponseValueOff)
}

// TsNs reads the ts_ns field (proto field 4, int64).
func (t SubscribeRecordResponse) TsNs() int64 {
	return t.o.Int64(subscribeRecordResponseTsNsOff)
}

// Error reads the error field (proto field 5, string).
func (t SubscribeRecordResponse) Error() string {
	return t.o.Text(subscribeRecordResponseErrorOff)
}

// IsEndOfStream reads the is_end_of_stream field (proto field 6, bool).
func (t SubscribeRecordResponse) IsEndOfStream() bool {
	return t.o.Bool(subscribeRecordResponseIsEndOfStreamOff)
}

// IsEndOfTopic reads the is_end_of_topic field (proto field 7, bool).
func (t SubscribeRecordResponse) IsEndOfTopic() bool {
	return t.o.Bool(subscribeRecordResponseIsEndOfTopicOff)
}

// Offset reads the offset field (proto field 8, int64). It is the sequential
// offset within the partition.
func (t SubscribeRecordResponse) Offset() int64 {
	return t.o.Int64(subscribeRecordResponseOffsetOff)
}

// SubscribeRecordResponseInput collects the field values for
// NewSubscribeRecordResponse. Value is the child message's own ZAP buffer (from
// mq_schemawire.NewRecordValue); nil encodes an absent value.
type SubscribeRecordResponseInput struct {
	Key           []byte
	Value         []byte
	TsNs          int64
	Error         string
	IsEndOfStream bool
	IsEndOfTopic  bool
	Offset        int64
}

// NewSubscribeRecordResponse builds a ZAP-encoded SubscribeRecordResponse
// message from in and returns the bytes.
func NewSubscribeRecordResponse(in SubscribeRecordResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(subscribeRecordResponseSize)
	ob.SetBytes(subscribeRecordResponseKeyOff, in.Key)
	ob.SetBytes(subscribeRecordResponseValueOff, in.Value)
	ob.SetInt64(subscribeRecordResponseTsNsOff, in.TsNs)
	ob.SetText(subscribeRecordResponseErrorOff, in.Error)
	ob.SetBool(subscribeRecordResponseIsEndOfStreamOff, in.IsEndOfStream)
	ob.SetBool(subscribeRecordResponseIsEndOfTopicOff, in.IsEndOfTopic)
	ob.SetInt64(subscribeRecordResponseOffsetOff, in.Offset)
	ob.FinishAsRoot()
	return b.Finish()
}
