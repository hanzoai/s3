package s3server

import (
	"net/http"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/storage/needle"
)

// ClusterStatusResult is the JSON shape served at /cluster/status. It replaces
// the raft StatusHandler: there is no leader, so it reports the pinned writer.
type ClusterStatusResult struct {
	IsWriter    bool               `json:"IsWriter,omitempty"`
	Writer      pb.ServerAddress   `json:"Writer,omitempty"`
	Members     []pb.ServerAddress `json:"Members,omitempty"`
	MaxVolumeId needle.VolumeId    `json:"MaxVolumeId,omitempty"`
}

func (ms *MasterServer) ClusterStatusHandler(w http.ResponseWriter, r *http.Request) {
	ret := ClusterStatusResult{
		IsWriter:    ms.Topo.IsWriter(),
		MaxVolumeId: ms.Topo.GetMaxVolumeId(),
	}
	if ms.Topo.Coordinator != nil {
		ret.Writer = ms.Topo.Coordinator.Writer()
		ret.Members = ms.Topo.Coordinator.Members()
	}
	writeJsonQuiet(w, r, http.StatusOK, ret)
}

// ClusterHealthzHandler is the writer's readiness probe: it reports Locked when
// the writer is holding an admin lock, mirroring the old raft HealthzHandler.
func (ms *MasterServer) ClusterHealthzHandler(w http.ResponseWriter, r *http.Request) {
	if ms.Topo.IsWriter() {
		isLocked, err := ms.Topo.IsChildLocked()
		if err != nil {
			glog.Errorf("ClusterHealthzHandler: %+v", err)
		}
		if isLocked {
			w.WriteHeader(http.StatusLocked)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
