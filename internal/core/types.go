package core

type Builder = func(Writer)

type Writer interface {
	Write(str string)
	Char(byt byte)

	Param(value []any)
	Value(v any, colWrite ColumnType)
	Column(table, base string)
	Table(table string)

	Split() Writer
	ToString() string
	Reset()
}

type Column interface {
	RawBase() string
	Base() Builder
	Expr() Builder
	Def() Builder
	Alias() Builder
}

type Clause interface {
	IsDeclared() bool
	Name() string
	Build(Writer, *Config) error
}
