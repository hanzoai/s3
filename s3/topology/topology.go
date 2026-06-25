package topology

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/storage/types"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/pb/master_pb"
	"github.com/hanzoai/s3/s3/sequence"
	"github.com/hanzoai/s3/s3/stats"
	"github.com/hanzoai/s3/s3/storage"
	"github.com/hanzoai/s3/s3/storage/needle"
	"github.com/hanzoai/s3/s3/storage/super_block"
	"github.com/hanzoai/s3/s3/util"
)

const (
	// WarmupPulseMultiplier is the number of heartbeat intervals to wait after
	// a leader change before treating volume lookup misses as definitive.
	WarmupPulseMultiplier = 3
)

type Topology struct {
	vacuumLockCounter int64
	NodeImpl

	collectionMap  *util.ConcurrentReadMap
	ecShardMap     map[needle.VolumeId]*EcShardLocations
	ecShardMapLock sync.RWMutex

	pulse int64

	volumeSizeLimit          uint64
	replicationAsMin         bool
	vacuumDisabledByOperator atomic.Bool // true when operator manually disables vacuum
	vacuumDisabledByPlugin   atomic.Bool // true when disabled by the vacuum plugin monitor
	adminServerConnectedFunc func() bool // optional callback to check admin server presence

	Sequence sequence.Sequencer

	chanFullVolumes    chan storage.VolumeInfo
	chanCrowdedVolumes chan storage.VolumeInfo

	Configuration *Configuration

	// Coordinator is the leaderless replacement for the raft leader + log:
	// it serializes volume-id / file-id allocation and pins the single writer
	// with no election. NewTopology installs an in-process LocalCoordinator;
	// command/master.go may swap in the schain-backed one for production.
	Coordinator Coordinator

	UuidAccessLock sync.RWMutex
	UuidMap        map[string][]string

	topologyId     string
	topologyIdLock sync.RWMutex

	lastLeaderChangeTime     time.Time
	hadVolumesAtLeaderChange bool
	lastLeaderChangeTimeLock sync.RWMutex

	// dataNodeIndex is an address -> *DataNode lookup so callers (e.g. the
	// Ping admission gate) do not have to walk every dc/rack/node tier on
	// every request. Keys use the canonical http form returned by
	// pb.ServerAddress.ToHttpAddress so a target like "1.2.3.4:8080" finds
	// the same node whether or not the grpc port suffix is present.
	dataNodeIndex     map[string]*DataNode
	dataNodeIndexLock sync.RWMutex
}

func NewTopology(id string, seq sequence.Sequencer, volumeSizeLimit uint64, pulse int, replicationAsMin bool) *Topology {
	t := &Topology{}
	t.id = NodeId(id)
	t.nodeType = "Topology"
	t.NodeImpl.value = t
	t.diskUsages = newDiskUsages()
	t.children = make(map[NodeId]Node)
	t.capacityReservations = newCapacityReservations()
	t.collectionMap = util.NewConcurrentReadMap()
	t.ecShardMap = make(map[needle.VolumeId]*EcShardLocations)
	t.pulse = int64(pulse)
	t.volumeSizeLimit = volumeSizeLimit
	t.replicationAsMin = replicationAsMin

	t.Sequence = seq

	t.chanFullVolumes = make(chan storage.VolumeInfo)
	t.chanCrowdedVolumes = make(chan storage.VolumeInfo)

	t.Configuration = &Configuration{}
	t.dataNodeIndex = make(map[string]*DataNode)

	// Default to the in-process pinned-writer coordinator: a standalone master
	// (and every test topology) is the sole writer and allocates freely.
	t.Coordinator = NewLocalCoordinator(t)

	return t
}

// LookupDataNodeByAddress returns the registered DataNode that serves addr,
// or nil if no such node has been observed. Lookup is O(1) and uses the
// canonical http form of the address so callers that pass either
// "host:port" or "host:port.grpc" find the same node.
func (t *Topology) LookupDataNodeByAddress(addr pb.ServerAddress) *DataNode {
	if addr == "" {
		return nil
	}
	t.dataNodeIndexLock.RLock()
	defer t.dataNodeIndexLock.RUnlock()
	if t.dataNodeIndex == nil {
		return nil
	}
	return t.dataNodeIndex[addr.ToHttpAddress()]
}

// registerDataNodeAddress records dn in the address index under its current
// http address. Callers must invoke unregisterDataNodeAddress with the prior
// address whenever a node's Ip or Port changes (e.g. k8s pod reschedule).
func (t *Topology) registerDataNodeAddress(dn *DataNode) {
	if dn == nil {
		return
	}
	key := dn.ServerAddress().ToHttpAddress()
	if key == "" {
		return
	}
	t.dataNodeIndexLock.Lock()
	defer t.dataNodeIndexLock.Unlock()
	if t.dataNodeIndex == nil {
		t.dataNodeIndex = make(map[string]*DataNode)
	}
	t.dataNodeIndex[key] = dn
}

// unregisterDataNodeAddress removes the index entry for addr, but only when
// the entry still points at dn. The conditional guard avoids dropping a
// freshly re-registered node whose address happens to alias the one being
// removed (e.g. legacy id transitions or a fast restart).
func (t *Topology) unregisterDataNodeAddress(addr pb.ServerAddress, dn *DataNode) {
	if addr == "" {
		return
	}
	key := addr.ToHttpAddress()
	if key == "" {
		return
	}
	t.dataNodeIndexLock.Lock()
	defer t.dataNodeIndexLock.Unlock()
	if existing, ok := t.dataNodeIndex[key]; ok && (dn == nil || existing == dn) {
		delete(t.dataNodeIndex, key)
	}
}

func (t *Topology) IsChildLocked() (bool, error) {
	if t.IsLocked() {
		return true, errors.New("topology is locked")
	}
	for _, dcNode := range t.Children() {
		if dcNode.IsLocked() {
			return true, fmt.Errorf("topology child %s is locked", dcNode.String())
		}
		for _, rackNode := range dcNode.Children() {
			if rackNode.IsLocked() {
				return true, fmt.Errorf("dc %s child %s is locked", dcNode.String(), rackNode.String())
			}
			for _, dataNode := range rackNode.Children() {
				if dataNode.IsLocked() {
					return true, fmt.Errorf("rack %s child %s is locked", rackNode.String(), dataNode.Id())
				}
			}
		}
	}
	return false, nil
}

// SetLastLeaderChangeTime records the time of the most recent leader transition.
// It also snapshots whether the topology already had known volumes at that
// moment. IsWarmingUp uses the snapshot instead of the live MaxVolumeId so a
// fresh cluster that happens to grow its first volume inside the warmup window
// does not retroactively flip into "warming up" state — there is no prior
// topology to wait for on a bootstrap.
func (t *Topology) SetLastLeaderChangeTime(ts time.Time) {
	hadVolumes := t.GetMaxVolumeId() > 0
	t.lastLeaderChangeTimeLock.Lock()
	defer t.lastLeaderChangeTimeLock.Unlock()
	t.lastLeaderChangeTime = ts
	t.hadVolumesAtLeaderChange = hadVolumes
}

// GetLastLeaderChangeTime returns the time of the most recent leader transition.
func (t *Topology) GetLastLeaderChangeTime() time.Time {
	t.lastLeaderChangeTimeLock.RLock()
	defer t.lastLeaderChangeTimeLock.RUnlock()
	return t.lastLeaderChangeTime
}

// IsWarmingUp returns true if the master recently became leader and may not yet
// have a complete topology. After a leader change or restart, volume servers need
// up to WarmupPulseMultiplier heartbeat intervals to reconnect and report their volumes.
// Returns false on a fresh cluster start — i.e. when no volumes existed at the
// time of the leader change — since there is no prior topology state to wait for.
// Checking the *live* MaxVolumeId here would make a bootstrapping cluster flip
// into warming-up the moment its first volume is grown, which manifested as a
// 15-second window of spurious Unavailable errors on AssignVolume for workloads
// that start writing immediately (see #8777).
func (t *Topology) IsWarmingUp() bool {
	t.lastLeaderChangeTimeLock.RLock()
	lastChange := t.lastLeaderChangeTime
	hadVolumes := t.hadVolumesAtLeaderChange
	t.lastLeaderChangeTimeLock.RUnlock()
	if !hadVolumes || lastChange.IsZero() {
		return false
	}
	return time.Since(lastChange) < t.WarmupDuration()
}

// WarmupDuration returns the configured warmup duration based on pulse interval.
func (t *Topology) WarmupDuration() time.Duration {
	return time.Duration(t.pulse*WarmupPulseMultiplier) * time.Second
}

// RemainingWarmupDuration returns how much warmup time is left, or 0 if not warming up.
func (t *Topology) RemainingWarmupDuration() time.Duration {
	if !t.IsWarmingUp() {
		return 0
	}
	remaining := t.WarmupDuration() - time.Since(t.GetLastLeaderChangeTime())
	if remaining < 0 {
		return 0
	}
	return remaining
}

// IsWriter reports whether this master is the pinned writer — the single
// member allowed to run the serialization-critical metadata ops (allocation,
// volume growth, vacuum, collection/volume admin). It replaces the raft
// IsLeader() gate; the writer is pinned deterministically by the Coordinator
// with no election.
func (t *Topology) IsWriter() bool {
	return t.Coordinator != nil && t.Coordinator.IsWriter()
}

// Writer returns the pinned-writer address (for read redirect / proxy), or ""
// if no membership is known yet. Replaces the raft Leader() lookup; there is no
// election to wait on, so it returns immediately.
func (t *Topology) Writer() (pb.ServerAddress, error) {
	if t.Coordinator == nil {
		return "", errors.New("coordinator not ready yet")
	}
	return t.Coordinator.Writer(), nil
}

// MaybeWriter is the non-erroring writer lookup. Kept distinct from Writer for
// call-site parity with the former MaybeLeader/Leader pair.
func (t *Topology) MaybeWriter() (pb.ServerAddress, error) {
	if t.Coordinator == nil {
		return "", errors.New("coordinator not ready yet")
	}
	return t.Coordinator.Writer(), nil
}

func (t *Topology) Lookup(collection string, vid needle.VolumeId) (dataNodes []*DataNode) {
	// maybe an issue if lots of collections?
	if collection == "" {
		for _, c := range t.collectionMap.Items() {
			if list := c.(*Collection).Lookup(vid); list != nil {
				return list
			}
		}
	} else {
		if c, ok := t.collectionMap.Find(collection); ok {
			return c.(*Collection).Lookup(vid)
		}
	}

	if locations, found := t.LookupEcShards(vid); found {
		for _, loc := range locations.Locations {
			dataNodes = append(dataNodes, loc...)
		}
		return dataNodes
	}

	return nil
}

// NextVolumeId allocates the next volume id through the Coordinator. The
// Coordinator — not a raft leader — serializes allocation and guarantees global
// uniqueness; on a non-writer member it returns ErrNotWriter and the caller
// redirects to the pinned writer.
func (t *Topology) NextVolumeId() (needle.VolumeId, error) {
	if t.Coordinator == nil {
		return 0, errors.New("coordinator not ready yet")
	}
	_, end, err := t.Coordinator.AllocateVolumeId(1)
	if err != nil {
		return 0, err
	}
	return end, nil
}

// DefaultNeedleSizeEstimate is the fallback per-file-ID size estimate when
// the client does not provide an expected data size.
const DefaultNeedleSizeEstimate uint64 = 1024 * 1024 // 1 MB

func (t *Topology) PickForWrite(requestedCount uint64, option *VolumeGrowOption, volumeLayout *VolumeLayout, expectedDataSize uint64) (fileId string, count uint64, volumeLocationList *VolumeLocationList, shouldGrow bool, err error) {
	var vid needle.VolumeId
	vid, count, volumeLocationList, shouldGrow, err = volumeLayout.PickForWrite(requestedCount, option)
	if err != nil {
		return "", 0, nil, shouldGrow, fmt.Errorf("failed to find writable volumes for collection:%s replication:%s ttl:%s error: %v", option.Collection, option.ReplicaPlacement.String(), option.Ttl.String(), err)
	}
	if volumeLocationList == nil || volumeLocationList.Length() == 0 {
		return "", 0, nil, shouldGrow, fmt.Errorf("%s available for collection:%s replication:%s ttl:%s", NoWritableVolumes, option.Collection, option.ReplicaPlacement.String(), option.Ttl.String())
	}
	// Track estimated assigned bytes to spread load between heartbeats.
	// Use the client hint if provided, otherwise fall back to 1MB estimate.
	sizePerFile := DefaultNeedleSizeEstimate
	if expectedDataSize > 0 {
		sizePerFile = expectedDataSize
	}
	pendingBytes := min(uint64(count)*sizePerFile, uint64(math.MaxInt64))
	if volumeLayout.RecordAssign(vid, int64(pendingBytes)) {
		volumeLayout.AdjustActiveVolumeCountForFull(vid)
	}
	nextFileId, err := t.Coordinator.NextFileId(requestedCount)
	if err != nil {
		return "", 0, nil, shouldGrow, err
	}
	fileId = needle.NewFileId(vid, nextFileId, rand.Uint32()).String()
	return fileId, count, volumeLocationList, shouldGrow, nil
}

func (t *Topology) GetVolumeLayout(collectionName string, rp *super_block.ReplicaPlacement, ttl *needle.TTL, diskType types.DiskType) *VolumeLayout {
	return t.collectionMap.Get(collectionName, func() interface{} {
		return NewCollection(collectionName, t.volumeSizeLimit, t.replicationAsMin)
	}).(*Collection).GetOrCreateVolumeLayout(rp, ttl, diskType)
}

func (t *Topology) ListCollections(includeNormalVolumes, includeEcVolumes bool) (ret []string) {
	found := make(map[string]bool)

	if includeNormalVolumes {
		t.collectionMap.RLock()
		for _, c := range t.collectionMap.Items() {
			found[c.(*Collection).Name] = true
		}
		t.collectionMap.RUnlock()
	}

	if includeEcVolumes {
		t.ecShardMapLock.RLock()
		for _, ecVolumeLocation := range t.ecShardMap {
			found[ecVolumeLocation.Collection] = true
		}
		t.ecShardMapLock.RUnlock()
	}

	for k := range found {
		ret = append(ret, k)
	}
	slices.Sort(ret)

	return ret
}

func (t *Topology) FindCollection(collectionName string) (*Collection, bool) {
	c, hasCollection := t.collectionMap.Find(collectionName)
	if !hasCollection {
		return nil, false
	}
	return c.(*Collection), hasCollection
}

func (t *Topology) DeleteCollection(collectionName string) {
	t.collectionMap.Delete(collectionName)
}

func (t *Topology) DeleteLayout(collectionName string, rp *super_block.ReplicaPlacement, ttl *needle.TTL, diskType types.DiskType) {
	collection, found := t.FindCollection(collectionName)
	if !found {
		return
	}
	collection.DeleteVolumeLayout(rp, ttl, diskType)
	if len(collection.storageType2VolumeLayout.Items()) == 0 {
		t.DeleteCollection(collectionName)
	}
}

func (t *Topology) RegisterVolumeLayout(v storage.VolumeInfo, dn *DataNode) {
	diskType := types.ToDiskType(v.DiskType)
	vl := t.GetVolumeLayout(v.Collection, v.ReplicaPlacement, v.Ttl, diskType)
	vl.RegisterVolume(&v, dn)
	vl.EnsureCorrectWritables(&v)
}

func (t *Topology) UnRegisterVolumeLayout(v storage.VolumeInfo, dn *DataNode) {
	glog.Infof("removing volume info: %+v from %v", v, dn.id)
	if v.ReplicaPlacement.GetCopyCount() > 1 {
		stats.MasterReplicaPlacementMismatch.WithLabelValues(v.Collection, v.Id.String()).Set(0)
	}
	diskType := types.ToDiskType(v.DiskType)
	volumeLayout := t.GetVolumeLayout(v.Collection, v.ReplicaPlacement, v.Ttl, diskType)
	volumeLayout.UnRegisterVolume(&v, dn)
	if volumeLayout.isEmpty() {
		t.DeleteLayout(v.Collection, v.ReplicaPlacement, v.Ttl, diskType)
	}
}

func (t *Topology) DataCenterExists(dcName string) bool {
	return dcName == "" || t.GetDataCenter(dcName) != nil
}

func (t *Topology) GetDataCenter(dcName string) (dc *DataCenter) {
	t.RLock()
	defer t.RUnlock()
	for _, c := range t.children {
		dc = c.(*DataCenter)
		if string(dc.Id()) == dcName {
			return dc
		}
	}
	return dc
}

func (t *Topology) GetOrCreateDataCenter(dcName string) *DataCenter {
	t.Lock()
	defer t.Unlock()
	for _, c := range t.children {
		dc := c.(*DataCenter)
		if string(dc.Id()) == dcName {
			return dc
		}
	}
	dc := NewDataCenter(dcName)
	t.doLinkChildNode(dc)
	return dc
}

func (t *Topology) ListDataCenters() (dcs []string) {
	t.RLock()
	defer t.RUnlock()
	for _, c := range t.children {
		dcs = append(dcs, string(c.(*DataCenter).Id()))
	}
	return dcs
}

func (t *Topology) ListDCAndRacks() (dcs map[NodeId][]NodeId) {
	t.RLock()
	defer t.RUnlock()
	dcs = make(map[NodeId][]NodeId)
	for _, dcNode := range t.children {
		dcNodeId := dcNode.(*DataCenter).Id()
		for _, rackNode := range dcNode.Children() {
			dcs[dcNodeId] = append(dcs[dcNodeId], rackNode.(*Rack).Id())
		}
	}
	return dcs
}

func (t *Topology) SyncDataNodeRegistration(volumes []*master_pb.VolumeInformationMessage, dn *DataNode) (newVolumes, deletedVolumes []storage.VolumeInfo) {
	// convert into in memory struct storage.VolumeInfo
	var volumeInfos []storage.VolumeInfo
	for _, v := range volumes {
		if vi, err := storage.NewVolumeInfo(v); err == nil {
			volumeInfos = append(volumeInfos, vi)
		} else {
			glog.V(0).Infof("Fail to convert joined volume information: %v", err)
		}
	}
	// find out the delta volumes
	var changedVolumes []storage.VolumeInfo
	newVolumes, deletedVolumes, changedVolumes = dn.UpdateVolumes(volumeInfos)
	for _, v := range newVolumes {
		t.RegisterVolumeLayout(v, dn)
	}
	for _, v := range deletedVolumes {
		t.UnRegisterVolumeLayout(v, dn)
	}
	for _, v := range changedVolumes {
		diskType := types.ToDiskType(v.DiskType)
		vl := t.GetVolumeLayout(v.Collection, v.ReplicaPlacement, v.Ttl, diskType)
		vl.EnsureCorrectWritables(&v)
	}
	// Update effective sizes for all reported volumes (decay pending estimates).
	// If decay brings a volume eagerly removed by RecordAssign back under the
	// writable threshold, restore the matching activeVolumeCount.
	for _, v := range volumeInfos {
		if v.ReplicaPlacement == nil {
			continue
		}
		diskType := types.ToDiskType(v.DiskType)
		vl := t.GetVolumeLayout(v.Collection, v.ReplicaPlacement, v.Ttl, diskType)
		// Self-heal: a volume reported by the data node but missing from the
		// lookup index is re-registered. This repairs the split left by a
		// disconnect/reconnect race, where UnRegisterDataNode dropped the volume
		// from vid2location but the reconnecting full heartbeat skipped it
		// (still in the disk map, so UpdateVolumes did not report it as new).
		// Without this, the volume stays visible in volume.list/admin UI yet
		// LookupVolume returns "volume id not found".
		if !vl.HasDataNode(v.Id, dn) {
			// vl is already resolved above; call it directly instead of
			// RegisterVolumeLayout, which would repeat the GetVolumeLayout lookup.
			vl.RegisterVolume(&v, dn)
			vl.EnsureCorrectWritables(&v)
		}
		if vl.UpdateVolumeSize(v.Id, v.Size, v.CompactRevision) {
			vl.AdjustActiveVolumeCountAfterRecovery(v.Id)
		}
	}
	return
}

func (t *Topology) IncrementalSyncDataNodeRegistration(newVolumes, deletedVolumes []*master_pb.VolumeShortInformationMessage, dn *DataNode) {
	var newVis, oldVis []storage.VolumeInfo
	for _, v := range newVolumes {
		vi, err := storage.NewVolumeInfoFromShort(v)
		if err != nil {
			glog.V(0).Infof("NewVolumeInfoFromShort %v: %v", v, err)
			continue
		}
		newVis = append(newVis, vi)
	}
	for _, v := range deletedVolumes {
		vi, err := storage.NewVolumeInfoFromShort(v)
		if err != nil {
			glog.V(0).Infof("NewVolumeInfoFromShort %v: %v", v, err)
			continue
		}
		oldVis = append(oldVis, vi)
	}
	dn.DeltaUpdateVolumes(newVis, oldVis)

	for _, vi := range newVis {
		t.RegisterVolumeLayout(vi, dn)
	}
	for _, vi := range oldVis {
		t.UnRegisterVolumeLayout(vi, dn)
	}

	return
}

func (t *Topology) DataNodeRegistration(dcName, rackName string, dn *DataNode) {
	if dn.Parent() != nil {
		return
	}
	// registration to topo
	dc := t.GetOrCreateDataCenter(dcName)
	rack := dc.GetOrCreateRack(rackName)
	rack.LinkChildNode(dn)
	glog.Infof("[%s] reLink To topo  ", dn.Id())
}

// IsVacuumDisabled returns true if vacuum is disabled by either the
// operator or the plugin monitor.
func (t *Topology) IsVacuumDisabled() bool {
	return t.vacuumDisabledByOperator.Load() || t.vacuumDisabledByPlugin.Load()
}

// DisableVacuum is called by the operator (shell command / manual RPC).
// Only sets the operator flag; does not affect the plugin flag.
func (t *Topology) DisableVacuum() {
	glog.V(0).Infof("DisableVacuum (by operator)")
	t.vacuumDisabledByOperator.Store(true)
}

// EnableVacuum is called by the operator (shell command / manual RPC).
// Only clears the operator flag; does not affect the plugin flag.
func (t *Topology) EnableVacuum() {
	glog.V(0).Infof("EnableVacuum (by operator)")
	t.vacuumDisabledByOperator.Store(false)
}

// DisableVacuumByPlugin is called by the admin server's vacuum monitor
// when a vacuum plugin worker connects. Only sets the plugin flag.
func (t *Topology) DisableVacuumByPlugin() {
	glog.V(0).Infof("DisableVacuum (by plugin worker)")
	t.vacuumDisabledByPlugin.Store(true)
}

// EnableVacuumByPlugin is called by the admin server's vacuum monitor
// when a vacuum plugin worker disconnects. Only clears the plugin flag.
func (t *Topology) EnableVacuumByPlugin() {
	glog.V(0).Infof("EnableVacuum (by plugin worker)")
	t.vacuumDisabledByPlugin.Store(false)
}

// IsVacuumDisabledByPlugin returns whether the plugin monitor has disabled vacuum.
func (t *Topology) IsVacuumDisabledByPlugin() bool {
	return t.vacuumDisabledByPlugin.Load()
}

// SetAdminServerConnectedFunc sets an optional callback used by the vacuum
// safety net to detect when the admin server has disconnected.
func (t *Topology) SetAdminServerConnectedFunc(f func() bool) {
	t.adminServerConnectedFunc = f
}

func (t *Topology) GetTopologyId() string {
	t.topologyIdLock.RLock()
	defer t.topologyIdLock.RUnlock()
	return t.topologyId
}

func (t *Topology) SetTopologyId(topologyId string) {
	t.topologyIdLock.Lock()
	defer t.topologyIdLock.Unlock()
	if topologyId == "" {
		return
	}
	if t.topologyId == "" {
		t.topologyId = topologyId
		return
	}
	if t.topologyId != topologyId {
		glog.Fatalf("Split-brain detected! Current TopologyId is %s, but received %s. Stopping to prevent data corruption.", t.topologyId, topologyId)
	}
}
