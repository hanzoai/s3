package gentest

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/hanzoai/s3/s3/zaprpc"
)

// echoImpl is a hand-written EchoServer (the generated interface). In the real
// migration this is the SeaweedFS service implementation.
type echoImpl struct{}

func (echoImpl) Ping(_ context.Context, req *PingRequest) (*PingResponse, error) {
	return &PingResponse{Msg: "pong:" + req.Msg, Seq: req.Seq + 1}, nil
}

func freePort(t *testing.T) int {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	p := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()
	return p
}

// TestGeneratedServiceOverZap proves the COMPLETE generated unary path:
// protoc-gen-zap server stub (RegisterEchoServer) + client stub (EchoClient)
// carrying generated message codecs over the real ZAP transport.
func TestGeneratedServiceOverZap(t *testing.T) {
	port := freePort(t)
	srv := zaprpc.NewServer("echo-server", port)
	RegisterEchoServer(srv, echoImpl{})
	if err := srv.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer srv.Stop()

	cli, err := zaprpc.NewClient("echo-client")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer cli.Stop()

	echo := NewEchoClient(cli, "127.0.0.1:"+strconv.Itoa(port))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := echo.Ping(ctx, &PingRequest{Msg: "hanzo", Seq: 41, Tags: []string{"a", "b"}})
	if err != nil {
		t.Fatalf("Ping: %v", err)
	}
	if resp.Msg != "pong:hanzo" || resp.Seq != 42 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}
