// Package luxdb implements the filer metadata store over luxfi/database, the
// lux-canonical key-value abstraction. The backend is selectable: "zapdb"
// (default, ZAP-native), "pebble", or "leveldb". The key encoding is
// byte-identical to the historical leveldb store so on-disk layout matches.
package luxdb

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/luxfi/database"
	"github.com/luxfi/database/factory"
	log "github.com/luxfi/log"

	"github.com/hanzoai/s3/s3/filer"
	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/pb/filer_pb"
	s3_util "github.com/hanzoai/s3/s3/util"
)

const (
	DIR_FILE_SEPARATOR = byte(0x00)

	// defaultBackend is the canonical metadata backend: ZAP-native zapdb.
	defaultBackend = "zapdb"
)

var (
	_ = filer.FilerStore(&LuxDBStore{})
	_ = filer.Debuggable(&LuxDBStore{})
)

func init() {
	filer.Stores = append(filer.Stores, &LuxDBStore{})
}

// LuxDBStore is a filer.FilerStore backed by a luxfi/database.Database.
type LuxDBStore struct {
	db database.Database
}

func (store *LuxDBStore) GetName() string {
	return "luxdb"
}

func (store *LuxDBStore) Initialize(configuration s3_util.Configuration, prefix string) (err error) {
	dir := configuration.GetString(prefix + "dir")
	backend := configuration.GetString(prefix + "backend")
	// Read as a STRING, not GetBool: the Configuration interface has no IsSet,
	// so a bool cannot distinguish "absent" from "explicitly false" — the exact
	// confusion that made this knob one-way in the backend. "" means leave the
	// backend default alone.
	return store.initialize(backend, dir, configuration.GetString(prefix+"syncWrites"))
}

func (store *LuxDBStore) initialize(backend, dir, syncWrites string) (err error) {
	if backend == "" {
		backend = defaultBackend
	}
	// "pebble" is the operator-facing alias for the luxfi/database "pebbledb"
	// backend name.
	if backend == "pebble" {
		backend = "pebbledb"
	}

	glog.V(0).Infof("filer luxdb store: backend=%s dir=%s", backend, dir)
	os.MkdirAll(dir, 0755)
	if err := s3_util.TestFolderWritable(dir); err != nil {
		return fmt.Errorf("check luxdb folder %s writable: %s", dir, err)
	}

	// Durability knob, off the hot path but decisive for it. The backend fsyncs
	// every commit by default, which is correct and is what we keep unless an
	// operator says otherwise — but on network-attached storage that fsync is
	// the dominant cost of metadata work. Measured on this filer: a single
	// metadata create is 11.39 ms on a cloud block volume against 1.41 ms for
	// the same binary on tmpfs, so the setting is worth an order of magnitude
	// to a deployment that can accept losing the last writes on a host crash.
	//
	// Passing nil config here is what made it unreachable, so an operator could
	// set filer.toml's syncWrites and nothing happened.
	var cfgJSON []byte
	{
		switch strings.ToLower(strings.TrimSpace(syncWrites)) {
		case "false", "0", "off", "no":
			cfgJSON = []byte(`{"syncWrites": false}`)
			glog.V(0).Infof("filer luxdb store: syncWrites=false (durability relaxed by config)")
		case "true", "1", "on", "yes":
			cfgJSON = []byte(`{"syncWrites": true}`)
		case "":
			// unset: keep the backend default
		default:
			return fmt.Errorf("filer luxdb: invalid syncWrites %q (want true or false)", syncWrites)
		}
	}

	// factory.New(name, dbPath, readOnly, config, gatherer, logger, metricsPrefix, meterDBRegName).
	// zapdb is always registered (default); pebbledb/leveldb register under
	// their build tags. A nil gatherer/logger keeps the store metric-free and
	// silent, matching the embedded single-process filer use.
	db, err := factory.New(backend, dir, false, cfgJSON, nil, log.Noop(), "", "")
	if err != nil {
		return fmt.Errorf("open luxdb backend %q at %s: %w", backend, dir, err)
	}
	store.db = db
	return nil
}

func (store *LuxDBStore) BeginTransaction(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
func (store *LuxDBStore) CommitTransaction(ctx context.Context) error {
	return nil
}
func (store *LuxDBStore) RollbackTransaction(ctx context.Context) error {
	return nil
}

func (store *LuxDBStore) InsertEntry(ctx context.Context, entry *filer.Entry) (err error) {
	key := genKey(entry.DirAndName())

	value, err := entry.EncodeAttributesAndChunks()
	if err != nil {
		return fmt.Errorf("encoding %s %+v: %v", entry.FullPath, entry.Attr, err)
	}

	if len(entry.GetChunks()) > filer.CountEntryChunksForGzip {
		value = s3_util.MaybeGzipData(value)
	}

	if err = store.db.Put(key, value); err != nil {
		return fmt.Errorf("persisting %s : %v", entry.FullPath, err)
	}

	return nil
}

func (store *LuxDBStore) UpdateEntry(ctx context.Context, entry *filer.Entry) (err error) {
	return store.InsertEntry(ctx, entry)
}

// BatchInsertEntries inserts multiple entries in a single batch write. This is
// more efficient than inserting entries one by one as it reduces the number of
// write operations and syncs to disk.
func (store *LuxDBStore) BatchInsertEntries(ctx context.Context, entries []*filer.Entry) error {
	if len(entries) == 0 {
		return nil
	}

	batch := store.db.NewBatch()

	for _, entry := range entries {
		key := genKey(entry.DirAndName())

		value, err := entry.EncodeAttributesAndChunks()
		if err != nil {
			return fmt.Errorf("encoding %s %+v: %w", entry.FullPath, entry.Attr, err)
		}

		if len(entry.GetChunks()) > filer.CountEntryChunksForGzip {
			value = s3_util.MaybeGzipData(value)
		}

		if err := batch.Put(key, value); err != nil {
			return fmt.Errorf("batch put %s: %w", entry.FullPath, err)
		}
	}

	if err := batch.Write(); err != nil {
		return fmt.Errorf("batch write: %w", err)
	}

	return nil
}

func (store *LuxDBStore) FindEntry(ctx context.Context, fullpath s3_util.FullPath) (entry *filer.Entry, err error) {
	key := genKey(fullpath.DirAndName())

	data, err := store.db.Get(key)

	if err == database.ErrNotFound {
		return nil, filer_pb.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get %s : %v", fullpath, err)
	}

	entry = &filer.Entry{
		FullPath: fullpath,
	}
	err = entry.DecodeAttributesAndChunks(s3_util.MaybeDecompressData(data))
	if err != nil {
		return entry, fmt.Errorf("decode %s : %v", entry.FullPath, err)
	}

	return entry, nil
}

func (store *LuxDBStore) DeleteEntry(ctx context.Context, fullpath s3_util.FullPath) (err error) {
	key := genKey(fullpath.DirAndName())

	if err = store.db.Delete(key); err != nil {
		return fmt.Errorf("delete %s : %v", fullpath, err)
	}

	return nil
}

func (store *LuxDBStore) DeleteFolderChildren(ctx context.Context, fullpath s3_util.FullPath) (err error) {

	batch := store.db.NewBatch()

	directoryPrefix := genDirectoryKeyPrefix(fullpath, "")
	iter := store.db.NewIteratorWithStart(directoryPrefix)
	for iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, directoryPrefix) {
			break
		}
		fileName := getNameFromKey(key)
		if fileName == "" {
			continue
		}
		if err := batch.Delete(genKey(string(fullpath), fileName)); err != nil {
			iter.Release()
			return fmt.Errorf("batch delete %s : %v", fullpath, err)
		}
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return fmt.Errorf("iterate %s : %v", fullpath, err)
	}

	if err = batch.Write(); err != nil {
		return fmt.Errorf("delete %s : %v", fullpath, err)
	}

	return nil
}

func (store *LuxDBStore) ListDirectoryEntries(ctx context.Context, dirPath s3_util.FullPath, startFileName string, includeStartFile bool, limit int64, eachEntryFunc filer.ListEachEntryFunc) (lastFileName string, err error) {
	return store.ListDirectoryPrefixedEntries(ctx, dirPath, startFileName, includeStartFile, limit, "", eachEntryFunc)
}

func (store *LuxDBStore) ListDirectoryPrefixedEntries(ctx context.Context, dirPath s3_util.FullPath, startFileName string, includeStartFile bool, limit int64, prefix string, eachEntryFunc filer.ListEachEntryFunc) (lastFileName string, err error) {

	directoryPrefix := genDirectoryKeyPrefix(dirPath, prefix)
	lastFileStart := directoryPrefix
	if startFileName != "" {
		lastFileStart = genDirectoryKeyPrefix(dirPath, startFileName)
	}

	iter := store.db.NewIteratorWithStart(lastFileStart)
	for iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, directoryPrefix) {
			break
		}
		fileName := getNameFromKey(key)
		if fileName == "" {
			continue
		}
		if fileName == startFileName && !includeStartFile {
			continue
		}
		limit--
		if limit < 0 {
			break
		}
		lastFileName = fileName
		entry := &filer.Entry{
			FullPath: s3_util.NewFullPath(string(dirPath), fileName),
		}
		if decodeErr := entry.DecodeAttributesAndChunks(s3_util.MaybeDecompressData(iter.Value())); decodeErr != nil {
			err = decodeErr
			glog.V(0).InfofCtx(ctx, "list %s : %v", entry.FullPath, err)
			break
		}

		resEachEntryFunc, resEachEntryFuncErr := eachEntryFunc(entry)
		if resEachEntryFuncErr != nil {
			err = fmt.Errorf("failed to process eachEntryFunc: %w", resEachEntryFuncErr)
			break
		}

		if !resEachEntryFunc {
			break
		}
	}
	iter.Release()
	if err == nil {
		if iterErr := iter.Error(); iterErr != nil {
			err = fmt.Errorf("iterate %s : %v", dirPath, iterErr)
		}
	}

	return lastFileName, err
}

func genKey(dirPath, fileName string) (key []byte) {
	key = []byte(dirPath)
	key = append(key, DIR_FILE_SEPARATOR)
	key = append(key, []byte(fileName)...)
	return key
}

func genDirectoryKeyPrefix(fullpath s3_util.FullPath, startFileName string) (keyPrefix []byte) {
	keyPrefix = []byte(string(fullpath))
	keyPrefix = append(keyPrefix, DIR_FILE_SEPARATOR)
	if len(startFileName) > 0 {
		keyPrefix = append(keyPrefix, []byte(startFileName)...)
	}
	return keyPrefix
}

func getNameFromKey(key []byte) string {

	sepIndex := len(key) - 1
	for sepIndex >= 0 && key[sepIndex] != DIR_FILE_SEPARATOR {
		sepIndex--
	}

	return string(key[sepIndex+1:])
}

func (store *LuxDBStore) Shutdown() {
	if store.db != nil {
		store.db.Close()
	}
}

func (store *LuxDBStore) Debug(writer io.Writer) {
	iter := store.db.NewIterator()
	for iter.Next() {
		key := iter.Key()
		fullName := bytes.Replace(key, []byte{DIR_FILE_SEPARATOR}, []byte{' '}, 1)
		fmt.Fprintf(writer, "%v\n", string(fullName))
	}
	iter.Release()
}
