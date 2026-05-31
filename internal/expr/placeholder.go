package expr

import "github.com/laacin/inyorm/internal/core"

// --- Entity

type Placeholder struct {
	ID       string
	Num      int
	lazy     bool
	paramIdx core.LazyVal[core.ParamIndex]
}

func NewPlaceholder(fn core.LazyVal[core.ParamIndex], lazy ...bool) *Placeholder {
	return &Placeholder{paramIdx: fn, lazy: core.GetLast(false, lazy)}
}
func NewPlaceholderEmpty(idx core.ParamIndex) *Placeholder {
	return &Placeholder{ID: idx.ID, Num: idx.Num}
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
