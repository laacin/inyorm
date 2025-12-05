package mapper_test

import (
	"reflect"
	"testing"

	"github.com/laacin/inyorm"
	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/mapper"
)

func run(t *testing.T, v []any, exp []string) {
	cols, err := mapper.GetColumns("col", v)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(cols, exp) {
		t.Errorf("columns doesn't match:\n - Expect: %v\n - Have: %v\n", exp, cols)
	}
}

func col(name string) inyorm.Column {
	return &column.Column[inyorm.Column]{BaseName: name}
}

func TestGetColumn(t *testing.T) {
	type User struct {
		ID      string `col:"id"`
		Account string `col:"account"`
	}

	t.Run("simple_struct", func(t *testing.T) {
		v := []any{User{}}
		exp := []string{"account", "id"}

		run(t, v, exp)
	})

	t.Run("embedded_structs", func(t *testing.T) {
		type Nested struct {
			User
			Name string `col:"name"`
			Age  string `col:"age"`
		}

		v := []any{Nested{}}
		exp := []string{"account", "age", "id", "name"}

		run(t, v, exp)
	})

	t.Run("many_structs", func(t *testing.T) {
		type (
			Str1 struct {
				Age int `col:"age"`
			}
			Str2 struct {
				Name string `col:"name"`
			}
			Str3 struct {
				Surname string `col:"surname"`
			}
		)
		v := []any{User{}, Str1{}, Str2{}, Str3{}}
		exp := []string{"account", "age", "id", "name", "surname"}

		run(t, v, exp)
	})

	t.Run("single_column", func(t *testing.T) {
		acc := col("account")

		v := []any{acc}
		exp := []string{"account"}

		run(t, v, exp)
	})

	t.Run("many_columns", func(t *testing.T) {
		var (
			acc = col("account")
			age = col("age")
			id  = col("id")
		)

		v := []any{acc, age, id}
		exp := []string{"account", "age", "id"}

		run(t, v, exp)
	})

	t.Run("single_map", func(t *testing.T) {
		m := map[string]any{
			"age":     123,
			"account": "asd",
		}

		v := []any{m}
		exp := []string{"account", "age"}

		run(t, v, exp)
	})

	t.Run("many_maps", func(t *testing.T) {
		m := []map[string]any{
			{"age": 123, "account": "asd"},
			{"name": "mary", "id": "uuid"},
		}

		v := []any{m[0], m[1]}
		exp := []string{"account", "age", "id", "name"}

		run(t, v, exp)
	})

	t.Run("mix", func(t *testing.T) {
		var (
			fname = col("firstname")
			lname = col("lastname")
		)

		m := map[string]any{"banned": false}

		v := []any{User{}, fname, m, lname}
		exp := []string{"account", "banned", "firstname", "id", "lastname"}

		run(t, v, exp)
	})

	t.Run("mix_with_ptr", func(t *testing.T) {
		var (
			fname = col("firstname")
			lname = col("lastname")
		)

		m := map[string]any{"banned": false}

		v := []any{&User{}, fname, &m, lname}
		exp := []string{"account", "banned", "firstname", "id", "lastname"}

		run(t, v, exp)
	})
}
