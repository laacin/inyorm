package writer

import "github.com/laacin/inyorm/intr/dialect"

func Normalize(value any) any {

	switch v := value.(type) {

	case string, dialect.Param, dialect.Column, dialect.Table:
		return v
	case *string:
		if v == nil {
			return nil
		}
		return *v
	case *dialect.Param:
		if v == nil {
			return nil
		}
		return *v
	case *dialect.Column:
		if v == nil {
			return nil
		}
		return *v
	case *dialect.Table:
		if v == nil {
			return nil
		}
		return *v

	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return v

	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)

	case *int:
		if v == nil {
			return nil
		}
		return *v

	case *int8:
		if v == nil {
			return nil
		}
		return *v

	case *int16:
		if v == nil {
			return nil
		}
		return int(*v)

	case *int32:
		if v == nil {
			return nil
		}
		return int(*v)

	case *int64:
		if v == nil {
			return nil
		}
		return *v

	case float32:
		return float64(v)
	case float64:
		return v

	case *float32:
		if v == nil {
			return nil
		}
		return float64(*v)

	case *float64:
		if v == nil {
			return nil
		}
		return *v

	case bool:
		return v

	case *bool:
		if v == nil {
			return nil
		}
		return *v

	default:
		return nil
	}
}
