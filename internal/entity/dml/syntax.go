package dml

import "github.com/laacin/inyorm/internal/entity/core"

type Syntax interface {
	ClauseSyntax
	StatementOrder
}

type ClauseSyntax interface {
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

type StatementOrder interface {
	SelectOrder() []ClauseKind
	InsertOrder() []ClauseKind
	UpdateOrder() []ClauseKind
	DeleteOrder() []ClauseKind
}
