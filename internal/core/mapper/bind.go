package mapper

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/laacin/inyorm/internal/core"
)

func (m *Mapper) Bind(rows core.Rows, scanner any) error {
	info := m.ReadKind(scanner)

	if !info.Ptr && info.Kind != core.KindMap {
		return errors.New("scanner must be a pointer")
	}

	switch info.Kind {
	case core.KindStruct:
		if info.Slice {
			return bindFromStructSlc(rows, scanner, info.Schema)
		}
		return bindFromStruct(rows, scanner, info.Schema)

	case core.KindMap:
		if info.Slice {
			return bindFromMapSlc(rows, scanner)
		}
		return bindFromMap(rows, scanner)

	default:
		if info.Slice {
			return bindFromPrimSlc(rows, scanner)
		}
		return bindFromPrim(rows, scanner)
	}
}

// helpers

func bindFromStruct(rows core.Rows, binder any, schema Schema) error {
	defer rows.Close()

	if !rows.Next() {
		return rows.Err()
	}

	cols, _ := rows.Columns()
	addrs := make([]any, len(cols))
	bind, _ := derefPtrVal(reflect.ValueOf(binder))

	for i, col := range cols {
		idx, ok := schema.GetIndex(col)
		if !ok {
			addrs[i] = new(any)
			continue
		}

		field := bind.FieldByIndex(idx)
		if !field.CanAddr() {
			return fmt.Errorf("field %s is not addressable", col)
		}

		addrs[i] = field.Addr().Interface()
	}

	if err := rows.Scan(addrs...); err != nil {
		return err
	}

	return rows.Err()
}

func bindFromStructSlc(rows core.Rows, binder any, schema Schema) error {
	defer rows.Close()

	cols, _ := rows.Columns()
	args := make([]any, len(cols))

	slc, _ := derefPtrVal(reflect.ValueOf(binder))
	ln := slc.Len()

	i := 0
	for rows.Next() {
		if i < ln {
			// update elem
			for ci, col := range cols {
				if idx, ok := schema.GetIndex(col); ok {
					elem, _ := derefPtrVal(slc.Index(i))
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
			typ, ptrs := derefPtrTyp(slc.Type().Elem())
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

	return rows.Err()
}

func bindFromMap(rows core.Rows, binder any) error {
	defer rows.Close()

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

	bind, _ := derefPtrVal(reflect.ValueOf(binder))
	for i, col := range cols {
		bind.SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(tmp[i]))
	}

	return rows.Err()
}

func bindFromMapSlc(rows core.Rows, binder any) error {
	defer rows.Close()

	cols, _ := rows.Columns()
	slc, _ := derefPtrVal(reflect.ValueOf(binder))
	args := make([]any, len(cols))
	typ, _ := derefPtrTyp(slc.Type().Elem())
	ln := slc.Len()

	i := 0
	for rows.Next() {
		tmp := make([]any, len(cols))
		for i := range args {
			args[i] = &tmp[i]
		}

		if err := rows.Scan(args...); err != nil {
			return err
		}

		if i < ln {
			elem, _ := derefPtrVal(slc.Index(i))
			for ci, col := range cols {
				elem.SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(tmp[ci]))
			}
		} else {
			dummy := reflect.MakeMap(typ)
			for ci, col := range cols {
				dummy.SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(tmp[ci]))
			}

			slc.Set(reflect.Append(slc, dummy))
		}

		if err := rows.Scan(args...); err != nil {
			return err
		}

		i++
	}

	if i < ln {
		slc.Set(slc.Slice(0, i))
	}

	return rows.Err()
}

func bindFromPrim(rows core.Rows, value any) error {
	defer rows.Close()

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

	if err := rows.Scan(args...); err != nil {
		return err
	}

	return rows.Err()
}

func bindFromPrimSlc(rows core.Rows, value any) error {
	defer rows.Close()

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

	return rows.Err()
}

// helpers
func normalizePrimSlc(value any) ([]any, error) {
	slc, _ := derefPtrVal(reflect.ValueOf(value))

	vals := make([]any, slc.Len())
	for i := range slc.Len() {
		elem, _ := derefPtrVal(slc.Index(i))

		if !elem.CanAddr() {
			return nil, errors.New("passed no adressable value")
		}

		vals[i] = elem.Addr().Interface()
	}
	return vals, nil
}

func normalizePrim(value any) (any, error) {
	val, _ := derefPtrVal(reflect.ValueOf(value))
	if !val.CanAddr() {
		return nil, errors.New("passed no adressable value")
	}
	return val.Addr().Interface(), nil
}
