package expression

import "github.com/laacin/inyorm/internal/entity"

type CaseSwitchImpl[Self, Next any] struct {
	entity.CaseSwitch
	current *entity.CaseWhen
}

func (c *CaseSwitchImpl[Self, Next]) Start(cond any) Self {
	c.Cond = cond
	return any(c).(Self)
}

func (c *CaseSwitchImpl[Self, Next]) When(when any) Next {
	c.current = &entity.CaseWhen{When: when}
	return any(c).(Next)
}

func (c *CaseSwitchImpl[Self, Next]) Then(then any) Self {
	c.current.Then = then
	c.Whens = append(c.Whens, *c.current)
	return any(c).(Self)
}

func (c *CaseSwitchImpl[Self, Next]) Else(els any) {
	c.Els = els
}

// --- Build

func (c *CaseSwitchImpl[Self, Next]) Build() entity.Value {
	return &c.CaseSwitch
}
