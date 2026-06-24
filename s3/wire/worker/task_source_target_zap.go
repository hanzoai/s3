// Code generated from worker.proto; DO NOT EDIT.

package workerwire

import (
	zap "github.com/zap-proto/go"
)

// --- TaskSource ---
//
// Fixed payload: node(8 ptr) disk_id(4) rack(8 ptr) data_center(8 ptr)
// volume_id(4) shard_ids(8 list ptr) estimated_size(8).

const (
	taskSourceNodeOff          = 0
	taskSourceDiskIDOff        = 8
	taskSourceRackOff          = 12
	taskSourceDataCenterOff    = 20
	taskSourceVolumeIDOff      = 28
	taskSourceShardIDsOff      = 32
	taskSourceEstimatedSizeOff = 40
	taskSourceSize             = 48
)

// TaskSource is a zero-copy view into a ZAP-encoded TaskSource message.
type TaskSource struct{ o zap.Object }

// WrapTaskSource parses b and returns a typed view.
func WrapTaskSource(b []byte) (TaskSource, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TaskSource{}, err
	}
	return TaskSource{o: m.Root()}, nil
}

func (t TaskSource) Node() string          { return t.o.Text(taskSourceNodeOff) }
func (t TaskSource) DiskID() uint32        { return t.o.Uint32(taskSourceDiskIDOff) }
func (t TaskSource) Rack() string          { return t.o.Text(taskSourceRackOff) }
func (t TaskSource) DataCenter() string    { return t.o.Text(taskSourceDataCenterOff) }
func (t TaskSource) VolumeID() uint32      { return t.o.Uint32(taskSourceVolumeIDOff) }
func (t TaskSource) EstimatedSize() uint64 { return t.o.Uint64(taskSourceEstimatedSizeOff) }

// ShardIDsLen returns the number of shard_ids elements.
func (t TaskSource) ShardIDsLen() int { return t.o.List(taskSourceShardIDsOff).Len() }

// ShardID returns the i-th shard_ids element.
func (t TaskSource) ShardID(i int) uint32 { return t.o.List(taskSourceShardIDsOff).Uint32(i) }

// TaskSourceInput collects the field values for NewTaskSource.
type TaskSourceInput struct {
	Node          string
	DiskID        uint32
	Rack          string
	DataCenter    string
	VolumeID      uint32
	ShardIDs      []uint32
	EstimatedSize uint64
}

// NewTaskSource builds a ZAP-encoded TaskSource message from in.
func NewTaskSource(in TaskSourceInput) []byte {
	b := zap.NewBuilder(256)
	shardOff, shardLen := 0, 0
	if len(in.ShardIDs) > 0 {
		lb := b.StartList(4)
		for _, v := range in.ShardIDs {
			lb.AddUint32(v)
		}
		shardOff, shardLen = lb.Finish()
	}
	ob := b.StartObject(taskSourceSize)
	ob.SetText(taskSourceNodeOff, in.Node)
	ob.SetUint32(taskSourceDiskIDOff, in.DiskID)
	ob.SetText(taskSourceRackOff, in.Rack)
	ob.SetText(taskSourceDataCenterOff, in.DataCenter)
	ob.SetUint32(taskSourceVolumeIDOff, in.VolumeID)
	ob.SetList(taskSourceShardIDsOff, shardOff, shardLen)
	ob.SetUint64(taskSourceEstimatedSizeOff, in.EstimatedSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TaskTarget ---
//
// Fixed payload: node(8 ptr) disk_id(4) rack(8 ptr) data_center(8 ptr)
// volume_id(4) shard_ids(8 list ptr) estimated_size(8).

const (
	taskTargetNodeOff          = 0
	taskTargetDiskIDOff        = 8
	taskTargetRackOff          = 12
	taskTargetDataCenterOff    = 20
	taskTargetVolumeIDOff      = 28
	taskTargetShardIDsOff      = 32
	taskTargetEstimatedSizeOff = 40
	taskTargetSize             = 48
)

// TaskTarget is a zero-copy view into a ZAP-encoded TaskTarget message.
type TaskTarget struct{ o zap.Object }

// WrapTaskTarget parses b and returns a typed view.
func WrapTaskTarget(b []byte) (TaskTarget, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TaskTarget{}, err
	}
	return TaskTarget{o: m.Root()}, nil
}

func (t TaskTarget) Node() string          { return t.o.Text(taskTargetNodeOff) }
func (t TaskTarget) DiskID() uint32        { return t.o.Uint32(taskTargetDiskIDOff) }
func (t TaskTarget) Rack() string          { return t.o.Text(taskTargetRackOff) }
func (t TaskTarget) DataCenter() string    { return t.o.Text(taskTargetDataCenterOff) }
func (t TaskTarget) VolumeID() uint32      { return t.o.Uint32(taskTargetVolumeIDOff) }
func (t TaskTarget) EstimatedSize() uint64 { return t.o.Uint64(taskTargetEstimatedSizeOff) }

// ShardIDsLen returns the number of shard_ids elements.
func (t TaskTarget) ShardIDsLen() int { return t.o.List(taskTargetShardIDsOff).Len() }

// ShardID returns the i-th shard_ids element.
func (t TaskTarget) ShardID(i int) uint32 { return t.o.List(taskTargetShardIDsOff).Uint32(i) }

// TaskTargetInput collects the field values for NewTaskTarget.
type TaskTargetInput struct {
	Node          string
	DiskID        uint32
	Rack          string
	DataCenter    string
	VolumeID      uint32
	ShardIDs      []uint32
	EstimatedSize uint64
}

// NewTaskTarget builds a ZAP-encoded TaskTarget message from in.
func NewTaskTarget(in TaskTargetInput) []byte {
	b := zap.NewBuilder(256)
	shardOff, shardLen := 0, 0
	if len(in.ShardIDs) > 0 {
		lb := b.StartList(4)
		for _, v := range in.ShardIDs {
			lb.AddUint32(v)
		}
		shardOff, shardLen = lb.Finish()
	}
	ob := b.StartObject(taskTargetSize)
	ob.SetText(taskTargetNodeOff, in.Node)
	ob.SetUint32(taskTargetDiskIDOff, in.DiskID)
	ob.SetText(taskTargetRackOff, in.Rack)
	ob.SetText(taskTargetDataCenterOff, in.DataCenter)
	ob.SetUint32(taskTargetVolumeIDOff, in.VolumeID)
	ob.SetList(taskTargetShardIDsOff, shardOff, shardLen)
	ob.SetUint64(taskTargetEstimatedSizeOff, in.EstimatedSize)
	ob.FinishAsRoot()
	return b.Finish()
}
