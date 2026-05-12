package clause

import (
	"github.com/laacin/inyorm/internal/entity/api"
	"github.com/laacin/inyorm/internal/entity/dml"
	"github.com/laacin/inyorm/internal/impl/expression"
)

type JoinImpl struct {
	declared bool
	emb      dml.Join
	current  *dml.JoinSegment
	conds    []*expression.ConditionImpl
}

func (c *JoinImpl) Join(table any) api.JoinNext {
	c.declared = true
	c.current = &dml.JoinSegment{
		Type:  dml.JoinInner,
		Table: table,
	}
	return c
}

func (c *JoinImpl) Left() api.JoinEnd {
	c.current.Type = dml.JoinLeft
	return c
}
func (c *JoinImpl) Right() api.JoinEnd {
	c.current.Type = dml.JoinRight
	return c
}
func (c *JoinImpl) Full() api.JoinEnd {
	c.current.Type = dml.JoinFull
	return c
}
func (c *JoinImpl) Cross() {
	c.current.Type = dml.JoinCross
	c.emb.Joins = append(c.emb.Joins, *c.current)
	c.conds = append(c.conds, nil)
}

func (c *JoinImpl) On(ident any) api.Condition {
	cond := &expression.ConditionImpl{}
	c.emb.Joins = append(c.emb.Joins, *c.current)
	c.conds = append(c.conds, cond)
	return cond.Start(ident)
}

// --- Build

func (c *JoinImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *JoinImpl) Kind() dml.ClauseKind {
	return dml.ClauseJoin
}

func (c *JoinImpl) Build() (dml.Clause, error) {
	for i, cond := range c.conds {
		if cond == nil {
			c.emb.Joins[i].Cond = nil
		}
		c.emb.Joins[i].Cond = cond.Build().(*dml.Condition)
	}
	return &c.emb, nil
}
