package mapper

import (
	"slices"

	"github.com/laacin/inyorm/internal/core"
)

func GetColumns(tag string, vals []any) ([]string, error) {
	var cols []string

	for _, v := range vals {
		s := getSchema(tag, v)

		if s.Slc {
			return nil, ErrInvalidSchema
		}

		switch s.Type {

		case typeAny, typeColumn:
			colsFromCol(&cols, v)

		case typeStruct:
			colsFromStruct(&cols, s.Index)

		case typeMap:
			if s.Slc {
				return nil, ErrInvalidSchema
			}
			colsFromMap(&cols, v, s.Ptr)

		default:
			return nil, ErrInvalidSchema
		}

		if len(cols) == 0 {
			return nil, ErrNoColumns
		}

		for _, col := range cols {
			if col == "" || col == "*" {
				return nil, ErrInvalidColumn
			}
		}
	}

	slices.Sort(cols)
	return cols, nil
}

// -- internal

func colsFromCol(cols *[]string, v any) {
	if col, ok := v.(core.Column); ok {
		*cols = append(*cols, col.RawBase())
	}
}

func colsFromStruct(cols *[]string, fieldInfo []fieldInfo) {
	for _, col := range fieldInfo {
		*cols = append(*cols, col.name)
	}
}

func colsFromMap(cols *[]string, v any, ptr bool) {
	var m map[string]any
	if ptr {
		m = *v.(*map[string]any)
	} else {
		m = v.(map[string]any)
	}

	for k := range m {
		*cols = append(*cols, k)
	}
}
