package clause_test

import (
	"testing"

	"github.com/laacin/inyorm/clause"
	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

func NewJoin() (*clause.Join, core.ColExpr, func(t *testing.T, cls string)) {
	stmt := writer.NewStatement("", "users")
	var c core.ColExpr = &column.ColExpr{Writer: stmt.Writer}
	cls := &clause.Join{}

	run := func(t *testing.T, clause string) {
		w := stmt.Writer()
		cls.Build(w)
		if val := w.ToString(); val != clause {
			t.Errorf("\nmismatch result:\nExpect:\n%s\nHave:\n%s\n", clause, val)
		}
	}

	return cls, c, run
}

func TestJoin(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cls, c, run := NewJoin()

		id := c.Col("id", "users")

		cls.Join("INNER", "posts")
		cls.On(c.Col("user_id", "posts")).Equal(id)

		exp := "INNER JOIN posts b ON (b.user_id = a.id)"
		run(t, exp)
	})

	t.Run("many_to_many", func(t *testing.T) {
		cls, c, run := NewJoin()
		id := c.Col("id", "users")

		cls.Join("INNER", "posts")
		cls.On(c.Col("user_id", "posts")).Equal(id)

		cls.Join("INNER", "user_roles")
		cls.On(c.Col("user_id", "user_roles")).Equal(id)

		cls.Join("INNER", "roles")
		cls.On(c.Col("id", "roles")).Equal(c.Col("role_id", "user_roles"))

		exp := "INNER JOIN posts b ON (b.user_id = a.id)"
		exp += " "
		exp += "INNER JOIN user_roles c ON (c.user_id = a.id)"
		exp += " "
		exp += "INNER JOIN roles d ON (d.id = c.role_id)"
		run(t, exp)
	})

	t.Run("complex_many_to_many_with_conditions", func(t *testing.T) {
		cls, c, run := NewJoin()

		var (
			id             = c.Col("id", "users")
			active         = c.Col("active", "users")
			roleName       = c.Col("name", "roles")
			postPublished  = c.Col("published", "posts")
			userIDPosts    = c.Col("user_id", "posts")
			userIDUserRole = c.Col("user_id", "user_roles")
			assignedAt     = c.Col("assigned_at", "user_roles")
			roleID         = c.Col("role_id", "user_roles")
			roleIDRoles    = c.Col("id", "roles")
		)

		cls.Join("INNER", "posts")
		on := cls.On(userIDPosts)
		on.Equal(id)
		on.And(postPublished)
		on.Equal(true)

		cls.Join("INNER", "user_roles")
		on = cls.On(userIDUserRole)
		on.Equal(id)
		on.And(assignedAt)
		on.Not()
		on.IsNull()

		cls.Join("INNER", "roles")
		on = cls.On(roleIDRoles)
		on.Equal(roleID)
		on.And(roleName)
		on.In([]any{"admin", "editor", "manager"})
		on.And(active)
		on.Equal(true)

		exp := "INNER JOIN posts b ON (b.user_id = a.id AND b.published = 1)"
		exp += " "
		exp += "INNER JOIN user_roles c ON (c.user_id = a.id AND c.assigned_at IS NOT NULL)"
		exp += " "
		exp += "INNER JOIN roles d ON (d.id = c.role_id AND d.name IN ('admin', 'editor', 'manager') AND a.active = 1)"
		run(t, exp)
	})

	t.Run("cross_join", func(t *testing.T) {
		cls, _, run := NewJoin()
		cls.Join("CROSS", "building")
		run(t, "CROSS JOIN building b")
	})
}
