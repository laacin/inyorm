package clause

import "github.com/laacin/inyorm/internal/core"

type From struct {
	value string
}

func (f *From) Name() core.ClauseType { return core.ClsTypFrom }
func (f *From) IsDeclared() bool      { return f != nil }
func (f *From) Build(w core.Writer) {
	w.Write("FROM")
	w.Char(' ')
	w.Table(f.value)
}

// -- Methods

func (f *From) From(from string) {
	f.value = from
}
