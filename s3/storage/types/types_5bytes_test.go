//go:build 5BytesOffset

package types

import (
	"testing"
)

func TestOffsetBytesRoundTripFiveBytes(t *testing.T) {
	// offsets above 32 GiB carry a non-zero b4 byte; exercise it through
	// serialization, which Uint32ToOffset (uint32-only) cannot reach
	for _, actual := range []int64{33 << 30, 64 << 30, 1 << 40, int64(MaxPossibleVolumeSize) - int64(NeedlePaddingSize)} {
		offset := ToOffset(actual)
		bytes := make([]byte, OffsetSize)
		OffsetToBytes(bytes, offset)
		got := BytesToOffset(bytes)
		if got != offset {
			t.Errorf("offset for actual %d byte round trip: got %v", actual, got)
			continue
		}
		if got.ToActualOffset() != actual {
			t.Errorf("actual offset %d round trip: got %d", actual, got.ToActualOffset())
		}
	}
}
