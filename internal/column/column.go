package column

type Column[Self, Value any] struct {
	Table    string
	BaseName string
	alias    string
	value    string
	Builder  columnBuilder[Self, Value]
}

// -- Aggregation

func (c *Column[Self, Value]) Count(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Builder.wAggr(dist, countAggr)
	return any(c).(Self)
}

func (c *Column[Self, Value]) Sum(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Builder.wAggr(dist, sumAggr)
	return any(c).(Self)
}

func (c *Column[Self, Value]) Min(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Builder.wAggr(dist, minAggr)
	return any(c).(Self)
}

func (c *Column[Self, Value]) Max(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Builder.wAggr(dist, maxAggr)
	return any(c).(Self)
}

func (c *Column[Self, Value]) Avg(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Builder.wAggr(dist, avgAggr)
	return any(c).(Self)
}

// -- Operation

func (c *Column[Self, Value]) Add(v Value) Self {
	c.Builder.wOp(addOp, v)
	return any(c).(Self)
}

func (c *Column[Self, Value]) Sub(v Value) Self {
	c.Builder.wOp(subOp, v)
	return any(c).(Self)
}

func (c *Column[Self, Value]) Mul(v Value) Self {
	c.Builder.wOp(mulOp, v)
	return any(c).(Self)
}

func (c *Column[Self, Value]) Div(v Value) Self {
	c.Builder.wOp(divOp, v)
	return any(c).(Self)
}

func (c *Column[Self, Value]) Mod(v Value) Self {
	c.Builder.wOp(modOp, v)
	return any(c).(Self)
}

func (c *Column[Self, Value]) Wrap() Self {
	c.Builder.wWrap()
	return any(c).(Self)
}

// -- Scalar

func (c *Column[Self, Value]) Lower() Self {
	c.Builder.wFunc(lowerFunc)
	return any(c).(Self)
}

func (c *Column[Self, Value]) Upper() Self {
	c.Builder.wFunc(upperFunc)
	return any(c).(Self)
}

func (c *Column[Self, Value]) Trim() Self {
	c.Builder.wFunc(trimFunc)
	return any(c).(Self)
}

func (c *Column[Self, Value]) Round() Self {
	c.Builder.wFunc(roundFunc)
	return any(c).(Self)
}

func (c *Column[Self, Value]) Abs() Self {
	c.Builder.wFunc(absFunc)
	return any(c).(Self)
}

// -- Alias

func (c *Column[Self, Value]) As(value string) Self {
	c.Builder.wAs(value)
	return any(c).(Self)
}
