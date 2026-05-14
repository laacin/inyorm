package exprimpl

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/ir/expr"
)

type CaseSearchImpl struct {
	expr.CaseSearch
	current *expr.CaseWhen
}

func (c *CaseSearchImpl) When(when any) api.CaseNext {
	c.current = &expr.CaseWhen{When: when}
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

func (c *CaseSearchImpl) Build() expr.Value {
	return &c.CaseSearch
}
