package clause

import (
	"strings"

	"github.com/laacin/inyorm/expression"
	"github.com/laacin/inyorm/internal"
)

type WhereBuilder struct {
	Table       string
	Expressions []*expression.ExprStart
	Ph          *internal.PlaceholderGen
}

func (w *WhereBuilder) NewExpr(field string, table ...string) expression.ExpressionStart {
	tbl := w.Table
	if len(table) > 0 {
		tbl = table[0]
	}

	expr := &expression.ExprStart{Ph: w.Ph}
	w.Expressions = append(w.Expressions, expr)
	return expr.Start(field, tbl)
}

func (w *WhereBuilder) Build(sb *strings.Builder) []any {
	sb.WriteString("WHERE ")
	for n, expr := range w.Expressions {
		if n > 0 {
			sb.WriteByte(' ')
			sb.WriteString(string(internal.And))
			sb.WriteByte(' ')
		}

		sb.WriteByte('(')
		for i, seg := range expr.Segments {
			if i > 0 {
				sb.WriteByte(' ')
				sb.WriteString(expr.End.Connectors[i-1])
				sb.WriteByte(' ')
			}

			sb.WriteString(seg.Reference)
			sb.WriteByte('.')
			sb.WriteString(seg.Identifier)
			sb.WriteByte(' ')
			sb.WriteString(seg.Argument)
		}
		sb.WriteByte(')')
	}

	return w.Ph.Values()
}
