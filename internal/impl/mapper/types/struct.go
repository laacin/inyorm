package types

import (
	"reflect"

	"github.com/laacin/inyorm/internal/entity/core"
)

type FieldInfo struct {
	Name  string
	Tag   string
	Index []int
}

func readStruct(t reflect.Type) []FieldInfo {
	infos := []FieldInfo{}
	for field := range t.Fields() {
		typ, _ := derefPtrTyp(field.Type)

		if field.Anonymous && typ.Kind() == reflect.Struct {
			baseIndex := append([]int(nil), field.Index...)

			for _, fi := range readStruct(typ) {
				infos = append(infos, FieldInfo{
					Name:  fi.Name,
					Index: append(baseIndex, fi.Index...),
					Tag:   fi.Tag,
				})
			}
			continue
		}

		infos = append(infos, FieldInfo{
			Name:  field.Name,
			Tag:   field.Tag.Get(core.TAG),
			Index: field.Index,
		})
	}

	return infos
}
