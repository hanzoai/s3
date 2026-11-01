package types

import (
	"testing"
)

func TestNeedleIdBytesRoundTrip(t *testing.T) {
	for _, id := range []uint64{0, 1, 0xff, 0xdeadbeef, 1<<32 - 1, 1 << 32, 1 << 63, 1<<64 - 1} {
		bytes := make([]byte, NeedleIdSize)
		NeedleIdToBytes(bytes, NeedleId(id))
		got := BytesToNeedleId(bytes)
		if got != NeedleId(id) {
			t.Errorf("needle id %d round trip: got %d", id, got)
		}
	}
}

func TestNeedleIdStringAndParse(t *testing.T) {
	for _, id := range []uint64{0, 1, 0xabc, 0xdeadbeef, 1<<64 - 1} {
		s := NeedleId(id).String()
		got, err := ParseNeedleId(s)
		if err != nil {
			t.Errorf("ParseNeedleId(%q) error: %v", s, err)
			continue
		}
		if got != NeedleId(id) {
			t.Errorf("needle id %d string round trip: got %d", id, got)
		}
	}
}

func TestParseNeedleIdInvalid(t *testing.T) {
	for _, s := range []string{"", "xyz", "12xz", "-1", "0x1f"} {
		if _, err := ParseNeedleId(s); err == nil {
			t.Errorf("ParseNeedleId(%q) expected error, got nil", s)
		}
	}
}

func TestNeedleIdFileId(t *testing.T) {
	// volume 3, needle 0x1, zero cookie placeholder
	if got, want := NeedleId(1).FileId(3), "3,100000000"; got != want {
		t.Errorf("FileId: want %q, got %q", want, got)
	}
	if got, want := NeedleId(0x1234).FileId(7), "7,123400000000"; got != want {
		t.Errorf("FileId: want %q, got %q", want, got)
	}
}

func TestOffsetBytesRoundTrip(t *testing.T) {
	for _, v := range []uint32{0, 1, 0xff, 0x01020304, 1 << 31, 1<<32 - 1} {
		bytes := make([]byte, OffsetSize)
		OffsetToBytes(bytes, Uint32ToOffset(v))
		got := BytesToOffset(bytes)
		if got != Uint32ToOffset(v) {
			t.Errorf("offset %d round trip: got %v", v, got)
		}
	}
}

func TestOffsetIsZero(t *testing.T) {
	if !Uint32ToOffset(0).IsZero() {
		t.Error("zero offset should be IsZero")
	}
	for _, v := range []uint32{1, 0xff000000, 1 << 31} {
		if Uint32ToOffset(v).IsZero() {
			t.Errorf("offset %d should not be IsZero", v)
		}
	}
}

func TestToOffsetActualOffsetRoundTrip(t *testing.T) {
	// offsets are stored divided by NeedlePaddingSize, so exact round trips
	// hold for multiples of the padding size
	for _, actual := range []int64{0, 8, 64, 8 * 1024, 8 * 1024 * 1024, int64(MaxPossibleVolumeSize) - int64(NeedlePaddingSize)} {
		got := ToOffset(actual).ToActualOffset()
		if got != actual {
			t.Errorf("actual offset %d round trip: got %d", actual, got)
		}
	}
}

func TestToOffsetTruncatesToPadding(t *testing.T) {
	// non-aligned offsets round down to the padding boundary
	if got, want := ToOffset(13).ToActualOffset(), int64(8); got != want {
		t.Errorf("ToOffset(13).ToActualOffset(): want %d, got %d", want, got)
	}
	if got, want := ToOffset(15).ToActualOffset(), int64(8); got != want {
		t.Errorf("ToOffset(15).ToActualOffset(): want %d, got %d", want, got)
	}
}

func TestSizeSemantics(t *testing.T) {
	testCases := []struct {
		size        Size
		isTombstone bool
		isDeleted   bool
		isValid     bool
		raw         uint32
	}{
		{Size(0), false, false, false, 0},
		{Size(1), false, false, true, 1},
		{Size(1024), false, false, true, 1024},
		{TombstoneFileSize, true, true, false, 0},
		{Size(-5), false, true, false, 5},
	}
	for _, tc := range testCases {
		if got := tc.size.IsTombstone(); got != tc.isTombstone {
			t.Errorf("Size(%d).IsTombstone(): want %v, got %v", tc.size, tc.isTombstone, got)
		}
		if got := tc.size.IsDeleted(); got != tc.isDeleted {
			t.Errorf("Size(%d).IsDeleted(): want %v, got %v", tc.size, tc.isDeleted, got)
		}
		if got := tc.size.IsValid(); got != tc.isValid {
			t.Errorf("Size(%d).IsValid(): want %v, got %v", tc.size, tc.isValid, got)
		}
		if got := tc.size.Raw(); got != tc.raw {
			t.Errorf("Size(%d).Raw(): want %d, got %d", tc.size, tc.raw, got)
		}
	}
}

func TestSizeBytesRoundTrip(t *testing.T) {
	for _, s := range []Size{0, 1, 1024, TombstoneFileSize, Size(-5)} {
		bytes := make([]byte, SizeSize)
		SizeToBytes(bytes, s)
		got := BytesToSize(bytes)
		if got != s {
			t.Errorf("size %d round trip: got %d", s, got)
		}
	}
}

func TestCookieBytesRoundTrip(t *testing.T) {
	for _, c := range []uint32{0, 1, 0xdeadbeef, 1<<32 - 1} {
		bytes := make([]byte, CookieSize)
		CookieToBytes(bytes, Cookie(c))
		got := BytesToCookie(bytes)
		if got != Cookie(c) {
			t.Errorf("cookie %d round trip: got %d", c, got)
		}
	}
}

func TestParseCookie(t *testing.T) {
	got, err := ParseCookie("deadbeef")
	if err != nil {
		t.Fatalf("ParseCookie error: %v", err)
	}
	if got != Cookie(0xdeadbeef) {
		t.Errorf("ParseCookie: want %d, got %d", Cookie(0xdeadbeef), got)
	}
	for _, s := range []string{"", "xyz", "-1", "0x1f"} {
		if _, err := ParseCookie(s); err == nil {
			t.Errorf("ParseCookie(%q) expected error, got nil", s)
		}
	}
}

func TestToDiskType(t *testing.T) {
	testCases := []struct {
		input string
		want  DiskType
	}{
		{"", HardDriveType},
		{"hdd", HardDriveType},
		{"HDD", HardDriveType},
		{"ssd", SsdType},
		{"SSD", SsdType},
		{"Ssd", SsdType},
		{"nvme", DiskType("nvme")},
	}
	for _, tc := range testCases {
		if got := ToDiskType(tc.input); got != tc.want {
			t.Errorf("ToDiskType(%q): want %q, got %q", tc.input, tc.want, got)
		}
	}
}

func TestDiskTypeStrings(t *testing.T) {
	// HddType and HardDriveType both denote hard drives but are distinct
	// DiskType values ("hdd" vs ""); ToDiskType canonicalizes both to
	// HardDriveType. Callers that key maps by raw DiskType must account
	// for both spellings.
	var hdd DiskType = HddType
	if hdd == HardDriveType {
		t.Error("HddType and HardDriveType are distinct DiskType values")
	}
	if ToDiskType(string(HddType)) != HardDriveType {
		t.Error("ToDiskType should canonicalize HddType to HardDriveType")
	}
	if got := DiskType(SsdType).String(); got != "ssd" {
		t.Errorf("SsdType.String(): want %q, got %q", "ssd", got)
	}
	if got := HardDriveType.String(); got != "" {
		t.Errorf("HardDriveType.String(): want empty, got %q", got)
	}
	if got := HardDriveType.ReadableString(); got != HddType {
		t.Errorf("HardDriveType.ReadableString(): want %q, got %q", HddType, got)
	}
	if got := DiskType(SsdType).ReadableString(); got != "ssd" {
		t.Errorf("ssd ReadableString(): want %q, got %q", "ssd", got)
	}
}

func TestOnDiskLayoutConstants(t *testing.T) {
	if got, want := NeedleHeaderSize, CookieSize+NeedleIdSize+SizeSize; got != want {
		t.Errorf("NeedleHeaderSize: want %d, got %d", want, got)
	}
	if got, want := NeedleMapEntrySize, NeedleIdSize+OffsetSize+SizeSize; got != want {
		t.Errorf("NeedleMapEntrySize: want %d, got %d", want, got)
	}
}
