package stmt

import "strconv"

// Stringify returns a string representation of basic comparable types without using reflection.
func Stringify(value any) string {
	switch t := value.(type) {
	case string:
		return FilterColumn(t)
	case *string:
		if t == nil {
			return "NULL"
		}
		return FilterColumn(*t)

	case int:
		return strconv.Itoa(t)
	case *int:
		if t == nil {
			return "NULL"
		}
		return strconv.Itoa(*t)

	case int8:
		return strconv.Itoa(int(t))
	case *int8:
		if t == nil {
			return "NULL"
		}
		return strconv.Itoa(int(*t))

	case int16:
		return strconv.Itoa(int(t))
	case *int16:
		if t == nil {
			return "NULL"
		}
		return strconv.Itoa(int(*t))

	case int32:
		return strconv.Itoa(int(t))
	case *int32:
		if t == nil {
			return "NULL"
		}
		return strconv.Itoa(int(*t))

	case int64:
		return strconv.FormatInt(t, 10)
	case *int64:
		if t == nil {
			return "NULL"
		}
		return strconv.FormatInt(*t, 10)

	case uint:
		return strconv.FormatUint(uint64(t), 10)
	case *uint:
		if t == nil {
			return "NULL"
		}
		return strconv.FormatUint(uint64(*t), 10)

	case uint8:
		return strconv.FormatUint(uint64(t), 10)
	case *uint8:
		if t == nil {
			return "NULL"
		}
		return strconv.FormatUint(uint64(*t), 10)

	case uint16:
		return strconv.FormatUint(uint64(t), 10)
	case *uint16:
		if t == nil {
			return "NULL"
		}
		return strconv.FormatUint(uint64(*t), 10)

	case uint32:
		return strconv.FormatUint(uint64(t), 10)
	case *uint32:
		if t == nil {
			return "NULL"
		}
		return strconv.FormatUint(uint64(*t), 10)

	case uint64:
		return strconv.FormatUint(t, 10)
	case *uint64:
		if t == nil {
			return "NULL"
		}
		return strconv.FormatUint(*t, 10)

	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 32)
	case *float32:
		if t == nil {
			return "NULL"
		}
		return strconv.FormatFloat(float64(*t), 'f', -1, 32)

	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case *float64:
		if t == nil {
			return "NULL"
		}
		return strconv.FormatFloat(*t, 'f', -1, 64)

	case bool:
		if t {
			return "1"
		}
		return "0"
	case *bool:
		if t == nil {
			return "NULL"
		}
		if *t {
			return "1"
		}
		return "0"
	}

	return "NULL"
}
