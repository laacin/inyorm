package test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/laacin/inyorm"
	"github.com/laacin/inyorm/engine/sqlite"
)

type User struct {
	ID       int
	Account  string
	Password string
	Age      int
}

func Test(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := filepath.Join(cwd, "data_test.db")
	if !strings.HasSuffix(path, "inyorm/test/data_test.db") {
		panic("Wrong path")
	}

	fmt.Println(path)
	qe, err := inyorm.New(sqlite.Open(path))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := qe.Close(); err != nil {
			panic(err)
		}

		if err := os.Remove(path); err != nil {
			panic(err)
		}
	}()

	t.Run("create_table", func(t *testing.T) {
		stmt := qe.CreateTable("users", func(q inyorm.CreateTable, e inyorm.Expr) {
			q.Int("id").PrimaryKey().AutoIncrement()
			q.Text("account").Unique()
			q.Text("password")
			q.Int("age").Nullable()
		})

		if err := stmt.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("insert_values", func(t *testing.T) {
		stmt := qe.Insert("users", func(q inyorm.InsertQuery, e inyorm.Expr) {
			q.Insert(User{})
			q.Into(e.Table("users"))
			q.Values(User{
				ID:       1,
				Account:  "acc123",
				Password: "mysecret",
				Age:      21,
			})
		})

		if err := stmt.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get_values", func(t *testing.T) {
		stmt := qe.Select("users", func(q inyorm.SelectQuery, e inyorm.Expr) {
			q.Select(e.All())
			q.From(e.Table("users"))
			q.Where(e.Col("id")).Equal(e.Param(1))
			q.Limit(1)
		})

		var u User
		if err := stmt.Bind(&u).Run(); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(u, User{
			ID:       1,
			Account:  "acc123",
			Password: "mysecret",
			Age:      21,
		}) {
			t.Fatal("mismatch binding result")
		}
	})

	t.Run("create_table_with_constraints", func(t *testing.T) {
		stmt := qe.CreateTable("posts", func(q inyorm.CreateTable, e inyorm.Expr) {
			q.Int("id").PrimaryKey().AutoIncrement()
			q.Text("author_id")
			q.Text("title").Unique().Default("untitled")
			q.Text("description").Nullable()

			q.ForeignKey("author_id").To("id", "users").OnDel("cascade")
			q.Check(e.Col("description")).Not().Like("%someword%")
		})

		if err := stmt.Run(); err != nil {
			t.Fatal(err)
		}
	})
}
