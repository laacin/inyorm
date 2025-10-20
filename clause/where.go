package clause

import "strings"

// ----- Builder -----

type WhereBuilder struct {
	Expressions []*WhereStart
	Table       string
	Ph          *PlaceholderGen
}

func (w *WhereBuilder) NewExpr(field string, table ...string) *WhereStart {
	tbl := w.Table
	if len(table) > 0 {
		tbl = table[0]
	}

	expr := &WhereStart{ph: w.Ph}
	w.Expressions = append(w.Expressions, expr)
	return expr.start(field, tbl)
}

func (w *WhereBuilder) Build(sb *strings.Builder) []any {
	var values []any

	sb.WriteString("WHERE ")
	for n, expr := range w.Expressions {
		values = append(values, expr.values...)
		if n > 0 {
			sb.WriteString(" AND ")
		}

		sb.WriteByte('(')
		for i, col := range expr.column {

			if i > 0 {
				sb.WriteByte(' ')
				sb.WriteString(expr.ended.connectors[i-1])
				sb.WriteByte(' ')
			}

			sb.WriteString(col.table)
			sb.WriteByte('.')
			sb.WriteString(col.field)
			sb.WriteByte(' ')
			sb.WriteString(expr.operator[i])
		}
		sb.WriteByte(')')
	}

	return values
}

// ----- Clause -----

type WhereStart struct {
	ph       *PlaceholderGen
	negated  bool
	column   []ColumnRef
	operator []string
	values   []any
	ended    *WhereEnd
}

func (w *WhereStart) start(field string, from string) *WhereStart {
	w.column = append(w.column, ColumnRef{table: from, field: field})
	w.negated = false
	return w
}

func (w *WhereStart) end() *WhereEnd {
	if w.ended == nil {
		w.ended = &WhereEnd{arg: w}
	}
	return w.ended
}

func (w *WhereStart) Not() *WhereStart {
	w.negated = !w.negated
	return w
}

func (w *WhereStart) Equal(value any) *WhereEnd {
	w.operator = append(w.operator, opEqual(w.ph, w.negated))
	w.values = append(w.values, value)
	return w.end()
}

func (w *WhereStart) Like(value any) *WhereEnd {
	w.operator = append(w.operator, opLike(w.ph, w.negated))
	w.values = append(w.values, value)
	return w.end()
}

func (w *WhereStart) In(values ...any) *WhereEnd {
	w.operator = append(w.operator, opIn(w.ph, w.negated, len(values)))
	w.values = append(w.values, values...)
	return w.end()
}

func (w *WhereStart) Between(minV, maxV any) *WhereEnd {
	w.operator = append(w.operator, opBetween(w.ph, w.negated))
	w.values = append(w.values, minV, maxV)
	return w.end()
}

func (w *WhereStart) Greater(value any) *WhereEnd {
	w.operator = append(w.operator, opGreater(w.ph, w.negated))
	w.values = append(w.values, value)
	return w.end()
}

func (w *WhereStart) Less(value any) *WhereEnd {
	w.operator = append(w.operator, opLess(w.ph, w.negated))
	w.values = append(w.values, value)
	return w.end()
}

func (w *WhereStart) IsNull() *WhereEnd {
	w.operator = append(w.operator, opIsNull(w.negated))
	return w.end()
}

// -- END EXPRESSION
type WhereEnd struct {
	arg        *WhereStart
	connectors []string
}

// Or starts a new condition with logical OR.
// - If no arguments are provided, it reuses the previous field from the clause.
// - If two arguments are provided, the first is treated as the field name and the second as its alias.
func (w *WhereEnd) Or(fieldFrom ...string) *WhereStart {
	return w.nextCondition("OR", fieldFrom)
}

// And starts a new condition with logical AND.
// - If no arguments are provided, it reuses the previous field from the clause.
// - If two arguments are provided, the first is treated as the field name and the second as its alias.
func (w *WhereEnd) And(fieldFrom ...string) *WhereStart {
	return w.nextCondition("AND", fieldFrom)
}

func (w *WhereEnd) nextCondition(logical string, fieldFrom []string) *WhereStart {
	w.connectors = append(w.connectors, logical)
	index := w.arg.column[len(w.arg.column)-1]

	var field, table string
	if ln := len(fieldFrom); ln > 1 {
		field = fieldFrom[0]
		table = fieldFrom[1]
	} else if ln > 0 {
		field = fieldFrom[0]
		table = index.table
	} else {
		field = index.field
		table = index.table
	}

	return w.arg.start(field, table)
}

// ----- Helpers -----
type ColumnRef struct {
	field string
	table string
}

func opEqual(ph *PlaceholderGen, negated bool) string {
	op := "= "
	if negated {
		op = "!= "
	}

	return op + ph.Next()
}

func opLike(ph *PlaceholderGen, negated bool) string {
	op := "LIKE "
	if negated {
		op = "NOT LIKE "
	}

	return op + ph.Next()
}

func opIn(ph *PlaceholderGen, negated bool, count int) string {
	placeholders := make([]string, count)
	for i := range placeholders {
		placeholders[i] = ph.Next()
	}

	inClause := "(" + strings.Join(placeholders, ", ") + ")"
	if negated {
		return "NOT IN " + inClause
	}
	return "IN " + inClause
}

func opBetween(ph *PlaceholderGen, negated bool) string {
	op := "BETWEEN "
	if negated {
		op = "NOT BETWEEN "
	}

	return op + ph.Next() + " AND " + ph.Next()
}

func opGreater(ph *PlaceholderGen, negated bool) string {
	op := "> "
	if negated {
		op = "<= "
	}

	return op + ph.Next()
}

func opLess(ph *PlaceholderGen, negated bool) string {
	op := "< "
	if negated {
		op = ">= "
	}
	return op + ph.Next()
}

func opIsNull(negated bool) string {
	if negated {
		return "IS NOT NULL"
	}
	return "IS NULL"
}
