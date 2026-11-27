package discogc

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const (
	segMagic         = uint32(0x44474331) // "DGC1"
	needleMagic      = uint32(0x4E444C31) // "NDL1"
	segHeaderSize    = int64(64)
	needleHeaderSize = int64(32)
	smallAlign       = int64(8)
)

var (
	// ErrNotFound is returned when an object key does not exist for the tenant.
	ErrNotFound = errors.New("discogc: object not found")
	// ErrCorruptData is returned when an on-disk needle fails magic, CRC32, or SHA-256 verification.
	ErrCorruptData = errors.New("discogc: on-disk object checksum mismatch")
	// ErrClosed is returned when operating on a closed Store.
	ErrClosed = errors.New("discogc: store is closed")
)

type segment struct {
	id                 uint64
	tenantID           string
	path               string
	file               *os.File
	writeOffset        int64
	sealed             bool
	createdSeq         uint64
	lastModSeq         uint64
	liveBytes          int64
	deadUnpunchedBytes int64
	punchedBytes       int64
	objects            map[string]*ObjectMeta
	deadExtents        []Extent
	punchedExtents     []Extent
}

type tenantState struct {
	mu        sync.RWMutex
	id        string
	index     map[string]*ObjectMeta
	segments  map[uint64]*segment
	activeSeg *segment
}

// Store is a real OS file-backed, multi-tenant log-structured extent store
// supporting direct hole-punch discard and copy-compaction GC policies.
type Store struct {
	cfg     Config
	dir     string
	tempDir bool

	mu      sync.RWMutex
	closed  bool
	tenants map[string]*tenantState

	seqCounter atomic.Uint64
	segIDGen   atomic.Uint64

	// Global atomic / locked telemetry.
	statsMu                 sync.Mutex
	userBytesWritten        int64
	gcCopiedLiveBytes       int64
	metadataBytesWritten    int64
	punchHoleCalls          int64
	punchHoleRequestedBytes int64
	osBlocksReclaimedBytes  int64
	discardCycles           int64
	compactionCycles        int64
	punchHoleSupported      bool
	punchHoleProbed         bool
	fsBlockSize             int64

	peakPhysicalBytes int64
	sumPhysicalBytes  float64
	sumLiveBytes      float64
	spaceSamples      int64

	startTime   time.Time
	startProcIO ProcIOSnapshot
}

// Open creates a new file-backed Store with the supplied Config.
func Open(cfg Config) (*Store, error) {
	if cfg.Policy == "" {
		cfg.Policy = PolicyAdaptiveDiscard
	}
	if cfg.SegmentCapacityBytes <= 0 {
		cfg.SegmentCapacityBytes = 2 * 1024 * 1024
	}
	if cfg.DiscardGranularity <= 0 {
		cfg.DiscardGranularity = 4096
	}
	if cfg.MinDiscardExtent <= 0 {
		cfg.MinDiscardExtent = cfg.DiscardGranularity
	}
	if cfg.GarbageTriggerRatio <= 0 {
		cfg.GarbageTriggerRatio = 0.30
	}
	if cfg.SubblockFragFallbackRatio <= 0 {
		cfg.SubblockFragFallbackRatio = cfg.GarbageTriggerRatio
	}
	if cfg.SparseEmptyLiveThreshold <= 0 {
		cfg.SparseEmptyLiveThreshold = 0.05
	}
	if cfg.Drive.RawCapacityTB <= 0 {
		cfg.Drive = DefaultEnterpriseNVMe()
	}

	dir := cfg.Dir
	tempDir := false
	if dir == "" {
		var err error
		dir, err = os.MkdirTemp("", "discogc-store-*")
		if err != nil {
			return nil, fmt.Errorf("discogc: create temp dir: %w", err)
		}
		tempDir = true
	} else {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("discogc: mkdir %s: %w", dir, err)
		}
	}

	s := &Store{
		cfg:                cfg,
		dir:                dir,
		tempDir:            tempDir,
		tenants:            make(map[string]*tenantState),
		punchHoleSupported: true,
		fsBlockSize:        4096,
		startTime:          time.Now(),
		startProcIO:        ReadProcIO(),
	}
	return s, nil
}

// Close closes all open segment file descriptors and removes the temporary
// backing directory if one was auto-created by Open.
func (s *Store) Close() error {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil
	}
	s.closed = true
	tenants := make([]*tenantState, 0, len(s.tenants))
	for _, ts := range s.tenants {
		tenants = append(tenants, ts)
	}
	s.mu.Unlock()

	var firstErr error
	for _, ts := range tenants {
		ts.mu.Lock()
		for _, seg := range ts.segments {
			if seg.file != nil {
				if err := seg.file.Close(); err != nil && firstErr == nil {
					firstErr = err
				}
				seg.file = nil
			}
		}
		ts.mu.Unlock()
	}

	if s.tempDir {
		if err := os.RemoveAll(s.dir); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (s *Store) getOrCreateTenant(tenantID string) (*tenantState, error) {
	if tenantID == "" {
		tenantID = "default"
	}
	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return nil, ErrClosed
	}
	ts, ok := s.tenants[tenantID]
	s.mu.RUnlock()
	if ok {
		return ts, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil, ErrClosed
	}
	if ts, ok = s.tenants[tenantID]; ok {
		return ts, nil
	}
	ts = &tenantState{
		id:       tenantID,
		index:    make(map[string]*ObjectMeta),
		segments: make(map[uint64]*segment),
	}
	s.tenants[tenantID] = ts
	return ts, nil
}

func (s *Store) getTenant(tenantID string) (*tenantState, error) {
	if tenantID == "" {
		tenantID = "default"
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return nil, ErrClosed
	}
	ts := s.tenants[tenantID]
	return ts, nil
}

// Put writes or overwrites an object for tenantID and key. Overwriting an
// existing key marks the previous extent dead and may trigger policy GC when
// a segment seals or reaches its garbage threshold.
func (s *Store) Put(tenantID, key string, payload []byte) error {
	ts, err := s.getOrCreateTenant(tenantID)
	if err != nil {
		return err
	}

	digest := sha256.Sum256(payload)
	crc := crc32.ChecksumIEEE(payload)
	seq := s.seqCounter.Add(1)

	ts.mu.Lock()
	defer ts.mu.Unlock()

	seg, err := s.ensureActiveSegmentLocked(ts, int64(len(key))+int64(len(payload))+needleHeaderSize, seq)
	if err != nil {
		return err
	}

	meta, metaOverhead, err := s.appendNeedleLocked(seg, ts.id, key, payload, digest, crc, seq)
	if err != nil {
		return err
	}

	// If overwriting an existing key, mark old extent dead.
	var invalidatedSeg *segment
	if old, exists := ts.index[key]; exists {
		if oldSeg, ok := ts.segments[old.SegmentID]; ok {
			markDeadLocked(oldSeg, old, seq)
			invalidatedSeg = oldSeg
		}
	}

	ts.index[key] = meta
	seg.objects[key] = meta
	seg.liveBytes += meta.TotalLen
	seg.lastModSeq = seq

	s.statsMu.Lock()
	s.userBytesWritten += int64(len(payload))
	s.metadataBytesWritten += metaOverhead
	s.statsMu.Unlock()

	// Policy-specific incremental actions on overwrite.
	if invalidatedSeg != nil && s.cfg.Policy == PolicyBlindSubblockDiscard {
		// Blind negative control issues raw PunchHole immediately without alignment checks.
		if err := s.runBlindDiscardLocked(invalidatedSeg); err != nil {
			return err
		}
	}

	// Seal segment if full and run policy GC.
	if seg.writeOffset >= s.cfg.SegmentCapacityBytes {
		seg.sealed = true
		ts.activeSeg = nil
		if err := s.runGCLocked(ts, seq); err != nil {
			return err
		}
	} else if invalidatedSeg != nil && invalidatedSeg.sealed {
		if err := s.maybeRunSegmentGCLocked(ts, invalidatedSeg, seq); err != nil {
			return err
		}
	}

	return nil
}

// Get reads and verifies an object's payload from disk, checking needle magic,
// key match, CRC32, and SHA-256 digest.
func (s *Store) Get(tenantID, key string) ([]byte, error) {
	ts, err := s.getTenant(tenantID)
	if err != nil {
		return nil, err
	}
	if ts == nil {
		return nil, ErrNotFound
	}

	ts.mu.RLock()
	defer ts.mu.RUnlock()

	meta, ok := ts.index[key]
	if !ok {
		return nil, ErrNotFound
	}
	seg, ok := ts.segments[meta.SegmentID]
	if !ok || seg.file == nil {
		return nil, ErrNotFound
	}
	return readAndVerifyNeedle(seg.file, meta)
}

// Delete removes key for tenantID and marks its on-disk extent dead.
func (s *Store) Delete(tenantID, key string) (bool, error) {
	ts, err := s.getTenant(tenantID)
	if err != nil {
		return false, err
	}
	if ts == nil {
		return false, nil
	}

	seq := s.seqCounter.Add(1)

	ts.mu.Lock()
	defer ts.mu.Unlock()

	meta, ok := ts.index[key]
	if !ok {
		return false, nil
	}
	delete(ts.index, key)

	seg, ok := ts.segments[meta.SegmentID]
	if !ok {
		return true, nil
	}
	markDeadLocked(seg, meta, seq)

	if s.cfg.Policy == PolicyBlindSubblockDiscard {
		if err := s.runBlindDiscardLocked(seg); err != nil {
			return true, err
		}
		return true, nil
	}

	if err := s.maybeRunSegmentGCLocked(ts, seg, seq); err != nil {
		return true, err
	}
	return true, nil
}

// RunGC triggers a garbage-collection pass across all segments of tenantID
// (or all tenants if tenantID == "").
func (s *Store) RunGC(tenantID string) error {
	seq := s.seqCounter.Add(1)
	if tenantID != "" {
		ts, err := s.getTenant(tenantID)
		if err != nil || ts == nil {
			return err
		}
		ts.mu.Lock()
		defer ts.mu.Unlock()
		return s.runGCLocked(ts, seq)
	}

	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return ErrClosed
	}
	tenants := make([]*tenantState, 0, len(s.tenants))
	for _, ts := range s.tenants {
		tenants = append(tenants, ts)
	}
	s.mu.RUnlock()

	for _, ts := range tenants {
		ts.mu.Lock()
		err := s.runGCLocked(ts, seq)
		ts.mu.Unlock()
		if err != nil {
			return err
		}
	}
	return nil
}

// SealActive seals the current active segment for tenantID (or all tenants if
// empty), useful for deterministic end-of-epoch GC evaluation.
func (s *Store) SealActive(tenantID string) error {
	if tenantID != "" {
		ts, err := s.getTenant(tenantID)
		if err != nil || ts == nil {
			return err
		}
		ts.mu.Lock()
		if ts.activeSeg != nil {
			ts.activeSeg.sealed = true
			ts.activeSeg = nil
		}
		ts.mu.Unlock()
		return nil
	}

	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return ErrClosed
	}
	tenants := make([]*tenantState, 0, len(s.tenants))
	for _, ts := range s.tenants {
		tenants = append(tenants, ts)
	}
	s.mu.RUnlock()

	for _, ts := range tenants {
		ts.mu.Lock()
		if ts.activeSeg != nil {
			ts.activeSeg.sealed = true
			ts.activeSeg = nil
		}
		ts.mu.Unlock()
	}
	return nil
}

// VerifyAllLive reads every live object across all tenants from disk and
// asserts bit-exact SHA-256 integrity.
func (s *Store) VerifyAllLive() (int, error) {
	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return 0, ErrClosed
	}
	tenants := make([]*tenantState, 0, len(s.tenants))
	for _, ts := range s.tenants {
		tenants = append(tenants, ts)
	}
	s.mu.RUnlock()

	verified := 0
	for _, ts := range tenants {
		ts.mu.RLock()
		for _, meta := range ts.index {
			seg, ok := ts.segments[meta.SegmentID]
			if !ok || seg.file == nil {
				ts.mu.RUnlock()
				return verified, fmt.Errorf("discogc: missing segment %d for key %q", meta.SegmentID, meta.Key)
			}
			if _, err := readAndVerifyNeedle(seg.file, meta); err != nil {
				ts.mu.RUnlock()
				return verified, fmt.Errorf("discogc: verify key %q: %w", meta.Key, err)
			}
			verified++
		}
		ts.mu.RUnlock()
	}
	return verified, nil
}

// SequentialScanAll reads all live objects in segment-offset order for tenantID
// (or all tenants if empty), returning total payload bytes scanned and elapsed duration.
func (s *Store) SequentialScanAll(tenantID string) (int64, time.Duration, error) {
	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return 0, 0, ErrClosed
	}
	var tenants []*tenantState
	if tenantID != "" {
		if ts := s.tenants[tenantID]; ts != nil {
			tenants = append(tenants, ts)
		}
	} else {
		for _, ts := range s.tenants {
			tenants = append(tenants, ts)
		}
	}
	s.mu.RUnlock()

	start := time.Now()
	var totalBytes int64
	for _, ts := range tenants {
		ts.mu.RLock()
		segIDs := make([]uint64, 0, len(ts.segments))
		for id := range ts.segments {
			segIDs = append(segIDs, id)
		}
		sort.Slice(segIDs, func(i, j int) bool { return segIDs[i] < segIDs[j] })
		for _, id := range segIDs {
			seg := ts.segments[id]
			metas := make([]*ObjectMeta, 0, len(seg.objects))
			for _, m := range seg.objects {
				metas = append(metas, m)
			}
			sort.Slice(metas, func(i, j int) bool { return metas[i].Offset < metas[j].Offset })
			for _, m := range metas {
				data, err := readAndVerifyNeedle(seg.file, m)
				if err != nil {
					ts.mu.RUnlock()
					return totalBytes, time.Since(start), err
				}
				totalBytes += int64(len(data))
			}
		}
		ts.mu.RUnlock()
	}
	return totalBytes, time.Since(start), nil
}

// SampleSpace queries current OS physical disk usage and live logical bytes,
// updating time-averaged space amplification and peak physical disk usage.
func (s *Store) SampleSpace() (physicalBytes, liveBytes int64) {
	_, phys, live, _, _, _ := s.inspectDiskUsage()
	s.statsMu.Lock()
	if phys > s.peakPhysicalBytes {
		s.peakPhysicalBytes = phys
	}
	if live > 0 {
		s.sumPhysicalBytes += float64(phys)
		s.sumLiveBytes += float64(live)
		s.spaceSamples++
	}
	s.statsMu.Unlock()
	return phys, live
}

// Metrics snapshots all measured I/O, OS block counts, GC activity, and
// derived flash-wear / $/TB-month economics.
func (s *Store) Metrics(workloadName string) Metrics {
	logicalFileBytes, physBytes, liveBytes, liveObjs, activeSegs, totalFrags := s.inspectDiskUsage()
	endProcIO := ReadProcIO()
	wall := time.Since(s.startTime)

	s.statsMu.Lock()
	peakPhys := s.peakPhysicalBytes
	if physBytes > peakPhys {
		peakPhys = physBytes
	}
	var avgPhys, avgLive, avgSpaceAmp float64
	if s.spaceSamples > 0 && s.sumLiveBytes > 0 {
		avgPhys = s.sumPhysicalBytes / float64(s.spaceSamples)
		avgLive = s.sumLiveBytes / float64(s.spaceSamples)
		avgSpaceAmp = avgPhys / avgLive
	} else if liveBytes > 0 {
		avgPhys = float64(physBytes)
		avgLive = float64(liveBytes)
		avgSpaceAmp = float64(physBytes) / float64(liveBytes)
	}

	var finalSpaceAmp float64
	if liveBytes > 0 {
		finalSpaceAmp = float64(physBytes) / float64(liveBytes)
	}

	m := Metrics{
		Policy:                  s.cfg.Policy,
		Workload:                workloadName,
		UserBytesWritten:        s.userBytesWritten,
		GCCopiedLiveBytes:       s.gcCopiedLiveBytes,
		MetadataBytesWritten:    s.metadataBytesWritten,
		LiveObjects:             liveObjs,
		LiveBytes:               liveBytes,
		LogicalFileBytes:        logicalFileBytes,
		PhysicalDiskBytes:       physBytes,
		PeakPhysicalDiskBytes:   peakPhys,
		AvgPhysicalDiskBytes:    round4(avgPhys),
		AvgLiveBytes:            round4(avgLive),
		AvgSpaceAmplification:   round4(avgSpaceAmp),
		FinalSpaceAmplification: round4(finalSpaceAmp),
		ActiveSegments:          activeSegs,
		PunchHoleSupported:      s.punchHoleSupported,
		PunchHoleCalls:          s.punchHoleCalls,
		PunchHoleRequestedBytes: s.punchHoleRequestedBytes,
		OSBlocksReclaimedBytes:  s.osBlocksReclaimedBytes,
		DiscardCycles:           s.discardCycles,
		CompactionCycles:        s.compactionCycles,
		ExtentFragments:         totalFrags,
		FSBlockSizeBytes:        s.fsBlockSize,
		WallDuration:            wall,
		WallDurationMs:          round4(float64(wall.Microseconds()) / 1000.0),
	}
	s.statsMu.Unlock()

	if s.startProcIO.Supported && endProcIO.Supported {
		m.ProcIODelta = ProcIOSnapshot{
			Supported:           true,
			RChar:               endProcIO.RChar - s.startProcIO.RChar,
			WChar:               endProcIO.WChar - s.startProcIO.WChar,
			ReadBytes:           endProcIO.ReadBytes - s.startProcIO.ReadBytes,
			WriteBytes:          endProcIO.WriteBytes - s.startProcIO.WriteBytes,
			CancelledWriteBytes: endProcIO.CancelledWriteBytes - s.startProcIO.CancelledWriteBytes,
		}
	}

	m.FlashWear, m.Cost = ComputeEconomics(&m, s.cfg.Drive)
	return m
}

func (s *Store) inspectDiskUsage() (logicalFileBytes, physBytes, livePayloadBytes int64, liveObjs, activeSegs, totalFrags int) {
	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return
	}
	tenants := make([]*tenantState, 0, len(s.tenants))
	for _, ts := range s.tenants {
		tenants = append(tenants, ts)
	}
	s.mu.RUnlock()

	for _, ts := range tenants {
		ts.mu.RLock()
		liveObjs += len(ts.index)
		for _, meta := range ts.index {
			livePayloadBytes += meta.PayloadLen
		}
		activeSegs += len(ts.segments)
		for _, seg := range ts.segments {
			totalFrags += len(seg.punchedExtents)
			if seg.file != nil {
				l, p, blk, err := FileDiskUsage(seg.file)
				if err == nil {
					logicalFileBytes += l
					physBytes += p
					if blk > 0 {
						s.statsMu.Lock()
						s.fsBlockSize = blk
						s.statsMu.Unlock()
					}
				}
			}
		}
		ts.mu.RUnlock()
	}
	return
}

func (s *Store) ensureActiveSegmentLocked(ts *tenantState, needed int64, seq uint64) (*segment, error) {
	if ts.activeSeg != nil && !ts.activeSeg.sealed {
		if ts.activeSeg.writeOffset+needed <= s.cfg.SegmentCapacityBytes || ts.activeSeg.writeOffset <= segHeaderSize {
			return ts.activeSeg, nil
		}
		ts.activeSeg.sealed = true
		ts.activeSeg = nil
	}
	seg, err := s.createSegmentLocked(ts, seq)
	if err != nil {
		return nil, err
	}
	ts.activeSeg = seg
	return seg, nil
}

func (s *Store) createSegmentLocked(ts *tenantState, seq uint64) (*segment, error) {
	id := s.segIDGen.Add(1)
	safeTenant := sanitizeName(ts.id)
	path := filepath.Join(s.dir, fmt.Sprintf("seg_%s_%06d.dat", safeTenant, id))
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, fmt.Errorf("discogc: create segment %s: %w", path, err)
	}

	var hdr [segHeaderSize]byte
	binary.LittleEndian.PutUint32(hdr[0:4], segMagic)
	binary.LittleEndian.PutUint32(hdr[4:8], 1)
	binary.LittleEndian.PutUint64(hdr[8:16], id)
	th := sha256.Sum256([]byte(ts.id))
	copy(hdr[16:48], th[:])

	if _, err := f.WriteAt(hdr[:], 0); err != nil {
		f.Close()
		return nil, fmt.Errorf("discogc: write segment header: %w", err)
	}

	s.statsMu.Lock()
	s.metadataBytesWritten += segHeaderSize
	s.statsMu.Unlock()

	seg := &segment{
		id:          id,
		tenantID:    ts.id,
		path:        path,
		file:        f,
		writeOffset: segHeaderSize,
		createdSeq:  seq,
		lastModSeq:  seq,
		objects:     make(map[string]*ObjectMeta),
	}
	ts.segments[id] = seg
	return seg, nil
}

func (s *Store) appendNeedleLocked(seg *segment, tenantID, key string, payload []byte, digest [32]byte, crc uint32, seq uint64) (*ObjectMeta, int64, error) {
	rawLen := needleHeaderSize + int64(len(key)) + int64(len(payload))
	align := smallAlign
	if s.cfg.AlignObjectsAtOrAboveGranularity && rawLen >= s.cfg.DiscardGranularity {
		align = s.cfg.DiscardGranularity
	}

	startOffset := ceilAlign(seg.writeOffset, align)
	leadPad := startOffset - seg.writeOffset
	totalLen := ceilAlign(rawLen, align)
	tailPad := totalLen - rawLen

	buf := make([]byte, totalLen)
	binary.LittleEndian.PutUint32(buf[0:4], needleMagic)
	binary.LittleEndian.PutUint32(buf[4:8], uint32(len(key)))
	binary.LittleEndian.PutUint64(buf[8:16], uint64(len(payload)))
	binary.LittleEndian.PutUint64(buf[16:24], seq)
	binary.LittleEndian.PutUint32(buf[24:28], crc)
	binary.LittleEndian.PutUint32(buf[28:32], uint32(totalLen))
	copy(buf[needleHeaderSize:], key)
	copy(buf[needleHeaderSize+int64(len(key)):], payload)

	if _, err := seg.file.WriteAt(buf, startOffset); err != nil {
		return nil, 0, fmt.Errorf("discogc: write needle: %w", err)
	}
	if s.cfg.SyncOnWrite {
		if err := seg.file.Sync(); err != nil {
			return nil, 0, fmt.Errorf("discogc: sync segment: %w", err)
		}
	}

	seg.writeOffset = startOffset + totalLen

	// Record leading alignment gap as dead space so adjacent dead extents can
	// coalesce cleanly without leaving un-reclaimable gaps.
	if leadPad > 0 {
		seg.deadExtents = insertExtent(seg.deadExtents, Extent{Start: startOffset - leadPad, End: startOffset})
		seg.deadUnpunchedBytes += leadPad
	}

	meta := &ObjectMeta{
		TenantID:   tenantID,
		Key:        key,
		SegmentID:  seg.id,
		Offset:     startOffset,
		HeaderLen:  needleHeaderSize + int64(len(key)),
		PayloadLen: int64(len(payload)),
		TotalLen:   totalLen,
		SHA256:     digest,
		Seq:        seq,
	}
	metaOverhead := needleHeaderSize + int64(len(key)) + leadPad + tailPad
	return meta, metaOverhead, nil
}

func readAndVerifyNeedle(f *os.File, meta *ObjectMeta) ([]byte, error) {
	rawSize := meta.HeaderLen + meta.PayloadLen
	buf := make([]byte, rawSize)
	if _, err := f.ReadAt(buf, meta.Offset); err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("discogc: read needle at %d: %w", meta.Offset, err)
	}
	magic := binary.LittleEndian.Uint32(buf[0:4])
	keyLen := int64(binary.LittleEndian.Uint32(buf[4:8]))
	payloadLen := int64(binary.LittleEndian.Uint64(buf[8:16]))
	storedCRC := binary.LittleEndian.Uint32(buf[24:28])
	if magic != needleMagic || keyLen != int64(len(meta.Key)) || payloadLen != meta.PayloadLen {
		return nil, ErrCorruptData
	}
	if string(buf[needleHeaderSize:needleHeaderSize+keyLen]) != meta.Key {
		return nil, ErrCorruptData
	}
	payload := buf[needleHeaderSize+keyLen : needleHeaderSize+keyLen+payloadLen]
	if crc32.ChecksumIEEE(payload) != storedCRC {
		return nil, ErrCorruptData
	}
	if sha256.Sum256(payload) != meta.SHA256 {
		return nil, ErrCorruptData
	}
	return payload, nil
}

func markDeadLocked(seg *segment, meta *ObjectMeta, seq uint64) {
	if _, ok := seg.objects[meta.Key]; ok {
		delete(seg.objects, meta.Key)
		seg.liveBytes -= meta.TotalLen
		if seg.liveBytes < 0 {
			seg.liveBytes = 0
		}
	}
	seg.deadExtents = insertExtent(seg.deadExtents, Extent{Start: meta.Offset, End: meta.Offset + meta.TotalLen})
	seg.deadUnpunchedBytes += meta.TotalLen
	seg.lastModSeq = seq
}

func (s *Store) maybeRunSegmentGCLocked(ts *tenantState, seg *segment, seq uint64) error {
	switch s.cfg.Policy {
	case PolicyNoGC:
		return nil
	case PolicyAdaptiveDiscard:
		return s.gcAdaptiveSegmentLocked(ts, seg, seq)
	case PolicyGreedyCompaction:
		if seg.sealed && segmentDeadRatio(seg) >= s.cfg.GarbageTriggerRatio {
			return s.runGCLocked(ts, seq)
		}
		return nil
	case PolicyCostBenefitCompaction:
		// Defer to segment-seal / periodic GC passes so hot segments can absorb
		// additional overwrites before age-weighted cost-benefit scoring runs.
		return nil
	case PolicyBlindSubblockDiscard:
		return s.runBlindDiscardLocked(seg)
	default:
		return nil
	}
}

func (s *Store) runGCLocked(ts *tenantState, seq uint64) error {
	switch s.cfg.Policy {
	case PolicyNoGC:
		return nil
	case PolicyGreedyCompaction:
		return s.runGreedyCompactionLocked(ts, seq)
	case PolicyCostBenefitCompaction:
		return s.runCostBenefitCompactionLocked(ts, seq)
	case PolicyAdaptiveDiscard:
		// Copy slice of segments since compaction may remove entries from ts.segments.
		segs := make([]*segment, 0, len(ts.segments))
		for _, seg := range ts.segments {
			segs = append(segs, seg)
		}
		sort.Slice(segs, func(i, j int) bool { return segs[i].id < segs[j].id })
		for _, seg := range segs {
			if _, stillPresent := ts.segments[seg.id]; !stillPresent {
				continue
			}
			if err := s.gcAdaptiveSegmentLocked(ts, seg, seq); err != nil {
				return err
			}
		}
		return nil
	case PolicyBlindSubblockDiscard:
		for _, seg := range ts.segments {
			if err := s.runBlindDiscardLocked(seg); err != nil {
				return err
			}
		}
		return nil
	default:
		return nil
	}
}

func (s *Store) runGreedyCompactionLocked(ts *tenantState, seq uint64) error {
	var candidates []*segment
	for _, seg := range ts.segments {
		if !seg.sealed {
			continue
		}
		if segmentDeadRatio(seg) >= s.cfg.GarbageTriggerRatio {
			candidates = append(candidates, seg)
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	// Greedy: highest dead fraction first.
	sort.Slice(candidates, func(i, j int) bool {
		ri, rj := segmentDeadRatio(candidates[i]), segmentDeadRatio(candidates[j])
		if ri != rj {
			return ri > rj
		}
		return candidates[i].id < candidates[j].id
	})
	for _, seg := range candidates {
		if _, ok := ts.segments[seg.id]; !ok {
			continue
		}
		if err := s.compactSegmentLocked(ts, seg, seq); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) runCostBenefitCompactionLocked(ts *tenantState, seq uint64) error {
	type scoredSeg struct {
		seg   *segment
		score float64
	}
	var candidates []scoredSeg
	for _, seg := range ts.segments {
		if !seg.sealed {
			continue
		}
		deadRatio := segmentDeadRatio(seg)
		if deadRatio < s.cfg.GarbageTriggerRatio {
			continue
		}
		liveRatio := 1.0 - deadRatio
		age := float64(1)
		if seq > seg.lastModSeq {
			age = float64(seq - seg.lastModSeq)
		}
		// Rosenblum-Ousterhout LFS cost-benefit metric: ((1 - u) * age) / (1 + u)
		score := (deadRatio * age) / (1.0 + liveRatio)
		candidates = append(candidates, scoredSeg{seg: seg, score: score})
	}
	if len(candidates) == 0 {
		return nil
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].score != candidates[j].score {
			return candidates[i].score > candidates[j].score
		}
		return candidates[i].seg.id < candidates[j].seg.id
	})

	// Cost-benefit compacts at most half of eligible segments (at least 1) per
	// pass so hot segments have time to accumulate more overwrites before copying.
	limit := len(candidates) / 2
	if limit < 1 {
		limit = 1
	}
	for i := 0; i < limit; i++ {
		seg := candidates[i].seg
		if _, ok := ts.segments[seg.id]; !ok {
			continue
		}
		if err := s.compactSegmentLocked(ts, seg, seq); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) gcAdaptiveSegmentLocked(ts *tenantState, seg *segment, seq uint64) error {
	// Step 1: Reclaim any strictly interior block-aligned dead extents via
	// direct OS hole-punching with zero live-data copying.
	if len(seg.deadExtents) > 0 {
		var remaining []Extent
		didPunch := false
		for _, d := range seg.deadExtents {
			alignedStart := ceilAlign(d.Start, s.cfg.DiscardGranularity)
			alignedEnd := floorAlign(d.End, s.cfg.DiscardGranularity)
			punchLen := alignedEnd - alignedStart
			if punchLen >= s.cfg.MinDiscardExtent {
				reclaimed, supported, err := PunchHole(seg.file, alignedStart, punchLen)
				if err != nil {
					return fmt.Errorf("discogc: punch hole seg %d [%d,%d): %w", seg.id, alignedStart, alignedEnd, err)
				}
				s.statsMu.Lock()
				s.punchHoleCalls++
				s.punchHoleRequestedBytes += punchLen
				s.osBlocksReclaimedBytes += reclaimed
				if !supported {
					s.punchHoleSupported = false
				}
				s.punchHoleProbed = true
				s.statsMu.Unlock()

				if supported {
					didPunch = true
					seg.deadUnpunchedBytes -= punchLen
					if seg.deadUnpunchedBytes < 0 {
						seg.deadUnpunchedBytes = 0
					}
					seg.punchedBytes += punchLen
					seg.punchedExtents = insertExtent(seg.punchedExtents, Extent{Start: alignedStart, End: alignedEnd})
					// Preserve unaligned dead head/tail fragments for future coalescing.
					if alignedStart > d.Start {
						remaining = append(remaining, Extent{Start: d.Start, End: alignedStart})
					}
					if d.End > alignedEnd {
						remaining = append(remaining, Extent{Start: alignedEnd, End: d.End})
					}
					continue
				}
			}
			remaining = append(remaining, d)
		}
		seg.deadExtents = remaining
		if didPunch {
			if s.cfg.SyncOnWrite {
				if err := seg.file.Sync(); err != nil {
					return err
				}
			}
			s.statsMu.Lock()
			s.discardCycles++
			s.statsMu.Unlock()
		}
	}

	// Step 2: Check whether selective copy-compaction fallback is needed.
	// Extent-map fragmentation cap applies to both active and sealed segments so
	// hole count never grows unbounded inside a long-lived active segment.
	if s.cfg.MaxExtentFragments > 0 && len(seg.punchedExtents) > s.cfg.MaxExtentFragments {
		if ts.activeSeg == seg {
			seg.sealed = true
			ts.activeSeg = nil
		}
		return s.compactSegmentLocked(ts, seg, seq)
	}

	if !seg.sealed {
		return nil
	}

	// (a) Fully dead sealed segment: drop file directly with zero live copies!
	if len(seg.objects) == 0 {
		return s.dropEmptySegmentLocked(ts, seg)
	}

	denom := seg.writeOffset
	if denom <= 0 {
		return nil
	}

	// (b) Sub-block fragmentation fallback: un-reclaimable dead bytes exceed threshold.
	unreclaimableRatio := float64(seg.deadUnpunchedBytes) / float64(denom)
	if unreclaimableRatio >= s.cfg.SubblockFragFallbackRatio {
		return s.compactSegmentLocked(ts, seg, seq)
	}

	// (c) Near-empty sparse segment: compact the few remaining live objects so
	// the sparse segment file descriptor and inode metadata can be freed.
	liveRatio := float64(seg.liveBytes) / float64(denom)
	if liveRatio <= s.cfg.SparseEmptyLiveThreshold && seg.punchedBytes > 0 {
		return s.compactSegmentLocked(ts, seg, seq)
	}

	return nil
}

func (s *Store) runBlindDiscardLocked(seg *segment) error {
	if len(seg.deadExtents) == 0 {
		return nil
	}
	didCall := false
	for _, d := range seg.deadExtents {
		if d.Len() <= 0 {
			continue
		}
		// Negative control: calls PunchHole directly on raw unaligned dead extent
		// without checking DiscardGranularity and never falls back to compaction.
		reclaimed, supported, err := PunchHole(seg.file, d.Start, d.Len())
		if err != nil {
			return fmt.Errorf("discogc: blind punch hole seg %d [%d,%d): %w", seg.id, d.Start, d.End, err)
		}
		didCall = true
		s.statsMu.Lock()
		s.punchHoleCalls++
		s.punchHoleRequestedBytes += d.Len()
		s.osBlocksReclaimedBytes += reclaimed
		if !supported {
			s.punchHoleSupported = false
		}
		s.punchHoleProbed = true
		s.statsMu.Unlock()
		if reclaimed > 0 {
			seg.punchedBytes += reclaimed
			seg.punchedExtents = insertExtent(seg.punchedExtents, d)
		}
	}
	// Blind discard marks extents as attempted (does not re-attempt forever),
	// leaving un-reclaimed sub-block bytes trapped on disk because it never compacts.
	seg.deadExtents = nil
	if didCall {
		s.statsMu.Lock()
		s.discardCycles++
		s.statsMu.Unlock()
	}
	return nil
}

func (s *Store) dropEmptySegmentLocked(ts *tenantState, seg *segment) error {
	if ts.activeSeg == seg {
		ts.activeSeg = nil
	}
	delete(ts.segments, seg.id)
	if seg.file != nil {
		_ = seg.file.Close()
		seg.file = nil
	}
	_ = os.Remove(seg.path)
	return nil
}

func (s *Store) compactSegmentLocked(ts *tenantState, seg *segment, seq uint64) error {
	if len(seg.objects) == 0 {
		return s.dropEmptySegmentLocked(ts, seg)
	}

	// Deterministic copy order by offset.
	metas := make([]*ObjectMeta, 0, len(seg.objects))
	for _, m := range seg.objects {
		metas = append(metas, m)
	}
	sort.Slice(metas, func(i, j int) bool { return metas[i].Offset < metas[j].Offset })

	var copiedPayload int64
	var copiedMeta int64
	for _, oldMeta := range metas {
		payload, err := readAndVerifyNeedle(seg.file, oldMeta)
		if err != nil {
			return fmt.Errorf("discogc: compact read seg %d key %q: %w", seg.id, oldMeta.Key, err)
		}
		crc := crc32.ChecksumIEEE(payload)
		dst, err := s.ensureActiveSegmentLocked(ts, int64(len(oldMeta.Key))+int64(len(payload))+needleHeaderSize, seq)
		if err != nil {
			return err
		}
		if dst.id == seg.id {
			seg.sealed = true
			ts.activeSeg = nil
			dst, err = s.ensureActiveSegmentLocked(ts, int64(len(oldMeta.Key))+int64(len(payload))+needleHeaderSize, seq)
			if err != nil {
				return err
			}
		}

		newMeta, metaOverhead, err := s.appendNeedleLocked(dst, ts.id, oldMeta.Key, payload, oldMeta.SHA256, crc, oldMeta.Seq)
		if err != nil {
			return err
		}
		ts.index[oldMeta.Key] = newMeta
		dst.objects[oldMeta.Key] = newMeta
		dst.liveBytes += newMeta.TotalLen
		dst.lastModSeq = seq

		copiedPayload += int64(len(payload))
		copiedMeta += metaOverhead

		if dst.writeOffset >= s.cfg.SegmentCapacityBytes {
			dst.sealed = true
			ts.activeSeg = nil
		}
	}

	if err := s.dropEmptySegmentLocked(ts, seg); err != nil {
		return err
	}

	s.statsMu.Lock()
	s.gcCopiedLiveBytes += copiedPayload
	s.metadataBytesWritten += copiedMeta
	s.compactionCycles++
	s.statsMu.Unlock()
	return nil
}

func segmentDeadRatio(seg *segment) float64 {
	if seg.writeOffset <= segHeaderSize {
		return 0
	}
	dead := seg.deadUnpunchedBytes + seg.punchedBytes
	return float64(dead) / float64(seg.writeOffset)
}

func ceilAlign(n, align int64) int64 {
	if align <= 1 {
		return n
	}
	rem := n % align
	if rem == 0 {
		return n
	}
	return n + (align - rem)
}

func floorAlign(n, align int64) int64 {
	if align <= 1 {
		return n
	}
	return n - (n % align)
}

// insertExtent inserts e into a sorted list of extents and coalesces any
// overlapping or directly adjacent ranges.
func insertExtent(list []Extent, e Extent) []Extent {
	if e.Len() <= 0 {
		return list
	}
	list = append(list, e)
	sort.Slice(list, func(i, j int) bool {
		if list[i].Start != list[j].Start {
			return list[i].Start < list[j].Start
		}
		return list[i].End < list[j].End
	})
	merged := list[:0]
	for _, cur := range list {
		if len(merged) == 0 || cur.Start > merged[len(merged)-1].End {
			merged = append(merged, cur)
		} else if cur.End > merged[len(merged)-1].End {
			merged[len(merged)-1].End = cur.End
		}
	}
	return merged
}

func sanitizeName(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
			out = append(out, c)
		} else {
			out = append(out, '_')
		}
	}
	if len(out) == 0 {
		return "default"
	}
	return string(out)
}
