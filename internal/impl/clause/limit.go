package clause

import "github.com/laacin/inyorm/internal/ir/dml"

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

func (c *LimitImpl) Build() (dml.Clause, error) {
	return &c.emb, nil
}
