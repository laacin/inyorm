package clause

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/ir/dml"
)

type OrderByImpl struct {
	declared bool
	emb      dml.OrderBy
	segments []*dml.OrderSegment
	current  *dml.OrderSegment
}

func (c *OrderByImpl) OrderBy(value any) api.OrderByNext {
	c.declared = true
	seg := &dml.OrderSegment{Value: value}
	c.segments = append(c.segments, seg)
	c.current = seg
	return c
}

func (c *OrderByImpl) Desc() {
	c.current.Descending = true
}

// --- Build

func (c *OrderByImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *OrderByImpl) Kind() dml.ClauseKind {
	return dml.ClauseOrderBy
}

func (c *OrderByImpl) Build() (dml.Clause, error) {
	c.emb.Orders = make([]dml.OrderSegment, len(c.segments))
	for i, seg := range c.segments {
		c.emb.Orders[i] = *seg
	}

	return &c.emb, nil
}
