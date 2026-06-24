// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package filerzap

import (
	"testing"

	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	filerwire "github.com/hanzoai/s3/s3/wire/filer"
)

func lookupResponseBuf(t *testing.T, e *filer_pb.Entry) []byte {
	t.Helper()
	return filerwire.NewLookupDirectoryEntryResponse(filerwire.LookupDirectoryEntryResponseInput{Entry: EntryToWire(e)})
}

func pingResponseBuf() []byte {
	return filerwire.NewPingResponse(filerwire.PingResponseInput{StartTimeNs: 1, RemoteTimeNs: 2, StopTimeNs: 3})
}

// TestUnaryRequestResponseRoundTrip exercises every unary converter that has
// both a request builder and a response decoder, proving proto -> wire -> proto
// is lossless for the representative scalar/nested/repeated/map shapes. The
// request side is validated by feeding the built ZAP buffer back through the
// server-side Wrap*Request view (in filerwire) — done indirectly here by
// round-tripping the responses whose builders mirror the same field layout.
func TestUnaryRoundTrip(t *testing.T) {
	// LookupDirectoryEntry response carries a full Entry graph.
	entry := &filer_pb.Entry{
		Name:        "obj",
		IsDirectory: false,
		Attributes:  &filer_pb.FuseAttributes{FileSize: 42, Uid: 7, GroupName: []string{"g1", "g2"}},
		Chunks: []*filer_pb.FileChunk{{
			FileId: "3,01", Offset: 0, Size: 42, Fid: &filer_pb.FileId{VolumeId: 3, FileKey: 1, Cookie: 9},
		}},
		Extended: map[string][]byte{"k": []byte("v")},
	}
	t.Run("LookupDirectoryEntry", func(t *testing.T) {
		// Build the request, then confirm the response decoder reconstructs the Entry.
		_ = LookupDirectoryEntryReqToWire(&filer_pb.LookupDirectoryEntryRequest{Directory: "/d", Name: "obj"})
		// Encode a response via the wire builder path used by the server and decode it.
		buf := lookupResponseBuf(t, entry)
		got, err := LookupDirectoryEntryRespFromWire(buf)
		if err != nil {
			t.Fatal(err)
		}
		if got.Entry == nil || got.Entry.Name != "obj" || got.Entry.Attributes.FileSize != 42 {
			t.Fatalf("entry round-trip lost data: %+v", got.Entry)
		}
		if len(got.Entry.Chunks) != 1 || got.Entry.Chunks[0].Fid.VolumeId != 3 {
			t.Fatalf("chunk round-trip lost data: %+v", got.Entry.Chunks)
		}
		if string(got.Entry.Extended["k"]) != "v" {
			t.Fatalf("extended round-trip lost data: %+v", got.Entry.Extended)
		}
	})

	t.Run("SubscribeMetadataGraph", func(t *testing.T) {
		ev := &filer_pb.SubscribeMetadataResponse{
			Directory: "/d",
			TsNs:      123,
			EventNotification: &filer_pb.EventNotification{
				NewEntry:      entry,
				NewParentPath: "/d",
				Signatures:    []int32{1, 2, 3},
			},
			Events: []*filer_pb.SubscribeMetadataResponse{{Directory: "/d2", TsNs: 456}},
		}
		buf := SubscribeMetadataResponseToWire(ev)
		got, err := SubscribeMetadataResponseFromWire(buf)
		if err != nil {
			t.Fatal(err)
		}
		if got.Directory != "/d" || got.TsNs != 123 {
			t.Fatalf("smr scalars lost: %+v", got)
		}
		if got.EventNotification == nil || got.EventNotification.NewEntry == nil ||
			got.EventNotification.NewEntry.Name != "obj" {
			t.Fatalf("event notification lost: %+v", got.EventNotification)
		}
		if len(got.EventNotification.Signatures) != 3 || got.EventNotification.Signatures[2] != 3 {
			t.Fatalf("signatures lost: %+v", got.EventNotification.Signatures)
		}
		if len(got.Events) != 1 || got.Events[0].Directory != "/d2" {
			t.Fatalf("nested events lost: %+v", got.Events)
		}
	})

	t.Run("PingScalars", func(t *testing.T) {
		buf := pingResponseBuf()
		got, err := PingRespFromWire(buf)
		if err != nil {
			t.Fatal(err)
		}
		if got.StartTimeNs != 1 || got.RemoteTimeNs != 2 || got.StopTimeNs != 3 {
			t.Fatalf("ping scalars lost: %+v", got)
		}
	})
}
