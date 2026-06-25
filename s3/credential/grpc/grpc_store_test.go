package grpc

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/zap-proto/go/transport"

	"github.com/hanzoai/s3/s3/credential"
	_ "github.com/hanzoai/s3/s3/credential/memory"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/iam_pb"
	"github.com/hanzoai/s3/s3/security"
	s3server "github.com/hanzoai/s3/s3/server"
	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

// TestIamGrpcStore_RoundTrip pins the production client path: IamGrpcStore dials
// the filer's IAM service over the canonical ZAP transport (no gRPC, no
// Bearer-token metadata) at FilerAddress.ToIamZapAddress() and exercises the
// real handler backed by an in-memory credential manager. The legacy
// admin-signed-token contract is gone — connection-level security is the ZAP
// transport itself — so this drives the full handler (business logic over the
// wire), not an auth gate.
func TestIamGrpcStore_RoundTrip(t *testing.T) {
	// Real IamGrpcServer backed by an in-memory credential manager so we run
	// the full handler, not a stub.
	cm, err := credential.NewCredentialManager(credential.StoreTypeMemory, nil, "")
	if err != nil {
		t.Fatalf("NewCredentialManager: %v", err)
	}
	defer cm.Shutdown()

	iamSrv := s3server.NewIamGrpcServer(cm, security.SigningKey("unused-under-zap"))

	// The store dials FilerAddress.ToIamZapAddress() == grpcPort+10000. Bind the
	// IAM ZAP listener on an ephemeral port Z, then advertise a filer address
	// whose grpc port is Z-10000 so the store's +10000 offset lands back on Z.
	// macOS ephemeral ports live in the high range (49152+), so Z-10000 is
	// always a valid positive port — no fragile re-bind loop needed.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen on IAM ZAP addr: %v", err)
	}
	zapPort := ln.Addr().(*net.TCPAddr).Port
	srv := transport.Serve(ln, func(env []byte) ([]byte, error) {
		return iamwire.DispatchHanzoIdentityAccessManagement(iamSrv, env)
	})
	t.Cleanup(func() { _ = srv.Close() })

	grpcPort := zapPort - 10000
	if grpcPort < 1 {
		t.Fatalf("ephemeral ZAP port %d too low to derive grpc port", zapPort)
	}
	// pb.ServerAddress encodes the grpc port in the trailing ".N" segment; the
	// host:0 segment is the unused HTTP port.
	serverAddr := pb.ServerAddress(fmt.Sprintf("127.0.0.1:0.%d", grpcPort))

	store := &IamGrpcStore{}
	store.SetFilerAddressFunc(func() pb.ServerAddress { return serverAddr }, pb.DialOption{})

	ctx := context.Background()

	if _, err := store.ListUsers(ctx); err != nil {
		t.Fatalf("ListUsers on fresh store: %v", err)
	}

	if err := store.CreateUser(ctx, &iam_pb.Identity{Name: "alice"}); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	got, err := store.GetUser(ctx, "alice")
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if got == nil || got.Name != "alice" {
		t.Fatalf("GetUser returned %+v, want name=alice", got)
	}

	users, err := store.ListUsers(ctx)
	if err != nil {
		t.Fatalf("ListUsers after create: %v", err)
	}
	if len(users) != 1 || users[0] != "alice" {
		t.Fatalf("ListUsers returned %v, want [alice]", users)
	}
}
