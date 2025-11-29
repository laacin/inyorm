package schema

import (
	"errors"
	"reflect"
	"slices"
	"sync"
)

const (
	TypeUnknown SchemaType = iota
	TypeString
	TypeInt
	TypeUint
	TypeFloat
	TypeBool
	TypeMap
	TypeColumn
	TypeStruct
	TypeAny
)

type Schema struct {
	Type  SchemaType
	Slc   bool
	Ptr   bool
	Index []fieldInfo
}

func (s *Schema) IndexMap() map[string][]int {
	m := make(map[string][]int, len(s.Index))
	for _, info := range s.Index {
		m[info.name] = info.index
	}
	return m
}

var cache sync.Map

func GetSchema(tag string, v any) (*Schema, error) {
	elemType := reflect.TypeOf(v)

	if cached, ok := cache.Load(elemType); ok {
		return cached.(*Schema), nil
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
	if typ == TypeStruct {
		index = indexField(tag, t)
	}

	schm := &Schema{
		Type:  typ,
		Ptr:   ptr,
		Slc:   slc,
		Index: index,
	}

	cache.Store(elemType, schm)
	return schm, nil
}

func GetColumns(tag string, v any) ([]string, error) {
	s, err := GetSchema(tag, v)
	if err != nil {
		return nil, err
	}

	if s == nil {
		return nil, errors.New("missing schema")
	}

	var cols []string
	switch s.Type {
	case TypeColumn:
		cols = colsFromCol(v, s.Slc)

	case TypeString:
		cols = colsFromString(v, s.Slc, s.Ptr)

	case TypeStruct:
		cols = colsFromStruct(s.Index)

	case TypeMap:
		if s.Slc {
			return nil, errors.New("invalid schema type")
		}
		cols = colsFromMap(v, s.Ptr)

	default:
		return nil, errors.New("invalid schema type")
	}

	for _, col := range cols {
		if col == "" || col == "*" {
			return nil, errors.New("invalid column name")
		}
	}

	slices.Sort(cols)
	return cols, nil
}
