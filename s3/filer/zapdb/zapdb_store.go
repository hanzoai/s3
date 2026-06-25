package zapdb

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	zapdb "github.com/luxfi/zapdb"

	"github.com/hanzoai/s3/s3/filer"
	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/pb/filer_pb"
	s3_util "github.com/hanzoai/s3/s3/util"
)

const (
	// DIR_FILE_SEPARATOR delimits the directory prefix from the file name in a
	// key. Identical to the leveldb store's separator so the on-disk key layout
	// is byte-for-byte compatible.
	DIR_FILE_SEPARATOR = byte(0x00)
)

var (
	_ = filer.Debuggable(&ZapDBStore{})
)

func init() {
	filer.Stores = append(filer.Stores, &ZapDBStore{})
}

// ZapDBStore is a filer.FilerStore backed by zapdb (a transactional,
// LSM-tree KV store). It mirrors the leveldb store's key encoding exactly,
// so data written by either store shares an identical key layout.
type ZapDBStore struct {
	db *zapdb.DB
}

func (store *ZapDBStore) GetName() string {
	return "zapdb"
}

func (store *ZapDBStore) Initialize(configuration s3_util.Configuration, prefix string) (err error) {
	dir := configuration.GetString(prefix + "dir")
	return store.initialize(dir)
}

func (store *ZapDBStore) initialize(dir string) (err error) {
	glog.V(0).Infof("filer store dir: %s", dir)
	os.MkdirAll(dir, 0755)
	if err := s3_util.TestFolderWritable(dir); err != nil {
		return fmt.Errorf("Check Level Folder %s Writable: %s", dir, err)
	}

	opts := zapdb.DefaultOptions(dir).
		WithLogger(nil).
		WithSyncWrites(false)

	if store.db, err = zapdb.Open(opts); err != nil {
		glog.Infof("filer store open dir %s: %v", dir, err)
		return
	}
	return
}

func (store *ZapDBStore) BeginTransaction(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
func (store *ZapDBStore) CommitTransaction(ctx context.Context) error {
	return nil
}
func (store *ZapDBStore) RollbackTransaction(ctx context.Context) error {
	return nil
}

func (store *ZapDBStore) InsertEntry(ctx context.Context, entry *filer.Entry) (err error) {
	key := genKey(entry.DirAndName())

	value, err := entry.EncodeAttributesAndChunks()
	if err != nil {
		return fmt.Errorf("encoding %s %+v: %v", entry.FullPath, entry.Attr, err)
	}

	if len(entry.GetChunks()) > filer.CountEntryChunksForGzip {
		value = s3_util.MaybeGzipData(value)
	}

	err = store.db.Update(func(txn *zapdb.Txn) error {
		return txn.Set(key, value)
	})

	if err != nil {
		return fmt.Errorf("persisting %s : %v", entry.FullPath, err)
	}

	return nil
}

func (store *ZapDBStore) UpdateEntry(ctx context.Context, entry *filer.Entry) (err error) {

	return store.InsertEntry(ctx, entry)
}

// BatchInsertEntries inserts multiple entries in a single zapdb transaction.
// This is more efficient than inserting entries one by one as it reduces
// the number of write operations and syncs to disk.
func (store *ZapDBStore) BatchInsertEntries(ctx context.Context, entries []*filer.Entry) error {
	if len(entries) == 0 {
		return nil
	}

	err := store.db.Update(func(txn *zapdb.Txn) error {
		for _, entry := range entries {
			key := genKey(entry.DirAndName())

			value, err := entry.EncodeAttributesAndChunks()
			if err != nil {
				return fmt.Errorf("encoding %s %+v: %w", entry.FullPath, entry.Attr, err)
			}

			if len(entry.GetChunks()) > filer.CountEntryChunksForGzip {
				value = s3_util.MaybeGzipData(value)
			}

			if err := txn.Set(key, value); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("batch write: %w", err)
	}

	return nil
}

func (store *ZapDBStore) FindEntry(ctx context.Context, fullpath s3_util.FullPath) (entry *filer.Entry, err error) {
	key := genKey(fullpath.DirAndName())

	var data []byte
	err = store.db.View(func(txn *zapdb.Txn) error {
		item, getErr := txn.Get(key)
		if getErr != nil {
			return getErr
		}
		data, getErr = item.ValueCopy(nil)
		return getErr
	})

	if err == zapdb.ErrKeyNotFound {
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

func (store *ZapDBStore) DeleteEntry(ctx context.Context, fullpath s3_util.FullPath) (err error) {
	key := genKey(fullpath.DirAndName())

	err = store.db.Update(func(txn *zapdb.Txn) error {
		return txn.Delete(key)
	})
	if err != nil {
		return fmt.Errorf("delete %s : %v", fullpath, err)
	}

	return nil
}

func (store *ZapDBStore) DeleteFolderChildren(ctx context.Context, fullpath s3_util.FullPath) (err error) {

	directoryPrefix := genDirectoryKeyPrefix(fullpath, "")

	// Collect the child names under a read-only iterator first, then delete in
	// a separate write transaction. zapdb forbids mutating keys while iterating
	// within the same transaction, so the two phases must not share a txn.
	var fileNames []string
	err = store.db.View(func(txn *zapdb.Txn) error {
		opts := zapdb.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = directoryPrefix
		iter := txn.NewIterator(opts)
		defer iter.Close()
		for iter.Seek(directoryPrefix); iter.Valid(); iter.Next() {
			key := iter.Item().Key()
			if !bytes.HasPrefix(key, directoryPrefix) {
				break
			}
			fileName := getNameFromKey(key)
			if fileName == "" {
				continue
			}
			fileNames = append(fileNames, fileName)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("delete %s : %v", fullpath, err)
	}

	err = store.db.Update(func(txn *zapdb.Txn) error {
		for _, fileName := range fileNames {
			if err := txn.Delete(genKey(string(fullpath), fileName)); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("delete %s : %v", fullpath, err)
	}

	return nil
}

func (store *ZapDBStore) ListDirectoryEntries(ctx context.Context, dirPath s3_util.FullPath, startFileName string, includeStartFile bool, limit int64, eachEntryFunc filer.ListEachEntryFunc) (lastFileName string, err error) {
	return store.ListDirectoryPrefixedEntries(ctx, dirPath, startFileName, includeStartFile, limit, "", eachEntryFunc)
}

func (store *ZapDBStore) ListDirectoryPrefixedEntries(ctx context.Context, dirPath s3_util.FullPath, startFileName string, includeStartFile bool, limit int64, prefix string, eachEntryFunc filer.ListEachEntryFunc) (lastFileName string, err error) {

	directoryPrefix := genDirectoryKeyPrefix(dirPath, prefix)
	lastFileStart := directoryPrefix
	if startFileName != "" {
		lastFileStart = genDirectoryKeyPrefix(dirPath, startFileName)
	}

	err = store.db.View(func(txn *zapdb.Txn) error {
		opts := zapdb.DefaultIteratorOptions
		opts.Prefix = directoryPrefix
		iter := txn.NewIterator(opts)
		defer iter.Close()
		for iter.Seek(lastFileStart); iter.Valid(); iter.Next() {
			item := iter.Item()
			key := item.Key()
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
			value, valueErr := item.ValueCopy(nil)
			if valueErr != nil {
				return valueErr
			}
			entry := &filer.Entry{
				FullPath: s3_util.NewFullPath(string(dirPath), fileName),
			}
			if decodeErr := entry.DecodeAttributesAndChunks(s3_util.MaybeDecompressData(value)); decodeErr != nil {
				glog.V(0).InfofCtx(ctx, "list %s : %v", entry.FullPath, decodeErr)
				return decodeErr
			}

			resEachEntryFunc, resEachEntryFuncErr := eachEntryFunc(entry)
			if resEachEntryFuncErr != nil {
				return fmt.Errorf("failed to process eachEntryFunc: %w", resEachEntryFuncErr)
			}

			if !resEachEntryFunc {
				break
			}
		}
		return nil
	})

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

func (store *ZapDBStore) Shutdown() {
	store.db.Close()
}

func (store *ZapDBStore) Debug(writer io.Writer) {
	store.db.View(func(txn *zapdb.Txn) error {
		opts := zapdb.DefaultIteratorOptions
		opts.PrefetchValues = false
		iter := txn.NewIterator(opts)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			key := iter.Item().Key()
			fullName := bytes.Replace(key, []byte{DIR_FILE_SEPARATOR}, []byte{' '}, 1)
			fmt.Fprintf(writer, "%v\n", string(fullName))
		}
		return nil
	})
}
