package gentest

import "testing"

// Round-trips the CODE-GENERATED structs (remote_zap.go, emitted by
// protoc-gen-zap from remote.proto). Proves the generator output compiles and
// the codec it emits is correct across scalar/string/bool, nested message, and
// map<string,message> — i.e. the generator faithfully reproduces the
// hand-written reference in s3/zaprpc/codec_test.go.

func TestGeneratedRemoteStorageLocation(t *testing.T) {
	cases := []RemoteStorageLocation{
		{Name: "wasabi", Bucket: "hanzo-cold", Path: "/vfs/blocks", ListingCacheTtlSeconds: 300},
		{},
		{Name: "s3", Bucket: "b", Path: "/", ListingCacheTtlSeconds: -1},
	}
	for i, in := range cases {
		var out RemoteStorageLocation
		if err := out.UnmarshalZap(in.MarshalZap()); err != nil {
			t.Fatalf("case %d: %v", i, err)
		}
		if out != in {
			t.Fatalf("case %d mismatch:\n in=%+v\nout=%+v", i, in, out)
		}
	}
}

func TestGeneratedRemoteConfAllFields(t *testing.T) {
	in := RemoteConf{
		Type: "s3", Name: "wasabi", S3AccessKey: "AK", S3SecretKey: "SK",
		S3Region: "us-east-1", S3Endpoint: "s3.example", S3ForcePathStyle: true,
		S3SupportTagging: true, ContaboRegion: "eu", BackblazeKeyId: "bz",
	}
	var out RemoteConf
	if err := out.UnmarshalZap(in.MarshalZap()); err != nil {
		t.Fatal(err)
	}
	if out != in {
		t.Fatalf("RemoteConf mismatch:\n in=%+v\nout=%+v", in, out)
	}
}

func TestGeneratedMapAndNested(t *testing.T) {
	in := RemoteStorageMapping{
		Mappings: map[string]*RemoteStorageLocation{
			"primary": {Name: "wasabi", Bucket: "b1", Path: "/a", ListingCacheTtlSeconds: 60},
			"cold":    {Name: "storj", Bucket: "b2", Path: "/b"},
		},
		PrimaryBucketStorageName: "primary",
	}
	var out RemoteStorageMapping
	if err := out.UnmarshalZap(in.MarshalZap()); err != nil {
		t.Fatal(err)
	}
	if out.PrimaryBucketStorageName != in.PrimaryBucketStorageName || len(out.Mappings) != len(in.Mappings) {
		t.Fatalf("top-level mismatch: %+v", out)
	}
	for k, v := range in.Mappings {
		if got, ok := out.Mappings[k]; !ok || *got != *v {
			t.Fatalf("map[%q]: got %+v want %+v", k, got, v)
		}
	}
}
