package test

import (
	"context"
	"fmt"
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

type Profile struct {
	ID        int
	Firstname string
	Lastname  string
	Bio       string
}

func Test(t *testing.T) {
	db, err := inyorm.New(sqlite.Open(tmpSqliteFilePath))
	if err != nil {
		panic(err)
	}
	defer End(db, deleteSqliteFile)

	t.Run("create_table", func(t *testing.T) {
		stmt := db.CreateTable(func(q inyorm.CreateTable, e inyorm.Expr) {
			q.TableName("users")

			q.Int("id").PrimaryKey().AutoIncrement()
			q.String("account").Unique()
			q.String("password")
			q.Int("age").Nullable()
		})

		if err := stmt.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("insert_values", func(t *testing.T) {
		stmt := db.Insert(func(q inyorm.InsertQuery, e inyorm.Expr) {
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
		stmt := db.Select(func(q inyorm.SelectQuery, e inyorm.Expr) {
			q.Select(e.All())
			q.From(e.Table("users"))
			q.Where(e.Col("id")).Equal(e.Param(1))
			q.Limit(1)
		})

		var user User
		if err := stmt.Bind(&user).Run(); err != nil {
			t.Fatal(err)
		}

		AssertEqual(t, user, User{
			ID:       1,
			Account:  "acc123",
			Password: "mysecret",
			Age:      21,
		})
	})

	t.Run("create_table_with_constraints", func(t *testing.T) {
		stmt := db.CreateTable(func(q inyorm.CreateTable, e inyorm.Expr) {
			q.TableName("posts")

			q.Int("id").PrimaryKey().AutoIncrement()
			q.Int("author_id")
			q.String("title").Default("untitled")
			q.String("description").Nullable()

			q.ForeignKey("author_id").To("id", "users").OnDel("cascade")
			q.Check(e.Col("description")).Not().Like("%someword%")
		})

		if err := stmt.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("insert_post", func(t *testing.T) {
		stmt := db.Insert(func(q inyorm.InsertQuery, e inyorm.Expr) {
			q.Insert(Post{}).Ignore(e.Col("id"), e.Col("title"), e.Col("description"))
			q.Into(e.Table("posts"))
			q.Values(Post{AuthorID: 1})
		})

		if err := stmt.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get_post_with_relation", func(t *testing.T) {
		stmt := db.Select(func(q inyorm.SelectQuery, e inyorm.Expr) {
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

		AssertEqual(t, data, Data{
			User: User{
				ID:       1,
				Account:  "acc123",
				Password: "mysecret",
				Age:      21,
			},
			Title:       "untitled",
			Description: nil,
		})
	})

	t.Run("update_post", func(t *testing.T) {
		stmt := db.Update(func(q inyorm.UpdateQuery, e inyorm.Expr) {
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
		stmt := db.Insert(func(q inyorm.InsertQuery, e inyorm.Expr) {
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

		stmt := db.Insert(func(q inyorm.InsertQuery, e inyorm.Expr) {
			q.Insert(Post{}).Ignore(e.Col("id"))
			q.Into(e.Table("posts"))
			q.Values(posts)
		})

		if err := stmt.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get_post_count", func(t *testing.T) {
		stmt := db.Select(func(q inyorm.SelectQuery, e inyorm.Expr) {
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
		stmt := db.Select(func(q inyorm.SelectQuery, e inyorm.Expr) {
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

		stmt := db.Select(func(q inyorm.SelectQuery, e inyorm.Expr) {
			q.Select(e.All())
			q.From(e.Table("posts"))
		})

		var posts []Post
		if err := stmt.Bind(&posts).Run(); err != nil {
			t.Fatal(err)
		}

		for i := range posts {
			AssertEqual(t, posts[i], expect[i])
		}
	})

	t.Run("delete_all_posts", func(t *testing.T) {
		stmt := db.Delete(func(q inyorm.DeleteQuery, e inyorm.Expr) {
			q.Delete()
			q.From(e.Table("posts"))
			q.Where(e.Col("author_id")).Equal(1)
		})

		if err := stmt.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("expect_zero_values", func(t *testing.T) {
		stmt := db.Select(func(q inyorm.SelectQuery, e inyorm.Expr) {
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

	t.Run("create_table_profile", func(t *testing.T) {
		stmt := db.CreateTable(func(q inyorm.CreateTable, e inyorm.Expr) {
			q.TableName("profiles")

			q.Int("id").PrimaryKey()
			q.String("firstname")
			q.String("lastname")
			q.String("bio").Nullable()

			q.ForeignKey("id").To("id", "users").OnDel("cascade")
		})

		if err := stmt.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("run_tx", func(t *testing.T) {
		stmt1 := db.Insert(func(q inyorm.InsertQuery, e inyorm.Expr) {
			q.Insert(User{})
			q.Into(e.Table("users"))
			q.Values(User{
				ID:       2,
				Account:  "myacc321",
				Password: "myPassword321",
				Age:      32,
			})
		})

		stmt2 := db.Insert(func(q inyorm.InsertQuery, e inyorm.Expr) {
			q.Insert(Profile{})
			q.Into(e.Table("profiles"))
			q.Values(Profile{
				ID:        2,
				Firstname: "john",
				Lastname:  "doe",
				Bio:       "this is my bio",
			})
		})

		if err := db.RunTx(context.Background(), stmt1, stmt2); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get_tx_values", func(t *testing.T) {
		stmt := db.Select(func(q inyorm.SelectQuery, e inyorm.Expr) {
			q.Select(e.All(), e.All("profiles"))
			q.From(e.Table("users"))
			q.Join(e.Table("profiles")).On(e.Col("id", "profiles")).Equal(e.Col("id"))
			q.Where(e.Col("id")).Equal(e.Param(2))
			q.Limit(1)
		})

		type ProfileNoID struct {
			Firstname string
			Lastname  string
			Bio       string
		}

		type Data struct {
			User
			ProfileNoID // BUG: <- column name conflicts cause only one embedded struct to bind correctly
		}
		var data Data
		if err := stmt.Bind(&data).Run(); err != nil {
			t.Fatal(err)
		}

		AssertEqual(t, data, Data{
			User: User{
				ID:       2,
				Account:  "myacc321",
				Password: "myPassword321",
				Age:      32,
			},
			ProfileNoID: ProfileNoID{
				Firstname: "john",
				Lastname:  "doe",
				Bio:       "this is my bio",
			},
		})
	})

	t.Run("bind_between_tx", func(t *testing.T) {
		var user1, user2 User

		stmt1 := db.Select(func(q inyorm.SelectQuery, e inyorm.Expr) {
			q.Select(e.All())
			q.From(e.Table("users"))
			q.Where(e.Col("id")).Equal(e.Param(1))
			q.Limit(1)
		}).Bind(&user1)

		stmt2 := db.Select(func(q inyorm.SelectQuery, e inyorm.Expr) {
			q.Select(e.All())
			q.From(e.Table("users"))
			q.Where(e.Col("id")).Equal(e.Param(2))
			q.Limit(1)
		}).Bind(&user2)

		if err := db.RunTx(context.Background(), stmt1, stmt2); err != nil {
			t.Fatal(err)
		}

		AssertEqual(t, user1, User{
			ID:       1,
			Account:  "acc123",
			Password: "mysecret",
			Age:      21,
		})
		AssertEqual(t, user2, User{
			ID:       2,
			Account:  "myacc321",
			Password: "myPassword321",
			Age:      32,
		})
	})
}
