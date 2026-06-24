// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Package filerzap is the gRPC→ZAP migration layer for the filer service: it
// converts the protobuf filer aggregate types (filer_pb) to and from their
// zero-copy ZAP wire encodings (filerwire), so the filer can serve and be
// called over the canonical github.com/zap-proto/go transport instead of gRPC.
//
// The Entry graph (Entry → FuseAttributes, []FileChunk → FileId, RemoteEntry,
// extended map) is the single most-referenced shape in the filer's 34-method
// surface (345 callsites). These converters are defined ONCE here and reused by
// the filer ZAP server dispatch and the ZAP-backed client facade — proto↔wire
// mapping for the file-metadata graph lives in exactly one place.
//
// This is the foundation of the full engine gRPC→ZAP rip. Each filer RPC's
// request/response converters build on the Entry converters below.
package filerzap

import (
	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	filerwire "github.com/hanzoai/s3/s3/wire/filer"
)

// --- FileId ---

func fileIDToWire(f *filer_pb.FileId) []byte {
	if f == nil {
		return nil
	}
	return filerwire.NewFileId(filerwire.FileIdInput{
		VolumeID: f.VolumeId,
		FileKey:  f.FileKey,
		Cookie:   f.Cookie,
	})
}

func fileIDFromWire(v filerwire.FileId) *filer_pb.FileId {
	return &filer_pb.FileId{
		VolumeId: v.VolumeID(),
		FileKey:  v.FileKey(),
		Cookie:   v.Cookie(),
	}
}

// --- FuseAttributes ---

func fuseAttrToWire(a *filer_pb.FuseAttributes) []byte {
	if a == nil {
		return nil
	}
	return filerwire.NewFuseAttributes(filerwire.FuseAttributesInput{
		FileSize:      a.FileSize,
		Mtime:         a.Mtime,
		FileMode:      a.FileMode,
		Uid:           a.Uid,
		Gid:           a.Gid,
		Crtime:        a.Crtime,
		Mime:          a.Mime,
		TtlSec:        a.TtlSec,
		UserName:      a.UserName,
		GroupName:     a.GroupName,
		SymlinkTarget: a.SymlinkTarget,
		Md5:           a.Md5,
		Rdev:          a.Rdev,
		Inode:         a.Inode,
		Ctime:         a.Ctime,
		MtimeNs:       a.MtimeNs,
		CtimeNs:       a.CtimeNs,
		CrtimeNs:      a.CrtimeNs,
		Atime:         a.Atime,
		AtimeNs:       a.AtimeNs,
	})
}

func fuseAttrFromWire(v filerwire.FuseAttributes) *filer_pb.FuseAttributes {
	groups := make([]string, v.GroupNameLen())
	for i := range groups {
		groups[i] = v.GroupName(i)
	}
	return &filer_pb.FuseAttributes{
		FileSize:      v.FileSize(),
		Mtime:         v.Mtime(),
		FileMode:      v.FileMode(),
		Uid:           v.Uid(),
		Gid:           v.Gid(),
		Crtime:        v.Crtime(),
		Mime:          v.Mime(),
		TtlSec:        v.TtlSec(),
		UserName:      v.UserName(),
		GroupName:     groups,
		SymlinkTarget: v.SymlinkTarget(),
		Md5:           v.Md5(),
		Rdev:          v.Rdev(),
		Inode:         v.Inode(),
		Ctime:         v.Ctime(),
		MtimeNs:       v.MtimeNs(),
		CtimeNs:       v.CtimeNs(),
		CrtimeNs:      v.CrtimeNs(),
		Atime:         v.Atime(),
		AtimeNs:       v.AtimeNs(),
	}
}

// --- FileChunk ---

func fileChunkToWire(c *filer_pb.FileChunk) []byte {
	if c == nil {
		return nil
	}
	return filerwire.NewFileChunk(filerwire.FileChunkInput{
		FileID:          c.FileId,
		Offset:          c.Offset,
		Size:            c.Size,
		ModifiedTsNs:    c.ModifiedTsNs,
		ETag:            c.ETag,
		SourceFileID:    c.SourceFileId,
		Fid:             fileIDToWire(c.Fid),
		SourceFid:       fileIDToWire(c.SourceFid),
		CipherKey:       c.CipherKey,
		IsCompressed:    c.IsCompressed,
		IsChunkManifest: c.IsChunkManifest,
		SseType:         uint32(c.SseType),
		SseMetadata:     c.SseMetadata,
	})
}

func fileChunkFromWire(v filerwire.FileChunk) *filer_pb.FileChunk {
	c := &filer_pb.FileChunk{
		FileId:          v.FileID(),
		Offset:          v.Offset(),
		Size:            v.Size(),
		ModifiedTsNs:    v.ModifiedTsNs(),
		ETag:            v.ETag(),
		SourceFileId:    v.SourceFileID(),
		CipherKey:       v.CipherKey(),
		IsCompressed:    v.IsCompressed(),
		IsChunkManifest: v.IsChunkManifest(),
		SseType:         filer_pb.SSEType(v.SseType()),
		SseMetadata:     v.SseMetadata(),
	}
	if v.HasFid() {
		c.Fid = fileIDFromWire(v.Fid())
	}
	if v.HasSourceFid() {
		c.SourceFid = fileIDFromWire(v.SourceFid())
	}
	return c
}

// --- RemoteEntry ---

func remoteEntryToWire(r *filer_pb.RemoteEntry) []byte {
	if r == nil {
		return nil
	}
	return filerwire.NewRemoteEntry(filerwire.RemoteEntryInput{
		StorageName:       r.StorageName,
		LastLocalSyncTsNs: r.LastLocalSyncTsNs,
		RemoteETag:        r.RemoteETag,
		RemoteMtime:       r.RemoteMtime,
		RemoteSize:        r.RemoteSize,
	})
}

func remoteEntryFromWire(v filerwire.RemoteEntry) *filer_pb.RemoteEntry {
	return &filer_pb.RemoteEntry{
		StorageName:       v.StorageName(),
		LastLocalSyncTsNs: v.LastLocalSyncTsNs(),
		RemoteETag:        v.RemoteETag(),
		RemoteMtime:       v.RemoteMtime(),
		RemoteSize:        v.RemoteSize(),
	}
}

// --- Entry (the keystone) ---

// EntryToWire encodes a protobuf Entry as a zero-copy ZAP Entry buffer.
func EntryToWire(e *filer_pb.Entry) []byte {
	if e == nil {
		return nil
	}
	chunks := make([][]byte, len(e.Chunks))
	for i, c := range e.Chunks {
		chunks[i] = fileChunkToWire(c)
	}
	// extended map<string,bytes> → list of EntryExtendedEntry buffers. Map
	// iteration order is non-deterministic in Go; the wire list order is not
	// semantically meaningful (it round-trips back into a map), so this is safe.
	ext := make([][]byte, 0, len(e.Extended))
	for k, val := range e.Extended {
		ext = append(ext, filerwire.NewEntryExtendedEntry(filerwire.EntryExtendedEntryInput{
			Key:   k,
			Value: val,
		}))
	}
	return filerwire.NewEntry(filerwire.EntryInput{
		Name:               e.Name,
		IsDirectory:        e.IsDirectory,
		Chunks:             chunks,
		Attributes:         fuseAttrToWire(e.Attributes),
		Extended:           ext,
		HardLinkID:         e.HardLinkId,
		HardLinkCounter:    e.HardLinkCounter,
		Content:            e.Content,
		RemoteEntry:        remoteEntryToWire(e.RemoteEntry),
		Quota:              e.Quota,
		WormEnforcedAtTsNs: e.WormEnforcedAtTsNs,
	})
}

// EntryFromWire decodes a ZAP Entry buffer back into a protobuf Entry.
func EntryFromWire(buf []byte) (*filer_pb.Entry, error) {
	v, err := filerwire.WrapEntry(buf)
	if err != nil {
		return nil, err
	}
	return entryFromView(v), nil
}

// entryFromView reconstructs a protobuf Entry from a (possibly nested) Entry
// view — shared by EntryFromWire and any response converter carrying an Entry.
func entryFromView(v filerwire.Entry) *filer_pb.Entry {
	e := &filer_pb.Entry{
		Name:               v.Name(),
		IsDirectory:        v.IsDirectory(),
		HardLinkId:         v.HardLinkID(),
		HardLinkCounter:    v.HardLinkCounter(),
		Content:            v.Content(),
		Quota:              v.Quota(),
		WormEnforcedAtTsNs: v.WormEnforcedAtTsNs(),
	}
	if n := v.ChunksLen(); n > 0 {
		e.Chunks = make([]*filer_pb.FileChunk, n)
		for i := 0; i < n; i++ {
			e.Chunks[i] = fileChunkFromWire(v.Chunks(i))
		}
	}
	if v.HasAttributes() {
		e.Attributes = fuseAttrFromWire(v.Attributes())
	}
	if n := v.ExtendedLen(); n > 0 {
		e.Extended = make(map[string][]byte, n)
		for i := 0; i < n; i++ {
			kv := v.Extended(i)
			e.Extended[kv.Key()] = kv.Value()
		}
	}
	if v.HasRemoteEntry() {
		e.RemoteEntry = remoteEntryFromWire(v.RemoteEntry())
	}
	return e
}
