// Package lifecyclerpc is the single pb<->ZAP-wire bridge for the
// HanzoS3LifecycleInternal.LifecycleDelete RPC. It is the one place that
// knows how to encode the in-process s3_lifecycle_pb domain types into the
// generated s3_lifecyclewire envelopes and decode them back, so both the
// client adapters (shell, worker) and the server dispatch share exactly one
// mapping. The pb types remain the in-process control-plane domain model;
// only the wire crossing goes through ZAP.
package lifecyclerpc

import (
	"context"

	"github.com/hanzoai/s3/s3/pb/s3_lifecycle_pb"
	s3_lifecyclewire "github.com/hanzoai/s3/s3/wire/s3_lifecycle"
)

// EncodeRequest builds a ZAP LifecycleDeleteRequest envelope from the pb
// request.
func EncodeRequest(req *s3_lifecycle_pb.LifecycleDeleteRequest) []byte {
	if req == nil {
		return s3_lifecyclewire.NewLifecycleDeleteRequest(s3_lifecyclewire.LifecycleDeleteRequestInput{})
	}
	in := s3_lifecyclewire.LifecycleDeleteRequestInput{
		Bucket:               req.Bucket,
		ObjectPath:           req.ObjectPath,
		VersionID:            req.VersionId,
		RuleHash:             req.RuleHash,
		ActionKind:           uint32(req.ActionKind),
		StreamShard:          req.StreamShard,
		StreamDelaySeconds:   req.StreamDelaySeconds,
		StreamPositionTsNs:   req.StreamPositionTsNs,
		StreamPositionOffset: req.StreamPositionOffset,
		EngineSnapshotID:     req.EngineSnapshotId,
	}
	if id := req.ExpectedIdentity; id != nil {
		in.ExpectedIdentity = &s3_lifecyclewire.EntryIdentityInput{
			MtimeNs:      id.MtimeNs,
			Size:         id.Size,
			HeadFid:      id.HeadFid,
			ExtendedHash: id.ExtendedHash,
		}
	}
	return s3_lifecyclewire.NewLifecycleDeleteRequest(in)
}

// DecodeRequest parses a ZAP LifecycleDeleteRequest envelope into the pb
// request the server business logic consumes.
func DecodeRequest(b []byte) (*s3_lifecycle_pb.LifecycleDeleteRequest, error) {
	v, err := s3_lifecyclewire.WrapLifecycleDeleteRequest(b)
	if err != nil {
		return nil, err
	}
	req := &s3_lifecycle_pb.LifecycleDeleteRequest{
		Bucket:               v.Bucket(),
		ObjectPath:           v.ObjectPath(),
		VersionId:            v.VersionID(),
		RuleHash:             v.RuleHash(),
		ActionKind:           s3_lifecycle_pb.ActionKind(v.ActionKind()),
		StreamShard:          v.StreamShard(),
		StreamDelaySeconds:   v.StreamDelaySeconds(),
		StreamPositionTsNs:   v.StreamPositionTsNs(),
		StreamPositionOffset: v.StreamPositionOffset(),
		EngineSnapshotId:     v.EngineSnapshotID(),
	}
	if id := v.ExpectedIdentity(); id.HeadFid() != "" || id.MtimeNs() != 0 || id.Size() != 0 || len(id.ExtendedHash()) != 0 {
		req.ExpectedIdentity = &s3_lifecycle_pb.EntryIdentity{
			MtimeNs:      id.MtimeNs(),
			Size:         id.Size(),
			HeadFid:      id.HeadFid(),
			ExtendedHash: id.ExtendedHash(),
		}
	}
	return req, nil
}

// EncodeResponse builds a ZAP LifecycleDeleteResponse envelope from the pb
// response.
func EncodeResponse(resp *s3_lifecycle_pb.LifecycleDeleteResponse) []byte {
	if resp == nil {
		return s3_lifecyclewire.NewLifecycleDeleteResponse(s3_lifecyclewire.LifecycleDeleteResponseInput{})
	}
	return s3_lifecyclewire.NewLifecycleDeleteResponse(s3_lifecyclewire.LifecycleDeleteResponseInput{
		Outcome: uint32(resp.Outcome),
		Reason:  resp.Reason,
	})
}

// DecodeResponse parses a ZAP LifecycleDeleteResponse envelope into the pb
// response the dailyrun pipeline expects.
func DecodeResponse(b []byte) (*s3_lifecycle_pb.LifecycleDeleteResponse, error) {
	v, err := s3_lifecyclewire.WrapLifecycleDeleteResponse(b)
	if err != nil {
		return nil, err
	}
	return &s3_lifecycle_pb.LifecycleDeleteResponse{
		Outcome: s3_lifecycle_pb.LifecycleDeleteOutcome(v.Outcome()),
		Reason:  v.Reason(),
	}, nil
}

// Client issues LifecycleDelete over a ZAP wire client, presenting the
// in-process pb signature to callers.
func Client(c *s3_lifecyclewire.HanzoS3LifecycleInternalClient, req *s3_lifecycle_pb.LifecycleDeleteRequest) (*s3_lifecycle_pb.LifecycleDeleteResponse, error) {
	_, body, err := c.LifecycleDelete(EncodeRequest(req))
	if err != nil {
		return nil, err
	}
	return DecodeResponse(body)
}

// Dispatch is the ZAP handler contract for the server: decode the request,
// run handler (the pb-typed business method, with a fresh background context —
// the ZAP transport carries no deadline), encode the response.
func Dispatch(handler func(context.Context, *s3_lifecycle_pb.LifecycleDeleteRequest) (*s3_lifecycle_pb.LifecycleDeleteResponse, error), env []byte) ([]byte, error) {
	return s3_lifecyclewire.DispatchHanzoS3LifecycleInternal(serverHandler(handler), env)
}

type serverHandler func(context.Context, *s3_lifecycle_pb.LifecycleDeleteRequest) (*s3_lifecycle_pb.LifecycleDeleteResponse, error)

func (h serverHandler) LifecycleDelete(req []byte) ([]byte, error) {
	pbReq, err := DecodeRequest(req)
	if err != nil {
		return nil, err
	}
	resp, err := h(context.Background(), pbReq)
	if err != nil {
		return nil, err
	}
	return EncodeResponse(resp), nil
}
