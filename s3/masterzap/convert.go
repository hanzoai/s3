// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// convert.go is the shared leaf proto<->wire converter set for the Hanzo
// (master) RPC surface. For each shared message it provides a *ToWire builder
// (*master_pb.X -> ZAP buffer) and a *FromView reader (masterwire view ->
// *master_pb.X). The per-RPC request builders and response decoders in rpc.go
// (unary) and stream.go (streaming) compose these. This is the ONLY place
// proto<->wire mapping for the master's shared message graph lives — the
// ZAP-backed master_pb.HanzoClient (client.go) and the ZAP master dispatch reuse
// them. Built directly on the generated masterwire New*/Wrap* surface.

package masterzap

import (
	master_pb "github.com/hanzoai/s3/s3/pb/master_pb"
	volume_server_pb "github.com/hanzoai/s3/s3/pb/volume_server_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
	volume_serverwire "github.com/hanzoai/s3/s3/wire/volume_server"
)

// --- Location ---

func locationToWire(l *master_pb.Location) []byte {
	if l == nil {
		return nil
	}
	return masterwire.NewLocation(masterwire.LocationInput{
		Url: l.Url, PublicUrl: l.PublicUrl, GrpcPort: l.GrpcPort, DataCenter: l.DataCenter,
	})
}

func locationFromView(v masterwire.Location) *master_pb.Location {
	return &master_pb.Location{
		Url: v.Url(), PublicUrl: v.PublicUrl(), GrpcPort: v.GrpcPort(), DataCenter: v.DataCenter(),
	}
}

func locationsToWire(ls []*master_pb.Location) [][]byte {
	if len(ls) == 0 {
		return nil
	}
	out := make([][]byte, len(ls))
	for i, l := range ls {
		out[i] = locationToWire(l)
	}
	return out
}

// --- VolumeInformationMessage / VolumeShortInformationMessage / VolumeEcShardInformationMessage ---

func volumeInfoToWire(m *master_pb.VolumeInformationMessage) []byte {
	if m == nil {
		return nil
	}
	return masterwire.NewVolumeInformationMessage(masterwire.VolumeInformationMessageInput{
		Id:                m.Id,
		Size:              m.Size,
		Collection:        m.Collection,
		FileCount:         m.FileCount,
		DeleteCount:       m.DeleteCount,
		DeletedByteCount:  m.DeletedByteCount,
		ReadOnly:          m.ReadOnly,
		ReplicaPlacement:  m.ReplicaPlacement,
		Version:           m.Version,
		Ttl:               m.Ttl,
		CompactRevision:   m.CompactRevision,
		ModifiedAtSecond:  m.ModifiedAtSecond,
		RemoteStorageName: m.RemoteStorageName,
		RemoteStorageKey:  m.RemoteStorageKey,
		DiskType:          m.DiskType,
		DiskId:            m.DiskId,
	})
}

func volumeInfoFromView(v masterwire.VolumeInformationMessage) *master_pb.VolumeInformationMessage {
	return &master_pb.VolumeInformationMessage{
		Id:                v.Id(),
		Size:              v.Size(),
		Collection:        v.Collection(),
		FileCount:         v.FileCount(),
		DeleteCount:       v.DeleteCount(),
		DeletedByteCount:  v.DeletedByteCount(),
		ReadOnly:          v.ReadOnly(),
		ReplicaPlacement:  v.ReplicaPlacement(),
		Version:           v.Version(),
		Ttl:               v.Ttl(),
		CompactRevision:   v.CompactRevision(),
		ModifiedAtSecond:  v.ModifiedAtSecond(),
		RemoteStorageName: v.RemoteStorageName(),
		RemoteStorageKey:  v.RemoteStorageKey(),
		DiskType:          v.DiskType(),
		DiskId:            v.DiskId(),
	}
}

func volumeShortInfoToWire(m *master_pb.VolumeShortInformationMessage) []byte {
	if m == nil {
		return nil
	}
	return masterwire.NewVolumeShortInformationMessage(masterwire.VolumeShortInformationMessageInput{
		Id:               m.Id,
		Collection:       m.Collection,
		ReplicaPlacement: m.ReplicaPlacement,
		Version:          m.Version,
		Ttl:              m.Ttl,
		DiskType:         m.DiskType,
		DiskId:           m.DiskId,
	})
}

func volumeShortInfoFromView(v masterwire.VolumeShortInformationMessage) *master_pb.VolumeShortInformationMessage {
	return &master_pb.VolumeShortInformationMessage{
		Id:               v.Id(),
		Collection:       v.Collection(),
		ReplicaPlacement: v.ReplicaPlacement(),
		Version:          v.Version(),
		Ttl:              v.Ttl(),
		DiskType:         v.DiskType(),
		DiskId:           v.DiskId(),
	}
}

func ecShardInfoToWire(m *master_pb.VolumeEcShardInformationMessage) []byte {
	if m == nil {
		return nil
	}
	return masterwire.NewVolumeEcShardInformationMessage(masterwire.VolumeEcShardInformationMessageInput{
		Id:          m.Id,
		Collection:  m.Collection,
		EcIndexBits: m.EcIndexBits,
		DiskType:    m.DiskType,
		ExpireAtSec: m.ExpireAtSec,
		DiskId:      m.DiskId,
		ShardSizes:  m.ShardSizes,
		FileCount:   m.FileCount,
		DeleteCount: m.DeleteCount,
		EncodeTsNs:  m.EncodeTsNs,
	})
}

func ecShardInfoFromView(v masterwire.VolumeEcShardInformationMessage) *master_pb.VolumeEcShardInformationMessage {
	shards := make([]int64, v.ShardSizesLen())
	for i := range shards {
		shards[i] = v.ShardSizeAt(i)
	}
	return &master_pb.VolumeEcShardInformationMessage{
		Id:          v.Id(),
		Collection:  v.Collection(),
		EcIndexBits: v.EcIndexBits(),
		DiskType:    v.DiskType(),
		ExpireAtSec: v.ExpireAtSec(),
		DiskId:      v.DiskId(),
		ShardSizes:  shards,
		FileCount:   v.FileCount(),
		DeleteCount: v.DeleteCount(),
		EncodeTsNs:  v.EncodeTsNs(),
	}
}

// --- DiskInfo (and the string->DiskInfo map shared by every topology node) ---

func diskInfoToWire(d *master_pb.DiskInfo) []byte {
	if d == nil {
		return nil
	}
	volumeInfos := make([][]byte, len(d.VolumeInfos))
	for i, vi := range d.VolumeInfos {
		volumeInfos[i] = volumeInfoToWire(vi)
	}
	ecShardInfos := make([][]byte, len(d.EcShardInfos))
	for i, ei := range d.EcShardInfos {
		ecShardInfos[i] = ecShardInfoToWire(ei)
	}
	return masterwire.NewDiskInfo(masterwire.DiskInfoInput{
		Type:              d.Type,
		VolumeCount:       d.VolumeCount,
		MaxVolumeCount:    d.MaxVolumeCount,
		FreeVolumeCount:   d.FreeVolumeCount,
		ActiveVolumeCount: d.ActiveVolumeCount,
		VolumeInfos:       volumeInfos,
		EcShardInfos:      ecShardInfos,
		RemoteVolumeCount: d.RemoteVolumeCount,
		DiskId:            d.DiskId,
		Tags:              d.Tags,
	})
}

func diskInfoFromView(v masterwire.DiskInfo) *master_pb.DiskInfo {
	d := &master_pb.DiskInfo{
		Type:              v.Type(),
		VolumeCount:       v.VolumeCount(),
		MaxVolumeCount:    v.MaxVolumeCount(),
		FreeVolumeCount:   v.FreeVolumeCount(),
		ActiveVolumeCount: v.ActiveVolumeCount(),
		RemoteVolumeCount: v.RemoteVolumeCount(),
		DiskId:            v.DiskId(),
	}
	if n := v.VolumeInfosLen(); n > 0 {
		d.VolumeInfos = make([]*master_pb.VolumeInformationMessage, 0, n)
		for i := 0; i < n; i++ {
			if vi, ok := v.VolumeInfoAt(i); ok {
				d.VolumeInfos = append(d.VolumeInfos, volumeInfoFromView(vi))
			}
		}
	}
	if n := v.EcShardInfosLen(); n > 0 {
		d.EcShardInfos = make([]*master_pb.VolumeEcShardInformationMessage, 0, n)
		for i := 0; i < n; i++ {
			if ei, ok := v.EcShardInfoAt(i); ok {
				d.EcShardInfos = append(d.EcShardInfos, ecShardInfoFromView(ei))
			}
		}
	}
	if n := v.TagsLen(); n > 0 {
		d.Tags = make([]string, n)
		for i := 0; i < n; i++ {
			d.Tags[i] = v.TagAt(i)
		}
	}
	return d
}

func diskInfoMapToWire(m map[string]*master_pb.DiskInfo) []masterwire.StringMsgEntry {
	if len(m) == 0 {
		return nil
	}
	out := make([]masterwire.StringMsgEntry, 0, len(m))
	for k, v := range m {
		out = append(out, masterwire.StringMsgEntry{Key: k, Value: diskInfoToWire(v)})
	}
	return out
}

func diskInfoMapFromView(n int, at func(int) masterwire.StringMsgEntry) map[string]*master_pb.DiskInfo {
	if n == 0 {
		return nil
	}
	out := make(map[string]*master_pb.DiskInfo, n)
	for i := 0; i < n; i++ {
		e := at(i)
		dv, err := masterwire.WrapDiskInfo(e.Value)
		if err != nil {
			continue
		}
		out[e.Key] = diskInfoFromView(dv)
	}
	return out
}

// --- topology graph: DataNodeInfo -> RackInfo -> DataCenterInfo -> TopologyInfo ---

func dataNodeInfoToWire(d *master_pb.DataNodeInfo) []byte {
	if d == nil {
		return nil
	}
	return masterwire.NewDataNodeInfo(masterwire.DataNodeInfoInput{
		Id:        d.Id,
		DiskInfos: diskInfoMapToWire(d.DiskInfos),
		GrpcPort:  d.GrpcPort,
		Address:   d.Address,
	})
}

func dataNodeInfoFromView(v masterwire.DataNodeInfo) *master_pb.DataNodeInfo {
	return &master_pb.DataNodeInfo{
		Id:        v.Id(),
		DiskInfos: diskInfoMapFromView(v.DiskInfosLen(), v.DiskInfoAt),
		GrpcPort:  v.GrpcPort(),
		Address:   v.Address(),
	}
}

func rackInfoToWire(r *master_pb.RackInfo) []byte {
	if r == nil {
		return nil
	}
	nodes := make([][]byte, len(r.DataNodeInfos))
	for i, n := range r.DataNodeInfos {
		nodes[i] = dataNodeInfoToWire(n)
	}
	return masterwire.NewRackInfo(masterwire.RackInfoInput{
		Id:            r.Id,
		DataNodeInfos: nodes,
		DiskInfos:     diskInfoMapToWire(r.DiskInfos),
	})
}

func rackInfoFromView(v masterwire.RackInfo) *master_pb.RackInfo {
	r := &master_pb.RackInfo{
		Id:        v.Id(),
		DiskInfos: diskInfoMapFromView(v.DiskInfosLen(), v.DiskInfoAt),
	}
	if n := v.DataNodeInfosLen(); n > 0 {
		r.DataNodeInfos = make([]*master_pb.DataNodeInfo, 0, n)
		for i := 0; i < n; i++ {
			if d, ok := v.DataNodeInfoAt(i); ok {
				r.DataNodeInfos = append(r.DataNodeInfos, dataNodeInfoFromView(d))
			}
		}
	}
	return r
}

func dataCenterInfoToWire(dc *master_pb.DataCenterInfo) []byte {
	if dc == nil {
		return nil
	}
	racks := make([][]byte, len(dc.RackInfos))
	for i, r := range dc.RackInfos {
		racks[i] = rackInfoToWire(r)
	}
	return masterwire.NewDataCenterInfo(masterwire.DataCenterInfoInput{
		Id:        dc.Id,
		RackInfos: racks,
		DiskInfos: diskInfoMapToWire(dc.DiskInfos),
	})
}

func dataCenterInfoFromView(v masterwire.DataCenterInfo) *master_pb.DataCenterInfo {
	dc := &master_pb.DataCenterInfo{
		Id:        v.Id(),
		DiskInfos: diskInfoMapFromView(v.DiskInfosLen(), v.DiskInfoAt),
	}
	if n := v.RackInfosLen(); n > 0 {
		dc.RackInfos = make([]*master_pb.RackInfo, 0, n)
		for i := 0; i < n; i++ {
			if r, ok := v.RackInfoAt(i); ok {
				dc.RackInfos = append(dc.RackInfos, rackInfoFromView(r))
			}
		}
	}
	return dc
}

func topologyInfoToWire(t *master_pb.TopologyInfo) []byte {
	if t == nil {
		return nil
	}
	dcs := make([][]byte, len(t.DataCenterInfos))
	for i, dc := range t.DataCenterInfos {
		dcs[i] = dataCenterInfoToWire(dc)
	}
	return masterwire.NewTopologyInfo(masterwire.TopologyInfoInput{
		Id:              t.Id,
		DataCenterInfos: dcs,
		DiskInfos:       diskInfoMapToWire(t.DiskInfos),
	})
}

func topologyInfoFromView(v masterwire.TopologyInfo) *master_pb.TopologyInfo {
	t := &master_pb.TopologyInfo{
		Id:        v.Id(),
		DiskInfos: diskInfoMapFromView(v.DiskInfosLen(), v.DiskInfoAt),
	}
	if n := v.DataCenterInfosLen(); n > 0 {
		t.DataCenterInfos = make([]*master_pb.DataCenterInfo, 0, n)
		for i := 0; i < n; i++ {
			if dc, ok := v.DataCenterInfoAt(i); ok {
				t.DataCenterInfos = append(t.DataCenterInfos, dataCenterInfoFromView(dc))
			}
		}
	}
	return t
}

// --- StorageBackend ---

func storageBackendToWire(s *master_pb.StorageBackend) []byte {
	if s == nil {
		return nil
	}
	props := make([]masterwire.StringStringEntry, 0, len(s.Properties))
	for k, v := range s.Properties {
		props = append(props, masterwire.StringStringEntry{Key: k, Value: v})
	}
	return masterwire.NewStorageBackend(masterwire.StorageBackendInput{
		Type: s.Type, Id: s.Id, Properties: props,
	})
}

func storageBackendsToWire(ss []*master_pb.StorageBackend) [][]byte {
	if len(ss) == 0 {
		return nil
	}
	out := make([][]byte, len(ss))
	for i, s := range ss {
		out[i] = storageBackendToWire(s)
	}
	return out
}

func storageBackendFromView(v masterwire.StorageBackend) *master_pb.StorageBackend {
	s := &master_pb.StorageBackend{Type: v.Type(), Id: v.Id()}
	if n := v.PropertiesLen(); n > 0 {
		s.Properties = make(map[string]string, n)
		for i := 0; i < n; i++ {
			e := v.PropertyAt(i)
			s.Properties[e.Key] = e.Value
		}
	}
	return s
}

// --- DiskTag ---

func diskTagToWire(d *master_pb.DiskTag) []byte {
	if d == nil {
		return nil
	}
	return masterwire.NewDiskTag(masterwire.DiskTagInput{DiskId: d.DiskId, Tags: d.Tags})
}

func diskTagFromView(v masterwire.DiskTag) *master_pb.DiskTag {
	tags := make([]string, v.TagsLen())
	for i := range tags {
		tags[i] = v.TagAt(i)
	}
	return &master_pb.DiskTag{DiskId: v.DiskId(), Tags: tags}
}

// --- VolumeServerState (Heartbeat.state) ---

func stateToWire(st *volume_server_pb.VolumeServerState) []byte {
	if st == nil {
		return nil
	}
	return volume_serverwire.NewVolumeServerState(volume_serverwire.VolumeServerStateInput{
		Maintenance: st.Maintenance,
		Version:     st.Version,
	})
}

func stateFromView(v volume_serverwire.VolumeServerState) *volume_server_pb.VolumeServerState {
	return &volume_server_pb.VolumeServerState{Maintenance: v.Maintenance(), Version: v.Version()}
}

// --- VolumeLocation / ClusterNodeUpdate / LockRingUpdate (KeepConnectedResponse) ---

func volumeLocationToWire(l *master_pb.VolumeLocation) []byte {
	if l == nil {
		return nil
	}
	return masterwire.NewVolumeLocation(masterwire.VolumeLocationInput{
		Url:           l.Url,
		PublicUrl:     l.PublicUrl,
		NewVids:       l.NewVids,
		DeletedVids:   l.DeletedVids,
		Leader:        l.Leader,
		DataCenter:    l.DataCenter,
		GrpcPort:      l.GrpcPort,
		NewEcVids:     l.NewEcVids,
		DeletedEcVids: l.DeletedEcVids,
	})
}

func volumeLocationFromView(v masterwire.VolumeLocation) *master_pb.VolumeLocation {
	uint32s := func(n int, at func(int) uint32) []uint32 {
		if n == 0 {
			return nil
		}
		out := make([]uint32, n)
		for i := 0; i < n; i++ {
			out[i] = at(i)
		}
		return out
	}
	return &master_pb.VolumeLocation{
		Url:           v.Url(),
		PublicUrl:     v.PublicUrl(),
		NewVids:       uint32s(v.NewVidsLen(), v.NewVidAt),
		DeletedVids:   uint32s(v.DeletedVidsLen(), v.DeletedVidAt),
		Leader:        v.Leader(),
		DataCenter:    v.DataCenter(),
		GrpcPort:      v.GrpcPort(),
		NewEcVids:     uint32s(v.NewEcVidsLen(), v.NewEcVidAt),
		DeletedEcVids: uint32s(v.DeletedEcVidsLen(), v.DeletedEcVidAt),
	}
}

func clusterNodeUpdateToWire(c *master_pb.ClusterNodeUpdate) []byte {
	if c == nil {
		return nil
	}
	return masterwire.NewClusterNodeUpdate(masterwire.ClusterNodeUpdateInput{
		NodeType: c.NodeType, Address: c.Address, IsAdd: c.IsAdd,
		FilerGroup: c.FilerGroup, CreatedAtNs: c.CreatedAtNs,
	})
}

func clusterNodeUpdateFromView(v masterwire.ClusterNodeUpdate) *master_pb.ClusterNodeUpdate {
	return &master_pb.ClusterNodeUpdate{
		NodeType: v.NodeType(), Address: v.Address(), IsAdd: v.IsAdd(),
		FilerGroup: v.FilerGroup(), CreatedAtNs: v.CreatedAtNs(),
	}
}

func lockRingUpdateToWire(l *master_pb.LockRingUpdate) []byte {
	if l == nil {
		return nil
	}
	return masterwire.NewLockRingUpdate(masterwire.LockRingUpdateInput{
		FilerGroup: l.FilerGroup, Servers: l.Servers, Version: l.Version,
	})
}

func lockRingUpdateFromView(v masterwire.LockRingUpdate) *master_pb.LockRingUpdate {
	l := &master_pb.LockRingUpdate{FilerGroup: v.FilerGroup(), Version: v.Version()}
	if n := v.ServersLen(); n > 0 {
		l.Servers = make([]string, n)
		for i := 0; i < n; i++ {
			l.Servers[i] = v.ServerAt(i)
		}
	}
	return l
}
