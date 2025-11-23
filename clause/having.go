package clause

import (
	"github.com/laacin/inyorm/internal/condition"
	"github.com/laacin/inyorm/internal/core"
)

type Having[Cond, CondNext any] struct {
	declared bool
	cond     *condition.Condition[Cond, CondNext]
}

func (cls *Having[Cond, CondNext]) Name() core.ClauseType { return core.ClsTypHaving }
func (cls *Having[Cond, CondNext]) IsDeclared() bool      { return cls != nil && cls.declared }
func (cls *Having[Cond, CondNext]) Build(w core.Writer, cfg *core.Config) {
	w.Write("HAVING")
	w.Char(' ')
	cls.cond.Build(w, cfg.ColWrite.Having)
}

// -- Methods

func (cls *Having[Cond, CondNext]) Having(on any) Cond {
	cls.declared = true
	cond := &condition.Condition[Cond, CondNext]{}
	cls.cond = cond
	return cond.Start(on)
}
