package types

import "reflect"

func Unwrap[T any](v any, ptr bool) T {
	if ptr {
		return *v.(*T)
	}
	return v.(T)
}

func UnwrapSlc[T any](v any, ptr bool) []T {
	if ptr {
		return *v.(*[]T)
	}
	return v.([]T)
}

func DerefPtrTyp(t reflect.Type) (reflect.Type, int) {
	count := 0
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
		count++
	}
	return t, count
}

func DerefPtrVal(v reflect.Value) (reflect.Value, int) {
	count := 0
	for v.IsValid() && v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
		count++
	}
	return v, count
}
