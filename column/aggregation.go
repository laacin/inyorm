package column

func (c *Column) Count() { wrapAggr(c, "COUNT") }
func (c *Column) Sum()   { wrapAggr(c, "SUM") }
func (c *Column) Avg()   { wrapAggr(c, "AVG") }
func (c *Column) Max()   { wrapAggr(c, "MAX") }
func (c *Column) Min()   { wrapAggr(c, "MIN") }

func wrapAggr(c *Column, aggr string) {
	*c = Column(aggr) + "(" + *c + ")"
}
