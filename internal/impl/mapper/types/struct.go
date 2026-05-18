package types

import (
	"reflect"

	"github.com/laacin/inyorm/internal/core"
)

func readStruct(t reflect.Type) StructSchema {
	infos := StructSchema(map[string]core.FieldResult{})

	for field := range t.Fields() {
		typ, _ := DerefPtrTyp(field.Type)

		if field.Anonymous && typ.Kind() == reflect.Struct {
			baseIndex := append([]int(nil), field.Index...)
			infos.merge(readStruct(typ), baseIndex)
			continue
		}

		infos.add(field.Name, field.Tag.Get(core.TAG), field.Index)
	}

	return infos
}
