package mapper

import (
	"reflect"
	"sync"

	"github.com/laacin/inyorm/internal/core"
)

type Mapper struct{}

func New() *Mapper {
	return &Mapper{}
}

var (
	kindCache  sync.Map
	customRefl = reflect.TypeFor[core.CustomKind]()
)

func (*Mapper) ReadKind(v any) core.KindInfo {
	typ := reflect.TypeOf(v)

	if cached, ok := kindCache.Load(typ); ok {
		return cached.(core.KindInfo)
	}

	info := core.KindInfo{}
	t := typ

	if t.Kind() == reflect.Pointer {
		t, _ = derefPtrTyp(typ)
		info.Ptr = true
	}

	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		t, _ = derefPtrTyp(t.Elem())
		info.Slice = true
	}

	info.Kind = whichKind(t)
	if info.Kind == core.KindStruct {
		info.Schema = readStruct(t)
	}

	kindCache.Store(typ, info)
	return info
}

// helpers

func whichKind(t reflect.Type) core.Kind {
	if t.Implements(customRefl) || reflect.PointerTo(t).Implements(customRefl) {
		return core.KindCustom
	}

	switch t.Kind() {
	case reflect.Struct:
		return core.KindStruct

	case reflect.Map:
		return core.KindMap

	case reflect.String:
		return core.KindString

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return core.KindInt

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return core.KindUint

	case reflect.Float32, reflect.Float64:
		return core.KindFloat

	case reflect.Bool:
		return core.KindBool

	case reflect.Interface:
		return core.KindAny

	default:
		return core.KindUnknown
	}
}

// --- Struct

func readStruct(t reflect.Type) Schema {
	schema := Schema(map[string]core.FieldInfo{})

	for field := range t.Fields() {
		typ := field.Type
		if typ.Kind() == reflect.Pointer {
			typ = typ.Elem()
		}

		if whichKind(typ) == core.KindUnknown {
			continue
		}

		if field.Anonymous && typ.Kind() == reflect.Struct {
			baseIndex := append([]int(nil), field.Index...)
			schema.Merge(baseIndex, readStruct(typ))
			continue
		}

		schema.Add(field.Name, field.Tag.Get(TAG), field.Index)
	}
	return schema
}

type Schema map[string]core.FieldInfo

func (s Schema) Add(fieldName, tag string, idx []int) {
	tagResult := ParseTag(fieldName, tag)

	s[tagResult.Name] = core.FieldInfo{
		Ignore: tagResult.Ignore,
		Index:  idx,
	}
}

func (s Schema) Merge(baseIndex []int, other Schema) {
	for name, info := range other {
		idx := append(baseIndex, info.Index...)

		s[name] = core.FieldInfo{
			Index:  idx,
			Ignore: info.Ignore,
		}
	}
}

func (s Schema) Get(name string) (core.FieldInfo, bool) {
	r, ok := s[name]
	return r, ok
}

func (s Schema) GetIndex(name string) ([]int, bool) {
	r, ok := s[name]
	return r.Index, ok
}

func (s Schema) IterNames(fn func(string)) {
	for name := range s {
		fn(name)
	}
}
