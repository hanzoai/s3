// Code generated from s3.proto; DO NOT EDIT.

package s3wire

import (
	zap "github.com/zap-proto/go"
)

// --- BucketMetadataTagsEntry ---
//
// Wire entry for one map<string, string> pair of BucketMetadata.tags (proto map
// field 1). The map is encoded as a repeated list of these entry objects; build
// each entry with NewBucketMetadataTagsEntry and pass the sub-buffers to
// BucketMetadataInput.Tags.

const (
	bucketMetadataTagsEntryKeyOff   = 0
	bucketMetadataTagsEntryValueOff = 8
	bucketMetadataTagsEntrySize     = 16
)

// BucketMetadataTagsEntry is a zero-copy view into one encoded map entry.
type BucketMetadataTagsEntry struct{ o zap.Object }

// WrapBucketMetadataTagsEntry parses b and returns a typed view.
func WrapBucketMetadataTagsEntry(b []byte) (BucketMetadataTagsEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return BucketMetadataTagsEntry{}, err
	}
	return BucketMetadataTagsEntry{o: m.Root()}, nil
}

// Key reads the map key (string).
func (t BucketMetadataTagsEntry) Key() string {
	return t.o.Text(bucketMetadataTagsEntryKeyOff)
}

// Value reads the map value (string).
func (t BucketMetadataTagsEntry) Value() string {
	return t.o.Text(bucketMetadataTagsEntryValueOff)
}

// BucketMetadataTagsEntryInput collects one map pair.
type BucketMetadataTagsEntryInput struct {
	Key   string
	Value string
}

// NewBucketMetadataTagsEntry builds one encoded map entry.
func NewBucketMetadataTagsEntry(in BucketMetadataTagsEntryInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(bucketMetadataTagsEntrySize)
	ob.SetText(bucketMetadataTagsEntryKeyOff, in.Key)
	ob.SetText(bucketMetadataTagsEntryValueOff, in.Value)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- BucketMetadata ---

const (
	bucketMetadataTagsOff       = 0
	bucketMetadataCorsOff       = 8
	bucketMetadataEncryptionOff = 12
	bucketMetadataSize          = 16
)

// BucketMetadata is a zero-copy view into a ZAP-encoded BucketMetadata message.
type BucketMetadata struct{ o zap.Object }

// WrapBucketMetadata parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapBucketMetadata(b []byte) (BucketMetadata, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return BucketMetadata{}, err
	}
	return BucketMetadata{o: m.Root()}, nil
}

// Tags reads the tags field (proto field 1, map<string, string>). Each element
// is an out-of-line BucketMetadataTagsEntry object; read element i with
// List.ObjectAt(i) and wrap it.
func (t BucketMetadata) Tags() zap.List { return t.o.List(bucketMetadataTagsOff) }

// Cors reads the cors field (proto field 2, CORSConfiguration nested message).
func (t BucketMetadata) Cors() CORSConfiguration {
	return CORSConfiguration{o: t.o.Object(bucketMetadataCorsOff)}
}

// Encryption reads the encryption field (proto field 3, EncryptionConfiguration
// nested message).
func (t BucketMetadata) Encryption() EncryptionConfiguration {
	return EncryptionConfiguration{o: t.o.Object(bucketMetadataEncryptionOff)}
}

// BucketMetadataInput collects the field values for NewBucketMetadata. Tags
// takes one pre-built entry sub-buffer (from NewBucketMetadataTagsEntry) per map
// pair; Cors and Encryption take a pre-built nested sub-buffer each (from
// NewCORSConfiguration / NewEncryptionConfiguration).
type BucketMetadataInput struct {
	Tags       [][]byte
	Cors       []byte
	Encryption []byte
}

// NewBucketMetadata builds a ZAP-encoded BucketMetadata message from in and
// returns the bytes.
func NewBucketMetadata(in BucketMetadataInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(bucketMetadataSize)
	embedObject(b, ob, bucketMetadataCorsOff, in.Cors)
	embedObject(b, ob, bucketMetadataEncryptionOff, in.Encryption)
	tagsLB := b.StartList(0)
	for _, elem := range in.Tags {
		tagsLB.AddObjectBytes(elem)
	}
	ob.SetList(bucketMetadataTagsOff, tagsLB.FinishOffset(), len(in.Tags))
	ob.FinishAsRoot()
	return b.Finish()
}
