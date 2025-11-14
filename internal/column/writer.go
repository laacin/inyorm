package column

import "github.com/laacin/inyorm/internal/core"

func wExpr(w core.Writer, c *Column) {
	c.builder.build(w.Split(), c)

	if c.value == "" {
		w.Column(c.table, c.base)
		return
	}
	w.Write(c.value)
}

func wAlias(w core.Writer, c *Column) {
	c.builder.build(w.Split(), c)

	if c.alias != "" {
		w.Write(c.alias)
		return
	}

	if c.value == "" {
		w.Column(c.table, c.base)
		return
	}
	w.Write(c.value)
}

func wDef(w core.Writer, c *Column) {
	c.builder.build(w.Split(), c)

	if c.value == "" {
		w.Column(c.table, c.base)
		return
	}

	w.Write(c.value)
	if c.alias != "" {
		w.Write(" AS ")
		w.Write(c.alias)
	}
}

// -- internal writers

type columnBuilder struct {
	exprs       []core.Builder
	aggregation core.Builder
	alias       string
}

func (ch *columnBuilder) first(w core.Writer, c *Column) {
	if c.value == "" {
		w.Column(c.table, c.base)
		return
	}
	w.Write(c.value)
}

func (ch *columnBuilder) build(w core.Writer, c *Column) {
	if ch.exprs == nil && ch.aggregation == nil && ch.alias == "" {
		return
	}

	if ch.exprs != nil || ch.aggregation != nil {
		ch.first(w, c)
	}
	if ch.exprs != nil {
		for _, expr := range ch.exprs {
			expr(w)
		}
		ch.exprs = nil
	}

	if ch.aggregation != nil {
		ch.aggregation(w)
		ch.aggregation = nil
	}

	if ch.alias != "" {
		c.alias = ch.alias
		ch.alias = ""
	}

	c.value = w.ToString()
	w.Reset()
}

func (cb *columnBuilder) wExpr(expr core.Builder) {
	cb.exprs = append(cb.exprs, expr)
}

func (cb *columnBuilder) wOp(arg byte, value any) {
	expr := func(w core.Writer) {
		prev := w.ToString()
		w.Reset()

		w.Write(prev)
		w.Char(' ')
		w.Char(arg)
		w.Char(' ')
		w.Value(value, core.ColumnIdentWriteOpt)
	}
	cb.exprs = append(cb.exprs, expr)
}

func (cb *columnBuilder) wFunc(arg string) {
	expr := func(w core.Writer) {
		prev := w.ToString()
		w.Reset()

		w.Write(arg)
		w.Char('(')
		w.Write(prev)
		w.Char(')')
	}
	cb.exprs = append(cb.exprs, expr)
}

func (cb *columnBuilder) wWrap() {
	expr := func(w core.Writer) {
		prev := w.ToString()
		w.Reset()

		w.Char('(')
		w.Write(prev)
		w.Char(')')
	}
	cb.exprs = append(cb.exprs, expr)
}

func (cb *columnBuilder) wAggr(distinct bool, aggr string) {
	expr := func(w core.Writer) {
		prev := w.ToString()
		w.Reset()

		w.Write(aggr)
		w.Char('(')
		if distinct {
			w.Write("DISTINCT ")
		}
		w.Write(prev)
		w.Char(')')
	}
	cb.aggregation = expr
}

func (cb *columnBuilder) wAs(name string) {
	cb.alias = name
}
