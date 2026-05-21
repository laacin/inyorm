package test

import (
	"fmt"
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

type Post struct {
	ID          int
	AuthorID    int
	Title       string
	Description string
}

func Test(t *testing.T) {
	db, err := inyorm.New(sqlite.Open(tmpSqliteFilePath))
	if err != nil {
		panic(err)
	}
	defer End(db, deleteSqliteFile)

	t.Run("create_table", func(t *testing.T) {
		stmt := db.CreateTable("users", func(q inyorm.CreateTable, e inyorm.Expr) {
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
		stmt := db.Insert("users", func(q inyorm.InsertQuery, e inyorm.Expr) {
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
		stmt := db.Select("users", func(q inyorm.SelectQuery, e inyorm.Expr) {
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
		stmt := db.CreateTable("posts", func(q inyorm.CreateTable, e inyorm.Expr) {
			q.Int("id").PrimaryKey().AutoIncrement()
			q.Int("author_id")
			q.Text("title").Unique().Default("untitled")
			q.Text("description").Nullable()

			q.ForeignKey("author_id").To("id", "users").OnDel("cascade")
			q.Check(e.Col("description")).Not().Like("%someword%")
		})

		if err := stmt.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("insert_post", func(t *testing.T) {
		stmt := db.Insert("posts", func(q inyorm.InsertQuery, e inyorm.Expr) {
			q.Insert(Post{}).Ignore(e.Col("id"), e.Col("title"), e.Col("description"))
			q.Into(e.Table("posts"))
			q.Values(Post{AuthorID: 1})
		})

		if err := stmt.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get_post_with_relation", func(t *testing.T) {
		stmt := db.Select("users", func(q inyorm.SelectQuery, e inyorm.Expr) {
			q.Select(e.All(), e.Col("title", "posts"), e.Col("description", "posts"))
			q.From(e.Table("users"))
			q.Join(e.Table("posts")).On(e.Col("author_id", "posts")).Equal(e.Col("id"))
			q.Where(e.Col("id")).Equal(e.Param(1))
			q.Limit(1)
		})

		type Data struct {
			User
			Title       string
			Description *string
		}

		var data Data
		if err := stmt.Bind(&data).Run(); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(data, Data{
			User: User{
				ID:       1,
				Account:  "acc123",
				Password: "mysecret",
				Age:      21,
			},
			Title:       "untitled",
			Description: nil,
		}) {
			t.Fatal("mismatch binding result")
		}
	})

	t.Run("update_post", func(t *testing.T) {
		stmt := db.Update("posts", func(q inyorm.UpdateQuery, e inyorm.Expr) {
			q.Update(e.Col("description"))
			q.Into(e.Table("posts"))
			q.Where(e.Col("author_id")).Equal(e.Param(1))
			q.Values(Post{Description: "My first post!"})
		})

		if err := stmt.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("try_to_insert_invalid_post", func(t *testing.T) {
		stmt := db.Insert("posts", func(q inyorm.InsertQuery, e inyorm.Expr) {
			q.Insert(Post{}).Ignore(e.Col("id"))
			q.Into(e.Table("posts"))
			q.Values(Post{
				AuthorID:    1,
				Title:       "This is my second post",
				Description: "this post contains someword, an invalid word",
			})
		})

		if err := stmt.Run(); err == nil || !strings.Contains(err.Error(), "CHECK constraint failed") {
			t.Fatalf("expected CHECK constraint failed, got: %v", err)
		}
	})

	t.Run("insert_a_few_posts", func(t *testing.T) {
		posts := make([]Post, 12)
		for i := range posts {
			posts[i].AuthorID = 1
			posts[i].Title = fmt.Sprintf("Title %d", i)
			posts[i].Description = fmt.Sprintf("Desc %d", i)
		}

		stmt := db.Insert("posts", func(q inyorm.InsertQuery, e inyorm.Expr) {
			q.Insert(Post{}).Ignore(e.Col("id"))
			q.Into(e.Table("posts"))
			q.Values(posts)
		})

		if err := stmt.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get_post_count", func(t *testing.T) {
		stmt := db.Select("users", func(q inyorm.SelectQuery, e inyorm.Expr) {
			q.Select(e.Col("id", "posts").Count())
			q.From(e.Table("users"))
			q.Join(e.Table("posts")).On(e.Col("author_id", "posts")).Equal(e.Col("id"))
			q.Where(e.Col("id")).Equal(e.Param(1))
			q.Limit(1)
		})

		var count int
		if err := stmt.Bind(&count).Run(); err != nil {
			t.Fatal(err)
		}

		if count != 13 {
			t.Fatalf("unexpected post count. got: %d", count)
		}
	})

	t.Run("get_custom_summary", func(t *testing.T) {
		stmt := db.Select("users", func(q inyorm.SelectQuery, e inyorm.Expr) {
			summary := e.Concat("The user with ID: ", e.Col("id"), ", Has: ", e.Col("id", "posts").Count(), " Posts.")

			q.Select(summary)
			q.From(e.Table("users"))
			q.Join(e.Table("posts")).On(e.Col("author_id", "posts")).Equal(e.Col("id"))
			q.Where(e.Col("id")).Equal(e.Param(1))
			q.Limit(1)
		})

		var result string
		if err := stmt.Bind(&result).Run(); err != nil {
			t.Fatal(err)
		}

		if result != "The user with ID: 1, Has: 13 Posts." {
			t.Fatalf("unexpected concat result. got: %s", result)
		}
	})

	t.Run("obtain_all_posts", func(t *testing.T) {
		expect := make([]Post, 13)
		for i := range expect {
			expect[i].ID = i + 1
			expect[i].AuthorID = 1

			if i == 0 {
				expect[i].Title = "untitled"
				expect[i].Description = "My first post!"
				continue
			}

			expect[i].Title = fmt.Sprintf("Title %d", i-1)
			expect[i].Description = fmt.Sprintf("Desc %d", i-1)
		}

		stmt := db.Select("posts", func(q inyorm.SelectQuery, e inyorm.Expr) {
			q.Select(e.All())
			q.From(e.Table("posts"))
		})

		var posts []Post
		if err := stmt.Bind(&posts).Run(); err != nil {
			t.Fatal(err)
		}

		for i := range posts {
			if !reflect.DeepEqual(posts[i], expect[i]) {
				t.Fatalf("mismatch results.\nexpect: %v\ngot: %v", expect[i], posts[i])
			}
		}
	})

	t.Run("delete_all_posts", func(t *testing.T) {
		stmt := db.Delete("posts", func(q inyorm.DeleteQuery, e inyorm.Expr) {
			q.Delete()
			q.From(e.Table("posts"))
			q.Where(e.Col("author_id")).Equal(1)
		})

		if err := stmt.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("expect_zero_values", func(t *testing.T) {
		stmt := db.Select("posts", func(q inyorm.SelectQuery, e inyorm.Expr) {
			q.Select(e.All())
			q.From(e.Table("posts"))
		})

		var posts []Post
		if err := stmt.Bind(&posts).Run(); err != nil {
			t.Fatal(err)
		}

		if len(posts) > 0 {
			t.Fatalf("expect zero values. got: %d", len(posts))
		}
	})
}
