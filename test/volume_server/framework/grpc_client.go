package framework

import (
	"testing"

	"github.com/hanzoai/s3/s3/pb/volume_server_pb"
	"github.com/hanzoai/s3/s3/volumezap"

	"github.com/zap-proto/go/transport"
)

// DialVolumeServer dials the volume server over the native ZAP transport (gRPC is
// gone) and returns the connection plus a volume_server_pb.VolumeServerClient
// backed by it (volumezap.New). The caller closes the returned conn.
func DialVolumeServer(t testing.TB, address string) (transport.Conn, volume_server_pb.VolumeServerClient) {
	t.Helper()

	conn, err := transport.Dial("tcp", address)
	if err != nil {
		t.Fatalf("dial volume ZAP %s: %v", address, err)
	}

	return conn, volumezap.New(conn, nil)
}
