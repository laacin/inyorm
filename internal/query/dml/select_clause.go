package dml

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
)

// --- Entity
type Select struct {
	Distinct bool
	Values   []any
}

// --- Builder

type SelectBuilder struct {
	declared bool
	emb      Select
}

// --- PUB API

func (b *SelectBuilder) Select(vals ...any) api.SelectNext {
	b.declared = true
	b.emb.Values = vals
	return b
}

func (b *SelectBuilder) Distinct() {
	b.emb.Distinct = true
}

// --- Build

func (*SelectBuilder) Kind() ClauseKind {
	return ClauseSelect
}

func (b *SelectBuilder) IsDeclared() bool {
	return b != nil && b.declared
}

func (b *SelectBuilder) Build(w core.InternalWriter, dial ClauseWriter) error {
	dial.WriteSelect(w, &b.emb)
	return nil
}
