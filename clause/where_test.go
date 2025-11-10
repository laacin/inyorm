package clause_test

import (
	"github.com/laacin/inyorm/clause"
	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
	"reflect"
	"testing"
)

func NewWhere(dialect ...string) (core.ClauseWhere, core.ColExpr, func(t *testing.T, cls string, vals []any)) {
	d := ""
	if len(dialect) > 0 {
		d = dialect[0]
	}

	stmt := writer.NewStatement(d, "default")
	stmt.SetFrom("default")
	var c core.ColExpr = &column.ColExpr{Statement: stmt}
	cls := &clause.WhereClause{}

	run := func(t *testing.T, clause string, vals []any) {
		w := stmt.Writer()
		cls.Build(w)
		if val := w.ToString(); val != clause {
			t.Errorf("\nmismatch result:\nExpect:\n%s\nHave:\n%s\n", clause, val)
		}

		values := stmt.Values()
		if !reflect.DeepEqual(vals, values) {
			t.Errorf("\nmissmatch values:\nExpect:\n%#v\nHave:\n%#v\n", vals, values)
		}
	}

	return cls, c, run
}

func TestWhere(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		cls, c, run := NewWhere()
		cls.Where(c.Col("firstname")).Not().Not().Equal("john").Or().Not().Equal("mary")

		exp := "WHERE (a.firstname = ? OR a.firstname <> ?)"
		run(t, exp, []any{"john", "mary"})
	})

	t.Run("basic_2", func(t *testing.T) {
		cls, c, run := NewWhere()
		cls.Where(c.Col("lastname")).Like("%alv%").And().Not().In("calvin", "malvina", "salvatore")

		exp := "WHERE (a.lastname LIKE ? AND a.lastname NOT IN (?, ?, ?))"
		run(t, exp, []any{"%alv%", "calvin", "malvina", "salvatore"})
	})

	t.Run("basic_3", func(t *testing.T) {
		cls, c, run := NewWhere()
		cls.Where(c.Col("age")).Between(17, 70).And().Not().Equal(45)

		exp := "WHERE (a.age BETWEEN ? AND ? AND a.age <> ?)"
		run(t, exp, []any{17, 70, 45})
	})

	t.Run("many_wheres", func(t *testing.T) {
		cls, c, run := NewWhere()
		cls.Where(c.Col("age")).Between(17, 70).And().Not().Equal(45)
		cls.Where(c.Col("firstname")).Like("%alv%").And().Not().In("calvin", "malvina", "salvatore")
		cls.Where(c.Col("lastname")).Not().Not().Equal("john").Or().Not().Equal("mary")
		cls.Where("literal").Not().IsNull()

		exp := "WHERE (a.age BETWEEN ? AND ? AND a.age <> ?)"
		exp += " AND "
		exp += "(a.firstname LIKE ? AND a.firstname NOT IN (?, ?, ?))"
		exp += " AND "
		exp += "(a.lastname = ? OR a.lastname <> ?)"
		exp += " AND "
		exp += "('literal' IS NOT NULL)"

		run(t, exp, []any{17, 70, 45, "%alv%", "calvin", "malvina", "salvatore", "john", "mary"})
	})

	t.Run("many_wheres_with_postgres_placeholder", func(t *testing.T) {
		cls, c, run := NewWhere(core.Postgres)
		cls.Where(c.Col("age")).Between(17, 70).And().Not().Equal(45)
		cls.Where(c.Col("firstname")).Like("%alv%").And().Not().In("calvin", "malvina", "salvatore")
		cls.Where(c.Col("lastname")).Not().Not().Equal("john").Or().Not().Equal("mary")
		cls.Where("literal").Not().IsNull()

		exp := "WHERE (a.age BETWEEN $1 AND $2 AND a.age <> $3)"
		exp += " AND "
		exp += "(a.firstname LIKE $4 AND a.firstname NOT IN ($5, $6, $7))"
		exp += " AND "
		exp += "(a.lastname = $8 OR a.lastname <> $9)"
		exp += " AND "
		exp += "('literal' IS NOT NULL)"

		run(t, exp, []any{17, 70, 45, "%alv%", "calvin", "malvina", "salvatore", "john", "mary"})
	})
}
