package s3server

import (
	"context"
	"testing"

	"github.com/hanzoai/s3/s3/security"
)

// TestCheckGrpcAdminAuth_AllowsZapTransport verifies the volume admin auth
// boundary under the native ZAP transport.
//
// Every volume RPC now arrives over ZAP, whose per-call context carries no peer
// address (unlike gRPC's peer.FromContext). On that transport the caller is
// authenticated at the TRANSPORT by PQ-mTLS (pb.ServerTLSConfig: required+verified
// client cert, allowed-CN gate, fail-closed without a CA), so the gRPC-era IP
// whitelist no longer gates admin RPCs — there is no source IP to match. This test
// pins that contract: with a guard configured but no gRPC peer in context (the
// ZAP shape), checkGrpcAdminAuth admits the call and defers authorization to the
// transport.
func TestCheckGrpcAdminAuth_AllowsZapTransport(t *testing.T) {
	vs := &VolumeServer{
		guard: security.NewGuard([]string{"10.0.0.1"}, "", 0, "", 0),
	}

	// A plain context with no gRPC peer is exactly what the ZAP transport hands
	// every RPC handler.
	if err := vs.checkGrpcAdminAuth(context.Background()); err != nil {
		t.Fatalf("checkGrpcAdminAuth over ZAP transport must admit (mTLS is the boundary), got %v", err)
	}

	// And with no guard at all (local/dev path) it must also admit.
	vsNoGuard := &VolumeServer{}
	if err := vsNoGuard.checkGrpcAdminAuth(context.Background()); err != nil {
		t.Fatalf("checkGrpcAdminAuth with no guard must admit, got %v", err)
	}
}
