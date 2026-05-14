package clause

import "github.com/laacin/inyorm/internal/ir/dml"

type GroupByImpl struct {
	declared bool
	emb      dml.GroupBy
}

func (c *GroupByImpl) GroupBy(values ...any) {
	c.declared = true
	c.emb.Values = values
}

// --- Build

func (c *GroupByImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *GroupByImpl) Kind() dml.ClauseKind {
	return dml.ClauseGroupBy
}

func (c *GroupByImpl) Build() (dml.Clause, error) {
	return &c.emb, nil
}
