package inyorm

import (
	"github.com/laacin/inyorm/clause"
	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/core"
)

type (
	Builder    = core.Builder
	Value      = any
	Identifier = any

	colBuilder = column.ColBuilder[
		Column, Condition, ConditionNext,
		Case, CaseNext,
	]

	clsSelect  = clause.Select[SelectNext]
	clsFrom    = clause.From
	clsJoin    = clause.Join[JoinNext, Condition, ConditionNext]
	clsWhere   = clause.Where[Condition, ConditionNext]
	clsGroupBy = clause.GroupBy
	clsHaving  = clause.Having[Condition, ConditionNext]
	clsOrderBy = clause.OrderBy[OrderByNext]
	clsLimit   = clause.Limit
	clsOffset  = clause.Offset
)

// ---- Column

type Column interface {
	Def() Builder
	Expr() Builder
	Alias() Builder
	Base() Builder

	Count(distinct ...bool) Column
	Sum(distinct ...bool) Column
	Min(distinct ...bool) Column
	Max(distinct ...bool) Column
	Avg(distinct ...bool) Column

	Add(v Value) Column
	Sub(v Value) Column
	Mul(v Value) Column
	Div(v Value) Column
	Mod(v Value) Column
	Wrap() Column

	Lower() Column
	Upper() Column
	Trim() Column
	Round() Column
	Abs() Column

	As(name string) Column
}

// ----- Condition -----

type Condition interface {
	Not() Condition

	Equal(v Value) ConditionNext
	Like(v Value) ConditionNext
	Greater(v Value) ConditionNext
	Less(v Value) ConditionNext
	In(v ...Value) ConditionNext
	Between(minV, maxV Value) ConditionNext
	IsNull() ConditionNext
}

type ConditionNext interface {
	And(ident Identifier) Condition
	Or(ident Identifier) Condition
}

// ----- Case -----

type Case interface {
	When(cond Identifier) CaseNext
	Else(els Value)
}

type CaseNext interface {
	Then(then Value) Case
}

// ----- Column Builder -----

type ColumnBuilder interface {
	Col(name string, table ...string) Column
	All() Column

	Ph() Builder
	Cond(ident Identifier) Condition

	Concat(v ...Identifier) Column
	Switch(cond Identifier, fn func(cs Case)) Column
	Search(fn func(cs Case)) Column
}

// ----- Clauses -----

// ----- SELECT

type Select interface {
	// Distinct writes DISTINCT in the SELECT clause
	//
	// @SQL: SELECT `DISTINCT` ...
	Distinct() SelectNext

	// Select writes the SELECT clause values
	//
	// @SQL: SELECT `DISTINCT?` `sel1`, `sel2`, `sel3` ... [SelectNext]
	Select(sel ...Identifier)
}

type SelectNext interface {
	// Select writes the SELECT clause values
	//
	// @SQL: SELECT `DISTINCT?` `sel1`, `sel2`, `sel3` ... [SelectNext]
	Select(sel ...Identifier)
}

// ----- FROM

type From interface {
	// From writes the FROM clause
	//
	// # This method is auto-generated for the statement’s default table.
	// Only use it for complex FROM clauses (such as subqueries)
	//
	// @SQL: FROM `table`
	From(table string)
}

// ----- JOIN

type Join interface {
	// Join writes the JOIN clause
	//
	// @SQL: INNER JOIN `table 'alias'` ... [JoinNext]
	Join(table string) JoinNext
}

type JoinNext interface {
	// On writes the join condition
	//
	// @SQL: [join] ... ON `on` ... [Condition]
	On(on Identifier) Condition
}

// ----- WHERE

type Where interface {
	// Where writes the WHERE clause
	//
	// # Can be called multiple times,
	// Conditions are combined using the logical "AND".
	// e.g: (cond1) AND (cond2) AND (cond3) ...
	//
	// @SQL: WHERE `ident` ... [Condition]
	Where(ident Identifier) Condition
}

// ----- GROUP BY

type GroupBy interface {
	// GroupBy writes the GROUP BY clause
	//
	// @SQL: GROUP BY `group1`, `group2`, `group3` ...
	GroupBy(group ...Identifier)
}

// ----- HAVING

type Having interface {
	// Having writes the HAVING clause
	//
	// @SQL: HAVING `having` ... [Condition]
	Having(having Value) Condition
}

// ----- ORDER BY

type OrderBy interface {
	// OrderBy writes the ORDER BY clause
	//
	// # Can be called multiple times for multiple orderings
	//
	// @SQL: ORDER BY `order` ... [OrderByNext]
	OrderBy(order Identifier) OrderByNext
}

type OrderByNext interface {
	// Desc sets the descending direction for the current order
	//
	// @SQL: [OrderBy] ... DESC
	Desc()
}

// ----- LIMIT

type Limit interface {
	// Limit writes the LIMIT clause value
	//
	// # Values less than 1 will be ignored
	//
	// @SQL: LIMIT `limit`
	Limit(limit int)
}

type Offset interface {
	// Offset writes the OFFSET clause value
	//
	// # Values less than 1 will be ignored
	//
	// @SQL: OFFSET `offset`
	Offset(offset int)
}

// ----- Statements -----

type SelectStmt interface {
	Select
	From
	Join
	Where
	GroupBy
	Having
	OrderBy
	Limit
	Offset
}
