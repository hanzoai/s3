package plugin

import (
	"context"
	"errors"
	"io"

	"github.com/hanzoai/s3/s3/pb/plugin_pb"
	pluginzapbridge "github.com/hanzoai/s3/s3/wire/plugin/zapbridge"
	"github.com/zap-proto/go/transport"
)

// WorkerStreamServer is the server view of one PluginControlService.WorkerStream
// stream: receive worker frames, send admin frames, observe the stream's
// lifetime via Context. It abstracts the underlying transport so WorkerStream's
// body is wire-agnostic.
type WorkerStreamServer interface {
	Recv() (*plugin_pb.WorkerToAdminMessage, error)
	Send(*plugin_pb.AdminToWorkerMessage) error
	Context() context.Context
}

// ServeWorkerStream is the transport.StreamHandler for the
// PluginControlService.WorkerStream ordinal: it adapts one ZAP stream into a
// WorkerStreamServer (decoding/encoding frames through pluginzapbridge) and
// runs the existing WorkerStream handler against it. The first worker frame
// arrives as the stream's init payload, so it is replayed as the first Recv.
func (r *Plugin) ServeWorkerStream(init []byte, s *transport.Stream) {
	_ = r.WorkerStream(&zapWorkerStream{
		stream: s,
		init:   init,
		hasIn:  len(init) > 0,
		ctx:    context.Background(),
	})
}

// zapWorkerStream adapts a transport.Stream to WorkerStreamServer, translating
// each frame between ZAP bytes and *plugin_pb.* via pluginzapbridge.
type zapWorkerStream struct {
	stream *transport.Stream
	init   []byte
	hasIn  bool
	ctx    context.Context
}

func (z *zapWorkerStream) Context() context.Context { return z.ctx }

func (z *zapWorkerStream) Recv() (*plugin_pb.WorkerToAdminMessage, error) {
	if z.hasIn {
		z.hasIn = false
		return pluginzapbridge.UnmarshalWorkerToAdmin(z.init)
	}
	body, err := z.stream.Recv()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, io.EOF
		}
		return nil, err
	}
	return pluginzapbridge.UnmarshalWorkerToAdmin(body)
}

func (z *zapWorkerStream) Send(msg *plugin_pb.AdminToWorkerMessage) error {
	return z.stream.Send(pluginzapbridge.MarshalAdminToWorker(msg))
}
