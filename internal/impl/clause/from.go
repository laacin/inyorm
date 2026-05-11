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

// --- Build

func (c *FromImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *FromImpl) Kind() entity.ClauseKind {
	return entity.ClauseFrom
}

func (c *FromImpl) Build() (entity.Clause, error) {
	return &c.emb, nil
}
