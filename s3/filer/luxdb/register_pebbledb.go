//go:build pebbledb

// luxfi/database registers leveldb in its factory under -tags=leveldb but ships
// no equivalent registration for pebbledb (factory/leveldb.go has no
// factory/pebbledb.go sibling). To make backend="pebble" selectable here, we
// register pebbledb in the factory under -tags=pebbledb, mirroring factory's own
// leveldb registration. zapdb (the default) needs no registration; it is always
// available.
package luxdb

import (
	"github.com/luxfi/database"
	"github.com/luxfi/database/factory"
	"github.com/luxfi/database/pebbledb"
	log "github.com/luxfi/log"
	"github.com/luxfi/metric"
)

func init() {
	factory.RegisterDatabase(pebbledb.Name, newPebbleDB)
}

func newPebbleDB(
	dbPath string,
	config []byte,
	logger log.Logger,
	registerer metric.Registerer,
	metricsPrefix string,
	readOnly bool,
) (database.Database, error) {
	// Defaults consistent with factory's leveldb sizing.
	cacheSize := 12 // MiB block cache
	handles := 1024
	return pebbledb.New(dbPath, cacheSize, handles, pebbledb.Name, readOnly)
}
