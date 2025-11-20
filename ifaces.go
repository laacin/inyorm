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
		CaseSwitch, CaseSearch, CaseNext,
		Identifier, Value,
	]

	clsSelect  = clause.Select[SelectNext, Identifier]
	clsFrom    = clause.From
	clsJoin    = clause.Join[JoinNext, Condition, ConditionNext, Identifier, Value]
	clsWhere   = clause.Where[Condition, ConditionNext, Identifier, Value]
	clsGroupBy = clause.GroupBy[Identifier]
	clsHaving  = clause.Having[Condition, ConditionNext, Identifier, Value]
	clsOrderBy = clause.OrderBy[OrderByNext, Identifier]
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

type CaseSwitch interface {
	When(v Value) CaseNext
	Else(v Value)
}

type CaseSearch interface {
	When(cond ConditionNext) CaseNext
	Else(v Value)
}

type CaseNext interface {
	Then(v Value) CaseSearch
}

// ----- Column Builder -----

type ColumnBuilder interface {
	Col(name string, table ...string) Column
	All() Column

	Ph() Builder
	Cond(ident Identifier) Condition

	Concat(v ...Identifier) Column
	Switch(cond Identifier, fn func(cs CaseSwitch)) Column
	Search(fn func(cs CaseSearch)) Column
}

// ----- Clauses -----

// -- SELECT

type Select interface {
	Distinct() SelectNext
	Select(sel ...Identifier)
}

type SelectNext interface {
	Select(sel ...Identifier)
}

// -- FROM

type From interface {
	From(table string)
}

// -- JOIN

type Join interface {
	Join(table string) JoinNext
}

type JoinNext interface {
	On(on Identifier) Condition
}

// -- WHERE

type Where interface {
	Where(ident Identifier) Condition
}

// -- GROUP BY

type GroupBy interface {
	GroupBy(group ...Identifier)
}

// -- HAVING

type Having interface {
	Having(having Value) Condition
}

// -- ORDER BY

type OrderBy interface {
	OrderBy(order Identifier) OrderByNext
}

type OrderByNext interface {
	Desc()
}

// -- LIMIT

type Limit interface {
	Limit(limit int)
}

type Offset interface {
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
