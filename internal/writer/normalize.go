package writer

import "github.com/laacin/inyorm/internal/entity"

func Normalize(value any) entity.Value {
	if value == nil {
		return &entity.Null{}
	}

	if v, ok := any(value).(entity.Value); ok {
		return v
	}

	if v, ok := any(value).(entity.ValueBuilder); ok {
		return v.Build()
	}

	switch v := value.(type) {
	case string:
		return entity.String(v)

	case *string:
		if v == nil {
			return &entity.Null{}
		}
		return entity.String(*v)

	case int:
		return entity.Number(v)
	case int8:
		return entity.Number(int(v))
	case int16:
		return entity.Number(int(v))
	case int32:
		return entity.Number(int(v))
	case int64:
		return entity.Number(int(v))

	case uint:
		return entity.Number(int(v))
	case uint8:
		return entity.Number(int(v))
	case uint16:
		return entity.Number(int(v))
	case uint32:
		return entity.Number(int(v))
	case uint64:
		return entity.Number(int(v))

	case *int:
		if v == nil {
			return &entity.Null{}
		}
		return entity.Number(*v)

	case *int8:
		if v == nil {
			return &entity.Null{}
		}
		return entity.Number(int(*v))

	case *int16:
		if v == nil {
			return &entity.Null{}
		}
		return entity.Number(int(*v))

	case *int32:
		if v == nil {
			return &entity.Null{}
		}
		return entity.Number(int(*v))

	case *int64:
		if v == nil {
			return &entity.Null{}
		}
		return entity.Number(int(*v))

	case float32:
		return entity.Float(float64(v))

	case float64:
		return entity.Float(v)

	case *float32:
		if v == nil {
			return &entity.Null{}
		}
		return entity.Float(float64(*v))

	case *float64:
		if v == nil {
			return &entity.Null{}
		}
		return entity.Float(*v)

	case bool:
		return entity.Bool(v)

	case *bool:
		if v == nil {
			return &entity.Null{}
		}
		return entity.Bool(*v)

	default:
		return &entity.Null{}
	}
}
