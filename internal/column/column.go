package column

import "github.com/laacin/inyorm/internal/core"

type Column struct {
	writer core.Writer
	custom bool
	base   string
	table  string
	alias  string
	expr   string
}

func (c *Column) Def() core.Builder {
	return func(w core.Writer) {
		wDef(w, c)
	}
}

func (c *Column) Expr() core.Builder {
	return func(w core.Writer) {
		wExpr(w, c)
	}
}

func (c *Column) Alias() core.Builder {
	return func(w core.Writer) {
		wAlias(w, c)
	}
}

func (c *Column) Base() core.Builder {
	return func(w core.Writer) {
		wBase(w, c)
	}
}

// -- Aggregation

func (c *Column) Count(distinct bool) { wAggr(c, distinct, countAggr) }
func (c *Column) Sum(distinct bool)   { wAggr(c, distinct, sumAggr) }
func (c *Column) Min(distinct bool)   { wAggr(c, distinct, minAggr) }
func (c *Column) Max(distinct bool)   { wAggr(c, distinct, maxAggr) }
func (c *Column) Avg(distinct bool)   { wAggr(c, distinct, avgAggr) }

// -- Operation

func (c *Column) Add(v any) { wOp(c, addOp, v) }
func (c *Column) Sub(v any) { wOp(c, subOp, v) }
func (c *Column) Mul(v any) { wOp(c, mulOp, v) }
func (c *Column) Div(v any) { wOp(c, divOp, v) }
func (c *Column) Mod(v any) { wOp(c, modOp, v) }
func (c *Column) Wrap()     { wWrap(c) }

// -- Scalar

func (c *Column) Lower() { wFunc(c, lowerFunc) }
func (c *Column) Upper() { wFunc(c, upperFunc) }
func (c *Column) Trim()  { wFunc(c, trimFunc) }
func (c *Column) Round() { wFunc(c, roundFunc) }
func (c *Column) Abs()   { wFunc(c, absFunc) }

// -- Alias

func (c *Column) As(value string) { wAs(c, value) }
