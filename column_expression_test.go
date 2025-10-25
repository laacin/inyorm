package inyorm_test

import (
	"testing"

	"github.com/laacin/inyorm"
)

func expect(t *testing.T, fld *inyorm.CustomField, expect string) {
	if val := fld.Select(); val != expect {
		t.Errorf("\nmismatch result:\nExpect:\n%s\nHave:\n%s\n", expect, val)
	}
}

func TestFieldBuilder(t *testing.T) {
	var ce inyorm.ColumnExpr
	var (
		fname  inyorm.Field = "firstname"
		lname  inyorm.Field = "lastname"
		prod   inyorm.Field = "product"
		qty    inyorm.Field = "quantity"
		price  inyorm.Field = "price"
		age    inyorm.Field = "age"
		banned inyorm.Field = "banned"

		initVal  inyorm.Field = "initial_value"
		finalVal inyorm.Field = "final_value"
	)

	t.Run("concat_field", func(t *testing.T) {
		fld := ce.New("concat", func(fb *inyorm.FB) *inyorm.NdF {
			return fb.Concat("Name: ", fname.Use(), " ", lname.Use())
		})

		exp := "CONCAT('Name: ', firstname, ' ', lastname)"
		expect(t, fld, exp)
	})

	t.Run("operation", func(t *testing.T) {
		fld := ce.New("operation", func(fb *inyorm.FB) *inyorm.NdF {
			return fb.Simple(prod.Use()).Mul(price.Use())
		})

		exp := "product * price"
		expect(t, fld, exp)
	})

	t.Run("search", func(t *testing.T) {
		fld := ce.New("search", func(fb *inyorm.FB) *inyorm.NdF {
			return fb.Search(func(cs *inyorm.CaseField[*inyorm.ExprEnd]) {
				expr := fb.Expr(age.Use()).Less(18)
				cs.When(expr).Then("kid").Else("adult")
			})
		})

		exp := "CASE WHEN (age < 18) THEN 'kid' ELSE 'adult' END"
		expect(t, fld, exp)
	})

	t.Run("switch", func(t *testing.T) {
		fld := ce.New("switch", func(fb *inyorm.FB) *inyorm.NdF {
			return fb.Switch(age.Use(), func(cs *inyorm.CaseField[inyorm.Value]) {
				cs.When(14).Then("14 years old")
				cs.When(17).Then("17 years old")
				cs.Else("other age")
			})
		})

		exp := "CASE age WHEN 14 THEN '14 years old' WHEN 17 THEN '17 years old' ELSE 'other age' END"
		expect(t, fld, exp)
	})

	t.Run("combine", func(t *testing.T) {
		fld := ce.New("combine", func(fb *inyorm.FB) *inyorm.NdF {
			total := fb.Simple(qty.Use()).Mul(price.Use())
			res := fb.Concat(
				"User: ", fname.Use(), " ", lname.Use(),
				" had a total of ", qty.Use(), " of ", prod.Use(),
				" worth ", total,
			)

			return fb.Search(func(cs *inyorm.CaseField[*inyorm.ExprEnd]) {
				cond := fb.Expr(age.Use()).Less(18).Or(banned.Use()).Equal(true)
				cs.When(cond).Then("invalid user").Else(res)
			})
		})

		exp := "CASE WHEN (age < 18) OR (banned = 1) THEN 'invalid user' " // FIX: invalid expression wrapper
		exp += "ELSE CONCAT('User: ', firstname, ' ', lastname, ' had a total of ', quantity, ' of ', product, ' worth ', quantity * price)"
		exp += " END"
		expect(t, fld, exp)
	})

	t.Run("variation", func(t *testing.T) {
		fld := ce.New("variation", func(fb *inyorm.FB) *inyorm.NdF {
			return fb.Simple(finalVal.Use()).Sub(initVal.Use()).Wrap().Div(initVal.Use()).Mul(100)
		})

		exp := "(final_value - initial_value) / initial_value * 100"
		expect(t, fld, exp)
	})
}
