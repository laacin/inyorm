package mapper_test

import (
	"reflect"
	"testing"

	"github.com/laacin/inyorm/internal/core/mapper"
)

func runTagTest(t *testing.T, v any, exp mapper.TagResult) {
	field := reflect.TypeOf(v).FieldByIndex([]int{0})
	result := mapper.ParseTag(field.Name, field.Tag.Get(mapper.TAG))

	if !reflect.DeepEqual(result, exp) {
		t.Fatalf("\nmismatch.\nExpect:\n%v\nHave:\n%v\n", exp, result)
	}
}

func TestParser(t *testing.T) {
	t.Run("no_tag", func(t *testing.T) {
		v := struct{ HashedPassword string }{}
		runTagTest(t, v, mapper.TagResult{Name: "hashed_password"})
	})

	t.Run("no_tag_with_upper_word", func(t *testing.T) {
		v := struct{ SQLName string }{}
		exp := mapper.TagResult{Name: "sql_name"}

		v2 := struct{ NameSQL string }{}
		exp2 := mapper.TagResult{Name: "name_sql"}

		runTagTest(t, v, exp)
		runTagTest(t, v2, exp2)
	})

	t.Run("weird_field_name", func(t *testing.T) {
		v := struct{ SQLNameIDFrom string }{}
		runTagTest(t, v, mapper.TagResult{Name: "sql_name_id_from"})
	})

	t.Run("keep_field_name", func(t *testing.T) {
		v := struct {
			SQLName string `inyorm:"col"`
		}{}

		runTagTest(t, v, mapper.TagResult{Name: "SQLName"})
	})

	t.Run("everything", func(t *testing.T) {
		v := struct {
			Account string `inyorm:"ignore, col:acc"`
		}{}

		runTagTest(t, v, mapper.TagResult{
			Name:   "acc",
			Ignore: true,
		})
	})

}
