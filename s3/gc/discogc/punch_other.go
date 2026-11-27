//go:build !linux

package discogc

import "os"

// PunchHole reports unsupported on non-Linux platforms where Linux
// FALLOC_FL_PUNCH_HOLE is unavailable.
func PunchHole(_ *os.File, _, _ int64) (reclaimedBytes int64, supported bool, err error) {
	return 0, false, nil
}

// FileDiskUsage falls back to os.File.Stat logical size on non-Linux platforms.
func FileDiskUsage(f *os.File) (logicalBytes, physicalBytes, blockSize int64, err error) {
	info, err := f.Stat()
	if err != nil {
		return 0, 0, 0, err
	}
	sz := info.Size()
	return sz, sz, 4096, nil
}

// ReadProcIO returns unsupported on non-Linux platforms.
func ReadProcIO() ProcIOSnapshot {
	return ProcIOSnapshot{Supported: false}
}
