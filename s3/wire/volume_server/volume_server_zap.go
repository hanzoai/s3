// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

import (
	"fmt"

	"github.com/zap-proto/go/rpc"
)

// Method ordinals for the VolumeServer service (stable 1-based wire ids,
// in volume_server.proto declaration order).
const (
	VolumeServerBatchDeleteOrdinal                 uint32 = 1
	VolumeServerVacuumVolumeCheckOrdinal           uint32 = 2
	VolumeServerVacuumVolumeCompactOrdinal         uint32 = 3
	VolumeServerVacuumVolumeCommitOrdinal          uint32 = 4
	VolumeServerVacuumVolumeCleanupOrdinal         uint32 = 5
	VolumeServerDeleteCollectionOrdinal            uint32 = 6
	VolumeServerAllocateVolumeOrdinal              uint32 = 7
	VolumeServerVolumeSyncStatusOrdinal            uint32 = 8
	VolumeServerVolumeIncrementalCopyOrdinal       uint32 = 9
	VolumeServerVolumeMountOrdinal                 uint32 = 10
	VolumeServerVolumeUnmountOrdinal               uint32 = 11
	VolumeServerVolumeDeleteOrdinal                uint32 = 12
	VolumeServerVolumeMarkReadonlyOrdinal          uint32 = 13
	VolumeServerVolumeMarkWritableOrdinal          uint32 = 14
	VolumeServerVolumeConfigureOrdinal             uint32 = 15
	VolumeServerVolumeStatusOrdinal                uint32 = 16
	VolumeServerGetStateOrdinal                    uint32 = 17
	VolumeServerSetStateOrdinal                    uint32 = 18
	VolumeServerVolumeCopyOrdinal                  uint32 = 19
	VolumeServerReadVolumeFileStatusOrdinal        uint32 = 20
	VolumeServerCopyFileOrdinal                    uint32 = 21
	VolumeServerReceiveFileOrdinal                 uint32 = 22
	VolumeServerReadNeedleBlobOrdinal              uint32 = 23
	VolumeServerReadNeedleMetaOrdinal              uint32 = 24
	VolumeServerWriteNeedleBlobOrdinal             uint32 = 25
	VolumeServerReadAllNeedlesOrdinal              uint32 = 26
	VolumeServerVolumeTailSenderOrdinal            uint32 = 27
	VolumeServerVolumeTailReceiverOrdinal          uint32 = 28
	VolumeServerVolumeEcShardsGenerateOrdinal      uint32 = 29
	VolumeServerVolumeEcShardsRebuildOrdinal       uint32 = 30
	VolumeServerVolumeEcShardsCopyOrdinal          uint32 = 31
	VolumeServerVolumeEcShardsDeleteOrdinal        uint32 = 32
	VolumeServerVolumeEcShardsMountOrdinal         uint32 = 33
	VolumeServerVolumeEcShardsUnmountOrdinal       uint32 = 34
	VolumeServerVolumeEcShardReadOrdinal           uint32 = 35
	VolumeServerVolumeEcBlobDeleteOrdinal          uint32 = 36
	VolumeServerVolumeEcShardsToVolumeOrdinal      uint32 = 37
	VolumeServerVolumeEcShardsInfoOrdinal          uint32 = 38
	VolumeServerVolumeTierMoveDatToRemoteOrdinal   uint32 = 39
	VolumeServerVolumeTierMoveDatFromRemoteOrdinal uint32 = 40
	VolumeServerVolumeServerStatusOrdinal          uint32 = 41
	VolumeServerVolumeServerLeaveOrdinal           uint32 = 42
	VolumeServerFetchAndWriteNeedleOrdinal         uint32 = 43
	VolumeServerScrubVolumeOrdinal                 uint32 = 44
	VolumeServerScrubEcVolumeOrdinal               uint32 = 45
	VolumeServerQueryOrdinal                       uint32 = 46
	VolumeServerVolumeNeedleStatusOrdinal          uint32 = 47
	VolumeServerPingOrdinal                        uint32 = 48
)

// VolumeServerChannel ships one Call envelope and awaits its correlated Response.
type VolumeServerChannel interface {
	Call(envelope []byte) (rpc.Response, error)
}

// VolumeServerClient is a typed RPC client for the VolumeServer service over a
// ZAP call channel. Each call takes a fresh PromiseID from sess; the pipelined
// "On" form of a method sets Target to a prior call's Promise so the server
// chains them.
type VolumeServerClient struct {
	ch   VolumeServerChannel
	cap  []byte
	sess *rpc.Session
}

// NewVolumeServerClient returns a client that issues calls over ch, attaching
// cap (which may be nil) to every request.
func NewVolumeServerClient(ch VolumeServerChannel, capability []byte) *VolumeServerClient {
	return &VolumeServerClient{ch: ch, cap: capability, sess: rpc.NewSession()}
}

func (c *VolumeServerClient) BatchDelete(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeBatchDelete(rpc.NoTarget, req)
}

// BatchDeleteOn issues BatchDelete as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) BatchDeleteOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeBatchDelete(on.ID, nil)
}

func (c *VolumeServerClient) invokeBatchDelete(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerBatchDeleteOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.BatchDelete: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VacuumVolumeCheck(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVacuumVolumeCheck(rpc.NoTarget, req)
}

// VacuumVolumeCheckOn issues VacuumVolumeCheck as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VacuumVolumeCheckOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVacuumVolumeCheck(on.ID, nil)
}

func (c *VolumeServerClient) invokeVacuumVolumeCheck(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVacuumVolumeCheckOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VacuumVolumeCheck: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: VacuumVolumeCompact is a server-streaming RPC (VolumeServer.VacuumVolumeCompact); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) VacuumVolumeCompact(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVacuumVolumeCompact(rpc.NoTarget, req)
}

// VacuumVolumeCompactOn issues VacuumVolumeCompact as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VacuumVolumeCompactOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVacuumVolumeCompact(on.ID, nil)
}

func (c *VolumeServerClient) invokeVacuumVolumeCompact(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVacuumVolumeCompactOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VacuumVolumeCompact: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VacuumVolumeCommit(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVacuumVolumeCommit(rpc.NoTarget, req)
}

// VacuumVolumeCommitOn issues VacuumVolumeCommit as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VacuumVolumeCommitOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVacuumVolumeCommit(on.ID, nil)
}

func (c *VolumeServerClient) invokeVacuumVolumeCommit(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVacuumVolumeCommitOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VacuumVolumeCommit: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VacuumVolumeCleanup(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVacuumVolumeCleanup(rpc.NoTarget, req)
}

// VacuumVolumeCleanupOn issues VacuumVolumeCleanup as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VacuumVolumeCleanupOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVacuumVolumeCleanup(on.ID, nil)
}

func (c *VolumeServerClient) invokeVacuumVolumeCleanup(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVacuumVolumeCleanupOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VacuumVolumeCleanup: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) DeleteCollection(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeDeleteCollection(rpc.NoTarget, req)
}

// DeleteCollectionOn issues DeleteCollection as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) DeleteCollectionOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeDeleteCollection(on.ID, nil)
}

func (c *VolumeServerClient) invokeDeleteCollection(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerDeleteCollectionOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.DeleteCollection: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) AllocateVolume(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeAllocateVolume(rpc.NoTarget, req)
}

// AllocateVolumeOn issues AllocateVolume as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) AllocateVolumeOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeAllocateVolume(on.ID, nil)
}

func (c *VolumeServerClient) invokeAllocateVolume(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerAllocateVolumeOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.AllocateVolume: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeSyncStatus(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeSyncStatus(rpc.NoTarget, req)
}

// VolumeSyncStatusOn issues VolumeSyncStatus as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeSyncStatusOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeSyncStatus(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeSyncStatus(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeSyncStatusOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeSyncStatus: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: VolumeIncrementalCopy is a server-streaming RPC (VolumeServer.VolumeIncrementalCopy); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) VolumeIncrementalCopy(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeIncrementalCopy(rpc.NoTarget, req)
}

// VolumeIncrementalCopyOn issues VolumeIncrementalCopy as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeIncrementalCopyOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeIncrementalCopy(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeIncrementalCopy(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeIncrementalCopyOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeIncrementalCopy: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeMount(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeMount(rpc.NoTarget, req)
}

// VolumeMountOn issues VolumeMount as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeMountOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeMount(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeMount(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeMountOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeMount: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeUnmount(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeUnmount(rpc.NoTarget, req)
}

// VolumeUnmountOn issues VolumeUnmount as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeUnmountOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeUnmount(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeUnmount(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeUnmountOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeUnmount: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeDelete(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeDelete(rpc.NoTarget, req)
}

// VolumeDeleteOn issues VolumeDelete as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeDeleteOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeDelete(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeDelete(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeDeleteOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeDelete: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeMarkReadonly(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeMarkReadonly(rpc.NoTarget, req)
}

// VolumeMarkReadonlyOn issues VolumeMarkReadonly as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeMarkReadonlyOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeMarkReadonly(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeMarkReadonly(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeMarkReadonlyOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeMarkReadonly: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeMarkWritable(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeMarkWritable(rpc.NoTarget, req)
}

// VolumeMarkWritableOn issues VolumeMarkWritable as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeMarkWritableOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeMarkWritable(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeMarkWritable(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeMarkWritableOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeMarkWritable: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeConfigure(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeConfigure(rpc.NoTarget, req)
}

// VolumeConfigureOn issues VolumeConfigure as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeConfigureOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeConfigure(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeConfigure(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeConfigureOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeConfigure: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeStatus(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeStatus(rpc.NoTarget, req)
}

// VolumeStatusOn issues VolumeStatus as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeStatusOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeStatus(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeStatus(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeStatusOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeStatus: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) GetState(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeGetState(rpc.NoTarget, req)
}

// GetStateOn issues GetState as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) GetStateOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeGetState(on.ID, nil)
}

func (c *VolumeServerClient) invokeGetState(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerGetStateOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.GetState: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) SetState(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeSetState(rpc.NoTarget, req)
}

// SetStateOn issues SetState as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) SetStateOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeSetState(on.ID, nil)
}

func (c *VolumeServerClient) invokeSetState(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerSetStateOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.SetState: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: VolumeCopy is a server-streaming RPC (VolumeServer.VolumeCopy); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) VolumeCopy(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeCopy(rpc.NoTarget, req)
}

// VolumeCopyOn issues VolumeCopy as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeCopyOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeCopy(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeCopy(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeCopyOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeCopy: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) ReadVolumeFileStatus(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeReadVolumeFileStatus(rpc.NoTarget, req)
}

// ReadVolumeFileStatusOn issues ReadVolumeFileStatus as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) ReadVolumeFileStatusOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeReadVolumeFileStatus(on.ID, nil)
}

func (c *VolumeServerClient) invokeReadVolumeFileStatus(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerReadVolumeFileStatusOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.ReadVolumeFileStatus: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: CopyFile is a server-streaming RPC (VolumeServer.CopyFile); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) CopyFile(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeCopyFile(rpc.NoTarget, req)
}

// CopyFileOn issues CopyFile as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) CopyFileOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeCopyFile(on.ID, nil)
}

func (c *VolumeServerClient) invokeCopyFile(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerCopyFileOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.CopyFile: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: ReceiveFile is a client-streaming RPC (VolumeServer.ReceiveFile); its body lands
// when the transport streaming primitive ships. The request element schema is
// emitted (ReceiveFileRequest) and the unary response is returned here.
func (c *VolumeServerClient) ReceiveFile(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeReceiveFile(rpc.NoTarget, req)
}

// ReceiveFileOn issues ReceiveFile as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) ReceiveFileOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeReceiveFile(on.ID, nil)
}

func (c *VolumeServerClient) invokeReceiveFile(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerReceiveFileOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.ReceiveFile: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) ReadNeedleBlob(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeReadNeedleBlob(rpc.NoTarget, req)
}

// ReadNeedleBlobOn issues ReadNeedleBlob as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) ReadNeedleBlobOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeReadNeedleBlob(on.ID, nil)
}

func (c *VolumeServerClient) invokeReadNeedleBlob(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerReadNeedleBlobOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.ReadNeedleBlob: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) ReadNeedleMeta(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeReadNeedleMeta(rpc.NoTarget, req)
}

// ReadNeedleMetaOn issues ReadNeedleMeta as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) ReadNeedleMetaOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeReadNeedleMeta(on.ID, nil)
}

func (c *VolumeServerClient) invokeReadNeedleMeta(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerReadNeedleMetaOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.ReadNeedleMeta: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) WriteNeedleBlob(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeWriteNeedleBlob(rpc.NoTarget, req)
}

// WriteNeedleBlobOn issues WriteNeedleBlob as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) WriteNeedleBlobOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeWriteNeedleBlob(on.ID, nil)
}

func (c *VolumeServerClient) invokeWriteNeedleBlob(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerWriteNeedleBlobOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.WriteNeedleBlob: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: ReadAllNeedles is a server-streaming RPC (VolumeServer.ReadAllNeedles); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) ReadAllNeedles(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeReadAllNeedles(rpc.NoTarget, req)
}

// ReadAllNeedlesOn issues ReadAllNeedles as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) ReadAllNeedlesOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeReadAllNeedles(on.ID, nil)
}

func (c *VolumeServerClient) invokeReadAllNeedles(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerReadAllNeedlesOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.ReadAllNeedles: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: VolumeTailSender is a server-streaming RPC (VolumeServer.VolumeTailSender); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) VolumeTailSender(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeTailSender(rpc.NoTarget, req)
}

// VolumeTailSenderOn issues VolumeTailSender as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeTailSenderOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeTailSender(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeTailSender(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeTailSenderOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeTailSender: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeTailReceiver(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeTailReceiver(rpc.NoTarget, req)
}

// VolumeTailReceiverOn issues VolumeTailReceiver as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeTailReceiverOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeTailReceiver(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeTailReceiver(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeTailReceiverOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeTailReceiver: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcShardsGenerate(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsGenerate(rpc.NoTarget, req)
}

// VolumeEcShardsGenerateOn issues VolumeEcShardsGenerate as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardsGenerateOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsGenerate(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardsGenerate(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeEcShardsGenerateOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsGenerate: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcShardsRebuild(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsRebuild(rpc.NoTarget, req)
}

// VolumeEcShardsRebuildOn issues VolumeEcShardsRebuild as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardsRebuildOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsRebuild(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardsRebuild(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeEcShardsRebuildOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsRebuild: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcShardsCopy(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsCopy(rpc.NoTarget, req)
}

// VolumeEcShardsCopyOn issues VolumeEcShardsCopy as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardsCopyOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsCopy(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardsCopy(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeEcShardsCopyOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsCopy: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcShardsDelete(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsDelete(rpc.NoTarget, req)
}

// VolumeEcShardsDeleteOn issues VolumeEcShardsDelete as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardsDeleteOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsDelete(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardsDelete(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeEcShardsDeleteOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsDelete: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcShardsMount(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsMount(rpc.NoTarget, req)
}

// VolumeEcShardsMountOn issues VolumeEcShardsMount as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardsMountOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsMount(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardsMount(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeEcShardsMountOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsMount: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcShardsUnmount(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsUnmount(rpc.NoTarget, req)
}

// VolumeEcShardsUnmountOn issues VolumeEcShardsUnmount as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardsUnmountOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsUnmount(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardsUnmount(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeEcShardsUnmountOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsUnmount: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: VolumeEcShardRead is a server-streaming RPC (VolumeServer.VolumeEcShardRead); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) VolumeEcShardRead(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardRead(rpc.NoTarget, req)
}

// VolumeEcShardReadOn issues VolumeEcShardRead as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardReadOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardRead(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardRead(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeEcShardReadOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardRead: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcBlobDelete(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcBlobDelete(rpc.NoTarget, req)
}

// VolumeEcBlobDeleteOn issues VolumeEcBlobDelete as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcBlobDeleteOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcBlobDelete(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcBlobDelete(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeEcBlobDeleteOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcBlobDelete: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcShardsToVolume(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsToVolume(rpc.NoTarget, req)
}

// VolumeEcShardsToVolumeOn issues VolumeEcShardsToVolume as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardsToVolumeOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsToVolume(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardsToVolume(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeEcShardsToVolumeOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsToVolume: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcShardsInfo(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsInfo(rpc.NoTarget, req)
}

// VolumeEcShardsInfoOn issues VolumeEcShardsInfo as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardsInfoOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsInfo(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardsInfo(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeEcShardsInfoOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsInfo: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: VolumeTierMoveDatToRemote is a server-streaming RPC (VolumeServer.VolumeTierMoveDatToRemote); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) VolumeTierMoveDatToRemote(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeTierMoveDatToRemote(rpc.NoTarget, req)
}

// VolumeTierMoveDatToRemoteOn issues VolumeTierMoveDatToRemote as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeTierMoveDatToRemoteOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeTierMoveDatToRemote(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeTierMoveDatToRemote(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeTierMoveDatToRemoteOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeTierMoveDatToRemote: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: VolumeTierMoveDatFromRemote is a server-streaming RPC (VolumeServer.VolumeTierMoveDatFromRemote); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) VolumeTierMoveDatFromRemote(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeTierMoveDatFromRemote(rpc.NoTarget, req)
}

// VolumeTierMoveDatFromRemoteOn issues VolumeTierMoveDatFromRemote as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeTierMoveDatFromRemoteOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeTierMoveDatFromRemote(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeTierMoveDatFromRemote(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeTierMoveDatFromRemoteOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeTierMoveDatFromRemote: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeServerStatus(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeServerStatus(rpc.NoTarget, req)
}

// VolumeServerStatusOn issues VolumeServerStatus as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeServerStatusOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeServerStatus(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeServerStatus(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeServerStatusOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeServerStatus: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeServerLeave(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeServerLeave(rpc.NoTarget, req)
}

// VolumeServerLeaveOn issues VolumeServerLeave as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeServerLeaveOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeServerLeave(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeServerLeave(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeServerLeaveOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeServerLeave: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) FetchAndWriteNeedle(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeFetchAndWriteNeedle(rpc.NoTarget, req)
}

// FetchAndWriteNeedleOn issues FetchAndWriteNeedle as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) FetchAndWriteNeedleOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeFetchAndWriteNeedle(on.ID, nil)
}

func (c *VolumeServerClient) invokeFetchAndWriteNeedle(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerFetchAndWriteNeedleOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.FetchAndWriteNeedle: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) ScrubVolume(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeScrubVolume(rpc.NoTarget, req)
}

// ScrubVolumeOn issues ScrubVolume as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) ScrubVolumeOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeScrubVolume(on.ID, nil)
}

func (c *VolumeServerClient) invokeScrubVolume(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerScrubVolumeOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.ScrubVolume: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) ScrubEcVolume(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeScrubEcVolume(rpc.NoTarget, req)
}

// ScrubEcVolumeOn issues ScrubEcVolume as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) ScrubEcVolumeOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeScrubEcVolume(on.ID, nil)
}

func (c *VolumeServerClient) invokeScrubEcVolume(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerScrubEcVolumeOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.ScrubEcVolume: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: Query is a server-streaming RPC (VolumeServer.Query); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) Query(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeQuery(rpc.NoTarget, req)
}

// QueryOn issues Query as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) QueryOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeQuery(on.ID, nil)
}

func (c *VolumeServerClient) invokeQuery(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerQueryOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.Query: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeNeedleStatus(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeNeedleStatus(rpc.NoTarget, req)
}

// VolumeNeedleStatusOn issues VolumeNeedleStatus as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeNeedleStatusOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeNeedleStatus(on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeNeedleStatus(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerVolumeNeedleStatusOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.VolumeNeedleStatus: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) Ping(req []byte) (rpc.Promise, []byte, error) {
	return c.invokePing(rpc.NoTarget, req)
}

// PingOn issues Ping as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) PingOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokePing(on.ID, nil)
}

func (c *VolumeServerClient) invokePing(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    VolumeServerPingOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("VolumeServer.Ping: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// VolumeServerHandler is the server contract for the VolumeServer service.
// Implement each method, then route requests to it with DispatchVolumeServer.
// Streaming methods (marked // STREAMING) take the same unary shape here until
// the transport streaming primitive ships.
type VolumeServerHandler interface {
	BatchDelete(req []byte) ([]byte, error)
	VacuumVolumeCheck(req []byte) ([]byte, error)
	VacuumVolumeCompact(req []byte) ([]byte, error) // STREAMING
	VacuumVolumeCommit(req []byte) ([]byte, error)
	VacuumVolumeCleanup(req []byte) ([]byte, error)
	DeleteCollection(req []byte) ([]byte, error)
	AllocateVolume(req []byte) ([]byte, error)
	VolumeSyncStatus(req []byte) ([]byte, error)
	VolumeIncrementalCopy(req []byte) ([]byte, error) // STREAMING
	VolumeMount(req []byte) ([]byte, error)
	VolumeUnmount(req []byte) ([]byte, error)
	VolumeDelete(req []byte) ([]byte, error)
	VolumeMarkReadonly(req []byte) ([]byte, error)
	VolumeMarkWritable(req []byte) ([]byte, error)
	VolumeConfigure(req []byte) ([]byte, error)
	VolumeStatus(req []byte) ([]byte, error)
	GetState(req []byte) ([]byte, error)
	SetState(req []byte) ([]byte, error)
	VolumeCopy(req []byte) ([]byte, error) // STREAMING
	ReadVolumeFileStatus(req []byte) ([]byte, error)
	CopyFile(req []byte) ([]byte, error)    // STREAMING
	ReceiveFile(req []byte) ([]byte, error) // STREAMING
	ReadNeedleBlob(req []byte) ([]byte, error)
	ReadNeedleMeta(req []byte) ([]byte, error)
	WriteNeedleBlob(req []byte) ([]byte, error)
	ReadAllNeedles(req []byte) ([]byte, error)   // STREAMING
	VolumeTailSender(req []byte) ([]byte, error) // STREAMING
	VolumeTailReceiver(req []byte) ([]byte, error)
	VolumeEcShardsGenerate(req []byte) ([]byte, error)
	VolumeEcShardsRebuild(req []byte) ([]byte, error)
	VolumeEcShardsCopy(req []byte) ([]byte, error)
	VolumeEcShardsDelete(req []byte) ([]byte, error)
	VolumeEcShardsMount(req []byte) ([]byte, error)
	VolumeEcShardsUnmount(req []byte) ([]byte, error)
	VolumeEcShardRead(req []byte) ([]byte, error) // STREAMING
	VolumeEcBlobDelete(req []byte) ([]byte, error)
	VolumeEcShardsToVolume(req []byte) ([]byte, error)
	VolumeEcShardsInfo(req []byte) ([]byte, error)
	VolumeTierMoveDatToRemote(req []byte) ([]byte, error)   // STREAMING
	VolumeTierMoveDatFromRemote(req []byte) ([]byte, error) // STREAMING
	VolumeServerStatus(req []byte) ([]byte, error)
	VolumeServerLeave(req []byte) ([]byte, error)
	FetchAndWriteNeedle(req []byte) ([]byte, error)
	ScrubVolume(req []byte) ([]byte, error)
	ScrubEcVolume(req []byte) ([]byte, error)
	Query(req []byte) ([]byte, error) // STREAMING
	VolumeNeedleStatus(req []byte) ([]byte, error)
	Ping(req []byte) ([]byte, error)
}

// DispatchVolumeServer decodes a Call envelope, routes it by method ordinal to
// h, and returns the response envelope. An unknown ordinal yields a
// StatusNotFound response; a handler error yields StatusInternal.
func DispatchVolumeServer(h VolumeServerHandler, envelope []byte) ([]byte, error) {
	call, err := rpc.ParseRequest(envelope)
	if err != nil {
		return nil, err
	}
	switch call.Method {
	case VolumeServerBatchDeleteOrdinal:
		body, err := h.BatchDelete(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVacuumVolumeCheckOrdinal:
		body, err := h.VacuumVolumeCheck(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVacuumVolumeCompactOrdinal:
		body, err := h.VacuumVolumeCompact(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVacuumVolumeCommitOrdinal:
		body, err := h.VacuumVolumeCommit(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVacuumVolumeCleanupOrdinal:
		body, err := h.VacuumVolumeCleanup(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerDeleteCollectionOrdinal:
		body, err := h.DeleteCollection(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerAllocateVolumeOrdinal:
		body, err := h.AllocateVolume(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeSyncStatusOrdinal:
		body, err := h.VolumeSyncStatus(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeIncrementalCopyOrdinal:
		body, err := h.VolumeIncrementalCopy(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeMountOrdinal:
		body, err := h.VolumeMount(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeUnmountOrdinal:
		body, err := h.VolumeUnmount(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeDeleteOrdinal:
		body, err := h.VolumeDelete(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeMarkReadonlyOrdinal:
		body, err := h.VolumeMarkReadonly(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeMarkWritableOrdinal:
		body, err := h.VolumeMarkWritable(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeConfigureOrdinal:
		body, err := h.VolumeConfigure(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeStatusOrdinal:
		body, err := h.VolumeStatus(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerGetStateOrdinal:
		body, err := h.GetState(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerSetStateOrdinal:
		body, err := h.SetState(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeCopyOrdinal:
		body, err := h.VolumeCopy(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerReadVolumeFileStatusOrdinal:
		body, err := h.ReadVolumeFileStatus(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerCopyFileOrdinal:
		body, err := h.CopyFile(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerReceiveFileOrdinal:
		body, err := h.ReceiveFile(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerReadNeedleBlobOrdinal:
		body, err := h.ReadNeedleBlob(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerReadNeedleMetaOrdinal:
		body, err := h.ReadNeedleMeta(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerWriteNeedleBlobOrdinal:
		body, err := h.WriteNeedleBlob(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerReadAllNeedlesOrdinal:
		body, err := h.ReadAllNeedles(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeTailSenderOrdinal:
		body, err := h.VolumeTailSender(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeTailReceiverOrdinal:
		body, err := h.VolumeTailReceiver(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardsGenerateOrdinal:
		body, err := h.VolumeEcShardsGenerate(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardsRebuildOrdinal:
		body, err := h.VolumeEcShardsRebuild(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardsCopyOrdinal:
		body, err := h.VolumeEcShardsCopy(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardsDeleteOrdinal:
		body, err := h.VolumeEcShardsDelete(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardsMountOrdinal:
		body, err := h.VolumeEcShardsMount(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardsUnmountOrdinal:
		body, err := h.VolumeEcShardsUnmount(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardReadOrdinal:
		body, err := h.VolumeEcShardRead(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcBlobDeleteOrdinal:
		body, err := h.VolumeEcBlobDelete(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardsToVolumeOrdinal:
		body, err := h.VolumeEcShardsToVolume(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardsInfoOrdinal:
		body, err := h.VolumeEcShardsInfo(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeTierMoveDatToRemoteOrdinal:
		body, err := h.VolumeTierMoveDatToRemote(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeTierMoveDatFromRemoteOrdinal:
		body, err := h.VolumeTierMoveDatFromRemote(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeServerStatusOrdinal:
		body, err := h.VolumeServerStatus(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeServerLeaveOrdinal:
		body, err := h.VolumeServerLeave(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerFetchAndWriteNeedleOrdinal:
		body, err := h.FetchAndWriteNeedle(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerScrubVolumeOrdinal:
		body, err := h.ScrubVolume(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerScrubEcVolumeOrdinal:
		body, err := h.ScrubEcVolume(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerQueryOrdinal:
		body, err := h.Query(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeNeedleStatusOrdinal:
		body, err := h.VolumeNeedleStatus(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerPingOrdinal:
		body, err := h.Ping(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	default:
		return rpc.BuildResponse(rpc.StatusNotFound, call.PromiseID, nil), nil
	}
}
