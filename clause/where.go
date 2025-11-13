package clause

import (
	"github.com/laacin/inyorm/internal/condition"
	"github.com/laacin/inyorm/internal/core"
)

type Where struct {
	exprs []*condition.Condition
}

func (whr *Where) Name() core.ClauseType { return core.ClsTypWhere }
func (whr *Where) IsDeclared() bool      { return whr != nil }
func (whr *Where) Build(w core.Writer) {
	w.Write("WHERE")
	w.Char(' ')
	for i, expr := range whr.exprs {
		if i > 0 {
			w.Write(" AND ")
		}

		expr.Build(
			w,
			core.WhereIdentWriteOpt,
			core.WhereValueWriteOpt,
		)
	}
}

// -- Methods

func (w *Where) Where(identifier any) core.Condition {
	cond := &condition.Condition{}
	w.exprs = append(w.exprs, cond)
	cond.Start(identifier)
	return cond
}
