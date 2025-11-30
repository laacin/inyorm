package mapper

import (
	"reflect"
	"sync"
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
