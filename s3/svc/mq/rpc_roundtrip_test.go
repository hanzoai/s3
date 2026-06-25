// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package mq

import (
	"context"
	"io"
	"testing"

	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	mq_pb "github.com/hanzoai/s3/s3/pb/mq_pb"
	schema_pb "github.com/hanzoai/s3/s3/pb/schema_pb"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
	"github.com/zap-proto/go/transport"
)

// fake is an in-memory Broker + StreamBroker. The unary side reconstructs each
// request from the zero-copy view and returns fixed values so the test can assert
// every field — scalars, nested Topic/Partition, the recursive RecordType, the
// repeated BrokerPartitionAssignment and DataMessage lists — survived the wire
// crossing. The streaming side serves PublishMessage (bidi echo) and
// GetUnflushedMessages (server-stream of LogEntry frames).
type fake struct {
	gotFilerGroup  string
	gotCfgTopic    *schema_pb.Topic
	gotCfgFields   int
	gotFetchTopic  *schema_pb.Topic
	gotFetchPart   *schema_pb.Partition
	gotPublishVals [][]byte
}

func (f *fake) FindBrokerLeader(req mq_brokerwire.FindBrokerLeaderRequest) (mq_brokerwire.FindBrokerLeaderResponseInput, error) {
	f.gotFilerGroup = req.FilerGroup()
	return mq_brokerwire.FindBrokerLeaderResponseInput{Broker: "broker-1.hanzo.svc:17777"}, nil
}

func (f *fake) ConfigureTopic(req mq_brokerwire.ConfigureTopicRequest) (mq_brokerwire.ConfigureTopicResponseInput, error) {
	f.gotCfgTopic = topicFromWire(req.Topic())
	if rt := recordTypeFromWire(req.MessageRecordType()); rt != nil {
		f.gotCfgFields = len(rt.Fields)
	}
	// Echo back one assignment + retention + the record type so the decoder side
	// is exercised too.
	return mq_brokerwire.ConfigureTopicResponseInput{
		BrokerPartitionAssignments: [][]byte{assignmentToWire(&mq_pb.BrokerPartitionAssignment{
			Partition: &schema_pb.Partition{RingSize: 8, RangeStart: 0, RangeStop: 7}, LeaderBroker: "ldr", FollowerBroker: "fol",
		})},
		Retention:         retentionToWire(&mq_pb.TopicRetention{RetentionSeconds: 3600, Enabled: true}),
		MessageRecordType: req.MessageRecordType(),
		KeyColumns:        []string{"id"},
		SchemaFormat:      "avro",
	}, nil
}

func (f *fake) FetchMessage(req mq_brokerwire.FetchMessageRequest) (mq_brokerwire.FetchMessageResponseInput, error) {
	f.gotFetchTopic = topicFromWire(req.Topic())
	f.gotFetchPart = partitionFromWire(req.Partition())
	m0 := dataMessageToWire(&mq_pb.DataMessage{Key: []byte("k0"), Value: []byte("v0"), TsNs: 1000})
	m1 := dataMessageToWire(&mq_pb.DataMessage{Key: []byte("k1"), Value: []byte("v1"), TsNs: 2000})
	return mq_brokerwire.FetchMessageResponseInput{
		Messages: [][]byte{m0, m1}, HighWaterMark: 42, NextOffset: req.StartOffset() + 2, EndOfPartition: true,
	}, nil
}

// Remaining unary methods: empty successful responses to complete the Broker.
func (f *fake) BalanceTopics(mq_brokerwire.BalanceTopicsRequest) (mq_brokerwire.BalanceTopicsResponseInput, error) {
	return mq_brokerwire.BalanceTopicsResponseInput{}, nil
}
func (f *fake) ListTopics(mq_brokerwire.ListTopicsRequest) (mq_brokerwire.ListTopicsResponseInput, error) {
	return mq_brokerwire.ListTopicsResponseInput{Topics: [][]byte{topicToWire(&schema_pb.Topic{Namespace: "ns", Name: "t"})}}, nil
}
func (f *fake) TopicExists(mq_brokerwire.TopicExistsRequest) (mq_brokerwire.TopicExistsResponseInput, error) {
	return mq_brokerwire.TopicExistsResponseInput{Exists: true}, nil
}
func (f *fake) LookupTopicBrokers(mq_brokerwire.LookupTopicBrokersRequest) (mq_brokerwire.LookupTopicBrokersResponseInput, error) {
	return mq_brokerwire.LookupTopicBrokersResponseInput{}, nil
}
func (f *fake) GetTopicConfiguration(mq_brokerwire.GetTopicConfigurationRequest) (mq_brokerwire.GetTopicConfigurationResponseInput, error) {
	return mq_brokerwire.GetTopicConfigurationResponseInput{}, nil
}
func (f *fake) GetTopicPublishers(mq_brokerwire.GetTopicPublishersRequest) (mq_brokerwire.GetTopicPublishersResponseInput, error) {
	return mq_brokerwire.GetTopicPublishersResponseInput{}, nil
}
func (f *fake) GetTopicSubscribers(mq_brokerwire.GetTopicSubscribersRequest) (mq_brokerwire.GetTopicSubscribersResponseInput, error) {
	return mq_brokerwire.GetTopicSubscribersResponseInput{}, nil
}
func (f *fake) AssignTopicPartitions(mq_brokerwire.AssignTopicPartitionsRequest) (mq_brokerwire.AssignTopicPartitionsResponseInput, error) {
	return mq_brokerwire.AssignTopicPartitionsResponseInput{}, nil
}
func (f *fake) ClosePublishers(mq_brokerwire.ClosePublishersRequest) (mq_brokerwire.ClosePublishersResponseInput, error) {
	return mq_brokerwire.ClosePublishersResponseInput{}, nil
}
func (f *fake) CloseSubscribers(mq_brokerwire.CloseSubscribersRequest) (mq_brokerwire.CloseSubscribersResponseInput, error) {
	return mq_brokerwire.CloseSubscribersResponseInput{}, nil
}
func (f *fake) GetPartitionRangeInfo(mq_brokerwire.GetPartitionRangeInfoRequest) (mq_brokerwire.GetPartitionRangeInfoResponseInput, error) {
	return mq_brokerwire.GetPartitionRangeInfoResponseInput{}, nil
}

// ServeStream serves the two streaming RPCs the test drives.
func (f *fake) ServeStream(method uint32, init []byte, s transport.Stream) {
	switch method {
	case mq_brokerwire.HanzoMessagingPublishMessageOrdinal:
		f.servePublish(s)
	case mq_brokerwire.HanzoMessagingGetUnflushedMessagesOrdinal:
		f.serveUnflushed(s)
	}
}

// servePublish reads every frame (init + data) off the stream and, for each data
// frame, records the value and acks with an incrementing offset.
func (f *fake) servePublish(s transport.Stream) {
	var n int64
	for {
		frame, err := s.Recv()
		if err != nil {
			return
		}
		req, err := mq_brokerwire.WrapPublishMessageRequest(frame)
		if err != nil {
			return
		}
		if req.WhichMessage() != mq_brokerwire.PublishMessageRequestMessageData {
			continue
		}
		f.gotPublishVals = append(f.gotPublishVals, append([]byte(nil), req.Data().Value()...))
		n++
		if err := s.Send(mq_brokerwire.NewPublishMessageResponse(mq_brokerwire.PublishMessageResponseInput{
			AckTsNs: req.Data().TsNs(), AssignedOffset: 100 + n,
		})); err != nil {
			return
		}
	}
}

// serveUnflushed streams two LogEntry frames then returns (the client sees EOF).
func (f *fake) serveUnflushed(s transport.Stream) {
	for i := 0; i < 2; i++ {
		if err := s.Send(mq_brokerwire.NewGetUnflushedMessagesResponse(mq_brokerwire.GetUnflushedMessagesResponseInput{
			Message: logEntryToWire(&filer_pb.LogEntry{TsNs: int64(i + 1), Offset: int64(i), Data: []byte("d")}),
		})); err != nil {
			return
		}
	}
}

func serve(t *testing.T) (*transport.Server, *fake) {
	t.Helper()
	f := &fake{}
	srv, err := mq_brokerwire.Serve("tcp", "127.0.0.1:0", f, f)
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	return srv, f
}

// TestUnaryRoundTrip proves the unary seam end to end over a real TCP socket:
// requests cross as ZAP envelopes carrying zero-copy mq_brokerwire payloads, the
// server dispatches to the backend, and responses come back the same way — no
// gRPC, no protobuf. It exercises a scalar RPC and a deeply-nested one (Topic +
// Retention + recursive RecordType + repeated assignments).
func TestUnaryRoundTrip(t *testing.T) {
	srv, f := serve(t)
	defer srv.Close()
	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()
	cli := New(conn, nil)
	ctx := context.Background()

	leader, err := cli.FindBrokerLeader(ctx, &mq_pb.FindBrokerLeaderRequest{FilerGroup: "fg-east"})
	if err != nil {
		t.Fatalf("FindBrokerLeader: %v", err)
	}
	if leader.Broker != "broker-1.hanzo.svc:17777" {
		t.Fatalf("FindBrokerLeader broker = %q", leader.Broker)
	}
	if f.gotFilerGroup != "fg-east" {
		t.Fatalf("server FilerGroup = %q", f.gotFilerGroup)
	}

	// ConfigureTopic: nested Topic, recursive RecordType (a record field whose type
	// is itself a record), Retention, and a repeated assignment back.
	rt := &schema_pb.RecordType{Fields: []*schema_pb.Field{
		{Name: "id", FieldIndex: 0, IsRequired: true, Type: &schema_pb.Type{Kind: &schema_pb.Type_ScalarType{ScalarType: schema_pb.ScalarType_INT64}}},
		{Name: "nested", FieldIndex: 1, Type: &schema_pb.Type{Kind: &schema_pb.Type_RecordType{RecordType: &schema_pb.RecordType{
			Fields: []*schema_pb.Field{{Name: "inner", Type: &schema_pb.Type{Kind: &schema_pb.Type_ScalarType{ScalarType: schema_pb.ScalarType_STRING}}}},
		}}}},
	}}
	cfg, err := cli.ConfigureTopic(ctx, &mq_pb.ConfigureTopicRequest{
		Topic:             &schema_pb.Topic{Namespace: "ns", Name: "orders"},
		PartitionCount:    4,
		Retention:         &mq_pb.TopicRetention{RetentionSeconds: 7200, Enabled: true},
		MessageRecordType: rt,
		KeyColumns:        []string{"id"},
		SchemaFormat:      "avro",
	})
	if err != nil {
		t.Fatalf("ConfigureTopic: %v", err)
	}
	if f.gotCfgTopic == nil || f.gotCfgTopic.Namespace != "ns" || f.gotCfgTopic.Name != "orders" {
		t.Fatalf("server Topic = %+v", f.gotCfgTopic)
	}
	if f.gotCfgFields != 2 {
		t.Fatalf("server RecordType fields = %d, want 2", f.gotCfgFields)
	}
	if len(cfg.BrokerPartitionAssignments) != 1 || cfg.BrokerPartitionAssignments[0].LeaderBroker != "ldr" ||
		cfg.BrokerPartitionAssignments[0].Partition.RingSize != 8 {
		t.Fatalf("assignment round-trip lost: %+v", cfg.BrokerPartitionAssignments)
	}
	if cfg.Retention == nil || cfg.Retention.RetentionSeconds != 3600 || !cfg.Retention.Enabled {
		t.Fatalf("retention round-trip lost: %+v", cfg.Retention)
	}
	if cfg.MessageRecordType == nil || len(cfg.MessageRecordType.Fields) != 2 {
		t.Fatalf("record type round-trip lost: %+v", cfg.MessageRecordType)
	}
	// The recursive field's nested record survived.
	nested := cfg.MessageRecordType.Fields[1].GetType().GetRecordType()
	if nested == nil || len(nested.Fields) != 1 || nested.Fields[0].Name != "inner" {
		t.Fatalf("recursive record type lost: %+v", nested)
	}

	// FetchMessage: nested Topic/Partition in, repeated DataMessage list out.
	fetch, err := cli.FetchMessage(ctx, &mq_pb.FetchMessageRequest{
		Topic:       &schema_pb.Topic{Namespace: "ns", Name: "orders"},
		Partition:   &schema_pb.Partition{RingSize: 16, RangeStart: 0, RangeStop: 15, UnixTimeNs: 99},
		StartOffset: 10,
	})
	if err != nil {
		t.Fatalf("FetchMessage: %v", err)
	}
	if f.gotFetchTopic == nil || f.gotFetchTopic.Name != "orders" {
		t.Fatalf("server fetch topic = %+v", f.gotFetchTopic)
	}
	if f.gotFetchPart == nil || f.gotFetchPart.RingSize != 16 || f.gotFetchPart.UnixTimeNs != 99 {
		t.Fatalf("server fetch partition = %+v", f.gotFetchPart)
	}
	if fetch.HighWaterMark != 42 || fetch.NextOffset != 12 || !fetch.EndOfPartition {
		t.Fatalf("fetch scalars: hwm=%d next=%d eop=%v", fetch.HighWaterMark, fetch.NextOffset, fetch.EndOfPartition)
	}
	if len(fetch.Messages) != 2 || string(fetch.Messages[0].Value) != "v0" || fetch.Messages[1].TsNs != 2000 {
		t.Fatalf("fetch messages round-trip lost: %+v", fetch.Messages)
	}
}

// TestPublishMessageStreamRoundTrip proves the bidi streaming seam: the client
// opens PublishMessage, Sends an init frame and two data frames, and Recvs an ack
// per data frame with incrementing offsets. Every frame is a zero-copy
// mq_brokerwire buffer over a transport.Stream — same doctrine as unary.
func TestPublishMessageStreamRoundTrip(t *testing.T) {
	srv, f := serve(t)
	defer srv.Close()
	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()
	cli := New(conn, nil)

	stream, err := cli.PublishMessage(context.Background())
	if err != nil {
		t.Fatalf("PublishMessage: %v", err)
	}
	if err := stream.Send(&mq_pb.PublishMessageRequest{Message: &mq_pb.PublishMessageRequest_Init{
		Init: &mq_pb.PublishMessageRequest_InitMessage{
			Topic: &schema_pb.Topic{Namespace: "ns", Name: "orders"}, PublisherName: "pub-A", AckInterval: 1,
		},
	}}); err != nil {
		t.Fatalf("Send init: %v", err)
	}

	wantValues := [][]byte{[]byte("payload-0"), []byte("payload-1")}
	for i, v := range wantValues {
		if err := stream.Send(&mq_pb.PublishMessageRequest{Message: &mq_pb.PublishMessageRequest_Data{
			Data: &mq_pb.DataMessage{Value: v, TsNs: int64(5000 + i)},
		}}); err != nil {
			t.Fatalf("Send data %d: %v", i, err)
		}
		resp, err := stream.Recv()
		if err != nil {
			t.Fatalf("Recv ack %d: %v", i, err)
		}
		if resp.AckTsNs != int64(5000+i) || resp.AssignedOffset != int64(101+i) {
			t.Fatalf("ack %d = ts=%d off=%d", i, resp.AckTsNs, resp.AssignedOffset)
		}
	}
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("CloseSend: %v", err)
	}
	if len(f.gotPublishVals) != 2 || string(f.gotPublishVals[0]) != "payload-0" {
		t.Fatalf("server values = %v", f.gotPublishVals)
	}
}

// TestGetUnflushedMessagesStreamRoundTrip proves the server-streaming seam: the
// request rides the opener and the client Recvs each LogEntry frame until EOF.
func TestGetUnflushedMessagesStreamRoundTrip(t *testing.T) {
	srv, _ := serve(t)
	defer srv.Close()
	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()
	cli := New(conn, nil)

	stream, err := cli.GetUnflushedMessages(context.Background(), &mq_pb.GetUnflushedMessagesRequest{
		Topic: &schema_pb.Topic{Namespace: "ns", Name: "orders"}, StartBufferOffset: 5,
	})
	if err != nil {
		t.Fatalf("GetUnflushedMessages: %v", err)
	}
	var got []int64
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Recv: %v", err)
		}
		if resp.Message == nil {
			t.Fatalf("nil LogEntry: %+v", resp)
		}
		got = append(got, resp.Message.TsNs)
	}
	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("streamed LogEntry ts = %v, want [1 2]", got)
	}
}
