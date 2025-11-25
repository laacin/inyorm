package clause_test

import "testing"

func TestJoin(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cls, c, run := NewJoin()

		id := c.Col("id", "users")

		cls.Join("posts").On(c.Col("user_id", "posts")).Equal(id)

		exp := "INNER JOIN posts x ON (x.user_id = x.id)"
		run(t, exp, nil)
	})

	t.Run("many_to_many", func(t *testing.T) {
		cls, c, run := NewJoin()
		id := c.Col("id", "users")

		cls.Join("posts").On(c.Col("user_id", "posts")).Equal(id)

		cls.Join("user_roles").On(c.Col("user_id", "user_roles")).Equal(id)

		cls.Join("roles").On(c.Col("id", "roles")).Equal(c.Col("role_id", "user_roles"))

		exp := "INNER JOIN posts x ON (x.user_id = x.id)"
		exp += " "
		exp += "INNER JOIN user_roles x ON (x.user_id = x.id)"
		exp += " "
		exp += "INNER JOIN roles x ON (x.id = x.role_id)"
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

		exp := "INNER JOIN posts x ON (x.user_id = x.id AND x.published = 1)"
		exp += " "
		exp += "INNER JOIN user_roles x ON (x.user_id = x.id AND x.assigned_at IS NOT NULL)"
		exp += " "
		exp += "INNER JOIN roles x ON (x.id = x.role_id AND x.name IN ('admin', 'editor', 'manager') AND x.active = 1)"
		run(t, exp, nil)
	})

	t.Run("different_join_types", func(t *testing.T) {
		cls, c, run := NewJoin()
		id := c.Col("id", "users")

		cls.Join("posts").Left().On(c.Col("user_id", "posts")).Equal(id)

		cls.Join("user_roles").Right().On(c.Col("user_id", "user_roles")).Equal(id)

		cls.Join("roles").Full().On(c.Col("id", "roles")).Equal(c.Col("role_id", "user_roles"))

		cls.Join("permissions").Cross()

		exp := "LEFT JOIN posts x ON (x.user_id = x.id)"
		exp += " "
		exp += "RIGHT JOIN user_roles x ON (x.user_id = x.id)"
		exp += " "
		exp += "FULL JOIN roles x ON (x.id = x.role_id)"
		exp += " "
		exp += "CROSS JOIN permissions x"

		run(t, exp, nil)
	})
}
