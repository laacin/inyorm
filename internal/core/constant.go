package core

// Dialects
const (
	Postgres = "postgresql"
	MySQL    = "mysql"
	SQLite   = "sqlite"
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
)
