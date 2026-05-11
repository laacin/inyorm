package clause

import "github.com/laacin/inyorm/internal/entity"

type DeleteImpl struct {
	declared bool
	emb      entity.Delete
}

func (c *DeleteImpl) Delete() {
	c.declared = true
}

// --- Build

func (c *DeleteImpl) Kind() entity.ClauseKind {
	return entity.ClauseDelete
}

func (c *DeleteImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *DeleteImpl) Build() (entity.Clause, error) {
	return &c.emb, nil
}
