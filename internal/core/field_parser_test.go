package core_test

import (
	"reflect"
	"testing"

	"github.com/laacin/inyorm/internal/core"
)

func run(t *testing.T, v any, exp core.FieldSchema) {
	field := reflect.TypeOf(v).FieldByIndex([]int{0})
	result := core.NewFieldSchema(field.Name, field.Tag.Get(core.TAG), field.Index)
	exp.Index = []int{0}

	if !reflect.DeepEqual(result, exp) {
		t.Fatalf("\nmismatch.\nExpect:\n%v\nHave:\n%v\n", exp, result)
	}
}

func TestParser(t *testing.T) {
	t.Run("no_tag", func(t *testing.T) {
		v := struct{ HashedPassword string }{}
		run(t, v, core.FieldSchema{Name: "hashed_password"})
	})

	t.Run("no_tag_with_upper_word", func(t *testing.T) {
		v := struct{ SQLName string }{}
		exp := core.FieldSchema{Name: "sql_name"}

		v2 := struct{ NameSQL string }{}
		exp2 := core.FieldSchema{Name: "name_sql"}

		run(t, v, exp)
		run(t, v2, exp2)
	})

	t.Run("weird_field_name", func(t *testing.T) {
		v := struct{ SQLNameIDFrom string }{}
		run(t, v, core.FieldSchema{Name: "sql_name_id_from"})
	})

	t.Run("keep_field_name", func(t *testing.T) {
		v := struct {
			SQLName string `inyorm:"col"`
		}{}

		run(t, v, core.FieldSchema{Name: "SQLName"})
	})

	t.Run("everything", func(t *testing.T) {
		v := struct {
			Account string `inyorm:"skip, col:acc"`
		}{}

		run(t, v, core.FieldSchema{
			Name: "acc",
			Skip: true,
		})
	})

}
