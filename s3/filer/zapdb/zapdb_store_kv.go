package zapdb

import (
	"context"
	"fmt"

	zapdb "github.com/luxfi/zapdb"

	"github.com/hanzoai/s3/s3/filer"
)

func (store *ZapDBStore) KvPut(ctx context.Context, key []byte, value []byte) (err error) {

	err = store.db.Update(func(txn *zapdb.Txn) error {
		return txn.Set(key, value)
	})

	if err != nil {
		return fmt.Errorf("kv put: %w", err)
	}

	return nil
}

func (store *ZapDBStore) KvGet(ctx context.Context, key []byte) (value []byte, err error) {

	err = store.db.View(func(txn *zapdb.Txn) error {
		item, getErr := txn.Get(key)
		if getErr != nil {
			return getErr
		}
		value, getErr = item.ValueCopy(nil)
		return getErr
	})

	if err == zapdb.ErrKeyNotFound {
		return nil, filer.ErrKvNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("kv get: %w", err)
	}

	return
}

func (store *ZapDBStore) KvDelete(ctx context.Context, key []byte) (err error) {

	err = store.db.Update(func(txn *zapdb.Txn) error {
		return txn.Delete(key)
	})

	if err != nil {
		return fmt.Errorf("kv delete: %w", err)
	}

	return nil
}
