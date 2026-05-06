package clause

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/impl/expression"
)

type HavingImpl[Cond, CondNext any] struct {
	declared bool
	emb      entity.Having
	cond     *expression.ConditionImpl[Cond, CondNext]
}

func (c *HavingImpl[Cond, CondNext]) Having(ident any) Cond {
	c.declared = true
	c.cond = &expression.ConditionImpl[Cond, CondNext]{}
	return c.cond.Start(ident)
}

// --- Build

func (c *HavingImpl[Cond, CondNext]) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *HavingImpl[Cond, CondNext]) Kind() entity.ClauseKind {
	return entity.ClauseHaving
}

func (c *HavingImpl[Cond, CondNext]) Build() entity.Clause {
	c.emb.Cond = *c.cond.Build().(*entity.Condition)
	return &c.emb
}
