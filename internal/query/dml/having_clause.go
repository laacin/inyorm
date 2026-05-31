package dml

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/query"
)

// --- Entity

type ClauseHaving struct {
	declared bool
	Cond     expr.Expr
}

// --- PUB API

func (c *ClauseHaving) Having(ident any) api.Cond {
	c.declared = true
	cond := expr.NewCond(ident)
	c.Cond = cond
	return cond
}

// --- Build

func (*ClauseHaving) Kind() ClauseKind {
	return ClauseKindHaving
}

func (c *ClauseHaving) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *ClauseHaving) Build(tools *query.Tools) error {
	return nil
}
