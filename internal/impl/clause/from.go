package clause

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/dml"
)

type FromImpl struct {
	declared bool
	emb      dml.From
}

func (c *FromImpl) From(from any) {
	c.declared = true
	c.emb.Value = from
}

// --- Build

func (c *FromImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *FromImpl) Kind() dml.ClauseKind {
	return dml.ClauseFrom
}

func (c *FromImpl) Build(w core.InternalWriter, dial dml.ClauseWriter) error {
	dial.WriteFrom(w, &c.emb)
	return nil
}
