package clause

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/impl/expression"
)

type WhereImpl[Cond, CondNext any] struct {
	declared bool
	emb      entity.Where
	conds    []*expression.ConditionImpl[Cond, CondNext]
	current  *expression.ConditionImpl[Cond, CondNext]
}

func (c *WhereImpl[Cond, CondNext]) Where(ident any) Cond {
	c.declared = true
	cond := &expression.ConditionImpl[Cond, CondNext]{}
	c.conds = append(c.conds, cond)
	return cond.Start(ident)
}

// --- Defer

func (c *WhereImpl[Cond, CondNext]) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *WhereImpl[Cond, CondNext]) Defer() entity.Clause {
	for _, cond := range c.conds {
		c.emb.Conds = append(c.emb.Conds, *cond.Deref().(*entity.Condition))
	}

	return &c.emb
}
