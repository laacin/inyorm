package dml

import "github.com/laacin/inyorm/internal/core"

// --- Entity

type Delete struct{}

// --- Builder

type DeleteBuilder struct {
	declared bool
	emb      Delete
}

// --- PUB API

func (b *DeleteBuilder) Delete() {
	b.declared = true
}

// --- Build

func (*DeleteBuilder) Kind() ClauseKind {
	return ClauseDelete
}

func (b *DeleteBuilder) IsDeclared() bool {
	return b != nil && b.declared
}

func (b *DeleteBuilder) Build(w core.InternalWriter, dial ClauseWriter) error {
	dial.WriteDelete(w, &b.emb)
	return nil
}
