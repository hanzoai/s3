package dailyrun

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/hanzoai/s3/s3/s3api/s3lifecycle"
	"github.com/hanzoai/s3/s3/s3api/s3lifecycle/router"
	s3_lifecyclewire "github.com/hanzoai/s3/s3/wire/s3_lifecycle"
)

// transportRetryAttempts is the in-run retry budget for transport
// errors (RPC dial failures, stream resets). Server-side outcomes
// (RETRY_LATER, BLOCKED) are NOT retried in-run — the server already
// classified them as "wait until next run." The retry exists only to
// paper over a single flake from the network.
const transportRetryAttempts = 3

// transportRetryInitial is the first sleep after a failed attempt.
// Doubled per attempt up to transportRetryMax. Phase 2 keeps these
// small because the daily run's own MaxRuntime is the larger cap.
const (
	transportRetryInitial = 200 * time.Millisecond
	transportRetryMax     = 5 * time.Second
)

// dispatchWithRetry sends one LifecycleDelete request, retrying on
// transport errors up to transportRetryAttempts times with exponential
// backoff. Returns the server's outcome (a LifecycleDeleteOutcome*
// constant) on success, or an error if all retries exhausted. Server
// outcomes (RETRY_LATER / BLOCKED) bypass the retry — the daily run's
// halt-on-failure semantics handles them.
func dispatchWithRetry(ctx context.Context, client LifecycleClient, m router.Match) (uint32, error) {
	in := buildDeleteRequest(m)
	backoff := transportRetryInitial
	var lastErr error
	for attempt := 1; attempt <= transportRetryAttempts; attempt++ {
		res, err := client.LifecycleDelete(ctx, in)
		if err == nil {
			return res.Outcome, nil
		}
		// Context cancellation is shutdown, not a transport flake.
		// Surface it immediately so the caller halts and tomorrow
		// resumes from the same cursor.
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return s3_lifecyclewire.LifecycleDeleteOutcomeUnspecified, err
		}
		lastErr = err
		if attempt == transportRetryAttempts {
			break
		}
		select {
		case <-ctx.Done():
			return s3_lifecyclewire.LifecycleDeleteOutcomeUnspecified, ctx.Err()
		case <-time.After(jitter(backoff)):
		}
		backoff *= 2
		if backoff > transportRetryMax {
			backoff = transportRetryMax
		}
	}
	return s3_lifecyclewire.LifecycleDeleteOutcomeUnspecified, lastErr
}

// jitter returns a duration in the range [d/2, d) using equal jitter.
// Prevents thundering herds when many daily-run workers retry simultaneously.
func jitter(d time.Duration) time.Duration {
	if d <= 0 {
		return 0
	}
	half := d / 2
	if half <= 0 {
		return d
	}
	return half + time.Duration(rand.Int63n(int64(half)))
}

// buildDeleteRequest constructs the LifecycleDelete RPC payload for a
// router Match. Mirrors dispatcher.dispatchOne's request shape — both
// targets the same server-side handler and the wire encoding must
// match exactly. Duplicated rather than shared because Phase 5
// deletes the dispatcher's call site and the cross-package dependency
// would just create a temporary import.
func buildDeleteRequest(m router.Match) s3_lifecyclewire.LifecycleDeleteRequestInput {
	rh := m.Key.RuleHash
	return s3_lifecyclewire.LifecycleDeleteRequestInput{
		Bucket:           m.Bucket,
		ObjectPath:       m.ObjectKey,
		VersionID:        m.VersionID,
		RuleHash:         rh[:],
		ActionKind:       wireActionKind(m.Key.ActionKind),
		ExpectedIdentity: wireIdentity(m.Identity),
	}
}

// wireActionKind maps the in-package ActionKind enum to the wire
// ActionKind* uint32 constant.
func wireActionKind(k s3lifecycle.ActionKind) uint32 {
	switch k {
	case s3lifecycle.ActionKindExpirationDays:
		return s3_lifecyclewire.ActionKindExpirationDays
	case s3lifecycle.ActionKindExpirationDate:
		return s3_lifecyclewire.ActionKindExpirationDate
	case s3lifecycle.ActionKindNoncurrentDays:
		return s3_lifecyclewire.ActionKindNoncurrentDays
	case s3lifecycle.ActionKindNewerNoncurrent:
		return s3_lifecyclewire.ActionKindNewerNoncurrent
	case s3lifecycle.ActionKindAbortMPU:
		return s3_lifecyclewire.ActionKindAbortMpu
	case s3lifecycle.ActionKindExpiredDeleteMarker:
		return s3_lifecyclewire.ActionKindExpiredDeleteMarker
	}
	return s3_lifecyclewire.ActionKindUnspecified
}

// wireIdentity maps the router CAS witness to a wire EntryIdentityInput;
// nil in yields nil so the server treats it as "no witness, bootstrap".
func wireIdentity(id *router.EntryIdentity) *s3_lifecyclewire.EntryIdentityInput {
	if id == nil {
		return nil
	}
	return &s3_lifecyclewire.EntryIdentityInput{
		MtimeNs:      id.MtimeNs,
		Size:         id.Size,
		HeadFid:      id.HeadFid,
		ExtendedHash: id.ExtendedHash,
	}
}

// outcomeLabel renders a LifecycleDeleteOutcome* constant as the stable
// Prometheus label string. The values match the proto enum names so
// existing dashboards/alerts keyed on these labels keep working across
// the gRPC->ZAP cut.
func outcomeLabel(o uint32) string {
	switch o {
	case s3_lifecyclewire.LifecycleDeleteOutcomeDone:
		return "DONE"
	case s3_lifecyclewire.LifecycleDeleteOutcomeNoopResolved:
		return "NOOP_RESOLVED"
	case s3_lifecyclewire.LifecycleDeleteOutcomeSkippedObjectLock:
		return "SKIPPED_OBJECT_LOCK"
	case s3_lifecyclewire.LifecycleDeleteOutcomeRetryLater:
		return "RETRY_LATER"
	case s3_lifecyclewire.LifecycleDeleteOutcomeBlocked:
		return "BLOCKED"
	default:
		return "LIFECYCLE_DELETE_OUTCOME_UNSPECIFIED"
	}
}
