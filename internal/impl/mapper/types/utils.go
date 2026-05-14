package types

import "reflect"

func derefPtrTyp(t reflect.Type) (reflect.Type, int) {
	count := 0
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
		count++
	}
	return t, count
}

func derefPtrVal(v reflect.Value) (reflect.Value, int) {
	count := 0
	for v.IsValid() && v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
		count++
	}
	return v, count
}
