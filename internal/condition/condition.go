package condition

type Condition struct {
	Exprs      []*expression
	Connectors []string
	Current    *expression
}

func (c *Condition) Start(ident any) {
	expr := &expression{identifier: ident}
	c.Exprs = append(c.Exprs, expr)
	c.Current = expr
}

func (c *Condition) Not()                   { c.Current.negated = !c.Current.negated }
func (c *Condition) Equal(value any)        { c.Current.addOne(equal, value) }
func (c *Condition) Like(value any)         { c.Current.addOne(like, value) }
func (c *Condition) In(values []any)        { c.Current.addMany(in, values) }
func (c *Condition) Between(val1, val2 any) { c.Current.addTwo(between, val1, val2) }
func (c *Condition) Greater(value any)      { c.Current.addOne(greater, value) }
func (c *Condition) Less(value any)         { c.Current.addOne(less, value) }
func (c *Condition) IsNull()                { c.Current.addZero(isNull) }

func (c *Condition) And(ident any) {
	c.Connectors = append(c.Connectors, and)
	c.Start(ident)
}
func (c *Condition) Or(ident any) {
	c.Connectors = append(c.Connectors, or)
	c.Start(ident)
}
