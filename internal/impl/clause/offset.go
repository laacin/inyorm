package clause

import "github.com/laacin/inyorm/internal/entity"

type OffsetImpl struct {
	declared bool
	emb      entity.Limit
}

func (c *OffsetImpl) Offset(value int) {
	if value > 0 {
		c.declared = true
		c.emb.ValueNumber = value
	}
}

// --- Deref

func (c *OffsetImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *OffsetImpl) Defer() entity.Clause {
	return &c.emb
}
