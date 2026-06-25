// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package mq_brokerwire

import (
	"bytes"
	"io"
	"sync"
	"testing"

	"github.com/zap-proto/go/transport"
)

// fakeBroker is an in-memory Broker + StreamBroker for the round-trip tests. The
// unary side records the last request it saw (read out of the zero-copy view)
// and returns fixed values, so the test can assert every field — scalars, the
// nested Topic/Partition opaque child buffers, and the repeated DataMessage list
// — survived the wire crossing. The streaming side echoes a PublishMessage ack
// per data frame it receives. Fields touched from both the server goroutine and
// the test goroutine are guarded by a mutex.
type fakeBroker struct {
	mu sync.Mutex

	gotFilerGroup string

	gotFetchTopic     []byte
	gotFetchPartition []byte
	gotFetchOffset    int64
	gotFetchMaxMsgs   int32
	gotFetchGroup     string

	// streaming: every value frame the server observed on PublishMessage.
	gotPublishValues [][]byte
}

// --- unary Broker ---

func (f *fakeBroker) FindBrokerLeader(req FindBrokerLeaderRequest) (FindBrokerLeaderResponseInput, error) {
	f.mu.Lock()
	f.gotFilerGroup = req.FilerGroup()
	f.mu.Unlock()
	return FindBrokerLeaderResponseInput{Broker: "broker-1.hanzo.svc:17777"}, nil
}

func (f *fakeBroker) FetchMessage(req FetchMessageRequest) (FetchMessageResponseInput, error) {
	f.mu.Lock()
	f.gotFetchTopic = append([]byte(nil), req.Topic()...)
	f.gotFetchPartition = append([]byte(nil), req.Partition()...)
	f.gotFetchOffset = req.StartOffset()
	f.gotFetchMaxMsgs = req.MaxMessages()
	f.gotFetchGroup = req.ConsumerGroup()
	f.mu.Unlock()

	// Return two messages as a repeated DataMessage list.
	msg0 := NewDataMessage(DataMessageInput{Key: []byte("k0"), Value: []byte("v0"), TsNs: 1000})
	msg1 := NewDataMessage(DataMessageInput{Key: []byte("k1"), Value: []byte("v1"), TsNs: 2000})
	return FetchMessageResponseInput{
		Messages:       [][]byte{msg0, msg1},
		HighWaterMark:  42,
		LogStartOffset: 7,
		EndOfPartition: true,
		NextOffset:     req.StartOffset() + 2,
	}, nil
}

// The remaining unary methods are unused by these tests; satisfy the interface
// with empty successful responses so fakeBroker is a complete Broker.
func (f *fakeBroker) BalanceTopics(BalanceTopicsRequest) (BalanceTopicsResponseInput, error) {
	return BalanceTopicsResponseInput{}, nil
}
func (f *fakeBroker) ListTopics(ListTopicsRequest) (ListTopicsResponseInput, error) {
	return ListTopicsResponseInput{}, nil
}
func (f *fakeBroker) TopicExists(TopicExistsRequest) (TopicExistsResponseInput, error) {
	return TopicExistsResponseInput{Exists: true}, nil
}
func (f *fakeBroker) ConfigureTopic(ConfigureTopicRequest) (ConfigureTopicResponseInput, error) {
	return ConfigureTopicResponseInput{}, nil
}
func (f *fakeBroker) LookupTopicBrokers(LookupTopicBrokersRequest) (LookupTopicBrokersResponseInput, error) {
	return LookupTopicBrokersResponseInput{}, nil
}
func (f *fakeBroker) GetTopicConfiguration(GetTopicConfigurationRequest) (GetTopicConfigurationResponseInput, error) {
	return GetTopicConfigurationResponseInput{}, nil
}
func (f *fakeBroker) GetTopicPublishers(GetTopicPublishersRequest) (GetTopicPublishersResponseInput, error) {
	return GetTopicPublishersResponseInput{}, nil
}
func (f *fakeBroker) GetTopicSubscribers(GetTopicSubscribersRequest) (GetTopicSubscribersResponseInput, error) {
	return GetTopicSubscribersResponseInput{}, nil
}
func (f *fakeBroker) AssignTopicPartitions(AssignTopicPartitionsRequest) (AssignTopicPartitionsResponseInput, error) {
	return AssignTopicPartitionsResponseInput{}, nil
}
func (f *fakeBroker) ClosePublishers(ClosePublishersRequest) (ClosePublishersResponseInput, error) {
	return ClosePublishersResponseInput{}, nil
}
func (f *fakeBroker) CloseSubscribers(CloseSubscribersRequest) (CloseSubscribersResponseInput, error) {
	return CloseSubscribersResponseInput{}, nil
}
func (f *fakeBroker) GetPartitionRangeInfo(GetPartitionRangeInfoRequest) (GetPartitionRangeInfoResponseInput, error) {
	return GetPartitionRangeInfoResponseInput{}, nil
}

// --- streaming StreamBroker ---

// ServeStream serves PublishMessage as a bidi stream: for the init frame it
// sends one ack with the publisher's assigned base offset, and for every data
// frame it records the value and sends a per-message ack. Other streaming
// ordinals are not exercised by these tests.
func (f *fakeBroker) ServeStream(method uint32, init []byte, s transport.Stream) {
	switch method {
	case HanzoMessagingPublishMessageOrdinal:
		f.servePublish(init, s)
	default:
		// unexercised streaming RPC: end immediately.
	}
}

func (f *fakeBroker) servePublish(init []byte, s transport.Stream) {
	// The opening frame is a PublishMessageRequest carrying the InitMessage
	// variant; ack it with the base offset.
	_ = init
	if err := s.Send(NewPublishMessageResponse(PublishMessageResponseInput{AssignedOffset: 100})); err != nil {
		return
	}
	for {
		frame, err := s.Recv()
		if err == io.EOF {
			return // client half-closed
		}
		if err != nil {
			return
		}
		req, err := WrapPublishMessageRequest(frame)
		if err != nil {
			return
		}
		if req.WhichMessage() == PublishMessageRequestMessageData {
			f.mu.Lock()
			f.gotPublishValues = append(f.gotPublishValues, append([]byte(nil), req.Data().Value()...))
			n := int64(len(f.gotPublishValues))
			f.mu.Unlock()
			if err := s.Send(NewPublishMessageResponse(PublishMessageResponseInput{
				AckTsNs:        req.Data().TsNs(),
				AssignedOffset: 100 + n,
			})); err != nil {
				return
			}
		}
	}
}

// TestUnaryRoundTrip proves FindBrokerLeader and FetchMessage over the canonical
// github.com/zap-proto/go transport across a real TCP socket: the requests cross
// the wire as ZAP RPC envelopes carrying zero-copy mq_brokerwire payloads
// (scalars, nested Topic/Partition opaque buffers, and a repeated DataMessage
// list), the server dispatches them to the backend, and the responses come back
// the same way — no HTTP, no protobuf, no struct marshaling.
func TestUnaryRoundTrip(t *testing.T) {
	b := &fakeBroker{}

	srv, err := Serve("tcp", "127.0.0.1:0", b, b)
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	// FindBrokerLeader: scalar string in, scalar string out.
	leaderResp, err := cli.FindBrokerLeader(FindBrokerLeaderRequestInput{FilerGroup: "fg-east"})
	if err != nil {
		t.Fatalf("FindBrokerLeader: %v", err)
	}
	if leaderResp.Broker() != "broker-1.hanzo.svc:17777" {
		t.Fatalf("FindBrokerLeader broker = %q", leaderResp.Broker())
	}

	// FetchMessage: nested Topic/Partition (built with their own wire pkg, here
	// opaque non-empty buffers) + repeated DataMessage list back.
	topic := []byte("topic-buffer-stand-in")
	partition := []byte("partition-buffer-stand-in")
	fetchResp, err := cli.FetchMessage(FetchMessageRequestInput{
		Topic:         topic,
		Partition:     partition,
		StartOffset:   10,
		MaxMessages:   16,
		ConsumerGroup: "cg-readers",
	})
	if err != nil {
		t.Fatalf("FetchMessage: %v", err)
	}

	// Response scalars.
	if fetchResp.HighWaterMark() != 42 || fetchResp.LogStartOffset() != 7 ||
		!fetchResp.EndOfPartition() || fetchResp.NextOffset() != 12 {
		t.Fatalf("FetchMessage scalars: hwm=%d lso=%d eop=%v next=%d",
			fetchResp.HighWaterMark(), fetchResp.LogStartOffset(),
			fetchResp.EndOfPartition(), fetchResp.NextOffset())
	}
	// Repeated DataMessage list survived: two messages, fields intact.
	if fetchResp.MessagesLen() != 2 {
		t.Fatalf("FetchMessage MessagesLen = %d, want 2", fetchResp.MessagesLen())
	}
	m0, ok0 := fetchResp.MessageAt(0)
	m1, ok1 := fetchResp.MessageAt(1)
	if !ok0 || !ok1 {
		t.Fatalf("FetchMessage MessageAt ok = %v,%v", ok0, ok1)
	}
	if !bytes.Equal(m0.Key(), []byte("k0")) || !bytes.Equal(m0.Value(), []byte("v0")) || m0.TsNs() != 1000 {
		t.Fatalf("msg0 = key=%q val=%q ts=%d", m0.Key(), m0.Value(), m0.TsNs())
	}
	if !bytes.Equal(m1.Key(), []byte("k1")) || !bytes.Equal(m1.Value(), []byte("v1")) || m1.TsNs() != 2000 {
		t.Fatalf("msg1 = key=%q val=%q ts=%d", m1.Key(), m1.Value(), m1.TsNs())
	}

	// Server observed every request field, including the nested opaque buffers.
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.gotFilerGroup != "fg-east" {
		t.Fatalf("server FilerGroup = %q", b.gotFilerGroup)
	}
	if !bytes.Equal(b.gotFetchTopic, topic) || !bytes.Equal(b.gotFetchPartition, partition) {
		t.Fatalf("server topic/partition mismatch: topic=%q partition=%q",
			b.gotFetchTopic, b.gotFetchPartition)
	}
	if b.gotFetchOffset != 10 || b.gotFetchMaxMsgs != 16 || b.gotFetchGroup != "cg-readers" {
		t.Fatalf("server fetch fields: off=%d maxmsgs=%d group=%q",
			b.gotFetchOffset, b.gotFetchMaxMsgs, b.gotFetchGroup)
	}
}

// TestPublishMessageStreamRoundTrip proves the hand-wired streaming path:
// PublishMessage is a bidi stream over transport.OpenStream / ListenStream. The
// client opens the stream with an InitMessage frame, streams two DataMessage
// frames, and half-closes; the server acks the init, records each value, and
// acks each data frame. Every frame is a zero-copy mq_brokerwire buffer — the
// same doctrine as the unary path, framed as stream messages. No HTTP, no
// protobuf.
func TestPublishMessageStreamRoundTrip(t *testing.T) {
	b := &fakeBroker{}

	srv, err := Serve("tcp", "127.0.0.1:0", b, b)
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	// Open the PublishMessage stream with the InitMessage variant.
	initVar := NewPublishMessageRequestInitMessage(PublishMessageRequestInitMessageInput{
		PublisherName: "pub-A",
		AckInterval:   1,
	})
	init := NewPublishMessageRequest(PublishMessageRequestInput{
		MessageWhich: PublishMessageRequestMessageInit,
		MessageValue: initVar,
	})
	stream, err := cli.OpenPublishMessage(init)
	if err != nil {
		t.Fatalf("OpenPublishMessage: %v", err)
	}

	// First inbound frame is the init ack carrying the base offset.
	ackBody, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv init ack: %v", err)
	}
	initAck, err := WrapPublishMessageResponse(ackBody)
	if err != nil {
		t.Fatalf("Wrap init ack: %v", err)
	}
	if initAck.AssignedOffset() != 100 {
		t.Fatalf("init ack AssignedOffset = %d, want 100", initAck.AssignedOffset())
	}

	// Stream two data frames; expect one ack each, with incrementing offsets.
	wantValues := [][]byte{[]byte("payload-0"), []byte("payload-1")}
	for i, v := range wantValues {
		dataVar := NewDataMessage(DataMessageInput{Value: v, TsNs: int64(5000 + i)})
		dataReq := NewPublishMessageRequest(PublishMessageRequestInput{
			MessageWhich: PublishMessageRequestMessageData,
			MessageValue: dataVar,
		})
		if err := stream.Send(dataReq); err != nil {
			t.Fatalf("Send data %d: %v", i, err)
		}
		respBody, err := stream.Recv()
		if err != nil {
			t.Fatalf("Recv data ack %d: %v", i, err)
		}
		resp, err := WrapPublishMessageResponse(respBody)
		if err != nil {
			t.Fatalf("Wrap data ack %d: %v", i, err)
		}
		if resp.AckTsNs() != int64(5000+i) {
			t.Fatalf("data ack %d AckTsNs = %d, want %d", i, resp.AckTsNs(), 5000+i)
		}
		if resp.AssignedOffset() != int64(101+i) {
			t.Fatalf("data ack %d AssignedOffset = %d, want %d", i, resp.AssignedOffset(), 101+i)
		}
	}

	// Half-close the send side; the server's Recv loop ends.
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("CloseSend: %v", err)
	}

	// The server observed both values in order.
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.gotPublishValues) != 2 {
		t.Fatalf("server saw %d values, want 2", len(b.gotPublishValues))
	}
	for i, v := range wantValues {
		if !bytes.Equal(b.gotPublishValues[i], v) {
			t.Fatalf("server value %d = %q, want %q", i, b.gotPublishValues[i], v)
		}
	}
}
