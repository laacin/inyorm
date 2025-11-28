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
		t.Error(err)
	}

	if query != exp {
		t.Errorf("mismatch query:\nExpect:\n%s\nHave:\n%s", exp, query)
	}

	if !reflect.DeepEqual(values, vals) {
		t.Errorf("mismatch values:\nExpect:\n%#v\nHave:\n%#v", values, vals)
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
}

func TestInsert(t *testing.T) {
	qe := inyorm.New("", nil, nil)

	type User struct {
		Account string `col:"account"`
		Age     int    `col:"age"`
	}

	t.Run("insert_one", func(t *testing.T) {
		q, _ := qe.NewInsert(context.Background(), "users")

		q.Insert(User{
			Account: "myacc",
			Age:     29,
		})

		exp := "INSERT INTO users (account, age) VALUES (?, ?)"
		run(t, q, exp, []any{"myacc", 29})
	})

	t.Run("insert_many", func(t *testing.T) {
		q, _ := qe.NewInsert(context.Background(), "users")

		q.Insert([]User{
			{Account: "acc1", Age: 10},
			{Account: "acc2", Age: 20},
			{Account: "acc3", Age: 30},
			{Account: "acc4", Age: 40},
			{Account: "acc5", Age: 50},
			{Account: "acc6", Age: 60},
		})

		exp := "INSERT INTO users (account, age) VALUES "
		exp += "(?, ?), (?, ?), (?, ?), (?, ?), (?, ?), (?, ?)"
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

		q.Update(Post{
			Title:       "something else",
			Description: "little description",
		})
		q.Where(c.Col("id")).Equal(c.Param(10))

		exp := "UPDATE posts SET description = ?, title = ? WHERE (id = ?)"
		run(t, q, exp, []any{"little description", "something else", 10})
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
