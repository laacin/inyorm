package dml

import "github.com/laacin/inyorm/internal/core"

type ClauseKind int

const (
	// Select statement
	ClauseKindSelect ClauseKind = iota
	ClauseKindFrom
	ClauseKindJoin
	ClauseKindWhere
	ClauseKindGroupBy
	ClauseKindHaving
	ClauseKindOrderBy
	ClauseKindLimit
	ClauseKindOffset

	// Insert statement
	ClauseKindInsert

	// Update statement
	ClauseKindUpdate

	// Delete statement
	ClauseKindDelete
)

type ClauseWriter interface {
	WriteSelect(core.Writer, *ClauseSelect)
	WriteFrom(core.Writer, *ClauseFrom)
	WriteJoin(core.Writer, *ClauseJoin)
	WriteWhere(core.Writer, *ClauseWhere)
	WriteGroupBy(core.Writer, *ClauseGroupBy)
	WriteHaving(core.Writer, *ClauseHaving)
	WriteOrderBy(core.Writer, *ClauseOrderBy)
	WriteLimit(core.Writer, *ClauseLimit)
	WriteOffset(core.Writer, *ClauseOffset)

	WriteInsertInto(core.Writer, *ClauseInsert)

	WriteUpdate(core.Writer, *ClauseUpdate)

	WriteDelete(core.Writer, *ClauseDelete)
}

type QueryOrder interface {
	SelectOrder() []ClauseKind
	InsertOrder() []ClauseKind
	UpdateOrder() []ClauseKind
	DeleteOrder() []ClauseKind
}

// --- Internal
type ClauseBuilder interface {
	IsDeclared() bool
	Kind() ClauseKind
	Build(core.InternalWriter, ClauseWriter) error
}
