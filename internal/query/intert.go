package query

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/query/dml"
)

type QueryInsert struct {
	dml.ClauseInsertInto
}

func (q *QueryInsert) Build(b *core.Builder) error {
	if q.ClauseInsertInto.IsDeclared() {
		if tbl, ok := q.ClauseInsertInto.Table.(*expr.Table); ok {
			b.Attach.MainRef = tbl.Value
		}

		q.ClauseInsertInto.Build(b)
	}

	return nil
}

func (q *QueryInsert) Render(w core.InternalWriter, dial Dialect) error {
	dial.WriteQueryInsert(w, q)
	return nil
}
