package expression

import (
	"strings"

	"github.com/laacin/inyorm/internal"
)

type ExprStart struct {
	Ph       *internal.PlaceholderGen
	Negated  bool
	Segments []*ExprSegment
	Current  *ExprSegment
	End      *ExprEnd
	Sb       strings.Builder
}

func (w *ExprStart) Start(identifier, reference string) ExpressionStart {
	if w.End == nil {
		w.End = &ExprEnd{ctx: w}
	}

	segment := &ExprSegment{Identifier: identifier, Reference: reference}
	w.Current = segment
	w.Segments = append(w.Segments, segment)
	w.Negated = false
	return w
}

func (w *ExprStart) Not() ExpressionStart {
	w.Negated = !w.Negated
	return w
}

func (w *ExprStart) Equal(value any) ExpressionEnd {
	w.Current.Argument = w.writeOp(internal.Equal, value)
	return w.End
}

func (w *ExprStart) Like(value any) ExpressionEnd {
	w.Current.Argument = w.writeOp(internal.Like, value)
	return w.End
}

func (w *ExprStart) In(values ...any) ExpressionEnd {
	w.Current.Argument = w.writeOp(internal.In, values...)
	return w.End
}

func (w *ExprStart) Between(minV, maxV any) ExpressionEnd {
	w.Current.Argument = w.writeOp(internal.Between, minV, maxV)
	return w.End
}

func (w *ExprStart) Greater(value any) ExpressionEnd {
	w.Current.Argument = w.writeOp(internal.Greater, value)
	return w.End
}

func (w *ExprStart) Less(value any) ExpressionEnd {
	w.Current.Argument = w.writeOp(internal.Less, value)
	return w.End
}

func (w *ExprStart) IsNull() ExpressionEnd {
	w.Current.Argument = w.writeOp(internal.IsNull)
	return w.End
}

// -- END EXPRESSION
type ExprEnd struct {
	ctx        *ExprStart
	Connectors []string
}

func (w *ExprEnd) Or(identifier ...string) ExpressionStart {
	return w.nextCondition(internal.Or, identifier)
}

func (w *ExprEnd) And(identifier ...string) ExpressionStart {
	return w.nextCondition(internal.And, identifier)
}

func (w *ExprEnd) nextCondition(logical internal.Operator, identifier []string) ExpressionStart {
	w.Connectors = append(w.Connectors, string(logical))
	prevIdent := w.ctx.Current.Identifier
	prevRef := w.ctx.Current.Reference

	var ident, ref string
	if ln := len(identifier); ln > 1 {
		ident = identifier[0]
		ref = identifier[1]
	} else if ln > 0 {
		ident = identifier[0]
		ref = prevRef
	} else {
		ident = prevIdent
		ref = prevRef
	}

	return w.ctx.Start(ident, ref)
}

// ----- Helpers -----
type ExprSegment struct {
	Identifier string
	Reference  string
	Argument   string
}

func (w *ExprStart) writeOp(kind internal.Operator, values ...any) string {
	w.Sb.Reset()
	w.Sb.WriteString(internal.GetOp(kind, w.Negated))
	if len(values) == 0 {
		return w.Sb.String()
	}

	w.Sb.WriteByte(' ')
	if kind == internal.Between {
		w.Ph.Write(&w.Sb, values[0])
		w.Sb.WriteString(" AND ")
		w.Ph.Write(&w.Sb, values[1])
	} else {
		w.Ph.Write(&w.Sb, values...)
	}
	return w.Sb.String()
}
