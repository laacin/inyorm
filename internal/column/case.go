package column

import "github.com/laacin/inyorm/internal/core"

type Case[T any] struct {
	args []*CaseNext[T]
	els  any
}

func (c *Case[T]) When(v T) core.CaseNext[T] {
	arg := &CaseNext[T]{ctx: c, when: v}
	c.args = append(c.args, arg)
	return arg
}

func (c *Case[T]) Else(v any) {
	c.els = v
}

type CaseNext[T any] struct {
	ctx  *Case[T]
	when T
	do   any
}

func (c *CaseNext[T]) Then(v any) core.Case[T] {
	c.do = v
	return c.ctx
}
