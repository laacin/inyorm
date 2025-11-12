package inyorm

// Value is a versatile type representing any element that can be used in a SQL expression
// within the ORM. It is essentially an alias for `any`, but semantically it signals
// that the value is intended to be interpreted in a SQL context.
//
// Value can be:
//   - SQL literals: strings, numbers, booleans, or nil. These will be rendered as
//     valid SQL constants, e.g., "example" → 'example', 0 → 0, true → 1, nil → NULL.
//   - ORM fields or columns, via (*ColumnExpr) methods. This ensures the value is treated
//     as a column reference rather than a literal.
type Value = any

func vMany(v []any) []any {
	newSlc := make([]any, len(v))
	for i, val := range v {
		newSlc[i] = vOne(val)
	}
	return newSlc
}

func vOne(v any) any {
	if col, ok := v.(*Column); ok {
		return col.wrap
	}

	if cond, ok := v.(*CondNext); ok {
		return cond.wrap
	}
	return v
}
