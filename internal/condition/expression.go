package condition

type expression[Ident, Value any] struct {
	identifier Ident
	operator   string
	negated    bool
	values     []Value
}

func (e *expression[Ident, Value]) addZero(op string) {
	e.operator = op
}

func (e *expression[Ident, Value]) addOne(op string, value Value) {
	e.operator = op
	e.values = []Value{value}
}

func (e *expression[Ident, Value]) addTwo(op string, val1, val2 Value) {
	e.operator = op
	e.values = []Value{val1, val2}
}

func (e *expression[Ident, Value]) addMany(op string, vals []Value) {
	e.operator = op
	e.values = vals
}
