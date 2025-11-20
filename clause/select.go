package clause

import "github.com/laacin/inyorm/internal/core"

type Select[Next any] struct {
	declared bool
	distinct bool
	targets  []any
}

func (cls *Select[Next]) Name() core.ClauseType { return core.ClsTypSelect }
func (cls *Select[Next]) IsDeclared() bool      { return cls != nil && cls.declared }
func (cls *Select[Next]) Build(w core.Writer) {
	w.Write("SELECT")
	w.Char(' ')

	if cls.distinct {
		w.Write("DISTINCT ")
	}

	for i, sel := range cls.targets {
		if i > 0 {
			w.Write(", ")
		}
		w.Identifier(sel, cls.Name())
	}
}

// -- Methods

func (cls *Select[Next]) Distinct() Next {
	cls.distinct = true
	return any(cls).(Next)
}

func (cls *Select[Next]) Select(sel ...any) {
	cls.declared = true
	cls.targets = sel
}
