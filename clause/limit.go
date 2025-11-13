package clause

import "github.com/laacin/inyorm/internal/core"

type Limit struct {
	limit int
}

func (l *Limit) Name() core.ClauseType { return core.ClsTypLimit }
func (l *Limit) IsDeclared() bool      { return l != nil }
func (l *Limit) Build(w core.Writer) {
	if l.limit > 0 {
		w.Write("LIMIT")
		w.Char(' ')
		w.Value(l.limit, core.LimitWriteOpt)
	}
}

// -- Methods

func (l *Limit) Limit(value int) {
	if value > 0 {
		l.limit = value
	}
}
