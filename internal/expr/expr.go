package expr

import "github.com/laacin/inyorm/internal/core"

type Expr interface {
	Kind() Kind
	Render(core.InternalWriter, Renderer, core.WritingMode)
}

type Renderer interface {
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

type Kind int

const (
	// Literals
	KindString Kind = iota
	KindNumber
	KindFloat
	KindBool
	KindNull

	// Specials
	KindPlaceholder
	KindCond
	KindConcat
	KindCaseSwitch
	KindCaseSearch

	// SQL Values
	KindTable
	KindCol
)
