package mapper

import (
	"reflect"

	"github.com/laacin/inyorm/internal/core"
)

const (
	typeUnknown = iota
	typeString
	typeInt
	typeUint
	typeFloat
	typeBool
	typeMap
	typeColumn
	typeStruct
	typeAny
)

var ColumnIface = reflect.TypeOf((*core.Column)(nil)).Elem()

func whichIs(t reflect.Type) int {
	if t.Implements(ColumnIface) {
		return typeColumn
	}

	knd := t.Kind()

	if knd == reflect.Struct {
		return typeStruct
	}

	if knd == reflect.Interface {
		return typeAny
	}

	if knd == reflect.Map {
		return typeMap
	}

	if knd == reflect.String {
		return typeString
	}

	if knd == reflect.Bool {
		return typeBool
	}

	if knd >= reflect.Int && knd <= reflect.Int64 {
		return typeInt
	}

	if knd >= reflect.Uint && knd <= reflect.Uint64 {
		return typeUint
	}

	if knd == reflect.Float32 || knd == reflect.Float64 {
		return typeFloat
	}

	return typeUnknown
}
