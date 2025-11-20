package column

type Case[Self, Next, Ident, Value any] struct {
	Exprs   []*caseExpr[Ident, Value]
	current *caseExpr[Ident, Value]
	Els     Value
}

func (c *Case[Self, Next, Ident, Value]) When(ident Ident) Next {
	expr := &caseExpr[Ident, Value]{Identifier: ident}
	c.Exprs = append(c.Exprs, expr)
	c.current = expr
	return any(c).(Next)
}

func (c *Case[Self, Next, Ident, Value]) Then(arg Value) Self {
	c.current.Argument = arg
	return any(c).(Self)
}

func (c *Case[Self, Next, Ident, Value]) Else(v Value) {
	c.Els = v
}

// -- Expression

type caseExpr[Ident, Value any] struct {
	Identifier Ident
	Argument   Value
}
