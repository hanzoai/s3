package s3api

import (
	"bytes"
	"testing"

	"github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/s3api/s3lifecycle"
	stats_collect "github.com/hanzoai/s3/s3/stats"
	s3_lifecyclewire "github.com/hanzoai/s3/s3/wire/s3_lifecycle"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// wireReq builds a zero-copy LifecycleDeleteRequest view from the given
// input — the test analogue of what the ZAP transport hands the server.
func wireReq(t *testing.T, in s3_lifecyclewire.LifecycleDeleteRequestInput) s3_lifecyclewire.LifecycleDeleteRequest {
	t.Helper()
	v, err := s3_lifecyclewire.WrapLifecycleDeleteRequest(s3_lifecyclewire.NewLifecycleDeleteRequest(in))
	if err != nil {
		t.Fatalf("wrap request: %v", err)
	}
	return v
}

func TestComputeEntryIdentity_BasicFields(t *testing.T) {
	entry := &filer_pb.Entry{
		Attributes: &filer_pb.FuseAttributes{Mtime: 1700000000, MtimeNs: 123, FileSize: 4096},
		Chunks: []*filer_pb.FileChunk{
			{FileId: "1,abc"},
			{FileId: "1,def"},
		},
	}
	id := computeEntryIdentity(entry)
	want := int64(1700000000)*int64(1e9) + int64(123)
	if id.MtimeNs != want {
		t.Fatalf("MtimeNs want %d, got %d", want, id.MtimeNs)
	}
	if id.Size != 4096 {
		t.Fatalf("Size want 4096, got %d", id.Size)
	}
	if id.HeadFid != "1,abc" {
		t.Fatalf("HeadFid want 1,abc, got %s", id.HeadFid)
	}
}

func TestComputeEntryIdentity_NilSafeMissingChunks(t *testing.T) {
	if got := computeEntryIdentity(nil); got != nil {
		t.Fatalf("nil entry should return nil, got %v", got)
	}
	id := computeEntryIdentity(&filer_pb.Entry{})
	if id == nil {
		t.Fatalf("entry with nil Attributes should still produce identity")
	}
	if id.HeadFid != "" {
		t.Fatalf("missing chunks should yield empty HeadFid, got %s", id.HeadFid)
	}
}

func TestHashExtended_OrderStable(t *testing.T) {
	a := map[string][]byte{"k1": []byte("v1"), "k2": []byte("v2")}
	b := map[string][]byte{"k2": []byte("v2"), "k1": []byte("v1")}
	if !bytes.Equal(s3lifecycle.HashExtended(a), s3lifecycle.HashExtended(b)) {
		t.Fatalf("hash should be insensitive to map iteration order")
	}
}

func TestHashExtended_DelimiterCollisionResistant(t *testing.T) {
	// Naively concatenated: "k1=v1k2v2" could collide with "k1=v1k" / "2v2".
	// Length-prefix encoding must keep them apart.
	a := map[string][]byte{"k1": []byte("v1"), "k2": []byte("v2")}
	b := map[string][]byte{"k1": []byte("v1k2v2")}
	if bytes.Equal(s3lifecycle.HashExtended(a), s3lifecycle.HashExtended(b)) {
		t.Fatalf("delimiter-forged Extended payloads must not collide")
	}
}

func TestHashExtended_NilEqualsEmpty(t *testing.T) {
	if got := s3lifecycle.HashExtended(nil); len(got) != 0 {
		t.Fatalf("nil should produce zero-length hash, got %d bytes", len(got))
	}
	if got := s3lifecycle.HashExtended(map[string][]byte{}); len(got) != 0 {
		t.Fatalf("empty map should produce zero-length hash, got %d bytes", len(got))
	}
}

func TestIdentityMatches_NilWantTreatedAsMatch(t *testing.T) {
	// Bootstrap callers that don't yet have an identity to CAS against
	// leave expected_identity unset; the server treats this as "no CAS".
	live := &entryIdentity{MtimeNs: 1, Size: 2}
	if !identityMatches(live, wireReq(t, s3_lifecyclewire.LifecycleDeleteRequestInput{})) {
		t.Fatalf("absent want should match")
	}
}

func TestIdentityMatches_NilLiveDoesNotMatch(t *testing.T) {
	req := wireReq(t, s3_lifecyclewire.LifecycleDeleteRequestInput{
		ExpectedIdentity: &s3_lifecyclewire.EntryIdentityInput{MtimeNs: 1},
	})
	if identityMatches(nil, req) {
		t.Fatalf("nil live should not match a populated want")
	}
}

func TestIdentityMatches_AllFieldsCompared(t *testing.T) {
	base := s3_lifecyclewire.EntryIdentityInput{MtimeNs: 100, Size: 2048, HeadFid: "1,abc", ExtendedHash: []byte{0x01, 0x02}}
	req := wireReq(t, s3_lifecyclewire.LifecycleDeleteRequestInput{ExpectedIdentity: &base})
	cases := []struct {
		name string
		live *entryIdentity
		want bool
	}{
		{"identical", &entryIdentity{MtimeNs: 100, Size: 2048, HeadFid: "1,abc", ExtendedHash: []byte{0x01, 0x02}}, true},
		{"mtime-drift", &entryIdentity{MtimeNs: 101, Size: 2048, HeadFid: "1,abc", ExtendedHash: []byte{0x01, 0x02}}, false},
		{"size-drift", &entryIdentity{MtimeNs: 100, Size: 2049, HeadFid: "1,abc", ExtendedHash: []byte{0x01, 0x02}}, false},
		{"fid-drift", &entryIdentity{MtimeNs: 100, Size: 2048, HeadFid: "1,xyz", ExtendedHash: []byte{0x01, 0x02}}, false},
		{"extended-drift", &entryIdentity{MtimeNs: 100, Size: 2048, HeadFid: "1,abc", ExtendedHash: []byte{0x03, 0x04}}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := identityMatches(c.live, req); got != c.want {
				t.Fatalf("want %v, got %v", c.want, got)
			}
		})
	}
}

func TestLifecycleDelete_RejectsEmptyRequest(t *testing.T) {
	s := &S3ApiServer{}
	res, err := s.LifecycleDelete(wireReq(t, s3_lifecyclewire.LifecycleDeleteRequestInput{}))
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if res.Outcome != s3_lifecyclewire.LifecycleDeleteOutcomeBlocked {
		t.Fatalf("empty request should be BLOCKED, got %v", res.Outcome)
	}
}

func TestLifecycleAbortMPU_RejectsTraversalUploadIDs(t *testing.T) {
	// "." and ".." pass the no-slash check but resolve to the bucket
	// root via util.JoinPath; they must be rejected before any rm call.
	s := &S3ApiServer{}
	cases := []string{
		"",
		".uploads",
		".uploads/",
		".uploads/.",
		".uploads/..",
		".uploads/u1/extra",
	}
	for _, p := range cases {
		t.Run(p, func(t *testing.T) {
			res, err := s.LifecycleDelete(wireReq(t, s3_lifecyclewire.LifecycleDeleteRequestInput{
				Bucket:     "bk",
				ObjectPath: p,
				ActionKind: s3_lifecyclewire.ActionKindAbortMpu,
			}))
			if err != nil {
				t.Fatalf("unexpected transport error: %v", err)
			}
			if res.Outcome != s3_lifecyclewire.LifecycleDeleteOutcomeBlocked {
				t.Fatalf("path %q: outcome=%v reason=%q, want BLOCKED",
					p, res.Outcome, res.Reason)
			}
		})
	}
}

func TestLifecycleDispatch_AbortMPUAfterFetchIsBlocked(t *testing.T) {
	// LifecycleDelete routes ABORT_MPU to lifecycleAbortMPU before
	// getObjectEntry; reaching lifecycleDispatch with ABORT_MPU means
	// some caller bypassed that route. Defensive BLOCKED so a
	// regression there can't accidentally rm a real object via the
	// expiration paths.
	s := &S3ApiServer{}
	res, err := s.lifecycleDispatch(nil, wireReq(t, s3_lifecyclewire.LifecycleDeleteRequestInput{
		Bucket:     "bk",
		ObjectPath: "k",
		ActionKind: s3_lifecyclewire.ActionKindAbortMpu,
	}), &filer_pb.Entry{Attributes: &filer_pb.FuseAttributes{}})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if res.Outcome != s3_lifecyclewire.LifecycleDeleteOutcomeBlocked {
		t.Fatalf("ABORT_MPU at dispatch should be BLOCKED, got %v reason=%q", res.Outcome, res.Reason)
	}
	if !contains(res.Reason, "ABORT_MPU dispatched after fetch") {
		t.Fatalf("reason should name the route bypass, got %q", res.Reason)
	}
}

func TestLifecycleDispatch_UnknownActionKindIsBlocked(t *testing.T) {
	// An ActionKind value the schema doesn't define yet must produce a
	// FATAL outcome rather than fall through to a default delete path.
	s := &S3ApiServer{}
	const bogus uint32 = 999
	res, err := s.lifecycleDispatch(nil, wireReq(t, s3_lifecyclewire.LifecycleDeleteRequestInput{
		Bucket:     "bk",
		ObjectPath: "k",
		ActionKind: bogus,
	}), &filer_pb.Entry{Attributes: &filer_pb.FuseAttributes{}})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if res.Outcome != s3_lifecyclewire.LifecycleDeleteOutcomeBlocked {
		t.Fatalf("unknown action kind should be BLOCKED, got %v reason=%q", res.Outcome, res.Reason)
	}
	if !contains(res.Reason, "unknown action_kind") {
		t.Fatalf("reason should name the unknown kind, got %q", res.Reason)
	}
}

func TestLifecycleDispatch_NoncurrentRequiresVersionID(t *testing.T) {
	// Noncurrent / EXPIRED_DELETE_MARKER target a specific version; an
	// empty version_id is a writer-side bug and must be rejected before
	// any filer call. This pinning keeps the early-return in place so
	// a refactor doesn't accidentally let the empty-version_id path
	// reach deleteSpecificObjectVersion.
	s := &S3ApiServer{}
	for _, kind := range []uint32{
		s3_lifecyclewire.ActionKindNoncurrentDays,
		s3_lifecyclewire.ActionKindNewerNoncurrent,
		s3_lifecyclewire.ActionKindExpiredDeleteMarker,
	} {
		t.Run(actionKindLabel(kind), func(t *testing.T) {
			res, err := s.lifecycleDispatch(nil, wireReq(t, s3_lifecyclewire.LifecycleDeleteRequestInput{
				Bucket:     "bk",
				ObjectPath: "k",
				ActionKind: kind,
				// VersionID intentionally empty
			}), &filer_pb.Entry{Attributes: &filer_pb.FuseAttributes{}})
			if err != nil {
				t.Fatalf("unexpected transport error: %v", err)
			}
			if res.Outcome != s3_lifecyclewire.LifecycleDeleteOutcomeBlocked {
				t.Fatalf("kind %v with empty version_id should be BLOCKED, got %v reason=%q",
					kind, res.Outcome, res.Reason)
			}
			if !contains(res.Reason, "version_id required") {
				t.Fatalf("reason should name the missing version_id, got %q", res.Reason)
			}
		})
	}
}

// contains is a tiny helper so the tests above don't pull in strings
// just for a substring check.
func contains(haystack, needle string) bool {
	if len(needle) == 0 {
		return true
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

func TestEntryUsesMetadataOnlyDelete(t *testing.T) {
	cases := []struct {
		name  string
		entry *filer_pb.Entry
		want  bool
	}{
		{
			name:  "nil entry",
			entry: nil,
			want:  false,
		},
		{
			name:  "nil attributes",
			entry: &filer_pb.Entry{},
			want:  false,
		},
		{
			name:  "TtlSec=0 (no per-write stamp)",
			entry: &filer_pb.Entry{Attributes: &filer_pb.FuseAttributes{TtlSec: 0}},
			want:  false,
		},
		{
			name:  "TtlSec>0 (PR 9377 stamped a fast-path TTL)",
			entry: &filer_pb.Entry{Attributes: &filer_pb.FuseAttributes{TtlSec: 86400}},
			want:  true,
		},
		{
			name:  "TtlSec<0 should not happen but must not flip the path on",
			entry: &filer_pb.Entry{Attributes: &filer_pb.FuseAttributes{TtlSec: -1}},
			want:  false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := entryUsesMetadataOnlyDelete(c.entry); got != c.want {
				t.Fatalf("want %v, got %v", c.want, got)
			}
		})
	}
}

func TestRecordMetadataOnlyIf_OnlyFiresWhenOn(t *testing.T) {
	// Counter must increment exactly once per (bucket, hex(rule_hash))
	// when on=true, and not at all when on=false. Other lifecycle paths
	// in the same suite share the global counter — use distinct bucket
	// names per test so series don't bleed.
	c := stats_collect.S3LifecycleMetadataOnlyCounter
	bucket := "bk-counter-on"
	hash := []byte{0xde, 0xad, 0xbe, 0xef, 0x01, 0x02, 0x03, 0x04}
	hexHash := "deadbeef01020304"

	before := testutil.ToFloat64(c.WithLabelValues(bucket, hexHash))
	recordMetadataOnlyIf(true, wireReq(t, s3_lifecyclewire.LifecycleDeleteRequestInput{
		Bucket:   bucket,
		RuleHash: hash,
	}))
	if got := testutil.ToFloat64(c.WithLabelValues(bucket, hexHash)); got != before+1 {
		t.Fatalf("on=true should bump by 1; before=%v after=%v", before, got)
	}

	beforeOff := testutil.ToFloat64(c.WithLabelValues("bk-counter-off", hexHash))
	recordMetadataOnlyIf(false, wireReq(t, s3_lifecyclewire.LifecycleDeleteRequestInput{
		Bucket:   "bk-counter-off",
		RuleHash: hash,
	}))
	if got := testutil.ToFloat64(c.WithLabelValues("bk-counter-off", hexHash)); got != beforeOff {
		t.Fatalf("on=false should not bump; before=%v after=%v", beforeOff, got)
	}
}

func TestRecordMetadataOnlyIf_EmptyRuleHashCollapsesToEmptyLabel(t *testing.T) {
	// Bootstrap or test paths may not stamp a rule hash; the label
	// must end up as an empty string rather than panicking on
	// hex.EncodeToString(nil).
	c := stats_collect.S3LifecycleMetadataOnlyCounter
	bucket := "bk-counter-emptyhash"
	before := testutil.ToFloat64(c.WithLabelValues(bucket, ""))
	recordMetadataOnlyIf(true, wireReq(t, s3_lifecyclewire.LifecycleDeleteRequestInput{Bucket: bucket}))
	if got := testutil.ToFloat64(c.WithLabelValues(bucket, "")); got != before+1 {
		t.Fatalf("nil rule_hash should produce empty-label series; before=%v after=%v", before, got)
	}
}
