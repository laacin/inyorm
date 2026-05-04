package dialect

// Clause name
type ClauseName int

const (
	// Insert Statement
	ClauseNameInsertInto ClauseName = iota

	// Select Statement
	ClauseNameSelect
	ClauseNameFrom
	ClauseNameJoin
	ClauseNameWhere
	ClauseNameGroupBy
	ClauseNameHaving
	ClauseNameOrderBy
	ClauseNameLimit
	ClauseNameOffset

	// Update Statement
	ClauseNameUpdate

	// Delete Statement
	ClauseNameDelete
)

// Clause types
type JoinType int

const (
	JoinInner JoinType = iota
	JoinLeft
	JoinRight
	JoinFull
	JoinCross
)

// ----- Clause building tools

type InsertIntoTools struct {
	Table   string
	Columns []string
	Rows    int
}

type SelectTools struct {
	Distinct bool
	Values   []any
}

type FromTools struct {
	Value any
}

type JoinTools struct {
	Type  JoinType
	Table string
	Cond  *Cond
} // Can be []JoinTools for multiple joins

type WhereTools struct {
	Conds []Cond
}

type GroupByTools struct {
	Values []any
}

type HavingTools struct {
	Cond Cond
}

type OrderByTools struct {
	Value      any
	Descending bool
} // Can be []OrderByTools for multiple orders

type LimitTools struct {
	ValueNumber int
}

type OffsetTools struct {
	ValueNumber int
}

type UpdateTools struct {
	Table   string
	Columns []string
}

type DeleteTools struct {
	Hard bool
}
