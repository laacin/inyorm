package clause

import "strings"

type WhereBuilder struct {
	expressions []*WhereStart
	table       string
	ph          *PlaceholderGen
}

func (w *WhereBuilder) Where(field string, from ...string) *WhereStart {
	if len(from) > 0 {
		w.table = from[0]
	}

	expr := &WhereStart{ph: w.ph}
	w.expressions = append(w.expressions, expr)
	return expr.start(field, w.table)
}

func (w *WhereBuilder) Build(sb *strings.Builder) []any {
	var values []any

	for n, expr := range w.expressions {
		values = append(values, expr.values...)
		if n > 0 {
			sb.WriteString(" AND ")
		}

		sb.WriteByte('(')
		for i := range expr.tables {
			if i > 0 {
				sb.WriteByte(' ')
				sb.WriteString(expr.ended.connectors[i-1])
				sb.WriteByte(' ')
			}
			sb.WriteString(string(expr.tables[i]))
			sb.WriteByte('.')
			sb.WriteString(string(expr.fields[i]))
			sb.WriteByte(' ')
			sb.WriteString(expr.operator[i])
		}
		sb.WriteByte(')')
	}

	return values
}
