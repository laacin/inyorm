package clause

import "github.com/laacin/inyorm/internal/entity"

type GroupByImpl struct {
	declared bool
	emb      entity.GroupBy
}

func (c *GroupByImpl) GroupBy(values ...any) {
	c.declared = true
	c.emb.Values = values
}

// --- Defer

func (c *GroupByImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *GroupByImpl) Defer() entity.Clause {
	return &c.emb
}
