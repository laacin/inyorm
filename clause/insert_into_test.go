package clause_test

import (
	"reflect"
	"testing"

	"github.com/laacin/inyorm/clause"
	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

func NewInsert(dialect ...string) (*clause.InsertInto, core.ColExpr, func(t *testing.T, cls string, vals []any)) {
	d := ""
	if len(dialect) > 0 {
		d = dialect[0]
	}

	stmt := writer.NewStatement(d, "users")
	var c core.ColExpr = &column.ColExpr{}
	cls := &clause.InsertInto{}

	run := func(t *testing.T, clause string, vals []any) {
		stmt.SetClauses([]core.Clause{cls}, writer.InsertOrder)
		statement, values := stmt.Build()

		if statement != clause {
			t.Errorf("\nmismatch result:\nExpect:\n%s\nHave:\n%s\n", clause, statement)
		}

		if !reflect.DeepEqual(vals, values) {
			t.Errorf("\nmissmatch values:\nExpect:\n%#v\nHave:\n%#v\n", vals, values)
		}
	}

	return cls, c, run
}

func TestInserInto(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cls, c, run := NewInsert("")

		var (
			fname = c.Col("firstname", "users")
			lname = c.Col("lastname", "users")
		)

		cls.Insert("users", []any{fname, lname})
		vals := []any{
			"John", "Doe",
			"Jane", "Smith",
			"Mike", "Brown",
		}
		cls.Values(vals)

		exp := "INSERT INTO users (firstname, lastname) VALUES (?, ?), (?, ?), (?, ?)"

		run(t, exp, vals)
	})
}
