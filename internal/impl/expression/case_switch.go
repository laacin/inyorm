package expression

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/entity/api"
)

type CaseSwitchImpl struct {
	entity.CaseSwitch
	current *entity.CaseWhen
}

func (c *CaseSwitchImpl) Start(cond any) api.Case {
	c.Cond = cond
	return c
}

func (c *CaseSwitchImpl) When(when any) api.CaseNext {
	c.current = &entity.CaseWhen{When: when}
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

func (c *CaseSwitchImpl) Build() entity.Value {
	return &c.CaseSwitch
}
