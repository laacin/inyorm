package dml

type ClauseKind int

const (
	// Select statement
	ClauseSelect ClauseKind = iota
	ClauseFrom
	ClauseJoin
	ClauseWhere
	ClauseGroupBy
	ClauseHaving
	ClauseOrderBy
	ClauseLimit
	ClauseOffset

	// Insert statement
	ClauseInsertInto

	// Update statement
	ClauseUpdate

	// Delete statement
	ClauseDelete
)

type StatementKind int

const (
	StatementSelect StatementKind = iota
	StatementInsert
	StatementUpdate
	StatementDelete
)
