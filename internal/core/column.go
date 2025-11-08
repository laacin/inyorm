package core

type Column interface {
	Def() Builder
	Ref() Builder

	Count(distinct ...bool) Column
	Sum(distinct ...bool) Column
	Min(distinct ...bool) Column
	Max(distinct ...bool) Column
	Avg(distinct ...bool) Column

	Add(v any) Column
	Sub(v any) Column
	Mul(v any) Column
	Div(v any) Column
	Mod(v any) Column
	Wrap() Column

	Lower() Column
	Upper() Column
	Trim() Column
	Round() Column
	Abs() Column

	As(v string) Column
}

type ColExpr interface {
	Col(v any, table ...string) Column
	All() Column
	Concat(v ...any) Column
	Condition(identifier any) Cond
	Switch(cond any, fn func(cs CaseSwitch)) Column
	Search(fn func(cs CaseSearch)) Column
}

// ---- Case condition

type (
	CaseSwitch = Case[any]
	CaseSearch = Case[CondNext]
)

type Case[T any] interface {
	When(v T) CaseNext[T]
	Else(v any)
}

type CaseNext[T any] interface {
	Then(v any) Case[T]
}
