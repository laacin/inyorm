package inyorm

import (
	"strings"

	"github.com/laacin/inyorm/internal/stmt"
)

// Expression represents the start of a conditional expression.
// It allows building chained conditions on columns or fields.
type Expression struct {
	ph       *stmt.PlaceholderGen
	negated  bool
	segments []*ExprSegment
	current  *ExprSegment
	end      *ExpressionEnd
	sb       strings.Builder
}

func (w *Expression) start(identifier string) *Expression {
	if w.end == nil {
		w.end = &ExpressionEnd{ctx: w}
	}

	segment := &ExprSegment{Identifier: identifier}
	w.current = segment
	w.segments = append(w.segments, segment)
	w.negated = false
	return w
}

// Not negates the current condition.
// If called before an operator, the operation will be negated.
func (w *Expression) Not() *Expression {
	w.negated = !w.negated
	return w
}

// Equal adds an equality condition: field = value.
// Returns ExpressionEnd to chain logical AND/OR conditions.
func (w *Expression) Equal(value any) *ExpressionEnd {
	w.current.Argument = w.writeOp(stmt.Equal, value)
	return w.end
}

// Like adds a pattern-matching condition: field LIKE value.
// Returns ExpressionEnd to chain logical conditions.
func (w *Expression) Like(value any) *ExpressionEnd {
	w.current.Argument = w.writeOp(stmt.Like, value)
	return w.end
}

// In adds a membership condition: field IN (values...).
// Returns ExpressionEnd to chain logical AND/OR conditions.
func (w *Expression) In(values ...any) *ExpressionEnd {
	w.current.Argument = w.writeOp(stmt.In, values...)
	return w.end
}

// Between adds a range condition: field BETWEEN minV AND maxV.
// Returns ExpressionEnd to chain logical conditions.
func (w *Expression) Between(minV, maxV any) *ExpressionEnd {
	w.current.Argument = w.writeOp(stmt.Between, minV, maxV)
	return w.end
}

// Greater adds a greater-than condition: field > value.
// Returns ExpressionEnd to chain logical AND/OR conditions.
func (w *Expression) Greater(value any) *ExpressionEnd {
	w.current.Argument = w.writeOp(stmt.Greater, value)
	return w.end
}

// Less adds a less-than condition: field < value.
// Returns ExpressionEnd to chain logical conditions.
func (w *Expression) Less(value any) *ExpressionEnd {
	w.current.Argument = w.writeOp(stmt.Less, value)
	return w.end
}

// IsNull adds a null-check condition: field IS NULL.
// Returns ExpressionEnd to chain logical AND/OR conditions.
func (w *Expression) IsNull() *ExpressionEnd {
	w.current.Argument = w.writeOp(stmt.IsNull)
	return w.end
}

// ExpressionEnd represents the end of a conditional expression,
// allowing logical connectors to continue building the query.
type ExpressionEnd struct {
	ctx        *Expression
	connectors []string
}

// Or starts a new condition with a logical OR.
// If identifiers are provided, the first is treated as the field to evaluate,
// and the second (optional) as its alias or table.
func (w *ExpressionEnd) Or(identifier ...string) *Expression {
	return w.nextCondition(stmt.Or, identifier)
}

// And starts a new condition with a logical AND.
// If identifiers are provided, the first is treated as the field to evaluate,
// and the second (optional) as its alias or table.
func (w *ExpressionEnd) And(identifier ...string) *Expression {
	return w.nextCondition(stmt.And, identifier)
}

func (w *ExpressionEnd) nextCondition(logical stmt.Operator, identifier []string) *Expression {
	w.connectors = append(w.connectors, string(logical))

	var ident string
	if ln := len(identifier); ln > 0 {
		ident = identifier[0]
	} else {
		ident = w.ctx.current.Identifier
	}

	return w.ctx.start(ident)
}

// ----- Helpers -----
type ExprSegment struct {
	Identifier string
	Argument   string
}

func (w *Expression) writeOp(kind stmt.Operator, values ...any) string {
	w.sb.Reset()
	w.sb.WriteString(stmt.GetOp(kind, w.negated))
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
