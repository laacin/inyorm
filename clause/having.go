package clause

import (
	"github.com/laacin/inyorm/internal/condition"
	"github.com/laacin/inyorm/internal/core"
)

type Having[Cond, CondNext, Ident, Value any] struct {
	declared bool
	cond     *condition.Condition[Cond, CondNext, Ident, Value]
}

func (cls *Having[Cond, CondNext, Ident, Value]) Name() core.ClauseType { return core.ClsTypHaving }
func (cls *Having[Cond, CondNext, Ident, Value]) IsDeclared() bool      { return cls != nil && cls.declared }
func (cls *Having[Cond, CondNext, Ident, Value]) Build(w core.Writer) {
	w.Write("HAVING")
	w.Char(' ')
	cls.cond.Build(w, cls.Name())
}

// -- Methods

func (cls *Having[Cond, CondNext, Ident, Value]) Having(on Ident) Cond {
	cls.declared = true
	cond := &condition.Condition[Cond, CondNext, Ident, Value]{}
	cls.cond = cond
	return cond.Start(on)
}
