package clause

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/entity/api"
)

type OrderByImpl struct {
	declared bool
	emb      entity.OrderBy
	segments []*entity.OrderSegment
	current  *entity.OrderSegment
}

func (c *OrderByImpl) OrderBy(value any) api.OrderByNext {
	c.declared = true
	seg := &entity.OrderSegment{Value: value}
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

func (c *OrderByImpl) Kind() entity.ClauseKind {
	return entity.ClauseOrderBy
}

func (c *OrderByImpl) Build() entity.Clause {
	c.emb.Orders = make([]entity.OrderSegment, len(c.segments))
	for i, seg := range c.segments {
		c.emb.Orders[i] = *seg
	}

	return &c.emb
}
