package s3server

import (
	"errors"

	"github.com/hanzoai/s3/s3/pb/volume_server_pb"
)

// Query is the experimental SELECT-over-needles RPC. It is not implemented; the
// server answers with an Unimplemented-tagged error (the ZAP code-name string
// convention) instead of streaming results. This replaces the behavior the
// retired grpc UnimplementedVolumeServerServer embed supplied.
func (vs *VolumeServer) Query(req *volume_server_pb.QueryRequest, stream volume_server_pb.VolumeServer_QueryServer) error {
	return errors.New("Unimplemented: method Query not implemented")
}
