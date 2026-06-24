// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	bootstrapStateLastScannedPathOff = 0
	bootstrapStateStartedAtNsOff     = 8
	bootstrapStateCompletedAtNsOff   = 16
	bootstrapStateSize               = 24
)

// BootstrapState is a zero-copy view into a ZAP-encoded BootstrapState message.
type BootstrapState struct{ o zap.Object }

// WrapBootstrapState parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapBootstrapState(b []byte) (BootstrapState, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return BootstrapState{}, err
	}
	return BootstrapState{o: m.Root()}, nil
}

func (t BootstrapState) LastScannedPath() string { return t.o.Text(bootstrapStateLastScannedPathOff) }
func (t BootstrapState) StartedAtNs() int64      { return t.o.Int64(bootstrapStateStartedAtNsOff) }
func (t BootstrapState) CompletedAtNs() int64    { return t.o.Int64(bootstrapStateCompletedAtNsOff) }

// BootstrapStateInput collects the field values for NewBootstrapState.
type BootstrapStateInput struct {
	LastScannedPath string
	StartedAtNs     int64
	CompletedAtNs   int64
}

// NewBootstrapState builds a ZAP-encoded BootstrapState message from in and returns the bytes.
func NewBootstrapState(in BootstrapStateInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(bootstrapStateSize)
	ob.SetText(bootstrapStateLastScannedPathOff, in.LastScannedPath)
	ob.SetInt64(bootstrapStateStartedAtNsOff, in.StartedAtNs)
	ob.SetInt64(bootstrapStateCompletedAtNsOff, in.CompletedAtNs)
	ob.FinishAsRoot()
	return b.Finish()
}
