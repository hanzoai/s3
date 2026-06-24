// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package zapsvc

import (
	"bytes"
	"testing"
)

// memStore is an in-memory ObjectStore for the round-trip test.
type memStore struct {
	objs map[string]struct {
		data []byte
		ct   string
		etag string
	}
}

func newMemStore() *memStore {
	return &memStore{objs: map[string]struct {
		data []byte
		ct   string
		etag string
	}{}}
}

func (m *memStore) Put(bucket, key string, data []byte, ct string) (string, error) {
	etag := "etag-" + key
	m.objs[bucket+"/"+key] = struct {
		data []byte
		ct   string
		etag string
	}{append([]byte(nil), data...), ct, etag}
	return etag, nil
}

func (m *memStore) Get(bucket, key string) ([]byte, string, string, error) {
	o := m.objs[bucket+"/"+key]
	return o.data, o.ct, o.etag, nil
}

// TestRoundTrip proves Put then Get over the canonical github.com/zap-proto/go
// transport across a real TCP socket: the bytes cross the wire as ZAP RPC
// envelopes carrying zero-copy objectwire payloads — no HTTP, no protobuf.
func TestRoundTrip(t *testing.T) {
	store := newMemStore()

	srv, err := Serve("tcp", "127.0.0.1:0", store)
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	want := []byte("the bytes are the message")
	etag, err := cli.PutObject("blobs", "team-go/x.bin", want, "application/octet-stream")
	if err != nil {
		t.Fatalf("PutObject: %v", err)
	}
	if etag != "etag-team-go/x.bin" {
		t.Fatalf("PutObject etag = %q", etag)
	}

	got, ct, getEtag, err := cli.GetObject("blobs", "team-go/x.bin")
	if err != nil {
		t.Fatalf("GetObject: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("GetObject data = %q, want %q", got, want)
	}
	if ct != "application/octet-stream" {
		t.Fatalf("GetObject contentType = %q", ct)
	}
	if getEtag != etag {
		t.Fatalf("GetObject etag = %q, want %q", getEtag, etag)
	}
}
