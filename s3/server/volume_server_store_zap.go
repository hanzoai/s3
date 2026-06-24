package s3server

import (
	"context"

	"github.com/hanzoai/s3/s3/pb/volume_server_pb"
	volume_serverwire "github.com/hanzoai/s3/s3/wire/volume_server"
)

// volume_server_store_zap.go adapts *VolumeServer (whose methods carry the gRPC
// (ctx, *pb.Req)->(*pb.Resp, error) shape) to the native ZAP
// volume_serverwire.VolumeServerStore contract: every unary method takes one
// request buffer (Wrap<Req>) and returns one response buffer (New<Resp>). The
// existing handler bodies stay in their lane on the pb domain model — the store
// adapter is a thin conv boundary (wire view -> pb request, pb response -> wire
// builder). This is the volume_server analogue of agentconv: it lets the ZAP
// transport drive the real volume server without gRPC, while the storage layer
// keeps its on-disk protobuf domain types untouched.
//
// Additive strangler: this binds the wire contract to the live handlers without
// removing the gRPC registration, so the tree stays green and the running server
// keeps serving until the registration is flipped to transport.ListenStream and
// the client callbacks move to the wire Client in lockstep.

// stateProtoToWire builds a wire VolumeServerState buffer from the pb state.
func stateProtoToWire(st *volume_server_pb.VolumeServerState) []byte {
	return volume_serverwire.NewVolumeServerState(volume_serverwire.VolumeServerStateInput{
		Maintenance: st.GetMaintenance(),
		Version:     st.GetVersion(),
	})
}

// stateWireToProto reads a wire VolumeServerState view into a pb state.
func stateWireToProto(v volume_serverwire.VolumeServerState) *volume_server_pb.VolumeServerState {
	return &volume_server_pb.VolumeServerState{
		Maintenance: v.Maintenance(),
		Version:     v.Version(),
	}
}

// volumeServerStore wraps *VolumeServer and satisfies
// volume_serverwire.VolumeServerStore. It is a distinct type (not *VolumeServer
// itself) because the live handlers keep the gRPC (ctx, *pb.Req)->(*pb.Resp)
// signatures, which collide name-for-name with the store's
// (req []byte)->([]byte) contract — exactly the split mq_agentwire's
// unaryHandler uses. Each method Wraps the request view, delegates to the real
// gRPC-shaped handler on the embedded *VolumeServer, then Builds the response
// buffer.
type volumeServerStore struct{ vs *VolumeServer }

// GetState implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) GetState(req []byte) ([]byte, error) {
	resp, err := s.vs.GetState(context.Background(), &volume_server_pb.GetStateRequest{})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewGetStateResponse(volume_serverwire.GetStateResponseInput{
		State: stateProtoToWire(resp.GetState()),
	}), nil
}

// SetState implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) SetState(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapSetStateRequest(req)
	if err != nil {
		return nil, err
	}
	var pbState *volume_server_pb.VolumeServerState
	if sv, ok := view.State(); ok {
		pbState = stateWireToProto(sv)
	}
	resp, err := s.vs.SetState(context.Background(), &volume_server_pb.SetStateRequest{State: pbState})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewSetStateResponse(volume_serverwire.SetStateResponseInput{
		State: stateProtoToWire(resp.GetState()),
	}), nil
}
