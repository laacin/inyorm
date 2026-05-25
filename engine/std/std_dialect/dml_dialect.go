package std_dialect

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/query"
	"github.com/laacin/inyorm/internal/query/dml"
)

// ---- CLAUSES -----

func (*Dialect) WriteClauseInsertInto(w core.Writer, cls *dml.ClauseInsertInto) {
	w.Write("INSERT INTO")
	w.Char(' ')

	w.Value(cls.Table, core.WriteBase)
	w.Char(' ')

	w.Char('(')
	for i, col := range cls.Cols {
		if i > 0 {
			w.Write(", ")
		}
		w.Write(col)
	}
	w.Char(')')

	perRow := len(cls.Vals) / cls.Rows

	w.Write(" VALUES ")
	for row := range cls.Rows {
		if row > 0 {
			w.Write(", ")
		}

		w.Char('(')
		for ci := range cls.Cols {
			if ci > 0 {
				w.Write(", ")
			}
			w.Value(cls.Vals[row*perRow+ci], core.WriteBase)
		}
		w.Char(')')
	}
}

func (*Dialect) WriteClauseSelect(w core.Writer, cls *dml.ClauseSelect) {
	w.Write("SELECT")
	w.Char(' ')

	if cls.Dist {
		w.Write("DISTINCT")
		w.Char(' ')
	}

	for i, val := range cls.Values {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(val, core.WriteDef)
	}
}

func (*Dialect) WriteClauseFrom(w core.Writer, cls *dml.ClauseFrom) {
	w.Write("FROM")
	w.Char(' ')
	w.Value(cls.Value, core.WriteDef)
}

func (*Dialect) WriteClauseJoin(w core.Writer, cls *dml.ClauseJoin) {
	for i, join := range cls.Segments {
		if i > 0 {
			w.Char(' ')
		}

		w.Write(joinKindMap[join.Kind])
		w.Write(" JOIN ")
		w.Value(join.Table, core.WriteDef)

		if join.Cond != nil {
			w.Write(" ON ")
			w.Value(join.Cond, core.WriteBase)
		}
	}
}

func (*Dialect) WriteClauseWhere(w core.Writer, cls *dml.ClauseWhere) {
	w.Write("WHERE")
	w.Char(' ')

	for i, cond := range cls.Conds {
		if i > 0 {
			w.Write(" AND ")
		}
		w.Value(cond, core.WriteExpr)
	}
}

func (*Dialect) WriteClauseGroupBy(w core.Writer, cls *dml.ClauseGroupBy) {
	w.Write("GROUP BY")
	w.Char(' ')

	for i, group := range cls.Values {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(group, core.WriteExpr)
	}
}

func (*Dialect) WriteClauseHaving(w core.Writer, cls *dml.ClauseHaving) {
	w.Write("HAVING")
	w.Char(' ')
	w.Value(cls.Cond, core.WriteExpr)
}

func (*Dialect) WriteClauseOrderBy(w core.Writer, cls *dml.ClauseOrderBy) {
	w.Write("ORDER BY")
	w.Char(' ')

	for i, ord := range cls.Segments {
		if i > 0 {
			w.Write(", ")
		}

		w.Value(ord.Value, core.WriteAlias)
		if ord.Descending {
			w.Char(' ')
			w.Write("DESC")
		}
	}
}

func (*Dialect) WriteClauseLimit(w core.Writer, cls *dml.ClauseLimit) {
	w.Write("LIMIT")
	w.Char(' ')
	w.Value(cls.ValueInt, core.WriteBase)
}

func (*Dialect) WriteClauseOffset(w core.Writer, cls *dml.ClauseOffset) {
	w.Write("OFFSET")
	w.Char(' ')
	w.Value(cls.ValueInt, core.WriteBase)
}

func (*Dialect) WriteClauseUpdate(w core.Writer, cls *dml.ClauseUpdate) {
	w.Write("UPDATE")
	w.Char(' ')

	w.Value(cls.Table, core.WriteBase)
	w.Write(" SET ")

	for i, col := range cls.Cols {
		if i > 0 {
			w.Write(", ")
		}

		w.Write(col)
		w.Write(" = ")
		w.Value(cls.Vals[i], core.WriteBase)
	}
}

func (*Dialect) WriteClauseDelete(w core.Writer, cls *dml.ClauseDelete) {
	w.Write("DELETE")
}

// ---- QUERY WRITER -----

func (s *Dialect) WriteQuerySelect(w core.Writer, q *query.QuerySelect) {
	if q.ClauseSelect.IsDeclared() {
		s.Self.WriteClauseSelect(w, &q.ClauseSelect)
	}

	if q.ClauseFrom.IsDeclared() {
		w.Char(' ')
		s.Self.WriteClauseFrom(w, &q.ClauseFrom)
	}

	if q.ClauseJoin.IsDeclared() {
		w.Char(' ')
		s.Self.WriteClauseJoin(w, &q.ClauseJoin)
	}

	if q.ClauseWhere.IsDeclared() {
		w.Char(' ')
		s.Self.WriteClauseWhere(w, &q.ClauseWhere)
	}

	if q.ClauseGroupBy.IsDeclared() {
		w.Char(' ')
		s.Self.WriteClauseGroupBy(w, &q.ClauseGroupBy)
	}

	if q.ClauseHaving.IsDeclared() {
		w.Char(' ')
		s.Self.WriteClauseHaving(w, &q.ClauseHaving)
	}

	if q.ClauseOrderBy.IsDeclared() {
		w.Char(' ')
		s.Self.WriteClauseOrderBy(w, &q.ClauseOrderBy)
	}

	if q.ClauseLimit.IsDeclared() {
		w.Char(' ')
		s.Self.WriteClauseLimit(w, &q.ClauseLimit)
	}

	if q.ClauseOffset.IsDeclared() {
		w.Char(' ')
		s.Self.WriteClauseOffset(w, &q.ClauseOffset)
	}
}

func (s *Dialect) WriteQueryInsert(w core.Writer, q *query.QueryInsert) {
	if q.ClauseInsertInto.IsDeclared() {
		s.Self.WriteClauseInsertInto(w, &q.ClauseInsertInto)
	}
}

func (s *Dialect) WriteQueryUpdate(w core.Writer, q *query.QueryUpdate) {
	if q.ClauseUpdate.IsDeclared() {
		s.Self.WriteClauseUpdate(w, &q.ClauseUpdate)
	}

	if q.ClauseWhere.IsDeclared() {
		w.Char(' ')
		s.Self.WriteClauseWhere(w, &q.ClauseWhere)
	}
}

func (s *Dialect) WriteQueryDelete(w core.Writer, q *query.QueryDelete) {
	if q.ClauseDelete.IsDeclared() {
		s.Self.WriteClauseDelete(w, &q.ClauseDelete)
	}

	if q.ClauseFrom.IsDeclared() {
		w.Char(' ')
		s.Self.WriteClauseFrom(w, &q.ClauseFrom)
	}

	if q.ClauseWhere.IsDeclared() {
		w.Char(' ')
		s.Self.WriteClauseWhere(w, &q.ClauseWhere)
	}
}

// --- Helpers
var joinKindMap = map[dml.JoinKind]string{
	dml.JoinInner: "INNER",
	dml.JoinLeft:  "LEFT",
	dml.JoinFull:  "FULL",
	dml.JoinCross: "CROSS",
}
