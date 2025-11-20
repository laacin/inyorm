package condition

type Condition[Self, Next any] struct {
	Exprs      []*expression
	Connectors []string
	Current    *expression
}

func (c *Condition[Self, Next]) Start(ident any) Self {
	expr := &expression{identifier: ident}
	c.Exprs = append(c.Exprs, expr)
	c.Current = expr
	return any(c).(Self)
}

func (c *Condition[Self, Next]) Not() Self {
	c.Current.negated = !c.Current.negated
	return any(c).(Self)
}

func (c *Condition[Self, Next]) Equal(value any) Next {
	c.Current.addOne(equal, value)
	return any(c).(Next)
}

func (c *Condition[Self, Next]) Like(value any) Next {
	c.Current.addOne(like, value)
	return any(c).(Next)
}

func (c *Condition[Self, Next]) In(values ...any) Next {
	c.Current.addMany(in, values)
	return any(c).(Next)
}

func (c *Condition[Self, Next]) Between(val1, val2 any) Next {
	c.Current.addTwo(between, val1, val2)
	return any(c).(Next)
}

func (c *Condition[Self, Next]) Greater(value any) Next {
	c.Current.addOne(greater, value)
	return any(c).(Next)
}

func (c *Condition[Self, Next]) Less(value any) Next {
	c.Current.addOne(less, value)
	return any(c).(Next)
}

func (c *Condition[Self, Next]) IsNull() Next {
	c.Current.addZero(isNull)
	return any(c).(Next)
}

func (c *Condition[Self, Next]) And(ident any) Self {
	c.Connectors = append(c.Connectors, and)
	return c.Start(ident)
}

func (c *Condition[Self, Next]) Or(ident any) Self {
	c.Connectors = append(c.Connectors, or)
	return c.Start(ident)
}
