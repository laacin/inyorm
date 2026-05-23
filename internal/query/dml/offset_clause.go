package dml

import "github.com/laacin/inyorm/internal/core"

// --- Entity

type Offset struct{ ValueInt int }

// --- Builder

type OffsetBuilder struct {
	declared bool
	emb      Offset
}

// --- PUB API

func (b *OffsetBuilder) Offset(v int) {
	if v > 0 {
		b.declared = true
		b.emb.ValueInt = v
	}
}

// --- Build

func (*OffsetBuilder) Kind() ClauseKind {
	return ClauseOffset
}

func (b *OffsetBuilder) IsDeclared() bool {
	return b != nil && b.declared
}

func (b *OffsetBuilder) Build(w core.InternalWriter, dial ClauseWriter) error {
	dial.WriteOffset(w, &b.emb)
	return nil
}
