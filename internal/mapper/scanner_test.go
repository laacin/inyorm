package mapper_test

import (
	"reflect"
	"testing"

	"github.com/laacin/inyorm/internal/mapper"
)

type MockRows struct {
	Cols []string
	Data [][]any
	i    int
}

func (m *MockRows) Columns() ([]string, error) {
	return m.Cols, nil
}

func (m *MockRows) Next() bool {
	if m.i >= len(m.Data) {
		return false
	}
	m.i++
	return true
}

func (m *MockRows) Scan(dest ...any) error {
	row := m.Data[m.i-1]
	for idx := range dest {
		if idx < len(row) {
			reflect.ValueOf(dest[idx]).Elem().Set(reflect.ValueOf(row[idx]))
		}
	}
	return nil
}

func newScanner(t *testing.T, rows *MockRows, binder any) func(expect any, expErr ...error) {
	err := mapper.Scan(rows, "inyorm", binder)

	return func(expect any, expErr ...error) {
		if err != nil && len(expErr) > 0 && expErr[0] != nil {
			if err.Error() != expErr[0].Error() {
				t.Fatalf("\nmismatch error:\nExpect:\n%s\nHave:\n%s\n", expErr[0].Error(), err.Error())
			}
			return

		} else if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(expect, binder) {
			t.Fatalf("\nmismatch result:\nExpect:\n%#v\nHave:\n%#v\n", expect, binder)
		}
	}
}

func TestScanner(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		type User struct {
			ID   int    `inyorm:"id"`
			Name string `inyorm:"name"`
		}

		rows := &MockRows{
			Cols: []string{"id", "name"},
			Data: [][]any{
				{1, "lautaro"},
			},
		}

		exp := User{
			ID:   1,
			Name: "lautaro",
		}

		var u User
		newScanner(t, rows, &u)(&exp)
	})

	t.Run("pointer_element", func(t *testing.T) {
		type User struct {
			ID   int    `inyorm:"id"`
			Name string `inyorm:"name"`
		}

		rows := &MockRows{
			Cols: []string{"id", "name"},
			Data: [][]any{
				{10, "ana"},
			},
		}

		exp := &User{
			ID:   10,
			Name: "ana",
		}

		var u User
		newScanner(t, rows, &u)(exp)
	})

	t.Run("missing_column", func(t *testing.T) {
		type User struct {
			ID   int    `inyorm:"id"`
			Name string `inyorm:"name"`
			Age  int    `inyorm:"age"`
		}

		rows := &MockRows{
			Cols: []string{"id", "name"},
			Data: [][]any{
				{5, "juan"},
			},
		}

		exp := User{
			ID:   5,
			Name: "juan",
			Age:  0,
		}

		var u User
		newScanner(t, rows, &u)(&exp)
	})

	t.Run("extra_column", func(t *testing.T) {
		type User struct {
			ID   int    `inyorm:"id"`
			Name string `inyorm:"name"`
		}

		rows := &MockRows{
			Cols: []string{"id", "name", "ignored"},
			Data: [][]any{
				{3, "maria", "x"},
			},
		}

		exp := User{
			ID:   3,
			Name: "maria",
		}

		var u User
		newScanner(t, rows, &u)(&exp)
	})

	t.Run("slice_structs", func(t *testing.T) {
		type User struct {
			ID   int    `inyorm:"id"`
			Name string `inyorm:"name"`
		}

		rows := &MockRows{
			Cols: []string{"id", "name"},
			Data: [][]any{
				{1, "a"},
				{2, "b"},
				{3, "c"},
			},
		}

		exp := []User{
			{ID: 1, Name: "a"},
			{ID: 2, Name: "b"},
			{ID: 3, Name: "c"},
		}

		var u []User
		newScanner(t, rows, &u)(&exp)
	})

	t.Run("array_structs", func(t *testing.T) {
		type User struct {
			ID   int    `inyorm:"id"`
			Name string `inyorm:"name"`
		}

		rows := &MockRows{
			Cols: []string{"id", "name"},
			Data: [][]any{
				{1, "a"},
				{2, "b"},
				{3, "c"},
			},
		}

		exp := [...]User{
			{ID: 1, Name: "a"},
			{ID: 2, Name: "b"},
			{ID: 3, Name: "c"},
		}

		var u [3]User
		newScanner(t, rows, &u)(exp, mapper.ErrUnexpectedType)
	})

	t.Run("slice_ptr_structs", func(t *testing.T) {
		type User struct {
			ID   int    `inyorm:"id"`
			Name string `inyorm:"name"`
		}

		rows := &MockRows{
			Cols: []string{"id", "name"},
			Data: [][]any{
				{7, "x"},
				{8, "y"},
			},
		}

		exp := []*User{
			{ID: 7, Name: "x"},
			{ID: 8, Name: "y"},
		}

		var u []*User
		newScanner(t, rows, &u)(&exp, mapper.ErrUnexpectedType)
	})

	t.Run("slice_preallocated", func(t *testing.T) {
		type User struct {
			ID   int    `inyorm:"id"`
			Name string `inyorm:"name"`
		}

		rows := &MockRows{
			Cols: []string{"id", "name"},
			Data: [][]any{
				{1, "a"},
				{2, "b"},
			},
		}

		u := make([]User, 1)
		u[0] = User{ID: 0, Name: "zero"}

		exp := []User{
			{ID: 1, Name: "a"},
			{ID: 2, Name: "b"},
		}

		newScanner(t, rows, &u)(&exp)
	})

	t.Run("struct_preallocated", func(t *testing.T) {
		type User struct {
			ID   int    `inyorm:"id"`
			Name string `inyorm:"name"`
		}

		rows := &MockRows{
			Cols: []string{"id", "name"},
			Data: [][]any{
				{100, "zzz"},
			},
		}

		u := User{ID: 1, Name: "old"}
		exp := User{ID: 100, Name: "zzz"}

		newScanner(t, rows, &u)(&exp)
	})

	t.Run("slice_less_capacity_than_rows", func(t *testing.T) {
		type User struct {
			ID   int    `inyorm:"id"`
			Name string `inyorm:"name"`
		}

		rows := &MockRows{
			Cols: []string{"id", "name"},
			Data: [][]any{
				{1, "a"},
				{2, "b"},
				{3, "c"},
			},
		}

		u := make([]User, 1)
		exp := []User{
			{ID: 1, Name: "a"},
			{ID: 2, Name: "b"},
			{ID: 3, Name: "c"},
		}

		newScanner(t, rows, &u)(&exp)
	})

	t.Run("slice_greater_capacity_than_rows", func(t *testing.T) {
		type User struct {
			ID   int    `inyorm:"id"`
			Name string `inyorm:"name"`
		}

		rows := &MockRows{
			Cols: []string{"id", "name"},
			Data: [][]any{
				{9, "x"},
			},
		}

		u := make([]User, 3)
		u[0] = User{ID: 1, Name: "a"}
		u[1] = User{ID: 2, Name: "b"}
		u[2] = User{ID: 3, Name: "c"}

		exp := []User{
			{ID: 9, Name: "x"},
		}
		newScanner(t, rows, &u)(&exp)
	})

	t.Run("single_map", func(t *testing.T) {
		rows := &MockRows{
			Cols: []string{"id", "name"},
			Data: [][]any{
				{9, "x"},
			},
		}

		m := make(map[string]any)
		exp := map[string]any{
			"id":   9,
			"name": "x",
		}
		newScanner(t, rows, m)(exp)
	})

	t.Run("many_maps", func(t *testing.T) {
		rows := &MockRows{
			Cols: []string{"id", "name", "age"},
			Data: [][]any{
				{1, "alice", 20},
				{2, "bob", 31},
				{3, "carol", 45},
				{4, "dan", 18},
				{5, "eva", 27},
				{6, "mike", 52},
			},
		}

		m := make([]map[string]any, 0)
		exp := []map[string]any{
			{"id": 1, "name": "alice", "age": 20},
			{"id": 2, "name": "bob", "age": 31},
			{"id": 3, "name": "carol", "age": 45},
			{"id": 4, "name": "dan", "age": 18},
			{"id": 5, "name": "eva", "age": 27},
			{"id": 6, "name": "mike", "age": 52},
		}

		newScanner(t, rows, &m)(&exp)
	})
}
