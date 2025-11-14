package mapper

import (
	"database/sql"
	"reflect"
)

func BindRows(rows *sql.Rows, columnTag string, values any) error {
	val, typ, err := resolveInput(false, values)
	if err != nil {
		return err
	}

	if typ == typStruct {
		return bindStruct(rows, columnTag, val)
	}

	if typ == typSlice || typ == typArray {
		return bindSlice(rows, columnTag, val)
	}

	return ErrInvalidType
}

func indexFields(columnTag string, typ reflect.Type) map[string][]int {
	fields := cachedFields(columnTag, typ)
	indexField := make(map[string][]int, len(fields))
	for _, f := range fields {
		indexField[f.name] = f.index
	}

	return indexField
}

func bindStruct(rows *sql.Rows, columnTag string, strct reflect.Value) error {
	indexMap := indexFields(columnTag, strct.Type())

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	args := make([]any, len(cols))
	for i, col := range cols {
		if index, exists := indexMap[col]; exists {
			args[i] = strct.FieldByIndex(index).Addr().Interface()
		} else {
			args[i] = new(any)
		}
	}

	if rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
	} else {
		return sql.ErrNoRows
	}

	return nil
}

func bindSlice(rows *sql.Rows, columnTag string, slice reflect.Value) error {
	typ := slice.Type().Elem()
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	indexMap := indexFields(columnTag, typ)
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	args := make([]any, len(cols))
	length := slice.Len()
	i := 0
	for rows.Next() {
		if i < length {
			// update elem
			for colIndex, col := range cols {
				if index, exists := indexMap[col]; exists {
					addr := slice.Index(i).FieldByIndex(index).Addr().Interface()
					args[colIndex] = addr
				} else {
					args[colIndex] = new(any)
				}
			}

			if err := rows.Scan(args...); err != nil {
				return err
			}
		} else {
			// create elem
			dummy := reflect.New(typ).Elem()

			for colIndex, col := range cols {
				if index, exists := indexMap[col]; exists {
					addr := dummy.FieldByIndex(index).Addr().Interface()
					args[colIndex] = addr
				} else {
					args[colIndex] = new(any)
				}
			}

			if err := rows.Scan(args...); err != nil {
				return err
			}
			slice.Set(reflect.Append(slice, dummy))
		}
		i++
	}

	return nil
}

// TODO: add support for primitive types and maps
