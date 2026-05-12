package clause

import "github.com/laacin/inyorm/internal/entity/dml"

type OffsetImpl struct {
	declared bool
	emb      dml.Offset
}

func (c *OffsetImpl) Offset(value int) {
	if value > 0 {
		c.declared = true
		c.emb.ValueNumber = value
	}
}

// --- Build

func (c *OffsetImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *OffsetImpl) Kind() dml.ClauseKind {
	return dml.ClauseOffset
}

func (c *OffsetImpl) Build() (dml.Clause, error) {
	return &c.emb, nil
}
