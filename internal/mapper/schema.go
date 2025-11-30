package mapper

import (
	"reflect"
	"sync"

	"github.com/laacin/inyorm/internal/core"
)

type schema struct {
	Type  int
	Slc   bool
	Ptr   bool
	Index []fieldInfo
}

func (s *schema) IndexMap() map[string][]int {
	m := make(map[string][]int, len(s.Index))
	for _, info := range s.Index {
		m[info.name] = info.index
	}
	return m
}

var cache sync.Map

func getSchema(tag string, v any) (*schema, error) {
	elemType := reflect.TypeOf(v)

	if cached, ok := cache.Load(elemType); ok {
		return cached.(*schema), nil
	}

	var (
		t     = elemType
		slc   bool
		ptr   bool
		index []fieldInfo
	)

	if t.Kind() == reflect.Pointer {
		ptr = true
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice {
		slc = true
		t = t.Elem()
	}

	typ := whichIs(t)
	if typ == typeStruct {
		index = indexField(tag, t)
	}

	schm := &schema{
		Type:  typ,
		Ptr:   ptr,
		Slc:   slc,
		Index: index,
	}

	cache.Store(elemType, schm)
	return schm, nil
}

// -- which type is it

const (
	typeUnknown = iota
	typeString
	typeInt
	typeUint
	typeFloat
	typeBool
	typeMap
	typeColumn
	typeStruct
	typeAny
)

var ColumnIface = reflect.TypeOf((*core.Column)(nil)).Elem()

func whichIs(t reflect.Type) int {
	if t.Implements(ColumnIface) {
		return typeColumn
	}

	knd := t.Kind()

	if knd == reflect.Struct {
		return typeStruct
	}

	if knd == reflect.Interface {
		return typeAny
	}

	if knd == reflect.Map {
		return typeMap
	}

	if knd == reflect.String {
		return typeString
	}

	if knd == reflect.Bool {
		return typeBool
	}

	if knd >= reflect.Int && knd <= reflect.Int64 {
		return typeInt
	}

	if knd >= reflect.Uint && knd <= reflect.Uint64 {
		return typeUint
	}

	if knd == reflect.Float32 || knd == reflect.Float64 {
		return typeFloat
	}

	return typeUnknown
}

// -- index fields for structs

type fieldInfo struct {
	name  string
	index []int
}

func indexField(tag string, t reflect.Type) []fieldInfo {
	var fields []fieldInfo
	for i := range t.NumField() {
		f := t.Field(i)

		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			nested := indexField(tag, f.Type)
			for _, nf := range nested {
				nf.index = append(f.Index, nf.index...)
				fields = append(fields, nf)
			}
		}

		name := f.Tag.Get(tag)
		if name == "" || name == "-" {
			continue
		}

		fields = append(fields, fieldInfo{
			name:  name,
			index: f.Index,
		})
	}

	return fields
}
