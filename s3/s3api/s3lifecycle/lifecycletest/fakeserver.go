// Package lifecycletest provides reusable test doubles for the lifecycle
// worker pipeline. The pieces here let component-level tests stand up the
// ZAP service boundary the worker dials at runtime without pulling in a
// real S3ApiServer or filer.
package lifecycletest

import (
	"sync"

	s3_lifecyclewire "github.com/hanzoai/s3/s3/wire/s3_lifecycle"
)

// Outcome is what the fake returns for a single LifecycleDelete call.
// Code is required (a LifecycleDeleteOutcome* constant); Reason is echoed
// verbatim into the LifecycleDeleteResult.
type Outcome struct {
	Code   uint32
	Reason string
}

// Done, NoopResolved, RetryLater, Blocked, SkippedObjectLock are
// constructors for the common outcomes; tests can also build Outcome
// values directly when they need a specific Reason.
func Done() Outcome {
	return Outcome{Code: s3_lifecyclewire.LifecycleDeleteOutcomeDone}
}
func NoopResolved(reason string) Outcome {
	return Outcome{Code: s3_lifecyclewire.LifecycleDeleteOutcomeNoopResolved, Reason: reason}
}
func RetryLater(reason string) Outcome {
	return Outcome{Code: s3_lifecyclewire.LifecycleDeleteOutcomeRetryLater, Reason: reason}
}
func Blocked(reason string) Outcome {
	return Outcome{Code: s3_lifecyclewire.LifecycleDeleteOutcomeBlocked, Reason: reason}
}
func SkippedObjectLock(reason string) Outcome {
	return Outcome{Code: s3_lifecyclewire.LifecycleDeleteOutcomeSkippedObjectLock, Reason: reason}
}

// FakeLifecycleServer implements s3_lifecyclewire.LifecycleDeleter. It
// returns per-key queued outcomes (FIFO) and falls back to Default when a
// key has no queued entry. Every received request is recorded; tests assert
// against Recorded() in arrival order.
//
// A non-nil Err short-circuits everything — LifecycleDelete returns (zero,
// Err) immediately, before the per-key lookup or the request recording. Use
// it to simulate transport failures.
//
// All methods are safe for concurrent use. Outcomes/Default may be set at
// construction or via Queue/SetDefault between calls; mid-call mutation is
// supported but ordering across that boundary is undefined.
type FakeLifecycleServer struct {
	mu       sync.Mutex
	queues   map[requestKey][]Outcome
	def      Outcome
	err      error
	received []s3_lifecyclewire.LifecycleDeleteRequestInput
}

// Compile-time check.
var _ s3_lifecyclewire.LifecycleDeleter = (*FakeLifecycleServer)(nil)

// requestKey is the map key for queues. A struct rather than a delimited
// string so bucket/object/versionId values containing "/" or "@" can't
// collide.
type requestKey struct {
	bucket, objectPath, versionId string
}

// NewFakeLifecycleServer returns a server whose Default outcome is DONE.
// Most tests want a different default; call SetDefault to change it.
func NewFakeLifecycleServer() *FakeLifecycleServer {
	return &FakeLifecycleServer{
		queues: map[requestKey][]Outcome{},
		def:    Done(),
	}
}

// Queue appends an outcome to the FIFO for (bucket, objectPath, versionId).
// Subsequent calls matching the same key return outcomes in the order they
// were queued; once the queue is drained, Default applies.
func (f *FakeLifecycleServer) Queue(bucket, objectPath, versionId string, outcome Outcome) {
	f.mu.Lock()
	defer f.mu.Unlock()
	k := key(bucket, objectPath, versionId)
	f.queues[k] = append(f.queues[k], outcome)
}

// SetDefault sets the outcome returned when no per-key queue entry remains.
func (f *FakeLifecycleServer) SetDefault(o Outcome) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.def = o
}

// SetError makes LifecycleDelete return (zero, err) on every call until
// SetError(nil) clears it. The request is not recorded while Err is set —
// transport-error tests should rely on the worker's own bookkeeping.
func (f *FakeLifecycleServer) SetError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.err = err
}

// Recorded returns a snapshot of every request the server has received
// (excluding calls that returned a transport error), in arrival order.
// LifecycleDeleteRequestInput is a value type, so the returned slice is
// independent of the fake's internal state.
func (f *FakeLifecycleServer) Recorded() []s3_lifecyclewire.LifecycleDeleteRequestInput {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]s3_lifecyclewire.LifecycleDeleteRequestInput, len(f.received))
	copy(out, f.received)
	return out
}

// LifecycleDelete is the ZAP service handler. It honors Err first, then
// dequeues the per-key outcome, falling back to Default. req is the
// zero-copy wire view; fields read out of it are copied into the recorded
// LifecycleDeleteRequestInput so the snapshot stays valid after the buffer
// is reused.
func (f *FakeLifecycleServer) LifecycleDelete(req s3_lifecyclewire.LifecycleDeleteRequest) (s3_lifecyclewire.LifecycleDeleteResult, error) {
	f.mu.Lock()
	if f.err != nil {
		err := f.err
		f.mu.Unlock()
		return s3_lifecyclewire.LifecycleDeleteResult{}, err
	}
	in := recordedInput(req)
	f.received = append(f.received, in)
	k := key(in.Bucket, in.ObjectPath, in.VersionID)
	q := f.queues[k]
	var out Outcome
	if len(q) > 0 {
		out = q[0]
		f.queues[k] = q[1:]
	} else {
		out = f.def
	}
	f.mu.Unlock()
	return s3_lifecyclewire.LifecycleDeleteResult{Outcome: out.Code, Reason: out.Reason}, nil
}

// recordedInput copies the wire request view into an owned
// LifecycleDeleteRequestInput so it survives buffer reuse.
func recordedInput(req s3_lifecyclewire.LifecycleDeleteRequest) s3_lifecyclewire.LifecycleDeleteRequestInput {
	in := s3_lifecyclewire.LifecycleDeleteRequestInput{
		Bucket:               req.Bucket(),
		ObjectPath:           req.ObjectPath(),
		VersionID:            req.VersionID(),
		RuleHash:             append([]byte(nil), req.RuleHash()...),
		ActionKind:           req.ActionKind(),
		StreamShard:          req.StreamShard(),
		StreamDelaySeconds:   req.StreamDelaySeconds(),
		StreamPositionTsNs:   req.StreamPositionTsNs(),
		StreamPositionOffset: req.StreamPositionOffset(),
		EngineSnapshotID:     req.EngineSnapshotID(),
	}
	if req.HasExpectedIdentity() {
		id := req.ExpectedIdentity()
		in.ExpectedIdentity = &s3_lifecyclewire.EntryIdentityInput{
			MtimeNs:      id.MtimeNs(),
			Size:         id.Size(),
			HeadFid:      id.HeadFid(),
			ExtendedHash: append([]byte(nil), id.ExtendedHash()...),
		}
	}
	return in
}

func key(bucket, objectPath, versionId string) requestKey {
	return requestKey{bucket: bucket, objectPath: objectPath, versionId: versionId}
}
