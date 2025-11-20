package condition

type Condition[Self, Next, Ident, Value any] struct {
	Exprs      []*expression[Ident, Value]
	Connectors []string
	Current    *expression[Ident, Value]
}

func (c *Condition[Self, Next, Ident, Value]) Start(ident Ident) Self {
	expr := &expression[Ident, Value]{identifier: ident}
	c.Exprs = append(c.Exprs, expr)
	c.Current = expr
	return any(c).(Self)
}

func (c *Condition[Self, Next, Ident, Value]) Not() Self {
	c.Current.negated = !c.Current.negated
	return any(c).(Self)
}

func (c *Condition[Self, Next, Ident, Value]) Equal(value Value) Next {
	c.Current.addOne(equal, value)
	return any(c).(Next)
}

func (c *Condition[Self, Next, Ident, Value]) Like(value Value) Next {
	c.Current.addOne(like, value)
	return any(c).(Next)
}

func (c *Condition[Self, Next, Ident, Value]) In(values ...Value) Next {
	c.Current.addMany(in, values)
	return any(c).(Next)
}

func (c *Condition[Self, Next, Ident, Value]) Between(val1, val2 Value) Next {
	c.Current.addTwo(between, val1, val2)
	return any(c).(Next)
}

func (c *Condition[Self, Next, Ident, Value]) Greater(value Value) Next {
	c.Current.addOne(greater, value)
	return any(c).(Next)
}

func (c *Condition[Self, Next, Ident, Value]) Less(value Value) Next {
	c.Current.addOne(less, value)
	return any(c).(Next)
}

func (c *Condition[Self, Next, Ident, Value]) IsNull() Next {
	c.Current.addZero(isNull)
	return any(c).(Next)
}

func (c *Condition[Self, Next, Ident, Value]) And(ident Ident) Self {
	c.Connectors = append(c.Connectors, and)
	return c.Start(ident)
}

func (c *Condition[Self, Next, Ident, Value]) Or(ident Ident) Self {
	c.Connectors = append(c.Connectors, or)
	return c.Start(ident)
}
