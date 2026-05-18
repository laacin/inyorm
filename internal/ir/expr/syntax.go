package expr

import "github.com/laacin/inyorm/internal/core"

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
