package dml

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/builder"
)

// --- Entity
type ClauseSelect struct {
	declared bool
	Dist     bool
	Values   []any
}

// --- PUB API

func (c *ClauseSelect) Select(vals ...any) api.SelectNext {
	c.declared = true
	c.Values = vals
	return c
}

func (c *ClauseSelect) Distinct() {
	c.Dist = true
}

// --- Build

func (*ClauseSelect) Kind() ClauseKind {
	return ClauseKindSelect
}

func (c *ClauseSelect) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *ClauseSelect) Build(b *builder.Builder) error {
	return nil
}
