package s3server

import (
	"context"
	"net"
	"testing"

	"github.com/hanzoai/s3/s3/security"

	"google.golang.org/grpc/peer"
)

// TestZapAdminAuthBypassesWhitelist is a RED-team proof that the ZAP transport
// path silently bypasses the volume server's admin IP whitelist.
//
// THREAT MODEL: an operator runs the volume server with a configured write
// whitelist (`-whiteList=10.0.0.1` or `guard.white_list`) and WITHOUT
// grpc.volume.cert/.key — i.e. plaintext ZAP (command/volume.go:591, the
// documented "loopback / dev" fallback). They reasonably expect destructive
// admin RPCs (DeleteCollection, AllocateVolume, VolumeDelete, EC destroys) to
// be restricted to whitelisted IPs, exactly as the legacy gRPC server enforced.
//
// REALITY: the ZAP backend (svc/volume/server.go:36) invokes every admin
// handler with context.Background(). checkGrpcAdminAuth sees peer.FromContext
// return ok=false and returns nil (allow) BEFORE consulting the whitelist
// (volume_grpc_admin.go:46-55). The whitelist is never checked on the ZAP
// path. Over plaintext ZAP there is no transport mTLS to compensate, so ANY
// unauthenticated TCP caller to the gRPC port can run destructive ops.
//
// This test reproduces the exact context the ZAP backend uses and asserts the
// SECURITY-CORRECT behavior (deny a non-whitelisted caller). It FAILS against
// the released v2026.3.0 code, proving the hole.
func TestZapAdminAuthBypassesWhitelist(t *testing.T) {
	// Guard with a non-empty whitelist that does NOT include the attacker.
	// isWriteActive=true, isEmptyWhiteList=false => the legacy gRPC path would
	// DENY any host not in {10.0.0.1}.
	g := security.NewGuard([]string{"10.0.0.1"}, "", 0, "", 0)
	if g.IsWhiteListed("203.0.113.7") {
		t.Fatal("precondition: attacker IP must not be whitelisted")
	}

	vs := &VolumeServer{guard: g}

	// This is byte-for-byte the context the ZAP backend hands to the admin
	// handlers: serverBackend{ctx: context.Background()} -> vs.DeleteCollection(ctx, ...).
	zapCtx := context.Background()

	if _, ok := peer.FromContext(zapCtx); ok {
		t.Fatal("precondition: ZAP backend ctx must carry no gRPC peer")
	}

	err := vs.checkGrpcAdminAuth(zapCtx)
	if err == nil {
		t.Fatalf("VULN CONFIRMED: checkGrpcAdminAuth ALLOWED an unauthenticated caller " +
			"over the ZAP path while a non-empty write whitelist was configured. " +
			"Over plaintext ZAP (no grpc.volume.cert) this lets any TCP client run " +
			"DeleteCollection/AllocateVolume/VolumeDelete/EC-destroy with zero auth.")
	}
}

// TestGrpcAdminAuthDeniesNonWhitelistedPeer is the control: the SAME guard, the
// SAME caller, but reached over a real gRPC peer (the legacy transport) is
// correctly DENIED. The contrast isolates the regression to the ZAP no-peer
// branch — the policy itself is intact, only the ZAP path skips it.
func TestGrpcAdminAuthDeniesNonWhitelistedPeer(t *testing.T) {
	g := security.NewGuard([]string{"10.0.0.1"}, "", 0, "", 0)
	vs := &VolumeServer{guard: g}

	grpcCtx := peer.NewContext(context.Background(), &peer.Peer{Addr: &net.TCPAddr{IP: net.ParseIP("203.0.113.7"), Port: 54321}})

	if err := vs.checkGrpcAdminAuth(grpcCtx); err == nil {
		t.Fatal("control failed: gRPC peer path should deny a non-whitelisted caller")
	}
}
