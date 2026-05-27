package dml

import "github.com/laacin/inyorm/internal/query"

// --- Entity

type ClauseDelete struct {
	declared bool
}

// --- PUB API

func (c *ClauseDelete) Delete() {
	c.declared = true
}

// --- Build

func (*ClauseDelete) Kind() ClauseKind {
	return ClauseKindDelete
}

func (c *ClauseDelete) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *ClauseDelete) Build(tools *query.Tools) error {
	return nil
}
