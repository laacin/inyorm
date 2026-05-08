package api

// ----- SELECT

// Select represents the SELECT clause of a statement,
// supporting DISTINCT and selectable identifiers.
type Select interface {
	// Select writes the SELECT clause values
	//
	// @SQL: SELECT `DISTINCT?` `sel1`, `sel2`, `sel3` ... [SelectNext]
	Select(sel ...Value) SelectNext
}

// SelectNext represents additional chained SELECT values
// when DISTINCT is already applied.
type SelectNext interface {
	// Distinct writes DISTINCT in the SELECT clause
	//
	// @SQL: SELECT `DISTINCT` ...
	Distinct()
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
	From(table Value)
}

// ----- JOIN

// Join represents the JOIN clause,
// initiating a table join operation.
type Join interface {
	// Join writes the JOIN clause
	//
	// @SQL: INNER JOIN `table 'alias'` ... [JoinNext]
	Join(table Value) JoinNext
}

// JoinNext represents the ON clause for a join,
// returning a condition builder.
type JoinNext interface {
	// Left changes the current join type to LEFT JOIN.
	//
	// @SQL: LEFT JOIN `table` ... [JoinEnd]
	Left() JoinEnd

	// Right changes the current join type to RIGHT JOIN.
	//
	// @SQL: RIGHT JOIN `table` ... [JoinEnd]
	Right() JoinEnd

	// Full changes the current join type to FULL JOIN.
	//
	// @SQL: FULL JOIN `table` ... [JoinEnd]
	Full() JoinEnd

	// Cross changes the current join type to CROSS JOIN.
	//
	// CROSS JOIN does not use an ON clause.
	//
	// @SQL: CROSS JOIN `table`
	Cross()

	// On writes the join condition
	//
	// @SQL: [join] ... ON `on` ... [Condition]
	On(on Value) Condition
}

type JoinEnd interface {
	// On writes the join condition
	//
	// @SQL: [join] ... ON `on` ... [Condition]
	On(on Value) Condition
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
	Where(ident Value) Condition
}

// ----- GROUP BY

// GroupBy represents the GROUP BY clause,
// grouping the selected rows.
type GroupBy interface {
	// GroupBy writes the GROUP BY clause
	//
	// @SQL: GROUP BY `group1`, `group2`, `group3` ...
	GroupBy(group ...Value)
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
	OrderBy(order Value) OrderByNext
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

// Insert represents the INSERT INTO clause, defining data to insert into a table.
type Insert interface {
	// Insert specifies the columns to insert and returns a Values handler.
	//
	// The argument can be:
	//   - a struct
	//   - a map[string]any
	//   - individual columns
	//   - any type where columns can be read dynamically
	//
	// Example SQL: INSERT INTO `table` (`columns`) VALUES (...)
	Insert(reference ...Value) Values

	// InsertIgnore behaves like Insert, but ignores specified columns from the reference.
	//
	// Useful when passing a struct with many fields and excluding some of them.
	InsertIgnore(reference Value, ignores ...Value) Values
}

// Update represents the UPDATE clause, defining data to update in a table.
type Update interface {
	// Update specifies the columns to update and returns a Values handler.
	//
	// The argument can be:
	//   - a struct
	//   - a map[string]any
	//   - individual columns
	//   - any type where columns can be read dynamically
	//
	// Example SQL: UPDATE `table` SET `column` = ...
	Update(reference ...Value) Values

	// UpdateIgnore behaves like Update, but ignores specified columns from the reference.
	//
	// Useful when passing a struct with many fields and excluding some of them.
	UpdateIgnore(reference Value, ignores ...Value) Values
}

// Values represents the VALUES clause for INSERT or the assignment values for UPDATE.
type Values interface {
	// Values sets the actual values to insert or update.
	//
	// Always provide the raw values; placeholders ("?") are handled internally.
	// You can omit this call if using prepared columns without values.
	Values(values Value)
}
