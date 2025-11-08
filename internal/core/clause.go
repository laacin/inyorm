package core

// Clauses name
const (
	ClsSelect  = "SELECT"
	ClsFrom    = "FROM"
	ClsJoin    = "JOIN"
	ClsWhere   = "WHERE"
	ClsGroupBy = "GROUP BY"
	ClsOrderBy = "ORDER BY"
	ClsLimit   = "LIMIT"
	ClsOffset  = "OFFSET"
)

type Clause interface {
	Name() string
	Build() func(w Writer)
}

// ---- Clauses

// Insert clause
type ClauseInsert interface {
	InsertInto(table string) ClauseValues
}

// Select clause
type ClauseSelect interface {
	Distinct() ClauseSelect
	Select(values ...any)
}

// From clause
type ClauseFrom interface {
	From(value any)
}

// Join clause
type ClauseJoin interface {
	Join(table string) ClauseOn
	JoinLeft(table string) ClauseOn
	JoinRight(table string) ClauseOn
	JoinFull(table string) ClauseOn
	JoinCross(table string)
}

// Where clause
type ClauseWhere interface {
	Where(identifier any) Cond
}

// Group by clause
type ClauseGroupBy interface {
	GroupBy(values ...any) ClauseHaving
}

// Order by clause
type ClauseOrderBy interface {
	OrderBy(values ...any)
}

// Limit clause
type ClauseLimit interface {
	Limit(value int)
}

// Offset clause
type ClauseOffset interface {
	Offset(value int)
}

// ---- Depending clauses

type ClauseValues interface {
	Values(v any)
}

type ClauseOn interface {
	On(identifier any) Cond
}

type ClauseHaving interface {
	Having(identifier any) Cond
}
