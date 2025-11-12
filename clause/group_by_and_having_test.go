package clause_test

import (
	"testing"

	"github.com/laacin/inyorm/clause"
	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

func NewGroupBy() (*clause.GroupBy, *clause.Having, core.ColExpr, func(t *testing.T, cls string)) {
	stmt := writer.NewStatement("", "users")
	var c core.ColExpr = &column.ColExpr{Writer: stmt.Writer}
	gb := &clause.GroupBy{}
	h := &clause.Having{}

	run := func(t *testing.T, clause string) {
		w := stmt.Writer()
		gb.Build(w)
		w.Char(' ')
		h.Build(w)
		if val := w.ToString(); val != clause {
			t.Errorf("\nmismatch result:\nExpect:\n%s\nHave:\n%s\n", clause, val)
		}
	}

	return gb, h, c, run
}

func TestGroupBy(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		gb, h, c, run := NewGroupBy()

		name := c.Col("name", "posts")
		name.Count(true)

		gb.GroupBy([]any{name.Base()})
		h.Having(name).Greater(10)

		run(t, "GROUP BY b.name HAVING (COUNT(DISTINCT b.name) > 10)")
	})

	t.Run("multiple_columns_with_function", func(t *testing.T) {
		gb, h, c, run := NewGroupBy()

		name := c.Col("name", "users")
		name.Lower()

		role := c.Col("role", "roles")
		role.Count(false)

		gb.GroupBy([]any{name, role.Base()})
		h.Having(role).Greater(5)

		run(t, "GROUP BY LOWER(a.name), b.role HAVING (COUNT(b.role) > 5)")
	})

	t.Run("lower_and_count", func(t *testing.T) {
		gb, h, c, run := NewGroupBy()

		name := c.Col("name", "users")
		name.Lower()

		email := c.Col("email", "users")
		email.Count(true)

		gb.GroupBy([]any{name})
		h.Having(email).Greater(3)

		run(t, "GROUP BY LOWER(a.name) HAVING (COUNT(DISTINCT a.email) > 3)")
	})

	t.Run("concat_and_max", func(t *testing.T) {
		gb, h, c, run := NewGroupBy()

		fname := c.Col("firstname", "users")
		lname := c.Col("lastname", "users")
		score := c.Col("score", "users")
		score.Max(false)

		con := c.Concat([]any{fname, " ", lname})
		gb.GroupBy([]any{con})
		h.Having(score).Less(80)

		run(t, "GROUP BY CONCAT(a.firstname, ' ', a.lastname) HAVING (MAX(a.score) < 80)")
	})

	t.Run("trim_and_min", func(t *testing.T) {
		gb, h, c, run := NewGroupBy()

		role := c.Col("role", "roles")
		role.Trim()

		points := c.Col("points", "users")
		points.Min(false)

		gb.GroupBy([]any{role})
		h.Having(points).Greater(10)

		run(t, "GROUP BY TRIM(b.role) HAVING (MIN(a.points) > 10)")
	})

	t.Run("round_and_avg", func(t *testing.T) {
		gb, h, c, run := NewGroupBy()

		price := c.Col("price", "products")
		price.Round()

		discount := c.Col("discount", "products")
		discount.Avg(false)

		gb.GroupBy([]any{price})
		h.Having(discount).Less(0.3)

		run(t, "GROUP BY ROUND(b.price) HAVING (AVG(b.discount) < 0.3)")
	})

	t.Run("nested_conditions", func(t *testing.T) {
		gb, h, c, run := NewGroupBy()

		role := c.Col("role", "roles")
		points := c.Col("points", "users")
		points.Sum(false)

		points2 := c.Col("points", "users")
		points2.Count(false)

		gb.GroupBy([]any{role})
		cond := h.Having(points)
		cond.Between(100, 1000)
		cond.Or(points2)
		cond.Greater(50)

		run(t, "GROUP BY b.role HAVING (SUM(a.points) BETWEEN 100 AND 1000 OR COUNT(a.points) > 50)")
	})

	t.Run("group_by_expression_and_alias", func(t *testing.T) {
		gb, h, c, run := NewGroupBy()

		score := c.Col("score", "stats")
		score.Round()

		score2 := c.Col("score", "stats")
		score2.Max(false)

		level := c.Col("level", "stats")

		con := c.Concat([]any{level, "_", score})
		gb.GroupBy([]any{con})
		h.Having(score2).Less(90)

		run(t, "GROUP BY CONCAT(b.level, '_', ROUND(b.score)) HAVING (MAX(b.score) < 90)")
	})
}
