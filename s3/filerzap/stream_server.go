// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// stream_server.go is the server-side half of the filer's streaming RPCs: a
// filerstream.Server that wraps the existing filer_pb.HanzoFilerServer so its
// server-streaming methods answer over transport.ListenStream. ListEntries —
// the one the IAM/policy stores use to enumerate a directory — is migrated; it
// drives the gRPC streaming server method through a Send-adapter that ships each
// proto response as a zero-copy wire frame. The other 5 streaming RPCs are
// stubbed (no-op) until their backends are wired.

package filerzap

import (
	"context"

	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	filerwire "github.com/hanzoai/s3/s3/wire/filer"
	"github.com/hanzoai/s3/s3/wire/filer/filerstream"
	"google.golang.org/grpc"
)

// ListEntriesRespToWire builds a zero-copy ListEntriesResponse wire frame.
func ListEntriesRespToWire(r *filer_pb.ListEntriesResponse) []byte {
	in := filerwire.ListEntriesResponseInput{SnapshotTsNs: r.SnapshotTsNs}
	if r.Entry != nil {
		in.Entry = EntryToWire(r.Entry)
	}
	return filerwire.NewListEntriesResponse(in)
}

// streamServerBackend adapts a filer_pb.HanzoFilerServer to filerstream.Server.
type streamServerBackend struct {
	fs  filer_pb.HanzoFilerServer
	ctx context.Context
}

// NewStreamServer returns a filerstream.Server that serves fs's streaming RPCs
// over ZAP. Pass it to filerstream.Handler / transport.ListenStream.
func NewStreamServer(fs filer_pb.HanzoFilerServer) filerstream.Server {
	return streamServerBackend{fs: fs, ctx: context.Background()}
}

func (b streamServerBackend) ListEntries(v filerwire.ListEntriesRequest, s *filerstream.ListEntriesStream) error {
	req := &filer_pb.ListEntriesRequest{
		Directory:          v.Directory(),
		Prefix:             v.Prefix(),
		StartFromFileName:  v.StartFromFileName(),
		InclusiveStartFrom: v.InclusiveStartFrom(),
		Limit:              v.Limit(),
		SnapshotTsNs:       v.SnapshotTsNs(),
	}
	return b.fs.ListEntries(req, &listEntriesSendAdapter{ctx: b.ctx, out: s})
}

// --- not yet migrated streaming RPCs ---

func (b streamServerBackend) StreamRenameEntry(filerwire.StreamRenameEntryRequest, *filerstream.StreamRenameEntryStream) error {
	return errFilerRPCNotMigrated
}
func (b streamServerBackend) TraverseBfsMetadata(filerwire.TraverseBfsMetadataRequest, *filerstream.TraverseBfsMetadataStream) error {
	return errFilerRPCNotMigrated
}
func (b streamServerBackend) SubscribeMetadata(filerwire.SubscribeMetadataRequest, *filerstream.SubscribeMetadataStream) error {
	return errFilerRPCNotMigrated
}
func (b streamServerBackend) SubscribeLocalMetadata(filerwire.SubscribeMetadataRequest, *filerstream.SubscribeMetadataStream) error {
	return errFilerRPCNotMigrated
}
func (b streamServerBackend) StreamMutateEntry(filerwire.StreamMutateEntryRequest, *filerstream.StreamMutateEntryStream) error {
	return errFilerRPCNotMigrated
}

// listEntriesSendAdapter implements grpc.ServerStreamingServer[ListEntriesResponse]
// (= filer_pb.HanzoFiler_ListEntriesServer): the existing fs.ListEntries calls
// Send for each entry; we ship it as a wire frame on the ZAP stream. The
// embedded grpc.ServerStream supplies the interface's remaining methods, which
// the filer engine does not call on this path.
type listEntriesSendAdapter struct {
	grpc.ServerStream
	ctx context.Context
	out *filerstream.ListEntriesStream
}

func (a *listEntriesSendAdapter) Send(resp *filer_pb.ListEntriesResponse) error {
	return a.out.Send(ListEntriesRespToWire(resp))
}

func (a *listEntriesSendAdapter) Context() context.Context { return a.ctx }
