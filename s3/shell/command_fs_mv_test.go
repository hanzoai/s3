package shell

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/hanzoai/s3/s3/svc/filer"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/pb/filerstub"
	filerwire "github.com/hanzoai/s3/s3/wire/filer"
	"github.com/hanzoai/s3/s3/wire/filer/filerstream"
	"github.com/zap-proto/go/transport"
)

type fsMvTestFilerServer struct {
	filerstub.FilerServer

	lookupReq *filer_pb.LookupDirectoryEntryRequest
	renameReq *filer_pb.AtomicRenameEntryRequest
}

func (s *fsMvTestFilerServer) LookupDirectoryEntry(_ context.Context, req *filer_pb.LookupDirectoryEntryRequest) (*filer_pb.LookupDirectoryEntryResponse, error) {
	s.lookupReq = req
	if req.Directory == "/dst" && req.Name == "dir" {
		return &filer_pb.LookupDirectoryEntryResponse{
			Entry: &filer_pb.Entry{
				Name:        "dir",
				IsDirectory: true,
				Attributes:  &filer_pb.FuseAttributes{},
			},
		}, nil
	}
	return nil, fmt.Errorf("NotFound: not found")
}

func (s *fsMvTestFilerServer) AtomicRenameEntry(_ context.Context, req *filer_pb.AtomicRenameEntryRequest) (*filer_pb.AtomicRenameEntryResponse, error) {
	s.renameReq = req
	return &filer_pb.AtomicRenameEntryResponse{}, nil
}

func TestFsMvMovesIntoExistingDestinationDirectory(t *testing.T) {
	filerServer := &fsMvTestFilerServer{}
	commandEnv, cleanup := newFsMvTestCommandEnv(t, filerServer)
	defer cleanup()

	var output bytes.Buffer
	err := (&commandFsMv{}).Do([]string{"/src/file", "/dst/dir"}, commandEnv, &output)
	if err != nil {
		t.Fatalf("fs.mv returned error: %v", err)
	}

	if filerServer.lookupReq == nil {
		t.Fatal("expected fs.mv to look up destination entry")
	}
	if filerServer.lookupReq.Directory != "/dst" || filerServer.lookupReq.Name != "dir" {
		t.Fatalf("destination lookup = directory %q name %q, want /dst dir", filerServer.lookupReq.Directory, filerServer.lookupReq.Name)
	}
	if filerServer.renameReq == nil {
		t.Fatal("expected fs.mv to issue rename")
	}
	if filerServer.renameReq.NewDirectory != "/dst/dir" || filerServer.renameReq.NewName != "file" {
		t.Fatalf("rename target = directory %q name %q, want /dst/dir file", filerServer.renameReq.NewDirectory, filerServer.renameReq.NewName)
	}
}

func newFsMvTestCommandEnv(t *testing.T, filerServer filer_pb.HanzoFilerServer) (*CommandEnv, func()) {
	t.Helper()

	// Serve the fake filer over the native ZAP transport — fs.mv dials the filer
	// via pb.WithGrpcFilerClient, which now opens a transport.Dial("tcp", ...) ZAP
	// connection, so the command exercises the real ZAP client+server path.
	srv, err := transport.ListenStream("tcp", "127.0.0.1:0",
		filerwire.Dispatch(filer.NewServerBackend(filerServer)),
		filerstream.Handler(filer.NewStreamServer(filerServer)))
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	// FilerAddress carries the real listener port as the ".grpcPort" suffix so
	// ToGrpcAddress() (which WithGrpcFilerClient dials) resolves to it exactly;
	// the leading ":0" public port is unused on this path.
	grpcPort := srv.Addr().(*net.TCPAddr).Port

	cleanup := func() {
		_ = srv.Close()
	}

	return &CommandEnv{
		option: &ShellOptions{
			FilerAddress:   pb.ServerAddress(fmt.Sprintf("127.0.0.1:0.%d", grpcPort)),
			GrpcDialOption: pb.DialOption{},
			Directory:      "/",
		},
	}, cleanup
}
