package column

type Case[Self, Next any] struct {
	Exprs   []*caseExpr
	current *caseExpr
	Els     any
}

func (c *Case[Self, Next]) When(ident any) Next {
	expr := &caseExpr{Identifier: ident}
	c.Exprs = append(c.Exprs, expr)
	c.current = expr
	return any(c).(Next)
}

func (c *Case[Self, Next]) Then(arg any) Self {
	c.current.Argument = arg
	return any(c).(Self)
}

func (c *Case[Self, Next]) Else(v any) {
	c.Els = v
}

// -- Expression

type caseExpr struct {
	Identifier any
	Argument   any
}
