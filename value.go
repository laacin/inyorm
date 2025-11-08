package inyorm

// Value is a versatile type representing any element that can be used in a SQL expression
// within the ORM. It is essentially an alias for `any`, but semantically it signals
// that the value is intended to be interpreted in a SQL context.
//
// Value can be:
//   - SQL literals: strings, numbers, booleans, or nil. These will be rendered as
//     valid SQL constants, e.g., "example" → 'example', 0 → 0, true → 1, nil → NULL.
//   - ORM fields or columns, via ColumnExpr.Col(). This ensures the value is treated
//     as a column reference rather than a literal.
//   - Expressions (*ExprEnd) created with the ORM’s expression builder. These
//     are evaluated depending on the logical context where they are used.
//
// Using Value makes ORM operations declarative and consistent, allowing you to
// pass either raw data or constructed expressions seamlessly when building queries.
type Value = any
