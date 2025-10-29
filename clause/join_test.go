package clause_test

import (
	"strings"
	"testing"

	"github.com/laacin/inyorm/clause"
)

func TestJoin(t *testing.T) {
	var sb strings.Builder
	expect := func(t *testing.T, expect string) {
		if val := sb.String(); val != expect {
			t.Errorf("\nmismatch result:\nExpect:\n%s\nHave:\n%s\n", expect, val)
		}
		sb.Reset()
	}

	t.Run("simple_inner_join", func(t *testing.T) {
		cls := clause.NewJoinBuilder("users", "u", "id")
		cls.Join("INNER", "posts", "p", "user_id")

		if errs := cls.Build(&sb); errs != nil {
			for _, err := range errs {
				t.Error(err)
			}
		}

		exp := "INNER JOIN posts p ON p.user_id = u.id"
		expect(t, exp)
	})

	t.Run("many_to_many", func(t *testing.T) {
		cls := clause.NewJoinBuilder("users", "u", "id")
		inter := cls.Many("INNER", "user_permissions", "up", map[string]string{
			"users":       "user_id",
			"permissions": "permission_id",
		})
		inter.Join("FULL", "permissions", "p", "id")

		if errs := cls.Build(&sb); errs != nil {
			for _, err := range errs {
				t.Error(err)
			}
		}

		exp := "INNER JOIN user_permissions up ON up.user_id = u.id"
		exp += " "
		exp += "FULL JOIN permissions p ON p.id = up.permission_id"
		expect(t, exp)
	})
}
