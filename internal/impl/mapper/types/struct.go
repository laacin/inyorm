package types

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/ddl"
)

type fieldInfo struct {
	Meta  ddl.ColumnMeta
	Index []int
}

func readStruct(t reflect.Type) []fieldInfo {
	infos := []fieldInfo{}

	for field := range t.Fields() {
		typ, _ := DerefPtrTyp(field.Type)

		if field.Anonymous && typ.Kind() == reflect.Struct {
			baseIndex := append([]int(nil), field.Index...)

			for _, fi := range readStruct(typ) {
				infos = append(infos, fieldInfo{
					Meta:  fi.Meta,
					Index: append(baseIndex, fi.Index...),
				})
			}
			continue
		}

		meta := ddl.ParseTag(field.Tag.Get(core.TAG))
		if meta.Name == "" {
			meta.Name = toSnake(field.Name)
		}

		infos = append(infos, fieldInfo{
			Meta:  meta,
			Index: field.Index,
		})
	}

	return infos
}

func toSnake(v string) string {
	if v == "" {
		return ""
	}

	var b strings.Builder
	runes := []rune(v)

	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := runes[i-1]

				if unicode.IsLower(prev) ||
					(unicode.IsUpper(prev) &&
						i+1 < len(runes) &&
						unicode.IsLower(runes[i+1])) {
					b.WriteByte('_')
				}
			}

			b.WriteRune(unicode.ToLower(r))
			continue
		}

		b.WriteRune(r)
	}

	return b.String()
}
