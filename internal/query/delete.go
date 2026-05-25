package query

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/query/dml"
)

type QueryDelete struct {
	dml.ClauseDelete
	dml.ClauseFrom
	dml.ClauseWhere
}

func (q *QueryDelete) Build(b *core.Builder) error {
	if q.ClauseDelete.IsDeclared() {
		q.ClauseDelete.Build(b)
	}

	if q.ClauseFrom.IsDeclared() {
		if tbl, ok := q.ClauseFrom.Value.(*expr.Table); ok {
			b.Attach.MainRef = tbl.Value
		}

		q.ClauseFrom.Build(b)
	}

	if q.ClauseWhere.IsDeclared() {
		q.ClauseWhere.Build(b)
	}

	return nil
}

func (q *QueryDelete) Render(w core.InternalWriter, dial Dialect) error {
	dial.WriteQueryDelete(w, q)
	return nil
}
