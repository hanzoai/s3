//go:build linux

package discogc

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const (
	// Linux fallocate(2) flags for hole punching while keeping file size constant.
	fallocFlKeepSize  = 0x01
	fallocFlPunchHole = 0x02
)

// PunchHole issues a real Linux fallocate(FALLOC_FL_PUNCH_HOLE|FALLOC_FL_KEEP_SIZE)
// syscall on f for the byte range [offset, offset+length) and returns the actual
// physical disk bytes reclaimed as measured by syscall.Stat_t.Blocks * 512.
//
// If the underlying filesystem does not support hole punching (EOPNOTSUPP / ENOSYS),
// PunchHole returns (0, false, nil) so callers can record an honest fallback.
func PunchHole(f *os.File, offset, length int64) (reclaimedBytes int64, supported bool, err error) {
	if f == nil || length <= 0 || offset < 0 {
		return 0, true, nil
	}
	_, beforePhys, _, err := FileDiskUsage(f)
	if err != nil {
		return 0, false, err
	}

	err = syscall.Fallocate(int(f.Fd()), fallocFlKeepSize|fallocFlPunchHole, offset, length)
	if err != nil {
		if errors.Is(err, syscall.EOPNOTSUPP) || errors.Is(err, syscall.ENOSYS) {
			return 0, false, nil
		}
		return 0, true, err
	}

	_, afterPhys, _, err := FileDiskUsage(f)
	if err != nil {
		return 0, true, err
	}
	if beforePhys > afterPhys {
		reclaimedBytes = beforePhys - afterPhys
	}
	return reclaimedBytes, true, nil
}

// FileDiskUsage queries the OS inode metadata via fstat(2) and returns the
// logical file size (st_size), real allocated physical disk bytes
// (st_blocks * 512), and preferred filesystem block size (st_blksize).
func FileDiskUsage(f *os.File) (logicalBytes, physicalBytes, blockSize int64, err error) {
	var st syscall.Stat_t
	if err := syscall.Fstat(int(f.Fd()), &st); err != nil {
		return 0, 0, 0, err
	}
	blkSize := int64(st.Blksize)
	if blkSize <= 0 {
		blkSize = 4096
	}
	return st.Size, st.Blocks * 512, blkSize, nil
}

// ReadProcIO reads real process-level I/O accounting from /proc/self/io.
func ReadProcIO() ProcIOSnapshot {
	f, err := os.Open("/proc/self/io")
	if err != nil {
		return ProcIOSnapshot{Supported: false}
	}
	defer f.Close()

	snap := ProcIOSnapshot{Supported: true}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		idx := strings.IndexByte(line, ':')
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val, err := strconv.ParseInt(strings.TrimSpace(line[idx+1:]), 10, 64)
		if err != nil {
			continue
		}
		switch key {
		case "rchar":
			snap.RChar = val
		case "wchar":
			snap.WChar = val
		case "read_bytes":
			snap.ReadBytes = val
		case "write_bytes":
			snap.WriteBytes = val
		case "cancelled_write_bytes":
			snap.CancelledWriteBytes = val
		}
	}
	return snap
}
