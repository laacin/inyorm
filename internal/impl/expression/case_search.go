package expression

import "github.com/laacin/inyorm/internal/entity"

type CaseSearchImpl[Self, Next any] struct {
	entity.CaseSearch
	Current *entity.CaseWhen
}

func (c *CaseSearchImpl[Self, Next]) When(when any) Next {
	c.Current = &entity.CaseWhen{When: when}
	return any(c).(Next)
}

func (c *CaseSearchImpl[Self, Next]) Then(then any) Self {
	c.Current.Then = then
	c.Whens = append(c.Whens, *c.Current)
	return any(c).(Self)
}

func (c *CaseSearchImpl[Self, Next]) Else(els any) {
	c.Els = els
}

// --- Deref
func (c *CaseSearchImpl[Self, Next]) Deref() entity.Value {
	return &c.CaseSearch
}
