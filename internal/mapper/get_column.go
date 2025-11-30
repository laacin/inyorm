package mapper

import (
	"slices"

	"github.com/laacin/inyorm/internal/core"
)

func GetColumns(tag string, v any) ([]string, error) {
	s := getSchema(tag, v)

	var cols []string
	switch s.Type {

	case typeAny, typeColumn:
		cols = colsFromCol(v, s.Ptr, s.Slc)

	case typeStruct:
		cols = colsFromStruct(s.Index)

	case typeMap:
		if s.Slc {
			return nil, ErrInvalidSchema
		}
		cols = colsFromMap(v, s.Ptr)

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

	slices.Sort(cols)
	return cols, nil
}

// -- internal

func colsFromCol(v any, ptr, slc bool) []string {
	if !slc {
		if col, ok := v.(core.Column); ok {
			return []string{col.RawBase()}
		}
		return nil
	}

	var s []any
	if ptr {
		s = *v.(*[]any)
	} else {
		s = v.([]any)
	}

	cols := make([]string, len(s))
	for i, elem := range s {
		col, ok := elem.(core.Column)
		if !ok {
			continue
		}
		cols[i] = col.RawBase()
	}

	return cols
}

func colsFromStruct(fieldInfo []fieldInfo) []string {
	cols := make([]string, len(fieldInfo))
	for i, col := range fieldInfo {
		cols[i] = col.name
	}
	return cols
}

func colsFromMap(v any, ptr bool) []string {
	var m map[string]any
	if ptr {
		m = *v.(*map[string]any)
	} else {
		m = v.(map[string]any)
	}

	cols := make([]string, 0, len(m))
	for k := range m {
		cols = append(cols, k)
	}

	return cols
}
