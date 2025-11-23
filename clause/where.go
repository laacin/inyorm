package clause

import (
	"github.com/laacin/inyorm/internal/condition"
	"github.com/laacin/inyorm/internal/core"
)

type Where[Cond, CondNext any] struct {
	declared bool
	exprs    []*condition.Condition[Cond, CondNext]
}

func (cls *Where[Cond, CondNext]) Name() core.ClauseType { return core.ClsTypWhere }
func (cls *Where[Cond, CondNext]) IsDeclared() bool      { return cls != nil && cls.declared }
func (cls *Where[Cond, CondNext]) Build(w core.Writer, cfg *core.Config) {
	w.Write("WHERE")
	w.Char(' ')
	for i, expr := range cls.exprs {
		if i > 0 {
			w.Write(" AND ")
		}
		expr.Build(w, cfg.ColWrite.Where)
	}
}

// -- Methods

func (cls *Where[Cond, CondNext]) Where(identifier any) Cond {
	cls.declared = true
	cond := &condition.Condition[Cond, CondNext]{}
	cls.exprs = append(cls.exprs, cond)
	return Cond(cond.Start(identifier))
}
