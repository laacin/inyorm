package core

type Cond interface {
	Not() Cond

	Equal(v any) CondNext
	Like(v any) CondNext
	In(v ...any) CondNext
	Between(minV, maxV any) CondNext

	Greater(v any) CondNext
	Less(v any) CondNext

	IsNull() CondNext
}

type CondNext interface {
	And(identifier ...any) Cond
	Or(identifier ...any) Cond
}
