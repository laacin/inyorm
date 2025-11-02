package inyorm

import (
	"github.com/laacin/inyorm/clause"
	"github.com/laacin/inyorm/column"
	"github.com/laacin/inyorm/expression"
)

// ----- Clause -----

type ManyToMany struct {
	mainTbl string
	inter   *clause.InterJoin
}

func (m *ManyToMany) Join(typ string, table Table) {
	if fk, exists := table.ForeignKeys[m.mainTbl]; exists {
		m.inter.Join(clause.JoinTyp(typ), table.Name, "alias", fk)
	}
}

// ----- Expression -----

func wrapExpr(value Value) *Expr {
	e := &expression.Expr{}
	e.Start(value)
	start := &Expr{wrap: e}
	end := &ExprEnd{wrap: e.End}
	start.end = end
	end.ctx = start
	return start
}

type Expr struct {
	wrap *expression.Expr
	end  *ExprEnd
}

// Not toggles the negation of the current segment.
func (e *Expr) Not() *Expr {
	e.wrap.Not()
	return e
}

// Equal adds an equality condition.
//
// @SQL: identifier = value
//
// @SQL (negated): identifier <> value
func (e *Expr) Equal(value Value) *ExprEnd {
	e.wrap.Equal(value)
	return e.end
}

// Greater adds a greater-than condition.
//
// @SQL: identifier > value
//
// @SQL (negated): identifier <= value
func (e *Expr) Greater(value Value) *ExprEnd {
	e.wrap.Greater(value)
	return e.end
}

// Less adds a less-than condition.
//
// @SQL: identifier < value
//
// @SQL (negated): identifier >= value
func (e *Expr) Less(value Value) *ExprEnd {
	e.wrap.Less(value)
	return e.end
}

// In adds an IN condition.
//
// @SQL: identifier IN (value1, value2, ...)
//
// @SQL (negated): identifier NOT IN (value1, value2, ...)
func (e *Expr) In(values ...Value) *ExprEnd {
	e.wrap.In(values)
	return e.end
}

// Between adds a BETWEEN condition.
//
// @SQL: identifier BETWEEN minValue AND maxValue
//
// @SQL (negated): identifier NOT BETWEEN minValue AND maxValue
func (e *Expr) Between(minV, maxV Value) *ExprEnd {
	e.wrap.Between(minV, maxV)
	return e.end
}

// IsNull adds an IS NULL condition.
//
// @SQL: identifier IS NULL
//
// @SQL (negated): identifier IS NOT NULL
func (e *Expr) IsNull() *ExprEnd {
	e.wrap.IsNull()
	return e.end
}

// Like adds a LIKE condition.
//
// @SQL: identifier LIKE value
//
// @SQL (negated): identifier NOT LIKE value
func (e *Expr) Like(value Value) *ExprEnd {
	e.wrap.Like(value)
	return e.end
}

type ExprEnd struct {
	ctx  *Expr
	wrap *expression.ExprEnd
}

// Or starts a new segment connected with the OR operator.
func (e *ExprEnd) Or(value ...Value) *Expr {
	e.wrap.Next(expression.Or, value)
	return e.ctx
}

// And starts a new segment connected with the AND operator.
func (e *ExprEnd) And(value ...Value) *Expr {
	e.wrap.Next(expression.And, value)
	return e.ctx
}

// ----- Case -----

type (
	CaseValue  interface{ *ExprEnd | Value }
	CaseSwitch = *Case[Value]
	CaseSearch = *Case[*ExprEnd]
)

func wrapCase[T CaseValue]() *Case[T] {
	cs := &column.Case{}
	return &Case[T]{wrap: cs}
}

// Represents a CASE field context used to define WHEN/THEN and ELSE clauses.
type Case[T CaseValue] struct {
	wrap *column.Case
}

// Defines a WHEN condition within a CASE clause.
// Returns a *Do[T] instance used to define the THEN action.
//
// SQL: WHEN v THEN ...
func (cs *Case[T]) When(value T) *Do[T] {
	switch val := any(value).(type) {
	case *ExprEnd:
		do := cs.wrap.When(val.wrap)
		return &Do[T]{ctx: cs, wrap: do}
	default:
		do := cs.wrap.When(val)
		return &Do[T]{ctx: cs, wrap: do}
	}
}

// Defines the ELSE branch of a CASE clause.
//
// SQL: ELSE v
func (cs *Case[T]) Else(value Value) {
	cs.wrap.Else(value)
}

// Represents a single WHEN/THEN pair within a CASE clause.
type Do[T CaseValue] struct {
	ctx  *Case[T]
	wrap *column.Do
}

// Then sets the THEN value for a WHEN condition in a CASE clause.
func (d *Do[T]) Then(value Value) *Case[T] {
	d.wrap.Then(value)
	return d.ctx
}
