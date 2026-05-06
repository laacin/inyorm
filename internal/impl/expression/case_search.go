package expression

import "github.com/laacin/inyorm/internal/entity"

type CaseSearchImpl[Self, Next any] struct {
	entity.CaseSearch
	current *entity.CaseWhen
}

func (c *CaseSearchImpl[Self, Next]) When(when any) Next {
	c.current = &entity.CaseWhen{When: when}
	return any(c).(Next)
}

func (c *CaseSearchImpl[Self, Next]) Then(then any) Self {
	c.current.Then = then
	c.Whens = append(c.Whens, *c.current)
	return any(c).(Self)
}

func (c *CaseSearchImpl[Self, Next]) Else(els any) {
	c.Els = els
}

// --- Deref
func (c *CaseSearchImpl[Self, Next]) Deref() entity.Value {
	return &c.CaseSearch
}
