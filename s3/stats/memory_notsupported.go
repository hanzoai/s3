//go:build !linux

package stats

import "github.com/hanzoai/s3/s3/pb/volume_server_pb"

func fillInMemStatus(status *volume_server_pb.MemStatus) {
	return
}
