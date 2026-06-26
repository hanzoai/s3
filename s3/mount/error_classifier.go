package mount

import (
	"strings"
	"syscall"

	"github.com/hanzoai/go-fuse/v2/fuse"
)

// grpcErrorToFuseStatus maps a filer RPC error to a FUSE status. The filer
// speaks ZAP now, which carries the error as a string; the server tags each
// failure with its code name (e.g. "NotFound: ..."), so we classify by matching
// that name. Order matters: check the specific code names before the generic
// transport fallback.
func grpcErrorToFuseStatus(err error) fuse.Status {
	if err == nil {
		return fuse.OK
	}

	s := err.Error()
	switch {
	case has(s, "Canceled"), has(s, "DeadlineExceeded"):
		return fuse.Status(syscall.ETIMEDOUT)
	case has(s, "Unavailable"), has(s, "ResourceExhausted"):
		return fuse.Status(syscall.EAGAIN)
	case has(s, "PermissionDenied"):
		return fuse.Status(syscall.EACCES)
	case has(s, "Unauthenticated"):
		return fuse.Status(syscall.EPERM)
	case has(s, "NotFound"):
		return fuse.ENOENT
	case has(s, "AlreadyExists"):
		return fuse.Status(syscall.EEXIST)
	case has(s, "InvalidArgument"):
		return fuse.EINVAL
	case has(s, "transport"):
		return fuse.Status(syscall.EAGAIN)
	}
	return fuse.EIO
}

// isRetryableFilerError reports whether a filer RPC error looks transient
// enough to retry. It takes a conservative whitelist approach: only errors
// that clearly describe a permanent application-level failure
// (NotFound/AlreadyExists/InvalidArgument/PermissionDenied/Unauthenticated/
// FailedPrecondition) short-circuit the retry loop. Everything else —
// transport errors, Canceled/Unavailable/ResourceExhausted, or errors with no
// code name — is treated as potentially transient and retried.
func isRetryableFilerError(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	switch {
	case has(s, "NotFound"),
		has(s, "AlreadyExists"),
		has(s, "InvalidArgument"),
		has(s, "PermissionDenied"),
		has(s, "Unauthenticated"),
		has(s, "FailedPrecondition"):
		return false
	}
	return true
}

func has(s, code string) bool { return strings.Contains(s, code) }
