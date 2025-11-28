package clause

import "github.com/laacin/inyorm/internal/core"

type From struct {
	declared bool
	value    string
}

func (cls *From) Name() string     { return "FROM" }
func (cls *From) IsDeclared() bool { return cls != nil && cls.declared }
func (cls *From) Build(w core.Writer, cfg *core.Config) error {
	w.Write("FROM")
	w.Char(' ')
	w.Table(cls.value)
	return nil
}

// -- Methods

func (cls *From) From(from string) {
	cls.declared = true
	cls.value = from
}
