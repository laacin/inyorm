package mapper

import (
	"reflect"

	"github.com/laacin/inyorm/internal/core"
)

func (m *Mapper) ReadCols(entries ...any) []string {
	c := collector(map[string]struct{}{})

	for _, entry := range entries {
		info := m.ReadKind(entry)

		switch info.Kind {
		case core.KindStruct:
			Schema(info.Schema).IterNames(func(s string) {
				c.Add(s)
			})

		case core.KindMap:
			colsFrom(entry, func(s map[string]any) {
				for k := range s {
					c.Add(k)
				}
			})

		case core.KindString:
			colsFrom(entry, func(s string) {
				c.Add(s)
			})

		case core.KindCustom:
			colsFrom(entry, func(s core.CustomKind) {
				c.Add(s.BaseName())
			})
		}
	}

	return c.ToSlice()
}

// --- helpers

func colsFrom[T any](v any, store func(T)) {
	switch val := v.(type) {
	case T:
		store(val)

	case *T:
		if val != nil {
			store(*val)
		}

	case []T:
		for _, item := range val {
			store(item)
		}

	case *[]T:
		if val != nil {
			for _, item := range *val {
				store(item)
			}
		}

	default:
		colsFromRefl(v, store)
	}
}

func colsFromRefl[T any](v any, store func(T)) {
	val, _ := derefPtrVal(reflect.ValueOf(v))

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := range val.Len() {
			elem, _ := derefPtrVal(val.Index(i))
			if assert, ok := reflect.TypeAssert[T](elem); ok {
				store(assert)
			}
		}
		return
	}

	if assert, ok := reflect.TypeAssert[T](val); ok {
		store(assert)
	}
}
