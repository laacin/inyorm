package mapper

import (
	"reflect"
	"sync"
)

// -- Resolve

func resolveRfl(v any, tag string, acceptValue bool) (reflect.Value, reflectResult) {
	info := reflectInfo(reflect.TypeOf(v), tag)
	if info.err != nil {
		return reflect.Value{}, info
	}

	if !info.ptr && !acceptValue {
		info.err = ErrPtrExpected
		return reflect.Value{}, info
	}

	val := reflect.ValueOf(v)
	if info.ptr {
		val = val.Elem()
	}

	if info.slc {
		return val, info
	}

	return val, info
}

// ---- Reflection info

type reflectResult struct {
	err   error
	slc   bool
	ptr   bool
	index map[string][]int
}

var cacheInfo sync.Map

func reflectInfo(t reflect.Type, tag string) reflectResult {
	if cached, ok := cacheInfo.Load(t); ok {
		return cached.(reflectResult)
	}

	info := reflectResult{}

	tmp := t
	if t.Kind() == reflect.Pointer {
		info.ptr = true
		tmp = tmp.Elem()
	}

	knd := tmp.Kind()
	if knd == reflect.Struct {
		idx := indexField(tag, tmp)
		info.index = mapIndexField(idx)
		cacheInfo.Store(t, info)
		return info
	}

	if knd == reflect.Slice {
		tmp = tmp.Elem()
		if tmp.Kind() == reflect.Struct {
			info.slc = true
			idx := indexField(tag, tmp)
			info.index = mapIndexField(idx)
			cacheInfo.Store(t, info)
			return info
		}
	}

	info.err = ErrUnexpectedType
	cacheInfo.Store(t, info)
	return info
}

// ---- Tag

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

func mapIndexField(idx []fieldInfo) map[string][]int {
	m := make(map[string][]int, len(idx))
	for _, v := range idx {
		m[v.name] = v.index
	}
	return m
}
