package dml

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/query"
)

// --- Entity

type ClauseWhere struct {
	declared bool
	Conds    []expr.Expr
}

// --- PUB API

func (c *ClauseWhere) Where(ident any) api.Cond {
	c.declared = true
	cond := expr.NewCond(ident)
	c.Conds = append(c.Conds, cond)
	return cond
}

// --- Build

func (*ClauseWhere) Kind() ClauseKind {
	return ClauseKindWhere
}

func (c *ClauseWhere) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *ClauseWhere) Build(tools *query.Tools) error {
	return nil
}
