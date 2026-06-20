//go:build !linux

package stats

import "github.com/hanzoai/s3/weed/pb/volume_server_pb"

func fillInMemStatus(status *volume_server_pb.MemStatus) {
	return
}
