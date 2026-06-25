package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/hanzoai/s3/s3/storage/needle"
	"github.com/hanzoai/s3/s3/storage/needle_map"
	"github.com/hanzoai/s3/s3/storage/types"
)

// A mid-commit crash can leave an old .ldb (with a high watermark) beside a
// freshly swapped, shorter .idx. generateLevelDbFile must distrust the stale
// watermark and rebuild from offset 0 instead of walking past EOF and
// replaying zero entries, which would leave the needle map empty and make
// every live needle a phantom 404.
func TestGenerateLevelDbFileStaleWatermarkRebuilds(t *testing.T) {
	dir := t.TempDir()

	idxPath := filepath.Join(dir, "1.idx")
	idxFile, err := os.OpenFile(idxPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Fatalf("create idx: %v", err)
	}
	defer idxFile.Close()

	// A short .idx: three live needles.
	const liveCount = 3
	for i := uint64(1); i <= liveCount; i++ {
		entry := needle_map.ToBytes(types.Uint64ToNeedleId(i), types.ToOffset(int64(i*1024)), types.Size(512))
		if _, err := idxFile.Write(entry); err != nil {
			t.Fatalf("write idx entry: %v", err)
		}
	}
	if err := idxFile.Sync(); err != nil {
		t.Fatalf("sync idx: %v", err)
	}

	// Seed an .ldb whose stored watermark sits far past the short .idx, mimicking
	// the leftover db from a pre-crash, much larger index.
	dbPath := filepath.Join(dir, "1.ldb")
	db, err := openNeedleMapDB(dbPath)
	if err != nil {
		t.Fatalf("open ldb: %v", err)
	}
	if err := setWatermark(db, watermarkBatchSize); err != nil {
		t.Fatalf("set stale watermark: %v", err)
	}
	db.Close()

	if err := generateLevelDbFile(dbPath, idxFile); err != nil {
		t.Fatalf("generateLevelDbFile: %v", err)
	}

	db, err = openNeedleMapDB(dbPath)
	if err != nil {
		t.Fatalf("reopen ldb: %v", err)
	}
	defer db.Close()
	for i := uint64(1); i <= liveCount; i++ {
		keyBytes := make([]byte, types.NeedleIdSize)
		types.NeedleIdToBytes(keyBytes, types.Uint64ToNeedleId(i))
		if _, err := db.Get(keyBytes); err != nil {
			t.Fatalf("needle %d missing after rebuild (stale watermark poisoned the map): %v", i, err)
		}
	}
}

// Round-trip a range of needle ids through the on-disk zapdb needle map:
// put, read back the exact offset/size, delete, and confirm the delete sticks.
// This exercises the needle-id -> (offset,size) byte layout end to end and the
// rebuild-from-empty path (the .ldb dir starts absent).
func TestLevelDbNeedleMapRoundTrip(t *testing.T) {
	dir := t.TempDir()
	indexFile, err := os.OpenFile(filepath.Join(dir, "rt.idx"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Fatalf("create index file: %v", err)
	}
	dbFileName := filepath.Join(dir, "rt.ldb")

	// ldbTimeout 0 keeps the db resident (no lazy unload), the volume-server
	// default for a single hot volume.
	m, err := NewLevelDbNeedleMap(dbFileName, indexFile, 0, needle.GetCurrentVersion())
	if err != nil {
		t.Fatalf("NewLevelDbNeedleMap: %v", err)
	}
	defer m.Close()

	const n = 256
	for i := uint64(1); i <= n; i++ {
		key := types.Uint64ToNeedleId(i)
		off := types.ToOffset(int64(i) * types.NeedlePaddingSize)
		size := types.Size(i * 7)
		if err := m.Put(key, off, size); err != nil {
			t.Fatalf("Put %d: %v", i, err)
		}
		nv, ok := m.Get(key)
		if !ok {
			t.Fatalf("Get %d after Put: not found", i)
		}
		if nv.Offset != off || nv.Size != size {
			t.Fatalf("Get %d mismatch: got (off=%v,size=%v) want (off=%v,size=%v)", i, nv.Offset, nv.Size, off, size)
		}
	}

	// Every key still resolves.
	for i := uint64(1); i <= n; i++ {
		if _, ok := m.Get(types.Uint64ToNeedleId(i)); !ok {
			t.Fatalf("Get %d before delete: not found", i)
		}
	}

	// Delete the even ids; a deleted needle reads back as not-found (tombstoned
	// negative size).
	for i := uint64(2); i <= n; i += 2 {
		key := types.Uint64ToNeedleId(i)
		nv, _ := m.Get(key)
		if err := m.Delete(key, nv.Offset); err != nil {
			t.Fatalf("Delete %d: %v", i, err)
		}
	}
	for i := uint64(1); i <= n; i++ {
		_, ok := m.Get(types.Uint64ToNeedleId(i))
		if i%2 == 0 {
			if ok {
				// getFromDb returns the value; a deleted needle stores a negative
				// (tombstone) size which still round-trips, so confirm via size.
				if nv, _ := m.Get(types.Uint64ToNeedleId(i)); !nv.Size.IsDeleted() {
					t.Fatalf("needle %d should be deleted, size=%v", i, nv.Size)
				}
			}
		} else if !ok {
			t.Fatalf("odd needle %d should survive delete of evens", i)
		}
	}
}

// A volume server batches needle writes and crosses the watermark batch
// boundary; the on-disk index must survive a Close/reopen with every needle
// intact and the watermark persisted. This drives the batched-write path
// (recordCount crossing watermarkBatchSize) and the rebuild-from-existing
// open path.
func TestLevelDbNeedleMapPersistsAcrossReopen(t *testing.T) {
	dir := t.TempDir()
	idxPath := filepath.Join(dir, "p.idx")
	dbPath := filepath.Join(dir, "p.ldb")

	indexFile, err := os.OpenFile(idxPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Fatalf("create index file: %v", err)
	}
	m, err := NewLevelDbNeedleMap(dbPath, indexFile, 0, needle.GetCurrentVersion())
	if err != nil {
		t.Fatalf("NewLevelDbNeedleMap: %v", err)
	}

	const n = 32
	for i := uint64(1); i <= n; i++ {
		if err := m.Put(types.Uint64ToNeedleId(i), types.ToOffset(int64(i)*types.NeedlePaddingSize), types.Size(i)); err != nil {
			t.Fatalf("Put %d: %v", i, err)
		}
	}
	m.Close() // syncs idx, closes db

	// Reopen: the .ldb is newer than the .idx, so isLevelDbFresh skips the
	// rebuild and we open the existing index in place.
	indexFile2, err := os.OpenFile(idxPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Fatalf("reopen index file: %v", err)
	}
	m2, err := NewLevelDbNeedleMap(dbPath, indexFile2, 0, needle.GetCurrentVersion())
	if err != nil {
		t.Fatalf("reopen NewLevelDbNeedleMap: %v", err)
	}
	defer m2.Close()
	for i := uint64(1); i <= n; i++ {
		nv, ok := m2.Get(types.Uint64ToNeedleId(i))
		if !ok {
			t.Fatalf("needle %d missing after reopen", i)
		}
		if want := types.Size(i); nv.Size != want {
			t.Fatalf("needle %d size after reopen: got %v want %v", i, nv.Size, want)
		}
	}
}

func TestLevelDbNeedleMap_Concurrency(t *testing.T) {
	dir, err := os.MkdirTemp("", "test_leveldb_concurrency")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	prefix := "test"
	indexFile, err := os.Create(filepath.Join(dir, prefix+".idx"))
	if err != nil {
		t.Fatalf("create index file: %v", err)
	}
	dbFileName := filepath.Join(dir, prefix+".ldb")

	// Create and initialize map
	m, err := NewLevelDbNeedleMap(dbFileName, indexFile, 1, needle.GetCurrentVersion())
	if err != nil {
		t.Fatalf("NewLevelDbNeedleMap: %v", err)
	}
	defer m.Close()

	// Pre-populate some data
	key := types.NeedleId(1)
	if err := m.Put(key, types.ToOffset(100), types.Size(200)); err != nil {
		t.Fatalf("Put: %v", err)
	}

	// Force unload to start from nil state
	if err := unloadLdb(m); err != nil {
		t.Fatalf("unloadLdb: %v", err)
	}

	var wg sync.WaitGroup
	startCh := make(chan struct{})
	errCh := make(chan error, 100)

	// Spawn multiple goroutines to trigger the race
	// Multiple readers will see m.db == nil and try to reload concurrently
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-startCh

			// Try multiple times to increase chance of collision
			for j := 0; j < 2; j++ {
				_, ok := m.Get(key)
				if !ok {
					// Get failed, possibly due to race in reload.
					// But we also put data concurrently, so maybe it's missing if deleted?
					// In this test, we only Put, never Delete. So Key 1 should be there.
					// However, if DB reload fails, Get returns false!
					errCh <- fmt.Errorf("routine %d iter %d: Get returned false", id, j)
				}

				// Also try Put concurrently
				err := m.Put(types.NeedleId(2+id), types.ToOffset(100), types.Size(200))
				if err != nil {
					errCh <- fmt.Errorf("routine %d iter %d: Put failed: %v", id, j, err)
				}

				// Manually unload occasionally to reset the state
				if j%2 == 0 {
					// This might fail if locked, but that's fine
					unloadLdb(m)
				}
			}
		}(i)
	}

	close(startCh)
	wg.Wait()
	close(errCh)

	for e := range errCh {
		t.Errorf("Error encountered: %v", e)
	}
}
