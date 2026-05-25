package query

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/query/dml"
)

type QuerySelect struct {
	dml.ClauseSelect
	dml.ClauseFrom
	dml.ClauseJoin
	dml.ClauseWhere
	dml.ClauseGroupBy
	dml.ClauseHaving
	dml.ClauseOrderBy
	dml.ClauseLimit
	dml.ClauseOffset
}

func (q *QuerySelect) Build(b *core.Builder) error {
	if q.ClauseSelect.IsDeclared() {
		q.ClauseSelect.Build(b)
	}

	if q.ClauseFrom.IsDeclared() {
		if tbl, ok := q.ClauseFrom.Value.(*expr.Table); ok {
			b.Attach.MainRef = tbl.Value
		}

		q.ClauseFrom.Build(b)
	}

	if q.ClauseJoin.IsDeclared() {
		b.Attach.UseAliases = true
		q.ClauseJoin.Build(b)
	}

	if q.ClauseWhere.IsDeclared() {
		q.ClauseWhere.Build(b)
	}

	if q.ClauseGroupBy.IsDeclared() {
		q.ClauseGroupBy.Build(b)
	}

	if q.ClauseHaving.IsDeclared() {
		q.ClauseHaving.Build(b)
	}

	if q.ClauseOrderBy.IsDeclared() {
		q.ClauseOrderBy.Build(b)
	}

	if q.ClauseLimit.IsDeclared() {
		q.ClauseLimit.Build(b)
	}

	if q.ClauseOffset.IsDeclared() {
		q.ClauseOffset.Build(b)
	}

	return nil
}

func (q *QuerySelect) Render(w core.InternalWriter, dial Dialect) error {
	dial.WriteQuerySelect(w, q)
	return nil
}
