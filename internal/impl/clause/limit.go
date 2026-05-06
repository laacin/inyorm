package clause

import "github.com/laacin/inyorm/internal/entity"

type LimitImpl struct {
	declared bool
	emb      entity.Limit
}

func (c *LimitImpl) Limit(value int) {
	if value > 0 {
		c.declared = true
		c.emb.ValueNumber = value
	}
}

// --- Deref

func (c *LimitImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *LimitImpl) Defer() entity.Clause {
	return &c.emb
}
