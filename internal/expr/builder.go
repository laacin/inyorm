package expr

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
)

type Builder struct {
	param core.ParamStore
	alias core.AliasStore
}

func NewBuilder(param core.ParamStore, alias core.AliasStore) *Builder {
	return &Builder{param: param, alias: alias}
}

// --- PUB API

func (e *Builder) Table(name string) any {
	tbl := &Table{}
	e.alias.Set(name)
	return tbl.Start(name, func() core.Reference {
		return e.alias.Get(name)
	})
}

func (e *Builder) Col(name string, ref ...string) api.Col {
	col := &Col{}
	return col.Start(name, func() core.Reference {
		return e.getRef(ref)
	})
}

func (e *Builder) All(ref ...string) api.Col {
	col := &Col{}
	return col.Start("*", func() core.Reference {
		return e.getRef(ref)
	})
}

func (e *Builder) Param(v any) any {
	return (&Placeholder{}).Start(func() core.ParamIndex {
		e.param.Store(v)
		return e.param.LastIndex(0)
	})
}

func (e *Builder) Lazy(id ...string) any {
	return (&Placeholder{}).StartLazy(func() core.ParamIndex {
		e.param.Lazy(getLast("", id))
		return e.param.LastIndex(0)
	})
}

func (e *Builder) Cond(ident any) api.Cond {
	cond := &Cond{}
	return cond.Start(ident)
}

func (e *Builder) Concat(values ...any) api.Col {
	col := &Col{}
	return col.StartFrom((&Concat{}).Start(values))
}

func (e *Builder) Switch(cond any, fn func(api.Case)) api.Col {
	cs := &Case{}
	fn(cs.StartSwitch(cond))

	col := &Col{}
	return col.StartFrom(cs)
}

func (e *Builder) Search(fn func(api.Case)) api.Col {
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

func (e *Builder) getRef(candidate []string) core.Reference {
	if len(candidate) > 0 {
		return e.alias.Get(candidate[0])
	}
	return e.alias.GetMain()
}
