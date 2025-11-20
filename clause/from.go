package clause

import "github.com/laacin/inyorm/internal/core"

type From struct {
	declared bool
	value    string
}

func (cls *From) Name() core.ClauseType { return core.ClsTypFrom }
func (cls *From) IsDeclared() bool      { return cls != nil && cls.declared }
func (cls *From) Build(w core.Writer) {
	w.Write("FROM")
	w.Char(' ')
	w.Table(cls.value)
}

// -- Methods

func (cls *From) From(from string) {
	cls.declared = true
	cls.value = from
}
