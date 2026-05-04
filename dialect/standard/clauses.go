package standard

import "github.com/laacin/inyorm/intr/dialect"

func (dial *DialectStandard) ClsInsertInto(w dialect.Writer, tools dialect.InsertIntoTools) {
	w.Write("INSERT INTO")
	w.Char(' ')

	w.Write(tools.Table)
	w.Char(' ')

	w.Char('(')
	for i, col := range tools.Columns {
		if i > 0 {
			w.Write(", ")
		}
		w.Write(col)
	}
	w.Char(')')

	w.Write(" VALUES ")
	for row := range tools.Rows {
		if row > 0 {
			w.Write(", ")
		}

		w.Char('(')
		for ci := range tools.Columns {
			if ci > 0 {
				w.Write(", ")
			}
			w.Char('?')
		}
		w.Char(')')
	}
}

func (dial *DialectStandard) ClsSelect(w dialect.Writer, tools dialect.SelectTools) {
	w.Write("SELECT")
	w.Char(' ')

	if tools.Distinct {
		w.Write("DISTINCT")
		w.Char(' ')
	}

	for i, val := range tools.Values {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(val, dialect.WriteDef)
	}
}

func (dial *DialectStandard) ClsFrom(w dialect.Writer, tools dialect.FromTools) {
	w.Write("FROM")
	w.Char(' ')

	switch tbl := tools.Value.(type) {
	case string:
		w.Write(tbl)
	default:
		panic("") // TODO: support sub-queries
	}
}

var joinTypeMap = map[dialect.JoinType]string{
	dialect.JoinInner: "INNER",
	dialect.JoinLeft:  "LEFT",
	dialect.JoinRight: "RIGHT",
	dialect.JoinFull:  "FULL",
	dialect.JoinCross: "CROSS",
}

func (dial *DialectStandard) ClsJoin(w dialect.Writer, tools []dialect.JoinTools) {
	for i, join := range tools {
		if i > 0 {
			w.Char(' ')
		}

		w.Write(joinTypeMap[join.Type]) // NOTE: could be fragile
		w.Write(" JOIN ")
		dial.Table(w, dialect.Table{Name: join.Table}, true)

		if join.Cond != nil {
			w.Write(" ON ")
			dial.Cond(w, *join.Cond, dialect.WriteBase)
		}
	}
}

func (dial *DialectStandard) ClsWhere(w dialect.Writer, tools dialect.WhereTools) {
	w.Write("WHERE")
	w.Char(' ')

	for i, cond := range tools.Conds {
		if i > 0 {
			w.Write(" AND ")
		}
		dial.Cond(w, cond, dialect.WriteExpr)
	}
}

func (dial *DialectStandard) ClsGroupBy(w dialect.Writer, tools dialect.GroupByTools) {
	w.Write("GROUP BY")
	w.Char(' ')

	for i, group := range tools.Values {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(group, dialect.WriteExpr)
	}
}

func (dial *DialectStandard) ClsHaving(w dialect.Writer, tools dialect.HavingTools) {
	w.Write("HAVING")
	w.Char(' ')
	dial.Cond(w, tools.Cond, dialect.WriteExpr)
}

func (dial *DialectStandard) ClsOrderBy(w dialect.Writer, tools []dialect.OrderByTools) {
	w.Write("ORDER BY")
	w.Char(' ')

	for i, ord := range tools {
		if i > 0 {
			w.Write(", ")
		}

		w.Value(ord.Value, dialect.WriteAlias)
		if ord.Descending {
			w.Char(' ')
			w.Write("DESC")
		}
	}
}

func (dial *DialectStandard) ClsLimit(w dialect.Writer, tools dialect.LimitTools) {
	w.Write("LIMIT")
	w.Char(' ')
	w.Value(tools.ValueNumber, dialect.WriteBase)
}

func (dial *DialectStandard) ClsOffset(w dialect.Writer, tools dialect.OffsetTools) {
	w.Write("OFFSET")
	w.Char(' ')
	w.Value(tools.ValueNumber, dialect.WriteBase)
}

func (dial *DialectStandard) ClsUpdate(w dialect.Writer, tools dialect.UpdateTools) {
	w.Write("UPDATE")
	w.Char(' ')

	w.Write(tools.Table)
	w.Write(" SET ")

	for i, col := range tools.Columns {
		if i > 0 {
			w.Write(", ")
		}

		w.Write(col)
		w.Write(" = ")
		w.Char('?')
	}
}

func (dial *DialectStandard) ClsDelete(w dialect.Writer, tools dialect.DeleteTools) {
	w.Write("DELETE")
}
