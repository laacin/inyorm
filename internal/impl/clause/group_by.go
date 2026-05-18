package clause

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/dml"
)

type GroupByImpl struct {
	declared bool
	emb      dml.GroupBy
}

func (c *GroupByImpl) GroupBy(values ...any) {
	c.declared = true
	c.emb.Values = values
}

// --- Build

func (c *GroupByImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *GroupByImpl) Kind() dml.ClauseKind {
	return dml.ClauseGroupBy
}

func (c *GroupByImpl) Build(w core.InternalWriter, dial dml.ClauseWriter) error {
	dial.WriteGroupBy(w, &c.emb)
	return nil
}
