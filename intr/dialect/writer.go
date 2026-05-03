package dialect

// Writer used by dialect
type Writer interface {
	Write(string)
	Char(byte)
	Value(v any, ctx ClauseName)
}

// Writer escentials
type ValueWriter interface {
	String(string) string
	Number(int) string
	Float(float64) string
	Bool(bool) string
	Null() string
	Param(num int) string

	Table(table Table, def bool) string

	ColDef(Column) string
	ColAlias(Column) string
	ColExpr(Column) string
	ColBase(Column) string

	Cond(Cond) string // must be wrapped
}
