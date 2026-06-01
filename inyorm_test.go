package inyorm_test

import (
	"reflect"
	"testing"

	"github.com/laacin/inyorm"
	"github.com/laacin/inyorm/engine/std"
)

func run(t *testing.T, stmt inyorm.Runner, exp string, vals []any) {
	query, values, err := stmt.Raw()
	if err != nil {
		t.Fatal(err)
	}

	if query != exp {
		t.Errorf("mismatch query:\nExpect:\n%s\nHave:\n%s", exp, query)
	}

	if !reflect.DeepEqual(values, vals) {
		t.Errorf("mismatch values:\nExpect:\n%#v\nHave:\n%#v", vals, values)
	}
}

func TestSelect(t *testing.T) {
	db, _ := inyorm.New(std.JustDialect())

	t.Run("simple", func(t *testing.T) {
		stmt := db.Select(func(q inyorm.SelectQuery, e inyorm.Expr) {
			q.Select(e.All())
			q.From(e.Table("users"))
			q.Where(e.Col("id")).Equal(e.Param("uuid"))
			q.Limit(1)
		})

		exp := "SELECT * FROM users WHERE (id = ?) LIMIT 1"

		run(t, stmt, exp, []any{"uuid"})
	})

	t.Run("pagination", func(t *testing.T) {
		stmt := db.Select(func(q inyorm.SelectQuery, e inyorm.Expr) {
			var (
				id      = e.Col("id")
				age     = e.Col("age")
				banned  = e.Col("banned")
				foreign = e.Col("user_id", "posts")
			)

			q.Select(e.All())
			q.From(e.Table("users"))
			q.Join(e.Table("posts")).On(foreign).Equal(id)
			q.Where(banned).IsNull().And(age).Greater(e.Param(17))
			q.OrderBy(age).Desc()
			q.Limit(100)
			q.Offset(20)
		})

		exp := "SELECT a.* "
		exp += "FROM users a "
		exp += "INNER JOIN posts b ON (b.user_id = a.id) "
		exp += "WHERE (a.banned IS NULL AND a.age > ?) "
		exp += "ORDER BY a.age DESC "
		exp += "LIMIT 100 OFFSET 20"

		run(t, stmt, exp, []any{17})
	})

	t.Run("complex", func(t *testing.T) {
		stmt := db.Select(func(q inyorm.SelectQuery, e inyorm.Expr) {
			var (
				banned  = e.Col("banned")
				fname   = e.Col("firstname")
				lname   = e.Col("lastname")
				age     = e.Col("age")
				postNum = e.Col("id", "posts").Count()
				role    = e.Col("name", "roles")
				lastLog = e.Col("last_login")

				id        = e.Col("id")
				postsFk   = e.Col("user_id", "posts")
				interUser = e.Col("user_id", "user_roles")
				interRole = e.Col("role_id", "user_roles")
				roleId    = e.Col("id", "roles")
			)

			success := e.Concat(
				"with role: ", role,
				" has ", postNum, " posts and",
				" his last login was: ", lastLog,
			)

			info := e.Search(func(cs inyorm.Case) {
				cs.When(e.Cond(banned).IsNull().And(banned)).Then(success)
				cs.Else(e.Concat("was banned at: ", banned))
			})

			result := e.Concat("User: ", fname, " ", lname, " ", info).As("user_info")

			q.Select(result)
			q.From(e.Table("users"))
			q.Join(e.Table("posts")).On(postsFk).Equal(id)
			q.Join(e.Table("user_roles")).On(interUser).Equal(id)
			q.Join(e.Table("roles")).On(roleId).Equal(interRole)
			q.Where(age).Greater(e.Param(17)).And(age).Less(e.Param(30))
			q.GroupBy(e.Col("id", "posts"))
			q.Having(postNum).Greater(10)
			q.OrderBy(age).Desc()
			q.Limit(100)
			q.Offset(20)
		})

		exp := "SELECT "
		exp += "CONCAT('User: ', a.firstname, ' ', a.lastname, ' ', "
		exp += "CASE WHEN (a.banned IS NULL) THEN "
		exp += "CONCAT('with role: ', d.name, ' has ', COUNT(b.id), ' posts and', ' his last login was: ', a.last_login) "
		exp += "ELSE CONCAT('was banned at: ', a.banned) END) AS user_info "
		exp += "FROM users a "
		exp += "INNER JOIN posts b ON (b.user_id = a.id) "
		exp += "INNER JOIN user_roles c ON (c.user_id = a.id) "
		exp += "INNER JOIN roles d ON (d.id = c.role_id) "
		exp += "WHERE (a.age > ? AND a.age < ?) "
		exp += "GROUP BY b.id "
		exp += "HAVING (COUNT(b.id) > 10) "
		exp += "ORDER BY a.age DESC "
		exp += "LIMIT 100 OFFSET 20"

		vals := []any{17, 30}

		run(t, stmt, exp, vals)
	})

	t.Run("complex_with_helper", func(t *testing.T) {
		stmt := db.Select(func(q inyorm.SelectQuery, e inyorm.Expr) {
			// helper
			isNull := func(cond inyorm.Col, then, els any) inyorm.Col {
				return e.Search(func(cs inyorm.Case) {
					cond := e.Cond(cond).IsNull()
					cs.When(cond).Then(then)
					cs.Else(els)
				})
			}

			var (
				// base columns
				fname     = e.Col("firstname")
				lname     = e.Col("lastname")
				lastLogin = e.Col("last_login")
				banned    = e.Col("banned")
				roleName  = e.Col("name", "roles")
				age       = e.Col("age")
				active    = e.Col("active")

				// summary columns
				posts    = e.Col("id", "posts").Count()
				comments = e.Col("id", "comments").Count()
				role     = isNull(roleName, "No role", roleName)
				lastLog  = isNull(lastLogin, "Never", lastLogin)
				status   = isNull(banned, "Active", e.Concat("Banned at: ", banned))

				// join keys
				userId        = e.Col("id")
				roleId        = e.Col("id", "roles")
				postUserFk    = e.Col("user_id", "posts")
				commentUserFk = e.Col("user_id", "comments")

				interUserFk = e.Col("user_id", "user_roles")
				interRoleFk = e.Col("role_id", "user_roles")

				// select columns
				totalPost    = e.Col("id", "posts").Count(true).As("total_posts")
				totalComment = e.Col("id", "comments").Count(true).As("total_comments")
				lastPost     = e.Col("created_at", "posts").Max().As("last_post_date")
				summary      = e.Concat(
					"User: ", fname, " ", lname, " | ",
					"Role: ", role, " | ",
					"Posts: ", posts, " | ",
					"Comments: ", comments, " | ",
					"Last login: ", lastLog, " | ",
					status,
				).As("user_summary")
			)

			// statement building
			q.Select(summary, totalPost, totalComment, lastPost)
			q.From(e.Table("users"))

			q.Join(e.Table("user_roles")).Left().On(interUserFk).Equal(userId)
			q.Join(e.Table("roles")).Left().On(roleId).Equal(interRoleFk)
			q.Join(e.Table("posts")).Left().On(postUserFk).Equal(userId)
			q.Join(e.Table("comments")).Left().On(commentUserFk).Equal(userId)

			q.Where(age).Between(18, 60).And(active).Equal(true)

			q.GroupBy(userId, fname, lname, lastLogin, banned, roleName)
			q.Having(posts).Greater(5)

			q.OrderBy(totalPost).Desc()
			q.OrderBy(age)
			q.Limit(50)
		})

		exp := "SELECT CONCAT("
		exp += "'User: ', a.firstname, ' ', a.lastname, ' | ', "
		exp += "'Role: ', CASE WHEN (c.name IS NULL) THEN 'No role' ELSE c.name END, ' | ', "
		exp += "'Posts: ', COUNT(d.id), ' | ', "
		exp += "'Comments: ', COUNT(e.id), ' | ', "
		exp += "'Last login: ', CASE WHEN (a.last_login IS NULL) THEN 'Never' ELSE a.last_login END, ' | ', "
		exp += "CASE WHEN (a.banned IS NULL) THEN 'Active' ELSE CONCAT('Banned at: ', a.banned) END"
		exp += ") AS user_summary, "
		exp += "COUNT(DISTINCT d.id) AS total_posts, "
		exp += "COUNT(DISTINCT e.id) AS total_comments, "
		exp += "MAX(d.created_at) AS last_post_date "
		exp += "FROM users a "
		exp += "LEFT JOIN user_roles b ON (b.user_id = a.id) "
		exp += "LEFT JOIN roles c ON (c.id = b.role_id) "
		exp += "LEFT JOIN posts d ON (d.user_id = a.id) "
		exp += "LEFT JOIN comments e ON (e.user_id = a.id) "
		exp += "WHERE (a.age BETWEEN 18 AND 60 AND a.active = 1) "
		exp += "GROUP BY a.id, a.firstname, a.lastname, a.last_login, a.banned, c.name "
		exp += "HAVING (COUNT(d.id) > 5) "
		exp += "ORDER BY total_posts DESC, a.age "
		exp += "LIMIT 50"

		run(t, stmt, exp, []any{})
	})
}

func TestInsert(t *testing.T) {
	db, _ := inyorm.New(std.JustDialect())

	type User struct {
		Account string
		Age     int
	}

	t.Run("insert_one", func(t *testing.T) {
		stmt := db.Insert(func(q inyorm.InsertQuery, e inyorm.Expr) {
			q.Insert(User{})
			q.Into(e.Table("users"))
			q.Values(User{
				Account: "myacc",
				Age:     29,
			})
		})

		exp := "INSERT INTO users (account, age) VALUES (?, ?)"
		run(t, stmt, exp, []any{"myacc", 29})
	})

	t.Run("insert_many", func(t *testing.T) {
		stmt := db.Insert(func(q inyorm.InsertQuery, e inyorm.Expr) {
			q.Insert(User{})
			q.Into(e.Table("users"))
			q.Values([]User{
				{Account: "acc1", Age: 10},
				{Account: "acc2", Age: 20},
				{Account: "acc3", Age: 30},
				{Account: "acc4", Age: 40},
				{Account: "acc5", Age: 50},
				{Account: "acc6", Age: 60},
			})
		})

		args := []any{
			"acc1", 10,
			"acc2", 20,
			"acc3", 30,
			"acc4", 40,
			"acc5", 50,
			"acc6", 60,
		}

		exp := "INSERT INTO users (account, age) VALUES "
		exp += "(?, ?), (?, ?), (?, ?), (?, ?), (?, ?), (?, ?)"
		run(t, stmt, exp, args)
	})

	t.Run("omit_values", func(t *testing.T) {
		vals := []map[string]any{
			{"account": "acc1", "age": 10, "active": true, "score": 100, "country": "AR"},
			{"account": "acc2", "age": 20, "active": false, "score": 200, "country": "US"},
			{"account": "acc3", "age": 30, "active": true, "score": 300, "country": "BR"},
			{"account": "acc4", "age": 40, "active": false, "score": 400, "country": "UK"},
			{"account": "acc5", "age": 50, "active": true, "score": 500, "country": "DE"},
			{"account": "acc6", "age": 60, "active": true, "score": 600, "country": "JP"},
		}

		stmt := db.Insert(func(q inyorm.InsertQuery, e inyorm.Expr) {
			q.Insert(e.Col("score"), e.Col("age"))
			q.Into(e.Table("users"))
			q.Values(&vals)
		})

		args := []any{
			10, 100,
			20, 200,
			30, 300,
			40, 400,
			50, 500,
			60, 600,
		}

		exp := "INSERT INTO users (age, score) VALUES "
		exp += "(?, ?), (?, ?), (?, ?), (?, ?), (?, ?), (?, ?)"
		run(t, stmt, exp, args)
	})

	t.Run("ignore_values", func(t *testing.T) {
		stmt := db.Insert(func(q inyorm.InsertQuery, e inyorm.Expr) {
			q.Insert(User{}).Ignore(e.Col("age"))
			q.Into(e.Table("users"))
			q.Values(User{
				Account: "acc123",
				Age:     42,
			})
		})

		args := []any{"acc123"}

		exp := "INSERT INTO users (account) VALUES (?)"
		run(t, stmt, exp, args)
	})
}

func TestUpdate(t *testing.T) {
	db, _ := inyorm.New(std.JustDialect())

	type Post struct {
		Title       string
		Description string
	}

	t.Run("update_one", func(t *testing.T) {
		stmt := db.Update(func(q inyorm.UpdateQuery, e inyorm.Expr) {
			q.Update(&Post{})
			q.Into(e.Table("posts"))
			q.Values(Post{
				Title:       "something else",
				Description: "little description",
			})
			q.Where(e.Col("id")).Equal(e.Param(10))
		})

		exp := "UPDATE posts SET description = ?, title = ? WHERE (id = ?)"
		run(t, stmt, exp, []any{"little description", "something else", 10})
	})

	t.Run("with_cols", func(t *testing.T) {
		stmt := db.Update(func(q inyorm.UpdateQuery, e inyorm.Expr) {
			q.Update(e.Col("title"), e.Col("description"))
			q.Into(e.Table("posts"))
			q.Values(Post{
				Title: "asd", Description: "dep",
			})
			q.Where(e.Col("id")).Equal(e.Param(10))
		})

		exp := "UPDATE posts SET description = ?, title = ? WHERE (id = ?)"
		run(t, stmt, exp, []any{"dep", "asd", 10})
	})

	t.Run("with_map", func(t *testing.T) {
		vals := map[string]any{
			"account": "acc123",
			"age":     56,
			"name":    "matias",
		}

		stmt := db.Update(func(q inyorm.UpdateQuery, e inyorm.Expr) {
			q.Update(vals)
			q.Into(e.Table("users"))
			q.Values(vals)
		})

		exp := "UPDATE users SET account = ?, age = ?, name = ?"
		run(t, stmt, exp, []any{"acc123", 56, "matias"})
	})

	t.Run("omit_values", func(t *testing.T) {
		vals := map[string]any{
			"account":  "acc123",
			"age":      56,
			"name":     "matias",
			"lastname": "doe",
			"id":       123,
		}

		stmt := db.Update(func(q inyorm.UpdateQuery, e inyorm.Expr) {
			var (
				name = e.Col("name")
				acc  = e.Col("account")
				age  = e.Col("age")
			)

			q.Update(name, acc, age)
			q.Into(e.Table("users"))
			q.Values(vals)
		})

		exp := "UPDATE users SET account = ?, age = ?, name = ?"
		run(t, stmt, exp, []any{"acc123", 56, "matias"})
	})
}

func TestDelete(t *testing.T) {
	db, _ := inyorm.New(std.JustDialect())

	t.Run("delete_one", func(t *testing.T) {
		stmt := db.Delete(func(q inyorm.DeleteQuery, e inyorm.Expr) {
			q.Delete()
			q.From(e.Table("comments"))
			q.Where(e.Col("id")).Equal(e.Param(12310))
		})

		exp := "DELETE FROM comments WHERE (id = ?)"
		run(t, stmt, exp, []any{12310})
	})
}

func TestCreateTable(t *testing.T) {
	db, _ := inyorm.New(std.JustDialect())

	t.Run("basic_table", func(t *testing.T) {
		stmt := db.CreateTable(func(q inyorm.CreateTable) {
			q.TableName("users")

			q.Int("id").PrimaryKey().AutoIncrement()
			q.String("account").Unique()
		})

		exp := "CREATE TABLE IF NOT EXISTS users ("
		exp += "id INTEGER PRIMARY KEY AUTOINCREMENT, "
		exp += "account TEXT UNIQUE NOT NULL"
		exp += ")"

		run(t, stmt, exp, []any{})
	})

	t.Run("with_constraints", func(t *testing.T) {
		stmt := db.CreateTable(func(q inyorm.CreateTable) {
			q.TableName("posts")

			q.String("id").PrimaryKey()
			q.String("author_id")
			q.String("title").Default("untitled")
			q.String("description").Nullable()

			q.ForeignKey("author_id").To("id", "users").OnDel("cascade")
			q.Check("description").Not().Like("%word%")
		})

		exp := "CREATE TABLE IF NOT EXISTS posts ("
		exp += "id TEXT PRIMARY KEY, "
		exp += "author_id TEXT NOT NULL, "
		exp += "title TEXT NOT NULL DEFAULT 'untitled', "
		exp += "description TEXT, "
		exp += "FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE, "
		exp += "CHECK (description NOT LIKE '%word%')"
		exp += ")"

		run(t, stmt, exp, []any{})
	})

	t.Run("relational_table", func(t *testing.T) {
		stmt := db.CreateTable(func(q inyorm.CreateTable) {
			q.TableName("user_roles")

			q.String("user_id")
			q.String("role_id")

			q.PrimaryKey("user_id", "role_id")
			q.ForeignKey("user_id").To("id", "users").OnDel("cascade")
			q.ForeignKey("role_id").To("id", "roles").OnDel("cascade")
		})

		exp := "CREATE TABLE IF NOT EXISTS user_roles ("
		exp += "user_id TEXT NOT NULL, "
		exp += "role_id TEXT NOT NULL, "
		exp += "PRIMARY KEY (user_id, role_id), "
		exp += "FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE, "
		exp += "FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE"
		exp += ")"

		run(t, stmt, exp, []any{})
	})
}
