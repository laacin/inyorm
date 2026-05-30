package mapper

import (
	"errors"
	"fmt"
	"reflect"
	"slices"

	"github.com/laacin/inyorm/internal/core"
)

// Errs
var (
	ErrMustVal     = errors.New("primitive input types must be values")
	ErrColMismatch = errors.New("col mismatch between columns and values")
	ErrColNotFound = func(col string) error { return fmt.Errorf("column %s not found", col) }
)

func (m *Mapper) ReadValues(cols []string, v any) ([]any, error) {
	info := m.ReadKind(v)
	val, _ := derefPtrVal(reflect.ValueOf(v))

	switch info.Kind {
	case core.KindStruct:
		if info.Slice {
			return valsFromStructSlc(cols, val, info.Schema)
		}
		return valsFromStruct(cols, val, info.Schema)

	case core.KindMap:
		if info.Slice {
			return valsFromMapSlc(cols, val)
		}
		return valsFromMap(cols, val)

	case core.KindString, core.KindInt, core.KindUint, core.KindBool, core.KindFloat:
		if info.Slice {
			return valsFromPrimSlc(val)
		}
		return valsFromPrim(val)

	case core.KindAny:
		if info.Slice {
			return valsFromPrimSlc(val)
		}
	}

	return nil, errors.New("something went wrong")
}

// helpers

func valsFromStruct(cols []string, val reflect.Value, schema Schema) ([]any, error) {
	args := make([]any, len(cols))

	for i, col := range cols {
		if idx, ok := schema.GetIndex(col); ok {
			args[i] = val.FieldByIndex(idx).Interface()
			continue
		}
		return nil, ErrColNotFound(col)
	}

	return args, nil
}

func valsFromStructSlc(cols []string, val reflect.Value, schema Schema) ([]any, error) {
	args := make([]any, len(cols)*val.Len())

	for i := range val.Len() {
		elem, _ := derefPtrVal(val.Index(i))

		for ci, col := range cols {
			if findex, ok := schema.GetIndex(col); ok {
				args[i*len(cols)+ci] = elem.FieldByIndex(findex).Interface()
				continue
			}
			return nil, ErrColNotFound(col)
		}
	}

	return args, nil
}

func valsFromMap(cols []string, val reflect.Value) ([]any, error) {
	args := make([]any, len(cols))

	for i, col := range cols {
		v, _ := derefPtrVal(val.MapIndex(reflect.ValueOf(col)))

		if !v.IsValid() {
			return nil, ErrColNotFound(col)
		}

		args[i] = v.Interface()
	}

	return args, nil
}

func valsFromMapSlc(cols []string, val reflect.Value) ([]any, error) {
	args := make([]any, len(cols)*val.Len())

	for i := range val.Len() {
		elem, _ := derefPtrVal(val.Index(i))

		for ci, col := range cols {
			v, _ := derefPtrVal(elem.MapIndex(reflect.ValueOf(col)))

			if !v.IsValid() {
				return nil, ErrColNotFound(col)
			}

			args[i*len(cols)+ci] = v.Interface()
		}
	}

	return args, nil
}

func valsFromPrim(val reflect.Value) ([]any, error) {
	return []any{val.Interface()}, nil
}

func valsFromPrimSlc(val reflect.Value) ([]any, error) {
	args := make([]any, val.Len())

	for i := range val.Len() {
		elem, _ := derefPtrVal(val.Index(i))

		if !slices.Contains(primitives, elem.Kind()) {
			return nil, errors.New("[]any must contains only primitive types")
		}

		args[i] = elem.Interface()
	}

	return args, nil
}

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
