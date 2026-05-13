// Package client — canonical client interface for Hanzo Storage.
//
//	import storage "github.com/hanzoai/storage/client"
//	var s storage.ObjectStore = hanzostorage.New(cfg)

package client

import (
	"context"
	"errors"
	"io"
	"time"
)

// Typed errors. Callers branch via errors.Is.
var (
	ErrObjectNotFound      = errors.New("storage: object not found")
	ErrBucketNotFound      = errors.New("storage: bucket not found")
	ErrPresignExpiryTooLong = errors.New("storage: presign expiry exceeds policy cap")
	ErrContentTypeMismatch = errors.New("storage: content type mismatch")
)

// ObjectStore is the bucket-scoped object-storage surface. Buckets
// are bound to org scope at construction; the bucket name visible
// here is the in-tenant leaf, not a globally unique name.
type ObjectStore interface {
	// Kind reports the backend identifier
	// (hanzo-storage | aws-s3 | gcs | r2 | b2 | memory).
	Kind() string

	// Put uploads an object. size=-1 selects multipart.
	Put(ctx context.Context, bucket, key string, r io.Reader, size int64, opts PutOpts) (*ObjectInfo, error)

	// Get streams the object body. Caller closes the reader.
	// ErrObjectNotFound on miss.
	Get(ctx context.Context, bucket, key string) (io.ReadCloser, *ObjectInfo, error)

	// Stat returns metadata without reading the body.
	Stat(ctx context.Context, bucket, key string) (*ObjectInfo, error)

	// Delete removes the object. Idempotent.
	Delete(ctx context.Context, bucket, key string) error

	// List paginates objects under prefix. Cursor-based.
	List(ctx context.Context, bucket, prefix string, opts ListOpts) (*ListPage, error)

	// Presign returns a time-limited URL the FE can use without
	// proxying through the server. PresignOpts narrow the URL's
	// authority — content-type pin, max-size, IP allowlist, max-uses.
	// Backends without per-condition support return
	// ErrPresignUnsupported.
	Presign(ctx context.Context, bucket, key string, opts PresignOpts) (string, error)
}

// PutOpts configures an upload. All fields optional.
type PutOpts struct {
	ContentType string
	// Encrypt requests server-side encryption with managed keys.
	Encrypt bool
	// SSEAlg + SSEKey override the default encryption.
	SSEAlg   string
	SSEKey   string
	Metadata map[string]string
	Tags     map[string]string
}

// ListOpts is the cursor-pagination shape.
type ListOpts struct {
	Cursor      string
	Limit       int // implementations cap at 1000
	IncludeDirs bool // pseudo-folders (CommonPrefixes)
}

// ListPage is one slice of a List call.
type ListPage struct {
	Items      []ObjectInfo
	NextCursor string
}

// PresignOpts narrows the authority of a presigned URL.
//
// Method is required (GET | PUT | DELETE). Expiry is required; the
// implementation refuses values above its policy cap with
// ErrPresignExpiryTooLong (typical caps: GET=24h, PUT=15m, DELETE=5m).
//
// All other fields are optional but strongly recommended:
//   - ContentTypeMustEqual pins the upload's Content-Type header.
//     Mismatch on upload returns S3-equivalent 403.
//   - MaxContentLength caps the upload byte size at the backend.
//   - AccessFromCIDR restricts the source IP that may use the URL.
//   - MaxUses limits how many times the URL may be redeemed.
type PresignOpts struct {
	Method               string // GET | PUT | DELETE
	Expiry               time.Duration
	ContentTypeMustEqual string
	MaxContentLength     int64
	AccessFromCIDR       []string
	MaxUses              int
	// ResponseHeaders inject response-header overrides on a GET.
	ResponseHeaders map[string]string
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
