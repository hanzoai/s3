package luxdb

import (
	"bytes"
	"context"
	"testing"

	"github.com/hanzoai/s3/s3/filer"
	"github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/util"
)

// runBackendRoundTrip exercises the full FilerStore + KV surface against the
// named backend. It is shared by the default (zapdb) test and the tag-gated
// pebble/leveldb tests so every selectable backend is proven, not just compiled.
func runBackendRoundTrip(t *testing.T, backend string) {
	t.Helper()
	store := &LuxDBStore{}
	if err := store.initialize(backend, t.TempDir()); err != nil {
		t.Fatalf("init %s: %v", backend, err)
	}
	defer store.Shutdown()
	ctx := context.Background()

	// Insert two entries under the same directory and one elsewhere.
	mk := func(path string) *filer.Entry {
		return &filer.Entry{FullPath: util.FullPath(path), Attr: filer.Attr{Mode: 0640}}
	}
	for _, p := range []string{"/d/a.txt", "/d/b.txt", "/other/c.txt"} {
		if err := store.InsertEntry(ctx, mk(p)); err != nil {
			t.Fatalf("%s insert %s: %v", backend, p, err)
		}
	}

	// FindEntry hits.
	if _, err := store.FindEntry(ctx, util.FullPath("/d/a.txt")); err != nil {
		t.Fatalf("%s find a.txt: %v", backend, err)
	}
	// FindEntry miss maps to filer_pb.ErrNotFound.
	if _, err := store.FindEntry(ctx, util.FullPath("/d/missing")); err != filer_pb.ErrNotFound {
		t.Fatalf("%s find missing: want ErrNotFound, got %v", backend, err)
	}

	// ListDirectoryEntries returns exactly the two children of /d in order.
	var names []string
	if _, err := store.ListDirectoryEntries(ctx, util.FullPath("/d"), "", true, 100,
		func(e *filer.Entry) (bool, error) {
			names = append(names, e.Name())
			return true, nil
		}); err != nil {
		t.Fatalf("%s list /d: %v", backend, err)
	}
	if len(names) != 2 || names[0] != "a.txt" || names[1] != "b.txt" {
		t.Fatalf("%s list /d = %v, want [a.txt b.txt]", backend, names)
	}

	// DeleteFolderChildren clears /d but not /other.
	if err := store.DeleteFolderChildren(ctx, util.FullPath("/d")); err != nil {
		t.Fatalf("%s delete children /d: %v", backend, err)
	}
	if _, err := store.FindEntry(ctx, util.FullPath("/d/a.txt")); err != filer_pb.ErrNotFound {
		t.Fatalf("%s a.txt should be gone: %v", backend, err)
	}
	if _, err := store.FindEntry(ctx, util.FullPath("/other/c.txt")); err != nil {
		t.Fatalf("%s c.txt should survive: %v", backend, err)
	}

	// KV surface.
	if err := store.KvPut(ctx, []byte("k"), []byte("v")); err != nil {
		t.Fatalf("%s kv put: %v", backend, err)
	}
	got, err := store.KvGet(ctx, []byte("k"))
	if err != nil || !bytes.Equal(got, []byte("v")) {
		t.Fatalf("%s kv get = %q, %v", backend, got, err)
	}
	if err := store.KvDelete(ctx, []byte("k")); err != nil {
		t.Fatalf("%s kv delete: %v", backend, err)
	}
	if _, err := store.KvGet(ctx, []byte("k")); err != filer.ErrKvNotFound {
		t.Fatalf("%s kv get after delete: want ErrKvNotFound, got %v", backend, err)
	}
}

// TestZapdbBackendRoundTrip proves the default backend on the default build.
func TestZapdbBackendRoundTrip(t *testing.T) {
	runBackendRoundTrip(t, "zapdb")
}

// TestDefaultBackendRoundTrip proves an empty backend string defaults to zapdb.
func TestDefaultBackendRoundTrip(t *testing.T) {
	runBackendRoundTrip(t, "")
}
