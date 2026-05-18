package exprimpl

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/expr"
)

type ConcatImpl struct {
	emb expr.Concat
}

func (c *ConcatImpl) Start(values []any) *ConcatImpl {
	c.emb.Values = values
	return c
}

// --- Build
func (c *ConcatImpl) Kind() expr.ExprKind {
	return expr.ExprConcat
}

func (c *ConcatImpl) Build(w core.InternalWriter, dial expr.ExprWriter, mode core.WritingMode) {
	dial.WriteConcat(w, &c.emb, mode)
}
