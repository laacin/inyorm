package condition

type exprSegment struct {
	Identifier any
	Operator   string
	Negated    bool
	Argument   []any
}

func (e *exprSegment) addArg(op string, arg any) {
	e.Operator = op
	e.Argument = []any{arg}
}

func (e *Condition) Start(identifier any) *Condition {
	if e.Next == nil {
		e.Next = &ConditionNext{ctx: e}
	}

	segment := &exprSegment{Identifier: identifier}
	e.current = segment
	e.segments = append(e.segments, segment)
	return e
}
