package version

import (
	"fmt"

	"github.com/hanzoai/s3/s3/stats"
	"github.com/hanzoai/s3/s3/util"
)

// Hanzo S3 versions on its own v1 line. The inherited 4.34 numbering tracked a
// different project, and it was never resolvable as a Go version either: the
// module path is github.com/hanzoai/s3 with no major suffix, so only v1.x.x
// tags can be required by a dependent — which is why importers had to pin
// pseudo-versions. One line now covers the module tag, the image tag, and the
// number the binary reports.
var (
	MAJOR_VERSION  = int32(1)
	MINOR_VERSION  = int32(0)
	PATCH_VERSION  = int32(4)
	VERSION_NUMBER = fmt.Sprintf("%d.%d.%d", MAJOR_VERSION, MINOR_VERSION, PATCH_VERSION)
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
