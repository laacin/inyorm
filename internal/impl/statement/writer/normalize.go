package writer

import "github.com/laacin/inyorm/internal/entity/dml"

func Normalize(value any) dml.Value {
	if value == nil {
		return &dml.Null{}
	}

	if v, ok := value.(dml.Value); ok {
		return v
	}

	if v, ok := value.(dml.ValueBuilder); ok {
		return v.Build()
	}

	switch v := value.(type) {
	case string:
		return dml.String(v)

	case *string:
		if v == nil {
			return &dml.Null{}
		}
		return dml.String(*v)

	case int:
		return dml.Number(v)
	case int8:
		return dml.Number(int(v))
	case int16:
		return dml.Number(int(v))
	case int32:
		return dml.Number(int(v))
	case int64:
		return dml.Number(int(v))

	case uint:
		return dml.Number(int(v))
	case uint8:
		return dml.Number(int(v))
	case uint16:
		return dml.Number(int(v))
	case uint32:
		return dml.Number(int(v))
	case uint64:
		return dml.Number(int(v))

	case *int:
		if v == nil {
			return &dml.Null{}
		}
		return dml.Number(*v)

	case *int8:
		if v == nil {
			return &dml.Null{}
		}
		return dml.Number(int(*v))

	case *int16:
		if v == nil {
			return &dml.Null{}
		}
		return dml.Number(int(*v))

	case *int32:
		if v == nil {
			return &dml.Null{}
		}
		return dml.Number(int(*v))

	case *int64:
		if v == nil {
			return &dml.Null{}
		}
		return dml.Number(int(*v))

	case float32:
		return dml.Float(float64(v))

	case float64:
		return dml.Float(v)

	case *float32:
		if v == nil {
			return &dml.Null{}
		}
		return dml.Float(float64(*v))

	case *float64:
		if v == nil {
			return &dml.Null{}
		}
		return dml.Float(*v)

	case bool:
		return dml.Bool(v)

	case *bool:
		if v == nil {
			return &dml.Null{}
		}
		return dml.Bool(*v)

	default:
		return &dml.Null{}
	}
}
