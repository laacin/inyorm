package statement

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/core/mapper"
)

type Binder[T any] struct {
	bind   any
	params core.ParamStore
	err    error

	chain T
}

func NewBinder[T any](params core.ParamStore, chain T) Binder[T] {
	return Binder[T]{chain: chain, params: params}
}

func (p *Binder[T]) Bind(v any) T {
	p.bind = v
	return p.chain
}

func (p *Binder[T]) Value(id string, v any) T {
	if p.err != nil {
		return p.chain
	}

	m := mapper.New()
	kind := m.ReadKind(v).Kind
	if kind == core.KindStruct || kind == core.KindMap {
		p.params.FillObject(func(cols []string) []any { // TODO: allow lazy objects to be referenced by ID
			vals, err := m.ReadValues(cols, v)
			p.err = err
			return vals
		})
		return p.chain
	}

	p.params.Fill(id, v)
	return p.chain
}

func (p *Binder[T]) Values(v ...any) T {
	if p.err != nil {
		return p.chain
	}

	for _, val := range v {
		m := mapper.New()
		kind := m.ReadKind(val).Kind
		if kind == core.KindStruct || kind == core.KindMap {
			p.params.FillObject(func(cols []string) []any {
				vals, err := m.ReadValues(cols, val)
				p.err = err
				return vals
			})
			return p.chain
		}

		p.params.Fill("", val)
	}
	return p.chain
}
