package entity

// --- Clauses

type Select struct {
	Dist   bool
	Values []any
}

type From struct {
	Value any
}

type Join struct {
	Joins []JoinSegment
}

type Where struct {
	Conds []Condition
}

type GroupBy struct {
	Values []any
}

type Having struct {
	Cond Condition
}

type OrderBy struct {
	Orders []OrderSegment
}

type Limit struct {
	ValueNumber int
}

type Offset struct {
	ValueNumber int
}

type InsertInto struct {
	Table string
	Cols  []string
	Rows  int
}

type Update struct {
	Table string
	Cols  []any
}

type Delete struct{}

// --- Kinds

func (c *Select) Kind() ClauseKind     { return ClauseSelect }
func (c *From) Kind() ClauseKind       { return ClauseFrom }
func (c *Join) Kind() ClauseKind       { return ClauseJoin }
func (c *Where) Kind() ClauseKind      { return ClauseWhere }
func (c *GroupBy) Kind() ClauseKind    { return ClauseGroupBy }
func (c *Having) Kind() ClauseKind     { return ClauseHaving }
func (c *OrderBy) Kind() ClauseKind    { return ClauseOrderBy }
func (c *Limit) Kind() ClauseKind      { return ClauseLimit }
func (c *Offset) Kind() ClauseKind     { return ClauseOffset }
func (c *InsertInto) Kind() ClauseKind { return ClauseInsertInto }
func (c *Update) Kind() ClauseKind     { return ClauseUpdate }
func (c *Delete) Kind() ClauseKind     { return ClauseDelete }

// --- Writes

func (c *Select) Write(w Writer, dial ClauseWriter)     { dial.WriteSelect(w, c) }
func (c *From) Write(w Writer, dial ClauseWriter)       { dial.WriteFrom(w, c) }
func (c *Join) Write(w Writer, dial ClauseWriter)       { dial.WriteJoin(w, c) }
func (c *Where) Write(w Writer, dial ClauseWriter)      { dial.WriteWhere(w, c) }
func (c *GroupBy) Write(w Writer, dial ClauseWriter)    { dial.WriteGroupBy(w, c) }
func (c *Having) Write(w Writer, dial ClauseWriter)     { dial.WriteHaving(w, c) }
func (c *OrderBy) Write(w Writer, dial ClauseWriter)    { dial.WriteOrderBy(w, c) }
func (c *Limit) Write(w Writer, dial ClauseWriter)      { dial.WriteLimit(w, c) }
func (c *Offset) Write(w Writer, dial ClauseWriter)     { dial.WriteOffset(w, c) }
func (c *InsertInto) Write(w Writer, dial ClauseWriter) { dial.WriteInsertInto(w, c) }
func (c *Update) Write(w Writer, dial ClauseWriter)     { dial.WriteUpdate(w, c) }
func (c *Delete) Write(w Writer, dial ClauseWriter)     { dial.WriteDelete(w, c) }

// --- Utilities

type JoinType int

const (
	JoinInner JoinType = iota
	JoinLeft
	JoinRight
	JoinFull
	JoinCross
)

type JoinSegment struct {
	Type  JoinType
	Table Table
	Cond  *Condition
}

type OrderSegment struct {
	Value      any
	Descending bool
}
