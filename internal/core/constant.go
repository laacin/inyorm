package core

// Dialects
const (
	Postgres = "postgres"
	MySQL    = "mysql"
	SQLite   = "sqlite3"
)

// ----- Column

type ColumnType int

const (
	ColTypUnset ColumnType = iota
	ColTypDef
	ColTypBase
	ColTypExpr
	ColTypAlias
)
