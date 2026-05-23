package dml

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
)

// --- Entity

type Join struct{ Segments []JoinSegment }

// --- Builder

type JoinBuilder struct {
	declared bool
	emb      Join
	current  *JoinSegment
}

// --- PUB API

func (b *JoinBuilder) Join(v any) api.JoinNext {
	b.declared = true
	b.current = &JoinSegment{
		Kind:  JoinInner,
		Table: v,
	}
	return b
}

func (b *JoinBuilder) Left() api.JoinEnd {
	b.current.Kind = JoinLeft
	return b
}
func (b *JoinBuilder) Full() api.JoinEnd {
	b.current.Kind = JoinFull
	return b
}
func (b *JoinBuilder) Cross() {
	b.current.Kind = JoinCross
	b.emb.Segments = append(b.emb.Segments, *b.current)
}

func (b *JoinBuilder) On(ident any) api.Cond {
	cond := &expr.CondBuilder{}
	b.current.Cond = cond
	b.emb.Segments = append(b.emb.Segments, *b.current)
	return cond.Start(ident)
}

// --- Build

func (*JoinBuilder) Kind() ClauseKind {
	return ClauseJoin
}

func (b *JoinBuilder) IsDeclared() bool {
	return b != nil && b.declared
}

func (b *JoinBuilder) Build(w core.InternalWriter, dial ClauseWriter) error {
	dial.WriteJoin(w, &b.emb)
	return nil
}

// --- Tools
type JoinSegment struct {
	Kind  JoinKind
	Table any
	Cond  expr.ExprBuilder
}

type JoinKind int

const (
	JoinInner JoinKind = iota
	JoinLeft
	JoinFull
	JoinCross
)
