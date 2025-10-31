package expr

type Expr struct {
	segments []*exprSegment
	current  *exprSegment
	end      *ExprEnd
}

// Not() toggles the negation of the current segment.
func (e *Expr) Not() *Expr {
	e.current.negated = !e.current.negated
	return e
}

// Equal() adds an equality condition.
//
// @SQL: identifier = value
//
// @SQL (negated): identifier <> value
func (e *Expr) Equal(value any) *ExprEnd {
	e.current.addArg(equal, value)
	return e.end
}

// Greater() adds a greater-than condition.
//
// @SQL: identifier > value
//
// @SQL (negated): identifier <= value
func (e *Expr) Greater(value any) *ExprEnd {
	e.current.addArg(greater, value)
	return e.end
}

// Less() adds a less-than condition.
//
// @SQL: identifier < value
//
// @SQL (negated): identifier >= value
func (e *Expr) Less(value any) *ExprEnd {
	e.current.addArg(less, value)
	return e.end
}

// In() adds an IN condition.
//
// @SQL: identifier IN (value1, value2, ...)
//
// @SQL (negated): identifier NOT IN (value1, value2, ...)
func (e *Expr) In(values ...any) *ExprEnd {
	e.current.op = in
	e.current.Argument = values
	return e.end
}

// Between() adds a BETWEEN condition.
//
// @SQL: identifier BETWEEN minValue AND maxValue
//
// @SQL (negated): identifier NOT BETWEEN minValue AND maxValue
func (e *Expr) Between(minV, maxV any) *ExprEnd {
	e.current.op = between
	e.current.Argument = []any{minV, maxV}
	return e.end
}

// IsNull() adds an IS NULL condition.
//
// @SQL: identifier IS NULL
//
// @SQL (negated): identifier IS NOT NULL
func (e *Expr) IsNull() *ExprEnd {
	e.current.op = isNull
	return e.end
}

// Like() adds a LIKE condition.
//
// @SQL: identifier LIKE value
//
// @SQL (negated): identifier NOT LIKE value
func (e *Expr) Like(value any) *ExprEnd {
	e.current.addArg(like, value)
	return e.end
}

type ExprEnd struct {
	ctx        *Expr
	connectors []string
}

// Or starts a new segment connected with the OR operator.
func (e *ExprEnd) Or(identifier ...any) *Expr {
	return e.nextCondition(or, identifier)
}

// And starts a new segment connected with the AND operator.
func (e *ExprEnd) And(identifier ...any) *Expr {
	return e.nextCondition(and, identifier)
}
