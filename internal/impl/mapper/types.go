package mapper

import (
	"reflect"

	"github.com/laacin/inyorm/internal/entity/api"
)

// --- Valid app types
type Kind int

const (
	// Native types
	KindString Kind = iota
	KindInt
	KindUint
	KindFloat
	KindBool
	KindStruct
	KindMap

	// App types
	KindColumn

	// Fallback
	KindAny
)

var column = reflect.TypeFor[api.Column]()

type TypeInfo struct {
	Kind Kind
	Ptr  int
	Slc  int
}

type FieldInfo struct {
	Name  string
	Tag   string
	Index []int
}

func ObtainInfo(t reflect.Type) TypeInfo {
	info := TypeInfo{}

	t, c := derefPtrTyp(t)
	info.Ptr = c

	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t, c = derefPtrTyp(t.Elem())
		info.Slc = c
	} else {
		info.Slc = -1
	}

	if t.Implements(column) {
		info.Kind = KindColumn
		return info
	}

	switch t.Kind() {
	case reflect.Struct:
		info.Kind = KindStruct

	case reflect.Map:
		info.Kind = KindMap

	case reflect.String:
		info.Kind = KindString

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		info.Kind = KindInt

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		info.Kind = KindUint

	case reflect.Float32, reflect.Float64:
		info.Kind = KindFloat

	case reflect.Bool:
		info.Kind = KindBool

	default:
		info.Kind = KindAny
	}

	return info
}

// ---- Helpers
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
