package dml

import "github.com/laacin/inyorm/internal/builder"

// --- Entity

type ClauseOffset struct {
	declared bool
	ValueInt int
}

// --- PUB API

func (c *ClauseOffset) Offset(v int) {
	if v > 0 {
		c.declared = true
		c.ValueInt = v
	}
}

// --- Build

func (*ClauseOffset) Kind() ClauseKind {
	return ClauseKindOffset
}

func (c *ClauseOffset) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *ClauseOffset) Build(b *builder.Builder) error {
	return nil
}
