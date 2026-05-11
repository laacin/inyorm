package clause

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/entity/api"
	"github.com/laacin/inyorm/internal/impl/expression"
)

type WhereImpl struct {
	declared bool
	emb      entity.Where
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

func (c *WhereImpl) Kind() entity.ClauseKind {
	return entity.ClauseWhere
}

func (c *WhereImpl) Build() (entity.Clause, error) {
	for _, cond := range c.conds {
		c.emb.Conds = append(c.emb.Conds, cond.Build().(*entity.Condition))
	}

	return &c.emb, nil
}
