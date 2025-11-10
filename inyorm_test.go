package inyorm_test

import (
	"reflect"
	"testing"

	"github.com/laacin/inyorm"
)

func run(t *testing.T, q *inyorm.Query, exp string, vals []any) {
	stmt, values := q.Build()
	if stmt != exp {
		t.Errorf("mismatch statement:\nExpect:\n%s\nHave:\n%s", exp, stmt)
	}

	if reflect.DeepEqual(values, vals) {
		t.Errorf("mismatch values:\nExpect:\n%#v\nHave:\n%#v", values, vals)
	}
}

func TestInyorm(t *testing.T) {
	qe := &inyorm.QueryEngine{Dialect: "mysql"}

	t.Run("simple", func(t *testing.T) {
		q, c := qe.New("users")

		q.Select(c.All())
		q.Where(c.Col("id")).Equal("uuid")
		q.Limit(1)
	})

	t.Run("pagination", func(t *testing.T) {
		q, c := qe.New("users")

		var (
			id      = c.Col("id")
			age     = c.Col("age")
			banned  = c.Col("banned")
			foreign = c.Col("user_id", "posts")
		)

		q.Select(c.All())
		q.Join("posts").On(foreign).Equal(id)
		q.Where(banned).IsNull().And(age).Greater(17)
		q.OrderBy(age).Desc()
		q.Limit(100).Offset(20)

		exp := "SELECT * "
		exp += "FROM users a "
		exp += "INNER JOIN posts b ON (b.user_id = a.id) "
		exp += "WHERE (a.banned IS NULL AND a.age > ?) "
		exp += "ORDER BY a.age DESC "
		exp += "LIMIT 100 OFFSET 20"

		run(t, q, exp, []any{17})
	})
}
