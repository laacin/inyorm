package clause

import "github.com/laacin/inyorm/internal/entity"

type OffsetImpl struct {
	declared bool
	emb      entity.Offset
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

func (c *OffsetImpl) Kind() entity.ClauseKind {
	return entity.ClauseOffset
}

func (c *OffsetImpl) Build() entity.Clause {
	return &c.emb
}
