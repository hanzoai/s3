// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// stream_server.go is the server-side half of the filer's streaming RPCs: a
// filerstream.Server that wraps the existing filer_pb.HanzoFilerServer so its
// streaming methods answer over transport.ListenStream. All 6 are implemented —
// ListEntries, StreamRenameEntry, TraverseBfsMetadata, SubscribeMetadata and
// SubscribeLocalMetadata (server-streams) drive the engine method through a
// Send-adapter that ships each proto response as a zero-copy wire frame;
// StreamMutateEntry (bidi) replays the open frame then pumps Recv/Send.
//
// Each handler runs against the per-stream s.Context() (zap-proto/go v1.6.2+):
// it is cancelled when the stream ends OR the peer drops the connection, so a
// long-lived subscription that goes idle and then disconnects is observed and
// the engine handler returns — no goroutine/subscription leak.

package filer

import (
	"context"

	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	filerwire "github.com/hanzoai/s3/s3/wire/filer"
	"github.com/hanzoai/s3/s3/wire/filer/filerstream"
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
	fs filer_pb.HanzoFilerServer
}

// NewStreamServer returns a filerstream.Server that serves fs's streaming RPCs
// over ZAP. Pass it to filerstream.Handler / transport.ListenStream.
func NewStreamServer(fs filer_pb.HanzoFilerServer) filerstream.Server {
	return streamServerBackend{fs: fs}
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
	return b.fs.ListEntries(req, &listEntriesSendAdapter{ctx: s.Context(), out: s})
}

func (b streamServerBackend) StreamRenameEntry(v filerwire.StreamRenameEntryRequest, s *filerstream.StreamRenameEntryStream) error {
	return b.fs.StreamRenameEntry(streamRenameEntryReqFromView(v), &streamRenameSendAdapter{ctx: s.Context(), out: s})
}
func (b streamServerBackend) TraverseBfsMetadata(v filerwire.TraverseBfsMetadataRequest, s *filerstream.TraverseBfsMetadataStream) error {
	return b.fs.TraverseBfsMetadata(traverseBfsMetadataReqFromView(v), &traverseBfsSendAdapter{ctx: s.Context(), out: s})
}
func (b streamServerBackend) SubscribeMetadata(v filerwire.SubscribeMetadataRequest, s *filerstream.SubscribeMetadataStream) error {
	return b.fs.SubscribeMetadata(subscribeMetadataReqFromView(v), &subscribeMetadataSendAdapter{ctx: s.Context(), out: s})
}
func (b streamServerBackend) SubscribeLocalMetadata(v filerwire.SubscribeMetadataRequest, s *filerstream.SubscribeMetadataStream) error {
	return b.fs.SubscribeLocalMetadata(subscribeMetadataReqFromView(v), &subscribeMetadataSendAdapter{ctx: s.Context(), out: s})
}
func (b streamServerBackend) StreamMutateEntry(init filerwire.StreamMutateEntryRequest, s *filerstream.StreamMutateEntryStream) error {
	return b.fs.StreamMutateEntry(&streamMutateAdapter{ctx: s.Context(), s: s, init: init, initPending: true})
}

// listEntriesSendAdapter implements filer_pb.HanzoFiler_ListEntriesServer: the
// existing fs.ListEntries calls Send for each entry; we ship it as a wire frame
// on the ZAP stream.
type listEntriesSendAdapter struct {
	ctx context.Context
	out *filerstream.ListEntriesStream
}

func (a *listEntriesSendAdapter) Send(resp *filer_pb.ListEntriesResponse) error {
	return a.out.Send(ListEntriesRespToWire(resp))
}

func (a *listEntriesSendAdapter) Context() context.Context { return a.ctx }

// --- server-stream Send adapters (proto Send -> wire frame on the ZAP stream) ---

type streamRenameSendAdapter struct {
	ctx context.Context
	out *filerstream.StreamRenameEntryStream
}

func (a *streamRenameSendAdapter) Send(resp *filer_pb.StreamRenameEntryResponse) error {
	return a.out.Send(StreamRenameEntryRespToWire(resp))
}
func (a *streamRenameSendAdapter) Context() context.Context { return a.ctx }

type traverseBfsSendAdapter struct {
	ctx context.Context
	out *filerstream.TraverseBfsMetadataStream
}

func (a *traverseBfsSendAdapter) Send(resp *filer_pb.TraverseBfsMetadataResponse) error {
	return a.out.Send(TraverseBfsMetadataRespToWire(resp))
}
func (a *traverseBfsSendAdapter) Context() context.Context { return a.ctx }

type subscribeMetadataSendAdapter struct {
	ctx context.Context
	out *filerstream.SubscribeMetadataStream
}

func (a *subscribeMetadataSendAdapter) Send(resp *filer_pb.SubscribeMetadataResponse) error {
	return a.out.Send(SubscribeMetadataResponseToWire(resp))
}
func (a *subscribeMetadataSendAdapter) Context() context.Context { return a.ctx }

// streamMutateAdapter implements filer_pb.HanzoFiler_StreamMutateEntryServer:
// the opener's first frame (init) is replayed by the first Recv, then subsequent
// frames come off the ZAP stream; Send ships replies.
type streamMutateAdapter struct {
	ctx         context.Context
	s           *filerstream.StreamMutateEntryStream
	init        filerwire.StreamMutateEntryRequest
	initPending bool
}

func (a *streamMutateAdapter) Recv() (*filer_pb.StreamMutateEntryRequest, error) {
	if a.initPending {
		a.initPending = false
		return streamMutateReqFromView(a.init), nil
	}
	v, err := a.s.Recv()
	if err != nil {
		return nil, err
	}
	return streamMutateReqFromView(v), nil
}
func (a *streamMutateAdapter) Send(resp *filer_pb.StreamMutateEntryResponse) error {
	return a.s.Send(StreamMutateEntryRespToWire(resp))
}
func (a *streamMutateAdapter) Context() context.Context { return a.ctx }

// --- streaming request view -> proto ---

func traverseBfsMetadataReqFromView(v filerwire.TraverseBfsMetadataRequest) *filer_pb.TraverseBfsMetadataRequest {
	req := &filer_pb.TraverseBfsMetadataRequest{Directory: v.Directory()}
	for i := 0; i < v.ExcludedPrefixesLen(); i++ {
		req.ExcludedPrefixes = append(req.ExcludedPrefixes, v.ExcludedPrefixe(i))
	}
	return req
}

func subscribeMetadataReqFromView(v filerwire.SubscribeMetadataRequest) *filer_pb.SubscribeMetadataRequest {
	req := &filer_pb.SubscribeMetadataRequest{
		ClientName:                   v.ClientName(),
		PathPrefix:                   v.PathPrefix(),
		SinceNs:                      v.SinceNs(),
		Signature:                    v.Signature(),
		ClientId:                     v.ClientID(),
		UntilNs:                      v.UntilNs(),
		ClientEpoch:                  v.ClientEpoch(),
		ClientSupportsBatching:       v.ClientSupportsBatching(),
		ClientSupportsMetadataChunks: v.ClientSupportsMetadataChunks(),
		ClientSupportsIdleHeartbeat:  v.ClientSupportsIdleHeartbeat(),
	}
	for i := 0; i < v.PathPrefixesLen(); i++ {
		req.PathPrefixes = append(req.PathPrefixes, v.PathPrefixe(i))
	}
	for i := 0; i < v.DirectoriesLen(); i++ {
		req.Directories = append(req.Directories, v.Directory(i))
	}
	return req
}

func streamMutateReqFromView(v filerwire.StreamMutateEntryRequest) *filer_pb.StreamMutateEntryRequest {
	req := &filer_pb.StreamMutateEntryRequest{RequestId: v.RequestID()}
	switch v.WhichRequest() {
	case filerwire.StreamMutateRequestCreate:
		req.Request = &filer_pb.StreamMutateEntryRequest_CreateRequest{CreateRequest: createEntryReqFromView(v.CreateRequest())}
	case filerwire.StreamMutateRequestUpdate:
		req.Request = &filer_pb.StreamMutateEntryRequest_UpdateRequest{UpdateRequest: updateEntryReqFromView(v.UpdateRequest())}
	case filerwire.StreamMutateRequestDelete:
		req.Request = &filer_pb.StreamMutateEntryRequest_DeleteRequest{DeleteRequest: deleteEntryReqFromView(v.DeleteRequest())}
	case filerwire.StreamMutateRequestRename:
		req.Request = &filer_pb.StreamMutateEntryRequest_RenameRequest{RenameRequest: streamRenameEntryReqFromView(v.RenameRequest())}
	}
	return req
}

// --- streaming response proto -> wire ---

func StreamRenameEntryRespToWire(r *filer_pb.StreamRenameEntryResponse) []byte {
	in := filerwire.StreamRenameEntryResponseInput{Directory: r.Directory, TsNs: r.TsNs}
	if r.EventNotification != nil {
		in.EventNotification = eventNotificationToWire(r.EventNotification)
	}
	return filerwire.NewStreamRenameEntryResponse(in)
}

func TraverseBfsMetadataRespToWire(r *filer_pb.TraverseBfsMetadataResponse) []byte {
	in := filerwire.TraverseBfsMetadataResponseInput{Directory: r.Directory}
	if r.Entry != nil {
		in.Entry = EntryToWire(r.Entry)
	}
	return filerwire.NewTraverseBfsMetadataResponse(in)
}

func StreamMutateEntryRespToWire(r *filer_pb.StreamMutateEntryResponse) []byte {
	in := filerwire.StreamMutateEntryResponseInput{
		RequestID: r.RequestId, IsLast: r.IsLast, Error: r.Error, Errno: r.Errno,
	}
	switch resp := r.Response.(type) {
	case *filer_pb.StreamMutateEntryResponse_CreateResponse:
		in.Which = filerwire.StreamMutateResponseCreate
		in.CreateResponse = filerwire.NewCreateEntryResponse(filerwire.CreateEntryResponseInput{
			Error: resp.CreateResponse.Error, ErrorCode: uint32(resp.CreateResponse.ErrorCode),
			MetadataEvent: SubscribeMetadataResponseToWire(resp.CreateResponse.MetadataEvent),
		})
	case *filer_pb.StreamMutateEntryResponse_UpdateResponse:
		in.Which = filerwire.StreamMutateResponseUpdate
		in.UpdateResponse = filerwire.NewUpdateEntryResponse(filerwire.UpdateEntryResponseInput{
			MetadataEvent: SubscribeMetadataResponseToWire(resp.UpdateResponse.MetadataEvent),
		})
	case *filer_pb.StreamMutateEntryResponse_DeleteResponse:
		in.Which = filerwire.StreamMutateResponseDelete
		in.DeleteResponse = filerwire.NewDeleteEntryResponse(filerwire.DeleteEntryResponseInput{
			Error: resp.DeleteResponse.Error, MetadataEvent: SubscribeMetadataResponseToWire(resp.DeleteResponse.MetadataEvent),
		})
	case *filer_pb.StreamMutateEntryResponse_RenameResponse:
		in.Which = filerwire.StreamMutateResponseRename
		in.RenameResponse = StreamRenameEntryRespToWire(resp.RenameResponse)
	}
	return filerwire.NewStreamMutateEntryResponse(in)
}
