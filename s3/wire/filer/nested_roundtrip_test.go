// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package filerwire

import (
	"bytes"
	"testing"
)

// TestNestedComposition is the worst-case zero-copy round-trip for the singular
// nested-object embedding (setNestedObject / t.o.Object): a deeply-nested Entry
// carrying
//   - a repeated FileChunk list (addObjectList path) where EACH chunk itself
//     holds a singular nested FileId (the embed-inside-a-list-element case),
//   - a singular nested FuseAttributes,
//   - a map<string,bytes> extended,
//   - a singular nested RemoteEntry,
// then embedded AGAIN as the singular Entry of a LookupDirectoryEntryResponse.
// It proves every internal relative pointer survives being relocated whole into
// the parent buffer (the fix in setNestedObject/embedMessage) across two embed
// levels and through a list element. No socket needed — this is a pure encode →
// Wrap → read assertion; the socket round-trip lives in
// service_adapter_roundtrip_test.go and the filerstream tests.
func TestNestedComposition(t *testing.T) {
	fid := NewFileId(FileIdInput{VolumeID: 7, FileKey: 0xDEADBEEF, Cookie: 0x1234})
	chunk0 := NewFileChunk(FileChunkInput{
		FileID:       "3,01abcd",
		Offset:       4096,
		Size:         1 << 20,
		ETag:         "chunk-etag-0",
		Fid:          fid,
		IsCompressed: true,
		SseType:      SSETypeSSEKMS,
	})
	chunk1 := NewFileChunk(FileChunkInput{
		FileID: "5,02ef01",
		Offset: 0,
		Size:   512,
		ETag:   "chunk-etag-1",
		// no Fid on this one: HasFid must be false after round-trip
	})

	attrs := NewFuseAttributes(FuseAttributesInput{
		FileSize: 1<<20 + 512,
		FileMode: 0o644,
		Uid:      1000,
		Gid:      1000,
		Mime:     "application/octet-stream",
		Inode:    42,
	})
	remote := NewRemoteEntry(RemoteEntryInput{
		StorageName: "s3-backend",
		RemoteETag:  "remote-etag",
		RemoteSize:  1 << 20,
	})
	ext0 := NewEntryExtendedEntry(EntryExtendedEntryInput{Key: "x-amz-meta-a", Value: []byte("alpha")})
	ext1 := NewEntryExtendedEntry(EntryExtendedEntryInput{Key: "x-amz-meta-b", Value: []byte("bravo")})

	content := []byte("the bytes are the message")
	entry := NewEntry(EntryInput{
		Name:        "deep.bin",
		IsDirectory: false,
		Chunks:      [][]byte{chunk0, chunk1},
		Attributes:  attrs,
		Extended:    [][]byte{ext0, ext1},
		Content:     content,
		RemoteEntry: remote,
		Quota:       -1,
	})

	// Embed the Entry as the singular nested field of a response (second level).
	respBuf := NewLookupDirectoryEntryResponse(LookupDirectoryEntryResponseInput{Entry: entry})

	resp, err := WrapLookupDirectoryEntryResponse(respBuf)
	if err != nil {
		t.Fatalf("WrapLookupDirectoryEntryResponse: %v", err)
	}
	if !resp.HasEntry() {
		t.Fatal("HasEntry = false, want true")
	}
	e := resp.Entry()

	if e.Name() != "deep.bin" {
		t.Errorf("Name = %q, want deep.bin", e.Name())
	}
	if !bytes.Equal(e.Content(), content) {
		t.Errorf("Content = %q, want %q", e.Content(), content)
	}
	if e.Quota() != -1 {
		t.Errorf("Quota = %d, want -1", e.Quota())
	}

	// --- repeated FileChunk, each with its OWN nested FileId ---
	if e.ChunksLen() != 2 {
		t.Fatalf("ChunksLen = %d, want 2", e.ChunksLen())
	}
	c0 := e.Chunks(0)
	if c0.FileID() != "3,01abcd" || c0.ETag() != "chunk-etag-0" {
		t.Errorf("chunk0 fileid/etag = %q/%q", c0.FileID(), c0.ETag())
	}
	if c0.Size() != 1<<20 || c0.Offset() != 4096 {
		t.Errorf("chunk0 size/offset = %d/%d", c0.Size(), c0.Offset())
	}
	if !c0.IsCompressed() {
		t.Error("chunk0 IsCompressed = false, want true")
	}
	if c0.SseType() != SSETypeSSEKMS {
		t.Errorf("chunk0 SseType = %d, want %d", c0.SseType(), SSETypeSSEKMS)
	}
	if !c0.HasFid() {
		t.Fatal("chunk0 HasFid = false, want true (nested-in-list-element)")
	}
	f := c0.Fid()
	if f.VolumeID() != 7 || f.FileKey() != 0xDEADBEEF || f.Cookie() != 0x1234 {
		t.Errorf("chunk0 fid = vol %d key %#x cookie %#x", f.VolumeID(), f.FileKey(), f.Cookie())
	}
	c1 := e.Chunks(1)
	if c1.FileID() != "5,02ef01" || c1.Size() != 512 {
		t.Errorf("chunk1 fileid/size = %q/%d", c1.FileID(), c1.Size())
	}
	if c1.HasFid() {
		t.Error("chunk1 HasFid = true, want false")
	}

	// --- singular nested FuseAttributes ---
	if !e.HasAttributes() {
		t.Fatal("HasAttributes = false, want true")
	}
	a := e.Attributes()
	if a.FileSize() != 1<<20+512 || a.FileMode() != 0o644 || a.Uid() != 1000 || a.Inode() != 42 {
		t.Errorf("attrs = size %d mode %o uid %d inode %d", a.FileSize(), a.FileMode(), a.Uid(), a.Inode())
	}
	if a.Mime() != "application/octet-stream" {
		t.Errorf("attrs mime = %q", a.Mime())
	}

	// --- map<string,bytes> extended ---
	if e.ExtendedLen() != 2 {
		t.Fatalf("ExtendedLen = %d, want 2", e.ExtendedLen())
	}
	got := map[string]string{}
	for i := 0; i < e.ExtendedLen(); i++ {
		kv := e.Extended(i)
		got[kv.Key()] = string(kv.Value())
	}
	if got["x-amz-meta-a"] != "alpha" || got["x-amz-meta-b"] != "bravo" {
		t.Errorf("extended = %v", got)
	}

	// --- singular nested RemoteEntry ---
	if !e.HasRemoteEntry() {
		t.Fatal("HasRemoteEntry = false, want true")
	}
	r := e.RemoteEntry()
	if r.StorageName() != "s3-backend" || r.RemoteETag() != "remote-etag" || r.RemoteSize() != 1<<20 {
		t.Errorf("remote = name %q etag %q size %d", r.StorageName(), r.RemoteETag(), r.RemoteSize())
	}
}
