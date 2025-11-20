package column

import "github.com/laacin/inyorm/internal/core"

type columnBuilder[Col any] struct {
	exprs       []core.Builder
	aggregation core.Builder
	alias       string
}

func (ch *columnBuilder[Col]) first(w core.Writer, c *Column[Col]) {
	if c.value == "" {
		w.Column(c.Table, c.BaseName)
		return
	}
	w.Write(c.value)
}

func (ch *columnBuilder[Col]) build(w core.Writer, c *Column[Col]) {
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

func (cb *columnBuilder[Col]) WExpr(expr core.Builder) {
	cb.exprs = append(cb.exprs, expr)
}

func (cb *columnBuilder[Col]) wOp(arg byte, value any) {
	expr := func(w core.Writer) {
		prev := w.ToString()
		w.Reset()

		w.Write(prev)
		w.Char(' ')
		w.Char(arg)
		w.Char(' ')
		inferColumn[Col](w, value)
	}
	cb.exprs = append(cb.exprs, expr)
}

func (cb *columnBuilder[Col]) wFunc(arg string) {
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

func (cb *columnBuilder[Col]) wWrap() {
	expr := func(w core.Writer) {
		prev := w.ToString()
		w.Reset()

		w.Char('(')
		w.Write(prev)
		w.Char(')')
	}
	cb.exprs = append(cb.exprs, expr)
}

func (cb *columnBuilder[Col]) wAggr(distinct bool, aggr string) {
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

func (cb *columnBuilder[Col]) wAs(name string) {
	cb.alias = name
}
