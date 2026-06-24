// Code generated from mq_broker.proto; DO NOT EDIT.

package mq_brokerwire

import (
	zap "github.com/zap-proto/go"
)

// Topic control-plane messages. schema_pb.Topic / schema_pb.RecordType singular
// fields are stored as the child's own ZAP buffer (8-byte tail pointer); build
// them with the mq_schema wire package and re-wrap with mq_schemawire.WrapTopic
// / WrapRecordType. Repeated BrokerPartitionAssignment is a same-package list of
// out-of-line objects (read with the typed *At accessor); repeated
// schema_pb.Topic is a list whose elements are raw Topic buffers (re-wrap with
// mq_schemawire.WrapTopic).

// --- TopicRetention ---

const (
	topicRetentionRetentionSecondsOff = 0 // int64
	topicRetentionEnabledOff          = 8 // bool
	topicRetentionSize                = 16
)

// TopicRetention is a zero-copy view into a ZAP-encoded TopicRetention message.
type TopicRetention struct{ o zap.Object }

// WrapTopicRetention parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapTopicRetention(b []byte) (TopicRetention, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TopicRetention{}, err
	}
	return TopicRetention{o: m.Root()}, nil
}

// RetentionSeconds reads the retention_seconds field (proto field 1, int64).
func (t TopicRetention) RetentionSeconds() int64 {
	return t.o.Int64(topicRetentionRetentionSecondsOff)
}

// Enabled reads the enabled field (proto field 2, bool).
func (t TopicRetention) Enabled() bool { return t.o.Bool(topicRetentionEnabledOff) }

// TopicRetentionInput collects the field values for NewTopicRetention.
type TopicRetentionInput struct {
	RetentionSeconds int64
	Enabled          bool
}

// NewTopicRetention builds a ZAP-encoded TopicRetention message from in and
// returns the bytes.
func NewTopicRetention(in TopicRetentionInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(topicRetentionSize)
	ob.SetInt64(topicRetentionRetentionSecondsOff, in.RetentionSeconds)
	ob.SetBool(topicRetentionEnabledOff, in.Enabled)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- BrokerPartitionAssignment ---

const (
	brokerPartitionAssignmentPartitionOff      = 0  // message schema_pb.Partition
	brokerPartitionAssignmentLeaderBrokerOff   = 8  // string
	brokerPartitionAssignmentFollowerBrokerOff = 16 // string
	brokerPartitionAssignmentSize              = 24
)

// BrokerPartitionAssignment is a zero-copy view into a ZAP-encoded
// BrokerPartitionAssignment message.
type BrokerPartitionAssignment struct{ o zap.Object }

// WrapBrokerPartitionAssignment parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapBrokerPartitionAssignment(b []byte) (BrokerPartitionAssignment, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return BrokerPartitionAssignment{}, err
	}
	return BrokerPartitionAssignment{o: m.Root()}, nil
}

// Partition reads the partition field (proto field 1, message
// schema_pb.Partition) as the child's raw ZAP buffer; nil when absent. Re-wrap
// with mq_schemawire.WrapPartition.
func (t BrokerPartitionAssignment) Partition() []byte {
	return t.o.Bytes(brokerPartitionAssignmentPartitionOff)
}

// LeaderBroker reads the leader_broker field (proto field 2, string).
func (t BrokerPartitionAssignment) LeaderBroker() string {
	return t.o.Text(brokerPartitionAssignmentLeaderBrokerOff)
}

// FollowerBroker reads the follower_broker field (proto field 3, string).
func (t BrokerPartitionAssignment) FollowerBroker() string {
	return t.o.Text(brokerPartitionAssignmentFollowerBrokerOff)
}

// BrokerPartitionAssignmentInput collects the field values for
// NewBrokerPartitionAssignment. Partition is the child message's own ZAP buffer
// (from mq_schemawire.NewPartition); nil encodes absent.
type BrokerPartitionAssignmentInput struct {
	Partition      []byte
	LeaderBroker   string
	FollowerBroker string
}

// NewBrokerPartitionAssignment builds a ZAP-encoded BrokerPartitionAssignment
// message from in and returns the bytes.
func NewBrokerPartitionAssignment(in BrokerPartitionAssignmentInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(brokerPartitionAssignmentSize)
	ob.SetBytes(brokerPartitionAssignmentPartitionOff, in.Partition)
	ob.SetText(brokerPartitionAssignmentLeaderBrokerOff, in.LeaderBroker)
	ob.SetText(brokerPartitionAssignmentFollowerBrokerOff, in.FollowerBroker)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ConfigureTopicRequest ---

const (
	configureTopicRequestTopicOff             = 0  // message schema_pb.Topic
	configureTopicRequestPartitionCountOff    = 8  // int32
	configureTopicRequestRetentionOff         = 16 // message TopicRetention
	configureTopicRequestMessageRecordTypeOff = 24 // message schema_pb.RecordType
	configureTopicRequestKeyColumnsOff        = 32 // repeated string
	configureTopicRequestSchemaFormatOff      = 40 // string
	configureTopicRequestSize                 = 48
)

// ConfigureTopicRequest is a zero-copy view into a ZAP-encoded
// ConfigureTopicRequest message.
type ConfigureTopicRequest struct{ o zap.Object }

// WrapConfigureTopicRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapConfigureTopicRequest(b []byte) (ConfigureTopicRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ConfigureTopicRequest{}, err
	}
	return ConfigureTopicRequest{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t ConfigureTopicRequest) Topic() []byte { return t.o.Bytes(configureTopicRequestTopicOff) }

// PartitionCount reads the partition_count field (proto field 2, int32).
func (t ConfigureTopicRequest) PartitionCount() int32 {
	return t.o.Int32(configureTopicRequestPartitionCountOff)
}

// Retention reads the retention field (proto field 3, message TopicRetention).
// The bool is false when the field is absent.
func (t ConfigureTopicRequest) Retention() (TopicRetention, bool) {
	b := t.o.Bytes(configureTopicRequestRetentionOff)
	if len(b) == 0 {
		return TopicRetention{}, false
	}
	r, err := WrapTopicRetention(b)
	if err != nil {
		return TopicRetention{}, false
	}
	return r, true
}

// MessageRecordType reads the message_record_type field (proto field 4, message
// schema_pb.RecordType) as the child's raw ZAP buffer; nil when absent. Re-wrap
// with mq_schemawire.WrapRecordType.
func (t ConfigureTopicRequest) MessageRecordType() []byte {
	return t.o.Bytes(configureTopicRequestMessageRecordTypeOff)
}

// KeyColumnsLen reports the number of key_columns elements (proto field 5,
// repeated string).
func (t ConfigureTopicRequest) KeyColumnsLen() int {
	return t.o.List(configureTopicRequestKeyColumnsOff).Len()
}

// KeyColumnAt returns the i-th key_columns element.
func (t ConfigureTopicRequest) KeyColumnAt(i int) string {
	return stringAt(t.o.List(configureTopicRequestKeyColumnsOff), i)
}

// SchemaFormat reads the schema_format field (proto field 6, string).
func (t ConfigureTopicRequest) SchemaFormat() string {
	return t.o.Text(configureTopicRequestSchemaFormatOff)
}

// ConfigureTopicRequestInput collects the field values for
// NewConfigureTopicRequest. Topic, Retention and MessageRecordType are each the
// child message's own ZAP buffer (from mq_schemawire.NewTopic /
// NewTopicRetention / mq_schemawire.NewRecordType); nil encodes absent.
type ConfigureTopicRequestInput struct {
	Topic             []byte
	PartitionCount    int32
	Retention         []byte
	MessageRecordType []byte
	KeyColumns        []string
	SchemaFormat      string
}

// NewConfigureTopicRequest builds a ZAP-encoded ConfigureTopicRequest message
// from in and returns the bytes.
func NewConfigureTopicRequest(in ConfigureTopicRequestInput) []byte {
	b := zap.NewBuilder(512)
	keyColumnsOff, keyColumnsLen := setStringList(b, in.KeyColumns)
	ob := b.StartObject(configureTopicRequestSize)
	ob.SetBytes(configureTopicRequestTopicOff, in.Topic)
	ob.SetInt32(configureTopicRequestPartitionCountOff, in.PartitionCount)
	ob.SetBytes(configureTopicRequestRetentionOff, in.Retention)
	ob.SetBytes(configureTopicRequestMessageRecordTypeOff, in.MessageRecordType)
	ob.SetList(configureTopicRequestKeyColumnsOff, keyColumnsOff, keyColumnsLen)
	ob.SetText(configureTopicRequestSchemaFormatOff, in.SchemaFormat)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ConfigureTopicResponse ---

const (
	configureTopicResponseBrokerPartitionAssignmentsOff = 0  // repeated BrokerPartitionAssignment
	configureTopicResponseRetentionOff                  = 8  // message TopicRetention
	configureTopicResponseMessageRecordTypeOff          = 16 // message schema_pb.RecordType
	configureTopicResponseKeyColumnsOff                 = 24 // repeated string
	configureTopicResponseSchemaFormatOff               = 32 // string
	configureTopicResponseSize                          = 40
)

// ConfigureTopicResponse is a zero-copy view into a ZAP-encoded
// ConfigureTopicResponse message.
type ConfigureTopicResponse struct{ o zap.Object }

// WrapConfigureTopicResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapConfigureTopicResponse(b []byte) (ConfigureTopicResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ConfigureTopicResponse{}, err
	}
	return ConfigureTopicResponse{o: m.Root()}, nil
}

// BrokerPartitionAssignmentsLen reports the number of
// broker_partition_assignments elements (proto field 2, repeated
// BrokerPartitionAssignment).
func (t ConfigureTopicResponse) BrokerPartitionAssignmentsLen() int {
	return t.o.List(configureTopicResponseBrokerPartitionAssignmentsOff).Len()
}

// BrokerPartitionAssignmentAt returns the i-th broker_partition_assignments
// element. The bool is false when i is out of range or it fails to parse.
func (t ConfigureTopicResponse) BrokerPartitionAssignmentAt(i int) (BrokerPartitionAssignment, bool) {
	o := t.o.List(configureTopicResponseBrokerPartitionAssignmentsOff).ObjectAt(i)
	if o.IsNull() {
		return BrokerPartitionAssignment{}, false
	}
	return BrokerPartitionAssignment{o: o}, true
}

// Retention reads the retention field (proto field 3, message TopicRetention).
// The bool is false when the field is absent.
func (t ConfigureTopicResponse) Retention() (TopicRetention, bool) {
	b := t.o.Bytes(configureTopicResponseRetentionOff)
	if len(b) == 0 {
		return TopicRetention{}, false
	}
	r, err := WrapTopicRetention(b)
	if err != nil {
		return TopicRetention{}, false
	}
	return r, true
}

// MessageRecordType reads the message_record_type field (proto field 4, message
// schema_pb.RecordType) as the child's raw ZAP buffer; nil when absent. Re-wrap
// with mq_schemawire.WrapRecordType.
func (t ConfigureTopicResponse) MessageRecordType() []byte {
	return t.o.Bytes(configureTopicResponseMessageRecordTypeOff)
}

// KeyColumnsLen reports the number of key_columns elements (proto field 5,
// repeated string).
func (t ConfigureTopicResponse) KeyColumnsLen() int {
	return t.o.List(configureTopicResponseKeyColumnsOff).Len()
}

// KeyColumnAt returns the i-th key_columns element.
func (t ConfigureTopicResponse) KeyColumnAt(i int) string {
	return stringAt(t.o.List(configureTopicResponseKeyColumnsOff), i)
}

// SchemaFormat reads the schema_format field (proto field 6, string).
func (t ConfigureTopicResponse) SchemaFormat() string {
	return t.o.Text(configureTopicResponseSchemaFormatOff)
}

// ConfigureTopicResponseInput collects the field values for
// NewConfigureTopicResponse. Each BrokerPartitionAssignments entry is a
// BrokerPartitionAssignment buffer (from NewBrokerPartitionAssignment).
// Retention and MessageRecordType are each the child message's own ZAP buffer.
type ConfigureTopicResponseInput struct {
	BrokerPartitionAssignments [][]byte
	Retention                  []byte
	MessageRecordType          []byte
	KeyColumns                 []string
	SchemaFormat               string
}

// NewConfigureTopicResponse builds a ZAP-encoded ConfigureTopicResponse message
// from in and returns the bytes.
func NewConfigureTopicResponse(in ConfigureTopicResponseInput) []byte {
	b := zap.NewBuilder(512)
	bpaLB := b.StartList(0)
	for _, e := range in.BrokerPartitionAssignments {
		bpaLB.AddObjectBytes(e)
	}
	bpaOff, bpaLen := bpaLB.Finish()
	keyColumnsOff, keyColumnsLen := setStringList(b, in.KeyColumns)
	ob := b.StartObject(configureTopicResponseSize)
	ob.SetList(configureTopicResponseBrokerPartitionAssignmentsOff, bpaOff, bpaLen)
	ob.SetBytes(configureTopicResponseRetentionOff, in.Retention)
	ob.SetBytes(configureTopicResponseMessageRecordTypeOff, in.MessageRecordType)
	ob.SetList(configureTopicResponseKeyColumnsOff, keyColumnsOff, keyColumnsLen)
	ob.SetText(configureTopicResponseSchemaFormatOff, in.SchemaFormat)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListTopicsRequest ---

const listTopicsRequestSize = 0

// ListTopicsRequest is a zero-copy view into a ZAP-encoded ListTopicsRequest
// message. It carries no fields.
type ListTopicsRequest struct{ o zap.Object }

// WrapListTopicsRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapListTopicsRequest(b []byte) (ListTopicsRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListTopicsRequest{}, err
	}
	return ListTopicsRequest{o: m.Root()}, nil
}

// ListTopicsRequestInput collects the field values for NewListTopicsRequest. It
// has no fields.
type ListTopicsRequestInput struct{}

// NewListTopicsRequest builds a ZAP-encoded ListTopicsRequest message and
// returns the bytes.
func NewListTopicsRequest(in ListTopicsRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(listTopicsRequestSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListTopicsResponse ---

const (
	listTopicsResponseTopicsOff = 0 // repeated schema_pb.Topic
	listTopicsResponseSize      = 8
)

// ListTopicsResponse is a zero-copy view into a ZAP-encoded ListTopicsResponse
// message.
type ListTopicsResponse struct{ o zap.Object }

// WrapListTopicsResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapListTopicsResponse(b []byte) (ListTopicsResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListTopicsResponse{}, err
	}
	return ListTopicsResponse{o: m.Root()}, nil
}

// TopicsLen reports the number of topics elements (proto field 1, repeated
// schema_pb.Topic).
func (t ListTopicsResponse) TopicsLen() int { return t.o.List(listTopicsResponseTopicsOff).Len() }

// TopicAt returns the i-th topics element as the child's raw ZAP buffer; nil
// when i is out of range. Re-wrap with mq_schemawire.WrapTopic.
func (t ListTopicsResponse) TopicAt(i int) []byte {
	return t.o.List(listTopicsResponseTopicsOff).BytesAt(i)
}

// ListTopicsResponseInput collects the field values for NewListTopicsResponse.
// Each Topics entry is a Topic buffer (from mq_schemawire.NewTopic).
type ListTopicsResponseInput struct {
	Topics [][]byte
}

// NewListTopicsResponse builds a ZAP-encoded ListTopicsResponse message from in
// and returns the bytes.
func NewListTopicsResponse(in ListTopicsResponseInput) []byte {
	b := zap.NewBuilder(256)
	lb := b.StartList(0)
	for _, e := range in.Topics {
		lb.AddObjectBytes(e)
	}
	listOff, listLen := lb.Finish()
	ob := b.StartObject(listTopicsResponseSize)
	ob.SetList(listTopicsResponseTopicsOff, listOff, listLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TopicExistsRequest ---

const (
	topicExistsRequestTopicOff = 0 // message schema_pb.Topic
	topicExistsRequestSize     = 8
)

// TopicExistsRequest is a zero-copy view into a ZAP-encoded TopicExistsRequest
// message.
type TopicExistsRequest struct{ o zap.Object }

// WrapTopicExistsRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapTopicExistsRequest(b []byte) (TopicExistsRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TopicExistsRequest{}, err
	}
	return TopicExistsRequest{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t TopicExistsRequest) Topic() []byte { return t.o.Bytes(topicExistsRequestTopicOff) }

// TopicExistsRequestInput collects the field values for NewTopicExistsRequest.
// Topic is the child message's own ZAP buffer (from mq_schemawire.NewTopic).
type TopicExistsRequestInput struct {
	Topic []byte
}

// NewTopicExistsRequest builds a ZAP-encoded TopicExistsRequest message from in
// and returns the bytes.
func NewTopicExistsRequest(in TopicExistsRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(topicExistsRequestSize)
	ob.SetBytes(topicExistsRequestTopicOff, in.Topic)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TopicExistsResponse ---

const (
	topicExistsResponseExistsOff = 0 // bool
	topicExistsResponseSize      = 1
)

// TopicExistsResponse is a zero-copy view into a ZAP-encoded TopicExistsResponse
// message.
type TopicExistsResponse struct{ o zap.Object }

// WrapTopicExistsResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapTopicExistsResponse(b []byte) (TopicExistsResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TopicExistsResponse{}, err
	}
	return TopicExistsResponse{o: m.Root()}, nil
}

// Exists reads the exists field (proto field 1, bool).
func (t TopicExistsResponse) Exists() bool { return t.o.Bool(topicExistsResponseExistsOff) }

// TopicExistsResponseInput collects the field values for
// NewTopicExistsResponse.
type TopicExistsResponseInput struct {
	Exists bool
}

// NewTopicExistsResponse builds a ZAP-encoded TopicExistsResponse message from
// in and returns the bytes.
func NewTopicExistsResponse(in TopicExistsResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(topicExistsResponseSize)
	ob.SetBool(topicExistsResponseExistsOff, in.Exists)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LookupTopicBrokersRequest ---

const (
	lookupTopicBrokersRequestTopicOff = 0 // message schema_pb.Topic
	lookupTopicBrokersRequestSize     = 8
)

// LookupTopicBrokersRequest is a zero-copy view into a ZAP-encoded
// LookupTopicBrokersRequest message.
type LookupTopicBrokersRequest struct{ o zap.Object }

// WrapLookupTopicBrokersRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapLookupTopicBrokersRequest(b []byte) (LookupTopicBrokersRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LookupTopicBrokersRequest{}, err
	}
	return LookupTopicBrokersRequest{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t LookupTopicBrokersRequest) Topic() []byte {
	return t.o.Bytes(lookupTopicBrokersRequestTopicOff)
}

// LookupTopicBrokersRequestInput collects the field values for
// NewLookupTopicBrokersRequest. Topic is the child message's own ZAP buffer
// (from mq_schemawire.NewTopic).
type LookupTopicBrokersRequestInput struct {
	Topic []byte
}

// NewLookupTopicBrokersRequest builds a ZAP-encoded LookupTopicBrokersRequest
// message from in and returns the bytes.
func NewLookupTopicBrokersRequest(in LookupTopicBrokersRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(lookupTopicBrokersRequestSize)
	ob.SetBytes(lookupTopicBrokersRequestTopicOff, in.Topic)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LookupTopicBrokersResponse ---

const (
	lookupTopicBrokersResponseTopicOff                      = 0 // message schema_pb.Topic
	lookupTopicBrokersResponseBrokerPartitionAssignmentsOff = 8 // repeated BrokerPartitionAssignment
	lookupTopicBrokersResponseSize                          = 16
)

// LookupTopicBrokersResponse is a zero-copy view into a ZAP-encoded
// LookupTopicBrokersResponse message.
type LookupTopicBrokersResponse struct{ o zap.Object }

// WrapLookupTopicBrokersResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapLookupTopicBrokersResponse(b []byte) (LookupTopicBrokersResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LookupTopicBrokersResponse{}, err
	}
	return LookupTopicBrokersResponse{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t LookupTopicBrokersResponse) Topic() []byte {
	return t.o.Bytes(lookupTopicBrokersResponseTopicOff)
}

// BrokerPartitionAssignmentsLen reports the number of
// broker_partition_assignments elements (proto field 2, repeated
// BrokerPartitionAssignment).
func (t LookupTopicBrokersResponse) BrokerPartitionAssignmentsLen() int {
	return t.o.List(lookupTopicBrokersResponseBrokerPartitionAssignmentsOff).Len()
}

// BrokerPartitionAssignmentAt returns the i-th broker_partition_assignments
// element. The bool is false when i is out of range or it fails to parse.
func (t LookupTopicBrokersResponse) BrokerPartitionAssignmentAt(i int) (BrokerPartitionAssignment, bool) {
	o := t.o.List(lookupTopicBrokersResponseBrokerPartitionAssignmentsOff).ObjectAt(i)
	if o.IsNull() {
		return BrokerPartitionAssignment{}, false
	}
	return BrokerPartitionAssignment{o: o}, true
}

// LookupTopicBrokersResponseInput collects the field values for
// NewLookupTopicBrokersResponse. Topic is the child message's own ZAP buffer
// (from mq_schemawire.NewTopic). Each BrokerPartitionAssignments entry is a
// BrokerPartitionAssignment buffer (from NewBrokerPartitionAssignment).
type LookupTopicBrokersResponseInput struct {
	Topic                      []byte
	BrokerPartitionAssignments [][]byte
}

// NewLookupTopicBrokersResponse builds a ZAP-encoded LookupTopicBrokersResponse
// message from in and returns the bytes.
func NewLookupTopicBrokersResponse(in LookupTopicBrokersResponseInput) []byte {
	b := zap.NewBuilder(256)
	lb := b.StartList(0)
	for _, e := range in.BrokerPartitionAssignments {
		lb.AddObjectBytes(e)
	}
	listOff, listLen := lb.Finish()
	ob := b.StartObject(lookupTopicBrokersResponseSize)
	ob.SetBytes(lookupTopicBrokersResponseTopicOff, in.Topic)
	ob.SetList(lookupTopicBrokersResponseBrokerPartitionAssignmentsOff, listOff, listLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetTopicConfigurationRequest ---

const (
	getTopicConfigurationRequestTopicOff = 0 // message schema_pb.Topic
	getTopicConfigurationRequestSize     = 8
)

// GetTopicConfigurationRequest is a zero-copy view into a ZAP-encoded
// GetTopicConfigurationRequest message.
type GetTopicConfigurationRequest struct{ o zap.Object }

// WrapGetTopicConfigurationRequest parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapGetTopicConfigurationRequest(b []byte) (GetTopicConfigurationRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetTopicConfigurationRequest{}, err
	}
	return GetTopicConfigurationRequest{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t GetTopicConfigurationRequest) Topic() []byte {
	return t.o.Bytes(getTopicConfigurationRequestTopicOff)
}

// GetTopicConfigurationRequestInput collects the field values for
// NewGetTopicConfigurationRequest. Topic is the child message's own ZAP buffer
// (from mq_schemawire.NewTopic).
type GetTopicConfigurationRequestInput struct {
	Topic []byte
}

// NewGetTopicConfigurationRequest builds a ZAP-encoded
// GetTopicConfigurationRequest message from in and returns the bytes.
func NewGetTopicConfigurationRequest(in GetTopicConfigurationRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getTopicConfigurationRequestSize)
	ob.SetBytes(getTopicConfigurationRequestTopicOff, in.Topic)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetTopicConfigurationResponse ---

const (
	getTopicConfigurationResponseTopicOff                      = 0  // message schema_pb.Topic
	getTopicConfigurationResponsePartitionCountOff             = 8  // int32
	getTopicConfigurationResponseBrokerPartitionAssignmentsOff = 16 // repeated BrokerPartitionAssignment
	getTopicConfigurationResponseCreatedAtNsOff                = 24 // int64
	getTopicConfigurationResponseLastUpdatedNsOff              = 32 // int64
	getTopicConfigurationResponseRetentionOff                  = 40 // message TopicRetention
	getTopicConfigurationResponseMessageRecordTypeOff          = 48 // message schema_pb.RecordType
	getTopicConfigurationResponseKeyColumnsOff                 = 56 // repeated string
	getTopicConfigurationResponseSchemaFormatOff               = 64 // string
	getTopicConfigurationResponseSize                          = 72
)

// GetTopicConfigurationResponse is a zero-copy view into a ZAP-encoded
// GetTopicConfigurationResponse message.
type GetTopicConfigurationResponse struct{ o zap.Object }

// WrapGetTopicConfigurationResponse parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapGetTopicConfigurationResponse(b []byte) (GetTopicConfigurationResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetTopicConfigurationResponse{}, err
	}
	return GetTopicConfigurationResponse{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t GetTopicConfigurationResponse) Topic() []byte {
	return t.o.Bytes(getTopicConfigurationResponseTopicOff)
}

// PartitionCount reads the partition_count field (proto field 2, int32).
func (t GetTopicConfigurationResponse) PartitionCount() int32 {
	return t.o.Int32(getTopicConfigurationResponsePartitionCountOff)
}

// BrokerPartitionAssignmentsLen reports the number of
// broker_partition_assignments elements (proto field 3, repeated
// BrokerPartitionAssignment).
func (t GetTopicConfigurationResponse) BrokerPartitionAssignmentsLen() int {
	return t.o.List(getTopicConfigurationResponseBrokerPartitionAssignmentsOff).Len()
}

// BrokerPartitionAssignmentAt returns the i-th broker_partition_assignments
// element. The bool is false when i is out of range or it fails to parse.
func (t GetTopicConfigurationResponse) BrokerPartitionAssignmentAt(i int) (BrokerPartitionAssignment, bool) {
	o := t.o.List(getTopicConfigurationResponseBrokerPartitionAssignmentsOff).ObjectAt(i)
	if o.IsNull() {
		return BrokerPartitionAssignment{}, false
	}
	return BrokerPartitionAssignment{o: o}, true
}

// CreatedAtNs reads the created_at_ns field (proto field 4, int64).
func (t GetTopicConfigurationResponse) CreatedAtNs() int64 {
	return t.o.Int64(getTopicConfigurationResponseCreatedAtNsOff)
}

// LastUpdatedNs reads the last_updated_ns field (proto field 5, int64).
func (t GetTopicConfigurationResponse) LastUpdatedNs() int64 {
	return t.o.Int64(getTopicConfigurationResponseLastUpdatedNsOff)
}

// Retention reads the retention field (proto field 6, message TopicRetention).
// The bool is false when the field is absent.
func (t GetTopicConfigurationResponse) Retention() (TopicRetention, bool) {
	b := t.o.Bytes(getTopicConfigurationResponseRetentionOff)
	if len(b) == 0 {
		return TopicRetention{}, false
	}
	r, err := WrapTopicRetention(b)
	if err != nil {
		return TopicRetention{}, false
	}
	return r, true
}

// MessageRecordType reads the message_record_type field (proto field 7, message
// schema_pb.RecordType) as the child's raw ZAP buffer; nil when absent. Re-wrap
// with mq_schemawire.WrapRecordType.
func (t GetTopicConfigurationResponse) MessageRecordType() []byte {
	return t.o.Bytes(getTopicConfigurationResponseMessageRecordTypeOff)
}

// KeyColumnsLen reports the number of key_columns elements (proto field 8,
// repeated string).
func (t GetTopicConfigurationResponse) KeyColumnsLen() int {
	return t.o.List(getTopicConfigurationResponseKeyColumnsOff).Len()
}

// KeyColumnAt returns the i-th key_columns element.
func (t GetTopicConfigurationResponse) KeyColumnAt(i int) string {
	return stringAt(t.o.List(getTopicConfigurationResponseKeyColumnsOff), i)
}

// SchemaFormat reads the schema_format field (proto field 9, string).
func (t GetTopicConfigurationResponse) SchemaFormat() string {
	return t.o.Text(getTopicConfigurationResponseSchemaFormatOff)
}

// GetTopicConfigurationResponseInput collects the field values for
// NewGetTopicConfigurationResponse. Topic, Retention and MessageRecordType are
// each the child message's own ZAP buffer; each BrokerPartitionAssignments
// entry is a BrokerPartitionAssignment buffer.
type GetTopicConfigurationResponseInput struct {
	Topic                      []byte
	PartitionCount             int32
	BrokerPartitionAssignments [][]byte
	CreatedAtNs                int64
	LastUpdatedNs              int64
	Retention                  []byte
	MessageRecordType          []byte
	KeyColumns                 []string
	SchemaFormat               string
}

// NewGetTopicConfigurationResponse builds a ZAP-encoded
// GetTopicConfigurationResponse message from in and returns the bytes.
func NewGetTopicConfigurationResponse(in GetTopicConfigurationResponseInput) []byte {
	b := zap.NewBuilder(512)
	bpaLB := b.StartList(0)
	for _, e := range in.BrokerPartitionAssignments {
		bpaLB.AddObjectBytes(e)
	}
	bpaOff, bpaLen := bpaLB.Finish()
	keyColumnsOff, keyColumnsLen := setStringList(b, in.KeyColumns)
	ob := b.StartObject(getTopicConfigurationResponseSize)
	ob.SetBytes(getTopicConfigurationResponseTopicOff, in.Topic)
	ob.SetInt32(getTopicConfigurationResponsePartitionCountOff, in.PartitionCount)
	ob.SetList(getTopicConfigurationResponseBrokerPartitionAssignmentsOff, bpaOff, bpaLen)
	ob.SetInt64(getTopicConfigurationResponseCreatedAtNsOff, in.CreatedAtNs)
	ob.SetInt64(getTopicConfigurationResponseLastUpdatedNsOff, in.LastUpdatedNs)
	ob.SetBytes(getTopicConfigurationResponseRetentionOff, in.Retention)
	ob.SetBytes(getTopicConfigurationResponseMessageRecordTypeOff, in.MessageRecordType)
	ob.SetList(getTopicConfigurationResponseKeyColumnsOff, keyColumnsOff, keyColumnsLen)
	ob.SetText(getTopicConfigurationResponseSchemaFormatOff, in.SchemaFormat)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- AssignTopicPartitionsRequest ---

const (
	assignTopicPartitionsRequestTopicOff                      = 0  // message schema_pb.Topic
	assignTopicPartitionsRequestBrokerPartitionAssignmentsOff = 8  // repeated BrokerPartitionAssignment
	assignTopicPartitionsRequestIsLeaderOff                   = 16 // bool
	assignTopicPartitionsRequestIsDrainingOff                 = 17 // bool
	assignTopicPartitionsRequestSize                          = 24
)

// AssignTopicPartitionsRequest is a zero-copy view into a ZAP-encoded
// AssignTopicPartitionsRequest message.
type AssignTopicPartitionsRequest struct{ o zap.Object }

// WrapAssignTopicPartitionsRequest parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapAssignTopicPartitionsRequest(b []byte) (AssignTopicPartitionsRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AssignTopicPartitionsRequest{}, err
	}
	return AssignTopicPartitionsRequest{o: m.Root()}, nil
}

// Topic reads the topic field (proto field 1, message schema_pb.Topic) as the
// child's raw ZAP buffer; nil when absent. Re-wrap with mq_schemawire.WrapTopic.
func (t AssignTopicPartitionsRequest) Topic() []byte {
	return t.o.Bytes(assignTopicPartitionsRequestTopicOff)
}

// BrokerPartitionAssignmentsLen reports the number of
// broker_partition_assignments elements (proto field 2, repeated
// BrokerPartitionAssignment).
func (t AssignTopicPartitionsRequest) BrokerPartitionAssignmentsLen() int {
	return t.o.List(assignTopicPartitionsRequestBrokerPartitionAssignmentsOff).Len()
}

// BrokerPartitionAssignmentAt returns the i-th broker_partition_assignments
// element. The bool is false when i is out of range or it fails to parse.
func (t AssignTopicPartitionsRequest) BrokerPartitionAssignmentAt(i int) (BrokerPartitionAssignment, bool) {
	o := t.o.List(assignTopicPartitionsRequestBrokerPartitionAssignmentsOff).ObjectAt(i)
	if o.IsNull() {
		return BrokerPartitionAssignment{}, false
	}
	return BrokerPartitionAssignment{o: o}, true
}

// IsLeader reads the is_leader field (proto field 3, bool).
func (t AssignTopicPartitionsRequest) IsLeader() bool {
	return t.o.Bool(assignTopicPartitionsRequestIsLeaderOff)
}

// IsDraining reads the is_draining field (proto field 4, bool).
func (t AssignTopicPartitionsRequest) IsDraining() bool {
	return t.o.Bool(assignTopicPartitionsRequestIsDrainingOff)
}

// AssignTopicPartitionsRequestInput collects the field values for
// NewAssignTopicPartitionsRequest. Topic is the child message's own ZAP buffer
// (from mq_schemawire.NewTopic). Each BrokerPartitionAssignments entry is a
// BrokerPartitionAssignment buffer (from NewBrokerPartitionAssignment).
type AssignTopicPartitionsRequestInput struct {
	Topic                      []byte
	BrokerPartitionAssignments [][]byte
	IsLeader                   bool
	IsDraining                 bool
}

// NewAssignTopicPartitionsRequest builds a ZAP-encoded
// AssignTopicPartitionsRequest message from in and returns the bytes.
func NewAssignTopicPartitionsRequest(in AssignTopicPartitionsRequestInput) []byte {
	b := zap.NewBuilder(256)
	lb := b.StartList(0)
	for _, e := range in.BrokerPartitionAssignments {
		lb.AddObjectBytes(e)
	}
	listOff, listLen := lb.Finish()
	ob := b.StartObject(assignTopicPartitionsRequestSize)
	ob.SetBytes(assignTopicPartitionsRequestTopicOff, in.Topic)
	ob.SetList(assignTopicPartitionsRequestBrokerPartitionAssignmentsOff, listOff, listLen)
	ob.SetBool(assignTopicPartitionsRequestIsLeaderOff, in.IsLeader)
	ob.SetBool(assignTopicPartitionsRequestIsDrainingOff, in.IsDraining)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- AssignTopicPartitionsResponse ---

const assignTopicPartitionsResponseSize = 0

// AssignTopicPartitionsResponse is a zero-copy view into a ZAP-encoded
// AssignTopicPartitionsResponse message. It carries no fields.
type AssignTopicPartitionsResponse struct{ o zap.Object }

// WrapAssignTopicPartitionsResponse parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapAssignTopicPartitionsResponse(b []byte) (AssignTopicPartitionsResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AssignTopicPartitionsResponse{}, err
	}
	return AssignTopicPartitionsResponse{o: m.Root()}, nil
}

// AssignTopicPartitionsResponseInput collects the field values for
// NewAssignTopicPartitionsResponse. It has no fields.
type AssignTopicPartitionsResponseInput struct{}

// NewAssignTopicPartitionsResponse builds a ZAP-encoded
// AssignTopicPartitionsResponse message and returns the bytes.
func NewAssignTopicPartitionsResponse(in AssignTopicPartitionsResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(assignTopicPartitionsResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}
