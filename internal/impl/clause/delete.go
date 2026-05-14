package clause

import "github.com/laacin/inyorm/internal/ir/dml"

type DeleteImpl struct {
	declared bool
	emb      dml.Delete
}

func (c *DeleteImpl) Delete() {
	c.declared = true
}

// --- Build

func (c *DeleteImpl) Kind() dml.ClauseKind {
	return dml.ClauseDelete
}

func (c *DeleteImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *DeleteImpl) Build() (dml.Clause, error) {
	return &c.emb, nil
}
