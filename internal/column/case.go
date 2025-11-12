package column

type Case struct {
	exprs   []*expression
	current *expression
	els     any
}

func (c *Case) When(v any) {
	expr := &expression{when: v}
	c.exprs = append(c.exprs, expr)
	c.current = expr
}
func (c *Case) Then(v any) { c.current.do = v }
func (c *Case) Else(v any) { c.els = v }

type expression struct {
	when any
	do   any
}
