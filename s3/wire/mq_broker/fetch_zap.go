// Code generated from mq_broker.proto; DO NOT EDIT.

package mq_brokerwire

import (
	zap "github.com/zap-proto/go"
)

// Stateless fetch, close, SQL-unflushed and partition-range messages.
// schema_pb.Topic / schema_pb.Partition singular fields and the filer_pb.LogEntry
// field are stored as the child's own ZAP buffer (8-byte tail pointer); re-wrap
// with mq_schemawire.WrapTopic / WrapPartition and filerwire.WrapLogEntry.
// Repeated DataMessage is a same-package list of out-of-line objects.

// --- FetchMessageRequest ---

const (
	fetchMessageRequestTopicOff         = 0  // message schema_pb.Topic
	fetchMessageRequestPartitionOff     = 8  // message schema_pb.Partition
	fetchMessageRequestStartOffsetOff   = 16 // int64
	fetchMessageRequestMaxBytesOff      = 24 // int32
	fetchMessageRequestMaxMessagesOff   = 28 // int32
	fetchMessageRequestMaxWaitMsOff     = 32 // int32
	fetchMessageRequestMinBytesOff      = 36 // int32
	fetchMessageRequestConsumerGroupOff = 40 // string
	fetchMessageRequestConsumerIdOff    = 48 // string
	fetchMessageRequestSize             = 56
)

// FetchMessageRequest is a zero-copy view into a ZAP-encoded FetchMessageRequest
// message.
type FetchMessageRequest struct{ o zap.Object }

// WrapFetchMessageRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapFetchMessageRequest(b []byte) (FetchMessageRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FetchMessageRequest{}, err
	}
	return FetchMessageRequest{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t FetchMessageRequest) Topic() []byte { return t.o.Bytes(fetchMessageRequestTopicOff) }

// Partition reads the partition field (proto field 2, message
// schema_pb.Partition) as the child's raw ZAP buffer; nil when absent. Re-wrap
// with mq_schemawire.WrapPartition.
func (t FetchMessageRequest) Partition() []byte { return t.o.Bytes(fetchMessageRequestPartitionOff) }

// StartOffset reads the start_offset field (proto field 3, int64).
func (t FetchMessageRequest) StartOffset() int64 {
	return t.o.Int64(fetchMessageRequestStartOffsetOff)
}

// MaxBytes reads the max_bytes field (proto field 4, int32).
func (t FetchMessageRequest) MaxBytes() int32 { return t.o.Int32(fetchMessageRequestMaxBytesOff) }

// MaxMessages reads the max_messages field (proto field 5, int32).
func (t FetchMessageRequest) MaxMessages() int32 {
	return t.o.Int32(fetchMessageRequestMaxMessagesOff)
}

// MaxWaitMs reads the max_wait_ms field (proto field 6, int32).
func (t FetchMessageRequest) MaxWaitMs() int32 { return t.o.Int32(fetchMessageRequestMaxWaitMsOff) }

// MinBytes reads the min_bytes field (proto field 7, int32).
func (t FetchMessageRequest) MinBytes() int32 { return t.o.Int32(fetchMessageRequestMinBytesOff) }

// ConsumerGroup reads the consumer_group field (proto field 8, string).
func (t FetchMessageRequest) ConsumerGroup() string {
	return t.o.Text(fetchMessageRequestConsumerGroupOff)
}

// ConsumerId reads the consumer_id field (proto field 9, string).
func (t FetchMessageRequest) ConsumerId() string { return t.o.Text(fetchMessageRequestConsumerIdOff) }

// FetchMessageRequestInput collects the field values for NewFetchMessageRequest.
// Topic and Partition are each the child message's own ZAP buffer (from
// mq_schemawire.NewTopic / NewPartition); nil encodes absent.
type FetchMessageRequestInput struct {
	Topic         []byte
	Partition     []byte
	StartOffset   int64
	MaxBytes      int32
	MaxMessages   int32
	MaxWaitMs     int32
	MinBytes      int32
	ConsumerGroup string
	ConsumerId    string
}

// NewFetchMessageRequest builds a ZAP-encoded FetchMessageRequest message from
// in and returns the bytes.
func NewFetchMessageRequest(in FetchMessageRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(fetchMessageRequestSize)
	ob.SetBytes(fetchMessageRequestTopicOff, in.Topic)
	ob.SetBytes(fetchMessageRequestPartitionOff, in.Partition)
	ob.SetInt64(fetchMessageRequestStartOffsetOff, in.StartOffset)
	ob.SetInt32(fetchMessageRequestMaxBytesOff, in.MaxBytes)
	ob.SetInt32(fetchMessageRequestMaxMessagesOff, in.MaxMessages)
	ob.SetInt32(fetchMessageRequestMaxWaitMsOff, in.MaxWaitMs)
	ob.SetInt32(fetchMessageRequestMinBytesOff, in.MinBytes)
	ob.SetText(fetchMessageRequestConsumerGroupOff, in.ConsumerGroup)
	ob.SetText(fetchMessageRequestConsumerIdOff, in.ConsumerId)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- FetchMessageResponse ---

const (
	fetchMessageResponseMessagesOff       = 0  // repeated DataMessage
	fetchMessageResponseHighWaterMarkOff  = 8  // int64
	fetchMessageResponseLogStartOffsetOff = 16 // int64
	fetchMessageResponseEndOfPartitionOff = 24 // bool
	fetchMessageResponseErrorOff          = 32 // string
	fetchMessageResponseErrorCodeOff      = 40 // int32
	fetchMessageResponseNextOffsetOff     = 48 // int64
	fetchMessageResponseSize              = 56
)

// FetchMessageResponse is a zero-copy view into a ZAP-encoded
// FetchMessageResponse message.
type FetchMessageResponse struct{ o zap.Object }

// WrapFetchMessageResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapFetchMessageResponse(b []byte) (FetchMessageResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FetchMessageResponse{}, err
	}
	return FetchMessageResponse{o: m.Root()}, nil
}

// MessagesLen reports the number of messages elements (proto field 1, repeated
// DataMessage).
func (t FetchMessageResponse) MessagesLen() int {
	return t.o.List(fetchMessageResponseMessagesOff).Len()
}

// MessageAt returns the i-th messages element. The bool is false when i is out
// of range or it fails to parse.
func (t FetchMessageResponse) MessageAt(i int) (DataMessage, bool) {
	o := t.o.List(fetchMessageResponseMessagesOff).ObjectAt(i)
	if o.IsNull() {
		return DataMessage{}, false
	}
	return DataMessage{o: o}, true
}

// HighWaterMark reads the high_water_mark field (proto field 2, int64).
func (t FetchMessageResponse) HighWaterMark() int64 {
	return t.o.Int64(fetchMessageResponseHighWaterMarkOff)
}

// LogStartOffset reads the log_start_offset field (proto field 3, int64).
func (t FetchMessageResponse) LogStartOffset() int64 {
	return t.o.Int64(fetchMessageResponseLogStartOffsetOff)
}

// EndOfPartition reads the end_of_partition field (proto field 4, bool).
func (t FetchMessageResponse) EndOfPartition() bool {
	return t.o.Bool(fetchMessageResponseEndOfPartitionOff)
}

// Error reads the error field (proto field 5, string).
func (t FetchMessageResponse) Error() string { return t.o.Text(fetchMessageResponseErrorOff) }

// ErrorCode reads the error_code field (proto field 6, int32).
func (t FetchMessageResponse) ErrorCode() int32 { return t.o.Int32(fetchMessageResponseErrorCodeOff) }

// NextOffset reads the next_offset field (proto field 7, int64).
func (t FetchMessageResponse) NextOffset() int64 {
	return t.o.Int64(fetchMessageResponseNextOffsetOff)
}

// FetchMessageResponseInput collects the field values for
// NewFetchMessageResponse. Each Messages entry is a DataMessage buffer (from
// NewDataMessage).
type FetchMessageResponseInput struct {
	Messages       [][]byte
	HighWaterMark  int64
	LogStartOffset int64
	EndOfPartition bool
	Error          string
	ErrorCode      int32
	NextOffset     int64
}

// NewFetchMessageResponse builds a ZAP-encoded FetchMessageResponse message from
// in and returns the bytes.
func NewFetchMessageResponse(in FetchMessageResponseInput) []byte {
	b := zap.NewBuilder(512)
	lb := b.StartList(0)
	for _, e := range in.Messages {
		lb.AddObjectBytes(e)
	}
	listOff, listLen := lb.Finish()
	ob := b.StartObject(fetchMessageResponseSize)
	ob.SetList(fetchMessageResponseMessagesOff, listOff, listLen)
	ob.SetInt64(fetchMessageResponseHighWaterMarkOff, in.HighWaterMark)
	ob.SetInt64(fetchMessageResponseLogStartOffsetOff, in.LogStartOffset)
	ob.SetBool(fetchMessageResponseEndOfPartitionOff, in.EndOfPartition)
	ob.SetText(fetchMessageResponseErrorOff, in.Error)
	ob.SetInt32(fetchMessageResponseErrorCodeOff, in.ErrorCode)
	ob.SetInt64(fetchMessageResponseNextOffsetOff, in.NextOffset)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ClosePublishersRequest ---

const (
	closePublishersRequestTopicOff      = 0 // message schema_pb.Topic
	closePublishersRequestUnixTimeNsOff = 8 // int64
	closePublishersRequestSize          = 16
)

// ClosePublishersRequest is a zero-copy view into a ZAP-encoded
// ClosePublishersRequest message.
type ClosePublishersRequest struct{ o zap.Object }

// WrapClosePublishersRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapClosePublishersRequest(b []byte) (ClosePublishersRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ClosePublishersRequest{}, err
	}
	return ClosePublishersRequest{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t ClosePublishersRequest) Topic() []byte { return t.o.Bytes(closePublishersRequestTopicOff) }

// UnixTimeNs reads the unix_time_ns field (proto field 2, int64).
func (t ClosePublishersRequest) UnixTimeNs() int64 {
	return t.o.Int64(closePublishersRequestUnixTimeNsOff)
}

// ClosePublishersRequestInput collects the field values for
// NewClosePublishersRequest. Topic is the child message's own ZAP buffer (from
// mq_schemawire.NewTopic).
type ClosePublishersRequestInput struct {
	Topic      []byte
	UnixTimeNs int64
}

// NewClosePublishersRequest builds a ZAP-encoded ClosePublishersRequest message
// from in and returns the bytes.
func NewClosePublishersRequest(in ClosePublishersRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(closePublishersRequestSize)
	ob.SetBytes(closePublishersRequestTopicOff, in.Topic)
	ob.SetInt64(closePublishersRequestUnixTimeNsOff, in.UnixTimeNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ClosePublishersResponse ---

const closePublishersResponseSize = 0

// ClosePublishersResponse is a zero-copy view into a ZAP-encoded
// ClosePublishersResponse message. It carries no fields.
type ClosePublishersResponse struct{ o zap.Object }

// WrapClosePublishersResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapClosePublishersResponse(b []byte) (ClosePublishersResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ClosePublishersResponse{}, err
	}
	return ClosePublishersResponse{o: m.Root()}, nil
}

// ClosePublishersResponseInput collects the field values for
// NewClosePublishersResponse. It has no fields.
type ClosePublishersResponseInput struct{}

// NewClosePublishersResponse builds a ZAP-encoded ClosePublishersResponse
// message and returns the bytes.
func NewClosePublishersResponse(in ClosePublishersResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(closePublishersResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- CloseSubscribersRequest ---

const (
	closeSubscribersRequestTopicOff      = 0 // message schema_pb.Topic
	closeSubscribersRequestUnixTimeNsOff = 8 // int64
	closeSubscribersRequestSize          = 16
)

// CloseSubscribersRequest is a zero-copy view into a ZAP-encoded
// CloseSubscribersRequest message.
type CloseSubscribersRequest struct{ o zap.Object }

// WrapCloseSubscribersRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapCloseSubscribersRequest(b []byte) (CloseSubscribersRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CloseSubscribersRequest{}, err
	}
	return CloseSubscribersRequest{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t CloseSubscribersRequest) Topic() []byte { return t.o.Bytes(closeSubscribersRequestTopicOff) }

// UnixTimeNs reads the unix_time_ns field (proto field 2, int64).
func (t CloseSubscribersRequest) UnixTimeNs() int64 {
	return t.o.Int64(closeSubscribersRequestUnixTimeNsOff)
}

// CloseSubscribersRequestInput collects the field values for
// NewCloseSubscribersRequest. Topic is the child message's own ZAP buffer (from
// mq_schemawire.NewTopic).
type CloseSubscribersRequestInput struct {
	Topic      []byte
	UnixTimeNs int64
}

// NewCloseSubscribersRequest builds a ZAP-encoded CloseSubscribersRequest
// message from in and returns the bytes.
func NewCloseSubscribersRequest(in CloseSubscribersRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(closeSubscribersRequestSize)
	ob.SetBytes(closeSubscribersRequestTopicOff, in.Topic)
	ob.SetInt64(closeSubscribersRequestUnixTimeNsOff, in.UnixTimeNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- CloseSubscribersResponse ---

const closeSubscribersResponseSize = 0

// CloseSubscribersResponse is a zero-copy view into a ZAP-encoded
// CloseSubscribersResponse message. It carries no fields.
type CloseSubscribersResponse struct{ o zap.Object }

// WrapCloseSubscribersResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapCloseSubscribersResponse(b []byte) (CloseSubscribersResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CloseSubscribersResponse{}, err
	}
	return CloseSubscribersResponse{o: m.Root()}, nil
}

// CloseSubscribersResponseInput collects the field values for
// NewCloseSubscribersResponse. It has no fields.
type CloseSubscribersResponseInput struct{}

// NewCloseSubscribersResponse builds a ZAP-encoded CloseSubscribersResponse
// message and returns the bytes.
func NewCloseSubscribersResponse(in CloseSubscribersResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(closeSubscribersResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetUnflushedMessagesRequest ---

const (
	getUnflushedMessagesRequestTopicOff             = 0  // message schema_pb.Topic
	getUnflushedMessagesRequestPartitionOff         = 8  // message schema_pb.Partition
	getUnflushedMessagesRequestStartBufferOffsetOff = 16 // int64
	getUnflushedMessagesRequestSize                 = 24
)

// GetUnflushedMessagesRequest is a zero-copy view into a ZAP-encoded
// GetUnflushedMessagesRequest message.
type GetUnflushedMessagesRequest struct{ o zap.Object }

// WrapGetUnflushedMessagesRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapGetUnflushedMessagesRequest(b []byte) (GetUnflushedMessagesRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetUnflushedMessagesRequest{}, err
	}
	return GetUnflushedMessagesRequest{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t GetUnflushedMessagesRequest) Topic() []byte {
	return t.o.Bytes(getUnflushedMessagesRequestTopicOff)
}

// Partition reads the partition field (proto field 2, message
// schema_pb.Partition) as the child's raw ZAP buffer; nil when absent. Re-wrap
// with mq_schemawire.WrapPartition.
func (t GetUnflushedMessagesRequest) Partition() []byte {
	return t.o.Bytes(getUnflushedMessagesRequestPartitionOff)
}

// StartBufferOffset reads the start_buffer_offset field (proto field 3, int64).
func (t GetUnflushedMessagesRequest) StartBufferOffset() int64 {
	return t.o.Int64(getUnflushedMessagesRequestStartBufferOffsetOff)
}

// GetUnflushedMessagesRequestInput collects the field values for
// NewGetUnflushedMessagesRequest. Topic and Partition are each the child
// message's own ZAP buffer (from mq_schemawire.NewTopic / NewPartition).
type GetUnflushedMessagesRequestInput struct {
	Topic             []byte
	Partition         []byte
	StartBufferOffset int64
}

// NewGetUnflushedMessagesRequest builds a ZAP-encoded GetUnflushedMessagesRequest
// message from in and returns the bytes.
func NewGetUnflushedMessagesRequest(in GetUnflushedMessagesRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getUnflushedMessagesRequestSize)
	ob.SetBytes(getUnflushedMessagesRequestTopicOff, in.Topic)
	ob.SetBytes(getUnflushedMessagesRequestPartitionOff, in.Partition)
	ob.SetInt64(getUnflushedMessagesRequestStartBufferOffsetOff, in.StartBufferOffset)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetUnflushedMessagesResponse ---

const (
	getUnflushedMessagesResponseMessageOff     = 0  // message filer_pb.LogEntry
	getUnflushedMessagesResponseErrorOff       = 8  // string
	getUnflushedMessagesResponseEndOfStreamOff = 16 // bool
	getUnflushedMessagesResponseSize           = 24
)

// GetUnflushedMessagesResponse is a zero-copy view into a ZAP-encoded
// GetUnflushedMessagesResponse message.
type GetUnflushedMessagesResponse struct{ o zap.Object }

// WrapGetUnflushedMessagesResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapGetUnflushedMessagesResponse(b []byte) (GetUnflushedMessagesResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetUnflushedMessagesResponse{}, err
	}
	return GetUnflushedMessagesResponse{o: m.Root()}, nil
}

// Message reads the message field (proto field 1, message filer_pb.LogEntry) as
// the child's raw ZAP buffer; nil when absent. Re-wrap with
// filerwire.WrapLogEntry.
func (t GetUnflushedMessagesResponse) Message() []byte {
	return t.o.Bytes(getUnflushedMessagesResponseMessageOff)
}

// Error reads the error field (proto field 2, string).
func (t GetUnflushedMessagesResponse) Error() string {
	return t.o.Text(getUnflushedMessagesResponseErrorOff)
}

// EndOfStream reads the end_of_stream field (proto field 3, bool).
func (t GetUnflushedMessagesResponse) EndOfStream() bool {
	return t.o.Bool(getUnflushedMessagesResponseEndOfStreamOff)
}

// GetUnflushedMessagesResponseInput collects the field values for
// NewGetUnflushedMessagesResponse. Message is the child message's own ZAP buffer
// (from filerwire.NewLogEntry); nil encodes absent.
type GetUnflushedMessagesResponseInput struct {
	Message     []byte
	Error       string
	EndOfStream bool
}

// NewGetUnflushedMessagesResponse builds a ZAP-encoded
// GetUnflushedMessagesResponse message from in and returns the bytes.
func NewGetUnflushedMessagesResponse(in GetUnflushedMessagesResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getUnflushedMessagesResponseSize)
	ob.SetBytes(getUnflushedMessagesResponseMessageOff, in.Message)
	ob.SetText(getUnflushedMessagesResponseErrorOff, in.Error)
	ob.SetBool(getUnflushedMessagesResponseEndOfStreamOff, in.EndOfStream)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- OffsetRangeInfo ---

const (
	offsetRangeInfoEarliestOffsetOff = 0  // int64
	offsetRangeInfoLatestOffsetOff   = 8  // int64
	offsetRangeInfoHighWaterMarkOff  = 16 // int64
	offsetRangeInfoSize              = 24
)

// OffsetRangeInfo is a zero-copy view into a ZAP-encoded OffsetRangeInfo
// message.
type OffsetRangeInfo struct{ o zap.Object }

// WrapOffsetRangeInfo parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapOffsetRangeInfo(b []byte) (OffsetRangeInfo, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return OffsetRangeInfo{}, err
	}
	return OffsetRangeInfo{o: m.Root()}, nil
}

// EarliestOffset reads the earliest_offset field (proto field 1, int64).
func (t OffsetRangeInfo) EarliestOffset() int64 {
	return t.o.Int64(offsetRangeInfoEarliestOffsetOff)
}

// LatestOffset reads the latest_offset field (proto field 2, int64).
func (t OffsetRangeInfo) LatestOffset() int64 { return t.o.Int64(offsetRangeInfoLatestOffsetOff) }

// HighWaterMark reads the high_water_mark field (proto field 3, int64).
func (t OffsetRangeInfo) HighWaterMark() int64 {
	return t.o.Int64(offsetRangeInfoHighWaterMarkOff)
}

// OffsetRangeInfoInput collects the field values for NewOffsetRangeInfo.
type OffsetRangeInfoInput struct {
	EarliestOffset int64
	LatestOffset   int64
	HighWaterMark  int64
}

// NewOffsetRangeInfo builds a ZAP-encoded OffsetRangeInfo message from in and
// returns the bytes.
func NewOffsetRangeInfo(in OffsetRangeInfoInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(offsetRangeInfoSize)
	ob.SetInt64(offsetRangeInfoEarliestOffsetOff, in.EarliestOffset)
	ob.SetInt64(offsetRangeInfoLatestOffsetOff, in.LatestOffset)
	ob.SetInt64(offsetRangeInfoHighWaterMarkOff, in.HighWaterMark)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TimestampRangeInfo ---

const (
	timestampRangeInfoEarliestTimestampNsOff = 0 // int64
	timestampRangeInfoLatestTimestampNsOff   = 8 // int64
	timestampRangeInfoSize                   = 16
)

// TimestampRangeInfo is a zero-copy view into a ZAP-encoded TimestampRangeInfo
// message.
type TimestampRangeInfo struct{ o zap.Object }

// WrapTimestampRangeInfo parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapTimestampRangeInfo(b []byte) (TimestampRangeInfo, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TimestampRangeInfo{}, err
	}
	return TimestampRangeInfo{o: m.Root()}, nil
}

// EarliestTimestampNs reads the earliest_timestamp_ns field (proto field 1,
// int64).
func (t TimestampRangeInfo) EarliestTimestampNs() int64 {
	return t.o.Int64(timestampRangeInfoEarliestTimestampNsOff)
}

// LatestTimestampNs reads the latest_timestamp_ns field (proto field 2, int64).
func (t TimestampRangeInfo) LatestTimestampNs() int64 {
	return t.o.Int64(timestampRangeInfoLatestTimestampNsOff)
}

// TimestampRangeInfoInput collects the field values for NewTimestampRangeInfo.
type TimestampRangeInfoInput struct {
	EarliestTimestampNs int64
	LatestTimestampNs   int64
}

// NewTimestampRangeInfo builds a ZAP-encoded TimestampRangeInfo message from in
// and returns the bytes.
func NewTimestampRangeInfo(in TimestampRangeInfoInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(timestampRangeInfoSize)
	ob.SetInt64(timestampRangeInfoEarliestTimestampNsOff, in.EarliestTimestampNs)
	ob.SetInt64(timestampRangeInfoLatestTimestampNsOff, in.LatestTimestampNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetPartitionRangeInfoRequest ---

const (
	getPartitionRangeInfoRequestTopicOff     = 0 // message schema_pb.Topic
	getPartitionRangeInfoRequestPartitionOff = 8 // message schema_pb.Partition
	getPartitionRangeInfoRequestSize         = 16
)

// GetPartitionRangeInfoRequest is a zero-copy view into a ZAP-encoded
// GetPartitionRangeInfoRequest message.
type GetPartitionRangeInfoRequest struct{ o zap.Object }

// WrapGetPartitionRangeInfoRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapGetPartitionRangeInfoRequest(b []byte) (GetPartitionRangeInfoRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetPartitionRangeInfoRequest{}, err
	}
	return GetPartitionRangeInfoRequest{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t GetPartitionRangeInfoRequest) Topic() []byte {
	return t.o.Bytes(getPartitionRangeInfoRequestTopicOff)
}

// Partition reads the partition field (proto field 2, message
// schema_pb.Partition) as the child's raw ZAP buffer; nil when absent. Re-wrap
// with mq_schemawire.WrapPartition.
func (t GetPartitionRangeInfoRequest) Partition() []byte {
	return t.o.Bytes(getPartitionRangeInfoRequestPartitionOff)
}

// GetPartitionRangeInfoRequestInput collects the field values for
// NewGetPartitionRangeInfoRequest. Topic and Partition are each the child
// message's own ZAP buffer (from mq_schemawire.NewTopic / NewPartition).
type GetPartitionRangeInfoRequestInput struct {
	Topic     []byte
	Partition []byte
}

// NewGetPartitionRangeInfoRequest builds a ZAP-encoded
// GetPartitionRangeInfoRequest message from in and returns the bytes.
func NewGetPartitionRangeInfoRequest(in GetPartitionRangeInfoRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getPartitionRangeInfoRequestSize)
	ob.SetBytes(getPartitionRangeInfoRequestTopicOff, in.Topic)
	ob.SetBytes(getPartitionRangeInfoRequestPartitionOff, in.Partition)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetPartitionRangeInfoResponse ---

const (
	getPartitionRangeInfoResponseOffsetRangeOff         = 0  // message OffsetRangeInfo
	getPartitionRangeInfoResponseTimestampRangeOff      = 8  // message TimestampRangeInfo
	getPartitionRangeInfoResponseRecordCountOff         = 16 // int64
	getPartitionRangeInfoResponseActiveSubscriptionsOff = 24 // int64
	getPartitionRangeInfoResponseErrorOff               = 32 // string
	getPartitionRangeInfoResponseSize                   = 40
)

// GetPartitionRangeInfoResponse is a zero-copy view into a ZAP-encoded
// GetPartitionRangeInfoResponse message.
type GetPartitionRangeInfoResponse struct{ o zap.Object }

// WrapGetPartitionRangeInfoResponse parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapGetPartitionRangeInfoResponse(b []byte) (GetPartitionRangeInfoResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetPartitionRangeInfoResponse{}, err
	}
	return GetPartitionRangeInfoResponse{o: m.Root()}, nil
}

// OffsetRange reads the offset_range field (proto field 1, message
// OffsetRangeInfo). The bool is false when the field is absent.
func (t GetPartitionRangeInfoResponse) OffsetRange() (OffsetRangeInfo, bool) {
	b := t.o.Bytes(getPartitionRangeInfoResponseOffsetRangeOff)
	if len(b) == 0 {
		return OffsetRangeInfo{}, false
	}
	r, err := WrapOffsetRangeInfo(b)
	if err != nil {
		return OffsetRangeInfo{}, false
	}
	return r, true
}

// TimestampRange reads the timestamp_range field (proto field 2, message
// TimestampRangeInfo). The bool is false when the field is absent.
func (t GetPartitionRangeInfoResponse) TimestampRange() (TimestampRangeInfo, bool) {
	b := t.o.Bytes(getPartitionRangeInfoResponseTimestampRangeOff)
	if len(b) == 0 {
		return TimestampRangeInfo{}, false
	}
	r, err := WrapTimestampRangeInfo(b)
	if err != nil {
		return TimestampRangeInfo{}, false
	}
	return r, true
}

// RecordCount reads the record_count field (proto field 10, int64).
func (t GetPartitionRangeInfoResponse) RecordCount() int64 {
	return t.o.Int64(getPartitionRangeInfoResponseRecordCountOff)
}

// ActiveSubscriptions reads the active_subscriptions field (proto field 11,
// int64).
func (t GetPartitionRangeInfoResponse) ActiveSubscriptions() int64 {
	return t.o.Int64(getPartitionRangeInfoResponseActiveSubscriptionsOff)
}

// Error reads the error field (proto field 12, string).
func (t GetPartitionRangeInfoResponse) Error() string {
	return t.o.Text(getPartitionRangeInfoResponseErrorOff)
}

// GetPartitionRangeInfoResponseInput collects the field values for
// NewGetPartitionRangeInfoResponse. OffsetRange and TimestampRange are each the
// child message's own ZAP buffer (from NewOffsetRangeInfo / NewTimestampRangeInfo).
type GetPartitionRangeInfoResponseInput struct {
	OffsetRange         []byte
	TimestampRange      []byte
	RecordCount         int64
	ActiveSubscriptions int64
	Error               string
}

// NewGetPartitionRangeInfoResponse builds a ZAP-encoded
// GetPartitionRangeInfoResponse message from in and returns the bytes.
func NewGetPartitionRangeInfoResponse(in GetPartitionRangeInfoResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getPartitionRangeInfoResponseSize)
	ob.SetBytes(getPartitionRangeInfoResponseOffsetRangeOff, in.OffsetRange)
	ob.SetBytes(getPartitionRangeInfoResponseTimestampRangeOff, in.TimestampRange)
	ob.SetInt64(getPartitionRangeInfoResponseRecordCountOff, in.RecordCount)
	ob.SetInt64(getPartitionRangeInfoResponseActiveSubscriptionsOff, in.ActiveSubscriptions)
	ob.SetText(getPartitionRangeInfoResponseErrorOff, in.Error)
	ob.FinishAsRoot()
	return b.Finish()
}
