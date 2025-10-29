package clause_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/laacin/inyorm"
	"github.com/laacin/inyorm/clause"
)

func TestWhere(t *testing.T) {
	var (
		sb    strings.Builder
		fname inyorm.Field = "firstname"
		lname inyorm.Field = "lastname"
		age   inyorm.Field = "age"
	)

	expect := func(t *testing.T, expect string) {
		if val := sb.String(); val != expect {
			t.Errorf("\nmismatch result:\nExpect:\n%s\nHave:\n%s\n", expect, val)
		}
		sb.Reset()
	}

	expectVals := func(t *testing.T, have, expect any) {
		if !reflect.DeepEqual(have, expect) {
			t.Errorf("\nmissmatch values:\nExpect:\n%#v\nHave:\n%s\n", expect, have)
		}
	}

	t.Run("basic", func(t *testing.T) {
		cls := clause.WhereClause{}
		cls.Where(fname.Use()).Not().Not().Equal("john").Or().Not().Equal("mary")

		exp := "WHERE (firstname = 'john' OR firstname <> 'mary')"
		cls.Build(&sb, nil)
		expect(t, exp)
	})

	t.Run("with_placeholders", func(t *testing.T) {
		ph := clause.Placeholder{}
		cls := clause.WhereClause{}
		cls.Where(lname.Use()).Like("%alv%").And().Not().In("calvin", "malvina", "salvatore")

		exp := "WHERE (lastname LIKE ? AND lastname NOT IN (?, ?, ?))"
		cls.Build(&sb, &ph)
		expect(t, exp)
		expectVals(t, ph.Values(), []any{"%alv%", "calvin", "malvina", "salvatore"})
	})

	t.Run("postgres_placeholders", func(t *testing.T) {
		ph := clause.Placeholder{Dialect: "postgres"}
		cls := clause.WhereClause{}
		cls.Where(age.Use()).Between(17, 70).And().Not().Equal(45)

		exp := "WHERE (age BETWEEN $1 AND $2 AND age <> $3)"
		cls.Build(&sb, &ph)
		expect(t, exp)
		expectVals(t, ph.Values(), []any{17, 70, 45})
	})

	t.Run("many_wheres", func(t *testing.T) {
		ph := clause.Placeholder{}
		cls := clause.WhereClause{}
		cls.Where(age.Use()).Between(17, 70).And().Not().Equal(45)
		cls.Where(lname.Use()).Like("%alv%").And().Not().In("calvin", "malvina", "salvatore")
		cls.Where(fname.Use()).Not().Not().Equal("john").Or().Not().Equal("mary")
		cls.Where("literal").Not().IsNull()

		exp := "WHERE (age BETWEEN ? AND ? AND age <> ?)"
		exp += " AND "
		exp += "(lastname LIKE ? AND lastname NOT IN (?, ?, ?))"
		exp += " AND "
		exp += "(firstname = ? OR firstname <> ?)"
		exp += " AND "
		exp += "('literal' IS NOT NULL)"

		cls.Build(&sb, &ph)
		expect(t, exp)
		expectVals(t, ph.Values(), []any{17, 70, 45, "%alv%", "calvin", "malvina", "salvatore", "john", "mary"})
	})
}
