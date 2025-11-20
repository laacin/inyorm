package clause_test

import "testing"

func TestGroupBy(t *testing.T) {
	t.Run("single_column", func(t *testing.T) {
		cls, c, run := NewGroupBy()

		name := c.Col("name", "posts")
		cls.GroupBy(name.Base())

		run(t, "GROUP BY name", nil)
	})

	t.Run("multiple_columns", func(t *testing.T) {
		cls, c, run := NewGroupBy()

		name := c.Col("name", "users")
		role := c.Col("role", "roles")

		cls.GroupBy(name.Base(), role.Base())

		run(t, "GROUP BY name, role", nil)
	})

	t.Run("expression_column", func(t *testing.T) {
		cls, c, run := NewGroupBy()

		fname := c.Col("firstname", "users")
		lname := c.Col("lastname", "users")
		con := c.Concat(fname, " ", lname)

		cls.GroupBy(con)

		run(t, "GROUP BY CONCAT(firstname, ' ', lastname)", nil)
	})

	t.Run("column_with_function", func(t *testing.T) {
		cls, c, run := NewGroupBy()

		price := c.Col("price", "products")
		price.Round()

		cls.GroupBy(price)

		run(t, "GROUP BY ROUND(price)", nil)
	})
}
