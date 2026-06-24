package lifecycletest

import (
	"errors"
	"sync"
	"testing"

	s3_lifecyclewire "github.com/hanzoai/s3/s3/wire/s3_lifecycle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// call drives the fake's LifecycleDeleter through a real wire request view,
// the same path the ZAP transport uses at runtime.
func call(t *testing.T, f *FakeLifecycleServer, in s3_lifecyclewire.LifecycleDeleteRequestInput) (s3_lifecyclewire.LifecycleDeleteResult, error) {
	t.Helper()
	v, err := s3_lifecyclewire.WrapLifecycleDeleteRequest(s3_lifecyclewire.NewLifecycleDeleteRequest(in))
	require.NoError(t, err)
	return f.LifecycleDelete(v)
}

func TestFake_DefaultIsDoneOutOfTheBox(t *testing.T) {
	// A test that doesn't queue anything should still get a non-error
	// response so it can exercise the worker's happy path.
	f := NewFakeLifecycleServer()
	res, err := call(t, f, s3_lifecyclewire.LifecycleDeleteRequestInput{
		Bucket:     "b",
		ObjectPath: "k",
	})
	require.NoError(t, err)
	assert.Equal(t, s3_lifecyclewire.LifecycleDeleteOutcomeDone, res.Outcome)
	assert.Equal(t, "", res.Reason)
}

func TestFake_QueuedOutcomesPopFIFO(t *testing.T) {
	// Per-key queue is FIFO and one-shot per entry; after the queue
	// drains, Default kicks in.
	f := NewFakeLifecycleServer()
	f.SetDefault(NoopResolved("nothing more queued"))
	f.Queue("b", "k", "", RetryLater("first"))
	f.Queue("b", "k", "", Blocked("second"))

	got := []uint32{}
	reasons := []string{}
	for i := 0; i < 3; i++ {
		res, err := call(t, f, s3_lifecyclewire.LifecycleDeleteRequestInput{
			Bucket: "b", ObjectPath: "k",
		})
		require.NoError(t, err)
		got = append(got, res.Outcome)
		reasons = append(reasons, res.Reason)
	}
	assert.Equal(t, []uint32{
		s3_lifecyclewire.LifecycleDeleteOutcomeRetryLater,
		s3_lifecyclewire.LifecycleDeleteOutcomeBlocked,
		s3_lifecyclewire.LifecycleDeleteOutcomeNoopResolved,
	}, got)
	assert.Equal(t, []string{"first", "second", "nothing more queued"}, reasons)
}

func TestFake_QueuesIsolatedByKey(t *testing.T) {
	// Queues are partitioned by (bucket, objectPath, versionId); a queued
	// outcome for one key must not bleed into another's lookup.
	f := NewFakeLifecycleServer()
	f.Queue("b", "a", "", Blocked("a-only"))
	f.Queue("b", "b", "", RetryLater("b-only"))

	respA, err := call(t, f, s3_lifecyclewire.LifecycleDeleteRequestInput{Bucket: "b", ObjectPath: "a"})
	require.NoError(t, err)
	assert.Equal(t, s3_lifecyclewire.LifecycleDeleteOutcomeBlocked, respA.Outcome)

	respB, err := call(t, f, s3_lifecyclewire.LifecycleDeleteRequestInput{Bucket: "b", ObjectPath: "b"})
	require.NoError(t, err)
	assert.Equal(t, s3_lifecyclewire.LifecycleDeleteOutcomeRetryLater, respB.Outcome)
}

func TestFake_VersionIDPartOfKey(t *testing.T) {
	// Two requests for the same bucket/objectPath but different
	// versionIds must address different queues.
	f := NewFakeLifecycleServer()
	f.Queue("b", "k", "v1", SkippedObjectLock("v1-locked"))
	f.Queue("b", "k", "v2", Done())

	respV1, err := call(t, f, s3_lifecyclewire.LifecycleDeleteRequestInput{Bucket: "b", ObjectPath: "k", VersionID: "v1"})
	require.NoError(t, err)
	respV2, err := call(t, f, s3_lifecyclewire.LifecycleDeleteRequestInput{Bucket: "b", ObjectPath: "k", VersionID: "v2"})
	require.NoError(t, err)
	assert.Equal(t, s3_lifecyclewire.LifecycleDeleteOutcomeSkippedObjectLock, respV1.Outcome)
	assert.Equal(t, s3_lifecyclewire.LifecycleDeleteOutcomeDone, respV2.Outcome)
}

func TestFake_KeyComponentsWithDelimitersDoNotCollide(t *testing.T) {
	// String-concatenation keys would have made these two requests
	// indistinguishable. The struct key keeps them separate.
	f := NewFakeLifecycleServer()
	f.Queue("b/k", "", "", Blocked("variant-a"))
	f.Queue("b", "k", "", Done())

	respA, err := call(t, f, s3_lifecyclewire.LifecycleDeleteRequestInput{Bucket: "b/k", ObjectPath: ""})
	require.NoError(t, err)
	respB, err := call(t, f, s3_lifecyclewire.LifecycleDeleteRequestInput{Bucket: "b", ObjectPath: "k"})
	require.NoError(t, err)
	assert.Equal(t, s3_lifecyclewire.LifecycleDeleteOutcomeBlocked, respA.Outcome)
	assert.Equal(t, s3_lifecyclewire.LifecycleDeleteOutcomeDone, respB.Outcome)
}

func TestFake_ErrShortCircuitsBeforeRecording(t *testing.T) {
	// Err makes LifecycleDelete return (zero, err) without recording the
	// request — transport-error tests rely on the worker's own
	// bookkeeping, not the fake's.
	f := NewFakeLifecycleServer()
	transportErr := errors.New("connection refused")
	f.SetError(transportErr)

	_, err := call(t, f, s3_lifecyclewire.LifecycleDeleteRequestInput{Bucket: "b", ObjectPath: "k"})
	assert.ErrorIs(t, err, transportErr)
	assert.Empty(t, f.Recorded(), "transport-error calls must not be recorded")

	// Clearing the error returns the server to normal behavior.
	f.SetError(nil)
	_, err2 := call(t, f, s3_lifecyclewire.LifecycleDeleteRequestInput{Bucket: "b", ObjectPath: "k"})
	require.NoError(t, err2)
	assert.Len(t, f.Recorded(), 1)
}

func TestFake_RecordsRequestsInOrder(t *testing.T) {
	// Recorded() preserves arrival order so tests can assert that
	// dispatch happened in the expected sequence.
	f := NewFakeLifecycleServer()
	for _, k := range []string{"k1", "k2", "k3"} {
		_, err := call(t, f, s3_lifecyclewire.LifecycleDeleteRequestInput{
			Bucket: "b", ObjectPath: k,
		})
		require.NoError(t, err)
	}
	rec := f.Recorded()
	require.Len(t, rec, 3)
	assert.Equal(t, "k1", rec[0].ObjectPath)
	assert.Equal(t, "k2", rec[1].ObjectPath)
	assert.Equal(t, "k3", rec[2].ObjectPath)
}

func TestFake_RecordedIsSnapshot(t *testing.T) {
	// Mutating the slice the caller got back must not bleed into the
	// fake's internal state — otherwise a flaky test could corrupt
	// bookkeeping for later assertions. LifecycleDeleteRequestInput is a
	// value type, so a copied slice is fully decoupled.
	f := NewFakeLifecycleServer()
	_, err := call(t, f, s3_lifecyclewire.LifecycleDeleteRequestInput{Bucket: "b", ObjectPath: "k"})
	require.NoError(t, err)

	snap := f.Recorded()
	require.Len(t, snap, 1)
	snap[0].ObjectPath = "mutated-by-caller"
	again := f.Recorded()
	require.Len(t, again, 1)
	assert.Equal(t, "k", again[0].ObjectPath, "internal record must survive caller-side mutation of the snapshot")
}

func TestFake_ConcurrentCallsSerializeWithoutDeadlock(t *testing.T) {
	// The dispatcher fans dispatch across many goroutines; the fake
	// must not livelock or drop records under concurrent load. Capture
	// each goroutine's err so a regression in concurrent paths surfaces
	// instead of being masked by length-only assertions.
	f := NewFakeLifecycleServer()
	const N = 64
	var wg sync.WaitGroup
	errCh := make(chan error, N)
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			_, err := call(t, f, s3_lifecyclewire.LifecycleDeleteRequestInput{Bucket: "b", ObjectPath: "k"})
			errCh <- err
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		require.NoError(t, err)
	}
	assert.Len(t, f.Recorded(), N)
}
