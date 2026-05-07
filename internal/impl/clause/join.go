package clause

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/entity/api"
	"github.com/laacin/inyorm/internal/impl/expression"
)

type JoinImpl struct {
	declared bool
	emb      entity.Join
	current  *entity.JoinSegment
	conds    []*expression.ConditionImpl
}

func (c *JoinImpl) Join(table any) api.JoinNext {
	c.declared = true
	c.current = &entity.JoinSegment{
		Type:  entity.JoinInner,
		Table: table,
	}
	return c
}

func (c *JoinImpl) Left() api.JoinEnd {
	c.current.Type = entity.JoinLeft
	return c
}
func (c *JoinImpl) Right() api.JoinEnd {
	c.current.Type = entity.JoinRight
	return c
}
func (c *JoinImpl) Full() api.JoinEnd {
	c.current.Type = entity.JoinFull
	return c
}
func (c *JoinImpl) Cross() {
	c.current.Type = entity.JoinCross
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

func (c *JoinImpl) Kind() entity.ClauseKind {
	return entity.ClauseJoin
}

func (c *JoinImpl) Build() entity.Clause {
	for i, cond := range c.conds {
		if cond == nil {
			c.emb.Joins[i].Cond = nil
		}
		c.emb.Joins[i].Cond = cond.Build().(*entity.Condition)
	}
	return &c.emb
}
