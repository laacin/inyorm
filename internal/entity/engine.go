package entity

import "context"

type Engine struct {
	Dialect Dialect
	Driver  Driver
	Err     error
}

// ---- Dialect

type Dialect interface {
	ValueSyntax
	ClauseSyntax
	StatementOrder
}

type ValueSyntax interface {
	// Literals
	WriteString(Writer, string)
	WriteNumber(Writer, int)
	WriteFloat(Writer, float64)
	WriteBool(Writer, bool)
	WriteNull(Writer)
	WriteWildcard(Writer)

	// Specials
	WritePlaceholder(Writer, int)
	WriteConcat(Writer, *Concat)
	WriteCondition(Writer, *Condition, WritingMode)
	WriteCaseSwitch(Writer, *CaseSwitch, WritingMode)
	WriteCaseSearch(Writer, *CaseSearch, WritingMode)

	// Table
	WriteTable(Writer, *Table)

	// Column
	WriteColBase(Writer, *Column)
	WriteColExpr(Writer, *Column)
	WriteColAlias(Writer, *Column)
	WriteColDef(Writer, *Column)
}

type ClauseSyntax interface {
	WriteSelect(Writer, *Select)
	WriteFrom(Writer, *From)
	WriteJoin(Writer, *Join)
	WriteWhere(Writer, *Where)
	WriteGroupBy(Writer, *GroupBy)
	WriteHaving(Writer, *Having)
	WriteOrderBy(Writer, *OrderBy)
	WriteLimit(Writer, *Limit)
	WriteOffset(Writer, *Offset)

	WriteInsertInto(Writer, *InsertInto)
	WriteUpdate(Writer, *Update)
	WriteDelete(Writer, *Delete)
}

type StatementOrder interface {
	SelectOrder() []ClauseKind
	InsertOrder() []ClauseKind
	UpdateOrder() []ClauseKind
	DeleteOrder() []ClauseKind
}

// --- Driver

type Driver interface {
	Connection
	Executor
}

type Connection interface {
	Close() error
}

type Executor interface {
	Exec(context.Context, string, ...any) error
	Query(context.Context, string, ...any) (Rows, error)
}

// dependencies

type Rows interface {
	Columns() ([]string, error)
	Next() bool
	Scan(...any) error
}
