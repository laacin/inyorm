package clause

import (
	"github.com/laacin/inyorm/internal/entity/api"
	"github.com/laacin/inyorm/internal/entity/dml"
	"github.com/laacin/inyorm/internal/entity/expr"
	"github.com/laacin/inyorm/internal/impl/expression"
)

type HavingImpl struct {
	declared bool
	emb      dml.Having
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

func (c *HavingImpl) Kind() dml.ClauseKind {
	return dml.ClauseHaving
}

func (c *HavingImpl) Build() (dml.Clause, error) {
	c.emb.Cond = c.cond.Build().(*expr.Condition)
	return &c.emb, nil
}
