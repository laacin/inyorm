package column

import "github.com/laacin/inyorm/internal/core"

type columnBuilder[Col, Value any] struct {
	exprs       []core.Builder
	aggregation core.Builder
	alias       string
}

func (ch *columnBuilder[Col, Value]) first(w core.Writer, c *Column[Col, Value]) {
	if c.value == "" {
		w.Column(c.Table, c.BaseName)
		return
	}
	w.Write(c.value)
}

func (ch *columnBuilder[Col, Value]) build(w core.Writer, c *Column[Col, Value]) {
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

func (cb *columnBuilder[Col, Value]) WExpr(expr core.Builder) {
	cb.exprs = append(cb.exprs, expr)
}

func (cb *columnBuilder[Col, Value]) wOp(arg byte, value any) {
	expr := func(w core.Writer) {
		prev := w.ToString()
		w.Reset()

		w.Write(prev)
		w.Char(' ')
		w.Char(arg)
		w.Char(' ')
		inferColumn[Col, Value](w, value)
	}
	cb.exprs = append(cb.exprs, expr)
}

func (cb *columnBuilder[Col, Value]) wFunc(arg string) {
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

func (cb *columnBuilder[Col, Value]) wWrap() {
	expr := func(w core.Writer) {
		prev := w.ToString()
		w.Reset()

		w.Char('(')
		w.Write(prev)
		w.Char(')')
	}
	cb.exprs = append(cb.exprs, expr)
}

func (cb *columnBuilder[Col, Value]) wAggr(distinct bool, aggr string) {
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

func (cb *columnBuilder[Col, Value]) wAs(name string) {
	cb.alias = name
}
