// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package filer

import (
	"context"
	"errors"
	"testing"

	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/pb/filerstub"
	filerwire "github.com/hanzoai/s3/s3/wire/filer"
	"github.com/zap-proto/go/transport"
)

// fakeFiler is a minimal in-memory filer_pb.HanzoFilerServer for the round-trip:
// the embedded filerstub.FilerServer supplies the methods this test doesn't exercise.
type fakeFiler struct {
	filerstub.FilerServer
	entries map[string]*filer_pb.Entry
}

func (f *fakeFiler) key(dir, name string) string { return dir + "/" + name }

func (f *fakeFiler) CreateEntry(_ context.Context, req *filer_pb.CreateEntryRequest) (*filer_pb.CreateEntryResponse, error) {
	f.entries[f.key(req.Directory, req.Entry.Name)] = req.Entry
	return &filer_pb.CreateEntryResponse{}, nil
}

func (f *fakeFiler) LookupDirectoryEntry(_ context.Context, req *filer_pb.LookupDirectoryEntryRequest) (*filer_pb.LookupDirectoryEntryResponse, error) {
	e, ok := f.entries[f.key(req.Directory, req.Name)]
	if !ok {
		return nil, filer_pb.ErrNotFound
	}
	return &filer_pb.LookupDirectoryEntryResponse{Entry: e}, nil
}

func (f *fakeFiler) DeleteEntry(_ context.Context, req *filer_pb.DeleteEntryRequest) (*filer_pb.DeleteEntryResponse, error) {
	delete(f.entries, f.key(req.Directory, req.Name))
	return &filer_pb.DeleteEntryResponse{}, nil
}

// TestServerBackendRoundTrip proves the filer answers over the canonical ZAP
// transport: a real filer_pb.HanzoFilerServer served via filerwire.Serve
// (NewServerBackend) and reached by the ZAP client (NewZapFilerClient over
// transport.Dial) — Create then Lookup then Delete, full Entry fidelity, no gRPC.
func TestServerBackendRoundTrip(t *testing.T) {
	fs := &fakeFiler{entries: map[string]*filer_pb.Entry{}}

	srv, err := filerwire.Serve("tcp", "127.0.0.1:0", NewServerBackend(fs))
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()
	client := NewZapFilerClient(conn, nil)

	entry := &filer_pb.Entry{
		Name:        "policy.json",
		IsDirectory: false,
		Content:     []byte(`{"effect":"allow"}`),
		Attributes:  &filer_pb.FuseAttributes{FileSize: 18, FileMode: 0o644, Uid: 1000},
		Extended:    map[string][]byte{"Content-Type": []byte("application/json")},
	}
	if _, err := client.CreateEntry(context.Background(), &filer_pb.CreateEntryRequest{
		Directory: "/etc/iam/policies", Entry: entry,
	}); err != nil {
		t.Fatalf("CreateEntry: %v", err)
	}

	got, err := client.LookupDirectoryEntry(context.Background(), &filer_pb.LookupDirectoryEntryRequest{
		Directory: "/etc/iam/policies", Name: "policy.json",
	})
	if err != nil {
		t.Fatalf("LookupDirectoryEntry: %v", err)
	}
	if got.Entry == nil || got.Entry.Name != entry.Name ||
		string(got.Entry.Content) != string(entry.Content) ||
		got.Entry.Attributes == nil || got.Entry.Attributes.FileSize != 18 ||
		string(got.Entry.Extended["Content-Type"]) != "application/json" {
		t.Fatalf("Lookup returned %+v, want fidelity with %+v", got.Entry, entry)
	}

	if _, err := client.DeleteEntry(context.Background(), &filer_pb.DeleteEntryRequest{
		Directory: "/etc/iam/policies", Name: "policy.json",
	}); err != nil {
		t.Fatalf("DeleteEntry: %v", err)
	}
	// Not-found is conveyed as a nil Entry on an OK response (the rpc envelope
	// carries no error type); filer_pb.LookupEntry maps that to ErrNotFound.
	afterDelete, err := client.LookupDirectoryEntry(context.Background(), &filer_pb.LookupDirectoryEntryRequest{
		Directory: "/etc/iam/policies", Name: "policy.json",
	})
	if err != nil {
		t.Fatalf("LookupDirectoryEntry after delete: %v", err)
	}
	if afterDelete.Entry != nil {
		t.Fatal("expected nil Entry (not-found) after delete")
	}
}

// TestLookupEntryMapsMissingToNotFound pins the contract callers depend on to
// tell "absent" from "present" over this transport. The raw client reports a
// missing name as an OK response carrying a nil Entry, so err alone says only
// that the filer was reached — reading it as existence reports every new name
// as already taken. filer_pb.LookupEntry is the mapping back to ErrNotFound.
func TestLookupEntryMapsMissingToNotFound(t *testing.T) {
	fs := &fakeFiler{entries: map[string]*filer_pb.Entry{}}

	srv, err := filerwire.Serve("tcp", "127.0.0.1:0", NewServerBackend(fs))
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()
	client := NewZapFilerClient(conn, nil)

	req := &filer_pb.LookupDirectoryEntryRequest{Directory: "/buckets", Name: "never-created"}

	raw, err := client.LookupDirectoryEntry(context.Background(), req)
	if err != nil {
		t.Fatalf("raw lookup of a missing name: got error %v, want nil", err)
	}
	if raw.GetEntry() != nil {
		t.Fatalf("raw lookup of a missing name: got Entry %+v, want nil", raw.GetEntry())
	}

	if _, err := filer_pb.LookupEntry(context.Background(), client, req); !errors.Is(err, filer_pb.ErrNotFound) {
		t.Fatalf("LookupEntry of a missing name: got %v, want ErrNotFound", err)
	}
}
