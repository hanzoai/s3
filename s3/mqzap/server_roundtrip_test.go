// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package mqzap

import (
	"bytes"
	"context"
	"io"
	"testing"

	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	mq_pb "github.com/hanzoai/s3/s3/pb/mq_pb"
	schema_pb "github.com/hanzoai/s3/s3/pb/schema_pb"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
	"github.com/zap-proto/go/transport"
)

// fakeEngine implements the engine-side contract mq_pb.HanzoMessagingServer (the
// gRPC-shaped interface the real *broker.MessageQueueBroker satisfies). The
// server-adapter round-trip test wraps it with NewServerBackend + NewStreamServer,
// serves it over the canonical ZAP transport via mq_brokerwire.Serve, and dials it
// with the ZAP-backed client (New). This exercises the WHOLE server path the
// command/mq_broker.go cutover relies on: request view -> *mq_pb -> engine ->
// *mq_pb response -> wire, with no protobuf framing or gRPC anywhere on the wire.
type fakeEngine struct {
	mq_pb.UnimplementedHanzoMessagingServer

	gotFilerGroup  string
	gotCfgTopic    *schema_pb.Topic
	gotCfgParts    int32
	gotFetchTopic  *schema_pb.Topic
	gotFetchPart   *schema_pb.Partition
	gotPublishVals [][]byte
	gotSubInitCG   string
}

func (e *fakeEngine) FindBrokerLeader(_ context.Context, req *mq_pb.FindBrokerLeaderRequest) (*mq_pb.FindBrokerLeaderResponse, error) {
	e.gotFilerGroup = req.FilerGroup
	return &mq_pb.FindBrokerLeaderResponse{Broker: "broker-7.hanzo.svc:17777"}, nil
}

func (e *fakeEngine) ConfigureTopic(_ context.Context, req *mq_pb.ConfigureTopicRequest) (*mq_pb.ConfigureTopicResponse, error) {
	e.gotCfgTopic = req.Topic
	e.gotCfgParts = req.PartitionCount
	return &mq_pb.ConfigureTopicResponse{
		BrokerPartitionAssignments: []*mq_pb.BrokerPartitionAssignment{
			{LeaderBroker: "leader-A", FollowerBroker: "follower-B"},
		},
		MessageRecordType: req.MessageRecordType,
		KeyColumns:        req.KeyColumns,
		SchemaFormat:      req.SchemaFormat,
		Retention:         &mq_pb.TopicRetention{RetentionSeconds: 3600, Enabled: true},
	}, nil
}

func (e *fakeEngine) FetchMessage(_ context.Context, req *mq_pb.FetchMessageRequest) (*mq_pb.FetchMessageResponse, error) {
	e.gotFetchTopic = req.Topic
	e.gotFetchPart = req.Partition
	return &mq_pb.FetchMessageResponse{
		Messages: []*mq_pb.DataMessage{
			{Key: []byte("k0"), Value: []byte("v0"), TsNs: 100},
			{Key: []byte("k1"), Value: []byte("v1"), TsNs: 200},
		},
		HighWaterMark:  9,
		LogStartOffset: 1,
		EndOfPartition: true,
		NextOffset:     req.StartOffset + 2,
	}, nil
}

func (e *fakeEngine) GetPartitionRangeInfo(_ context.Context, req *mq_pb.GetPartitionRangeInfoRequest) (*mq_pb.GetPartitionRangeInfoResponse, error) {
	return &mq_pb.GetPartitionRangeInfoResponse{
		OffsetRange:         &mq_pb.OffsetRangeInfo{EarliestOffset: 3, LatestOffset: 42, HighWaterMark: 43},
		TimestampRange:      &mq_pb.TimestampRangeInfo{EarliestTimestampNs: 1000, LatestTimestampNs: 5000},
		RecordCount:         40,
		ActiveSubscriptions: 2,
	}, nil
}

// PublishMessage (bidi): ack the init with a base offset, then echo each data
// frame's value (recorded) with an incrementing assigned offset.
func (e *fakeEngine) PublishMessage(stream mq_pb.HanzoMessaging_PublishMessageServer) error {
	first, err := stream.Recv()
	if err != nil {
		return err
	}
	if first.GetInit() == nil {
		return io.ErrUnexpectedEOF
	}
	if err := stream.Send(&mq_pb.PublishMessageResponse{AssignedOffset: 500}); err != nil {
		return err
	}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if d := req.GetData(); d != nil {
			e.gotPublishVals = append(e.gotPublishVals, append([]byte(nil), d.Value...))
			if err := stream.Send(&mq_pb.PublishMessageResponse{
				AckTsNs:        d.TsNs,
				AssignedOffset: int64(500 + len(e.gotPublishVals)),
			}); err != nil {
				return err
			}
		}
	}
}

// SubscribeMessage (bidi): record the init consumer group, then send one ctrl
// (end-of-stream) frame and return — proving the oneof response path.
func (e *fakeEngine) SubscribeMessage(stream mq_pb.HanzoMessaging_SubscribeMessageServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	if init := req.GetInit(); init != nil {
		e.gotSubInitCG = init.ConsumerGroup
	}
	return stream.Send(&mq_pb.SubscribeMessageResponse{
		Message: &mq_pb.SubscribeMessageResponse_Data{
			Data: &mq_pb.DataMessage{Key: []byte("dk"), Value: []byte("dv"), TsNs: 777},
		},
	})
}

// GetUnflushedMessages (server-stream): emit two LogEntry frames then end-of-stream.
func (e *fakeEngine) GetUnflushedMessages(req *mq_pb.GetUnflushedMessagesRequest, stream mq_pb.HanzoMessaging_GetUnflushedMessagesServer) error {
	for i := 0; i < 2; i++ {
		if err := stream.Send(&mq_pb.GetUnflushedMessagesResponse{
			Message: &filer_pb.LogEntry{TsNs: int64(1000 + i), Offset: int64(i), Data: []byte{byte(i)}},
		}); err != nil {
			return err
		}
	}
	return stream.Send(&mq_pb.GetUnflushedMessagesResponse{EndOfStream: true})
}

// serveFakeEngine wires the fake engine through the real server adapters and the
// canonical ZAP transport, returning a connected ZAP-backed client.
func serveFakeEngine(t *testing.T, e *fakeEngine) (mq_pb.HanzoMessagingClient, func()) {
	t.Helper()
	srv, err := mq_brokerwire.Serve("tcp", "127.0.0.1:0", NewServerBackend(e), NewStreamServer(e))
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		srv.Close()
		t.Fatalf("Dial: %v", err)
	}
	cli := New(conn, nil)
	return cli, func() { conn.Close(); srv.Close() }
}

// TestServerAdapterUnaryRoundTrip proves the server adapter for the unary RPCs:
// every request field reaches the engine and every response field returns,
// including nested Topic/Partition, the recursive Retention, the repeated
// assignment list, and the Offset/Timestamp range sub-messages.
func TestServerAdapterUnaryRoundTrip(t *testing.T) {
	e := &fakeEngine{}
	cli, stop := serveFakeEngine(t, e)
	defer stop()
	ctx := context.Background()

	leader, err := cli.FindBrokerLeader(ctx, &mq_pb.FindBrokerLeaderRequest{FilerGroup: "fg-west"})
	if err != nil {
		t.Fatalf("FindBrokerLeader: %v", err)
	}
	if leader.Broker != "broker-7.hanzo.svc:17777" {
		t.Fatalf("leader = %q", leader.Broker)
	}
	if e.gotFilerGroup != "fg-west" {
		t.Fatalf("engine FilerGroup = %q", e.gotFilerGroup)
	}

	cfg, err := cli.ConfigureTopic(ctx, &mq_pb.ConfigureTopicRequest{
		Topic:          &schema_pb.Topic{Namespace: "ns", Name: "events"},
		PartitionCount: 6,
		KeyColumns:     []string{"id"},
		SchemaFormat:   "avro",
	})
	if err != nil {
		t.Fatalf("ConfigureTopic: %v", err)
	}
	if e.gotCfgTopic.GetName() != "events" || e.gotCfgParts != 6 {
		t.Fatalf("engine cfg topic=%v parts=%d", e.gotCfgTopic, e.gotCfgParts)
	}
	if len(cfg.BrokerPartitionAssignments) != 1 ||
		cfg.BrokerPartitionAssignments[0].LeaderBroker != "leader-A" ||
		cfg.BrokerPartitionAssignments[0].FollowerBroker != "follower-B" {
		t.Fatalf("cfg assignments = %+v", cfg.BrokerPartitionAssignments)
	}
	if cfg.Retention == nil || cfg.Retention.RetentionSeconds != 3600 || !cfg.Retention.Enabled {
		t.Fatalf("cfg retention = %+v", cfg.Retention)
	}
	if len(cfg.KeyColumns) != 1 || cfg.KeyColumns[0] != "id" || cfg.SchemaFormat != "avro" {
		t.Fatalf("cfg keycols=%v fmt=%q", cfg.KeyColumns, cfg.SchemaFormat)
	}

	topic := &schema_pb.Topic{Namespace: "ns", Name: "events"}
	part := &schema_pb.Partition{RingSize: 1024, RangeStart: 0, RangeStop: 512, UnixTimeNs: 42}
	fetch, err := cli.FetchMessage(ctx, &mq_pb.FetchMessageRequest{Topic: topic, Partition: part, StartOffset: 10})
	if err != nil {
		t.Fatalf("FetchMessage: %v", err)
	}
	if e.gotFetchTopic.GetName() != "events" || e.gotFetchPart.GetRingSize() != 1024 || e.gotFetchPart.GetRangeStop() != 512 {
		t.Fatalf("engine fetch topic=%v part=%v", e.gotFetchTopic, e.gotFetchPart)
	}
	if fetch.HighWaterMark != 9 || fetch.NextOffset != 12 || !fetch.EndOfPartition || len(fetch.Messages) != 2 {
		t.Fatalf("fetch resp = %+v", fetch)
	}
	if !bytes.Equal(fetch.Messages[0].Value, []byte("v0")) || fetch.Messages[1].TsNs != 200 {
		t.Fatalf("fetch msgs = %+v", fetch.Messages)
	}

	rng, err := cli.GetPartitionRangeInfo(ctx, &mq_pb.GetPartitionRangeInfoRequest{Topic: topic, Partition: part})
	if err != nil {
		t.Fatalf("GetPartitionRangeInfo: %v", err)
	}
	if rng.OffsetRange.LatestOffset != 42 || rng.TimestampRange.LatestTimestampNs != 5000 ||
		rng.RecordCount != 40 || rng.ActiveSubscriptions != 2 {
		t.Fatalf("range = %+v", rng)
	}
}

// TestServerAdapterPublishStream proves the bidi server-stream path through the
// adapter: init ack, two data frames recorded by the engine, per-frame acks.
func TestServerAdapterPublishStream(t *testing.T) {
	e := &fakeEngine{}
	cli, stop := serveFakeEngine(t, e)
	defer stop()
	ctx := context.Background()

	stream, err := cli.PublishMessage(ctx)
	if err != nil {
		t.Fatalf("PublishMessage: %v", err)
	}
	if err := stream.Send(&mq_pb.PublishMessageRequest{
		Message: &mq_pb.PublishMessageRequest_Init{
			Init: &mq_pb.PublishMessageRequest_InitMessage{PublisherName: "pub-X"},
		},
	}); err != nil {
		t.Fatalf("send init: %v", err)
	}
	initAck, err := stream.Recv()
	if err != nil {
		t.Fatalf("recv init ack: %v", err)
	}
	if initAck.AssignedOffset != 500 {
		t.Fatalf("init ack offset = %d", initAck.AssignedOffset)
	}

	want := [][]byte{[]byte("p0"), []byte("p1")}
	for i, v := range want {
		if err := stream.Send(&mq_pb.PublishMessageRequest{
			Message: &mq_pb.PublishMessageRequest_Data{Data: &mq_pb.DataMessage{Value: v, TsNs: int64(900 + i)}},
		}); err != nil {
			t.Fatalf("send data %d: %v", i, err)
		}
		ack, err := stream.Recv()
		if err != nil {
			t.Fatalf("recv ack %d: %v", i, err)
		}
		if ack.AckTsNs != int64(900+i) || ack.AssignedOffset != int64(501+i) {
			t.Fatalf("ack %d = ts=%d off=%d", i, ack.AckTsNs, ack.AssignedOffset)
		}
	}
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("CloseSend: %v", err)
	}
	if len(e.gotPublishVals) != 2 || !bytes.Equal(e.gotPublishVals[0], []byte("p0")) || !bytes.Equal(e.gotPublishVals[1], []byte("p1")) {
		t.Fatalf("engine saw %q", e.gotPublishVals)
	}
}

// TestServerAdapterSubscribeStream proves a bidi RPC whose response is a oneof
// (SubscribeMessage Data variant) round-trips through the adapter.
func TestServerAdapterSubscribeStream(t *testing.T) {
	e := &fakeEngine{}
	cli, stop := serveFakeEngine(t, e)
	defer stop()
	ctx := context.Background()

	stream, err := cli.SubscribeMessage(ctx)
	if err != nil {
		t.Fatalf("SubscribeMessage: %v", err)
	}
	if err := stream.Send(&mq_pb.SubscribeMessageRequest{
		Message: &mq_pb.SubscribeMessageRequest_Init{
			Init: &mq_pb.SubscribeMessageRequest_InitMessage{ConsumerGroup: "cg-1", ConsumerId: "c-1"},
		},
	}); err != nil {
		t.Fatalf("send init: %v", err)
	}
	resp, err := stream.Recv()
	if err != nil {
		t.Fatalf("recv: %v", err)
	}
	d := resp.GetData()
	if d == nil || !bytes.Equal(d.Value, []byte("dv")) || d.TsNs != 777 {
		t.Fatalf("sub data = %+v", resp.Message)
	}
	if e.gotSubInitCG != "cg-1" {
		t.Fatalf("engine sub cg = %q", e.gotSubInitCG)
	}
	_ = stream.CloseSend()
}

// TestServerAdapterGetUnflushedStream proves the server-stream path: the request
// rides the opener and the engine streams two LogEntry frames then end-of-stream.
func TestServerAdapterGetUnflushedStream(t *testing.T) {
	e := &fakeEngine{}
	cli, stop := serveFakeEngine(t, e)
	defer stop()
	ctx := context.Background()

	stream, err := cli.GetUnflushedMessages(ctx, &mq_pb.GetUnflushedMessagesRequest{
		Topic:             &schema_pb.Topic{Namespace: "ns", Name: "t"},
		Partition:         &schema_pb.Partition{RingSize: 8},
		StartBufferOffset: 3,
	})
	if err != nil {
		t.Fatalf("GetUnflushedMessages: %v", err)
	}
	var got int
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("recv: %v", err)
		}
		if resp.EndOfStream {
			break
		}
		if resp.Message == nil {
			t.Fatalf("nil message frame %d", got)
		}
		if resp.Message.Offset != int64(got) {
			t.Fatalf("frame %d offset = %d", got, resp.Message.Offset)
		}
		got++
	}
	if got != 2 {
		t.Fatalf("received %d log entries, want 2", got)
	}
}
