package clause

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/entity/api"
)

type SelectImpl struct {
	declared bool
	emb      entity.Select
}

func (c *SelectImpl) Distinct() {
	c.declared = true
	c.emb.Distinct = true
}

func (c *SelectImpl) Select(values ...any) api.SelectNext {
	c.declared = true
	c.emb.Values = values
	return c
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
