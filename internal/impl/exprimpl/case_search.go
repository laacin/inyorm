package exprimpl

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/expr"
)

type CaseSearchImpl struct {
	emb     expr.CaseSearch
	current *expr.CaseWhen
}

func (c *CaseSearchImpl) When(when any) api.CaseNext {
	c.current = &expr.CaseWhen{When: when}
	return c
}

func (c *CaseSearchImpl) Then(then any) api.Case {
	c.current.Then = then
	c.emb.Whens = append(c.emb.Whens, *c.current)
	return c
}

func (c *CaseSearchImpl) Else(els any) {
	c.emb.Els = els
}

// --- Build
func (c *CaseSearchImpl) Kind() expr.ExprKind {
	return expr.ExprCaseSearch
}

func (c *CaseSearchImpl) Build(w core.InternalWriter, dial expr.ExprWriter, mode core.WritingMode) {
	dial.WriteCaseSearch(w, &c.emb, mode)
}
