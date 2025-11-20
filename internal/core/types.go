package core

type Builder = func(Writer)

type Writer interface {
	Write(str string)
	Char(byt byte)

	Placeholder()
	Identifier(v any, ctx ClauseType)
	Value(v any, ctx ClauseType)
	Column(table, base string)
	Table(table string)

	Split() Writer
	ToString() string
	Reset()
}

type Column interface {
	Base() Builder
	Expr() Builder
	Def() Builder
	Alias() Builder
}

type Clause interface {
	IsDeclared() bool
	Name() ClauseType
	Build(Writer)
}
