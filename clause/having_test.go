package clause_test

import "testing"

func TestHavingOnly(t *testing.T) {
	t.Run("simple_having", func(t *testing.T) {
		cls, c, run := NewHaving()

		name := c.Col("name", "posts").Count(true)

		cls.Having(name).Greater(10)

		run(t, "HAVING (COUNT(DISTINCT name) > 10)", nil)
	})

	t.Run("multiple_columns_with_function", func(t *testing.T) {
		cls, c, run := NewHaving()

		role := c.Col("role", "roles")
		role.Count()

		cls.Having(role).Greater(5)

		run(t, "HAVING (COUNT(role) > 5)", nil)
	})

	t.Run("lower_and_count", func(t *testing.T) {
		cls, c, run := NewHaving()

		email := c.Col("email", "users").Count(true)

		cls.Having(email).Greater(3)

		run(t, "HAVING (COUNT(DISTINCT email) > 3)", nil)
	})

	t.Run("concat_and_max", func(t *testing.T) {
		cls, c, run := NewHaving()

		score := c.Col("score", "users").Max()

		cls.Having(score).Less(80)

		run(t, "HAVING (MAX(score) < 80)", nil)
	})

	t.Run("trim_and_min", func(t *testing.T) {
		cls, c, run := NewHaving()

		points := c.Col("points", "users").Min()

		cls.Having(points).Greater(10)

		run(t, "HAVING (MIN(points) > 10)", nil)
	})

	t.Run("round_and_avg", func(t *testing.T) {
		cls, c, run := NewHaving()

		discount := c.Col("discount", "products").Avg()

		cls.Having(discount).Less(0.3)

		run(t, "HAVING (AVG(discount) < 0.3)", nil)
	})

	t.Run("nested_conditions", func(t *testing.T) {
		cls, c, run := NewHaving()

		points := c.Col("points", "users").Sum()
		points2 := c.Col("points", "users").Count()

		cls.Having(points).Between(100, 1000).Or(points2).Greater(50)

		run(t, "HAVING (SUM(points) BETWEEN 100 AND 1000 OR COUNT(points) > 50)", nil)
	})

	t.Run("group_by_expression_and_alias", func(t *testing.T) {
		cls, c, run := NewHaving()

		score2 := c.Col("score", "stats").Max()
		cls.Having(score2).Less(90)

		run(t, "HAVING (MAX(score) < 90)", nil)
	})
}
