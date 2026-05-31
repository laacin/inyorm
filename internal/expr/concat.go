package expr

import "github.com/laacin/inyorm/internal/core"

// --- Entity

type Concat struct{ Values []any }

func NewConcat(values []any) *Concat {
	return &Concat{Values: values}
}

// --- Build
func (*Concat) Kind() Kind { return KindConcat }

func (c *Concat) Render(w core.InternalWriter, dial Renderer, mode core.WritingMode) {
	dial.WriteExprConcat(w, c, mode)
}
