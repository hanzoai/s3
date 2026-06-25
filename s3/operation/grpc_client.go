package operation

import (
	"context"

	"google.golang.org/grpc"

	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/master_pb"
	"github.com/hanzoai/s3/s3/pb/volume_server_pb"
)

func WithVolumeServerClient(streamingMode bool, volumeServer pb.ServerAddress, grpcDialOption pb.DialOption, fn func(volume_server_pb.VolumeServerClient) error) error {

	return pb.WithGrpcClient(context.Background(), streamingMode, 0, func(grpcConnection *grpc.ClientConn) error {
		client := volume_server_pb.NewVolumeServerClient(grpcConnection)
		return fn(client)
	}, volumeServer.ToGrpcAddress(), false, grpcDialOption)

}

// WithMasterServerClient threads the caller's per-request context into the
// connection-invalidation decision, so a Canceled/DeadlineExceeded from the
// caller's own timeout does not invalidate the shared cached master connection.
// Pass context.Background() when there is no per-request deadline to honor.
func WithMasterServerClient(ctx context.Context, streamingMode bool, masterServer pb.ServerAddress, grpcDialOption pb.DialOption, fn func(masterClient master_pb.HanzoClient) error) error {
	return pb.WithMasterClient(ctx, streamingMode, masterServer, grpcDialOption, false, fn)
}
