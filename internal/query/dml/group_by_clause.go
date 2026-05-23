package dml

import "github.com/laacin/inyorm/internal/core"

// --- Entity

type GroupBy struct{ Values []any }

// --- Builder

type GroupByBuilder struct {
	declared bool
	emb      GroupBy
}

// --- PUB API

func (b *GroupByBuilder) GroupBy(vals ...any) {
	b.declared = true
	b.emb.Values = vals
}

// --- Build

func (*GroupByBuilder) Kind() ClauseKind {
	return ClauseGroupBy
}

func (b *GroupByBuilder) IsDeclared() bool {
	return b != nil && b.declared
}

func (b *GroupByBuilder) Build(w core.InternalWriter, dial ClauseWriter) error {
	dial.WriteGroupBy(w, &b.emb)
	return nil
}
