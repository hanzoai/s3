//go:build (linux || darwin || windows) && sqlite

package sqlite

import (
	"fmt"

	"github.com/hanzoai/s3/s3/filer/abstract_sql"
)

// SqlGenSqlite generates the SQL dialect used by the SQLite filer store.
// SQLite accepts the same backtick-quoted, positional-parameter SQL that the
// former shared MySQL generator emitted, so the sqlite store owns this small
// generator outright rather than importing a deleted sibling backend.
type SqlGenSqlite struct {
	CreateTableSqlTemplate string
	DropTableSqlTemplate   string
	UpsertQueryTemplate    string
	// ForceBinaryCollation forces byte ordering on a non-binary name column;
	// see GetSqlListExclusive / GetSqlListInclusive.
	ForceBinaryCollation bool
}

var (
	_ = abstract_sql.SqlGenerator(&SqlGenSqlite{})
)

func (gen *SqlGenSqlite) GetSqlInsert(tableName string) string {
	if gen.UpsertQueryTemplate != "" {
		return fmt.Sprintf(gen.UpsertQueryTemplate, tableName)
	} else {
		return fmt.Sprintf("INSERT INTO `%s` (`dirhash`,`name`,`directory`,`meta`) VALUES(?,?,?,?)", tableName)
	}
}

func (gen *SqlGenSqlite) GetSqlUpdate(tableName string) string {
	return fmt.Sprintf("UPDATE `%s` SET `meta` = ? WHERE `dirhash` = ? AND `name` = ? AND `directory` = ?", tableName)
}

func (gen *SqlGenSqlite) GetSqlFind(tableName string) string {
	return fmt.Sprintf("SELECT `meta` FROM `%s` WHERE `dirhash` = ? AND `name` = ? AND `directory` = ?", tableName)
}

func (gen *SqlGenSqlite) GetSqlDelete(tableName string) string {
	return fmt.Sprintf("DELETE FROM `%s` WHERE `dirhash` = ? AND `name` = ? AND `directory` = ?", tableName)
}

func (gen *SqlGenSqlite) GetSqlDeleteFolderChildren(tableName string) string {
	return fmt.Sprintf("DELETE FROM `%s` WHERE `dirhash` = ? AND `directory` = ?", tableName)
}

// nameExpr forces byte ordering on a non-binary column at the cost of a filesort.
func (gen *SqlGenSqlite) nameExpr() string {
	if gen.ForceBinaryCollation {
		return "BINARY `name`"
	}
	return "`name`"
}

func (gen *SqlGenSqlite) GetSqlListExclusive(tableName string) string {
	name := gen.nameExpr()
	return fmt.Sprintf("SELECT `name`, `meta` FROM `%s` WHERE `dirhash` = ? AND %s > ? AND `directory` = ? AND %s LIKE ? ORDER BY %s ASC LIMIT ?", tableName, name, name, name)
}

func (gen *SqlGenSqlite) GetSqlListInclusive(tableName string) string {
	name := gen.nameExpr()
	return fmt.Sprintf("SELECT `name`, `meta` FROM `%s` WHERE `dirhash` = ? AND %s >= ? AND `directory` = ? AND %s LIKE ? ORDER BY %s ASC LIMIT ?", tableName, name, name, name)
}

func (gen *SqlGenSqlite) GetSqlCreateTable(tableName string) string {
	return fmt.Sprintf(gen.CreateTableSqlTemplate, tableName)
}

func (gen *SqlGenSqlite) GetSqlDropTable(tableName string) string {
	return fmt.Sprintf(gen.DropTableSqlTemplate, tableName)
}
