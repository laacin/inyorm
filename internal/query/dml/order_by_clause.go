package dml

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
)

// --- Entity

type OrderBy struct{ Segments []OrderSegment }

type OrderByBuilder struct {
	declared bool
	emb      OrderBy
	current  *OrderSegment
}

// --- PUB API

func (b *OrderByBuilder) OrderBy(v any) api.OrderByNext {
	if b.current != nil {
		b.emb.Segments = append(b.emb.Segments, *b.current)
	}

	b.declared = true
	b.current = &OrderSegment{Value: v}
	return b
}

func (b *OrderByBuilder) Desc() {
	b.current.Descending = true
}

// --- Build

func (*OrderByBuilder) Kind() ClauseKind {
	return ClauseOrderBy
}

func (b *OrderByBuilder) IsDeclared() bool {
	return b != nil && b.declared
}

func (b *OrderByBuilder) Build(w core.InternalWriter, dial ClauseWriter) error {
	b.emb.Segments = append(b.emb.Segments, *b.current)
	dial.WriteOrderBy(w, &b.emb)
	return nil
}

// --- Tools

type OrderSegment struct {
	Value      any
	Descending bool
}
