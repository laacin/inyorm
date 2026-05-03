package dialect

// --- Main Dialect interface
type Dialect interface {
	ValueBuilder
	ColumnBuilder
	ClauseBuilder
}

// Writer used by dialect
type Writer interface {
	Write(string)
	Char(byte)

	// Column(Column, WritingMode)
	Value(v any, ctx ClauseName)

	Result() string
	Reset()
}

// Dialect essentials
type ValueBuilder interface {
	String(Writer, string)
	Number(Writer, int)
	Float(Writer, float64)
	Bool(Writer, bool)
	Null(Writer)

	Placeholder(w Writer, num int)
	Cond(Writer, Cond) // must be wrapped
}

type ColumnBuilder interface {
	Table(w Writer, table Table, def bool)
	ColDef(Writer, Column)
	ColAlias(Writer, Column)
	ColExpr(Writer, Column)
	ColBase(Writer, Column)
}

type ClauseBuilder interface {
	ClsInsertInto(Writer, InsertIntoTools)

	ClsSelect(Writer, SelectTools)
	ClsFrom(Writer, FromTools)
	ClsJoin(Writer, []JoinTools)
	ClsWhere(Writer, WhereTools)
	ClsGroupBy(Writer, GroupByTools)
	ClsHaving(Writer, HavingTools)
	ClsOrderBy(Writer, []OrderByTools)
	ClsLimit(Writer, LimitTools)
	ClsOffset(Writer, OffsetTools)

	ClsUpdate(Writer, UpdateTools)
	ClsDelete(WhereTools, DeleteTools)
}

// Writing Mode
type WritingMode int

const (
	WriteDef WritingMode = iota
	WriteBase
	WriteAlias // Column only
	WriteExpr  // Column only
)

type WritingModeConfig struct {
	None WritingMode

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
