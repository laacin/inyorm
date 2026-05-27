package expr

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
)

// --- Entity

type Case struct {
	Cond    any
	Whens   []CaseWhen
	Els     any
	kind    Kind
	current *CaseWhen
}

// start

func (c *Case) StartSwitch(cond any) *Case {
	c.kind = KindCaseSwitch
	c.Cond = cond
	return c
}
func (c *Case) StartSearch() *Case {
	c.kind = KindCaseSearch
	return c
}

// --- PUB API

func (c *Case) When(when any) api.CaseNext {
	c.current = &CaseWhen{When: when}
	return c
}

func (c *Case) Then(then any) api.Case {
	c.current.Then = then
	c.Whens = append(c.Whens, *c.current)
	return c
}

func (c *Case) Else(els any) {
	c.Els = els
}

// --- Build

func (c *Case) Kind() Kind { return c.kind }

func (c *Case) Render(w core.InternalWriter, dial Renderer, mode core.WritingMode) {
	switch c.kind {
	case KindCaseSwitch:
		dial.WriteExprCaseSwitch(w, c, mode)
	case KindCaseSearch:
		dial.WriteExprCaseSearch(w, c, mode)
	}
}

// --- Tools
type CaseWhen struct {
	When any
	Then any
}
