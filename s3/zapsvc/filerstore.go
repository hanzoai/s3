package zapsvc

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"time"

	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/filer_pb"
)

// objectRoot is where S3 objects live in the filer namespace:
// /buckets/<bucket>/<key>. Matches the filer layout the HTTP S3 API uses.
const objectRoot = "/buckets"

// FilerStore is the real ObjectStore backing the native ZAP S3 service: it
// persists objects to the filer (content inline in the entry, content-type +
// etag in Extended metadata). This is the first real backend — small/medium
// objects via inline content; the chunked-to-volumes path for large objects is
// the production follow-on. The filer hop is still gRPC (the engine's internal
// transport) until the master↔volume↔filer rip; the SERVICE boundary is ZAP.
type FilerStore struct {
	filers     []pb.ServerAddress
	dialOption pb.DialOption
}

// NewFilerStore targets the given filer(s).
func NewFilerStore(filers []pb.ServerAddress, dialOption pb.DialOption) *FilerStore {
	return &FilerStore{filers: filers, dialOption: dialOption}
}

func objectDirName(bucket, key string) (dir, name string) {
	full := objectRoot + "/" + bucket + "/" + strings.TrimPrefix(key, "/")
	i := strings.LastIndexByte(full, '/')
	return full[:i], full[i+1:]
}

// Put stores an object and returns its etag (md5 hex, single-part S3 semantics).
func (f *FilerStore) Put(bucket, key string, data []byte, contentType string) (string, error) {
	sum := md5.Sum(data)
	etag := hex.EncodeToString(sum[:])
	dir, name := objectDirName(bucket, key)
	err := pb.WithOneOfGrpcFilerClients(false, f.filers, f.dialOption, func(client filer_pb.HanzoFilerClient) error {
		return filer_pb.CreateEntry(context.Background(), client, &filer_pb.CreateEntryRequest{
			Directory: dir,
			Entry: &filer_pb.Entry{
				Name:        name,
				IsDirectory: false,
				Attributes: &filer_pb.FuseAttributes{
					Mtime:    time.Now().Unix(),
					Crtime:   time.Now().Unix(),
					FileMode: uint32(0644),
					FileSize: uint64(len(data)),
				},
				Content: data,
				Extended: map[string][]byte{
					"Content-Type": []byte(contentType),
					"ETag":         []byte(etag),
				},
			},
			SkipCheckParentDirectory: true,
		})
	})
	return etag, err
}

// Get fetches an object's content, content-type, and etag.
func (f *FilerStore) Get(bucket, key string) (data []byte, contentType, etag string, err error) {
	dir, name := objectDirName(bucket, key)
	err = pb.WithOneOfGrpcFilerClients(false, f.filers, f.dialOption, func(client filer_pb.HanzoFilerClient) error {
		resp, e := filer_pb.LookupEntry(context.Background(), client, &filer_pb.LookupDirectoryEntryRequest{
			Directory: dir,
			Name:      name,
		})
		if e != nil {
			if e == filer_pb.ErrNotFound {
				return nil // missing object -> empty, no error (caller decides)
			}
			return e
		}
		data = resp.Entry.Content
		if ct := resp.Entry.Extended["Content-Type"]; len(ct) > 0 {
			contentType = string(ct)
		}
		if et := resp.Entry.Extended["ETag"]; len(et) > 0 {
			etag = string(et)
		}
		return nil
	})
	return
}
