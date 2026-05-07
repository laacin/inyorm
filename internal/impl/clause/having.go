package clause

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/entity/api"
	"github.com/laacin/inyorm/internal/impl/expression"
)

type HavingImpl struct {
	declared bool
	emb      entity.Having
	cond     *expression.ConditionImpl
}

func (c *HavingImpl) Having(ident any) api.Condition {
	c.declared = true
	c.cond = &expression.ConditionImpl{}
	return c.cond.Start(ident)
}

// --- Build

func (c *HavingImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *HavingImpl) Kind() entity.ClauseKind {
	return entity.ClauseHaving
}

func (c *HavingImpl) Build() entity.Clause {
	c.emb.Cond = c.cond.Build().(*entity.Condition)
	return &c.emb
}
