package column_test

import (
	"testing"

	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

func New(defaultTable string) (core.ColExpr, func(t *testing.T, col core.Column, ref, def string)) {
	stmt := writer.NewStatement("")
	stmt.SetFrom(defaultTable)
	var c core.ColExpr = &column.ColExpr{Statement: stmt}

	run := func(t *testing.T, col core.Column, ref, def string) {
		r := stmt.Writer()
		d := stmt.Writer()

		col.Def()(d)
		col.Ref()(r)

		compare := func(name string, have, expect string) {
			if have != expect {
				t.Errorf("\nmismatch on %s:\nExpect:\n%s\nHave:\n%s\n", name, expect, have)
			}
		}

		compare("reference", r.ToString(), ref)
		compare("definition", d.ToString(), def)
	}
	return c, run
}

func TestColumn(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		c, run := New("users")

		name := c.Col("name")
		run(t, name, "a.name", "a.name")

		name.As("n")
		run(t, name, "n", "name AS n")

		name.Count()
		run(t, name, "n", "COUNT(name) AS n")
	})

	t.Run("operation", func(t *testing.T) {
		c, run := New("products")

		var (
			stock = c.Col("stock")
			price = c.Col("price")
		)

		total := stock.Mul(price).As("total")
		run(t, total, "total", "a.stock * a.price AS total")
	})

	t.Run("variation", func(t *testing.T) {
		c, run := New("users")

		var (
			initPrice  = c.Col("initial_price", "products")
			finalPrice = c.Col("final_price", "products")
		)
		result := finalPrice.Sub(initPrice).Wrap().Div(initPrice).Mul(100).As("variation")
		run(t, result, "variation", "(b.final_price - b.initial_price) / b.initial_price * 100 AS variation")
	})

	t.Run("scalar", func(t *testing.T) {
		c, run := New("users")

		firstname := c.Col("firstname")
		firstname.Lower().As("fname")
		run(t, firstname, "fname", "LOWER(a.firstname) AS fname")
	})

	t.Run("aggregation", func(t *testing.T) {
		c, run := New("users")

		all := c.All().Count()
		run(t, all, "COUNT(*)", "COUNT(*)")
	})

	t.Run("concat", func(t *testing.T) {
		c, run := New("users")

		var (
			fname = c.Col("firstname")
			lname = c.Col("lastname")
		)
		fullname := c.Concat(fname, " ", lname).As("fullname")
		run(t, fullname, "fullname", "CONCAT(a.firstname, ' ', a.lastname) AS fullname")
	})

	t.Run("switch", func(t *testing.T) {
		c, run := New("users")

		banned := c.Col("banned")
		cs := c.Switch(banned, func(cs core.CaseSwitch) {
			cs.When(true).Then("invalid")
			cs.Else("valid")
		}).As("is_valid")
		run(t, cs, "is_valid", "CASE a.banned WHEN 1 THEN 'invalid' ELSE 'valid' END AS is_valid")
	})

	t.Run("search", func(t *testing.T) {
		c, run := New("users")

		age := c.Col("age")
		cond := c.Condition(age).Less(18)
		valid := c.Search(func(cs core.CaseSearch) {
			cs.When(cond).Then(false).Else(true)
		}).As("valid")
		run(t, valid, "valid", "CASE WHEN (a.age < 18) THEN 0 ELSE 1 END AS valid")
	})

	t.Run("combine", func(t *testing.T) {
		c, run := New("users")

		var (
			banned  = c.Col("banned")
			fname   = c.Col("firstname")
			lname   = c.Col("lastname")
			postNum = c.Col("id", "posts").Count()
			role    = c.Col("name", "roles")
			lastLog = c.Col("last_login")
		)

		success := c.Concat(
			"with role: ", role,
			" has ", postNum, " posts and",
			" his last login was: ", lastLog,
		)

		cond := c.Condition(banned).IsNull()
		info := c.Search(func(cs core.CaseSearch) {
			cs.When(cond).Then(success)
			cs.Else(c.Concat("was banned at: ", banned))
		})

		result := c.Concat("User: ", fname, " ", lname, " ", info).As("user_info")

		exp := "CONCAT('User: ', a.firstname, ' ', a.lastname, ' ', "
		exp += "CASE WHEN (a.banned IS NULL) THEN "
		exp += "CONCAT('with role: ', c.name, ' has ', COUNT(b.id), ' posts and', ' his last login was: ', a.last_login)"
		exp += " ELSE CONCAT('was banned at: ', a.banned) END) AS user_info"
		run(t, result, "user_info", exp)
	})
}
