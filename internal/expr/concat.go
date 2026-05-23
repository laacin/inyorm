package expr

import "github.com/laacin/inyorm/internal/core"

// --- Entity

type Concat struct{ Values []any }

// --- Builder

type ConcatBuilder struct{ emb Concat }

// start
func (c *ConcatBuilder) Start(values []any) *ConcatBuilder {
	c.emb.Values = values
	return c
}

// --- Build
func (c *ConcatBuilder) Kind() ExprKind {
	return ExprConcat
}

func (c *ConcatBuilder) Build(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	dial.WriteConcat(w, &c.emb, mode)
}
