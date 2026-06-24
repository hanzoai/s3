// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package s3_lifecyclewire

import (
	"bytes"
	"sync"
	"testing"
)

// fakeDeleter is an in-memory LifecycleDeleter for the round-trip test. It
// records the last request it saw (read out of the zero-copy view) and returns
// a fixed verdict, so the test can assert that every field — including the
// enum action_kind and the nested EntryIdentity CAS witness — survived the wire
// crossing intact. The recorded fields are written on the server goroutine and
// read on the test goroutine, so a mutex guards them (the response round-trip
// orders the call's completion, but not the field writes).
type fakeDeleter struct {
	mu sync.Mutex

	gotBucket       string
	gotObjectPath   string
	gotVersionID    string
	gotRuleHash     []byte
	gotActionKind   uint32
	gotShard        string
	gotDelay        int64
	gotTsNs         int64
	gotOffset       int64
	gotSnapshotID   uint64
	gotIdentMtime   int64
	gotIdentSize    int64
	gotIdentHeadFid string
	gotIdentHash    []byte

	outcome uint32
	reason  string
}

func (f *fakeDeleter) LifecycleDelete(req LifecycleDeleteRequest) (LifecycleDeleteResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	// Copy out of the zero-copy view (it aliases the request buffer).
	f.gotBucket = req.Bucket()
	f.gotObjectPath = req.ObjectPath()
	f.gotVersionID = req.VersionID()
	f.gotRuleHash = append([]byte(nil), req.RuleHash()...)
	f.gotActionKind = req.ActionKind()
	f.gotShard = req.StreamShard()
	f.gotDelay = req.StreamDelaySeconds()
	f.gotTsNs = req.StreamPositionTsNs()
	f.gotOffset = req.StreamPositionOffset()
	f.gotSnapshotID = req.EngineSnapshotID()

	id := req.ExpectedIdentity()
	f.gotIdentMtime = id.MtimeNs()
	f.gotIdentSize = id.Size()
	f.gotIdentHeadFid = id.HeadFid()
	f.gotIdentHash = append([]byte(nil), id.ExtendedHash()...)

	return LifecycleDeleteResult{Outcome: f.outcome, Reason: f.reason}, nil
}

// TestLifecycleDeleteRoundTrip proves LifecycleDelete over the canonical
// github.com/zap-proto/go transport across a real TCP socket: the request
// crosses the wire as a ZAP RPC envelope carrying the zero-copy
// LifecycleDeleteRequest payload (scalars + enum + nested EntryIdentity), the
// server dispatches it to the backend, and the LifecycleDeleteResponse comes
// back the same way — no HTTP, no protobuf, no struct marshaling.
func TestLifecycleDeleteRoundTrip(t *testing.T) {
	deleter := &fakeDeleter{
		outcome: LifecycleDeleteOutcomeRetryLater,
		reason:  "TRANSPORT_ERROR: filer unavailable",
	}

	srv, err := Serve("tcp", "127.0.0.1:0", deleter)
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	wantHash := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	wantExtHash := bytes.Repeat([]byte{0xab}, 32)
	in := LifecycleDeleteRequestInput{
		Bucket:               "blobs",
		ObjectPath:           "team-go/x.bin",
		VersionID:            "v-7",
		RuleHash:             wantHash,
		ActionKind:           ActionKindNoncurrentDays,
		StreamShard:          "filer-3",
		StreamDelaySeconds:   86400,
		StreamPositionTsNs:   1_700_000_000_000_000_000,
		StreamPositionOffset: 4242,
		EngineSnapshotID:     0xDEADBEEFCAFEBABE,
		ExpectedIdentity: &EntryIdentityInput{
			MtimeNs:      1_699_999_999_000_000_000,
			Size:         1 << 30,
			HeadFid:      "3,01637037d6",
			ExtendedHash: wantExtHash,
		},
	}

	res, err := cli.LifecycleDelete(in)
	if err != nil {
		t.Fatalf("LifecycleDelete: %v", err)
	}

	// Snapshot the server-recorded fields under the lock before asserting
	// (they were written on the server goroutine).
	deleter.mu.Lock()
	defer deleter.mu.Unlock()

	// Response survived the wire (outcome enum + reason).
	if res.Outcome != LifecycleDeleteOutcomeRetryLater {
		t.Fatalf("Outcome = %d, want %d", res.Outcome, LifecycleDeleteOutcomeRetryLater)
	}
	if res.Reason != "TRANSPORT_ERROR: filer unavailable" {
		t.Fatalf("Reason = %q", res.Reason)
	}

	// Request fields the server observed all match what the client sent.
	if deleter.gotBucket != "blobs" || deleter.gotObjectPath != "team-go/x.bin" || deleter.gotVersionID != "v-7" {
		t.Fatalf("routing mismatch: bucket=%q path=%q version=%q",
			deleter.gotBucket, deleter.gotObjectPath, deleter.gotVersionID)
	}
	if !bytes.Equal(deleter.gotRuleHash, wantHash) {
		t.Fatalf("RuleHash = %x, want %x", deleter.gotRuleHash, wantHash)
	}
	if deleter.gotActionKind != ActionKindNoncurrentDays {
		t.Fatalf("ActionKind = %d, want %d", deleter.gotActionKind, ActionKindNoncurrentDays)
	}
	if deleter.gotShard != "filer-3" || deleter.gotDelay != 86400 ||
		deleter.gotTsNs != 1_700_000_000_000_000_000 || deleter.gotOffset != 4242 {
		t.Fatalf("stream ctx mismatch: shard=%q delay=%d ts=%d off=%d",
			deleter.gotShard, deleter.gotDelay, deleter.gotTsNs, deleter.gotOffset)
	}
	if deleter.gotSnapshotID != 0xDEADBEEFCAFEBABE {
		t.Fatalf("EngineSnapshotID = %#x, want %#x", deleter.gotSnapshotID, uint64(0xDEADBEEFCAFEBABE))
	}

	// Nested EntryIdentity CAS witness round-tripped field-for-field.
	if deleter.gotIdentMtime != 1_699_999_999_000_000_000 || deleter.gotIdentSize != 1<<30 ||
		deleter.gotIdentHeadFid != "3,01637037d6" || !bytes.Equal(deleter.gotIdentHash, wantExtHash) {
		t.Fatalf("EntryIdentity mismatch: mtime=%d size=%d fid=%q hashlen=%d",
			deleter.gotIdentMtime, deleter.gotIdentSize, deleter.gotIdentHeadFid, len(deleter.gotIdentHash))
	}
}
