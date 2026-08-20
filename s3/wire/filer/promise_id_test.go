// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package filerwire

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/zap-proto/go/rpc"
	"github.com/zap-proto/go/transport"
)

// echoServer answers every unary request with StatusOK under the caller's own
// PromiseID, pausing first so concurrent calls are genuinely in flight together
// (a serialised pair would not exercise the demultiplexer at all).
func echoServer(t *testing.T) *transport.Server {
	t.Helper()
	srv, err := transport.ListenStream("tcp", "127.0.0.1:0", func(env []byte) ([]byte, error) {
		call, err := rpc.ParseRequest(env)
		if err != nil {
			return nil, err
		}
		time.Sleep(100 * time.Millisecond)
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, []byte("ok")), nil
	}, nil)
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { _ = srv.Close() })
	return srv
}

// TestPromiseIDIsPerConnNotPerClient pins the invariant that makes pooled
// connections safe: a PromiseID identifies an in-flight call *on a connection*,
// so it must be allocated by the connection.
//
// pb.WithGrpcFilerClient builds a fresh client on every call while filerPool
// hands back one shared conn per filer address. When each client numbered its
// own calls from a private rpc.Session, every client restarted at 1, two
// concurrent calls landed on the same key in the transport's in-flight table,
// and the second registration evicted the first. The evicted caller's response
// was then dropped on arrival and it blocked for the life of its context —
// which for the S3 startup path was context.Background(), i.e. forever. That
// deadlocked NewS3ApiServer before it ever bound the S3 API port.
func TestPromiseIDIsPerConnNotPerClient(t *testing.T) {
	srv := echoServer(t)

	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	// Two independently constructed clients over the ONE pooled conn, exactly
	// as WithGrpcFilerClient does for two concurrent filer RPCs.
	clients := []*HanzoFilerClient{
		NewHanzoFilerClient(conn, nil),
		NewHanzoFilerClient(conn, nil),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		mu       sync.Mutex
		promises []uint32
		wg       sync.WaitGroup
	)
	for i, c := range clients {
		wg.Add(1)
		go func(i int, c *HanzoFilerClient) {
			defer wg.Done()
			p, body, err := c.LookupDirectoryEntry(ctx, nil)
			if err != nil {
				// Before the fix this is "context deadline exceeded": the
				// response came back and was discarded as unclaimed.
				t.Errorf("client %d: %v", i, err)
				return
			}
			if string(body) != "ok" {
				t.Errorf("client %d: body = %q, want %q", i, body, "ok")
			}
			mu.Lock()
			promises = append(promises, p.ID)
			mu.Unlock()
		}(i, c)
	}
	wg.Wait()

	if len(promises) != len(clients) {
		t.Fatalf("got %d responses, want %d", len(promises), len(clients))
	}
	if promises[0] == promises[1] {
		t.Fatalf("both calls used PromiseID %d; ids must be unique per connection", promises[0])
	}
}

// TestPromiseIDDoesNotCollideWithStreamIDs guards the other half of the same
// namespace: transport.OpenStream also draws stream ids from the connection
// counter, so unary calls must draw from it too or a call and a stream opened
// on one conn could share an id.
func TestPromiseIDDoesNotCollideWithStreamIDs(t *testing.T) {
	srv := echoServer(t)

	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	before := conn.NextPromiseID()
	p, _, err := NewHanzoFilerClient(conn, nil).LookupDirectoryEntry(ctx, nil)
	if err != nil {
		t.Fatalf("call: %v", err)
	}
	if p.ID <= before {
		t.Fatalf("call took PromiseID %d, not advanced past the connection counter at %d", p.ID, before)
	}
}
