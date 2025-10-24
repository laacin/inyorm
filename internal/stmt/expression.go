package stmt

import "strings"

type Expression struct {
	Ph       *PlaceholderGen
	Negated  bool
	Segments []*ExprSegment
	Current  *ExprSegment
	End      *ExpressionEnd
	Sb       strings.Builder
}

func (w *Expression) Start(identifier string) inyorm.Expression {
	if w.End == nil {
		w.End = &ExpressionEnd{ctx: w}
	}

	segment := &ExprSegment{Identifier: identifier}
	w.Current = segment
	w.Segments = append(w.Segments, segment)
	w.Negated = false
	return w
}

func (w *Expression) Not() inyorm.Expression {
	w.Negated = !w.Negated
	return w
}

func (w *Expression) Equal(value any) inyorm.ExpressionEnd {
	w.Current.Argument = w.writeOp(Equal, value)
	return w.End
}

func (w *Expression) Like(value any) inyorm.ExpressionEnd {
	w.Current.Argument = w.writeOp(Like, value)
	return w.End
}

func (w *Expression) In(values ...any) inyorm.ExpressionEnd {
	w.Current.Argument = w.writeOp(In, values...)
	return w.End
}

func (w *Expression) Between(minV, maxV any) inyorm.ExpressionEnd {
	w.Current.Argument = w.writeOp(Between, minV, maxV)
	return w.End
}

func (w *Expression) Greater(value any) inyorm.ExpressionEnd {
	w.Current.Argument = w.writeOp(Greater, value)
	return w.End
}

func (w *Expression) Less(value any) inyorm.ExpressionEnd {
	w.Current.Argument = w.writeOp(Less, value)
	return w.End
}

func (w *Expression) IsNull() inyorm.ExpressionEnd {
	w.Current.Argument = w.writeOp(IsNull)
	return w.End
}

// -- END EXPRESSION
type ExpressionEnd struct {
	ctx        *Expression
	Connectors []string
}

func (w *ExpressionEnd) Or(identifier ...string) inyorm.Expression {
	return w.nextCondition(Or, identifier)
}

func (w *ExpressionEnd) And(identifier ...string) inyorm.Expression {
	return w.nextCondition(And, identifier)
}

func (w *ExpressionEnd) nextCondition(logical Operator, identifier []string) inyorm.Expression {
	w.Connectors = append(w.Connectors, string(logical))

	var ident string
	if ln := len(identifier); ln > 0 {
		ident = identifier[0]
	} else {
		ident = w.ctx.Current.Identifier
	}

	return w.ctx.Start(ident)
}

// ----- Helpers -----
type ExprSegment struct {
	Identifier string
	Argument   string
}

func (w *Expression) writeOp(kind Operator, values ...any) string {
	w.Sb.Reset()
	w.Sb.WriteString(GetOp(kind, w.Negated))
	if len(values) == 0 {
		return w.Sb.String()
	}

	w.Sb.WriteByte(' ')
	if kind == Between {
		w.Ph.Write(&w.Sb, values[0])
		w.Sb.WriteString(" AND ")
		w.Ph.Write(&w.Sb, values[1])
	} else {
		w.Ph.Write(&w.Sb, values...)
	}
	return w.Sb.String()
}
