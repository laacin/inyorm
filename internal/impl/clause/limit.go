package clause

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/dml"
)

type LimitImpl struct {
	declared bool
	emb      dml.Limit
}

func (c *LimitImpl) Limit(value int) {
	if value > 0 {
		c.declared = true
		c.emb.ValueNumber = value
	}
}

// --- Build

func (c *LimitImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *LimitImpl) Kind() dml.ClauseKind {
	return dml.ClauseLimit
}

func (c *LimitImpl) Build(w core.InternalWriter, dial dml.ClauseWriter) error {
	dial.WriteLimit(w, &c.emb)
	return nil
}
