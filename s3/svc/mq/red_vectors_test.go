// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// red_vectors_test.go: adversarial round-trip coverage for the MQ broker ZAP rip
// vectors Blue self-flagged MEDIUM/LOW and left untested:
//   - PublisherToPubBalancer BrokerStats/TopicPartitionStats map + nested oneof
//   - SubscriberToSubCoordinator Assignment AND UnAssignment oneof response variants
//   - SubscribeFollowMe client-stream SendAndClose half-close
//   - PublishFollowMe bidi (Init/Data/Flush/Close variants)
//   - empty/None oneof + nil-Partition / nil-Topic paths
//
// These wrap the SAME server adapters (NewServerBackend + NewStreamServer) over the
// real ZAP transport as server_roundtrip_test.go's fakeEngine.

package mq

import (
	"bytes"
	"context"
	"io"
	"testing"

	mq_pb "github.com/hanzoai/s3/s3/pb/mq_pb"
	schema_pb "github.com/hanzoai/s3/s3/pb/schema_pb"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
	"github.com/zap-proto/go/transport"
)

// redEngine implements the four streaming RPCs the base fakeEngine leaves on
// Unimplemented, plus records what it saw so the test can assert the request
// direction round-tripped too.
type redEngine struct {
	mq_pb.UnimplementedHanzoMessagingServer

	gotBalancerInitBroker string
	gotBalancerStats      *mq_pb.BrokerStats

	gotSubCoordInitCG string
	gotAckAssignPart  *schema_pb.Partition
	subCoordResponses []*mq_pb.SubscriberToSubCoordinatorResponse

	gotFollowMeInitCG string
	gotFollowMeAcks   []int64
	gotFollowMeClose  bool

	gotPubFollowInit  *schema_pb.Topic
	gotPubFollowData  [][]byte
	gotPubFollowFlush []int64
	gotPubFollowClose bool
}

// PublisherToPubBalancer: read Init then a Stats frame, record both, reply empty.
func (e *redEngine) PublisherToPubBalancer(stream mq_pb.HanzoMessaging_PublisherToPubBalancerServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		switch m := req.Message.(type) {
		case *mq_pb.PublisherToPubBalancerRequest_Init:
			e.gotBalancerInitBroker = m.Init.Broker
		case *mq_pb.PublisherToPubBalancerRequest_Stats:
			e.gotBalancerStats = m.Stats
			if err := stream.Send(&mq_pb.PublisherToPubBalancerResponse{}); err != nil {
				return err
			}
		}
	}
}

// SubscriberToSubCoordinator: record Init, then emit an Assignment response AND an
// UnAssignment response so the test asserts both oneof variants survive the wire.
func (e *redEngine) SubscriberToSubCoordinator(stream mq_pb.HanzoMessaging_SubscriberToSubCoordinatorServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		switch m := req.Message.(type) {
		case *mq_pb.SubscriberToSubCoordinatorRequest_Init:
			e.gotSubCoordInitCG = m.Init.ConsumerGroup
			// Assignment variant: a non-nil PartitionAssignment with a Partition.
			if err := stream.Send(&mq_pb.SubscriberToSubCoordinatorResponse{
				Message: &mq_pb.SubscriberToSubCoordinatorResponse_Assignment_{
					Assignment: &mq_pb.SubscriberToSubCoordinatorResponse_Assignment{
						PartitionAssignment: &mq_pb.BrokerPartitionAssignment{
							Partition:    &schema_pb.Partition{RingSize: 2048, RangeStart: 7, RangeStop: 9, UnixTimeNs: 11},
							LeaderBroker: "lead-X", FollowerBroker: "foll-Y",
						},
					},
				},
			}); err != nil {
				return err
			}
			// UnAssignment variant: a Partition.
			if err := stream.Send(&mq_pb.SubscriberToSubCoordinatorResponse{
				Message: &mq_pb.SubscriberToSubCoordinatorResponse_UnAssignment_{
					UnAssignment: &mq_pb.SubscriberToSubCoordinatorResponse_UnAssignment{
						Partition: &schema_pb.Partition{RingSize: 4096, RangeStart: 1, RangeStop: 2, UnixTimeNs: 3},
					},
				},
			}); err != nil {
				return err
			}
		case *mq_pb.SubscriberToSubCoordinatorRequest_AckAssignment:
			e.gotAckAssignPart = m.AckAssignment.Partition
		}
	}
}

// SubscribeFollowMe (client-stream): Recv Init + Acks + Close until EOF, then
// SendAndClose the single summary.
func (e *redEngine) SubscribeFollowMe(stream mq_pb.HanzoMessaging_SubscribeFollowMeServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&mq_pb.SubscribeFollowMeResponse{AckTsNs: 999})
		}
		if err != nil {
			return err
		}
		switch m := req.Message.(type) {
		case *mq_pb.SubscribeFollowMeRequest_Init:
			e.gotFollowMeInitCG = m.Init.ConsumerGroup
		case *mq_pb.SubscribeFollowMeRequest_Ack:
			e.gotFollowMeAcks = append(e.gotFollowMeAcks, m.Ack.TsNs)
		case *mq_pb.SubscribeFollowMeRequest_Close:
			e.gotFollowMeClose = true
		}
	}
}

// PublishFollowMe (bidi): record every variant, ack each data frame.
func (e *redEngine) PublishFollowMe(stream mq_pb.HanzoMessaging_PublishFollowMeServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		switch m := req.Message.(type) {
		case *mq_pb.PublishFollowMeRequest_Init:
			e.gotPubFollowInit = m.Init.Topic
		case *mq_pb.PublishFollowMeRequest_Data:
			e.gotPubFollowData = append(e.gotPubFollowData, append([]byte(nil), m.Data.Value...))
			if err := stream.Send(&mq_pb.PublishFollowMeResponse{AckTsNs: m.Data.TsNs}); err != nil {
				return err
			}
		case *mq_pb.PublishFollowMeRequest_Flush:
			e.gotPubFollowFlush = append(e.gotPubFollowFlush, m.Flush.TsNs)
		case *mq_pb.PublishFollowMeRequest_Close:
			e.gotPubFollowClose = true
		}
	}
}

func serveRedEngine(t *testing.T, e *redEngine) (mq_pb.HanzoMessagingClient, func()) {
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
	return New(conn, nil), func() { conn.Close(); srv.Close() }
}

// Vector 1: PublisherToPubBalancer BrokerStats map + nested TopicPartitionStats.
func TestRedPublisherToPubBalancerStats(t *testing.T) {
	e := &redEngine{}
	cli, stop := serveRedEngine(t, e)
	defer stop()

	stream, err := cli.PublisherToPubBalancer(context.Background())
	if err != nil {
		t.Fatalf("PublisherToPubBalancer: %v", err)
	}
	if err := stream.Send(&mq_pb.PublisherToPubBalancerRequest{
		Message: &mq_pb.PublisherToPubBalancerRequest_Init{
			Init: &mq_pb.PublisherToPubBalancerRequest_InitMessage{Broker: "broker-init.hanzo.svc:17777"},
		},
	}); err != nil {
		t.Fatalf("send init: %v", err)
	}
	stats := &mq_pb.BrokerStats{
		CpuUsagePercent: 73,
		Stats: map[string]*mq_pb.TopicPartitionStats{
			"ns.topic-a/0": {
				Topic:           &schema_pb.Topic{Namespace: "ns", Name: "topic-a"},
				Partition:       &schema_pb.Partition{RingSize: 1024, RangeStart: 0, RangeStop: 512, UnixTimeNs: 5},
				PublisherCount:  3,
				SubscriberCount: 4,
				Follower:        "follower-z",
			},
			"ns.topic-b/1": {
				Topic:           &schema_pb.Topic{Namespace: "ns", Name: "topic-b"},
				Partition:       &schema_pb.Partition{RingSize: 2048, RangeStart: 512, RangeStop: 1024, UnixTimeNs: 6},
				PublisherCount:  1,
				SubscriberCount: 2,
				Follower:        "follower-w",
			},
		},
	}
	if err := stream.Send(&mq_pb.PublisherToPubBalancerRequest{
		Message: &mq_pb.PublisherToPubBalancerRequest_Stats{Stats: stats},
	}); err != nil {
		t.Fatalf("send stats: %v", err)
	}
	if _, err := stream.Recv(); err != nil {
		t.Fatalf("recv ack: %v", err)
	}
	_ = stream.CloseSend()

	if e.gotBalancerInitBroker != "broker-init.hanzo.svc:17777" {
		t.Fatalf("init broker = %q", e.gotBalancerInitBroker)
	}
	got := e.gotBalancerStats
	if got == nil {
		t.Fatal("engine never saw Stats frame (oneof collapsed to None?)")
	}
	if got.CpuUsagePercent != 73 {
		t.Fatalf("cpu = %v", got.CpuUsagePercent)
	}
	if len(got.Stats) != 2 {
		t.Fatalf("stats map len = %d, want 2: %+v", len(got.Stats), got.Stats)
	}
	a := got.Stats["ns.topic-a/0"]
	if a == nil || a.Topic.GetName() != "topic-a" || a.Partition.GetRingSize() != 1024 ||
		a.PublisherCount != 3 || a.SubscriberCount != 4 || a.Follower != "follower-z" {
		t.Fatalf("topic-a stats = %+v", a)
	}
	b := got.Stats["ns.topic-b/1"]
	if b == nil || b.Topic.GetName() != "topic-b" || b.Partition.GetRangeStop() != 1024 || b.Follower != "follower-w" {
		t.Fatalf("topic-b stats = %+v", b)
	}
}

// Vector 2: SubscriberToSubCoordinator Assignment AND UnAssignment response variants.
func TestRedSubscriberToSubCoordinatorOneof(t *testing.T) {
	e := &redEngine{}
	cli, stop := serveRedEngine(t, e)
	defer stop()

	stream, err := cli.SubscriberToSubCoordinator(context.Background())
	if err != nil {
		t.Fatalf("SubscriberToSubCoordinator: %v", err)
	}
	if err := stream.Send(&mq_pb.SubscriberToSubCoordinatorRequest{
		Message: &mq_pb.SubscriberToSubCoordinatorRequest_Init{
			Init: &mq_pb.SubscriberToSubCoordinatorRequest_InitMessage{
				ConsumerGroup: "cg-red", ConsumerGroupInstanceId: "inst-1",
				Topic:             &schema_pb.Topic{Namespace: "ns", Name: "evt"},
				MaxPartitionCount: 8, RebalanceSeconds: 30,
			},
		},
	}); err != nil {
		t.Fatalf("send init: %v", err)
	}

	// First response: Assignment variant.
	r1, err := stream.Recv()
	if err != nil {
		t.Fatalf("recv assignment: %v", err)
	}
	asg := r1.GetAssignment()
	if asg == nil {
		t.Fatalf("variant collapsed: want Assignment, got %T", r1.Message)
	}
	if asg.PartitionAssignment == nil {
		t.Fatal("Assignment.PartitionAssignment nil")
	}
	if asg.PartitionAssignment.LeaderBroker != "lead-X" || asg.PartitionAssignment.FollowerBroker != "foll-Y" {
		t.Fatalf("assignment brokers = %+v", asg.PartitionAssignment)
	}
	if p := asg.PartitionAssignment.Partition; p == nil || p.RingSize != 2048 || p.RangeStart != 7 || p.RangeStop != 9 || p.UnixTimeNs != 11 {
		t.Fatalf("assignment partition = %+v", p)
	}

	// Second response: UnAssignment variant.
	r2, err := stream.Recv()
	if err != nil {
		t.Fatalf("recv unassignment: %v", err)
	}
	un := r2.GetUnAssignment()
	if un == nil {
		t.Fatalf("variant collapsed: want UnAssignment, got %T", r2.Message)
	}
	if p := un.Partition; p == nil || p.RingSize != 4096 || p.RangeStart != 1 || p.RangeStop != 2 || p.UnixTimeNs != 3 {
		t.Fatalf("unassignment partition = %+v", p)
	}

	// Request direction: AckAssignment with a Partition.
	if err := stream.Send(&mq_pb.SubscriberToSubCoordinatorRequest{
		Message: &mq_pb.SubscriberToSubCoordinatorRequest_AckAssignment{
			AckAssignment: &mq_pb.SubscriberToSubCoordinatorRequest_AckAssignmentMessage{
				Partition: &schema_pb.Partition{RingSize: 64, RangeStart: 4, RangeStop: 5, UnixTimeNs: 6},
			},
		},
	}); err != nil {
		t.Fatalf("send ackassign: %v", err)
	}
	_ = stream.CloseSend()

	// Give the engine the EOF-driven loop a moment by draining.
	for {
		if _, err := stream.Recv(); err != nil {
			break
		}
	}
	if e.gotSubCoordInitCG != "cg-red" {
		t.Fatalf("init cg = %q", e.gotSubCoordInitCG)
	}
	if p := e.gotAckAssignPart; p == nil || p.RingSize != 64 || p.RangeStop != 5 {
		t.Fatalf("ackassign partition = %+v", p)
	}
}

// Vector 3: SubscribeFollowMe client-stream SendAndClose half-close.
func TestRedSubscribeFollowMeClientStream(t *testing.T) {
	e := &redEngine{}
	cli, stop := serveRedEngine(t, e)
	defer stop()

	stream, err := cli.SubscribeFollowMe(context.Background())
	if err != nil {
		t.Fatalf("SubscribeFollowMe: %v", err)
	}
	if err := stream.Send(&mq_pb.SubscribeFollowMeRequest{
		Message: &mq_pb.SubscribeFollowMeRequest_Init{
			Init: &mq_pb.SubscribeFollowMeRequest_InitMessage{
				Topic: &schema_pb.Topic{Namespace: "ns", Name: "t"}, ConsumerGroup: "cg-fm",
			},
		},
	}); err != nil {
		t.Fatalf("send init: %v", err)
	}
	for _, ts := range []int64{100, 200, 300} {
		if err := stream.Send(&mq_pb.SubscribeFollowMeRequest{
			Message: &mq_pb.SubscribeFollowMeRequest_Ack{
				Ack: &mq_pb.SubscribeFollowMeRequest_AckMessage{TsNs: ts},
			},
		}); err != nil {
			t.Fatalf("send ack %d: %v", ts, err)
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("CloseAndRecv: %v", err)
	}
	if resp.AckTsNs != 999 {
		t.Fatalf("summary ack = %d, want 999", resp.AckTsNs)
	}
	if e.gotFollowMeInitCG != "cg-fm" {
		t.Fatalf("followme init cg = %q", e.gotFollowMeInitCG)
	}
	if len(e.gotFollowMeAcks) != 3 || e.gotFollowMeAcks[0] != 100 || e.gotFollowMeAcks[2] != 300 {
		t.Fatalf("followme acks = %v", e.gotFollowMeAcks)
	}
}

// Vector 4: PublishFollowMe bidi (Init/Data/Flush/Close variants).
func TestRedPublishFollowMeBidi(t *testing.T) {
	e := &redEngine{}
	cli, stop := serveRedEngine(t, e)
	defer stop()

	stream, err := cli.PublishFollowMe(context.Background())
	if err != nil {
		t.Fatalf("PublishFollowMe: %v", err)
	}
	if err := stream.Send(&mq_pb.PublishFollowMeRequest{
		Message: &mq_pb.PublishFollowMeRequest_Init{
			Init: &mq_pb.PublishFollowMeRequest_InitMessage{
				Topic:     &schema_pb.Topic{Namespace: "ns", Name: "follow"},
				Partition: &schema_pb.Partition{RingSize: 8},
			},
		},
	}); err != nil {
		t.Fatalf("send init: %v", err)
	}
	for i, v := range [][]byte{[]byte("f0"), []byte("f1")} {
		if err := stream.Send(&mq_pb.PublishFollowMeRequest{
			Message: &mq_pb.PublishFollowMeRequest_Data{Data: &mq_pb.DataMessage{Value: v, TsNs: int64(10 + i)}},
		}); err != nil {
			t.Fatalf("send data %d: %v", i, err)
		}
		ack, err := stream.Recv()
		if err != nil {
			t.Fatalf("recv ack %d: %v", i, err)
		}
		if ack.AckTsNs != int64(10+i) {
			t.Fatalf("ack %d = %d", i, ack.AckTsNs)
		}
	}
	if err := stream.Send(&mq_pb.PublishFollowMeRequest{
		Message: &mq_pb.PublishFollowMeRequest_Flush{
			Flush: &mq_pb.PublishFollowMeRequest_FlushMessage{TsNs: 555},
		},
	}); err != nil {
		t.Fatalf("send flush: %v", err)
	}
	if err := stream.Send(&mq_pb.PublishFollowMeRequest{
		Message: &mq_pb.PublishFollowMeRequest_Close{Close: &mq_pb.PublishFollowMeRequest_CloseMessage{}},
	}); err != nil {
		t.Fatalf("send close: %v", err)
	}
	_ = stream.CloseSend()
	// Drain to let the engine process flush+close before assertions.
	for {
		if _, err := stream.Recv(); err != nil {
			break
		}
	}
	if e.gotPubFollowInit.GetName() != "follow" {
		t.Fatalf("pubfollow init = %+v", e.gotPubFollowInit)
	}
	if len(e.gotPubFollowData) != 2 || !bytes.Equal(e.gotPubFollowData[0], []byte("f0")) {
		t.Fatalf("pubfollow data = %q", e.gotPubFollowData)
	}
	if len(e.gotPubFollowFlush) != 1 || e.gotPubFollowFlush[0] != 555 {
		t.Fatalf("pubfollow flush = %v", e.gotPubFollowFlush)
	}
	if !e.gotPubFollowClose {
		t.Fatal("pubfollow close not seen")
	}
}

// Vector 5a: empty/None oneof on the request side — a request with NO variant set
// must NOT crash the server adapter (it should encode MessageNone and the engine
// switch falls through). We send a SubscribeMessage with a nil Message after init.
func TestRedNoneOneofRequestSurvives(t *testing.T) {
	e := &redEngine{}
	cli, stop := serveRedEngine(t, e)
	defer stop()

	// PublisherToPubBalancer with an empty (no-variant) request frame: must not panic
	// the server; the engine simply sees neither Init nor Stats.
	stream, err := cli.PublisherToPubBalancer(context.Background())
	if err != nil {
		t.Fatalf("PublisherToPubBalancer: %v", err)
	}
	if err := stream.Send(&mq_pb.PublisherToPubBalancerRequest{}); err != nil {
		t.Fatalf("send empty: %v", err)
	}
	// Follow with a real Init so we can confirm the stream is still alive.
	if err := stream.Send(&mq_pb.PublisherToPubBalancerRequest{
		Message: &mq_pb.PublisherToPubBalancerRequest_Init{
			Init: &mq_pb.PublisherToPubBalancerRequest_InitMessage{Broker: "after-empty"},
		},
	}); err != nil {
		t.Fatalf("send init after empty: %v", err)
	}
	_ = stream.CloseSend()
	for {
		if _, err := stream.Recv(); err != nil {
			break
		}
	}
	if e.gotBalancerInitBroker != "after-empty" {
		t.Fatalf("stream died after empty oneof; init broker = %q", e.gotBalancerInitBroker)
	}
}

// Vector 5b: nil-Partition / nil-Topic in a oneof payload must round-trip as nil,
// not as a zero-value struct. SubscriberToSubCoordinator AckAssignment with a nil
// Partition.
func TestRedNilPartitionInOneof(t *testing.T) {
	e := &redEngine{}
	cli, stop := serveRedEngine(t, e)
	defer stop()

	stream, err := cli.SubscriberToSubCoordinator(context.Background())
	if err != nil {
		t.Fatalf("SubscriberToSubCoordinator: %v", err)
	}
	// Init (engine will send two responses we drain).
	if err := stream.Send(&mq_pb.SubscriberToSubCoordinatorRequest{
		Message: &mq_pb.SubscriberToSubCoordinatorRequest_Init{
			Init: &mq_pb.SubscriberToSubCoordinatorRequest_InitMessage{ConsumerGroup: "cg-nil"},
		},
	}); err != nil {
		t.Fatalf("send init: %v", err)
	}
	_, _ = stream.Recv()
	_, _ = stream.Recv()
	// AckAssignment with nil Partition.
	if err := stream.Send(&mq_pb.SubscriberToSubCoordinatorRequest{
		Message: &mq_pb.SubscriberToSubCoordinatorRequest_AckAssignment{
			AckAssignment: &mq_pb.SubscriberToSubCoordinatorRequest_AckAssignmentMessage{Partition: nil},
		},
	}); err != nil {
		t.Fatalf("send ackassign nil part: %v", err)
	}
	_ = stream.CloseSend()
	for {
		if _, err := stream.Recv(); err != nil {
			break
		}
	}
	if e.gotAckAssignPart != nil {
		t.Fatalf("nil Partition round-tripped as non-nil: %+v", e.gotAckAssignPart)
	}
}
