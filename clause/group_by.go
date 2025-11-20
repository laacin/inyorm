package clause

import "github.com/laacin/inyorm/internal/core"

type GroupBy[Ident any] struct {
	declared bool
	groups   []Ident
}

func (cls *GroupBy[Ident]) Name() core.ClauseType { return core.ClsTypGroupBy }
func (cls *GroupBy[Ident]) IsDeclared() bool      { return cls != nil && cls.declared }
func (cls *GroupBy[Ident]) Build(w core.Writer) {
	w.Write("GROUP BY")
	w.Char(' ')
	for i, group := range cls.groups {
		if i > 0 {
			w.Write(", ")
		}
		w.Identifier(group, cls.Name())
	}
}

// -- Methods

func (cls *GroupBy[Ident]) GroupBy(value ...Ident) {
	cls.declared = true
	cls.groups = value
}
