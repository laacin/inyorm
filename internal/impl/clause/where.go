package clause

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/impl/exprimpl"
	"github.com/laacin/inyorm/internal/ir/dml"
)

type WhereImpl struct {
	declared bool
	emb      dml.Where
	current  *exprimpl.ConditionImpl
}

func (c *WhereImpl) Where(ident any) api.Condition {
	c.declared = true
	cond := &exprimpl.ConditionImpl{}
	c.emb.Conds = append(c.emb.Conds, cond)
	return cond.Start(ident)
}

// --- Build

func (c *WhereImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *WhereImpl) Kind() dml.ClauseKind {
	return dml.ClauseWhere
}

func (c *WhereImpl) Build(w core.InternalWriter, dial dml.ClauseWriter) error {
	dial.WriteWhere(w, &c.emb)
	return nil
}
