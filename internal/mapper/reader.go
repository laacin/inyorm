package mapper

import (
	"reflect"
	"sort"
)

type ReadResult struct {
	Args    []any
	Columns []string
	Rows    int
}

func Read(tag string, v any) (ReadResult, error) {
	typ, ptr, slc := resolve(v)
	if typ != typReflect && ptr {
		return ReadResult{}, ErrValueExpected
	}

	switch typ {
	case typPrimitive:
		if slc {
			return readSlcOfPrim(v.([]any)), nil
		}
		return readPrim(v), nil

	case typMap:
		if slc {
			return readSlcOfMap(v.([]map[string]any))
		}
		return readMap(v.(map[string]any)), nil
	}

	val, info := resolveRfl(v, tag, true)
	if info.err != nil {
		return ReadResult{}, info.err
	}

	if info.slc {
		return readSlcOfStruct(val, info.index)
	}
	return readStruct(val, info.index)
}

func readPrim(v any) ReadResult {
	return ReadResult{Args: []any{v}}
}

func readSlcOfPrim(v []any) ReadResult {
	return ReadResult{Args: v}
}

func readMap(mp map[string]any) ReadResult {
	cols := make([]string, len(mp))
	args := make([]any, len(mp))

	i := 0
	for k, v := range mp {
		cols[i] = k
		args[i] = v
		i++
	}

	return ReadResult{
		Columns: cols,
		Args:    args,
		Rows:    1,
	}
}

func readSlcOfMap(mp []map[string]any) (ReadResult, error) {
	rows := len(mp)
	if rows == 0 {
		return ReadResult{}, ErrEmptySlice
	}

	first := mp[0]
	colNum := len(first)

	cols := make([]string, 0, colNum)
	for k := range first {
		cols = append(cols, k)
	}
	sort.Strings(cols)

	args := make([]any, colNum*rows)
	for idx, m := range mp {
		for i, col := range cols {
			v, ok := m[col]
			if !ok {
				return ReadResult{}, ErrColumnMismatch
			}

			args[idx*colNum+i] = v
		}
	}

	return ReadResult{
		Rows:    rows,
		Columns: cols,
		Args:    args,
	}, nil
}

func readStruct(v reflect.Value, indexField map[string][]int) (ReadResult, error) {
	colNum := len(indexField)

	cols := make([]string, 0, colNum)
	args := make([]any, colNum)
	for k := range indexField {
		cols = append(cols, k)
	}
	sort.Strings(cols)

	for i, col := range cols {
		idx, ok := indexField[col]
		if !ok {
			return ReadResult{}, ErrColumnMismatch
		}

		args[i] = v.FieldByIndex(idx).Interface()
	}

	return ReadResult{
		Rows:    1,
		Columns: cols,
		Args:    args,
	}, nil
}

func readSlcOfStruct(v reflect.Value, indexField map[string][]int) (ReadResult, error) {
	rows := v.Len()
	if rows == 0 {
		return ReadResult{}, ErrEmptySlice
	}

	colNum := len(indexField)
	cols := make([]string, 0, colNum)
	args := make([]any, rows*colNum)

	for k := range indexField {
		cols = append(cols, k)
	}
	sort.Strings(cols)

	for row := range rows {
		item := v.Index(row)
		for i, col := range cols {
			idx, ok := indexField[col]
			if !ok {
				return ReadResult{}, ErrColumnMismatch
			}
			args[row*colNum+i] = item.FieldByIndex(idx).Interface()
		}
	}

	return ReadResult{
		Rows:    rows,
		Columns: cols,
		Args:    args,
	}, nil
}
