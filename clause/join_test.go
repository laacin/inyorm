package clause_test

import "testing"

func TestJoin(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cls, c, run := NewJoin()

		id := c.Col("id", "users")

		cls.Join("posts").On(c.Col("user_id", "posts")).Equal(id)

		exp := "INNER JOIN posts b ON (b.user_id = a.id)"
		run(t, exp, nil)
	})

	t.Run("many_to_many", func(t *testing.T) {
		cls, c, run := NewJoin()
		id := c.Col("id", "users")

		cls.Join("posts").On(c.Col("user_id", "posts")).Equal(id)

		cls.Join("user_roles").On(c.Col("user_id", "user_roles")).Equal(id)

		cls.Join("roles").On(c.Col("id", "roles")).Equal(c.Col("role_id", "user_roles"))

		exp := "INNER JOIN posts b ON (b.user_id = a.id)"
		exp += " "
		exp += "INNER JOIN user_roles c ON (c.user_id = a.id)"
		exp += " "
		exp += "INNER JOIN roles d ON (d.id = c.role_id)"
		run(t, exp, nil)
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

		cls.Join("posts").On(userIDPosts).
			Equal(id).And(postPublished).Equal(true)

		cls.Join("user_roles").On(userIDUserRole).
			Equal(id).And(assignedAt).Not().IsNull()

		cls.Join("roles").On(roleIDRoles).
			Equal(roleID).
			And(roleName).
			In([]any{"admin", "editor", "manager"}).
			And(active).
			Equal(true)

		exp := "INNER JOIN posts b ON (b.user_id = a.id AND b.published = 1)"
		exp += " "
		exp += "INNER JOIN user_roles c ON (c.user_id = a.id AND c.assigned_at IS NOT NULL)"
		exp += " "
		exp += "INNER JOIN roles d ON (d.id = c.role_id AND d.name IN ('admin', 'editor', 'manager') AND a.active = 1)"
		run(t, exp, nil)
	})
}
