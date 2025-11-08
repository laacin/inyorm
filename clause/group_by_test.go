package clause_test

import (
	"testing"

	"github.com/laacin/inyorm/clause"
	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

func NewGroupBy() (core.ClauseGroupBy, core.ColExpr, func(t *testing.T, cls string)) {
	stmt := writer.NewStatement("")
	stmt.SetFrom("default")
	var c core.ColExpr = &column.ColExpr{Statement: stmt}
	cls := &clause.GroupByClause{}

	run := func(t *testing.T, clause string) {
		w := stmt.Writer()
		cls.Build()(w)
		if val := w.ToString(); val != clause {
			t.Errorf("\nmismatch result:\nExpect:\n%s\nHave:\n%s\n", clause, val)
		}
	}

	return cls, c, run
}

func TestGroupBy(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cls, c, run := NewGroupBy()
		name := c.Col("name", "posts")
		cls.GroupBy(name).Having(name.Count(true)).Greater(10)

		run(t, "GROUP BY b.name HAVING (COUNT(DISTINCT b.name) > 10)")
	})

	t.Run("multiple_columns_with_function", func(t *testing.T) {
		cls, c, run := NewGroupBy()
		name := c.Col("name")
		role := c.Col("role", "roles")
		cls.GroupBy(name.Lower(), role).Having(role.Count()).Greater(5)

		run(t, "GROUP BY LOWER(a.name), b.role HAVING (COUNT(b.role) > 5)")
	})

	t.Run("lower_and_count", func(t *testing.T) {
		cls, c, run := NewGroupBy()
		name := c.Col("name")
		email := c.Col("email")
		cls.GroupBy(name.Lower()).Having(email.Count(true)).Greater(3)

		run(t, "GROUP BY LOWER(a.name) HAVING (COUNT(DISTINCT a.email) > 3)")
	})

	t.Run("concat_and_max", func(t *testing.T) {
		cls, c, run := NewGroupBy()
		fname := c.Col("firstname")
		lname := c.Col("lastname")
		score := c.Col("score")
		cls.GroupBy(c.Concat(fname, " ", lname)).Having(score.Max()).Less(80)

		run(t, "GROUP BY CONCAT(a.firstname, ' ', a.lastname) HAVING (MAX(a.score) < 80)")
	})

	t.Run("trim_and_min", func(t *testing.T) {
		cls, c, run := NewGroupBy()
		role := c.Col("role", "roles")
		points := c.Col("points")
		cls.GroupBy(role.Trim()).Having(points.Min()).Greater(10)

		run(t, "GROUP BY TRIM(b.role) HAVING (MIN(a.points) > 10)")
	})

	t.Run("round_and_avg", func(t *testing.T) {
		cls, c, run := NewGroupBy()
		price := c.Col("price", "products")
		discount := c.Col("discount", "products")
		cls.GroupBy(price.Round()).Having(discount.Avg()).Less(0.3)

		run(t, "GROUP BY ROUND(b.price) HAVING (AVG(b.discount) < 0.3)")
	})

	t.Run("nested_conditions", func(t *testing.T) {
		cls, c, run := NewGroupBy()
		role := c.Col("role", "roles")
		points := c.Col("points")
		points2 := c.Col("points")
		cls.GroupBy(role).Having(points.Sum()).Between(100, 1000).Or(points2.Count()).Greater(50)

		run(t, "GROUP BY b.role HAVING (SUM(a.points) BETWEEN 100 AND 1000 OR COUNT(a.points) > 50)")
	})

	t.Run("group_by_expression_and_alias", func(t *testing.T) {
		cls, c, run := NewGroupBy()
		score := c.Col("score", "stats")
		score2 := c.Col("score", "stats")
		level := c.Col("level", "stats")
		cls.GroupBy(c.Concat(level, "_", score.Round())).Having(score2.Max()).Less(90)

		run(t, "GROUP BY CONCAT(b.level, '_', ROUND(b.score)) HAVING (MAX(b.score) < 90)")
	})
}
