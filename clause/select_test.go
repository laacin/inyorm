package clause_test

import (
	"strings"
	"testing"

	"github.com/laacin/inyorm/clause"
)

func TestSelect(t *testing.T) {
	expect := func(t *testing.T, exp, have string) {
		if exp != have {
			t.Errorf("mismatch clause:\nExpect: %s\nHave: %s\n", exp, have)
		}
	}

	t.Run("basic_select", func(t *testing.T) {
		sel := clause.SelectBuilder{}
		var sb strings.Builder

		sel.New("firstname")
		sel.New("lastname")
		sel.Build(&sb)

		expClause := "firstname, lastname"
		expect(t, expClause, sb.String())
	})

	t.Run("select_with_alias", func(t *testing.T) {
		sel := clause.SelectBuilder{}
		var sb strings.Builder

		sel.New("firstname", "fn")
		sel.New("lastname", "ln")
		sel.Build(&sb)

		expClause := "firstname AS fn, lastname AS ln"
		expect(t, expClause, sb.String())
	})

	t.Run("select_with_alias", func(t *testing.T) {
		sel := clause.SelectBuilder{}
		var sb strings.Builder

		sel.New("firstname", "fn")
		sel.New("lastname", "ln")
		sel.New("age")
		sel.New("email", "mail")
		sel.Build(&sb)

		expClause := "firstname AS fn, lastname AS ln, age, email AS mail"
		expect(t, expClause, sb.String())
	})
}
