package dml

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
)

// --- Entity

type Where struct{ Conds []expr.ExprBuilder }

// --- Builder

type WhereBuilder struct {
	declared bool
	emb      Where
	current  *expr.CondBuilder
}

// --- PUB API

func (b *WhereBuilder) Where(ident any) api.Cond {
	b.declared = true
	cond := &expr.CondBuilder{}
	b.emb.Conds = append(b.emb.Conds, cond)
	return cond.Start(ident)
}

// --- Build

func (*WhereBuilder) Kind() ClauseKind {
	return ClauseWhere
}

func (b *WhereBuilder) IsDeclared() bool {
	return b != nil && b.declared
}

func (b *WhereBuilder) Build(w core.InternalWriter, dial ClauseWriter) error {
	dial.WriteWhere(w, &b.emb)
	return nil
}
