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
	typ, ptr, slc := resolve(v)
	if typ != typReflect && typ != typMap && !ptr {
		return ErrPtrExpected
	}

	switch typ {
	case typPrimitive:
		if slc {
			return scanSlcOfPrim(rows, v.(*[]any))
		}
		return scanPrim(rows, v)

	case typMap:
		if slc {
			return scanSlcOfMap(rows, v.(*[]map[string]any))
		}
		var m map[string]any
		if ptr {
			m = *v.(*map[string]any)
		} else {
			m = v.(map[string]any)
		}
		return scanMap(rows, m)

	}

	val, info := resolveRfl(v, tag, false)
	if info.err != nil {
		return info.err
	}

	if info.slc {
		return scanSlcOfStruct(rows, info.index, val)
	}

	return scanStruct(rows, info.index, val)
}

func scanPrim(rows RowScanner, v any) error {
	if rows.Next() {
		return rows.Scan(v)
	} else {
		return sql.ErrNoRows
	}
}

// BUG: doesn't work
func scanSlcOfPrim(rows RowScanner, v *[]any) error {
	m := make([]any, len(*v))
	for i := range m {
		x := *v
		m[i] = x[i]
	}

	if !rows.Next() {
		return sql.ErrNoRows
	}

	if err := rows.Scan(m...); err != nil {
		return err
	}

	return nil
}

func scanMap(rows RowScanner, mp map[string]any) error {
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

func scanSlcOfMap(rows RowScanner, mp *[]map[string]any) error {
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

func scanStruct(rows RowScanner, indexFields map[string][]int, strct reflect.Value) error {
	cols, _ := rows.Columns()
	args := make([]any, len(cols))

	for i, col := range cols {
		if index, exists := indexFields[col]; exists {
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

func scanSlcOfStruct(rows RowScanner, indexField map[string][]int, slice reflect.Value) error {
	typ := slice.Type().Elem()
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	cols, _ := rows.Columns()
	args := make([]any, len(cols))
	length := slice.Len()

	i := 0
	for rows.Next() {
		if i < length {
			// update elem
			for colIndex, col := range cols {
				if index, exists := indexField[col]; exists {
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
			slice.Set(reflect.Append(slice, dummy))
		}
		i++
	}

	if i < length {
		slice.Set(slice.Slice(0, i))
	}

	return nil
}
