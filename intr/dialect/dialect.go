package dialect

// --- Main Dialect interface
type Dialect interface {
	ValueBuilder
	ColumnBuilder
	ClauseBuilder
}

// Writer used by a dialect
type Writer interface {
	Write(string)
	Char(byte)

	// Writes any value based on clause context
	Value(v any, mode WritingMode) // writing mode used if is required
	GetTableRef(string) (ref byte, shouldBeUsed bool)

	Result() string
	Reset()
}

// Dialect essentials (Builders)
type InternalBuilder interface {
	// Writes a condition (must be wrapped)
	Cond(Writer, Cond)
}

type ValueBuilder interface {
	String(Writer, string)
	Number(Writer, int)
	Float(Writer, float64)
	Bool(Writer, bool)
	Null(Writer)

	// Writes a placeholder based on position
	Placeholder(Writer, int)
}

type ColumnBuilder interface {
	Table(Writer, Table)

	// Writing modes
	ColDef(Writer, Column)
	ColAlias(Writer, Column)
	ColExpr(Writer, Column)
	ColBase(Writer, Column)

	// Specials columns
	ColWildcard(Writer, Table)
	ColConcat(Writer, []any)
	ColSwitch(w Writer, cond any, cas CaseCond)
	ColSearch(Writer, CaseCond) // Case identifier must be a Cond
}

type ClauseBuilder interface {
	// Insert Statement
	ClsInsertInto(Writer, InsertIntoTools)

	// Select Statement
	ClsSelect(Writer, SelectTools)
	ClsFrom(Writer, FromTools)
	ClsJoin(Writer, []JoinTools)
	ClsWhere(Writer, WhereTools)
	ClsGroupBy(Writer, GroupByTools)
	ClsHaving(Writer, HavingTools)
	ClsOrderBy(Writer, []OrderByTools)
	ClsLimit(Writer, LimitTools)
	ClsOffset(Writer, OffsetTools)

	// Update Statement
	ClsUpdate(Writer, UpdateTools)

	// Delete Statement
	ClsDelete(Writer, DeleteTools)
}

// Writing modes
type WritingMode int

const (
	WriteDef WritingMode = iota
	WriteBase
	WriteAlias // Column only
	WriteExpr  // Column only
)

type WritingModeConfig struct {
	Default WritingMode

	InsertInto WritingMode

	Select  WritingMode
	From    WritingMode
	Join    WritingMode
	Where   WritingMode
	GroupBy WritingMode
	Having  WritingMode
	OrderBy WritingMode
	Limit   WritingMode
	Offset  WritingMode

	Update WritingMode

	Delete WritingMode
}
