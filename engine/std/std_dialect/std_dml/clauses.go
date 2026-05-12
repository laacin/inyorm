package std_dml

import (
	"github.com/laacin/inyorm/internal/entity/core"
	"github.com/laacin/inyorm/internal/entity/dml"
)

func (*DmlSyntax) WriteInsertInto(w core.Writer, cls *dml.InsertInto) {
	w.Write("INSERT INTO")
	w.Char(' ')

	w.Value(cls.Table, core.WriteDef)
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
			w.Value(cls.Values[row*perRow+ci], core.WriteDef)
		}
		w.Char(')')
	}
}

func (*DmlSyntax) WriteSelect(w core.Writer, cls *dml.Select) {
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
		w.Value(val, core.WriteDef)
	}
}

func (*DmlSyntax) WriteFrom(w core.Writer, cls *dml.From) {
	w.Write("FROM")
	w.Char(' ')
	w.Value(cls.Value, core.WriteDef)
}

var joinTypeMap = map[dml.JoinType]string{
	dml.JoinInner: "INNER",
	dml.JoinLeft:  "LEFT",
	dml.JoinRight: "RIGHT",
	dml.JoinFull:  "FULL",
	dml.JoinCross: "CROSS",
}

func (*DmlSyntax) WriteJoin(w core.Writer, cls *dml.Join) {
	for i, join := range cls.Joins {
		if i > 0 {
			w.Char(' ')
		}

		w.Write(joinTypeMap[join.Type])
		w.Write(" JOIN ")
		w.Value(join.Table, core.WriteDef)

		if join.Cond != nil {
			w.Write(" ON ")
			w.Value(join.Cond, core.WriteBase)
		}
	}
}

func (*DmlSyntax) WriteWhere(w core.Writer, cls *dml.Where) {
	w.Write("WHERE")
	w.Char(' ')

	for i, cond := range cls.Conds {
		if i > 0 {
			w.Write(" AND ")
		}
		w.Value(cond, core.WriteExpr)
	}
}

func (*DmlSyntax) WriteGroupBy(w core.Writer, cls *dml.GroupBy) {
	w.Write("GROUP BY")
	w.Char(' ')

	for i, group := range cls.Values {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(group, core.WriteExpr)
	}
}

func (*DmlSyntax) WriteHaving(w core.Writer, cls *dml.Having) {
	w.Write("HAVING")
	w.Char(' ')
	w.Value(cls.Cond, core.WriteExpr)
}

func (*DmlSyntax) WriteOrderBy(w core.Writer, cls *dml.OrderBy) {
	w.Write("ORDER BY")
	w.Char(' ')

	for i, ord := range cls.Orders {
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

func (*DmlSyntax) WriteLimit(w core.Writer, cls *dml.Limit) {
	w.Write("LIMIT")
	w.Char(' ')
	w.Value(cls.ValueNumber, core.WriteBase)
}

func (*DmlSyntax) WriteOffset(w core.Writer, cls *dml.Offset) {
	w.Write("OFFSET")
	w.Char(' ')
	w.Value(cls.ValueNumber, core.WriteBase)
}

func (*DmlSyntax) WriteUpdate(w core.Writer, cls *dml.Update) {
	w.Write("UPDATE")
	w.Char(' ')

	w.Value(cls.Table, core.WriteDef)
	w.Write(" SET ")

	for i, col := range cls.Cols {
		if i > 0 {
			w.Write(", ")
		}

		w.Write(col)
		w.Write(" = ")
		w.Value(cls.Values[i], core.WriteDef)
	}
}

func (*DmlSyntax) WriteDelete(w core.Writer, cls *dml.Delete) {
	w.Write("DELETE")
}
