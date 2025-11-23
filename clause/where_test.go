package clause_test

import (
	"testing"

	"github.com/laacin/inyorm/internal/core"
)

func TestWhere(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		cls, c, run := NewWhere()
		name := c.Col("firstname", "users")

		cond := cls.Where(name).Not().Not().Equal(c.Param("john"))
		cond.Or(name).Not().Equal(c.Param("mary"))

		exp := "WHERE (firstname = ? OR firstname <> ?)"
		run(t, exp, []any{"john", "mary"})
	})

	t.Run("basic_2", func(t *testing.T) {
		cls, c, run := NewWhere()
		lname := c.Col("lastname", "users")

		in := []any{c.Param("calvin"), c.Param("malvina"), c.Param("salvatore")}
		cls.Where(lname).Like(c.Param("%alv%")).And(lname).Not().In(in)

		exp := "WHERE (lastname LIKE ? AND lastname NOT IN (?, ?, ?))"
		run(t, exp, []any{"%alv%", "calvin", "malvina", "salvatore"})
	})

	t.Run("basic_3", func(t *testing.T) {
		cls, c, run := NewWhere()
		age := c.Col("age", "users")
		cond := cls.Where(age).Between(c.Param(17), c.Param(70))
		cond.And(age).Not().Equal(c.Param(45))

		exp := "WHERE (age BETWEEN ? AND ? AND age <> ?)"
		run(t, exp, []any{17, 70, 45})
	})

	t.Run("many_wheres", func(t *testing.T) {
		cls, c, run := NewWhere()
		var (
			age   = c.Col("age", "users")
			fname = c.Col("firstname", "users")
			lname = c.Col("lastname", "users")
		)

		in := []any{c.Param("calvin"), c.Param("malvina"), c.Param("salvatore")}

		cls.Where(age).Between(c.Param(17), c.Param(70)).And(age).Not().Equal(c.Param(45))
		cls.Where(fname).Like(c.Param("%alv%")).And(fname).Not().In(in)
		cls.Where(lname).Equal(c.Param("john")).Or(lname).Not().Equal(c.Param("mary"))
		cls.Where("literal").Not().IsNull()

		exp := "WHERE (age BETWEEN ? AND ? AND age <> ?)"
		exp += " AND "
		exp += "(firstname LIKE ? AND firstname NOT IN (?, ?, ?))"
		exp += " AND "
		exp += "(lastname = ? OR lastname <> ?)"
		exp += " AND "
		exp += "('literal' IS NOT NULL)"

		run(t, exp, []any{17, 70, 45, "%alv%", "calvin", "malvina", "salvatore", "john", "mary"})
	})

	t.Run("many_wheres_with_postgres_placeholder", func(t *testing.T) {
		cls, c, run := NewWhere(core.Postgres)

		var (
			age   = c.Col("age", "users")
			fname = c.Col("firstname", "users")
			lname = c.Col("lastname", "users")
		)

		in := []any{c.Param("calvin"), c.Param("malvina"), c.Param("salvatore")}

		cls.Where(age).Between(c.Param(17), c.Param(70)).And(age).Not().Equal(c.Param(45))
		cls.Where(fname).Like(c.Param("%alv%")).And(fname).Not().In(in)
		cls.Where(lname).Not().Not().Equal(c.Param("john")).Or(lname).Not().Equal(c.Param("mary"))
		cls.Where("literal").Not().IsNull()

		exp := "WHERE (age BETWEEN $1 AND $2 AND age <> $3)"
		exp += " AND "
		exp += "(firstname LIKE $4 AND firstname NOT IN ($5, $6, $7))"
		exp += " AND "
		exp += "(lastname = $8 OR lastname <> $9)"
		exp += " AND "
		exp += "('literal' IS NOT NULL)"

		run(t, exp, []any{17, 70, 45, "%alv%", "calvin", "malvina", "salvatore", "john", "mary"})
	})
}
