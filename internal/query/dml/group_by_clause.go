package dml

import "github.com/laacin/inyorm/internal/builder"

// --- Entity

type ClauseGroupBy struct {
	declared bool
	Values   []any
}

// --- PUB API

func (c *ClauseGroupBy) GroupBy(vals ...any) {
	c.declared = true
	c.Values = vals
}

// --- Build

func (*ClauseGroupBy) Kind() ClauseKind {
	return ClauseKindGroupBy
}

func (c *ClauseGroupBy) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *ClauseGroupBy) Build(b *builder.Builder) error {
	return nil
}
