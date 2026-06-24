// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
	zap "github.com/zap-proto/go"
)

// --- RemoteEntry ---

const (
	remoteEntryStorageNameOff       = 0
	remoteEntryLastLocalSyncTsNsOff = 8
	remoteEntryRemoteETagOff        = 16
	remoteEntryRemoteMtimeOff       = 24
	remoteEntryRemoteSizeOff        = 32
	remoteEntrySize                 = 40
)

// RemoteEntry is a zero-copy view into a ZAP-encoded RemoteEntry message.
type RemoteEntry struct{ o zap.Object }

// WrapRemoteEntry parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapRemoteEntry(b []byte) (RemoteEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RemoteEntry{}, err
	}
	return RemoteEntry{o: m.Root()}, nil
}

func (t RemoteEntry) StorageName() string      { return t.o.Text(remoteEntryStorageNameOff) }
func (t RemoteEntry) LastLocalSyncTsNs() int64 { return t.o.Int64(remoteEntryLastLocalSyncTsNsOff) }
func (t RemoteEntry) RemoteETag() string       { return t.o.Text(remoteEntryRemoteETagOff) }
func (t RemoteEntry) RemoteMtime() int64       { return t.o.Int64(remoteEntryRemoteMtimeOff) }
func (t RemoteEntry) RemoteSize() int64        { return t.o.Int64(remoteEntryRemoteSizeOff) }

// RemoteEntryInput collects the field values for NewRemoteEntry.
type RemoteEntryInput struct {
	StorageName       string
	LastLocalSyncTsNs int64
	RemoteETag        string
	RemoteMtime       int64
	RemoteSize        int64
}

// NewRemoteEntry builds a ZAP-encoded RemoteEntry message from in.
func NewRemoteEntry(in RemoteEntryInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(remoteEntrySize)
	ob.SetText(remoteEntryStorageNameOff, in.StorageName)
	ob.SetInt64(remoteEntryLastLocalSyncTsNsOff, in.LastLocalSyncTsNs)
	ob.SetText(remoteEntryRemoteETagOff, in.RemoteETag)
	ob.SetInt64(remoteEntryRemoteMtimeOff, in.RemoteMtime)
	ob.SetInt64(remoteEntryRemoteSizeOff, in.RemoteSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- FileId ---

const (
	fileIDVolumeIDOff = 0  // volume_id (field 1, uint32)
	fileIDFileKeyOff  = 4  // file_key (field 2, uint64)
	fileIDCookieOff   = 12 // cookie (field 3, fixed32)
	fileIDSize        = 16
)

// FileId is a zero-copy view into a ZAP-encoded FileId message.
type FileId struct{ o zap.Object }

// WrapFileId parses b and returns a typed view.
func WrapFileId(b []byte) (FileId, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FileId{}, err
	}
	return FileId{o: m.Root()}, nil
}

func (t FileId) VolumeID() uint32 { return t.o.Uint32(fileIDVolumeIDOff) }
func (t FileId) FileKey() uint64  { return t.o.Uint64(fileIDFileKeyOff) }
func (t FileId) Cookie() uint32   { return t.o.Uint32(fileIDCookieOff) }

// FileIdInput collects the field values for NewFileId.
type FileIdInput struct {
	VolumeID uint32
	FileKey  uint64
	Cookie   uint32
}

// NewFileId builds a ZAP-encoded FileId message from in.
func NewFileId(in FileIdInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(fileIDSize)
	ob.SetUint32(fileIDVolumeIDOff, in.VolumeID)
	ob.SetUint64(fileIDFileKeyOff, in.FileKey)
	ob.SetUint32(fileIDCookieOff, in.Cookie)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- FileChunk ---

const (
	fileChunkFileIDOff          = 0  // file_id (field 1, string)
	fileChunkOffsetOff          = 8  // offset (field 2, int64)
	fileChunkSizeFieldOff       = 16 // size (field 3, uint64)
	fileChunkModifiedTsNsOff    = 24 // modified_ts_ns (field 4, int64)
	fileChunkETagOff            = 32 // e_tag (field 5, string)
	fileChunkSourceFileIDOff    = 40 // source_file_id (field 6, string)
	fileChunkFidOff             = 48 // fid (field 7, FileId object)
	fileChunkSourceFidOff       = 52 // source_fid (field 8, FileId object)
	fileChunkCipherKeyOff       = 56 // cipher_key (field 9, bytes)
	fileChunkIsCompressedOff    = 64 // is_compressed (field 10, bool)
	fileChunkIsChunkManifestOff = 65 // is_chunk_manifest (field 11, bool)
	fileChunkSseTypeOff         = 68 // sse_type (field 12, enum SSEType)
	fileChunkSseMetadataOff     = 72 // sse_metadata (field 13, bytes)
	fileChunkSize               = 80
)

// FileChunk is a zero-copy view into a ZAP-encoded FileChunk message.
type FileChunk struct{ o zap.Object }

// WrapFileChunk parses b and returns a typed view.
func WrapFileChunk(b []byte) (FileChunk, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FileChunk{}, err
	}
	return FileChunk{o: m.Root()}, nil
}

func (t FileChunk) FileID() string       { return t.o.Text(fileChunkFileIDOff) }
func (t FileChunk) Offset() int64        { return t.o.Int64(fileChunkOffsetOff) }
func (t FileChunk) Size() uint64         { return t.o.Uint64(fileChunkSizeFieldOff) }
func (t FileChunk) ModifiedTsNs() int64  { return t.o.Int64(fileChunkModifiedTsNsOff) }
func (t FileChunk) ETag() string         { return t.o.Text(fileChunkETagOff) }
func (t FileChunk) SourceFileID() string { return t.o.Text(fileChunkSourceFileIDOff) }
func (t FileChunk) Fid() FileId          { return FileId{o: t.o.Object(fileChunkFidOff)} }
func (t FileChunk) SourceFid() FileId    { return FileId{o: t.o.Object(fileChunkSourceFidOff)} }
func (t FileChunk) CipherKey() []byte    { return t.o.Bytes(fileChunkCipherKeyOff) }
func (t FileChunk) IsCompressed() bool   { return t.o.Bool(fileChunkIsCompressedOff) }
func (t FileChunk) IsChunkManifest() bool {
	return t.o.Bool(fileChunkIsChunkManifestOff)
}

// SseType reads the sse_type enum (see SSEType* constants).
func (t FileChunk) SseType() uint32     { return t.o.Uint32(fileChunkSseTypeOff) }
func (t FileChunk) SseMetadata() []byte { return t.o.Bytes(fileChunkSseMetadataOff) }

// HasFid reports whether the fid nested message is present.
func (t FileChunk) HasFid() bool { return !t.o.Object(fileChunkFidOff).IsNull() }

// HasSourceFid reports whether the source_fid nested message is present.
func (t FileChunk) HasSourceFid() bool { return !t.o.Object(fileChunkSourceFidOff).IsNull() }

// FileChunkInput collects the field values for NewFileChunk. Fid and SourceFid,
// when non-nil, are standalone FileId buffers (from NewFileId).
type FileChunkInput struct {
	FileID          string
	Offset          int64
	Size            uint64
	ModifiedTsNs    int64
	ETag            string
	SourceFileID    string
	Fid             []byte
	SourceFid       []byte
	CipherKey       []byte
	IsCompressed    bool
	IsChunkManifest bool
	SseType         uint32
	SseMetadata     []byte
}

// NewFileChunk builds a ZAP-encoded FileChunk message from in.
func NewFileChunk(in FileChunkInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(fileChunkSize)
	ob.SetText(fileChunkFileIDOff, in.FileID)
	ob.SetInt64(fileChunkOffsetOff, in.Offset)
	ob.SetUint64(fileChunkSizeFieldOff, in.Size)
	ob.SetInt64(fileChunkModifiedTsNsOff, in.ModifiedTsNs)
	ob.SetText(fileChunkETagOff, in.ETag)
	ob.SetText(fileChunkSourceFileIDOff, in.SourceFileID)
	setNestedObject(b, ob, fileChunkFidOff, in.Fid)
	setNestedObject(b, ob, fileChunkSourceFidOff, in.SourceFid)
	ob.SetBytes(fileChunkCipherKeyOff, in.CipherKey)
	ob.SetBool(fileChunkIsCompressedOff, in.IsCompressed)
	ob.SetBool(fileChunkIsChunkManifestOff, in.IsChunkManifest)
	ob.SetUint32(fileChunkSseTypeOff, in.SseType)
	ob.SetBytes(fileChunkSseMetadataOff, in.SseMetadata)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- FileChunkManifest ---

const (
	fileChunkManifestChunksOff = 0
	fileChunkManifestSize      = 8
)

// FileChunkManifest is a zero-copy view into a ZAP-encoded FileChunkManifest message.
type FileChunkManifest struct{ o zap.Object }

// WrapFileChunkManifest parses b and returns a typed view.
func WrapFileChunkManifest(b []byte) (FileChunkManifest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FileChunkManifest{}, err
	}
	return FileChunkManifest{o: m.Root()}, nil
}

// ChunksLen returns the number of chunks elements (proto field 1, repeated FileChunk).
func (t FileChunkManifest) ChunksLen() int { return t.o.List(fileChunkManifestChunksOff).Len() }

// Chunks returns the i-th chunks element as a typed FileChunk view.
func (t FileChunkManifest) Chunks(i int) FileChunk {
	return FileChunk{o: t.o.List(fileChunkManifestChunksOff).ObjectAt(i)}
}

// FileChunkManifestInput collects the field values for NewFileChunkManifest.
// Each Chunks element is a standalone FileChunk buffer (from NewFileChunk).
type FileChunkManifestInput struct {
	Chunks [][]byte
}

// NewFileChunkManifest builds a ZAP-encoded FileChunkManifest message from in.
func NewFileChunkManifest(in FileChunkManifestInput) []byte {
	b := zap.NewBuilder(256)
	chunksOff, chunksLen := addObjectList(b, in.Chunks)
	ob := b.StartObject(fileChunkManifestSize)
	ob.SetList(fileChunkManifestChunksOff, chunksOff, chunksLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- FuseAttributes ---

const (
	fuseAttributesFileSizeOff      = 0   // file_size (field 1, uint64)
	fuseAttributesMtimeOff         = 8   // mtime (field 2, int64)
	fuseAttributesFileModeOff      = 16  // file_mode (field 3, uint32)
	fuseAttributesUidOff           = 20  // uid (field 4, uint32)
	fuseAttributesGidOff           = 24  // gid (field 5, uint32)
	fuseAttributesCrtimeOff        = 32  // crtime (field 6, int64)
	fuseAttributesMimeOff          = 40  // mime (field 7, string)
	fuseAttributesTtlSecOff        = 48  // ttl_sec (field 10, int32)
	fuseAttributesUserNameOff      = 56  // user_name (field 11, string)
	fuseAttributesGroupNameOff     = 64  // group_name (field 12, repeated string)
	fuseAttributesSymlinkTargetOff = 72  // symlink_target (field 13, string)
	fuseAttributesMd5Off           = 80  // md5 (field 14, bytes)
	fuseAttributesRdevOff          = 88  // rdev (field 16, uint32)
	fuseAttributesInodeOff         = 96  // inode (field 17, uint64)
	fuseAttributesCtimeOff         = 104 // ctime (field 18, int64)
	fuseAttributesMtimeNsOff       = 112 // mtime_ns (field 19, int32)
	fuseAttributesCtimeNsOff       = 116 // ctime_ns (field 20, int32)
	fuseAttributesCrtimeNsOff      = 120 // crtime_ns (field 21, int32)
	fuseAttributesAtimeOff         = 128 // atime (field 22, int64)
	fuseAttributesAtimeNsOff       = 136 // atime_ns (field 23, int32)
	fuseAttributesSize             = 144
)

// FuseAttributes is a zero-copy view into a ZAP-encoded FuseAttributes message.
type FuseAttributes struct{ o zap.Object }

// WrapFuseAttributes parses b and returns a typed view.
func WrapFuseAttributes(b []byte) (FuseAttributes, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FuseAttributes{}, err
	}
	return FuseAttributes{o: m.Root()}, nil
}

func (t FuseAttributes) FileSize() uint64 { return t.o.Uint64(fuseAttributesFileSizeOff) }
func (t FuseAttributes) Mtime() int64     { return t.o.Int64(fuseAttributesMtimeOff) }
func (t FuseAttributes) FileMode() uint32 { return t.o.Uint32(fuseAttributesFileModeOff) }
func (t FuseAttributes) Uid() uint32      { return t.o.Uint32(fuseAttributesUidOff) }
func (t FuseAttributes) Gid() uint32      { return t.o.Uint32(fuseAttributesGidOff) }
func (t FuseAttributes) Crtime() int64    { return t.o.Int64(fuseAttributesCrtimeOff) }
func (t FuseAttributes) Mime() string     { return t.o.Text(fuseAttributesMimeOff) }
func (t FuseAttributes) TtlSec() int32    { return t.o.Int32(fuseAttributesTtlSecOff) }
func (t FuseAttributes) UserName() string { return t.o.Text(fuseAttributesUserNameOff) }
func (t FuseAttributes) SymlinkTarget() string {
	return t.o.Text(fuseAttributesSymlinkTargetOff)
}
func (t FuseAttributes) Md5() []byte     { return t.o.Bytes(fuseAttributesMd5Off) }
func (t FuseAttributes) Rdev() uint32    { return t.o.Uint32(fuseAttributesRdevOff) }
func (t FuseAttributes) Inode() uint64   { return t.o.Uint64(fuseAttributesInodeOff) }
func (t FuseAttributes) Ctime() int64    { return t.o.Int64(fuseAttributesCtimeOff) }
func (t FuseAttributes) MtimeNs() int32  { return t.o.Int32(fuseAttributesMtimeNsOff) }
func (t FuseAttributes) CtimeNs() int32  { return t.o.Int32(fuseAttributesCtimeNsOff) }
func (t FuseAttributes) CrtimeNs() int32 { return t.o.Int32(fuseAttributesCrtimeNsOff) }
func (t FuseAttributes) Atime() int64    { return t.o.Int64(fuseAttributesAtimeOff) }
func (t FuseAttributes) AtimeNs() int32  { return t.o.Int32(fuseAttributesAtimeNsOff) }

// GroupNameLen returns the number of group_name elements (proto field 12,
// repeated string).
func (t FuseAttributes) GroupNameLen() int {
	return t.o.List(fuseAttributesGroupNameOff).Len()
}

// GroupName returns the i-th group_name element.
func (t FuseAttributes) GroupName(i int) string {
	return elemText(t.o.List(fuseAttributesGroupNameOff), i)
}

// FuseAttributesInput collects the field values for NewFuseAttributes.
type FuseAttributesInput struct {
	FileSize      uint64
	Mtime         int64
	FileMode      uint32
	Uid           uint32
	Gid           uint32
	Crtime        int64
	Mime          string
	TtlSec        int32
	UserName      string
	GroupName     []string
	SymlinkTarget string
	Md5           []byte
	Rdev          uint32
	Inode         uint64
	Ctime         int64
	MtimeNs       int32
	CtimeNs       int32
	CrtimeNs      int32
	Atime         int64
	AtimeNs       int32
}

// NewFuseAttributes builds a ZAP-encoded FuseAttributes message from in.
func NewFuseAttributes(in FuseAttributesInput) []byte {
	b := zap.NewBuilder(256)
	groupOff, groupLen := addStringList(b, in.GroupName)
	ob := b.StartObject(fuseAttributesSize)
	ob.SetUint64(fuseAttributesFileSizeOff, in.FileSize)
	ob.SetInt64(fuseAttributesMtimeOff, in.Mtime)
	ob.SetUint32(fuseAttributesFileModeOff, in.FileMode)
	ob.SetUint32(fuseAttributesUidOff, in.Uid)
	ob.SetUint32(fuseAttributesGidOff, in.Gid)
	ob.SetInt64(fuseAttributesCrtimeOff, in.Crtime)
	ob.SetText(fuseAttributesMimeOff, in.Mime)
	ob.SetInt32(fuseAttributesTtlSecOff, in.TtlSec)
	ob.SetText(fuseAttributesUserNameOff, in.UserName)
	ob.SetList(fuseAttributesGroupNameOff, groupOff, groupLen)
	ob.SetText(fuseAttributesSymlinkTargetOff, in.SymlinkTarget)
	ob.SetBytes(fuseAttributesMd5Off, in.Md5)
	ob.SetUint32(fuseAttributesRdevOff, in.Rdev)
	ob.SetUint64(fuseAttributesInodeOff, in.Inode)
	ob.SetInt64(fuseAttributesCtimeOff, in.Ctime)
	ob.SetInt32(fuseAttributesMtimeNsOff, in.MtimeNs)
	ob.SetInt32(fuseAttributesCtimeNsOff, in.CtimeNs)
	ob.SetInt32(fuseAttributesCrtimeNsOff, in.CrtimeNs)
	ob.SetInt64(fuseAttributesAtimeOff, in.Atime)
	ob.SetInt32(fuseAttributesAtimeNsOff, in.AtimeNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- EntryExtendedEntry (one entry of Entry.extended map<string, bytes>) ---

const (
	entryExtendedEntryKeyOff   = 0
	entryExtendedEntryValueOff = 8
	entryExtendedEntrySize     = 16
)

// EntryExtendedEntry is a zero-copy view of one map<string,bytes> pair of
// Entry.extended (proto field 5).
type EntryExtendedEntry struct{ o zap.Object }

// WrapEntryExtendedEntry parses b and returns a typed view.
func WrapEntryExtendedEntry(b []byte) (EntryExtendedEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return EntryExtendedEntry{}, err
	}
	return EntryExtendedEntry{o: m.Root()}, nil
}

func (t EntryExtendedEntry) Key() string   { return t.o.Text(entryExtendedEntryKeyOff) }
func (t EntryExtendedEntry) Value() []byte { return t.o.Bytes(entryExtendedEntryValueOff) }

// EntryExtendedEntryInput collects one map<string,bytes> pair.
type EntryExtendedEntryInput struct {
	Key   string
	Value []byte
}

// NewEntryExtendedEntry builds one encoded Entry.extended map entry.
func NewEntryExtendedEntry(in EntryExtendedEntryInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(entryExtendedEntrySize)
	ob.SetText(entryExtendedEntryKeyOff, in.Key)
	ob.SetBytes(entryExtendedEntryValueOff, in.Value)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- Entry ---

const (
	entryNameOff               = 0  // name (field 1, string)
	entryIsDirectoryOff        = 8  // is_directory (field 2, bool)
	entryChunksOff             = 16 // chunks (field 3, repeated FileChunk)
	entryAttributesOff         = 24 // attributes (field 4, FuseAttributes object)
	entryExtendedOff           = 32 // extended (field 5, map<string, bytes>)
	entryHardLinkIDOff         = 40 // hard_link_id (field 7, bytes)
	entryHardLinkCounterOff    = 48 // hard_link_counter (field 8, int32)
	entryContentOff            = 56 // content (field 9, bytes)
	entryRemoteEntryOff        = 64 // remote_entry (field 10, RemoteEntry object)
	entryQuotaOff              = 72 // quota (field 11, int64)
	entryWormEnforcedAtTsNsOff = 80 // worm_enforced_at_ts_ns (field 12, int64)
	entrySize                  = 88
)

// Entry is a zero-copy view into a ZAP-encoded Entry message.
type Entry struct{ o zap.Object }

// WrapEntry parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapEntry(b []byte) (Entry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Entry{}, err
	}
	return Entry{o: m.Root()}, nil
}

func (t Entry) Name() string           { return t.o.Text(entryNameOff) }
func (t Entry) IsDirectory() bool      { return t.o.Bool(entryIsDirectoryOff) }
func (t Entry) HardLinkID() []byte     { return t.o.Bytes(entryHardLinkIDOff) }
func (t Entry) HardLinkCounter() int32 { return t.o.Int32(entryHardLinkCounterOff) }
func (t Entry) Content() []byte        { return t.o.Bytes(entryContentOff) }
func (t Entry) Quota() int64           { return t.o.Int64(entryQuotaOff) }
func (t Entry) WormEnforcedAtTsNs() int64 {
	return t.o.Int64(entryWormEnforcedAtTsNsOff)
}

// ChunksLen returns the number of chunks elements (proto field 3, repeated FileChunk).
func (t Entry) ChunksLen() int { return t.o.List(entryChunksOff).Len() }

// Chunks returns the i-th chunks element as a typed FileChunk view.
func (t Entry) Chunks(i int) FileChunk {
	return FileChunk{o: t.o.List(entryChunksOff).ObjectAt(i)}
}

// Attributes returns the nested FuseAttributes (null-object if unset).
func (t Entry) Attributes() FuseAttributes {
	return FuseAttributes{o: t.o.Object(entryAttributesOff)}
}

// HasAttributes reports whether the attributes field is present.
func (t Entry) HasAttributes() bool { return !t.o.Object(entryAttributesOff).IsNull() }

// ExtendedLen returns the number of extended map entries (proto field 5,
// map<string, bytes>).
func (t Entry) ExtendedLen() int { return t.o.List(entryExtendedOff).Len() }

// Extended returns the i-th extended map entry as a typed EntryExtendedEntry.
func (t Entry) Extended(i int) EntryExtendedEntry {
	return EntryExtendedEntry{o: t.o.List(entryExtendedOff).ObjectAt(i)}
}

// RemoteEntry returns the nested RemoteEntry (null-object if unset).
func (t Entry) RemoteEntry() RemoteEntry {
	return RemoteEntry{o: t.o.Object(entryRemoteEntryOff)}
}

// HasRemoteEntry reports whether the remote_entry field is present.
func (t Entry) HasRemoteEntry() bool { return !t.o.Object(entryRemoteEntryOff).IsNull() }

// EntryInput collects the field values for NewEntry. Chunks elements are
// standalone FileChunk buffers; Attributes/RemoteEntry are standalone buffers;
// each Extended element is a pre-built EntryExtendedEntry buffer.
type EntryInput struct {
	Name               string
	IsDirectory        bool
	Chunks             [][]byte
	Attributes         []byte
	Extended           [][]byte
	HardLinkID         []byte
	HardLinkCounter    int32
	Content            []byte
	RemoteEntry        []byte
	Quota              int64
	WormEnforcedAtTsNs int64
}

// NewEntry builds a ZAP-encoded Entry message from in.
func NewEntry(in EntryInput) []byte {
	b := zap.NewBuilder(512)
	chunksOff, chunksLen := addObjectList(b, in.Chunks)
	extOff, extLen := addObjectList(b, in.Extended)
	ob := b.StartObject(entrySize)
	ob.SetText(entryNameOff, in.Name)
	ob.SetBool(entryIsDirectoryOff, in.IsDirectory)
	ob.SetList(entryChunksOff, chunksOff, chunksLen)
	setNestedObject(b, ob, entryAttributesOff, in.Attributes)
	ob.SetList(entryExtendedOff, extOff, extLen)
	ob.SetBytes(entryHardLinkIDOff, in.HardLinkID)
	ob.SetInt32(entryHardLinkCounterOff, in.HardLinkCounter)
	ob.SetBytes(entryContentOff, in.Content)
	setNestedObject(b, ob, entryRemoteEntryOff, in.RemoteEntry)
	ob.SetInt64(entryQuotaOff, in.Quota)
	ob.SetInt64(entryWormEnforcedAtTsNsOff, in.WormEnforcedAtTsNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- FullEntry ---

const (
	fullEntryDirOff   = 0 // dir (field 1, string)
	fullEntryEntryOff = 8 // entry (field 2, Entry object)
	fullEntrySize     = 16
)

// FullEntry is a zero-copy view into a ZAP-encoded FullEntry message.
type FullEntry struct{ o zap.Object }

// WrapFullEntry parses b and returns a typed view.
func WrapFullEntry(b []byte) (FullEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FullEntry{}, err
	}
	return FullEntry{o: m.Root()}, nil
}

func (t FullEntry) Dir() string  { return t.o.Text(fullEntryDirOff) }
func (t FullEntry) Entry() Entry { return Entry{o: t.o.Object(fullEntryEntryOff)} }

// HasEntry reports whether the entry field is present.
func (t FullEntry) HasEntry() bool { return !t.o.Object(fullEntryEntryOff).IsNull() }

// FullEntryInput collects the field values for NewFullEntry. Entry, when
// non-nil, is a standalone Entry buffer (from NewEntry).
type FullEntryInput struct {
	Dir   string
	Entry []byte
}

// NewFullEntry builds a ZAP-encoded FullEntry message from in.
func NewFullEntry(in FullEntryInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(fullEntrySize)
	ob.SetText(fullEntryDirOff, in.Dir)
	setNestedObject(b, ob, fullEntryEntryOff, in.Entry)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- EventNotification ---

const (
	eventNotificationOldEntryOff           = 0  // old_entry (field 1, Entry object)
	eventNotificationNewEntryOff           = 4  // new_entry (field 2, Entry object)
	eventNotificationDeleteChunksOff       = 8  // delete_chunks (field 3, bool)
	eventNotificationNewParentPathOff      = 16 // new_parent_path (field 4, string)
	eventNotificationIsFromOtherClusterOff = 24 // is_from_other_cluster (field 5, bool)
	eventNotificationSignaturesOff         = 32 // signatures (field 6, repeated int32)
	eventNotificationSize                  = 40
)

// EventNotification is a zero-copy view into a ZAP-encoded EventNotification message.
type EventNotification struct{ o zap.Object }

// WrapEventNotification parses b and returns a typed view.
func WrapEventNotification(b []byte) (EventNotification, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return EventNotification{}, err
	}
	return EventNotification{o: m.Root()}, nil
}

func (t EventNotification) OldEntry() Entry {
	return Entry{o: t.o.Object(eventNotificationOldEntryOff)}
}
func (t EventNotification) NewEntry() Entry {
	return Entry{o: t.o.Object(eventNotificationNewEntryOff)}
}
func (t EventNotification) DeleteChunks() bool {
	return t.o.Bool(eventNotificationDeleteChunksOff)
}
func (t EventNotification) NewParentPath() string {
	return t.o.Text(eventNotificationNewParentPathOff)
}
func (t EventNotification) IsFromOtherCluster() bool {
	return t.o.Bool(eventNotificationIsFromOtherClusterOff)
}

// HasOldEntry reports whether the old_entry field is present.
func (t EventNotification) HasOldEntry() bool {
	return !t.o.Object(eventNotificationOldEntryOff).IsNull()
}

// HasNewEntry reports whether the new_entry field is present.
func (t EventNotification) HasNewEntry() bool {
	return !t.o.Object(eventNotificationNewEntryOff).IsNull()
}

// SignaturesLen returns the number of signatures elements (proto field 6,
// repeated int32).
func (t EventNotification) SignaturesLen() int {
	return t.o.List(eventNotificationSignaturesOff).Len()
}

// Signature returns the i-th signatures element.
func (t EventNotification) Signature(i int) int32 {
	return int32(t.o.List(eventNotificationSignaturesOff).Uint32(i))
}

// EventNotificationInput collects the field values for NewEventNotification.
// OldEntry/NewEntry, when non-nil, are standalone Entry buffers (from NewEntry).
type EventNotificationInput struct {
	OldEntry           []byte
	NewEntry           []byte
	DeleteChunks       bool
	NewParentPath      string
	IsFromOtherCluster bool
	Signatures         []int32
}

// NewEventNotification builds a ZAP-encoded EventNotification message from in.
func NewEventNotification(in EventNotificationInput) []byte {
	b := zap.NewBuilder(256)
	sigOff, sigLen := addInt32List(b, in.Signatures)
	ob := b.StartObject(eventNotificationSize)
	setNestedObject(b, ob, eventNotificationOldEntryOff, in.OldEntry)
	setNestedObject(b, ob, eventNotificationNewEntryOff, in.NewEntry)
	ob.SetBool(eventNotificationDeleteChunksOff, in.DeleteChunks)
	ob.SetText(eventNotificationNewParentPathOff, in.NewParentPath)
	ob.SetBool(eventNotificationIsFromOtherClusterOff, in.IsFromOtherCluster)
	ob.SetList(eventNotificationSignaturesOff, sigOff, sigLen)
	ob.FinishAsRoot()
	return b.Finish()
}
