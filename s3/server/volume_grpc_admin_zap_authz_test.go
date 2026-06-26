package s3server

import (
	"context"
	"testing"

	"github.com/hanzoai/s3/s3/security"
)

// TestZapAdminAuthDeniesWhenWhitelistConfigured is the regression guard for the
// v2026.3.1 hotfix that CLOSED the ZAP admin authz bypass (Red CRITICAL).
//
// THREAT MODEL: an operator runs the volume server with a configured write
// whitelist (`-whiteList=10.0.0.1` or `guard.white_list`) and WITHOUT
// grpc.volume.cert/.key — i.e. plaintext ZAP (command/volume.go, the documented
// "loopback / dev" fallback). They reasonably expect destructive admin RPCs
// (DeleteCollection, AllocateVolume, VolumeDelete, EC destroys) to be restricted,
// exactly as the legacy gRPC server enforced.
//
// REALITY of the ZAP transport: every admin RPC arrives over the native ZAP mesh,
// whose per-call context carries NO peer address (gRPC's peer.FromContext is gone
// with the volume gRPC server). The pre-hotfix code returned nil (allow) on that
// peerless path BEFORE consulting the whitelist, so any unauthenticated TCP caller
// could run destructive ops over plaintext ZAP — a categorical authz bypass.
//
// FIXED: checkGrpcAdminAuth now FAILS CLOSED — with a write whitelist configured
// it cannot honor an IP restriction it can no longer evaluate, so it DENIES rather
// than silently bypass. This test reproduces the exact context the ZAP backend
// hands the admin handlers (context.Background) and asserts the denial holds.
func TestZapAdminAuthDeniesWhenWhitelistConfigured(t *testing.T) {
	// Guard with a non-empty whitelist that does NOT include the attacker:
	// isWriteActive=true, isEmptyWhiteList=false.
	g := security.NewGuard([]string{"10.0.0.1"}, "", 0, "", 0)
	if g.IsWhiteListed("203.0.113.7") {
		t.Fatal("precondition: attacker IP must not be whitelisted")
	}

	vs := &VolumeServer{guard: g}

	// Byte-for-byte the context the ZAP backend hands the admin handlers:
	// serverBackend{ctx: context.Background()} -> vs.DeleteCollection(ctx, ...).
	// It carries no peer, exactly as the ZAP mesh does.
	zapCtx := context.Background()

	if err := vs.checkGrpcAdminAuth(zapCtx); err == nil {
		t.Fatalf("BYPASS REOPENED: checkGrpcAdminAuth ALLOWED an unauthenticated caller " +
			"over the ZAP path while a non-empty write whitelist was configured. " +
			"Over plaintext ZAP (no grpc.volume.cert) this lets any TCP client run " +
			"DeleteCollection/AllocateVolume/VolumeDelete/EC-destroy with zero auth.")
	}
}

// TestZapAdminAuthAllowsWhenUnrestricted is the complement: the fail-closed gate
// must NOT over-deny the legitimate dev / loopback path. With no guard at all, or
// with an empty whitelist (allow-all, IsWhiteListed("")==true), the admin RPCs are
// allowed — authentication on this path is the transport's PQ-mTLS allowed-CN gate
// (grpc.volume.allowed_commonNames), not a per-call IP.
func TestZapAdminAuthAllowsWhenUnrestricted(t *testing.T) {
	zapCtx := context.Background()

	// No guard configured (local / dev path).
	vsNoGuard := &VolumeServer{guard: nil}
	if err := vsNoGuard.checkGrpcAdminAuth(zapCtx); err != nil {
		t.Fatalf("nil guard must allow (local/dev path), got: %v", err)
	}

	// Guard with an empty whitelist (allow-all).
	vsOpen := &VolumeServer{guard: security.NewGuard(nil, "", 0, "", 0)}
	if err := vsOpen.checkGrpcAdminAuth(zapCtx); err != nil {
		t.Fatalf("empty whitelist must allow (allow-all), got: %v", err)
	}
}
