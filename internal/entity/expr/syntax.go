package expr

import "github.com/laacin/inyorm/internal/entity/core"

type Syntax interface {
	ValueSyntax
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
