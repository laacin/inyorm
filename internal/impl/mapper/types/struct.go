package types

import (
	"reflect"

	"github.com/laacin/inyorm/internal/core"
)

func readStruct(t reflect.Type) core.StructInfo {
	info := core.NewStructInfo()

	for field := range t.Fields() {
		typ, _ := DerefPtrTyp(field.Type)

		if field.Anonymous && typ.Kind() == reflect.Struct {
			baseIndex := append([]int(nil), field.Index...)
			info.Merge(baseIndex, readStruct(typ))
			continue
		}

		info.Add(field.Name, field.Tag.Get(core.TAG), field.Index)
	}

	return info
}
