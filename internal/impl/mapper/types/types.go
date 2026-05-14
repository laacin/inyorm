package types

import (
	"reflect"
	"sync"

	"github.com/laacin/inyorm/internal/api"
)

type TypeInfo struct {
	Kind Kind
	Ptr  int
	Slc  int
	Fis  []FieldInfo // nil if Kind != KindStruct
}

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

var (
	column    = reflect.TypeFor[api.Column]()
	infoCache sync.Map
)

func ReadInfo(typ reflect.Type) TypeInfo {
	if cache, ok := infoCache.Load(typ); ok {
		return cache.(TypeInfo)
	}

	info := TypeInfo{}

	t, c := derefPtrTyp(typ)
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
		info.Fi = readStruct(t)

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

	infoCache.Store(t, info)
	return info
}
