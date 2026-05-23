package dml

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
)

// --- Entity

type Having struct{ Cond expr.ExprBuilder }

// --- Builder

type HavingBuilder struct {
	declared bool
	emb      Having
}

// --- PUB API

func (b *HavingBuilder) Having(ident any) api.Cond {
	b.declared = true
	cond := &expr.CondBuilder{}
	b.emb.Cond = cond
	return cond.Start(ident)
}

// --- Build

func (*HavingBuilder) Kind() ClauseKind {
	return ClauseHaving
}

func (b *HavingBuilder) IsDeclared() bool {
	return b != nil && b.declared
}

func (b *HavingBuilder) Build(w core.InternalWriter, dial ClauseWriter) error {
	dial.WriteHaving(w, &b.emb)
	return nil
}
