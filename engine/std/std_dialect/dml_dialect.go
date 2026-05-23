package std_dialect

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/query/dml"
)

// ---- CLAUSES -----

func (*Dialect) WriteInsertInto(w core.Writer, cls *dml.ClauseInsert) {
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

func (*Dialect) WriteSelect(w core.Writer, cls *dml.ClauseSelect) {
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

func (*Dialect) WriteFrom(w core.Writer, cls *dml.ClauseFrom) {
	w.Write("FROM")
	w.Char(' ')
	w.Value(cls.Value, core.WriteDef)
}

var joinKindMap = map[dml.JoinKind]string{
	dml.JoinInner: "INNER",
	dml.JoinLeft:  "LEFT",
	dml.JoinFull:  "FULL",
	dml.JoinCross: "CROSS",
}

func (*Dialect) WriteJoin(w core.Writer, cls *dml.ClauseJoin) {
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

func (*Dialect) WriteWhere(w core.Writer, cls *dml.ClauseWhere) {
	w.Write("WHERE")
	w.Char(' ')

	for i, cond := range cls.Conds {
		if i > 0 {
			w.Write(" AND ")
		}
		w.Value(cond, core.WriteExpr)
	}
}

func (*Dialect) WriteGroupBy(w core.Writer, cls *dml.ClauseGroupBy) {
	w.Write("GROUP BY")
	w.Char(' ')

	for i, group := range cls.Values {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(group, core.WriteExpr)
	}
}

func (*Dialect) WriteHaving(w core.Writer, cls *dml.ClauseHaving) {
	w.Write("HAVING")
	w.Char(' ')
	w.Value(cls.Cond, core.WriteExpr)
}

func (*Dialect) WriteOrderBy(w core.Writer, cls *dml.ClauseOrderBy) {
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

func (*Dialect) WriteLimit(w core.Writer, cls *dml.ClauseLimit) {
	w.Write("LIMIT")
	w.Char(' ')
	w.Value(cls.ValueInt, core.WriteBase)
}

func (*Dialect) WriteOffset(w core.Writer, cls *dml.ClauseOffset) {
	w.Write("OFFSET")
	w.Char(' ')
	w.Value(cls.ValueInt, core.WriteBase)
}

func (*Dialect) WriteUpdate(w core.Writer, cls *dml.ClauseUpdate) {
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

func (*Dialect) WriteDelete(w core.Writer, cls *dml.ClauseDelete) {
	w.Write("DELETE")
}

// ---- QUERY WRITER -----

func (s *Dialect) WriteQuerySelect(w core.Writer, q *dml.SelectQuery) {
	if q.ClauseSelect.IsDeclared() {
		q.ClauseSelect.Build()
		s.Self.WriteSelect(w, &q.ClauseSelect)
	}

	if q.ClauseFrom.IsDeclared() {
		q.ClauseFrom.Build()
		w.Char(' ')
		s.Self.WriteFrom(w, &q.ClauseFrom)
	}

	if q.ClauseJoin.IsDeclared() {
		q.ClauseJoin.Build()
		w.Char(' ')
		s.Self.WriteJoin(w, &q.ClauseJoin)
	}

	if q.ClauseWhere.IsDeclared() {
		q.ClauseWhere.Build()
		w.Char(' ')
		s.Self.WriteWhere(w, &q.ClauseWhere)
	}

	if q.ClauseGroupBy.IsDeclared() {
		q.ClauseGroupBy.Build()
		w.Char(' ')
		s.Self.WriteGroupBy(w, &q.ClauseGroupBy)
	}

	if q.ClauseHaving.IsDeclared() {
		q.ClauseHaving.Build()
		w.Char(' ')
		s.Self.WriteHaving(w, &q.ClauseHaving)
	}

	if q.ClauseOrderBy.IsDeclared() {
		q.ClauseOrderBy.Build()
		w.Char(' ')
		s.Self.WriteOrderBy(w, &q.ClauseOrderBy)
	}

	if q.ClauseLimit.IsDeclared() {
		q.ClauseLimit.Build()
		w.Char(' ')
		s.Self.WriteLimit(w, &q.ClauseLimit)
	}

	if q.ClauseOffset.IsDeclared() {
		q.ClauseOffset.Build()
		w.Char(' ')
		s.Self.WriteOffset(w, &q.ClauseOffset)
	}
}

func (s *Dialect) WriteQueryInsert(w core.Writer, q *dml.InsertQuery) {
	if q.ClauseInsert.IsDeclared() {
		q.ClauseInsert.Build()
		s.Self.WriteInsertInto(w, &q.ClauseInsert)
	}
}

func (s *Dialect) WriteQueryUpdate(w core.Writer, q *dml.UpdateQuery) {
	if q.ClauseUpdate.IsDeclared() {
		q.ClauseUpdate.Build()
		s.Self.WriteUpdate(w, &q.ClauseUpdate)
	}

	if q.ClauseWhere.IsDeclared() {
		q.ClauseWhere.Build()
		w.Char(' ')
		s.Self.WriteWhere(w, &q.ClauseWhere)
	}
}

func (s *Dialect) WriteQueryDelete(w core.Writer, q *dml.DeleteQuery) {
	if q.ClauseDelete.IsDeclared() {
		q.ClauseDelete.Build()
		s.Self.WriteDelete(w, &q.ClauseDelete)
	}

	if q.ClauseFrom.IsDeclared() {
		q.ClauseFrom.Build()
		w.Char(' ')
		s.Self.WriteFrom(w, &q.ClauseFrom)
	}

	if q.ClauseWhere.IsDeclared() {
		q.ClauseWhere.Build()
		w.Char(' ')
		s.Self.WriteWhere(w, &q.ClauseWhere)
	}
}

// func (*Dialect) SelectOrder() []dml.ClauseKind {
// 	return []dml.ClauseKind{
// 		dml.ClauseKindSelect,
// 		dml.ClauseKindFrom,
// 		dml.ClauseKindJoin,
// 		dml.ClauseKindWhere,
// 		dml.ClauseKindGroupBy,
// 		dml.ClauseKindHaving,
// 		dml.ClauseKindOrderBy,
// 		dml.ClauseKindLimit,
// 		dml.ClauseKindOffset,
// 	}
// }
//
// func (*Dialect) InsertOrder() []dml.ClauseKind {
// 	return []dml.ClauseKind{
// 		dml.ClauseKindInsert,
// 	}
// }
//
// func (*Dialect) UpdateOrder() []dml.ClauseKind {
// 	return []dml.ClauseKind{
// 		dml.ClauseKindUpdate,
// 		dml.ClauseKindWhere,
// 	}
// }
//
// func (*Dialect) DeleteOrder() []dml.ClauseKind {
// 	return []dml.ClauseKind{
// 		dml.ClauseKindDelete,
// 		dml.ClauseKindFrom,
// 		dml.ClauseKindWhere,
// 	}
// }
