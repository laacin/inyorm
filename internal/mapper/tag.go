package mapper

import (
	"reflect"
	"sync"
)

type fieldInfo struct {
	name  string
	index []int
}

var fieldCache sync.Map

func cachedFields(columnTag string, typ reflect.Type) []fieldInfo {
	if cached, ok := fieldCache.Load(typ); ok {
		return cached.([]fieldInfo)
	}

	var fields []fieldInfo
	for i := range typ.NumField() {
		f := typ.Field(i)

		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			nested := cachedFields(columnTag, f.Type)
			for _, nf := range nested {
				nf.index = append(nf.index, nf.index...)
				fields = append(fields, nf)
			}
			continue
		}

		tagVal := f.Tag.Get(columnTag)
		if tagVal == "" || tagVal == "-" {
			continue
		}
		fields = append(fields, fieldInfo{tagVal, f.Index})
	}

	fieldCache.Store(typ, fields)
	return fields
}
