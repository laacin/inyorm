package clause

import "github.com/laacin/inyorm/internal/core"

type Select[Next, Ident any] struct {
	declared bool
	distinct bool
	targets  []Ident
}

func (cls *Select[Next, Ident]) Name() core.ClauseType { return core.ClsTypSelect }
func (cls *Select[Next, Ident]) IsDeclared() bool      { return cls != nil && cls.declared }
func (cls *Select[Next, Ident]) Build(w core.Writer) {
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

func (cls *Select[Next, Ident]) Distinct() Next {
	cls.distinct = true
	return any(cls).(Next)
}

func (cls *Select[Next, Ident]) Select(sel ...Ident) {
	cls.declared = true
	cls.targets = sel
}
