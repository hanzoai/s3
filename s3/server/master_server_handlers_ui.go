package s3server

import (
	"net/http"
	"time"

	"github.com/hanzoai/s3/s3/util/version"

	ui "github.com/hanzoai/s3/s3/server/master_ui"
	"github.com/hanzoai/s3/s3/stats"
)

func (ms *MasterServer) uiStatusHandler(w http.ResponseWriter, r *http.Request) {
	infos := make(map[string]interface{})
	infos["Up Time"] = time.Since(startTime).Truncate(time.Second).String()
	infos["Max Volume Id"] = ms.Topo.GetMaxVolumeId()

	if ms.consensus == nil {
		return
	}

	args := struct {
		Version           string
		Topology          interface{}
		RaftServer        ConsensusStatusView
		Stats             map[string]interface{}
		Counters          *stats.ServerStats
		VolumeSizeLimitMB uint32
	}{
		version.Version(),
		ms.Topo.ToInfo(),
		ms.consensus.StatusView(),
		infos,
		serverStats,
		ms.option.VolumeSizeLimitMB,
	}
	ui.StatusTpl.Execute(w, args)
}
