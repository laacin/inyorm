package exprimpl

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/ir/expr"
)

type CaseSwitchImpl struct {
	expr.CaseSwitch
	current *expr.CaseWhen
}

func (c *CaseSwitchImpl) Start(cond any) api.Case {
	c.Cond = cond
	return c
}

func (c *CaseSwitchImpl) When(when any) api.CaseNext {
	c.current = &expr.CaseWhen{When: when}
	return c
}

func (c *CaseSwitchImpl) Then(then any) api.Case {
	c.current.Then = then
	c.Whens = append(c.Whens, *c.current)
	return c
}

func (c *CaseSwitchImpl) Else(els any) {
	c.Els = els
}

// --- Build

func (c *CaseSwitchImpl) Build() expr.Value {
	return &c.CaseSwitch
}
