package expr

import "github.com/laacin/inyorm/internal/core"

type ExprKind int

const (
	// Literals
	ExprString ExprKind = iota
	ExprNumber
	ExprFloat
	ExprBool
	ExprNull

	// Specials
	ExprParam
	ExprCond
	ExprConcat
	ExprCaseSwitch
	ExprCaseSearch

	// SQL Values
	ExprTable
	ExprCol
)

type ExprWriter interface {
	// Literals
	WriteString(core.Writer, string)
	WriteNumber(core.Writer, int)
	WriteFloat(core.Writer, float64)
	WriteBool(core.Writer, bool)
	WriteNull(core.Writer)

	// Specials
	WritePlaceholder(core.Writer)
	WriteConcat(core.Writer, *Concat, core.WritingMode)
	WriteCond(core.Writer, *Cond, core.WritingMode)
	WriteCaseSwitch(core.Writer, *Case, core.WritingMode)
	WriteCaseSearch(core.Writer, *Case, core.WritingMode)

	// Table
	WriteTable(core.Writer, *Table)

	// Column
	WriteColBase(core.Writer, *Col)
	WriteColExpr(core.Writer, *Col)
	WriteColAlias(core.Writer, *Col)
	WriteColDef(core.Writer, *Col)

	WriteColAggr(core.Writer, *ColAggr)
	WriteColScalar(core.Writer, *ColScalar)
	WriteColArith(core.Writer, *ColArith)
	WriteColWrap(core.Writer)
}

// --- Internal
type ExprBuilder interface {
	Kind() ExprKind
	Build(core.InternalWriter, ExprWriter, core.WritingMode)
}
