package types

import (
	"reflect"

	"github.com/laacin/inyorm/internal/core"
)

func Unwrap[T any](v any, ptr bool) T {
	if ptr {
		return *v.(*T)
	}
	return v.(T)
}

func UnwrapSlc[T any](v any, ptr bool) []T {
	if ptr {
		return *v.(*[]T)
	}
	return v.([]T)
}

func DerefPtrTyp(t reflect.Type) (reflect.Type, int) {
	count := 0
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
		count++
	}
	return t, count
}

func DerefPtrVal(v reflect.Value) (reflect.Value, int) {
	count := 0
	for v.IsValid() && v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
		count++
	}
	return v, count
}

type StructSchema map[string]core.FieldResult

func (s StructSchema) Get(name string) (core.FieldResult, bool) {
	r, ok := s[name]
	return r, ok
}

func (s StructSchema) GetIndex(name string) ([]int, bool) {
	if r, ok := s[name]; ok {
		return r.Index, true
	}
	return nil, false
}

func (s StructSchema) Len() int {
	return len(s)
}

func (s StructSchema) IterNames(fn func(string)) {
	for k := range s {
		fn(k)
	}
}

// internal methods
func (s StructSchema) add(name, tag string, idx []int) {
	r := core.ParseField(name, tag, idx)
	s[r.Name] = r
}

func (s StructSchema) merge(other StructSchema, baseIndex []int) {
	for name, r := range other {
		if _, ok := s[name]; !ok {
			r.Index = append(baseIndex, r.Index...)
			s[name] = r
		}
	}
}
