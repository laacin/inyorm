package dml

import "github.com/laacin/inyorm/internal/ir/expr"

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
	Conds []expr.ExprBuilder
}

type GroupBy struct {
	Values []any
}

type Having struct {
	Cond expr.ExprBuilder
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
	Cond  expr.ExprBuilder
}

type OrderSegment struct {
	Value      any
	Descending bool
}
