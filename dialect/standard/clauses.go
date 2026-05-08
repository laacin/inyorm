package standard

import "github.com/laacin/inyorm/internal/entity"

func (dial *DialectStandard) WriteInsertInto(w entity.Writer, cls *entity.InsertInto) {
	w.Write("INSERT INTO")
	w.Char(' ')

	w.Value(cls.Table, entity.WriteDef)
	w.Char(' ')

	w.Char('(')
	for i, col := range cls.Cols {
		if i > 0 {
			w.Write(", ")
		}
		w.Write(col)
	}
	w.Char(')')

	perRow := len(cls.Values) / cls.Rows

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
			w.Value(cls.Values[row*perRow+ci], entity.WriteDef)
		}
		w.Char(')')
	}
}

func (dial *DialectStandard) WriteSelect(w entity.Writer, cls *entity.Select) {
	w.Write("SELECT")
	w.Char(' ')

	if cls.Distinct {
		w.Write("DISTINCT")
		w.Char(' ')
	}

	for i, val := range cls.Values {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(val, entity.WriteDef)
	}
}

func (dial *DialectStandard) WriteFrom(w entity.Writer, cls *entity.From) {
	w.Write("FROM")
	w.Char(' ')
	w.Value(cls.Value, entity.WriteDef)
}

var joinTypeMap = map[entity.JoinType]string{
	entity.JoinInner: "INNER",
	entity.JoinLeft:  "LEFT",
	entity.JoinRight: "RIGHT",
	entity.JoinFull:  "FULL",
	entity.JoinCross: "CROSS",
}

func (dial *DialectStandard) WriteJoin(w entity.Writer, cls *entity.Join) {
	for i, join := range cls.Joins {
		if i > 0 {
			w.Char(' ')
		}

		w.Write(joinTypeMap[join.Type])
		w.Write(" JOIN ")
		w.Value(join.Table, entity.WriteDef)

		if join.Cond != nil {
			w.Write(" ON ")
			w.Value(join.Cond, entity.WriteBase)
		}
	}
}

func (dial *DialectStandard) WriteWhere(w entity.Writer, cls *entity.Where) {
	w.Write("WHERE")
	w.Char(' ')

	for i, cond := range cls.Conds {
		if i > 0 {
			w.Write(" AND ")
		}
		w.Value(cond, entity.WriteExpr)
	}
}

func (dial *DialectStandard) WriteGroupBy(w entity.Writer, cls *entity.GroupBy) {
	w.Write("GROUP BY")
	w.Char(' ')

	for i, group := range cls.Values {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(group, entity.WriteExpr)
	}
}

func (dial *DialectStandard) WriteHaving(w entity.Writer, cls *entity.Having) {
	w.Write("HAVING")
	w.Char(' ')
	w.Value(cls.Cond, entity.WriteExpr)
}

func (dial *DialectStandard) WriteOrderBy(w entity.Writer, cls *entity.OrderBy) {
	w.Write("ORDER BY")
	w.Char(' ')

	for i, ord := range cls.Orders {
		if i > 0 {
			w.Write(", ")
		}

		w.Value(ord.Value, entity.WriteAlias)
		if ord.Descending {
			w.Char(' ')
			w.Write("DESC")
		}
	}
}

func (dial *DialectStandard) WriteLimit(w entity.Writer, cls *entity.Limit) {
	w.Write("LIMIT")
	w.Char(' ')
	w.Value(cls.ValueNumber, entity.WriteBase)
}

func (dial *DialectStandard) WriteOffset(w entity.Writer, cls *entity.Offset) {
	w.Write("OFFSET")
	w.Char(' ')
	w.Value(cls.ValueNumber, entity.WriteBase)
}

func (dial *DialectStandard) WriteUpdate(w entity.Writer, cls *entity.Update) {
	w.Write("UPDATE")
	w.Char(' ')

	w.Value(cls.Table, entity.WriteDef)
	w.Write(" SET ")

	for i, col := range cls.Cols {
		if i > 0 {
			w.Write(", ")
		}

		w.Write(col)
		w.Write(" = ")
		w.Value(cls.Values[i], entity.WriteDef)
	}
}

func (dial *DialectStandard) WriteDelete(w entity.Writer, cls *entity.Delete) {
	w.Write("DELETE")
}
