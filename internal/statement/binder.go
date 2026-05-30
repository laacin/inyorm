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

func (p *Binder[T]) Values(v any, id ...string) T {
	if p.err != nil {
		return p.chain
	}

	m := mapper.New()
	if m.ReadKind(v).Kind == core.KindStruct {
		p.params.FillObj(func(cols []string) []any {
			vals, err := m.ReadValues(cols, v)
			if err != nil {
				p.err = err
				return nil
			}

			return vals
		})
		return p.chain
	}

	p.params.Fill(core.GetLast("", id), v)
	return p.chain
}
