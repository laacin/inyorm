package expr

func Parse(value any) Expr {
	if value == nil {
		return Null{}
	}

	if expr, ok := value.(Expr); ok {
		return expr
	}

	switch v := value.(type) {
	case string:
		return String(v)

	case *string:
		if v == nil {
			return Null{}
		}
		return String(*v)

	case int:
		return Int(v)
	case int8:
		return Int(int(v))
	case int16:
		return Int(int(v))
	case int32:
		return Int(int(v))
	case int64:
		return Int(int(v))

	case uint:
		return Int(int(v))
	case uint8:
		return Int(int(v))
	case uint16:
		return Int(int(v))
	case uint32:
		return Int(int(v))
	case uint64:
		return Int(int(v))

	case *int:
		if v == nil {
			return Null{}
		}
		return Int(*v)

	case *int8:
		if v == nil {
			return Null{}
		}
		return Int(int(*v))

	case *int16:
		if v == nil {
			return Null{}
		}
		return Int(int(*v))

	case *int32:
		if v == nil {
			return Null{}
		}
		return Int(int(*v))

	case *int64:
		if v == nil {
			return Null{}
		}
		return Int(int(*v))

	case float32:
		return Float(float64(v))

	case float64:
		return Float(v)

	case *float32:
		if v == nil {
			return Null{}
		}
		return Float(float64(*v))

	case *float64:
		if v == nil {
			return Null{}
		}
		return Float(*v)

	case bool:
		return Bool(v)

	case *bool:
		if v == nil {
			return Null{}
		}
		return Bool(*v)

	default:
		return Null{}
	}
}
