package expression

import (
	"github.com/laacin/inyorm/internal/entity/api"
	"github.com/laacin/inyorm/internal/entity/dml"
)

type CaseSearchImpl struct {
	dml.CaseSearch
	current *dml.CaseWhen
}

func (c *CaseSearchImpl) When(when any) api.CaseNext {
	c.current = &dml.CaseWhen{When: when}
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

func (c *CaseSearchImpl) Build() dml.Value {
	return &c.CaseSearch
}
