package condition

import "github.com/laacin/inyorm/internal/core"

func (c *Condition[Self, Next]) Build(w core.Writer, ctx core.ClauseType) {
	w.Char('(')
	for i, expr := range c.Exprs {
		if !expr.closed {
			continue
		}

		if i > 0 {
			w.Char(' ')
			w.Write(c.Connectors[i-1])
			w.Char(' ')
		}

		w.Identifier(expr.identifier, ctx)
		w.Char(' ')
		w.Write(getOp(expr.operator, expr.negated))
		switch expr.operator {
		case isNull:
		case between:
			w.Char(' ')
			w.Value(expr.values[0], ctx)
			w.Char(' ')
			w.Write(and)
			w.Char(' ')
			w.Value(expr.values[1], ctx)
		case in:
			w.Char(' ')
			w.Char('(')
			for i, v := range expr.values {
				if i > 0 {
					w.Write(", ")
				}
				w.Value(v, ctx)
			}
			w.Char(')')

		default:
			w.Char(' ')
			w.Value(expr.values[0], ctx)
		}
	}
	w.Char(')')
}
