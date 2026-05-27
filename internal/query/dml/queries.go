package dml

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/query"
)

// --- Constructors

func NewSelect(rend Renderer) *QuerySelect { return &QuerySelect{rend: rend} }
func NewInsert(rend Renderer) *QueryInsert { return &QueryInsert{rend: rend} }
func NewUpdate(rend Renderer) *QueryUpdate { return &QueryUpdate{rend: rend} }
func NewDelete(rend Renderer) *QueryDelete { return &QueryDelete{rend: rend} }

// ---- SELECT -----

type QuerySelect struct {
	rend Renderer
	helper

	ClauseSelect
	ClauseFrom
	ClauseJoin
	ClauseWhere
	ClauseGroupBy
	ClauseHaving
	ClauseOrderBy
	ClauseLimit
	ClauseOffset
}

func (q *QuerySelect) Build(tools *query.Tools) error {
	q.must(tools, &q.ClauseSelect)
	q.must(tools, &q.ClauseFrom, func() {
		if tbl, ok := q.ClauseFrom.Value.(*expr.Table); ok {
			tools.Aliases.SetMain(tbl.Name)
		}
	})

	q.build(tools, &q.ClauseJoin, func() {
		tools.Aliases.Enable()
	})

	q.build(tools, &q.ClauseWhere)
	q.build(tools, &q.ClauseGroupBy)
	q.build(tools, &q.ClauseHaving)
	q.build(tools, &q.ClauseOrderBy)
	q.build(tools, &q.ClauseLimit)
	q.build(tools, &q.ClauseOffset)

	return q.end()
}

func (q *QuerySelect) Render(w core.InternalWriter) error {
	q.rend.WriteQuerySelect(w, q)
	return nil
}

// ---- INSERT -----

type QueryInsert struct {
	rend Renderer
	helper

	ClauseInsertInto
}

func (q *QueryInsert) Api() api.InsertQuery {
	return q
}

func (q *QueryInsert) Build(tools *query.Tools) error {
	q.must(tools, &q.ClauseInsertInto, func() {
		if tbl, ok := q.ClauseInsertInto.Table.(*expr.Table); ok {
			tools.Aliases.SetMain(tbl.Name)
		}
	})

	return q.end()
}

func (q *QueryInsert) Render(w core.InternalWriter) error {
	q.rend.WriteQueryInsert(w, q)
	return nil
}

// ---- UPDATE -----

type QueryUpdate struct {
	rend Renderer
	helper

	ClauseUpdate
	ClauseWhere
}

func (q *QueryUpdate) Api() api.UpdateQuery {
	return q
}

func (q *QueryUpdate) Build(tools *query.Tools) error {
	q.must(tools, &q.ClauseUpdate, func() {
		if tbl, ok := q.ClauseUpdate.Table.(*expr.Table); ok {
			tools.Aliases.SetMain(tbl.Name)
		}
	})
	q.must(tools, &q.ClauseUpdate)

	return nil
}

func (q *QueryUpdate) Render(w core.InternalWriter) error {
	q.rend.WriteQueryUpdate(w, q)
	return nil
}

// ---- DELETE -----

type QueryDelete struct {
	rend Renderer
	helper

	ClauseDelete
	ClauseFrom
	ClauseWhere
}

func (q *QueryDelete) Api() api.DeleteQuery {
	return q
}

func (q *QueryDelete) Build(tools *query.Tools) error {
	q.must(tools, &q.ClauseDelete)

	q.must(tools, &q.ClauseFrom, func() {
		if tbl, ok := q.ClauseFrom.Value.(*expr.Table); ok {
			tools.Aliases.SetMain(tbl.Name)
		}
	})

	q.must(tools, &q.ClauseWhere)
	return nil
}

func (q *QueryDelete) Render(w core.InternalWriter) error {
	q.rend.WriteQueryDelete(w, q)
	return nil
}
