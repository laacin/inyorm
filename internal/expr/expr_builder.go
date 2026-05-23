package expr

import "github.com/laacin/inyorm/internal/api"

type ExprBuilderImpl struct{ Ref string }

func (e *ExprBuilderImpl) Start(defaultTable string) *ExprBuilderImpl {
	e.Ref = defaultTable
	return e
}

func (e *ExprBuilderImpl) Table(name string) any {
	tbl := &Table{}
	return tbl.Start(name)
}

func (e *ExprBuilderImpl) Col(name string, ref ...string) api.Col {
	col := &Col{}
	return col.Start(name, getLast(e.Ref, ref))
}

func (e *ExprBuilderImpl) All(ref ...string) api.Col {
	col := &Col{}
	return col.Start("*", getLast(e.Ref, ref))
}

func (e *ExprBuilderImpl) Param(value ...any) any {
	param := &Param{}
	return param.Start(len(value) > 0, getLast(nil, value))
}

func (e *ExprBuilderImpl) Cond(ident any) api.Cond {
	cond := &Cond{}
	return cond.Start(ident)
}

func (e *ExprBuilderImpl) Concat(values ...any) api.Col {
	col := &Col{}
	return col.StartFrom((&Concat{}).Start(values))
}

func (e *ExprBuilderImpl) Switch(cond any, fn func(api.Case)) api.Col {
	cs := &Case{}
	fn(cs.StartSwitch(cond))

	col := &Col{}
	return col.StartFrom(cs)
}

func (e *ExprBuilderImpl) Search(fn func(api.Case)) api.Col {
	cs := &Case{}
	fn(cs.StartSearch())

	col := &Col{}
	return col.StartFrom(cs)

}

// --- Helpers
func getLast[T any](prev T, candidate []T) T {
	if len(candidate) > 0 {
		return candidate[0]
	}
	return prev
}
