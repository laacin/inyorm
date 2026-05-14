package clause

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/impl/exprimpl"
	"github.com/laacin/inyorm/internal/ir/dml"
	"github.com/laacin/inyorm/internal/ir/expr"
)

type HavingImpl struct {
	declared bool
	emb      dml.Having
	cond     *exprimpl.ConditionImpl
}

func (c *HavingImpl) Having(ident any) api.Condition {
	c.declared = true
	c.cond = &exprimpl.ConditionImpl{}
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
