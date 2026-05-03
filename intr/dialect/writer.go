package dialect

// Writer used by dialect
type Writer interface {
	Write(string)
	Char(byte)
	Value(v any, ctx ClauseName)
}

// Writer escentials
type ValueWriter interface {
	String(Writer, string)
	Number(Writer, int)
	Float(Writer, float64)
	Bool(Writer, bool)
	Null(Writer)
	Placeholder(w Writer, num int)

	Table(w Writer, table Table, def bool)

	ColDef(Writer, Column)
	ColAlias(Writer, Column)
	ColExpr(Writer, Column)
	ColBase(Writer, Column)

	Cond(Writer, Cond) // must be wrapped
}
