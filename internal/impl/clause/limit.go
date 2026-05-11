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

// --- Build

func (c *LimitImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *LimitImpl) Kind() entity.ClauseKind {
	return entity.ClauseLimit
}

func (c *LimitImpl) Build() (entity.Clause, error) {
	return &c.emb, nil
}
