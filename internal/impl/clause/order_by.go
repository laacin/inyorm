package clause

import "github.com/laacin/inyorm/internal/entity"

type OrderByImpl[Next any] struct {
	declared bool
	emb      entity.OrderBy
	segments []*entity.OrderSegment
	current  *entity.OrderSegment
}

func (c *OrderByImpl[Next]) OrderBy(value any) Next {
	c.declared = true
	seg := &entity.OrderSegment{Value: value}
	c.segments = append(c.segments, seg)
	c.current = seg
	return any(c).(Next)
}

func (c *OrderByImpl[Next]) Desc() {
	c.current.Descending = true
}

// --- Deref

func (c *OrderByImpl[Next]) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *OrderByImpl[Next]) Defer() entity.Clause {
	c.emb.Orders = make([]entity.OrderSegment, len(c.segments))
	for i, seg := range c.segments {
		c.emb.Orders[i] = *seg
	}

	return &c.emb
}
