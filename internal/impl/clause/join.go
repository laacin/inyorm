package clause

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/impl/expression"
)

type JoinImpl[Next, End, Cond, CondNext any] struct {
	declared bool
	emb      entity.Join
	current  *entity.JoinSegment
	conds    []*expression.ConditionImpl[Cond, CondNext]
}

func (c *JoinImpl[Next, End, Cond, CondNext]) Join(table string) Next {
	c.declared = true
	c.current = &entity.JoinSegment{
		Type:  entity.JoinInner,
		Table: entity.Table{Value: table},
	}
	return any(c).(Next)
}

func (c *JoinImpl[Next, End, Cond, CondNext]) Left() End {
	c.current.Type = entity.JoinLeft
	return any(c).(End)
}
func (c *JoinImpl[Next, End, Cond, CondNext]) Right() End {
	c.current.Type = entity.JoinRight
	return any(c).(End)
}
func (c *JoinImpl[Next, End, Cond, CondNext]) Full() End {
	c.current.Type = entity.JoinFull
	return any(c).(End)
}
func (c *JoinImpl[Next, End, Cond, CondNext]) Cross() {
	c.current.Type = entity.JoinCross
	c.emb.Joins = append(c.emb.Joins, *c.current)
	c.conds = append(c.conds, nil)
}

func (c *JoinImpl[Next, End, Cond, CondNext]) On(ident any) Cond {
	cond := &expression.ConditionImpl[Cond, CondNext]{}
	c.emb.Joins = append(c.emb.Joins, *c.current)
	c.conds = append(c.conds, cond)
	return cond.Start(ident)
}

// --- Build

func (c *JoinImpl[Next, End, Cond, CondNext]) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *JoinImpl[Next, End, Cond, CondNext]) Kind() entity.ClauseKind {
	return entity.ClauseJoin
}

func (c *JoinImpl[Next, End, Cond, CondNext]) Build() entity.Clause {
	for i, cond := range c.conds {
		if cond == nil {
			c.emb.Joins[i].Cond = nil
		}
		c.emb.Joins[i].Cond = cond.Build().(*entity.Condition)
	}
	return &c.emb
}
