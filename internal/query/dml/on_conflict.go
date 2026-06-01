package dml

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/query"
)

// --- Entity

type ClauseOnConflict struct {
	declared bool
	Action   OnConflictAction
	On       []any
	Update   []any
}

// --- PUB API

func (c *ClauseOnConflict) OnConflict(ident ...any) api.OnConflictNext {
	c.On = ident
	return c
}

func (c *ClauseOnConflict) DoNothing() {
	c.declared = true
	c.Action = OnConflictDoNothing
}

func (c *ClauseOnConflict) DoUpdate(update ...any) {
	c.declared = true
	c.Action = OnConflictDoUpdate
	c.Update = update
}

// --- Build

func (*ClauseOnConflict) Kind() ClauseKind {
	return ClauseKindOnConflict
}

func (c *ClauseOnConflict) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *ClauseOnConflict) Build(tools *query.Tools) error {
	return nil
}

// helpers

type OnConflictAction int

const (
	OnConflictDoNothing OnConflictAction = iota
	OnConflictDoUpdate
)
