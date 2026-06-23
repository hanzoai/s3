package zaprpc

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"
)

// freePort returns an OS-assigned free TCP port on loopback.
func freePort(t testing.TB) int {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("freePort: %v", err)
	}
	p := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()
	return p
}

const codeEchoRSL uint8 = 1

// echoRSL decodes a RemoteStorageLocation request and re-encodes it as the
// response — the smallest end-to-end exercise of codec + ZAP transport, the
// shape every generated unary stub will take.
func echoRSL(_ context.Context, req []byte) ([]byte, error) {
	var m remoteStorageLocation
	if err := m.UnmarshalZap(req); err != nil {
		return nil, err
	}
	return m.MarshalZap(), nil
}

func startEchoServer(t testing.TB) (addr string, stop func()) {
	port := freePort(t)
	srv := NewServer("s3-test-server", port)
	srv.Handle(codeEchoRSL, echoRSL)
	if err := srv.Start(); err != nil {
		t.Fatalf("server Start: %v", err)
	}
	return "127.0.0.1:" + strconv.Itoa(port), srv.Stop
}

func TestUnaryRPCRoundTrip(t *testing.T) {
	addr, stopSrv := startEchoServer(t)
	defer stopSrv()

	cli, err := NewClient("s3-test-client")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer cli.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	in := remoteStorageLocation{Name: "wasabi", Bucket: "hanzo-cold", Path: "/vfs/blocks", ListingCacheTTLSeconds: 300}
	respBytes, err := cli.Call(ctx, addr, codeEchoRSL, in.MarshalZap())
	if err != nil {
		t.Fatalf("Call: %v", err)
	}
	var out remoteStorageLocation
	if err := out.UnmarshalZap(respBytes); err != nil {
		t.Fatalf("UnmarshalZap response: %v", err)
	}
	if out != in {
		t.Fatalf("round-trip over the wire mismatch:\n in=%+v\nout=%+v", in, out)
	}
}

// TestUnaryRPCManyCalls exercises connection reuse (memoized peer) + a few
// payload shapes across repeated calls on one client.
func TestUnaryRPCManyCalls(t *testing.T) {
	addr, stopSrv := startEchoServer(t)
	defer stopSrv()
	cli, err := NewClient("s3-test-client-2")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer cli.Stop()
	ctx := context.Background()

	for i := 0; i < 50; i++ {
		in := remoteStorageLocation{Name: "n" + strconv.Itoa(i), Bucket: "b", Path: "/p", ListingCacheTTLSeconds: int32(i)}
		resp, err := cli.Call(ctx, addr, codeEchoRSL, in.MarshalZap())
		if err != nil {
			t.Fatalf("call %d: %v", i, err)
		}
		var out remoteStorageLocation
		if err := out.UnmarshalZap(resp); err != nil || out != in {
			t.Fatalf("call %d mismatch: err=%v in=%+v out=%+v", i, err, in, out)
		}
	}
}

// BenchmarkUnaryRPC measures the full client→server→client round trip
// (frame + zap transport + unframe + codec on both ends).
func BenchmarkUnaryRPC(b *testing.B) {
	addr, stopSrv := startEchoServer(b)
	defer stopSrv()
	cli, err := NewClient("s3-bench-client")
	if err != nil {
		b.Fatalf("NewClient: %v", err)
	}
	defer cli.Stop()
	ctx := context.Background()
	in := (&remoteStorageLocation{Name: "wasabi", Bucket: "hanzo-cold", Path: "/vfs/blocks", ListingCacheTTLSeconds: 300}).MarshalZap()
	// warm the connection (dial + handshake) so it's out of the measured loop
	if _, err := cli.Call(ctx, addr, codeEchoRSL, in); err != nil {
		b.Fatalf("warmup: %v", err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := cli.Call(ctx, addr, codeEchoRSL, in); err != nil {
			b.Fatalf("call: %v", err)
		}
	}
}

// BenchmarkCodecMarshal isolates the codec (no transport) for comparison.
func BenchmarkCodecMarshal(b *testing.B) {
	m := &remoteStorageLocation{Name: "wasabi", Bucket: "hanzo-cold", Path: "/vfs/blocks", ListingCacheTTLSeconds: 300}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out remoteStorageLocation
		_ = out.UnmarshalZap(m.MarshalZap())
	}
}
