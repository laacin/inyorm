package inyorm

// ----- SCHEMA -----

type Table interface {
	TableName() string
	PKey() string
	FKeys() map[string]string
}

type JoinTable interface {
	Table() string
	Keys() map[string]string
}

type Column interface {
	Name() string
	As(name string)

	Count()
	Sum()
	Avg()
	Max()
	Min()

	Concat(values ...string)
	Substring(start, end int)
	Upper()
	Lower()
	Trim()

	// Date()
	// DateTrunc()
	// Extract()
	Now()

	Switch(target any, fn func(w Case[string]))
	Search(target any) Case[func(Expression)]
}

type Case[T any] interface {
	When(value T) Then[T]
	Else(do any)
}

type Then[T any] interface {
	Then(do any) Case[T]
}

type Operation interface {
	Add()
	Sub()
	Mul()
	Div()
	Mod()

	Get()
}

// ----- EXPRESSION -----

// // Expression represents the start of a conditional expression.
// // It allows building chained conditions on columns or fields.
// type Expression interface {
// 	// Not negates the current condition.
// 	// If called before an operator, the operation will be negated.
// 	Not() Expression
//
// 	// Equal adds an equality condition: field = value.
// 	// Returns ExpressionEnd to chain logical AND/OR conditions.
// 	Equal(value any) ExpressionEnd
//
// 	// Like adds a pattern-matching condition: field LIKE value.
// 	// Returns ExpressionEnd to chain logical conditions.
// 	Like(value any) ExpressionEnd
//
// 	// In adds a membership condition: field IN (values...).
// 	// Returns ExpressionEnd to chain logical AND/OR conditions.
// 	In(values ...any) ExpressionEnd
//
// 	// Between adds a range condition: field BETWEEN minV AND maxV.
// 	// Returns ExpressionEnd to chain logical conditions.
// 	Between(minV, maxV any) ExpressionEnd
//
// 	// Greater adds a greater-than condition: field > value.
// 	// Returns ExpressionEnd to chain logical AND/OR conditions.
// 	Greater(value any) ExpressionEnd
//
// 	// Less adds a less-than condition: field < value.
// 	// Returns ExpressionEnd to chain logical conditions.
// 	Less(value any) ExpressionEnd
//
// 	// IsNull adds a null-check condition: field IS NULL.
// 	// Returns ExpressionEnd to chain logical AND/OR conditions.
// 	IsNull() ExpressionEnd
// }
//
// // ExpressionEnd represents the end of a conditional expression,
// // allowing logical connectors to continue building the query.
// type ExpressionEnd interface {
// 	// Or starts a new condition with a logical OR.
// 	// If identifiers are provided, the first is treated as the field to evaluate,
// 	// and the second (optional) as its alias or table.
// 	Or(identifier ...string) Expression
//
// 	// And starts a new condition with a logical AND.
// 	// If identifiers are provided, the first is treated as the field to evaluate,
// 	// and the second (optional) as its alias or table.
// 	And(identifier ...string) Expression
// }

type Query interface {
	Select(selects ...string)
	Join(func(j Join))
	Where(column string) Expression
	GroupBy(columns ...string) GroupBy
	Limit(num int)
	Offset(num int)
}

// --- Clauses

type Join interface {
	Simple(tbl Table)
	Many(interTbl JoinTable, tbls ...Table)
}

type GroupBy interface {
	Having(value any)
}
