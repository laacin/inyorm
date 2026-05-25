package dml

import "github.com/laacin/inyorm/internal/core"

// ---- SELECT QUERY -----

type QuerySelect struct {
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

func (q *QuerySelect) Build(b *core.Builder) error {
	if q.ClauseSelect.IsDeclared() {
		q.ClauseSelect.Build(b)
	}

	if q.ClauseFrom.IsDeclared() {
		q.ClauseFrom.Build(b)
	}

	if q.ClauseJoin.IsDeclared() {
		q.ClauseJoin.Build(b)
	}

	if q.ClauseWhere.IsDeclared() {
		q.ClauseWhere.Build(b)
	}

	if q.ClauseGroupBy.IsDeclared() {
		q.ClauseGroupBy.Build(b)
	}

	if q.ClauseHaving.IsDeclared() {
		q.ClauseHaving.Build(b)
	}

	if q.ClauseOrderBy.IsDeclared() {
		q.ClauseOrderBy.Build(b)
	}

	if q.ClauseLimit.IsDeclared() {
		q.ClauseLimit.Build(b)
	}

	if q.ClauseOffset.IsDeclared() {
		q.ClauseOffset.Build(b)
	}

	return nil
}

func (q *QuerySelect) Render(w core.InternalWriter, dial QueryWriter) error {
	dial.WriteQuerySelect(w, q)
	return nil
}

// ----- INSERT QUERY ------

type QueryInsert struct {
	ClauseInsertInto
}

func (q *QueryInsert) Build(b *core.Builder) error {
	if q.ClauseInsertInto.IsDeclared() {
		q.ClauseInsertInto.Build(b)
	}

	return nil
}

func (q *QueryInsert) Render(w core.InternalWriter, dial QueryWriter) error {
	dial.WriteQueryInsert(w, q)
	return nil
}

// ----- UPDATE QUERY ------

type QueryUpdate struct {
	ClauseUpdate
	ClauseWhere
}

func (q *QueryUpdate) Build(b *core.Builder) error {
	if q.ClauseUpdate.IsDeclared() {
		q.ClauseUpdate.Build(b)
	}

	if q.ClauseWhere.IsDeclared() {
		q.ClauseWhere.Build(b)
	}

	return nil
}

func (q *QueryUpdate) Render(w core.InternalWriter, dial QueryWriter) error {
	dial.WriteQueryUpdate(w, q)
	return nil
}

// ----- DELETE QUERY ------

type QueryDelete struct {
	ClauseDelete
	ClauseFrom
	ClauseWhere
}

func (q *QueryDelete) Build(b *core.Builder) error {
	if q.ClauseDelete.IsDeclared() {
		q.ClauseDelete.Build(b)
	}

	if q.ClauseFrom.IsDeclared() {
		q.ClauseFrom.Build(b)
	}

	if q.ClauseWhere.IsDeclared() {
		q.ClauseWhere.Build(b)
	}

	return nil
}

func (q *QueryDelete) Render(w core.InternalWriter, dial QueryWriter) error {
	dial.WriteQueryDelete(w, q)
	return nil
}
