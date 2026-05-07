package api

// ExprBuilder provides constructors for columns, placeholders,
// conditions, concatenations, and CASE-based expressions.
type ExprBuilder[Cs any] interface {
	Table(name string) Table

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

	// Param writes a placeholder
	//
	// You can omit the value for lazy parameters, useful in prepared statements.
	//
	// @SQL: ? / $[number]
	Param(value ...Value) Parameter

	// Cond starts a condition block with the given identifier
	//
	// @SQL: `ident` ... [Condition]
	Cond(ident Value) Condition

	// Concat writes CONCAT(...)
	//
	// @SQL: CONCAT(val1, val2, ...)
	Concat(v ...Value) Column

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
	Switch(cond Value, fn func(cs Cs)) Column

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
	Search(fn func(cs Cs)) Column
}
