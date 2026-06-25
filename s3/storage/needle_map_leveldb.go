package storage

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/luxfi/database"
	"github.com/luxfi/database/factory"
	log "github.com/luxfi/log"

	"github.com/hanzoai/s3/s3/storage/idx"
	"github.com/hanzoai/s3/s3/storage/needle"
	"github.com/hanzoai/s3/s3/util"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/storage/needle_map"
	. "github.com/hanzoai/s3/s3/storage/types"
)

// needleMapBackend is the lux-canonical key-value backend for the on-disk
// volume needle-map index. ZAP-native zapdb. The needle-id -> (offset,size)
// byte layout and the watermark/recovery semantics are byte-identical to the
// historical leveldb index, so existing .ldb directories are unaffected.
const needleMapBackend = "zapdb"

// mark it every watermarkBatchSize operations
const watermarkBatchSize = 10000

var watermarkKey = []byte("idx_entry_watermark")

// openNeedleMapDB opens (creating if absent) a zapdb-backed
// database.Database at dir. A nil gatherer/logger keeps the store metric-free
// and silent, matching the embedded single-process volume-server use.
// BadgerDB (zapdb) recovers automatically on open, subsuming leveldb's
// explicit RecoverFile path.
func openNeedleMapDB(dir string) (database.Database, error) {
	// factory.New(name, dbPath, readOnly, config, gatherer, logger, metricsPrefix, meterDBRegName).
	return factory.New(needleMapBackend, dir, false, nil, nil, log.Noop(), "", "")
}

type LevelDbNeedleMap struct {
	baseNeedleMapper
	dbFileName    string
	db            database.Database
	ldbAccessLock sync.RWMutex
	exitChan      chan bool
	// no need to use atomic
	accessFlag int64
	ldbTimeout int64
	// recordCount is bumped on every Put/Delete, which run under a shared
	// RLock (lazy-load path) or no lock at all (resident path), so it must be
	// atomic. Add returns the post-increment value, giving each writer a
	// consistent count for its watermark math.
	recordCount atomic.Uint64
}

func NewLevelDbNeedleMap(dbFileName string, indexFile *os.File, ldbTimeout int64, version needle.Version) (m *LevelDbNeedleMap, err error) {
	m = &LevelDbNeedleMap{dbFileName: dbFileName}
	m.indexFile = indexFile
	if !isLevelDbFresh(dbFileName, indexFile) {
		glog.V(1).Infof("Start to Generate %s from %s", dbFileName, indexFile.Name())
		generateLevelDbFile(dbFileName, indexFile)
		glog.V(1).Infof("Finished Generating %s from %s", dbFileName, indexFile.Name())
	}
	if stat, err := indexFile.Stat(); err != nil {
		glog.Fatalf("stat file %s: %v", indexFile.Name(), err)
	} else {
		m.indexFileOffset = stat.Size()
	}
	glog.V(1).Infof("Opening %s...", dbFileName)

	if m.ldbTimeout == 0 {
		if m.db, err = openNeedleMapDB(dbFileName); err != nil {
			return
		}
		glog.V(1).Infof("Loading %s... , watermark: %d", dbFileName, getWatermark(m.db))
		recordCount := uint64(m.indexFileOffset / NeedleMapEntrySize)
		m.recordCount.Store(recordCount)
		watermark := (recordCount / watermarkBatchSize) * watermarkBatchSize
		err = setWatermark(m.db, watermark)
		if err != nil {
			glog.Fatalf("set watermark for %s error: %s\n", dbFileName, err)
			return
		}
	}
	mm, indexLoadError := newNeedleMapMetricFromIndexFile(indexFile, version)
	if indexLoadError != nil {
		return nil, indexLoadError
	}
	m.mapMetric = *mm
	m.ldbTimeout = ldbTimeout
	if m.ldbTimeout > 0 {
		m.exitChan = make(chan bool, 1)
		m.accessFlag = 0
		go lazyLoadingRoutine(m)
	}
	return
}

// isLevelDbFresh reports whether the on-disk index db is newer than the .idx
// it must mirror, so regeneration can be skipped. We always write the .idx
// first, so a db whose newest file post-dates the .idx is consistent.
//
// Returning false (rebuild) is always safe — it just replays the .idx — while
// a wrong true would trust a stale index. So any uncertainty (missing dir,
// empty dir, stat error) returns false and forces a rebuild from the .idx.
func isLevelDbFresh(dbFileName string, indexFile *os.File) bool {
	dbModTime, ok := newestModTime(dbFileName)
	if !ok {
		return false
	}
	indexStat, indexStatErr := indexFile.Stat()
	if indexStatErr != nil {
		glog.V(0).Infof("Can not stat file: %v", indexStatErr)
		return false
	}
	return dbModTime.After(indexStat.ModTime())
}

// newestModTime returns the most recent modification time among the regular
// files inside the db directory, and whether the directory holds any. An
// empty or absent directory yields (zero, false).
func newestModTime(dir string) (time.Time, bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return time.Time{}, false
	}
	var newest time.Time
	var found bool
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		found = true
		if info.ModTime().After(newest) {
			newest = info.ModTime()
		}
	}
	return newest, found
}

func generateLevelDbFile(dbFileName string, indexFile *os.File) error {
	db, err := openNeedleMapDB(dbFileName)
	if err != nil {
		return err
	}
	defer db.Close()

	watermark := getWatermark(db)
	if stat, err := indexFile.Stat(); err != nil {
		glog.Fatalf("stat file %s: %v", indexFile.Name(), err)
		return err
	} else {
		// A watermark past the end of the .idx means the .ldb is stale relative
		// to the index it must mirror (e.g. an interrupted compaction left the
		// old .ldb beside a freshly swapped, shorter .idx). Trusting it would
		// replay zero entries and silently poison the needle map, so rebuild
		// from offset 0 instead.
		// Compare in entries, not bytes: watermark*NeedleMapEntrySize can
		// overflow uint64 for a corrupted watermark and wrap past the size check.
		if watermark > uint64(stat.Size())/NeedleMapEntrySize {
			glog.Warningf("stale watermark %d for %s (filesize %d); rebuilding leveldb from start", watermark, dbFileName, stat.Size())
			watermark = 0
		}
		glog.V(1).Infof("generateLevelDbFile %s, watermark %d, num of entries:%d", dbFileName, watermark, (uint64(stat.Size())-watermark*NeedleMapEntrySize)/NeedleMapEntrySize)
	}
	return idx.WalkIndexFile(indexFile, watermark, func(key NeedleId, offset Offset, size Size) error {
		if !offset.IsZero() && !size.IsDeleted() {
			levelDbWrite(db, key, offset, size, false, 0)
		} else {
			levelDbDelete(db, key)
		}
		return nil
	})
}

func (m *LevelDbNeedleMap) Get(key NeedleId) (element *needle_map.NeedleValue, ok bool) {
	if m.ldbTimeout > 0 {
		if err := m.ensureLdbLoaded(); err != nil {
			return nil, false
		}
		defer m.ldbAccessLock.RUnlock()
	}
	return m.getFromDb(key)
}

func (m *LevelDbNeedleMap) getFromDb(key NeedleId) (element *needle_map.NeedleValue, ok bool) {
	bytes := make([]byte, NeedleIdSize)
	NeedleIdToBytes(bytes[0:NeedleIdSize], key)
	data, err := m.db.Get(bytes)
	if err != nil || len(data) != OffsetSize+SizeSize {
		return nil, false
	}
	offset := BytesToOffset(data[0:OffsetSize])
	size := BytesToSize(data[OffsetSize : OffsetSize+SizeSize])
	return &needle_map.NeedleValue{Key: key, Offset: offset, Size: size}, true
}

func (m *LevelDbNeedleMap) Put(key NeedleId, offset Offset, size Size) error {
	var oldSize Size
	var watermark uint64
	if m.ldbTimeout > 0 {
		if err := m.ensureLdbLoaded(); err != nil {
			return err
		}
		defer m.ldbAccessLock.RUnlock()
	}
	if oldNeedle, ok := m.getFromDb(key); ok {
		oldSize = oldNeedle.Size
	}
	m.logPut(key, oldSize, size)
	// write to index file first
	if err := m.appendToIndexFile(key, offset, size); err != nil {
		return fmt.Errorf("cannot write to indexfile %s: %v", m.indexFile.Name(), err)
	}
	recordCount := m.recordCount.Add(1)
	if recordCount%watermarkBatchSize != 0 {
		watermark = 0
	} else {
		watermark = (recordCount / watermarkBatchSize) * watermarkBatchSize
		glog.V(1).Infof("put cnt:%d for %s,watermark: %d", recordCount, m.dbFileName, watermark)
	}
	return levelDbWrite(m.db, key, offset, size, watermark == 0, watermark)
}

func getWatermark(db database.Database) uint64 {
	data, err := db.Get(watermarkKey)
	if err != nil || len(data) != 8 {
		glog.V(1).Infof("read previous watermark from db: %v, %d", err, len(data))
		return 0
	}
	return util.BytesToUint64(data)
}

func setWatermark(db database.Database, watermark uint64) error {
	glog.V(3).Infof("set watermark %d", watermark)
	var wmBytes = make([]byte, 8)
	util.Uint64toBytes(wmBytes, watermark)
	if err := db.Put(watermarkKey, wmBytes); err != nil {
		return fmt.Errorf("failed to setWatermark: %w", err)
	}
	return nil
}

func levelDbWrite(db database.Database, key NeedleId, offset Offset, size Size, updateWatermark bool, watermark uint64) error {

	bytes := needle_map.ToBytes(key, offset, size)

	if err := db.Put(bytes[0:NeedleIdSize], bytes[NeedleIdSize:NeedleIdSize+OffsetSize+SizeSize]); err != nil {
		return fmt.Errorf("failed to write leveldb: %w", err)
	}
	// set watermark
	if updateWatermark {
		return setWatermark(db, watermark)
	}
	return nil
}

func levelDbDelete(db database.Database, key NeedleId) error {
	bytes := make([]byte, NeedleIdSize)
	NeedleIdToBytes(bytes, key)
	return db.Delete(bytes)
}

func (m *LevelDbNeedleMap) Delete(key NeedleId, offset Offset) error {
	var watermark uint64
	if m.ldbTimeout > 0 {
		if err := m.ensureLdbLoaded(); err != nil {
			return err
		}
		defer m.ldbAccessLock.RUnlock()
	}
	oldNeedle, found := m.getFromDb(key)
	if !found || oldNeedle.Size.IsDeleted() {
		return nil
	}
	m.logDelete(oldNeedle.Size)

	// write to index file first
	if err := m.appendToIndexFile(key, offset, TombstoneFileSize); err != nil {
		return err
	}
	recordCount := m.recordCount.Add(1)
	if recordCount%watermarkBatchSize != 0 {
		watermark = 0
	} else {
		watermark = (recordCount / watermarkBatchSize) * watermarkBatchSize
	}
	return levelDbWrite(m.db, key, oldNeedle.Offset, -oldNeedle.Size, watermark == 0, watermark)
}

func (m *LevelDbNeedleMap) Close() {
	if m.indexFile != nil {
		indexFileName := m.indexFile.Name()
		if err := m.indexFile.Sync(); err != nil {
			glog.Warningf("sync file %s failed: %v", indexFileName, err)
		}
		if err := m.indexFile.Close(); err != nil {
			glog.Warningf("close index file %s failed: %v", indexFileName, err)
		}
	}

	if m.db != nil {
		if err := m.db.Close(); err != nil {
			glog.Warningf("close levelDB failed: %v", err)
		}
	}
	if m.ldbTimeout > 0 {
		m.exitChan <- true
	}
}

func (m *LevelDbNeedleMap) Destroy() error {
	m.Close()
	os.Remove(m.indexFile.Name())
	return os.RemoveAll(m.dbFileName)
}

func (m *LevelDbNeedleMap) UpdateNeedleMap(v *Volume, indexFile *os.File, ldbTimeout int64) error {
	if v.nm != nil {
		v.nm.Close()
		v.nm = nil
	}
	defer func() {
		if v.tmpNm != nil {
			v.tmpNm.Close()
			v.tmpNm = nil
		}
	}()
	levelDbFile := v.FileName(".ldb")
	m.indexFile = indexFile
	err := os.RemoveAll(levelDbFile)
	if err != nil {
		return err
	}
	if err = os.Rename(v.FileName(".cpldb"), levelDbFile); err != nil {
		return fmt.Errorf("rename %s: %v", levelDbFile, err)
	}

	db, err := openNeedleMapDB(levelDbFile)
	if err != nil {
		return err
	}
	m.db = db

	stat, e := indexFile.Stat()
	if e != nil {
		glog.Fatalf("stat file %s: %v", indexFile.Name(), e)
		return e
	}
	m.indexFileOffset = stat.Size()
	recordCount := uint64(stat.Size() / NeedleMapEntrySize)
	m.recordCount.Store(recordCount)

	//set watermark
	watermark := (recordCount / watermarkBatchSize) * watermarkBatchSize
	err = setWatermark(db, uint64(watermark))
	if err != nil {
		glog.Fatalf("setting watermark failed %s: %v", indexFile.Name(), err)
		return err
	}
	v.nm = m
	v.tmpNm = nil
	m.ldbTimeout = ldbTimeout
	if m.ldbTimeout > 0 {
		m.exitChan = make(chan bool, 1)
		m.accessFlag = 0
		go lazyLoadingRoutine(m)
	}
	return e
}

func (m *LevelDbNeedleMap) DoOffsetLoading(v *Volume, indexFile *os.File, startFrom uint64) (err error) {
	glog.V(0).Infof("loading idx to leveldb from offset %d for file: %s", startFrom, indexFile.Name())
	version := needle.GetCurrentVersion()
	if v != nil {
		version = v.Version()
	}
	dbFileName := v.FileName(".cpldb")
	db, dbErr := openNeedleMapDB(dbFileName)
	defer func() {
		if dbErr == nil {
			db.Close()
		}
		if err != nil {
			os.RemoveAll(dbFileName)
		}

	}()
	if dbErr != nil {
		return dbErr
	}

	err = idx.WalkIndexFile(indexFile, startFrom, func(key NeedleId, offset Offset, size Size) (e error) {
		m.mapMetric.FileCounter++
		m.mapMetric.MaybeSetMaxNeedleEnd(offset, size, version)
		bytes := make([]byte, NeedleIdSize)
		NeedleIdToBytes(bytes[0:NeedleIdSize], key)
		// fresh loading
		if startFrom == 0 {
			m.mapMetric.FileByteCounter += uint64(size)
			e = levelDbWrite(db, key, offset, size, false, 0)
			return e
		}
		// increment loading
		data, err := db.Get(bytes)
		if err != nil {
			if err != database.ErrNotFound {
				// unexpected error
				return err
			}
			// new needle, unlikely happen
			m.mapMetric.FileByteCounter += uint64(size)
			e = levelDbWrite(db, key, offset, size, false, 0)
		} else {
			// needle is found
			oldSize := BytesToSize(data[OffsetSize : OffsetSize+SizeSize])
			oldOffset := BytesToOffset(data[0:OffsetSize])
			if !offset.IsZero() && !size.IsDeleted() {
				// updated needle
				m.mapMetric.FileByteCounter += uint64(size)
				if !oldOffset.IsZero() && !oldSize.IsDeleted() {
					m.mapMetric.DeletionCounter++
					m.mapMetric.DeletionByteCounter += uint64(oldSize)
				}
				e = levelDbWrite(db, key, offset, size, false, 0)
			} else {
				// deleted needle
				m.mapMetric.DeletionCounter++
				m.mapMetric.DeletionByteCounter += uint64(oldSize)
				e = levelDbDelete(db, key)
			}
		}
		return e
	})
	return err
}

func (m *LevelDbNeedleMap) ensureLdbLoaded() error {
	for {
		m.ldbAccessLock.RLock()
		if m.db != nil {
			return nil
		}
		m.ldbAccessLock.RUnlock()
		m.ldbAccessLock.Lock()
		if m.db == nil {
			if err := reloadLdb(m); err != nil {
				m.ldbAccessLock.Unlock()
				return err
			}
		}
		m.ldbAccessLock.Unlock()
	}
}

func reloadLdb(m *LevelDbNeedleMap) (err error) {
	if m.db != nil {
		return nil
	}
	glog.V(1).Infof("reloading leveldb %s", m.dbFileName)
	m.accessFlag = 1
	if m.db, err = openNeedleMapDB(m.dbFileName); err != nil {
		glog.Fatalf("open %s failed:%v", m.dbFileName, err)
		return err
	}
	return nil
}

func unloadLdb(m *LevelDbNeedleMap) (err error) {
	m.ldbAccessLock.Lock()
	defer m.ldbAccessLock.Unlock()
	if m.db != nil {
		glog.V(1).Infof("reached max idle count, unload leveldb, %s", m.dbFileName)
		m.db.Close()
		m.db = nil
	}
	return nil
}

func lazyLoadingRoutine(m *LevelDbNeedleMap) (err error) {
	glog.V(1).Infof("lazyLoadingRoutine %s", m.dbFileName)
	var accessRecord int64
	accessRecord = 1
	for {
		select {
		case exit := <-m.exitChan:
			if exit {
				glog.V(1).Infof("exit from lazyLoadingRoutine")
				return nil
			}
		case <-time.After(time.Hour * 1):
			glog.V(1).Infof("timeout %s", m.dbFileName)
			if m.accessFlag == 0 {
				accessRecord++
				glog.V(1).Infof("accessRecord++")
				if accessRecord >= m.ldbTimeout {
					unloadLdb(m)
				}
			} else {
				glog.V(1).Infof("reset accessRecord %s", m.dbFileName)
				// reset accessRecord
				accessRecord = 0
				m.accessFlag = 0
			}
			continue
		}
	}
}
