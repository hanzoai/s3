// Code generated from mq_broker.proto; DO NOT EDIT.

package mq_brokerwire

import (
	zap "github.com/zap-proto/go"
)

// Control-plane and balancer messages. schema_pb.Topic / schema_pb.Partition
// singular fields are stored as the child's own ZAP buffer (8-byte tail
// pointer); build them with the mq_schema wire package and re-wrap with
// mq_schemawire.WrapTopic / WrapPartition.

// --- FindBrokerLeaderRequest ---

const (
	findBrokerLeaderRequestFilerGroupOff = 0 // string
	findBrokerLeaderRequestSize          = 8
)

// FindBrokerLeaderRequest is a zero-copy view into a ZAP-encoded
// FindBrokerLeaderRequest message.
type FindBrokerLeaderRequest struct{ o zap.Object }

// WrapFindBrokerLeaderRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapFindBrokerLeaderRequest(b []byte) (FindBrokerLeaderRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FindBrokerLeaderRequest{}, err
	}
	return FindBrokerLeaderRequest{o: m.Root()}, nil
}

// FilerGroup reads the filer_group field (proto field 1, string).
func (t FindBrokerLeaderRequest) FilerGroup() string {
	return t.o.Text(findBrokerLeaderRequestFilerGroupOff)
}

// FindBrokerLeaderRequestInput collects the field values for
// NewFindBrokerLeaderRequest.
type FindBrokerLeaderRequestInput struct {
	FilerGroup string
}

// NewFindBrokerLeaderRequest builds a ZAP-encoded FindBrokerLeaderRequest
// message from in and returns the bytes.
func NewFindBrokerLeaderRequest(in FindBrokerLeaderRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(findBrokerLeaderRequestSize)
	ob.SetText(findBrokerLeaderRequestFilerGroupOff, in.FilerGroup)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- FindBrokerLeaderResponse ---

const (
	findBrokerLeaderResponseBrokerOff = 0 // string
	findBrokerLeaderResponseSize      = 8
)

// FindBrokerLeaderResponse is a zero-copy view into a ZAP-encoded
// FindBrokerLeaderResponse message.
type FindBrokerLeaderResponse struct{ o zap.Object }

// WrapFindBrokerLeaderResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapFindBrokerLeaderResponse(b []byte) (FindBrokerLeaderResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FindBrokerLeaderResponse{}, err
	}
	return FindBrokerLeaderResponse{o: m.Root()}, nil
}

// Broker reads the broker field (proto field 1, string).
func (t FindBrokerLeaderResponse) Broker() string {
	return t.o.Text(findBrokerLeaderResponseBrokerOff)
}

// FindBrokerLeaderResponseInput collects the field values for
// NewFindBrokerLeaderResponse.
type FindBrokerLeaderResponseInput struct {
	Broker string
}

// NewFindBrokerLeaderResponse builds a ZAP-encoded FindBrokerLeaderResponse
// message from in and returns the bytes.
func NewFindBrokerLeaderResponse(in FindBrokerLeaderResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(findBrokerLeaderResponseSize)
	ob.SetText(findBrokerLeaderResponseBrokerOff, in.Broker)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TopicPartitionStats ---

const (
	topicPartitionStatsTopicOff           = 0  // message schema_pb.Topic
	topicPartitionStatsPartitionOff       = 8  // message schema_pb.Partition
	topicPartitionStatsPublisherCountOff  = 16 // int32
	topicPartitionStatsSubscriberCountOff = 20 // int32
	topicPartitionStatsFollowerOff        = 24 // string
	topicPartitionStatsSize               = 32
)

// TopicPartitionStats is a zero-copy view into a ZAP-encoded
// TopicPartitionStats message.
type TopicPartitionStats struct{ o zap.Object }

// WrapTopicPartitionStats parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapTopicPartitionStats(b []byte) (TopicPartitionStats, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TopicPartitionStats{}, err
	}
	return TopicPartitionStats{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t TopicPartitionStats) Topic() []byte { return t.o.Bytes(topicPartitionStatsTopicOff) }

// Partition reads the partition field (proto field 2, message
// schema_pb.Partition) as the child's raw ZAP buffer; nil when absent. Re-wrap
// with mq_schemawire.WrapPartition.
func (t TopicPartitionStats) Partition() []byte { return t.o.Bytes(topicPartitionStatsPartitionOff) }

// PublisherCount reads the publisher_count field (proto field 3, int32).
func (t TopicPartitionStats) PublisherCount() int32 {
	return t.o.Int32(topicPartitionStatsPublisherCountOff)
}

// SubscriberCount reads the subscriber_count field (proto field 4, int32).
func (t TopicPartitionStats) SubscriberCount() int32 {
	return t.o.Int32(topicPartitionStatsSubscriberCountOff)
}

// Follower reads the follower field (proto field 5, string).
func (t TopicPartitionStats) Follower() string { return t.o.Text(topicPartitionStatsFollowerOff) }

// TopicPartitionStatsInput collects the field values for
// NewTopicPartitionStats. Topic and Partition are each the child message's own
// ZAP buffer (from mq_schemawire.NewTopic / NewPartition); nil encodes absent.
type TopicPartitionStatsInput struct {
	Topic           []byte
	Partition       []byte
	PublisherCount  int32
	SubscriberCount int32
	Follower        string
}

// NewTopicPartitionStats builds a ZAP-encoded TopicPartitionStats message from
// in and returns the bytes.
func NewTopicPartitionStats(in TopicPartitionStatsInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(topicPartitionStatsSize)
	ob.SetBytes(topicPartitionStatsTopicOff, in.Topic)
	ob.SetBytes(topicPartitionStatsPartitionOff, in.Partition)
	ob.SetInt32(topicPartitionStatsPublisherCountOff, in.PublisherCount)
	ob.SetInt32(topicPartitionStatsSubscriberCountOff, in.SubscriberCount)
	ob.SetText(topicPartitionStatsFollowerOff, in.Follower)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- BrokerStats ---

const (
	brokerStatsCpuUsagePercentOff = 0 // int32
	brokerStatsStatsOff           = 8 // map<string,TopicPartitionStats>
	brokerStatsSize               = 16
)

// BrokerStats is a zero-copy view into a ZAP-encoded BrokerStats message.
type BrokerStats struct{ o zap.Object }

// WrapBrokerStats parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapBrokerStats(b []byte) (BrokerStats, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return BrokerStats{}, err
	}
	return BrokerStats{o: m.Root()}, nil
}

// CpuUsagePercent reads the cpu_usage_percent field (proto field 1, int32).
func (t BrokerStats) CpuUsagePercent() int32 { return t.o.Int32(brokerStatsCpuUsagePercentOff) }

// StatsLen reports the number of stats entries (proto field 2,
// map<string,TopicPartitionStats>).
func (t BrokerStats) StatsLen() int { return t.o.List(brokerStatsStatsOff).Len() }

// StatsAt returns the i-th stats entry. The entry Value is a raw
// TopicPartitionStats buffer (re-wrap with WrapTopicPartitionStats).
func (t BrokerStats) StatsAt(i int) StringMsgEntry {
	return stringMsgEntryAt(t.o.List(brokerStatsStatsOff), i)
}

// BrokerStatsInput collects the field values for NewBrokerStats. Each Stats
// entry's Value is a TopicPartitionStats buffer (from NewTopicPartitionStats).
type BrokerStatsInput struct {
	CpuUsagePercent int32
	Stats           []StringMsgEntry
}

// NewBrokerStats builds a ZAP-encoded BrokerStats message from in and returns
// the bytes.
func NewBrokerStats(in BrokerStatsInput) []byte {
	b := zap.NewBuilder(256)
	statsOff, statsLen := setStringMsgMap(b, in.Stats)
	ob := b.StartObject(brokerStatsSize)
	ob.SetInt32(brokerStatsCpuUsagePercentOff, in.CpuUsagePercent)
	ob.SetList(brokerStatsStatsOff, statsOff, statsLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PublisherToPubBalancerRequestInitMessage ---
// (oneof message variant `init` of PublisherToPubBalancerRequest)

const (
	publisherToPubBalancerRequestInitMessageBrokerOff = 0 // string
	publisherToPubBalancerRequestInitMessageSize      = 8
)

// PublisherToPubBalancerRequestInitMessage is a zero-copy view into a
// ZAP-encoded PublisherToPubBalancerRequest.InitMessage message.
type PublisherToPubBalancerRequestInitMessage struct{ o zap.Object }

// WrapPublisherToPubBalancerRequestInitMessage parses b and returns a typed
// view. Returns an error if the wire-level checks (magic, version, size) fail.
func WrapPublisherToPubBalancerRequestInitMessage(b []byte) (PublisherToPubBalancerRequestInitMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PublisherToPubBalancerRequestInitMessage{}, err
	}
	return PublisherToPubBalancerRequestInitMessage{o: m.Root()}, nil
}

// Broker reads the broker field (proto field 1, string).
func (t PublisherToPubBalancerRequestInitMessage) Broker() string {
	return t.o.Text(publisherToPubBalancerRequestInitMessageBrokerOff)
}

// PublisherToPubBalancerRequestInitMessageInput collects the field values for
// NewPublisherToPubBalancerRequestInitMessage.
type PublisherToPubBalancerRequestInitMessageInput struct {
	Broker string
}

// NewPublisherToPubBalancerRequestInitMessage builds a ZAP-encoded
// PublisherToPubBalancerRequest.InitMessage message from in and returns the
// bytes.
func NewPublisherToPubBalancerRequestInitMessage(in PublisherToPubBalancerRequestInitMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(publisherToPubBalancerRequestInitMessageSize)
	ob.SetText(publisherToPubBalancerRequestInitMessageBrokerOff, in.Broker)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PublisherToPubBalancerRequest ---

// PublisherToPubBalancerRequest.message oneof tags (proto field numbers, stable
// wire ids). Zero means no variant is set. The selected variant message is
// embedded as a nested object pointed to by message_value.
const (
	PublisherToPubBalancerRequestMessageNone  uint32 = 0
	PublisherToPubBalancerRequestMessageInit  uint32 = 1
	PublisherToPubBalancerRequestMessageStats uint32 = 2
)

const (
	publisherToPubBalancerRequestMessageWhichOff = 0 // oneof tag (uint32)
	publisherToPubBalancerRequestMessageValueOff = 4 // embedded variant (obj ptr)
	publisherToPubBalancerRequestSize            = 8
)

// PublisherToPubBalancerRequest is a zero-copy view into a ZAP-encoded
// PublisherToPubBalancerRequest message.
type PublisherToPubBalancerRequest struct{ o zap.Object }

// WrapPublisherToPubBalancerRequest parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapPublisherToPubBalancerRequest(b []byte) (PublisherToPubBalancerRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PublisherToPubBalancerRequest{}, err
	}
	return PublisherToPubBalancerRequest{o: m.Root()}, nil
}

// WhichMessage returns the set oneof tag
// (PublisherToPubBalancerRequestMessage* constants).
func (t PublisherToPubBalancerRequest) WhichMessage() uint32 {
	return t.o.Uint32(publisherToPubBalancerRequestMessageWhichOff)
}

// Init returns the embedded init variant.
func (t PublisherToPubBalancerRequest) Init() PublisherToPubBalancerRequestInitMessage {
	return PublisherToPubBalancerRequestInitMessage{o: t.o.Object(publisherToPubBalancerRequestMessageValueOff)}
}

// Stats returns the embedded stats variant.
func (t PublisherToPubBalancerRequest) Stats() BrokerStats {
	return BrokerStats{o: t.o.Object(publisherToPubBalancerRequestMessageValueOff)}
}

// PublisherToPubBalancerRequestInput collects the field values for
// NewPublisherToPubBalancerRequest. MessageWhich selects the oneof variant;
// MessageValue is that variant's standalone buffer.
type PublisherToPubBalancerRequestInput struct {
	MessageWhich uint32
	MessageValue []byte
}

// NewPublisherToPubBalancerRequest builds a ZAP-encoded
// PublisherToPubBalancerRequest message from in and returns the bytes.
func NewPublisherToPubBalancerRequest(in PublisherToPubBalancerRequestInput) []byte {
	b := zap.NewBuilder(512)
	valueOff := 0
	if len(in.MessageValue) > 0 {
		valueOff = embedMessage(b, in.MessageValue)
	}
	ob := b.StartObject(publisherToPubBalancerRequestSize)
	ob.SetUint32(publisherToPubBalancerRequestMessageWhichOff, in.MessageWhich)
	ob.SetObject(publisherToPubBalancerRequestMessageValueOff, valueOff)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PublisherToPubBalancerResponse ---

const publisherToPubBalancerResponseSize = 0

// PublisherToPubBalancerResponse is a zero-copy view into a ZAP-encoded
// PublisherToPubBalancerResponse message. It carries no fields.
type PublisherToPubBalancerResponse struct{ o zap.Object }

// WrapPublisherToPubBalancerResponse parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapPublisherToPubBalancerResponse(b []byte) (PublisherToPubBalancerResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PublisherToPubBalancerResponse{}, err
	}
	return PublisherToPubBalancerResponse{o: m.Root()}, nil
}

// PublisherToPubBalancerResponseInput collects the field values for
// NewPublisherToPubBalancerResponse. It has no fields.
type PublisherToPubBalancerResponseInput struct{}

// NewPublisherToPubBalancerResponse builds a ZAP-encoded
// PublisherToPubBalancerResponse message and returns the bytes.
func NewPublisherToPubBalancerResponse(in PublisherToPubBalancerResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(publisherToPubBalancerResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- BalanceTopicsRequest ---

const balanceTopicsRequestSize = 0

// BalanceTopicsRequest is a zero-copy view into a ZAP-encoded
// BalanceTopicsRequest message. It carries no fields.
type BalanceTopicsRequest struct{ o zap.Object }

// WrapBalanceTopicsRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapBalanceTopicsRequest(b []byte) (BalanceTopicsRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return BalanceTopicsRequest{}, err
	}
	return BalanceTopicsRequest{o: m.Root()}, nil
}

// BalanceTopicsRequestInput collects the field values for
// NewBalanceTopicsRequest. It has no fields.
type BalanceTopicsRequestInput struct{}

// NewBalanceTopicsRequest builds a ZAP-encoded BalanceTopicsRequest message and
// returns the bytes.
func NewBalanceTopicsRequest(in BalanceTopicsRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(balanceTopicsRequestSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- BalanceTopicsResponse ---

const balanceTopicsResponseSize = 0

// BalanceTopicsResponse is a zero-copy view into a ZAP-encoded
// BalanceTopicsResponse message. It carries no fields.
type BalanceTopicsResponse struct{ o zap.Object }

// WrapBalanceTopicsResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapBalanceTopicsResponse(b []byte) (BalanceTopicsResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return BalanceTopicsResponse{}, err
	}
	return BalanceTopicsResponse{o: m.Root()}, nil
}

// BalanceTopicsResponseInput collects the field values for
// NewBalanceTopicsResponse. It has no fields.
type BalanceTopicsResponseInput struct{}

// NewBalanceTopicsResponse builds a ZAP-encoded BalanceTopicsResponse message
// and returns the bytes.
func NewBalanceTopicsResponse(in BalanceTopicsResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(balanceTopicsResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}
