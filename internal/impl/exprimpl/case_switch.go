package exprimpl

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/expr"
)

type CaseSwitchImpl struct {
	emb     expr.CaseSwitch
	current *expr.CaseWhen
}

func (c *CaseSwitchImpl) Start(cond any) api.Case {
	c.emb.Cond = cond
	return c
}

func (c *CaseSwitchImpl) When(when any) api.CaseNext {
	c.current = &expr.CaseWhen{When: when}
	return c
}

func (c *CaseSwitchImpl) Then(then any) api.Case {
	c.current.Then = then
	c.emb.Whens = append(c.emb.Whens, *c.current)
	return c
}

func (c *CaseSwitchImpl) Else(els any) {
	c.emb.Els = els
}

// --- Build
func (c *CaseSwitchImpl) Kind() expr.ExprKind {
	return expr.ExprCaseSwitch
}

func (c *CaseSwitchImpl) Build(w core.InternalWriter, dial expr.ExprWriter, mode core.WritingMode) {
	dial.WriteCaseSwitch(w, &c.emb, mode)
}
