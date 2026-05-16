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
	args := make([]any, len(cols))
	for i, col := range cols {
		if idx, ok := schema.GetIndex(col); ok {
			args[i] = val.FieldByIndex(idx).Interface()
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
	args := make([]any, len(cols)*val.Len())

	for i := range val.Len() {
		elem, _ := types.DerefPtrVal(val.Index(i))

		for ci, col := range cols {
			if findex, ok := schema.GetIndex(col); ok {
				args[i*len(cols)+ci] = elem.FieldByIndex(findex).Interface()
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
	m, _ := reflect.TypeAssert[map[string]any](val)
	args := make([]any, len(cols))

	for i, col := range cols {
		if v, ok := m[col]; ok {
			args[i] = v
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

func valsByMapSlc(cols []string, val reflect.Value) (*Result, error) {
	args := make([]any, len(cols)*val.Len())

	for i := range val.Len() {
		elem, _ := types.DerefPtrVal(val.Index(i))
		m, _ := reflect.TypeAssert[map[string]any](elem)

		for ci, col := range cols {
			if rslt, ok := m[col]; ok {
				args[i*len(cols)+ci] = rslt
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

func valsByPrim(cols []string, val reflect.Value) (*Result, error) {
	return &Result{
		Rows:    1,
		Columns: cols,
		Args:    []any{val.Interface()},
	}, nil
}

func valsByPrimSlc(cols []string, val reflect.Value) (*Result, error) {
	args := make([]any, val.Len())

	for i := range val.Len() {
		elem, _ := types.DerefPtrVal(val.Index(i))
		if !slices.Contains(primitives, elem.Kind()) {
			return nil, errors.New("[]any must contains only primitive types")
		}

		args[i] = elem.Interface()
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
