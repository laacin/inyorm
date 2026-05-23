package expr

import "github.com/laacin/inyorm/internal/core"

// --- Entity

type Param struct {
	Value any
	Store bool
}

// --- Builder

type ParamBuilder struct{ emb Param }

// start

func (p *ParamBuilder) Start(store bool, value any) *ParamBuilder {
	p.emb.Store = store
	p.emb.Value = value
	return p
}

// --- Build

func (p *ParamBuilder) Kind() ExprKind {
	return ExprParam
}

func (p *ParamBuilder) Build(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	if p.emb.Store {
		w.PushValue(p.emb.Value)
	}
	w.IncValueCount()
	dial.WritePlaceholder(w)
}
