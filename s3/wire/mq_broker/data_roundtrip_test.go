// Pure-wire (path-B) certification for the streaming DATA path: DataMessage
// and its optional ControlMessage child. These are the live message-queue
// payloads carried by the bidi PublishMessage/SubscribeMessage and the
// leader/follower PublishFollowMe/SubscribeFollowMe streams — the highest-risk
// surface of the broker RPC cutover, where a framing or half-close error
// "silently drops or reorders messages" (per the cutover playbook).
//
// Certifying the codec here, BEFORE the streaming handlers are ported, proves
// the bytes carry binary key/value verbatim, preserve TsNs, and signal the
// publisher half-close (Ctrl.IsClose) presence correctly — so the eventual
// stream port can trust the payload layer and focus only on the transport
// Send/Recv/CloseSend loop.
package mq_brokerwire

import (
	"bytes"
	"testing"
)

// TestDataMessageRoundTrip proves binary key/value (including embedded NUL
// bytes and empty slices) and the int64 timestamp survive build -> wrap with
// byte-for-byte fidelity. Reordering or truncating these is the exact failure
// mode that drops queue records.
func TestDataMessageRoundTrip(t *testing.T) {
	cases := []DataMessageInput{
		{Key: []byte("k1"), Value: []byte("hello"), TsNs: 1},
		{Key: []byte{0x00, 0xff, 0x00}, Value: []byte{0xde, 0xad, 0x00, 0xbe, 0xef}, TsNs: 1<<62},
		{Key: nil, Value: nil, TsNs: 0},
		{Key: []byte(""), Value: bytes.Repeat([]byte("x"), 4096), TsNs: -1},
	}
	for i, in := range cases {
		buf := NewDataMessage(in)
		v, err := WrapDataMessage(buf)
		if err != nil {
			t.Fatalf("case %d WrapDataMessage: %v", i, err)
		}
		if !bytes.Equal(v.Key(), in.Key) && !(len(v.Key()) == 0 && len(in.Key) == 0) {
			t.Fatalf("case %d Key = %x, want %x", i, v.Key(), in.Key)
		}
		if !bytes.Equal(v.Value(), in.Value) && !(len(v.Value()) == 0 && len(in.Value) == 0) {
			t.Fatalf("case %d Value len = %d, want %d", i, len(v.Value()), len(in.Value))
		}
		if v.TsNs() != in.TsNs {
			t.Fatalf("case %d TsNs = %d, want %d", i, v.TsNs(), in.TsNs)
		}
		// A bare data message carries no control frame.
		if _, ok := v.Ctrl(); ok {
			t.Fatalf("case %d: bare DataMessage reported a Ctrl child", i)
		}
	}
}

// TestDataMessageControlClose proves the publisher half-close signal — a
// DataMessage whose Ctrl child has IsClose=true — round-trips with the
// presence bool set and the flag readable. The streaming pub loop relies on
// exactly this to terminate the partition cleanly without losing in-flight
// records, so the codec must never drop or misreport it.
func TestDataMessageControlClose(t *testing.T) {
	ctrl := NewControlMessage(ControlMessageInput{IsClose: true, PublisherName: "pub-7"})
	buf := NewDataMessage(DataMessageInput{TsNs: 99, Ctrl: ctrl})

	v, err := WrapDataMessage(buf)
	if err != nil {
		t.Fatalf("WrapDataMessage: %v", err)
	}
	got, ok := v.Ctrl()
	if !ok {
		t.Fatalf("Ctrl present bool = false, want true (half-close signal lost)")
	}
	if !got.IsClose() {
		t.Fatalf("Ctrl.IsClose = false, want true")
	}
	if got.PublisherName() != "pub-7" {
		t.Fatalf("Ctrl.PublisherName = %q, want %q", got.PublisherName(), "pub-7")
	}
	if v.TsNs() != 99 {
		t.Fatalf("TsNs = %d, want 99", v.TsNs())
	}
}

// TestDataMessageControlAbsentIsNotClose guards the inverse: a non-control
// data record must NOT look like a close. A false positive here would
// truncate a live stream.
func TestDataMessageControlAbsentIsNotClose(t *testing.T) {
	buf := NewDataMessage(DataMessageInput{Key: []byte("k"), Value: []byte("v"), TsNs: 5})
	v, err := WrapDataMessage(buf)
	if err != nil {
		t.Fatalf("WrapDataMessage: %v", err)
	}
	if _, ok := v.Ctrl(); ok {
		t.Fatalf("data record without Ctrl reported a control frame (would truncate stream)")
	}
}
