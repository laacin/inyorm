package clause

import "github.com/laacin/inyorm/internal/core"

type LimitClause struct {
	limit int
}

func (l *LimitClause) Name() core.ClauseType {
	return core.ClsTypLimit
}

func (l *LimitClause) IsDeclared() bool { return l != nil }

func (l *LimitClause) Build(w core.Writer) {
	if l.limit > 0 {
		w.Write("LIMIT ")
		w.Value(l.limit, core.WriterOpts{})
	}
}

// -- Methods

func (l *LimitClause) Limit(value int) {
	if value > 0 {
		l.limit = value
	}
}
