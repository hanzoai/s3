package mount

// filer_client_adapter.go is the FUSE mount's seam onto the native ZAP
// transport. The adapter itself — a filer_pb.HanzoFilerClient that routes every
// unary and streaming RPC over a github.com/zap-proto/go transport.Conn — lives
// in package filerzap (filerzap.NewZapFilerClient), re-exported by package pb
// (pb.NewZapFilerClient), so the contract is bridged in exactly ONE place and
// shared with the package-level filer helpers (pb.WithFilerClient and the
// metadata subscribers). mount keeps the short newFilerClientAdapter name its
// dial sites use; it is a thin alias onto the canonical pb constructor.

import (
	"github.com/hanzoai/s3/s3/pb"
	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"

	"github.com/zap-proto/go/transport"
)

// newFilerClientAdapter wraps an established transport.Conn as a
// filer_pb.HanzoFilerClient backed by ZAP. The caller owns conn's lifecycle.
func newFilerClientAdapter(conn transport.Conn) filer_pb.HanzoFilerClient {
	return pb.NewZapFilerClient(conn)
}
