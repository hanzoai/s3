// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// convert.go is the proto<->wire converter set for the 27 unary HanzoFiler RPCs
// (the streaming RPCs live in stream.go). For each RPC it provides:
//
//	<Rpc>RequestToWire(*filer_pb.<Rpc>Request) []byte    — build the ZAP request buffer
//	<Rpc>ResponseFromWire([]byte) (*filer_pb.<Rpc>Response, error) — decode the reply
//
// plus the shared message converters (WriteCondition, ObjectMutation, PosixLockRange,
// MountInfo, the metadata-event graph: SubscribeMetadataResponse -> EventNotification
// -> Entry, LogFileChunkRef). These are the ONLY place proto<->wire mapping for the
// filer's request/response surface lives — the ZAP-backed filer_pb.HanzoFilerClient
// (client.go) and the ZAP filer dispatch reuse them. Built on the Entry converters in
// entry.go.

package filerzap

import (
	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	filerwire "github.com/hanzoai/s3/s3/wire/filer"
)

// --- shared leaf converters ---

func chunksToWire(cs []*filer_pb.FileChunk) [][]byte {
	if len(cs) == 0 {
		return nil
	}
	out := make([][]byte, len(cs))
	for i, c := range cs {
		out[i] = fileChunkToWire(c)
	}
	return out
}

func chunksFromView(n int, at func(int) filerwire.FileChunk) []*filer_pb.FileChunk {
	if n == 0 {
		return nil
	}
	out := make([]*filer_pb.FileChunk, n)
	for i := 0; i < n; i++ {
		out[i] = fileChunkFromWire(at(i))
	}
	return out
}

func extendedToWire(m map[string][]byte) [][]byte {
	if len(m) == 0 {
		return nil
	}
	out := make([][]byte, 0, len(m))
	for k, v := range m {
		out = append(out, filerwire.NewEntryExtendedEntry(filerwire.EntryExtendedEntryInput{Key: k, Value: v}))
	}
	return out
}

// --- WriteCondition ---

func writeConditionToWire(c *filer_pb.WriteCondition) []byte {
	if c == nil {
		return nil
	}
	clauses := make([][]byte, len(c.Clauses))
	for i, cl := range c.Clauses {
		clauses[i] = filerwire.NewWriteConditionClause(filerwire.WriteConditionClauseInput{
			Kind:      uint32(cl.Kind),
			Etags:     cl.Etags,
			UnixTime:  cl.UnixTime,
			AllowWeak: cl.AllowWeak,
			ExtKey:    cl.ExtKey,
			ExtValue:  cl.ExtValue,
			GateKey:   cl.GateKey,
			GateValue: cl.GateValue,
		})
	}
	return filerwire.NewWriteCondition(filerwire.WriteConditionInput{Clauses: clauses})
}

// --- Recompute ---

func recomputeToWire(r *filer_pb.Recompute) []byte {
	if r == nil {
		return nil
	}
	copyExt := make([][]byte, 0, len(r.CopyExtended))
	for k, v := range r.CopyExtended {
		copyExt = append(copyExt, filerwire.NewEntryExtendedEntry(filerwire.EntryExtendedEntryInput{Key: k, Value: []byte(v)}))
	}
	return filerwire.NewRecompute(filerwire.RecomputeInput{
		ScanDir:      r.ScanDir,
		Descending:   r.Descending,
		CopyExtended: copyExt,
		NameToKey:    r.NameToKey,
		SizeToKey:    r.SizeToKey,
		MtimeToKey:   r.MtimeToKey,
		DemoteKey:    r.DemoteKey,
		DemoteValue:  r.DemoteValue,
		ExcludeName:  r.ExcludeName,
	})
}

// --- ObjectMutation ---

func objectMutationToWire(m *filer_pb.ObjectMutation) []byte {
	if m == nil {
		return nil
	}
	return filerwire.NewObjectMutation(filerwire.ObjectMutationInput{
		Type:           uint32(m.Type),
		Directory:      m.Directory,
		Name:           m.Name,
		Entry:          EntryToWire(m.Entry),
		SetExtended:    extendedToWire(m.SetExtended),
		DeleteExtended: m.DeleteExtended,
		IsDeleteData:   m.IsDeleteData,
		IsRecursive:    m.IsRecursive,
		Recompute:      recomputeToWire(m.Recompute),
		SetContent:     m.SetContent,
		Content:        m.Content,
		TouchMtime:     m.TouchMtime,
	})
}

// --- PosixLockRange ---

func posixLockRangeToWire(r *filer_pb.PosixLockRange) []byte {
	if r == nil {
		return nil
	}
	return filerwire.NewPosixLockRange(filerwire.PosixLockRangeInput{
		Start:   r.Start,
		End:     r.End,
		Type:    r.Type,
		Sid:     r.Sid,
		Owner:   r.Owner,
		Pid:     r.Pid,
		IsFlock: r.IsFlock,
	})
}

func posixLockRangeFromView(v filerwire.PosixLockRange) *filer_pb.PosixLockRange {
	return &filer_pb.PosixLockRange{
		Start:   v.Start(),
		End:     v.End(),
		Type:    v.Type(),
		Sid:     v.Sid(),
		Owner:   v.Owner(),
		Pid:     v.Pid(),
		IsFlock: v.IsFlock(),
	}
}

// --- metadata-event graph ---

func eventNotificationToWire(e *filer_pb.EventNotification) []byte {
	if e == nil {
		return nil
	}
	return filerwire.NewEventNotification(filerwire.EventNotificationInput{
		OldEntry:           EntryToWire(e.OldEntry),
		NewEntry:           EntryToWire(e.NewEntry),
		DeleteChunks:       e.DeleteChunks,
		NewParentPath:      e.NewParentPath,
		IsFromOtherCluster: e.IsFromOtherCluster,
		Signatures:         e.Signatures,
	})
}

func eventNotificationFromView(v filerwire.EventNotification) *filer_pb.EventNotification {
	e := &filer_pb.EventNotification{
		DeleteChunks:       v.DeleteChunks(),
		NewParentPath:      v.NewParentPath(),
		IsFromOtherCluster: v.IsFromOtherCluster(),
	}
	if v.HasOldEntry() {
		e.OldEntry = entryFromView(v.OldEntry())
	}
	if v.HasNewEntry() {
		e.NewEntry = entryFromView(v.NewEntry())
	}
	if n := v.SignaturesLen(); n > 0 {
		e.Signatures = make([]int32, n)
		for i := 0; i < n; i++ {
			e.Signatures[i] = v.Signature(i)
		}
	}
	return e
}

func logFileChunkRefToWire(r *filer_pb.LogFileChunkRef) []byte {
	if r == nil {
		return nil
	}
	return filerwire.NewLogFileChunkRef(filerwire.LogFileChunkRefInput{
		Chunks:   chunksToWire(r.Chunks),
		FileTsNs: r.FileTsNs,
		FilerID:  r.FilerId,
	})
}

func logFileChunkRefFromView(v filerwire.LogFileChunkRef) *filer_pb.LogFileChunkRef {
	return &filer_pb.LogFileChunkRef{
		Chunks:   chunksFromView(v.ChunksLen(), v.Chunks),
		FileTsNs: v.FileTsNs(),
		FilerId:  v.FilerID(),
	}
}

// SubscribeMetadataResponseToWire encodes the recursive metadata-event aggregate.
func SubscribeMetadataResponseToWire(r *filer_pb.SubscribeMetadataResponse) []byte {
	if r == nil {
		return nil
	}
	events := make([][]byte, len(r.Events))
	for i, ev := range r.Events {
		events[i] = SubscribeMetadataResponseToWire(ev)
	}
	refs := make([][]byte, len(r.LogFileRefs))
	for i, ref := range r.LogFileRefs {
		refs[i] = logFileChunkRefToWire(ref)
	}
	return filerwire.NewSubscribeMetadataResponse(filerwire.SubscribeMetadataResponseInput{
		Directory:         r.Directory,
		EventNotification: eventNotificationToWire(r.EventNotification),
		TsNs:              r.TsNs,
		Events:            events,
		LogFileRefs:       refs,
	})
}

// SubscribeMetadataResponseFromWire decodes a metadata-event buffer.
func SubscribeMetadataResponseFromWire(buf []byte) (*filer_pb.SubscribeMetadataResponse, error) {
	v, err := filerwire.WrapSubscribeMetadataResponse(buf)
	if err != nil {
		return nil, err
	}
	return subscribeMetadataResponseFromView(v), nil
}

func subscribeMetadataResponseFromView(v filerwire.SubscribeMetadataResponse) *filer_pb.SubscribeMetadataResponse {
	r := &filer_pb.SubscribeMetadataResponse{
		Directory: v.Directory(),
		TsNs:      v.TsNs(),
	}
	if v.HasEventNotification() {
		r.EventNotification = eventNotificationFromView(v.EventNotification())
	}
	if n := v.EventsLen(); n > 0 {
		r.Events = make([]*filer_pb.SubscribeMetadataResponse, n)
		for i := 0; i < n; i++ {
			r.Events[i] = subscribeMetadataResponseFromView(v.Events(i))
		}
	}
	if n := v.LogFileRefsLen(); n > 0 {
		r.LogFileRefs = make([]*filer_pb.LogFileChunkRef, n)
		for i := 0; i < n; i++ {
			r.LogFileRefs[i] = logFileChunkRefFromView(v.LogFileRef(i))
		}
	}
	return r
}

// metadataEventFromView is the optional-nested form used by Create/Update/Delete/Cache
// responses that embed a metadata_event field.
func metadataEventFromView(has bool, v filerwire.SubscribeMetadataResponse) *filer_pb.SubscribeMetadataResponse {
	if !has {
		return nil
	}
	return subscribeMetadataResponseFromView(v)
}
