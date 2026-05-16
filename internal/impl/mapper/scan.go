package mapper

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/laacin/inyorm/internal/impl/mapper/types"
	"github.com/laacin/inyorm/internal/ir/driver"
)

func Scan(rows driver.Rows, scanner any) error {
	info := types.ReadInfo(reflect.TypeOf(scanner))

	if !info.IsPtr() && info.Kind != types.KindMap {
		return errors.New("scanner must be a pointer")
	}

	switch info.Kind {
	case types.KindStruct:
		if info.IsSlc() {
			return scanByStructSlc(rows, scanner, info.Schema)
		}
		return scanByStruct(rows, scanner, info.Schema)

	case types.KindMap:
		if info.IsSlc() {
			return scanByMapSlc(rows, scanner)
		}
		return scanByMap(rows, scanner)

	default:
		if info.IsSlc() {
			return scanByPrimSlc(rows, scanner)
		}
		return scanByPrim(rows, scanner)
	}
}

// --- Scanners
func scanByStruct(rows driver.Rows, value any, schema types.StructSchema) error {
	if !rows.Next() {
		return nil
	}

	cols, _ := rows.Columns()
	addrs := make([]any, len(cols))
	val, _ := types.DerefPtrVal(reflect.ValueOf(value))

	for i, col := range cols {
		idx, ok := schema.GetIndex(col)
		if !ok {
			addrs[i] = new(any)
			continue
		}

		field := val.FieldByIndex(idx)
		if !field.CanAddr() {
			return fmt.Errorf("field %s is not addressable", col)
		}

		addrs[i] = field.Addr().Interface()
	}

	return rows.Scan(addrs...)
}

func scanByStructSlc(rows driver.Rows, value any, schema types.StructSchema) error {
	cols, _ := rows.Columns()

	args := make([]any, len(cols))

	slc, _ := types.DerefPtrVal(reflect.ValueOf(value))
	ln := slc.Len()

	i := 0
	for rows.Next() {
		if i < ln {
			// update elem
			for ci, col := range cols {
				if idx, ok := schema.GetIndex(col); ok {
					elem, _ := types.DerefPtrVal(slc.Index(i))
					field := elem.FieldByIndex(idx)
					if !field.CanAddr() {
						return fmt.Errorf("field %s is not addressable", col)
					}

					args[ci] = field.Addr().Interface()
				} else {
					args[ci] = new(any)
				}
			}

			if err := rows.Scan(args...); err != nil {
				return err
			}

		} else {
			// create elem
			typ, ptrs := types.DerefPtrTyp(slc.Type().Elem())
			dummy := reflect.New(typ).Elem()

			for ci, col := range cols {
				if idx, ok := schema.GetIndex(col); ok {
					field := dummy.FieldByIndex(idx)
					if !field.CanAddr() {
						return fmt.Errorf("field %s is not addressable", col)
					}

					args[ci] = field.Addr().Interface()
				} else {
					args[ci] = new(any)
				}
			}

			if err := rows.Scan(args...); err != nil {
				return err
			}

			for range ptrs {
				dummy = dummy.Addr()
			}

			slc.Set(reflect.Append(slc, dummy))
		}
		i++
	}

	if i < ln {
		slc.Set(slc.Slice(0, i))
	}

	return nil
}

func scanByMap(rows driver.Rows, value any) error {
	mp, ok := value.(map[string]any)
	if !ok {
		return errors.New("map scanning must receive map[string]any or *[]map[string]any")
	}

	if !rows.Next() {
		return nil
	}

	cols, _ := rows.Columns()

	args := make([]any, len(cols))
	tmp := make([]any, len(cols))
	for i := range args {
		args[i] = &tmp[i]
	}

	if err := rows.Scan(args...); err != nil {
		return err
	}

	for i, col := range cols {
		mp[col] = tmp[i]
	}

	return nil
}

func scanByMapSlc(rows driver.Rows, value any) error {
	maps, ok := value.(*[]map[string]any)
	if !ok {
		return errors.New("map scanning must receive map[string]any or *[]map[string]any")
	}

	ln := len(*maps)
	cols, _ := rows.Columns()
	args := make([]any, len(cols))
	tmp := make([]any, len(cols))
	for i := range args {
		args[i] = &tmp[i]
	}

	i := 0
	for rows.Next() {
		if i < ln {
			if err := rows.Scan(args...); err != nil {
				return err
			}

			for i, col := range cols {
				(*maps)[i][col] = tmp[i]
			}
		} else {
			if err := rows.Scan(args...); err != nil {
				return err
			}

			mp := make(map[string]any, len(cols))
			for i, col := range cols {
				mp[col] = tmp[i]
			}
			*maps = append(*maps, mp)
		}
		i++
	}

	if i < ln {
		*maps = (*maps)[:i]
	}

	return nil
}

func scanByPrim(rows driver.Rows, value any) error {
	if !rows.Next() {
		return nil
	}

	val, err := normalizePrim(value)
	if err != nil {
		return err
	}

	cols, _ := rows.Columns()
	if len(cols) == 1 {
		return rows.Scan(val)
	}

	args := make([]any, len(cols))
	for i := range args {
		if i == 0 {
			args[i] = val
			continue
		}
		args[i] = new(any)
	}
	return rows.Scan(args...)
}

func scanByPrimSlc(rows driver.Rows, value any) error {
	vals, err := normalizePrimSlc(value)
	if err != nil {
		return err
	}

	cols, _ := rows.Columns()
	args := make([]any, len(cols))

	current := 0
	for rows.Next() {
		for i := range cols {
			if len(vals) > current {
				args[i] = vals[current]
			} else {
				args[i] = new(any)
			}
			current++
		}

		if err := rows.Scan(args...); err != nil {
			return err
		}
	}

	return nil
}

// helpers
func normalizePrimSlc(value any) ([]any, error) {
	slc, _ := types.DerefPtrVal(reflect.ValueOf(value))

	vals := make([]any, slc.Len())
	for i := range slc.Len() {
		elem, _ := types.DerefPtrVal(slc.Index(i))

		if !elem.CanAddr() {
			return nil, errors.New("passed no adressable value")
		}

		vals[i] = elem.Addr().Interface()
	}
	return vals, nil
}

func normalizePrim(value any) (any, error) {
	val, _ := types.DerefPtrVal(reflect.ValueOf(value))
	if !val.CanAddr() {
		return nil, errors.New("passed no adressable value")
	}
	return val.Addr().Interface(), nil
}
