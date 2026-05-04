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

	New() Writer
	Result() string
	Reset()
}
type WriterFunc = func(Writer)

// Dialect essentials (Builders)
type InternalBuilder interface {
	// Writes a condition (must be wrapped)
	Cond(Writer, Cond)
}

type ValueBuilder interface {
	String(string) string
	Number(int) string
	Float(float64) string
	Bool(bool) string
	Null() string

	// Writes a placeholder based on position
	Placeholder(number int) string
	Wildcard() string
}

type ColumnBuilder interface {
	Table(Writer, Table)

	// Writing modes
	ColWriteDef(Writer, Column)
	ColWriteAlias(Writer, Column)
	ColWriteExpr(Writer, Column)
	ColWriteBase(Writer, Column)

	// Specials columns
	ColConcat([]any) WriterFunc
	ColSwitch(cond any, cas CaseCond) WriterFunc
	ColSearch(CaseCond) WriterFunc // Case identifier must be a Cond
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
	WriteAlias
	WriteExpr
)
