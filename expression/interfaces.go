package expression

// ----- EXPRESSION -----

// ExpressionStart represents the start of a conditional expression.
// It allows building chained conditions on columns or fields.
type ExpressionStart interface {
	// Not negates the current condition.
	// If called before an operator, the operation will be negated.
	Not() ExpressionStart

	// Equal adds an equality condition: field = value.
	// Returns ExpressionEnd to chain logical AND/OR conditions.
	Equal(value any) ExpressionEnd

	// Like adds a pattern-matching condition: field LIKE value.
	// Returns ExpressionEnd to chain logical conditions.
	Like(value any) ExpressionEnd

	// In adds a membership condition: field IN (values...).
	// Returns ExpressionEnd to chain logical AND/OR conditions.
	In(values ...any) ExpressionEnd

	// Between adds a range condition: field BETWEEN minV AND maxV.
	// Returns ExpressionEnd to chain logical conditions.
	Between(minV, maxV any) ExpressionEnd

	// Greater adds a greater-than condition: field > value.
	// Returns ExpressionEnd to chain logical AND/OR conditions.
	Greater(value any) ExpressionEnd

	// Less adds a less-than condition: field < value.
	// Returns ExpressionEnd to chain logical conditions.
	Less(value any) ExpressionEnd

	// IsNull adds a null-check condition: field IS NULL.
	// Returns ExpressionEnd to chain logical AND/OR conditions.
	IsNull() ExpressionEnd
}

// ExpressionEnd represents the end of a conditional expression,
// allowing logical connectors to continue building the query.
type ExpressionEnd interface {
	// Or starts a new condition with a logical OR.
	// If identifiers are provided, the first is treated as the field to evaluate,
	// and the second (optional) as its alias or table.
	Or(identifier ...string) ExpressionStart

	// And starts a new condition with a logical AND.
	// If identifiers are provided, the first is treated as the field to evaluate,
	// and the second (optional) as its alias or table.
	And(identifier ...string) ExpressionStart
}

// ----- AGGREGATION -----

type ColumnExpressions interface {
	Count(value string)
	Sum(value string)
	Avg(value string)
	Max(value string)
	Min(value string)

	Concat(values ...string)
	Substring(value string, start, end int)
	Upper(value string)
	Lower(value string)
	Trim(value string)

	Date(value string)
	DateTrunc(value string)
	Extract(value string)
	Now()
}
