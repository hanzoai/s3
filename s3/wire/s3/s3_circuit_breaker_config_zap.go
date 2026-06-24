// Code generated from s3.proto; DO NOT EDIT.

package s3wire

import (
	zap "github.com/zap-proto/go"
)

// --- S3CircuitBreakerConfigBucketsEntry ---
//
// Wire entry for one map<string, S3CircuitBreakerOptions> pair of
// S3CircuitBreakerConfig.buckets (proto map field 2). The map is encoded as a
// repeated list of these entry objects; build each entry with
// NewS3CircuitBreakerConfigBucketsEntry and pass the sub-buffers to
// S3CircuitBreakerConfigInput.Buckets.

const (
	s3CircuitBreakerConfigBucketsEntryKeyOff   = 0
	s3CircuitBreakerConfigBucketsEntryValueOff = 8
	s3CircuitBreakerConfigBucketsEntrySize     = 12
)

// S3CircuitBreakerConfigBucketsEntry is a zero-copy view into one encoded
// map entry.
type S3CircuitBreakerConfigBucketsEntry struct{ o zap.Object }

// WrapS3CircuitBreakerConfigBucketsEntry parses b and returns a typed view.
func WrapS3CircuitBreakerConfigBucketsEntry(b []byte) (S3CircuitBreakerConfigBucketsEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return S3CircuitBreakerConfigBucketsEntry{}, err
	}
	return S3CircuitBreakerConfigBucketsEntry{o: m.Root()}, nil
}

// Key reads the map key (string).
func (t S3CircuitBreakerConfigBucketsEntry) Key() string {
	return t.o.Text(s3CircuitBreakerConfigBucketsEntryKeyOff)
}

// Value reads the map value (S3CircuitBreakerOptions nested message).
func (t S3CircuitBreakerConfigBucketsEntry) Value() S3CircuitBreakerOptions {
	return S3CircuitBreakerOptions{o: t.o.Object(s3CircuitBreakerConfigBucketsEntryValueOff)}
}

// S3CircuitBreakerConfigBucketsEntryInput collects one map pair. Value takes a
// pre-built S3CircuitBreakerOptions sub-buffer (from
// NewS3CircuitBreakerOptions).
type S3CircuitBreakerConfigBucketsEntryInput struct {
	Key   string
	Value []byte
}

// NewS3CircuitBreakerConfigBucketsEntry builds one encoded map entry.
func NewS3CircuitBreakerConfigBucketsEntry(in S3CircuitBreakerConfigBucketsEntryInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(s3CircuitBreakerConfigBucketsEntrySize)
	ob.SetText(s3CircuitBreakerConfigBucketsEntryKeyOff, in.Key)
	embedObject(b, ob, s3CircuitBreakerConfigBucketsEntryValueOff, in.Value)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- S3CircuitBreakerConfig ---

const (
	s3CircuitBreakerConfigGlobalOff  = 0
	s3CircuitBreakerConfigBucketsOff = 4
	s3CircuitBreakerConfigSize       = 12
)

// S3CircuitBreakerConfig is a zero-copy view into a ZAP-encoded
// S3CircuitBreakerConfig message.
type S3CircuitBreakerConfig struct{ o zap.Object }

// WrapS3CircuitBreakerConfig parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapS3CircuitBreakerConfig(b []byte) (S3CircuitBreakerConfig, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return S3CircuitBreakerConfig{}, err
	}
	return S3CircuitBreakerConfig{o: m.Root()}, nil
}

// Global reads the global field (proto field 1, S3CircuitBreakerOptions nested
// message).
func (t S3CircuitBreakerConfig) Global() S3CircuitBreakerOptions {
	return S3CircuitBreakerOptions{o: t.o.Object(s3CircuitBreakerConfigGlobalOff)}
}

// Buckets reads the buckets field (proto field 2, map<string,
// S3CircuitBreakerOptions>). Each element is an out-of-line
// S3CircuitBreakerConfigBucketsEntry object; read element i with
// List.ObjectAt(i) and wrap it.
func (t S3CircuitBreakerConfig) Buckets() zap.List {
	return t.o.List(s3CircuitBreakerConfigBucketsOff)
}

// S3CircuitBreakerConfigInput collects the field values for
// NewS3CircuitBreakerConfig. Global takes a pre-built S3CircuitBreakerOptions
// sub-buffer; Buckets takes one pre-built entry sub-buffer (from
// NewS3CircuitBreakerConfigBucketsEntry) per map pair.
type S3CircuitBreakerConfigInput struct {
	Global  []byte
	Buckets [][]byte
}

// NewS3CircuitBreakerConfig builds a ZAP-encoded S3CircuitBreakerConfig message
// from in and returns the bytes.
func NewS3CircuitBreakerConfig(in S3CircuitBreakerConfigInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(s3CircuitBreakerConfigSize)
	embedObject(b, ob, s3CircuitBreakerConfigGlobalOff, in.Global)
	bucketsLB := b.StartList(0)
	for _, elem := range in.Buckets {
		bucketsLB.AddObjectBytes(elem)
	}
	ob.SetList(s3CircuitBreakerConfigBucketsOff, bucketsLB.FinishOffset(), len(in.Buckets))
	ob.FinishAsRoot()
	return b.Finish()
}
