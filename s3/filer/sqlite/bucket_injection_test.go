//go:build sqlite

package sqlite

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

// newProdStore builds a SqliteStore exactly the way Initialize() does in
// production: same CREATE TABLE / upsert templates, SupportBucketTable=true.
func newProdStore(t *testing.T) *SqliteStore {
	t.Helper()
	store := &SqliteStore{}
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
	if err := store.initialize(filepath.Join(t.TempDir(), "filer.db"), createTable, upsertQuery); err != nil {
		t.Fatalf("initialize store: %v", err)
	}
	return store
}

func tableExists(t *testing.T, store *SqliteStore, name string) bool {
	t.Helper()
	var got string
	err := store.DB.QueryRowContext(context.Background(),
		"SELECT name FROM sqlite_master WHERE type='table' AND name=?", name).Scan(&got)
	return err == nil && got == name
}

// TestOnBucketDeletionInjection drives the production OnBucketDeletion metadata
// handler with a bucket name that breaks out of the DROP TABLE `%s` identifier.
// The bucket name reaches this path from a filer metadata event
// (message.OldEntry.Name for a directory under /buckets/); a peer filer or the
// unvalidated filer mkdir/gRPC create path can supply an arbitrary name.
//
// The emitted SQL is logged. The assertion is behavioral: a bystander "victim"
// table must SURVIVE. Before the fix the stacked DROP runs and this FAILS (red);
// after the fix the invalid name is refused before any SQL runs (green).
func TestOnBucketDeletionInjection(t *testing.T) {
	store := newProdStore(t)
	ctx := context.Background()

	// A legitimately-created bucket table, so the first stacked DROP has a real
	// target and the injected second statement is reached.
	store.OnBucketCreation("realbucket")

	if _, err := store.DB.ExecContext(ctx, "CREATE TABLE victim(secret TEXT)"); err != nil {
		t.Fatalf("seed victim: %v", err)
	}

	payload := "realbucket`;drop table `victim"
	t.Logf("EMITTED DROP SQL: %s", store.SqlGenerator.GetSqlDropTable(payload))

	store.OnBucketDeletion(payload)

	if !tableExists(t, store, "victim") {
		t.Fatalf("victim table was dropped: SQL injection via OnBucketDeletion executed")
	}
}

// TestOnBucketCreationInjection is the CREATE TABLE `%s` analogue.
func TestOnBucketCreationInjection(t *testing.T) {
	store := newProdStore(t)
	ctx := context.Background()

	if _, err := store.DB.ExecContext(ctx, "CREATE TABLE victim2(secret TEXT)"); err != nil {
		t.Fatalf("seed victim2: %v", err)
	}

	payload := `x"(a int);DROP TABLE victim2;--`
	t.Logf("EMITTED CREATE SQL: %s", store.SqlGenerator.GetSqlCreateTable(payload))

	store.OnBucketCreation(payload)

	if !tableExists(t, store, "victim2") {
		t.Fatalf("victim2 table was dropped: SQL injection via OnBucketCreation executed")
	}
}

// TestOnBucketLifecycleValidNameStillWorks guards against an over-broad fix:
// a well-formed bucket name must still create and drop its per-bucket table.
func TestOnBucketLifecycleValidNameStillWorks(t *testing.T) {
	store := newProdStore(t)

	store.OnBucketCreation("my-valid-bucket")
	if !tableExists(t, store, "my-valid-bucket") {
		t.Fatalf("valid bucket table was not created")
	}

	store.OnBucketDeletion("my-valid-bucket")
	if tableExists(t, store, "my-valid-bucket") {
		t.Fatalf("valid bucket table was not dropped")
	}

	// A name that must never be treated as a bucket table.
	if strings.EqualFold("filemeta", "FILEMETA") {
		store.OnBucketDeletion("filemeta")
		if !tableExists(t, store, "filemeta") {
			t.Fatalf("default filemeta table must never be droppable via OnBucketDeletion")
		}
	}
}
