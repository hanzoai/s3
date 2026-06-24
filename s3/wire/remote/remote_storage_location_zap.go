// Code generated from remote.proto; DO NOT EDIT.

package remotewire

import (
	zap "github.com/zap-proto/go"
)

// --- RemoteStorageLocation ---
//
// RemoteStorageLocation mirrors remote_pb.RemoteStorageLocation: a remote
// store coordinate (name/bucket/path) plus an on-demand listing-cache TTL.
//
// This is the zap-proto/go (zero-cgo) schema. The hand-written
// s3/wire/remote_storage_location_zap.go in the parent `wire` package is the
// github.com/luxfi/zap (typed-View) variant for that package's call path; both
// describe the same proto message, one per runtime.
//
// Wire layout: three string tail pointers (8B each) then the int32 TTL (4B).
const (
	remoteStorageLocationNameOff                   = 0
	remoteStorageLocationBucketOff                 = 8
	remoteStorageLocationPathOff                   = 16
	remoteStorageLocationListingCacheTTLSecondsOff = 24
	remoteStorageLocationSize                      = 28
)

// RemoteStorageLocation is a zero-copy view into a ZAP-encoded
// RemoteStorageLocation message.
type RemoteStorageLocation struct{ o zap.Object }

// WrapRemoteStorageLocation parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapRemoteStorageLocation(b []byte) (RemoteStorageLocation, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RemoteStorageLocation{}, err
	}
	return RemoteStorageLocation{o: m.Root()}, nil
}

func (t RemoteStorageLocation) Name() string   { return t.o.Text(remoteStorageLocationNameOff) }
func (t RemoteStorageLocation) Bucket() string { return t.o.Text(remoteStorageLocationBucketOff) }
func (t RemoteStorageLocation) Path() string   { return t.o.Text(remoteStorageLocationPathOff) }

// ListingCacheTTLSeconds reads the listing_cache_ttl_seconds field (proto field
// 4, int32). 0 = disabled; >0 enables on-demand directory listing with this TTL
// in seconds.
func (t RemoteStorageLocation) ListingCacheTTLSeconds() int32 {
	return t.o.Int32(remoteStorageLocationListingCacheTTLSecondsOff)
}

// RemoteStorageLocationInput collects the field values for
// NewRemoteStorageLocation.
type RemoteStorageLocationInput struct {
	Name                   string
	Bucket                 string
	Path                   string
	ListingCacheTTLSeconds int32
}

// NewRemoteStorageLocation builds a ZAP-encoded RemoteStorageLocation message
// from in and returns the bytes.
func NewRemoteStorageLocation(in RemoteStorageLocationInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(remoteStorageLocationSize)
	ob.SetText(remoteStorageLocationNameOff, in.Name)
	ob.SetText(remoteStorageLocationBucketOff, in.Bucket)
	ob.SetText(remoteStorageLocationPathOff, in.Path)
	ob.SetInt32(remoteStorageLocationListingCacheTTLSecondsOff, in.ListingCacheTTLSeconds)
	ob.FinishAsRoot()
	return b.Finish()
}
