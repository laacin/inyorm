package query

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/query/dml"
)

type QueryUpdate struct {
	dml.ClauseUpdate
	dml.ClauseWhere
}

func (q *QueryUpdate) Build(b *core.Builder) error {
	if q.ClauseUpdate.IsDeclared() {
		if tbl, ok := q.ClauseUpdate.Table.(*expr.Table); ok {
			b.Attach.MainRef = tbl.Value
		}

		q.ClauseUpdate.Build(b)
	}

	if q.ClauseWhere.IsDeclared() {
		q.ClauseWhere.Build(b)
	}

	return nil
}

func (q *QueryUpdate) Render(w core.InternalWriter, dial Dialect) error {
	dial.WriteQueryUpdate(w, q)
	return nil
}
