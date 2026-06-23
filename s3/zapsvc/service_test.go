package zapsvc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/hanzoai/s3/s3/zaprpc"
)

// memStore is the in-memory ObjectStore for the reference proof. The real S3
// engine implements the same interface.
type memStore struct {
	mu   sync.Mutex
	objs map[string]obj
}
type obj struct {
	data []byte
	ct   string
	etag string
}

func newMemStore() *memStore { return &memStore{objs: map[string]obj{}} }

func (m *memStore) Put(bucket, key string, data []byte, ct string) (string, error) {
	sum := sha256.Sum256(data)
	etag := hex.EncodeToString(sum[:])
	m.mu.Lock()
	m.objs[bucket+"/"+key] = obj{data: append([]byte(nil), data...), ct: ct, etag: etag}
	m.mu.Unlock()
	return etag, nil
}
func (m *memStore) Get(bucket, key string) ([]byte, string, string, error) {
	m.mu.Lock()
	o := m.objs[bucket+"/"+key]
	m.mu.Unlock()
	return o.data, o.ct, o.etag, nil
}

func freePort(t testing.TB) int {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	p := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()
	return p
}

// TestS3OverZAP proves the platform-law reference: an internal caller does
// PutObject + GetObject against the S3 service entirely over native ZAP —
// zero-copy schemas on the wire, no HTTP, no gRPC.
func TestS3OverZAP(t *testing.T) {
	port := freePort(t)
	srv := zaprpc.NewServer("s3-service", port)
	Register(srv, newMemStore())
	if err := srv.Start(); err != nil {
		t.Fatalf("server Start: %v", err)
	}
	defer srv.Stop()

	rpc, err := zaprpc.NewClient("s3-caller")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer rpc.Stop()
	s3 := NewClient(rpc, "127.0.0.1:"+strconv.Itoa(port))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := []byte("hanzo-s3-over-zap: the bytes are the message")
	etag, err := s3.PutObject(ctx, "hanzo", "vfs/blocks/abc", body, "application/octet-stream")
	if err != nil {
		t.Fatalf("PutObject: %v", err)
	}
	if etag == "" {
		t.Fatal("empty etag")
	}

	got, ct, getEtag, err := s3.GetObject(ctx, "hanzo", "vfs/blocks/abc")
	if err != nil {
		t.Fatalf("GetObject: %v", err)
	}
	if string(got) != string(body) {
		t.Fatalf("body mismatch:\n got=%q\nwant=%q", got, body)
	}
	if ct != "application/octet-stream" {
		t.Fatalf("contentType = %q", ct)
	}
	if getEtag != etag {
		t.Fatalf("etag mismatch put=%q get=%q", etag, getEtag)
	}
}

// TestS3OverZAP_LargeAndEmpty exercises a big body + a missing key over ZAP.
func TestS3OverZAP_LargeAndEmpty(t *testing.T) {
	port := freePort(t)
	srv := zaprpc.NewServer("s3-service-2", port)
	Register(srv, newMemStore())
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}
	defer srv.Stop()
	rpc, err := zaprpc.NewClient("s3-caller-2")
	if err != nil {
		t.Fatal(err)
	}
	defer rpc.Stop()
	s3 := NewClient(rpc, "127.0.0.1:"+strconv.Itoa(port))
	ctx := context.Background()

	big := make([]byte, 256*1024)
	for i := range big {
		big[i] = byte(i)
	}
	if _, err := s3.PutObject(ctx, "b", "big", big, "x"); err != nil {
		t.Fatalf("put big: %v", err)
	}
	got, _, _, err := s3.GetObject(ctx, "b", "big")
	if err != nil || len(got) != len(big) {
		t.Fatalf("get big: err=%v len=%d", err, len(got))
	}
	// missing key -> empty object, no error
	miss, _, _, err := s3.GetObject(ctx, "b", "nope")
	if err != nil || len(miss) != 0 {
		t.Fatalf("get missing: err=%v len=%d", err, len(miss))
	}
}
