package clause

import (
	"strings"

	"github.com/laacin/inyorm/internal/stmt"
)

type Expr struct {
	segments []*exprSegment
	current  *exprSegment
	end      *ExprEnd
}

type exprSegment struct {
	Identifier any
	op         sqlOp
	negated    bool
	Argument   []any
}

func (e *exprSegment) addArg(op sqlOp, arg any) {
	e.op = op
	e.Argument = []any{arg}
}

func (e *Expr) Start(identifier any) *Expr {
	if e.end == nil {
		e.end = &ExprEnd{ctx: e}
	}

	segment := &exprSegment{Identifier: identifier}
	e.current = segment
	e.segments = append(e.segments, segment)
	return e
}

func (e *Expr) Not() *Expr {
	e.current.negated = !e.current.negated
	return e
}

func (e *Expr) Equal(value any) *ExprEnd {
	e.current.addArg(equal, value)
	return e.end
}

func (e *Expr) Greater(value any) *ExprEnd {
	e.current.addArg(greater, value)
	return e.end
}

func (e *Expr) Less(value any) *ExprEnd {
	e.current.addArg(less, value)
	return e.end
}

func (e *Expr) In(values ...any) *ExprEnd {
	e.current.op = in
	e.current.Argument = values
	return e.end
}

func (e *Expr) Between(minV, maxV any) *ExprEnd {
	e.current.op = between
	e.current.Argument = []any{minV, maxV}
	return e.end
}

func (e *Expr) IsNull() *ExprEnd {
	e.current.op = isNull
	return e.end
}

func (e *Expr) Like(value any) *ExprEnd {
	e.current.addArg(like, value)
	return e.end
}

type ExprEnd struct {
	ctx        *Expr
	connectors []string
}

func (e *ExprEnd) Or(identifier ...any) *Expr {
	return e.nextCondition(or, identifier)
}

func (e *ExprEnd) And(identifier ...any) *Expr {
	return e.nextCondition(and, identifier)
}

func (e *ExprEnd) nextCondition(logical sqlOp, identifier []any) *Expr {
	e.connectors = append(e.connectors, string(logical))

	var ident any
	if ln := len(identifier); ln > 0 {
		ident = identifier[0]
	} else {
		ident = e.ctx.current.Identifier
	}

	return e.ctx.Start(ident)
}

// ----- Build
func (e *Expr) build(sb *strings.Builder, ph *Placeholder) {
	sb.WriteByte('(')
	for i, seg := range e.segments {
		if i > 0 {
			sb.WriteByte(' ')
			sb.WriteString(e.end.connectors[i-1])
			sb.WriteByte(' ')
		}

		sb.WriteString(stmt.Stringify(seg.Identifier))
		sb.WriteByte(' ')
		writeArg(sb, ph, seg.op, seg.negated, seg.Argument)
	}
	sb.WriteByte(')')
}

type sqlOp string

const (
	equal    sqlOp = "="
	notEqual sqlOp = "<>"

	greater    sqlOp = ">"
	notGreater sqlOp = "<="

	less    sqlOp = "<"
	notLess sqlOp = ">="

	and sqlOp = "AND"
	or  sqlOp = "OR"

	in    sqlOp = "IN"
	notIn sqlOp = "NOT IN"

	between    sqlOp = "BETWEEN"
	notBetween sqlOp = "NOT BETWEEN"

	isNull    sqlOp = "IS NULL"
	isNotNull sqlOp = "IS NOT NULL"

	like    sqlOp = "LIKE"
	notLike sqlOp = "NOT LIKE"
)

var negations = map[sqlOp]sqlOp{
	equal:    notEqual,
	notEqual: equal,

	greater:    notGreater,
	notGreater: greater,

	less:    notLess,
	notLess: less,

	in:    notIn,
	notIn: in,

	between:    notBetween,
	notBetween: between,

	isNull:    isNotNull,
	isNotNull: isNull,

	like:    notLike,
	notLike: like,
}

func writeArg(sb *strings.Builder, ph *Placeholder, kind sqlOp, negated bool, values []any) {
	var op string
	if negated {
		if val, ok := negations[kind]; ok {
			op = string(val)
		}
	} else {
		op = string(kind)
	}

	if kind == isNull {
		sb.WriteString(op)
		return
	}

	switch kind {
	case isNull:
		sb.WriteString(op)
	case between:
		sb.WriteString(op)
		sb.WriteByte(' ')
		ph.Write(sb, values[0])
		sb.WriteByte(' ')
		sb.WriteString(string(and))
		sb.WriteByte(' ')
		ph.Write(sb, values[1])
	default:
		sb.WriteString(op)
		sb.WriteByte(' ')
		ph.Write(sb, values...)
	}
}
