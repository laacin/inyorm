package column

type Column[Self any] struct {
	Table    string
	BaseName string
	alias    string
	value    string
	Builder  columnBuilder[Self]
}

// -- Aggregation

func (c *Column[Self]) Count(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Builder.wAggr(dist, countAggr)
	return any(c).(Self)
}

func (c *Column[Self]) Sum(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Builder.wAggr(dist, sumAggr)
	return any(c).(Self)
}

func (c *Column[Self]) Min(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Builder.wAggr(dist, minAggr)
	return any(c).(Self)
}

func (c *Column[Self]) Max(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Builder.wAggr(dist, maxAggr)
	return any(c).(Self)
}

func (c *Column[Self]) Avg(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Builder.wAggr(dist, avgAggr)
	return any(c).(Self)
}

// -- Operation

func (c *Column[Self]) Add(v any) Self {
	c.Builder.wOp(addOp, v)
	return any(c).(Self)
}

func (c *Column[Self]) Sub(v any) Self {
	c.Builder.wOp(subOp, v)
	return any(c).(Self)
}

func (c *Column[Self]) Mul(v any) Self {
	c.Builder.wOp(mulOp, v)
	return any(c).(Self)
}

func (c *Column[Self]) Div(v any) Self {
	c.Builder.wOp(divOp, v)
	return any(c).(Self)
}

func (c *Column[Self]) Mod(v any) Self {
	c.Builder.wOp(modOp, v)
	return any(c).(Self)
}

func (c *Column[Self]) Wrap() Self {
	c.Builder.wWrap()
	return any(c).(Self)
}

// -- Scalar

func (c *Column[Self]) Lower() Self {
	c.Builder.wFunc(lowerFunc)
	return any(c).(Self)
}

func (c *Column[Self]) Upper() Self {
	c.Builder.wFunc(upperFunc)
	return any(c).(Self)
}

func (c *Column[Self]) Trim() Self {
	c.Builder.wFunc(trimFunc)
	return any(c).(Self)
}

func (c *Column[Self]) Round() Self {
	c.Builder.wFunc(roundFunc)
	return any(c).(Self)
}

func (c *Column[Self]) Abs() Self {
	c.Builder.wFunc(absFunc)
	return any(c).(Self)
}

// -- Alias

func (c *Column[Self]) As(value string) Self {
	c.Builder.wAs(value)
	return any(c).(Self)
}
