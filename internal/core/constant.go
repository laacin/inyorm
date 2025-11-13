package core

// Dialects
const (
	Postgres  = "postgresql"
	MySQL     = "mysql"
	SQLite    = "sqlite"
	Oracle    = "oracle"
	SQLServer = "sqlserver"
)

// ----- Column

type ColumnType int

const (
	ColTypDef ColumnType = iota
	ColTypBase
	ColTypExpr
	ColTypAlias
)

// ----- Clauses

type ClauseType int

const (
	ClsTypSelect ClauseType = iota
	ClsTypFrom
	ClsTypJoin
	ClsTypWhere
	ClsTypGroupBy
	ClsTypHaving
	ClsTypOrderBy
	ClsTypLimit
	ClsTypOffset
	ClsTypInsertInto
)
