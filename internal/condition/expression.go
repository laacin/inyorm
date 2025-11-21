package condition

type expression struct {
	identifier any
	operator   string
	negated    bool
	values     []any
	closed     bool
}

func (e *expression) addZero(op string) {
	e.closed = true
	e.operator = op
}

func (e *expression) addOne(op string, value any) {
	e.closed = true
	e.operator = op
	e.values = []any{value}
}

func (e *expression) addTwo(op string, val1, val2 any) {
	e.closed = true
	e.operator = op
	e.values = []any{val1, val2}
}

func (e *expression) addMany(op string, vals []any) {
	e.closed = true
	e.operator = op
	e.values = vals
}
