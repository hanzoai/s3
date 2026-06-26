package s3server

import (
	"context"
	"testing"

	"github.com/hanzoai/s3/s3/security"
)

// TestCheckGrpcAdminAuth_FailsClosedOverZap pins the v2026.3.1 hotfix contract for
// the volume admin auth boundary under the native ZAP transport.
//
// Every volume RPC now arrives over ZAP, whose per-call context carries no peer
// address (unlike gRPC's peer.FromContext, gone with the volume gRPC server), so
// there is no source IP to match against the guard's write whitelist. The boundary
// FAILS CLOSED: when a write whitelist is configured but the caller cannot be
// IP-verified over ZAP, the admin RPC is DENIED (authenticated ZAP admin authz is
// the transport's PQ-mTLS allowed-CN gate, not IP). Only an empty whitelist
// (allow-all / dev) or no guard admits. The prior code returned nil here BEFORE
// consulting the whitelist — a categorical authz bypass on every ZAP admin RPC.
func TestCheckGrpcAdminAuth_FailsClosedOverZap(t *testing.T) {
	// A write whitelist is configured and the ZAP call has no peer to verify:
	// must DENY (this is the closed bypass).
	vs := &VolumeServer{
		guard: security.NewGuard([]string{"10.0.0.1"}, "", 0, "", 0),
	}
	if err := vs.checkGrpcAdminAuth(context.Background()); err == nil {
		t.Fatal("checkGrpcAdminAuth must DENY a ZAP admin RPC when a write whitelist is " +
			"configured and no peer IP can be verified (the v2026.3.1 bypass fix); got nil")
	}

	// Empty whitelist (allow-all / dev): IsWhiteListed("")==true, so admit.
	vsAllowAll := &VolumeServer{
		guard: security.NewGuard(nil, "", 0, "", 0),
	}
	if err := vsAllowAll.checkGrpcAdminAuth(context.Background()); err != nil {
		t.Fatalf("checkGrpcAdminAuth with an empty whitelist (allow-all) must admit, got %v", err)
	}

	// No guard at all (local / dev path): admit.
	vsNoGuard := &VolumeServer{}
	if err := vsNoGuard.checkGrpcAdminAuth(context.Background()); err != nil {
		t.Fatalf("checkGrpcAdminAuth with no guard must admit, got %v", err)
	}
}
