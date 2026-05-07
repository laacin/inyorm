package api

// SelectStmt represents a full SELECT statement
type SelectStmt interface {
	// Executor
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

// InsertStmt represents a full INSERT statement
type InsertStmt interface {
	// Executor
	Insert
}

// UpdateStmt represents a full UPDATE statement
type UpdateStmt interface {
	// Executor
	Update
	Where
}

// DeleteStmt represents a full DELETE statement
type DeleteStmt interface {
	// Executor
	From
	Where
}
