// Code generated from worker.proto; DO NOT EDIT.

package workerwire

import (
	zap "github.com/zap-proto/go"
)

// --- VacuumTaskParams ---
//
// Fixed payload: garbage_threshold(8 double) force_vacuum(1 bool)
// batch_size(4) working_dir(8 ptr) verify_checksum(1 bool).

const (
	vacuumTaskParamsGarbageThresholdOff = 0
	vacuumTaskParamsForceVacuumOff      = 8
	vacuumTaskParamsBatchSizeOff        = 12
	vacuumTaskParamsWorkingDirOff       = 16
	vacuumTaskParamsVerifyChecksumOff   = 24
	vacuumTaskParamsSize                = 25
)

// VacuumTaskParams is a zero-copy view into a ZAP-encoded VacuumTaskParams message.
type VacuumTaskParams struct{ o zap.Object }

// WrapVacuumTaskParams parses b and returns a typed view.
func WrapVacuumTaskParams(b []byte) (VacuumTaskParams, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VacuumTaskParams{}, err
	}
	return VacuumTaskParams{o: m.Root()}, nil
}

func (t VacuumTaskParams) GarbageThreshold() float64 {
	return t.o.Float64(vacuumTaskParamsGarbageThresholdOff)
}
func (t VacuumTaskParams) ForceVacuum() bool    { return t.o.Bool(vacuumTaskParamsForceVacuumOff) }
func (t VacuumTaskParams) BatchSize() int32     { return t.o.Int32(vacuumTaskParamsBatchSizeOff) }
func (t VacuumTaskParams) WorkingDir() string   { return t.o.Text(vacuumTaskParamsWorkingDirOff) }
func (t VacuumTaskParams) VerifyChecksum() bool { return t.o.Bool(vacuumTaskParamsVerifyChecksumOff) }

// VacuumTaskParamsInput collects the field values for NewVacuumTaskParams.
type VacuumTaskParamsInput struct {
	GarbageThreshold float64
	ForceVacuum      bool
	BatchSize        int32
	WorkingDir       string
	VerifyChecksum   bool
}

// NewVacuumTaskParams builds a ZAP-encoded VacuumTaskParams message from in.
func NewVacuumTaskParams(in VacuumTaskParamsInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(vacuumTaskParamsSize)
	ob.SetFloat64(vacuumTaskParamsGarbageThresholdOff, in.GarbageThreshold)
	ob.SetBool(vacuumTaskParamsForceVacuumOff, in.ForceVacuum)
	ob.SetInt32(vacuumTaskParamsBatchSizeOff, in.BatchSize)
	ob.SetText(vacuumTaskParamsWorkingDirOff, in.WorkingDir)
	ob.SetBool(vacuumTaskParamsVerifyChecksumOff, in.VerifyChecksum)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ErasureCodingTaskParams ---
//
// Field 7 is reserved. Fixed payload: estimated_shard_size(8)
// data_shards(4) parity_shards(4) working_dir(8 ptr) master_client(8 ptr)
// cleanup_source(1) source_disk_type(8 ptr) encode_ts_ns(8).

const (
	erasureCodingTaskParamsEstimatedShardSizeOff = 0
	erasureCodingTaskParamsDataShardsOff         = 8
	erasureCodingTaskParamsParityShardsOff       = 12
	erasureCodingTaskParamsWorkingDirOff         = 16
	erasureCodingTaskParamsMasterClientOff       = 24
	erasureCodingTaskParamsCleanupSourceOff      = 32
	erasureCodingTaskParamsSourceDiskTypeOff     = 33
	erasureCodingTaskParamsEncodeTsNsOff         = 41
	erasureCodingTaskParamsSize                  = 49
)

// ErasureCodingTaskParams is a zero-copy view into a ZAP-encoded ErasureCodingTaskParams message.
type ErasureCodingTaskParams struct{ o zap.Object }

// WrapErasureCodingTaskParams parses b and returns a typed view.
func WrapErasureCodingTaskParams(b []byte) (ErasureCodingTaskParams, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ErasureCodingTaskParams{}, err
	}
	return ErasureCodingTaskParams{o: m.Root()}, nil
}

func (t ErasureCodingTaskParams) EstimatedShardSize() uint64 {
	return t.o.Uint64(erasureCodingTaskParamsEstimatedShardSizeOff)
}
func (t ErasureCodingTaskParams) DataShards() int32 {
	return t.o.Int32(erasureCodingTaskParamsDataShardsOff)
}
func (t ErasureCodingTaskParams) ParityShards() int32 {
	return t.o.Int32(erasureCodingTaskParamsParityShardsOff)
}
func (t ErasureCodingTaskParams) WorkingDir() string {
	return t.o.Text(erasureCodingTaskParamsWorkingDirOff)
}
func (t ErasureCodingTaskParams) MasterClient() string {
	return t.o.Text(erasureCodingTaskParamsMasterClientOff)
}
func (t ErasureCodingTaskParams) CleanupSource() bool {
	return t.o.Bool(erasureCodingTaskParamsCleanupSourceOff)
}
func (t ErasureCodingTaskParams) SourceDiskType() string {
	return t.o.Text(erasureCodingTaskParamsSourceDiskTypeOff)
}
func (t ErasureCodingTaskParams) EncodeTsNs() int64 {
	return t.o.Int64(erasureCodingTaskParamsEncodeTsNsOff)
}

// ErasureCodingTaskParamsInput collects the field values for NewErasureCodingTaskParams.
type ErasureCodingTaskParamsInput struct {
	EstimatedShardSize uint64
	DataShards         int32
	ParityShards       int32
	WorkingDir         string
	MasterClient       string
	CleanupSource      bool
	SourceDiskType     string
	EncodeTsNs         int64
}

// NewErasureCodingTaskParams builds a ZAP-encoded ErasureCodingTaskParams message from in.
func NewErasureCodingTaskParams(in ErasureCodingTaskParamsInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(erasureCodingTaskParamsSize)
	ob.SetUint64(erasureCodingTaskParamsEstimatedShardSizeOff, in.EstimatedShardSize)
	ob.SetInt32(erasureCodingTaskParamsDataShardsOff, in.DataShards)
	ob.SetInt32(erasureCodingTaskParamsParityShardsOff, in.ParityShards)
	ob.SetText(erasureCodingTaskParamsWorkingDirOff, in.WorkingDir)
	ob.SetText(erasureCodingTaskParamsMasterClientOff, in.MasterClient)
	ob.SetBool(erasureCodingTaskParamsCleanupSourceOff, in.CleanupSource)
	ob.SetText(erasureCodingTaskParamsSourceDiskTypeOff, in.SourceDiskType)
	ob.SetInt64(erasureCodingTaskParamsEncodeTsNsOff, in.EncodeTsNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- BalanceMoveSpec ---
//
// Fixed payload: volume_id(4) source_node(8 ptr) target_node(8 ptr)
// collection(8 ptr) volume_size(8).

const (
	balanceMoveSpecVolumeIDOff   = 0
	balanceMoveSpecSourceNodeOff = 8
	balanceMoveSpecTargetNodeOff = 16
	balanceMoveSpecCollectionOff = 24
	balanceMoveSpecVolumeSizeOff = 32
	balanceMoveSpecSize          = 40
)

// BalanceMoveSpec is a zero-copy view into a ZAP-encoded BalanceMoveSpec message.
type BalanceMoveSpec struct{ o zap.Object }

// WrapBalanceMoveSpec parses b and returns a typed view.
func WrapBalanceMoveSpec(b []byte) (BalanceMoveSpec, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return BalanceMoveSpec{}, err
	}
	return BalanceMoveSpec{o: m.Root()}, nil
}

func (t BalanceMoveSpec) VolumeID() uint32   { return t.o.Uint32(balanceMoveSpecVolumeIDOff) }
func (t BalanceMoveSpec) SourceNode() string { return t.o.Text(balanceMoveSpecSourceNodeOff) }
func (t BalanceMoveSpec) TargetNode() string { return t.o.Text(balanceMoveSpecTargetNodeOff) }
func (t BalanceMoveSpec) Collection() string { return t.o.Text(balanceMoveSpecCollectionOff) }
func (t BalanceMoveSpec) VolumeSize() uint64 { return t.o.Uint64(balanceMoveSpecVolumeSizeOff) }

// BalanceMoveSpecInput collects the field values for NewBalanceMoveSpec.
type BalanceMoveSpecInput struct {
	VolumeID   uint32
	SourceNode string
	TargetNode string
	Collection string
	VolumeSize uint64
}

// NewBalanceMoveSpec builds a ZAP-encoded BalanceMoveSpec message from in.
func NewBalanceMoveSpec(in BalanceMoveSpecInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(balanceMoveSpecSize)
	ob.SetUint32(balanceMoveSpecVolumeIDOff, in.VolumeID)
	ob.SetText(balanceMoveSpecSourceNodeOff, in.SourceNode)
	ob.SetText(balanceMoveSpecTargetNodeOff, in.TargetNode)
	ob.SetText(balanceMoveSpecCollectionOff, in.Collection)
	ob.SetUint64(balanceMoveSpecVolumeSizeOff, in.VolumeSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- BalanceTaskParams ---
//
// Fixed payload: force_move(1) timeout_seconds(4) max_concurrent_moves(4)
// moves(8 list ptr). moves is a repeated BalanceMoveSpec (out-of-line list).

const (
	balanceTaskParamsForceMoveOff          = 0
	balanceTaskParamsTimeoutSecondsOff     = 4
	balanceTaskParamsMaxConcurrentMovesOff = 8
	balanceTaskParamsMovesOff              = 16
	balanceTaskParamsSize                  = 24
)

// BalanceTaskParams is a zero-copy view into a ZAP-encoded BalanceTaskParams message.
type BalanceTaskParams struct{ o zap.Object }

// WrapBalanceTaskParams parses b and returns a typed view.
func WrapBalanceTaskParams(b []byte) (BalanceTaskParams, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return BalanceTaskParams{}, err
	}
	return BalanceTaskParams{o: m.Root()}, nil
}

func (t BalanceTaskParams) ForceMove() bool { return t.o.Bool(balanceTaskParamsForceMoveOff) }
func (t BalanceTaskParams) TimeoutSeconds() int32 {
	return t.o.Int32(balanceTaskParamsTimeoutSecondsOff)
}
func (t BalanceTaskParams) MaxConcurrentMoves() int32 {
	return t.o.Int32(balanceTaskParamsMaxConcurrentMovesOff)
}

// MovesLen returns the number of moves elements.
func (t BalanceTaskParams) MovesLen() int { return t.o.List(balanceTaskParamsMovesOff).Len() }

// Moves returns the i-th moves element as a typed BalanceMoveSpec view.
func (t BalanceTaskParams) Moves(i int) BalanceMoveSpec {
	return BalanceMoveSpec{o: t.o.List(balanceTaskParamsMovesOff).ObjectAt(i)}
}

// BalanceTaskParamsInput collects the field values for NewBalanceTaskParams.
// Moves elements are standalone BalanceMoveSpec buffers (NewBalanceMoveSpec).
type BalanceTaskParamsInput struct {
	ForceMove          bool
	TimeoutSeconds     int32
	MaxConcurrentMoves int32
	Moves              [][]byte
}

// NewBalanceTaskParams builds a ZAP-encoded BalanceTaskParams message from in.
func NewBalanceTaskParams(in BalanceTaskParamsInput) []byte {
	b := zap.NewBuilder(256)
	movesOff, movesLen := 0, 0
	if len(in.Moves) > 0 {
		el := startElemList(b)
		for _, m := range in.Moves {
			el.addElem(m)
		}
		movesOff, movesLen = el.finish()
	}
	ob := b.StartObject(balanceTaskParamsSize)
	ob.SetBool(balanceTaskParamsForceMoveOff, in.ForceMove)
	ob.SetInt32(balanceTaskParamsTimeoutSecondsOff, in.TimeoutSeconds)
	ob.SetInt32(balanceTaskParamsMaxConcurrentMovesOff, in.MaxConcurrentMoves)
	ob.SetList(balanceTaskParamsMovesOff, movesOff, movesLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ReplicationTaskParams ---
//
// Fixed payload: replica_count(4) verify_consistency(1).

const (
	replicationTaskParamsReplicaCountOff      = 0
	replicationTaskParamsVerifyConsistencyOff = 4
	replicationTaskParamsSize                 = 5
)

// ReplicationTaskParams is a zero-copy view into a ZAP-encoded ReplicationTaskParams message.
type ReplicationTaskParams struct{ o zap.Object }

// WrapReplicationTaskParams parses b and returns a typed view.
func WrapReplicationTaskParams(b []byte) (ReplicationTaskParams, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReplicationTaskParams{}, err
	}
	return ReplicationTaskParams{o: m.Root()}, nil
}

func (t ReplicationTaskParams) ReplicaCount() int32 {
	return t.o.Int32(replicationTaskParamsReplicaCountOff)
}
func (t ReplicationTaskParams) VerifyConsistency() bool {
	return t.o.Bool(replicationTaskParamsVerifyConsistencyOff)
}

// ReplicationTaskParamsInput collects the field values for NewReplicationTaskParams.
type ReplicationTaskParamsInput struct {
	ReplicaCount      int32
	VerifyConsistency bool
}

// NewReplicationTaskParams builds a ZAP-encoded ReplicationTaskParams message from in.
func NewReplicationTaskParams(in ReplicationTaskParamsInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(replicationTaskParamsSize)
	ob.SetInt32(replicationTaskParamsReplicaCountOff, in.ReplicaCount)
	ob.SetBool(replicationTaskParamsVerifyConsistencyOff, in.VerifyConsistency)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- EcShardMoveSpec ---
//
// Fixed payload: volume_id(4) shard_id(4) collection(8 ptr) source_node(8 ptr)
// source_disk_id(4) target_node(8 ptr) target_disk_id(4).

const (
	ecShardMoveSpecVolumeIDOff     = 0
	ecShardMoveSpecShardIDOff      = 4
	ecShardMoveSpecCollectionOff   = 8
	ecShardMoveSpecSourceNodeOff   = 16
	ecShardMoveSpecSourceDiskIDOff = 24
	ecShardMoveSpecTargetNodeOff   = 28
	ecShardMoveSpecTargetDiskIDOff = 36
	ecShardMoveSpecSize            = 40
)

// EcShardMoveSpec is a zero-copy view into a ZAP-encoded EcShardMoveSpec message.
type EcShardMoveSpec struct{ o zap.Object }

// WrapEcShardMoveSpec parses b and returns a typed view.
func WrapEcShardMoveSpec(b []byte) (EcShardMoveSpec, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return EcShardMoveSpec{}, err
	}
	return EcShardMoveSpec{o: m.Root()}, nil
}

func (t EcShardMoveSpec) VolumeID() uint32     { return t.o.Uint32(ecShardMoveSpecVolumeIDOff) }
func (t EcShardMoveSpec) ShardID() uint32      { return t.o.Uint32(ecShardMoveSpecShardIDOff) }
func (t EcShardMoveSpec) Collection() string   { return t.o.Text(ecShardMoveSpecCollectionOff) }
func (t EcShardMoveSpec) SourceNode() string   { return t.o.Text(ecShardMoveSpecSourceNodeOff) }
func (t EcShardMoveSpec) SourceDiskID() uint32 { return t.o.Uint32(ecShardMoveSpecSourceDiskIDOff) }
func (t EcShardMoveSpec) TargetNode() string   { return t.o.Text(ecShardMoveSpecTargetNodeOff) }
func (t EcShardMoveSpec) TargetDiskID() uint32 { return t.o.Uint32(ecShardMoveSpecTargetDiskIDOff) }

// EcShardMoveSpecInput collects the field values for NewEcShardMoveSpec.
type EcShardMoveSpecInput struct {
	VolumeID     uint32
	ShardID      uint32
	Collection   string
	SourceNode   string
	SourceDiskID uint32
	TargetNode   string
	TargetDiskID uint32
}

// NewEcShardMoveSpec builds a ZAP-encoded EcShardMoveSpec message from in.
func NewEcShardMoveSpec(in EcShardMoveSpecInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(ecShardMoveSpecSize)
	ob.SetUint32(ecShardMoveSpecVolumeIDOff, in.VolumeID)
	ob.SetUint32(ecShardMoveSpecShardIDOff, in.ShardID)
	ob.SetText(ecShardMoveSpecCollectionOff, in.Collection)
	ob.SetText(ecShardMoveSpecSourceNodeOff, in.SourceNode)
	ob.SetUint32(ecShardMoveSpecSourceDiskIDOff, in.SourceDiskID)
	ob.SetText(ecShardMoveSpecTargetNodeOff, in.TargetNode)
	ob.SetUint32(ecShardMoveSpecTargetDiskIDOff, in.TargetDiskID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- EcBalanceTaskParams ---
//
// Fixed payload: disk_type(8 ptr) max_parallelization(4) timeout_seconds(4)
// moves(8 list ptr). moves is a repeated EcShardMoveSpec (out-of-line list).

const (
	ecBalanceTaskParamsDiskTypeOff           = 0
	ecBalanceTaskParamsMaxParallelizationOff = 8
	ecBalanceTaskParamsTimeoutSecondsOff     = 12
	ecBalanceTaskParamsMovesOff              = 16
	ecBalanceTaskParamsSize                  = 24
)

// EcBalanceTaskParams is a zero-copy view into a ZAP-encoded EcBalanceTaskParams message.
type EcBalanceTaskParams struct{ o zap.Object }

// WrapEcBalanceTaskParams parses b and returns a typed view.
func WrapEcBalanceTaskParams(b []byte) (EcBalanceTaskParams, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return EcBalanceTaskParams{}, err
	}
	return EcBalanceTaskParams{o: m.Root()}, nil
}

func (t EcBalanceTaskParams) DiskType() string { return t.o.Text(ecBalanceTaskParamsDiskTypeOff) }
func (t EcBalanceTaskParams) MaxParallelization() int32 {
	return t.o.Int32(ecBalanceTaskParamsMaxParallelizationOff)
}
func (t EcBalanceTaskParams) TimeoutSeconds() int32 {
	return t.o.Int32(ecBalanceTaskParamsTimeoutSecondsOff)
}

// MovesLen returns the number of moves elements.
func (t EcBalanceTaskParams) MovesLen() int { return t.o.List(ecBalanceTaskParamsMovesOff).Len() }

// Moves returns the i-th moves element as a typed EcShardMoveSpec view.
func (t EcBalanceTaskParams) Moves(i int) EcShardMoveSpec {
	return EcShardMoveSpec{o: t.o.List(ecBalanceTaskParamsMovesOff).ObjectAt(i)}
}

// EcBalanceTaskParamsInput collects the field values for NewEcBalanceTaskParams.
// Moves elements are standalone EcShardMoveSpec buffers (NewEcShardMoveSpec).
type EcBalanceTaskParamsInput struct {
	DiskType           string
	MaxParallelization int32
	TimeoutSeconds     int32
	Moves              [][]byte
}

// NewEcBalanceTaskParams builds a ZAP-encoded EcBalanceTaskParams message from in.
func NewEcBalanceTaskParams(in EcBalanceTaskParamsInput) []byte {
	b := zap.NewBuilder(256)
	movesOff, movesLen := 0, 0
	if len(in.Moves) > 0 {
		el := startElemList(b)
		for _, m := range in.Moves {
			el.addElem(m)
		}
		movesOff, movesLen = el.finish()
	}
	ob := b.StartObject(ecBalanceTaskParamsSize)
	ob.SetText(ecBalanceTaskParamsDiskTypeOff, in.DiskType)
	ob.SetInt32(ecBalanceTaskParamsMaxParallelizationOff, in.MaxParallelization)
	ob.SetInt32(ecBalanceTaskParamsTimeoutSecondsOff, in.TimeoutSeconds)
	ob.SetList(ecBalanceTaskParamsMovesOff, movesOff, movesLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ContinuationHint ---
//
// Fixed payload: last_scanned_path(8 ptr) last_position_ns(8).

const (
	continuationHintLastScannedPathOff = 0
	continuationHintLastPositionNsOff  = 8
	continuationHintSize               = 16
)

// ContinuationHint is a zero-copy view into a ZAP-encoded ContinuationHint message.
type ContinuationHint struct{ o zap.Object }

// WrapContinuationHint parses b and returns a typed view.
func WrapContinuationHint(b []byte) (ContinuationHint, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ContinuationHint{}, err
	}
	return ContinuationHint{o: m.Root()}, nil
}

func (t ContinuationHint) LastScannedPath() string {
	return t.o.Text(continuationHintLastScannedPathOff)
}
func (t ContinuationHint) LastPositionNs() int64 {
	return t.o.Int64(continuationHintLastPositionNsOff)
}

// ContinuationHintInput collects the field values for NewContinuationHint.
type ContinuationHintInput struct {
	LastScannedPath string
	LastPositionNs  int64
}

// NewContinuationHint builds a ZAP-encoded ContinuationHint message from in.
func NewContinuationHint(in ContinuationHintInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(continuationHintSize)
	ob.SetText(continuationHintLastScannedPathOff, in.LastScannedPath)
	ob.SetInt64(continuationHintLastPositionNsOff, in.LastPositionNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- S3LifecycleParams ---
//
// Fixed payload: subtype(4 enum) bucket(8 ptr) rule_hash(8 bytes ptr)
// force(1) batch_time_budget_ns(8) batch_event_budget(4) continuation(4 obj ptr)
// shard_id(4). continuation is a singular ContinuationHint nested object.

const (
	s3LifecycleParamsSubtypeOff           = 0
	s3LifecycleParamsBucketOff            = 8
	s3LifecycleParamsRuleHashOff          = 16
	s3LifecycleParamsForceOff             = 24
	s3LifecycleParamsBatchTimeBudgetNsOff = 32
	s3LifecycleParamsBatchEventBudgetOff  = 40
	s3LifecycleParamsContinuationOff      = 44
	s3LifecycleParamsShardIDOff           = 48
	s3LifecycleParamsSize                 = 52
)

// S3LifecycleParams is a zero-copy view into a ZAP-encoded S3LifecycleParams message.
type S3LifecycleParams struct{ o zap.Object }

// WrapS3LifecycleParams parses b and returns a typed view.
func WrapS3LifecycleParams(b []byte) (S3LifecycleParams, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return S3LifecycleParams{}, err
	}
	return S3LifecycleParams{o: m.Root()}, nil
}

// Subtype reads the subtype enum (see S3LifecycleParamsSubtype* constants).
func (t S3LifecycleParams) Subtype() uint32  { return t.o.Uint32(s3LifecycleParamsSubtypeOff) }
func (t S3LifecycleParams) Bucket() string   { return t.o.Text(s3LifecycleParamsBucketOff) }
func (t S3LifecycleParams) RuleHash() []byte { return t.o.Bytes(s3LifecycleParamsRuleHashOff) }
func (t S3LifecycleParams) Force() bool      { return t.o.Bool(s3LifecycleParamsForceOff) }
func (t S3LifecycleParams) BatchTimeBudgetNs() int64 {
	return t.o.Int64(s3LifecycleParamsBatchTimeBudgetNsOff)
}
func (t S3LifecycleParams) BatchEventBudget() int32 {
	return t.o.Int32(s3LifecycleParamsBatchEventBudgetOff)
}
func (t S3LifecycleParams) ShardID() int32 { return t.o.Int32(s3LifecycleParamsShardIDOff) }

// Continuation returns the nested ContinuationHint (null-object if unset).
func (t S3LifecycleParams) Continuation() ContinuationHint {
	return ContinuationHint{o: t.o.Object(s3LifecycleParamsContinuationOff)}
}

// HasContinuation reports whether the continuation field is present.
func (t S3LifecycleParams) HasContinuation() bool {
	return !t.o.Object(s3LifecycleParamsContinuationOff).IsNull()
}

// S3LifecycleParamsInput collects the field values for NewS3LifecycleParams.
// Continuation, when non-nil, is a standalone ContinuationHint buffer.
type S3LifecycleParamsInput struct {
	Subtype           uint32
	Bucket            string
	RuleHash          []byte
	Force             bool
	BatchTimeBudgetNs int64
	BatchEventBudget  int32
	Continuation      []byte
	ShardID           int32
}

// NewS3LifecycleParams builds a ZAP-encoded S3LifecycleParams message from in.
func NewS3LifecycleParams(in S3LifecycleParamsInput) []byte {
	b := zap.NewBuilder(256)
	contOff := 0
	if len(in.Continuation) > 0 {
		contOff = embedMessage(b, in.Continuation)
	}
	ob := b.StartObject(s3LifecycleParamsSize)
	ob.SetUint32(s3LifecycleParamsSubtypeOff, in.Subtype)
	ob.SetText(s3LifecycleParamsBucketOff, in.Bucket)
	ob.SetBytes(s3LifecycleParamsRuleHashOff, in.RuleHash)
	ob.SetBool(s3LifecycleParamsForceOff, in.Force)
	ob.SetInt64(s3LifecycleParamsBatchTimeBudgetNsOff, in.BatchTimeBudgetNs)
	ob.SetInt32(s3LifecycleParamsBatchEventBudgetOff, in.BatchEventBudget)
	ob.SetObject(s3LifecycleParamsContinuationOff, contOff)
	ob.SetInt32(s3LifecycleParamsShardIDOff, in.ShardID)
	ob.FinishAsRoot()
	return b.Finish()
}
