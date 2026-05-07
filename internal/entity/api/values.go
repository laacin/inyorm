package api

type (
	Value   = any
	Scanner = any
)

// --- Table

type Table any

// ---- Column

// Column represents a selectable or computable SQL column.
// Provides builders for expressions, aliases, aggregates,
// arithmetic operations, and transformations.
type Column interface {

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

// ---- Param ----

type Parameter = any

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
	And(ident Value) Condition

	// Or links the next condition with OR
	//
	// @SQL: [Condition] ... OR `ident` ... [Condition]
	Or(ident Value) Condition
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
