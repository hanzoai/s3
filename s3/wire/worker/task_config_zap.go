// Code generated from worker.proto; DO NOT EDIT.

package workerwire

import (
	zap "github.com/zap-proto/go"
)

// --- VacuumTaskConfig ---
//
// Fixed payload: garbage_threshold(8 double) min_volume_age_hours(4)
// min_interval_seconds(4).

const (
	vacuumTaskConfigGarbageThresholdOff   = 0
	vacuumTaskConfigMinVolumeAgeHoursOff  = 8
	vacuumTaskConfigMinIntervalSecondsOff = 12
	vacuumTaskConfigSize                  = 16
)

// VacuumTaskConfig is a zero-copy view into a ZAP-encoded VacuumTaskConfig message.
type VacuumTaskConfig struct{ o zap.Object }

// WrapVacuumTaskConfig parses b and returns a typed view.
func WrapVacuumTaskConfig(b []byte) (VacuumTaskConfig, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VacuumTaskConfig{}, err
	}
	return VacuumTaskConfig{o: m.Root()}, nil
}

func (t VacuumTaskConfig) GarbageThreshold() float64 {
	return t.o.Float64(vacuumTaskConfigGarbageThresholdOff)
}
func (t VacuumTaskConfig) MinVolumeAgeHours() int32 {
	return t.o.Int32(vacuumTaskConfigMinVolumeAgeHoursOff)
}
func (t VacuumTaskConfig) MinIntervalSeconds() int32 {
	return t.o.Int32(vacuumTaskConfigMinIntervalSecondsOff)
}

// VacuumTaskConfigInput collects the field values for NewVacuumTaskConfig.
type VacuumTaskConfigInput struct {
	GarbageThreshold   float64
	MinVolumeAgeHours  int32
	MinIntervalSeconds int32
}

// NewVacuumTaskConfig builds a ZAP-encoded VacuumTaskConfig message from in.
func NewVacuumTaskConfig(in VacuumTaskConfigInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(vacuumTaskConfigSize)
	ob.SetFloat64(vacuumTaskConfigGarbageThresholdOff, in.GarbageThreshold)
	ob.SetInt32(vacuumTaskConfigMinVolumeAgeHoursOff, in.MinVolumeAgeHours)
	ob.SetInt32(vacuumTaskConfigMinIntervalSecondsOff, in.MinIntervalSeconds)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ErasureCodingTaskConfig ---
//
// Fixed payload: fullness_ratio(8 double) quiet_for_seconds(4)
// min_volume_size_mb(4) collection_filter(8 ptr) preferred_tags(8 list ptr)
// replica_placement(8 ptr). preferred_tags is repeated string.

const (
	erasureCodingTaskConfigFullnessRatioOff    = 0
	erasureCodingTaskConfigQuietForSecondsOff  = 8
	erasureCodingTaskConfigMinVolumeSizeMbOff  = 12
	erasureCodingTaskConfigCollectionFilterOff = 16
	erasureCodingTaskConfigPreferredTagsOff    = 24
	erasureCodingTaskConfigReplicaPlacementOff = 32
	erasureCodingTaskConfigSize                = 40
)

// ErasureCodingTaskConfig is a zero-copy view into a ZAP-encoded ErasureCodingTaskConfig message.
type ErasureCodingTaskConfig struct{ o zap.Object }

// WrapErasureCodingTaskConfig parses b and returns a typed view.
func WrapErasureCodingTaskConfig(b []byte) (ErasureCodingTaskConfig, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ErasureCodingTaskConfig{}, err
	}
	return ErasureCodingTaskConfig{o: m.Root()}, nil
}

func (t ErasureCodingTaskConfig) FullnessRatio() float64 {
	return t.o.Float64(erasureCodingTaskConfigFullnessRatioOff)
}
func (t ErasureCodingTaskConfig) QuietForSeconds() int32 {
	return t.o.Int32(erasureCodingTaskConfigQuietForSecondsOff)
}
func (t ErasureCodingTaskConfig) MinVolumeSizeMb() int32 {
	return t.o.Int32(erasureCodingTaskConfigMinVolumeSizeMbOff)
}
func (t ErasureCodingTaskConfig) CollectionFilter() string {
	return t.o.Text(erasureCodingTaskConfigCollectionFilterOff)
}
func (t ErasureCodingTaskConfig) ReplicaPlacement() string {
	return t.o.Text(erasureCodingTaskConfigReplicaPlacementOff)
}

// PreferredTagsLen returns the number of preferred_tags elements.
func (t ErasureCodingTaskConfig) PreferredTagsLen() int {
	return t.o.List(erasureCodingTaskConfigPreferredTagsOff).Len()
}

// PreferredTag returns the i-th preferred_tags element.
func (t ErasureCodingTaskConfig) PreferredTag(i int) string {
	return elemText(t.o.List(erasureCodingTaskConfigPreferredTagsOff), i)
}

// ErasureCodingTaskConfigInput collects the field values for NewErasureCodingTaskConfig.
type ErasureCodingTaskConfigInput struct {
	FullnessRatio    float64
	QuietForSeconds  int32
	MinVolumeSizeMb  int32
	CollectionFilter string
	PreferredTags    []string
	ReplicaPlacement string
}

// NewErasureCodingTaskConfig builds a ZAP-encoded ErasureCodingTaskConfig message from in.
func NewErasureCodingTaskConfig(in ErasureCodingTaskConfigInput) []byte {
	b := zap.NewBuilder(256)
	tagsOff, tagsLen := 0, 0
	if len(in.PreferredTags) > 0 {
		el := startElemList(b)
		for _, s := range in.PreferredTags {
			el.addText(s)
		}
		tagsOff, tagsLen = el.finish()
	}
	ob := b.StartObject(erasureCodingTaskConfigSize)
	ob.SetFloat64(erasureCodingTaskConfigFullnessRatioOff, in.FullnessRatio)
	ob.SetInt32(erasureCodingTaskConfigQuietForSecondsOff, in.QuietForSeconds)
	ob.SetInt32(erasureCodingTaskConfigMinVolumeSizeMbOff, in.MinVolumeSizeMb)
	ob.SetText(erasureCodingTaskConfigCollectionFilterOff, in.CollectionFilter)
	ob.SetList(erasureCodingTaskConfigPreferredTagsOff, tagsOff, tagsLen)
	ob.SetText(erasureCodingTaskConfigReplicaPlacementOff, in.ReplicaPlacement)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- BalanceTaskConfig ---
//
// Fixed payload: imbalance_threshold(8 double) min_server_count(4).

const (
	balanceTaskConfigImbalanceThresholdOff = 0
	balanceTaskConfigMinServerCountOff     = 8
	balanceTaskConfigSize                  = 12
)

// BalanceTaskConfig is a zero-copy view into a ZAP-encoded BalanceTaskConfig message.
type BalanceTaskConfig struct{ o zap.Object }

// WrapBalanceTaskConfig parses b and returns a typed view.
func WrapBalanceTaskConfig(b []byte) (BalanceTaskConfig, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return BalanceTaskConfig{}, err
	}
	return BalanceTaskConfig{o: m.Root()}, nil
}

func (t BalanceTaskConfig) ImbalanceThreshold() float64 {
	return t.o.Float64(balanceTaskConfigImbalanceThresholdOff)
}
func (t BalanceTaskConfig) MinServerCount() int32 {
	return t.o.Int32(balanceTaskConfigMinServerCountOff)
}

// BalanceTaskConfigInput collects the field values for NewBalanceTaskConfig.
type BalanceTaskConfigInput struct {
	ImbalanceThreshold float64
	MinServerCount     int32
}

// NewBalanceTaskConfig builds a ZAP-encoded BalanceTaskConfig message from in.
func NewBalanceTaskConfig(in BalanceTaskConfigInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(balanceTaskConfigSize)
	ob.SetFloat64(balanceTaskConfigImbalanceThresholdOff, in.ImbalanceThreshold)
	ob.SetInt32(balanceTaskConfigMinServerCountOff, in.MinServerCount)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ReplicationTaskConfig ---
//
// Fixed payload: target_replica_count(4).

const (
	replicationTaskConfigTargetReplicaCountOff = 0
	replicationTaskConfigSize                  = 4
)

// ReplicationTaskConfig is a zero-copy view into a ZAP-encoded ReplicationTaskConfig message.
type ReplicationTaskConfig struct{ o zap.Object }

// WrapReplicationTaskConfig parses b and returns a typed view.
func WrapReplicationTaskConfig(b []byte) (ReplicationTaskConfig, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReplicationTaskConfig{}, err
	}
	return ReplicationTaskConfig{o: m.Root()}, nil
}

func (t ReplicationTaskConfig) TargetReplicaCount() int32 {
	return t.o.Int32(replicationTaskConfigTargetReplicaCountOff)
}

// ReplicationTaskConfigInput collects the field values for NewReplicationTaskConfig.
type ReplicationTaskConfigInput struct {
	TargetReplicaCount int32
}

// NewReplicationTaskConfig builds a ZAP-encoded ReplicationTaskConfig message from in.
func NewReplicationTaskConfig(in ReplicationTaskConfigInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(replicationTaskConfigSize)
	ob.SetInt32(replicationTaskConfigTargetReplicaCountOff, in.TargetReplicaCount)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- EcBalanceTaskConfig ---
//
// Fixed payload: imbalance_threshold(8 double) min_server_count(4)
// collection_filter(8 ptr) disk_type(8 ptr) preferred_tags(8 list ptr)
// replica_placement(8 ptr). preferred_tags is repeated string.

const (
	ecBalanceTaskConfigImbalanceThresholdOff = 0
	ecBalanceTaskConfigMinServerCountOff     = 8
	ecBalanceTaskConfigCollectionFilterOff   = 12
	ecBalanceTaskConfigDiskTypeOff           = 20
	ecBalanceTaskConfigPreferredTagsOff      = 28
	ecBalanceTaskConfigReplicaPlacementOff   = 36
	ecBalanceTaskConfigSize                  = 44
)

// EcBalanceTaskConfig is a zero-copy view into a ZAP-encoded EcBalanceTaskConfig message.
type EcBalanceTaskConfig struct{ o zap.Object }

// WrapEcBalanceTaskConfig parses b and returns a typed view.
func WrapEcBalanceTaskConfig(b []byte) (EcBalanceTaskConfig, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return EcBalanceTaskConfig{}, err
	}
	return EcBalanceTaskConfig{o: m.Root()}, nil
}

func (t EcBalanceTaskConfig) ImbalanceThreshold() float64 {
	return t.o.Float64(ecBalanceTaskConfigImbalanceThresholdOff)
}
func (t EcBalanceTaskConfig) MinServerCount() int32 {
	return t.o.Int32(ecBalanceTaskConfigMinServerCountOff)
}
func (t EcBalanceTaskConfig) CollectionFilter() string {
	return t.o.Text(ecBalanceTaskConfigCollectionFilterOff)
}
func (t EcBalanceTaskConfig) DiskType() string {
	return t.o.Text(ecBalanceTaskConfigDiskTypeOff)
}
func (t EcBalanceTaskConfig) ReplicaPlacement() string {
	return t.o.Text(ecBalanceTaskConfigReplicaPlacementOff)
}

// PreferredTagsLen returns the number of preferred_tags elements.
func (t EcBalanceTaskConfig) PreferredTagsLen() int {
	return t.o.List(ecBalanceTaskConfigPreferredTagsOff).Len()
}

// PreferredTag returns the i-th preferred_tags element.
func (t EcBalanceTaskConfig) PreferredTag(i int) string {
	return elemText(t.o.List(ecBalanceTaskConfigPreferredTagsOff), i)
}

// EcBalanceTaskConfigInput collects the field values for NewEcBalanceTaskConfig.
type EcBalanceTaskConfigInput struct {
	ImbalanceThreshold float64
	MinServerCount     int32
	CollectionFilter   string
	DiskType           string
	PreferredTags      []string
	ReplicaPlacement   string
}

// NewEcBalanceTaskConfig builds a ZAP-encoded EcBalanceTaskConfig message from in.
func NewEcBalanceTaskConfig(in EcBalanceTaskConfigInput) []byte {
	b := zap.NewBuilder(256)
	tagsOff, tagsLen := 0, 0
	if len(in.PreferredTags) > 0 {
		el := startElemList(b)
		for _, s := range in.PreferredTags {
			el.addText(s)
		}
		tagsOff, tagsLen = el.finish()
	}
	ob := b.StartObject(ecBalanceTaskConfigSize)
	ob.SetFloat64(ecBalanceTaskConfigImbalanceThresholdOff, in.ImbalanceThreshold)
	ob.SetInt32(ecBalanceTaskConfigMinServerCountOff, in.MinServerCount)
	ob.SetText(ecBalanceTaskConfigCollectionFilterOff, in.CollectionFilter)
	ob.SetText(ecBalanceTaskConfigDiskTypeOff, in.DiskType)
	ob.SetList(ecBalanceTaskConfigPreferredTagsOff, tagsOff, tagsLen)
	ob.SetText(ecBalanceTaskConfigReplicaPlacementOff, in.ReplicaPlacement)
	ob.FinishAsRoot()
	return b.Finish()
}
