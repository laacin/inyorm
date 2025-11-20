package clause_test

import "testing"

func TestSelect(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cls, c, run := NewSelect()
		cls.Select("active", c.Col("name", "users"))

		run(t, "SELECT 'active', name", nil)
	})

	t.Run("distinct", func(t *testing.T) {
		cls, c, run := NewSelect()
		cls.Distinct().Select(c.Col("age", "users"))

		run(t, "SELECT DISTINCT age", nil)
	})
}
