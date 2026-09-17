package filer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	util_http "github.com/hanzoai/s3/s3/util/http"
)

// These pin one property: a download that failed is not remembered. The error
// it holds was true when it was fetched -- a volume still registering, a server
// shedding load -- and a later read of the same file id must fetch again rather
// than be answered with that stale failure for the life of the process.
//
// It is worth its own file because the failure it guards against is invisible
// from outside the process: the lookup behind the cached error heals, every
// fresh caller succeeds, and only reads of the specific file ids that were
// unlucky once keep failing -- with an error message that is minutes old.

var initHTTPOnce sync.Once

func initHTTP() { initHTTPOnce.Do(util_http.InitGlobalHttpClient) }

// chunkServer serves payload for any path, after answering the first `fail`
// requests with 429 -- a status the direct fetch path treats as non-retryable,
// so the failure lands on the first attempt instead of after a backoff.
func chunkServer(t *testing.T, payload []byte, fail int32) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hits.Add(1) <= fail {
			http.Error(w, "shed", http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Length", "")
		_, _ = w.Write(payload)
	}))
	t.Cleanup(srv.Close)
	return srv, &hits
}

func payload(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(i % 251)
	}
	return b
}

func TestReaderCacheRetriesAFailedLookup(t *testing.T) {
	initHTTP()
	want := payload(512)
	srv, _ := chunkServer(t, want, 0)

	var calls atomic.Int32
	lookup := func(ctx context.Context, fileId string) ([]string, error) {
		if calls.Add(1) == 1 {
			return nil, errors.New("volume 78 not found") // the registration race
		}
		return []string{srv.URL + "/" + fileId}, nil
	}
	rc := NewReaderCache(10, newMockChunkCacheForReaderCache(), lookup)
	defer rc.destroy()

	buf := make([]byte, len(want))
	if _, err := rc.ReadChunkAt(context.Background(), buf, "78,294ca1219f760a", nil, false, 0, len(want), false); err == nil {
		t.Fatal("first read: want the lookup failure, got nil")
	}

	n, err := rc.ReadChunkAt(context.Background(), buf, "78,294ca1219f760a", nil, false, 0, len(want), false)
	if err != nil {
		t.Fatalf("second read: the lookup has healed, so the read must too; got %v", err)
	}
	if n != len(want) || !bytes.Equal(buf, want) {
		t.Fatalf("second read: got %d bytes, want the %d-byte chunk", n, len(want))
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("lookup called %d times, want 2: the failure must not be replayed from the cache", got)
	}

	// Recovered, the chunk is data like any other and is reused, not re-fetched.
	again := make([]byte, len(want))
	if n, err := rc.ReadChunkAt(context.Background(), again, "78,294ca1219f760a", nil, false, 0, len(want), false); err != nil || n != len(want) || !bytes.Equal(again, want) {
		t.Fatalf("third read: n=%d err=%v, want the cached chunk", n, err)
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("lookup called %d times after recovery, want still 2", got)
	}
}

// A lookup that answers with no locations is a failure too, and is retried.
func TestReaderCacheRetriesAnEmptyLookup(t *testing.T) {
	initHTTP()
	want := payload(200)
	srv, _ := chunkServer(t, want, 0)

	var calls atomic.Int32
	lookup := func(ctx context.Context, fileId string) ([]string, error) {
		if calls.Add(1) == 1 {
			return nil, nil
		}
		return []string{srv.URL + "/" + fileId}, nil
	}
	rc := NewReaderCache(10, newMockChunkCacheForReaderCache(), lookup)
	defer rc.destroy()

	buf := make([]byte, len(want))
	_, err := rc.ReadChunkAt(context.Background(), buf, "76,empty", nil, false, 0, len(want), false)
	if err == nil || !strings.Contains(err.Error(), "urls not found") {
		t.Fatalf("first read: want an explicit urls-not-found failure, got %v", err)
	}
	if n, err := rc.ReadChunkAt(context.Background(), buf, "76,empty", nil, false, 0, len(want), false); err != nil || n != len(want) {
		t.Fatalf("second read: n=%d err=%v", n, err)
	}
}

// A prefetch that failed must not stand between a later demand read and a
// fresh fetch of the same file id.
func TestFailedPrefetchIsRetriedOnDemand(t *testing.T) {
	initHTTP()
	want := payload(96)
	srv, _ := chunkServer(t, want, 0)

	var healed atomic.Bool
	var calls atomic.Int32
	lookup := func(ctx context.Context, fileId string) ([]string, error) {
		calls.Add(1)
		if !healed.Load() {
			return nil, errors.New("volume not registered yet")
		}
		return []string{srv.URL + "/" + fileId}, nil
	}
	rc := NewReaderCache(10, newMockChunkCacheForReaderCache(), lookup)
	defer rc.destroy()

	const fid = "79,prefetch"
	prefetch(t, rc, fid, len(want))

	healed.Store(true)
	buf := make([]byte, len(want))
	if n, err := rc.ReadChunkAt(context.Background(), buf, fid, nil, false, 0, len(want), false); err != nil || n != len(want) || !bytes.Equal(buf, want) {
		t.Fatalf("demand read after a failed prefetch: n=%d err=%v", n, err)
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("lookup called %d times, want 2 (the failed prefetch, then one retry)", got)
	}

	again := make([]byte, len(want))
	if n, err := rc.ReadChunkAt(context.Background(), again, fid, nil, false, 0, len(want), false); err != nil || n != len(want) || !bytes.Equal(again, want) {
		t.Fatalf("read after recovery: n=%d err=%v", n, err)
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("lookup called %d times after recovery, want still 2: the recovered chunk is reused", got)
	}
}

// Failed prefetches that fill every slot must not starve a new file: they are
// completed failures now, so the size sweep reclaims them.
func TestFailedPrefetchesDoNotStarveTheLimit(t *testing.T) {
	initHTTP()
	want := payload(48)
	srv, _ := chunkServer(t, want, 0)

	lookup := func(ctx context.Context, fileId string) ([]string, error) {
		if fileId == "fresh" {
			return []string{srv.URL + "/fresh"}, nil
		}
		return nil, errors.New("volume not registered yet")
	}
	const limit = 3
	rc := NewReaderCache(limit, newMockChunkCacheForReaderCache(), lookup)
	defer rc.destroy()

	for i := 0; i < limit; i++ {
		prefetch(t, rc, fmt.Sprintf("78,fail-%d", i), len(want))
	}

	buf := make([]byte, len(want))
	if n, err := rc.ReadChunkAt(context.Background(), buf, "fresh", nil, false, 0, len(want), false); err != nil || n != len(want) || !bytes.Equal(buf, want) {
		t.Fatalf("a new file with every slot held by a failure: n=%d err=%v", n, err)
	}
	rc.Lock()
	held := len(rc.downloaders)
	rc.Unlock()
	if held > limit {
		t.Fatalf("%d downloaders against a limit of %d: the failed prefetches were not reclaimed", held, limit)
	}
}

// prefetch starts a background download for fid and waits for it to finish
// FAILING -- on the cacher's own completion signal, not a sleep -- so the read
// that follows meets a completed error rather than joining a download in flight.
func prefetch(t *testing.T, rc *ReaderCache, fid string, size int) {
	t.Helper()
	rc.MaybeCache(&Interval[*ChunkView]{Value: &ChunkView{
		FileId: fid, ViewSize: uint64(size), ChunkSize: uint64(size),
	}}, 1)
	rc.Lock()
	c := rc.downloaders[fid]
	rc.Unlock()
	if c == nil {
		t.Fatalf("prefetch of %s did not start", fid)
	}
	select {
	case <-c.done:
	case <-time.After(5 * time.Second):
		t.Fatalf("prefetch of %s never completed", fid)
	}
	if !c.hasCompletedError() {
		t.Fatalf("prefetch of %s was expected to fail", fid)
	}
}

func TestReaderCacheRetriesAFailedFetch(t *testing.T) {
	initHTTP()
	want := payload(300)
	srv, hits := chunkServer(t, want, 1)

	lookup := func(ctx context.Context, fileId string) ([]string, error) {
		return []string{srv.URL + "/" + fileId}, nil
	}
	rc := NewReaderCache(10, newMockChunkCacheForReaderCache(), lookup)
	defer rc.destroy()

	buf := make([]byte, len(want))
	if _, err := rc.ReadChunkAt(context.Background(), buf, "79,291ffbf5acab21", nil, false, 0, len(want), false); err == nil {
		t.Fatal("first read: want the fetch failure, got nil")
	}
	n, err := rc.ReadChunkAt(context.Background(), buf, "79,291ffbf5acab21", nil, false, 0, len(want), false)
	if err != nil || n != len(want) || !bytes.Equal(buf, want) {
		t.Fatalf("second read: got n=%d err=%v, want the full chunk", n, err)
	}
	if got := hits.Load(); got != 2 {
		t.Fatalf("volume server asked %d times, want 2", got)
	}
}

// A failed download used to carry no completion stamp, and the size sweep
// evicts only completed entries -- so failures piled up past the limit and could
// never be reclaimed.
func TestFailedDownloadsDoNotHoldTheLimit(t *testing.T) {
	initHTTP()
	want := payload(64)
	srv, _ := chunkServer(t, want, 0)

	lookup := func(ctx context.Context, fileId string) ([]string, error) {
		if fileId == "ok" {
			return []string{srv.URL + "/ok"}, nil
		}
		return nil, errors.New("not found")
	}
	const limit = 2
	rc := NewReaderCache(limit, newMockChunkCacheForReaderCache(), lookup)
	defer rc.destroy()

	buf := make([]byte, len(want))
	for _, fid := range []string{"bad-1", "bad-2"} {
		if _, err := rc.ReadChunkAt(context.Background(), buf, fid, nil, false, 0, len(want), false); err == nil {
			t.Fatalf("%s: want a failure", fid)
		}
	}
	if _, err := rc.ReadChunkAt(context.Background(), buf, "ok", nil, false, 0, len(want), false); err != nil {
		t.Fatalf("ok: %v", err)
	}

	rc.Lock()
	n := len(rc.downloaders)
	rc.Unlock()
	if n > limit {
		t.Fatalf("%d downloaders held against a limit of %d: failed entries were not reclaimable", n, limit)
	}
}

func TestReaderCacheReusesASuccessfulChunk(t *testing.T) {
	initHTTP()
	want := payload(128)
	srv, hits := chunkServer(t, want, 0)

	var calls atomic.Int32
	lookup := func(ctx context.Context, fileId string) ([]string, error) {
		calls.Add(1)
		return []string{srv.URL + "/" + fileId}, nil
	}
	rc := NewReaderCache(10, newMockChunkCacheForReaderCache(), lookup)
	defer rc.destroy()

	for i := 0; i < 3; i++ {
		buf := make([]byte, len(want))
		n, err := rc.ReadChunkAt(context.Background(), buf, "18,abc", nil, false, 0, len(want), false)
		if err != nil || n != len(want) {
			t.Fatalf("read %d: n=%d err=%v", i, n, err)
		}
	}
	if calls.Load() != 1 || hits.Load() != 1 {
		t.Fatalf("lookup=%d fetch=%d, want 1 and 1: a success is data and is reused", calls.Load(), hits.Load())
	}
}

// Readers, evictions and UnCache racing on one file id. Run with -race: the
// point is the synchronization, and a pass without the detector proves little.
func TestConcurrentRetryAndEviction(t *testing.T) {
	initHTTP()
	want := payload(256)
	srv, _ := chunkServer(t, want, 0)

	var calls atomic.Int32
	lookup := func(ctx context.Context, fileId string) ([]string, error) {
		if calls.Add(1)%3 == 1 {
			return nil, errors.New("flaky lookup")
		}
		return []string{srv.URL + "/" + fileId}, nil
	}
	rc := NewReaderCache(4, newMockChunkCacheForReaderCache(), lookup)
	defer rc.destroy()

	const fid = "78,race"
	var wg sync.WaitGroup
	var ok atomic.Int32
	for g := 0; g < 16; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 40; i++ {
				buf := make([]byte, len(want))
				n, err := rc.ReadChunkAt(context.Background(), buf, fid, nil, false, 0, len(want), false)
				if err == nil {
					if n != len(want) || !bytes.Equal(buf, want) {
						t.Errorf("a successful read returned %d wrong bytes", n)
						return
					}
					ok.Add(1)
				}
				if i%7 == 0 {
					rc.UnCache(fid)
				}
			}
		}()
	}
	wg.Wait()
	if ok.Load() == 0 {
		t.Fatal("no read ever succeeded: failures are still being replayed")
	}
}
