package mapper

import (
	"reflect"
	"sync"
)

type valueTyp int

const (
	typStruct valueTyp = iota
	typArray
	typSlice
	typPrimitive
)

var resolveCache sync.Map

func resolveInput(acceptValue bool, v any) (reflect.Value, valueTyp, error) {
	if v == nil {
		return reflect.Value{}, 0, ErrExpectedPointer
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return reflect.Value{}, 0, ErrExpectedPointer
		}
		val = val.Elem()

	} else if !acceptValue {
		return reflect.Value{}, 0, ErrExpectedPointer
	}

	typ := val.Type()

	if cached, ok := resolveCache.Load(typ); ok {
		return val, cached.(valueTyp), nil
	}

	switch typ.Kind() {
	case reflect.Struct:
		resolveCache.Store(typ, typStruct)
		return val, typStruct, nil

	case reflect.Slice:
		if err := resolveSlice(val); err != nil {
			return reflect.Value{}, 0, err
		}
		resolveCache.Store(typ, typSlice)
		return val, typSlice, nil

	case reflect.Array:
		if err := resolveSlice(val); err != nil {
			return reflect.Value{}, 0, err
		}
		resolveCache.Store(typ, typArray)
		return val, typArray, nil

	default:
		resolveCache.Store(typ, typPrimitive)
		return val, typPrimitive, nil
	}

}

func resolveSlice(val reflect.Value) error {
	stTyp := val.Type().Elem()

	switch stTyp.Kind() {
	case reflect.Struct: // continue
	case reflect.Pointer:
		return ErrSlicePtr
	default:
		return ErrExpectedSlice
	}

	if length := val.Len(); length > 0 {
		for i := range length {
			item := val.Index(i)
			if item.Kind() == reflect.Pointer && item.IsNil() {
				return ErrExpectedPointer
			}

			if stTyp != val.Index(i).Type() {
				return ErrMixedSliceElementTypes
			}
		}
	}

	return nil
}
