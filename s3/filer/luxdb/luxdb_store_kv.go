package luxdb

import (
	"context"
	"fmt"

	"github.com/luxfi/database"

	"github.com/hanzoai/s3/s3/filer"
)

func (store *LuxDBStore) KvPut(ctx context.Context, key []byte, value []byte) (err error) {

	if err = store.db.Put(key, value); err != nil {
		return fmt.Errorf("kv put: %w", err)
	}

	return nil
}

func (store *LuxDBStore) KvGet(ctx context.Context, key []byte) (value []byte, err error) {

	value, err = store.db.Get(key)

	if err == database.ErrNotFound {
		return nil, filer.ErrKvNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("kv get: %w", err)
	}

	return
}

func (store *LuxDBStore) KvDelete(ctx context.Context, key []byte) (err error) {

	if err = store.db.Delete(key); err != nil {
		return fmt.Errorf("kv delete: %w", err)
	}

	return nil
}
