package clause

import (
	"github.com/laacin/inyorm/internal/condition"
	"github.com/laacin/inyorm/internal/core"
)

type WhereClause struct {
	exprs []*condition.Condition
}

func (w *WhereClause) Name() core.ClauseType {
	return core.ClsTypWhere
}

func (w *WhereClause) IsDeclared() bool { return w != nil }

func (wc *WhereClause) Build(w core.Writer) {
	w.Write("WHERE ")
	for i, expr := range wc.exprs {
		if i > 0 {
			w.Write(" AND ")
		}
		expr.Build(w, core.WriterOpts{
			Placeholder: true,
			ColType:     core.ColTypExpr,
		})
	}
}

// -- Methods

func (w *WhereClause) Where(identifier any) core.Cond {
	cond := &condition.Condition{}
	w.exprs = append(w.exprs, cond)
	return cond.Start(identifier)
}
