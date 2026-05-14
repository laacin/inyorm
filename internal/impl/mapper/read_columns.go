package mapper

import "reflect"

func GetColumns(v any) ([]string, error) {
	info := ObtainInfo(reflect.TypeOf(v))
}
