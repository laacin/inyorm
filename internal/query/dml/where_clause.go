package dml

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/expr"
)

// --- Entity

type ClauseWhere struct {
	declared bool
	Conds    []expr.ExprBuilder
	current  expr.ExprBuilder
}

// --- PUB API

func (c *ClauseWhere) Where(ident any) api.Cond {
	c.declared = true
	cond := &expr.Cond{}
	c.Conds = append(c.Conds, cond)
	return cond.Start(ident)
}

// --- Build

func (*ClauseWhere) Kind() ClauseKind {
	return ClauseKindWhere
}

func (c *ClauseWhere) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *ClauseWhere) Build() error {
	return nil
}
