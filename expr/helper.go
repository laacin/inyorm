package expr

type exprSegment struct {
	Identifier any
	op         sqlOp
	negated    bool
	Argument   []any
}

func NewExpr(identifier any) *Expr {
	expr := &Expr{}
	return expr.start(identifier)
}

func (e *exprSegment) addArg(op sqlOp, arg any) {
	e.op = op
	e.Argument = []any{arg}
}

func (e *Expr) start(identifier any) *Expr {
	if e.end == nil {
		e.end = &ExprEnd{ctx: e}
	}

	segment := &exprSegment{Identifier: identifier}
	e.current = segment
	e.segments = append(e.segments, segment)
	return e
}

func (e *ExprEnd) nextCondition(logical sqlOp, identifier []any) *Expr {
	e.connectors = append(e.connectors, string(logical))

	var ident any
	if ln := len(identifier); ln > 0 {
		ident = identifier[0]
	} else {
		ident = e.ctx.current.Identifier
	}

	return e.ctx.start(ident)
}
