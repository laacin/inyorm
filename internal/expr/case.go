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
	kind    ExprKind
	current *CaseWhen
}

// start

func (c *Case) StartSwitch(cond any) *Case {
	c.kind = ExprKindCaseSwitch
	c.Cond = cond
	return c
}
func (c *Case) StartSearch() *Case {
	c.kind = ExprKindCaseSearch
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

func (c *Case) Kind() ExprKind {
	return c.kind
}

func (c *Case) Render(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	switch c.kind {
	case ExprKindCaseSwitch:
		dial.WriteExprCaseSwitch(w, c, mode)
	case ExprKindCaseSearch:
		dial.WriteExprCaseSearch(w, c, mode)
	}
}

// --- Tools
type CaseWhen struct {
	When any
	Then any
}
