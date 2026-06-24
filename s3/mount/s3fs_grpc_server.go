package mount

import (
	"fmt"

	"github.com/hanzoai/s3/s3/glog"
	mountwire "github.com/hanzoai/s3/s3/wire/mount"
)

func (wfs *WFS) Configure(req []byte) ([]byte, error) {
	request, err := mountwire.WrapConfigureRequest(req)
	if err != nil {
		return nil, err
	}
	if wfs.option.Collection == "" {
		return nil, fmt.Errorf("mount quota only works when mounted to a new folder with a collection")
	}
	glog.V(0).Infof("quota changed from %d to %d", wfs.option.Quota, request.CollectionCapacity())
	wfs.option.Quota = request.CollectionCapacity()
	return mountwire.NewConfigureResponse(mountwire.ConfigureResponseInput{}), nil
}
