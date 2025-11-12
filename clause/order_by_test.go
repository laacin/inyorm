package clause_test

import (
	"testing"

	"github.com/laacin/inyorm/clause"
	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

func NewOrderBy() (*clause.OrderBy, core.ColExpr, func(t *testing.T, cls string)) {
	stmt := writer.NewStatement("", "users")
	var c core.ColExpr = &column.ColExpr{Writer: stmt.Writer}
	cls := &clause.OrderBy{}

	run := func(t *testing.T, clause string) {
		w := stmt.Writer()
		cls.Build(w)
		if val := w.ToString(); val != clause {
			t.Errorf("\nmismatch result:\nExpect:\n%s\nHave:\n%s\n", clause, val)
		}
	}

	return cls, c, run
}

func TestOrderBy(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cls, c, run := NewOrderBy()

		cls.OrderBy(c.Col("firstname", "users"))
		cls.OrderBy(c.Col("age", "users"))

		run(t, "ORDER BY a.firstname, a.age")
	})

	t.Run("descending", func(t *testing.T) {
		cls, c, run := NewOrderBy()

		cls.OrderBy(c.Col("age", "users"))
		cls.Desc()

		cls.OrderBy(c.Col("lastname", "users"))
		cls.Desc()

		run(t, "ORDER BY a.age DESC, a.lastname DESC")
	})

	t.Run("mix", func(t *testing.T) {
		cls, c, run := NewOrderBy()

		var (
			postNum = c.Col("id", "posts")
			age     = c.Col("age", "users")
		)

		postNum.Count(false)

		cls.OrderBy(postNum)
		cls.Desc()

		cls.OrderBy(age)

		run(t, "ORDER BY COUNT(b.id) DESC, a.age")
	})
}
