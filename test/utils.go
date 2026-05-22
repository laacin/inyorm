package test

import (
	"reflect"
	"testing"
)

func AssertEqual[T any](t *testing.T, got, expected T) {
	t.Helper()

	if reflect.DeepEqual(expected, got) {
		return
	}

	exp := reflect.ValueOf(expected)
	g := reflect.ValueOf(got)

	if exp.Kind() == reflect.Struct && g.Kind() == reflect.Struct {
		typ := exp.Type()

		for i := 0; i < exp.NumField(); i++ {
			field := typ.Field(i)

			ev := exp.Field(i).Interface()
			gv := g.Field(i).Interface()

			if !reflect.DeepEqual(ev, gv) {
				t.Errorf(
					"field %s mismatch:\nexpected: %#v\ngot:      %#v",
					field.Name,
					ev,
					gv,
				)
			}
		}

		return
	}

	t.Errorf(
		"values are different:\nexpected: %#v\ngot:      %#v",
		expected,
		got,
	)
}
