// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// stream.go holds the per-frame request builders and response decoders for the
// HanzoFiler STREAMING RPCs (ListEntries, StreamRenameEntry, StreamMutateEntry,
// SubscribeMetadata). They compose the leaf converters in convert.go/entry.go
// and the unary builders/decoders in rpc.go, so the streaming seam carries the
// same *filer_pb message graph the FUSE engine already speaks while every frame
// crosses the socket as a filerwire buffer. Naming mirrors rpc.go:
// <Rpc>ReqToWire(*pb.<Rpc>Request) []byte and <Rpc>RespFromWire([]byte) (*pb.<Rpc>Response, error).

package filer

import (
	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	filerwire "github.com/hanzoai/s3/s3/wire/filer"
)

// ListEntries (server stream)

func ListEntriesReqToWire(r *filer_pb.ListEntriesRequest) []byte {
	return filerwire.NewListEntriesRequest(filerwire.ListEntriesRequestInput{
		Directory:          r.Directory,
		Prefix:             r.Prefix,
		StartFromFileName:  r.StartFromFileName,
		InclusiveStartFrom: r.InclusiveStartFrom,
		Limit:              r.Limit,
		SnapshotTsNs:       r.SnapshotTsNs,
	})
}

func ListEntriesRespFromWire(b []byte) (*filer_pb.ListEntriesResponse, error) {
	v, err := filerwire.WrapListEntriesResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &filer_pb.ListEntriesResponse{SnapshotTsNs: v.SnapshotTsNs()}
	if v.HasEntry() {
		resp.Entry = entryFromView(v.Entry())
	}
	return resp, nil
}

// SubscribeMetadata (server stream) — request only; the response reuses
// SubscribeMetadataResponseFromWire from convert.go.

func SubscribeMetadataReqToWire(r *filer_pb.SubscribeMetadataRequest) []byte {
	return filerwire.NewSubscribeMetadataRequest(filerwire.SubscribeMetadataRequestInput{
		ClientName:                   r.ClientName,
		PathPrefix:                   r.PathPrefix,
		SinceNs:                      r.SinceNs,
		Signature:                    r.Signature,
		PathPrefixes:                 r.PathPrefixes,
		ClientID:                     r.ClientId,
		UntilNs:                      r.UntilNs,
		ClientEpoch:                  r.ClientEpoch,
		Directories:                  r.Directories,
		ClientSupportsBatching:       r.ClientSupportsBatching,
		ClientSupportsMetadataChunks: r.ClientSupportsMetadataChunks,
		ClientSupportsIdleHeartbeat:  r.ClientSupportsIdleHeartbeat,
	})
}

// TraverseBfsMetadata (server stream) — request only; mount does not consume
// the response frames, so no decoder is needed here.

func TraverseBfsMetadataReqToWire(r *filer_pb.TraverseBfsMetadataRequest) []byte {
	return filerwire.NewTraverseBfsMetadataRequest(filerwire.TraverseBfsMetadataRequestInput{
		Directory:        r.Directory,
		ExcludedPrefixes: r.ExcludedPrefixes,
	})
}

// StreamRenameEntry (server stream)

func StreamRenameEntryReqToWire(r *filer_pb.StreamRenameEntryRequest) []byte {
	return filerwire.NewStreamRenameEntryRequest(filerwire.StreamRenameEntryRequestInput{
		OldDirectory: r.OldDirectory,
		OldName:      r.OldName,
		NewDirectory: r.NewDirectory,
		NewName:      r.NewName,
		Signatures:   r.Signatures,
	})
}

func StreamRenameEntryRespFromWire(b []byte) (*filer_pb.StreamRenameEntryResponse, error) {
	v, err := filerwire.WrapStreamRenameEntryResponse(b)
	if err != nil {
		return nil, err
	}
	return streamRenameEntryRespFromView(v), nil
}

func streamRenameEntryRespFromView(v filerwire.StreamRenameEntryResponse) *filer_pb.StreamRenameEntryResponse {
	resp := &filer_pb.StreamRenameEntryResponse{
		Directory: v.Directory(),
		TsNs:      v.TsNs(),
	}
	if v.HasEventNotification() {
		resp.EventNotification = eventNotificationFromView(v.EventNotification())
	}
	return resp
}

// StreamMutateEntry (bidi stream) — request carries a oneof of the four unary
// mutation requests; response carries a oneof of the four unary responses plus
// request_id/is_last/error/errno correlation fields.

func StreamMutateEntryReqToWire(r *filer_pb.StreamMutateEntryRequest) []byte {
	in := filerwire.StreamMutateEntryRequestInput{RequestID: r.RequestId}
	switch req := r.Request.(type) {
	case *filer_pb.StreamMutateEntryRequest_CreateRequest:
		in.Which = filerwire.StreamMutateRequestCreate
		in.CreateRequest = CreateEntryReqToWire(req.CreateRequest)
	case *filer_pb.StreamMutateEntryRequest_UpdateRequest:
		in.Which = filerwire.StreamMutateRequestUpdate
		in.UpdateRequest = UpdateEntryReqToWire(req.UpdateRequest)
	case *filer_pb.StreamMutateEntryRequest_DeleteRequest:
		in.Which = filerwire.StreamMutateRequestDelete
		in.DeleteRequest = DeleteEntryReqToWire(req.DeleteRequest)
	case *filer_pb.StreamMutateEntryRequest_RenameRequest:
		in.Which = filerwire.StreamMutateRequestRename
		in.RenameRequest = StreamRenameEntryReqToWire(req.RenameRequest)
	}
	return filerwire.NewStreamMutateEntryRequest(in)
}

func StreamMutateEntryRespFromWire(b []byte) (*filer_pb.StreamMutateEntryResponse, error) {
	v, err := filerwire.WrapStreamMutateEntryResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &filer_pb.StreamMutateEntryResponse{
		RequestId: v.RequestID(),
		IsLast:    v.IsLast(),
		Error:     v.Error(),
		Errno:     v.Errno(),
	}
	switch v.WhichResponse() {
	case filerwire.StreamMutateResponseCreate:
		cr := v.CreateResponse()
		resp.Response = &filer_pb.StreamMutateEntryResponse_CreateResponse{CreateResponse: &filer_pb.CreateEntryResponse{
			Error:         cr.Error(),
			ErrorCode:     filer_pb.FilerError(cr.ErrorCode()),
			MetadataEvent: metadataEventFromView(cr.HasMetadataEvent(), cr.MetadataEvent()),
		}}
	case filerwire.StreamMutateResponseUpdate:
		ur := v.UpdateResponse()
		resp.Response = &filer_pb.StreamMutateEntryResponse_UpdateResponse{UpdateResponse: &filer_pb.UpdateEntryResponse{
			MetadataEvent: metadataEventFromView(ur.HasMetadataEvent(), ur.MetadataEvent()),
		}}
	case filerwire.StreamMutateResponseDelete:
		dr := v.DeleteResponse()
		resp.Response = &filer_pb.StreamMutateEntryResponse_DeleteResponse{DeleteResponse: &filer_pb.DeleteEntryResponse{
			Error:         dr.Error(),
			MetadataEvent: metadataEventFromView(dr.HasMetadataEvent(), dr.MetadataEvent()),
		}}
	case filerwire.StreamMutateResponseRename:
		resp.Response = &filer_pb.StreamMutateEntryResponse_RenameResponse{
			RenameResponse: streamRenameEntryRespFromView(v.RenameResponse()),
		}
	}
	return resp, nil
}
