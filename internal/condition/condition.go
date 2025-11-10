package condition

import "github.com/laacin/inyorm/internal/core"

type Condition struct {
	segments []*exprSegment
	current  *exprSegment
	Next     *ConditionNext
}

func (e *Condition) Not() core.Cond {
	e.current.Negated = !e.current.Negated
	return e
}

func (e *Condition) Equal(value any) core.CondNext {
	e.current.addArg(equal, value)
	return e.Next
}

func (e *Condition) Greater(value any) core.CondNext {
	e.current.addArg(greater, value)
	return e.Next
}

func (e *Condition) Less(value any) core.CondNext {
	e.current.addArg(less, value)
	return e.Next
}

func (e *Condition) In(values ...any) core.CondNext {
	e.current.Operator = in
	e.current.Argument = values
	return e.Next
}

func (e *Condition) Between(minV, maxV any) core.CondNext {
	e.current.Operator = between
	e.current.Argument = []any{minV, maxV}
	return e.Next
}

func (e *Condition) IsNull() core.CondNext {
	e.current.Operator = isNull
	return e.Next
}

func (e *Condition) Like(value any) core.CondNext {
	e.current.addArg(like, value)
	return e.Next
}

type ConditionNext struct {
	ctx        *Condition
	connectors []string
}

func (e *ConditionNext) And(identifier ...any) core.Cond {
	e.Next(and, identifier)
	return e.ctx
}

func (e *ConditionNext) Or(identifier ...any) core.Cond {
	e.Next(or, identifier)
	return e.ctx
}

func (e *ConditionNext) Next(logical string, identifier []any) *Condition {
	e.connectors = append(e.connectors, logical)

	var ident any
	if ln := len(identifier); ln > 0 {
		ident = identifier[0]
	} else {
		ident = e.ctx.current.Identifier
	}

	return e.ctx.Start(ident)
}
