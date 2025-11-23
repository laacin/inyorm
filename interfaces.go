package inyorm

import (
	"github.com/laacin/inyorm/clause"
	"github.com/laacin/inyorm/exec"
	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/core"
)

// ---- internal
type (
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

	executor = exec.Executor[Prepare]
)

type (
	Builder    = core.Builder
	Value      = any
	Identifier = any
	Binder     = any
)

// ---- Column

// Column represents a selectable or computable SQL column.
// Provides builders for expressions, aliases, aggregates,
// arithmetic operations, and transformations.
type Column interface {
	// Def writes the definition of the column
	//
	// e.g: if we have this column being built: COUNT(id) AS count:
	//  - Def()   would write: COUNT(id) AS count; Fallbacks: Expr() -> Base()
	//  - Expr()  would write: COUNT(id);          Fallbacks: Base()
	//  - Alias() would write: count;              Fallbacks: Expr() -> Base()
	//  - Base()  would write: id                  Fallbacks: none
	//
	// @SQL: used to build column definitions
	Def() Builder

	// Expr writes the expression of the column
	//
	// e.g: if we have this column being built: COUNT(id) AS count:
	//  - Def()   would write: COUNT(id) AS count; Fallbacks: Expr() -> Base()
	//  - Expr()  would write: COUNT(id);          Fallbacks: Base()
	//  - Alias() would write: count;              Fallbacks: Expr() -> Base()
	//  - Base()  would write: id                  Fallbacks: none
	Expr() Builder

	// Alias writes the column alias
	//
	// e.g: if we have this column being built: COUNT(id) AS count:
	//  - Def()   would write: COUNT(id) AS count; Fallbacks: Expr() -> Base()
	//  - Expr()  would write: COUNT(id);          Fallbacks: Base()
	//  - Alias() would write: count;              Fallbacks: Expr() -> Base()
	//  - Base()  would write: id                  Fallbacks: none
	Alias() Builder

	// Base writes the intrinsic value of the column
	//
	// If this column was generated using a complex builder
	// such as Concat() or Search(), this method does nothing,
	// and that may break the statement.
	//
	// e.g: if we have this column being built: COUNT(id) AS count:
	//  - Def()   would write: COUNT(id) AS count; Fallbacks: Expr() -> Base()
	//  - Expr()  would write: COUNT(id);          Fallbacks: Base()
	//  - Alias() would write: count;              Fallbacks: Expr() -> Base()
	//  - Base()  would write: id                  Fallbacks: none
	Base() Builder

	// Count writes COUNT(column) or COUNT(DISTINCT column)
	//
	// @SQL: COUNT([SELF]) / COUNT(DISTINCT [SELF])
	Count(distinct ...bool) Column

	// Sum writes SUM(column) or SUM(DISTINCT column)
	//
	// @SQL: SUM([SELF])
	Sum(distinct ...bool) Column

	// Min writes MIN(column)
	//
	// @SQL: MIN([SELF])
	Min(distinct ...bool) Column

	// Max writes MAX(column)
	//
	// @SQL: MAX([SELF])
	Max(distinct ...bool) Column

	// Avg writes AVG(column)
	//
	// @SQL: AVG([SELF])
	Avg(distinct ...bool) Column

	// Add writes column + value
	//
	// @SQL: [SELF] + value
	Add(v Value) Column

	// Sub writes column - value
	//
	// @SQL: [SELF] - value
	Sub(v Value) Column

	// Mul writes column * value
	//
	// @SQL: [SELF] * value
	Mul(v Value) Column

	// Div writes column / value
	//
	// @SQL: [SELF] / value
	Div(v Value) Column

	// Mod writes column % value
	//
	// @SQL: [SELF] % value
	Mod(v Value) Column

	// Wrap wraps the entire column expression in parentheses
	//
	// @SQL: ([SELF])
	Wrap() Column

	// Lower writes LOWER(column)
	//
	// @SQL: LOWER([SELF])
	Lower() Column

	// Upper writes UPPER(column)
	//
	// @SQL: UPPER([SELF])
	Upper() Column

	// Trim writes TRIM(column)
	//
	// @SQL: TRIM([SELF])
	Trim() Column

	// Round writes ROUND(column)
	//
	// @SQL: ROUND([SELF])
	Round() Column

	// Abs writes ABS(column)
	//
	// @SQL: ABS([SELF])
	Abs() Column

	// As writes an alias for the column
	//
	// @SQL: [SELF] AS `name`
	As(name string) Column
}

// ----- Condition -----

// Condition represents a SQL condition used in WHERE and HAVING.
// It supports comparisons, pattern matching, ranges,
// null checks, and negation.
type Condition interface {
	// Not negates the current condition
	//
	// @SQL: NOT (`cond`)
	Not() Condition

	// Equal writes column = value
	//
	// @SQL: `col` = value
	Equal(v Value) ConditionNext

	// Like writes column LIKE value
	//
	// @SQL: `col` LIKE value
	Like(v Value) ConditionNext

	// Greater writes column > value
	//
	// @SQL: `col` > value
	Greater(v Value) ConditionNext

	// Less writes column < value
	//
	// @SQL: `col` < value
	Less(v Value) ConditionNext

	// In writes column IN (...)
	//
	// @SQL: `col` IN (value1, value2, ...)
	In(v []Value) ConditionNext

	// Between writes column BETWEEN min AND max
	//
	// @SQL: `col` BETWEEN min AND max
	Between(minV, maxV Value) ConditionNext

	// IsNull writes column IS NULL
	//
	// @SQL: `col` IS NULL
	IsNull() ConditionNext
}

// ConditionNext represents the continuation of a condition,
// allowing logical chaining with AND / OR.
type ConditionNext interface {
	// And links the next condition with AND
	//
	// @SQL: [Condition] ... AND `ident` ... [Condition]
	And(ident Identifier) Condition

	// Or links the next condition with OR
	//
	// @SQL: [Condition] ... OR `ident` ... [Condition]
	Or(ident Identifier) Condition
}

// ----- Case -----

// Case represents a SQL CASE expression,
// allowing conditional branching inside statements.
type Case interface {
	// When starts a CASE WHEN block
	//
	// @SQL: CASE WHEN `when`
	When(when Value) CaseNext

	// Else writes the ELSE value
	//
	// @SQL: ELSE value END
	Else(els Value)
}

// CaseNext represents the THEN part of a CASE WHEN branch,
// allowing chained WHEN/THEN blocks.
type CaseNext interface {
	// Then writes THEN value
	//
	// @SQL: THEN value
	Then(then Value) Case
}

// ----- Column Builder -----

// ColumnBuilder provides constructors for columns, placeholders,
// conditions, concatenations, and CASE-based expressions.
type ColumnBuilder interface {
	// Col creates a new column reference
	//
	// @SQL: `internal_alias`.`col`
	Col(name string, table ...string) Column

	// All writes the wildcard *
	//
	// You can specify a table to provide context when joins exist.
	// If not specified, the default table will be used.
	//
	// @SQL: *
	All(table ...string) Column

	// Ph writes a placeholder
	//
	// @SQL: ? / $[number]
	Ph() Builder

	// Cond starts a condition block with the given identifier
	//
	// @SQL: `ident` ... [Condition]
	Cond(ident Identifier) Condition

	// Concat writes CONCAT(...)
	//
	// @SQL: CONCAT(val1, val2, ...)
	Concat(v ...Identifier) Column

	// Switch writes a simple CASE expression for this column.
	//
	// The argument `cond` is evaluated and compared against literals
	// in each WHEN branch.
	//
	// @SQL:
	//  CASE `cond`
	//    WHEN compartion1 THEN ...
	//    WHEN comparation2 THEN ...
	//    ELSE ...
	//  END
	Switch(cond Identifier, fn func(cs Case)) Column

	// Search writes a searched CASE expression for this column.
	//
	// Each WHEN branch can use Cond() to define a boolean condition.
	//
	// @SQL:
	//  CASE
	//    WHEN boolean1 THEN ...
	//    WHEN boolean2 THEN ...
	//    ELSE ...
	//  END
	Search(fn func(cs Case)) Column
}

// ----- Clauses -----

// ----- SELECT

// Select represents the SELECT clause of a statement,
// supporting DISTINCT and selectable identifiers.
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

// SelectNext represents additional chained SELECT values
// when DISTINCT is already applied.
type SelectNext interface {
	// Select writes the SELECT clause values
	//
	// @SQL: SELECT `DISTINCT?` `sel1`, `sel2`, `sel3` ... [SelectNext]
	Select(sel ...Identifier)
}

// ----- FROM

// From represents the FROM clause,
// defining the statement's source table or subquery.
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

// Join represents the JOIN clause,
// initiating a table join operation.
type Join interface {
	// Join writes the JOIN clause
	//
	// @SQL: INNER JOIN `table 'alias'` ... [JoinNext]
	Join(table string) JoinNext
}

// JoinNext represents the ON clause for a join,
// returning a condition builder.
type JoinNext interface {
	// On writes the join condition
	//
	// @SQL: [join] ... ON `on` ... [Condition]
	On(on Identifier) Condition
}

// ----- WHERE

// Where represents the WHERE clause,
// supporting multiple AND-combined conditions.
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

// GroupBy represents the GROUP BY clause,
// grouping the selected rows.
type GroupBy interface {
	// GroupBy writes the GROUP BY clause
	//
	// @SQL: GROUP BY `group1`, `group2`, `group3` ...
	GroupBy(group ...Identifier)
}

// ----- HAVING

// Having represents the HAVING clause,
// filtering grouped results using conditions.
type Having interface {
	// Having writes the HAVING clause
	//
	// @SQL: HAVING `having` ... [Condition]
	Having(having Value) Condition
}

// ----- ORDER BY

// OrderBy represents the ORDER BY clause,
// sorting rows by one or more identifiers.
type OrderBy interface {
	// OrderBy writes the ORDER BY clause
	//
	// # Can be called multiple times for multiple orderings
	//
	// @SQL: ORDER BY `order` ... [OrderByNext]
	OrderBy(order Identifier) OrderByNext
}

// OrderByNext represents ordering modifiers,
// such as descending direction.
type OrderByNext interface {
	// Desc writes the descending direction for the current order
	//
	// @SQL: [OrderBy] ... DESC
	Desc()
}

// ----- LIMIT

// Limit represents the LIMIT clause,
// restricting the maximum number of returned rows.
type Limit interface {
	// Limit writes the LIMIT clause value
	//
	// # Values less than 1 will be ignored
	//
	// @SQL: LIMIT `limit`
	Limit(limit int)
}

// Offset represents the OFFSET clause,
// skipping a number of rows before returning results.
type Offset interface {
	// Offset writes the OFFSET clause value
	//
	// # Values less than 1 will be ignored
	//
	// @SQL: OFFSET `offset`
	Offset(offset int)
}

type Returning interface {
	Returning(ident []Identifier)
}

// ----- Statements -----

// SelectStmt represents a full SELECT statement
type SelectStmt interface {
	Executor
	Select
	From
	Join
	Where
	GroupBy
	Having
	OrderBy
	Limit
	Offset
	Raw() (string, []any)
}

type InsertStmt interface {
	Insert(values Binder)
	Into(table string)
	Returning
	Raw() (string, []any)
}

type UpdateStmt interface {
	Update(table string)
	Set(value Binder)
	From
	Where
	Returning
	Raw() (string, []any)
}

type DeleteStmt interface {
	From
	Where
	Returning
}

// ----- Executor -----

type Execute interface {
	Run() error
	Find(binder Binder) error
}

type Prepare interface {
	Run(args []Value) error
	Find(args []Value, binder Binder) error
}

type Executor interface {
	Execute
	Prepare(fn func(exec Prepare) error) error
}
