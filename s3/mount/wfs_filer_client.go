package mount

import (
	"sync/atomic"

	"github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/util"

	"github.com/zap-proto/go/transport"
)

var _ = filer_pb.FilerClient(&WFS{})

// WithFilerClient dials the next filer over the native ZAP transport and runs fn
// with a filer_pb.HanzoFilerClient backed by that connection (see
// filerClientAdapter). It rotates through option.FilerAddresses on failure, the
// same retry policy the gRPC path used. The connection is closed when fn
// returns; the adapter pools nothing, matching the prior per-call dial.
func (wfs *WFS) WithFilerClient(streamingMode bool, fn func(filer_pb.HanzoFilerClient) error) (err error) {

	return util.Retry("filer zap", func() error {

		i := atomic.LoadInt32(&wfs.option.filerIndex)
		n := len(wfs.option.FilerAddresses)
		for x := 0; x < n; x++ {

			filerGrpcAddress := wfs.option.FilerAddresses[i].ToGrpcAddress()
			conn, dialErr := transport.Dial("tcp", filerGrpcAddress)
			if dialErr != nil {
				err = dialErr
			} else {
				err = fn(newFilerClientAdapter(conn))
				_ = conn.Close()
			}

			if err == nil {
				atomic.StoreInt32(&wfs.option.filerIndex, i)
				return nil
			}

			i++
			if i >= int32(n) {
				i = 0
			}

		}
		return err
	})

}

func (wfs *WFS) AdjustedUrl(location *filer_pb.Location) string {
	if wfs.option.VolumeServerAccess == "publicUrl" {
		return location.PublicUrl
	}
	return location.Url
}

func (wfs *WFS) GetDataCenter() string {
	return wfs.option.DataCenter
}
