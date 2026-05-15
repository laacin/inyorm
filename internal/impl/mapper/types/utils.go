package types

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/laacin/inyorm/internal/ir/ddl"
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

func toSnake(v string) string {
	if v == "" {
		return ""
	}

	var b strings.Builder
	runes := []rune(v)

	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := runes[i-1]

				if unicode.IsLower(prev) ||
					(unicode.IsUpper(prev) &&
						i+1 < len(runes) &&
						unicode.IsLower(runes[i+1])) {
					b.WriteByte('_')
				}
			}

			b.WriteRune(unicode.ToLower(r))
			continue
		}

		b.WriteRune(r)
	}

	return b.String()
}

type (
	fieldSchema struct {
		Meta  ddl.ColumnMeta
		Index []int
	}
	StructSchema map[string]fieldSchema
)

func (s StructSchema) GetMeta(name string) (ddl.ColumnMeta, bool) {
	fi, ok := s[name]
	return fi.Meta, ok
}

func (s StructSchema) GetIndex(name string) ([]int, bool) {
	if fi, ok := s[name]; ok {
		return fi.Index, true
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
	meta := ddl.ParseTag(tag)
	if meta.Name == "" {
		meta.Name = toSnake(name)
	}

	s[meta.Name] = fieldSchema{Meta: meta, Index: idx}
}

func (s StructSchema) merge(other StructSchema, baseIndex []int) {
	for name, fi := range other {
		if _, ok := s[name]; !ok {
			fi.Index = append(baseIndex, fi.Index...)
			s[name] = fi
		}
	}
}
