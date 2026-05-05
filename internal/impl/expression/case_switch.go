package expression

import "github.com/laacin/inyorm/internal/entity"

type CaseSwitchImpl[Self, Next any] struct {
	entity.CaseSwitch
	Current *entity.CaseWhen
}

func (c *CaseSwitchImpl[Self, Next]) Start(cond any) Self {
	c.Cond = cond
	return any(c).(Self)
}

func (c *CaseSwitchImpl[Self, Next]) When(when any) Next {
	c.Current = &entity.CaseWhen{When: when}
	return any(c).(Next)
}

func (c *CaseSwitchImpl[Self, Next]) Then(then any) Self {
	c.Current.Then = then
	c.Whens = append(c.Whens, *c.Current)
	return any(c).(Self)
}

func (c *CaseSwitchImpl[Self, Next]) Else(els any) {
	c.Els = els
}

// --- Deref
func (c *CaseSwitchImpl[Self, Next]) Deref() entity.Value {
	return &c.CaseSwitch
}
