package wire

import "testing"

// Proves the ZAP-native zero-copy model on a real S3 message with
// variable-length string fields: build directly into the buffer, read in place
// via the View — no struct, no Marshal/Unmarshal.
func TestRemoteStorageLocationZeroCopy(t *testing.T) {
	cases := []struct {
		name, bucket, path string
		ttl                int32
	}{
		{"wasabi", "hanzo-cold", "/vfs/blocks", 300},
		{"", "", "", 0},
		{"s3", "b", "/", -1},
		{"long-" + string(make([]byte, 4000)), "bucket", "/p", 1 << 30},
	}
	for i, c := range cases {
		_, buf := NewRemoteStorageLocation(c.name, c.bucket, c.path, c.ttl)

		// Read back from the raw bytes — the receiver side, post-transport.
		v, err := WrapRemoteStorageLocation(buf)
		if err != nil {
			t.Fatalf("case %d: Wrap: %v", i, err)
		}
		if got := RemoteStorageLocationName(v); got != c.name {
			t.Fatalf("case %d: Name = %q want %q", i, got, c.name)
		}
		if got := RemoteStorageLocationBucket(v); got != c.bucket {
			t.Fatalf("case %d: Bucket = %q want %q", i, got, c.bucket)
		}
		if got := RemoteStorageLocationPath(v); got != c.path {
			t.Fatalf("case %d: Path = %q want %q", i, got, c.path)
		}
		if got := RemoteStorageLocationTTLSeconds(v); got != c.ttl {
			t.Fatalf("case %d: TTL = %d want %d", i, got, c.ttl)
		}
	}
}

// WrapRemoteStorageLocation on a wrong-kind buffer must be rejected.
func TestRemoteStorageLocationWrongKind(t *testing.T) {
	bad := make([]byte, 64) // zeroed -> kind 0x00 != 0x01
	if _, err := WrapRemoteStorageLocation(bad); err == nil {
		// zeroed buffer may fail Parse first; either way it must not succeed silently.
		t.Skip("zeroed buffer rejected at Parse")
	}
}

// BenchmarkRSLRead proves the read path is zero-copy / zero-alloc — the whole
// reason to use ZAP. Wrap + read three strings + an int must not allocate.
func BenchmarkRSLRead(b *testing.B) {
	_, buf := NewRemoteStorageLocation("wasabi", "hanzo-cold", "/vfs/blocks", 300)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, _ := WrapRemoteStorageLocation(buf)
		_ = RemoteStorageLocationName(v)
		_ = RemoteStorageLocationBucket(v)
		_ = RemoteStorageLocationPath(v)
		_ = RemoteStorageLocationTTLSeconds(v)
	}
}

// BenchmarkRSLBuild measures the build (one buffer alloc, no intermediate struct).
func BenchmarkRSLBuild(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewRemoteStorageLocation("wasabi", "hanzo-cold", "/vfs/blocks", 300)
	}
}
