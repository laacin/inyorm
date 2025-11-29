package mapper

import (
	"reflect"

	"github.com/laacin/inyorm/internal/mapper/schema"
)

type ReadResult struct {
	Rows    int
	Columns []string
	Args    []any
}

func Read(tag string, columns, v any) (*ReadResult, error) {
	cols, err := schema.GetColumns(tag, columns)
	if err != nil {
		return nil, err
	}

	s, err := schema.GetSchema(tag, v)
	if err != nil {
		return nil, err
	}

	switch s.Type {
	case schema.TypeString, schema.TypeInt, schema.TypeUint,
		schema.TypeFloat, schema.TypeBool:
		return readPrim(cols, v, s.Slc, s.Ptr)

	case schema.TypeAny:
		return readAny(cols, v, s.Slc, s.Ptr)

	case schema.TypeMap:
		if s.Slc {
			return readSlcOfMap(cols, v, s.Ptr)
		}
		return readMap(cols, v, s.Ptr)

	case schema.TypeStruct:
		if s.Slc {
			return readSlcOfStruct(cols, v, s.IndexMap(), s.Ptr)
		}
		return readStruct(cols, v, s.IndexMap(), s.Ptr)

	default:
		return nil, ErrUnexpectedType
	}
}

// -- internal

func readPrim(cols []string, v any, slc, ptr bool) (*ReadResult, error) {
	if slc || ptr {
		return nil, ErrUnexpectedType
	}

	if len(cols) != 1 {
		return nil, ErrColumnMismatch
	}

	return &ReadResult{
		Rows:    1,
		Columns: cols,
		Args:    []any{v},
	}, nil
}

func readAny(cols []string, v any, slc, ptr bool) (*ReadResult, error) {
	if !slc {
		return readPrim(cols, v, slc, ptr)
	}

	val := *(v).(*[]any)
	colNum, valNum := len(cols), len(val)
	if valNum == 0 {
		return nil, ErrEmptySlice
	}

	return &ReadResult{
		Rows:    valNum / colNum,
		Columns: cols,
		Args:    val,
	}, nil
}

func readMap(cols []string, v any, ptr bool) (*ReadResult, error) {
	var mp map[string]any
	if ptr {
		mp = *(v).(*map[string]any)
	} else {
		mp = v.(map[string]any)
	}

	args := make([]any, len(mp))

	for i, col := range cols {
		args[i] = mp[col]
	}

	return &ReadResult{
		Columns: cols,
		Args:    args,
		Rows:    1,
	}, nil
}

func readSlcOfMap(cols []string, v any, ptr bool) (*ReadResult, error) {
	var mp []map[string]any
	if ptr {
		mp = *(v).(*[]map[string]any)
	} else {
		mp = v.([]map[string]any)
	}

	colNum, rows := len(cols), len(mp)
	if rows == 0 {
		return nil, ErrEmptySlice
	}

	args := make([]any, colNum*rows)

	for idx, m := range mp {
		for i, col := range cols {
			v, ok := m[col]
			if !ok {
				return nil, ErrColumnMismatch
			}

			args[idx*colNum+i] = v
		}
	}

	return &ReadResult{
		Rows:    rows,
		Columns: cols,
		Args:    args,
	}, nil
}

func readStruct(cols []string, v any, indexField map[string][]int, ptr bool) (*ReadResult, error) {
	val := reflect.ValueOf(v)
	if ptr {
		val = val.Elem()
	}

	colNum := len(indexField)
	args := make([]any, colNum)

	for i, col := range cols {
		idx, ok := indexField[col]
		if !ok {
			return nil, ErrColumnMismatch
		}

		args[i] = val.FieldByIndex(idx).Interface()
	}

	return &ReadResult{
		Rows:    1,
		Columns: cols,
		Args:    args,
	}, nil
}

func readSlcOfStruct(cols []string, v any, indexField map[string][]int, ptr bool) (*ReadResult, error) {
	val := reflect.ValueOf(v)
	if ptr {
		val = val.Elem()
	}

	rows := val.Len()
	if rows == 0 {
		return nil, ErrEmptySlice
	}

	colNum := len(indexField)
	args := make([]any, rows*colNum)

	for row := range rows {
		item := val.Index(row)
		for i, col := range cols {
			idx, ok := indexField[col]
			if !ok {
				return nil, ErrColumnMismatch
			}
			args[row*colNum+i] = item.FieldByIndex(idx).Interface()
		}
	}

	return &ReadResult{
		Rows:    rows,
		Columns: cols,
		Args:    args,
	}, nil
}
