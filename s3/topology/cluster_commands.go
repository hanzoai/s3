package topology

import (
	"github.com/hanzoai/s3/s3/storage/needle"
)

// MaxVolumeIdCommand is the master's entire replicated state machine: a monotonic
// max-volume-id plus the cluster's topology id. It is committed through the
// consensus engine (s3server.ConsensusServer) as a JSON payload and applied, in
// total order, on every node. It is a plain serializable value — the consensus
// engine owns finalization and apply, so this type carries no engine callbacks.
type MaxVolumeIdCommand struct {
	MaxVolumeId needle.VolumeId `json:"maxVolumeId"`
	TopologyId  string          `json:"topologyId"`
}

func NewMaxVolumeIdCommand(value needle.VolumeId, topologyId string) *MaxVolumeIdCommand {
	return &MaxVolumeIdCommand{
		MaxVolumeId: value,
		TopologyId:  topologyId,
	}
}
