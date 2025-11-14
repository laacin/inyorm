package mapper

import "reflect"

func ReadValues(columnTag string, v any) (int, []string, []any, error) {
	val, typ, err := resolveInput(true, v)
	if err != nil {
		return 0, nil, nil, err
	}

	switch typ {
	case typStruct:
		row, columns, values := readStruct(columnTag, val)
		return row, columns, values, nil

	case typSlice, typArray:
		rows, columns, values := readSlice(columnTag, val)
		if rows == 0 {
			return 0, nil, nil, ErrEmptyValues
		}
		return rows, columns, values, nil
	}

	return 0, nil, nil, ErrInvalidType
}

func readStruct(columnTag string, val reflect.Value) (int, []string, []any) {
	fields := cachedFields(columnTag, val.Type())

	columns := make([]string, 0, len(fields))
	values := make([]any, 0, len(fields))

	for _, f := range fields {
		target := val.FieldByIndex(f.index)
		if target.Kind() == reflect.Pointer {

			if target.IsNil() {
				values = append(values, nil)
				columns = append(columns, f.name)
				continue
			}
			target = target.Elem()
		}

		values = append(values, target.Interface())
		columns = append(columns, f.name)
	}

	return 1, columns, values
}

func readSlice(columnTag string, val reflect.Value) (int, []string, []any) {
	rows := val.Len()
	if rows == 0 {
		return 0, nil, nil
	}

	dummy := val.Index(0).Type()
	if dummy.Kind() == reflect.Pointer {
		dummy = dummy.Elem()
	}

	fields := cachedFields(columnTag, dummy)
	fieldsNum := len(fields)
	if fieldsNum < 1 {
		return 0, nil, nil
	}

	columns := make([]string, 0, fieldsNum)
	for _, f := range fields {
		columns = append(columns, f.name)
	}

	values := make([]any, 0, rows*fieldsNum)
	for i := range rows {
		elem := val.Index(i)
		if elem.Kind() == reflect.Pointer {
			if elem.IsNil() {
				continue
			}
			elem = elem.Elem()
		}

		for _, f := range fields {
			target := elem.FieldByIndex(f.index)

			if target.Kind() == reflect.Pointer {
				if target.IsNil() {
					values = append(values, nil)
					continue
				}
				target = target.Elem()
			}

			values = append(values, target.Interface())
		}
	}

	return rows, columns, values
}
