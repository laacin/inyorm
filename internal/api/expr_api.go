package api

// --- Col

type Col interface {
	Count(dist ...bool) Col
	Sum(dist ...bool) Col
	Min(dist ...bool) Col
	Max(dist ...bool) Col
	Avg(dist ...bool) Col

	Add(v any) Col
	Sub(v any) Col
	Mul(v any) Col
	Div(v any) Col
	Mod(v any) Col
	Wrap() Col

	Lower() Col
	Upper() Col
	Trim() Col
	Round() Col
	Abs() Col

	As(alias string) Col
}

// --- Cond

type Cond interface {
	Not() Cond

	Equal(v any) CondNext
	Like(v any) CondNext
	In(vals []any) CondNext
	Between(v1, v2 any) CondNext
	Greater(v any) CondNext
	Less(v any) CondNext
	IsNull() CondNext
}

type CondNext interface {
	Or(ident any) Cond
	And(ident any) Cond
}

// --- Case

type Case interface {
	When(v any) CaseNext
	Else(v any)
}

type CaseNext interface {
	Then(v any) Case
}
