package archive

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type memRemoteS3 struct {
	mu          sync.RWMutex
	buckets     map[string]map[string][]byte
	failNextPut atomic.Int32
}

func newMemRemoteS3() *memRemoteS3 {
	return &memRemoteS3{
		buckets: make(map[string]map[string][]byte),
	}
}

func (m *memRemoteS3) PutObject(_ context.Context, bucket, key string, body []byte) error {
	if rem := m.failNextPut.Load(); rem > 0 {
		m.failNextPut.Add(-1)
		return errors.New("simulated remote S3 503 SlowDown")
	}
	cp := make([]byte, len(body))
	copy(cp, body)
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.buckets[bucket]
	if !ok {
		b = make(map[string][]byte)
		m.buckets[bucket] = b
	}
	b[key] = cp
	return nil
}

func (m *memRemoteS3) GetObject(_ context.Context, bucket, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, ok := m.buckets[bucket]
	if !ok {
		return nil, fmt.Errorf("bucket %q not found", bucket)
	}
	val, ok := b[key]
	if !ok {
		return nil, fmt.Errorf("key %q not found in %q", key, bucket)
	}
	cp := make([]byte, len(val))
	copy(cp, val)
	return cp, nil
}

func (m *memRemoteS3) DeleteObject(_ context.Context, bucket, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if b, ok := m.buckets[bucket]; ok {
		delete(b, key)
	}
	return nil
}

func (m *memRemoteS3) ListObjects(_ context.Context, bucket, prefix string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, ok := m.buckets[bucket]
	if !ok {
		return nil, nil
	}
	var out []string
	for k := range b {
		if strings.HasPrefix(k, prefix) {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out, nil
}

type simHotCluster struct {
	mu      sync.RWMutex
	clockNs atomic.Int64
	objects map[string]map[string][]byte
	events  map[string][]MutationEvent
}

func newSimHotCluster(start time.Time) *simHotCluster {
	c := &simHotCluster{
		objects: make(map[string]map[string][]byte),
		events:  make(map[string][]MutationEvent),
	}
	c.clockNs.Store(start.UnixNano())
	return c
}

func (c *simHotCluster) Put(bucket, key string, data []byte, mtime time.Time) int64 {
	cp := make([]byte, len(data))
	copy(cp, data)

	c.mu.Lock()
	defer c.mu.Unlock()
	ts := c.clockNs.Add(1_000)
	if _, ok := c.objects[bucket]; !ok {
		c.objects[bucket] = make(map[string][]byte)
	}
	c.objects[bucket][key] = cp
	c.events[bucket] = append(c.events[bucket], MutationEvent{
		TsNs:     ts,
		Key:      key,
		Deleted:  false,
		MtimeSec: mtime.Unix(),
		Data:     cp,
	})
	return ts
}

func (c *simHotCluster) Delete(bucket, key string, mtime time.Time) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	ts := c.clockNs.Add(1_000)
	if b, ok := c.objects[bucket]; ok {
		delete(b, key)
	}
	c.events[bucket] = append(c.events[bucket], MutationEvent{
		TsNs:     ts,
		Key:      key,
		Deleted:  true,
		MtimeSec: mtime.Unix(),
	})
	return ts
}

func (c *simHotCluster) KillHotVolume(bucket string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.objects, bucket)
	delete(c.events, bucket)
}

func (c *simHotCluster) SnapshotLive(bucket string) map[string][]byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[string][]byte)
	for k, v := range c.objects[bucket] {
		cp := make([]byte, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out
}

func (c *simHotCluster) ReadEvents(_ context.Context, bucket string, afterTsNs, upToTsNs int64) ([]MutationEvent, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var matched []MutationEvent
	for _, ev := range c.events[bucket] {
		if ev.TsNs > afterTsNs && ev.TsNs <= upToTsNs {
			matched = append(matched, ev)
		}
	}
	return matched, nil
}

func (c *simHotCluster) LatestEventTsNs(_ context.Context, bucket string) (int64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	evs := c.events[bucket]
	if len(evs) == 0 {
		return 0, nil
	}
	return evs[len(evs)-1].TsNs, nil
}

func (c *simHotCluster) WriteRestoredObject(_ context.Context, bucket, key string, body []byte) error {
	cp := make([]byte, len(body))
	copy(cp, body)
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.objects[bucket]; !ok {
		c.objects[bucket] = make(map[string][]byte)
	}
	c.objects[bucket][key] = cp
	return nil
}

func (c *simHotCluster) DeleteRestoredObject(_ context.Context, bucket, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if b, ok := c.objects[bucket]; ok {
		delete(b, key)
	}
	return nil
}

type memOffsetStore struct {
	mu      sync.Mutex
	offsets map[int32]int64
}

func newMemOffsetStore() *memOffsetStore {
	return &memOffsetStore{offsets: make(map[int32]int64)}
}

func (s *memOffsetStore) GetOffset(_ context.Context, id int32) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.offsets[id], nil
}

func (s *memOffsetStore) SetOffset(_ context.Context, id int32, tsNs int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.offsets[id] = tsNs
	return nil
}

func TestArchiveScheduleParsingAndIntervals(t *testing.T) {
	cases := []struct {
		in   string
		want Schedule
		dur  time.Duration
	}{
		{"off", ScheduleOff, 0},
		{"HOURLY", ScheduleHourly, time.Hour},
		{"daily", ScheduleDaily, 24 * time.Hour},
		{" weekly ", ScheduleWeekly, 7 * 24 * time.Hour},
	}
	for _, tc := range cases {
		got, err := ParseSchedule(tc.in)
		if err != nil {
			t.Fatalf("ParseSchedule(%q) unexpected error: %v", tc.in, err)
		}
		if got != tc.want || got.Interval() != tc.dur {
			t.Fatalf("ParseSchedule(%q) = (%q, %v), want (%q, %v)", tc.in, got, got.Interval(), tc.want, tc.dur)
		}
	}
	if _, err := ParseSchedule("every-5-seconds"); err == nil {
		t.Fatal("expected error for unsupported schedule")
	}
}

func TestBucketPolicyRejectsTraversalAndInvalidNames(t *testing.T) {
	badBuckets := []string{"", "../etc", "org/sub", "org\\sub", "org\x00bad", "."}
	for _, b := range badBuckets {
		pol := BucketPolicy{
			Bucket:       b,
			RemoteBucket: "archive-vault",
			Schedule:     ScheduleDaily,
		}
		if err := pol.Validate(); err == nil {
			t.Fatalf("expected policy with bucket %q to fail validation", b)
		}
	}
}

func TestScheduleCadenceAndPerOrgTenantMetering(t *testing.T) {
	ctx := context.Background()
	t0 := time.Date(2026, 11, 18, 0, 0, 0, 0, time.UTC)
	hot := newSimHotCluster(t0)
	remote := newMemRemoteS3()
	offsets := newMemOffsetStore()
	sched := NewScheduler(hot, remote, offsets)

	if err := sched.UpsertPolicy(BucketPolicy{
		Bucket:       "org-alpha",
		RemoteBucket: "cold-s3",
		RemotePrefix: "archives",
		Schedule:     ScheduleHourly,
		PropagateDel: true,
	}, t0); err != nil {
		t.Fatal(err)
	}
	if err := sched.UpsertPolicy(BucketPolicy{
		Bucket:       "org-beta",
		RemoteBucket: "cold-s3",
		RemotePrefix: "archives",
		Schedule:     ScheduleDaily,
		PropagateDel: true,
	}, t0); err != nil {
		t.Fatal(err)
	}
	if err := sched.UpsertPolicy(BucketPolicy{
		Bucket:       "org-disabled",
		RemoteBucket: "cold-s3",
		Schedule:     ScheduleOff,
	}, t0); err != nil {
		t.Fatal(err)
	}

	hot.Put("org-alpha", "models/v1.bin", bytes.Repeat([]byte("A"), 4096), t0)
	hot.Put("org-beta", "models/v1.bin", bytes.Repeat([]byte("B"), 8192), t0)
	hot.Put("org-disabled", "models/v1.bin", bytes.Repeat([]byte("D"), 2048), t0)

	ran, err := sched.RunDue(ctx, t0.Add(30*time.Minute))
	if err != nil || ran != 0 {
		t.Fatalf("at +30m: ran=%d err=%v, want 0", ran, err)
	}

	ran, err = sched.RunDue(ctx, t0.Add(61*time.Minute))
	if err != nil || ran != 1 {
		t.Fatalf("at +61m: ran=%d err=%v, want 1", ran, err)
	}

	alphaM, _ := sched.Metrics("org-alpha")
	betaM, _ := sched.Metrics("org-beta")
	disM, _ := sched.Metrics("org-disabled")
	if alphaM.PassesCompleted != 1 || alphaM.ArchivedBytes < 8192 || alphaM.PutRequests != 3 {
		t.Fatalf("unexpected org-alpha metrics after pass 1: %+v", alphaM)
	}
	if betaM.PassesCompleted != 0 || betaM.PutRequests != 0 {
		t.Fatalf("org-beta ran prematurely: %+v", betaM)
	}
	if disM.PassesCompleted != 0 || disM.PutRequests != 0 {
		t.Fatalf("org-disabled should never run: %+v", disM)
	}

	prevBytes := alphaM.ArchivedBytes
	prevPuts := alphaM.PutRequests
	ran, err = sched.RunDue(ctx, t0.Add(122*time.Minute))
	if err != nil || ran != 1 {
		t.Fatalf("at +122m: ran=%d err=%v, want 1", ran, err)
	}
	alphaM2, _ := sched.Metrics("org-alpha")
	if alphaM2.ArchivedBytes != prevBytes || alphaM2.PutRequests != prevPuts {
		t.Fatalf("idle incremental pass re-uploaded data: before=%+v after=%+v", alphaM, alphaM2)
	}
}

func TestCrashSafetyDoesNotAdvanceOffsetOnRemoteFailure(t *testing.T) {
	ctx := context.Background()
	t0 := time.Date(2026, 11, 18, 1, 0, 0, 0, time.UTC)
	hot := newSimHotCluster(t0)
	remote := newMemRemoteS3()
	offsets := newMemOffsetStore()
	sched := NewScheduler(hot, remote, offsets)

	pol := BucketPolicy{
		Bucket:       "org-resilient",
		RemoteBucket: "cold-s3",
		Schedule:     ScheduleHourly,
		PropagateDel: true,
	}
	if err := sched.UpsertPolicy(pol, t0); err != nil {
		t.Fatal(err)
	}

	hot.Put("org-resilient", "ckpt/step-1.bin", []byte("step-1-payload"), t0)
	remote.failNextPut.Store(1)

	if err := sched.RunBucketPass(ctx, "org-resilient", t0.Add(time.Minute)); err == nil {
		t.Fatal("expected pass failure when remote S3 returns 503")
	}
	storedOff, _ := offsets.GetOffset(ctx, pol.SinkCheckpointID())
	if storedOff != 0 {
		t.Fatalf("failed pass must not advance stored offset, got %d", storedOff)
	}

	if err := sched.RunBucketPass(ctx, "org-resilient", t0.Add(2*time.Minute)); err != nil {
		t.Fatalf("retry pass failed: %v", err)
	}
	hot.KillHotVolume("org-resilient")
	n, err := sched.RestoreBucket(ctx, "org-resilient")
	if err != nil || n != 1 {
		t.Fatalf("restore after retry failed: n=%d err=%v", n, err)
	}
	restored := hot.SnapshotLive("org-resilient")
	if string(restored["ckpt/step-1.bin"]) != "step-1-payload" {
		t.Fatalf("restored payload mismatch: %q", string(restored["ckpt/step-1.bin"]))
	}
}

func TestRetentionPrunesExpiredPartitionsWhilePreservingRestore(t *testing.T) {
	ctx := context.Background()
	day1 := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	hot := newSimHotCluster(day1)
	remote := newMemRemoteS3()
	offsets := newMemOffsetStore()
	sched := NewScheduler(hot, remote, offsets)

	if err := sched.UpsertPolicy(BucketPolicy{
		Bucket:        "org-ret",
		RemoteBucket:  "cold-s3",
		Schedule:      ScheduleDaily,
		RetentionDays: 7,
		PropagateDel:  true,
	}, day1); err != nil {
		t.Fatal(err)
	}

	hot.Put("org-ret", "long-lived.txt", []byte("immortal-bytes"), day1)
	if err := sched.RunBucketPass(ctx, "org-ret", day1.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}

	day12 := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	hot.clockNs.Store(day12.UnixNano())
	hot.Put("org-ret", "recent.txt", []byte("fresh-bytes"), day12)
	if err := sched.RunBucketPass(ctx, "org-ret", day12.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}

	day15 := time.Date(2026, 11, 15, 12, 0, 0, 0, time.UTC)
	pruned, err := sched.PruneRetention(ctx, "org-ret", day15)
	if err != nil {
		t.Fatalf("PruneRetention failed: %v", err)
	}
	if pruned < 2 {
		t.Fatalf("expected expired 2026-11-01 data+log keys to be pruned, got %d", pruned)
	}

	hot.KillHotVolume("org-ret")
	n, err := sched.RestoreBucket(ctx, "org-ret")
	if err != nil || n != 2 {
		t.Fatalf("RestoreBucket after retention prune: n=%d err=%v", n, err)
	}
	got := hot.SnapshotLive("org-ret")
	if string(got["long-lived.txt"]) != "immortal-bytes" || string(got["recent.txt"]) != "fresh-bytes" {
		t.Fatalf("unexpected restored contents after retention prune: %v", got)
	}
}

func TestPerOrgArchiveHotVolumeKillAndByteIdenticalRestoreWithConcurrentWorkload(t *testing.T) {
	ctx := context.Background()
	t0 := time.Date(2026, 11, 18, 2, 0, 0, 0, time.UTC)

	const numObjects = 32
	const objSize = 256 * 1024 // 256 KiB
	payloads := make([][]byte, numObjects)
	for i := 0; i < numObjects; i++ {
		buf := make([]byte, objSize)
		for j := 0; j < len(buf); j++ {
			buf[j] = byte((i*31 + j*17 + 13) & 0xff)
		}
		payloads[i] = buf
	}

	nvmeDir := t.TempDir()
	nvmeStart := time.Now()
	for i := 0; i < numObjects; i++ {
		fp := filepath.Join(nvmeDir, fmt.Sprintf("obj-%03d.bin", i))
		f, err := os.Create(fp)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.Write(payloads[i]); err != nil {
			_ = f.Close()
			t.Fatal(err)
		}
		if err := f.Sync(); err != nil {
			_ = f.Close()
			t.Fatal(err)
		}
		_ = f.Close()
	}
	nvmeElapsed := time.Since(nvmeStart)
	totalMiB := float64(numObjects*objSize) / (1024.0 * 1024.0)
	nvmeMiBs := totalMiB / nvmeElapsed.Seconds()

	hot := newSimHotCluster(t0)
	remote := newMemRemoteS3()
	offsets := newMemOffsetStore()
	sched := NewScheduler(hot, remote, offsets)

	if err := sched.UpsertPolicy(BucketPolicy{
		Bucket:        "org-prod",
		RemoteBucket:  "hanzo-archive-s3",
		RemotePrefix:  "hip1342",
		Schedule:      ScheduleHourly,
		RetentionDays: 30,
		PropagateDel:  true,
	}, t0); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	writeStart := time.Now()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < numObjects; i++ {
			key := fmt.Sprintf("datasets/shard-%03d.bin", i)
			hot.Put("org-prod", key, payloads[i], t0.Add(time.Duration(i)*time.Second))
			if i == 8 {
				hot.Delete("org-prod", "datasets/shard-001.bin", t0.Add(9*time.Second))
			}
			if i == 12 {
				updated := bytes.Repeat([]byte{0xAB}, objSize)
				hot.Put("org-prod", "datasets/shard-002.bin", updated, t0.Add(13*time.Second))
			}
		}
	}()

	for p := 1; p <= 3; p++ {
		time.Sleep(2 * time.Millisecond)
		nowTick := time.Unix(0, hot.clockNs.Load()+500)
		if err := sched.RunBucketPass(ctx, "org-prod", nowTick); err != nil {
			t.Fatalf("concurrent pass %d failed: %v", p, err)
		}
	}
	wg.Wait()

	finalTick := time.Unix(0, hot.clockNs.Load()+1_000)
	if err := sched.RunBucketPass(ctx, "org-prod", finalTick); err != nil {
		t.Fatalf("final pass failed: %v", err)
	}
	totalArchiveElapsed := time.Since(writeStart)

	expectedLive := hot.SnapshotLive("org-prod")
	expectedHashes := make(map[string][32]byte, len(expectedLive))
	for k, v := range expectedLive {
		expectedHashes[k] = sha256.Sum256(v)
	}
	if _, deletedPresent := expectedLive["datasets/shard-001.bin"]; deletedPresent {
		t.Fatal("shard-001.bin should have been deleted on hot store")
	}

	hot.KillHotVolume("org-prod")
	if len(hot.SnapshotLive("org-prod")) != 0 {
		t.Fatal("hot volume must be empty immediately after KillHotVolume")
	}

	restoreStart := time.Now()
	restoredCount, err := sched.RestoreBucket(ctx, "org-prod")
	restoreElapsed := time.Since(restoreStart)
	if err != nil {
		t.Fatalf("RestoreBucket failed: %v", err)
	}
	if restoredCount != len(expectedLive) {
		t.Fatalf("restored object count = %d, want %d", restoredCount, len(expectedLive))
	}

	restoredLive := hot.SnapshotLive("org-prod")
	for k, wantSum := range expectedHashes {
		gotBytes, ok := restoredLive[k]
		if !ok {
			t.Fatalf("missing restored key %q", k)
		}
		gotSum := sha256.Sum256(gotBytes)
		if gotSum != wantSum {
			t.Fatalf("byte-identical SHA-256 mismatch for %q", k)
		}
	}
	if _, ghost := restoredLive["datasets/shard-001.bin"]; ghost {
		t.Fatal("deleted object shard-001.bin must not reappear after restore")
	}

	m, _ := sched.Metrics("org-prod")
	archiveMiBs := (float64(m.ArchivedBytes) / (1024.0 * 1024.0)) / totalArchiveElapsed.Seconds()
	restoreMiBs := (float64(m.RestoredBytes) / (1024.0 * 1024.0)) / restoreElapsed.Seconds()
	t.Logf("BENCHMARK_SUMMARY plain_nvme_fsync_mibs=%.2f archive_pass_mibs=%.2f restore_mibs=%.2f archived_bytes=%d restored_bytes=%d puts=%d deletes=%d gets=%d lists=%d lag_ns=%d",
		nvmeMiBs, archiveMiBs, restoreMiBs, m.ArchivedBytes, m.RestoredBytes, m.PutRequests, m.DeleteRequests, m.GetRequests, m.ListRequests, m.LastLagNs)
}
