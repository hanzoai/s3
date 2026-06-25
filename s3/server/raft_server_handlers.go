package s3server

import (
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/storage/needle"
)

// ClusterStatusResult is the JSON body of GET /cluster/status. The handler now
// lives on ConsensusServer (consensus_server_master.go); the shape is unchanged
// from the Raft era.
type ClusterStatusResult struct {
	IsLeader    bool             `json:"IsLeader,omitempty"`
	Leader      pb.ServerAddress `json:"Leader,omitempty"`
	Peers       []string         `json:"Peers,omitempty"`
	MaxVolumeId needle.VolumeId  `json:"MaxVolumeId,omitempty"`
}
