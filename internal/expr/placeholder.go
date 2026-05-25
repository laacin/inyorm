package expr

import "github.com/laacin/inyorm/internal/core"

// --- Entity

type Placeholder struct {
	ID       string
	Num      int
	lazy     bool
	onRender func() core.ParamIndex
}

// start

func (p *Placeholder) Start(fn func() core.ParamIndex) *Placeholder {
	p.onRender = fn
	return p
}

func (p *Placeholder) StartEmpty(idx core.ParamIndex) *Placeholder {
	p.ID = idx.ID
	p.Num = idx.Num
	return p
}

func (p *Placeholder) StartLazy(fn func() core.ParamIndex) *Placeholder {
	p.lazy = true
	p.onRender = fn
	return p
}

// --- Build
func (p *Placeholder) IsLazy() bool {
	return p.lazy
}

func (p *Placeholder) Kind() ExprKind {
	return ExprKindPlaceholder
}

func (p *Placeholder) Render(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	if p.onRender != nil {
		idx := p.onRender()
		p.ID = idx.ID
		p.Num = idx.Num
	}

	dial.WriteExprPlaceholder(w, p)
}
