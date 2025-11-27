package mapper

type typValue int

const (
	typReflect typValue = iota
	typPrimitive
	typMap
)

func resolve(t any) (typ typValue, ptr, slc bool) {
	if is, ptr, slc := isPrimitive(t); is {
		return typPrimitive, ptr, slc
	}

	if is, ptr, slc := isMap(t); is {
		return typMap, ptr, slc
	}

	return
}

func isMap(t any) (is, ptr, slc bool) {
	switch t.(type) {
	case map[string]any:
		return true, false, false

	case []map[string]any:
		return true, false, true

	case []*map[string]any:
		return true, true, true

	default:
		return
	}
}

func isPrimitive(t any) (is, ptr, slc bool) {
	switch t.(type) {

	case string,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64, bool:
		return true, false, false

	case *string,
		*int, *int8, *int16, *int32, *int64,
		*uint, *uint8, *uint16, *uint32, *uint64,
		*float32, *float64, *bool:
		return true, true, false

	case []string,
		[]int, []int8, []int16, []int32, []int64,
		[]uint, []uint8, []uint16, []uint32, []uint64,
		[]float32, []float64, []bool, []any:
		return true, false, true

	case *[]string,
		*[]int, *[]int8, *[]int16, *[]int32, *[]int64,
		*[]uint, *[]uint8, *[]uint16, *[]uint32, *[]uint64,
		*[]float32, *[]float64, *[]bool, *[]any:
		return true, true, true

	default:
		return
	}
}
