package core

type Column interface {
	Def(Writer)
	Expr(Writer)
	Alias(Writer)
	Base(Writer)

	Count(distinct bool)
	Sum(distinct bool)
	Min(distinct bool)
	Max(distinct bool)
	Avg(distinct bool)

	Add(v any)
	Sub(v any)
	Mul(v any)
	Div(v any)
	Mod(v any)
	Wrap()

	Lower()
	Upper()
	Trim()
	Round()
	Abs()

	As(v string)
}

type ColExpr interface {
	Col(col, table string) Column
	All() Column
	Concat(v []any) Column
	Condition(identifier any) Condition
	Switch(cond any, cs Case) Column
	Search(cs Case) Column
}

// ---- Case condition

type Case interface {
	When(v any)
	Then(v any)
	Else(v any)
}
