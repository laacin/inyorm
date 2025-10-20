package clause_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/laacin/inyorm/clause"
)

func TestWhere(t *testing.T) {
	expect := func(t *testing.T, exp, have string) {
		if exp != have {
			t.Errorf("mismatch clause:\nExpect: %s\nHave: %s\n", exp, have)
		}
	}

	expectValues := func(t *testing.T, exp, have []any) {
		if !reflect.DeepEqual(exp, have) {
			t.Errorf("mismatch values:\nExpect: %#v\nHave: %#v\n", exp, have)
		}
	}

	t.Run("Equal_and_Not", func(t *testing.T) {
		ph := &clause.PlaceholderGen{Kind: clause.Simple}
		wb := &clause.WhereBuilder{Ph: ph, Table: "main_table"}

		wb.NewExpr("column1").Equal(42).Or().Not().Equal(50)
		var sb strings.Builder

		values := wb.Build(&sb)
		expCls := "WHERE "
		expCls += "(main_table.column1 = ? OR main_table.column1 != ?)"
		valuesExp := []any{42, 50}

		expect(t, expCls, sb.String())
		expectValues(t, valuesExp, values)
	})

	t.Run("Like_and_NotLike", func(t *testing.T) {
		ph := &clause.PlaceholderGen{Kind: clause.Simple}
		wb := &clause.WhereBuilder{Ph: ph, Table: "main_table"}

		wb.NewExpr("name").Like("%abc%").And("surname").Not().Like("%xyz%")
		var sb strings.Builder

		values := wb.Build(&sb)
		expCls := "WHERE "
		expCls += "(main_table.name LIKE ? AND main_table.surname NOT LIKE ?)"
		valuesExp := []any{"%abc%", "%xyz%"}

		expect(t, expCls, sb.String())
		expectValues(t, valuesExp, values)
	})

	t.Run("In_and_NotIn", func(t *testing.T) {
		ph := &clause.PlaceholderGen{Kind: clause.Simple}
		wb := &clause.WhereBuilder{Ph: ph, Table: "main_table"}

		wb.NewExpr("age").In(20, 25, 30).Or().Not().In(35, 40)
		var sb strings.Builder

		values := wb.Build(&sb)
		expCls := "WHERE "
		expCls += "(main_table.age IN (?, ?, ?) OR main_table.age NOT IN (?, ?))"
		valuesExp := []any{20, 25, 30, 35, 40}

		expect(t, expCls, sb.String())
		expectValues(t, valuesExp, values)
	})

	t.Run("Between_and_NotBetween", func(t *testing.T) {
		ph := &clause.PlaceholderGen{Kind: clause.Simple}
		wb := &clause.WhereBuilder{Ph: ph, Table: "main_table"}

		wb.NewExpr("score").Between(50, 100).And("level").Not().Between(10, 20)
		var sb strings.Builder

		values := wb.Build(&sb)
		expCls := "WHERE "
		expCls += "(main_table.score BETWEEN ? AND ? AND main_table.level NOT BETWEEN ? AND ?)"
		valuesExp := []any{50, 100, 10, 20}

		expect(t, expCls, sb.String())
		expectValues(t, valuesExp, values)
	})

	t.Run("Greater_and_Less", func(t *testing.T) {
		ph := &clause.PlaceholderGen{Kind: clause.Simple}
		wb := &clause.WhereBuilder{Ph: ph, Table: "main_table"}

		wb.NewExpr("price").Greater(100).And("quantity").Less(50)
		var sb strings.Builder

		values := wb.Build(&sb)
		expCls := "WHERE "
		expCls += "(main_table.price > ? AND main_table.quantity < ?)"
		valuesExp := []any{100, 50}

		expect(t, expCls, sb.String())
		expectValues(t, valuesExp, values)
	})

	t.Run("Greater_and_NotLess", func(t *testing.T) {
		ph := &clause.PlaceholderGen{Kind: clause.Simple}
		wb := &clause.WhereBuilder{Ph: ph, Table: "main_table"}

		wb.NewExpr("weight").Greater(10).Or().Not().Less(5)
		var sb strings.Builder

		values := wb.Build(&sb)
		expCls := "WHERE "
		expCls += "(main_table.weight > ? OR main_table.weight >= ?)"
		valuesExp := []any{10, 5}

		expect(t, expCls, sb.String())
		expectValues(t, valuesExp, values)
	})

	t.Run("IsNull_and_NotIsNull", func(t *testing.T) {
		ph := &clause.PlaceholderGen{Kind: clause.Simple}
		wb := &clause.WhereBuilder{Ph: ph, Table: "main_table"}

		wb.NewExpr("deleted_at").IsNull().Or("updated_at").Not().IsNull()
		var sb strings.Builder

		values := wb.Build(&sb)
		expCls := "WHERE "
		expCls += "(main_table.deleted_at IS NULL OR main_table.updated_at IS NOT NULL)"
		var valuesExp []any

		expect(t, expCls, sb.String())
		expectValues(t, valuesExp, values)
	})

	t.Run("Mix_all_methods_some_aliases", func(t *testing.T) {
		ph := &clause.PlaceholderGen{Kind: clause.Simple}
		wb := &clause.WhereBuilder{Ph: ph, Table: "main_table"}

		wb.NewExpr("column1").Equal(1).Or().Not().Greater(5)
		wb.NewExpr("column2", "t2").Like("%x%").And("column3").Not().Like("%y%")
		wb.NewExpr("column4").In(1, 2, 3).And("column5", "t5").Not().In(4, 5)
		wb.NewExpr("column6").Between(10, 20).And("column7").Not().Between(30, 40)
		wb.NewExpr("column8", "t8").Greater(100).And("column9").Less(50)
		wb.NewExpr("column10").IsNull().Or("column11", "t11").Not().IsNull()
		var sb strings.Builder

		values := wb.Build(&sb)
		expCls := "WHERE "
		expCls += "(main_table.column1 = ? OR main_table.column1 <= ?) AND "
		expCls += "(t2.column2 LIKE ? AND t2.column3 NOT LIKE ?) AND "
		expCls += "(main_table.column4 IN (?, ?, ?) AND t5.column5 NOT IN (?, ?)) AND "
		expCls += "(main_table.column6 BETWEEN ? AND ? AND main_table.column7 NOT BETWEEN ? AND ?) AND "
		expCls += "(t8.column8 > ? AND t8.column9 < ?) AND "
		expCls += "(main_table.column10 IS NULL OR t11.column11 IS NOT NULL)"

		valuesExp := []any{
			1, 5,
			"%x%", "%y%",
			1, 2, 3, 4, 5,
			10, 20, 30, 40,
			100, 50,
		}

		expect(t, expCls, sb.String())
		expectValues(t, valuesExp, values)
	})

	t.Run("Mix_all_methods_some_aliases_numeric_placeholders", func(t *testing.T) {
		ph := &clause.PlaceholderGen{Kind: clause.Numbered}
		wb := &clause.WhereBuilder{Ph: ph, Table: "main_table"}

		wb.NewExpr("column1").Equal(1).Or().Not().Greater(5)
		wb.NewExpr("column2", "t2").Like("%x%").And("column3").Not().Like("%y%")
		wb.NewExpr("column4").In(1, 2, 3).And("column5", "t5").Not().In(4, 5)
		wb.NewExpr("column6").Between(10, 20).And("column7").Not().Between(30, 40)
		wb.NewExpr("column8", "t8").Greater(100).And("column9").Less(50)
		wb.NewExpr("column10").IsNull().Or("column11", "t11").Not().IsNull()
		var sb strings.Builder

		values := wb.Build(&sb)
		expCls := "WHERE "
		expCls += "(main_table.column1 = $1 OR main_table.column1 <= $2) AND "
		expCls += "(t2.column2 LIKE $3 AND t2.column3 NOT LIKE $4) AND "
		expCls += "(main_table.column4 IN ($5, $6, $7) AND t5.column5 NOT IN ($8, $9)) AND "
		expCls += "(main_table.column6 BETWEEN $10 AND $11 AND main_table.column7 NOT BETWEEN $12 AND $13) AND "
		expCls += "(t8.column8 > $14 AND t8.column9 < $15) AND "
		expCls += "(main_table.column10 IS NULL OR t11.column11 IS NOT NULL)"

		valuesExp := []any{
			1, 5,
			"%x%", "%y%",
			1, 2, 3, 4, 5,
			10, 20, 30, 40,
			100, 50,
		}

		expect(t, expCls, sb.String())
		expectValues(t, valuesExp, values)
	})
}
