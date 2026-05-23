package expr

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
)

// --- Entity

type Case struct {
	Cond  any
	Whens []CaseWhen
	Els   any
}

// --- Builder
type CaseBuilder struct {
	kind    ExprKind
	emb     Case
	current *CaseWhen
}

// start
func (c *CaseBuilder) StartSwitch(cond any) *CaseBuilder {
	c.kind = ExprCaseSwitch
	c.emb.Cond = cond
	return c
}
func (c *CaseBuilder) StartSearch() *CaseBuilder {
	c.kind = ExprCaseSearch
	return c
}

// --- PUB API

func (c *CaseBuilder) When(when any) api.CaseNext {
	c.current = &CaseWhen{When: when}
	return c
}

func (c *CaseBuilder) Then(then any) api.Case {
	c.current.Then = then
	c.emb.Whens = append(c.emb.Whens, *c.current)
	return c
}

func (c *CaseBuilder) Else(els any) {
	c.emb.Els = els
}

// --- Build

func (c *CaseBuilder) Kind() ExprKind {
	return c.kind
}

func (c *CaseBuilder) Build(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	switch c.kind {
	case ExprCaseSwitch:
		dial.WriteCaseSwitch(w, &c.emb, mode)
	case ExprCaseSearch:
		dial.WriteCaseSearch(w, &c.emb, mode)
	}
}

// --- Tools
type CaseWhen struct {
	When any
	Then any
}
