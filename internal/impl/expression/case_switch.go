package expression

import (
	"github.com/laacin/inyorm/internal/entity/api"
	"github.com/laacin/inyorm/internal/entity/dml"
)

type CaseSwitchImpl struct {
	dml.CaseSwitch
	current *dml.CaseWhen
}

func (c *CaseSwitchImpl) Start(cond any) api.Case {
	c.Cond = cond
	return c
}

func (c *CaseSwitchImpl) When(when any) api.CaseNext {
	c.current = &dml.CaseWhen{When: when}
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

func (c *CaseSwitchImpl) Build() dml.Value {
	return &c.CaseSwitch
}
