// Code generated from master.proto; DO NOT EDIT.

package masterwire

import (
	zap "github.com/zap-proto/go"
)

// --- VolumeInformationMessage ---

const (
	volumeInformationMessageIdOff                = 0
	volumeInformationMessageSizeOff              = 4
	volumeInformationMessageCollectionOff        = 12
	volumeInformationMessageFileCountOff         = 20
	volumeInformationMessageDeleteCountOff       = 28
	volumeInformationMessageDeletedByteCountOff  = 36
	volumeInformationMessageReadOnlyOff          = 44
	volumeInformationMessageReplicaPlacementOff  = 45
	volumeInformationMessageVersionOff           = 49
	volumeInformationMessageTtlOff               = 53
	volumeInformationMessageCompactRevisionOff   = 57
	volumeInformationMessageModifiedAtSecondOff  = 61
	volumeInformationMessageRemoteStorageNameOff = 69
	volumeInformationMessageRemoteStorageKeyOff  = 77
	volumeInformationMessageDiskTypeOff          = 85
	volumeInformationMessageDiskIdOff            = 93
	volumeInformationMessageSize                 = 97
)

// VolumeInformationMessage is a zero-copy view into a ZAP-encoded
// VolumeInformationMessage message.
type VolumeInformationMessage struct{ o zap.Object }

// WrapVolumeInformationMessage parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeInformationMessage(b []byte) (VolumeInformationMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeInformationMessage{}, err
	}
	return VolumeInformationMessage{o: m.Root()}, nil
}

// Id reads the id field (proto field 1, uint32).
func (t VolumeInformationMessage) Id() uint32 { return t.o.Uint32(volumeInformationMessageIdOff) }

// Size reads the size field (proto field 2, uint64).
func (t VolumeInformationMessage) Size() uint64 {
	return t.o.Uint64(volumeInformationMessageSizeOff)
}

// Collection reads the collection field (proto field 3, string).
func (t VolumeInformationMessage) Collection() string {
	return t.o.Text(volumeInformationMessageCollectionOff)
}

// FileCount reads the file_count field (proto field 4, uint64).
func (t VolumeInformationMessage) FileCount() uint64 {
	return t.o.Uint64(volumeInformationMessageFileCountOff)
}

// DeleteCount reads the delete_count field (proto field 5, uint64).
func (t VolumeInformationMessage) DeleteCount() uint64 {
	return t.o.Uint64(volumeInformationMessageDeleteCountOff)
}

// DeletedByteCount reads the deleted_byte_count field (proto field 6, uint64).
func (t VolumeInformationMessage) DeletedByteCount() uint64 {
	return t.o.Uint64(volumeInformationMessageDeletedByteCountOff)
}

// ReadOnly reads the read_only field (proto field 7, bool).
func (t VolumeInformationMessage) ReadOnly() bool {
	return t.o.Bool(volumeInformationMessageReadOnlyOff)
}

// ReplicaPlacement reads the replica_placement field (proto field 8, uint32).
func (t VolumeInformationMessage) ReplicaPlacement() uint32 {
	return t.o.Uint32(volumeInformationMessageReplicaPlacementOff)
}

// Version reads the version field (proto field 9, uint32).
func (t VolumeInformationMessage) Version() uint32 {
	return t.o.Uint32(volumeInformationMessageVersionOff)
}

// Ttl reads the ttl field (proto field 10, uint32).
func (t VolumeInformationMessage) Ttl() uint32 { return t.o.Uint32(volumeInformationMessageTtlOff) }

// CompactRevision reads the compact_revision field (proto field 11, uint32).
func (t VolumeInformationMessage) CompactRevision() uint32 {
	return t.o.Uint32(volumeInformationMessageCompactRevisionOff)
}

// ModifiedAtSecond reads the modified_at_second field (proto field 12, int64).
func (t VolumeInformationMessage) ModifiedAtSecond() int64 {
	return t.o.Int64(volumeInformationMessageModifiedAtSecondOff)
}

// RemoteStorageName reads the remote_storage_name field (proto field 13, string).
func (t VolumeInformationMessage) RemoteStorageName() string {
	return t.o.Text(volumeInformationMessageRemoteStorageNameOff)
}

// RemoteStorageKey reads the remote_storage_key field (proto field 14, string).
func (t VolumeInformationMessage) RemoteStorageKey() string {
	return t.o.Text(volumeInformationMessageRemoteStorageKeyOff)
}

// DiskType reads the disk_type field (proto field 15, string).
func (t VolumeInformationMessage) DiskType() string {
	return t.o.Text(volumeInformationMessageDiskTypeOff)
}

// DiskId reads the disk_id field (proto field 16, uint32).
func (t VolumeInformationMessage) DiskId() uint32 {
	return t.o.Uint32(volumeInformationMessageDiskIdOff)
}

// VolumeInformationMessageInput collects the field values for
// NewVolumeInformationMessage.
type VolumeInformationMessageInput struct {
	Id                uint32
	Size              uint64
	Collection        string
	FileCount         uint64
	DeleteCount       uint64
	DeletedByteCount  uint64
	ReadOnly          bool
	ReplicaPlacement  uint32
	Version           uint32
	Ttl               uint32
	CompactRevision   uint32
	ModifiedAtSecond  int64
	RemoteStorageName string
	RemoteStorageKey  string
	DiskType          string
	DiskId            uint32
}

// NewVolumeInformationMessage builds a ZAP-encoded VolumeInformationMessage
// message from in and returns the bytes.
func NewVolumeInformationMessage(in VolumeInformationMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeInformationMessageSize)
	ob.SetUint32(volumeInformationMessageIdOff, in.Id)
	ob.SetUint64(volumeInformationMessageSizeOff, in.Size)
	ob.SetText(volumeInformationMessageCollectionOff, in.Collection)
	ob.SetUint64(volumeInformationMessageFileCountOff, in.FileCount)
	ob.SetUint64(volumeInformationMessageDeleteCountOff, in.DeleteCount)
	ob.SetUint64(volumeInformationMessageDeletedByteCountOff, in.DeletedByteCount)
	ob.SetBool(volumeInformationMessageReadOnlyOff, in.ReadOnly)
	ob.SetUint32(volumeInformationMessageReplicaPlacementOff, in.ReplicaPlacement)
	ob.SetUint32(volumeInformationMessageVersionOff, in.Version)
	ob.SetUint32(volumeInformationMessageTtlOff, in.Ttl)
	ob.SetUint32(volumeInformationMessageCompactRevisionOff, in.CompactRevision)
	ob.SetInt64(volumeInformationMessageModifiedAtSecondOff, in.ModifiedAtSecond)
	ob.SetText(volumeInformationMessageRemoteStorageNameOff, in.RemoteStorageName)
	ob.SetText(volumeInformationMessageRemoteStorageKeyOff, in.RemoteStorageKey)
	ob.SetText(volumeInformationMessageDiskTypeOff, in.DiskType)
	ob.SetUint32(volumeInformationMessageDiskIdOff, in.DiskId)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeShortInformationMessage ---

const (
	volumeShortInformationMessageIdOff               = 0
	volumeShortInformationMessageCollectionOff       = 4
	volumeShortInformationMessageReplicaPlacementOff = 12
	volumeShortInformationMessageVersionOff          = 16
	volumeShortInformationMessageTtlOff              = 20
	volumeShortInformationMessageDiskTypeOff         = 24
	volumeShortInformationMessageDiskIdOff           = 32
	volumeShortInformationMessageSize                = 36
)

// VolumeShortInformationMessage is a zero-copy view into a ZAP-encoded
// VolumeShortInformationMessage message.
type VolumeShortInformationMessage struct{ o zap.Object }

// WrapVolumeShortInformationMessage parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapVolumeShortInformationMessage(b []byte) (VolumeShortInformationMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeShortInformationMessage{}, err
	}
	return VolumeShortInformationMessage{o: m.Root()}, nil
}

// Id reads the id field (proto field 1, uint32).
func (t VolumeShortInformationMessage) Id() uint32 {
	return t.o.Uint32(volumeShortInformationMessageIdOff)
}

// Collection reads the collection field (proto field 3, string).
func (t VolumeShortInformationMessage) Collection() string {
	return t.o.Text(volumeShortInformationMessageCollectionOff)
}

// ReplicaPlacement reads the replica_placement field (proto field 8, uint32).
func (t VolumeShortInformationMessage) ReplicaPlacement() uint32 {
	return t.o.Uint32(volumeShortInformationMessageReplicaPlacementOff)
}

// Version reads the version field (proto field 9, uint32).
func (t VolumeShortInformationMessage) Version() uint32 {
	return t.o.Uint32(volumeShortInformationMessageVersionOff)
}

// Ttl reads the ttl field (proto field 10, uint32).
func (t VolumeShortInformationMessage) Ttl() uint32 {
	return t.o.Uint32(volumeShortInformationMessageTtlOff)
}

// DiskType reads the disk_type field (proto field 15, string).
func (t VolumeShortInformationMessage) DiskType() string {
	return t.o.Text(volumeShortInformationMessageDiskTypeOff)
}

// DiskId reads the disk_id field (proto field 16, uint32).
func (t VolumeShortInformationMessage) DiskId() uint32 {
	return t.o.Uint32(volumeShortInformationMessageDiskIdOff)
}

// VolumeShortInformationMessageInput collects the field values for
// NewVolumeShortInformationMessage.
type VolumeShortInformationMessageInput struct {
	Id               uint32
	Collection       string
	ReplicaPlacement uint32
	Version          uint32
	Ttl              uint32
	DiskType         string
	DiskId           uint32
}

// NewVolumeShortInformationMessage builds a ZAP-encoded
// VolumeShortInformationMessage message from in and returns the bytes.
func NewVolumeShortInformationMessage(in VolumeShortInformationMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeShortInformationMessageSize)
	ob.SetUint32(volumeShortInformationMessageIdOff, in.Id)
	ob.SetText(volumeShortInformationMessageCollectionOff, in.Collection)
	ob.SetUint32(volumeShortInformationMessageReplicaPlacementOff, in.ReplicaPlacement)
	ob.SetUint32(volumeShortInformationMessageVersionOff, in.Version)
	ob.SetUint32(volumeShortInformationMessageTtlOff, in.Ttl)
	ob.SetText(volumeShortInformationMessageDiskTypeOff, in.DiskType)
	ob.SetUint32(volumeShortInformationMessageDiskIdOff, in.DiskId)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardInformationMessage ---

const (
	volumeEcShardInformationMessageIdOff          = 0
	volumeEcShardInformationMessageCollectionOff  = 4
	volumeEcShardInformationMessageEcIndexBitsOff = 12
	volumeEcShardInformationMessageDiskTypeOff    = 16
	volumeEcShardInformationMessageExpireAtSecOff = 24
	volumeEcShardInformationMessageDiskIdOff      = 32
	volumeEcShardInformationMessageShardSizesOff  = 36
	volumeEcShardInformationMessageFileCountOff   = 44
	volumeEcShardInformationMessageDeleteCountOff = 52
	volumeEcShardInformationMessageEncodeTsNsOff  = 60
	volumeEcShardInformationMessageSize           = 68
)

// VolumeEcShardInformationMessage is a zero-copy view into a ZAP-encoded
// VolumeEcShardInformationMessage message.
type VolumeEcShardInformationMessage struct{ o zap.Object }

// WrapVolumeEcShardInformationMessage parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardInformationMessage(b []byte) (VolumeEcShardInformationMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardInformationMessage{}, err
	}
	return VolumeEcShardInformationMessage{o: m.Root()}, nil
}

// Id reads the id field (proto field 1, uint32).
func (t VolumeEcShardInformationMessage) Id() uint32 {
	return t.o.Uint32(volumeEcShardInformationMessageIdOff)
}

// Collection reads the collection field (proto field 2, string).
func (t VolumeEcShardInformationMessage) Collection() string {
	return t.o.Text(volumeEcShardInformationMessageCollectionOff)
}

// EcIndexBits reads the ec_index_bits field (proto field 3, uint32).
func (t VolumeEcShardInformationMessage) EcIndexBits() uint32 {
	return t.o.Uint32(volumeEcShardInformationMessageEcIndexBitsOff)
}

// DiskType reads the disk_type field (proto field 4, string).
func (t VolumeEcShardInformationMessage) DiskType() string {
	return t.o.Text(volumeEcShardInformationMessageDiskTypeOff)
}

// ExpireAtSec reads the expire_at_sec field (proto field 5, uint64).
func (t VolumeEcShardInformationMessage) ExpireAtSec() uint64 {
	return t.o.Uint64(volumeEcShardInformationMessageExpireAtSecOff)
}

// DiskId reads the disk_id field (proto field 6, uint32).
func (t VolumeEcShardInformationMessage) DiskId() uint32 {
	return t.o.Uint32(volumeEcShardInformationMessageDiskIdOff)
}

// ShardSizesLen reports the number of shard_sizes elements (proto field 7,
// repeated int64).
func (t VolumeEcShardInformationMessage) ShardSizesLen() int {
	return t.o.List(volumeEcShardInformationMessageShardSizesOff).Len()
}

// ShardSizeAt returns the i-th shard_sizes element.
func (t VolumeEcShardInformationMessage) ShardSizeAt(i int) int64 {
	return int64At(t.o.List(volumeEcShardInformationMessageShardSizesOff), i)
}

// FileCount reads the file_count field (proto field 8, uint64).
func (t VolumeEcShardInformationMessage) FileCount() uint64 {
	return t.o.Uint64(volumeEcShardInformationMessageFileCountOff)
}

// DeleteCount reads the delete_count field (proto field 9, uint64).
func (t VolumeEcShardInformationMessage) DeleteCount() uint64 {
	return t.o.Uint64(volumeEcShardInformationMessageDeleteCountOff)
}

// EncodeTsNs reads the encode_ts_ns field (proto field 14, int64).
func (t VolumeEcShardInformationMessage) EncodeTsNs() int64 {
	return t.o.Int64(volumeEcShardInformationMessageEncodeTsNsOff)
}

// VolumeEcShardInformationMessageInput collects the field values for
// NewVolumeEcShardInformationMessage.
type VolumeEcShardInformationMessageInput struct {
	Id          uint32
	Collection  string
	EcIndexBits uint32
	DiskType    string
	ExpireAtSec uint64
	DiskId      uint32
	ShardSizes  []int64
	FileCount   uint64
	DeleteCount uint64
	EncodeTsNs  int64
}

// NewVolumeEcShardInformationMessage builds a ZAP-encoded
// VolumeEcShardInformationMessage message from in and returns the bytes.
func NewVolumeEcShardInformationMessage(in VolumeEcShardInformationMessageInput) []byte {
	b := zap.NewBuilder(256)
	shardSizesOff, shardSizesLen := setInt64List(b, in.ShardSizes)
	ob := b.StartObject(volumeEcShardInformationMessageSize)
	ob.SetUint32(volumeEcShardInformationMessageIdOff, in.Id)
	ob.SetText(volumeEcShardInformationMessageCollectionOff, in.Collection)
	ob.SetUint32(volumeEcShardInformationMessageEcIndexBitsOff, in.EcIndexBits)
	ob.SetText(volumeEcShardInformationMessageDiskTypeOff, in.DiskType)
	ob.SetUint64(volumeEcShardInformationMessageExpireAtSecOff, in.ExpireAtSec)
	ob.SetUint32(volumeEcShardInformationMessageDiskIdOff, in.DiskId)
	ob.SetList(volumeEcShardInformationMessageShardSizesOff, shardSizesOff, shardSizesLen)
	ob.SetUint64(volumeEcShardInformationMessageFileCountOff, in.FileCount)
	ob.SetUint64(volumeEcShardInformationMessageDeleteCountOff, in.DeleteCount)
	ob.SetInt64(volumeEcShardInformationMessageEncodeTsNsOff, in.EncodeTsNs)
	ob.FinishAsRoot()
	return b.Finish()
}
