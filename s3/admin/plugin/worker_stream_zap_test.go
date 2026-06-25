package plugin

import (
	"testing"
	"time"

	"github.com/hanzoai/s3/s3/pb/plugin_pb"
	pluginwire "github.com/hanzoai/s3/s3/wire/plugin"
	pluginzapbridge "github.com/hanzoai/s3/s3/wire/plugin/zapbridge"
	"github.com/zap-proto/go/transport"
)

// TestServeWorkerStream_EndToEnd drives the plugin's PluginControlService over
// the real ZAP transport: it starts the plugin stream server, dials it as a
// worker, opens the WorkerStream, sends a hello, and asserts the admin replies
// with an AdminHello that accepts the connection — exercising the full cut
// (transport.ListenStream -> ServeWorkerStream -> bridge -> WorkerStream and
// back through the bridge on the client side).
func TestServeWorkerStream_EndToEnd(t *testing.T) {
	p, err := New(Options{})
	if err != nil {
		t.Fatalf("New plugin: %v", err)
	}
	defer p.Shutdown()

	srv, err := transport.ListenStream("tcp", "127.0.0.1:0", nil, func(method uint32, init []byte, s transport.Stream) {
		if method != pluginwire.PluginControlServiceWorkerStreamOrdinal {
			return
		}
		p.ServeWorkerStream(init, s)
	})
	if err != nil {
		t.Fatalf("ListenStream: %v", err)
	}
	defer srv.Close()

	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()

	stream, err := conn.OpenStream(pluginwire.PluginControlServiceWorkerStreamOrdinal, nil)
	if err != nil {
		t.Fatalf("OpenStream: %v", err)
	}

	hello := &plugin_pb.WorkerToAdminMessage{
		WorkerId: "w-e2e",
		Body: &plugin_pb.WorkerToAdminMessage_Hello{Hello: &plugin_pb.WorkerHello{
			WorkerId: "w-e2e",
			Address:  "127.0.0.1:7000",
			Capabilities: []*plugin_pb.JobTypeCapability{
				{JobType: "vacuum", CanDetect: true, CanExecute: true},
			},
		}},
	}
	if err := stream.Send(pluginzapbridge.MarshalWorkerToAdmin(hello)); err != nil {
		t.Fatalf("Send hello: %v", err)
	}

	type recvd struct {
		msg *plugin_pb.AdminToWorkerMessage
		err error
	}
	ch := make(chan recvd, 1)
	go func() {
		env, rErr := stream.Recv()
		if rErr != nil {
			ch <- recvd{err: rErr}
			return
		}
		m, dErr := pluginzapbridge.UnmarshalAdminToWorker(env)
		ch <- recvd{msg: m, err: dErr}
	}()

	select {
	case got := <-ch:
		if got.err != nil {
			t.Fatalf("Recv admin hello: %v", got.err)
		}
		ah := got.msg.GetHello()
		if ah == nil {
			t.Fatalf("first admin frame is not a hello: %v", got.msg)
		}
		if !ah.GetAccepted() {
			t.Fatalf("admin hello did not accept worker: %v", ah)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for admin hello")
	}

	if !p.HasCapableWorker("vacuum") {
		t.Fatal("admin did not register the worker's vacuum capability from hello")
	}
}
