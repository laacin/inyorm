package condition

import "github.com/laacin/inyorm/internal/core"

func (c *Condition) Build(w core.Writer, opts core.WriterOpts) {
	var identOpt, valOpt core.WriterOpts
	identOpt = core.WriterOpts{ColType: opts.ColType}
	valOpt = core.WriterOpts{Placeholder: opts.Placeholder}

	w.Char('(')
	for i, expr := range c.Exprs {
		if i > 0 {
			w.Char(' ')
			w.Write(c.Connectors[i-1])
			w.Char(' ')
		}
		w.Value(expr.identifier, identOpt)
		w.Char(' ')
		w.Write(getOp(expr.operator, expr.negated))
		switch expr.operator {
		case isNull:
		case between:
			w.Char(' ')
			w.Value(expr.values[0], valOpt)
			w.Char(' ')
			w.Write(and)
			w.Char(' ')
			w.Value(expr.values[1], valOpt)
		case in:
			w.Char(' ')
			w.Char('(')
			for i, v := range expr.values {
				if i > 0 {
					w.Write(", ")
				}
				w.Value(v, valOpt)
			}
			w.Char(')')

		default:
			w.Char(' ')
			w.Value(expr.values[0], valOpt)
		}
	}
	w.Char(')')
}
