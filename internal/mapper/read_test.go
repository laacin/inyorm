package mapper_test

import (
	"reflect"
	"testing"

	"github.com/laacin/inyorm/internal/mapper"
)

func newTest(t *testing.T, v any) func(rows int, cols []string, vals []any) {
	r, columns, values, err := mapper.ReadValues("inyorm", v)
	if err != nil {
		t.Fatal(err)
	}
	return func(rows int, cols []string, vals []any) {
		if r != rows {
			t.Fatalf("\nmismatch result:\nExpect:\n%d\nHave:\n%d\n", rows, r)
		}
		if !reflect.DeepEqual(columns, cols) {
			t.Fatalf("\nmismatch result:\nExpect:\n%#v\nHave:\n%#v\n", cols, columns)
		}
		if !reflect.DeepEqual(values, vals) {
			t.Fatalf("\nmismatch result:\nExpect:\n%#v\nHave:\n%#v\n", vals, values)
		}
	}
}

type User struct {
	Account   string `inyorm:"account"`
	Age       int    `inyorm:"age"`
	Firstname string `inyorm:"firstname"`
	Lastname  string `inyorm:"lastname"`
}

type UserPartial struct {
	Account   string  `inyorm:"account"`
	Age       int     `inyorm:"age"`
	Firstname *string `inyorm:"firstname"`
	Lastname  *string `inyorm:"lastname"`
}

func TestRead(t *testing.T) {
	t.Run("one", func(t *testing.T) {
		u := &User{
			Account:   "acc",
			Age:       17,
			Firstname: "max",
			Lastname:  "porter",
		}
		cols := []string{"account", "age", "firstname", "lastname"}
		vals := []any{"acc", 17, "max", "porter"}

		run := newTest(t, u)
		run(1, cols, vals)
	})

	t.Run("many", func(t *testing.T) {
		users := []User{
			{
				Account:   "acc1",
				Age:       20,
				Firstname: "john",
				Lastname:  "doe",
			},
			{
				Account:   "acc2",
				Age:       30,
				Firstname: "jane",
				Lastname:  "smith",
			},
			{
				Account:   "acc3",
				Age:       41,
				Firstname: "mark",
				Lastname:  "stone",
			},
			{
				Account:   "acc4",
				Age:       28,
				Firstname: "lucas",
				Lastname:  "brown",
			},
			{
				Account:   "acc5",
				Age:       33,
				Firstname: "alice",
				Lastname:  "white",
			},
		}

		cols := []string{"account", "age", "firstname", "lastname"}
		vals := []any{
			"acc1", 20, "john", "doe",
			"acc2", 30, "jane", "smith",
			"acc3", 41, "mark", "stone",
			"acc4", 28, "lucas", "brown",
			"acc5", 33, "alice", "white",
		}

		run := newTest(t, users)
		run(5, cols, vals)
	})

	t.Run("with_nils", func(t *testing.T) {
		u := UserPartial{
			Account: "acc",
			Age:     21,
		}

		cols := []string{"account", "age", "firstname", "lastname"}
		vals := []any{"acc", 21, nil, nil}

		run := newTest(t, u)
		run(1, cols, vals)
	})
}
