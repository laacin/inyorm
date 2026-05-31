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
	e.alias.Set(name)
	return NewTable(name, func() core.Reference {
		return e.alias.Get(name)
	})
}

func (e *Builder) Col(name string, ref ...string) api.Col {
	return NewCol(name, func() core.Reference {
		return e.getRef(ref)
	})
}

func (e *Builder) All(ref ...string) api.Col {
	return NewCol("*", func() core.Reference {
		return e.getRef(ref)
	})
}

func (e *Builder) Param(v any) any {
	return NewPlaceholder(func() core.ParamIndex {
		e.param.Store(v)
		return e.param.LastIndex(0)
	})
}

func (e *Builder) Lazy(id ...string) any {
	return NewPlaceholder(func() core.ParamIndex {
		e.param.Lazy(core.GetLast("", id))
		return e.param.LastIndex(0)
	}, true)
}

func (e *Builder) Cond(ident any) api.Cond {
	return NewCond(ident)
}

func (e *Builder) Concat(values ...any) api.Col {
	return NewColFrom(NewConcat(values))
}

func (e *Builder) Switch(cond any, fn func(api.Case)) api.Col {
	cs := NewCaseSwitch(cond)
	fn(cs)

	return NewColFrom(cs)
}

func (e *Builder) Search(fn func(api.Case)) api.Col {
	cs := NewCaseSearch()
	fn(cs)

	return NewColFrom(cs)

}

// --- Helpers
func (e *Builder) getRef(candidate []string) core.Reference {
	if len(candidate) > 0 {
		return e.alias.Get(candidate[0])
	}
	return e.alias.GetMain()
}
