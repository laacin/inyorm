package expr

import "github.com/laacin/inyorm/internal/core"

// --- Entity

type Concat struct{ Values []any }

// start
func (c *Concat) Start(values []any) *Concat {
	c.Values = values
	return c
}

// --- Build
func (c *Concat) Kind() ExprKind {
	return ExprConcat
}

func (c *Concat) Build(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	dial.WriteConcat(w, c, mode)
}
