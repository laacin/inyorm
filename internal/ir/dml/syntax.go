package dml

import "github.com/laacin/inyorm/internal/core"

type Syntax interface {
	ClauseWriter
	QueryOrder
}

type ClauseWriter interface {
	WriteSelect(core.Writer, *Select)
	WriteFrom(core.Writer, *From)
	WriteJoin(core.Writer, *Join)
	WriteWhere(core.Writer, *Where)
	WriteGroupBy(core.Writer, *GroupBy)
	WriteHaving(core.Writer, *Having)
	WriteOrderBy(core.Writer, *OrderBy)
	WriteLimit(core.Writer, *Limit)
	WriteOffset(core.Writer, *Offset)

	WriteInsertInto(core.Writer, *InsertInto)

	WriteUpdate(core.Writer, *Update)

	WriteDelete(core.Writer, *Delete)
}

type QueryOrder interface {
	SelectOrder() []ClauseKind
	InsertOrder() []ClauseKind
	UpdateOrder() []ClauseKind
	DeleteOrder() []ClauseKind
}

// --- Internal
type ClauseBuilder interface {
	IsDeclared() bool
	Kind() ClauseKind
	Build(core.InternalWriter, ClauseWriter) error
}

type QueryBuilder interface {
	Kind() QueryKind
	Build() (string, []any, error)
}
