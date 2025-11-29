package schema

import "reflect"

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
