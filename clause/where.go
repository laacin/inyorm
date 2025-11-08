package clause

import (
	"github.com/laacin/inyorm/internal/condition"
	"github.com/laacin/inyorm/internal/core"
)

type WhereClause struct {
	exprs []*condition.Condition
}

func (w *WhereClause) Name() string {
	return core.ClsWhere
}

func (wc *WhereClause) Build() core.Builder {
	return func(w core.Writer) {
		w.Write("WHERE ")
		for i, expr := range wc.exprs {
			if i > 0 {
				w.Write(" AND ")
			}
			expr.Build(w, &core.ValueOpts{Placeholder: true})
		}
	}
}

// -- Methods

func (w *WhereClause) Where(identifier any) core.Cond {
	cond := &condition.Condition{}
	w.exprs = append(w.exprs, cond)
	return cond.Start(identifier)
}
