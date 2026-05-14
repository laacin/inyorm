package clause

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/impl/exprimpl"
	"github.com/laacin/inyorm/internal/ir/dml"
	"github.com/laacin/inyorm/internal/ir/expr"
)

type WhereImpl struct {
	declared bool
	emb      dml.Where
	conds    []*exprimpl.ConditionImpl
	current  *exprimpl.ConditionImpl
}

func (c *WhereImpl) Where(ident any) api.Condition {
	c.declared = true
	cond := &exprimpl.ConditionImpl{}
	c.conds = append(c.conds, cond)
	return cond.Start(ident)
}

// --- Build

func (c *WhereImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *WhereImpl) Kind() dml.ClauseKind {
	return dml.ClauseWhere
}

func (c *WhereImpl) Build() (dml.Clause, error) {
	for _, cond := range c.conds {
		c.emb.Conds = append(c.emb.Conds, cond.Build().(*expr.Condition))
	}

	return &c.emb, nil
}
