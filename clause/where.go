package clause

import (
	"strings"

	"github.com/laacin/inyorm/internal/stmt"
)

type WhereBuilder struct {
	Table       string
	Expressions []*stmt.Expression
	Ph          *stmt.PlaceholderGen
}

// func (w *WhereBuilder) Where(identifier string) iface.Expression {
// 	expr := &stmt.Expression{Ph: w.Ph}
// 	w.Expressions = append(w.Expressions, expr)
// 	return expr.Start(identifier)
// }

func (w *WhereBuilder) Build(sb *strings.Builder) []any {
	sb.WriteString("WHERE ")
	for n, expr := range w.Expressions {
		if n > 0 {
			sb.WriteByte(' ')
			sb.WriteString(string(stmt.And))
			sb.WriteByte(' ')
		}

		sb.WriteByte('(')
		for i, seg := range expr.Segments {
			if i > 0 {
				sb.WriteByte(' ')
				sb.WriteString(expr.End.Connectors[i-1])
				sb.WriteByte(' ')
			}

			sb.WriteString(stmt.SetColumn(seg.Identifier))
			sb.WriteByte(' ')
			sb.WriteString(seg.Argument)
		}
		sb.WriteByte(')')
	}

	return w.Ph.Values()
}
