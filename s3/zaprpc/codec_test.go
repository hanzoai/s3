package zaprpc

import (
	"bytes"
	"testing"
)

// remoteStorageLocation mirrors remote_pb.RemoteStorageLocation and is the
// hand-written reference for what cmd/protoc-gen-zap emits: plain struct +
// MarshalZap/UnmarshalZap over the hand-rolled codec. Fields carry their proto
// field numbers.
type remoteStorageLocation struct {
	Name                   string // 1
	Bucket                 string // 2
	Path                   string // 3
	ListingCacheTTLSeconds int32  // 4
}

func (m *remoteStorageLocation) MarshalZap() []byte {
	w := NewWriter()
	w.String(1, m.Name)
	w.String(2, m.Bucket)
	w.String(3, m.Path)
	w.Int32(4, m.ListingCacheTTLSeconds)
	return w.Result()
}

func (m *remoteStorageLocation) UnmarshalZap(b []byte) error {
	r := NewReader(b)
	for r.Next() {
		switch r.Field() {
		case 1:
			m.Name = r.String()
		case 2:
			m.Bucket = r.String()
		case 3:
			m.Path = r.String()
		case 4:
			m.ListingCacheTTLSeconds = r.Int32()
		default:
			r.Skip()
		}
	}
	return r.Err()
}

// remoteStorageMapping mirrors remote_pb.RemoteStorageMapping: a map field
// (encoded as repeated key-value entries) plus a scalar, exercising nested
// messages and repeated fields — the bulk of real SeaweedFS message shapes.
type remoteStorageMapping struct {
	Mappings                 map[string]*remoteStorageLocation // 1 (map<string,msg>)
	PrimaryBucketStorageName string                            // 2
}

func (m *remoteStorageMapping) MarshalZap() []byte {
	w := NewWriter()
	for k, v := range m.Mappings {
		// map entry = embedded message { key=1 (string), value=2 (message) }
		ew := NewWriter()
		ew.String(1, k)
		ew.Message(2, v.MarshalZap())
		w.Message(1, ew.Result())
	}
	w.String(2, m.PrimaryBucketStorageName)
	return w.Result()
}

func (m *remoteStorageMapping) UnmarshalZap(b []byte) error {
	r := NewReader(b)
	for r.Next() {
		switch r.Field() {
		case 1:
			m.ensure()
			er := NewReader(r.Message())
			var k string
			v := &remoteStorageLocation{}
			for er.Next() {
				switch er.Field() {
				case 1:
					k = er.String()
				case 2:
					_ = v.UnmarshalZap(er.Message())
				default:
					er.Skip()
				}
			}
			if er.Err() != nil {
				return er.Err()
			}
			m.Mappings[k] = v
		case 2:
			m.PrimaryBucketStorageName = r.String()
		default:
			r.Skip()
		}
	}
	return r.Err()
}

func (m *remoteStorageMapping) ensure() {
	if m.Mappings == nil {
		m.Mappings = map[string]*remoteStorageLocation{}
	}
}

func TestCodecScalarStringRoundTrip(t *testing.T) {
	cases := []remoteStorageLocation{
		{Name: "wasabi", Bucket: "hanzo-cold", Path: "/vfs/blocks", ListingCacheTTLSeconds: 300},
		{},                                  // all-zero -> empty wire -> back to all-zero
		{Name: "s3", Bucket: "b", Path: "/", ListingCacheTTLSeconds: -1}, // negative int32 must round-trip
		{Name: string(bytes.Repeat([]byte("x"), 5000)), ListingCacheTTLSeconds: 1 << 30},
	}
	for i, in := range cases {
		var out remoteStorageLocation
		if err := out.UnmarshalZap(in.MarshalZap()); err != nil {
			t.Fatalf("case %d: %v", i, err)
		}
		if out != in {
			t.Fatalf("case %d mismatch:\n in=%+v\nout=%+v", i, in, out)
		}
	}
	// all-zero must encode to zero bytes (proto3 semantics).
	if w := (&remoteStorageLocation{}).MarshalZap(); len(w) != 0 {
		t.Fatalf("all-zero should marshal empty, got %d bytes", len(w))
	}
}

func TestCodecNestedAndMapRoundTrip(t *testing.T) {
	in := remoteStorageMapping{
		Mappings: map[string]*remoteStorageLocation{
			"primary": {Name: "wasabi", Bucket: "b1", Path: "/a", ListingCacheTTLSeconds: 60},
			"cold":    {Name: "storj", Bucket: "b2", Path: "/b"},
		},
		PrimaryBucketStorageName: "primary",
	}
	var out remoteStorageMapping
	if err := out.UnmarshalZap(in.MarshalZap()); err != nil {
		t.Fatal(err)
	}
	if out.PrimaryBucketStorageName != in.PrimaryBucketStorageName || len(out.Mappings) != len(in.Mappings) {
		t.Fatalf("top-level mismatch: %+v", out)
	}
	for k, v := range in.Mappings {
		got, ok := out.Mappings[k]
		if !ok || *got != *v {
			t.Fatalf("map[%q] mismatch: got %+v want %+v", k, got, v)
		}
	}
}

// TestCodecForwardCompat proves a reader skips fields it doesn't know — the
// invariant that makes rolling cluster upgrades safe.
func TestCodecForwardCompat(t *testing.T) {
	// A "v+1" writer emits an extra field 99 (a string) the "v" reader ignores.
	w := NewWriter()
	w.String(1, "keepme")
	w.String(99, "future-field-the-old-node-never-heard-of")
	w.Int32(4, 7)
	var out remoteStorageLocation
	if err := out.UnmarshalZap(w.Result()); err != nil {
		t.Fatal(err)
	}
	if out.Name != "keepme" || out.ListingCacheTTLSeconds != 7 {
		t.Fatalf("forward-compat: known fields lost: %+v", out)
	}
}
