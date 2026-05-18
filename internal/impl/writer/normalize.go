package writer

import "github.com/laacin/inyorm/internal/ir/expr"

func Normalize(value any) expr.ExprBuilder {
	if value == nil {
		return &expr.Null{}
	}

	if v, ok := value.(expr.ExprBuilder); ok {
		return v
	}

	switch v := value.(type) {
	case string:
		return expr.String(v)

	case *string:
		if v == nil {
			return &expr.Null{}
		}
		return expr.String(*v)

	case int:
		return expr.Number(v)
	case int8:
		return expr.Number(int(v))
	case int16:
		return expr.Number(int(v))
	case int32:
		return expr.Number(int(v))
	case int64:
		return expr.Number(int(v))

	case uint:
		return expr.Number(int(v))
	case uint8:
		return expr.Number(int(v))
	case uint16:
		return expr.Number(int(v))
	case uint32:
		return expr.Number(int(v))
	case uint64:
		return expr.Number(int(v))

	case *int:
		if v == nil {
			return &expr.Null{}
		}
		return expr.Number(*v)

	case *int8:
		if v == nil {
			return &expr.Null{}
		}
		return expr.Number(int(*v))

	case *int16:
		if v == nil {
			return &expr.Null{}
		}
		return expr.Number(int(*v))

	case *int32:
		if v == nil {
			return &expr.Null{}
		}
		return expr.Number(int(*v))

	case *int64:
		if v == nil {
			return &expr.Null{}
		}
		return expr.Number(int(*v))

	case float32:
		return expr.Float(float64(v))

	case float64:
		return expr.Float(v)

	case *float32:
		if v == nil {
			return &expr.Null{}
		}
		return expr.Float(float64(*v))

	case *float64:
		if v == nil {
			return &expr.Null{}
		}
		return expr.Float(*v)

	case bool:
		return expr.Bool(v)

	case *bool:
		if v == nil {
			return &expr.Null{}
		}
		return expr.Bool(*v)

	default:
		return &expr.Null{}
	}
}
