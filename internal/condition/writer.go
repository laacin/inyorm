package condition

import "github.com/laacin/inyorm/internal/core"

func (e *Condition) Build(w core.Writer, opts *core.ValueOpts) {
	var identOpt, valOpt *core.ValueOpts
	if opts != nil {
		identOpt = &core.ValueOpts{Definition: opts.Definition}
		valOpt = &core.ValueOpts{Placeholder: opts.Placeholder}
	}

	w.Char('(')
	for i, seg := range e.segments {
		if i > 0 {
			w.Char(' ')
			w.Write(e.Next.connectors[i-1])
			w.Char(' ')
		}
		w.Value(seg.Identifier, identOpt)
		w.Char(' ')
		w.Write(getOp(seg.Operator, seg.Negated))
		switch seg.Operator {
		case isNull:
		case between:
			w.Char(' ')
			w.Value(seg.Argument[0], valOpt)
			w.Char(' ')
			w.Write(string(And))
			w.Char(' ')
			w.Value(seg.Argument[1], valOpt)
		case in:
			w.Char(' ')
			w.Char('(')
			for i, v := range seg.Argument {
				if i > 0 {
					w.Write(", ")
				}
				w.Value(v, valOpt)
			}
			w.Char(')')

		default:
			w.Char(' ')
			w.Value(seg.Argument[0], valOpt)
		}
	}
	w.Char(')')
}

func (e *ConditionNext) Build(w core.Writer, opts *core.ValueOpts) {
	e.ctx.Build(w, opts)
}
