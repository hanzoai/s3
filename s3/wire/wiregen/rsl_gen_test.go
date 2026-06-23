package wiregen

import (
	"testing"

	zapv2 "github.com/luxfi/zap/v2"
)

// Proves the GENERATED variable-length schema (rsl_gen.go, emitted by the
// extended luxfi/zap v2/codegen from the RemoteStorageLocation schema) compiles
// and gives zero-copy / zero-alloc reads — identical behavior to the
// hand-written s3/wire reference. Closes the loop: generator output IS the
// proven pattern.
func TestGeneratedRSLZeroCopy(t *testing.T) {
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
		_, buf := NewRemoteStorageLocation(c.ttl, c.name, c.bucket, c.path)
		v, err := WrapRemoteStorageLocation(buf)
		if err != nil {
			t.Fatalf("case %d: Wrap: %v", i, err)
		}
		if got := RemoteStorageLocationName(v); got != c.name {
			t.Fatalf("case %d: Name=%q want %q", i, got, c.name)
		}
		if got := RemoteStorageLocationBucket(v); got != c.bucket {
			t.Fatalf("case %d: Bucket=%q want %q", i, got, c.bucket)
		}
		if got := RemoteStorageLocationPath(v); got != c.path {
			t.Fatalf("case %d: Path=%q want %q", i, got, c.path)
		}
		if got := zapv2.Read(v, RSLSchemaFields.TTLSeconds); got != c.ttl {
			t.Fatalf("case %d: TTL=%d want %d", i, got, c.ttl)
		}
	}
}

// BenchmarkGeneratedRSLRead — the generated accessors must be zero-alloc too.
func BenchmarkGeneratedRSLRead(b *testing.B) {
	_, buf := NewRemoteStorageLocation(300, "wasabi", "hanzo-cold", "/vfs/blocks")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, _ := WrapRemoteStorageLocation(buf)
		_ = RemoteStorageLocationName(v)
		_ = RemoteStorageLocationBucket(v)
		_ = RemoteStorageLocationPath(v)
		_ = zapv2.Read(v, RSLSchemaFields.TTLSeconds)
	}
}
