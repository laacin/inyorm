package dml

import "github.com/laacin/inyorm/internal/builder"

// --- Entity

type ClauseLimit struct {
	declared bool
	ValueInt int
}

// --- PUB API

func (c *ClauseLimit) Limit(v int) {
	if v > 0 {
		c.declared = true
		c.ValueInt = v
	}
}

// --- Build

func (*ClauseLimit) Kind() ClauseKind {
	return ClauseKindLimit
}

func (c *ClauseLimit) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *ClauseLimit) Build(b *builder.Builder) error {
	return nil
}
