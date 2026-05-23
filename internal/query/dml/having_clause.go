package dml

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
)

// --- Entity

type ClauseHaving struct {
	declared bool
	Cond     expr.ExprBuilder
}

// --- PUB API

func (c *ClauseHaving) Having(ident any) api.Cond {
	c.declared = true
	cond := &expr.Cond{}
	c.Cond = cond
	return cond.Start(ident)
}

// --- Build

func (*ClauseHaving) Kind() ClauseKind {
	return ClauseKindHaving
}

func (c *ClauseHaving) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *ClauseHaving) Build(w core.InternalWriter, dial ClauseWriter) error {
	dial.WriteHaving(w, c)
	return nil
}
