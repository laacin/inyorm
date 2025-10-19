package clause_test

import (
	"strings"
	"testing"

	"github.com/laacin/inyorm/clause"
)

func TestJoin(t *testing.T) {
	expect := func(t *testing.T, exp, have string) {
		if exp != have {
			t.Errorf("mismatch clause:\nExpect: %s\nHave: %s\n", exp, have)
		}
	}

	t.Run("n-1_join", func(t *testing.T) {
		jb := clause.JoinBuilder{Table: "main", PrimaryKey: "id"}
		var sb strings.Builder

		jb.Simple(clause.InnerJoin, "other", "main_id")
		jb.Simple(clause.LeftJoin, "another", "foreign_id")

		errs := jb.Build(&sb)
		for _, err := range errs {
			if err != nil {
				t.Error(err)
				return
			}
		}

		expCls := "INNER JOIN other other ON other.main_id = main.id "
		expCls += "LEFT JOIN another another ON another.foreign_id = main.id"

		expect(t, expCls, sb.String())
	})

	t.Run("with_intermediate_table", func(t *testing.T) {
		jb := clause.JoinBuilder{Table: "main", PrimaryKey: "id"}
		var sb strings.Builder

		keys := map[string]string{
			"main":    "main_id",
			"other":   "other_id",
			"another": "another_id",
		}

		rel := jb.ManyToMany(clause.InnerJoin, "relation", keys)
		rel.Join(clause.InnerJoin, "other", "id")
		rel.Join(clause.InnerJoin, "another", "id")

		errs := jb.Build(&sb)
		for _, err := range errs {
			if err != nil {
				t.Error(err)
			}
		}

		expCls := "INNER JOIN relation relation ON relation.main_id = main.id "
		expCls += "INNER JOIN other other ON other.id = relation.other_id "
		expCls += "INNER JOIN another another ON another.id = relation.another_id"

		expect(t, expCls, sb.String())
	})

	t.Run("complex_many_to_many_with_self_and_cross", func(t *testing.T) {
		jb := clause.JoinBuilder{Table: "users", PrimaryKey: "id"}
		var sb strings.Builder

		jb.Simple(clause.LeftJoin, "profiles", "user_id")
		jb.Simple(clause.RightJoin, "addresses", "user_id")
		jb.Simple(clause.CrossJoin, "logs", "")

		keys := map[string]string{
			"users":  "user_id",
			"groups": "group_id",
			"roles":  "role_id",
		}

		rel := jb.ManyToMany(clause.InnerJoin, "memberships", keys)
		rel.Join(clause.LeftJoin, "groups", "id")
		rel.Join(clause.InnerJoin, "roles", "id")
		rel.Join(clause.InnerJoin, "users", "id")

		errs := jb.Build(&sb)
		for _, err := range errs {
			if err != nil {
				t.Error(err)
			}
		}

		expCls := "INNER JOIN memberships memberships ON memberships.user_id = users.id "
		expCls += "LEFT JOIN groups groups ON groups.id = memberships.group_id "
		expCls += "INNER JOIN roles roles ON roles.id = memberships.role_id "
		expCls += "INNER JOIN users users ON users.id = memberships.user_id "

		expCls += "LEFT JOIN profiles profiles ON profiles.user_id = users.id "
		expCls += "RIGHT JOIN addresses addresses ON addresses.user_id = users.id "
		expCls += "CROSS JOIN logs"

		expect(t, expCls, sb.String())
	})
}
