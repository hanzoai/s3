package version

import (
	"fmt"

	"github.com/hanzoai/s3/s3/stats"
	"github.com/hanzoai/s3/s3/util"
)

var (
	MAJOR_VERSION  = int32(4)
	MINOR_VERSION  = int32(34)
	VERSION_NUMBER = fmt.Sprintf("%d.%02d", MAJOR_VERSION, MINOR_VERSION)
	VERSION        = util.SizeLimit + " " + VERSION_NUMBER
	COMMIT         = ""
)

func init() {
	// Set version info in stats for Prometheus metrics
	stats.SetVersionInfo(VERSION_NUMBER, COMMIT, util.SizeLimit)
}

func Version() string {
	return VERSION + " " + COMMIT
}
