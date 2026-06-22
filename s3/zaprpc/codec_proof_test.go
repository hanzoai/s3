package zaprpc

import (
	"testing"

	"github.com/luxfi/zap"
)

// proofRemoteStorageLocation mirrors remote_pb.RemoteStorageLocation:
//
//	string name = 1; string bucket = 2; string path = 3;
//	int32  listing_cache_ttl_seconds = 4;
//
// It is the hand-written reference for what cmd/protoc-gen-zap must emit:
// plain struct + MarshalZap/UnmarshalZap over the luxfi/zap wire format. The
// fixed object payload assigns each field a byte offset (var-len string fields
// take an 8-byte {relOffset,length} tail pointer; the int32 sits inline).
type proofRemoteStorageLocation struct {
	Name                    string
	Bucket                  string
	Path                    string
	ListingCacheTtlSeconds  int32
}

// Field offsets in the fixed object payload, assigned by the (future) codegen.
const (
	rslOffName   = 0  // string  -> 8 bytes (relOffset u32 + length u32)
	rslOffBucket = 8  // string  -> 8 bytes
	rslOffPath   = 16 // string  -> 8 bytes
	rslOffTTL    = 24 // int32   -> 4 bytes
	rslDataSize  = 28
)

func (m *proofRemoteStorageLocation) MarshalZap() []byte {
	b := zap.NewBuilder(64)
	ob := b.StartObject(rslDataSize)
	ob.SetText(rslOffName, m.Name)
	ob.SetText(rslOffBucket, m.Bucket)
	ob.SetText(rslOffPath, m.Path)
	ob.SetInt32(rslOffTTL, m.ListingCacheTtlSeconds)
	ob.FinishAsRoot()
	return b.Finish()
}

func (m *proofRemoteStorageLocation) UnmarshalZap(data []byte) error {
	msg, err := zap.Parse(data)
	if err != nil {
		return err
	}
	o := msg.Root()
	m.Name = o.Text(rslOffName)
	m.Bucket = o.Text(rslOffBucket)
	m.Path = o.Text(rslOffPath)
	m.ListingCacheTtlSeconds = o.Int32(rslOffTTL)
	return nil
}

func TestZapCodecRoundTrip(t *testing.T) {
	cases := []proofRemoteStorageLocation{
		{Name: "wasabi", Bucket: "hanzo-cold", Path: "/vfs/blocks", ListingCacheTtlSeconds: 300},
		{Name: "", Bucket: "", Path: "", ListingCacheTtlSeconds: 0},                       // all-zero
		{Name: "s3", Bucket: "b", Path: "/", ListingCacheTtlSeconds: -1},                  // negative int32
		{Name: "long-name-" + string(make([]byte, 200)), Bucket: "x", Path: "y", ListingCacheTtlSeconds: 1 << 30}, // long string + big int
	}
	for i, in := range cases {
		wire := in.MarshalZap()
		var out proofRemoteStorageLocation
		if err := out.UnmarshalZap(wire); err != nil {
			t.Fatalf("case %d: UnmarshalZap: %v", i, err)
		}
		if out != in {
			t.Fatalf("case %d: round-trip mismatch\n in: %+v\nout: %+v\nwire: %d bytes", i, in, out, len(wire))
		}
	}
}
