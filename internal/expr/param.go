package expr

import "github.com/laacin/inyorm/internal/core"

// --- Entity

type Param struct {
	lazy  bool
	Value any
	Ref   string
	Cols  []string
}

// start
func (p *Param) Start(value any) *Param {
	p.Value = value
	return p
}

func (p *Param) Lazy(ref string) *Param {
	p.Ref = ref
	return p
}

// --- Build

func (p *Param) Kind() ExprKind {
	return ExprKindParam
}

func (p *Param) Render(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	w.PushValue(p.Value)
	dial.WriteExprPlaceholder(w)
}
