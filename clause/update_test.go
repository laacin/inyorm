package clause_test

import "testing"

func TestUpdate(t *testing.T) {
	type User struct {
		Account string `col:"account"`
		Age     int    `col:"age"`
		Active  bool   `col:"active"`
		Score   int    `col:"score"`
		Country string `col:"country"`
	}

	t.Run("full", func(t *testing.T) {
		cls, _, run := NewUpdate("users")

		cls.Update(User{}).Values(User{
			Account: "acc1",
			Age:     28,
			Active:  true,
			Score:   250,
			Country: "US",
		})

		exp := "UPDATE users SET account = ?, active = ?, age = ?, country = ?, score = ?"
		vals := []any{"acc1", true, 28, "US", 250}
		run(t, exp, vals)
	})

	t.Run("partial_with_omission", func(t *testing.T) {
		cls, c, run := NewUpdate("users")

		cls.Update(c.Col("country")).Values(User{
			Country: "US",
		})

		exp := "UPDATE users SET country = ?"
		vals := []any{"US"}
		run(t, exp, vals)
	})

	t.Run("partial_with_literal", func(t *testing.T) {
		cls, c, run := NewUpdate("users")

		cls.Update(c.Col("country")).Values("US")

		exp := "UPDATE users SET country = ?"
		vals := []any{"US"}
		run(t, exp, vals)
	})

	t.Run("full_with_map", func(t *testing.T) {
		cls, _, run := NewUpdate("users")

		m := map[string]any{
			"account": "acc1",
			"age":     28,
			"active":  true,
			"score":   250,
			"country": "US",
		}

		cls.Update(m).Values(m)

		exp := "UPDATE users SET account = ?, active = ?, age = ?, country = ?, score = ?"
		vals := []any{"acc1", true, 28, "US", 250}
		run(t, exp, vals)
	})

	t.Run("partial_with_map", func(t *testing.T) {
		cls, c, run := NewUpdate("users")

		m := map[string]any{
			"country": "US",
		}

		cls.Update(c.Col("country")).Values(m)

		exp := "UPDATE users SET country = ?"
		vals := []any{"US"}
		run(t, exp, vals)
	})
}
