package pb

// grpc_raft_holdout.go is the SINGLE remaining gRPC surface in Hanzo S3.
//
// Every application RPC — master, volume, filer, broker, object — now rides the
// native ZAP transport (github.com/zap-proto/go); the WithXClient helpers in
// grpc_client_server.go dial it and the *_pb packages are message-only. The one
// thing still on gRPC is the MASTER'S RAFT CONSENSUS TRANSPORT:
//
//   - seaweedfs/raft serves its RPCs via protobuf.RegisterRaftServer on a
//     *grpc.Server (command/master.go), and
//   - hashicorp/raft uses github.com/Jille/raft-grpc-transport, which dials
//     peers with []grpc.DialOption (server/raft_hashicorp.go).
//
// Both are a Raft-engine internal wire, not an S3 service contract, so they are
// migrated separately (server/consensus_server.go already implements the
// ZAP-native, post-quantum Lux-consensus replacement; wiring it in is a live
// consensus-engine swap gated on multi-master cluster-formation testing). Until
// that lands, the primitives the Raft transport needs live HERE, isolated and
// documented, so the rest of package pb is gRPC-free:
//
//   - NewGrpcServer / ServeGrpcOnLocalSocket: construct and (optionally) Unix-
//     socket-serve the master's raft *grpc.Server.
//   - DialOption.Grpc(): adapt the gRPC-free DialOption to the one grpc.DialOption
//     the Jille transport and seaweedfs raft transporter consume.

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/util/request_id"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

const (
	// Max_Message_Size bounds a single raft RPC frame on the master grpc server.
	Max_Message_Size = 1 << 30 // 1 GB

	// gRPC keepalive settings - must be consistent between client and server.
	GrpcKeepAliveTime        = 60 * time.Second // ping interval when no activity
	GrpcKeepAliveTimeout     = 20 * time.Second // ping timeout
	GrpcKeepAliveMinimumTime = 20 * time.Second // minimum interval between client pings (enforcement)
)

// NewGrpcServer builds the master's raft consensus *grpc.Server with the shared
// keepalive/window/message-size tuning and a request-id unary interceptor. It is
// the only grpc.Server left in S3; everything else serves over ZAP.
func NewGrpcServer(opts ...grpc.ServerOption) *grpc.Server {
	var options []grpc.ServerOption
	options = append(options,
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    GrpcKeepAliveTime,    // server pings client if no activity for this long
			Timeout: GrpcKeepAliveTimeout, // ping timeout
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             GrpcKeepAliveMinimumTime, // min time a client should wait before sending a ping
			PermitWithoutStream: true,
		}),
		grpc.MaxRecvMsgSize(Max_Message_Size),
		grpc.MaxSendMsgSize(Max_Message_Size),
		grpc.MaxConcurrentStreams(1000),          // Allow more concurrent streams
		grpc.InitialWindowSize(16*1024*1024),     // 16MB initial window
		grpc.InitialConnWindowSize(16*1024*1024), // 16MB connection window
		grpc.MaxHeaderListSize(8*1024*1024),      // 8MB header list limit
		grpc.UnaryInterceptor(requestIDUnaryInterceptor()),
	)
	for _, opt := range opts {
		if opt != nil {
			options = append(options, opt)
		}
	}
	return grpc.NewServer(options...)
}

// ServeGrpcOnLocalSocket starts serving the raft grpc server on a Unix socket
// if one is registered for the given port (see RegisterLocalGrpcSocket).
func ServeGrpcOnLocalSocket(grpcServer *grpc.Server, grpcPort int) {
	socketPath := GetLocalGrpcSocket(grpcPort)
	if socketPath == "" {
		return
	}
	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		glog.Warningf("Failed to remove old gRPC socket %s: %v", socketPath, err)
	}
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		glog.Errorf("Failed to listen on gRPC Unix socket %s: %v", socketPath, err)
		return
	}
	glog.V(0).Infof("gRPC also listening on Unix socket %s", socketPath)
	go func() {
		if err := grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			glog.Errorf("gRPC Unix socket server error on %s: %v", socketPath, err)
		}
		os.Remove(socketPath)
	}()
}

// requestIDUnaryInterceptor propagates the x-amz-request-id across the raft grpc
// server: it reads any inbound request id from gRPC metadata (or mints one),
// stores it in the handler context, and mirrors it on the outgoing context and
// the response trailer.
func requestIDUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		var reqID string
		if incomingMd, ok := metadata.FromIncomingContext(ctx); ok {
			if idList := incomingMd.Get(request_id.AmzRequestIDHeader); len(idList) > 0 {
				reqID = idList[0]
			}
		}
		if reqID == "" {
			reqID = uuid.New().String()
		}
		ctx = request_id.Set(ctx, reqID)
		ctx = metadata.AppendToOutgoingContext(ctx, request_id.AmzRequestIDHeader, reqID)
		grpc.SetTrailer(ctx, metadata.Pairs(request_id.AmzRequestIDHeader, reqID))
		return handler(ctx, req)
	}
}

// Grpc adapts the gRPC-free DialOption to the single grpc credential the master
// raft transport needs: TLS when configured, plaintext otherwise. The ZAP path
// (filer/master/volume/broker) ignores DialOption — its TLS is configured in the
// transport layer (DialTLS + PQTLSConfig).
func (o DialOption) Grpc() grpc.DialOption {
	if o.TLS != nil {
		return grpc.WithTransportCredentials(credentials.NewTLS(o.TLS))
	}
	return grpc.WithTransportCredentials(insecure.NewCredentials())
}
