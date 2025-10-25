package stmt

import "strconv"

const (
	falseValue = "0"
	trueValue  = "1"
	nullValue  = "NULL"
)

type number interface {
	~int | ~int8 | ~int16 |
		~int32 | ~int64 | ~uint |
		~uint8 | ~uint16 | ~uint32 |
		~uint64
}

func handlePtr[T any](v *T, fn func(T) string) string {
	if v == nil {
		return nullValue
	}
	return fn(*v)
}

func handleInt[T number](v T) string {
	return strconv.Itoa(int(v))
}

func handleFloat[T ~float32 | ~float64](v T) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}

func handleBool(v bool) string {
	if !v {
		return falseValue
	}
	return trueValue
}

// Stringify returns a string representation of basic comparable types without using reflection.
func Stringify(value any) string {
	switch t := value.(type) {

	case string:
		return FilterColumn(t)
	case *string:
		return handlePtr(t, FilterColumn)

	case int:
		return handleInt(t)
	case *int:
		return handlePtr(t, handleInt)

	case int8:
		return handleInt(t)
	case *int8:
		return handlePtr(t, handleInt)

	case int16:
		return handleInt(t)
	case *int16:
		return handlePtr(t, handleInt)

	case int32:
		return handleInt(t)
	case *int32:
		return handlePtr(t, handleInt)

	case int64:
		return handleInt(t)
	case *int64:
		return handlePtr(t, handleInt)

	case uint:
		return handleInt(t)
	case *uint:
		return handlePtr(t, handleInt)

	case uint8:
		return handleInt(t)
	case *uint8:
		return handlePtr(t, handleInt)

	case uint16:
		return handleInt(t)
	case *uint16:
		return handlePtr(t, handleInt)

	case uint32:
		return handleInt(t)
	case *uint32:
		return handlePtr(t, handleInt)

	case uint64:
		return handleInt(t)
	case *uint64:
		return handlePtr(t, handleInt)

	case float32:
		return handleFloat(t)
	case *float32:
		return handlePtr(t, handleFloat)

	case float64:
		return handleFloat(t)
	case *float64:
		return handlePtr(t, handleFloat)

	case bool:
		return handleBool(t)
	case *bool:
		return handlePtr(t, handleBool)
	}

	return nullValue
}
