package clause

import "github.com/laacin/inyorm/internal/core"

type FromClause struct {
	value string
}

func (f *FromClause) Name() core.ClauseType {
	return core.ClsTypFrom
}

func (f *FromClause) IsDeclared() bool { return f != nil }

func (f *FromClause) Build(w core.Writer) {
	w.Write("FROM ")
	w.Table(f.value)
}

// -- Methods

func (f *FromClause) From(from string) {
	f.value = from
}
