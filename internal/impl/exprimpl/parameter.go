package exprimpl

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/expr"
)

type ParameterImpl struct {
	emb expr.Parameter
}

func (p *ParameterImpl) Start(store bool, value any) api.Parameter {
	p.emb.Store = store
	p.emb.Value = value
	return p
}

// --- Build

func (p *ParameterImpl) Kind() expr.ExprKind {
	return expr.ExprParameter
}

func (p *ParameterImpl) Build(w core.InternalWriter, dial expr.ExprWriter, mode core.WritingMode) {
	if p.emb.Store {
		w.PushValue(p.emb.Value)
	}
	w.IncValueCount()
	dial.WritePlaceholder(w)
}
