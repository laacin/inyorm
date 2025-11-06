package column_test

import (
	"testing"

	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

func TestColumn(t *testing.T) {
	stmt := &writer.Statement{}
	var c core.ColExpr = &column.ColExpr{Statement: stmt, DefaultTable: "default"}

	expect := func(t *testing.T, fn func(core.Writer), expect string) {
		w := stmt.Writer()
		fn(w)
		if val := w.ToString(); val != expect {
			t.Errorf("\nmismatch result:\nExpect:\n%s\nHave:\n%s\n", expect, val)
		}
	}

	t.Run("simple", func(t *testing.T) {
		name := c.Col("name")
		expect(t, name.Def(), "a.name")

		name.As("n")
		expect(t, name.Def(), "name AS n")

		name.Count()
		expect(t, name.Def(), "COUNT(name) AS n")
		expect(t, name.Ref(), "n")
	})

	t.Run("operation", func(t *testing.T) {
		var (
			stock = c.Col("stock")
			price = c.Col("price")
		)
		total := stock.Mul(price).As("total")
		expect(t, total.Def(), "a.stock * a.price AS total")
		expect(t, total.Ref(), "total")
	})

	t.Run("variation", func(t *testing.T) {
		var (
			initPrice  = c.Col("initial_price", "products")
			finalPrice = c.Col("final_price", "products")
		)
		result := finalPrice.Sub(initPrice).Wrap().Div(initPrice).Mul(100).As("variation")
		expect(t, result.Def(), "(b.final_price - b.initial_price) / b.initial_price * 100 AS variation")
		expect(t, result.Ref(), "variation")
	})

	t.Run("scalar", func(t *testing.T) {
		firstname := c.Col("firstname")
		firstname.Lower().As("fname")
		expect(t, firstname.Def(), "LOWER(a.firstname) AS fname")
		expect(t, firstname.Ref(), "fname")
	})

	t.Run("aggregation", func(t *testing.T) {
		all := c.All().Count()
		expect(t, all.Def(), "COUNT(*)")
	})

	t.Run("concat", func(t *testing.T) {
		var (
			fname = c.Col("firstname")
			lname = c.Col("lastname")
		)
		fullName := c.Concat(fname, " ", lname).As("full_name")
		expect(t, fullName.Def(), "CONCAT(a.firstname, ' ', a.lastname) AS full_name")
		expect(t, fullName.Ref(), "full_name")
	})

	t.Run("switch", func(t *testing.T) {
		banned := c.Col("banned")
		cs := c.Switch(banned, func(cs core.CaseSwitch) {
			cs.When(true).Then("invalid")
			cs.Else("valid")
		}).As("is_valid")
		expect(t, cs.Def(), "CASE a.banned WHEN 1 THEN 'invalid' ELSE 'valid' END AS is_valid")
		expect(t, cs.Ref(), "is_valid")
	})

	t.Run("search", func(t *testing.T) {
		age := c.Col("age")
		cond := c.Condition(age).Less(18)
		valid := c.Search(func(cs core.CaseSearch) {
			cs.When(cond).Then(false).Else(true)
		}).As("valid")
		expect(t, valid.Def(), "CASE WHEN (a.age < 18) THEN 0 ELSE 1 END AS valid")
		expect(t, valid.Ref(), "valid")
	})

	t.Run("combine", func(t *testing.T) {
		var (
			banned  = c.Col("banned")
			fname   = c.Col("firstname")
			lname   = c.Col("lastname")
			role    = c.Col("name", "roles")
			postNum = c.Col("id", "posts").Count()
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
		exp += "CONCAT('with role: ', c.name, ' has ', COUNT(d.id), ' posts and', ' his last login was: ', a.last_login)"
		exp += " ELSE CONCAT('was banned at: ', a.banned) END) AS user_info"
		expect(t, result.Def(), exp)
		expect(t, result.Ref(), "user_info")
	})
}
