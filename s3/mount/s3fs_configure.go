package mount

import (
	"fmt"

	"github.com/hanzoai/s3/s3/glog"
)

// Configure implements mountwire.Mounter — it applies a quota change pushed over
// the native ZAP mount socket (formerly the HanzoMount gRPC Configure RPC).
func (wfs *WFS) Configure(collectionCapacity int64) error {
	if wfs.option.Collection == "" {
		return fmt.Errorf("mount quota only works when mounted to a new folder with a collection")
	}
	glog.V(0).Infof("quota changed from %d to %d", wfs.option.Quota, collectionCapacity)
	wfs.option.Quota = collectionCapacity
	return nil
}
