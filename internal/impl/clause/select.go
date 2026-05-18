package clause

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/dml"
)

type SelectImpl struct {
	declared bool
	emb      dml.Select
}

func (c *SelectImpl) Distinct() {
	c.declared = true
	c.emb.Distinct = true
}

func (c *SelectImpl) Select(values ...any) api.SelectNext {
	c.declared = true
	c.emb.Values = values
	return c
}

// --- Build

func (c *SelectImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *SelectImpl) Kind() dml.ClauseKind {
	return dml.ClauseSelect
}

func (c *SelectImpl) Build(w core.InternalWriter, dial dml.ClauseWriter) error {
	dial.WriteSelect(w, &c.emb)
	return nil
}
