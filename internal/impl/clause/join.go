package clause

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/impl/exprimpl"
	"github.com/laacin/inyorm/internal/ir/dml"
)

type JoinImpl struct {
	declared bool
	emb      dml.Join
	current  *dml.JoinSegment
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
}

func (c *JoinImpl) On(ident any) api.Condition {
	cond := &exprimpl.ConditionImpl{}
	c.current.Cond = cond
	c.emb.Joins = append(c.emb.Joins, *c.current)
	return cond.Start(ident)
}

// --- Build

func (c *JoinImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *JoinImpl) Kind() dml.ClauseKind {
	return dml.ClauseJoin
}

func (c *JoinImpl) Build(w core.InternalWriter, dial dml.ClauseWriter) error {
	dial.WriteJoin(w, &c.emb)
	return nil
}
