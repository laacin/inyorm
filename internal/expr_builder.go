package internal

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/expr"
)

type ExprBuilder struct{ Ref string }

func (e *ExprBuilder) Start(defaultTable string) *ExprBuilder {
	e.Ref = defaultTable
	return e
}

func (e *ExprBuilder) Table(name string) any {
	tbl := &expr.Table{}
	return tbl.Start(name)
}

func (e *ExprBuilder) Col(name string, ref ...string) api.Col {
	col := &expr.Col{}
	return col.Start(name, getLast(e.Ref, ref))
}

func (e *ExprBuilder) All(ref ...string) api.Col {
	col := &expr.Col{}
	return col.Start("*", getLast(e.Ref, ref))
}

func (e *ExprBuilder) Param(value ...any) any {
	param := &expr.Param{}
	return param.Start(len(value) > 0, getLast(nil, value))
}

func (e *ExprBuilder) Cond(ident any) api.Cond {
	cond := &expr.Cond{}
	return cond.Start(ident)
}

func (e *ExprBuilder) Concat(values ...any) api.Col {
	col := &expr.Col{}
	return col.StartFrom((&expr.Concat{}).Start(values))
}

func (e *ExprBuilder) Switch(cond any, fn func(api.Case)) api.Col {
	cs := &expr.Case{}
	fn(cs.StartSwitch(cond))

	col := &expr.Col{}
	return col.StartFrom(cs)
}

func (e *ExprBuilder) Search(fn func(api.Case)) api.Col {
	cs := &expr.Case{}
	fn(cs.StartSearch())

	col := &expr.Col{}
	return col.StartFrom(cs)

}

// --- Helpers
func getLast[T any](prev T, candidate []T) T {
	if len(candidate) > 0 {
		return candidate[0]
	}
	return prev
}
