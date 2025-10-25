package inyorm

import (
	"strings"

	"github.com/laacin/inyorm/internal/stmt"
)

// Expr represents the start of a conditional expression.
// It allows building chained conditions on columns or fields.
type Expr struct {
	ph       *stmt.PlaceholderGen
	negated  bool
	segments []*exprSegment
	current  *exprSegment
	end      *ExprEnd
	sb       strings.Builder
}

func (w *Expr) start(identifier any) *Expr {
	if w.end == nil {
		w.end = &ExprEnd{ctx: w}
	}

	segment := &exprSegment{Identifier: identifier}
	w.current = segment
	w.segments = append(w.segments, segment)
	w.negated = false
	return w
}

// Not negates the current condition.
// If called before an operator, the operation will be negated.
func (w *Expr) Not() *Expr {
	w.negated = !w.negated
	return w
}

// Equal adds an equality condition: field = value.
// Returns ExprEnd to chain logical AND/OR conditions.
func (w *Expr) Equal(value Value) *ExprEnd {
	w.current.Argument = w.writeOp(stmt.Equal, value)
	return w.end
}

// Like adds a pattern-matching condition: field LIKE value.
// Returns ExprEnd to chain logical conditions.
func (w *Expr) Like(value Value) *ExprEnd {
	w.current.Argument = w.writeOp(stmt.Like, value)
	return w.end
}

// In adds a membership condition: field IN (values...).
// Returns ExprEnd to chain logical AND/OR conditions.
func (w *Expr) In(values ...any) *ExprEnd {
	w.current.Argument = w.writeOp(stmt.In, values...)
	return w.end
}

// Between adds a range condition: field BETWEEN minV AND maxV.
// Returns ExprEnd to chain logical conditions.
func (w *Expr) Between(minV, maxV any) *ExprEnd {
	w.current.Argument = w.writeOp(stmt.Between, minV, maxV)
	return w.end
}

// Greater adds a greater-than condition: field > value.
// Returns ExprEnd to chain logical AND/OR conditions.
func (w *Expr) Greater(value any) *ExprEnd {
	w.current.Argument = w.writeOp(stmt.Greater, value)
	return w.end
}

// Less adds a less-than condition: field < value.
// Returns ExprEnd to chain logical conditions.
func (w *Expr) Less(value any) *ExprEnd {
	w.current.Argument = w.writeOp(stmt.Less, value)
	return w.end
}

// IsNull adds a null-check condition: field IS NULL.
// Returns ExprEnd to chain logical AND/OR conditions.
func (w *Expr) IsNull() *ExprEnd {
	w.current.Argument = w.writeOp(stmt.IsNull)
	return w.end
}

// ExprEnd represents the end of a conditional expression,
// allowing logical connectors to continue building the query.
type ExprEnd struct {
	ctx        *Expr
	connectors []string
}

// Or starts a new condition with a logical OR.
// If identifiers are provided, the first is treated as the field to evaluate,
// and the second (optional) as its alias or table.
func (w *ExprEnd) Or(identifier ...any) *Expr {
	return w.nextCondition(stmt.Or, identifier)
}

// And starts a new condition with a logical AND.
// If identifiers are provided, the first is treated as the field to evaluate,
// and the second (optional) as its alias or table.
func (w *ExprEnd) And(identifier ...any) *Expr {
	return w.nextCondition(stmt.And, identifier)
}

func (w *ExprEnd) nextCondition(logical stmt.SqlOp, identifier []any) *Expr {
	w.connectors = append(w.connectors, string(logical))

	var ident any
	if ln := len(identifier); ln > 0 {
		ident = identifier[0]
	} else {
		ident = w.ctx.current.Identifier
	}

	return w.ctx.start(ident)
}

// ----- Helpers -----
type exprSegment struct {
	Identifier any
	Argument   string
}

func (w *Expr) writeOp(kind stmt.SqlOp, values ...any) string {
	w.sb.Reset()
	w.sb.WriteString(stmt.GetSqlOp(kind, w.negated))
	if len(values) == 0 {
		return w.sb.String()
	}

	w.sb.WriteByte(' ')
	if kind == stmt.Between {
		w.ph.Write(&w.sb, values[0])
		w.sb.WriteString(" AND ")
		w.ph.Write(&w.sb, values[1])
	} else {
		w.ph.Write(&w.sb, values...)
	}
	return w.sb.String()
}

func (w *Expr) build(sb *strings.Builder) {
	for i, seg := range w.segments {
		if i > 0 {
			con := w.end.connectors[i-1]
			sb.WriteByte(' ')
			sb.WriteString(con)
			sb.WriteByte(' ')
		}

		sb.WriteByte('(')
		sb.WriteString(stmt.Stringify(seg.Identifier))
		sb.WriteByte(' ')
		sb.WriteString(seg.Argument)
		sb.WriteByte(')')
	}
}
