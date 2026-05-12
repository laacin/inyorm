package dml

import "github.com/laacin/inyorm/internal/entity/core"

type Dialect interface {
	ValueSyntax
	ClauseSyntax
	StatementOrder
}

type ValueSyntax interface {
	// Literals
	WriteString(core.Writer, string)
	WriteNumber(core.Writer, int)
	WriteFloat(core.Writer, float64)
	WriteBool(core.Writer, bool)
	WriteNull(core.Writer)
	WriteWildcard(core.Writer)

	// Specials
	WritePlaceholder(core.Writer, int)
	WriteConcat(core.Writer, *Concat)
	WriteCondition(core.Writer, *Condition, core.WritingMode)
	WriteCaseSwitch(core.Writer, *CaseSwitch, core.WritingMode)
	WriteCaseSearch(core.Writer, *CaseSearch, core.WritingMode)

	// Table
	WriteTable(core.Writer, *Table)

	// Column
	WriteColBase(core.Writer, *Column)
	WriteColExpr(core.Writer, *Column)
	WriteColAlias(core.Writer, *Column)
	WriteColDef(core.Writer, *Column)
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
