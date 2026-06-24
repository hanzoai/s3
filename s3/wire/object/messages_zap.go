// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// Zero-copy ZAP wire messages for the HanzoS3Object service, on the canonical
// github.com/zap-proto/go runtime (same lineage as the iam/filer/master wire
// packages). Each variable-length field (string/bytes) occupies an 8-byte tail
// pointer slot in the fixed object region; offsets run in declaration order.
// The bytes ARE the message — no struct, no marshal, no protobuf.

package objectwire

import (
	zap "github.com/zap-proto/go"
)

// --- GetObjectRequest ---

const (
	getObjectRequestBucketOff = 0
	getObjectRequestKeyOff    = 8
	getObjectRequestSize      = 16
)

// GetObjectRequest is a zero-copy view into a ZAP-encoded GetObjectRequest.
type GetObjectRequest struct{ o zap.Object }

// WrapGetObjectRequest parses b and returns a typed view.
func WrapGetObjectRequest(b []byte) (GetObjectRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetObjectRequest{}, err
	}
	return GetObjectRequest{o: m.Root()}, nil
}

func (t GetObjectRequest) Bucket() string { return t.o.Text(getObjectRequestBucketOff) }
func (t GetObjectRequest) Key() string    { return t.o.Text(getObjectRequestKeyOff) }

// GetObjectRequestInput collects the field values for NewGetObjectRequest.
type GetObjectRequestInput struct {
	Bucket string
	Key    string
}

// NewGetObjectRequest builds a ZAP-encoded GetObjectRequest and returns the bytes.
func NewGetObjectRequest(in GetObjectRequestInput) []byte {
	b := zap.NewBuilder(128)
	ob := b.StartObject(getObjectRequestSize)
	ob.SetText(getObjectRequestBucketOff, in.Bucket)
	ob.SetText(getObjectRequestKeyOff, in.Key)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetObjectResponse ---

const (
	getObjectResponseDataOff        = 0
	getObjectResponseContentTypeOff = 8
	getObjectResponseEtagOff        = 16
	getObjectResponseSize           = 24
)

// GetObjectResponse is a zero-copy view into a ZAP-encoded GetObjectResponse.
type GetObjectResponse struct{ o zap.Object }

// WrapGetObjectResponse parses b and returns a typed view.
func WrapGetObjectResponse(b []byte) (GetObjectResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetObjectResponse{}, err
	}
	return GetObjectResponse{o: m.Root()}, nil
}

func (t GetObjectResponse) Data() []byte        { return t.o.Bytes(getObjectResponseDataOff) }
func (t GetObjectResponse) ContentType() string { return t.o.Text(getObjectResponseContentTypeOff) }
func (t GetObjectResponse) Etag() string        { return t.o.Text(getObjectResponseEtagOff) }

// GetObjectResponseInput collects the field values for NewGetObjectResponse.
type GetObjectResponseInput struct {
	Data        []byte
	ContentType string
	Etag        string
}

// NewGetObjectResponse builds a ZAP-encoded GetObjectResponse and returns the bytes.
func NewGetObjectResponse(in GetObjectResponseInput) []byte {
	b := zap.NewBuilder(len(in.Data) + 128)
	ob := b.StartObject(getObjectResponseSize)
	ob.SetBytes(getObjectResponseDataOff, in.Data)
	ob.SetText(getObjectResponseContentTypeOff, in.ContentType)
	ob.SetText(getObjectResponseEtagOff, in.Etag)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PutObjectRequest ---

const (
	putObjectRequestBucketOff      = 0
	putObjectRequestKeyOff         = 8
	putObjectRequestDataOff        = 16
	putObjectRequestContentTypeOff = 24
	putObjectRequestSize           = 32
)

// PutObjectRequest is a zero-copy view into a ZAP-encoded PutObjectRequest.
type PutObjectRequest struct{ o zap.Object }

// WrapPutObjectRequest parses b and returns a typed view.
func WrapPutObjectRequest(b []byte) (PutObjectRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PutObjectRequest{}, err
	}
	return PutObjectRequest{o: m.Root()}, nil
}

func (t PutObjectRequest) Bucket() string      { return t.o.Text(putObjectRequestBucketOff) }
func (t PutObjectRequest) Key() string         { return t.o.Text(putObjectRequestKeyOff) }
func (t PutObjectRequest) Data() []byte        { return t.o.Bytes(putObjectRequestDataOff) }
func (t PutObjectRequest) ContentType() string { return t.o.Text(putObjectRequestContentTypeOff) }

// PutObjectRequestInput collects the field values for NewPutObjectRequest.
type PutObjectRequestInput struct {
	Bucket      string
	Key         string
	Data        []byte
	ContentType string
}

// NewPutObjectRequest builds a ZAP-encoded PutObjectRequest and returns the bytes.
func NewPutObjectRequest(in PutObjectRequestInput) []byte {
	b := zap.NewBuilder(len(in.Data) + 128)
	ob := b.StartObject(putObjectRequestSize)
	ob.SetText(putObjectRequestBucketOff, in.Bucket)
	ob.SetText(putObjectRequestKeyOff, in.Key)
	ob.SetBytes(putObjectRequestDataOff, in.Data)
	ob.SetText(putObjectRequestContentTypeOff, in.ContentType)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PutObjectResponse ---

const (
	putObjectResponseEtagOff = 0
	putObjectResponseSize    = 8
)

// PutObjectResponse is a zero-copy view into a ZAP-encoded PutObjectResponse.
type PutObjectResponse struct{ o zap.Object }

// WrapPutObjectResponse parses b and returns a typed view.
func WrapPutObjectResponse(b []byte) (PutObjectResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PutObjectResponse{}, err
	}
	return PutObjectResponse{o: m.Root()}, nil
}

func (t PutObjectResponse) Etag() string { return t.o.Text(putObjectResponseEtagOff) }

// PutObjectResponseInput collects the field values for NewPutObjectResponse.
type PutObjectResponseInput struct {
	Etag string
}

// NewPutObjectResponse builds a ZAP-encoded PutObjectResponse and returns the bytes.
func NewPutObjectResponse(in PutObjectResponseInput) []byte {
	b := zap.NewBuilder(128)
	ob := b.StartObject(putObjectResponseSize)
	ob.SetText(putObjectResponseEtagOff, in.Etag)
	ob.FinishAsRoot()
	return b.Finish()
}
