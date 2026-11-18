// Package archive implements per-organization / per-bucket scheduled incremental
// backup, metering, retention pruning, and SHA-256 verified bucket restore (HIP-1342).
package archive

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"path"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// KeyPrefix namespaces per-bucket incremental metadata checkpoints in filer KV.
	KeyPrefix = "archive."
	// FilerBucketsRoot is the canonical filer path root holding S3 buckets.
	FilerBucketsRoot = "/buckets"
)

// Schedule defines the per-bucket incremental archive cadence (HIP-1342).
type Schedule string

const (
	ScheduleOff    Schedule = "off"
	ScheduleHourly Schedule = "hourly"
	ScheduleDaily  Schedule = "daily"
	ScheduleWeekly Schedule = "weekly"
)

// ParseSchedule parses and validates a schedule cadence string.
func ParseSchedule(raw string) (Schedule, error) {
	switch Schedule(strings.ToLower(strings.TrimSpace(raw))) {
	case ScheduleOff:
		return ScheduleOff, nil
	case ScheduleHourly:
		return ScheduleHourly, nil
	case ScheduleDaily:
		return ScheduleDaily, nil
	case ScheduleWeekly:
		return ScheduleWeekly, nil
	default:
		return "", fmt.Errorf("invalid archive schedule %q: must be off, hourly, daily, or weekly", raw)
	}
}

// Interval returns the duration between scheduled passes, or 0 when disabled.
func (s Schedule) Interval() time.Duration {
	switch s {
	case ScheduleHourly:
		return time.Hour
	case ScheduleDaily:
		return 24 * time.Hour
	case ScheduleWeekly:
		return 7 * 24 * time.Hour
	default:
		return 0
	}
}

// ValidateBucketName enforces strict single-bucket naming so a tenant policy
// cannot traverse outside /buckets/<bucket>/.
func ValidateBucketName(bucket string) error {
	if bucket == "" {
		return errors.New("bucket name must not be empty")
	}
	if bucket == "." || bucket == ".." || strings.Contains(bucket, "..") {
		return fmt.Errorf("invalid bucket name %q: path traversal is forbidden", bucket)
	}
	if strings.ContainsAny(bucket, "/\\\x00") {
		return fmt.Errorf("invalid bucket name %q: path separators and NUL bytes are forbidden", bucket)
	}
	for i := 0; i < len(bucket); i++ {
		ch := bucket[i]
		if ch < 0x20 || ch > 0x7e {
			return fmt.Errorf("invalid bucket name %q: non-printable ASCII byte at %d", bucket, i)
		}
	}
	return nil
}

// BucketPolicy configures scheduled incremental backup and retention for one bucket.
type BucketPolicy struct {
	Bucket        string
	RemoteBucket  string
	RemotePrefix  string
	Schedule      Schedule
	RetentionDays int
	PropagateDel  bool
}

// Validate checks policy invariants before scheduling or running a pass.
func (p BucketPolicy) Validate() error {
	if err := ValidateBucketName(p.Bucket); err != nil {
		return err
	}
	if err := ValidateBucketName(p.RemoteBucket); err != nil {
		return fmt.Errorf("remote bucket: %w", err)
	}
	if strings.Contains(p.RemotePrefix, "..") || strings.ContainsRune(p.RemotePrefix, '\x00') {
		return fmt.Errorf("invalid remote prefix %q", p.RemotePrefix)
	}
	if _, err := ParseSchedule(string(p.Schedule)); err != nil {
		return err
	}
	if p.RetentionDays < 0 {
		return fmt.Errorf("retentionDays must be >= 0, got %d", p.RetentionDays)
	}
	return nil
}

// FilerBucketPath returns the canonical filer prefix `/buckets/<bucket>/`.
func (p BucketPolicy) FilerBucketPath() string {
	return FilerBucketsRoot + "/" + p.Bucket + "/"
}

// SinkCheckpointID computes a deterministic 32-bit checkpoint key ID for filer KV offset storage.
func (p BucketPolicy) SinkCheckpointID() int32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(p.Bucket))
	_, _ = h.Write([]byte("->"))
	_, _ = h.Write([]byte(p.RemoteBucket))
	_, _ = h.Write([]byte(":"))
	_, _ = h.Write([]byte(strings.Trim(p.RemotePrefix, "/")))
	return int32(h.Sum32() & 0x7fffffff)
}

// BucketMetrics holds a snapshot for per-org/per-bucket billing and observability.
type BucketMetrics struct {
	Bucket             string
	ArchivedBytes      int64
	RestoredBytes      int64
	PutRequests        int64
	DeleteRequests     int64
	ListRequests       int64
	GetRequests        int64
	PassesCompleted    int64
	PassesFailed       int64
	RetainedPrunedKeys int64
	LastOffsetTsNs     int64
	LastLagNs          int64
	LastPassDurationNs int64
}

type bucketMeterCounters struct {
	archivedBytes      atomic.Int64
	restoredBytes      atomic.Int64
	putRequests        atomic.Int64
	deleteRequests     atomic.Int64
	listRequests       atomic.Int64
	getRequests        atomic.Int64
	passesCompleted    atomic.Int64
	passesFailed       atomic.Int64
	retainedPrunedKeys atomic.Int64
	lastOffsetTsNs     atomic.Int64
	lastLagNs          atomic.Int64
	lastPassDurationNs atomic.Int64
}

func (c *bucketMeterCounters) snapshot(bucket string) BucketMetrics {
	return BucketMetrics{
		Bucket:             bucket,
		ArchivedBytes:      c.archivedBytes.Load(),
		RestoredBytes:      c.restoredBytes.Load(),
		PutRequests:        c.putRequests.Load(),
		DeleteRequests:     c.deleteRequests.Load(),
		ListRequests:       c.listRequests.Load(),
		GetRequests:        c.getRequests.Load(),
		PassesCompleted:    c.passesCompleted.Load(),
		PassesFailed:       c.passesFailed.Load(),
		RetainedPrunedKeys: c.retainedPrunedKeys.Load(),
		LastOffsetTsNs:     c.lastOffsetTsNs.Load(),
		LastLagNs:          c.lastLagNs.Load(),
		LastPassDurationNs: c.lastPassDurationNs.Load(),
	}
}

// MutationEvent represents one ordered filer metadata mutation inside a bucket.
type MutationEvent struct {
	TsNs      int64
	Key       string // key relative to /buckets/<bucket>/
	Deleted   bool
	MtimeSec  int64
	Data      []byte
	ContentMD string
}

// ManifestRecord is stored alongside incremental partitions to allow exact
// point-in-time / latest-state bucket reconstruction after hot-volume failure.
type ManifestRecord struct {
	Key       string `json:"key"`
	TsNs      int64  `json:"ts_ns"`
	Deleted   bool   `json:"deleted"`
	Size      int64  `json:"size"`
	Sha256Hex string `json:"sha256_hex"`
	ObjectKey string `json:"object_key"`
	Partition string `json:"partition"`
}

// RemoteObjectStore abstracts the remote S3 bucket used by incremental archive passes and restores.
type RemoteObjectStore interface {
	PutObject(ctx context.Context, bucket, key string, body []byte) error
	GetObject(ctx context.Context, bucket, key string) ([]byte, error)
	DeleteObject(ctx context.Context, bucket, key string) error
	ListObjects(ctx context.Context, bucket, prefix string) ([]string, error)
}

// HotBucketStore abstracts a hot filer/volume bucket data plane for incremental scanning and restore.
type HotBucketStore interface {
	ReadEvents(ctx context.Context, bucket string, afterTsNs, upToTsNs int64) ([]MutationEvent, error)
	LatestEventTsNs(ctx context.Context, bucket string) (int64, error)
	WriteRestoredObject(ctx context.Context, bucket, key string, body []byte) error
	DeleteRestoredObject(ctx context.Context, bucket, key string) error
}

// OffsetStore persists per-bucket metadata offsets across scheduler restarts.
type OffsetStore interface {
	GetOffset(ctx context.Context, checkpointID int32) (int64, error)
	SetOffset(ctx context.Context, checkpointID int32, tsNs int64) error
}

// Scheduler coordinates per-bucket scheduled incremental passes, metering, retention, and restore.
type Scheduler struct {
	mu       sync.RWMutex
	policies map[string]BucketPolicy
	nextDue  map[string]time.Time
	meters   map[string]*bucketMeterCounters

	hot     HotBucketStore
	remote  RemoteObjectStore
	offsets OffsetStore
}

// NewScheduler constructs a scheduler bound to hot, remote, and offset stores.
func NewScheduler(hot HotBucketStore, remote RemoteObjectStore, offsets OffsetStore) *Scheduler {
	return &Scheduler{
		policies: make(map[string]BucketPolicy),
		nextDue:  make(map[string]time.Time),
		meters:   make(map[string]*bucketMeterCounters),
		hot:      hot,
		remote:   remote,
		offsets:  offsets,
	}
}

// UpsertPolicy registers or updates a bucket's archive schedule and retention policy.
func (s *Scheduler) UpsertPolicy(policy BucketPolicy, now time.Time) error {
	if err := policy.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.policies[policy.Bucket] = policy
	if _, ok := s.meters[policy.Bucket]; !ok {
		s.meters[policy.Bucket] = &bucketMeterCounters{}
	}
	interval := policy.Schedule.Interval()
	if interval > 0 {
		s.nextDue[policy.Bucket] = now.Add(interval)
	} else {
		delete(s.nextDue, policy.Bucket)
	}
	return nil
}

// Metrics returns a point-in-time copy of metered counters for `bucket`.
func (s *Scheduler) Metrics(bucket string) (BucketMetrics, bool) {
	s.mu.RLock()
	meter, ok := s.meters[bucket]
	s.mu.RUnlock()
	if !ok {
		return BucketMetrics{}, false
	}
	return meter.snapshot(bucket), true
}

// RunDue executes incremental passes (plus retention pruning) for all buckets whose schedule is due at `now`.
func (s *Scheduler) RunDue(ctx context.Context, now time.Time) (int, error) {
	s.mu.RLock()
	var dueBuckets []BucketPolicy
	for _, pol := range s.policies {
		if pol.Schedule == ScheduleOff {
			continue
		}
		dueAt, ok := s.nextDue[pol.Bucket]
		if ok && !now.Before(dueAt) {
			dueBuckets = append(dueBuckets, pol)
		}
	}
	s.mu.RUnlock()

	sort.Slice(dueBuckets, func(i, j int) bool {
		return dueBuckets[i].Bucket < dueBuckets[j].Bucket
	})

	ran := 0
	var firstErr error
	for _, pol := range dueBuckets {
		if err := s.RunBucketPass(ctx, pol.Bucket, now); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if pol.RetentionDays > 0 {
			if _, err := s.PruneRetention(ctx, pol.Bucket, now); err != nil && firstErr == nil {
				firstErr = err
			}
		}
		s.mu.Lock()
		if cur, ok := s.policies[pol.Bucket]; ok && cur.Schedule != ScheduleOff {
			s.nextDue[pol.Bucket] = now.Add(cur.Schedule.Interval())
		}
		s.mu.Unlock()
		ran++
	}
	return ran, firstErr
}

func (s *Scheduler) remoteKeyPrefix(pol BucketPolicy) string {
	trimmed := strings.Trim(pol.RemotePrefix, "/")
	if trimmed == "" {
		return pol.Bucket
	}
	return trimmed + "/" + pol.Bucket
}

// RunBucketPass executes a single bounded incremental pass from the stored metadata offset up to `now`.
func (s *Scheduler) RunBucketPass(ctx context.Context, bucket string, now time.Time) error {
	s.mu.RLock()
	pol, ok := s.policies[bucket]
	meter := s.meters[bucket]
	s.mu.RUnlock()
	if !ok {
		return fmt.Errorf("no archive policy configured for bucket %q", bucket)
	}

	passStart := time.Now()
	checkpointID := pol.SinkCheckpointID()
	lastOffset, err := s.offsets.GetOffset(ctx, checkpointID)
	if err != nil {
		meter.passesFailed.Add(1)
		return fmt.Errorf("read offset for bucket %q: %w", bucket, err)
	}

	passBoundTsNs := now.UnixNano()
	if passBoundTsNs <= lastOffset {
		meter.passesCompleted.Add(1)
		meter.lastPassDurationNs.Store(time.Since(passStart).Nanoseconds())
		return nil
	}

	events, err := s.hot.ReadEvents(ctx, bucket, lastOffset, passBoundTsNs)
	if err != nil {
		meter.passesFailed.Add(1)
		return fmt.Errorf("read metadata events for bucket %q: %w", bucket, err)
	}

	sort.SliceStable(events, func(i, j int) bool {
		return events[i].TsNs < events[j].TsNs
	})

	basePrefix := s.remoteKeyPrefix(pol)
	maxEventTsNs := lastOffset

	for idx, ev := range events {
		cleanRel, err := SanitizeObjectRelativeKey(ev.Key)
		if err != nil {
			meter.passesFailed.Add(1)
			return fmt.Errorf("unsafe object key in bucket %q: %w", bucket, err)
		}

		partDate := time.Unix(0, ev.TsNs).UTC().Format("2006-01-02")
		if ev.MtimeSec > 0 {
			partDate = time.Unix(ev.MtimeSec, 0).UTC().Format("2006-01-02")
		}

		var rec ManifestRecord
		rec.Key = cleanRel
		rec.TsNs = ev.TsNs
		rec.Deleted = ev.Deleted
		rec.Partition = partDate

		if !ev.Deleted {
			sum := sha256.Sum256(ev.Data)
			rec.Sha256Hex = hex.EncodeToString(sum[:])
			rec.Size = int64(len(ev.Data))
			dataObjKey := fmt.Sprintf("%s/data/%s/%020d_%06d/%s", basePrefix, partDate, ev.TsNs, idx, cleanRel)
			if err := s.remote.PutObject(ctx, pol.RemoteBucket, dataObjKey, ev.Data); err != nil {
				meter.passesFailed.Add(1)
				return fmt.Errorf("upload archived object %q: %w", dataObjKey, err)
			}
			meter.putRequests.Add(1)
			meter.archivedBytes.Add(int64(len(ev.Data)))
			rec.ObjectKey = dataObjKey

			// Maintain live mirror object key for direct remote.mount access and post-retention fallback.
			liveObjKey := fmt.Sprintf("%s/current/%s", basePrefix, cleanRel)
			if err := s.remote.PutObject(ctx, pol.RemoteBucket, liveObjKey, ev.Data); err != nil {
				meter.passesFailed.Add(1)
				return fmt.Errorf("upload current mirror object %q: %w", liveObjKey, err)
			}
			meter.putRequests.Add(1)
			meter.archivedBytes.Add(int64(len(ev.Data)))
		} else if pol.PropagateDel {
			liveObjKey := fmt.Sprintf("%s/current/%s", basePrefix, cleanRel)
			if err := s.remote.DeleteObject(ctx, pol.RemoteBucket, liveObjKey); err != nil {
				meter.passesFailed.Add(1)
				return fmt.Errorf("delete mirrored object %q: %w", liveObjKey, err)
			}
			meter.deleteRequests.Add(1)
		}

		recBytes, err := json.Marshal(rec)
		if err != nil {
			meter.passesFailed.Add(1)
			return err
		}
		metaKey := fmt.Sprintf("%s/log/%s/%020d_%06d.json", basePrefix, partDate, ev.TsNs, idx)
		if err := s.remote.PutObject(ctx, pol.RemoteBucket, metaKey, recBytes); err != nil {
			meter.passesFailed.Add(1)
			return fmt.Errorf("upload manifest log %q: %w", metaKey, err)
		}
		meter.putRequests.Add(1)
		meter.archivedBytes.Add(int64(len(recBytes)))

		if ev.TsNs > maxEventTsNs {
			maxEventTsNs = ev.TsNs
		}
	}

	if len(events) > 0 && maxEventTsNs > lastOffset {
		if err := s.offsets.SetOffset(ctx, checkpointID, maxEventTsNs); err != nil {
			meter.passesFailed.Add(1)
			return fmt.Errorf("persist archive offset for bucket %q: %w", bucket, err)
		}
	}
	newOffset := maxEventTsNs

	latestHotTsNs, err := s.hot.LatestEventTsNs(ctx, bucket)
	if err == nil && latestHotTsNs > newOffset {
		meter.lastLagNs.Store(latestHotTsNs - newOffset)
	} else {
		meter.lastLagNs.Store(0)
	}
	meter.lastOffsetTsNs.Store(newOffset)
	meter.passesCompleted.Add(1)
	meter.lastPassDurationNs.Store(time.Since(passStart).Nanoseconds())
	return nil
}

// PruneRetention deletes date-partitioned incremental objects (`data/YYYY-MM-DD/` and `log/YYYY-MM-DD/`)
// strictly older than `now - RetentionDays*24h`, keeping `current/` untouched so live restore always works.
func (s *Scheduler) PruneRetention(ctx context.Context, bucket string, now time.Time) (int, error) {
	s.mu.RLock()
	pol, ok := s.policies[bucket]
	meter := s.meters[bucket]
	s.mu.RUnlock()
	if !ok {
		return 0, fmt.Errorf("no archive policy configured for bucket %q", bucket)
	}
	if pol.RetentionDays <= 0 {
		return 0, nil
	}

	cutoffDate := now.UTC().Add(-time.Duration(pol.RetentionDays) * 24 * time.Hour).Format("2006-01-02")
	basePrefix := s.remoteKeyPrefix(pol)
	pruned := 0

	for _, kind := range []string{"data", "log"} {
		prefix := fmt.Sprintf("%s/%s/", basePrefix, kind)
		keys, err := s.remote.ListObjects(ctx, pol.RemoteBucket, prefix)
		meter.listRequests.Add(1)
		if err != nil {
			return pruned, fmt.Errorf("list %s objects for retention: %w", kind, err)
		}
		for _, k := range keys {
			rel := strings.TrimPrefix(k, prefix)
			parts := strings.SplitN(rel, "/", 2)
			if len(parts) < 2 || len(parts[0]) != 10 {
				continue
			}
			partDate := parts[0]
			if _, parseErr := time.Parse("2006-01-02", partDate); parseErr != nil {
				continue
			}
			if partDate < cutoffDate {
				if err := s.remote.DeleteObject(ctx, pol.RemoteBucket, k); err != nil {
					return pruned, fmt.Errorf("delete expired archive object %q: %w", k, err)
				}
				meter.deleteRequests.Add(1)
				meter.retainedPrunedKeys.Add(1)
				pruned++
			}
		}
	}
	return pruned, nil
}

// RestoreBucket reconstructs a bucket into the hot store from its remote S3 archive,
// verifying SHA-256 digests when restoring via manifest logs and falling back to `current/`
// when older incremental logs were pruned by retention.
func (s *Scheduler) RestoreBucket(ctx context.Context, bucket string) (int, error) {
	s.mu.RLock()
	pol, ok := s.policies[bucket]
	meter := s.meters[bucket]
	s.mu.RUnlock()
	if !ok {
		return 0, fmt.Errorf("no archive policy configured for bucket %q", bucket)
	}

	basePrefix := s.remoteKeyPrefix(pol)
	logPrefix := basePrefix + "/log/"
	logKeys, err := s.remote.ListObjects(ctx, pol.RemoteBucket, logPrefix)
	meter.listRequests.Add(1)
	if err != nil {
		return 0, fmt.Errorf("list manifest logs for restore: %w", err)
	}
	sort.Strings(logKeys)

	latestByKey := make(map[string]ManifestRecord)
	for _, lk := range logKeys {
		raw, err := s.remote.GetObject(ctx, pol.RemoteBucket, lk)
		meter.getRequests.Add(1)
		if err != nil {
			return 0, fmt.Errorf("read manifest %q: %w", lk, err)
		}
		var rec ManifestRecord
		if err := json.Unmarshal(raw, &rec); err != nil {
			return 0, fmt.Errorf("decode manifest %q: %w", lk, err)
		}
		prev, exists := latestByKey[rec.Key]
		if !exists || rec.TsNs >= prev.TsNs {
			latestByKey[rec.Key] = rec
		}
	}

	currentPrefix := basePrefix + "/current/"
	currentKeys, err := s.remote.ListObjects(ctx, pol.RemoteBucket, currentPrefix)
	meter.listRequests.Add(1)
	if err != nil {
		return 0, fmt.Errorf("list current mirror for restore: %w", err)
	}

	restoredCount := 0
	seenLive := make(map[string]bool)

	sortedRel := make([]string, 0, len(latestByKey))
	for rel := range latestByKey {
		sortedRel = append(sortedRel, rel)
	}
	sort.Strings(sortedRel)

	for _, rel := range sortedRel {
		rec := latestByKey[rel]
		if rec.Deleted {
			if err := s.hot.DeleteRestoredObject(ctx, bucket, rel); err != nil {
				return restoredCount, err
			}
			continue
		}
		payload, err := s.remote.GetObject(ctx, pol.RemoteBucket, rec.ObjectKey)
		meter.getRequests.Add(1)
		if err != nil {
			fallbackKey := currentPrefix + rel
			payload, err = s.remote.GetObject(ctx, pol.RemoteBucket, fallbackKey)
			meter.getRequests.Add(1)
			if err != nil {
				return restoredCount, fmt.Errorf("fetch archived payload for %q: %w", rel, err)
			}
		}
		sum := sha256.Sum256(payload)
		gotHex := hex.EncodeToString(sum[:])
		if rec.Sha256Hex != "" && gotHex != rec.Sha256Hex {
			return restoredCount, fmt.Errorf("archive integrity mismatch for %q: got sha256 %s, want %s", rel, gotHex, rec.Sha256Hex)
		}
		if err := s.hot.WriteRestoredObject(ctx, bucket, rel, payload); err != nil {
			return restoredCount, fmt.Errorf("write restored object %q: %w", rel, err)
		}
		meter.restoredBytes.Add(int64(len(payload)))
		seenLive[rel] = true
		restoredCount++
	}

	sort.Strings(currentKeys)
	for _, ck := range currentKeys {
		rel := strings.TrimPrefix(ck, currentPrefix)
		cleanRel, err := SanitizeObjectRelativeKey(rel)
		if err != nil || seenLive[cleanRel] {
			continue
		}
		if rec, inLog := latestByKey[cleanRel]; inLog && rec.Deleted {
			continue
		}
		payload, err := s.remote.GetObject(ctx, pol.RemoteBucket, ck)
		meter.getRequests.Add(1)
		if err != nil {
			return restoredCount, fmt.Errorf("fetch current mirror %q: %w", ck, err)
		}
		if err := s.hot.WriteRestoredObject(ctx, bucket, cleanRel, payload); err != nil {
			return restoredCount, fmt.Errorf("write restored mirror object %q: %w", cleanRel, err)
		}
		meter.restoredBytes.Add(int64(len(payload)))
		restoredCount++
	}

	return restoredCount, nil
}

// SanitizeObjectRelativeKey rejects empty, NUL, backslash, or traversal path segments inside a bucket.
func SanitizeObjectRelativeKey(raw string) (string, error) {
	trimmed := strings.TrimPrefix(raw, "/")
	if trimmed == "" {
		return "", errors.New("empty object key")
	}
	if strings.ContainsRune(trimmed, '\x00') || strings.Contains(trimmed, "\\") {
		return "", fmt.Errorf("invalid characters in key %q", raw)
	}
	for _, seg := range strings.Split(trimmed, "/") {
		if seg == "" || seg == "." || seg == ".." {
			return "", fmt.Errorf("disallowed path segment %q in key %q", seg, raw)
		}
	}
	cleaned := path.Clean(trimmed)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", fmt.Errorf("escaped path %q", raw)
	}
	return cleaned, nil
}
