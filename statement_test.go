package inyorm_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/laacin/inyorm"
)

func run(t *testing.T, stmt any, exp string, vals []any) {
	query, values, err := stmt.(interface{ Raw() (string, []any, error) }).Raw()
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
	qe := inyorm.New("", nil, nil)

	t.Run("simple", func(t *testing.T) {
		q, c := qe.NewSelect(context.Background(), "users")

		q.Select(c.All())
		q.Where(c.Col("id")).Equal(c.Param("uuid"))
		q.Limit(1)

		exp := "SELECT * FROM users WHERE (id = ?) LIMIT 1"

		run(t, q, exp, []any{"uuid"})
	})

	t.Run("pagination", func(t *testing.T) {
		q, c := qe.NewSelect(context.Background(), "users")

		var (
			id      = c.Col("id")
			age     = c.Col("age")
			banned  = c.Col("banned")
			foreign = c.Col("user_id", "posts")
		)

		q.Select(c.All())
		q.From("users")
		q.Join("posts").On(foreign).Equal(id)
		q.Where(banned).IsNull().And(age).Greater(c.Param(17))
		q.OrderBy(age).Desc()
		q.Limit(100)
		q.Offset(20)

		exp := "SELECT a.* "
		exp += "FROM users a "
		exp += "INNER JOIN posts b ON (b.user_id = a.id) "
		exp += "WHERE (a.banned IS NULL AND a.age > ?) "
		exp += "ORDER BY a.age DESC "
		exp += "LIMIT 100 OFFSET 20"

		run(t, q, exp, []any{17})
	})

	t.Run("complex", func(t *testing.T) {
		q, c := qe.NewSelect(context.Background(), "users")

		var (
			banned  = c.Col("banned")
			fname   = c.Col("firstname")
			lname   = c.Col("lastname")
			age     = c.Col("age")
			postNum = c.Col("id", "posts").Count()
			role    = c.Col("name", "roles")
			lastLog = c.Col("last_login")

			id        = c.Col("id")
			postsFk   = c.Col("user_id", "posts")
			interUser = c.Col("user_id", "user_roles")
			interRole = c.Col("role_id", "user_roles")
			roleId    = c.Col("id", "roles")
		)

		success := c.Concat(
			"with role: ", role,
			" has ", postNum, " posts and",
			" his last login was: ", lastLog,
		)

		info := c.Search(func(cs inyorm.Case) {
			cs.When(c.Cond(banned).IsNull().And(banned)).Then(success)
			cs.Else(c.Concat("was banned at: ", banned))
		})

		result := c.Concat("User: ", fname, " ", lname, " ", info).As("user_info")

		q.Select(result)
		q.From("users")
		q.Join("posts").On(postsFk).Equal(id)
		q.Join("user_roles").On(interUser).Equal(id)
		q.Join("roles").On(roleId).Equal(interRole)
		q.Where(age).Greater(c.Param(17)).And(age).Less(c.Param(30))
		q.GroupBy(postNum.Base())
		q.Having(postNum).Greater(10)
		q.OrderBy(age).Desc()
		q.Limit(100)
		q.Offset(20)

		exp := "SELECT "
		exp += "CONCAT('User: ', a.firstname, ' ', a.lastname, ' ', "
		exp += "CASE WHEN (a.banned IS NULL) THEN "
		exp += "CONCAT('with role: ', b.name, ' has ', COUNT(c.id), ' posts and', ' his last login was: ', a.last_login) "
		exp += "ELSE CONCAT('was banned at: ', a.banned) END) AS user_info "
		exp += "FROM users a "
		exp += "INNER JOIN posts c ON (c.user_id = a.id) "
		exp += "INNER JOIN user_roles d ON (d.user_id = a.id) "
		exp += "INNER JOIN roles b ON (b.id = d.role_id) "
		exp += "WHERE (a.age > ? AND a.age < ?) "
		exp += "GROUP BY c.id "
		exp += "HAVING (COUNT(c.id) > 10) "
		exp += "ORDER BY a.age DESC "
		exp += "LIMIT 100 OFFSET 20"

		vals := []any{17, 30}

		run(t, q, exp, vals)
	})

	t.Run("complex_with_helper", func(t *testing.T) {
		q, c := qe.NewSelect(context.Background(), "users")

		// helper
		isNull := func(cond inyorm.Column, then, els any) inyorm.Column {
			return c.Search(func(cs inyorm.Case) {
				cond := c.Cond(cond).IsNull()
				cs.When(cond).Then(then)
				cs.Else(els)
			})
		}

		var (
			// base columns
			fname     = c.Col("firstname")
			lname     = c.Col("lastname")
			lastLogin = c.Col("last_login")
			banned    = c.Col("banned")
			roleName  = c.Col("name", "roles")
			age       = c.Col("age")
			active    = c.Col("active")

			// summary columns
			posts    = c.Col("id", "posts").Count()
			comments = c.Col("id", "comments").Count()
			role     = isNull(roleName, "No role", roleName)
			lastLog  = isNull(lastLogin, "Never", lastLogin)
			status   = isNull(banned, "Active", c.Concat("Banned at: ", banned))

			// join keys
			userId        = c.Col("id")
			roleId        = c.Col("id", "roles")
			postUserFk    = c.Col("user_id", "posts")
			commentUserFk = c.Col("user_id", "comments")

			interUserFk = c.Col("user_id", "user_roles")
			interRoleFk = c.Col("role_id", "user_roles")

			// select columns
			totalPost    = c.Col("id", "posts").Count(true).As("total_posts")
			totalComment = c.Col("id", "comments").Count(true).As("total_comments")
			lastPost     = c.Col("created_at", "posts").Max().As("last_post_date")
			summary      = c.Concat(
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

		q.Join("user_roles").Left().On(interUserFk).Equal(userId)
		q.Join("roles").Left().On(roleId).Equal(interRoleFk)
		q.Join("posts").Left().On(postUserFk).Equal(userId)
		q.Join("comments").Left().On(commentUserFk).Equal(userId)

		q.Where(age).Between(18, 60).And(active).Equal(true)

		q.GroupBy(userId, fname, lname, lastLogin, banned, roleName)
		q.Having(posts).Greater(5)

		q.OrderBy(totalPost).Desc()
		q.OrderBy(age)
		q.Limit(50)

		exp := "SELECT CONCAT("
		exp += "'User: ', a.firstname, ' ', a.lastname, ' | ', "
		exp += "'Role: ', CASE WHEN (b.name IS NULL) THEN 'No role' ELSE b.name END, ' | ', "
		exp += "'Posts: ', COUNT(c.id), ' | ', "
		exp += "'Comments: ', COUNT(d.id), ' | ', "
		exp += "'Last login: ', CASE WHEN (a.last_login IS NULL) THEN 'Never' ELSE a.last_login END, ' | ', "
		exp += "CASE WHEN (a.banned IS NULL) THEN 'Active' ELSE CONCAT('Banned at: ', a.banned) END"
		exp += ") AS user_summary, "
		exp += "COUNT(DISTINCT c.id) AS total_posts, "
		exp += "COUNT(DISTINCT d.id) AS total_comments, "
		exp += "MAX(c.created_at) AS last_post_date "
		exp += "FROM users a "
		exp += "LEFT JOIN user_roles e ON (e.user_id = a.id) "
		exp += "LEFT JOIN roles b ON (b.id = e.role_id) "
		exp += "LEFT JOIN posts c ON (c.user_id = a.id) "
		exp += "LEFT JOIN comments d ON (d.user_id = a.id) "
		exp += "WHERE (a.age BETWEEN 18 AND 60 AND a.active = 1) "
		exp += "GROUP BY a.id, a.firstname, a.lastname, a.last_login, a.banned, b.name "
		exp += "HAVING (COUNT(c.id) > 5) "
		exp += "ORDER BY total_posts DESC, a.age "
		exp += "LIMIT 50"

		run(t, q, exp, nil)
	})
}

func TestInsert(t *testing.T) {
	qe := inyorm.New("", nil, nil)

	type User struct {
		Account string `col:"account"`
		Age     int    `col:"age"`
	}

	t.Run("insert_one", func(t *testing.T) {
		q, _ := qe.NewInsert(context.Background(), "users")

		q.Insert(User{}).Values(User{
			Account: "myacc",
			Age:     29,
		})

		exp := "INSERT INTO users (account, age) VALUES (?, ?)"
		run(t, q, exp, []any{"myacc", 29})
	})

	t.Run("insert_many", func(t *testing.T) {
		q, _ := qe.NewInsert(context.Background(), "users")

		q.Insert(User{}).Values([]User{
			{Account: "acc1", Age: 10},
			{Account: "acc2", Age: 20},
			{Account: "acc3", Age: 30},
			{Account: "acc4", Age: 40},
			{Account: "acc5", Age: 50},
			{Account: "acc6", Age: 60},
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
		run(t, q, exp, args)
	})

	t.Run("omit_values", func(t *testing.T) {
		q, c := qe.NewInsert(context.Background(), "users")

		vals := []map[string]any{
			{"account": "acc1", "age": 10, "active": true, "score": 100, "country": "AR"},
			{"account": "acc2", "age": 20, "active": false, "score": 200, "country": "US"},
			{"account": "acc3", "age": 30, "active": true, "score": 300, "country": "BR"},
			{"account": "acc4", "age": 40, "active": false, "score": 400, "country": "UK"},
			{"account": "acc5", "age": 50, "active": true, "score": 500, "country": "DE"},
			{"account": "acc6", "age": 60, "active": true, "score": 600, "country": "JP"},
		}

		q.Insert(c.Col("score"), c.Col("age")).Values(&vals)

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
		run(t, q, exp, args)
	})
}

func TestUpdate(t *testing.T) {
	qe := inyorm.New("", nil, nil)

	type Post struct {
		Title       string `col:"title"`
		Description string `col:"description"`
	}

	t.Run("update_one", func(t *testing.T) {
		q, c := qe.NewUpdate(context.Background(), "posts")
		q.Update(&Post{}).Values(Post{
			Title:       "something else",
			Description: "little description",
		})
		q.Where(c.Col("id")).Equal(c.Param(10))

		exp := "UPDATE posts SET description = ?, title = ? WHERE (id = ?)"
		run(t, q, exp, []any{"little description", "something else", 10})
	})

	t.Run("with_cols", func(t *testing.T) {
		q, c := qe.NewUpdate(context.Background(), "posts")

		q.Update(c.Col("title"), c.Col("description")).Values(Post{
			Title: "asd", Description: "dep",
		})
		q.Where(c.Col("id")).Equal(c.Param(10))

		exp := "UPDATE posts SET description = ?, title = ? WHERE (id = ?)"
		run(t, q, exp, []any{"dep", "asd", 10})
	})

	t.Run("with_map", func(t *testing.T) {
		q, _ := qe.NewUpdate(context.Background(), "users")
		vals := map[string]any{
			"account": "acc123",
			"age":     56,
			"name":    "matias",
		}
		q.Update(vals).Values(vals)

		exp := "UPDATE users SET account = ?, age = ?, name = ?"
		run(t, q, exp, []any{"acc123", 56, "matias"})
	})

	t.Run("omit_values", func(t *testing.T) {
		q, c := qe.NewUpdate(context.Background(), "users")
		vals := map[string]any{
			"account":  "acc123",
			"age":      56,
			"name":     "matias",
			"lastname": "doe",
			"id":       123,
		}

		var (
			name = c.Col("name")
			acc  = c.Col("account")
			age  = c.Col("age")
		)

		q.Update(name, acc, age).Values(vals)

		exp := "UPDATE users SET account = ?, age = ?, name = ?"
		run(t, q, exp, []any{"acc123", 56, "matias"})
	})
}

func TestDelete(t *testing.T) {
	qe := inyorm.New("", nil, nil)

	t.Run("delete_one", func(t *testing.T) {
		q, c := qe.NewDelete(context.Background(), "comments")

		q.Where(c.Col("id")).Equal(c.Param(12310))

		exp := "DELETE FROM comments WHERE (id = ?)"
		run(t, q, exp, []any{12310})
	})
}
