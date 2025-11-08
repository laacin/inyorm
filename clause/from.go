package clause

import "github.com/laacin/inyorm/internal/core"

type FromClause struct {
	value string
}

func (f *FromClause) Name() string {
	return core.ClsFrom
}

func (f *FromClause) Build() core.Builder {
	return func(w core.Writer) {
		w.Write("FROM ")
		w.Table(f.value)
	}
}

// -- Methods

func (f *FromClause) From(from string) {
	f.value = from
}
