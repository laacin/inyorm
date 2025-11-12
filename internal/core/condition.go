package core

type Condition interface {
	Start(ident any)

	Not()

	Equal(v any)
	Like(v any)
	In(v []any)
	Between(val1, val2 any)

	Greater(v any)
	Less(v any)

	IsNull()

	And(ident any)
	Or(ident any)
}
