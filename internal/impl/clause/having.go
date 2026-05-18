package clause

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/impl/exprimpl"
	"github.com/laacin/inyorm/internal/ir/dml"
)

type HavingImpl struct {
	declared bool
	emb      dml.Having
	cond     *exprimpl.ConditionImpl
}

func (c *HavingImpl) Having(ident any) api.Condition {
	c.declared = true
	cond := &exprimpl.ConditionImpl{}
	c.emb.Cond = cond
	return cond.Start(ident)
}

// --- Build

func (c *HavingImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *HavingImpl) Kind() dml.ClauseKind {
	return dml.ClauseHaving
}

func (c *HavingImpl) Build(w core.InternalWriter, dial dml.ClauseWriter) error {
	dial.WriteHaving(w, &c.emb)
	return nil
}
