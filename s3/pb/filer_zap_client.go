package pb

// filer_zap_client.go re-exports the canonical ZAP-backed filer_pb.HanzoFilerClient
// for package-pb callers (WithFilerClient / WithGrpcFilerClient / the mount FUSE
// adapter). The bridge implementation lives in exactly ONE place —
// filer.NewZapFilerClient — which owns the per-RPC wire converters and the
// filerstream streaming adapters. This wrapper keeps the long-standing
// pb.NewZapFilerClient(conn) call shape (no capability token) without duplicating
// any conversion logic.

import (
	"github.com/zap-proto/go/transport"

	"github.com/hanzoai/s3/s3/svc/filer"
	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
)

// NewZapFilerClient wraps an established transport.Conn as a
// filer_pb.HanzoFilerClient that routes every call over ZAP. It delegates to
// filer.NewZapFilerClient with no capability token; the caller owns conn's
// lifecycle.
func NewZapFilerClient(conn transport.Conn) filer_pb.HanzoFilerClient {
	return filer.NewZapFilerClient(conn, nil)
}
