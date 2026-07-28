package s3api

import (
	"testing"

	"github.com/hanzoai/s3/s3/s3api/s3err"
)

// fakeObjectWriteLock reports a fixed held/not-held state, which is the only
// thing withObjectWriteLock is allowed to act on.
type fakeObjectWriteLock struct {
	held     bool
	released bool
}

func (f *fakeObjectWriteLock) IsLocked() bool { return f.held }
func (f *fakeObjectWriteLock) StopShortLivedLock() error {
	f.released = true
	return nil
}

// A precondition is a compare-and-swap: it is only meaningful if nothing can
// write between the check and the commit. When the lock that would provide that
// guarantee is not held, writing anyway turns a conditional PUT into an
// unconditional one — silently. Two racers both pass the check, both commit, and
// the store has admitted two winners for one CAS.
//
// So the write must NOT happen. This is the regression that made a caller's
// atomicity probe report "update-race admitted 2 writers (want 1)".
func TestWithObjectWriteLock_PreconditionRefusesWhenLockNotHeld(t *testing.T) {
	for _, tc := range []struct {
		name string
		lock func(bucket, object string) objectWriteLock
	}{
		{"no lock configured", nil},
		{"lock constructor returns nil", func(string, string) objectWriteLock { return nil }},
		{"lock object exists but is NOT held", func(string, string) objectWriteLock {
			return &fakeObjectWriteLock{held: false}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s3a := &S3ApiServer{newObjectWriteLock: tc.lock}

			preconditionRan, wrote := false, false
			code := s3a.withObjectWriteLock("b", "o", true,
				func() s3err.ErrorCode { preconditionRan = true; return s3err.ErrNone },
				func() s3err.ErrorCode { wrote = true; return s3err.ErrNone },
			)

			if wrote {
				t.Fatal("wrote with a precondition and no held lock — the conditional PUT is not atomic")
			}
			if code == s3err.ErrNone {
				t.Fatalf("returned ErrNone; a refusal must be visible to the caller, got %v", code)
			}
			_ = preconditionRan // may or may not run; the load-bearing fact is that fn did not
		})
	}
}

// The held case still does the work, in order: precondition first, then write.
func TestWithObjectWriteLock_HeldLockRunsPreconditionThenWrite(t *testing.T) {
	lock := &fakeObjectWriteLock{held: true}
	s3a := &S3ApiServer{newObjectWriteLock: func(string, string) objectWriteLock { return lock }}

	var order []string
	code := s3a.withObjectWriteLock("b", "o", true,
		func() s3err.ErrorCode { order = append(order, "precondition"); return s3err.ErrNone },
		func() s3err.ErrorCode { order = append(order, "write"); return s3err.ErrNone },
	)
	if code != s3err.ErrNone {
		t.Fatalf("code = %v, want ErrNone", code)
	}
	if len(order) != 2 || order[0] != "precondition" || order[1] != "write" {
		t.Fatalf("order = %v, want [precondition write]", order)
	}
	if !lock.released {
		t.Error("lock was not released")
	}
}

// A failing precondition must still stop the write — the 412 path.
func TestWithObjectWriteLock_FailedPreconditionBlocksWrite(t *testing.T) {
	s3a := &S3ApiServer{newObjectWriteLock: func(string, string) objectWriteLock {
		return &fakeObjectWriteLock{held: true}
	}}
	wrote := false
	code := s3a.withObjectWriteLock("b", "o", true,
		func() s3err.ErrorCode { return s3err.ErrPreconditionFailed },
		func() s3err.ErrorCode { wrote = true; return s3err.ErrNone },
	)
	if wrote {
		t.Fatal("wrote despite a failed precondition")
	}
	if code != s3err.ErrPreconditionFailed {
		t.Fatalf("code = %v, want ErrPreconditionFailed", code)
	}
}

// The regression this signature exists to prevent, and it is not hypothetical —
// it was nearly shipped. The real caller builds preconditionFn UNCONDITIONALLY and
// lets it no-op when the request carries no conditional headers, so inferring
// "is this conditional?" from `preconditionFn != nil` reads TRUE on every single
// PUT. That would demand a held lock for ordinary writes and fail the whole object
// path closed. A non-nil precondition with needsSerialization=false must still
// write, lock or no lock.
func TestWithObjectWriteLock_NonNilPreconditionIsNotItselfAConditionalRequest(t *testing.T) {
	s3a := &S3ApiServer{newObjectWriteLock: func(string, string) objectWriteLock {
		return &fakeObjectWriteLock{held: false} // no lock available
	}}
	wrote := false
	code := s3a.withObjectWriteLock("b", "o", false,
		func() s3err.ErrorCode { return s3err.ErrNone }, // always-present, no-op closure
		func() s3err.ErrorCode { wrote = true; return s3err.ErrNone },
	)
	if !wrote || code != s3err.ErrNone {
		t.Fatalf("plain PUT blocked (wrote=%v code=%v) — a no-op precondition closure must not make a write conditional", wrote, code)
	}
}

// UNCONDITIONAL writes must be unaffected. There is nothing for them to be atomic
// with — last-writer-wins is the defined S3 semantic — so requiring a lock would
// serialize the whole object path and fail writes for no benefit. This is the
// blast-radius guard on the fix above.
func TestWithObjectWriteLock_UnconditionalWriteNeedsNoLock(t *testing.T) {
	for _, tc := range []struct {
		name string
		lock func(bucket, object string) objectWriteLock
	}{
		{"no lock configured", nil},
		{"lock constructor returns nil", func(string, string) objectWriteLock { return nil }},
		{"lock not held", func(string, string) objectWriteLock { return &fakeObjectWriteLock{held: false} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s3a := &S3ApiServer{newObjectWriteLock: tc.lock}
			wrote := false
			code := s3a.withObjectWriteLock("b", "o", false, nil,
				func() s3err.ErrorCode { wrote = true; return s3err.ErrNone })
			if !wrote || code != s3err.ErrNone {
				t.Fatalf("unconditional write blocked (wrote=%v code=%v) — plain PUTs must not need the lock", wrote, code)
			}
		})
	}
}
