// Package wire holds Hanzo S3's message types as ZAP schemas — zero-copy,
// protobuf-free, no _pb. These replace the generated *_pb packages. This file
// is the hand-written reference for the variable-length schema shape the
// extended v2/codegen will emit; once the generator handles strings/lists/
// nested it produces files exactly like this.
package wire

import (
	"github.com/luxfi/zap"
	zapv2 "github.com/luxfi/zap/v2"
)

// KindRemoteStorageLocation is the wire discriminator (byte at object offset 0).
const KindRemoteStorageLocation zapv2.KindByte = 0x01

// SizeRemoteStorageLocation is the fixed object payload: kind(1)+pad, ttl(4),
// and three string tail-pointers (8 each = {relOffset u32, length u32}). The
// string bytes themselves live in the tail after the fixed payload.
const SizeRemoteStorageLocation = 32

const (
	offRSLKind   = 0
	offRSLTTL    = 4
	offRSLName   = 8
	offRSLBucket = 16
	offRSLPath   = 24
)

// RemoteStorageLocation is the v2 schema marker (Kind/Size/Name fold to literals).
type RemoteStorageLocation struct{}

func (RemoteStorageLocation) Kind() zapv2.KindByte { return KindRemoteStorageLocation }
func (RemoteStorageLocation) Size() int            { return SizeRemoteStorageLocation }
func (RemoteStorageLocation) Name() string         { return "RemoteStorageLocation" }

// NewRemoteStorageLocation builds the message directly into a ZAP buffer — one
// allocation, no intermediate struct. Strings are written to the tail; the
// fixed payload holds their pointers.
func NewRemoteStorageLocation(name, bucket, path string, ttlSeconds int32) (zapv2.View[RemoteStorageLocation], []byte) {
	b := zap.NewBuilder(zap.HeaderSize + SizeRemoteStorageLocation + len(name) + len(bucket) + len(path))
	ob := b.StartObject(SizeRemoteStorageLocation)
	ob.SetUint8(offRSLKind, uint8(KindRemoteStorageLocation))
	ob.SetInt32(offRSLTTL, ttlSeconds)
	ob.SetText(offRSLName, name)
	ob.SetText(offRSLBucket, bucket)
	ob.SetText(offRSLPath, path)
	rootOff := ob.FinishAsRoot()
	buf := b.Finish()
	end := rootOff + SizeRemoteStorageLocation
	if end > len(buf) {
		end = len(buf)
	}
	return zapv2.AsView[RemoteStorageLocation](zapv2.RawFromSlices(buf, rootOff, end)), buf
}

// WrapRemoteStorageLocation validates the frame + kind and returns a zero-copy
// typed view over the bytes. No deserialization — the bytes ARE the message.
func WrapRemoteStorageLocation(b []byte) (zapv2.View[RemoteStorageLocation], error) {
	msg, err := zap.Parse(b)
	if err != nil {
		return zapv2.View[RemoteStorageLocation]{}, err
	}
	root := msg.Root()
	if got := root.Uint8(offRSLKind); got != uint8(KindRemoteStorageLocation) {
		return zapv2.View[RemoteStorageLocation]{}, zapv2.NewSchemaError(
			KindRemoteStorageLocation, zapv2.KindByte(got), "RemoteStorageLocation")
	}
	data := msg.Bytes()
	rootOff := root.Offset()
	end := rootOff + SizeRemoteStorageLocation
	if end > len(data) {
		end = len(data)
	}
	return zapv2.AsView[RemoteStorageLocation](zapv2.RawFromSlices(data, rootOff, end)), nil
}

// obj re-anchors a View on its underlying zap.Object for variable-length reads.
func rslObj(v zapv2.View[RemoteStorageLocation]) zap.Object {
	return zap.WrapBuffer(v.Bytes()).RootObjectAt(int(zapv2.RootOff(v)))
}

// Field accessors — zero-copy reads straight off the buffer.

func RemoteStorageLocationTTLSeconds(v zapv2.View[RemoteStorageLocation]) int32 {
	return rslObj(v).Int32(offRSLTTL)
}
func RemoteStorageLocationName(v zapv2.View[RemoteStorageLocation]) string {
	return rslObj(v).Text(offRSLName)
}
func RemoteStorageLocationBucket(v zapv2.View[RemoteStorageLocation]) string {
	return rslObj(v).Text(offRSLBucket)
}
func RemoteStorageLocationPath(v zapv2.View[RemoteStorageLocation]) string {
	return rslObj(v).Text(offRSLPath)
}
