package clause

import (
	"github.com/laacin/inyorm/internal/condition"
	"github.com/laacin/inyorm/internal/core"
)

type Having struct {
	cond *condition.Condition
}

func (h *Having) Name() core.ClauseType { return core.ClsTypHaving }
func (h *Having) IsDeclared() bool      { return h != nil }
func (h *Having) Build(w core.Writer) {
	w.Write("HAVING")
	w.Char(' ')
	h.cond.Build(w, core.WriterOpts{ColType: core.ColTypExpr})
}

// -- Methods

func (h *Having) Having(on any) core.Condition {
	cond := &condition.Condition{}
	h.cond = cond
	cond.Start(on)
	return cond
}
