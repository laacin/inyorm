package dml

import "github.com/laacin/inyorm/internal/core"

// --- Entity

type From struct{ Value any }

// --- Builder

type FromBuilder struct {
	declared bool
	emb      From
}

// --- PUB API

func (b *FromBuilder) From(from any) {
	b.declared = true
	b.emb.Value = from
}

// --- Build

func (*FromBuilder) Kind() ClauseKind {
	return ClauseFrom
}

func (b *FromBuilder) IsDeclared() bool {
	return b != nil && b.declared
}

func (b *FromBuilder) Build(w core.InternalWriter, dial ClauseWriter) error {
	dial.WriteFrom(w, &b.emb)
	return nil
}
