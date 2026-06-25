package pb

import (
	"runtime"
	"testing"
)

// TestResolveLocalGrpcSocket_RemotePortCollision is a regression test for
// issue #9254. A `s3 server` process registers a Unix socket for its
// in-process volume server on host A. A standalone `s3 volume` on host B
// happens to use the same gRPC port. Dials from the master to host B must
// continue out over TCP — they must NOT be hijacked into host A's local
// socket on the basis of port match alone.
//
// The connection-invalidation tests that previously shared this file covered
// the gRPC cached-*grpc.ClientConn machinery (shouldInvalidateConnection /
// isClientSideMarshalError). That machinery was removed with the gRPC→ZAP rip:
// every application RPC now rides a pooled ZAP transport.Conn (see the WithX
// helpers), which has no shared-ClientConn-teardown failure mode to guard
// against. The socket registry below stays — it routes both the ZAP same-host
// dial fast-path and the master raft holdout.
func TestResolveLocalGrpcSocket_RemotePortCollision(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-socket routing is disabled on Windows (#9430)")
	}
	// Snapshot and restore global state so the test does not leak into others.
	localGrpcSocketsLock.Lock()
	prevSockets := localGrpcSockets
	prevHosts := localGrpcHosts
	localGrpcSockets = make(map[int]string)
	localGrpcHosts = make(map[int]map[string]struct{})
	localGrpcSocketsLock.Unlock()
	t.Cleanup(func() {
		localGrpcSocketsLock.Lock()
		localGrpcSockets = prevSockets
		localGrpcHosts = prevHosts
		localGrpcSocketsLock.Unlock()
	})

	const localHost = "10.0.0.2"
	const remoteHost = "10.0.0.3"
	const collidingPort = 17334
	const socketPath = "/tmp/hanzo-volume-grpc-17334.sock"

	RegisterLocalGrpcSocket(localHost, collidingPort, socketPath)

	cases := []struct {
		name    string
		address string
		want    string
	}{
		{"local advertised host routes to socket", localHost + ":17334", socketPath},
		{"loopback v4 routes to socket", "127.0.0.1:17334", socketPath},
		{"localhost routes to socket", "localhost:17334", socketPath},
		{"loopback v6 routes to socket", "[::1]:17334", socketPath},
		{"empty host (bare port) routes to socket", ":17334", socketPath},
		{"remote host with same port stays on TCP", remoteHost + ":17334", ""},
		{"unrelated host with same port stays on TCP", "192.168.1.5:17334", ""},
		{"unregistered port stays on TCP", localHost + ":17335", ""},
		{"malformed address stays on TCP", "not-a-host-port", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := resolveLocalGrpcSocket(tc.address); got != tc.want {
				t.Fatalf("resolveLocalGrpcSocket(%q) = %q, want %q", tc.address, got, tc.want)
			}
		})
	}
}
