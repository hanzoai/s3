// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package mountwire

import (
	"errors"
	"testing"
)

// memMounter is an in-memory Mounter for the round-trip test. It records the
// last collectionCapacity it was asked to configure.
type memMounter struct {
	gotCapacity int64
	called      bool
	failWith    error
}

func (m *memMounter) Configure(collectionCapacity int64) error {
	if m.failWith != nil {
		return m.failWith
	}
	m.called = true
	m.gotCapacity = collectionCapacity
	return nil
}

// TestRoundTrip proves Configure over the canonical github.com/zap-proto/go
// transport across a real TCP socket: the request crosses the wire as a ZAP RPC
// envelope carrying a zero-copy mountwire.ConfigureRequest, and the empty
// ConfigureResponse comes back — no HTTP, no protobuf, no gRPC.
func TestRoundTrip(t *testing.T) {
	mounter := &memMounter{}

	srv, err := Serve("tcp", "127.0.0.1:0", mounter)
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	const want int64 = 1 << 40 // exercises the full int64 width across the wire.
	if err := cli.Configure(want); err != nil {
		t.Fatalf("Configure: %v", err)
	}
	if !mounter.called {
		t.Fatalf("backend Configure was not invoked")
	}
	if mounter.gotCapacity != want {
		t.Fatalf("Configure collectionCapacity = %d, want %d", mounter.gotCapacity, want)
	}
}

// TestConfigureError proves a backend error surfaces to the client as a
// non-nil error (StatusInternal on the wire), not a silent success.
func TestConfigureError(t *testing.T) {
	mounter := &memMounter{failWith: errors.New("boom")}

	srv, err := Serve("tcp", "127.0.0.1:0", mounter)
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	if err := cli.Configure(7); err == nil {
		t.Fatalf("Configure: expected error, got nil")
	}
}

// TestMessageRoundTrip proves the zero-copy ConfigureRequest builder/view pair
// is byte-faithful without any transport: New then Wrap returns the same value.
func TestMessageRoundTrip(t *testing.T) {
	const want int64 = -42 // negative value confirms signed int64 round-trip.
	buf := NewConfigureRequest(ConfigureRequestInput{CollectionCapacity: want})
	v, err := WrapConfigureRequest(buf)
	if err != nil {
		t.Fatalf("WrapConfigureRequest: %v", err)
	}
	if got := v.CollectionCapacity(); got != want {
		t.Fatalf("CollectionCapacity = %d, want %d", got, want)
	}

	// The empty ConfigureResponse must also build and wrap cleanly.
	if _, err := WrapConfigureResponse(NewConfigureResponse(ConfigureResponseInput{})); err != nil {
		t.Fatalf("WrapConfigureResponse: %v", err)
	}
}
