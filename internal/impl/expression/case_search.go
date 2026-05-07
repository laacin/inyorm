package expression

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/entity/api"
)

type CaseSearchImpl struct {
	entity.CaseSearch
	current *entity.CaseWhen
}

func (c *CaseSearchImpl) When(when any) api.CaseNext {
	c.current = &entity.CaseWhen{When: when}
	return c
}

func (c *CaseSearchImpl) Then(then any) api.Case {
	c.current.Then = then
	c.Whens = append(c.Whens, *c.current)
	return c
}

func (c *CaseSearchImpl) Else(els any) {
	c.Els = els
}

// --- Build

func (c *CaseSearchImpl) Build() entity.Value {
	return &c.CaseSearch
}
