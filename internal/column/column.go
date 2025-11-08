package column

import "github.com/laacin/inyorm/internal/core"

type Column struct {
	Writer      core.Writer
	Type        ColumnType
	Value       string
	Alias       string
	Aggregation func()
}

func (c *Column) Def() func(w core.Writer) {
	return func(w core.Writer) {
		if c.Aggregation != nil {
			c.Aggregation()
		}

		switch c.Type {
		case normalCol:
			if c.Alias != "" {
				w.ColRef(c.Alias)
				w.Char('.')
			}
			w.Write(c.Value)

		case customCol:
			w.Write(c.Value)
			if c.Alias != "" {
				w.Write(" AS ")
				w.Write(c.Alias)
			}
		}
	}
}

func (c *Column) Ref() func(w core.Writer) {
	return func(w core.Writer) {
		switch c.Type {
		case normalCol:
			if c.Alias != "" {
				w.ColRef(c.Alias)
				w.Char('.')
			}
			w.Write(c.Value)

		case customCol:
			if c.Alias == "" {
				w.Write(c.Value)
				return
			}
			w.Write(c.Alias)
		}
	}
}

// -- Aggregation

func (c *Column) Count(distinct ...bool) core.Column {
	dist := len(distinct) > 0 && distinct[0]
	wAggr(c, dist, countAggr)
	return c
}

func (c *Column) Sum(distinct ...bool) core.Column {
	dist := len(distinct) > 0 && distinct[0]
	wAggr(c, dist, sumAggr)
	return c
}

func (c *Column) Min(distinct ...bool) core.Column {
	dist := len(distinct) > 0 && distinct[0]
	wAggr(c, dist, minAggr)
	return c
}

func (c *Column) Max(distinct ...bool) core.Column {
	dist := len(distinct) > 0 && distinct[0]
	wAggr(c, dist, maxAggr)
	return c
}

func (c *Column) Avg(distinct ...bool) core.Column {
	dist := len(distinct) > 0 && distinct[0]
	wAggr(c, dist, avgAggr)
	return c
}

// -- Operation

func (c *Column) Add(v any) core.Column {
	wOp(c, addOp, v)
	return c
}

func (c *Column) Sub(v any) core.Column {
	wOp(c, subOp, v)
	return c
}

func (c *Column) Mul(v any) core.Column {
	wOp(c, mulOp, v)
	return c
}

func (c *Column) Div(v any) core.Column {
	wOp(c, divOp, v)
	return c
}

func (c *Column) Mod(v any) core.Column {
	wOp(c, modOp, v)
	return c
}

func (c *Column) Wrap() core.Column {
	wWrap(c)
	return c
}

// -- Scalar

func (c *Column) Lower() core.Column {
	wFunc(c, lowerFunc)
	return c
}

func (c *Column) Upper() core.Column {
	wFunc(c, upperFunc)
	return c
}

func (c *Column) Trim() core.Column {
	wFunc(c, trimFunc)
	return c
}

func (c *Column) Round() core.Column {
	wFunc(c, roundFunc)
	return c
}

func (c *Column) Abs() core.Column {
	wFunc(c, absFunc)
	return c
}

// -- Alias

func (c *Column) As(value string) core.Column {
	wAs(c, value)
	return c
}
