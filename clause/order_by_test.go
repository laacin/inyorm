package clause_test

import "testing"

func TestOrderBy(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cls, c, run := NewOrderBy()

		cls.OrderBy(c.Col("firstname", "users"))
		cls.OrderBy(c.Col("age", "users"))

		run(t, "ORDER BY firstname, age", nil)
	})

	t.Run("descending", func(t *testing.T) {
		cls, c, run := NewOrderBy()

		cls.OrderBy(c.Col("age", "users")).Desc()
		cls.OrderBy(c.Col("lastname", "users")).Desc()

		run(t, "ORDER BY age DESC, lastname DESC", nil)
	})

	t.Run("mix", func(t *testing.T) {
		cls, c, run := NewOrderBy()

		var (
			postNum = c.Col("id", "posts")
			age     = c.Col("age", "users")
		)

		postNum.Count(false)

		cls.OrderBy(postNum).Desc()
		cls.OrderBy(age)

		run(t, "ORDER BY COUNT(id) DESC, age", nil)
	})
}
