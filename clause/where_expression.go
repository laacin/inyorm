package clause

type WhereStart struct {
	ph       *PlaceholderGen
	negated  bool
	tables   []string
	fields   []string
	operator []string
	values   []any
	ended    *WhereEnd
}

func (w *WhereStart) start(field string, from string) *WhereStart {
	w.fields = append(w.fields, field)
	w.tables = append(w.tables, from)
	w.negated = false
	return w
}

func (w *WhereStart) end() *WhereEnd {
	if w.ended == nil {
		w.ended = &WhereEnd{arg: w}
	}
	return w.ended
}

func (w *WhereStart) Not() *WhereStart {
	w.negated = !w.negated
	return w
}

func (w *WhereStart) Equal(value any) *WhereEnd {
	w.operator = append(w.operator, opEqual(w.ph, w.negated))
	w.values = append(w.values, value)
	return w.end()
}

func (w *WhereStart) Like(value any) *WhereEnd {
	w.operator = append(w.operator, opLike(w.ph, w.negated))
	w.values = append(w.values, value)
	return w.end()
}

func (w *WhereStart) In(values ...any) *WhereEnd {
	w.operator = append(w.operator, opIn(w.ph, w.negated, len(values)))
	w.values = append(w.values, values...)
	return w.end()
}

func (w *WhereStart) Between(minV, maxV any) *WhereEnd {
	w.operator = append(w.operator, opBetween(w.ph, w.negated))
	w.values = append(w.values, minV, maxV)
	return w.end()
}

func (w *WhereStart) Greater(value any) *WhereEnd {
	w.operator = append(w.operator, opGreater(w.ph, w.negated))
	w.values = append(w.values, value)
	return w.end()
}

func (w *WhereStart) Less(value any) *WhereEnd {
	w.operator = append(w.operator, opLess(w.ph, w.negated))
	w.values = append(w.values, value)
	return w.end()
}

func (w *WhereStart) IsNull() *WhereEnd {
	w.operator = append(w.operator, opIsNull(w.negated))
	return w.end()
}

// ---- END EXPRESSION
type WhereEnd struct {
	arg        *WhereStart
	connectors []string
}

// Or starts a new condition with logical OR.
// - If no arguments are provided, it reuses the previous field from the clouse.
// - If two arguments are provided, the first is treated as the field name and the second as its alias.
func (w *WhereEnd) Or(fieldFrom ...string) *WhereStart {
	return w.nextCondition("OR", fieldFrom)
}

// And starts a new condition with logical AND.
// - If no arguments are provided, it reuses the previous field from the clouse.
// - If two arguments are provided, the first is treated as the field name and the second as its alias.
func (w *WhereEnd) And(fieldFrom ...string) *WhereStart {
	return w.nextCondition("AND", fieldFrom)
}

func (w *WhereEnd) nextCondition(logical string, fieldFrom []string) *WhereStart {
	w.connectors = append(w.connectors, logical)
	var field, table string

	if ln := len(fieldFrom); ln > 1 {
		field = fieldFrom[0]
		table = fieldFrom[1]
	} else if ln > 0 {
		field = fieldFrom[0]
		table = w.arg.tables[len(w.arg.tables)-1]
	} else {
		field = w.arg.fields[len(w.arg.fields)-1]
		table = w.arg.tables[len(w.arg.tables)-1]
	}

	return w.arg.start(field, table)
}
