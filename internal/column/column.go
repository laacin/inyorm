package column

import "github.com/laacin/inyorm/internal/core"

type Column struct {
	table   string
	base    string
	value   string
	alias   string
	builder columnBuilder
}

func (c *Column) Def(w core.Writer) {
	wDef(w, c)
}

func (c *Column) Expr(w core.Writer) {
	wExpr(w, c)
}

func (c *Column) Alias(w core.Writer) {
	wAlias(w, c)
}

func (c *Column) Base(w core.Writer) {
	w.Column(c.table, c.base)
}

// -- Aggregation

func (c *Column) Count(distinct bool) { c.builder.wAggr(distinct, countAggr) }
func (c *Column) Sum(distinct bool)   { c.builder.wAggr(distinct, sumAggr) }
func (c *Column) Min(distinct bool)   { c.builder.wAggr(distinct, minAggr) }
func (c *Column) Max(distinct bool)   { c.builder.wAggr(distinct, maxAggr) }
func (c *Column) Avg(distinct bool)   { c.builder.wAggr(distinct, avgAggr) }

// -- Operation

func (c *Column) Add(v any) { c.builder.wOp(addOp, v) }
func (c *Column) Sub(v any) { c.builder.wOp(subOp, v) }
func (c *Column) Mul(v any) { c.builder.wOp(mulOp, v) }
func (c *Column) Div(v any) { c.builder.wOp(divOp, v) }
func (c *Column) Mod(v any) { c.builder.wOp(modOp, v) }
func (c *Column) Wrap()     { c.builder.wWrap() }

// -- Scalar

func (c *Column) Lower() { c.builder.wFunc(lowerFunc) }
func (c *Column) Upper() { c.builder.wFunc(upperFunc) }
func (c *Column) Trim()  { c.builder.wFunc(trimFunc) }
func (c *Column) Round() { c.builder.wFunc(roundFunc) }
func (c *Column) Abs()   { c.builder.wFunc(absFunc) }

// -- Alias

func (c *Column) As(value string) { c.builder.wAs(value) }
