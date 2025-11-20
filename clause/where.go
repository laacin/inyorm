package clause

import (
	"github.com/laacin/inyorm/internal/condition"
	"github.com/laacin/inyorm/internal/core"
)

type Where[Cond, CondNext, Ident, Value any] struct {
	declared bool
	exprs    []*condition.Condition[Cond, CondNext, Ident, Value]
}

func (cls *Where[Cond, CondNext, Ident, Value]) Name() core.ClauseType { return core.ClsTypWhere }
func (cls *Where[Cond, CondNext, Ident, Value]) IsDeclared() bool      { return cls != nil && cls.declared }
func (cls *Where[Cond, CondNext, Ident, Value]) Build(w core.Writer) {
	w.Write("WHERE")
	w.Char(' ')
	for i, expr := range cls.exprs {
		if i > 0 {
			w.Write(" AND ")
		}
		expr.Build(w, cls.Name())
	}
}

// -- Methods

func (cls *Where[Cond, CondNext, Ident, Value]) Where(identifier Ident) Cond {
	cls.declared = true
	cond := &condition.Condition[Cond, CondNext, Ident, Value]{}
	cls.exprs = append(cls.exprs, cond)
	return Cond(cond.Start(identifier))
}
