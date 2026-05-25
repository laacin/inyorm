package dml

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
)

// --- Entity

type ClauseOrderBy struct {
	declared bool
	Segments []OrderSegment
	current  *OrderSegment
}

// --- PUB API

func (c *ClauseOrderBy) OrderBy(v any) api.OrderByNext {
	if c.current != nil {
		c.Segments = append(c.Segments, *c.current)
	}

	c.declared = true
	c.current = &OrderSegment{Value: v}
	return c
}

func (c *ClauseOrderBy) Desc() {
	c.current.Descending = true
}

// --- Build

func (*ClauseOrderBy) Kind() ClauseKind {
	return ClauseKindOrderBy
}

func (c *ClauseOrderBy) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *ClauseOrderBy) Build(b *core.Builder) error {
	if c.current != nil {
		c.Segments = append(c.Segments, *c.current)
		c.current = nil
	}
	return nil
}

// --- Tools

type OrderSegment struct {
	Value      any
	Descending bool
}
