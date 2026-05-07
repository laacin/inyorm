package clause

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/entity/api"
)

type SelectImpl struct {
	declared bool
	emb      entity.Select
}

func (c *SelectImpl) Distinct() api.SelectNext {
	c.declared = true
	c.emb.Distinct = true
	return c
}

func (c *SelectImpl) Select(values ...any) {
	c.declared = true
	c.emb.Values = values
}

// --- Build

func (c *SelectImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *SelectImpl) Kind() entity.ClauseKind {
	return entity.ClauseSelect
}

func (c *SelectImpl) Build() entity.Clause {
	return &c.emb
}
