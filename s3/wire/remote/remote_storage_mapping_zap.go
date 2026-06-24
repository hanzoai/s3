// Code generated from remote.proto; DO NOT EDIT.

package remotewire

import (
	zap "github.com/zap-proto/go"
)

// --- RemoteStorageMapping ---
//
// RemoteStorageMapping mirrors remote_pb.RemoteStorageMapping:
//
//	map<string, RemoteStorageLocation> mappings           = 1;
//	string                             primary_bucket_... = 2;
//
// The proto map is modeled per the ZAP field-shape rule for map<K,V>: a list of
// {key,value} entry objects. Each entry (RemoteStorageMappingEntry) is an
// out-of-line element carrying its own tail, so the list uses the
// variable-element encoding (StartList(0) + AddObjectBytes / List.ObjectAt).
//
// The entry's value is the singular message RemoteStorageLocation. It is held
// as a self-contained ZAP sub-buffer (an out-of-line object that keeps its own
// tail), re-parsed on read via WrapRemoteStorageLocation — the same shape
// List.ObjectAt expects for nested-struct list elements.

// --- RemoteStorageMappingEntry (synthetic map<string,RemoteStorageLocation> entry) ---

const (
	remoteStorageMappingEntryKeyOff   = 0 // string  (proto map key, field 1)
	remoteStorageMappingEntryValueOff = 8 // bytes   (proto map value, field 2: a RemoteStorageLocation sub-buffer)
	remoteStorageMappingEntrySize     = 16
)

// RemoteStorageMappingEntry is a zero-copy view into one map entry.
type RemoteStorageMappingEntry struct{ o zap.Object }

// Key reads the entry key (proto map key, string).
func (t RemoteStorageMappingEntry) Key() string { return t.o.Text(remoteStorageMappingEntryKeyOff) }

// ValueBytes returns the raw RemoteStorageLocation sub-buffer for this entry's
// value (zero-copy; aliases the parent buffer).
func (t RemoteStorageMappingEntry) ValueBytes() []byte {
	return t.o.Bytes(remoteStorageMappingEntryValueOff)
}

// Value parses the entry value into a typed RemoteStorageLocation view.
func (t RemoteStorageMappingEntry) Value() (RemoteStorageLocation, error) {
	return WrapRemoteStorageLocation(t.ValueBytes())
}

// RemoteStorageMappingEntryInput collects one map entry for NewRemoteStorageMapping.
type RemoteStorageMappingEntryInput struct {
	Key   string
	Value RemoteStorageLocationInput
}

// --- RemoteStorageMapping ---

const (
	remoteStorageMappingMappingsOff                 = 0 // list  (8B: relOffset + length)
	remoteStorageMappingPrimaryBucketStorageNameOff = 8 // string (proto field 2)
	remoteStorageMappingSize                        = 16
)

// RemoteStorageMapping is a zero-copy view into a ZAP-encoded
// RemoteStorageMapping message.
type RemoteStorageMapping struct{ o zap.Object }

// WrapRemoteStorageMapping parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapRemoteStorageMapping(b []byte) (RemoteStorageMapping, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RemoteStorageMapping{}, err
	}
	return RemoteStorageMapping{o: m.Root()}, nil
}

// PrimaryBucketStorageName reads the primary_bucket_storage_name field (proto
// field 2, string).
func (t RemoteStorageMapping) PrimaryBucketStorageName() string {
	return t.o.Text(remoteStorageMappingPrimaryBucketStorageNameOff)
}

// MappingsLen returns the number of map entries.
func (t RemoteStorageMapping) MappingsLen() int {
	return t.o.List(remoteStorageMappingMappingsOff).Len()
}

// MappingsAt returns the i-th map entry as a RemoteStorageMappingEntry view.
// Returns a zero entry if i is out of range or the element fails to parse.
func (t RemoteStorageMapping) MappingsAt(i int) RemoteStorageMappingEntry {
	return RemoteStorageMappingEntry{o: t.o.List(remoteStorageMappingMappingsOff).ObjectAt(i)}
}

// EachMapping walks the map entries once (O(N)), calling fn with the index and
// each entry view in order. fn returns true to continue, false to stop early.
func (t RemoteStorageMapping) EachMapping(fn func(i int, e RemoteStorageMappingEntry) bool) {
	t.o.List(remoteStorageMappingMappingsOff).Each(func(i int, elem zap.Object) bool {
		return fn(i, RemoteStorageMappingEntry{o: elem})
	})
}

// RemoteStorageMappingInput collects the field values for NewRemoteStorageMapping.
type RemoteStorageMappingInput struct {
	Mappings                 []RemoteStorageMappingEntryInput
	PrimaryBucketStorageName string
}

// newRemoteStorageMappingEntry builds one self-contained map-entry sub-buffer:
// the key string and the value RemoteStorageLocation (itself a self-contained
// sub-buffer held as bytes).
func newRemoteStorageMappingEntry(in RemoteStorageMappingEntryInput) []byte {
	value := NewRemoteStorageLocation(in.Value)
	b := zap.NewBuilder(256)
	ob := b.StartObject(remoteStorageMappingEntrySize)
	ob.SetText(remoteStorageMappingEntryKeyOff, in.Key)
	ob.SetBytes(remoteStorageMappingEntryValueOff, value)
	ob.FinishAsRoot()
	return b.Finish()
}

// NewRemoteStorageMapping builds a ZAP-encoded RemoteStorageMapping message from
// in and returns the bytes. The mappings list is built first (each entry is an
// out-of-line sub-buffer appended via AddObjectBytes), then its offset/length
// are stored in the root object alongside the primary-bucket string.
func NewRemoteStorageMapping(in RemoteStorageMappingInput) []byte {
	b := zap.NewBuilder(512)

	// Build the variable-element list of map entries first so its bytes land in
	// the data segment ahead of the root object's pointer cell.
	lb := b.StartList(0)
	for i := range in.Mappings {
		lb.AddObjectBytes(newRemoteStorageMappingEntry(in.Mappings[i]))
	}
	listOff, listLen := lb.Finish()

	ob := b.StartObject(remoteStorageMappingSize)
	ob.SetList(remoteStorageMappingMappingsOff, listOff, listLen)
	ob.SetText(remoteStorageMappingPrimaryBucketStorageNameOff, in.PrimaryBucketStorageName)
	ob.FinishAsRoot()
	return b.Finish()
}
