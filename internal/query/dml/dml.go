package dml

import (
	"fmt"

	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/query"
)

type Renderer interface {
	// --- Queries
	WriteQuerySelect(core.Writer, *QuerySelect)
	WriteQueryInsert(core.Writer, *QueryInsert)
	WriteQueryUpdate(core.Writer, *QueryUpdate)
	WriteQueryDelete(core.Writer, *QueryDelete)

	// --- Dep
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

func (k ClauseKind) String() string {
	switch k {
	case ClauseKindSelect:
		return "SELECT"
	case ClauseKindFrom:
		return "FROM"
	case ClauseKindJoin:
		return "JOIN"
	case ClauseKindWhere:
		return "WHERE"
	case ClauseKindGroupBy:
		return "GROUP BY"
	case ClauseKindHaving:
		return "HAVING"
	case ClauseKindOrderBy:
		return "ORDER BY"
	case ClauseKindLimit:
		return "LIMIT"
	case ClauseKindOffset:
		return "OFFSET"
	case ClauseKindInsertInto:
		return "INSERT INTO"
	case ClauseKindUpdate:
		return "UPDATE"
	case ClauseKindDelete:
		return "DELETE"
	default:
		return ""
	}
}

type Clause interface {
	Kind() ClauseKind
	IsDeclared() bool
	Build(*query.Tools) error
}

// --- Helpers

type helper struct {
	errs []error
}

func (q helper) must(tools *query.Tools, cls Clause, cb ...func()) {
	if !cls.IsDeclared() {
		q.errs = append(q.errs, fmt.Errorf("clause %s must be declared", cls.Kind()))
		return
	}

	if len(cb) > 0 {
		cb[0]()
	}

	if err := cls.Build(tools); err != nil {
		q.errs = append(q.errs, err)
	}
}

func (q helper) build(tools *query.Tools, cls Clause, cb ...func()) {
	if cls.IsDeclared() {
		if len(cb) > 0 {
			cb[0]()
		}

		if err := cls.Build(tools); err != nil {
			q.errs = append(q.errs, err)
		}
	}
}

func (q helper) end() error {
	for _, err := range q.errs {
		if err != nil {
			return err
		}
	}
	return nil
}
