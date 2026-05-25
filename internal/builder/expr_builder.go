package builder

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
)

type ExprBuilder struct {
	param core.ParamStore
	Ref   string
}

// start
func (e *ExprBuilder) Start(paramStore core.ParamStore) *ExprBuilder {
	e.param = paramStore
	return e
}

func (e *ExprBuilder) AttachRef(ref string) *ExprBuilder {
	e.Ref = ref
	return e
}

// --- PUB API

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

func (e *ExprBuilder) Param(v any) any {
	return (&expr.Placeholder{}).Start(func() core.ParamIndex {
		e.param.Store(v)
		return e.param.LastIndex(0)
	})
}

func (e *ExprBuilder) Lazy(id ...string) any {
	return (&expr.Placeholder{}).StartLazy(func() core.ParamIndex {
		e.param.Lazy(getLast("", id))
		return e.param.LastIndex(0)
	})
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
