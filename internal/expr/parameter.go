package expr

import "github.com/laacin/inyorm/internal/core"

// --- Entity

type Param struct {
	Value any
	Store bool
}

// start

func (p *Param) Start(store bool, value any) *Param {
	p.Store = store
	p.Value = value
	return p
}

// --- Build

func (p *Param) Kind() ExprKind {
	return ExprKindParam
}

func (p *Param) Render(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	if p.Store {
		w.PushValue(p.Value)
	}
	w.IncValueCount()
	dial.WriteExprPlaceholder(w)
}
