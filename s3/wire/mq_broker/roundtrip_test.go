// Pure-wire (path-B) roundtrip certification for the HanzoMessaging ZAP
// codecs. NO protobuf, NO mq_pb, NO proto.Marshal, NO bridge — every
// message is built via New<Msg>(input) -> []byte and read via Wrap<Msg>(buf)
// -> view, exactly the zero-copy contract the broker RPC cutover relies on.
//
// This test is the gate the migration playbook names ("the wire roundtrip
// test"): it certifies that the generated mq_brokerwire codecs (and their
// mq_schemawire children) survive build -> wrap -> read, and that the
// DispatchHanzoMessaging request/response envelope path is sound, BEFORE any
// handler or client is moved off gRPC.
package mq_brokerwire

import (
	"testing"

	mq_schemawire "github.com/hanzoai/s3/s3/wire/mq_schema"
	"github.com/zap-proto/go/rpc"
)

// TestFindBrokerLeaderRoundTrip exercises the simplest unary pair: a single
// scalar string in the request and the response. This is the canonical
// bridge-free shape (no schema children).
func TestFindBrokerLeaderRoundTrip(t *testing.T) {
	reqBuf := NewFindBrokerLeaderRequest(FindBrokerLeaderRequestInput{FilerGroup: "group-A"})
	req, err := WrapFindBrokerLeaderRequest(reqBuf)
	if err != nil {
		t.Fatalf("WrapFindBrokerLeaderRequest: %v", err)
	}
	if got := req.FilerGroup(); got != "group-A" {
		t.Fatalf("FilerGroup = %q, want %q", got, "group-A")
	}

	respBuf := NewFindBrokerLeaderResponse(FindBrokerLeaderResponseInput{Broker: "broker-1:17777"})
	resp, err := WrapFindBrokerLeaderResponse(respBuf)
	if err != nil {
		t.Fatalf("WrapFindBrokerLeaderResponse: %v", err)
	}
	if got := resp.Broker(); got != "broker-1:17777" {
		t.Fatalf("Broker = %q, want %q", got, "broker-1:17777")
	}
}

// TestListTopicsNestedSchemaRoundTrip exercises the outer-message +
// mq_schemawire child composition: a repeated schema_pb.Topic carried as
// opaque child ZAP buffers and re-wrapped on the read side. This proves the
// two wire packages compose with zero copy and no protobuf in between.
func TestListTopicsNestedSchemaRoundTrip(t *testing.T) {
	want := []mq_schemawire.TopicInput{
		{Namespace: "ns1", Name: "orders"},
		{Namespace: "ns2", Name: "events"},
		{Namespace: "", Name: "anon"},
	}

	children := make([][]byte, len(want))
	for i, w := range want {
		children[i] = mq_schemawire.NewTopic(w)
	}

	respBuf := NewListTopicsResponse(ListTopicsResponseInput{Topics: children})
	resp, err := WrapListTopicsResponse(respBuf)
	if err != nil {
		t.Fatalf("WrapListTopicsResponse: %v", err)
	}
	if got := resp.TopicsLen(); got != len(want) {
		t.Fatalf("TopicsLen = %d, want %d", got, len(want))
	}
	for i := range want {
		childBuf := resp.TopicAt(i)
		if childBuf == nil {
			t.Fatalf("TopicAt(%d) = nil", i)
		}
		topic, err := mq_schemawire.WrapTopic(childBuf)
		if err != nil {
			t.Fatalf("WrapTopic(%d): %v", i, err)
		}
		if topic.Namespace() != want[i].Namespace || topic.Name() != want[i].Name {
			t.Fatalf("topic[%d] = {%q,%q}, want {%q,%q}",
				i, topic.Namespace(), topic.Name(), want[i].Namespace, want[i].Name)
		}
	}

	// Out-of-range index must be nil, never a panic.
	if resp.TopicAt(len(want)) != nil {
		t.Fatalf("TopicAt(out-of-range) should be nil")
	}
}

// TestListTopicsEmpty proves an empty repeated field round-trips as length 0,
// not a parse error.
func TestListTopicsEmpty(t *testing.T) {
	respBuf := NewListTopicsResponse(ListTopicsResponseInput{Topics: nil})
	resp, err := WrapListTopicsResponse(respBuf)
	if err != nil {
		t.Fatalf("WrapListTopicsResponse(empty): %v", err)
	}
	if got := resp.TopicsLen(); got != 0 {
		t.Fatalf("TopicsLen(empty) = %d, want 0", got)
	}
}

// fakeHandler is a minimal wire-native HanzoMessagingHandler used to prove the
// DispatchHanzoMessaging request/response envelope path end-to-end with zero
// protobuf. Only the methods this test invokes do real work; the rest satisfy
// the interface and report "not implemented" via an empty response.
type fakeHandler struct {
	gotFilerGroup string
}

func notImpl([]byte) ([]byte, error) { return []byte{}, nil }

func (h *fakeHandler) FindBrokerLeader(req []byte) ([]byte, error) {
	v, err := WrapFindBrokerLeaderRequest(req)
	if err != nil {
		return nil, err
	}
	h.gotFilerGroup = v.FilerGroup()
	return NewFindBrokerLeaderResponse(FindBrokerLeaderResponseInput{Broker: "leader:" + v.FilerGroup()}), nil
}

func (h *fakeHandler) PublisherToPubBalancer(req []byte) ([]byte, error)     { return notImpl(req) }
func (h *fakeHandler) BalanceTopics(req []byte) ([]byte, error)              { return notImpl(req) }
func (h *fakeHandler) ListTopics(req []byte) ([]byte, error)                 { return notImpl(req) }
func (h *fakeHandler) TopicExists(req []byte) ([]byte, error)               { return notImpl(req) }
func (h *fakeHandler) ConfigureTopic(req []byte) ([]byte, error)            { return notImpl(req) }
func (h *fakeHandler) LookupTopicBrokers(req []byte) ([]byte, error)        { return notImpl(req) }
func (h *fakeHandler) GetTopicConfiguration(req []byte) ([]byte, error)     { return notImpl(req) }
func (h *fakeHandler) GetTopicPublishers(req []byte) ([]byte, error)        { return notImpl(req) }
func (h *fakeHandler) GetTopicSubscribers(req []byte) ([]byte, error)       { return notImpl(req) }
func (h *fakeHandler) AssignTopicPartitions(req []byte) ([]byte, error)     { return notImpl(req) }
func (h *fakeHandler) ClosePublishers(req []byte) ([]byte, error)           { return notImpl(req) }
func (h *fakeHandler) CloseSubscribers(req []byte) ([]byte, error)          { return notImpl(req) }
func (h *fakeHandler) SubscriberToSubCoordinator(req []byte) ([]byte, error) { return notImpl(req) }
func (h *fakeHandler) PublishMessage(req []byte) ([]byte, error)            { return notImpl(req) }
func (h *fakeHandler) SubscribeMessage(req []byte) ([]byte, error)          { return notImpl(req) }
func (h *fakeHandler) PublishFollowMe(req []byte) ([]byte, error)           { return notImpl(req) }
func (h *fakeHandler) SubscribeFollowMe(req []byte) ([]byte, error)         { return notImpl(req) }
func (h *fakeHandler) FetchMessage(req []byte) ([]byte, error)              { return notImpl(req) }
func (h *fakeHandler) GetUnflushedMessages(req []byte) ([]byte, error)      { return notImpl(req) }
func (h *fakeHandler) GetPartitionRangeInfo(req []byte) ([]byte, error)     { return notImpl(req) }

// TestDispatchEnvelopeRoundTrip drives a request envelope through
// DispatchHanzoMessaging and parses the response envelope — the exact server
// path transport.Serve will run, proving routing-by-ordinal and the
// request/response framing work with wire-native payloads only.
func TestDispatchEnvelopeRoundTrip(t *testing.T) {
	h := &fakeHandler{}

	payload := NewFindBrokerLeaderRequest(FindBrokerLeaderRequestInput{FilerGroup: "g7"})
	envelope := rpc.BuildRequest(rpc.Call{
		Method:    HanzoMessagingFindBrokerLeaderOrdinal,
		PromiseID: 42,
		Target:    rpc.NoTarget,
		Payload:   payload,
	})

	respEnvelope, err := DispatchHanzoMessaging(h, envelope)
	if err != nil {
		t.Fatalf("DispatchHanzoMessaging: %v", err)
	}
	if h.gotFilerGroup != "g7" {
		t.Fatalf("handler saw FilerGroup %q, want %q", h.gotFilerGroup, "g7")
	}

	resp, err := rpc.ParseResponse(respEnvelope)
	if err != nil {
		t.Fatalf("ParseResponse: %v", err)
	}
	if resp.Status != rpc.StatusOK {
		t.Fatalf("response status = %d, want %d", resp.Status, rpc.StatusOK)
	}
	if resp.PromiseID != 42 {
		t.Fatalf("response PromiseID = %d, want 42", resp.PromiseID)
	}
	out, err := WrapFindBrokerLeaderResponse(resp.Body)
	if err != nil {
		t.Fatalf("WrapFindBrokerLeaderResponse(resp.Body): %v", err)
	}
	if got := out.Broker(); got != "leader:g7" {
		t.Fatalf("Broker = %q, want %q", got, "leader:g7")
	}
}

// TestDispatchUnknownOrdinal proves an unrouted method yields StatusNotFound,
// not a panic or a silent OK — fail-closed routing.
func TestDispatchUnknownOrdinal(t *testing.T) {
	h := &fakeHandler{}
	envelope := rpc.BuildRequest(rpc.Call{
		Method:    9999, // no such ordinal
		PromiseID: 7,
		Target:    rpc.NoTarget,
	})
	respEnvelope, err := DispatchHanzoMessaging(h, envelope)
	if err != nil {
		t.Fatalf("DispatchHanzoMessaging(unknown): %v", err)
	}
	resp, err := rpc.ParseResponse(respEnvelope)
	if err != nil {
		t.Fatalf("ParseResponse(unknown): %v", err)
	}
	if resp.Status != rpc.StatusNotFound {
		t.Fatalf("unknown-ordinal status = %d, want %d (StatusNotFound)", resp.Status, rpc.StatusNotFound)
	}
}
