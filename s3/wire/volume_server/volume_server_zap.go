// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

import (
	"context"
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
	// CallContext is Call that also aborts when ctx is done (transport.Conn
	// satisfies both).
	CallContext(ctx context.Context, envelope []byte) (rpc.Response, error)
	// NextPromiseID allocates a call's PromiseID from the CONNECTION rather than
	// from a per-client counter. PromiseID is the connection's demultiplexing
	// key: the transport keys its in-flight table by it, and OpenStream already
	// allocates stream IDs the same way. A client that numbered its own calls
	// would restart at 1 on every construction, so two clients sharing a pooled
	// conn collide — the transport overwrites the first waiter's slot and its
	// response is dropped, blocking that caller for as long as its context runs.
	NextPromiseID() uint32
}

// VolumeServerClient is a typed RPC client for the VolumeServer service over a
// ZAP call channel. Each call takes a fresh PromiseID from sess; the pipelined
// "On" form of a method sets Target to a prior call's Promise so the server
// chains them.
type VolumeServerClient struct {
	ch   VolumeServerChannel
	cap  []byte
}

// NewVolumeServerClient returns a client that issues calls over ch, attaching
// cap (which may be nil) to every request.
func NewVolumeServerClient(ch VolumeServerChannel, capability []byte) *VolumeServerClient {
	return &VolumeServerClient{ch: ch, cap: capability}
}

func (c *VolumeServerClient) BatchDelete(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeBatchDelete(ctx, rpc.NoTarget, req)
}

// BatchDeleteOn issues BatchDelete as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) BatchDeleteOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeBatchDelete(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeBatchDelete(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.BatchDelete: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.BatchDelete: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VacuumVolumeCheck(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVacuumVolumeCheck(ctx, rpc.NoTarget, req)
}

// VacuumVolumeCheckOn issues VacuumVolumeCheck as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VacuumVolumeCheckOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVacuumVolumeCheck(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVacuumVolumeCheck(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VacuumVolumeCheck: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VacuumVolumeCheck: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: VacuumVolumeCompact is a server-streaming RPC (VolumeServer.VacuumVolumeCompact); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) VacuumVolumeCompact(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVacuumVolumeCompact(ctx, rpc.NoTarget, req)
}

// VacuumVolumeCompactOn issues VacuumVolumeCompact as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VacuumVolumeCompactOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVacuumVolumeCompact(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVacuumVolumeCompact(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VacuumVolumeCompact: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VacuumVolumeCompact: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VacuumVolumeCommit(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVacuumVolumeCommit(ctx, rpc.NoTarget, req)
}

// VacuumVolumeCommitOn issues VacuumVolumeCommit as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VacuumVolumeCommitOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVacuumVolumeCommit(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVacuumVolumeCommit(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VacuumVolumeCommit: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VacuumVolumeCommit: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VacuumVolumeCleanup(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVacuumVolumeCleanup(ctx, rpc.NoTarget, req)
}

// VacuumVolumeCleanupOn issues VacuumVolumeCleanup as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VacuumVolumeCleanupOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVacuumVolumeCleanup(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVacuumVolumeCleanup(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VacuumVolumeCleanup: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VacuumVolumeCleanup: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) DeleteCollection(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeDeleteCollection(ctx, rpc.NoTarget, req)
}

// DeleteCollectionOn issues DeleteCollection as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) DeleteCollectionOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeDeleteCollection(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeDeleteCollection(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.DeleteCollection: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.DeleteCollection: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) AllocateVolume(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeAllocateVolume(ctx, rpc.NoTarget, req)
}

// AllocateVolumeOn issues AllocateVolume as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) AllocateVolumeOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeAllocateVolume(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeAllocateVolume(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.AllocateVolume: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.AllocateVolume: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeSyncStatus(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeSyncStatus(ctx, rpc.NoTarget, req)
}

// VolumeSyncStatusOn issues VolumeSyncStatus as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeSyncStatusOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeSyncStatus(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeSyncStatus(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeSyncStatus: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeSyncStatus: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: VolumeIncrementalCopy is a server-streaming RPC (VolumeServer.VolumeIncrementalCopy); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) VolumeIncrementalCopy(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeIncrementalCopy(ctx, rpc.NoTarget, req)
}

// VolumeIncrementalCopyOn issues VolumeIncrementalCopy as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeIncrementalCopyOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeIncrementalCopy(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeIncrementalCopy(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeIncrementalCopy: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeIncrementalCopy: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeMount(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeMount(ctx, rpc.NoTarget, req)
}

// VolumeMountOn issues VolumeMount as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeMountOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeMount(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeMount(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeMount: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeMount: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeUnmount(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeUnmount(ctx, rpc.NoTarget, req)
}

// VolumeUnmountOn issues VolumeUnmount as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeUnmountOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeUnmount(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeUnmount(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeUnmount: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeUnmount: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeDelete(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeDelete(ctx, rpc.NoTarget, req)
}

// VolumeDeleteOn issues VolumeDelete as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeDeleteOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeDelete(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeDelete(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeDelete: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeDelete: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeMarkReadonly(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeMarkReadonly(ctx, rpc.NoTarget, req)
}

// VolumeMarkReadonlyOn issues VolumeMarkReadonly as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeMarkReadonlyOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeMarkReadonly(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeMarkReadonly(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeMarkReadonly: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeMarkReadonly: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeMarkWritable(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeMarkWritable(ctx, rpc.NoTarget, req)
}

// VolumeMarkWritableOn issues VolumeMarkWritable as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeMarkWritableOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeMarkWritable(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeMarkWritable(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeMarkWritable: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeMarkWritable: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeConfigure(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeConfigure(ctx, rpc.NoTarget, req)
}

// VolumeConfigureOn issues VolumeConfigure as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeConfigureOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeConfigure(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeConfigure(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeConfigure: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeConfigure: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeStatus(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeStatus(ctx, rpc.NoTarget, req)
}

// VolumeStatusOn issues VolumeStatus as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeStatusOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeStatus(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeStatus(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeStatus: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeStatus: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) GetState(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeGetState(ctx, rpc.NoTarget, req)
}

// GetStateOn issues GetState as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) GetStateOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeGetState(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeGetState(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.GetState: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.GetState: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) SetState(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeSetState(ctx, rpc.NoTarget, req)
}

// SetStateOn issues SetState as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) SetStateOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeSetState(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeSetState(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.SetState: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.SetState: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: VolumeCopy is a server-streaming RPC (VolumeServer.VolumeCopy); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) VolumeCopy(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeCopy(ctx, rpc.NoTarget, req)
}

// VolumeCopyOn issues VolumeCopy as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeCopyOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeCopy(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeCopy(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeCopy: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeCopy: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) ReadVolumeFileStatus(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeReadVolumeFileStatus(ctx, rpc.NoTarget, req)
}

// ReadVolumeFileStatusOn issues ReadVolumeFileStatus as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) ReadVolumeFileStatusOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeReadVolumeFileStatus(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeReadVolumeFileStatus(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.ReadVolumeFileStatus: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.ReadVolumeFileStatus: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: CopyFile is a server-streaming RPC (VolumeServer.CopyFile); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) CopyFile(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeCopyFile(ctx, rpc.NoTarget, req)
}

// CopyFileOn issues CopyFile as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) CopyFileOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeCopyFile(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeCopyFile(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.CopyFile: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.CopyFile: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: ReceiveFile is a client-streaming RPC (VolumeServer.ReceiveFile); its body lands
// when the transport streaming primitive ships. The request element schema is
// emitted (ReceiveFileRequest) and the unary response is returned here.
func (c *VolumeServerClient) ReceiveFile(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeReceiveFile(ctx, rpc.NoTarget, req)
}

// ReceiveFileOn issues ReceiveFile as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) ReceiveFileOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeReceiveFile(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeReceiveFile(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.ReceiveFile: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.ReceiveFile: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) ReadNeedleBlob(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeReadNeedleBlob(ctx, rpc.NoTarget, req)
}

// ReadNeedleBlobOn issues ReadNeedleBlob as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) ReadNeedleBlobOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeReadNeedleBlob(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeReadNeedleBlob(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.ReadNeedleBlob: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.ReadNeedleBlob: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) ReadNeedleMeta(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeReadNeedleMeta(ctx, rpc.NoTarget, req)
}

// ReadNeedleMetaOn issues ReadNeedleMeta as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) ReadNeedleMetaOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeReadNeedleMeta(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeReadNeedleMeta(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.ReadNeedleMeta: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.ReadNeedleMeta: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) WriteNeedleBlob(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeWriteNeedleBlob(ctx, rpc.NoTarget, req)
}

// WriteNeedleBlobOn issues WriteNeedleBlob as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) WriteNeedleBlobOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeWriteNeedleBlob(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeWriteNeedleBlob(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.WriteNeedleBlob: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.WriteNeedleBlob: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: ReadAllNeedles is a server-streaming RPC (VolumeServer.ReadAllNeedles); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) ReadAllNeedles(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeReadAllNeedles(ctx, rpc.NoTarget, req)
}

// ReadAllNeedlesOn issues ReadAllNeedles as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) ReadAllNeedlesOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeReadAllNeedles(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeReadAllNeedles(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.ReadAllNeedles: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.ReadAllNeedles: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: VolumeTailSender is a server-streaming RPC (VolumeServer.VolumeTailSender); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) VolumeTailSender(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeTailSender(ctx, rpc.NoTarget, req)
}

// VolumeTailSenderOn issues VolumeTailSender as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeTailSenderOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeTailSender(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeTailSender(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeTailSender: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeTailSender: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeTailReceiver(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeTailReceiver(ctx, rpc.NoTarget, req)
}

// VolumeTailReceiverOn issues VolumeTailReceiver as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeTailReceiverOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeTailReceiver(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeTailReceiver(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeTailReceiver: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeTailReceiver: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcShardsGenerate(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsGenerate(ctx, rpc.NoTarget, req)
}

// VolumeEcShardsGenerateOn issues VolumeEcShardsGenerate as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardsGenerateOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsGenerate(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardsGenerate(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsGenerate: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsGenerate: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcShardsRebuild(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsRebuild(ctx, rpc.NoTarget, req)
}

// VolumeEcShardsRebuildOn issues VolumeEcShardsRebuild as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardsRebuildOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsRebuild(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardsRebuild(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsRebuild: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsRebuild: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcShardsCopy(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsCopy(ctx, rpc.NoTarget, req)
}

// VolumeEcShardsCopyOn issues VolumeEcShardsCopy as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardsCopyOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsCopy(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardsCopy(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsCopy: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsCopy: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcShardsDelete(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsDelete(ctx, rpc.NoTarget, req)
}

// VolumeEcShardsDeleteOn issues VolumeEcShardsDelete as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardsDeleteOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsDelete(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardsDelete(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsDelete: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsDelete: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcShardsMount(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsMount(ctx, rpc.NoTarget, req)
}

// VolumeEcShardsMountOn issues VolumeEcShardsMount as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardsMountOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsMount(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardsMount(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsMount: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsMount: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcShardsUnmount(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsUnmount(ctx, rpc.NoTarget, req)
}

// VolumeEcShardsUnmountOn issues VolumeEcShardsUnmount as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardsUnmountOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsUnmount(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardsUnmount(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsUnmount: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsUnmount: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: VolumeEcShardRead is a server-streaming RPC (VolumeServer.VolumeEcShardRead); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) VolumeEcShardRead(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardRead(ctx, rpc.NoTarget, req)
}

// VolumeEcShardReadOn issues VolumeEcShardRead as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardReadOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardRead(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardRead(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardRead: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardRead: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcBlobDelete(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcBlobDelete(ctx, rpc.NoTarget, req)
}

// VolumeEcBlobDeleteOn issues VolumeEcBlobDelete as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcBlobDeleteOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcBlobDelete(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcBlobDelete(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeEcBlobDelete: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcBlobDelete: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcShardsToVolume(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsToVolume(ctx, rpc.NoTarget, req)
}

// VolumeEcShardsToVolumeOn issues VolumeEcShardsToVolume as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardsToVolumeOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsToVolume(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardsToVolume(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsToVolume: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsToVolume: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeEcShardsInfo(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsInfo(ctx, rpc.NoTarget, req)
}

// VolumeEcShardsInfoOn issues VolumeEcShardsInfo as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeEcShardsInfoOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeEcShardsInfo(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeEcShardsInfo(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsInfo: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeEcShardsInfo: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: VolumeTierMoveDatToRemote is a server-streaming RPC (VolumeServer.VolumeTierMoveDatToRemote); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) VolumeTierMoveDatToRemote(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeTierMoveDatToRemote(ctx, rpc.NoTarget, req)
}

// VolumeTierMoveDatToRemoteOn issues VolumeTierMoveDatToRemote as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeTierMoveDatToRemoteOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeTierMoveDatToRemote(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeTierMoveDatToRemote(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeTierMoveDatToRemote: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeTierMoveDatToRemote: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: VolumeTierMoveDatFromRemote is a server-streaming RPC (VolumeServer.VolumeTierMoveDatFromRemote); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) VolumeTierMoveDatFromRemote(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeTierMoveDatFromRemote(ctx, rpc.NoTarget, req)
}

// VolumeTierMoveDatFromRemoteOn issues VolumeTierMoveDatFromRemote as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeTierMoveDatFromRemoteOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeTierMoveDatFromRemote(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeTierMoveDatFromRemote(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeTierMoveDatFromRemote: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeTierMoveDatFromRemote: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeServerStatus(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeServerStatus(ctx, rpc.NoTarget, req)
}

// VolumeServerStatusOn issues VolumeServerStatus as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeServerStatusOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeServerStatus(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeServerStatus(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeServerStatus: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeServerStatus: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeServerLeave(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeServerLeave(ctx, rpc.NoTarget, req)
}

// VolumeServerLeaveOn issues VolumeServerLeave as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeServerLeaveOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeServerLeave(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeServerLeave(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeServerLeave: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeServerLeave: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) FetchAndWriteNeedle(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeFetchAndWriteNeedle(ctx, rpc.NoTarget, req)
}

// FetchAndWriteNeedleOn issues FetchAndWriteNeedle as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) FetchAndWriteNeedleOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeFetchAndWriteNeedle(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeFetchAndWriteNeedle(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.FetchAndWriteNeedle: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.FetchAndWriteNeedle: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) ScrubVolume(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeScrubVolume(ctx, rpc.NoTarget, req)
}

// ScrubVolumeOn issues ScrubVolume as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) ScrubVolumeOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeScrubVolume(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeScrubVolume(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.ScrubVolume: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.ScrubVolume: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) ScrubEcVolume(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeScrubEcVolume(ctx, rpc.NoTarget, req)
}

// ScrubEcVolumeOn issues ScrubEcVolume as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) ScrubEcVolumeOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeScrubEcVolume(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeScrubEcVolume(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.ScrubEcVolume: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.ScrubEcVolume: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// STREAMING: Query is a server-streaming RPC (VolumeServer.Query); its body lands
// when the transport streaming primitive ships. The per-message response schema
// is emitted; this stub issues the unary request envelope.
func (c *VolumeServerClient) Query(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeQuery(ctx, rpc.NoTarget, req)
}

// QueryOn issues Query as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) QueryOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeQuery(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeQuery(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.Query: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.Query: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) VolumeNeedleStatus(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokeVolumeNeedleStatus(ctx, rpc.NoTarget, req)
}

// VolumeNeedleStatusOn issues VolumeNeedleStatus as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) VolumeNeedleStatusOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeVolumeNeedleStatus(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokeVolumeNeedleStatus(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.VolumeNeedleStatus: %s", resp.Body)
		}
		return p, nil, fmt.Errorf("VolumeServer.VolumeNeedleStatus: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *VolumeServerClient) Ping(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.invokePing(ctx, rpc.NoTarget, req)
}

// PingOn issues Ping as a dependent call pipelined on the answer of on: the
// server substitutes on's resolved result for this call's payload before
// dispatch, so it ships without waiting for on to round-trip.
func (c *VolumeServerClient) PingOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokePing(ctx, on.ID, nil)
}

func (c *VolumeServerClient) invokePing(ctx context.Context, target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
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
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("VolumeServer.Ping: %s", resp.Body)
		}
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
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVacuumVolumeCheckOrdinal:
		body, err := h.VacuumVolumeCheck(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVacuumVolumeCompactOrdinal:
		body, err := h.VacuumVolumeCompact(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVacuumVolumeCommitOrdinal:
		body, err := h.VacuumVolumeCommit(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVacuumVolumeCleanupOrdinal:
		body, err := h.VacuumVolumeCleanup(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerDeleteCollectionOrdinal:
		body, err := h.DeleteCollection(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerAllocateVolumeOrdinal:
		body, err := h.AllocateVolume(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeSyncStatusOrdinal:
		body, err := h.VolumeSyncStatus(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeIncrementalCopyOrdinal:
		body, err := h.VolumeIncrementalCopy(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeMountOrdinal:
		body, err := h.VolumeMount(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeUnmountOrdinal:
		body, err := h.VolumeUnmount(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeDeleteOrdinal:
		body, err := h.VolumeDelete(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeMarkReadonlyOrdinal:
		body, err := h.VolumeMarkReadonly(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeMarkWritableOrdinal:
		body, err := h.VolumeMarkWritable(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeConfigureOrdinal:
		body, err := h.VolumeConfigure(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeStatusOrdinal:
		body, err := h.VolumeStatus(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerGetStateOrdinal:
		body, err := h.GetState(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerSetStateOrdinal:
		body, err := h.SetState(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeCopyOrdinal:
		body, err := h.VolumeCopy(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerReadVolumeFileStatusOrdinal:
		body, err := h.ReadVolumeFileStatus(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerCopyFileOrdinal:
		body, err := h.CopyFile(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerReceiveFileOrdinal:
		body, err := h.ReceiveFile(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerReadNeedleBlobOrdinal:
		body, err := h.ReadNeedleBlob(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerReadNeedleMetaOrdinal:
		body, err := h.ReadNeedleMeta(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerWriteNeedleBlobOrdinal:
		body, err := h.WriteNeedleBlob(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerReadAllNeedlesOrdinal:
		body, err := h.ReadAllNeedles(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeTailSenderOrdinal:
		body, err := h.VolumeTailSender(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeTailReceiverOrdinal:
		body, err := h.VolumeTailReceiver(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardsGenerateOrdinal:
		body, err := h.VolumeEcShardsGenerate(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardsRebuildOrdinal:
		body, err := h.VolumeEcShardsRebuild(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardsCopyOrdinal:
		body, err := h.VolumeEcShardsCopy(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardsDeleteOrdinal:
		body, err := h.VolumeEcShardsDelete(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardsMountOrdinal:
		body, err := h.VolumeEcShardsMount(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardsUnmountOrdinal:
		body, err := h.VolumeEcShardsUnmount(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardReadOrdinal:
		body, err := h.VolumeEcShardRead(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcBlobDeleteOrdinal:
		body, err := h.VolumeEcBlobDelete(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardsToVolumeOrdinal:
		body, err := h.VolumeEcShardsToVolume(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeEcShardsInfoOrdinal:
		body, err := h.VolumeEcShardsInfo(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeTierMoveDatToRemoteOrdinal:
		body, err := h.VolumeTierMoveDatToRemote(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeTierMoveDatFromRemoteOrdinal:
		body, err := h.VolumeTierMoveDatFromRemote(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeServerStatusOrdinal:
		body, err := h.VolumeServerStatus(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeServerLeaveOrdinal:
		body, err := h.VolumeServerLeave(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerFetchAndWriteNeedleOrdinal:
		body, err := h.FetchAndWriteNeedle(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerScrubVolumeOrdinal:
		body, err := h.ScrubVolume(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerScrubEcVolumeOrdinal:
		body, err := h.ScrubEcVolume(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerQueryOrdinal:
		body, err := h.Query(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerVolumeNeedleStatusOrdinal:
		body, err := h.VolumeNeedleStatus(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case VolumeServerPingOrdinal:
		body, err := h.Ping(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	default:
		return rpc.BuildResponse(rpc.StatusNotFound, call.PromiseID, nil), nil
	}
}
