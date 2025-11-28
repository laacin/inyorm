package clause

import "github.com/laacin/inyorm/internal/core"

type Select[Next any] struct {
	declared bool
	distinct bool
	targets  []any
}

func (cls *Select[Next]) Name() string     { return "SELECT" }
func (cls *Select[Next]) IsDeclared() bool { return cls != nil && cls.declared }
func (cls *Select[Next]) Build(w core.Writer, cfg *core.Config) error {
	w.Write("SELECT")
	w.Char(' ')

	if cls.distinct {
		w.Write("DISTINCT ")
	}

	for i, sel := range cls.targets {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(sel, cfg.ColWrite.Select)
	}
	return nil
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
