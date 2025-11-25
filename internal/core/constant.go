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

// ----- Clauses

type ClauseType int

const (
	ClsTypUnset ClauseType = iota
	ClsTypSelect
	ClsTypFrom
	ClsTypJoin
	ClsTypWhere
	ClsTypGroupBy
	ClsTypHaving
	ClsTypOrderBy
	ClsTypLimit
	ClsTypOffset
	ClsTypInsert
	ClsTypUpdate
	ClsTypDelete
)
