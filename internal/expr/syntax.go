package expr

import "github.com/laacin/inyorm/internal/core"

type ExprKind int

const (
	// Literals
	ExprKindString ExprKind = iota
	ExprKindNumber
	ExprKindFloat
	ExprKindBool
	ExprKindNull

	// Specials
	ExprKindPlaceholder
	ExprKindCond
	ExprKindConcat
	ExprKindCaseSwitch
	ExprKindCaseSearch

	// SQL Values
	ExprKindTable
	ExprKindCol
)

type ExprWriter interface {
	// Literals
	WriteLitString(core.Writer, string)
	WriteLitInt(core.Writer, int)
	WriteLitFloat(core.Writer, float64)
	WriteLitBool(core.Writer, bool)
	WriteLitNull(core.Writer)

	// Specials
	WriteExprPlaceholder(core.Writer, *Placeholder)
	WriteExprConcat(core.Writer, *Concat, core.WritingMode)
	WriteExprCond(core.Writer, *Cond, core.WritingMode)
	WriteExprCaseSwitch(core.Writer, *Case, core.WritingMode)
	WriteExprCaseSearch(core.Writer, *Case, core.WritingMode)

	// Table
	WriteExprTable(core.Writer, *Table)

	// Column
	WriteExprColBase(core.Writer, *Col)
	WriteExprColExpr(core.Writer, *Col)
	WriteExprColAlias(core.Writer, *Col)
	WriteExprColDef(core.Writer, *Col)

	WriteExprColAggr(core.Writer, *ColAggr)
	WriteExprColScalar(core.Writer, *ColScalar)
	WriteExprColArith(core.Writer, *ColArith)
	WriteExprColWrap(core.Writer)
}

type Expr interface {
	Kind() ExprKind
	Render(core.InternalWriter, ExprWriter, core.WritingMode)
}
