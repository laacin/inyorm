package column_test

import (
	"testing"

	"github.com/laacin/inyorm/column"
	"github.com/laacin/inyorm/iface"
)

func TestColumn(t *testing.T) {

	t.Run("switch_column", func(t *testing.T) {
		col := column.Column("column")
		col.Switch(col.Name(), func(c iface.Case[string]) {
			c.When("ok").Then(1)
			c.When("wrong").Then("valid")
			c.Else(nil)
		})

		got := string(col)
		expect := "CASE column WHEN 'ok' THEN 1 WHEN 'wrong' THEN 'valid' ELSE NULL END"
		if got != expect {
			t.Errorf("mismatch column\n Expect:\n %s\n Have:\n %s\n", expect, got)
		}
	})
}
