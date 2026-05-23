package dml

import "github.com/laacin/inyorm/internal/core"

// --- Entity

type Limit struct{ ValueInt int }

// --- Builder

type LimitBuilder struct {
	declared bool
	emb      Limit
}

// --- PUB API

func (b *LimitBuilder) Limit(v int) {
	if v > 0 {
		b.declared = true
		b.emb.ValueInt = v
	}
}

// --- Build

func (*LimitBuilder) Kind() ClauseKind {
	return ClauseLimit
}

func (b *LimitBuilder) IsDeclared() bool {
	return b != nil && b.declared
}

func (b *LimitBuilder) Build(w core.InternalWriter, dial ClauseWriter) error {
	dial.WriteLimit(w, &b.emb)
	return nil
}
