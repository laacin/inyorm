package dml

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/expr"
)

// --- Clauses

type Select struct {
	Distinct bool
	Values   []any
}

type From struct {
	Value any
}

type Join struct {
	Joins []JoinSegment
}

type Where struct {
	Conds []*expr.Condition
}

type GroupBy struct {
	Values []any
}

type Having struct {
	Cond *expr.Condition
}

type OrderBy struct {
	Orders []OrderSegment
}

type Limit struct {
	ValueNumber int
}

type Offset struct {
	ValueNumber int
}

type InsertInto struct {
	Table  any
	Cols   []string
	Rows   int
	Values []any
}

type Update struct {
	Table  any
	Cols   []string
	Values []any
}

type Delete struct{}

// --- Kinds

func (c *Select) Kind() ClauseKind     { return ClauseSelect }
func (c *From) Kind() ClauseKind       { return ClauseFrom }
func (c *Join) Kind() ClauseKind       { return ClauseJoin }
func (c *Where) Kind() ClauseKind      { return ClauseWhere }
func (c *GroupBy) Kind() ClauseKind    { return ClauseGroupBy }
func (c *Having) Kind() ClauseKind     { return ClauseHaving }
func (c *OrderBy) Kind() ClauseKind    { return ClauseOrderBy }
func (c *Limit) Kind() ClauseKind      { return ClauseLimit }
func (c *Offset) Kind() ClauseKind     { return ClauseOffset }
func (c *InsertInto) Kind() ClauseKind { return ClauseInsertInto }
func (c *Update) Kind() ClauseKind     { return ClauseUpdate }
func (c *Delete) Kind() ClauseKind     { return ClauseDelete }

// --- Writes

func (c *Select) Write(w core.InternalWriter, dial ClauseSyntax)     { dial.WriteSelect(w, c) }
func (c *From) Write(w core.InternalWriter, dial ClauseSyntax)       { dial.WriteFrom(w, c) }
func (c *Join) Write(w core.InternalWriter, dial ClauseSyntax)       { dial.WriteJoin(w, c) }
func (c *Where) Write(w core.InternalWriter, dial ClauseSyntax)      { dial.WriteWhere(w, c) }
func (c *GroupBy) Write(w core.InternalWriter, dial ClauseSyntax)    { dial.WriteGroupBy(w, c) }
func (c *Having) Write(w core.InternalWriter, dial ClauseSyntax)     { dial.WriteHaving(w, c) }
func (c *OrderBy) Write(w core.InternalWriter, dial ClauseSyntax)    { dial.WriteOrderBy(w, c) }
func (c *Limit) Write(w core.InternalWriter, dial ClauseSyntax)      { dial.WriteLimit(w, c) }
func (c *Offset) Write(w core.InternalWriter, dial ClauseSyntax)     { dial.WriteOffset(w, c) }
func (c *InsertInto) Write(w core.InternalWriter, dial ClauseSyntax) { dial.WriteInsertInto(w, c) }
func (c *Update) Write(w core.InternalWriter, dial ClauseSyntax)     { dial.WriteUpdate(w, c) }
func (c *Delete) Write(w core.InternalWriter, dial ClauseSyntax)     { dial.WriteDelete(w, c) }

// --- Utilities

type JoinType int

const (
	JoinInner JoinType = iota
	JoinLeft
	JoinRight
	JoinFull
	JoinCross
)

type JoinSegment struct {
	Type  JoinType
	Table any
	Cond  *expr.Condition
}

type OrderSegment struct {
	Value      any
	Descending bool
}
