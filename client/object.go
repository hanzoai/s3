// Package client — canonical client interface for Hanzo Storage.
//
//	import storage "github.com/hanzoai/storage/client"
//	var s storage.ObjectStore = hanzostorage.New(cfg)

package client

import (
	"context"
	"io"
	"time"
)

// ObjectStore is the bucket-scoped object-storage surface.
type ObjectStore interface {
	// Kind reports the backend identifier
	// (hanzo-storage | aws-s3 | gcs | r2 | b2 | memory).
	Kind() string

	// Put uploads an object. size=-1 selects multipart for unknown lengths.
	Put(ctx context.Context, bucket, key string, r io.Reader, size int64, opts PutOpts) (*ObjectInfo, error)

	// Get streams the object body. Caller closes the reader.
	Get(ctx context.Context, bucket, key string) (io.ReadCloser, *ObjectInfo, error)

	// Stat returns metadata without reading the body.
	Stat(ctx context.Context, bucket, key string) (*ObjectInfo, error)

	// Delete removes the object. Idempotent.
	Delete(ctx context.Context, bucket, key string) error

	// List returns objects under prefix.
	List(ctx context.Context, bucket, prefix string, opts ListOpts) ([]ObjectInfo, error)

	// Presign returns a time-limited URL the FE can use to GET or PUT
	// directly without proxying through the server.
	Presign(ctx context.Context, bucket, key, method string, expiry time.Duration) (string, error)
}

// PutOpts configures an upload. All fields optional.
type PutOpts struct {
	ContentType string
	Encrypt     bool
	// SSEAlg + SSEKey override the default encryption.
	SSEAlg   string
	SSEKey   string
	Metadata map[string]string
	Tags     map[string]string
}

// ListOpts configures a list call.
type ListOpts struct {
	MaxResults  int
	StartAfter  string
	IncludeDirs bool // pseudo-folders (CommonPrefixes)
}

// ObjectInfo is the metadata sidecar.
type ObjectInfo struct {
	Bucket       string
	Key          string
	Size         int64
	ETag         string
	ContentType  string
	LastModified time.Time
	Metadata     map[string]string
	Tags         map[string]string
	VersionID    string
}
