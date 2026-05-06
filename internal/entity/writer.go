package entity

// Writing modes
type WritingMode int

const (
	WriteDef WritingMode = iota
	WriteBase
	WriteAlias
	WriteExpr
)

// Writer used by a dialect
type Writer interface {
	Write(string)
	Char(byte)

	Value(v any, mode WritingMode)
	ValueCount() int
	GetRef(table string) (byte, bool)

	New() Writer
	Result() string
	Reset()
}

type WriterFunc = func(Writer)

// --- Dialect

type Dialect interface {
	ValueWriter
	ClauseWriter
	StatementBuilder
}

type ValueWriter interface {
	// Literals
	WriteString(Writer, string)
	WriteNumber(Writer, int)
	WriteFloat(Writer, float64)
	WriteBool(Writer, bool)
	WriteNull(Writer)
	WriteWildcard(Writer)

	// Specials
	WritePlaceholder(Writer)
	WriteConcat(Writer, *Concat)
	WriteCondition(Writer, *Condition, WritingMode)
	WriteCaseSwitch(Writer, *CaseSwitch, WritingMode)
	WriteCaseSearch(Writer, *CaseSearch, WritingMode)

	// Table
	WriteTable(Writer, *Table)

	// Column
	WriteColBase(Writer, *Column)
	WriteColExpr(Writer, *Column)
	WriteColAlias(Writer, *Column)
	WriteColDef(Writer, *Column)
}

type ClauseWriter interface {
	WriteSelect(Writer, *Select)
	WriteFrom(Writer, *From)
	WriteJoin(Writer, *Join)
	WriteWhere(Writer, *Where)
	WriteGroupBy(Writer, *GroupBy)
	WriteHaving(Writer, *Having)
	WriteOrderBy(Writer, *OrderBy)
	WriteLimit(Writer, *Limit)
	WriteOffset(Writer, *Offset)

	WriteInsertInto(Writer, *InsertInto)
	WriteUpdate(Writer, *Update)
	WriteDelete(Writer, *Delete)
}

type StatementBuilder interface {
	SelectOrder() []ClauseKind
	InsertOrder() []ClauseKind
	UpdateOrder() []ClauseKind
	DeleteOrder() []ClauseKind
}
