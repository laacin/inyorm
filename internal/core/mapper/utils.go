package mapper

import (
	"reflect"
	"slices"
)

func derefPtrTyp(t reflect.Type) (reflect.Type, int) {
	count := 0
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
		count++
	}
	return t, count
}

func derefPtrVal(v reflect.Value) (reflect.Value, int) {
	count := 0
	for v.IsValid() && v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
		count++
	}
	return v, count
}

type collector map[string]struct{}

func (c collector) Add(col string) {
	if col != "" {
		c[col] = struct{}{}
	}
}

func (c collector) ToSlice() []string {
	out := make([]string, 0, len(c))

	for col := range c {
		out = append(out, col)
	}

	slices.Sort(out)
	return out
}
