package condition

import "github.com/laacin/inyorm/internal/core"

func (e *Condition) Build(w core.Writer, placeholder bool) {
	w.Char('(')
	for i, seg := range e.segments {
		if i > 0 {
			w.Char(' ')
			w.Write(e.Next.connectors[i-1])
			w.Char(' ')
		}
		w.Value(seg.Identifier, core.ValueOpts{})
		w.Char(' ')
		w.Write(getOp(seg.Operator, seg.Negated))
		opts := core.ValueOpts{Placeholder: placeholder}
		switch seg.Operator {
		case isNull:
		case between:
			w.Char(' ')
			w.Value(seg.Argument[0], opts)
			w.Char(' ')
			w.Write(string(And))
			w.Char(' ')
			w.Value(seg.Argument[1], opts)
		case in:
			w.Char(' ')
			w.Char('(')
			for i, v := range seg.Argument {
				if i > 0 {
					w.Write(", ")
				}
				w.Value(v, opts)
			}
			w.Char(')')

		default:
			w.Char(' ')
			w.Value(seg.Argument[0], opts)
		}
	}
	w.Char(')')
}

func (e *ConditionNext) Build(w core.Writer, placeholder bool) {
	e.ctx.Build(w, placeholder)
}
