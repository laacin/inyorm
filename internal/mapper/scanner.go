package mapper

import (
	"database/sql"
	"reflect"
)

type RowScanner interface {
	Columns() ([]string, error)
	Next() bool
	Scan(...any) error
}

func Scan(rows RowScanner, tag string, v any) error {
	s := getSchema(tag, v)

	if s.Type != typeMap && !s.Slc && !s.Ptr {
		return ErrPtrExpected
	}

	switch s.Type {
	case typeString, typeInt, typeUint,
		typeFloat, typeBool:
		return scanPrim(rows, v)

	case typeStruct:
		if s.Slc {
			return scanSlcOfStruct(rows, s.IndexMap(), v)
		}
		return scanStruct(rows, s.IndexMap(), v)

	case typeMap:
		if s.Slc {
			return scanSlcOfMap(rows, v)
		}
		return scanMap(rows, v, s.Ptr)

	default:
		return ErrUnexpectedType
	}
}

// -- internal

func scanPrim(rows RowScanner, v any) error {
	if rows.Next() {
		return rows.Scan(v)
	} else {
		return sql.ErrNoRows
	}
}

func scanMap(rows RowScanner, v any, ptr bool) error {
	var mp map[string]any
	if ptr {
		mp = *(v).(*map[string]any)
	} else {
		mp = v.(map[string]any)
	}

	cols, _ := rows.Columns()
	args := make([]any, len(cols))
	for i := range args {
		var x any
		args[i] = &x
	}

	if rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
	} else {
		return sql.ErrNoRows
	}

	for i, col := range cols {
		mp[col] = *(args[i].(*any))
	}

	return nil
}

func scanSlcOfMap(rows RowScanner, v any) error {
	mp := v.(*[]map[string]any)

	cols, _ := rows.Columns()
	length := len(*mp)
	args := make([]any, len(cols))
	for idx := range args {
		var x any
		args[idx] = &x
	}

	i := 0
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}

		rowMap := make(map[string]any, len(cols))
		for idx, col := range cols {
			rowMap[col] = *(args[idx].(*any))
		}

		if i < length {
			(*mp)[i] = rowMap
		} else {
			*mp = append(*mp, rowMap)
		}
		i++
	}

	*mp = (*mp)[:i]
	return nil
}

func scanStruct(rows RowScanner, indexFields map[string][]int, v any) error {
	val := reflect.ValueOf(v).Elem()

	cols, _ := rows.Columns()
	args := make([]any, len(cols))

	for i, col := range cols {
		if index, exists := indexFields[col]; exists {
			args[i] = val.FieldByIndex(index).Addr().Interface()
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

func scanSlcOfStruct(rows RowScanner, indexField map[string][]int, v any) error {
	val := reflect.ValueOf(v).Elem()

	typ := val.Type().Elem()
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	cols, _ := rows.Columns()
	args := make([]any, len(cols))
	length := val.Len()

	i := 0
	for rows.Next() {
		if i < length {
			// update elem
			for colIndex, col := range cols {
				if index, exists := indexField[col]; exists {
					addr := val.Index(i).FieldByIndex(index).Addr().Interface()
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
				if index, exists := indexField[col]; exists {
					addr := dummy.FieldByIndex(index).Addr().Interface()
					args[colIndex] = addr
				} else {
					args[colIndex] = new(any)
				}
			}

			if err := rows.Scan(args...); err != nil {
				return err
			}
			val.Set(reflect.Append(val, dummy))
		}
		i++
	}

	if i < length {
		val.Set(val.Slice(0, i))
	}

	return nil
}
