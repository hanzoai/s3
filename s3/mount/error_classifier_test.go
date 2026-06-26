package mount

import (
	"errors"
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/hanzoai/go-fuse/v2/fuse"
)

// A Canceled error from a closing filer connection must surface as ETIMEDOUT (a
// retryable hint for FUSE callers) rather than EIO, even when wrapped by a
// caller. wfs_save.go wraps with %w; the ZAP classifier matches the "Canceled"
// code name the filer tags onto the error string.
func TestGrpcErrorToFuseStatusUnwrapsCanceledThroughFmtErrorf(t *testing.T) {
	grpcErr := fmt.Errorf("Canceled: grpc: the client connection is closing")

	wrapped := fmt.Errorf("UpdateEntry dir /some/path: %w", grpcErr)

	got := grpcErrorToFuseStatus(wrapped)
	want := fuse.Status(syscall.ETIMEDOUT)
	if got != want {
		t.Fatalf("grpcErrorToFuseStatus(canceled wrapped with %%w) = %v, want %v", got, want)
	}
}

// The filer speaks ZAP now: the error is a string carrying the "Canceled" code
// name, so classification is by code-name match and no longer depends on the
// wrap verb. A Canceled error wrapped with %v (which flattens the old gRPC
// status object) classifies identically to %w — both ETIMEDOUT. This guards the
// verb-agnostic ZAP contract so a future change can't reintroduce status-object
// unwrapping that would make %v drop to EIO again.
func TestGrpcErrorToFuseStatusClassifiesCanceledThroughPercentV(t *testing.T) {
	grpcErr := fmt.Errorf("Canceled: grpc: the client connection is closing")

	wrapped := fmt.Errorf("UpdateEntry dir /some/path: %v", grpcErr)

	got := grpcErrorToFuseStatus(wrapped)
	want := fuse.Status(syscall.ETIMEDOUT)
	if got != want {
		t.Fatalf("grpcErrorToFuseStatus(canceled wrapped with %%v) = %v, want %v (ZAP classifies by code name, verb-agnostic)", got, want)
	}
}

func TestIsRetryableFilerError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"canceled", fmt.Errorf("Canceled: grpc: the client connection is closing"), true},
		{"unavailable", fmt.Errorf("Unavailable: connection refused"), true},
		{"deadline_exceeded", fmt.Errorf("DeadlineExceeded: deadline exceeded"), true},
		{"resource_exhausted", fmt.Errorf("ResourceExhausted: too many concurrent requests"), true},
		{"internal", fmt.Errorf("Internal: server crashed"), true},
		{"not_found", fmt.Errorf("NotFound: entry missing"), false},
		{"already_exists", fmt.Errorf("AlreadyExists: duplicate"), false},
		{"invalid_argument", fmt.Errorf("InvalidArgument: bad request"), false},
		{"permission_denied", fmt.Errorf("PermissionDenied: no access"), false},
		{"unauthenticated", fmt.Errorf("Unauthenticated: missing creds"), false},
		{"failed_precondition", fmt.Errorf("FailedPrecondition: not empty"), false},
		{"plain_error_retries", errors.New("random network glitch"), true},
		{"wrapped_canceled_still_retries", fmt.Errorf("ctx: %w", fmt.Errorf("Canceled: closing")), true},
		{"wrapped_not_found_still_skipped", fmt.Errorf("ctx: %w", fmt.Errorf("NotFound: gone")), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isRetryableFilerError(tc.err); got != tc.want {
				t.Fatalf("isRetryableFilerError(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

// retryMetadataFlushIf must short-circuit on non-retryable errors so that
// synchronous FUSE ops (chmod/utimes/xattr) don't hang for ~7s on ENOENT/
// EACCES/EINVAL.
func TestRetryMetadataFlushIfShortCircuitsOnPermanentError(t *testing.T) {
	originalSleep := metadataFlushSleep
	t.Cleanup(func() {
		metadataFlushSleep = originalSleep
	})
	metadataFlushSleep = func(_ time.Duration) {
		t.Fatal("sleep should not be called when shouldRetry returns false")
	}

	attempts := 0
	permanent := fmt.Errorf("NotFound: entry missing")
	err := retryMetadataFlushIf(func() error {
		attempts++
		return permanent
	}, isRetryableFilerError, nil)

	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1 (permanent error should short-circuit)", attempts)
	}
	if !errors.Is(err, permanent) {
		t.Fatalf("err = %v, want permanent sentinel", err)
	}
}

// Transient errors must keep retrying up to the attempt cap even when a
// predicate is supplied.
func TestRetryMetadataFlushIfRetriesTransientErrors(t *testing.T) {
	originalSleep := metadataFlushSleep
	t.Cleanup(func() {
		metadataFlushSleep = originalSleep
	})
	metadataFlushSleep = func(_ time.Duration) {}

	attempts := 0
	transient := fmt.Errorf("Canceled: grpc: the client connection is closing")
	err := retryMetadataFlushIf(func() error {
		attempts++
		return transient
	}, isRetryableFilerError, nil)

	if attempts != metadataFlushRetries+1 {
		t.Fatalf("attempts = %d, want %d", attempts, metadataFlushRetries+1)
	}
	if !errors.Is(err, transient) {
		t.Fatalf("err = %v, want transient sentinel", err)
	}
}
