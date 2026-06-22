//go:build netbsd || plan9

package stats

import "github.com/hanzoai/s3/s3/pb/volume_server_pb"

func fillInDiskStatus(status *volume_server_pb.DiskStatus) error {
	return nil
}
