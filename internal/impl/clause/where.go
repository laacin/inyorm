package clause

import (
	"github.com/laacin/inyorm/internal/entity/api"
	"github.com/laacin/inyorm/internal/entity/dml"
	"github.com/laacin/inyorm/internal/impl/expression"
)

type WhereImpl struct {
	declared bool
	emb      dml.Where
	conds    []*expression.ConditionImpl
	current  *expression.ConditionImpl
}

func (c *WhereImpl) Where(ident any) api.Condition {
	c.declared = true
	cond := &expression.ConditionImpl{}
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
		c.emb.Conds = append(c.emb.Conds, cond.Build().(*dml.Condition))
	}

	return &c.emb, nil
}
