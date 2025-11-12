package inyorm

import (
	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/core"
)

type (
	Condition struct {
		wrap core.Condition
		next *CondNext
	}
	CondNext struct {
		wrap core.Condition
		ctx  *Condition
	}

	Case[T any] struct {
		wrap core.Case
		next *CaseNext[T]
	}
	CaseNext[T any] struct {
		wrap core.Case
		ctx  *Case[T]
	}
	CaseSwitch = *Case[Value]
	CaseSearch = *Case[*CondNext]
)

// ----- Constructors

func wrapCondition(wrap core.Condition) *Condition {
	start := &Condition{wrap: wrap}
	next := &CondNext{ctx: start, wrap: wrap}
	start.next = next
	return start
}

func newCase[T any]() *Case[T] {
	cs := &column.Case{}
	start := &Case[T]{wrap: cs}
	next := &CaseNext[T]{ctx: start, wrap: cs}
	start.next = next
	return start
}

// ----- Condition -----

func (c *Condition) Not() *Condition               { c.wrap.Not(); return c }
func (c *Condition) Equal(value Value) *CondNext   { c.wrap.Equal(vOne(value)); return c.next }
func (c *Condition) Greater(value Value) *CondNext { c.wrap.Greater(vOne(value)); return c.next }
func (c *Condition) Less(value Value) *CondNext    { c.wrap.Less(vOne(value)); return c.next }
func (c *Condition) In(values ...Value) *CondNext  { c.wrap.In(vMany(values)); return c.next }
func (c *Condition) Between(minV, maxV Value) *CondNext {
	c.wrap.Between(vOne(minV), vOne(maxV))
	return c.next
}
func (c *Condition) IsNull() *CondNext { c.wrap.IsNull(); return c.next }

func (c *CondNext) And(identifier Value) *Condition { c.wrap.And(vOne(identifier)); return c.ctx }
func (c *CondNext) Or(identifier Value) *Condition  { c.wrap.Or(vOne(identifier)); return c.ctx }

// ----- Case -----

func (c *Case[T]) When(value T) *CaseNext[T] { c.wrap.When(vOne(value)); return c.next }
func (c *Case[T]) Else(value Value)          { c.wrap.Else(vOne(value)) }

func (c *CaseNext[T]) Then(value Value) *Case[T] { c.wrap.Then(vOne(value)); return c.ctx }
