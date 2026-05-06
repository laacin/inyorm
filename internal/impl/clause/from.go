package clause

import "github.com/laacin/inyorm/internal/entity"

type FromImpl struct {
	declared bool
	emb      entity.From
}

func (c *FromImpl) From(from any) {
	c.declared = true
	c.emb.Value = from
}

// --- Deref

func (c *FromImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *FromImpl) Defer() entity.Clause {
	return &c.emb
}
