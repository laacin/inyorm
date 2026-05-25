package dml

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
)

// --- Entity

type ClauseJoin struct {
	declared bool
	Segments []JoinSegment
	current  *JoinSegment
}

// --- PUB API

func (c *ClauseJoin) Join(v any) api.JoinNext {
	c.declared = true
	c.current = &JoinSegment{
		Kind:  JoinInner,
		Table: v,
	}
	return c
}

func (c *ClauseJoin) Left() api.JoinEnd {
	c.current.Kind = JoinLeft
	return c
}
func (c *ClauseJoin) Full() api.JoinEnd {
	c.current.Kind = JoinFull
	return c
}
func (c *ClauseJoin) Cross() {
	c.current.Kind = JoinCross
	c.Segments = append(c.Segments, *c.current)
}

func (c *ClauseJoin) On(ident any) api.Cond {
	cond := &expr.Cond{}
	c.current.Cond = cond
	c.Segments = append(c.Segments, *c.current)
	return cond.Start(ident)
}

// --- Build

func (*ClauseJoin) Kind() ClauseKind {
	return ClauseKindJoin
}

func (c *ClauseJoin) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *ClauseJoin) Build(b *core.Builder) error {
	return nil
}

// --- Tools
type JoinSegment struct {
	Kind  JoinKind
	Table any
	Cond  expr.Expr
}

type JoinKind int

const (
	JoinInner JoinKind = iota
	JoinLeft
	JoinFull
	JoinCross
)
