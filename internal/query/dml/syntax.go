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
	ClauseKindInsertInto

	// Update statement
	ClauseKindUpdate

	// Delete statement
	ClauseKindDelete
)

type ClauseWriter interface {
	WriteClauseSelect(core.Writer, *ClauseSelect)
	WriteClauseFrom(core.Writer, *ClauseFrom)
	WriteClauseJoin(core.Writer, *ClauseJoin)
	WriteClauseWhere(core.Writer, *ClauseWhere)
	WriteClauseGroupBy(core.Writer, *ClauseGroupBy)
	WriteClauseHaving(core.Writer, *ClauseHaving)
	WriteClauseOrderBy(core.Writer, *ClauseOrderBy)
	WriteClauseLimit(core.Writer, *ClauseLimit)
	WriteClauseOffset(core.Writer, *ClauseOffset)

	WriteClauseInsertInto(core.Writer, *ClauseInsertInto)

	WriteClauseUpdate(core.Writer, *ClauseUpdate)

	WriteClauseDelete(core.Writer, *ClauseDelete)
}
