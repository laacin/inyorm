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

// --- Build

func (c *GroupByImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *GroupByImpl) Kind() entity.ClauseKind {
	return entity.ClauseGroupBy
}

func (c *GroupByImpl) Build() (entity.Clause, error) {
	return &c.emb, nil
}
