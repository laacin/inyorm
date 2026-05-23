package mapper_test

import (
	"reflect"
	"testing"

	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/impl/mapper"
)

func run(t *testing.T, v []any, exp []string) {
	cols := mapper.ReadColumns(v)
	if !reflect.DeepEqual(cols, exp) {
		t.Errorf("columns doesn't match:\n - Expect: %v\n - Have: %v\n", exp, cols)
	}
}

func col(name string) *expr.Col {
	col := &expr.Col{}
	return col.Start(name, "")
}

func TestGetColumn(t *testing.T) {
	type User struct {
		ID      string
		Account string
	}

	t.Run("simple_struct", func(t *testing.T) {
		v := []any{User{}}
		exp := []string{"account", "id"}

		run(t, v, exp)
	})

	t.Run("embedded_structs", func(t *testing.T) {
		type Nested struct {
			User
			Name string
			Age  string
		}

		v := []any{Nested{}}
		exp := []string{"account", "age", "id", "name"}

		run(t, v, exp)
	})

	t.Run("many_structs", func(t *testing.T) {
		type (
			Str1 struct {
				Age int
			}
			Str2 struct {
				Name string
			}
			Str3 struct {
				Surname string
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

	t.Run("duplicated_columns", func(t *testing.T) {
		m := map[string]any{
			"id": 1,
		}

		v := []any{
			User{},
			m,
			col("yet"),
		}

		exp := []string{"account", "id", "yet"}

		run(t, v, exp)
	})

	t.Run("slice_of_columns", func(t *testing.T) {
		v := []any{[]*expr.Col{col("name"), col("age")}}

		exp := []string{"age", "name"}

		run(t, v, exp)
	})

	t.Run("ptr_slice_of_columns", func(t *testing.T) {
		slc := []*expr.Col{col("name"), col("age")}

		v := []any{&slc}
		exp := []string{"age", "name"}

		run(t, v, exp)
	})

	t.Run("slice_of_maps", func(t *testing.T) {
		v := []any{
			[]map[string]any{
				{"name": "john"},
				{"age": 10},
			},
		}

		exp := []string{"age", "name"}

		run(t, v, exp)
	})

	t.Run("ptr_slice_of_maps", func(t *testing.T) {
		slc := []map[string]any{
			{"name": "john"},
			{"age": 10},
		}

		v := []any{&slc}
		exp := []string{"age", "name"}

		run(t, v, exp)
	})

	t.Run("slice_of_ptr_maps", func(t *testing.T) {
		m1 := map[string]any{"name": "john"}
		m2 := map[string]any{"age": 10}

		v := []any{
			[]*map[string]any{&m1, &m2},
		}

		exp := []string{"age", "name"}

		run(t, v, exp)
	})

	t.Run("nil_ptr_map_inside_slice", func(t *testing.T) {
		m := map[string]any{"name": "john"}

		v := []any{
			[]*map[string]any{&m, nil},
		}

		exp := []string{"name"}

		run(t, v, exp)
	})

	t.Run("slice_of_strings", func(t *testing.T) {
		v := []any{
			[]string{"name", "age"},
		}

		exp := []string{"age", "name"}

		run(t, v, exp)
	})

	t.Run("ptr_slice_of_strings", func(t *testing.T) {
		slc := []string{"name", "age"}

		v := []any{&slc}
		exp := []string{"age", "name"}

		run(t, v, exp)
	})

	t.Run("slice_of_ptr_strings", func(t *testing.T) {
		name := "name"
		age := "age"

		v := []any{
			[]*string{&name, &age},
		}

		exp := []string{"age", "name"}

		run(t, v, exp)
	})

	t.Run("nil_ptr_string_inside_slice", func(t *testing.T) {
		name := "name"

		v := []any{
			[]*string{&name, nil},
		}

		exp := []string{"name"}

		run(t, v, exp)
	})

	t.Run("name_in_tag", func(t *testing.T) {
		type Dummy struct {
			Account  string `inyorm:"col:acc"`
			Password string `inyorm:"col:pw"`
		}

		v := []any{Dummy{}}

		exp := []string{"acc", "pw"}
		run(t, v, exp)
	})
}
