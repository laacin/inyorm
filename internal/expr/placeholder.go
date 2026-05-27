package expr

import "github.com/laacin/inyorm/internal/core"

// --- Entity

type Placeholder struct {
	ID       string
	Num      int
	lazy     bool
	paramIdx core.LazyVal[core.ParamIndex]
}

// start

func (p *Placeholder) Start(fn func() core.ParamIndex) *Placeholder {
	p.paramIdx = fn
	return p
}

func (p *Placeholder) StartEmpty(idx core.ParamIndex) *Placeholder {
	p.ID = idx.ID
	p.Num = idx.Num
	return p
}

func (p *Placeholder) StartLazy(fn func() core.ParamIndex) *Placeholder {
	p.lazy = true
	p.paramIdx = fn
	return p
}

// --- Build
func (p *Placeholder) IsLazy() bool {
	return p.lazy
}

func (p *Placeholder) Kind() Kind { return KindPlaceholder }

func (p *Placeholder) Render(w core.InternalWriter, dial Renderer, mode core.WritingMode) {
	if p.paramIdx != nil {
		idx := p.paramIdx()
		p.ID = idx.ID
		p.Num = idx.Num
	}

	dial.WriteExprPlaceholder(w, p)
}
