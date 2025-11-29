package schema

import (
	"reflect"

	"github.com/laacin/inyorm/internal/core"
)

type SchemaType int

var ColumnIface = reflect.TypeOf((*core.Column)(nil)).Elem()

func whichIs(t reflect.Type) SchemaType {
	if t.Implements(ColumnIface) {
		return TypeColumn
	}

	knd := t.Kind()

	if knd == reflect.Struct {
		return TypeStruct
	}

	if knd == reflect.Interface {
		return TypeAny
	}

	if knd == reflect.Map {
		return TypeMap
	}

	if knd == reflect.String {
		return TypeString
	}

	if knd == reflect.Bool {
		return TypeBool
	}

	if knd >= reflect.Int && knd <= reflect.Int64 {
		return TypeInt
	}

	if knd >= reflect.Uint && knd <= reflect.Uint64 {
		return TypeUint
	}

	if knd == reflect.Float32 || knd == reflect.Float64 {
		return TypeFloat
	}

	return TypeUnknown
}
