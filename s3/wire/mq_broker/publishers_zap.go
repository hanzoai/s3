// Code generated from mq_broker.proto; DO NOT EDIT.

package mq_brokerwire

import (
	zap "github.com/zap-proto/go"
)

// Publisher / subscriber introspection messages. schema_pb.Partition singular
// fields are stored as the child's own ZAP buffer (8-byte tail pointer); re-wrap
// with mq_schemawire.WrapPartition.

// --- TopicPublisher ---

const (
	topicPublisherPublisherNameOff       = 0  // string
	topicPublisherClientIdOff            = 8  // string
	topicPublisherPartitionOff           = 16 // message schema_pb.Partition
	topicPublisherConnectTimeNsOff       = 24 // int64
	topicPublisherLastSeenTimeNsOff      = 32 // int64
	topicPublisherBrokerOff              = 40 // string
	topicPublisherIsActiveOff            = 48 // bool
	topicPublisherLastPublishedOffsetOff = 56 // int64
	topicPublisherLastAckedOffsetOff     = 64 // int64
	topicPublisherSize                   = 72
)

// TopicPublisher is a zero-copy view into a ZAP-encoded TopicPublisher message.
type TopicPublisher struct{ o zap.Object }

// WrapTopicPublisher parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapTopicPublisher(b []byte) (TopicPublisher, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TopicPublisher{}, err
	}
	return TopicPublisher{o: m.Root()}, nil
}

// PublisherName reads the publisher_name field (proto field 1, string).
func (t TopicPublisher) PublisherName() string { return t.o.Text(topicPublisherPublisherNameOff) }

// ClientId reads the client_id field (proto field 2, string).
func (t TopicPublisher) ClientId() string { return t.o.Text(topicPublisherClientIdOff) }

// Partition reads the partition field (proto field 3, message
// schema_pb.Partition) as the child's raw ZAP buffer; nil when absent. Re-wrap
// with mq_schemawire.WrapPartition.
func (t TopicPublisher) Partition() []byte { return t.o.Bytes(topicPublisherPartitionOff) }

// ConnectTimeNs reads the connect_time_ns field (proto field 4, int64).
func (t TopicPublisher) ConnectTimeNs() int64 { return t.o.Int64(topicPublisherConnectTimeNsOff) }

// LastSeenTimeNs reads the last_seen_time_ns field (proto field 5, int64).
func (t TopicPublisher) LastSeenTimeNs() int64 { return t.o.Int64(topicPublisherLastSeenTimeNsOff) }

// Broker reads the broker field (proto field 6, string).
func (t TopicPublisher) Broker() string { return t.o.Text(topicPublisherBrokerOff) }

// IsActive reads the is_active field (proto field 7, bool).
func (t TopicPublisher) IsActive() bool { return t.o.Bool(topicPublisherIsActiveOff) }

// LastPublishedOffset reads the last_published_offset field (proto field 8,
// int64).
func (t TopicPublisher) LastPublishedOffset() int64 {
	return t.o.Int64(topicPublisherLastPublishedOffsetOff)
}

// LastAckedOffset reads the last_acked_offset field (proto field 9, int64).
func (t TopicPublisher) LastAckedOffset() int64 {
	return t.o.Int64(topicPublisherLastAckedOffsetOff)
}

// TopicPublisherInput collects the field values for NewTopicPublisher.
// Partition is the child message's own ZAP buffer (from
// mq_schemawire.NewPartition); nil encodes absent.
type TopicPublisherInput struct {
	PublisherName       string
	ClientId            string
	Partition           []byte
	ConnectTimeNs       int64
	LastSeenTimeNs      int64
	Broker              string
	IsActive            bool
	LastPublishedOffset int64
	LastAckedOffset     int64
}

// NewTopicPublisher builds a ZAP-encoded TopicPublisher message from in and
// returns the bytes.
func NewTopicPublisher(in TopicPublisherInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(topicPublisherSize)
	ob.SetText(topicPublisherPublisherNameOff, in.PublisherName)
	ob.SetText(topicPublisherClientIdOff, in.ClientId)
	ob.SetBytes(topicPublisherPartitionOff, in.Partition)
	ob.SetInt64(topicPublisherConnectTimeNsOff, in.ConnectTimeNs)
	ob.SetInt64(topicPublisherLastSeenTimeNsOff, in.LastSeenTimeNs)
	ob.SetText(topicPublisherBrokerOff, in.Broker)
	ob.SetBool(topicPublisherIsActiveOff, in.IsActive)
	ob.SetInt64(topicPublisherLastPublishedOffsetOff, in.LastPublishedOffset)
	ob.SetInt64(topicPublisherLastAckedOffsetOff, in.LastAckedOffset)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TopicSubscriber ---

const (
	topicSubscriberConsumerGroupOff      = 0  // string
	topicSubscriberConsumerIdOff         = 8  // string
	topicSubscriberClientIdOff           = 16 // string
	topicSubscriberPartitionOff          = 24 // message schema_pb.Partition
	topicSubscriberConnectTimeNsOff      = 32 // int64
	topicSubscriberLastSeenTimeNsOff     = 40 // int64
	topicSubscriberBrokerOff             = 48 // string
	topicSubscriberIsActiveOff           = 56 // bool
	topicSubscriberCurrentOffsetOff      = 64 // int64
	topicSubscriberLastReceivedOffsetOff = 72 // int64
	topicSubscriberSize                  = 80
)

// TopicSubscriber is a zero-copy view into a ZAP-encoded TopicSubscriber
// message.
type TopicSubscriber struct{ o zap.Object }

// WrapTopicSubscriber parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapTopicSubscriber(b []byte) (TopicSubscriber, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TopicSubscriber{}, err
	}
	return TopicSubscriber{o: m.Root()}, nil
}

// ConsumerGroup reads the consumer_group field (proto field 1, string).
func (t TopicSubscriber) ConsumerGroup() string { return t.o.Text(topicSubscriberConsumerGroupOff) }

// ConsumerId reads the consumer_id field (proto field 2, string).
func (t TopicSubscriber) ConsumerId() string { return t.o.Text(topicSubscriberConsumerIdOff) }

// ClientId reads the client_id field (proto field 3, string).
func (t TopicSubscriber) ClientId() string { return t.o.Text(topicSubscriberClientIdOff) }

// Partition reads the partition field (proto field 4, message
// schema_pb.Partition) as the child's raw ZAP buffer; nil when absent. Re-wrap
// with mq_schemawire.WrapPartition.
func (t TopicSubscriber) Partition() []byte { return t.o.Bytes(topicSubscriberPartitionOff) }

// ConnectTimeNs reads the connect_time_ns field (proto field 5, int64).
func (t TopicSubscriber) ConnectTimeNs() int64 { return t.o.Int64(topicSubscriberConnectTimeNsOff) }

// LastSeenTimeNs reads the last_seen_time_ns field (proto field 6, int64).
func (t TopicSubscriber) LastSeenTimeNs() int64 {
	return t.o.Int64(topicSubscriberLastSeenTimeNsOff)
}

// Broker reads the broker field (proto field 7, string).
func (t TopicSubscriber) Broker() string { return t.o.Text(topicSubscriberBrokerOff) }

// IsActive reads the is_active field (proto field 8, bool).
func (t TopicSubscriber) IsActive() bool { return t.o.Bool(topicSubscriberIsActiveOff) }

// CurrentOffset reads the current_offset field (proto field 9, int64).
func (t TopicSubscriber) CurrentOffset() int64 { return t.o.Int64(topicSubscriberCurrentOffsetOff) }

// LastReceivedOffset reads the last_received_offset field (proto field 10,
// int64).
func (t TopicSubscriber) LastReceivedOffset() int64 {
	return t.o.Int64(topicSubscriberLastReceivedOffsetOff)
}

// TopicSubscriberInput collects the field values for NewTopicSubscriber.
// Partition is the child message's own ZAP buffer (from
// mq_schemawire.NewPartition); nil encodes absent.
type TopicSubscriberInput struct {
	ConsumerGroup      string
	ConsumerId         string
	ClientId           string
	Partition          []byte
	ConnectTimeNs      int64
	LastSeenTimeNs     int64
	Broker             string
	IsActive           bool
	CurrentOffset      int64
	LastReceivedOffset int64
}

// NewTopicSubscriber builds a ZAP-encoded TopicSubscriber message from in and
// returns the bytes.
func NewTopicSubscriber(in TopicSubscriberInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(topicSubscriberSize)
	ob.SetText(topicSubscriberConsumerGroupOff, in.ConsumerGroup)
	ob.SetText(topicSubscriberConsumerIdOff, in.ConsumerId)
	ob.SetText(topicSubscriberClientIdOff, in.ClientId)
	ob.SetBytes(topicSubscriberPartitionOff, in.Partition)
	ob.SetInt64(topicSubscriberConnectTimeNsOff, in.ConnectTimeNs)
	ob.SetInt64(topicSubscriberLastSeenTimeNsOff, in.LastSeenTimeNs)
	ob.SetText(topicSubscriberBrokerOff, in.Broker)
	ob.SetBool(topicSubscriberIsActiveOff, in.IsActive)
	ob.SetInt64(topicSubscriberCurrentOffsetOff, in.CurrentOffset)
	ob.SetInt64(topicSubscriberLastReceivedOffsetOff, in.LastReceivedOffset)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetTopicPublishersRequest ---

const (
	getTopicPublishersRequestTopicOff = 0 // message schema_pb.Topic
	getTopicPublishersRequestSize     = 8
)

// GetTopicPublishersRequest is a zero-copy view into a ZAP-encoded
// GetTopicPublishersRequest message.
type GetTopicPublishersRequest struct{ o zap.Object }

// WrapGetTopicPublishersRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapGetTopicPublishersRequest(b []byte) (GetTopicPublishersRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetTopicPublishersRequest{}, err
	}
	return GetTopicPublishersRequest{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t GetTopicPublishersRequest) Topic() []byte {
	return t.o.Bytes(getTopicPublishersRequestTopicOff)
}

// GetTopicPublishersRequestInput collects the field values for
// NewGetTopicPublishersRequest. Topic is the child message's own ZAP buffer
// (from mq_schemawire.NewTopic).
type GetTopicPublishersRequestInput struct {
	Topic []byte
}

// NewGetTopicPublishersRequest builds a ZAP-encoded GetTopicPublishersRequest
// message from in and returns the bytes.
func NewGetTopicPublishersRequest(in GetTopicPublishersRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getTopicPublishersRequestSize)
	ob.SetBytes(getTopicPublishersRequestTopicOff, in.Topic)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetTopicPublishersResponse ---

const (
	getTopicPublishersResponsePublishersOff = 0 // repeated TopicPublisher
	getTopicPublishersResponseSize          = 8
)

// GetTopicPublishersResponse is a zero-copy view into a ZAP-encoded
// GetTopicPublishersResponse message.
type GetTopicPublishersResponse struct{ o zap.Object }

// WrapGetTopicPublishersResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapGetTopicPublishersResponse(b []byte) (GetTopicPublishersResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetTopicPublishersResponse{}, err
	}
	return GetTopicPublishersResponse{o: m.Root()}, nil
}

// PublishersLen reports the number of publishers elements (proto field 1,
// repeated TopicPublisher).
func (t GetTopicPublishersResponse) PublishersLen() int {
	return t.o.List(getTopicPublishersResponsePublishersOff).Len()
}

// PublisherAt returns the i-th publishers element. The bool is false when i is
// out of range or it fails to parse.
func (t GetTopicPublishersResponse) PublisherAt(i int) (TopicPublisher, bool) {
	o := t.o.List(getTopicPublishersResponsePublishersOff).ObjectAt(i)
	if o.IsNull() {
		return TopicPublisher{}, false
	}
	return TopicPublisher{o: o}, true
}

// GetTopicPublishersResponseInput collects the field values for
// NewGetTopicPublishersResponse. Each Publishers entry is a TopicPublisher
// buffer (from NewTopicPublisher).
type GetTopicPublishersResponseInput struct {
	Publishers [][]byte
}

// NewGetTopicPublishersResponse builds a ZAP-encoded GetTopicPublishersResponse
// message from in and returns the bytes.
func NewGetTopicPublishersResponse(in GetTopicPublishersResponseInput) []byte {
	b := zap.NewBuilder(256)
	lb := b.StartList(0)
	for _, e := range in.Publishers {
		lb.AddObjectBytes(e)
	}
	listOff, listLen := lb.Finish()
	ob := b.StartObject(getTopicPublishersResponseSize)
	ob.SetList(getTopicPublishersResponsePublishersOff, listOff, listLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetTopicSubscribersRequest ---

const (
	getTopicSubscribersRequestTopicOff = 0 // message schema_pb.Topic
	getTopicSubscribersRequestSize     = 8
)

// GetTopicSubscribersRequest is a zero-copy view into a ZAP-encoded
// GetTopicSubscribersRequest message.
type GetTopicSubscribersRequest struct{ o zap.Object }

// WrapGetTopicSubscribersRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapGetTopicSubscribersRequest(b []byte) (GetTopicSubscribersRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetTopicSubscribersRequest{}, err
	}
	return GetTopicSubscribersRequest{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t GetTopicSubscribersRequest) Topic() []byte {
	return t.o.Bytes(getTopicSubscribersRequestTopicOff)
}

// GetTopicSubscribersRequestInput collects the field values for
// NewGetTopicSubscribersRequest. Topic is the child message's own ZAP buffer
// (from mq_schemawire.NewTopic).
type GetTopicSubscribersRequestInput struct {
	Topic []byte
}

// NewGetTopicSubscribersRequest builds a ZAP-encoded GetTopicSubscribersRequest
// message from in and returns the bytes.
func NewGetTopicSubscribersRequest(in GetTopicSubscribersRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getTopicSubscribersRequestSize)
	ob.SetBytes(getTopicSubscribersRequestTopicOff, in.Topic)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetTopicSubscribersResponse ---

const (
	getTopicSubscribersResponseSubscribersOff = 0 // repeated TopicSubscriber
	getTopicSubscribersResponseSize           = 8
)

// GetTopicSubscribersResponse is a zero-copy view into a ZAP-encoded
// GetTopicSubscribersResponse message.
type GetTopicSubscribersResponse struct{ o zap.Object }

// WrapGetTopicSubscribersResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapGetTopicSubscribersResponse(b []byte) (GetTopicSubscribersResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetTopicSubscribersResponse{}, err
	}
	return GetTopicSubscribersResponse{o: m.Root()}, nil
}

// SubscribersLen reports the number of subscribers elements (proto field 1,
// repeated TopicSubscriber).
func (t GetTopicSubscribersResponse) SubscribersLen() int {
	return t.o.List(getTopicSubscribersResponseSubscribersOff).Len()
}

// SubscriberAt returns the i-th subscribers element. The bool is false when i
// is out of range or it fails to parse.
func (t GetTopicSubscribersResponse) SubscriberAt(i int) (TopicSubscriber, bool) {
	o := t.o.List(getTopicSubscribersResponseSubscribersOff).ObjectAt(i)
	if o.IsNull() {
		return TopicSubscriber{}, false
	}
	return TopicSubscriber{o: o}, true
}

// GetTopicSubscribersResponseInput collects the field values for
// NewGetTopicSubscribersResponse. Each Subscribers entry is a TopicSubscriber
// buffer (from NewTopicSubscriber).
type GetTopicSubscribersResponseInput struct {
	Subscribers [][]byte
}

// NewGetTopicSubscribersResponse builds a ZAP-encoded
// GetTopicSubscribersResponse message from in and returns the bytes.
func NewGetTopicSubscribersResponse(in GetTopicSubscribersResponseInput) []byte {
	b := zap.NewBuilder(256)
	lb := b.StartList(0)
	for _, e := range in.Subscribers {
		lb.AddObjectBytes(e)
	}
	listOff, listLen := lb.Finish()
	ob := b.StartObject(getTopicSubscribersResponseSize)
	ob.SetList(getTopicSubscribersResponseSubscribersOff, listOff, listLen)
	ob.FinishAsRoot()
	return b.Finish()
}
