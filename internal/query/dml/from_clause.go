package dml

import "github.com/laacin/inyorm/internal/core"

// --- Entity

type ClauseFrom struct {
	declared bool
	Value    any
}

// --- PUB API

func (c *ClauseFrom) From(from any) {
	c.declared = true
	c.Value = from
}

// --- Build

func (*ClauseFrom) Kind() ClauseKind {
	return ClauseKindFrom
}

func (c *ClauseFrom) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *ClauseFrom) Build(w core.InternalWriter, dial ClauseWriter) error {
	dial.WriteFrom(w, c)
	return nil
}
