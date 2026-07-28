//go:build (linux || darwin || windows) && sqlite

// limited GOOS due to modernc.org/libc/unistd, which hanzoai/sqlite's pure-Go
// backend sits on.

package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hanzoai/s3/s3/filer"
	"github.com/hanzoai/s3/s3/filer/abstract_sql"
	"github.com/hanzoai/s3/s3/util"

	// github.com/hanzoai/sqlite is the ONE Hanzo SQLite driver, and it is the
	// facade — never an engine directly. It registers the "sqlite" driver name
	// this store opens, and with it the hanzoai/sqlcipher codec VFS, so an
	// encrypted store opens on the same path as a plain one. Importing
	// modernc.org/sqlite here instead skipped the codec entirely and pinned a
	// second engine into go.mod alongside the one the rest of the fleet links.
	_ "github.com/hanzoai/sqlite"
)

func init() {
	filer.Stores = append(filer.Stores, &SqliteStore{})
}

type SqliteStore struct {
	abstract_sql.AbstractSqlStore
}

func (store *SqliteStore) GetName() string {
	return "sqlite"
}

func (store *SqliteStore) Initialize(configuration util.Configuration, prefix string) (err error) {
	dbFile := configuration.GetString(prefix + "dbFile")
	createTable := `CREATE TABLE IF NOT EXISTS "%s" (
		dirhash BIGINT,
		name VARCHAR(1000),
		directory TEXT,
		meta BLOB,
		PRIMARY KEY (dirhash, name)
	) WITHOUT ROWID;`
	upsertQuery := `INSERT INTO "%s"(dirhash,name,directory,meta)VALUES(?,?,?,?)
	ON CONFLICT(dirhash,name) DO UPDATE SET
		directory=excluded.directory,
		meta=excluded.meta;
	`
	return store.initialize(
		dbFile,
		createTable,
		upsertQuery,
	)
}

func (store *SqliteStore) initialize(dbFile, createTable, upsertQuery string) (err error) {

	store.SupportBucketTable = true
	store.SqlGenerator = &SqlGenSqlite{
		CreateTableSqlTemplate: createTable,
		DropTableSqlTemplate:   "drop table `%s`",
		UpsertQueryTemplate:    upsertQuery,
	}

	var dbErr error
	store.DB, dbErr = sql.Open("sqlite", dbFile)
	if dbErr != nil {
		if store.DB != nil {
			store.DB.Close()
			store.DB = nil
		}
		return fmt.Errorf("can not connect to %s error:%v", dbFile, dbErr)
	}

	if err = store.DB.Ping(); err != nil {
		return fmt.Errorf("connect to %s error:%v", dbFile, err)
	}

	store.DB.SetMaxOpenConns(1)

	if err = store.CreateTable(context.Background(), abstract_sql.DEFAULT_TABLE); err != nil {
		return fmt.Errorf("init table %s: %v", abstract_sql.DEFAULT_TABLE, err)
	}

	return nil
}
