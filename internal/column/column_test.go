package column_test

import (
	"testing"

	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

type Dummy struct {
	Base  string
	Expr  string
	Alias string
	Def   string
}

func New(defaultTable string) (core.ColExpr, func(t *testing.T, col core.Column, d Dummy)) {
	stmt := writer.NewStatement("", defaultTable)
	var c core.ColExpr = &column.ColExpr{Writer: stmt.Writer}

	run := func(t *testing.T, col core.Column, dummy Dummy) {
		d := stmt.Writer()
		b := stmt.Writer()
		e := stmt.Writer()
		a := stmt.Writer()

		col.Def()(d)
		col.Alias()(a)
		col.Expr()(e)
		col.Base()(b)

		compare := func(name string, have, expect string) {
			if have != expect {
				t.Errorf("\nmismatch on %s:\nExpect:\n%s\nHave:\n%s\n", name, expect, have)
			}
		}

		compare("definition", d.ToString(), dummy.Def)
		compare("alias", a.ToString(), dummy.Alias)
		compare("expr", e.ToString(), dummy.Expr)
		compare("base", b.ToString(), dummy.Base)
	}
	return c, run
}

func TestColumn(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		c, run := New("users")

		name := c.Col("name", "users")
		run(t, name, Dummy{
			Base:  "a.name",
			Alias: "a.name",
			Expr:  "a.name",
			Def:   "a.name",
		})
	})

	t.Run("operation", func(t *testing.T) {
		c, run := New("products")

		var (
			stock = c.Col("stock", "products")
			price = c.Col("price", "products")
		)

		total := stock
		total.Mul(price)
		total.As("total")
		run(t, total, Dummy{
			Base:  "a.stock",
			Def:   "a.stock * a.price AS total",
			Expr:  "a.stock * a.price",
			Alias: "total",
		})
	})

	t.Run("variation", func(t *testing.T) {
		c, run := New("users")

		var (
			initPrice  = c.Col("initial_price", "products")
			finalPrice = c.Col("final_price", "products")
		)
		result := finalPrice
		result.Sub(initPrice)
		result.Wrap()
		result.Div(initPrice)
		result.Mul(100)
		result.As("variation")
		run(t, result, Dummy{
			Base:  "b.final_price",
			Alias: "variation",
			Def:   "(b.final_price - b.initial_price) / b.initial_price * 100 AS variation",
			Expr:  "(b.final_price - b.initial_price) / b.initial_price * 100",
		})
	})

	t.Run("scalar", func(t *testing.T) {
		c, run := New("users")

		firstname := c.Col("firstname", "users")
		firstname.Lower()
		firstname.As("fname")
		run(t, firstname, Dummy{
			Base:  "a.firstname",
			Expr:  "LOWER(a.firstname)",
			Def:   "LOWER(a.firstname) AS fname",
			Alias: "fname",
		})
	})

	t.Run("aggregation", func(t *testing.T) {
		c, run := New("users")

		all := c.All()
		all.Count(false)
		run(t, all, Dummy{
			Base:  "*",
			Expr:  "COUNT(*)",
			Alias: "COUNT(*)",
			Def:   "COUNT(*)",
		})
	})

	t.Run("concat", func(t *testing.T) {
		c, run := New("users")

		var (
			fname = c.Col("firstname", "users")
			lname = c.Col("lastname", "users")
		)
		fullname := c.Concat([]any{fname, " ", lname})
		fullname.As("fullname")
		run(t, fullname, Dummy{
			Def:   "CONCAT(a.firstname, ' ', a.lastname) AS fullname",
			Expr:  "CONCAT(a.firstname, ' ', a.lastname)",
			Alias: "fullname",
		})
	})

	t.Run("switch", func(t *testing.T) {
		c, run := New("users")

		banned := c.Col("banned", "users")

		cs := &column.Case{}
		cs.When(true)
		cs.Then("invalid")
		cs.Else("valid")

		info := c.Switch(banned, cs)
		info.As("is_valid")

		run(t, info, Dummy{
			Def:   "CASE a.banned WHEN 1 THEN 'invalid' ELSE 'valid' END AS is_valid",
			Expr:  "CASE a.banned WHEN 1 THEN 'invalid' ELSE 'valid' END",
			Alias: "is_valid",
		})
	})

	t.Run("search", func(t *testing.T) {
		c, run := New("users")

		age := c.Col("age", "users")
		cond := c.Condition(age)
		cond.Less(18)
		cs := &column.Case{}
		cs.When(cond)
		cs.Then(false)
		cs.Else(true)
		valid := c.Search(cs)
		valid.As("valid")
		run(t, valid, Dummy{
			Def:   "CASE WHEN (a.age < 18) THEN 0 ELSE 1 END AS valid",
			Expr:  "CASE WHEN (a.age < 18) THEN 0 ELSE 1 END",
			Alias: "valid",
		})
	})

	t.Run("combine", func(t *testing.T) {
		c, run := New("users")

		var (
			banned  = c.Col("banned", "users")
			fname   = c.Col("firstname", "users")
			lname   = c.Col("lastname", "users")
			postNum = c.Col("id", "posts")
			role    = c.Col("name", "roles")
			lastLog = c.Col("last_login", "users")
		)
		postNum.Count(false)

		success := c.Concat([]any{
			"with role: ", role,
			" has ", postNum, " posts and",
			" his last login was: ", lastLog,
		})

		cond := c.Condition(banned)
		cond.IsNull()
		cs := &column.Case{}
		cs.When(cond)
		cs.Then(success)
		cs.Else(c.Concat([]any{"was banned at: ", banned}))
		info := c.Search(cs)

		result := c.Concat([]any{"User: ", fname, " ", lname, " ", info})
		result.As("user_info")

		exp := "CONCAT('User: ', a.firstname, ' ', a.lastname, ' ', "
		exp += "CASE WHEN (a.banned IS NULL) THEN "
		exp += "CONCAT('with role: ', c.name, ' has ', COUNT(b.id), ' posts and', ' his last login was: ', a.last_login)"
		exp += " ELSE CONCAT('was banned at: ', a.banned) END)"
		run(t, result, Dummy{
			Expr:  exp,
			Def:   exp + " AS user_info",
			Alias: "user_info",
		})
	})
}
