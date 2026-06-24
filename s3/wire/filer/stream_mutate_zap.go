// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
	zap "github.com/zap-proto/go"
)

// --- StreamMutateEntryRequest (oneof request) ---

// StreamMutateEntryRequest carries a request_id plus a oneof "request": WhichRequest
// reports the set member (0 = none, else the proto field number), and each
// member has its own nested-object slot.
const (
	streamMutateEntryRequestRequestIDOff = 0  // request_id (field 1, uint64)
	streamMutateEntryRequestWhichOff     = 8  // u32 oneof tag (0 = none)
	streamMutateEntryRequestCreateOff    = 12 // create_request (field 2, CreateEntryRequest object)
	streamMutateEntryRequestUpdateOff    = 16 // update_request (field 3, UpdateEntryRequest object)
	streamMutateEntryRequestDeleteOff    = 20 // delete_request (field 4, DeleteEntryRequest object)
	streamMutateEntryRequestRenameOff    = 24 // rename_request (field 5, StreamRenameEntryRequest object)
	streamMutateEntryRequestSize         = 28
)

// StreamMutateEntryRequest is a zero-copy view into a ZAP-encoded StreamMutateEntryRequest message.
type StreamMutateEntryRequest struct{ o zap.Object }

// WrapStreamMutateEntryRequest parses b and returns a typed view.
func WrapStreamMutateEntryRequest(b []byte) (StreamMutateEntryRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return StreamMutateEntryRequest{}, err
	}
	return StreamMutateEntryRequest{o: m.Root()}, nil
}

func (t StreamMutateEntryRequest) RequestID() uint64 {
	return t.o.Uint64(streamMutateEntryRequestRequestIDOff)
}

// WhichRequest reports which oneof member is set (one of the StreamMutateRequest*
// consts; 0 = none).
func (t StreamMutateEntryRequest) WhichRequest() uint32 {
	return t.o.Uint32(streamMutateEntryRequestWhichOff)
}

// CreateRequest returns the create_request variant (null-object unless WhichRequest
// == StreamMutateRequestCreate).
func (t StreamMutateEntryRequest) CreateRequest() CreateEntryRequest {
	return CreateEntryRequest{o: t.o.Object(streamMutateEntryRequestCreateOff)}
}

// UpdateRequest returns the update_request variant.
func (t StreamMutateEntryRequest) UpdateRequest() UpdateEntryRequest {
	return UpdateEntryRequest{o: t.o.Object(streamMutateEntryRequestUpdateOff)}
}

// DeleteRequest returns the delete_request variant.
func (t StreamMutateEntryRequest) DeleteRequest() DeleteEntryRequest {
	return DeleteEntryRequest{o: t.o.Object(streamMutateEntryRequestDeleteOff)}
}

// RenameRequest returns the rename_request variant.
func (t StreamMutateEntryRequest) RenameRequest() StreamRenameEntryRequest {
	return StreamRenameEntryRequest{o: t.o.Object(streamMutateEntryRequestRenameOff)}
}

// StreamMutateEntryRequestInput collects the field values for
// NewStreamMutateEntryRequest. Set Which to the selected member and supply that
// member's standalone sub-buffer.
type StreamMutateEntryRequestInput struct {
	RequestID     uint64
	Which         uint32
	CreateRequest []byte
	UpdateRequest []byte
	DeleteRequest []byte
	RenameRequest []byte
}

// NewStreamMutateEntryRequest builds a ZAP-encoded StreamMutateEntryRequest message
// from in. Only the member selected by in.Which is written.
func NewStreamMutateEntryRequest(in StreamMutateEntryRequestInput) []byte {
	b := zap.NewBuilder(512)
	ob := b.StartObject(streamMutateEntryRequestSize)
	ob.SetUint64(streamMutateEntryRequestRequestIDOff, in.RequestID)
	ob.SetUint32(streamMutateEntryRequestWhichOff, in.Which)
	switch in.Which {
	case StreamMutateRequestCreate:
		setNestedObject(b, ob, streamMutateEntryRequestCreateOff, in.CreateRequest)
	case StreamMutateRequestUpdate:
		setNestedObject(b, ob, streamMutateEntryRequestUpdateOff, in.UpdateRequest)
	case StreamMutateRequestDelete:
		setNestedObject(b, ob, streamMutateEntryRequestDeleteOff, in.DeleteRequest)
	case StreamMutateRequestRename:
		setNestedObject(b, ob, streamMutateEntryRequestRenameOff, in.RenameRequest)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- StreamMutateEntryResponse (oneof response) ---

// StreamMutateEntryResponse carries a request_id + is_last + error/errno plus a
// oneof "response": WhichResponse reports the set member (0 = none, else the proto
// field number), and each member has its own nested-object slot.
const (
	streamMutateEntryResponseRequestIDOff = 0  // request_id (field 1, uint64)
	streamMutateEntryResponseIsLastOff    = 8  // is_last (field 2, bool)
	streamMutateEntryResponseWhichOff     = 12 // u32 oneof tag (0 = none)
	streamMutateEntryResponseCreateOff    = 16 // create_response (field 3, CreateEntryResponse object)
	streamMutateEntryResponseUpdateOff    = 20 // update_response (field 4, UpdateEntryResponse object)
	streamMutateEntryResponseDeleteOff    = 24 // delete_response (field 5, DeleteEntryResponse object)
	streamMutateEntryResponseRenameOff    = 28 // rename_response (field 6, StreamRenameEntryResponse object)
	streamMutateEntryResponseErrorOff     = 32 // error (field 7, string)
	streamMutateEntryResponseErrnoOff     = 40 // errno (field 8, int32)
	streamMutateEntryResponseSize         = 48
)

// StreamMutateEntryResponse is a zero-copy view into a ZAP-encoded StreamMutateEntryResponse message.
type StreamMutateEntryResponse struct{ o zap.Object }

// WrapStreamMutateEntryResponse parses b and returns a typed view.
func WrapStreamMutateEntryResponse(b []byte) (StreamMutateEntryResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return StreamMutateEntryResponse{}, err
	}
	return StreamMutateEntryResponse{o: m.Root()}, nil
}

func (t StreamMutateEntryResponse) RequestID() uint64 {
	return t.o.Uint64(streamMutateEntryResponseRequestIDOff)
}
func (t StreamMutateEntryResponse) IsLast() bool { return t.o.Bool(streamMutateEntryResponseIsLastOff) }
func (t StreamMutateEntryResponse) Error() string {
	return t.o.Text(streamMutateEntryResponseErrorOff)
}
func (t StreamMutateEntryResponse) Errno() int32 { return t.o.Int32(streamMutateEntryResponseErrnoOff) }

// WhichResponse reports which oneof member is set (one of the StreamMutateResponse*
// consts; 0 = none).
func (t StreamMutateEntryResponse) WhichResponse() uint32 {
	return t.o.Uint32(streamMutateEntryResponseWhichOff)
}

// CreateResponse returns the create_response variant (null-object unless
// WhichResponse == StreamMutateResponseCreate).
func (t StreamMutateEntryResponse) CreateResponse() CreateEntryResponse {
	return CreateEntryResponse{o: t.o.Object(streamMutateEntryResponseCreateOff)}
}

// UpdateResponse returns the update_response variant.
func (t StreamMutateEntryResponse) UpdateResponse() UpdateEntryResponse {
	return UpdateEntryResponse{o: t.o.Object(streamMutateEntryResponseUpdateOff)}
}

// DeleteResponse returns the delete_response variant.
func (t StreamMutateEntryResponse) DeleteResponse() DeleteEntryResponse {
	return DeleteEntryResponse{o: t.o.Object(streamMutateEntryResponseDeleteOff)}
}

// RenameResponse returns the rename_response variant.
func (t StreamMutateEntryResponse) RenameResponse() StreamRenameEntryResponse {
	return StreamRenameEntryResponse{o: t.o.Object(streamMutateEntryResponseRenameOff)}
}

// StreamMutateEntryResponseInput collects the field values for
// NewStreamMutateEntryResponse. Set Which to the selected member and supply that
// member's standalone sub-buffer.
type StreamMutateEntryResponseInput struct {
	RequestID      uint64
	IsLast         bool
	Which          uint32
	CreateResponse []byte
	UpdateResponse []byte
	DeleteResponse []byte
	RenameResponse []byte
	Error          string
	Errno          int32
}

// NewStreamMutateEntryResponse builds a ZAP-encoded StreamMutateEntryResponse message
// from in. Only the member selected by in.Which is written.
func NewStreamMutateEntryResponse(in StreamMutateEntryResponseInput) []byte {
	b := zap.NewBuilder(512)
	ob := b.StartObject(streamMutateEntryResponseSize)
	ob.SetUint64(streamMutateEntryResponseRequestIDOff, in.RequestID)
	ob.SetBool(streamMutateEntryResponseIsLastOff, in.IsLast)
	ob.SetUint32(streamMutateEntryResponseWhichOff, in.Which)
	switch in.Which {
	case StreamMutateResponseCreate:
		setNestedObject(b, ob, streamMutateEntryResponseCreateOff, in.CreateResponse)
	case StreamMutateResponseUpdate:
		setNestedObject(b, ob, streamMutateEntryResponseUpdateOff, in.UpdateResponse)
	case StreamMutateResponseDelete:
		setNestedObject(b, ob, streamMutateEntryResponseDeleteOff, in.DeleteResponse)
	case StreamMutateResponseRename:
		setNestedObject(b, ob, streamMutateEntryResponseRenameOff, in.RenameResponse)
	}
	ob.SetText(streamMutateEntryResponseErrorOff, in.Error)
	ob.SetInt32(streamMutateEntryResponseErrnoOff, in.Errno)
	ob.FinishAsRoot()
	return b.Finish()
}
