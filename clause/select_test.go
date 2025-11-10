package clause_test

import (
	"testing"

	"github.com/laacin/inyorm/clause"
	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

func NewSelect() (core.ClauseSelect, core.ColExpr, func(t *testing.T, cls string)) {
	stmt := writer.NewStatement("", "default")
	stmt.SetFrom("default")
	var c core.ColExpr = &column.ColExpr{Statement: stmt}
	cls := &clause.SelectClause{}

	run := func(t *testing.T, clause string) {
		w := stmt.Writer()
		cls.Build(w)
		if val := w.ToString(); val != clause {
			t.Errorf("\nmismatch result:\nExpect:\n%s\nHave:\n%s\n", clause, val)
		}
	}

	return cls, c, run
}

func TestSelect(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cls, c, run := NewSelect()
		cls.Select("active", c.Col("name"))

		run(t, "SELECT 'active', a.name")
	})

	t.Run("distinct", func(t *testing.T) {
		cls, c, run := NewSelect()
		cls.Distinct().Select(c.Col("age"))

		run(t, "SELECT DISTINCT a.age")
	})
}
