package mapper

import (
	"errors"
	"fmt"
	"reflect"
	"slices"

	"github.com/laacin/inyorm/internal/impl/mapper/types"
)

// Errs
var (
	ErrMustVal     = errors.New("primitive input types must be values")
	ErrColMismatch = errors.New("col mismatch between columns and values")
	ErrColNotFound = func(col string) error { return fmt.Errorf("column %s not found", col) }
)

type Result struct {
	Rows    int
	Columns []string
	Args    []any
}

func ReadValues(cols []string, values any) (*Result, error) {
	val := reflect.ValueOf(values)
	info := types.ReadInfo(val.Type())
	val, _ = types.DerefPtrVal(val)

	switch info.Kind {
	case types.KindStruct:
		if info.IsSlc() {
			return valsByStructSlc(cols, val, info.Schema)
		}
		return valsByStruct(cols, val, info.Schema)

	case types.KindMap:
		if info.IsSlc() {
			return valsByMapSlc(cols, val)
		}
		return valsByMap(cols, val)

	case types.KindString, types.KindInt, types.KindUint, types.KindBool, types.KindFloat:
		if info.IsSlc() {
			return valsByPrimSlc(cols, val)
		}
		return valsByPrim(cols, val)

	case types.KindAny:
		if info.IsSlc() {
			return valsByPrimSlc(cols, val)
		}
	}

	return nil, errors.New("something went wrong")
}

// Readers
func valsByStruct(cols []string, val reflect.Value, schema types.StructSchema) (*Result, error) {
	if schema.Len()%len(cols) != 0 {
		return nil, ErrColMismatch
	}

	args := make([]any, 0, schema.Len())
	for _, col := range cols {
		if idx, ok := schema.GetIndex(col); ok {
			args = append(args, val.FieldByIndex(idx).Interface())
			continue
		}
		return nil, ErrColNotFound(col)
	}

	return &Result{
		Rows:    1,
		Columns: cols,
		Args:    args,
	}, nil
}

func valsByStructSlc(cols []string, val reflect.Value, schema types.StructSchema) (*Result, error) {
	if (schema.Len()*val.Len())%len(cols) != 0 {
		return nil, ErrColMismatch
	}

	args := make([]any, 0, val.Len()*schema.Len())

	for i := range val.Len() {
		elem, _ := types.DerefPtrVal(val.Index(i))

		for _, col := range cols {
			if findex, ok := schema.GetIndex(col); ok {
				args = append(args, elem.FieldByIndex(findex).Interface())
				continue
			}
			return nil, ErrColNotFound(col)
		}
	}

	return &Result{
		Rows:    len(args) / len(cols),
		Columns: cols,
		Args:    args,
	}, nil
}

func valsByMap(cols []string, val reflect.Value) (*Result, error) {
	if len(cols) != val.Len() {
		return nil, ErrColMismatch
	}

	m, _ := reflect.TypeAssert[map[string]any](val)
	args := make([]any, 0, val.Len())
	for _, col := range cols {
		args = append(args, m[col])
	}

	return &Result{
		Rows:    1,
		Columns: cols,
		Args:    args,
	}, nil
}

func valsByMapSlc(cols []string, val reflect.Value) (*Result, error) {
	args := []any{}

	for i := range val.Len() {
		elem, _ := types.DerefPtrVal(val.Index(i))
		m, _ := reflect.TypeAssert[map[string]any](elem)

		arg := make([]any, 0, len(m))
		for _, col := range cols {
			arg = append(arg, m[col])
		}

		args = append(args, arg...)
	}

	if len(args)%len(cols) != 0 {
		return nil, ErrColMismatch
	}

	return &Result{
		Rows:    len(args) / len(cols),
		Columns: cols,
		Args:    args,
	}, nil
}

func valsByPrim(cols []string, val reflect.Value) (*Result, error) {
	if len(cols) != 1 {
		return nil, ErrColMismatch
	}

	return &Result{
		Rows:    1,
		Columns: cols,
		Args:    []any{val.Interface()},
	}, nil
}

func valsByPrimSlc(cols []string, val reflect.Value) (*Result, error) {
	if val.Len()%len(cols) != 0 {
		return nil, ErrColMismatch
	}

	args := make([]any, 0, val.Len())
	for i := range val.Len() {
		elem, _ := types.DerefPtrVal(val.Index(i))
		if !slices.Contains(primitives, elem.Kind()) {
			return nil, errors.New("[]any must contains only primitive types")
		}

		args = append(args, elem.Interface())
	}

	return &Result{
		Rows:    val.Len() / len(cols),
		Columns: cols,
		Args:    args,
	}, nil
}

// Helpers
var primitives = []reflect.Kind{
	reflect.Bool,

	reflect.Int,
	reflect.Int8,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,

	reflect.Uint,
	reflect.Uint8,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,

	reflect.Float32,
	reflect.Float64,

	reflect.String,
}
