package mapper

import "slices"

func GetColumns(tag string, v any) ([]string, error) {
	s, err := getSchema(tag, v)
	if err != nil {
		return nil, err
	}

	var cols []string
	switch s.Type {
	case typeColumn:
		cols = colsFromCol(v, s.Slc)

	case typeString:
		cols = colsFromString(v, s.Slc, s.Ptr)

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

	for _, col := range cols {
		if col == "" || col == "*" {
			return nil, ErrInvalidColumn
		}
	}

	slices.Sort(cols)
	return cols, nil
}

// -- internal

func colsFromCol(v any, slc bool) []string {
	type column interface{ RawBase() string }

	if !slc {
		col := v.(column)
		return []string{col.RawBase()}
	}

	cols := v.([]column)
	result := make([]string, len(cols))
	for i, col := range cols {
		result[i] = col.RawBase()
	}

	return result
}

func colsFromString(v any, slc, ptr bool) []string {
	if !slc {
		var str string
		if ptr {
			str = *(v).(*string)
		} else {
			str = v.(string)
		}

		return []string{str}
	}

	var cols []string
	if ptr {
		cols = *(v).(*[]string)
	} else {
		cols = v.([]string)
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
		m = *(v).(*map[string]any)
	} else {
		m = v.(map[string]any)
	}

	cols := make([]string, 0, len(m))
	for k := range m {
		cols = append(cols, k)
	}

	return cols
}
